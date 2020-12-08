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
)

func TestNewBaseTestgridTemplateData(t *testing.T) {
	SetupForTesting()
	data := newBaseTestgridTemplateData("foo")
	if diff := cmp.Diff(data.TestGroupName, "foo"); diff != "" {
		t.Errorf("(-got +want)\n%s", diff)
	}
}

func TestTestGridMetaDataGet(t *testing.T) {
	SetupForTesting()
	data := NewTestGridMetaData()
	jobDetails := data.Get("foo")
	if diff := cmp.Diff(jobDetails, data.md["foo"]); diff != "" {
		t.Errorf("(-got +want\n%s", diff)
	}
}

func TestTestGridMetaDataEnsureExists(t *testing.T) {
	SetupForTesting()
	data := NewTestGridMetaData()
	out := data.EnsureExists("foo")
	if out {
		t.Errorf("foo did not exist but function returned true")
	}
	if _, exists := data.md["foo"]; !exists {
		t.Errorf("foo should have been added but was not")
	}
	out = data.EnsureExists("foo")
	if !out {
		t.Errorf("foo existed but the function returned false")
	}
	if diff := cmp.Diff(data.projNames, []string{"foo"}); diff != "" {
		t.Errorf("(-got +want\n%s", diff)
	}
}

func TestTestGridMetaDataEnsureRepo(t *testing.T) {
	SetupForTesting()
	data := NewTestGridMetaData()
	out := data.EnsureRepo("proj-name", "repo-name")
	if out {
		t.Errorf("repo did not exist but function returned true")
	}
	if data.repoNames[0] != "repo-name" {
		t.Errorf("Should have added repo-name but did not")
	}
	out = data.EnsureRepo("proj-name", "repo-name")
	if !out {
		t.Errorf("repo existed but function returned false")
	}
}

func TestTestGridMetaDataGenerateTestGridSection(t *testing.T) {
	SetupForTesting()
	data := NewTestGridMetaData()
	data.projNames = []string{"project-a", "project-b"}
	data.repoNames = []string{"repo-1", "repo-2", "repo-3"}
	data.md["project-a"] = JobDetailMap{
		"repo-1": []string{"job-1a", "job-1b"},
		"repo-2": []string{"job-2a", "job-2b"},
	}
	data.md["project-b"] = JobDetailMap{
		"repo-3": []string{"job-3a", "job-3b"},
	}
	skipReleasedProj := false
	outputs := []string{}
	generator := func(proj, repo string, jobs []string) {
		outputs = append(outputs, fmt.Sprintf("%s %s %v", proj, repo, jobs))
	}
	data.generateTestGridSection("section-name", generator, skipReleasedProj)
	expected := []string{
		"project-a repo-1 [job-1a job-1b]",
		"project-a repo-2 [job-2a job-2b]",
		"project-b repo-3 [job-3a job-3b]",
	}
	if diff := cmp.Diff(outputs, expected); diff != "" {
		t.Errorf("(-got +want): \n%s", diff)
	}
}

func TestTestGridMetaDataGenerateNonAlignedTestGroups(t *testing.T) {
	SetupForTesting()
	data := NewTestGridMetaData()
	data.nonAligned = []NonAlignedTestGroup{
		{
			CIJobName: "ci-job-name",
			Extra:     map[string]string{},
		},
	}
	data.generateNonAlignedTestGroups()
	if len(GetOutput()) == 0 {
		t.Errorf("No output")
	}
	if logFatalCalls != 0 {
		t.Errorf("LogFatal was called.")
	}
}

func TestTestGridMetaDataAddNonAlignedTest(t *testing.T) {
	SetupForTesting()
	data := NewTestGridMetaData()
	data.AddNonAlignedTest(NonAlignedTestGroup{})
	if len(data.nonAligned) != 1 {
		t.Errorf("Test was not appended.")
	}
}

func TestGetGcsLogDir(t *testing.T) {
	SetupForTesting()
	GCSBucket = "gcs-bucket"
	LogsDir = "logs-dir"
	expected := "gcs-bucket/logs-dir/tg-name"
	if diff := cmp.Diff(getGcsLogDir("tg-name"), expected); diff != "" {
		t.Errorf("(-got +want): \n%s", diff)
	}
}

