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
	"os"
	"path"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"text/template"
	"time"

	"gopkg.in/yaml.v3"
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
	Go113                  bool
	Go114                  bool
	Go112Branches          []string
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
type sectionGenerator func(string, string, *prowJob)

// stringArrayFlag is the content of a multi-value flag.
type stringArrayFlag []string

var (
	// Values used in the jobs that can be changed through command-line flags.
	output                       *os.File
	prowHost                     string
	testGridHost                 string
	gubernatorHost               string
	gcsBucket                    string
	testGridGcsBucket            string
	logsDir                      string
	presubmitLogsDir             string
	testAccount                  string
	nightlyAccount               string
	releaseAccount               string
	flakesreporterDockerImage    string
	prowversionbumperDockerImage string
	prowconfigupdaterDockerImage string
	githubCommenterDockerImage   string
	coverageDockerImage          string
	prowTestsDockerImage         string
	backupsDockerImage           string
	presubmitScript              string
	releaseScript                string
	webhookAPICoverageScript     string

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
	repositories []repositoryData

	// Map which sections of the config.yaml were written to stdout.
	sectionMap map[string]bool

	// To be used to flag that outputConfig() emitted data.
	emittedOutput bool

	projNameRegex = regexp.MustCompile(`.+-[0-9\.]+$`)
)

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
	data.ReleaseGcs = strings.Replace(repo, data.OrgName+"/", "knative-releases/", 1)
	data.AlwaysRun = true
	data.Image = prowTestsDockerImage
	data.ServiceAccount = testAccount
	data.Command = ""
	data.Args = make([]string, 0)
	data.Volumes = make([]string, 0)
	data.VolumeMounts = make([]string, 0)
	data.Resources = make([]string, 0)
	data.Env = make([]string, 0)
	data.ExtraRefs = []string{"- org: " + data.OrgName, "  repo: " + data.RepoName}
	data.Labels = make([]string, 0)
	data.Optional = ""
	return data
}

// Config parsers.

// parseBasicJobConfigOverrides updates the given baseProwJobTemplateData with any base option present in the given config.
func parseBasicJobConfigOverrides(data *baseProwJobTemplateData, pj *prowJob) {
	(*data).ExtraRefs = append((*data).ExtraRefs, "  base_ref: "+(*data).RepoBranch)
	var needDotdev, needGo113, needGo114 bool

	(*data).SkipBranches = pj.SkipBranches
	(*data).Branches = pj.Branches
	(*data).Args = pj.Args
	(*data).Timeout = pj.Timeout
	(*data).Command = pj.Command
	(*data).NeedsMonitor = pj.NeedsMonitor
	if pj.NeedsDind {
		setupDockerInDockerForJob(data)
	}
	(*data).AlwaysRun = pj.AlwaysRun
	// This is a good candidate of using default
	if pj.DotDev {
		needDotdev = true
		for i, repo := range repositories {
			if path.Base(repo.Name) == (*data).RepoName {
				repositories[i].DotDev = true
			}
		}
	}
	if pj.Go113 {
		needGo113 = true
		for i, repo := range repositories {
			if path.Base(repo.Name) == (*data).RepoName {
				repositories[i].Go113 = true
			}
		}
	}
	if pj.Go114 {
		needGo114 = true
		for i, repo := range repositories {
			if path.Base(repo.Name) == (*data).RepoName {
				repositories[i].Go114 = true
				data.RepoNameForJob = (*data).RepoName + getGo114ID()
			}
		}
	}
	if pj.Performance {
		for i, repo := range repositories {
			if path.Base(repo.Name) == (*data).RepoName {
				repositories[i].EnablePerformanceTests = true
			}
		}
	}
	if len(pj.Go112Branches) > 0 {
		for i, repo := range repositories {
			if path.Base(repo.Name) == (*data).RepoName {
				repositories[i].Go112Branches = pj.Go112Branches
			}
		}
	}

	addExtraEnvVarsToJob(pj.EnvVars, data)

	if pj.Optional {
		(*data).Optional = "optional: true"
	}
	setReourcesReqForJob(pj.Resources, data)

	if needDotdev {
		(*data).PathAlias = "path_alias: knative.dev/" + (*data).RepoName
		(*data).ExtraRefs = append((*data).ExtraRefs, "  "+(*data).PathAlias)
	}
	if needGo113 {
		(*data).Image = getGo113ImageName((*data).Image)
	}
	if needGo114 {
		(*data).Image = getGo114ImageName((*data).Image)
	}
	// Override any values if provided by command-line flags.
	if timeoutOverride > 0 {
		(*data).Timeout = timeoutOverride
	}
}

