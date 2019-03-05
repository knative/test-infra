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

package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"strings"
	"sync"
	"time"

	"net/http"
	"net/url"

	"github.com/knative/test-infra/shared/testgrid"
)

const (
	knativeBotName          = "Knative bot"
	slackChatPostMessageURL = "https://slack.com/api/chat.postMessage"
	// default filter for testgrid link
	testgridFilter = "exclude-non-failed-tests=20"
)

var (
	// jobNameTestgridURLMap contains harded coded mapping of job name: Testgrid tab URL relative to base URL
	jobNameTestgridURLMap = map[string]string{
		"ci-knative-serving-continuous": "knative-serving#continuous",
	}
	// slackChannelsMap defines mapping of repo: slack channels
	slackChannelsMap = map[string][]slackChannel{
		"serving": []slackChannel{{"test", "CA1DTGZ2N"}},
	}
)

type SlackClient struct {
	userName  string
	tokenStr  string
	iconEmoji *string
}

// slackChannel contains channel logical name and Slack identity
type slackChannel struct {
	name, identity string
}

// newSlackClient reads token file and stores it for later authentication
func newSlackClient(slackTokenPath string) (*SlackClient, error) {
	b, err := ioutil.ReadFile(slackTokenPath)
	if err != nil {
		return nil, err
	}
	return &SlackClient{
		userName: knativeBotName,
		tokenStr: strings.TrimSpace(string(b)),
	}, nil
}

// postMessage does http post
func (c *SlackClient) postMessage(url string, uv *url.Values) error {
	resp, err := http.PostForm(url, *uv)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	t, _ := ioutil.ReadAll(resp.Body)
	if resp.StatusCode != 200 {
		return fmt.Errorf("http response code is not 200 '%s'", string(t))
	}
	// response code could also be 200 if channel doesn't exist, parse response body to find out
	var b struct {
		OK bool `json:"ok"`
	}
	if err = json.Unmarshal(t, &b); nil != err || !b.OK {
		return fmt.Errorf("response not ok '%s'", string(t))
	}
	return nil
}

// writeSlackMessage posts text to channel
func (c *SlackClient) writeSlackMessage(text, channel string) error {
	uv := &url.Values{}
	uv.Add("username", c.userName)
	uv.Add("token", c.tokenStr)
	if nil != c.iconEmoji {
		uv.Add("icon_emoji", *c.iconEmoji)
	}
	uv.Add("channel", channel)
	uv.Add("text", text)

	return c.postMessage(slackChatPostMessageURL, uv)
}

// getTestgridTabURL gets Testgrid URL for giving job
func getTestgridTabURL(jobName string) (string, error) {
	url, ok := jobNameTestgridURLMap[jobName]
	if !ok {
		return "", fmt.Errorf("cannot find Testgrid tab for job '%s'", jobName)
	}
	return fmt.Sprintf("%s/%s&%s", testgrid.BaseURL, url, testgridFilter), nil
}

// createSlackMessageForRepo creates slack message layout from RepoData
func createSlackMessageForRepo(rd *RepoData) string {
	flakyTests := getFlakyTests(rd)
	message := fmt.Sprintf("As of %s, there are %d flaky tests in '%s'",
		time.Unix(*rd.LastBuildStartTime, 0).String(), rd.Config.Repo, len(flakyTests))
	for _, testName := range flakyTests {
		message += fmt.Sprintf("\n>- %s", testName)
		// TODO(chaodaiG): try adding Github issues when PR #506 merged
	}
	if testgirdTabURL, err := getTestgridTabURL(rd.Config.Name); nil != err {
		log.Println(err) // don't fail as this could be optional
	} else {
		message += fmt.Sprintf("\nSee Testgrid tab for most recent flaky tests: %s", testgirdTabURL)
	}
	return message
}

func sendSlackNotifications(repoDataAll []*RepoData, c *SlackClient, dryrun *bool) error {
	var allErrs []error
	for _, rd := range repoDataAll {
		channels, ok := slackChannelsMap[rd.Config.Repo]
		if !ok {
			log.Printf("cannot find Slack channel for repo '%s', skip Slack notifiction", rd.Config.Repo)
			continue
		}

		ch := make(chan bool, len(channels))
		wg := sync.WaitGroup{}
		for _, channel := range channels {
			wg.Add(1)
			go func(wg *sync.WaitGroup) {
				message := createSlackMessageForRepo(rd)
				if err := run(
					fmt.Sprintf("post Slack message for repo '%s' in channel '%s'", rd.Config.Repo, channel.name),
					func() error {
						return c.writeSlackMessage(message, channel.identity)
					},
					dryrun); nil != err {
					allErrs = append(allErrs, err)
					log.Printf("failed sending notification to Slack channel '%s': '%v'", channel.name, err)
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
