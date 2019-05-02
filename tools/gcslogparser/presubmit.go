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
	"path"
	"strconv"
	"sync"

	"github.com/knative/test-infra/shared/prow"
)

type prInfo struct {
	repoName string
	ID       int
}

func (c *Client) processPR(wg *sync.WaitGroup) {
	for {
		select {
		case pr := <-c.PrChan:
			for _, j := range prow.GetJobsFromPullRequest(pr.repoName, pr.ID) {
				if len(c.JobFilter) > 0 && !sliceContains(c.JobFilter, j.Name) {
					continue
				}
				c.JobChan <- j
			}
			wg.Done()
		}
	}
}

func (c *Client) feedPresubmitJobsFromRepo(repoName string) {
	wg := &sync.WaitGroup{}

	for i := 0; i < 100; i++ {
		go c.processPR(wg)
	}

	for _, pr := range prow.GetPullRequestsFromRepo(repoName) {
		if ID, _ := strconv.Atoi(path.Base(pr)); -1 != ID {
			wg.Add(1)
			c.PrChan <- prInfo{
				repoName: repoName,
				ID:       ID,
			}
		}
	}

	wg.Wait()
}
