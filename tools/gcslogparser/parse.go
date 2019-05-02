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
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"sync"

	"github.com/knative/test-infra/shared/prow"
)

type logInfo struct {
	build prow.Build
	l     string
}

func (c *Client) processJob(wgBuild *sync.WaitGroup) {
	for {
		select {
		case j := <-c.JobChan:
			// fmt.Println(j.StoragePath)
			for _, buildID := range j.GetBuildIDs() {
				wgBuild.Add(1)
				c.BuildChan <- *(j.NewBuild(buildID))
			}
		}
	}
}

func (c *Client) processBuild(wgBuild, wgLog *sync.WaitGroup) {
	for {
		select {
		case b := <-c.BuildChan:
			// fmt.Println(b.StoragePath)
			if b.FinishTime != nil && *b.FinishTime > c.StartDate.Unix() {
				wgLog.Add(1)
				content, _ := b.ReadFile("build-log.txt")
				c.LogChan <- logInfo{
					build: b,
					l:     string(content),
				}
			}
			wgBuild.Done()
		}
	}
}

func (c *Client) parseLog(wgLog *sync.WaitGroup) {
	for {
		select {
		case l := <-c.LogChan:
			if c.Parser(l.l) {
				fmt.Println(l.build.StoragePath)
			} else {
				fmt.Println("not found")
			}
			wgLog.Done()
		}
	}
}

// ParseRepo parses all jobs within a repo, including Presubmit and Postsubmit
func (c *Client) ParseRepo(repoName string) {
	c.refreshChans()

	wgBuild := &sync.WaitGroup{}
	wgLog := &sync.WaitGroup{}

	for i := 0; i < 500; i++ {
		go c.processJob(wgBuild)
	}
	for i := 0; i < 500; i++ {
		go c.processBuild(wgBuild, wgLog)
	}
	for i := 0; i < 500; i++ {
		go c.parseLog(wgLog)
	}

	c.feedPresubmitJobsFromRepo(repoName)
	c.feedPostsubmitJobsFromRepo(repoName)

	wgBuild.Wait()
	wgLog.Wait()
}

func main() {
	serviceAccount := flag.String("service-account", os.Getenv("GOOGLE_APPLICATION_CREDENTIALS"), "JSON key file for GCS service account")
	// repoNames := flag.String("repo", "test-infra", "repo to be downloaded")
	startDate := flag.String("start-date", "2019-02-22", "cut off date to be analyzed")
	dryrun := flag.Bool("dry-run", false, "dry run switch")
	flag.Parse()

	if nil != dryrun && true == *dryrun {
		log.Printf("running in [dry run mode]")
	}

	c, _ := NewClient(*serviceAccount, func(s string) bool {
		return strings.Contains(s, "Running coverage for PR")
	})

	defer c.cleanup()
	c.CleanupOnInterrupt()

	c.Dryrun = *dryrun
	c.SetStartDate(*startDate)
	c.ParseRepo("test-infra")
}
