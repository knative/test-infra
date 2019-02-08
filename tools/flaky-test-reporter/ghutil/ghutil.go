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

// ghutil.go provides generic functions supporting Github operations
// It helps overcome Github rate limit errors

package ghutil

import (
	"context"
	"io/ioutil"
	"strings"
	"log"
	"fmt"
	"time"

	"golang.org/x/oauth2"
	"github.com/google/go-github/github"
)

const (
	maxRetryCount = 5
	tokenReserve  = 50

	// IssueOpenState is the state of open github issue
	IssueOpenState IssueStateEnum = "open"
	// IssueCloseState is the state of closed github issue
	IssueCloseState IssueStateEnum = "closed"
	// IssueAllState is the state for all, useful when querying issues
	IssueAllState  IssueStateEnum = "all"
)

// IssueStateEnum represents different states of Github Issues
type IssueStateEnum string

// client for github authentication
// without explicitly authenticate this client can still be used for anonymous operations
var client *github.Client
var ctx = context.Background()

// Authenticate explicitly authenticate github with giving token
func Authenticate(tokenFilePath *string) error {
	b, err := ioutil.ReadFile(*tokenFilePath)
	if err != nil {
		return err
	}
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: strings.TrimSpace(string(b))},
	)
	client = github.NewClient(oauth2.NewClient(ctx, ts))
	return nil
}

// GetCurrentUser gets authenticated user
func GetCurrentUser() (*github.User, error) {
	var res *github.User
	_, err := retry(
		"getting current user",
		maxRetryCount,
		func() (*github.Response, error) {
			var resp *github.Response
			var err error
			res, resp, err = client.Users.Get(ctx, "")
			return resp, err
		},
	)
	return res, err
}

// ListRepos lists repos under org
func ListRepos(org string) ([]string, error) {
	var res []string
	options := &github.ListOptions{}
	genericList, err := depaginate(
		"listing repos",
		maxRetryCount,
		options,
		func() ([]interface{}, *github.Response, error) {
			page, resp, err := client.Repositories.List(ctx, org, nil)
			var interfaceList []interface{}
			if nil == err {
				for _, repo := range page {
					interfaceList = append(interfaceList, repo)
				}
			}
			return interfaceList, resp, err
		},
	)
	for _, repo := range genericList {
		res = append(res, repo.(*github.Repository).GetName())
	}
	return res, err
}

// ListIssuesByRepo lists issues within given repo, filters by labels if provided
func ListIssuesByRepo(org, repo string, labels []string) ([]*github.Issue, error) {
	issueListOptions := github.IssueListByRepoOptions{
		State:  string(IssueAllState),
	}
	if len(labels) > 0 {
		issueListOptions.Labels = labels
	}
	
	var res []*github.Issue
	options := &github.ListOptions{}
	genericList, err := depaginate(
		fmt.Sprintf("listing issues with label '%v'", labels),
		maxRetryCount,
		options,
		func() ([]interface{}, *github.Response, error) {
			page, resp, err := client.Issues.ListByRepo(ctx, org, repo, &issueListOptions)
			var interfaceList []interface{}
			if nil == err {
				for _, issue := range page {
					interfaceList = append(interfaceList, issue)
				}
			}
			return interfaceList, resp, err
		},
	)
	for _, issue := range genericList {
		res = append(res, issue.(*github.Issue))
	}
	return res, err
}

// CreateIssue creates issue
func CreateIssue(org, repo, title, body string) (*github.Issue, error) {
	issue := &github.IssueRequest{
		Title: &title,
		Body:  &body,
	}

	var res *github.Issue
	_, err := retry(
		fmt.Sprintf("creating issue '%s %s' '%s'", org, repo, title),
		maxRetryCount,
		func() (*github.Response, error) {
			var resp *github.Response
			var err error
			res, resp, err = client.Issues.Create(ctx, org, repo, issue)
			return resp, err
		},
	)
	return res, err
}

// CloseIssue closes issue
func CloseIssue(org, repo string, issueNumber int) error {
	return updateIssueState(org, repo, IssueCloseState, issueNumber)
}

// ReopenIssue reopen issue
func ReopenIssue(org, repo string, issueNumber int) error {
	return updateIssueState(org, repo, IssueOpenState, issueNumber)
}

// ListComments gets all comments from issue
func ListComments(org, repo string, issueNumber int) ([]*github.IssueComment, error) {
	var res []*github.IssueComment
	options := &github.ListOptions{}
	genericList, err := depaginate(
		fmt.Sprintf("listing comment for issue '%s %s %d'", org, repo, issueNumber),
		maxRetryCount,
		options,
		func() ([]interface{}, *github.Response, error) {
			page, resp, err := client.Issues.ListComments(ctx, org, repo, issueNumber, nil)
			var interfaceList []interface{}
			if nil == err {
				for _, issue := range page {
					interfaceList = append(interfaceList, issue)
				}
			}
			return interfaceList, resp, err
		},
	)
	for _, issue := range genericList {
		res = append(res, issue.(*github.IssueComment))
	}
	return res, err
}

