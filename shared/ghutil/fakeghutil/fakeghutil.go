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
	"github.com/knative/test-infra/shared/ghutil"
)

// FakeGithubClient is a faked client, implements all functions of ghutil.GithubClientInterface
type FakeGithubClient struct {
	User     *github.User
	Repos    []string
	Issues   map[string]map[int]*github.Issue       // map of repo: map of issueNumber: issues
	Comments map[int]map[int64]*github.IssueComment // map of issueNumber: map of commentID: comments

	NextNumber int    // number to be assigned to next newly created issue/comment
	BaseURL    string // base URL of Github
}

// NewFakeGithubClient creates a FakeGithubClient and initialize it's maps
func NewFakeGithubClient() *FakeGithubClient {
	return &FakeGithubClient{
		Issues:   make(map[string]map[int]*github.Issue),
		Comments: make(map[int]map[int64]*github.IssueComment),
		BaseURL:  "fakeurl",
	}
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
	var issues []*github.Issue
	for _, issue := range fgc.Issues[repo] {
		labelMap := make(map[string]bool)
		for _, label := range issue.Labels {
			labelMap[*label.Name] = true
		}
		missingLabel := false
		for _, label := range labels {
			if _, ok := labelMap[label]; !ok {
				missingLabel = true
				break
			}
		}
		if !missingLabel {
			issues = append(issues, issue)
		}
	}
	return issues, nil
}

// CreateIssue creates issue
func (fgc *FakeGithubClient) CreateIssue(org, repo, title, body string) (*github.Issue, error) {
	issueNumber := fgc.getNextNumber()
	stateStr := string(ghutil.IssueOpenState)
	repoURL := fmt.Sprintf("%s/%s/%s", fgc.BaseURL, org, repo)
	url := fmt.Sprintf("%s/%d", repoURL, issueNumber)
	newIssue := &github.Issue{
		Title:         &title,
		Body:          &body,
		Number:        &issueNumber,
		State:         &stateStr,
		URL:           &url,
		RepositoryURL: &repoURL,
	}
	if _, ok := fgc.Issues[repo]; !ok {
		fgc.Issues[repo] = make(map[int]*github.Issue)
	}
	fgc.Issues[repo][issueNumber] = newIssue
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
	var comments []*github.IssueComment
	for _, comment := range fgc.Comments[issueNumber] {
		comments = append(comments, comment)
	}
	return comments, nil
}

// GetComment gets comment by comment ID
func (fgc *FakeGithubClient) GetComment(org, repo string, commentID int64) (*github.IssueComment, error) {
	for _, comments := range fgc.Comments {
		if comment, ok := comments[commentID]; ok {
			return comment, nil
		}
	}
	return nil, fmt.Errorf("cannot find comment")
}

// CreateComment adds comment to issue
func (fgc *FakeGithubClient) CreateComment(org, repo string, issueNumber int, commentBody string) (*github.IssueComment, error) {
	commentID := int64(fgc.getNextNumber())
	newComment := &github.IssueComment{
		ID:   &commentID,
		Body: &commentBody,
		User: fgc.User,
	}
	if _, ok := fgc.Comments[issueNumber]; !ok {
		fgc.Comments[issueNumber] = make(map[int64]*github.IssueComment)
	}
	fgc.Comments[issueNumber][commentID] = newComment
	return newComment, nil
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
	targetIssue := fgc.Issues[repo][issueNumber]
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
	targetIssue := fgc.Issues[repo][issueNumber]
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

func (fgc *FakeGithubClient) updateIssueState(org, repo string, state ghutil.IssueStateEnum, issueNumber int) error {
	targetIssue := fgc.Issues[repo][issueNumber]
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
