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
	buildsCount   = 10
	// Minimal number of results to be counted as valid results for each testcase, this is an arbitrary number
	requiredCount = 8
	// Don't do anything if found more than 1% tests flaky, this is an arbitrary number
	threshold     = 0.01

	org           = "knative"
)

var (
	jobConfigs = []JobConfig{
		{"ci-knative-serving-continuous", "serving", prow.PostsubmitJob}, // CI flow for serving repo
	}
	// Temporarily creating issues under "test-infra" for better management
	// TODO(chaodaiG): repo for issue same as the src of the test
	repoIssueMap = map[string]string{
		"serving": "test-infra",
	}
)