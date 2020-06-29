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
	"regexp"
	"strings"

	"github.com/google/go-github/v27/github"
	"knative.dev/test-infra/pkg/ghutil"
)

func latestReleaseBranch(gc ghutil.GithubOperations, repo string) (string, error) {
	parts := strings.Split(repo, "/")
	if len(parts) != 2 {
		return "", fmt.Errorf("repo name %q should be in the form of [ORG]/[REPO]", repo)
	}
	branches, err := gc.ListBranches(parts[0], parts[1])
	if err != nil {
		return "", fmt.Errorf("failed listing branches for repo %q: %w", repo, err)
	}
	return filterLatest(branches), nil
}

// filterLatest returns latest release branch in the form of
// [MAJOR].[MINOR], if there is no valid release branch exist in the form of
// `release-[MAJOR]-[MINOR]`, then it returns ""
func filterLatest(branches []*github.Branch) string {
	var (
		reReleaseBranch = regexp.MustCompile(`^release\-(\d+\.\d+)$`)
		latest          = ""
	)

	for _, branch := range branches {
		if matches := reReleaseBranch.FindStringSubmatch(*branch.Name); len(matches) > 1 {
			release := matches[1]
			if latest == "" || versionComp(release, latest) > 0 {
				latest = release
			}
		}
	}

	return latest
}
