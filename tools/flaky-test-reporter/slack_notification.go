/*
Copyright 2018 The Knative Authors

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

// This implementation is based on assumption that data is already collected

package main

import (
	"errors"
	"io/ioutil"
	"strings"
	"strconv"
	"log"
	"fmt"

	"net/http"
	"net/url"

	"github.com/knative/test-infra/shared/testgrid"
)

const (
	slackChatPostMessageURL = "https://slack.com/api/chat.postMessage"
)

var (
	slackTokenStr string
	slackUsername string
	slackIconEmoji *string
)

var repoSlackChannelMap = map[string][]string{
	"serving": []string{"api"},
}

func getError(errs []error) error {
	if 0 == len(errs) {
		return nil
	}
	var errStrs []string
	for _, err := range errs {
		errStrs = append(errStrs, err.Error())
	}
	return fmt.Errorf(strings.Join(errStrs, "\n"))
}

// authenticateSlack reads token file and stores it for later authentication
func authenticateSlack(slackTokenPath *string) error {
	b, err := ioutil.ReadFile(*slackTokenPath)
	if err != nil {
		return err
	}
	slackTokenStr = strings.TrimSpace(string(b))
	return nil
}

func setUsername(userName string) {
	slackUsername = userName
}

func setIconEmoji(iconEmoji string) {
	slackIconEmoji = &iconEmoji
}

// postMessage does http post
func postMessage(url string, uv *url.Values) ([]byte, error) {
	resp, err := http.PostForm(url, *uv)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		t, _ := ioutil.ReadAll(resp.Body)
		return nil, errors.New(string(t))
	}
	t, _ := ioutil.ReadAll(resp.Body)
	return t, nil
}

// writeSlackMessage adds text to channel
func writeSlackMessage(text, channel string) error {
	uv := &url.Values{}
	uv.Add("username", slackUsername)
	uv.Add("token", slackTokenStr)
	if nil != slackIconEmoji {
		uv.Add("icon_emoji", *slackIconEmoji)
	}
	uv.Add("channel", channel)
	uv.Add("text", text)

	_, err := postMessage(slackChatPostMessageURL, uv)
	return err
}

// createSlackMessageForRepo creates slack message layout from RepoData
// <[datetime]>, active flaky tests [count]:
// 	[TestName#1]
//	[...]
// 	[TestName#n]
// See [Testgrid tab] for up-to-date views
func createSlackMessageForRepo(rd *RepoData) string {
	message := ""
	flakyTests := rd.getFlakyTests()
	message += fmt.Sprintf("<%s>, active flaky tests [%d]:", strconv.Itoa(int(*rd.LastBuildStartTime)), len(flakyTests))
	for testName := range flakyTests {
		message += fmt.Sprintf("\n\t%s", testName)
	}
	message += fmt.Sprintf("\nSee Testgrid tab for most recent flaky tests:\n%s", testgrid.GetTabURL(rd.Config.Name))
	return message
}

func sendSlackNotifications(repoDataAll []*RepoData, dryrun *bool) error {
	var allErrs []error
	for _, rd := range repoDataAll {
		channels, ok := repoSlackChannelMap[rd.Config.Repo]
		if !ok || len(channels) == 0 {
			errMsg := fmt.Sprintf("cannot find Slack channel for repo '%s', skip Slack notifiction", rd.Config.Repo)
			allErrs = append(allErrs, fmt.Errorf(errMsg))
			log.Printf(errMsg)
			continue
		}
		for _, channel := range channels {
			if err := writeSlackMessage(createSlackMessageForRepo(rd), channel); nil != err {
				allErrs = append(allErrs, err)
				log.Printf("failed sending notification to Slack channel '%s'", channel)
			}
		}
	}
	return getError(allErrs)
}
