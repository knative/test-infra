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
	"github.com/google/go-github/v27/github"
	"knative.dev/test-infra/pkg/ghutil/fakeghutil"
)

func TestLatestReleaseBranch(t *testing.T) {
	SetupForTesting()
	fgc := fakeghutil.NewFakeGithubClient()

	names := []string{
		"release-0.1",
		"release-1.0",
		"release-3.4",
		"release-0.2",
	}

	branches := []*github.Branch{}
	for i := range names {
		branches = append(branches, &github.Branch{Name: &names[i]})
	}

	fgc.Branches = map[string][]*github.Branch{
		"my-repo": branches,
	}

	_, err := latestReleaseBranch(fgc, "no slash")
	if err == nil {
		t.Errorf("Format was not ORG/REPO, expected error.")
	}
	latest, _ := latestReleaseBranch(fgc, "my-org/my-repo")
	if diff := cmp.Diff(latest, "3.4"); diff != "" {
		t.Errorf("Did not find latest version (-got +want)\n%s", diff)
	}
}

func TestFilterLatest(t *testing.T) {
	SetupForTesting()
	names := []string{
		"release-0.1",
		"release-1.0",
		"release-3.4",
		"release-0.2",
	}

	branches := []*github.Branch{}
	for i := range names {
		branches = append(branches, &github.Branch{Name: &names[i]})
	}

	res := filterLatest(branches)
	if diff := cmp.Diff(res, "3.4"); diff != "" {
		t.Errorf("Did not find latest version (-got +want)\n%s", diff)
	}
}
