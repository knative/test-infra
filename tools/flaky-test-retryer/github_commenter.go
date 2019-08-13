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
	"sort"
	"strconv"
	"strings"

	"github.com/google/go-github/github"
	"knative.dev/test-infra/shared/ghutil"
)

const (
	maxRetries            = 3
	maxFailedTestsToPrint = 8
)

var (
	identifier      = "<!--AUTOMATED-FLAKY-RETRYER-->"
	commentTemplate = "%s\nThe following jobs failed due to test flakiness:\n\nTest name | Triggers | Retries\n--- | --- | ---\n%s\n\n%s"
)

// GithubClient wraps the ghutil Github client
type GithubClient struct {
	ghutil.GithubOperations
	ID     int64
	Dryrun bool
}

// entry holds all of the relevant information for a retried job
type entry struct {
	oldLinks string
	retries  int
}

// NewGithubClient builds us a GitHub client based on the token file passed in
func NewGithubClient(githubAccount string, dryrun bool) (*GithubClient, error) {
	ghc, err := ghutil.NewGithubClient(githubAccount)
	if err != nil {
		return nil, err
	}
	user, err := ghc.GetGithubUser()
	if err != nil {
		return nil, err
	}
	return &GithubClient{ghc, user.GetID(), dryrun}, nil
}

// PostComment posts a new comment on the PR specified in JobData, retrying the job that triggered it.
// The comment body is dynamically built based on previous retry comments on this PR, and any old
// comments are removed before the new one is posted.
func (gc *GithubClient) PostComment(jd *JobData, outliers []string) error {
	oldComment, err := gc.getOldComment(jd.Refs[0].Org, jd.Refs[0].Repo, jd.Refs[0].Pulls[0].Number)
	if err != nil {
		return err
	}
	oldEntries, err := parseEntries(oldComment)
	if err != nil {
		return err
	}
	if _, ok := oldEntries[jd.JobName]; !ok {
		oldEntries[jd.JobName] = &entry{}
	}
	newComment := buildNewComment(jd, oldEntries, outliers)
	if gc.Dryrun {
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
	var match *github.IssueComment
	for _, comment := range comments {
		found, err := regexp.Match(identifier, []byte(comment.GetBody()))
		if err != nil {
			return nil, err
		}
		if found && *comment.GetUser().ID == gc.ID {
			if match == nil {
				match = comment
			} else {
				return nil, fmt.Errorf("more than one comment on PR")
			}
		}
	}

	return match, nil
}

// parseEntries collects retry information from the given comment, so we can reuse it in
// a new comment.
func parseEntries(comment *github.IssueComment) (map[string]*entry, error) {
	entries := make(map[string]*entry)
	if comment == nil {
		return entries, nil
	}
	re := regexp.MustCompile(`.* \| \d`)
	entryStrings := re.FindAll([]byte(comment.GetBody()), -1)
	for _, e := range entryStrings {
		fields := strings.Split(string(e), " | ")
		if len(fields) != 3 {
			return nil, fmt.Errorf("invalid number of table entries")
		}
		retry, err := strconv.Atoi(strings.Split(fields[2], "/")[0])
		if err != nil {
			return nil, err
		}
		entries[fields[0]] = &entry{
			oldLinks: fields[1],
			retries:  retry,
		}
	}
	return entries, nil
}

// buildNewComment takes the old entry data, the job we are processing, and any outlying
// non-flaky tests, building a comment body based on these parameters.
func buildNewComment(jd *JobData, entries map[string]*entry, outliers []string) string {
	var cmd string
	var entryString []string
	if entries[jd.JobName].retries >= maxRetries {
		cmd = buildOutOfRetriesString(jd.JobName)
		logWithPrefix(jd, "expended all %d retries\n", maxRetries)
	} else if len(outliers) > 0 {
		cmd = buildNoRetryString(jd.JobName, outliers)
		logWithPrefix(jd, "%d failed tests are not flaky, cannot retry\n", len(outliers))
	} else {
		cmd = buildRetryString(jd.JobName, entries)
		logWithPrefix(jd, "all failed tests are flaky, triggering retry\n")
	}
	// print in sorted order so we can actually unit test the results
	var keys []string
	for test := range entries {
		keys = append(keys, test)
	}
	sort.Strings(keys)
	for _, test := range keys {
		entryString = append(entryString, fmt.Sprintf("%s | %s | %d/%d", test, buildLinks(entries[test].oldLinks, jd.URL, jd.RunID), entries[test].retries, maxRetries))
	}
	return fmt.Sprintf(commentTemplate, identifier, strings.Join(entryString, "\n"), cmd)
}

// buildLinks constructs a Markdown-formatted list of URLs
func buildLinks(oldLinks, newLink, id string) string {
	if oldLinks == "" {
		return fmt.Sprintf("[%s](%s)", id, newLink)
	}
	return fmt.Sprintf("%s<br>[%s](%s)", oldLinks, id, newLink)
}

// buildRetryString increments the retry counter and generates a /test string if we have
// more retries available.
func buildRetryString(job string, entries map[string]*entry) string {
	if entries[job].retries++; entries[job].retries <= maxRetries {
		return fmt.Sprintf("Automatically retrying...\n/test %s", job)
	}
	return ""
}

// buildNoRetryString formats the tests that prevent us from retrying into a truncated list.
func buildNoRetryString(job string, outliers []string) string {
	noRetryFmt := "Failed non-flaky tests preventing automatic retry of %s:\n\n```\n%s\n```%s"
	extraFailedTests := ""

	lastIndex := len(outliers)
	if len(outliers) > maxFailedTestsToPrint {
		lastIndex = maxFailedTestsToPrint
		extraFailedTests = fmt.Sprintf("\n\nand %d more.", len(outliers)-maxFailedTestsToPrint)
	}
	return fmt.Sprintf(noRetryFmt, job, strings.Join(outliers[:lastIndex], "\n"), extraFailedTests)
}

//buildOutOfRetriesString notifies the author that the job has been retriggered maxRetries times
// while still failing.
func buildOutOfRetriesString(job string) string {
	return fmt.Sprintf("Job %s expended all %d retries without success.", job, maxRetries)
}
