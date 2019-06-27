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

// flaky-test-retryer detects failed integration jobs on new pull requests,
// determines if they failed due to flaky tests, posts comments describing the
// issue, and retries them until they succeed.

package main

import (
	"log"

	"github.com/knative/test-infra/tools/flaky-test-reporter/jsonreport"
	// TODO: remove this import once "k8s.io/test-infra" import problems are fixed
	// https://github.com/test-infra/test-infra/issues/912
	"github.com/knative/test-infra/tools/monitoring/prowapi"
)

// TODO(TrevorFarrelly): supported repos depend on flaky-test-reporter's
// configuration. Unify the reporter and retryer configs somehow?
var repos = []string{"serving"}

func validStatus(msg *prowapi.ReportMessage) bool {
	return msg.Status == prowapi.FailureState
}

func validType(msg *prowapi.ReportMessage) bool {
	return msg.JobType == prowapi.PresubmitJob
}

func validRepo(msg *prowapi.ReportMessage) bool {
	if len(msg.Refs) > 0 {
		for _, repo := range repos {
			if msg.Refs[0].Repo == repo {
				return true
			}
		}
	}
	return false
}

// spawn a goroutine for each valid message received
func handleDriver(msg *prowapi.ReportMessage) {
	if validStatus(msg) && validType(msg) && validRepo(msg) {
		go handleMessage(msg)
	}
}

func handleMessage(msg *prowapi.ReportMessage) {
	log.Printf("Message fit all criteria - Starting analysis\n")
	// verify job failed due to tests, not build issues. Get current flaky tests,
	// cross-reference failed tests with flaky tests, retest if the failed tests are flaky.
}
