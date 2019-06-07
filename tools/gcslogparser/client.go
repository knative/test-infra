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
	"os"
	"os/signal"
	"sync"
	"time"

	"github.com/knative/test-infra/shared/prow"
)

type Parser struct {
	StartDate time.Time             // Earliest date to be analyzed, i.e. "2019-02-22"
	logParser func(s string) string // logParser function
	jobFilter []string              // Jobs to be parsed. If not provided will parse all jobs
	PrChan    chan prInfo           // For PR use only, make it here so it's easier to cleanup
	jobChan   chan prow.Job
	buildChan chan buildInfo
	wgPR      sync.WaitGroup
	wgJob     sync.WaitGroup
	wgBuild   sync.WaitGroup

	found     [][]string
	processed []string

	mutex *sync.Mutex
}

type buildInfo struct {
	job prow.Job
	ID  int
}

func NewParser(serviceAccount string) (*Parser, error) {
	if err := prow.Initialize(serviceAccount); nil != err { // Explicit authenticate with gcs Parser
		return nil, fmt.Errorf("Failed authenticating GCS: '%v'", err)
	}

	c := &Parser{}
	c.mutex = &sync.Mutex{}

	c.PrChan = make(chan prInfo, 100)
	c.jobChan = make(chan prow.Job, 100)
	c.buildChan = make(chan buildInfo, 10000)

	for i := 0; i < 100; i++ {
		go c.jobListener()
	}
	for i := 0; i < 1000; i++ {
		go c.buildListener()
	}

	return c, nil
}

func (c *Parser) wait() {
	c.wgPR.Wait()
	c.wgJob.Wait()
	c.wgBuild.Wait()
}

func (c *Parser) setStartDate(startDate string) error {
	tt, err := time.Parse("2006-01-02", startDate)
	if nil != err {
		return fmt.Errorf("invalid start date string, expecing format YYYY-MM-DD: '%v'", err)
	}
	c.StartDate = tt
	return nil
}

func (c *Parser) jobListener() {
	for {
		select {
		case j := <-c.jobChan:
			for _, buildID := range j.GetBuildIDs() {
				c.wgBuild.Add(1)
				c.buildChan <- buildInfo{
					job: j,
					ID:  buildID,
				}
			}
			c.wgJob.Done()
		}
	}
}

func (c *Parser) buildListener() {
	for {
		select {
		case b := <-c.buildChan:
			build := b.job.NewBuild(b.ID)
			if build.FinishTime != nil && *build.FinishTime > c.StartDate.Unix() {
				content, _ := build.ReadFile("build-log.txt")
				found := c.logParser(string(content))
				c.mutex.Lock()
				c.processed = append(c.processed, build.StoragePath)
				if "" != found {
					c.found = append(c.found, []string{found, time.Unix(*build.StartTime, 0).String(), build.StoragePath})
				}
				c.mutex.Unlock()
			}
			c.wgBuild.Done()
		}
	}
}

func (c *Parser) cleanup() {
	if c.PrChan != nil {
		close(c.PrChan)
	}
	if c.jobChan != nil {
		close(c.jobChan)
	}
	if c.buildChan != nil {
		close(c.buildChan)
	}
}

// CleanupOnInterrupt will execute the function cleanup if an interrupt signal is caught
func (c *Parser) CleanupOnInterrupt() {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt)
	go func() {
		for range ch {
			c.cleanup()
			os.Exit(1)
		}
	}()
}

func sliceContains(sl []string, target string) bool {
	for _, s := range sl {
		if s == target {
			return true
		}
	}
	return false
}
