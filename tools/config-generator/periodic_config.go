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

// data definitions that are used for the config file generation of periodic prow jobs

package main

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"hash/fnv"
	"log"
	"path"
	"strings"

	"gopkg.in/yaml.v2"
)

const (
	// Template for periodic test/release jobs.
	periodicTestJob = "prow_periodic_test_job.yaml"

	// Template for periodic custom jobs.
	periodicCustomJob = "prow_periodic_custom_job.yaml"

	// Cron strings for key jobs
	goCoveragePeriodicJobCron          = "0 1 * * *"   // Run at 01:00 every day
	recreatePerfClusterPeriodicJobCron = "30 07 * * *" // Run at 00:30PST every day (07:30 UTC)
	updatePerfClusterPeriodicJobCron   = "5 * * * *"   // Run every hour
)

// periodicJobTemplateData contains data about a periodic Prow job.
type periodicJobTemplateData struct {
	Base            baseProwJobTemplateData
	PeriodicJobName string
	CronString      string
	PeriodicCommand []string
}

func (p periodicJobTemplateData) Clone() periodicJobTemplateData {
	var r periodicJobTemplateData
	var err error
	buff := new(bytes.Buffer)
	enc := gob.NewEncoder(buff)
	dec := gob.NewDecoder(buff)
	if err = enc.Encode(&p); err != nil {
		panic(err)
	}
	if err = dec.Decode(&r); err != nil {
		panic(err)
	}
	return r
}

func getUTCtime(i int) int {
	r := i + 7
	if r > 23 {
		return r - 24
	}
	return r
}

func calculateMinuteOffset(str ...string) int {
	h := fnv.New32a()
	for _, s := range str {
		h.Write([]byte(s))
	}
	return int(h.Sum32()) % 60
}

// Generate cron string based on job type, offset generated from jobname
// instead of assign random value to ensure consistency among runs,
// timeout is used for determining how many hours apart
func generateCron(jobType, jobName, repoName string, timeout int) string {
	minutesOffset := calculateMinuteOffset(jobType, jobName)
	// Determines hourly job inteval based on timeout
	hours := int((timeout+5)/60) + 1 // Allow at least 5 minutes between runs
	hourCron := fmt.Sprintf("%d * * * *", minutesOffset)
	if hours > 1 {
		hourCron = fmt.Sprintf("%d */%d * * *", minutesOffset, hours)
	}
	daily := func(pacificHour int) string {
		return fmt.Sprintf("%d %d * * *", minutesOffset, getUTCtime(pacificHour))
	}
	weekly := func(pacificHour, dayOfWeek int) string {
		return fmt.Sprintf("%d %d * * %d", minutesOffset, getUTCtime(pacificHour), dayOfWeek)
	}

	var res string
	switch jobType {
	case "continuous", "custom-job", "auto-release": // As much as every hour
		res = hourCron
	case "branch-ci":
		res = daily(1) // 1 AM
	case "nightly":
		res = daily(2) // 2 AM
	case "dot-release":
		if strings.HasSuffix(repoName, "-operator") {
			// Every Tuesday noon
			res = weekly(12, 2)
		} else {
			// Every Tuesday 2 AM
			res = weekly(2, 2)
		}
	case "webhook-apicoverage":
		res = daily(2) // 2 AM
	default:
		log.Printf("job type not supported for cron generation '%s'", jobName)
	}
	return res
}

