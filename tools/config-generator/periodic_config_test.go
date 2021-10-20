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
	"fmt"
	"testing"

	"github.com/google/go-cmp/cmp"
	"gopkg.in/yaml.v2"
	"knative.dev/test-infra/tools/config-generator/unstructured"
)

func TestClone(t *testing.T) {
	SetupForTesting()
	base := baseProwJobTemplateData{OrgName: "org-name"}
	data := periodicJobTemplateData{
		Base:            base,
		PeriodicJobName: "periodic-job-name",
		CronString:      "cron-string",
		PeriodicCommand: []string{"string-a", "string-b"},
	}
	if diff := cmp.Diff(data.Clone(), data); diff != "" {
		t.Fatalf("Incorrect output for empty string: (-got +want)\n%s", diff)
	}
}

func TestGetUTCtime(t *testing.T) {
	SetupForTesting()
	for i := 0; i < 24; i++ {
		utcTime := getUTCtime(i)
		expected := (i + 7) % 24
		if utcTime != expected {
			t.Fatalf("Expected %d, got %d", expected, utcTime)
		}
	}
}

func TestCalculateMinuteOffset(t *testing.T) {
	SetupForTesting()
	out1 := calculateMinuteOffset("foo")
	out2 := calculateMinuteOffset("foo")
	if diff := cmp.Diff(out1, out2); diff != "" {
		t.Fatalf("Same input should always yield same offset")
	}
}

func TestGenerateCron(t *testing.T) {
	SetupForTesting()
	jobName := "job-name"
	tests := []struct {
		jobType  string
		repoName string
		timeout  int
		expected string
	}{
		{
			jobType:  "not-supported",
			expected: "",
		},
		{
			jobType:  "continuous",
			timeout:  54,
			expected: fmt.Sprintf("%d * * * *", calculateMinuteOffset("continuous", jobName)),
		},
		{
			jobType:  "continuous",
			timeout:  55,
			expected: fmt.Sprintf("%d */2 * * *", calculateMinuteOffset("continuous", jobName)),
		},
		{
			jobType:  "continuous",
			timeout:  60 + 55,
			expected: fmt.Sprintf("%d */3 * * *", calculateMinuteOffset("continuous", jobName)),
		},
		{
			jobType:  "custom-job",
			timeout:  54,
			expected: fmt.Sprintf("%d * * * *", calculateMinuteOffset("custom-job", jobName)),
		},
		{
			jobType:  "auto-release",
			timeout:  54,
			expected: fmt.Sprintf("%d * * * *", calculateMinuteOffset("auto-release", jobName)),
		},
		{
			jobType:  "branch-ci",
			expected: fmt.Sprintf("%d 8 * * *", calculateMinuteOffset("branch-ci", jobName)),
		},
		{
			jobType:  "nightly",
			expected: fmt.Sprintf("%d 9 * * *", calculateMinuteOffset("nightly", jobName)),
		},
		{
			jobType:  "dot-release",
			repoName: "foo",
			expected: fmt.Sprintf("%d 9 * * 2", calculateMinuteOffset("dot-release", jobName)),
		},
		{
			jobType:  "dot-release",
			repoName: "foo-operator",
			expected: fmt.Sprintf("%d 19 * * 2", calculateMinuteOffset("dot-release", jobName)),
		},
	}
	for _, tc := range tests {
		out := generateCron(tc.jobType, jobName, tc.repoName, tc.timeout)
		if diff := cmp.Diff(out, tc.expected); diff != "" {
			t.Fatalf("For jobType %v and timeout %d: (-got +want)\n%s", tc.jobType, tc.timeout, diff)
		}
	}
}

func TestGeneratePeriodic(t *testing.T) {
	title := "title"
	repoName := "repoName"
	tests := []struct {
		jobType    string
		assertions []unstructured.Assertion
	}{
		{jobType: "continuous"},
		{jobType: "nightly", assertions: []unstructured.Assertion{hasProperArgs(title, []string{
			"./hack/release.sh",
			"--publish", "--tag-release",
		})}},
		{jobType: "branch-ci"},
		{jobType: "dot-release", assertions: []unstructured.Assertion{hasProperArgs(title, []string{
			"./hack/release.sh",
			"--dot-release", "--release-gcs", repoName,
			"--release-gcr", "gcr.io/knative-releases",
			"--github-token", "/etc/hub-token/token",
		})}},
		{jobType: "auto-release", assertions: []unstructured.Assertion{hasProperArgs(title, []string{
			"./hack/release.sh",
			"--auto-release", "--release-gcs", repoName,
			"--release-gcr", "gcr.io/knative-releases",
			"--github-token", "/etc/hub-token/token",
		})}},
	}
	var periodicConfig yaml.MapSlice
	oldReleaseScript := releaseScript
	defer func() {
		releaseScript = oldReleaseScript
	}()
	for _, tc := range tests {
		tc := tc
		t.Run(tc.jobType, func(t *testing.T) {
			SetupForTesting()
			releaseScript = "./hack/release.sh"
			periodicConfig = yaml.MapSlice{{Key: tc.jobType, Value: true}}
			generatePeriodic(title, repoName, periodicConfig)
			out := GetOutput()
			outputLen := len(out)
			if outputLen == 0 {
				t.Fatal("No output")
			}
			if logFatalCalls != 0 {
				t.Fatal("LogFatal was called")
			}
			un := make(map[interface{}]interface{})
			err := yaml.Unmarshal([]byte(out), &un)
			if err != nil {
				t.Fatal(err)
			}
			for _, assertion := range tc.assertions {
				err = assertion(un)
				if err != nil {
					t.Fatal(err)
				}
			}
		})
	}
}

func TestGenerateGoCoveragePeriodic(t *testing.T) {
	SetupForTesting()
	repositories = []repositoryData{
		{
			Name:                "repo-name",
			EnableGoCoverage:    true,
			GoCoverageThreshold: 80,
		},
	}
	generateGoCoveragePeriodic("title", "repo-name", nil)
	if len(GetOutput()) == 0 {
		t.Fatalf("No output")
	}
	if logFatalCalls != 0 {
		t.Fatalf("LogFatal was called.")
	}
}

func hasProperArgs(title string, want []string) unstructured.Assertion {
	assert := unstructured.EqualsStringSlice(want)
	query := fmt.Sprintf("%s.0.spec.containers.0.args", title)
	return queryAndAssert(query, assert)
}

func queryAndAssert(query string, assert unstructured.Assertion) unstructured.Assertion {
	return func(un interface{}) error {
		questioner := unstructured.NewQuestioner(un)
		val, err := questioner.Query(query)
		if err != nil {
			return err
		}
		return assert(val)
	}
}
