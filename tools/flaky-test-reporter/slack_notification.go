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

// slack_notification.go sends notifications to slack channels

package main

import (
	"fmt"
	"log"
	"sync"
	"time"

	"knative.dev/test-infra/shared/slackutil"
	"knative.dev/test-infra/shared/testgrid"
)

const (
	knativeBotName = "Knative Testgrid Robot"
	// default filter for testgrid link
	testgridFilter = "exclude-non-failed-tests=20"
)

// createSlackMessageForRepo creates slack message layout from RepoData
func createSlackMessageForRepo(rd *RepoData, flakyIssuesMap map[string][]*flakyIssue) string {
	flakyTests := getFlakyTests(rd)
	message := fmt.Sprintf("As of %s, there are %d flaky tests in '%s' from repo '%s'",
		time.Unix(*rd.LastBuildStartTime, 0).String(), len(flakyTests), rd.Config.Name, rd.Config.Repo)
	if "" == rd.Config.IssueRepo {
		message += fmt.Sprintf("\n(Job is marked to not create GitHub issues)")
	}
	if flakyRateAboveThreshold(rd) { // Don't list each test as this can be huge
		flakyRate := getFlakyRate(rd)
		message += fmt.Sprintf("\n>- skip displaying all tests as flaky rate above threshold")
		if flakyIssues, ok := flakyIssuesMap[getBulkIssueIdentity(rd, flakyRate)]; ok && "" != rd.Config.IssueRepo {
			// When flaky rate is above threshold, there is only one issue created,
			// so there is only one element in flakyIssues
			for _, fi := range flakyIssues {
				message += fmt.Sprintf("\t%s", fi.issue.GetHTMLURL())
			}
		}
	} else {
		for _, testFullName := range flakyTests {
			message += fmt.Sprintf("\n>- %s", testFullName)
			if flakyIssues, ok := flakyIssuesMap[getIdentityForTest(testFullName, rd.Config.Repo)]; ok && "" != rd.Config.IssueRepo {
				for _, fi := range flakyIssues {
					message += fmt.Sprintf("\t%s", fi.issue.GetHTMLURL())
				}
			}
		}
	}

	if testgridTabURL, err := testgrid.GetTestgridTabURL(rd.Config.Name, []string{testgridFilter}); nil != err {
		log.Println(err) // don't fail as this could be optional
	} else {
		message += fmt.Sprintf("\nSee Testgrid for up-to-date flaky tests information: %s", testgridTabURL)
	}
	return message
}

func sendSlackNotifications(repoDataAll []*RepoData, c slackutil.Operations, gih *GithubIssueHandler, dryrun bool) error {
	var allErrs []error
	flakyIssuesMap, err := gih.getFlakyIssues()
	if nil != err { // failure here will make message missing Github issues link, which should not prevent notification
		allErrs = append(allErrs, err)
		log.Println("Warning: cannot get flaky Github issues: ", err)
	}
	for _, rd := range repoDataAll {
		channels := rd.Config.SlackChannels
		if len(channels) == 0 {
			log.Printf("cannot find Slack channel for job '%s' in repo '%s', skipping Slack notification", rd.Config.Name, rd.Config.Repo)
			continue
		}
		ch := make(chan bool, len(channels))
		wg := sync.WaitGroup{}
		for i := range channels {
			wg.Add(1)
			channel := channels[i]
			go func(wg *sync.WaitGroup) {
				message := createSlackMessageForRepo(rd, flakyIssuesMap)
				if err := run(
					fmt.Sprintf("post Slack message for job '%s' from repo '%s' in channel '%s'", rd.Config.Name, rd.Config.Repo, channel.Name),
					func() error {
						return c.Post(message, channel.Identity)
					},
					dryrun,
				); nil != err {
					allErrs = append(allErrs, err)
					log.Printf("failed sending notification to Slack channel '%s': '%v'", channel.Name, err)
				}
				if dryrun {
					log.Printf("[dry run] Slack message not sent. See it below:\n%s\n\n", message)
				}
				ch <- true
				wg.Done()
			}(&wg)
		}
		wg.Wait()
		close(ch)
	}
	return combineErrors(allErrs)
}
