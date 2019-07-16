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
	log.Printf("Job %s: fit all criteria - Starting analysis", jd.String())

	failedTests, err := jd.getFailedTests()
	if err != nil {
		log.Printf("Job %s: could not get failed tests: %v", jd.String(), err)
		return
	}
	if len(failedTests) == 0 {
		log.Printf("Job %s: no failed tests, skipping\n", jd.String())
		return
	}
	log.Printf("Job %s: got %d failed tests\n", jd.String(), len(failedTests))

	flakyTests, err := jd.getFlakyTests()
	if err != nil {
		log.Printf("Job %s: could not get flaky tests: %v", jd.String(), err)
		return
	}
	log.Printf("Job %s: got %d flaky tests from today's report\n", jd.String(), len(flakyTests))

	if outliers := getNonFlakyTests(failedTests, flakyTests); len(outliers) > 0 {
		log.Printf("Job %s: %d of %d failed tests are not flaky, cannot retry\n", jd.String(), len(outliers), len(failedTests))
		// TODO: Post GitHub comment describing why we cannot retry, listing the
		// non-flaky failed tests that the developer needs to fix. Logic will be in
		// github_commenter.go
		return
	}
	log.Printf("Job %s: all failed tests are flaky, triggering retry\n", jd.String())
	// TODO: Post GitHub comment stating as such, and trigger the job. Do not post
	// comment if we are out of retries. Logic will be in github_commenter.go
}