func TestGetTestgroupExtras(t *testing.T) {
	SetupForTesting()
	defaultProjectName := "project-name"
	tests := []struct {
		ProjName string
		JobName  string
		Expected map[string]string
	}{
		{
			ProjName: "proj-name-1.2.3",
			JobName:  "continuous",
			Expected: map[string]string{
				"num_failures_to_alert": "3",
				"alert_options":         "\n    alert_mail_to_addresses: \"serverless-engprod-sea@google.com\"",
			},
		},
		{
			JobName: "continuous",
			Expected: map[string]string{
				"alert_stale_results_hours": "3",
			},
		},
		{
			JobName: "dot-release",
			Expected: map[string]string{
				"num_failures_to_alert":     "1",
				"alert_options":             "\n    alert_mail_to_addresses: \"serverless-engprod-sea@google.com\"",
				"alert_stale_results_hours": "170",
			},
		},
		{
			JobName: "auto-release",
			Expected: map[string]string{
				"num_failures_to_alert": "1",
				"alert_options":         "\n    alert_mail_to_addresses: \"serverless-engprod-sea@google.com\"",
			},
		},
		{
			JobName: "nightly",
			Expected: map[string]string{
				"num_failures_to_alert": "1",
				"alert_options":         "\n    alert_mail_to_addresses: \"serverless-engprod-sea@google.com\"",
			},
		},
		{
			JobName: "test-coverage",
			Expected: map[string]string{
				"short_text_metric": "coverage",
			},
		},
		{
			JobName:  "some-other-job-name",
			Expected: map[string]string{"alert_stale_results_hours": "3"},
		},
	}

	for _, test := range tests {
		projName := test.ProjName
		if projName == "" {
			projName = defaultProjectName
		}

		out := getTestgroupExtras(test.ProjName, test.JobName)
		if diff := cmp.Diff(out, test.Expected); diff != "" {
			t.Errorf("(-got +want): \n%s", diff)
		}
	}
}

func TestGenerateProwJobAnnotations(t *testing.T) {
	SetupForTesting()
	tgExtras := map[string]string{
		"alert_stale_results_hours": "48",
		"alert_options":             "\n    alert_mail_to_addresses: \"foo-bar@google.com\"",
		"num_failures_to_alert":     "3",
		"short_text_metric":         "coverage",
	}
	expected := []string{
		"  testgrid-dashboards: repo-name",
		"  testgrid-tab-name: job-name",
		"  testgrid-alert-stale-results-hours: \"48\"",
		"  testgrid-in-cell-metric: coverage",
		"  testgrid-alert-email: \"foo-bar@google.com\"",
		"  testgrid-num-failures-to-alert: \"3\"",
	}
	annotations := generateProwJobAnnotations("repo-name", "job-name", tgExtras)
	if diff := cmp.Diff(annotations, expected); diff != "" {
		t.Errorf("(-got +want): \n%s", diff)
	}
}

func TestFmtDashboardAnnotation(t *testing.T) {
	if diff := cmp.Diff(fmtDashboardAnnotation("dashboardName"), "  testgrid-dashboards: dashboardName"); diff != "" {
		t.Errorf("(-got +want): \n%s", diff)
	}
}

func TestFmtTabAnnotation(t *testing.T) {
	if diff := cmp.Diff(fmtTabAnnotation("tabName"), "  testgrid-tab-name: tabName"); diff != "" {
		t.Errorf("(-got +want): \n%s", diff)
	}
}

func TestTestGridMetaDataGenerateTestGroup(t *testing.T) {
	SetupForTesting()
	data := NewTestGridMetaData()
	projName := "proj-name"
	repoName := "repo-name"
	jobNames := []string{"continuous", "dot-release", "webhook-api-coverage", "test-coverage", "default"}
	data.generateTestGroup(projName, repoName, jobNames)
	if len(GetOutput()) == 0 {
		t.Errorf("No output")
	}
	if logFatalCalls != 0 {
		t.Errorf("LogFatal was called.")
	}
}

func TestExecuteTestGroupTemplate(t *testing.T) {
	SetupForTesting()
	executeTestGroupTemplate("tg-name", "gcs-log-dir", map[string]string{})
	if len(GetOutput()) == 0 {
		t.Errorf("No output")
	}
	if logFatalCalls != 0 {
		t.Errorf("LogFatal was called.")
	}
}