// generatePeriodic generates periodic job configs for the given repo and configuration.
// Normally it generates one job per call
// But if it is continuous or branch-ci job, it generates a second job for beta testing of new prow-tests images
func generatePeriodic(title string, repoName string, periodicConfig yaml.MapSlice) {
	var data periodicJobTemplateData
	data.Base = newbaseProwJobTemplateData(repoName)
	jobNameSuffix := ""
	jobTemplate := readTemplate(periodicTestJob)
	jobType := ""
	isMonitoredJob := false
	isContinuousJob := false

	// Parse the input yaml and set values data based on them
	for i, item := range periodicConfig {
		switch item.Key {
		case "continuous":
			if !getBool(item.Value) {
				return
			}
			jobType = getString(item.Key)
			jobNameSuffix = "continuous"
			isContinuousJob = true
			isMonitoredJob = true
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
			isMonitoredJob = true
		case "branch-ci":
			if !getBool(item.Value) {
				return
			}
			jobType = getString(item.Key)
			jobNameSuffix = "continuous"
			isContinuousJob = true
			data.Base.Command = releaseScript
			data.Base.Args = releaseLocal
			setupDockerInDockerForJob(&data.Base)
			// TODO(adrcunha): Consider reducing the timeout in the future.
			data.Base.Timeout = 180
			isMonitoredJob = true
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
			addVolumeToJob(&data.Base, "/etc/hub-token", "hub-token", true, "")
			data.Base.Timeout = 90
			isMonitoredJob = true
		case "custom-job":
			jobType = getString(item.Key)
			jobNameSuffix = getString(item.Value)
			data.Base.Timeout = 100
		case "cron":
			data.CronString = getString(item.Value)
		case "release":
			version := getString(item.Value)
			jobNameSuffix = version + "-" + jobNameSuffix
			data.Base.RepoBranch = "release-" + version
			if jobType == "dot-release" {
				data.Base.Args = append(data.Base.Args, "--branch release-"+version)
			}
			isMonitoredJob = true
		case "webhook-apicoverage":
			if !getBool(item.Value) {
				return
			}
			jobType = getString(item.Key)
			jobNameSuffix = "webhook-apicoverage"
			data.Base.Command = webhookAPICoverageScript
			data.Base.addEnvToJob("SYSTEM_NAMESPACE", data.Base.RepoNameForJob)
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
	if isMonitoredJob {
		addMonitoringPubsubLabelsToJob(&data.Base, data.PeriodicJobName)
	}
	if data.CronString == "" {
		data.CronString = generateCron(jobType, data.PeriodicJobName, data.Base.RepoName, data.Base.Timeout)
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
		data.Base.addEnvToJob("GOOGLE_APPLICATION_CREDENTIALS", data.Base.ServiceAccount)
		data.Base.addEnvToJob("E2E_CLUSTER_REGION", "us-central1")
	}
	if data.Base.RepoBranch != "" && data.Base.RepoBranch != "master" {
		// If it's a release version, add env var PULL_BASE_REF as ref name of the base branch.
		// The reason for having it is in https://github.com/knative/test-infra/issues/780.
		data.Base.addEnvToJob("PULL_BASE_REF", data.Base.RepoBranch)
	}
	addExtraEnvVarsToJob(extraEnvVars, &data.Base)
	configureServiceAccountForJob(&data.Base)

	// This is where the data actually gets written out
	executeJobTemplate("periodic", jobTemplate, title, repoName, data.PeriodicJobName, false, data)

	// If job is a continuous run, add a duplicate for pre-release testing of new prow-tests image
	// It will (mostly) run less often than source job
	if isContinuousJob {
		betaData := data.Clone()

		// Change the name and image
		betaData.PeriodicJobName += "-beta-prow-tests"
		betaData.Base.Image = strings.ReplaceAll(betaData.Base.Image, ":stable", ":beta")

		// Run 2 or 3 times a day because prow-tests beta testing has different desired interval than the underlying job
		hours := []int{getUTCtime(1), getUTCtime(4)}
		if jobType == "continuous" { // as opposed to branch-ci
			// These jobs run 8-24 times per day, so it matters more if they break
			// So test them slightly more often
			hours = append(hours, getUTCtime(15))
		}
		var hoursStr []string
		for _, h := range hours {
			hoursStr = append(hoursStr, fmt.Sprint(h))
		}
		betaData.CronString = fmt.Sprintf("%d %s * * *",
			calculateMinuteOffset(jobType, betaData.PeriodicJobName),
			strings.Join(hoursStr, ","))

		// Write out our duplicate job
		executeJobTemplate("periodic", jobTemplate, title, repoName, betaData.PeriodicJobName, false, betaData)

		// Setup TestGrid here?
		// Each job becomes one of "test_groups"
		// Then we want our own "dashboard" separate from others
		// With each one of the jobs (aka "test_groups") in the single dashboard group
		metaData.AddNonAlignedTest(NonAlignedTestGroup{
			DashboardGroup: "prow-tests",
			DashboardName:  "beta-prow-tests",
			HumanTabName:   data.PeriodicJobName, // this is purposefully not betaData, so the display name is the original CI job name
			CIJobName:      betaData.PeriodicJobName,
			Extra:          nil,
		})
	}
}

// generateGoCoveragePeriodic generates the go coverage periodic job config for the given repo (configuration is ignored).
func generateGoCoveragePeriodic(title string, repoName string, _ yaml.MapSlice) {
	var repo *repositoryData
	// Find a repository entry where repo name matches and Go Coverage is enabled
	for i, repoI := range repositories {
		if repoName != repoI.Name || !repoI.EnableGoCoverage {
			continue
		}
		repo = &repositories[i]
		break
	}
	if repo != nil && repo.EnableGoCoverage {
		repo.Processed = true
		var data periodicJobTemplateData
		data.Base = newbaseProwJobTemplateData(repoName)
		data.Base.Image = coverageDockerImage
		data.PeriodicJobName = fmt.Sprintf("ci-%s-go-coverage", data.Base.RepoNameForJob)
		data.CronString = goCoveragePeriodicJobCron
		data.Base.GoCoverageThreshold = repo.GoCoverageThreshold
		data.Base.Command = "/coverage"
		data.Base.Args = []string{
			"--artifacts=$(ARTIFACTS)",
			fmt.Sprintf("--cov-threshold-percentage=%d", data.Base.GoCoverageThreshold)}
		data.Base.ServiceAccount = ""
		data.Base.ExtraRefs = append(data.Base.ExtraRefs, "  base_ref: "+data.Base.RepoBranch)
		if repo.DotDev {
			data.Base.ExtraRefs = append(data.Base.ExtraRefs, "  path_alias: knative.dev/"+path.Base(repoName))
		}
		if repo.Go114 {
			data.Base.SetGoVersion(GoVersion{1, 14})
		}
		addExtraEnvVarsToJob(extraEnvVars, &data.Base)
		addMonitoringPubsubLabelsToJob(&data.Base, data.PeriodicJobName)
		configureServiceAccountForJob(&data.Base)
		executeJobTemplate("periodic go coverage", readTemplate(periodicCustomJob), title, repoName, data.PeriodicJobName, false, data)
	}
}
