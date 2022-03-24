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

// release.go fetches the releases on GitHub and compares them with the existing
// releases configured with Prow jobs.

package pkg

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	"github.com/google/go-github/v32/github"
	"k8s.io/apimachinery/pkg/util/sets"
	"knative.dev/test-infra/pkg/ghutil"
)

const (
	maxReleaseBranches      = 4
	releaseBranchNamePrefix = "release-"
)

// UpdateReleaseBranchConfig updates the config files for release branches.
func UpdateReleaseBranchConfig(gc ghutil.GithubOperations, configRootPath string) error {
	repoReleaseMap, err := collectRepoReleases(configRootPath)
	if err != nil {
		return fmt.Errorf("error collecting repo releases: %w", err)
	}

	if err := updateProwJobsForReleases(configRootPath, repoReleaseMap, gc); err != nil {
		return fmt.Errorf("error updating Prow jobs for releases: %w", err)
	}

	return nil
}

// collectRepoReleases iterates all the Prow job meta config files, and collects
// the release branch names configured for each repo.
func collectRepoReleases(configRootPath string) (map[string]sets.String, error) {
	orgRepoReleaseMap := map[string]sets.String{}
	if err := filepath.WalkDir(configRootPath, func(path string, d os.DirEntry, err error) error {
		// Skip directory and base config file.
		if d.IsDir() || d.Name() == ".base.yaml" {
			return nil
		}

		log.Printf("Reading jobs from %q", path)
		jobsConfig := mustReadJobsConfig(path)
		org := jobsConfig.Org
		repo := jobsConfig.Repo
		orgRepo := org + "/" + repo

		if _, ok := orgRepoReleaseMap[orgRepo]; !ok {
			orgRepoReleaseMap[orgRepo] = sets.NewString()
		}
		branchName := jobsConfig.Branches[0]
		if isReleaseBranch(branchName) {
			orgRepoReleaseMap[orgRepo].Insert(branchName)
		}

		return nil
	}); err != nil {
		return orgRepoReleaseMap, err
	}

	return orgRepoReleaseMap, nil
}

// updateProwJobsForReleases updates the Prow jobs config for the latest release for
// each repo, if needed.
func updateProwJobsForReleases(configRootPath string, orgRepoReleaseMap map[string]sets.String, gc ghutil.GithubOperations) error {
	for orgRepo, releaseSet := range orgRepoReleaseMap {
		org := strings.Split(orgRepo, "/")[0]
		repo := strings.Split(orgRepo, "/")[1]

		latest, err := latestReleaseBranch(gc, org, repo)
		if err != nil {
			return fmt.Errorf("error getting the latest release branch for %s: %w", orgRepo, err)
		}
		// Skip if there is no release branch.
		if latest == "" {
			continue
		}

		releases := releaseSet.List()
		sortReleases(releases)
		log.Printf("Existing releases for %s/%s: %v", org, repo, releases)
		log.Printf("Latest release for %s/%s: %s", org, repo, latest)
		// Skip if the latest release has already been configured.
		if len(releases) != 0 && releases[len(releases)-1] == latest {
			log.Printf("%s is already added for %s/%s:%v, skipping it", latest, org, repo, releases)
			continue
		}

		releaseToAdd := latest
		releaseToRemove := ""
		// If the number of releases is already maximum, remove the earliest one.
		if len(releases) == maxReleaseBranches {
			releaseToRemove = releases[0]
		}

		if err := syncProwJobsForRelease(configRootPath, org, repo, releaseToRemove, releaseToAdd); err != nil {
			return fmt.Errorf("error syncing Prow jobs for %s: %w", orgRepo, err)
		}
	}

	return nil
}

// latestReleaseBranch fetches the branch names for the give org/repo and
// returns the latest release branch name.
func latestReleaseBranch(gc ghutil.GithubOperations, org, repo string) (string, error) {
	branches, err := gc.ListBranches(org, repo)
	if err != nil {
		return "", fmt.Errorf("failed listing branches for '%s/%s': %w", org, repo, err)
	}

	// TODO: delete
	branchNames := []string{}
	for _, branch := range branches {
		branchNames = append(branchNames, *branch.Name)
	}
	log.Printf("######### All branches for %s/%s: %v", org, repo, branchNames)
	return filterLatest(branches), nil
}

// filterLatest returns latest release branch in the form of
// `release-[MAJOR].[MINOR]``, if there is no valid release branch exist in the form of
// `release-[MAJOR]-[MINOR]`, then it returns ""
func filterLatest(branches []*github.Branch) string {
	var latest string

	for _, branch := range branches {
		if isReleaseBranch(*branch.Name) {
			if latest == "" || versionComp(*branch.Name, latest) > 0 {
				latest = *branch.Name
			}
		}
	}

	return latest
}

func isReleaseBranch(branchName string) bool {
	return strings.HasPrefix(branchName, releaseBranchNamePrefix)
}

func sortReleases(releases []string) {
	sort.Slice(releases, func(i, j int) bool {
		return versionComp(releases[i], releases[j]) < 0
	})
}

func versionComp(release1, release2 string) int {
	v1 := strings.TrimPrefix(release1, releaseBranchNamePrefix)
	v2 := strings.TrimPrefix(release2, releaseBranchNamePrefix)
	leftMajor, leftMinor := majorMinor(v1)
	rightMajor, rightMinor := majorMinor(v2)
	if leftMajor != rightMajor {
		return leftMajor - rightMajor
	}
	if leftMinor != rightMinor {
		return leftMinor - rightMinor
	}
	return 0
}

func majorMinor(s string) (int, int) {
	parts := strings.Split(s, ".")
	if len(parts) != 2 {
		log.Fatalf("Version string has to be in the form of [MAJOR].[MINOR]: %q", s)
	}
	return mustInt(parts[0]), mustInt(parts[1])
}

func mustInt(s string) int {
	r, err := strconv.Atoi(s)
	if err != nil {
		log.Fatalf("Failed to parse int %q: %v", s, err)
	}
	return r
}
