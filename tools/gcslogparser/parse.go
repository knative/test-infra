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

package gcslogparser

import (
	"fmt"
	"sync"

	"github.com/knative/test-infra/shared/prow"
)

type jobInfo struct {
	repoName string
	jobName  string
}

type buildInfo struct {
	job *prow.Job
	ID  int // Build ID
}

type logInfo struct {
	build *prow.Build
	l     string
}

func (c *Client) processJob(jobChan chan prow.Job, buildChan chan prow.Build, wgBuild *sync.WaitGroup) {
	for {
		select {
		case pr := <-prChan:
			// fmt.Printf("Process PR '%v'\n", pr)
			for _, j := range prow.ListJobsFromPR(pr.repoName, pr.ID) {
				// fmt.Printf("job: %s\n", j.Name)
				for _, buildID := range prow.NewJob(j.Name, prow.PresubmitJob, pr.repoName, j.PullID).GetBuildIDs() {
					wgBuild.Add(1)
					// fmt.Printf("build: %s %d\n", j.Name, buildID)
					buildChan <- buildInfo{job: j, ID: buildID}
				}
			}
		}
	}
}

func (c *Client) processBuild(buildChan chan prow.Build, logChan chan logInfo, startTimestamp int64, wgBuild, wgLog *sync.WaitGroup) {
	for {
		select {
		case b := <-buildChan:
			// fmt.Print(b.job, b.ID)
			build := b.job.NewBuild(b.ID)
			if build.FinishTime != nil && *build.FinishTime > startTimestamp {
				wgLog.Add(1)
				content, _ := build.ReadFile("build-log.txt")
				logChan <- logInfo{
					build: build,
					l:     string(content),
				}
			}
			wgBuild.Done()
		}
	}
}

func (c *Client) parseLog(logChan chan logInfo, wgLog *sync.WaitGroup) {
	for {
		select {
		case li := <-logChan:
			if c.Parser(li.l) {
				fmt.Println(li.build.StoragePath)
			}
			wgLog.Done()
		}
	}
}

// ParseRepo parses all jobs within a repo, including Presubmit and Postsubmit
func (c *Client) ParseRepo(repoName, f func() bool, startTime string) {
	jobChan := make(chan prow.Job, 50)
	buildChan := make(chan prow.Build, 500)
	logChan := make(chan logInfo, 500)

	wgBuild := &sync.WaitGroup{}
	wgLog := &sync.WaitGroup{}

	defer func() {
		close(jobChan)
		close(buildChan)
		close(logChan)
	}()

	for i := 0; i < 500; i++ {
		go c.processJob(jobChan, buildChan, wgBuild)
	}
	for i := 0; i < 500; i++ {
		go c.processBuild(buildChan, logChan, startTimestamp, wgBuild, wgLog)
	}
	for i := 0; i < 500; i++ {
		go c.parseLog(logChan, wgLog)
	}

	c.feedPresubmitJobsFromRepo(repoName)
	c.feedPostsubmitJobsFromRepo(repoName)

	wgBuild.Wait()
	wgLog.Wait()
}
