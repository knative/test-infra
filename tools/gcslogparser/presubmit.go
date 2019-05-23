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
	"path"
	"strconv"

	"github.com/knative/test-infra/shared/prow"
)

type prInfo struct {
	repoName string
	ID       int
}

func (c *Parser) processPR() {
	for {
		select {
		case pr := <-c.PrChan:
			for _, j := range prow.GetJobsFromPullRequest(pr.repoName, pr.ID) {
				if len(c.jobFilter) > 0 && !sliceContains(c.jobFilter, j.Name) {
					continue
				}
				c.wgJob.Add(1)
				c.jobChan <- j
			}
			c.wgPR.Done()
		}
	}
}

func (c *Parser) feedPresubmitJobsFromRepo(repoName string) {
	for i := 0; i < 100; i++ {
		go c.processPR()
	}

	for _, pr := range prow.GetPullRequestsFromRepo(repoName) {
		if ID, _ := strconv.Atoi(path.Base(pr)); -1 != ID {
			c.wgPR.Add(1)
			c.PrChan <- prInfo{
				repoName: repoName,
				ID:       ID,
			}
		}
	}
}
