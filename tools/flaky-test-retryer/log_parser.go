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
	"log"
	"path/filepath"
	"strings"

	"github.com/knative/test-infra/shared/junit"
	"github.com/knative/test-infra/shared/prow"
	"github.com/knative/test-infra/tools/flaky-test-reporter/jsonreport"
	// TODO: remove this import once "k8s.io/test-infra" import problems are fixed
	// https://github.com/test-infra/test-infra/issues/912
	"github.com/knative/test-infra/tools/monitoring/prowapi"
)

// InitLogParser configures jsonreport's dependencies.
func InitLogParser(serviceAccount string) error {
	return jsonreport.Initialize(serviceAccount)
}

// JobData contains stripped-down information about a failed job, and local caches
// of its failed tests and the flaky report it is referencing.
type JobData struct {
	Name         string
	Type         prowapi.ProwJobType
	Status       prowapi.ProwJobState
	Org          string
	Repo         string
	Pull         int
	failedTests  []string
	flakyReports []jsonreport.Report
}

// NewJobData creates a JobData object for the given message, and returns an error
// if it does not fit the criteria we have set for it.
func NewJobData(msg *prowapi.ReportMessage) *JobData {
	jd := &JobData{
		Name:   msg.JobName,
		Type:   msg.JobType,
		Status: msg.Status,
	}
	// add repo data if it exists
	// msg.Refs' first element is always the repo that that triggered this job, later elements
	// are other dependencies the job needed that were not already in the main repository.
	if len(msg.Refs) > 0 {
		jd.Org = msg.Refs[0].Org
		jd.Repo = msg.Refs[0].Repo
		// add pull ID if it exists
		if len(msg.Refs[0].Pulls) > 0 {
			jd.Pull = msg.Refs[0].Pulls[0].Number
		}
	}
	return jd
}

// IsSupported checks to make sure the message can be processed with the current flaky
// test information
func (jd *JobData) IsSupported() bool {
	prefix := fmt.Sprintf("Job did not fit criteria")
	if jd.Status != prowapi.FailureState {
		log.Printf("%s: message did not signal a failure: %v\n", prefix, jd.Status)
		return false
	}
	// check type
	if jd.Type != prowapi.PresubmitJob {
		log.Printf("%s: message did not originate from presubmit: %v\n", prefix, jd.Type)
		return false
	}
	// check repo
	repos, err := jd.getReportRepos()
	if err != nil {
		log.Printf("%s: error getting reporter's repositories: %v\n", prefix, err)
		return false
	}
	expRepo := false
	for _, repo := range repos {
		if jd.Repo == repo {
			expRepo = true
			break
		}
	}
	if !expRepo {
		log.Printf("%s: message's repo is not being analyzed by flaky test reporter: '%v'\n", prefix, jd.Repo)
		return false
	}
	// make sure pull ID exists
	if jd.Pull == 0 {
		log.Printf("%s: message does not have any pull IDs\n", prefix)
		return false
	}
	return true
}

// Logf prefixes a call to log.Printf with information about the job it is being called on
func (jd *JobData) Logf(format string, a ...interface{}) {
	input := append([]interface{}{jd.Repo, jd.Pull, jd.Name}, a...)
	log.Printf("%s/pull/%d: %s: "+format, input...)
}

// getFailedTests gets all the tests that failed in the given job.
func (jd *JobData) getFailedTests() ([]string, error) {
	// use cache if it is populated
	if jd.failedTests != nil {
		return jd.failedTests, nil
	}
	job := prow.NewJob(jd.Name, string(jd.Type), jd.Repo, jd.Pull)
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
	jd.failedTests = tests
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

// getFlakyTests gets the latest flaky tests that could affect this job
func (jd *JobData) getFlakyTests() ([]string, error) {
	return jd.parseFlakyLog(func(report jsonreport.Report, result *[]string) {
		if report.Repo == jd.Repo {
			*result = report.Flaky
		}
	})
}

// getReportRepos gets all of the repositories where we collect flaky tests.
func (jd *JobData) getReportRepos() ([]string, error) {
	return jd.parseFlakyLog(func(report jsonreport.Report, result *[]string) {
		*result = append(*result, report.Repo)
	})
}

// parseFlakyLog reads the latest flaky test report and returns filtered results based
// on the function the caller passes in.
func (jd *JobData) parseFlakyLog(f func(report jsonreport.Report, result *[]string)) ([]string, error) {
	var results []string
	var err error
	// populate cache if it is empty
	if jd.flakyReports == nil {
		jd.flakyReports, err = jsonreport.GetFlakyTestReport("", -1)
	}
	if err == nil && len(jd.flakyReports) > 0 {
		for _, r := range jd.flakyReports {
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