// GetComment gets comment by comment ID
func GetComment(org, repo string, commentID int64) (*github.IssueComment, error) {
	var res *github.IssueComment
	_, err := retry(
		fmt.Sprintf("getting comment '%s %s %d'", org, repo, commentID),
		maxRetryCount,
		func() (*github.Response, error) {
			var resp *github.Response
			var err error
			res, resp, err = client.Issues.GetComment(ctx, org, repo, commentID)
			return resp, err
		},
	)
	return res, err
}

// CreateComment adds comment to issue
func CreateComment(org, repo string, issueNumber int, commentBody string) error {
	comment := &github.IssueComment{
		Body: &commentBody,
	}
	_, err := retry(
		fmt.Sprintf("commenting issue '%s %s %d'", org, repo, issueNumber),
		maxRetryCount,
		func() (*github.Response, error) {
			_, resp, err := client.Issues.CreateComment(ctx, org, repo, issueNumber, comment)
			return resp, err
		},
	)
	return err
}

// EditComment edits comment by replacing with provided comment
func EditComment(org, repo string, commentID int64, commentBody string) error {
	comment := &github.IssueComment{
		Body: &commentBody,
	}
	_, err := retry(
		fmt.Sprintf("editing comment '%s %s %d'", org, repo, commentID),
		maxRetryCount,
		func() (*github.Response, error) {
			_, resp, err := client.Issues.EditComment(ctx, org, repo, commentID, comment)
			return resp, err
		},
	)
	return err
}

// AddLabelsToIssue adds label on issue
func AddLabelsToIssue(org, repo string, issueNumber int, labels []string) error {
	_, err := retry(
		fmt.Sprintf("add labels '%v' to '%s %s %d'", labels, org, repo, issueNumber),
		maxRetryCount,
		func() (*github.Response, error) {
			_, resp, err := client.Issues.AddLabelsToIssue(ctx, org, repo, issueNumber, labels)
			return resp, err
		},
	)
	return err
}

// RemoveLabelForIssue removes given label for issue
func RemoveLabelForIssue(org, repo string, issueNumber int, label string) error {
	_, err := retry(
		fmt.Sprintf("remove label '%s' from '%s %s %d'", label, org, repo, issueNumber),
		maxRetryCount,
		func() (*github.Response, error) {
			return client.Issues.RemoveLabelForIssue(ctx, org, repo, issueNumber, label)
		},
	)
	return err
}

func updateIssueState(org, repo string, state IssueStateEnum, issueNumber int) error {
	stateString := string(state)
	issueRequest := &github.IssueRequest{
		State: &stateString,
	}
	_, err := retry(
		fmt.Sprintf("applying '%s' action on issue '%s %s %d'", stateString, org, repo, issueNumber),
		maxRetryCount,
		func() (*github.Response, error) {
			_, resp, err := client.Issues.Edit(ctx, org, repo, issueNumber, issueRequest)
			return resp, err
		},
	)
	return err
}

func waitForRateReset(r *github.Rate) {
	if r.Remaining <= tokenReserve {
		sleepDuration := time.Until(r.Reset.Time) + (time.Second * 10)
		if sleepDuration > 0 {
			log.Printf("--Rate Limiting-- GitHub tokens reached minimum reserve %d. Sleeping %ds until reset.\n", tokenReserve, sleepDuration)
			time.Sleep(sleepDuration)
		}
	}
}

// Github API has a rate limit, retry waits until rate limit reset if request failed with RateLimitError,
// then retry maxRetries times until succeed
func retry(message string, maxRetries int, call func() (*github.Response, error)) (*github.Response, error) {
	var err error
	var resp *github.Response

	for retryCount := 0; retryCount <= maxRetries; retryCount++ {
		if resp, err = call(); nil == err {
			return resp, nil
		}
		switch err := err.(type) {
			case *github.RateLimitError:
				waitForRateReset(&err.Rate)
			default:
				return resp, err
		}
		log.Printf("error %s: %v. Will retry.\n", message, err)
	}
	return resp, err
}

// depaginate adds depagination on top of the retry, in case list exceeds rate limit
func depaginate(message string, maxRetries int, options *github.ListOptions, call func() ([]interface{}, *github.Response, error)) ([]interface{}, error) {
	var allItems []interface{}
	wrapper := func() (*github.Response, error) {
		items, resp, err := call()
		if err == nil {
			allItems = append(allItems, items...)
		}
		return resp, err
	}

	options.Page = 1
	options.PerPage = 100
	lastPage := 1
	for ; options.Page <= lastPage; options.Page++ {
		resp, err := retry(message, maxRetries, wrapper)
		if err != nil {
			return allItems, fmt.Errorf("error while depaginating page %d/%d: %v", options.Page, lastPage, err)
		}
		if resp.LastPage > 0 {
			lastPage = resp.LastPage
		}
	}
	return allItems, nil
}
