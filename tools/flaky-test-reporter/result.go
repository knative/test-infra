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

// result.go contains structs and functions for shared data

package main

import (
	"strings"
	"log"
	"fmt"
	"encoding/json"
	"io/ioutil"
	"path"
	"path/filepath"
	"sort"

	"github.com/knative/test-infra/shared/prow"
	"github.com/knative/test-infra/shared/junit"
)

// RepoData struct contains all configurations and test results for a repo
type RepoData struct {
	Config             *JobConfig
	TestStats          map[string]*TestStat // key is test full name
	BuildIDs           []int // all build IDs scanned in this run
	LastBuildStartTime *int64 // timestamp, determines how fresh the data is
}

// JobConfig is initial configuration for a given repo, defines which job to scan
type JobConfig struct {
	Name  string
	Repo  string
	Type  string
}

// TestStat represents test results of a single testcase across all builds,
// Passed, Skipped and Failed contains buildIDs with corresponding results
type TestStat struct {
	TestName string
	Passed   []int
	Skipped  []int
	Failed   []int
}

func (ts *TestStat) isFlaky() bool {
	return ts.hasEnoughRuns() && len(ts.Failed) > 0 && len(ts.Passed) != 0
}

func (ts *TestStat) isPassed() bool {
	return ts.hasEnoughRuns() && len(ts.Failed) == 0
}

func (ts *TestStat) hasEnoughRuns() bool {
	return len(ts.Passed) + len(ts.Failed) >= requiredCount
}

func getFlakyRate(testStats map[string]TestStat) (float32, error) {
	totalCount := len(testStats)
	if 0 == totalCount {
		return 0.0, nil
	}
	flakyCount := 0
	for _, ts := range testStats {
		if ts.isFlaky() {
			flakyCount++
		}
	}
	return float32(flakyCount)/float32(totalCount), nil
}

// createArtifactForRepo marshals RepoData into json format and stores it in a json file,
// under local artifacts directory
func createArtifactForRepo(rd *RepoData) error {
	outFilePath := path.Join(prow.GetLocalArtifactsDir(), rd.Config.Repo + ".json")
	contents, err := json.Marshal(*rd)
	if nil != err {
		return err
	}
	return ioutil.WriteFile(outFilePath, contents, 0644)
}

// addSuiteToRepoData adds all testCase from suite into RepoData
func addSuiteToRepoData(suite *junit.TestSuite, buildID int, rd *RepoData) {
	if nil == rd.TestStats {
		rd.TestStats = make(map[string]*TestStat)
	}
	for _, testCase := range suite.TestCases {
		testFullName := fmt.Sprintf("%s.%s", suite.Name, testCase.Name)
		if _, ok := rd.TestStats[testFullName]; !ok {
			rd.TestStats[testFullName] = &TestStat{TestName: testFullName}
		}
		switch testCase.GetTestStatus() {
			case junit.Passed:
				rd.TestStats[testFullName].Passed = append(rd.TestStats[testFullName].Passed, buildID)
			case junit.Skipped:
				rd.TestStats[testFullName].Skipped = append(rd.TestStats[testFullName].Skipped, buildID)
			case junit.Failed:
				rd.TestStats[testFullName].Failed = append(rd.TestStats[testFullName].Failed, buildID)
		}
	}
}

// getCombinedResultsForBuild gets all junit results from a build,
// and converts each one into a junit TestSuites struct
func getCombinedResultsForBuild(build *prow.Build) ([]*junit.TestSuites, error) {
	var allSuites []*junit.TestSuites
	for _, artifact := range build.GetArtifacts() {
		_, fileName := filepath.Split(artifact)
		if ! strings.HasPrefix(fileName, "junit_") || ! strings.HasSuffix(fileName, ".xml") {
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

// collectTestResultsForRepo collects test results, build IDs from all builds,
// as well as LastBuildStartTime, and stores them in RepoData
func collectTestResultsForRepo(jc *JobConfig) (*RepoData, error) {
	rd := &RepoData{Config: jc}
	job := prow.NewJob(jc.Name, jc.Type, jc.Repo, 0)
	if !job.PathExists() {
		return rd, fmt.Errorf("job path not exist '%s'", jc.Name)
	}
	builds := getLatestFinishedBuilds(job, buildsCount)
	
	log.Printf("latest builds: ")
	for i, build := range builds {
		log.Printf("\t%d", build.BuildID)
		rd.BuildIDs = append(rd.BuildIDs, build.BuildID)
		if 0 == i { // This is the latest build as builds are sorted by start time in descending order
			rd.LastBuildStartTime = build.StartTime
		}
		combinedResults, err := getCombinedResultsForBuild(&build)
		if nil != err {
			return nil, err
		}
		for _, suites := range combinedResults {
			for _, suite := range suites.Suites {
				addSuiteToRepoData(&suite, build.BuildID, rd)
			}
		}
	}
	return rd, nil
}

// getLatestFinishedBuilds is an inexpensive way of listing latest finished builds, in comparing to
// the GetLatestBuilds function from prow package, as it doesn't precompute start/finish time before sorting.
// This function takes the assumption that build IDs are always incremental integers
func getLatestFinishedBuilds(job *prow.Job, count int) []prow.Build {
	var builds []prow.Build
	buildIDs := job.GetBuildIDs()
	fmt.Println(len(buildIDs))
	sort.Sort(sort.Reverse(sort.IntSlice(buildIDs)))
	for _, buildID := range buildIDs {
		if len(builds) >= count {
			break
		}
		build := job.NewBuild(buildID)
		if nil != build.FinishTime {
			builds = append(builds, *build)
		}
	}
	return builds
}
