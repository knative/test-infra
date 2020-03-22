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
	"path/filepath"
	"strings"

	"github.com/google/go-github/github"
	"knative.dev/pkg/test/cmd"

	"knative.dev/test-infra/tools/prow-config-updater/config"
)

type Client struct {
	githubmainhandler *GitHubMainHandler
	githubcommenter   *GitHubCommenter
	forkOrgName       string
	pr                *github.PullRequest
	files             []string
	dryrun            bool
}

func (cli *Client) initialize() error {
	pr, err := cli.githubmainhandler.getLatestPullRequest()
	if err != nil {
		return fmt.Errorf("error getting the latest PR number: %v", err)
	}
	cli.pr = pr

	fs, err := cli.githubmainhandler.getChangedFiles(*pr.Number)
	if err != nil {
		return fmt.Errorf("error getting changed files in PR %q: %v", *pr.Number, err)
	}
	cli.files = fs

	return nil
}

func (cli *Client) runProwConfigUpdate() error {
	// If no staging process if needed, we can directly update production Prow configs.
	if !cli.needsStaging() {
		if err := cli.updateProw(config.ProdProwEnv); err != nil {
			return fmt.Errorf("error updating production Prow configs: %v", err)
		}
	} else {
		if err := cli.startStaging(); err != nil {
			// Best effort, won't fail the process if the comment fails.
			cli.githubcommenter.commentOnStagingFailure(*cli.pr.Number, err)
			return fmt.Errorf("error running Prow staging process: %v", err)
		} else {
			// Best effort, won't fail the process if the comment fails.
			cli.githubcommenter.commentOnStagingSuccess(*cli.pr.Number)
		}

		newpr, err := cli.rollOutToProd()
		if err != nil {
			// Best effort, won't fail the process if the comment fails.
			cli.githubcommenter.commentOnRolloutFailure(*cli.pr.Number, err)
			return fmt.Errorf("error rolling out the staging change to production: %v", err)
		} else {
			// Best effort, won't fail the process if the comment fails.
			cli.githubcommenter.commentOnRolloutSuccess(*cli.pr.Number, *newpr.Number)
		}
	}
	return nil
}

// Decide if we need staging process by checking the PR.
func (cli *Client) needsStaging() bool {
	// If the PR is created by the main github bot, we should be confident to blindly update production Prow configs.
	if *cli.pr.User.Login == cli.githubmainhandler.info.UserID {
		return false
	}
	// If any key config files for staging Prow are changed, the staging process will be needed.
	fs := collectRelevantFiles(cli.files, config.StagingProwKeyConfigPaths)
	return len(fs) != 0
}

// Update Prow with the changed config files and send message for the update result.
func (cli *Client) updateProw(env config.ProwEnv) error {
	updatedFiles, err := cli.doProwUpdate(env)
	prnumber := *cli.pr.Number
	if err == nil {
		// Best effort, won't fail the process if the comment fails.
		cli.githubcommenter.commentOnUpdateConfigsSuccess(prnumber, config.ProdProwEnv, updatedFiles)
	} else {
		// Best effort, won't fail the process if the comment fails.
		cli.githubcommenter.commentOnUpdateConfigsFailure(prnumber, config.ProdProwEnv, updatedFiles, err)
	}
	return err
}

// Start running the staging process to update and test staging Prow.
func (cli *Client) startStaging() error {
	// Update staging Prow.
	if err := cli.updateProw(config.StagingProwEnv); err != nil {
		return fmt.Errorf("error updating staging Prow configs: %v", err)
	}
	// Create pull request in the fork repository for the testing of staging Prow.
	fpr, err := cli.githubmainhandler.createForkPullRequest(cli.forkOrgName)
	if err != nil {
		return fmt.Errorf("error creating pull request in the fork repository: %v", err)
	}
	// Create comment on the fork pull request to get it tested by staging Prow and merged.
	forkprnumber := *fpr.Number
	if err = cli.githubcommenter.mergeForkPullRequest(cli.forkOrgName, forkprnumber); err != nil {
		return fmt.Errorf("error creating comment on the fork pull request: %v", err)
	}
	// Wait for the fork pull request to be automatically merged by staging Prow.
	if err := cli.githubmainhandler.waitForForkPullRequestMerged(cli.forkOrgName, forkprnumber); err != nil {
		return fmt.Errorf("error waiting for the fork pull request to be merged: %v", err)
	}
	return nil
}

// Filter out all files that are under the given paths.
func collectRelevantFiles(files []string, paths []string) []string {
	rfs := make([]string, 0)
	for _, f := range files {
		for _, p := range paths {
			if !strings.HasSuffix(p, string(filepath.Separator)) {
				p = p + string(filepath.Separator)
			}
			if strings.HasPrefix(f, p) {
				rfs = append(rfs, f)
			}
		}
	}
	return rfs
}

// Run the `make` command to update Prow configs.
// path is the Prow config root folder.
func (cli *Client) doProwUpdate(env config.ProwEnv) ([]string, error) {
	relevantFiles := make([]string, 0)
	switch env {
	case config.ProdProwEnv:
		relevantFiles = append(relevantFiles, collectRelevantFiles(cli.files, config.ProdProwConfigPaths)...)
	case config.StagingProwEnv:
		relevantFiles = append(relevantFiles, collectRelevantFiles(cli.files, config.StagingProwConfigPaths)...)
	default:
		return nil, fmt.Errorf("unsupported Prow environement: %q, cannot make the update", env)
	}
	if len(relevantFiles) != 0 {
		if err := config.UpdateProw(env, cli.dryrun); err != nil {
			return nil, fmt.Errorf("error updating Prow configs for %q environment: %v", env, err)
		}
	}

	// For production Prow, we also need to update Testgrid config if it's changed.
	tfs := collectRelevantFiles(cli.files, []string{config.ProdTestgridConfigPath})
	if len(tfs) != 0 {
		relevantFiles = append(relevantFiles, tfs...)
		if err := config.UpdateTestgrid(env, cli.dryrun); err != nil {
			return nil, fmt.Errorf("error updating Testgrid configs for %q environment: %v", env, err)
		}
	}
	return relevantFiles, nil
}

// Roll out the staging config files to production.
func (cli *Client) rollOutToProd() (*github.PullRequest, error) {
	// Copy staging config files to production.
	for i, stagingPath := range config.StagingProwKeyConfigPaths {
		cpCmd := fmt.Sprintf("cp -r %s/* %s", stagingPath, config.ProdProwKeyConfigPaths[i])
		if _, err := cmd.RunCommand(cpCmd); err != nil {
			return nil, fmt.Errorf("error copying staging config files to production: %v", err)
		}
	}

	// Try generating new config files.
	if err := config.GenerateConfigFiles(); err != nil {
		return nil, fmt.Errorf("error generating Prow config files for production: %v", err)
	}

	// Create a pull request to update production Prow.
	commitMsg := fmt.Sprintf("roll out staging Prow change in #%d to production", *cli.pr.Number)
	body := fmt.Sprintf(
		"This is a PR auto-synced from #%d, it will be automatically merged after all tests pass.",
		*cli.pr.Number)
	pr, err := cli.githubmainhandler.createAutoMergePullRequest(commitMsg, body)
	if err != nil {
		return nil, fmt.Errorf("error creating pull request to roll out staging Prow to production: %v", err)
	}
	return pr, nil
}
