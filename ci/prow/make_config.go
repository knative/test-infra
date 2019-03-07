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
	"strings"
	"text/template"
	"time"

	"gopkg.in/yaml.v2"
)

const (
	// Manifests generated by ko are indented by 2 spaces.
	baseIndent = "  "
	// Cron strings for key jobs
	goCoveragePeriodicJobCron = "0 1 * * *"  // Run at 01:00 every day
	cleanupPeriodicJobCron    = "0 19 * * 1" // Run at 11:00PST/12:00PST every Monday (19:00 UTC)
	flakytoolPeriodicJobCron  = "0 12 * * *" // Run at 4:00PST/5:00PST every day (12:00 UTC)
	backupPeriodicJobCron     = "15 9 * * *" // Run at 02:15PST every day (09:15 UTC)
)

// repositoryData contains basic data about each Knative repository.
type repositoryData struct {
	Name                string
	EnableGoCoverage    bool
	GoCoverageThreshold int
}

// baseProwJobTemplateData contains basic data about a Prow job.
type baseProwJobTemplateData struct {
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
}

// presubmitJobTemplateData contains data about a presubmit Prow job.
type presubmitJobTemplateData struct {
	Base                 baseProwJobTemplateData
	PresubmitJobName     string
	PresubmitPullJobName string
	PresubmitPostJobName string
	PresubmitCommand     []string
}

// periodicJobTemplateData contains data about a periodic Prow job.
type periodicJobTemplateData struct {
	Base            baseProwJobTemplateData
	PeriodicJobName string
	CronString      string
	PeriodicCommand []string
}

// postsubmitJobTemplateData contains data about a postsubmit Prow job.
type postsubmitJobTemplateData struct {
	Base              baseProwJobTemplateData
	PostsubmitJobName string
}

// sectionGenerator is a function that generates Prow job configs given a slice of a yaml file with configs.
type sectionGenerator func(string, string, yaml.MapSlice)

// stringArrayFlag is the content of a multi-value flag.
type stringArrayFlag []string

var (
	// Array constants used throughout the jobs.
	allPresubmitTests = []string{"--all-tests", "--emit-metrics"}
	releaseNightly    = []string{"--publish", "--tag-release"}
	releaseLocal      = []string{"--nopublish", "--notag-release"}

	// Values used in the jobs that can be changed through command-line flags.
	gcsBucket                string
	logsDir                  string
	presubmitLogsDir         string
	testAccount              string
	nightlyAccount           string
	releaseAccount           string
	flakytoolDockerImage     string
	coverageDockerImage      string
	prowTestsDockerImage     string
	presubmitScript          string
	releaseScript            string
	performanceScript        string
	webhookAPICoverageScript string
	cleanupScript            string

	// Overrides and behavior changes through command-line flags.
	repositoryOverride string
	jobNameFilter      string
	preCommand         string
	extraEnvVars       stringArrayFlag

	// List of Knative repositories.
	repositories []repositoryData

	// Map which sections of the config.yaml were written to stdout.
	sectionMap map[string]bool
)

// Templates for config generation.
// TODO(adrcunha): eliminate redundant templates (e.g., latency job) by factoring them as standard jobs (e.g., periodic job).

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

# ############################################################
# ####                                                    ####
# #### THIS FILE IS AUTOMATICALLY GENERATED. DO NOT EDIT. ####
# ####     USE "make config" TO REGENERATE THIS FILE.     ####
# ####                                                    ####
# ############################################################

