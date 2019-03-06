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

// fakeghutil.go fakes GithubClient for testing purpose

package fakeghutil

import (
	"fmt"

	"github.com/google/go-github/github"
	"github.com/knative/test-infra/tools/flaky-test-reporter/ghutil"
)

// FakeGithubClient is a faked client, implements all functions of ghutil.GithubClientInterface
type FakeGithubClient struct {
	User     *github.User
	Repos    []string
	Issues   map[string][]*github.Issue     // map of repo: slice of issues
	Comments map[int][]*github.IssueComment // map of issueNumber: slice of comments

	NextNumber int // number to be assigned to next newly created issue/comment
}

// GetGithubUser gets current authenticated user
func (fgc *FakeGithubClient) GetGithubUser() (*github.User, error) {
	return fgc.User, nil
}

// ListRepos lists repos under org
func (fgc *FakeGithubClient) ListRepos(org string) ([]string, error) {
	return fgc.Repos, nil
}

// ListIssuesByRepo lists issues within given repo, filters by labels if provided
func (fgc *FakeGithubClient) ListIssuesByRepo(org, repo string, labels []string) ([]*github.Issue, error) {
	return fgc.Issues[repo], nil
}

// CreateIssue creates issue
func (fgc *FakeGithubClient) CreateIssue(org, repo, title, body string) (*github.Issue, error) {
	issueNumber := fgc.getNextNumber()
	newIssue := &github.Issue{
		Title:  &title,
		Body:   &body,
		Number: &issueNumber,
	}
	fgc.ReopenIssue(org, repo, issueNumber)
	fgc.Issues[repo] = append(fgc.Issues[repo], newIssue)
	return newIssue, nil
}

// CloseIssue closes issue
func (fgc *FakeGithubClient) CloseIssue(org, repo string, issueNumber int) error {
	return fgc.updateIssueState(org, repo, ghutil.IssueCloseState, issueNumber)
}

// ReopenIssue reopen issue
func (fgc *FakeGithubClient) ReopenIssue(org, repo string, issueNumber int) error {
	return fgc.updateIssueState(org, repo, ghutil.IssueOpenState, issueNumber)
}

// ListComments gets all comments from issue
func (fgc *FakeGithubClient) ListComments(org, repo string, issueNumber int) ([]*github.IssueComment, error) {
	return fgc.Comments[issueNumber], nil
}

// GetComment gets comment by comment ID
func (fgc *FakeGithubClient) GetComment(org, repo string, commentID int64) (*github.IssueComment, error) {
	for _, comments := range fgc.Comments {
		for _, comment := range comments {
			if *comment.ID == commentID {
				return comment, nil
			}
		}
	}
	return nil, fmt.Errorf("cannot find comment")
}

// CreateComment adds comment to issue
func (fgc *FakeGithubClient) CreateComment(org, repo string, issueNumber int, commentBody string) error {
	commentID := int64(fgc.getNextNumber())
	newComment := &github.IssueComment{
		ID:   &commentID,
		Body: &commentBody,
	}
	fgc.Comments[issueNumber] = append(fgc.Comments[issueNumber], newComment)
	return nil
}

// EditComment edits comment by replacing with provided comment
func (fgc *FakeGithubClient) EditComment(org, repo string, commentID int64, commentBody string) error {
	comment, err := fgc.GetComment(org, repo, commentID)
	if nil != err {
		return err
	}
	comment.Body = &commentBody
	return nil
}

// AddLabelsToIssue adds label on issue
func (fgc *FakeGithubClient) AddLabelsToIssue(org, repo string, issueNumber int, labels []string) error {
	targetIssue := fgc.getIssue(org, repo, issueNumber)
	if nil == targetIssue {
		return fmt.Errorf("cannot find issue")
	}
	for _, label := range labels {
		targetIssue.Labels = append(targetIssue.Labels, github.Label{
			Name: &label,
		})
	}
	return nil
}

// RemoveLabelForIssue removes given label for issue
func (fgc *FakeGithubClient) RemoveLabelForIssue(org, repo string, issueNumber int, label string) error {
	targetIssue := fgc.getIssue(org, repo, issueNumber)
	if nil == targetIssue {
		return fmt.Errorf("cannot find issue")
	}
	targetI := -1
	for i, l := range targetIssue.Labels {
		if *l.Name == label {
			targetI = i
		}
	}
	if -1 == targetI {
		return fmt.Errorf("cannot find label")
	}
	targetIssue.Labels = append(targetIssue.Labels[:targetI], targetIssue.Labels[targetI+1:]...)
	return nil
}

func (fgc *FakeGithubClient) getIssue(org, repo string, issueNumber int) *github.Issue {
	var targetIssue *github.Issue
	for _, issue := range fgc.Issues[repo] {
		if *issue.Number == issueNumber {
			targetIssue = issue
		}
	}
	return targetIssue
}

func (fgc *FakeGithubClient) updateIssueState(org, repo string, state ghutil.IssueStateEnum, issueNumber int) error {
	targetIssue := fgc.getIssue(org, repo, issueNumber)
	if nil == targetIssue {
		return fmt.Errorf("cannot find issue")
	}
	stateStr := string(state)
	targetIssue.State = &stateStr
	return nil
}

func (fgc *FakeGithubClient) getNextNumber() int {
	fgc.NextNumber++
	return fgc.NextNumber
}
