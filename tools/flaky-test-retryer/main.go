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

// flaky-test-retryer detects failed integration jobs on new pull requests,
// determines if they failed due to flaky tests, posts comments describing the
// issue, and retries them until they succeed.

package main

import (
	"context"
	"flag"
	"log"

	"github.com/knative/test-infra/shared/ghutil"
	"github.com/knative/test-infra/tools/monitoring/subscriber"
)

const (
	projectName = "knative-tests"
	pubsubTopic = "test-infra-monitoring-sub"
)

func main() {
	githubAccount := flag.String("github-account", "", "Token file for Github authentication")
	flag.Parse()

	ctx := context.Background()

	// will be used later
	_, err := ghutil.NewGithubClient(*githubAccount)
	if nil != err {
		log.Fatalf("Cannot setup GitHub client: '%v'", err)
	}
	pubsubClient, err := subscriber.NewSubscriberClient(ctx, projectName, pubsubTopic)
	if nil != err {
		log.Fatalf("Coud not create Pub/Sub client: '%v'", err)
	}

	for {
		pubsubClient.ReceiveMessageAckAll(ctx, handleDriver)
	}
}
