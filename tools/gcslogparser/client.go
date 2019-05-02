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
	"time"

	"github.com/knative/test-infra/shared/prow"
)

type Client struct {
	StartDate time.Time         // Earliest date to be analyzed, i.e. "2019-02-22"
	Parser    func(string) bool // Parser function
	JobFilter []string          // Jobs to be parsed. If not provided will parse all jobs
	Dryrun    bool
}

func NewClient(serviceAccount string, parser func(string) bool) (*Client, error) {
	if err := prow.Initialize(serviceAccount); nil != err { // Explicit authenticate with gcs Client
		return nil, fmt.Errorf("Failed authenticating GCS: '%v'", err)
	}

	return &Client{Parser: parser}, nil
}

func (c *Client) SetStartDate(startDate string) error {
	tt, err := time.Parse("2006-01-02", startDate)
	if nil != err {
		return fmt.Errorf("invalid start date string, expecing format YYYY-MM-DD: '%v'", err)
	}
	c.StartDate = tt
}
