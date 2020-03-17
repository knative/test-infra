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
	"flag"
	"log"
	"strings"

	"knative.dev/pkg/test/ghutil"
	"knative.dev/pkg/test/prow"

	"knative.dev/test-infra/tools/prow-config-updater/config"
)

func main() {
	githubTokenPath := flag.String("github-token", "", "Github token file path for authenticating with Github")
	flag.Parse()

	ec, err := prow.GetEnvConfig()
	if err != nil {
		log.Fatalf("Error getting environment variables for Prow: %v", err)
	}

	// We only check for presubmit jobs.
	if ec.JobType == prow.PresubmitJob {
		gc, err := ghutil.NewGithubClient(*githubTokenPath)
		if err != nil {
			log.Fatalf("Cannot authenticate to github: %v", err)
		}

		pn := ec.PullNumber
		org := ec.RepoOwner
		repo := ec.RepoName
		pr, err := gc.GetPullRequest(org, repo, int(pn))
		if err != nil {
			log.Fatalf("Cannot find the pull request %d: %v", int(pn), err)
		}

		// If the PR is created by the bot, skip the check.
		if *pr.User.Name == config.ProwBotName {
			return
		}

		files, err := gc.ListFiles(org, repo, int(pn))
		if err != nil {
			log.Fatalf("Cannot find files changed in this PR: %v", err)
		}

		// Collect all files that are not allowed to change directly by users.
		bannedFiles := make([]string, 0)
		for _, file := range files {
			for _, p := range config.ProdProwKeyConfigPaths {
				fileName := file.GetFilename()
				if strings.HasPrefix(fileName, p) {
					bannedFiles = append(bannedFiles, fileName)
				}
			}
		}

		// TODO(chizhg): should not allow other production Prow config files to be changed if staging is needed.

		// If any of the production Prow key config files are changed, report the error.
		if len(bannedFiles) != 0 {
			log.Fatalf(
				"Directly changing the production Prow cluster config and templates is not allowed, please revert:\n%s",
				strings.Join(bannedFiles, "\n"))
		}
	}
}
