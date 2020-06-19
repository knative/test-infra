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
	"io/ioutil"
	"log"
	"os"
	"sort"
	"strconv"
	"strings"

	"gopkg.in/yaml.v2"
	"knative.dev/pkg/test/ghutil"
)

const (
	maxReleaseBranches = 4
)

func upgradeReleaseBranchesTemplate(configfileName string, gc ghutil.GithubOperations) error {
	config := yaml.MapSlice{}
	info, err := os.Lstat(configfileName)
	if err != nil {
		return fmt.Errorf("failed stats file %q: %w", configfileName, err)
	}
	content, err := ioutil.ReadFile(configfileName)
	if err != nil {
		return fmt.Errorf("cannot read file %q: %w", configfileName, err)
	}
	if err = yaml.Unmarshal(content, &config); err != nil {
		return fmt.Errorf("cannot parse config %q: %w", configfileName, err)
	}
	for i, repos := range config {
		if repos.Key != "presubmits" {
			config[i].Value, err = getReposMap(gc, repos.Value)
			if err != nil {
				return err
			}
		}
	}

	updated, err := yaml.Marshal(&config)
	// This shouldn't happen, just catch it in case
	if err != nil {
		return fmt.Errorf("failed marshal modified content: %w", err)
	}
	return ioutil.WriteFile(configfileName, updated, info.Mode())
}

func getReposMap(gc ghutil.GithubOperations, val interface{}) (interface{}, error) {
	reposMap := getMapSlice(val)
	for j, repo := range reposMap {
		var (
			ciBranches        []string
			releaseBranches   []string
			skipCiUpdate      bool
			skipReleaseUpdate bool
		)
		repoName := getString(repo.Key)
		latest, err := latestReleaseBranch(gc, repoName)
		if err != nil {
			return nil, fmt.Errorf("failed getting latest release branches: %w", err)
		}
		if latest == "" {
			continue
		}

		log.Printf("Latest branch for repo %q is %q", repoName, latest)

		repoConfigs := getInterfaceArray(repo.Value)
		for _, repoConfig := range repoConfigs {
			jobConfig := getMapSlice(repoConfig)
			ciBranch, releaseBranch := getBranch(jobConfig)
			if ciBranch != "" {
				ciBranches = append(ciBranches, ciBranch)
				if ciBranch == latest {
					skipCiUpdate = true
				}
			}
			if releaseBranch != "" {
				releaseBranches = append(releaseBranches, releaseBranch)
				if releaseBranch == latest {
					skipReleaseUpdate = true
				}
			}
		}

		if !skipCiUpdate && len(ciBranches) > 0 {
			repoConfigs = updateConfigForJob(repoConfigs, ciBranches, latest,
				func(jobConfig yaml.MapSlice) string {
					branch, _ := getBranch(jobConfig)
					return branch
				})
		}

		if !skipReleaseUpdate && len(releaseBranches) > 0 {
			repoConfigs = updateConfigForJob(repoConfigs, releaseBranches, latest,
				func(jobConfig yaml.MapSlice) string {
					_, branch := getBranch(jobConfig)
					return branch
				})
		}

		reposMap[j].Value = repoConfigs
	}
	return reposMap, nil
}

func updateConfigForJob(repoConfigs []interface{}, branches []string, latest string,
	getBranchForJob func(yaml.MapSlice) string) []interface{} {

	var oldestBranchToSupport = "0.0"
	sortFunc(branches)
	if len(branches) >= maxReleaseBranches-1 {
		oldestBranchToSupport = branches[maxReleaseBranches-2]
	}
	var updatedRepoConfigs []interface{}
	for _, repoConfig := range repoConfigs {
		jobConfig := getMapSlice(repoConfig)
		branch := getBranchForJob(jobConfig)
		if branch == "" {
			updatedRepoConfigs = append(updatedRepoConfigs, jobConfig)
			continue
		}
		if versionComp(branch, oldestBranchToSupport) < 0 {
			log.Printf("Skipping %q for %q", branch, oldestBranchToSupport)
			continue
		}
		updatedRepoConfigs = append(updatedRepoConfigs, jobConfig)
		if branch == branches[0] {
			var next yaml.MapSlice
			for _, item := range jobConfig {
				val := item.Value
				if item.Key == "release" {
					val = latest
				}
				next = append(next, yaml.MapItem{Key: item.Key, Value: val})
			}
			updatedRepoConfigs = append(updatedRepoConfigs, next)
		}
	}

	return updatedRepoConfigs
}

func getBranch(jobConfig yaml.MapSlice) (ciBranch string, releaseBranch string) {
	var (
		branch     string
		isBranchCi bool
		isRelease  bool
	)
	for _, item := range jobConfig {
		switch item.Key {
		case "branch-ci":
			isBranchCi = true
		case "dot-release":
			isRelease = true
		case "release":
			branch = getString(item.Value)
		}
	}
	if branch == "" {
		return
	}
	if isBranchCi {
		ciBranch = branch
	} else if isRelease {
		releaseBranch = branch
	}

	return
}

func sortFunc(strSlice []string) {
	sort.Slice(strSlice, func(i, j int) bool {
		return versionComp(strSlice[i], strSlice[j]) > 0
	})
}

func versionComp(v1, v2 string) int {
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

func mustInt(s string) int {
	r, err := strconv.Atoi(s)
	if err != nil {
		log.Fatalf("Failed to parse int %q: %v", s, err)
	}
	return r
}

func majorMinor(s string) (int, int) {
	parts := strings.Split(s, ".")
	if len(parts) != 2 {
		log.Fatalf("Version string has to be in the form of [MAJOR].[MINOR]: %q", s)
	}
	return mustInt(parts[0]), mustInt(parts[1])
}