plank:
  job_url_template: 'https://gubernator.knative.dev/build/[[.GcsBucket]]/{{if or (eq .Spec.Type "presubmit") (eq .Spec.Type "batch")}}[[.PresubmitLogsDir]]/pull{{with .Spec.Refs}}/{{.Org}}_{{.Repo}}{{end}}{{else}}[[.LogsDir]]{{end}}{{if eq .Spec.Type "presubmit"}}/{{with index .Spec.Refs.Pulls 0}}{{.Number}}{{end}}{{else if eq .Spec.Type "batch"}}/batch{{end}}/{{.Spec.Job}}/{{.Status.BuildID}}/'
  report_template: '[Full PR test history](https://gubernator.knative.dev/pr/{{.Spec.Refs.Org}}_{{.Spec.Refs.Repo}}/{{with index .Spec.Refs.Pulls 0}}{{.Number}}{{end}}). [Your PR dashboard](https://gubernator.knative.dev/pr/{{with index .Spec.Refs.Pulls 0}}{{.Author}}{{end}}).'
  pod_pending_timeout: 60m
  default_decoration_config:
    timeout: 7200000000000 # 2h
    grace_period: 15000000000 # 15s
    utility_images:
      clonerefs: "gcr.io/k8s-prow/clonerefs@sha256:b62ba1f379ac19c5ec9ee7bcab14d3f0b3c31cea9cdd4bc491e98e2c5f346c07"
      initupload: "gcr.io/k8s-prow/initupload@sha256:58f89f2aae68f7dc46aaf05c7e8204c4f26b53ec9ce30353d1c27ce44a60d121"
      entrypoint: "gcr.io/k8s-prow/entrypoint:v20180512-0255926d1"
      sidecar: "gcr.io/k8s-prow/sidecar@sha256:8807b2565f4d2699920542fcf890878824b1ede4198d7ff46bca53feb064ed44"
    gcs_configuration:
      bucket: "[[.GcsBucket]]"
      path_strategy: "explicit"
    gcs_credentials_secret: "test-account"

prowjob_namespace: default
pod_namespace: test-pods
log_level: info

branch-protection:
  orgs:
    knative:
      # Protect all branches in knative
      # This means all prow jobs with "always_run" set are required
      # to pass before tide can merge the PR.
      # Currently this is manually enabled by the knative org admins,
      # but it's stated here for documentation and reference purposes.
      protect: true
      # Admins can overrule checks
      enforce_admins: false

tide:
  queries:
  - repos:
    - knative/build
    - knative/build-pipeline
    - knative/build-templates
    - knative/observability
    - knative/serving
    - knative/eventing
    - knative/eventing-sources
    - knative/docs
    - knative/test-infra
    - knative/pkg
    - knative/caching
    - knative/website
    labels:
    - lgtm
    - approved
    missingLabels:
    - do-not-merge/hold
    - do-not-merge/work-in-progress
  merge_method:
    knative: squash
    knative/build-pipeline: rebase
  target_url: https://prow.knative.dev/tide
`

	// presubmitJob is the template for presubmit jobs.
	presubmitJob = `
  - name: [[.PresubmitPullJobName]]
    agent: kubernetes
    context: [[.PresubmitPullJobName]]
    always_run: [[.Base.AlwaysRun]]
    rerun_command: "/test [[.PresubmitPullJobName]]"
    trigger: "(?m)^/test (all|[[.PresubmitPullJobName]]),?(\\s+|$)"
    [[indent_array_section 4 "skip_branches" .Base.SkipBranches]]
    spec:
      containers:
      - image: [[.Base.Image]]
        imagePullPolicy: Always
        args:
        - "--scenario=kubernetes_execute_bazel"
        - "--clean"
        - "--job=$(JOB_NAME)"
        - "--repo=github.com/$(REPO_OWNER)/$(REPO_NAME)=$(PULL_REFS)"
        - "--root=/go/src"
        - "--service-account=[[.Base.ServiceAccount]]"
        - "--upload=[[.Base.GcsPresubmitLogDir]]"
        - "--" # end bootstrap args, scenario args below
        - "--" # end kubernetes_execute_bazel flags (consider following flags as text)
        [[indent_array 8 .PresubmitCommand]]
        [[indent_section 10 "securityContext" .Base.SecurityContext]]
        [[indent_section 8 "volumeMounts" .Base.VolumeMounts]]
        [[indent_section 8 "env" .Base.Env]]
      [[indent_section 6 "volumes" .Base.Volumes]]
