package main

import (
	"fmt"
	"strings"

	"knative.dev/pkg/test/ghutil"
	"knative.dev/pkg/test/helpers"

	"knative.dev/test-infra/tools/prow-config-updater/config"
)

const (
	successUpdatedConfigsCommentTemplate = "Updated %s configs with below files being modified: \n%s"
	failureUpdatedConfigsCommentTemplate = "Failed updating %s configs with below files being modified: \n%s\nSee error detail below: %v"
	successStagingCommentTemplate        = "The staging process is completed and succeeded.\nCreated an auto-merge PR to sync the changes to production Prow."
	failureStagingCommentTemplate        = "The staging process failed with the error below:\n%v\n\nPlease check if there is anything going wrong."
)

type GitHubCommenter struct {
	client *ghutil.GithubClient
	dryrun bool
}

// Comment on the pull request for the Prow update success result.
func (ghc *GitHubCommenter) commentOnUpdateConfigsSuccess(prnumber int, env config.ProwEnv, files []string) error {
	comment := fmt.Sprintf(successUpdatedConfigsCommentTemplate, env, fileListCommentString(files))
	return helpers.Run(
		fmt.Sprintf("Creating 'update success' comment %q on PR %q in %s/%s",
			comment, prnumber, config.OrgName, config.RepoName),
		func() error {
			_, err := ghc.client.CreateComment(config.OrgName, config.RepoName, prnumber, comment)
			return err
		},
		ghc.dryrun,
	)
}

// Comment on the pull request for the Prow update failure result.
func (ghc *GitHubCommenter) commentOnUpdateConfigsFailure(prnumber int, env config.ProwEnv, files []string, err error) error {
	comment := fmt.Sprintf(failureUpdatedConfigsCommentTemplate, env, fileListCommentString(files), err)
	return helpers.Run(
		fmt.Sprintf("Creating 'update failure' comment %q on PR %q in %s/%s",
			comment, prnumber, config.OrgName, config.RepoName),
		func() error {
			_, err := ghc.client.CreateComment(config.OrgName, config.RepoName, prnumber, comment)
			return err
		},
		ghc.dryrun,
	)
}

// Comment on the pull request to get it tested and merged by Prow.
func (ghc *GitHubCommenter) commentToMergePullRequest(prnumber int) error {
	commentToAdd := "/ok-to-test\n/lgtm\n/approve"
	return helpers.Run(
		fmt.Sprintf("Add comment %q on the pull request", commentToAdd),
		func() error {
			_, err := ghc.client.CreateComment(config.ForkOrgName, config.RepoName, prnumber, commentToAdd)
			return err
		},
		ghc.dryrun,
	)
}

// Comment on the pull request for the staging success result.
func (ghc *GitHubCommenter) commentOnStagingSuccess(prnumber int, err error) error {
	return helpers.Run(
		fmt.Sprintf("Creating 'staging success' comment %q on PR %q in %s/%s",
			successStagingCommentTemplate, prnumber, config.OrgName, config.RepoName),
		func() error {
			_, err := ghc.client.CreateComment(config.OrgName, config.RepoName, prnumber, successStagingCommentTemplate)
			return err
		},
		ghc.dryrun,
	)
}

// Comment on the pull request for the staging success result.
func (ghc *GitHubCommenter) commentOnStagingFailure(prnumber int, err error) error {
	comment := fmt.Sprintf(failureStagingCommentTemplate, err)
	return helpers.Run(
		fmt.Sprintf("Creating 'staging failure' comment %q on PR %q in %s/%s",
			comment, prnumber, config.OrgName, config.RepoName),
		func() error {
			_, err := ghc.client.CreateComment(config.OrgName, config.RepoName, prnumber, comment)
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
