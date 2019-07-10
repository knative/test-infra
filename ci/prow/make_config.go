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
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"os"
	"path"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"text/template"
	"time"

	yaml "gopkg.in/yaml.v2"
)

const (
	// Manifests generated by ko are indented by 2 spaces.
	baseIndent  = "  "
	templateDir = "templates"

	// ##########################################################
	// ############## prow configuration templates ##############
	// ##########################################################
	// generalProwConfig contains config-wide definitions.
	generalProwConfig = "prow_config_header.yaml"

	// presubmitJob is the template for presubmit jobs.
	presubmitJob = "prow_presubmit_job.yaml"

	// presubmitGoCoverageJob is the template for go coverage presubmit jobs.
	presubmitGoCoverageJob = "prow_presubmit_gocoverate_job.yaml"

	// goCoveragePostsubmitJob is the template for the go postsubmit coverage job.
	goCoveragePostsubmitJob = "prow_postsubmit_gocoverage_job.yaml"
)

// repositoryData contains basic data about each Knative repository.
type repositoryData struct {
	Name                string
	EnableGoCoverage    bool
	GoCoverageThreshold int
	Processed           bool
	DotDev              bool
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
	DecorationConfig    []string
	ExtraRefs           []string
	Command             string
	Args                []string
	Env                 []string
	Volumes             []string
	VolumeMounts        []string
	Timeout             int
	AlwaysRun           bool
	LogsDir             string
	PresubmitLogsDir    string
	TestAccount         string
	ServiceAccount      string
	ReleaseGcs          string
	GoCoverageThreshold int
	Image               string
	Year                int
	Labels              []string
	PathAlias           string
}

// ####################################################################################################
// ################ data definitions that are used for the prow config file generation ################
// ####################################################################################################
// presubmitJobTemplateData contains data about a presubmit Prow job.
type presubmitJobTemplateData struct {
	Base                 baseProwJobTemplateData
	PresubmitJobName     string
	PresubmitPullJobName string
	PresubmitPostJobName string
	PresubmitCommand     []string
}

// postsubmitJobTemplateData contains data about a postsubmit Prow job.
type postsubmitJobTemplateData struct {
	Base              baseProwJobTemplateData
	PostsubmitJobName string
}

// sectionGenerator is a function that generates Prow job configs given a slice of a yaml file with configs.
type sectionGenerator func(string, string, yaml.MapSlice)

// newJobNeeded is a function that determined if we need to add a new job for this repository.
type newJobNeeded func(repositoryData) bool

// stringArrayFlag is the content of a multi-value flag.
type stringArrayFlag []string

var (
	// Values used in the jobs that can be changed through command-line flags.
	output                       *os.File
	gcsBucket                    string
	logsDir                      string
	presubmitLogsDir             string
	testAccount                  string
	nightlyAccount               string
	releaseAccount               string
	flakesreporterDockerImage    string
	prowversionbumperDockerImage string
	coverageDockerImage          string
	clearalertsDockerImage       string
	prowTestsDockerImage         string
	presubmitScript              string
	releaseScript                string
	performanceScript            string
	webhookAPICoverageScript     string
	cleanupScript                string

	// #########################################################################
	// ############## data used for generating prow configuration ##############
	// #########################################################################
	// Array constants used throughout the jobs.
	allPresubmitTests = []string{"--all-tests", "--emit-metrics"}
	releaseNightly    = []string{"--publish", "--tag-release"}
	releaseLocal      = []string{"--nopublish", "--notag-release"}

	// Overrides and behavior changes through command-line flags.
	repositoryOverride string
	jobNameFilter      string
	preCommand         string
	extraEnvVars       stringArrayFlag
	timeoutOverride    int

	// List of Knative repositories.
	repositories []repositoryData

	// Map which sections of the config.yaml were written to stdout.
	sectionMap map[string]bool
)

