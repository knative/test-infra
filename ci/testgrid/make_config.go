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

// The make_config tool generates a full Testgrid config for the Knative project,
// with input from a yaml file with key definitions.

package main

import (
	"bytes"
	"flag"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"strings"
	"time"

	cg "github.com/knative/test-infra/shared/configgenerator"
	"gopkg.in/yaml.v2"
)

var (
	// Values used in the jobs that can be changed through command-line flags.
	gcsBucket string
	logsDir   string

	goCoverageMap map[string]bool
	// save the repo names when parsing the config file, for the purpose of maintaining the output sequence
	repoNames []string
)

// baseTestgridTemplateData contains basic data about the testgrid config file.
type baseTestgridTemplateData struct {
	TestGroupName string
	Year          int
}

// testGroupTemplateData contains data about a test group
type testGroupTemplateData struct {
	Base      baseTestgridTemplateData
	GcsLogDir string
	Extras    map[string]string
}

// dashboardTabTemplateData contains data about a dashboard tab
type dashboardTabTemplateData struct {
	Base        baseTestgridTemplateData
	Name        string
	BaseOptions string
	Extras      map[string]string
}

const (

	// generalConfig contains config-wide definitions.
	generalConfig = `
# Copyright [[.Year]] The Knative Authors
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.
	
# Default testgroup and dashboardtab, please do not change them
default_test_group:
  days_of_results: 14            # Number of days of test results to gather and serve
  tests_name_policy: 2           # Replace the name of the test
  ignore_pending: false          # Show in-progress tests
  column_header:
  - configuration_value: Commit  # Shows the commit number on column header
  - configuration_value: infra-commit
  num_columns_recent: 10         # The number of columns to consider "recent" for a variety of purposes
  use_kubernetes_client: true    # ** This field is deprecated and should always be true **
  is_external: true              # ** This field is deprecated and should always be true **
  alert_stale_results_hours: 24  # Alert if tests haven't run for a day
  num_failures_to_alert: 3       # Consider a test failed if it has 3 or more consecutive failures
  num_passes_to_disable_alert: 1 # Consider a failing test passing if it has 1 or more consecutive passes

default_dashboard_tab:
  open_test_template:            # The URL template to visit after clicking on a cell
    url: https://gubernator.knative.dev/build/<gcs_prefix>/<changelist>
  file_bug_template:             # The URL template to visit when filing a bug
    url: https://github.com/knative/serving/issues/new
    options:
    - key: title
      value: 'Test "<test-name>" failed'
    - key: body
      value: <test-url>
  attach_bug_template:           # The URL template to visit when attaching a bug
    url:                         # Empty
    options:                     # Empty
  # Text to show in the about menu as a link to another view of the results
  results_text: See these results in Gubernator
  results_url_template:          # The URL template to visit after clicking
    url: https://gubernator.knative.dev/builds/<gcs_prefix>
  # URL for regression search links.
  code_search_path: github.com/knative/serving/search
  num_columns_recent: 10
  code_search_url_template:      # The URL template to visit when searching for changelists
    url: https://github.com/knative/serving/compare/<start-custom-0>...<end-custom-0>
  alert_options:
	alert_mail_to_addresses: 'knative-productivity-dev@googlegroups.com'
	`

	// testGroupTemplate is the template for the test group config
	testGroupTemplate = `
- name: [[.Base.TestGroupName]]
  gcs_prefix: [[.GcsLogDir]]
  [[indent_map 2 .Extras]]
	`

	// dashboardTabTemplate is the template for the dashboard tab config
	dashboardTabTemplate = `
	- name: [[.Name]]
      test_group_name: [[.Base.TestGroupName]]
      base_options: '[[.BaseOptions]]'
	  [[indent_map 2 .Extras]]
	`
)

// newTestgridTemplateData returns a testgridTemplateData type with its initial, default values.
func newBaseTestgridTemplateData(testGroupName string) baseTestgridTemplateData {
	var data baseTestgridTemplateData
	data.Year = time.Now().Year()
	data.TestGroupName = testGroupName
	return data
}

// executeTemplate outputs the given template with the given data.
func executeTemplate(name, templ string, data interface{}) {
	var res bytes.Buffer
	// comment out this part since we do not need it for now, but might need in the future
	funcMap := template.FuncMap{
		"indent_section":       cg.IndentSection,
		"indent_array_section": cg.IndentArraySection,
		"indent_array":         cg.IndentArray,
		"indent_map":           cg.IndentMap,
	}
	t := template.Must(template.New(name).Funcs(funcMap).Delims("[[", "]]").Parse(templ))
	if err := t.Execute(&res, data); err != nil {
		log.Fatalf("Error in template %s: %v", name, err)
	}
	for _, line := range strings.Split(res.String(), "\n") {
		cg.OutputConfig(line)
	}
}