`

	// presubmitGoCoverageJob is the template for go coverage presubmit jobs.
	presubmitGoCoverageJob = `
  - name: [[.PresubmitPullJobName]]
    agent: kubernetes
    context: [[.PresubmitPullJobName]]
    always_run: [[.Base.AlwaysRun]]
    rerun_command: "/test [[.PresubmitPullJobName]]"
    trigger: "(?m)^/test (all|[[.PresubmitPullJobName]]),?(\\s+|$)"
    optional: true
    decorate: true
    clone_uri: [[.Base.CloneURI]]
    spec:
      containers:
      - image: [[.Base.Image]]
        imagePullPolicy: Always
        command:
        - "/coverage"
        args:
        - "--postsubmit-gcs-bucket=[[.Base.GcsBucket]]"
        - "--postsubmit-job-name=[[.PresubmitPostJobName]]"
        - "--artifacts=$(ARTIFACTS)"
        - "--profile-name=coverage_profile.txt"
        - "--cov-target=."
        - "--cov-threshold-percentage=[[.Base.GoCoverageThreshold]]"
        - "--github-token=/etc/covbot-token/token"
        [[indent_section 8 "volumeMounts" .Base.VolumeMounts]]
        [[indent_section 8 "env" .Base.Env]]
      [[indent_section 6 "volumes" .Base.Volumes]]
`

	// periodicTestJob is the template for periodic test/release jobs.
	periodicTestJob = `
- cron: "[[.CronString]]"
  name: [[.PeriodicJobName]]
  agent: kubernetes
  spec:
    containers:
    - image: [[.Base.Image]]
      imagePullPolicy: Always
      args:
      - "--scenario=kubernetes_execute_bazel"
      - "--clean"
      - "--job=$(JOB_NAME)"
      - "--repo=[[repo .Base]]"
      - "--root=/go/src"
      - "--service-account=[[.Base.ServiceAccount]]"
      - "--upload=[[.Base.GcsLogDir]]"
      - "--timeout=[[.Base.Timeout]]" # Avoid overrun
      - "--" # end bootstrap args, scenario args below
      - "--" # end kubernetes_execute_bazel flags (consider following flags as text)
      [[indent_array 6 .PeriodicCommand]]
      [[indent_section 8 "securityContext" .Base.SecurityContext]]
      [[indent_section 6 "volumeMounts" .Base.VolumeMounts]]
      [[indent_section 6 "env" .Base.Env]]
    [[indent_section 4 "volumes" .Base.Volumes]]
`

	// periodicCustomJob is the template for periodic custom jobs.
	periodicCustomJob = `
- cron: "[[.CronString]]"
  name: [[.PeriodicJobName]]
  agent: kubernetes
  decorate: true
  [[indent_section 4 "decoration_config" .Base.DecorationConfig]]
  [[indent_section 2 "extra_refs" .Base.ExtraRefs]]
  spec:
    containers:
    - image: [[.Base.Image]]
      imagePullPolicy: Always
      command:
      - "[[.Base.Command]]"
      [[indent_array_section 6 "args" .Base.Args]]
      [[indent_section 6 "volumeMounts" .Base.VolumeMounts]]
      [[indent_section 6 "env" .Base.Env]]
    [[indent_section 4 "volumes" .Base.Volumes]]
`

	// goCoveragePostsubmitJob is the template for the go postsubmit coverage job.
	goCoveragePostsubmitJob = `
  - name: [[.PostsubmitJobName]]
    branches:
    - master
    agent: kubernetes
    decorate: true
    clone_uri: [[.Base.CloneURI]]
    spec:
      containers:
      - image: [[.Base.Image]]
        imagePullPolicy: Always
        command:
        - "/coverage"
        args:
        - "--artifacts=$(ARTIFACTS)"
        - "--profile-name=coverage_profile.txt"
        - "--cov-target=."
        - "--cov-threshold-percentage=0"
