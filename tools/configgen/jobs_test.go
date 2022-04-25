/*
Copyright 2022 The Knative Authors

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

// jobs_test.go runs basic validations for the meta Prow job config files.

package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"istio.io/test-infra/tools/prowgen/pkg/spec"
	"k8s.io/apimachinery/pkg/util/sets"
	"sigs.k8s.io/yaml"
)

const jobsConfigPath = "../../prow/jobs_config"

func TestOrgRepo(t *testing.T) {
	var errStrs strings.Builder
	if err := filepath.WalkDir(jobsConfigPath, func(path string, d os.DirEntry, err error) error {
		t.Logf("Validating org and repo for %q", path)
		// Skip directory, base config file and other unrelated files.
		if d.IsDir() || d.Name() == ".base.yaml" || !strings.HasSuffix(path, ".yaml") {
			return nil
		}

		jobs := mustReadJobsConfig(t, path)

		org := jobs.Org
		parentDir := filepath.Base(filepath.Dir(path))
		if parentDir != org {
			errStrs.WriteString(fmt.Sprintf("Config file %q must be under %q folder.\n", path, org))
		}

		if len(jobs.Branches) != 1 {
			errStrs.WriteString(fmt.Sprintf("Config file %q must only have one branch configured but got %v.\n", path, jobs.Branches))
		}

		repo := jobs.Repo
		branch := jobs.Branches[0]
		repoBranch := repo
		if branch != "main" {
			repoBranch = repo + "-" + branch
		}
		if strings.TrimSuffix(d.Name(), ".yaml") != repoBranch {
			errStrs.WriteString(fmt.Sprintf("Config file %q must be named as %q.\n", path, repoBranch+".yaml"))
		}

		return nil
	}); err != nil {
		t.Fatalf("Error walking dir %q: %v", jobsConfigPath, err)
	}

	if errStrs.Len() != 0 {
		t.Fatalf("Error validating org and repo:\n%s", errStrs.String())
	}
}

func TestReleaseJobs(t *testing.T) {
	var errStrs strings.Builder
	if err := filepath.WalkDir(jobsConfigPath, func(path string, d os.DirEntry, err error) error {
		t.Logf("Validating release jobs for %q", path)
		// Skip directory, base config file and other unrelated files.
		if d.IsDir() || d.Name() == ".base.yaml" || !strings.HasSuffix(path, ".yaml") {
			return nil
		}

		jobs := mustReadJobsConfig(t, path)

		for _, job := range jobs.Jobs {
			if job.Name == "nightly" {
				if job.Interval != "" || job.Cron != "" {
					errStrs.WriteString(fmt.Sprintf("cron is supposed to be auto-generated, do not add for nightly job in %q\n", path))
				}
				reqs := sets.NewString(job.Requirements...)
				if !reqs.Has("nightly") {
					errStrs.WriteString(fmt.Sprintf("nightly requirement is required for nightly job in %q\n", path))
				}

				excludedReqs := sets.NewString(job.ExcludedRequirements...)
				if !excludedReqs.Has("gcp") {
					errStrs.WriteString(fmt.Sprintf("gcp requirement cannot be set for nightly job in %q\n", path))
				}
			}

			if job.Name == "release" {
				if job.Interval != "" || job.Cron != "" {
					errStrs.WriteString(fmt.Sprintf("cron is supposed to be auto-generated, do not add it for release job in %q\n", path))
				}
				reqs := sets.NewString(job.Requirements...)
				if !reqs.Has("release") || !reqs.Has("release-dev") {
					errStrs.WriteString(fmt.Sprintf("release requirement is required for release job in %q\n", path))
				}

				excludedReqs := sets.NewString(job.ExcludedRequirements...)
				if !excludedReqs.Has("gcp") {
					errStrs.WriteString(fmt.Sprintf("gcp requirement cannot be set for release job in %q\n", path))
				}

				commandArgs := append(job.Command, job.Args...)
				for i, arg := range commandArgs {
					if arg == "--release-gcs" {
						if commandArgs[i+1] != "knative-releases/"+jobs.Repo {
							errStrs.WriteString(fmt.Sprintf("--release-gcs must be set to %q for release job in %q\n", "knative-releases/"+jobs.Repo, path))
						}
					}
				}
			}
		}

		return nil
	}); err != nil {
		t.Fatalf("Error walking dir %q: %v", jobsConfigPath, err)
	}

	if errStrs.Len() != 0 {
		t.Fatalf("Error validating release jobs:\n%s", errStrs.String())
	}
}

func mustReadJobsConfig(t *testing.T, file string) spec.JobsConfig {
	t.Helper()
	yamlFile, err := ioutil.ReadFile(file)
	if err != nil {
		t.Fatalf("Failed to read %q: %v", file, err)
	}
	jobsConfig := spec.JobsConfig{}
	if err := yaml.Unmarshal(yamlFile, &jobsConfig); err != nil {
		t.Fatalf("Failed to unmarshal %q: %v", file, err)
	}

	return jobsConfig
}
