/*
Copyright 2022 The Kubernetes Authors.

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

package mage

import (
	"fmt"

	"github.com/carolynvs/magex/shx"
)

// getVersion gets a description of the commit, e.g. v0.30.1 (latest) or v0.30.1-32-gfe72ff73 (canary)
func getVersion() (string, error) {
	version, err := shx.Output("git", "describe", "--tags", "--match=v*")
	if err != nil {
		return "", err
	}
	if version != "" {
		return version, nil
	}

	// repo without any tags in it
	return "v0.0.0", nil
}

// getCommit gets the hash of the current commit
func getCommit() (string, error) {
	return shx.Output("git", "rev-parse", "--short", "HEAD")
}

// getGitState gets the state of the git repository
func getGitState() string {
	_, err := shx.Output("git", "diff", "--quiet")
	if err != nil {
		return "dirty"
	}

	return "clean"
}

// getBuildDateTime gets the build date and time
func getBuildDateTime() (string, error) {
	result, err := shx.Output("git", "log", "-1", "--pretty=%ct")
	if err != nil {
		return "", err
	}
	if result != "" {
		sourceDateEpoch := fmt.Sprintf("@%s", result)
		date, err := shx.Output("date", "-u", "-d", sourceDateEpoch, "+%Y-%m-%dT%H:%M:%SZ")
		if err != nil {
			return "", err
		}
		return date, nil
	}

	return shx.Output("date", "+%Y-%m-%dT%H:%M:%SZ")
}

// GenerateLDFlags returns the string to use in the `-ldflags` flag.
func GenerateLDFlags() (string, error) {
	pkg := "sigs.k8s.io/release-utils/version"
	version, err := getVersion()
	if err != nil {
		return "", err
	}
	commit, err := getCommit()
	if err != nil {
		return "", err
	}
	buildTime, err := getBuildDateTime()
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("-X %[1]s.gitVersion=%[2]s -X %[1]s.gitCommit=%[3]s -X %[1]s.gitTreeState=%[4]s -X %[1]s.buildDate=%[5]s",
		pkg, version, commit, getGitState(), buildTime), nil
}
