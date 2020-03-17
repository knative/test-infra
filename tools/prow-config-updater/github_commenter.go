/*
Copyright 2020 The Knative Authors

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
	failureUpdatedConfigsCommentTemplate = "Failed updating %s configs with below files being modified: \n%s\nSee error details below: %v"
	successStagingCommentTemplate        = "The staging process is completed and succeeded.\nWill create an auto-merge PR to roll out the changes to production."
	failureStagingCommentTemplate        = "The staging process failed with the error below:\n%v\n\nPlease check if there is anything going wrong."
	successRolloutCommentTemplate        = "Created #%d to roll out the staging changes to production."
	failureRolloutCommentTemplate        = "Failed to create an auto-merge PR to roll out the staging change to production, please check the error below:\n%v"
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

// Try to get the fork pull request tested and merged by staging Prow.
func (ghc *GitHubCommenter) tryMergeForkPullRequest(forkOrgName string, prnumber int) error {
	commentToAdd := "/ok-to-test\n/lgtm\n/approve"
	labelToAdd := []string{"cla: yes"}
	return helpers.Run(
		fmt.Sprintf("Add comment %q on the pull request", commentToAdd),
		func() error {
			errs := make([]error, 2)
			_, errs[0] = ghc.client.CreateComment(forkOrgName, config.RepoName, prnumber, commentToAdd)
			errs[1] = ghc.client.AddLabelsToIssue(forkOrgName, config.RepoName, prnumber, labelToAdd)
			return helpers.CombineErrors(errs)
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

// Comment on the pull request for the staging failure result.
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

// Comment on the pull request for the rollout success result.
func (ghc *GitHubCommenter) commentOnRolloutSuccess(prnumber, newprnumber int) error {
	comment := fmt.Sprintf(successRolloutCommentTemplate, newprnumber)
	return helpers.Run(
		fmt.Sprintf("Creating 'rollout success' comment %q on PR %q in %s/%s",
			comment, prnumber, config.OrgName, config.RepoName),
		func() error {
			_, err := ghc.client.CreateComment(config.OrgName, config.RepoName, prnumber, comment)
			return err
		},
		ghc.dryrun,
	)
}

// Comment on the pull request for the rollout failure result.
func (ghc *GitHubCommenter) commentOnRolloutFailure(prnumber int, err error) error {
	comment := fmt.Sprintf(failureRolloutCommentTemplate, err)
	return helpers.Run(
		fmt.Sprintf("Creating 'rollout failure' comment %q on PR %q in %s/%s",
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
