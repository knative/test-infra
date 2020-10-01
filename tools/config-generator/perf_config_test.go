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

	"github.com/google/go-cmp/cmp"
)

func TestGeneratePerfClusterUpdatePeriodicJobs(t *testing.T) {
	SetupForTesting()
	repositories = []repositoryData{
		{
			Name:                   "enabled-repo",
			EnablePerformanceTests: true,
		},
	}
	generatePerfClusterUpdatePeriodicJobs()
	if logFatalCalls != 0 || len(GetOutput()) == 0 {
		t.Errorf("Expected job to be written without errors")
	}

	SetupForTesting()
	repositories = []repositoryData{
		{
			Name:                   "disabled-repo",
			EnablePerformanceTests: false,
		},
	}
	generatePerfClusterUpdatePeriodicJobs()
	if len(GetOutput()) != 0 {
		t.Errorf("Expected nothing to be written")
	}
}

func TestGeneratePerfClusterPostsubmitJob(t *testing.T) {
	SetupForTesting()
	generatePerfClusterPostsubmitJob(repositoryData{Name: "my-repo"})
	if logFatalCalls != 0 || len(GetOutput()) == 0 {
		t.Errorf("Expected job to be written without errors")
	}
}

func TestPerfClusterPeriodicJob(t *testing.T) {
	SetupForTesting()
	repoData := repositoryData{Name: "my-repo"}
	perfClusterPeriodicJob("postfix", "cronString", "command", []string{"arg1", "arg2"}, repoData, "sa")

	if logFatalCalls != 0 || len(GetOutput()) == 0 {
		t.Errorf("Expected job to be written without errors")
	}
}

func TestPerfClusterReconcilePostsubmitJob(t *testing.T) {
	SetupForTesting()
	repoData := repositoryData{Name: "my-repo"}
	perfClusterReconcilePostsubmitJob("postfix", "command", []string{"arg1", "arg2"}, repoData, "sa")

	if logFatalCalls != 0 || len(GetOutput()) == 0 {
		t.Errorf("Expected job to be written without errors")
	}
}

func TestPerfClusterBaseProwJob(t *testing.T) {
	SetupForTesting()
	command := "command"
	args := []string{"arg1", "arg2"}
	repoName := "org-name/repo-name"
	sa := "foo"
	res := perfClusterBaseProwJob(command, args, repoName, sa)

	if diff := cmp.Diff(res.Command, command); diff != "" {
		t.Errorf("Incorrect command: (-got +want)\n%s", diff)
	}
	if diff := cmp.Diff(res.Args, args); diff != "" {
		t.Errorf("Incorrect args: (-got +want)\n%s", diff)
	}
	if diff := cmp.Diff(res.Command, command); diff != "" {
		t.Errorf("Incorrect command: (-got +want)\n%s", diff)
	}
	if len(res.Env) != 8 {
		t.Errorf("Expected 8 environments, got %d", len(res.Env))
	}
}
