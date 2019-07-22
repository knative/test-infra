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

// github_commenter.go finds the relevant pull requests for the failed jobs that
// triggered the retryer and posts comments to it, either retrying the test or
// telling the contributors why we cannot retry.

package main

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/google/go-github/github"
	"github.com/knative/test-infra/shared/ghutil"
)

const (
	maxRetries = 3
	maxFailedTestsToPrint = 8
)

var (
	identifier = "<!--AUTOMATED-FLAKY-RETRYER-->"
	commentTemplate = "%s\nThe following tests are currently flaky. Running them again to verify..."+
		"\n\nTest name | Retries\n--- | ---\n%s\n\n%s"
)

// GithubClient wraps the ghutil Github client
type GithubClient struct {
	*ghutil.GithubClient
	Login string
}

// NewGithubClient builds us a GitHub client based on the token file passed in
func NewGithubClient(githubAccount string) (*GithubClient, error) {
	ghc, err := ghutil.NewGithubClient(githubAccount)
	if err != nil {
		return nil, err
	}
	user, err := ghc.GetGithubUser()
	if err != nil {
		return nil, err
	}
	return &GithubClient{ghc, *user.Login}, nil
}

// PostComment posts a new comment on the PR specified in JobData, retrying the job that triggered it.
// The comment body is dynamically built based on previous retry comments on this PR, and any old
// comments are removed before the new one is posted.
func (gc *GithubClient) PostComment(jd *JobData, outliers []string, dryrun bool) error {
	oldComment, err := gc.getOldComment(jd.Refs[0].Org, jd.Refs[0].Repo, jd.Refs[0].Pulls[0].Number)
	if err != nil {
		return err
	}
	oldEntries, err := parseEntries(oldComment)
	if err != nil {
		return err
	}
	if _, ok := oldEntries[jd.JobName]; !ok {
		oldEntries[jd.JobName] = 0
	}
	newComment, canRetry := buildNewComment(jd, oldEntries, outliers)
	if !canRetry {
		return fmt.Errorf("expended all %d retries", maxRetries)
	}
	if dryrun {
		logWithPrefix(jd, "[dry run] Comment not updated. See it here:\n%s\n", newComment)
		return nil
	}
	if oldComment != nil {
		if err := gc.DeleteComment(jd.Refs[0].Org, jd.Refs[0].Repo, oldComment.GetID()); err != nil {
			return err
		}
	}
	_, err = gc.CreateComment(jd.Refs[0].Org, jd.Refs[0].Repo, jd.Refs[0].Pulls[0].Number, newComment)
	return err
}

// getOldComment queries the GitHub PR specified and gets the comment made by us. If no comment
// is found, we do not error, since we will be creating a new one anyways.
func (gc *GithubClient) getOldComment(org, repo string, pull int) (*github.IssueComment, error) {
	comments, err := gc.ListComments(org, repo, pull)
	if err != nil {
		return nil, err
	}
	// this robot should only leave one comment in a PR. Get it by finding the comment with the identifier
	for _, comment := range comments {
		match, err := regexp.Match(identifier, []byte(comment.GetBody()))
		if err != nil {
			return nil, err
		}
		if match && *comment.GetUser().Login == gc.Login {
			return comment, nil
		}
	}
	return nil, nil
}

// parseEntries collects retry information from the given comment, so we can reuse it in
// a new comment.
func parseEntries(comment *github.IssueComment) (map[string]int, error) {
	entries := map[string]int{}
	if comment == nil {
		return entries, nil
	}
	re := regexp.MustCompile(`.* \| \d`)
	entryStrings := re.FindAll([]byte(comment.GetBody()), -1)
	for _, e := range entryStrings {
		fields := strings.Split(string(e), " | ")
		retry, err := strconv.Atoi(strings.TrimSuffix(fields[1], "/3"))
		if err != nil {
			return nil, err
		}
		entries[strings.Trim(fields[0], "`")] = retry
	}
	return entries, nil
}

// buildNewComment takes the old entry data, the job we are processing, and any outlying
// non-flaky tests, building a comment body based on these parameters.
func buildNewComment(jd *JobData, entries map[string]int, outliers []string) (string, bool) {
	var cmd string
	var entryString []string
	if len(outliers) > 0 {
		cmd = buildNoRetryString(jd.JobName, outliers)
		logWithPrefix(jd, "%d failed tests are not flaky, cannot retry\n", len(outliers))
	} else {
		cmd = buildRetryString(jd.JobName, entries)
		logWithPrefix(jd, "all failed tests are flaky, triggering retry\n")
	}
	for test, retry := range entries {
		entryString = append(entryString, fmt.Sprintf("%s | %d/%d", test, retry, maxRetries))
	}
	return fmt.Sprintf(commentTemplate, identifier, strings.Join(entryString, "\n"), cmd), entries[jd.JobName] <= maxRetries
}

// buildRetryString increments the retry counter and generates a /test string if we have
// more retries available.
func buildRetryString(job string, entries map[string]int) string {
	entries[job]++
	if entries[job] <= maxRetries {
		return fmt.Sprintf("Automatically retrying...\n/test %s", job)
	}
	return ""
}

// buildNoRetryString formats the tests that prevent us from retrying into a list of 10
// top entries and
func buildNoRetryString(job string, outliers []string) string {
	noRetryFmt := "Failed non-flaky tests preventing automatic retry of %s:\n\n```\n%s\n```%s"
	extraFailedTests := ""

	lastIndex := len(outliers)
	if len(outliers) > maxFailedTestsToPrint {
		lastIndex = maxFailedTestsToPrint
		extraFailedTests = fmt.Sprintf("\n\nand %d more.", len(outliers) - maxFailedTestsToPrint)
	}
	return fmt.Sprintf(noRetryFmt, job, strings.Join(outliers[:lastIndex], "\n"), extraFailedTests)
}
