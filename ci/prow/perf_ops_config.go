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

// generatePerfClusterUpdatePeriodicJobs generates periodic jobs to update clusters
// that run performance testing benchmarks
func generatePerfClusterUpdatePeriodicJobs() {
	// Generate periodic performance jobs for serving
	perfClusterUpdatePeriodicJob(
		"ci-knative-serving-recreate-clusters",
		recreatePerfClusterPeriodicJobCron,
		"./test/performance/tools/recreate_clusters.sh",
		[]string{},
		"serving",
		"performance-test",
	)
	perfClusterUpdatePeriodicJob(
		"ci-knative-serving-update-clusters",
		updatePerfClusterPeriodicJobCron,
		"./test/performance/tools/update_clusters.sh",
		[]string{},
		"serving",
		"performance-test",
	)

	// Generate periodic performance jobs for eventing
	perfClusterUpdatePeriodicJob(
		"ci-knative-eventing-recreate-clusters",
		recreatePerfClusterPeriodicJobCron,
		"./test/performance/performance-tests.sh",
		[]string{"--recreate-clusters"},
		"eventing",
		"eventing-performance-test",
	)
	perfClusterUpdatePeriodicJob(
		"ci-knative-eventing-update-clusters",
		updatePerfClusterPeriodicJobCron,
		"./test/performance/performance-tests.sh",
		[]string{"--update-clusters"},
		"eventing",
		"eventing-performance-test",
	)
}

// generatePerfClusterReconcilePostsubmitJob generates postsubmit job for the
// repo to reconcile clusters that run performance testing benchmarks.
func generatePerfClusterReconcilePostsubmitJob(repo string) {
	// TODO(chizhg): also add postsubmit jobs for other repos after they set up
	// performance test infra.
	if repo == "knative/eventing" {
		perfClusterReconcilePostsubmitJob(
			"post-knative-eventing-reconcile-clusters",
			"./test/performance/performance-tests.sh",
			[]string{"--reconcile-clusters"},
			"eventing",
			"eventing-performance-test",
		)
	}
}

func perfClusterUpdatePeriodicJob(jobName, cronString, command string, args []string, repo, sa string) {
	var data periodicJobTemplateData
	data.Base = perfClusterBaseProwJob(jobName, command, args, repo, sa)
	data.PeriodicJobName = jobName
	data.CronString = cronString
	data.PeriodicCommand = createCommand(data.Base)
	executeJobTemplate(jobName, readTemplate(periodicTestJob), "periodics", "knative/"+repo, jobName, false, data)
}

func perfClusterReconcilePostsubmitJob(jobName, command string, args []string, repo, sa string) {
	var data postsubmitJobTemplateData
	data.Base = perfClusterBaseProwJob(jobName, command, args, repo, sa)
	data.PostsubmitJobName = jobName
	data.PostsubmitCommand = createCommand(data.Base)
	executeJobTemplate(jobName, readTemplate(perfOpsPostsubmitJob), "postsubmits", "knative/"+repo, jobName, true, data)
}

func perfClusterBaseProwJob(jobName, command string, args []string, repoName, sa string) baseProwJobTemplateData {
	fullRepoName := "knative/" + repoName
	base := newbaseProwJobTemplateData(fullRepoName)
	for _, repo := range repositories {
		if fullRepoName == repo.Name && repo.Go113 {
			base.Image = getGo113ImageName(base.Image)
			break
		}
	}
	base.PathAlias = "path_alias: knative.dev/" + repoName
	base.ExtraRefs = append(base.ExtraRefs, "  base_ref: "+base.RepoBranch)
	base.ExtraRefs = append(base.ExtraRefs, "  path_alias: knative.dev/"+repoName)

	base.Command = command
	base.Args = args
	addVolumeToJob(&base, "/etc/performance-test", sa, true, "")
	addEnvToJob(&base, "GOOGLE_APPLICATION_CREDENTIALS", "/etc/performance-test/service-account.json")
	// TODO(chizhg): remove PERF_TEST_GOOGLE_APPLICATION_CREDENTIALS once
	// serving is also using the new performance-tests.sh
	addEnvToJob(&base, "PERF_TEST_GOOGLE_APPLICATION_CREDENTIALS", "/etc/performance-test/service-account.json")
	addEnvToJob(&base, "GITHUB_TOKEN", "/etc/performance-test/github-token")
	addEnvToJob(&base, "SLACK_READ_TOKEN", "/etc/performance-test/slack-read-token")
	addEnvToJob(&base, "SLACK_WRITE_TOKEN", "/etc/performance-test/slack-write-token")
	addMonitoringPubsubLabelsToJob(&base, jobName)
	return base
}
