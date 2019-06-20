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

package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"path"

	"github.com/knative/test-infra/shared/prow"
)

type JSONReport struct {
	Flaky []string `json:"flaky"`
}

func (jr *JSONReport) writeToArtifactsDir(repo string, dryrun bool) error {
	artifactsDir := prow.GetLocalArtifactsDir()
	// no need to create the directory, it is guaranteed to exist since artifacts were generated earlier
	outFilePath := path.Join(artifactsDir, repo, "all-flaky-tests.json")
	contents, err := json.Marshal(jr)
	if nil != err {
		return err
	}
	return run(
		fmt.Sprintf("Create JSON report for repo '%s'", repo),
		func() error {
			return ioutil.WriteFile(outFilePath, contents, 0644)
		},
		dryrun,
	)
}

// we want all flaky tests across all jobs for a given repo. There can be duplicate flaky tests
// if multiple jobs use the same test cases. We remove duplicates by using a map as a "set".
func getFlakyTestSet(repo string, repoDataAll []*RepoData) []string {
	testSet := make(map[string]bool)
	for _, rd := range repoDataAll {
		if rd.Config.Repo == repo {
			for _, test := range getFlakyTests(rd) {
				testSet[test] = true
			}
		}
	}
	var flakyTests []string
	for test := range testSet {
		flakyTests = append(flakyTests, test)
	}
	return flakyTests
}

func writeFlakyTestsToJSON(repoDataAll []*RepoData, dryrun bool) error {
	for repo := range jobConfigs {
		output := JSONReport{Flaky: getFlakyTestSet(repo, repoDataAll)}
		err := output.writeToArtifactsDir(repo, dryrun)
		if err != nil {
			return err
		}
		if dryrun {
			log.Printf("[dry run] JSON report not written to disk.\n")
		}
	}
	return nil
}
