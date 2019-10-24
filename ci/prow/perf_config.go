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

// data definitions that are used for the config file generation of performance
// tests cluster maintanence jobs.

package main

import (
	"fmt"
	"log"
)

// generatePerfClusterUpdatePeriodicJobs generates periodic jobs to update clusters
// that run performance testing benchmarks
func generatePerfClusterUpdatePeriodicJobs() {
	for _, repo := range repositories {
		if repo.EnablePerformanceTests {
			perfClusterUpdatePeriodicJob(
				"recreate-clusters",
				recreatePerfClusterPeriodicJobCron,
				"./test/performance/performance-tests.sh",
				[]string{"--recreate-clusters"},
				repo.Name,
				perfClusterJobSecret(repo.Name),
			)
			perfClusterUpdatePeriodicJob(
				"update-clusters",
				updatePerfClusterPeriodicJobCron,
				"./test/performance/performance-tests.sh",
				[]string{"--update-clusters"},
				repo.Name,
				perfClusterJobSecret(repo.Name),
			)
		}
	}
}

// generatePerfClusterReconcilePostsubmitJob generates postsubmit job for the
// repo to reconcile clusters that run performance testing benchmarks.
func generatePerfClusterReconcilePostsubmitJob(repo repositoryData) {
	perfClusterReconcilePostsubmitJob(
		"reconcile-clusters",
		"./test/performance/performance-tests.sh",
		[]string{"--reconcile-clusters"},
		repo.Name,
		perfClusterJobSecret(repo.Name),
	)
}

// perfClusterJobSecret returns the secret we need to mount for the given repo
func perfClusterJobSecret(fullRepoName string) string {
	// secret that needs to be mounted for the job
	var secret string
	switch fullRepoName {
	case "knative/serving":
		secret = "performance-test"
	case "knative/eventing":
		secret = "eventing-performance-test"
	default:
		log.Fatalf("Secret hasn't been added for repo %q for performance testing", fullRepoName)
	}
	return secret
}

func perfClusterUpdatePeriodicJob(jobNamePostFix, cronString, command string, args []string, repo, sa string) {
	var data periodicJobTemplateData
	data.Base = perfClusterBaseProwJob(command, args, repo, sa)
	data.Base.ExtraRefs = append(data.Base.ExtraRefs, "  base_ref: "+data.Base.RepoBranch)
	data.Base.ExtraRefs = append(data.Base.ExtraRefs, "  path_alias: knative.dev/"+data.Base.RepoName)
	data.PeriodicJobName = fmt.Sprintf("ci-%s-%s", data.Base.RepoNameForJob, jobNamePostFix)
	data.CronString = cronString
	data.PeriodicCommand = createCommand(data.Base)
	addMonitoringPubsubLabelsToJob(&data.Base, data.PeriodicJobName)
	executeJobTemplate("performance tests periodic", readTemplate(periodicTestJob),
		"periodics", repo, data.PeriodicJobName, false, data)
}

func perfClusterReconcilePostsubmitJob(jobNamePostFix, command string, args []string, repo, sa string) {
	var data postsubmitJobTemplateData
	data.Base = perfClusterBaseProwJob(command, args, repo, sa)
	data.Base.PathAlias = "path_alias: knative.dev/" + data.Base.RepoName
	data.PostsubmitJobName = fmt.Sprintf("post-%s-%s", data.Base.RepoNameForJob, jobNamePostFix)
	data.PostsubmitCommand = createCommand(data.Base)
	addMonitoringPubsubLabelsToJob(&data.Base, data.PostsubmitJobName)
	executeJobTemplate("performance tests postsubmit", readTemplate(perfPostsubmitJob),
		"postsubmits", repo, data.PostsubmitJobName, true, data)
}

func perfClusterBaseProwJob(command string, args []string, fullRepoName, sa string) baseProwJobTemplateData {
	base := newbaseProwJobTemplateData(fullRepoName)
	for _, repo := range repositories {
		if fullRepoName == repo.Name && repo.Go113 {
			base.Image = getGo113ImageName(base.Image)
			break
		}
	}

	base.Command = command
	base.Args = args
	addVolumeToJob(&base, "/etc/performance-test", sa, true, "")
	addEnvToJob(&base, "GOOGLE_APPLICATION_CREDENTIALS", "/etc/performance-test/service-account.json")
	addEnvToJob(&base, "GITHUB_TOKEN", "/etc/performance-test/github-token")
	addEnvToJob(&base, "SLACK_READ_TOKEN", "/etc/performance-test/slack-read-token")
	addEnvToJob(&base, "SLACK_WRITE_TOKEN", "/etc/performance-test/slack-write-token")
	return base
}