// get the job data from the original yaml data, now the jobName can be "presubmits" or "periodic"
func parseJob(config yaml.MapSlice, jobName string) yaml.MapSlice {
	for _, section := range config {
		if section.Key == jobName {
			return cg.GetMapSlice(section.Value)
		}
	}

	log.Fatalf("The metadata misses %s configuration, cannot continue.", jobName)
	return nil
}

// construct a map, indicating which repo is enabled for go coverage check
func parseGoCoverageMap(presubmitJob yaml.MapSlice) map[string]bool {
	goCoverageMap := make(map[string]bool)
	for _, repo := range presubmitJob {
		repoName := strings.Split(cg.GetString(repo.Key), "/")[1]
		goCoverageMap[repoName] = false
		for _, jobConfig := range cg.GetInterfaceArray(repo.Value) {
			for _, item := range cg.GetMapSlice(jobConfig) {
				if item.Key == "go-coverage" {
					goCoverageMap[repoName] = cg.GetBool(item.Value)
					break
				}
			}
		}
	}

	return goCoverageMap
}

// collect the meta data from the original yaml data, which we can then use for building the test groups and dashboards config
func collectMetaData(periodicJob yaml.MapSlice, metaData *map[string]map[string][]string) {
	for _, repo := range periodicJob {
		rawName := cg.GetString(repo.Key)
		projName := strings.Split(rawName, "/")[0]
		repoName := strings.Split(rawName, "/")[1]
		jobDetailMap := addProjAndRepoIfNeed(metaData, projName, repoName)

		// parse job configs
		for _, conf := range cg.GetInterfaceArray(repo.Value) {
			jobDetailMap = (*metaData)[projName]
			jobConfig := cg.GetMapSlice(conf)
			enabled := false
			jobName := ""
			releaseVersion := ""
			for _, item := range jobConfig {
				switch item.Key {
				case "continuous", "dot-release", "auto-release", "performance", "latency", "api-coverage", "nightly":
					if cg.GetBool(item.Value) {
						enabled = true
						jobName = cg.GetString(item.Key)
					}
				case "branch-ci":
					enabled = cg.GetBool(item.Value)
					jobName = "continuous"
				case "release":
					releaseVersion = cg.GetString(item.Value)
				case "custom-job":
					enabled = true
					jobName = cg.GetString(item.Value)
				default:
					continue
				}
			}
			// add job types for the corresponding repos, if needed
			if enabled {
				if releaseVersion != "" {
					releaseProjName := fmt.Sprintf("%s-%s", projName, releaseVersion)
					jobDetailMap = addProjAndRepoIfNeed(metaData, releaseProjName, repoName)
				}
				newJobTypes := append(jobDetailMap[repoName], jobName)
				jobDetailMap[repoName] = newJobTypes
			}
		}
		addTestCoverageJobIfNeeded(&jobDetailMap, repoName)
	}
}

// add the project and repo if they are new in the metaData map, then return the jobDetailMap
func addProjAndRepoIfNeed(metaData *map[string]map[string][]string, projName string, repoName string) map[string][]string {
	// add project in the metaData
	if _, exists := (*metaData)[projName]; !exists {
		(*metaData)[projName] = make(map[string][]string)
	}

	// add repo in the project
	jobDetailMap := (*metaData)[projName]
	if _, exists := jobDetailMap[repoName]; !exists {
		if !cg.ContainsStr(repoNames, repoName) {
			repoNames = append(repoNames, repoName)
		}
		jobDetailMap[repoName] = make([]string, 0)
	}
	return jobDetailMap
}

// if the repo has go coverage check, add test-coverage job for this repo
func addTestCoverageJobIfNeeded(jobDetailMap *map[string][]string, repoName string) {
	if goCoverageMap[repoName] {
		newJobTypes := append((*jobDetailMap)[repoName], "test-coverage")
		(*jobDetailMap)[repoName] = newJobTypes
	}
}

func generateTestGroups(data yaml.MapSlice, goCoverageMap map[string]bool) {
	cg.OutputConfig("test_groups:")
	for _, repo := range data {
		repoName := cg.GetString(repo.Key)
		repoNameForConfig := strings.Replace(repoName, "/", "-", -1)
		for _, item := range cg.GetInterfaceArray(repo.Value) {
			jobConfig := cg.GetMapSlice(item)
			generateTestGroup(repoNameForConfig, jobConfig)
		}

		if goCoverageMap[repoName] {
			generateCoverageTestGroup(repoNameForConfig)
		}

		fmt.Println()
	}
}