`
)

// Yaml parsing helpers.

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
	data.RepoName = strings.Replace(repo, "knative/", "", 1)
	data.RepoNameForJob = strings.Replace(repo, "/", "-", -1)
	data.GcsBucket = gcsBucket
	data.RepoURI = "github.com/" + repo
	data.CloneURI = fmt.Sprintf("\"https://%s.git\"", data.RepoURI)
	data.GcsLogDir = fmt.Sprintf("gs://%s/%s", gcsBucket, logsDir)
	data.GcsPresubmitLogDir = fmt.Sprintf("gs://%s/%s", gcsBucket, presubmitLogsDir)
	data.Year = time.Now().Year()
	data.PresubmitLogsDir = presubmitLogsDir
	data.LogsDir = logsDir
	data.ReleaseGcs = strings.Replace(repo, "knative/", "knative-releases/", 1)
	data.AlwaysRun = true
	data.Image = prowTestsDockerImage
	data.ServiceAccount = testAccount
	data.Command = ""
	data.Args = make([]string, 0)
	data.Volumes = make([]string, 0)
	data.VolumeMounts = make([]string, 0)
	data.Env = make([]string, 0)
	data.ExtraRefs = []string{"- org: knative", "  repo: " + data.RepoName, "  base_ref: master", "  clone_uri: " + data.CloneURI}
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
	(*data).Env = append((*data).Env, []string{"- name: " + key, "  value: " + value}...)
}

// addVolumeToJob adds the given mount path as volume for the job.
func addVolumeToJob(data *baseProwJobTemplateData, mountPath, name string, isSecret bool) {
	(*data).VolumeMounts = append((*data).VolumeMounts, []string{"- name: " + name, "  mountPath: " + mountPath}...)
	if isSecret {
		(*data).VolumeMounts = append((*data).VolumeMounts, "  readOnly: true")
	}
	s := []string{"- name: " + name}
	if isSecret {
		s = append(s, []string{"  secret:", "    secretName: " + name}...)
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
	addVolumeToJob(data, "/etc/"+name, name, true)
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
	addVolumeToJob(data, "/docker-graph", "docker-graph", false)
	addEnvToJob(data, "DOCKER_IN_DOCKER_ENABLED", "\"true\"")
	(*data).SecurityContext = []string{"privileged: true"}
}

// Config parsers.

// parseBasicJobConfigOverrides updates the given baseProwJobTemplateData with any base option present in the given config.
func parseBasicJobConfigOverrides(data *baseProwJobTemplateData, config yaml.MapSlice) {
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
		case "needs-dind":
			if getBool(item.Value) {
				setupDockerInDockerForJob(data)
			}
		case "always_run":
			(*data).AlwaysRun = getBool(item.Value)
		case nil: // already processed
			continue
		default:
			log.Fatalf("Unknown entry %q for job", item.Key)
			continue
		}
		// Knock-out the item, signalling it was already parsed.
		config[i] = yaml.MapItem{}
	}
}

// generatePresubmit generates all presubmit job configs for the given repo and configuration.
func generatePresubmit(title string, repoName string, presubmitConfig yaml.MapSlice) {
	var data presubmitJobTemplateData
	data.Base = newbaseProwJobTemplateData(repoName)
	data.Base.Command = presubmitScript
	data.Base.GoCoverageThreshold = 80
	jobTemplate := presubmitJob
	repoData := repositoryData{Name: repoName, EnableGoCoverage: false, GoCoverageThreshold: data.Base.GoCoverageThreshold}
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
		case "go-coverage":
			if !getBool(item.Value) {
				return
			}
			jobTemplate = presubmitGoCoverageJob
			data.PresubmitJobName = data.Base.RepoNameForJob + "-go-coverage"
			data.Base.Image = coverageDockerImage
			data.Base.ServiceAccount = ""
			repoData.EnableGoCoverage = true
			addVolumeToJob(&data.Base, "/etc/covbot-token", "covbot-token", true)
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
	parseBasicJobConfigOverrides(&data.Base, presubmitConfig)
	repositories = append(repositories, repoData)
	data.PresubmitCommand = createCommand(data.Base)
	data.PresubmitPullJobName = "pull-" + data.PresubmitJobName
	data.PresubmitPostJobName = "post-" + data.PresubmitJobName
	if data.Base.ServiceAccount != "" {
		addEnvToJob(&data.Base, "GOOGLE_APPLICATION_CREDENTIALS", data.Base.ServiceAccount)
		addEnvToJob(&data.Base, "E2E_CLUSTER_REGION", "us-west1")
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
		template := strings.Replace(presubmitGoCoverageJob, "(all|", "(", 1)
		executeJobTemplate("presubmit", template, title, repoName, data.PresubmitPullJobName, true, data)
	}
}

// generatePeriodic generates all periodic job configs for the given repo and configuration.
func generatePeriodic(title string, repoName string, periodicConfig yaml.MapSlice) {
	var data periodicJobTemplateData
	data.Base = newbaseProwJobTemplateData(repoName)
	jobNameSuffix := ""
	jobTemplate := periodicTestJob
	jobType := ""
	for i, item := range periodicConfig {
		switch item.Key {
		case "continuous":
			if !getBool(item.Value) {
				return
			}
			jobType = getString(item.Key)
			jobNameSuffix = "continuous"
			// Use default command and arguments if none given.
			if data.Base.Command == "" {
				data.Base.Command = presubmitScript
			}
			if len(data.Base.Args) == 0 {
				data.Base.Args = allPresubmitTests
			}
		case "nightly":
			if !getBool(item.Value) {
				return
			}
			jobType = getString(item.Key)
			jobNameSuffix = "nightly-release"
			data.Base.ServiceAccount = nightlyAccount
			data.Base.Command = releaseScript
			data.Base.Args = releaseNightly
			data.Base.Timeout = 90
		case "branch-ci":
			if !getBool(item.Value) {
				return
			}
			jobType = getString(item.Key)
			jobNameSuffix = "continuous"
			data.Base.Command = releaseScript
			data.Base.Args = releaseLocal
			setupDockerInDockerForJob(&data.Base)
			data.Base.Timeout = 90
		case "dot-release", "auto-release":
			if !getBool(item.Value) {
				return
			}
			jobType = getString(item.Key)
			jobNameSuffix = getString(item.Key)
			data.Base.ServiceAccount = releaseAccount
			data.Base.Command = releaseScript
			data.Base.Args = []string{
				"--" + jobNameSuffix,
				"--release-gcs " + data.Base.ReleaseGcs,
				"--release-gcr gcr.io/knative-releases",
				"--github-token /etc/hub-token/token"}
			addVolumeToJob(&data.Base, "/etc/hub-token", "hub-token", true)
			data.Base.Timeout = 90
		case "performance":
			if !getBool(item.Value) {
				return
			}
			jobType = getString(item.Key)
			jobNameSuffix = "performance"
			data.Base.Command = performanceScript
		case "latency":
			if !getBool(item.Value) {
				return
			}
			jobType = getString(item.Key)
			jobTemplate = periodicCustomJob
			jobNameSuffix = "latency"
			data.Base.Image = "gcr.io/knative-tests/test-infra/metrics:latest"
			data.Base.Command = "/metrics"
			data.Base.Args = []string{
				fmt.Sprintf("--source-directory=ci-%s-continuous", data.Base.RepoNameForJob),
				"--artifacts-dir=$(ARTIFACTS)",
				"--service-account=" + data.Base.ServiceAccount}
		case "api-coverage":
			if !getBool(item.Value) {
				return
			}
			jobType = getString(item.Key)
			jobTemplate = periodicCustomJob
			jobNameSuffix = "api-coverage"
			data.Base.Image = "gcr.io/knative-tests/test-infra/apicoverage:latest"
			data.Base.Command = "/apicoverage"
			data.Base.Args = []string{
				"--artifacts-dir=$(ARTIFACTS)",
				"--service-account=" + data.Base.ServiceAccount}
		case "custom-job":
			jobType = getString(item.Key)
			jobNameSuffix = getString(item.Value)
		case "cron":
			data.CronString = getString(item.Value)
		case "release":
			version := getString(item.Value)
			jobNameSuffix = version + "-" + jobNameSuffix
			data.Base.RepoBranch = "release-" + version
		case "webhook-apicoverage":
			if !getBool(item.Value) {
				return
			}
			jobNameSuffix = "webhook-apicoverage"
			data.Base.Command = webhookAPICoverageScript
		default:
			continue
		}
		// Knock-out the item, signalling it was already parsed.
		periodicConfig[i] = yaml.MapItem{}
	}
	parseBasicJobConfigOverrides(&data.Base, periodicConfig)
	data.PeriodicJobName = fmt.Sprintf("ci-%s", data.Base.RepoNameForJob)
	if jobNameSuffix != "" {
		data.PeriodicJobName += "-" + jobNameSuffix
	}
	// Ensure required data exist.
	if data.CronString == "" {
		log.Fatalf("Job %q is missing cron string", data.PeriodicJobName)
	}
	if len(data.Base.Args) == 0 && data.Base.Command == "" {
		log.Fatalf("Job %q is missing command", data.PeriodicJobName)
	}
	if jobType == "branch-ci" && data.Base.RepoBranch == "" {
		log.Fatalf("%q jobs are intended to be used on release branches", jobType)
	}
	// Generate config itself.
	data.PeriodicCommand = createCommand(data.Base)
	if data.Base.ServiceAccount != "" {
		addEnvToJob(&data.Base, "GOOGLE_APPLICATION_CREDENTIALS", data.Base.ServiceAccount)
		addEnvToJob(&data.Base, "E2E_CLUSTER_REGION", "us-west1")
	}
	addExtraEnvVarsToJob(&data.Base)
	configureServiceAccountForJob(&data.Base)
	executeJobTemplate("periodic", jobTemplate, title, repoName, data.PeriodicJobName, false, data)
}

// generateCleanupPeriodicJob generates the cleanup job config.
func generateCleanupPeriodicJob() {
	var data periodicJobTemplateData
	data.Base = newbaseProwJobTemplateData("knative/test-infra")
	data.PeriodicJobName = "ci-knative-cleanup"
	data.CronString = cleanupPeriodicJobCron
	data.Base.DecorationConfig = []string{"timeout: 28800000000000"} // 8 hours
	data.Base.Command = cleanupScript
	data.Base.Args = []string{
		"delete-old-gcr-images",
		"--project-resource-yaml ci/prow/boskos/resources.yaml",
		"--days-to-keep 30",
		"--service-account " + data.Base.ServiceAccount,
		"--artifacts $(ARTIFACTS)"}
	addExtraEnvVarsToJob(&data.Base)
	configureServiceAccountForJob(&data.Base)
	executeJobTemplate("periodic cleanup", periodicCustomJob, "presubmits", "", data.PeriodicJobName, false, data)
}

// generateFlakytoolPeriodicJob generates the cleanup job config.
func generateFlakytoolPeriodicJob() {
	var data periodicJobTemplateData
	data.Base = newbaseProwJobTemplateData("knative/test-infra")
	data.Base.Image = flakytoolDockerImage
	data.PeriodicJobName = "ci-knative-flakytool"
	data.CronString = flakytoolPeriodicJobCron
	data.Base.Command = "/flaky-test-reporter"
	data.Base.Args = []string{
		"--service-account " + data.Base.ServiceAccount,
		"--github-account /etc/github-token/token",
		"--slack-account /etc/slack-token/token"}
	addExtraEnvVarsToJob(&data.Base)
	configureServiceAccountForJob(&data.Base)
	addVolumeToJob(&data.Base, "/etc/github-token", "github-token", true)
	addVolumeToJob(&data.Base, "/etc/slack-token", "slack-token", true)
	executeJobTemplate("periodic flakytool", periodicCustomJob, "presubmits", "", data.PeriodicJobName, false, data)
}

// generateBackupPeriodicJob generates the backup job config.
func generateBackupPeriodicJob() {
	var data periodicJobTemplateData
	data.Base = newbaseProwJobTemplateData("none/unused")
	data.Base.ServiceAccount = "/etc/backup-account/service-account.json"
	data.Base.Image = "gcr.io/knative-tests/test-infra/backups:latest"
	data.PeriodicJobName = "ci-knative-backup-artifacts"
	data.CronString = backupPeriodicJobCron
	data.Base.Command = "/backup.sh"
	data.Base.Args = []string{data.Base.ServiceAccount}
	data.Base.ExtraRefs = []string{} // no repo clone required
	addExtraEnvVarsToJob(&data.Base)
	configureServiceAccountForJob(&data.Base)
	executeJobTemplate("periodic backup", periodicCustomJob, "presubmits", "", data.PeriodicJobName, false, data)
}

// generateGoCoveragePeriodic generates the go coverage periodic job config for the given repo (configuration is ignored).
func generateGoCoveragePeriodic(title string, repoName string, periodicConfig yaml.MapSlice) {
	for _, repo := range repositories {
		if repoName != repo.Name || !repo.EnableGoCoverage {
			continue
		}
		var data periodicJobTemplateData
		data.Base = newbaseProwJobTemplateData(repoName)
		data.Base.Image = coverageDockerImage
		data.PeriodicJobName = fmt.Sprintf("ci-%s-go-coverage", data.Base.RepoNameForJob)
		data.CronString = goCoveragePeriodicJobCron
		data.Base.GoCoverageThreshold = repo.GoCoverageThreshold
		data.Base.Command = "/coverage"
		data.Base.Args = []string{
			"--artifacts=$(ARTIFACTS)",
			"--profile-name=coverage_profile.txt",
			"--cov-target=.",
			fmt.Sprintf("--cov-threshold-percentage=%d", data.Base.GoCoverageThreshold)}
		data.Base.ServiceAccount = ""
		addExtraEnvVarsToJob(&data.Base)
		configureServiceAccountForJob(&data.Base)
		executeJobTemplate("periodic go coverage", periodicCustomJob, title, repoName, data.PeriodicJobName, false, data)
		return
	}
}

// generateGoCoveragePostsubmits generates the go coverage postsubmit job configs for all repos.
func generateGoCoveragePostsubmits() {
	for i := range repositories { // Keep order for predictable output.
		if !repositories[i].EnableGoCoverage {
			continue
		}
		repo := repositories[i].Name
		var data postsubmitJobTemplateData
		data.Base = newbaseProwJobTemplateData(repo)
		data.Base.Image = coverageDockerImage
		data.PostsubmitJobName = fmt.Sprintf("post-%s-go-coverage", data.Base.RepoNameForJob)
		addExtraEnvVarsToJob(&data.Base)
		configureServiceAccountForJob(&data.Base)
		executeJobTemplate("postsubmit go coverage", goCoveragePostsubmitJob, "postsubmits", repo, data.PostsubmitJobName, true, data)
		// TODO(adrcunha): remove once the coverage-dev job isn't necessary anymore.
		// Generate config for post-knative-serving-go-coverage-dev right after post-knative-serving-go-coverage
		if data.PostsubmitJobName == "post-knative-serving-go-coverage" {
			data.PostsubmitJobName += "-dev"
			data.Base.Image = strings.Replace(data.Base.Image, "coverage:latest", "coverage-dev:latest-dev", -1)
			executeJobTemplate("presubmit", goCoveragePostsubmitJob, "postsubmits", repo, data.PostsubmitJobName, false, data)
		}
	}
}

// parseSection generate the configs form a given section of the input yaml file.
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

// quote returns the given string quoted if it's not a key/value pair or already quoted.
func quote(s string) string {
	if strings.Contains(s, "\"") || strings.Contains(s, ": ") || strings.HasSuffix(s, ":") {
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

// outputConfig outputs the given line, if not empty, to stdout.
func outputConfig(line string) {
	s := strings.TrimSpace(line)
	if s != "" {
		fmt.Println(line)
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

// main is the script entry point.
func main() {
	// Parse flags and sanity check them.
	var includeConfig = flag.Bool("include-config", true, "Whether to include general configuration (e.g., plank) in the generated config")
	flag.StringVar(&gcsBucket, "gcs-bucket", "knative-prow", "GCS bucket to upload the logs to")
	flag.StringVar(&logsDir, "logs-dir", "logs", "Path in the GCS bucket to upload logs of periodic and post-submit jobs")
	flag.StringVar(&presubmitLogsDir, "presubmit-logs-dir", "pr-logs", "Path in the GCS bucket to upload logs of pre-submit jobs")
	flag.StringVar(&testAccount, "test-account", "/etc/test-account/service-account.json", "Path to the service account JSON for test jobs")
	flag.StringVar(&nightlyAccount, "nightly-account", "/etc/nightly-account/service-account.json", "Path to the service account JSON for nightly release jobs")
	flag.StringVar(&releaseAccount, "release-account", "/etc/release-account/service-account.json", "Path to the service account JSON for release jobs")
	flag.StringVar(&flakytoolDockerImage, "flaky-test-reporter-docker", "gcr.io/knative-tests/test-infra/flaky-test-reporter:latest", "Docker image for flaky test reporting tool")
	flag.StringVar(&coverageDockerImage, "coverage-docker", "gcr.io/knative-tests/test-infra/coverage:latest", "Docker image for coverage tool")
	flag.StringVar(&prowTestsDockerImage, "prow-tests-docker", "gcr.io/knative-tests/test-infra/prow-tests:stable", "prow-tests docker image")
	flag.StringVar(&presubmitScript, "presubmit-script", "./test/presubmit-tests.sh", "Executable for running presubmit tests")
	flag.StringVar(&releaseScript, "release-script", "./hack/release.sh", "Executable for creating releases")
	flag.StringVar(&performanceScript, "performance-script", "./test/performance-tests.sh", "Executable for running performance tests")
	flag.StringVar(&webhookAPICoverageScript, "webhookAPICoverageScript", "./test/apicoverage.sh", "Executable for running webhook apicoverage tool")
	flag.StringVar(&cleanupScript, "cleanup-script", "./tools/cleanup/cleanup.sh", "Executable for running the cleanup tasks")
	flag.StringVar(&repositoryOverride, "repo-override", "", "Repository path (github.com/foo/bar[=branch]) to use instead for a job")
	flag.StringVar(&jobNameFilter, "job-filter", "", "Generate only this job, instead of all jobs")
	flag.StringVar(&preCommand, "pre-command", "", "Executable for running instead of the real command of a job")
	flag.Var(&extraEnvVars, "extra-env", "Extra environment variables (key=value) to add to a job")
	flag.Parse()
	if len(flag.Args()) != 1 {
		log.Fatal("Pass the config file as parameter")
	}
	// We use MapSlice instead of maps to keep key order and create predictable output.
	config := yaml.MapSlice{}
	repositories = make([]repositoryData, 0)
	sectionMap = make(map[string]bool)

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
	if *includeConfig {
		executeTemplate("general config", generalConfig, newbaseProwJobTemplateData(""))
	}
	parseSection(config, "presubmits", generatePresubmit, nil)
	parseSection(config, "periodics", generatePeriodic, generateGoCoveragePeriodic)
	generateCleanupPeriodicJob()
	generateFlakytoolPeriodicJob()
	generateBackupPeriodicJob()
	generateGoCoveragePostsubmits()
}
