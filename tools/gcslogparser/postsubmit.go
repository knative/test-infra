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

// flaky-test-reporter collects test results from continuous flows,
// identifies flaky tests, tracking flaky tests related github issues,
// and sends slack notifications.

package main

import (
	"github.com/knative/test-infra/shared/prow"
)

func (c *Client) feedPostsubmitJobsFromRepo(repoName string, jobChan chan prow.Job) {
	for _, j := range prow.GetPostsubmitJobsFromRepo(repoName) {
		if len(c.JobFilter) > 0 && !c.JobFilter.Contains(j.Name) {
			continue
		}
		jobChan <- j
	}
}
