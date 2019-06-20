/*
Copyright 2019 The Knative Authors

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

// flaky-test-reporter collects test results from continuous flows,
// identifies flaky tests, tracking flaky tests related github issues,
// and sends slack notifications.

package main

import (
	"flag"
	"log"
	"os"

	"github.com/knative/test-infra/shared/prow"
)

func main() {
	serviceAccount := flag.String("service-account", os.Getenv("GOOGLE_APPLICATION_CREDENTIALS"), "JSON key file for GCS service account")
	githubAccount := flag.String("github-account", "", "Token file for Github authentication")
	slackAccount := flag.String("slack-account", "", "slack secret file for authenticating with Slack")
	dryrun := flag.Bool("dry-run", false, "dry run switch")
	flag.Parse()

	if nil != dryrun && true == *dryrun {
		log.Printf("running in [dry run mode]")
	}

	if err := prow.Initialize(*serviceAccount); nil != err { // Explicit authenticate with gcs Client
		log.Fatalf("Failed authenticating GCS: '%v'", err)
	}
	ghi, err := Setup(*githubAccount)
	if err != nil {
		log.Fatalf("Cannot setup github: %v", err)
	}
	slackClient, err := newSlackClient(*slackAccount)
	if nil != err {
		log.Fatalf("Failed authenticating Slack: '%v'", err)
	}

	var repoDataAll []*RepoData
	// Clean up local artifacts directory, this will be used later for artifacts uploads
	err = os.RemoveAll(prow.GetLocalArtifactsDir()) // this function returns nil if path not found
	if nil != err {
		log.Fatalf("Failed removing local artifacts directory: %v", err)
	}
	for repoName, jobList := range jobConfigs {
		for _, jc := range jobList {
			jc.Repo = repoName
			log.Printf("collecting results for job '%s' in repo '%s'\n", jc.Name, jc.Repo)
			rd, err := collectTestResultsForRepo(jc)
			if nil != err {
				log.Fatalf("Error collecting results for job '%s' in repo '%s': %v", jc.Name, jc.Repo, err)
			}
			if err = createArtifactForRepo(rd); nil != err {
				log.Fatalf("Error creating artifacts for job '%s' in repo '%s': %v", jc.Name, jc.Repo, err)
			}
			repoDataAll = append(repoDataAll, rd)
		}
	}

	// Errors that could result in inaccuracy reporting would be treated with fast fail by processGithubIssues,
	// so any errors returned are github opeations error, which in most cases wouldn't happen, but in case it
	// happens, it should fail the job after Slack notification
	githubErr := ghi.processGithubIssues(repoDataAll, *dryrun)
	slackErr := sendSlackNotifications(repoDataAll, slackClient, ghi, *dryrun)
	jsonErr := writeFlakyTestsToJSON(repoDataAll, *dryrun)
	if nil != githubErr {
		log.Printf("Github step failures:\n%v", githubErr)
	}
	if nil != slackErr {
		log.Printf("Slack step failures:\n%v", slackErr)
	}
	if nil != jsonErr {
		log.Printf("JSON step failures:\n%v", jsonErr)
	}
	if nil != githubErr || nil != slackErr || nil != jsonErr { // Fail this job if there is any error
		os.Exit(1)
	}
}
