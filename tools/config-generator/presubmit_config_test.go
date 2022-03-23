/*
Copyright 2020 The Knative Authors

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
	"testing"

	"gopkg.in/yaml.v2"
)

func TestGeneratePresubmit(t *testing.T) {
	SetupForTesting()
	title := "title"
	repoName := "repoName"
	items := []yaml.MapItem{
		{Key: "build-tests", Value: true},
		{Key: "unit-tests", Value: true},
		{Key: "integration-tests", Value: true},
		{Key: "go-coverage", Value: true},
		{Key: "custom-test", Value: "foo"},
		{Key: "go-coverage-threshold", Value: 80},
		{Key: "run-if-changed", Value: "foo"},
	}
	var presubmitConfig yaml.MapSlice
	for _, item := range items {
		presubmitConfig = yaml.MapSlice{item}
		generatePresubmit(title, repoName, presubmitConfig)
		outputLen := len(GetOutput())
		if outputLen == 0 {
			t.Errorf("Failure for key %s: No output", item.Key)
		}
		if logFatalCalls != 0 {
			t.Errorf("Failure for key %s: LogFatal was called.", item.Key)
		}
		SetupForTesting()
	}
}
