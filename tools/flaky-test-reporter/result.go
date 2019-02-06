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

// result.go contains structs and functions for common results

package main

import (
	"strings"
	"strconv"
	"log"
	"fmt"
	"io/ioutil"

	"github.com/knative/test-infra/shared/junit"
)

const (
	buildsCount         = 10 // Builds to be analyzed
	requiredCount       = 8 // Minimal number of results to be counted as valid results for each testcase
	threshold           = 0.05 // Don't do anything if found more than 5% tests flaky
)

// repoData struct contains all configurations and test results for a repo
type repoData struct {
	config *jobConfig
	testStats []*testStat
	buildIDs  []int // all build IDs scanned in this run
	lastBuildStartTime *int64 // useful for 
}

// jobConfig is initial configuration for a given repo, defines which job to scan
type jobConfig struct {
	name  string
	repo  string
	t     string // Use shorthand as "type" is a reserved name
}

// testStat represents test results of a single testcase across all builds,
// passed, skipped and failed contains buildIDs with corresponding results
type testStat struct {
	testName string
	passed   []int
	skipped  []int
	failed   []int
}

func (ts *testStat) isValid() bool {
	return len(ts.passed) + len(ts.failed) >= requiredCount
}

func (ts *testStat) isFlaky() bool {
	return ts.isValid() && len(ts.failed) > 0
}

func (ts *testStat) isPassed() bool {
	return ts.isValid() && len(ts.failed) == 0
}

func readFromJunitFile(filePath string) *junit.TestSuites {
	contents, err := ioutil.ReadFile(filePath)
	if nil != err {
		log.Fatalf("%v", err)
		return nil
	}
	suites, err := junit.UnMarshal(contents)
	if nil != err {
		log.Fatalf("%v", err)
		return nil
	}
	return suites
}
