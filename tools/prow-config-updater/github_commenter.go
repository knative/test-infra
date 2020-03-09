package main

import (
	"fmt"
	"strings"

	"github.com/google/go-github/github"
	"knative.dev/pkg/test/ghutil"
	"knative.dev/pkg/test/helpers"

	"knative.dev/test-infra/tools/prow-config-updater/config"
)

const (
	successUpdatedConfigsCommentTemplate = "Updated %s configs with below files being modified: \n%s"
	failureUpdatedConfigsCommentTemplate = "Failed updating %s configs with below files being modified: \n%s\nSee error detail below: %v"
)

type GitHubCommenter struct {
	client      *ghutil.GithubClient
	pullrequest *github.PullRequest
	dryrun      bool
}

// Comment on the pull request for the Prow update success result.
func (ghc *GitHubCommenter) commentOnUpdateSuccess(env config.ProwEnv, files []string) error {
	comment := fmt.Sprintf(successUpdatedConfigsCommentTemplate, env, fileListCommentString(files))
	return helpers.Run(
		fmt.Sprintf("Creating 'update success' comment %q on PR %q in %s/%s",
			comment, *ghc.pullrequest.Number, config.OrgName, config.RepoName),
		func() error {
			_, err := ghc.client.CreateComment(config.OrgName, config.RepoName, *ghc.pullrequest.Number, comment)
			return err
		},
		ghc.dryrun,
	)
}

// Comment on the pull request for the Prow update failure result.
func (ghc *GitHubCommenter) commentOnUpdateFailure(env config.ProwEnv, files []string, err error) error {
	comment := fmt.Sprintf(failureUpdatedConfigsCommentTemplate, env, fileListCommentString(files), err)
	return helpers.Run(
		fmt.Sprintf("Creating 'update failure' comment %q on PR %q in %s/%s",
			comment, *ghc.pullrequest.Number, config.OrgName, config.RepoName),
		func() error {
			_, err := ghc.client.CreateComment(config.OrgName, config.RepoName, *ghc.pullrequest.Number, comment)
			return err
		},
		ghc.dryrun,
	)
}

// Comment on the pull request to get it tested and merged by Prow.
func (ghc *GitHubCommenter) commentToMergePullRequest() error {
	commentToAdd := "/ok-to-test\n/lgtm\n/approve"
	return helpers.Run(
		fmt.Sprintf("Add comment %q on the fork pull request", commentToAdd),
		func() error {
			_, err := ghc.client.CreateComment(config.ProwBotName, config.RepoName, *pr.Number, commentToAdd)
			return err
		},
		ghc.dryrun,
	)
}

func fileListCommentString(files []string) string {
	res := make([]string, len(files))
	for i, f := range files {
		res[i] = "- `" + f + "`"
	}
	return strings.Join(res, "\n")
}
