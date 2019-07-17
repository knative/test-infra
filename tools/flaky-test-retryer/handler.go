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

// handler.go contains most of the main logic for the flaky-test-retryer. Listen for
// incoming Pubsub messages, verify that the message we received is one we want to
// process, compare flaky and failed tests, and trigger retests if necessary.

package main

import (
	"context"
	"fmt"
	"log"

	"github.com/knative/test-infra/tools/monitoring/subscriber"
	// TODO: remove this import once "k8s.io/test-infra" import problems are fixed
	// https://github.com/test-infra/test-infra/issues/912
	"github.com/knative/test-infra/tools/monitoring/prowapi"
)

// HandlerClient wraps the other clients we need when processing failed jobs.
type HandlerClient struct {
	context.Context
	pubsub *subscriber.Client
	github *GithubClient
}

// NewHandlerClient gives us a handler where we can listen for Pubsub messages and
// post comments on GitHub.
func NewHandlerClient(githubAccount string) (*HandlerClient, error) {
	ctx := context.Background()
	githubClient, err := NewGithubClient(githubAccount)
	if err != nil {
		return nil, fmt.Errorf("Github client: %v", err)
	}
	pubsubClient, err := subscriber.NewSubscriberClient(ctx, projectName, pubsubTopic)
	if err != nil {
		return nil, fmt.Errorf("Pubsub client: %v", err)
	}
	return &HandlerClient{
		ctx,
		pubsubClient,
		githubClient,
	}, nil
}

// Listen scans for incoming Pubsub messages, spawning a new goroutine for each
// one that fits our criteria.
func (hc *HandlerClient) Listen() {
	log.Printf("Listening for failed jobs...\n")
	for {
		hc.pubsub.ReceiveMessageAckAll(hc, func(msg *prowapi.ReportMessage) {
			data := NewJobData(msg)
			if data.IsSupported() {
				go hc.HandleJob(data)
			}
		})
	}
}

// HandleMessage gets the job's failed tests and the current flaky tests,
// compares them, and triggers a retest if all the failed tests are flaky.
func (hc *HandlerClient) HandleJob(jd *JobData) {
	logWithPrefix(jd, "fit all criteria - Starting analysis\n")

	failedTests, err := jd.getFailedTests()
	if err != nil {
		logWithPrefix(jd, "could not get failed tests: %v", err)
		return
	}
	if len(failedTests) == 0 {
		logWithPrefix(jd, "no failed tests, skipping\n")
		return
	}
	logWithPrefix(jd, "got %d failed tests", len(failedTests))

	flakyTests, err := jd.getFlakyTests()
	if err != nil {
		logWithPrefix(jd, "could not get flaky tests: %v", err)
		return
	}
	logWithPrefix(jd, "got %d flaky tests from today's report\n", len(flakyTests))

	if outliers := getNonFlakyTests(failedTests, flakyTests); len(outliers) > 0 {
		logWithPrefix(jd, "%d of %d failed tests are not flaky, cannot retry\n", len(outliers), len(failedTests))
		// TODO: Post GitHub comment describing why we cannot retry, listing the
		// non-flaky failed tests that the developer needs to fix. Logic will be in
		// github_commenter.go
		return
	}
	logWithPrefix(jd, "all failed tests are flaky, triggering retry\n")
	// TODO: Post GitHub comment stating as such, and trigger the job. Do not post
	// comment if we are out of retries. Logic will be in github_commenter.go
}

// logWithPrefix wraps a call to log.Printf, prefixing the arguments with details
// about the job passed in.
func logWithPrefix(jd *JobData, format string, a ...interface{}) {
	input := append([]interface{}{jd.Refs[0].Repo, jd.Refs[0].Pulls[0].Number, jd.JobName}, a...)
	log.Printf("%s/pull/%d: %s: "+format, input...)
}
