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

func NewHandlerClient(githubClient *GithubClient) (*HandlerClient, error) {
	ctx := context.Background()
	pubsubClient, err := subscriber.NewSubscriberClient(ctx, projectName, pubsubTopic)
	return &HandlerClient{
		ctx,
		pubsubClient,
		githubClient,
	}, err
}

// expectedMsg checks that the message we received is one we want to process.
func expectedMsg(msg *prowapi.ReportMessage) bool {
	repos, err := getReportRepos()
	if err != nil {
		log.Printf("Failed getting reporter's repos: %v", err)
		return false
	}
	expRepo := false
	if len(msg.Refs) > 0 {
		for _, repo := range repos {
			if msg.Refs[0].Repo == repo {
				expRepo = true
				break
			}
		}
	}
	return expRepo && msg.Status == prowapi.FailureState && msg.JobType == prowapi.PresubmitJob
}

// Listen scans for incoming Pubsub messages, spawning a new goroutine for each
// one that fits our criteria.
func (hc *HandlerClient) Listen() {
	log.Printf("Listening for failed jobs...\n")
	for {
		hc.pubsub.ReceiveMessageAckAll(hc, func(msg *prowapi.ReportMessage) {
			if expectedMsg(msg) {
				go hc.HandleMessage(msg)
			} else {
				log.Println("Job did not fit criteria - skipping")
			}
		})
	}
}

// HandleMessage gets the job's failed tests and the current flaky tests,
// compares them, and triggers a retest if all the failed tests are flaky.
func (hc *HandlerClient) HandleMessage(msg *prowapi.ReportMessage) {
	log.Println("Job fit all criteria - Starting analysis")
	// TODO: ensure job failed due to tests, get job's failed tests, get current
	//       flaky tests, cross-reference failed tests and flaky tests, retry tests
	//       if all failed tests are flaky.
}
