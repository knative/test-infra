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

// log_parser.go collects failed tests from jobs that triggered the retryer and
// finds flaky tests that are relevant to that failed job.

package main

import (
	"github.com/knative/test-infra/tools/flaky-test-reporter/jsonreport"
)

// InitLogParser configures jsonreport's dependencies
func InitLogParser(serviceAccount string) error {
	return jsonreport.Initialize(serviceAccount)
}

// getReportRepos gets all of the repositories where we are reporting flaky tests
// TODO: combine this with the function for getting flaky tests, as they will be
//       almost exactly the same thing.
func getReportRepos() ([]string, error) {
	var repos []string
	if reports, err := jsonreport.GetFlakyTestReport("", -1); err == nil && len(reports) > 0 {
		for _, r := range reports {
			repos = append(repos, r.Repo)
		}
	}
	return repos, err
}
