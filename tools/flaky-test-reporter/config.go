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

// config.go contains configurations for flaky tests reporting

package main

import (
	"github.com/knative/test-infra/shared/prow"
)

const (
	// Builds to be analyzed, this is an arbitrary number
	buildsCount = 10
	// Minimal number of results to be counted as valid results for each testcase, this is an arbitrary number
	requiredCount = 8
	// Don't do anything if found more than 1% tests flaky, this is an arbitrary number
	threshold = 0.01

	org = "knative"
)

var (
	// jobConfigs lists all repos and jobs to analyze within those repos
	jobConfigs = map[string][]JobConfig{
		// CI flows for serving repo
		"serving": {{Name: "ci-knative-serving-continuous", Type: prow.PostsubmitJob, PostIssue: true},
			{Name: "ci-knative-serving-istio-1.0.7-mesh", Type: prow.PostsubmitJob, PostIssue: false},
			{Name: "ci-knative-serving-istio-1.0.7-no-mesh", Type: prow.PostsubmitJob, PostIssue: false},
			{Name: "ci-knative-serving-istio-1.1.2-mesh", Type: prow.PostsubmitJob, PostIssue: false},
			{Name: "ci-knative-serving-istio-1.1.2-no-mesh", Type: prow.PostsubmitJob, PostIssue: false},
		},
	}
	// slackChannelsMap lists which Slack channel to post results in for each job in repo
	slackChannelsMap = map[string]map[string][]slackChannel{
		// channel mapping for serving repo
		// "CA4DNJ9A4" => serving-api
		// "CA9RHBGJX" => networking
		"serving": {"ci-knative-serving-continuous": {{"api", "CA4DNJ9A4"}},
			"ci-knative-serving-istio-1.0.7-mesh":    {{"networking", "CA9RHBGJX"}},
			"ci-knative-serving-istio-1.0.7-no-mesh": {{"networking", "CA9RHBGJX"}},
			"ci-knative-serving-istio-1.1.2-mesh":    {{"networking", "CA9RHBGJX"}},
			"ci-knative-serving-istio-1.1.2-no-mesh": {{"networking", "CA9RHBGJX"}},
		},
	}

	githubIssueMap = map[string]string{
		"serving": "serving",
	}
)