func generateTestGroup(repoNameForConfig string, jobConfig yaml.MapSlice) {
	enabled := false
	testGroupName := ""
	gcsLogDir := ""
	extras := make(map[string]string)
	for _, item := range jobConfig {
		switch item.Key {
		case "continuous", "dot-release", "auto-release", "performance", "latency", "api-coverage", "nightly":
			enabled = cg.GetBool(item.Value)
			testGroupName = fmt.Sprintf("ci-%s-%s", repoNameForConfig, item.Key)
			gcsLogDir = fmt.Sprintf("%s/%s/ci-%s-%s", gcsBucket, logsDir, repoNameForConfig, item.Key)

			if item.Key == "nightly" {
				testGroupName += "-release"
				gcsLogDir += "-release"
			}

			// TODO: confirm if they are needed or not
			if item.Key == "latency" {
				extras["short_text_metric"] = "latency"
			}
			if item.Key == "api-coverage" {
				extras["short_text_metric"] = "api_coverage"
			}
			if item.Key == "performance" {
				extras["short_text_metric"] = "perf_latency"
			}
		case "branch-ci":
			enabled = cg.GetBool(item.Value)
		case "release":
			releaseVersion := cg.GetString(item.Value)
			testGroupName = fmt.Sprintf("ci-%s-%s-%s", repoNameForConfig, releaseVersion, "continuous")
			gcsLogDir = fmt.Sprintf("%s/%s/ci-%s-%s-%s", gcsBucket, logsDir, repoNameForConfig, releaseVersion, "continuous")
		case "custom-job":
			enabled = true
			customJobName := cg.GetString(item.Value)
			testGroupName = fmt.Sprintf("ci-%s-%s", repoNameForConfig, customJobName)
			gcsLogDir = fmt.Sprintf("%s/%s/ci-%s-%s", gcsBucket, logsDir, repoNameForConfig, customJobName)
			extras["alert_stale_results_hours"] = "168"
		default:
			continue
		}
	}

	if !enabled {
		return
	}
	executeTestGroupTemplate(testGroupName, gcsLogDir, extras)
}

func generateCoverageTestGroup(repoNameForConfig string) {
	testGroupName := fmt.Sprintf("pull-%s-%s", repoNameForConfig, "test-coverage")
	gcsLogDir := fmt.Sprintf("%s/%s/ci-%s-%s", gcsBucket, logsDir, repoNameForConfig, "go-coverage")
	extras := make(map[string]string)
	extras["short_text_metric"] = "coverage"
	executeTestGroupTemplate(testGroupName, gcsLogDir, extras)
}

func executeTestGroupTemplate(testGroupName string, gcsLogDir string, extras map[string]string) {
	var data testGroupTemplateData
	data.Base.TestGroupName = testGroupName
	data.GcsLogDir = gcsLogDir
	data.Extras = extras
	executeTemplate("test group", testGroupTemplate, data)
}

func generateDashboards(data yaml.MapSlice, goCoverageMap map[string]bool) {
	cg.OutputConfig("dashboards:")

}

func generateDashboardGroups() {

}

// main is the script entry point.
func main() {
	// Parse flags and sanity check them.
	var includeConfig = flag.Bool("include-config", false, "Whether to include general configuration (e.g., plank) in the generated config")
	flag.StringVar(&gcsBucket, "gcs-bucket", "knative-prow", "GCS bucket to upload the logs to")
	flag.StringVar(&logsDir, "logs-dir", "logs", "Path in the GCS bucket to upload logs of periodic and post-submit jobs")

	flag.Parse()
	if len(flag.Args()) != 1 {
		log.Fatal("Pass the config file as parameter")
	}
	// We use MapSlice instead of maps to keep key order and create predictable output.
	config := yaml.MapSlice{}
	// repositories = make([]repositoryData, 0)
	// sectionMap = make(map[string]bool)

	// Read input config.
	name := flag.Arg(0)
	content, err := ioutil.ReadFile(name)
	if err != nil {
		log.Fatalf("Cannot read file %q: %v", name, err)
	}
	if err = yaml.Unmarshal(content, &config); err != nil {
		log.Fatalf("Cannot parse config %q: %v", name, err)
	}

	// Generate Testgrid config.
	if *includeConfig {
		// executeTemplate("general config", generalConfig, newBaseTestgridTemplateData(""))
	}

	presubmitJobData := parseJob(config, "presubmits")
	goCoverageMap = parseGoCoverageMap(presubmitJobData)
	// key is the main project version, value is another map containing job details
	//     for the job detail map, key is the repo name, value is the list of job types, like continuous, latency, nightly, and etc.
	metaDataMap := make(map[string]map[string][]string)
	periodicJobData := parseJob(config, "periodics")
	collectMetaData(periodicJobData, &metaDataMap)
	fmt.Println(metaDataMap)
	fmt.Println(repoNames)

	// generateTestGroups(periodicJobData, goCoverageMap)
	// generateDashboards(periodicJobData, goCoverageMap)
	// generateDashboardGroups()
}