// getProwConfigData gets some basic, general data for the Prow config.
func getProwConfigData(ac allConfig) prowConfigTemplateData {
	var data prowConfigTemplateData
	data.Year = time.Now().Year()
	data.ProwHost = prowHost
	data.TestGridHost = testGridHost
	data.GubernatorHost = gubernatorHost
	data.GcsBucket = gcsBucket
	data.TestGridGcsBucket = testGridGcsBucket
	data.PresubmitLogsDir = presubmitLogsDir
	data.LogsDir = logsDir
	data.TideRepos = make([]string, 0)
	// Repos enabled for tide are all those that have presubmit jobs.
	for _, rc := range ac.Periodics {
		data.TideRepos = appendIfUnique(data.TideRepos, rc.Repo)
		if strings.HasSuffix(rc.Repo, "test-infra") {
			data.TestInfraRepo = rc.Repo
		}
	}

	// Sort repos to make output stable.
	sort.Strings(data.TideRepos)
	return data
}

// parseSection generate the configs from a given section of the input yaml file.
func parseSection(rcs []repoConfig, title string, generate sectionGenerator, finalize sectionGenerator) {
	for _, rc := range rcs {
		for _, pj := range rc.Jobs {
			pj := pj
			generate(title, rc.Repo, &pj)
		}
		if finalize != nil {
			finalize(title, rc.Repo, nil)
		}
	}
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

// executeJobTemplateWrapper takes in consideration of repo settings, decides how many varianats of the
// same job needs to be generated and generates them.
func executeJobTemplateWrapper(repoName string, data interface{}, generateOneJob func(data interface{})) {
	var sbs []specialBranchLogic

	switch data.(type) {
	case *postsubmitJobTemplateData:
		if strings.HasSuffix(data.(*postsubmitJobTemplateData).PostsubmitJobName, "go-coverage") {
			generateOneJob(data)
			return
		}
	}

	var go112Branches []string
	// Find out if Go112Branches is set in repo settings
	for _, repo := range repositories {
		if repo.Name == repoName {
			if len(repo.Go112Branches) > 0 {
				go112Branches = repo.Go112Branches
			}
		}
	}
	if len(go112Branches) > 0 {
		sbs = append(sbs, specialBranchLogic{
			branches: go112Branches,
			opsNew: func(base *baseProwJobTemplateData) {
				base.Image = getGo113ImageName(base.Image)
			},
			restore: func(base *baseProwJobTemplateData) {
				base.Image = getGo112ImageName(base.Image)
			},
		})
	} else {
		base := getBase(data)
		base.Image = getGo113ImageName(base.Image)
	}

	if len(sbs) == 0 { // Generate single job if there is no special branch logic
		generateOneJob(data)
	} else {
		recursiveSBL(repoName, data, generateOneJob, sbs, 0)
	}
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

// main is the script entry point.
func main() {
	// Parse flags and sanity check them.
	prowConfigOutput := ""
	prowJobsConfigOutput := ""
	testgridConfigOutput := ""
	var env = flag.String("env", "prow", "The name of the environment, can be prow or prow-staging")
	var generateProwConfig = flag.Bool("generate-prow-config", true, "Whether to generate the prow configuration file from the template")
	var generatePluginsConfig = flag.Bool("generate-plugins-config", true, "Whether to generate the plugins configuration file from the template")
	// var generateTestgridConfig = flag.Bool("generate-testgrid-config", true, "Whether to generate the testgrid config from the template file")
	var generateMaintenanceJobs = flag.Bool("generate-maintenance-jobs", true, "Whether to generate the maintenance periodic jobs (e.g. backup)")
	var includeConfig = flag.Bool("include-config", true, "Whether to include general configuration (e.g., plank) in the generated config")
	var dockerImagesBase = flag.String("image-docker", "gcr.io/knative-tests/test-infra", "Default registry for the docker images used by the jobs")
	flag.StringVar(&prowConfigOutput, "prow-config-output", "", "The destination for the prow config output, default to be stdout")
	flag.StringVar(&prowJobsConfigOutput, "prow-jobs-config-output", "", "The destination for the prow jobs config output, default to be stdout")
	var pluginsConfigOutput = flag.String("plugins-config-output", "", "The destination for the plugins config output, default to be stdout")
	flag.StringVar(&testgridConfigOutput, "testgrid-config-output", "", "The destination for the testgrid config output, default to be stdout")
	flag.StringVar(&prowHost, "prow-host", "https://prow.knative.dev", "Prow host, including HTTP protocol")
	flag.StringVar(&testGridHost, "testgrid-host", "https://testgrid.knative.dev", "TestGrid host, including HTTP protocol")
	flag.StringVar(&gubernatorHost, "gubernator-host", "https://gubernator.knative.dev", "Gubernator host, including HTTP protocol")
	flag.StringVar(&gcsBucket, "gcs-bucket", "knative-prow", "GCS bucket to upload the logs to")
	flag.StringVar(&testGridGcsBucket, "testgrid-gcs-bucket", "knative-testgrid", "TestGrid GCS bucket")
	flag.StringVar(&logsDir, "logs-dir", "logs", "Path in the GCS bucket to upload logs of periodic and post-submit jobs")
	flag.StringVar(&presubmitLogsDir, "presubmit-logs-dir", "pr-logs", "Path in the GCS bucket to upload logs of pre-submit jobs")
	flag.StringVar(&testAccount, "test-account", "/etc/test-account/service-account.json", "Path to the service account JSON for test jobs")
	flag.StringVar(&nightlyAccount, "nightly-account", "/etc/nightly-account/service-account.json", "Path to the service account JSON for nightly release jobs")
	flag.StringVar(&releaseAccount, "release-account", "/etc/release-account/service-account.json", "Path to the service account JSON for release jobs")
	var flakesreporterDockerImageName = flag.String("flaky-test-reporter-docker", "flaky-test-reporter:latest", "Docker image for flaky test reporting tool")
	var prowversionbumperDockerImageName = flag.String("prow-auto-bumper", "prow-auto-bumper:latest", "Docker image for Prow version bumping tool")
	var prowconfigupdaterDockerImageName = flag.String("prow-config-updater", "prow-config-updater:latest", "Docker image for Prow config updater tool")
	var coverageDockerImageName = flag.String("coverage-docker", "coverage-go112:latest", "Docker image for coverage tool")
	var prowTestsDockerImageName = flag.String("prow-tests-docker", "prow-tests-go112:stable", "prow-tests docker image")
	var backupsDockerImageName = flag.String("backups-docker", "backups:latest", "Docker image for the backups job")
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

	flakesreporterDockerImage = path.Join(*dockerImagesBase, *flakesreporterDockerImageName)
	prowversionbumperDockerImage = path.Join(*dockerImagesBase, *prowversionbumperDockerImageName)
	prowconfigupdaterDockerImage = path.Join(*dockerImagesBase, *prowconfigupdaterDockerImageName)
	coverageDockerImage = path.Join(*dockerImagesBase, *coverageDockerImageName)
	prowTestsDockerImage = path.Join(*dockerImagesBase, *prowTestsDockerImageName)
	backupsDockerImage = path.Join(*dockerImagesBase, *backupsDockerImageName)

	// Read input config.
	name := flag.Arg(0)
	content, err := ioutil.ReadFile(name)
	if err != nil {
		log.Fatalf("Cannot read file %q: %v", name, err)
	}

	var ac allConfig
	if err = yaml.Unmarshal(content, &ac); err != nil {
		log.Fatalf("Cannot parse config %q: %v", name, err)
	}

	// Inject top level configs into presubmit, postsubmit, and periodic configs
	ac.init()

	prowConfigData := getProwConfigData(ac)

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
		parseSection(ac.Presubmits, "presubmits", generatePresubmit, nil)
		parseSection(ac.Postsubmits, "postsubmits", generatePostsubmit, nil /*, generateGoCoveragePeriodic*/)
		parseSection(ac.Periodics, "periodics", generatePeriodic, nil /*, generateGoCoveragePeriodic*/)

		// for _, repo := range repositories { // Keep order for predictable output.
		// 	if !repo.Processed && repo.EnableGoCoverage {
		// 		generateGoCoveragePeriodic("periodics", repo.Name, nil)
		// 	}
		// }
		// Intentionally removed as Chi is working on disabling generating these
		// if *generateMaintenanceJobs {
		// 	generateCleanupPeriodicJob()
		// 	generateFlakytoolPeriodicJob()
		// 	generateVersionBumpertoolPeriodicJob()
		// 	generateBackupPeriodicJob()
		// 	generateIssueTrackerPeriodicJobs()
		// 	generatePerfClusterUpdatePeriodicJobs()
		// }

		// TODO(chizhg): These should be defined as postsubmit in config_knative.yaml
		for _, repo := range repositories {
			if repo.EnableGoCoverage {
				// generateGoCoveragePostsubmit("postsubmits", repo.Name, nil)
				if repo.Name == "knative/test-infra" && *generateMaintenanceJobs {
					generateConfigUpdaterToolPostsubmitJob()
				}
			}
			if repo.EnablePerformanceTests {
				generatePerfClusterPostsubmitJob(repo)
			}
		}
	}

	// Worry about testgrid later
	// // config object is modified when we generate prow config, so we'll need to reload it here
	// if err = yaml.Unmarshal(content, &config); err != nil {
	// 	log.Fatalf("Cannot parse config %q: %v", name, err)
	// }
	// // Generate Testgrid config.
	// if *generateTestgridConfig {
	// 	setOutput(testgridConfigOutput)

	// 	if *includeConfig {
	// 		executeTemplate("general header", readTemplate(commonHeaderConfig), newBaseTestgridTemplateData(""))
	// 		executeTemplate("general config", readTemplate(generalTestgridConfig), newBaseTestgridTemplateData(""))
	// 	}

	// 	presubmitJobData := parseJob(config, "presubmits")
	// 	goCoverageMap = parseGoCoverageMap(presubmitJobData)

	// 	periodicJobData := parseJob(config, "periodics")
	// 	collectMetaData(periodicJobData)

	// 	generateTestGridSection("test_groups", generateTestGroup, false)
	// 	generateTestGridSection("dashboards", generateDashboard, true)
	// 	generateDashboardsForReleases()
	// 	generateDashboardGroups()
	// }
}
