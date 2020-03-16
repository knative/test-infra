package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/google/go-github/github"
	"k8s.io/apimachinery/pkg/util/wait"
	"knative.dev/pkg/test/cmd"
	"knative.dev/pkg/test/ghutil"

	"knative.dev/test-infra/tools/prow-config-updater/config"
)

// Get the latest pull request created in this repository.
// TODO(chizhg): get rid of this hack once Prow supports setting PR number as an env var for postsubmit jobs.
func getLatestPullRequest(gc *ghutil.GithubClient) (*github.PullRequest, error) {
	// Use git command to get the latest commit ID.
	ci, err := cmd.RunCommand("git rev-parse HEAD")
	if err != nil {
		return nil, fmt.Errorf("error getting the last commit ID: %v", err)
	}
	// As we always use squash in merging PRs, we can get the pull request with the commit ID.
	pr, err := gc.GetPullRequestByCommitID(config.OrgName, config.RepoName, ci)
	if err != nil {
		return nil, fmt.Errorf("error getting the PR with commit ID %q: %v", ci, err)
	}
	return pr, nil
}

// Get the list of changed file names with the given PR number.
func getChangedFiles(gc *ghutil.GithubClient, pn int) ([]string, error) {
	fs, err := gc.ListFiles(config.OrgName, config.RepoName, pn)
	if err != nil {
		return nil, fmt.Errorf("error listing the changed files for PR %q: %v", pn, err)
	}

	fns := make([]string, len(fs))
	for i := range fns {
		fns[i] = *fs[i].Filename
	}

	return fns, nil
}

// Use the pull bot (https://github.com/wei/pull) to create a PR in the fork repository.
func createForkPullRequest(gc *ghutil.GithubClient) (*github.PullRequest, error) {
	// The endpoint to manually trigger pull bot to create the pull request in the fork.
	pullTriggerEndpoint := fmt.Sprintf("https://pull.git.ci/process/%s/%s", config.ForkOrgName, config.RepoName)
	resp, err := http.Get(pullTriggerEndpoint)
	if err != nil {
		return nil, fmt.Errorf("error creating the pull request in the fork repository: %v", err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading the response for creating the pull request in the fork repository: %v", err)
	}
	if resp.StatusCode != http.StatusOK || string(body) != "Success" {
		return nil, fmt.Errorf("error creating the pull request in the fork repository: "+
			"status code is %d, body is %s", resp.StatusCode, body)
	}
	return findForkPullRequest(gc)
}

// Find the pull request created by the pull bot (should be exactly one pull request, otherwise there must be an error)
func findForkPullRequest(gc *ghutil.GithubClient) (*github.PullRequest, error) {
	prs, _, err := gc.Client.PullRequests.List(
		context.Background(), config.ForkOrgName, config.RepoName,
		&github.PullRequestListOptions{
			State: "open",
			Head:  "knative:master",
		})
	orgRepoName := fmt.Sprintf("%s/%s", config.ForkOrgName, config.RepoName)
	if err != nil {
		return nil, fmt.Errorf("error listing pull request in repository %s: %v", orgRepoName, err)
	}
	if len(prs) != 1 {
		return nil, fmt.Errorf("expected one pull request in repository %s but found %d", orgRepoName, len(prs))
	}
	return prs[0], nil
}

func waitForForkPullRequestMerged(gc *ghutil.GithubClient, pn int) error {
	interval := 10 * time.Second
	timeout := 20 * time.Minute
	return wait.PollImmediate(interval, timeout, func() (bool, error) {
		pr, err := gc.GetPullRequest(config.ForkOrgName, config.RepoName, pn)
		if err != nil {
			return false, err
		}
		if !*pr.Merged {
			return false, nil
		}
		return true, nil
	})
}

func createAutoMergePullRequest() error {
	// TODO(chizhg): create a pull request with the auto-merge label.
	return nil
}