func TestGenerateDashboard(t *testing.T) {
	SetupForTesting()
	projName := "proj-name"
	repoName := "repo-name"
	jobNames := []string{"continuous", "dot-release", "webhook-api-coverage", "nightly", "test-coverage", "default"}
	generateDashboard(projName, repoName, jobNames)
	if len(GetOutput()) == 0 {
		t.Errorf("No output")
	}
	if logFatalCalls != 0 {
		t.Errorf("LogFatal was called.")
	}
}

func TestExecuteDashboardTabTemplate(t *testing.T) {
	SetupForTesting()
	executeDashboardTabTemplate("tab-name", "tg-name", "base-opts", map[string]string{})
	if len(GetOutput()) == 0 {
		t.Errorf("No output")
	}
	if logFatalCalls != 0 {
		t.Errorf("LogFatal was called.")
	}
}

func TestGetTestGroupName(t *testing.T) {
	SetupForTesting()
	out := getTestGroupName("foo", "bar")
	expected := "ci-foo-bar"
	if diff := cmp.Diff(out, expected); diff != "" {
		t.Errorf("(-got +want): \n%s", diff)
	}

	out = getTestGroupName("foo", "nightly")
	expected = "ci-foo-nightly-release"
	if diff := cmp.Diff(out, expected); diff != "" {
		t.Errorf("(-got +want): \n%s", diff)
	}
}

func TestGenerateNonAlignedDashboards(t *testing.T) {
	SetupForTesting()
	data := NewTestGridMetaData()
	data.AddNonAlignedTest(NonAlignedTestGroup{
		DashboardName: "dashboard-name",
		HumanTabName:  "human-tab-name",
		CIJobName:     "ci-job-name",
		BaseOptions:   "base-opts",
	})
	data.generateNonAlignedDashboards()
	if len(GetOutput()) == 0 {
		t.Errorf("No output")
	}
	if logFatalCalls != 0 {
		t.Errorf("LogFatal was called.")
	}
}

func TestGenerateDashboardsForReleases(t *testing.T) {
	SetupForTesting()
	data := NewTestGridMetaData()
	data.projNames = []string{"project-a", "project-b-2.0"}
	data.repoNames = []string{"repo-1", "repo-2", "repo-3"}
	data.md["project-a"] = JobDetailMap{
		"repo-1": []string{"job-1a", "job-1b"},
		"repo-2": []string{"job-2a", "job-2b"},
	}
	data.md["project-b"] = JobDetailMap{
		"repo-3": []string{"job-3a", "job-3b"},
	}
	data.generateDashboardsForReleases()
	if len(GetOutput()) == 0 {
		t.Errorf("No output")
	}
	if logFatalCalls != 0 {
		t.Errorf("LogFatal was called.")
	}
}

func TestGenerateNonAlignedDashboardGroups(t *testing.T) {
	SetupForTesting()
	data := NewTestGridMetaData()
	data.nonAligned = []NonAlignedTestGroup{
		{
			DashboardName:  "dashboard-name",
			DashboardGroup: "dashboard-group",
		},
	}
	data.generateNonAlignedDashboardGroups()
	if len(GetOutput()) == 0 {
		t.Errorf("No output")
	}
	if logFatalCalls != 0 {
		t.Errorf("LogFatal was called.")
	}
}

func TestGenerateDashboardGroups(t *testing.T) {
	SetupForTesting()
	data := NewTestGridMetaData()
	data.projNames = []string{"project-a", "project-b-2.0"}
	data.repoNames = []string{"repo-1", "repo-2", "repo-3"}
	data.md["project-a"] = JobDetailMap{
		"repo-1": []string{"job-1a", "job-1b"},
		"repo-2": []string{"job-2a", "job-2b"},
	}
	data.md["project-b"] = JobDetailMap{
		"repo-3": []string{"job-3a", "job-3b"},
	}
	data.generateDashboardGroups()
	if len(GetOutput()) == 0 {
		t.Errorf("No output")
	}
	if logFatalCalls != 0 {
		t.Errorf("LogFatal was called.")
	}
}

func TestExecuteDashboardGroupTemplate(t *testing.T) {
	SetupForTesting()
	executeDashboardGroupTemplate("group-name", []string{"repo1", "repo2"})
	if len(GetOutput()) == 0 {
		t.Errorf("No output")
	}
	if logFatalCalls != 0 {
		t.Errorf("LogFatal was called.")
	}
}
