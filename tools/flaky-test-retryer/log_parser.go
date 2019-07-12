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
	"fmt"
	"path/filepath"
	"strings"

	"github.com/knative/test-infra/shared/junit"
	"github.com/knative/test-infra/shared/prow"
	"github.com/knative/test-infra/tools/flaky-test-reporter/jsonreport"
)

// InitLogParser configures jsonreport's dependencies.
func InitLogParser(serviceAccount string) error {
	return jsonreport.Initialize(serviceAccount)
}

// getFailedTests gets all the tests that failed in the given job.
func getFailedTests(jobName, jobType, repo string, pull int) ([]string, error) {
	job := prow.NewJob(jobName, jobType, repo, pull)
	buildID, err := job.GetLatestBuildNumber()
	if err != nil {
		return nil, err
	}
	build := job.NewBuild(buildID)
	results, err := GetCombinedResultsForBuild(build)
	if err != nil {
		return nil, err
	}
	var tests []string
	for _, suites := range results {
		for _, suite := range suites.Suites {
			for _, test := range suite.TestCases {
				if test.GetTestStatus() == junit.Failed {
					tests = append(tests, fmt.Sprintf("%s.%s", suite.Name, test.Name))
				}
			}
		}
	}
	return tests, nil
}

// TODO: This function is a direct copy-paste of the function in
// tools/flaky-test-reporter/result.go. Refactor it out into a shared library.

// GetCombinedResultsForBuild gets all junit results from a build,
// and converts each one into a junit TestSuites struct
func GetCombinedResultsForBuild(build *prow.Build) ([]*junit.TestSuites, error) {
	var allSuites []*junit.TestSuites
	for _, artifact := range build.GetArtifacts() {
		_, fileName := filepath.Split(artifact)
		if !strings.HasPrefix(fileName, "junit_") || !strings.HasSuffix(fileName, ".xml") {
			continue
		}
		relPath, _ := filepath.Rel(build.StoragePath, artifact)
		contents, err := build.ReadFile(relPath)
		if nil != err {
			return nil, err
		}
		if suites, err := junit.UnMarshal(contents); nil != err {
			return nil, err
		} else {
			allSuites = append(allSuites, suites)
		}
	}
	return allSuites, nil
}

// getFlakyTests gets the latest flaky tests for the given repository.
func getFlakyTests(repo string) ([]string, error) {
	return parseFlakyLog(func(report jsonreport.Report, result *[]string) {
		if report.Repo == repo {
			*result = report.Flaky
		}
	})
}

// getReportRepos gets all of the repositories where we collect flaky tests.
func getReportRepos() ([]string, error) {
	return parseFlakyLog(func(report jsonreport.Report, result *[]string) {
		*result = append(*result, report.Repo)
	})
}

// parseFlakyLog reads the latest flaky test report and returns filtered results based
// on the function the caller passes in.
func parseFlakyLog(f func(report jsonreport.Report, result *[]string)) ([]string, error) {
	var results []string
	reports, err := jsonreport.GetFlakyTestReport("", -1)
	if err == nil && len(reports) > 0 {
		for _, r := range reports {
			f(r, &results)
		}
	}
	return results, err
}

// compareTests compares lists of failed and flaky tests, and returns any outlying failed
// tests, i.e. tests that failed that are NOT flaky.
func getNonFlakyTests(failedTests, flakyTests []string) []string {
	flakyMap := map[string]bool{}
	for _, flaky := range flakyTests {
		flakyMap[flaky] = true
	}
	var notFlaky []string
	for _, failed := range failedTests {
		if _, ok := flakyMap[failed]; !ok {
			notFlaky = append(notFlaky, failed)
		}
	}
	return notFlaky
}
