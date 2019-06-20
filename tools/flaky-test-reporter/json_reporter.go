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
	"fmt"
	"log"
	"sync"

	"github.com/TrevorFarrelly/test-infra/shared/jsonreport"
)

func getFlakyTestSet(repoDataAll []*RepoData) map[string]map[string]bool {
	// use a map for each repo as a "set", to eliminate duplicates
	flakyTestSets := map[string]map[string]bool{}
	for _, rd := range repoDataAll {
		if flakyTestSet[rd.Config.Repo] == nil {
			flakyTestSet[rd.Config.Repo] = map[string]bool{}
		}
		for _, test := range getFlakyTests(rd) {
			flakyTestSet[rd.Config.Repo][test] = true
		}
	}
	return flakyTestSet
}

func writeFlakyTestsToJSON(repoDataAll []*RepoData, dryrun bool) error {
	var allErrs []error
	flakyTestSets := getFlakyTestSet(repoDataAll)
	ch := make(chan bool, len(flakyTestSets))
	wg := sync.WaitGroup{}
	for repo, tests := range flakyTestSets {
		wg.Add(1)
		go func(wg *sync.WaitGroup) {
			report := jsonreport.NewReport(repo, nil)
			for test := range tests {
				report.Flaky = append(report.Flaky, test)
			}
			if err := run(
				fmt.Sprintf("writing JSON report for repo '%s'", repo),
				func() error {
					return report.WriteToArtifactsDir()
				},
				dryrun); nil != err {
				allErrs = append(allErrs, err)
				log.Printf("failed writing JSON report for repo '%s': '%v'", err)
			}
			if dryrun {
				log.Printf("[dry run] JSON report not written. See it below:\n%s\n\n", report)
			}
			ch <- true
			wg.Done()
		}(&wg)
	}
	wg.Wait()
	close(ch)
	return combineErrors(allErrs)
}
