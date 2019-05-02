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
	"fmt"
	"os"
	"os/signal"
	"time"

	"github.com/knative/test-infra/shared/prow"
)

type Client struct {
	StartDate time.Time         // Earliest date to be analyzed, i.e. "2019-02-22"
	Parser    func(string) bool // Parser function
	JobFilter []string          // Jobs to be parsed. If not provided will parse all jobs
	PrChan    chan prInfo
	JobChan   chan prow.Job
	BuildChan chan prow.Build
	LogChan   chan logInfo
	Dryrun    bool
}

func NewClient(serviceAccount string, parser func(string) bool) (*Client, error) {
	if err := prow.Initialize(serviceAccount); nil != err { // Explicit authenticate with gcs Client
		return nil, fmt.Errorf("Failed authenticating GCS: '%v'", err)
	}

	c := &Client{Parser: parser}
	c.refreshChans()

	return c, nil
}

func (c *Client) SetStartDate(startDate string) error {
	tt, err := time.Parse("2006-01-02", startDate)
	if nil != err {
		return fmt.Errorf("invalid start date string, expecing format YYYY-MM-DD: '%v'", err)
	}
	c.StartDate = tt
	return nil
}

func sliceContains(sl []string, target string) bool {
	for _, s := range sl {
		if s == target {
			return true
		}
	}
	return false
}

func (c *Client) refreshChans() {
	c.cleanup()
	c.PrChan = make(chan prInfo)
	c.JobChan = make(chan prow.Job, 50)
	c.BuildChan = make(chan prow.Build, 500)
	c.LogChan = make(chan logInfo, 500)
}

func (c *Client) cleanup() {
	if c.PrChan != nil {
		close(c.PrChan)
	}
	if c.JobChan != nil {
		close(c.JobChan)
	}
	if c.BuildChan != nil {
		close(c.BuildChan)
	}
	if c.LogChan != nil {
		close(c.LogChan)
	}
}

// CleanupOnInterrupt will execute the function cleanup if an interrupt signal is caught
func (c *Client) CleanupOnInterrupt() {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt)
	go func() {
		for range ch {
			c.cleanup()
			os.Exit(1)
		}
	}()
}