// Generate cron string based on job type, offset generated from jobname
// instead of assign random value to ensure consistency among runs,
// timeout is used for determining how many hours apart
func generateCron(jobType, jobName string, timeout int) string {
	getUTCtime := func(i int) int { return i + 7 }
	// Sums the ascii valus of all letters in a jobname,
	// this value is used for deriving offset after hour
	var sum float64
	for _, c := range jobType + jobName {
		sum += float64(c)
	}
	// Divide 60 minutes into 6 buckets
	bucket := int(math.Mod(sum, 6))
	// Offset in bucket, range from 0-9, first mod with 11(a random prime number)
	// to ensure every digit has a chance (i.e., if bucket is 0, sum has to be multiply of 6,
	// so mod by 10 can only return even number)
	offsetInBucket := int(math.Mod(math.Mod(sum, 11), 10))
	minutesOffset := bucket*10 + offsetInBucket
	// Determines hourly job inteval based on timeout
	hours := int((timeout+5)/60) + 1 // Allow at least 5 minutes between runs
	hourCron := fmt.Sprintf("%d * * * *", minutesOffset)
	if hours > 1 {
		hourCron = fmt.Sprintf("%d */%d * * *", minutesOffset, hours)
	}
	dayCron := fmt.Sprintf("%d %%d * * *", minutesOffset)    // hour
	weekCron := fmt.Sprintf("%d %%d * * %%d", minutesOffset) // hour, weekday

	var res string
	switch jobType {
	case "continuous", "custom-job", "auto-release": // Every hour
		res = fmt.Sprintf(hourCron)
	case "branch-ci": // Every day 1-2 PST
		res = fmt.Sprintf(dayCron, getUTCtime(1))
	case "nightly": // Every day 2-3 PST
		res = fmt.Sprintf(dayCron, getUTCtime(2))
	case "dot-release": // Every Tuesday 2-3 PST
		res = fmt.Sprintf(weekCron, getUTCtime(2), 2)
	case "latency": // Every day 1-2 PST
		res = fmt.Sprintf(dayCron, getUTCtime(1))
	case "performance": // Every day 1-2 PST
		res = fmt.Sprintf(dayCron, getUTCtime(1))
	case "performance-mesh": // Every day 3-4 PST
		res = fmt.Sprintf(dayCron, getUTCtime(3))
	case "webhook-apicoverage": // Every day 2-3 PST
		res = fmt.Sprintf(dayCron, getUTCtime(2))
	default:
		log.Printf("job type not supported for cron generation '%s'", jobName)
	}
	return res
}

// Yaml parsing helpers.

