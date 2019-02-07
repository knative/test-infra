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
	return ts.isValid() && len(ts.Failed) > 0
}

func (ts *TestStat) isPassed() bool {
	return ts.isValid() && len(ts.Failed) == 0
}

func (ts *TestStat) isValid() bool {
	return len(ts.Passed) + len(ts.Failed) >= requiredCount
}

func (rd *RepoData) getFlakyTests() map[string]*TestStat {
	res := make(map[string]*TestStat)
	for testName, ts := range rd.TestStats {
		if ts.isFlaky() {
			res[testName] = rd.TestStats[testName]
		}
	}
	return res
}
