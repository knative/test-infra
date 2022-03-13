/*
Copyright 2022 The Knative Authors

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

package pkg

import (
	"fmt"
	"hash/fnv"

	"istio.io/test-infra/tools/prowgen/pkg/spec"
)

const (
	// default timeout is 2 hours
	defaultTimeout = 120

	// name of the main branch
	mainBranchName = "main"

	// type of periodic Prow job
	periodicProwJobType = "periodic"
)

// addSchedule calculates and adds schedule for periodic Prow jobs to try to
// distribute the workloads evenly to avoid overloading Prow.
func addSchedule(jobsConfig spec.JobsConfig) spec.JobsConfig {
	org := jobsConfig.Org
	repo := jobsConfig.Repo
	branch := jobsConfig.Branches[0]
	for i, job := range jobsConfig.Jobs {
		// Only add the calculated cron schedule if both Schedule and Cron are
		// empty.
		if hasPeriodic(job.Types) && job.Interval == "" && job.Cron == "" {
			var timeout int
			if job.Timeout != nil {
				timeout = int(job.Timeout.Minutes())
			}
			if timeout == 0 {
				timeout = defaultTimeout
			}
			job.Cron = generateCron(org, repo, branch, job.Name, timeout)
		}
		jobsConfig.Jobs[i] = job
	}

	return jobsConfig
}

func hasPeriodic(pjTypes []string) bool {
	for _, tp := range pjTypes {
		if tp == periodicProwJobType {
			return true
		}
	}
	return false
}

// Generate cron string based on job type, offset generated from jobname
// instead of assign random value to ensure consistency among runs,
// timeout is used for determining how many hours apart
func generateCron(org, repo, branch, jobName string, timeout int) string {
	hourOffset := calculateHourOffset(org, repo, branch, jobName)
	minutesOffset := calculateMinuteOffset(org, repo, branch, jobName)
	// Determines hourly job inteval based on timeout
	hours := int((timeout+5)/60) + 1
	hourCron := fmt.Sprintf("%d */%d * * *", minutesOffset, hours*3)
	daily := func(pacificHour int) string {
		return fmt.Sprintf("%d %d * * *", minutesOffset, utcTime(pacificHour))
	}
	weekly := func(pacificHour, dayOfWeek int) string {
		return fmt.Sprintf("%d %d * * %d", minutesOffset, utcTime(pacificHour), dayOfWeek)
	}

	var res string
	switch jobName {
	case "continuous":
		if branch == mainBranchName {
			res = hourCron // Multiple times per day for main branch continuous Prow jobs
		} else {
			res = daily(hourOffset) // Random hour in the day for release branch continuous Prow jobs
		}
	case "nightly":
		res = daily(2) // nightlys run at 2 AM
	case "release":
		if branch == mainBranchName {
			res = hourCron // auto-release for main branch runs multiple times per day
		} else {
			res = weekly(2, 2) // dot-release for release branches runs every Tuesday 2 AM
		}
	default:
		if repo == "serving" {
			res = hourCron // Multiple times per day for knative/serving periodic Prow jobs
		} else {
			res = daily(hourOffset) // Random hour in the day for other periodic Prow jobs
		}
	}
	return res
}

func utcTime(i int) int {
	r := i + 7
	if r > 23 {
		return r - 24
	}
	return r
}

func calculateMinuteOffset(str ...string) int {
	return calculateHash(str...) % 60
}

func calculateHourOffset(str ...string) int {
	return calculateHash(str...) % 24
}

func calculateHash(str ...string) int {
	h := fnv.New32a()
	for _, s := range str {
		h.Write([]byte(s))
	}
	return int(h.Sum32())
}