// read template yaml file content
func readTemplate(fp string) string {
	if _, ok := templatesCache[fp]; !ok {
		content, err := ioutil.ReadFile(path.Join(templateDir, fp))
		if nil != err {
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

// Config generation functions.

// newbaseProwJobTemplateData returns a baseProwJobTemplateData type with its initial, default values.
func newbaseProwJobTemplateData(repo string) baseProwJobTemplateData {
	var data baseProwJobTemplateData
	data.Timeout = 50
	data.OrgName = strings.Split(repo, "/")[0]
	data.RepoName = strings.Replace(repo, data.OrgName+"/", "", 1)
	data.RepoNameForJob = strings.ToLower(strings.Replace(repo, "/", "-", -1))
	data.RepoBranch = "master" // Default to be master, will override later for other branches
	data.GcsBucket = gcsBucket
	data.RepoURI = "github.com/" + repo
	data.CloneURI = fmt.Sprintf("\"https://%s.git\"", data.RepoURI)
	data.GcsLogDir = fmt.Sprintf("gs://%s/%s", gcsBucket, logsDir)
	data.GcsPresubmitLogDir = fmt.Sprintf("gs://%s/%s", gcsBucket, presubmitLogsDir)
	data.Year = time.Now().Year()
	data.PresubmitLogsDir = presubmitLogsDir
	data.LogsDir = logsDir
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

// addEnvToJob adds the given key/pair environment variable to the job.
func addEnvToJob(data *baseProwJobTemplateData, key, value string) {
	// Value should always be string. Add quotes if we get a number
	if isNum(value) {
		value = "\"" + value + "\""
	}

	(*data).Env = append((*data).Env, []string{"- name: " + key, "  value: " + value}...)
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

// addExtraEnvVarsToJob adds any extra environment variables (defined on command-line) to a job.
func addExtraEnvVarsToJob(data *baseProwJobTemplateData) {
	for _, env := range extraEnvVars {
		pair := strings.Split(env, "=")
		if len(pair) != 2 {
			log.Fatalf("Environment variable %q is expected to be \"key=value\"", env)
		}
		addEnvToJob(data, pair[0], pair[1])
	}
}

// setupDockerInDockerForJob enables docker-in-docker for the given job.
func setupDockerInDockerForJob(data *baseProwJobTemplateData) {
	addVolumeToJob(data, "/docker-graph", "docker-graph", false, "")
	addEnvToJob(data, "DOCKER_IN_DOCKER_ENABLED", "\"true\"")
	(*data).SecurityContext = []string{"privileged: true"}
}

// Config parsers.

// parseBasicJobConfigOverrides updates the given baseProwJobTemplateData with any base option present in the given config.
func parseBasicJobConfigOverrides(data *baseProwJobTemplateData, config yaml.MapSlice) {
	(*data).ExtraRefs = append((*data).ExtraRefs, "  base_ref: "+(*data).RepoBranch)
	for i, item := range config {
		switch item.Key {
		case "skip_branches":
			(*data).SkipBranches = getStringArray(item.Value)
		case "args":
			(*data).Args = getStringArray(item.Value)
		case "timeout":
			(*data).Timeout = getInt(item.Value)
		case "command":
			(*data).Command = getString(item.Value)
		case "full-command":
			parts := strings.Split(getString(item.Value), " ")
			(*data).Command = parts[0]
			(*data).Args = parts[1:]
		case "needs-dind":
			if getBool(item.Value) {
				setupDockerInDockerForJob(data)
			}
		case "always_run":
			(*data).AlwaysRun = getBool(item.Value)
		case "dot-dev":
			for i, repo := range repositories {
				if path.Base(repo.Name) == (*data).RepoName {
					repositories[i].DotDev = true
				}
			}
		case nil: // already processed
			continue
		default:
			log.Fatalf("Unknown entry %q for job", item.Key)
			continue
		}
		// Knock-out the item, signalling it was already parsed.
		config[i] = yaml.MapItem{}
	}
	for _, repo := range repositories {
		if path.Base(repo.Name) == (*data).RepoName && repo.DotDev {
			(*data).PathAlias = "path_alias: knative.dev/" + (*data).RepoName
			(*data).ExtraRefs = append((*data).ExtraRefs, "  "+(*data).PathAlias)
			break
		}
	}
	// Override any values if provided by command-line flags.
	if timeoutOverride > 0 {
		(*data).Timeout = timeoutOverride
	}
}

// generatePresubmit generates all presubmit job configs for the given repo and configuration.
func generatePresubmit(title string, repoName string, presubmitConfig yaml.MapSlice) {
	var data presubmitJobTemplateData
	data.Base = newbaseProwJobTemplateData(repoName)
	data.Base.Command = presubmitScript
	data.Base.GoCoverageThreshold = 50
	jobTemplate := readTemplate(presubmitJob)
	repoData := repositoryData{Name: repoName, EnableGoCoverage: false, GoCoverageThreshold: data.Base.GoCoverageThreshold}
	isMonitoredJob := false
	for i, item := range presubmitConfig {
		switch item.Key {
		case "build-tests", "unit-tests", "integration-tests":
			if !getBool(item.Value) {
				return
			}
			jobName := getString(item.Key)
			data.PresubmitJobName = data.Base.RepoNameForJob + "-" + jobName
			// Use default arguments if none given.
			if len(data.Base.Args) == 0 {
				data.Base.Args = []string{"--" + jobName}
			}

			if item.Key == "integration-tests" {
				isMonitoredJob = true
			}
		case "go-coverage":
			if !getBool(item.Value) {
				return
			}
			jobTemplate = readTemplate(presubmitGoCoverageJob)
			data.PresubmitJobName = data.Base.RepoNameForJob + "-go-coverage"
			data.Base.Image = coverageDockerImage
			data.Base.ServiceAccount = ""
			repoData.EnableGoCoverage = true
			addVolumeToJob(&data.Base, "/etc/covbot-token", "covbot-token", true, "")
		case "custom-test":
			data.PresubmitJobName = data.Base.RepoNameForJob + "-" + getString(item.Value)
		case "go-coverage-threshold":
			data.Base.GoCoverageThreshold = getInt(item.Value)
			repoData.GoCoverageThreshold = data.Base.GoCoverageThreshold
		default:
			continue
		}
		// Knock-out the item, signalling it was already parsed.
		presubmitConfig[i] = yaml.MapItem{}
	}
	repositories = append(repositories, repoData)
	parseBasicJobConfigOverrides(&data.Base, presubmitConfig)
	data.PresubmitCommand = createCommand(data.Base)
	data.PresubmitPullJobName = "pull-" + data.PresubmitJobName
	data.PresubmitPostJobName = "post-" + data.PresubmitJobName
	if data.Base.ServiceAccount != "" {
		addEnvToJob(&data.Base, "GOOGLE_APPLICATION_CREDENTIALS", data.Base.ServiceAccount)
		addEnvToJob(&data.Base, "E2E_CLUSTER_REGION", "us-central1")
	}
	if isMonitoredJob {
		addMonitoringPubsubLabelsToJob(&data.Base, data.PresubmitPullJobName)
	}
	addExtraEnvVarsToJob(&data.Base)
	configureServiceAccountForJob(&data.Base)
	executeJobTemplate("presubmit", jobTemplate, title, repoName, data.PresubmitPullJobName, true, data)
	// TODO(adrcunha): remove once the coverage-dev job isn't necessary anymore.
	// Generate config for pull-knative-serving-go-coverage-dev right after pull-knative-serving-go-coverage
	if data.PresubmitPullJobName == "pull-knative-serving-go-coverage" {
		data.PresubmitPullJobName += "-dev"
		data.Base.AlwaysRun = false
		data.Base.Image = strings.Replace(data.Base.Image, "coverage:latest", "coverage-dev:latest-dev", -1)
		template := strings.Replace(readTemplate(presubmitGoCoverageJob), "(all|", "(", 1)
		executeJobTemplate("presubmit", template, title, repoName, data.PresubmitPullJobName, true, data)
	}
}

// generateGoCoveragePostsubmit generates the go coverage postsubmit job config for the given repo.
func generateGoCoveragePostsubmit(title, repoName string, _ yaml.MapSlice) {
	var data postsubmitJobTemplateData
	data.Base = newbaseProwJobTemplateData(repoName)
	data.Base.Image = coverageDockerImage
	data.PostsubmitJobName = fmt.Sprintf("post-%s-go-coverage", data.Base.RepoNameForJob)
	for _, repo := range repositories {
		if repo.Name == repoName && repo.DotDev {
			data.Base.PathAlias = "path_alias: knative.dev/" + path.Base(repoName)
		}
	}
	addExtraEnvVarsToJob(&data.Base)
	configureServiceAccountForJob(&data.Base)
	executeJobTemplate("postsubmit go coverage", readTemplate(goCoveragePostsubmitJob), "postsubmits", repoName, data.PostsubmitJobName, true, data)
	// TODO(adrcunha): remove once the coverage-dev job isn't necessary anymore.
	// Generate config for post-knative-serving-go-coverage-dev right after post-knative-serving-go-coverage
	if data.PostsubmitJobName == "post-knative-serving-go-coverage" {
		data.PostsubmitJobName += "-dev"
		data.Base.Image = strings.Replace(data.Base.Image, "coverage:latest", "coverage-dev:latest-dev", -1)
		executeJobTemplate("presubmit", readTemplate(goCoveragePostsubmitJob), "postsubmits", repoName, data.PostsubmitJobName, false, data)
	}
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

// generateOtherJobConfigs generates job config with the generator if new job is required for it.
func generateOtherJobConfigs(title string, newJobNeeded newJobNeeded, generate sectionGenerator) {
	for i := range repositories { // Keep order for predictable output.
		if !newJobNeeded(repositories[i]) {
			continue
		}
		generate(title, repositories[i].Name, nil)
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
	s := strings.TrimSpace(line)
	if s != "" {
		fmt.Fprintln(output, line)
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
		jobDetailMap := addProjAndRepoIfNeed(projName, repoName)

		// parse job configs
		for _, conf := range getInterfaceArray(repo.Value) {
			jobDetailMap = metaData[projName]
			jobConfig := getMapSlice(conf)
			enabled := false
			jobName := ""
			releaseVersion := ""
			for _, item := range jobConfig {
				switch item.Key {
				case "continuous", "dot-release", "auto-release", "performance", "performance-mesh", "latency", "nightly":
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
					jobDetailMap = addProjAndRepoIfNeed(releaseProjName, repoName)
				}
				newJobTypes := append(jobDetailMap[repoName], jobName)
				jobDetailMap[repoName] = newJobTypes
			}
		}
		addTestCoverageJobIfNeeded(&jobDetailMap, repoName)
	}

	// add test coverage jobs for the repos that haven't been handled
	addRemainingTestCoverageJobs()
}

// addProjAndRepoIfNeed adds the project and repo if they are new in the metaData map, then return the jobDetailMap
func addProjAndRepoIfNeed(projName string, repoName string) map[string][]string {
	// add project in the metaData
	if _, exists := metaData[projName]; !exists {
		metaData[projName] = make(map[string][]string)
		if !strExists(projNames, projName) {
			projNames = append(projNames, projName)
		}
	}

	// add repo in the project
	jobDetailMap := metaData[projName]
	if _, exists := jobDetailMap[repoName]; !exists {
		if !strExists(repoNames, repoName) {
			repoNames = append(repoNames, repoName)
		}
		jobDetailMap[repoName] = make([]string, 0)
	}
	return jobDetailMap
}

// addTestCoverageJobIfNeeded adds test-coverage job for the repo if it has go coverage check
func addTestCoverageJobIfNeeded(jobDetailMap *map[string][]string, repoName string) {
	if goCoverageMap[repoName] {
		newJobTypes := append((*jobDetailMap)[repoName], "test-coverage")
		(*jobDetailMap)[repoName] = newJobTypes
		// delete this repoName from the goCoverageMap to avoid it being processed again when we
		// call the function addRemainingTestCoverageJobs
		delete(goCoverageMap, repoName)
	}
}

// addRemainingTestCoverageJobs adds test-coverage jobs for the repos that haven't been processed.
func addRemainingTestCoverageJobs() {
	// handle repos that only have go coverage
	for repoName, hasGoCoverage := range goCoverageMap {
		if hasGoCoverage {
			jobDetailMap := addProjAndRepoIfNeed(projNames[0], repoName)
			jobDetailMap[repoName] = []string{"test-coverage"}
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
	return regexp.MustCompile(`.+-[0-9\.]+$`).FindString(projName) != ""
}

// setOutput set the given file as the output target, then all the output will be written to this file
func setOutput(fileName string) {
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
	testgridConfigOutput := ""
	var generateProwConfig = flag.Bool("generate-prow-config", true, "Whether to generate the prow configuration file from the template")
	flag.StringVar(&prowConfigOutput, "prow-config-output", "", "The destination for the prow config output, default to be stdout")
	var generateTestgridConfig = flag.Bool("generate-testgrid-config", true, "Whether to generate the testgrid config from the template file")
	flag.StringVar(&testgridConfigOutput, "testgrid-config-output", "", "The destination for the testgrid config output, default to be stdout")

	var includeConfig = flag.Bool("include-config", true, "Whether to include general configuration (e.g., plank) in the generated config")
	flag.StringVar(&gcsBucket, "gcs-bucket", "knative-prow", "GCS bucket to upload the logs to")
	flag.StringVar(&logsDir, "logs-dir", "logs", "Path in the GCS bucket to upload logs of periodic and post-submit jobs")
	flag.StringVar(&presubmitLogsDir, "presubmit-logs-dir", "pr-logs", "Path in the GCS bucket to upload logs of pre-submit jobs")
	flag.StringVar(&testAccount, "test-account", "/etc/test-account/service-account.json", "Path to the service account JSON for test jobs")
	flag.StringVar(&nightlyAccount, "nightly-account", "/etc/nightly-account/service-account.json", "Path to the service account JSON for nightly release jobs")
	flag.StringVar(&releaseAccount, "release-account", "/etc/release-account/service-account.json", "Path to the service account JSON for release jobs")
	flag.StringVar(&flakesreporterDockerImage, "flaky-test-reporter-docker", "gcr.io/knative-tests/test-infra/flaky-test-reporter:latest", "Docker image for flaky test reporting tool")
	flag.StringVar(&prowversionbumperDockerImage, "prow-auto-bumper", "gcr.io/knative-tests/test-infra/prow-auto-bumper:latest", "Docker image for Prow version bumping tool")
	flag.StringVar(&coverageDockerImage, "coverage-docker", "gcr.io/knative-tests/test-infra/coverage:latest", "Docker image for coverage tool")
	flag.StringVar(&clearalertsDockerImage, "clear-alerts", "gcr.io/knative-tests/test-infra/monitoring/clear-alerts:latest", "Docker image for clearing alerts in test-infra monitoring")
	flag.StringVar(&prowTestsDockerImage, "prow-tests-docker", "gcr.io/knative-tests/test-infra/prow-tests:stable", "prow-tests docker image")
	flag.StringVar(&presubmitScript, "presubmit-script", "./test/presubmit-tests.sh", "Executable for running presubmit tests")
	flag.StringVar(&releaseScript, "release-script", "./hack/release.sh", "Executable for creating releases")
	flag.StringVar(&performanceScript, "performance-script", "./test/performance-tests.sh", "Executable for running performance tests")
	flag.StringVar(&webhookAPICoverageScript, "webhookAPICoverageScript", "./test/apicoverage.sh", "Executable for running webhook apicoverage tool")
	flag.StringVar(&cleanupScript, "cleanup-script", "./tools/cleanup/cleanup.sh", "Executable for running the cleanup tasks")
	flag.StringVar(&repositoryOverride, "repo-override", "", "Repository path (github.com/foo/bar[=branch]) to use instead for a job")
	flag.IntVar(&timeoutOverride, "timeout-override", 0, "Timeout (in minutes) to use instead for a job")
	flag.StringVar(&jobNameFilter, "job-filter", "", "Generate only this job, instead of all jobs")
	flag.StringVar(&preCommand, "pre-command", "", "Executable for running instead of the real command of a job")
	flag.Var(&extraEnvVars, "extra-env", "Extra environment variables (key=value) to add to a job")
	flag.Parse()
	if len(flag.Args()) != 1 {
		log.Fatal("Pass the config file as parameter")
	}
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

	// Generate Prow config.
	if *generateProwConfig {
		output = os.Stdout
		if prowConfigOutput != "" {
			setOutput(prowConfigOutput)
		}
		repositories = make([]repositoryData, 0)
		sectionMap = make(map[string]bool)
		if *includeConfig {
			executeTemplate("general config", readTemplate(generalProwConfig), newbaseProwJobTemplateData(""))
		}
		parseSection(config, "presubmits", generatePresubmit, nil)
		parseSection(config, "periodics", generatePeriodic, generateGoCoveragePeriodic)
		generateOtherJobConfigs("periodics", func(repo repositoryData) bool {
			return !repo.Processed && repo.EnableGoCoverage
		}, generateGoCoveragePeriodic)
		generateCleanupPeriodicJob()
		generateClearAlertsPeriodicJob()
		generateFlakytoolPeriodicJob()
		generateVersionBumpertoolPeriodicJob()
		generateBackupPeriodicJob()
		generateOtherJobConfigs("postsubmits", func(repo repositoryData) bool {
			return repo.EnableGoCoverage
		}, generateGoCoveragePostsubmit)
	}

	// config object is modified when we generate prow config, so we'll need to reload it here
	if err = yaml.Unmarshal(content, &config); err != nil {
		log.Fatalf("Cannot parse config %q: %v", name, err)
	}
	// Generate Testgrid config.
	if *generateTestgridConfig {
		output = os.Stdout
		if testgridConfigOutput != "" {
			setOutput(testgridConfigOutput)
		}

		if *includeConfig {
			executeTemplate("general config", readTemplate(generalTestgridConfig), newBaseTestgridTemplateData(""))
		}

		presubmitJobData := parseJob(config, "presubmits")
		goCoverageMap = parseGoCoverageMap(presubmitJobData)

		periodicJobData := parseJob(config, "periodics")
		collectMetaData(periodicJobData)

		generateTestGridSection("test_groups", generateTestGroup, false)
		generateTestGridSection("dashboards", generateDashboard, true)
		generateDashboardsForReleases()
		generateDashboardGroups()
	}
}
