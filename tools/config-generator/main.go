/*
Copyright 2019 The Knative Authors

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// The make_config tool generates a full Prow config for the Knative project,
// with input from a yaml file with key definitions.

package main

import (
	"bytes"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"runtime"
	"sort"
	"strconv"
	"strings"
	"text/template"
	"time"

	"gopkg.in/yaml.v2"
	"k8s.io/apimachinery/pkg/util/sets"
)

const (
	// Manifests generated by ko are indented by 2 spaces.
	baseIndent  = "  "
	templateDir = "templates"

	// ##########################################################
	// ############## prow configuration templates ##############
	// ##########################################################
	// commonHeaderConfig contains common header definitions.
	commonHeaderConfig = "common_header.yaml"

	// generalProwConfig contains config-wide definitions.
	generalProwConfig = "prow_config.yaml"

	// pluginsConfig is the template for the plugins YAML file.
	pluginsConfig = "prow_plugins.yaml"
)

// repositoryData contains basic data about each Knative repository.
type repositoryData struct {
	Name                   string
	EnablePerformanceTests bool
	EnableGoCoverage       bool
	GoCoverageThreshold    int
	Processed              bool
	DotDev                 bool
}

// prowConfigTemplateData contains basic data about Prow.
type prowConfigTemplateData struct {
	Year              int
	GcsBucket         string
	PresubmitLogsDir  string
	LogsDir           string
	ProwHost          string
	TestGridHost      string
	GubernatorHost    string
	TestGridGcsBucket string
	TideRepos         []string
	ManagedRepos      []string
	ManagedOrgs       []string
	JobConfigPath     string
	TestInfraRepo     string
}

// baseProwJobTemplateData contains basic data about a Prow job.
type baseProwJobTemplateData struct {
	OrgName             string
	RepoName            string
	RepoNameForJob      string
	GcsBucket           string
	GcsLogDir           string
	GcsPresubmitLogDir  string
	RepoURI             string
	RepoBranch          string
	CloneURI            string
	SecurityContext     []string
	SkipBranches        []string
	Branches            []string
	DecorationConfig    []string
	ExtraRefs           []string
	Command             string
	Args                []string
	Env                 []string
	Volumes             []string
	VolumeMounts        []string
	Resources           []string
	Timeout             int
	AlwaysRun           bool
	TestAccount         string
	ServiceAccount      string
	ReleaseGcs          string
	GoCoverageThreshold int
	Image               string
	Labels              []string
	PathAlias           string
	RunIfChanged        string
	Cluster             string
	Optional            string
	NeedsMonitor        bool
}

// ####################################################################################################
// ################ data definitions that are used for the prow config file generation ################
// ####################################################################################################

// sectionGenerator is a function that generates Prow job configs given a slice of a yaml file with configs.
type sectionGenerator func(string, string, yaml.MapSlice)

// stringArrayFlag is the content of a multi-value flag.
type stringArrayFlag []string

var (
	// Values used in the jobs that can be changed through command-line flags.
	// TODO: these should be CapsCase
	// ... until they are not global
	output                     *os.File
	prowHost                   string
	testGridHost               string
	gubernatorHost             string
	GCSBucket                  string
	testGridGcsBucket          string
	LogsDir                    string
	presubmitLogsDir           string
	testAccount                string
	nightlyAccount             string
	releaseAccount             string
	githubCommenterDockerImage string
	prowTestsDockerImage       string
	presubmitScript            string
	releaseScript              string
	webhookAPICoverageScript   string

	// #########################################################################
	// ############## data used for generating prow configuration ##############
	// #########################################################################
	// Array constants used throughout the jobs.
	allPresubmitTests = []string{"--all-tests"}
	releaseNightly    = []string{"--publish", "--tag-release"}
	releaseLocal      = []string{"--nopublish", "--notag-release"}

	// Overrides and behavior changes through command-line flags.
	repositoryOverride string
	jobNameFilter      string
	preCommand         string
	extraEnvVars       stringArrayFlag
	timeoutOverride    int

	// List of Knative repositories.
	// Not guaranteed unique by any value of the struct
	repositories []repositoryData

	// Map which sections of the config.yaml were written to stdout.
	sectionMap map[string]bool

	// To be used to flag that outputConfig() emitted data.
	emittedOutput bool

	releaseRegex = regexp.MustCompile(`.+-[0-9\.]+$`)
)

// Yaml parsing helpers.

// read template yaml file content
func readTemplate(fp string) string {
	if _, ok := templatesCache[fp]; !ok {
		// get the directory of the currently running file
		_, f, _, _ := runtime.Caller(0)
		content, err := ioutil.ReadFile(path.Join(path.Dir(f), templateDir, fp))
		if err != nil {
			log.Fatalf("Failed read file '%s': '%v'", fp, err)
		}
		templatesCache[fp] = string(content)
	}
	return templatesCache[fp]
}

// getString casts the given interface (expected string) as string.
// An array of length 1 is also considered a single string.
func getString(s interface{}) string {
	if _, ok := s.([]interface{}); ok {
		values := getStringArray(s)
		if len(values) == 1 {
			return values[0]
		}
		log.Fatalf("Entry %v is not a string or string array of size 1", s)
	}
	if str, ok := s.(string); ok {
		return str
	}
	log.Fatalf("Entry %v is not a string", s)
	return ""
}

// getInt casts the given interface (expected int) as int.
func getInt(s interface{}) int {
	if value, ok := s.(int); ok {
		return value
	}
	log.Fatalf("Entry %v is not an integer", s)
	return 0
}

// getBool casts the given interface (expected bool) as bool.
func getBool(s interface{}) bool {
	if value, ok := s.(bool); ok {
		return value
	}
	log.Fatalf("Entry %v is not a boolean", s)
	return false
}

// getInterfaceArray casts the given interface (expected interface array) as interface array.
func getInterfaceArray(s interface{}) []interface{} {
	if interfaceArray, ok := s.([]interface{}); ok {
		return interfaceArray
	}
	log.Fatalf("Entry %v is not an interface array", s)
	return nil
}

// getStringArray casts the given interface (expected string array) as string array.
func getStringArray(s interface{}) []string {
	interfaceArray := getInterfaceArray(s)
	strArray := make([]string, len(interfaceArray))
	for i := range interfaceArray {
		strArray[i] = getString(interfaceArray[i])
	}
	return strArray
}

// getMapSlice casts the given interface (expected MapSlice) as MapSlice.
func getMapSlice(m interface{}) yaml.MapSlice {
	if mm, ok := m.(yaml.MapSlice); ok {
		return mm
	}
	log.Fatalf("Entry %v is not a yaml.MapSlice", m)
	return nil
}

// appendIfUnique appends an element to an array of strings, unless it's already present.
func appendIfUnique(a1 []string, e2 string) []string {
	var res []string
	res = append(res, a1...)
	for _, e1 := range a1 {
		if e1 == e2 {
			return res
		}
	}
	return append(res, e2)
}

func combineSlices(a1 []string, a2 []string) []string {
	var res []string
	res = append(res, a1...)
	for _, e2 := range a2 {
		res = appendIfUnique(res, e2)
	}
	return res
}

// intersectSlices returns intersect of 2 slices
func intersectSlices(a1, a2 []string) []string {
	var res []string
	s1 := sets.NewString(a1...)
	for _, e2 := range a2 {
		if s1.Has(e2) {
			res = append(res, e2)
		}
	}
	return res
}

// exclusiveSlices returns elements in a1 but not in a2
func exclusiveSlices(a1, a2 []string) []string {
	var res []string
	s2 := sets.NewString(a2...)
	for _, e1 := range a1 {
		if !s2.Has(e1) {
			res = append(res, e1)
		}
	}
	return res
}

// Consolidate whitelisted and skipped branches with newly added
// whitelisted/skipped. To make the logic easier to maintain, this function
// makes the assumption that the outcome follows these rules:
//   - Special branch logics always apply on master and future branches
// Based on the previous rule, if there is a special branch logic, the 2 Prow
// jobs that serves different branches become:
//   - Standard job definition:
//		- whitelisted: [release-0.1]
//		- skipped: []
//   - Standard job definition + branch special logic #1:
//		- whitelisted: []
//		- skipped: [release-0.1]
// And when there is a new special logic comes up with different list of release
// branches to exclude, for example [release-0.1, release-0.2], then the desired
// outcome becomes:
//   - Standard job definition:
//		- whitelisted: [release-0.1]
//		- skipped: []
//   - Standard job definition + branch special logic #1: (This will never run)
//		- whitelisted: []
//		- skipped: []
//   - Standard job definition + branch special logic #2:
//		- whitelisted: [release-0.2]
//		- skipped: []
//   - Standard job definition + branch special logic #1 + branch special logic #2:
//		- whitelisted: []
//		- skipped: [release-0.1, release-0.2]
// Noted that only jobs with all special branch logics have something in
// skipped, while all other jobs only have whitelisted. This rule also applies
// when there is a third branch specific logic and so on.
// This function takes the logic above, and determines whether generate
// whitelisted or skipped as output.
func consolidateBranches(whitelisted []string, skipped []string, newWhitelisted []string, newSkipped []string) ([]string, []string) {
	var combinedWhitelisted, combinedSkipped []string

	// Do the legacy part(old branches):
	if len(newWhitelisted) > 0 {
		if len(skipped) > 0 {
			// - if previous is skipped(latest), then minus the skipped from current
			// branches, as we want to run exclusive on branches supported currently
			combinedWhitelisted = exclusiveSlices(newWhitelisted, skipped)
		} else if len(whitelisted) > 0 {
			// - if previous is include, then find their intersections, as these are the
			// real supported branches
			combinedWhitelisted = intersectSlices(newWhitelisted, whitelisted)
		} else {
			combinedWhitelisted = newWhitelisted
		}
	} else if len(newSkipped) > 0 { // Then do the pos part(latest)
		if len(skipped) > 0 {
			// - if previous is skipped(latest), then find the combination, as we want to
			// skip all non-supported
			combinedSkipped = combineSlices(newSkipped, skipped)
		} else if len(whitelisted) > 0 {
			// - if previous is include, then minus current branches from included
			combinedWhitelisted = exclusiveSlices(whitelisted, newSkipped)
		} else {
			combinedSkipped = newSkipped
		}
	}
	return combinedWhitelisted, combinedSkipped
}

// Config generation functions.

// newbaseProwJobTemplateData returns a baseProwJobTemplateData type with its initial, default values.
func newbaseProwJobTemplateData(repo string) baseProwJobTemplateData {
	var data baseProwJobTemplateData
	data.Timeout = 50
	data.OrgName = strings.Split(repo, "/")[0]
	data.RepoName = strings.Replace(repo, data.OrgName+"/", "", 1)
	data.RepoNameForJob = strings.ToLower(strings.Replace(repo, "/", "-", -1))
	data.RepoBranch = "master" // Default to be master, will override later for other branches
	data.GcsBucket = GCSBucket
	data.RepoURI = "github.com/" + repo
	data.CloneURI = fmt.Sprintf("\"https://%s.git\"", data.RepoURI)
	data.GcsLogDir = fmt.Sprintf("gs://%s/%s", GCSBucket, LogsDir)
	data.GcsPresubmitLogDir = fmt.Sprintf("gs://%s/%s", GCSBucket, presubmitLogsDir)
	data.ReleaseGcs = strings.Replace(repo, data.OrgName+"/", "knative-releases/", 1)
	data.AlwaysRun = true
	data.Image = prowTestsDockerImage
	data.ServiceAccount = testAccount
	data.Command = ""
	data.Args = make([]string, 0)
	data.Volumes = make([]string, 0)
	data.VolumeMounts = make([]string, 0)
	data.Env = make([]string, 0)
	data.ExtraRefs = []string{"- org: " + data.OrgName, "  repo: " + data.RepoName}
	data.Labels = make([]string, 0)
	data.Optional = ""

	// The build cluster for prod and staging Prow are both named
	// `build-knative` in kubeconfig in the cluster
	data.Cluster = "cluster: \"build-knative\""
	return data
}

// General helpers.

// createCommand returns an array with the command to run and its arguments.
func createCommand(data baseProwJobTemplateData) []string {
	c := []string{data.Command}
	// Prefix the pre-command if present.
	if preCommand != "" {
		c = append([]string{preCommand}, c...)
	}
	return append(c, data.Args...)
}

func envNameToKey(key string) string {
	return "- name: " + key
}

func envValueToValue(value string) string {
	return "  value: " + value
}

// addEnvToJob adds the given key/pair environment variable to the job.
func (data *baseProwJobTemplateData) addEnvToJob(key, value string) {
	// Value should always be string. Add quotes if we get a number
	if isNum(value) {
		value = "\"" + value + "\""
	}

	data.Env = append(data.Env, envNameToKey(key), envValueToValue(value))
}

// addLabelToJob adds extra labels to a job
func addLabelToJob(data *baseProwJobTemplateData, key, value string) {
	(*data).Labels = append((*data).Labels, []string{key + ": " + value}...)
}

// addPubsubLabelsToJob adds the pubsub labels so the prow job message will be picked up by test-infra monitoring
func addMonitoringPubsubLabelsToJob(data *baseProwJobTemplateData, runID string) {
	addLabelToJob(data, "prow.k8s.io/pubsub.project", "knative-tests")
	addLabelToJob(data, "prow.k8s.io/pubsub.topic", "knative-monitoring")
	addLabelToJob(data, "prow.k8s.io/pubsub.runID", runID)
}

// addVolumeToJob adds the given mount path as volume for the job.
func addVolumeToJob(data *baseProwJobTemplateData, mountPath, name string, isSecret bool, defaultMode string) {
	(*data).VolumeMounts = append((*data).VolumeMounts, []string{"- name: " + name, "  mountPath: " + mountPath}...)
	if isSecret {
		(*data).VolumeMounts = append((*data).VolumeMounts, "  readOnly: true")
	}
	s := []string{"- name: " + name}
	if isSecret {
		arr := []string{"  secret:", "    secretName: " + name}
		if len(defaultMode) > 0 {
			arr = append(arr, "    defaultMode: "+defaultMode)
		}
		s = append(s, arr...)
	} else {
		s = append(s, "  emptyDir: {}")
	}
	(*data).Volumes = append((*data).Volumes, s...)
}

// configureServiceAccountForJob adds the necessary volumes for the service account for the job.
func configureServiceAccountForJob(data *baseProwJobTemplateData) {
	if data.ServiceAccount == "" {
		return
	}
	p := strings.Split(data.ServiceAccount, "/")
	if len(p) != 4 || p[0] != "" || p[1] != "etc" || p[3] != "service-account.json" {
		log.Fatalf("Service account path %q is expected to be \"/etc/<name>/service-account.json\"", data.ServiceAccount)
	}
	name := p[2]
	addVolumeToJob(data, "/etc/"+name, name, true, "")
}

// addExtraEnvVarsToJob adds extra environment variables to a job.
func addExtraEnvVarsToJob(envVars []string, data *baseProwJobTemplateData) {
	for _, env := range envVars {
		pair := strings.SplitN(env, "=", 2)
		if len(pair) != 2 {
			log.Fatalf("Environment variable %q is expected to be \"key=value\"", env)
		}
		data.addEnvToJob(pair[0], pair[1])
	}
}

// setupDockerInDockerForJob enables docker-in-docker for the given job.
func setupDockerInDockerForJob(data *baseProwJobTemplateData) {
	addVolumeToJob(data, "/docker-graph", "docker-graph", false, "")
	data.addEnvToJob("DOCKER_IN_DOCKER_ENABLED", "\"true\"")
	(*data).SecurityContext = []string{"privileged: true"}
}

// setResourcesReqForJob sets resource requirement for job
func setResourcesReqForJob(res yaml.MapSlice, data *baseProwJobTemplateData) {
	data.Resources = nil
	for _, val := range res {
		data.Resources = append(data.Resources, fmt.Sprintf("  %s:", getString(val.Key)))
		for _, item := range getMapSlice(val.Value) {
			data.Resources = append(data.Resources, fmt.Sprintf("    %s: %s", getString(item.Key), getString(item.Value)))
		}
	}
}

// Config parsers.

// parseBasicJobConfigOverrides updates the given baseProwJobTemplateData with any base option present in the given config.
func parseBasicJobConfigOverrides(data *baseProwJobTemplateData, config yaml.MapSlice) {
	(*data).ExtraRefs = append((*data).ExtraRefs, "  base_ref: "+(*data).RepoBranch)
	var needDotdev bool
	for i, item := range config {
		switch item.Key {
		case "skip_branches":
			(*data).SkipBranches = getStringArray(item.Value)
		case "branches":
			(*data).Branches = getStringArray(item.Value)
		case "args":
			(*data).Args = getStringArray(item.Value)
		case "timeout":
			(*data).Timeout = getInt(item.Value)
		case "command":
			(*data).Command = getString(item.Value)
		case "needs-monitor":
			(*data).NeedsMonitor = true
		case "needs-dind":
			if getBool(item.Value) {
				setupDockerInDockerForJob(data)
			}
		case "always_run":
			(*data).AlwaysRun = getBool(item.Value)
		case "dot-dev":
			needDotdev = true
			for i, repo := range repositories {
				if path.Base(repo.Name) == (*data).RepoName {
					repositories[i].DotDev = true
				}
			}
		case "performance":
			for i, repo := range repositories {
				if path.Base(repo.Name) == (*data).RepoName {
					repositories[i].EnablePerformanceTests = true
				}
			}
		case "env-vars":
			addExtraEnvVarsToJob(getStringArray(item.Value), data)
		case "optional":
			(*data).Optional = "optional: true"
		case "resources":
			setResourcesReqForJob(getMapSlice(item.Value), data)
		case nil: // already processed
			continue
		default:
			log.Fatalf("Unknown entry %q for job", item.Key)
			continue
		}
		// Knock-out the item, signalling it was already parsed.
		config[i] = yaml.MapItem{}
	}

	if needDotdev {
		(*data).PathAlias = "path_alias: knative.dev/" + (*data).RepoName
		(*data).ExtraRefs = append((*data).ExtraRefs, "  "+(*data).PathAlias)
	}
	// Override any values if provided by command-line flags.
	if timeoutOverride > 0 {
		(*data).Timeout = timeoutOverride
	}
}

// getProwConfigData gets some basic, general data for the Prow config.
func getProwConfigData(config yaml.MapSlice) prowConfigTemplateData {
	var data prowConfigTemplateData
	data.Year = time.Now().Year()
	data.ProwHost = prowHost
	data.TestGridHost = testGridHost
	data.GubernatorHost = gubernatorHost
	data.GcsBucket = GCSBucket
	data.TestGridGcsBucket = testGridGcsBucket
	data.PresubmitLogsDir = presubmitLogsDir
	data.LogsDir = LogsDir
	data.TideRepos = make([]string, 0)
	data.ManagedRepos = make([]string, 0)
	data.ManagedOrgs = make([]string, 0)
	// Repos enabled for tide are all those that have presubmit jobs.
	var isProd bool
	for _, section := range config {
		if section.Key != "presubmits" {
			continue
		}
		for _, repo := range getMapSlice(section.Value) {
			orgRepoName := getString(repo.Key)
			data.TideRepos = appendIfUnique(data.TideRepos, orgRepoName)
			if strings.HasSuffix(orgRepoName, "test-infra") {
				data.TestInfraRepo = orgRepoName
			}
			if strings.HasPrefix(orgRepoName, "knative/") {
				isProd = true
			}
		}
	}
	// TODO: Remove all of these once prow core and plugins configs are not
	// generated any more
	if isProd {
		data.ManagedOrgs = []string{"knative", "knative-sandbox"}
		data.ManagedRepos = []string{"google/knative-gcp"}
		data.JobConfigPath = "config/prod/prow/jobs/*.yaml"
	} else {
		data.ManagedOrgs = []string{"knative-prow-robot"}
		data.JobConfigPath = "config/prod/staging/jobs/*.yaml"
	}
	// Sort repos to make output stable.
	sort.Strings(data.TideRepos)
	sort.Strings(data.ManagedOrgs)
	sort.Strings(data.ManagedRepos)
	return data
}

// parseSection generate the configs from a given section of the input yaml file.
func parseSection(config yaml.MapSlice, title string, generate sectionGenerator, finalize sectionGenerator) {
	for _, section := range config {
		if section.Key != title {
			continue
		}
		for _, repo := range getMapSlice(section.Value) {
			repoName := getString(repo.Key)
			for _, jobConfig := range getInterfaceArray(repo.Value) {
				generate(title, repoName, getMapSlice(jobConfig))
			}
			if finalize != nil {
				finalize(title, repoName, nil)
			}
		}
	}
}

// Template helpers.

// gitHubRepo returns the correct reference for the GitHub repository.
func gitHubRepo(data baseProwJobTemplateData) string {
	if repositoryOverride != "" {
		return repositoryOverride
	}
	s := data.RepoURI
	if data.RepoBranch != "" {
		s += "=" + data.RepoBranch
	}
	return s
}

// isNum checks if the given string is a valid number
func isNum(s string) bool {
	_, err := strconv.ParseFloat(s, 64)
	return err == nil
}

// quote returns the given string quoted if it's not a number, or not a key/value pair, or already quoted.
func quote(s string) string {
	if isNum(s) {
		return s
	}
	if strings.HasPrefix(s, "'") || strings.HasPrefix(s, "\"") || strings.Contains(s, ": ") || strings.HasSuffix(s, ":") {
		return s
	}
	return "\"" + s + "\""
}

// indentBase is a helper function which returns the given array indented.
func indentBase(indentation int, prefix string, indentFirstLine bool, array []string) string {
	s := ""
	if len(array) == 0 {
		return s
	}
	indent := strings.Repeat(" ", indentation)
	for i := 0; i < len(array); i++ {
		if i > 0 || indentFirstLine {
			s += indent
		}
		s += prefix + quote(array[i]) + "\n"
	}
	return s
}

// indentArray returns the given array indented, prefixed by "-".
func indentArray(indentation int, array []string) string {
	return indentBase(indentation, "- ", false, array)
}

// indentKeys returns the given array of key/value pairs indented.
func indentKeys(indentation int, array []string) string {
	return indentBase(indentation, "", false, array)
}

// indentSectionBase is a helper function which returns the given array of key/value pairs indented inside a section.
func indentSectionBase(indentation int, title string, prefix string, array []string) string {
	keys := indentBase(indentation, prefix, true, array)
	if keys == "" {
		return keys
	}
	return title + ":\n" + keys
}

// indentArraySection returns the given array indented inside a section.
func indentArraySection(indentation int, title string, array []string) string {
	return indentSectionBase(indentation, title, "- ", array)
}

// indentSection returns the given array of key/value pairs indented inside a section.
func indentSection(indentation int, title string, array []string) string {
	return indentSectionBase(indentation, title, "", array)
}

// indentMap returns the given map indented, with each key/value separated by ": "
func indentMap(indentation int, mp map[string]string) string {
	// Extract map keys to keep order consistent.
	keys := make([]string, 0, len(mp))
	for key := range mp {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	arr := make([]string, len(mp))
	for i := 0; i < len(mp); i++ {
		arr[i] = keys[i] + ": " + quote(mp[keys[i]])
	}
	return indentBase(indentation, "", false, arr)
}

// outputConfig outputs the given line, if not empty, to stdout.
func outputConfig(line string) {
	if strings.TrimSpace(line) != "" {
		fmt.Fprintln(output, strings.TrimRight(line, " "))
		emittedOutput = true
	}
}

// strExists checks if the given string exists in the array
func strExists(arr []string, str string) bool {
	for _, s := range arr {
		if str == s {
			return true
		}
	}
	return false
}

type specialBranchLogic struct {
	branches []string
	// create new job data based on branches
	opsNew  func(*baseProwJobTemplateData)
	restore func(*baseProwJobTemplateData)
}

// getBase casts data into baseProwJobTemplateData and returns it
func getBase(data interface{}) *baseProwJobTemplateData {
	var base *baseProwJobTemplateData
	switch v := data.(type) {
	case *presubmitJobTemplateData:
		base = &data.(*presubmitJobTemplateData).Base
	case *postsubmitJobTemplateData:
		base = &data.(*postsubmitJobTemplateData).Base
	default:
		log.Fatalf("Unrecognized job template type: '%v'", v)
	}
	return base
}

// recursiveSBL recursively going through specialBranchLogic, and generate job
// at last. Use `i` to keeps track of current index in sbs to be used
func recursiveSBL(repoName string, data interface{}, generateOneJob func(data interface{}), sbs []specialBranchLogic, i int) {
	// Base case, all special branch logics have been applied
	if i == len(sbs) {
		// If there is no branch left, this job shouldn't be generated at all
		if len(getBase(data).Branches) > 0 || len(getBase(data).SkipBranches) > 0 {
			generateOneJob(data)
		}
		return
	}

	sb := sbs[i]
	base := getBase(data)

	origBranches, origSkipBranches := base.Branches, base.SkipBranches
	// Do legacy branches first
	base.Branches, base.SkipBranches = consolidateBranches(origBranches, origSkipBranches, sb.branches, []string{})
	recursiveSBL(repoName, data, generateOneJob, sbs, i+1)
	// Then do latest branches
	base.Branches, base.SkipBranches = consolidateBranches(origBranches, origSkipBranches, []string{}, sb.branches)
	sb.opsNew(base)
	recursiveSBL(repoName, data, generateOneJob, sbs, i+1)
	sb.restore(base)
}

// executeJobTemplateWrapper takes in consideration of repo settings, decides how many variants of the
// same job needs to be generated and generates them.
func executeJobTemplateWrapper(repoName string, data interface{}, generateOneJob func(data interface{})) {
	switch data.(type) {
	case *postsubmitJobTemplateData:
		if strings.HasSuffix(data.(*postsubmitJobTemplateData).PostsubmitJobName, "go-coverage") {
			generateOneJob(data)
			return
		}
	}
	generateOneJob(data)
}

// executeTemplate outputs the given job template with the given data, respecting any filtering.
func executeJobTemplate(name, templ, title, repoName, jobName string, groupByRepo bool, data interface{}) {
	if jobNameFilter != "" && jobNameFilter != jobName {
		return
	}
	if !sectionMap[title] {
		outputConfig(title + ":")
		sectionMap[title] = true
	}
	if groupByRepo {
		if !sectionMap[title+repoName] {
			outputConfig(baseIndent + repoName + ":")
			sectionMap[title+repoName] = true
		}
	}
	executeTemplate(name, templ, data)
}

// executeTemplate outputs the given template with the given data.
func executeTemplate(name, templ string, data interface{}) {
	var res bytes.Buffer
	funcMap := template.FuncMap{
		"indent_section":       indentSection,
		"indent_array_section": indentArraySection,
		"indent_array":         indentArray,
		"indent_keys":          indentKeys,
		"indent_map":           indentMap,
		"repo":                 gitHubRepo,
	}
	t := template.Must(template.New(name).Funcs(funcMap).Delims("[[", "]]").Parse(templ))
	if err := t.Execute(&res, data); err != nil {
		log.Fatalf("Error in template %s: %v", name, err)
	}
	for _, line := range strings.Split(res.String(), "\n") {
		outputConfig(line)
	}
}

// Multi-value flag parser.

func (a *stringArrayFlag) String() string {
	return strings.Join(*a, ", ")
}

func (a *stringArrayFlag) Set(value string) error {
	*a = append(*a, value)
	return nil
}

// parseJob gets the job data from the original yaml data, now the jobName can be "presubmits" or "periodic"
func parseJob(config yaml.MapSlice, jobName string) yaml.MapSlice {
	for _, section := range config {
		if section.Key == jobName {
			return getMapSlice(section.Value)
		}
	}

	log.Fatalf("The metadata misses %s configuration, cannot continue.", jobName)
	return nil
}

// parseGoCoverageMap constructs a map, indicating which repo is enabled for go coverage check
func parseGoCoverageMap(presubmitJob yaml.MapSlice) map[string]bool {
	goCoverageMap := make(map[string]bool)
	for _, repo := range presubmitJob {
		repoName := strings.Split(getString(repo.Key), "/")[1]
		goCoverageMap[repoName] = false
		for _, jobConfig := range getInterfaceArray(repo.Value) {
			for _, item := range getMapSlice(jobConfig) {
				if item.Key == "go-coverage" {
					goCoverageMap[repoName] = getBool(item.Value)
					break
				}
			}
		}
	}

	return goCoverageMap
}

// collectMetaData collects the meta data from the original yaml data, which can be then used for building the test groups and dashboards config
func collectMetaData(periodicJob yaml.MapSlice) {
	for _, repo := range periodicJob {
		rawName := getString(repo.Key)
		projName := strings.Split(rawName, "/")[0]
		repoName := strings.Split(rawName, "/")[1]
		jobDetailMap := metaData.Get(projName)
		metaData.EnsureRepo(projName, repoName)

		// parse job configs
		for _, conf := range getInterfaceArray(repo.Value) {
			jobDetailMap = metaData.Get(projName)
			jobConfig := getMapSlice(conf)
			enabled := false
			jobName := ""
			releaseVersion := ""
			for _, item := range jobConfig {
				switch item.Key {
				case "continuous", "dot-release", "auto-release", "performance",
					"nightly", "webhook-apicoverage":
					if getBool(item.Value) {
						enabled = true
						jobName = getString(item.Key)
					}
				case "branch-ci":
					enabled = getBool(item.Value)
					jobName = "continuous"
				case "release":
					releaseVersion = getString(item.Value)
				case "custom-job":
					enabled = true
					jobName = getString(item.Value)
				default:
					// continue here since we do not need to care about other entries, like cron, command, etc.
					continue
				}
			}
			// add job types for the corresponding repos, if needed
			if enabled {
				// if it's a job for a release branch
				if releaseVersion != "" {
					releaseProjName := fmt.Sprintf("%s-%s", projName, releaseVersion)

					// TODO: Why do we assign?
					jobDetailMap = metaData.Get(releaseProjName)
				}
				jobDetailMap.Add(repoName, jobName)
			}
		}
		updateTestCoverageJobDataIfNeeded(jobDetailMap, repoName)
	}

	// add test coverage jobs for the repos that haven't been handled
	addRemainingTestCoverageJobs()
}

// updateTestCoverageJobDataIfNeeded adds test-coverage job data for the repo if it has go coverage check
func updateTestCoverageJobDataIfNeeded(jobDetailMap JobDetailMap, repoName string) {
	if goCoverageMap[repoName] {
		jobDetailMap.Add(repoName, "test-coverage")
		// delete this repoName from the goCoverageMap to avoid it being processed again when we
		// call the function addRemainingTestCoverageJobs
		delete(goCoverageMap, repoName)
	}
}

// addRemainingTestCoverageJobs adds test-coverage jobs data for the repos that haven't been processed.
func addRemainingTestCoverageJobs() {
	// handle repos that only have go coverage
	for repoName, hasGoCoverage := range goCoverageMap {
		if hasGoCoverage {
			jobDetailMap := metaData.Get(metaData.projNames[0]) // TODO: WTF why projNames[0] !??!?!?!?
			jobDetailMap.Add(repoName, "test-coverage")
		}
	}
}

// buildProjRepoStr builds the projRepoStr used in the config file with projName and repoName
func buildProjRepoStr(projName string, repoName string) string {
	projVersion := ""
	if strings.Contains(projName, "-") {
		projNameAndVersion := strings.Split(projName, "-")
		projName = projNameAndVersion[0]
		projVersion = projNameAndVersion[1]
	}
	projRepoStr := repoName
	if projVersion != "" {
		projRepoStr += ("-" + projVersion)
	}
	projRepoStr = projName + "-" + projRepoStr
	return strings.ToLower(projRepoStr)
}

// isReleased returns true for project name that has version
func isReleased(projName string) bool {
	return releaseRegex.FindString(projName) != ""
}

// setOutput set the given file as the output target, then all the output will be written to this file
func setOutput(fileName string) {
	output = os.Stdout
	if fileName == "" {
		return
	}
	configFile, err := os.OpenFile(fileName, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		log.Fatalf("Cannot create the configuration file %q: %v", fileName, err)
	}
	configFile.Truncate(0)
	configFile.Seek(0, 0)
	output = configFile
}

// main is the script entry point.
func main() {
	// Parse flags and sanity check them.
	prowConfigOutput := ""
	prowJobsConfigOutput := ""
	testgridConfigOutput := ""
	var env = flag.String("env", "prow", "The name of the environment, can be prow or prow-staging")
	var generateProwConfig = flag.Bool("generate-prow-config", true, "Whether to generate the prow configuration file from the template")
	var generatePluginsConfig = flag.Bool("generate-plugins-config", true, "Whether to generate the plugins configuration file from the template")
	var generateTestgridConfig = flag.Bool("generate-testgrid-config", true, "Whether to generate the testgrid config from the template file")
	var generateMaintenanceJobs = flag.Bool("generate-maintenance-jobs", true, "Whether to generate the maintenance periodic jobs (e.g. backup)")
	// TODO: remove these options and just generate everything
	if !(*generateProwConfig && *generatePluginsConfig && *generateTestgridConfig && *generateMaintenanceJobs) {
		panic(errors.New("Must enable all generators"))
	}
	var includeConfig = flag.Bool("include-config", true, "Whether to include general configuration (e.g., plank) in the generated config")
	var dockerImagesBase = flag.String("image-docker", "gcr.io/knative-tests/test-infra", "Default registry for the docker images used by the jobs")
	flag.StringVar(&prowConfigOutput, "prow-config-output", "", "The destination for the prow config output, default to be stdout")
	flag.StringVar(&prowJobsConfigOutput, "prow-jobs-config-output", "", "The destination for the prow jobs config output, default to be stdout")
	var pluginsConfigOutput = flag.String("plugins-config-output", "", "The destination for the plugins config output, default to be stdout")
	flag.StringVar(&testgridConfigOutput, "testgrid-config-output", "", "The destination for the testgrid config output, default to be stdout")
	flag.StringVar(&prowHost, "prow-host", "https://prow.knative.dev", "Prow host, including HTTP protocol")
	flag.StringVar(&testGridHost, "testgrid-host", "https://testgrid.knative.dev", "TestGrid host, including HTTP protocol")
	flag.StringVar(&gubernatorHost, "gubernator-host", "https://gubernator.knative.dev", "Gubernator host, including HTTP protocol")
	flag.StringVar(&GCSBucket, "gcs-bucket", "knative-prow", "GCS bucket to upload the logs to")
	flag.StringVar(&testGridGcsBucket, "testgrid-gcs-bucket", "knative-testgrid", "TestGrid GCS bucket")
	flag.StringVar(&LogsDir, "logs-dir", "logs", "Path in the GCS bucket to upload logs of periodic and post-submit jobs")
	flag.StringVar(&presubmitLogsDir, "presubmit-logs-dir", "pr-logs", "Path in the GCS bucket to upload logs of pre-submit jobs")
	flag.StringVar(&testAccount, "test-account", "/etc/test-account/service-account.json", "Path to the service account JSON for test jobs")
	flag.StringVar(&nightlyAccount, "nightly-account", "/etc/nightly-account/service-account.json", "Path to the service account JSON for nightly release jobs")
	flag.StringVar(&releaseAccount, "release-account", "/etc/release-account/service-account.json", "Path to the service account JSON for release jobs")
	var prowTestsDockerImageName = flag.String("prow-tests-docker", "prow-tests:stable", "prow-tests docker image")
	flag.StringVar(&githubCommenterDockerImage, "github-commenter-docker", "gcr.io/k8s-prow/commenter:v20190731-e3f7b9853", "github commenter docker image")
	flag.StringVar(&presubmitScript, "presubmit-script", "./test/presubmit-tests.sh", "Executable for running presubmit tests")
	flag.StringVar(&releaseScript, "release-script", "./hack/release.sh", "Executable for creating releases")
	flag.StringVar(&webhookAPICoverageScript, "webhook-api-coverage-script", "./test/apicoverage.sh", "Executable for running webhook apicoverage tool")
	flag.StringVar(&repositoryOverride, "repo-override", "", "Repository path (github.com/foo/bar[=branch]) to use instead for a job")
	flag.IntVar(&timeoutOverride, "timeout-override", 0, "Timeout (in minutes) to use instead for a job")
	flag.StringVar(&jobNameFilter, "job-filter", "", "Generate only this job, instead of all jobs")
	flag.StringVar(&preCommand, "pre-command", "", "Executable for running instead of the real command of a job")
	flag.Var(&extraEnvVars, "extra-env", "Extra environment variables (key=value) to add to a job")
	flag.Parse()
	if len(flag.Args()) != 1 {
		log.Fatal("Pass the config file as parameter")
	}

	prowTestsDockerImage = path.Join(*dockerImagesBase, *prowTestsDockerImageName)

	// We use MapSlice instead of maps to keep key order and create predictable output.
	config := yaml.MapSlice{}

	// Read input config.
	name := flag.Arg(0)
	content, err := ioutil.ReadFile(name)
	if err != nil {
		log.Fatalf("Cannot read file %q: %v", name, err)
	}
	if err = yaml.Unmarshal(content, &config); err != nil {
		log.Fatalf("Cannot parse config %q: %v", name, err)
	}

	prowConfigData := getProwConfigData(config)

	if *generatePluginsConfig {
		setOutput(*pluginsConfigOutput)
		executeTemplate("plugins config", readTemplate(filepath.Join(*env, pluginsConfig)), prowConfigData)
	}

	// Generate Prow config.
	if *generateProwConfig {
		setOutput(prowConfigOutput)
		repositories = make([]repositoryData, 0)
		sectionMap = make(map[string]bool)
		if *includeConfig {
			executeTemplate("general header", readTemplate(commonHeaderConfig), prowConfigData)
			executeTemplate("general config", readTemplate(filepath.Join(*env, generalProwConfig)), prowConfigData)
		}

		setOutput(prowJobsConfigOutput)
		executeTemplate("general header", readTemplate(commonHeaderConfig), prowConfigData)
		parseSection(config, "presubmits", generatePresubmit, nil)
		parseSection(config, "periodics", generatePeriodic, generateGoCoveragePeriodic)
		for _, repo := range repositories { // Keep order for predictable output.
			if !repo.Processed && repo.EnableGoCoverage {
				generateGoCoveragePeriodic("periodics", repo.Name, nil)
			}
		}
		generatePerfClusterUpdatePeriodicJobs()
		if *generateMaintenanceJobs {
			generateIssueTrackerPeriodicJobs()
		}
		for _, repo := range repositories {
			if repo.EnableGoCoverage {
				generateGoCoveragePostsubmit("postsubmits", repo.Name, nil)
			}
			if repo.EnablePerformanceTests {
				generatePerfClusterPostsubmitJob(repo)
			}
		}
	}

	// config object is modified when we generate prow config, so we'll need to reload it here
	if err = yaml.Unmarshal(content, &config); err != nil {
		log.Fatalf("Cannot parse config %q: %v", name, err)
	}
	// Generate Testgrid config.
	if *generateTestgridConfig {
		setOutput(testgridConfigOutput)

		if *includeConfig {
			executeTemplate("general header", readTemplate(commonHeaderConfig), newBaseTestgridTemplateData(""))
			executeTemplate("general config", readTemplate(generalTestgridConfig), newBaseTestgridTemplateData(""))
		}

		presubmitJobData := parseJob(config, "presubmits")
		goCoverageMap = parseGoCoverageMap(presubmitJobData)

		periodicJobData := parseJob(config, "periodics")
		collectMetaData(periodicJobData)

		// log.Print(spew.Sdump(metaData))

		// These generate "test_groups:"
		metaData.generateTestGridSection("test_groups", metaData.generateTestGroup, false)
		metaData.generateNonAlignedTestGroups()

		// These generate "dashboards:"
		metaData.generateTestGridSection("dashboards", generateDashboard, true)
		metaData.generateDashboardsForReleases()
		metaData.generateNonAlignedDashboards()

		// These generate "dashboard_groups:"
		metaData.generateDashboardGroups()
		metaData.generateNonAlignedDashboardGroups()
	}
}
