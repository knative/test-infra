package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/google/go-github/github"
	"knative.dev/pkg/test/cmd"
	"knative.dev/pkg/test/ghutil"
	"knative.dev/pkg/test/helpers"

	"knative.dev/test-infra/shared/common"
	"knative.dev/test-infra/tools/prow-config-updater/config"
)

type ConfigUpdater struct {
	githubclient *ghutil.GithubClient
	pullrequest  *github.PullRequest
	files        []string
	dryrun       bool
}

const (
	updateProwCommandTemplate     = "make -C %s update-prow-cluster"
	updateTestgridCommandTemplate = "make -C %s update-testgrid-config"
)

func runProwConfigUpdate(gc *ghutil.GithubClient, pr *github.PullRequest, fs []string, dryrun bool) error {
	if err := common.CDToRootDir(); err != nil {
		return err
	}

	cu := ConfigUpdater{
		githubclient: gc,
		pullrequest:  pr,
		files:        fs,
		dryrun:       dryrun,
	}
	// If the PR is created by the Prow bot, we should be confident to blindly update production Prow configs.
	if *pr.User.Name == config.ProwBotName {
		if err := cu.updateProw(config.ProdProwEnv); err != nil {
			return fmt.Errorf("error updating production Prow configs: %v", err)
		}
	} else {
		// Check if we need staging process for updating the config files.
		if cu.needsStaging() {
            if err := cu.startStaging(); err != nil {
                return fmt.Errorf("error running Prow staging process: %v", err)
            }
            // TODO(chizhg): create an auto-merge pull request to update production Prow configs.
		} else {
            // TODO(chizhg): update production Prow.
		}
	}
	return nil
}

// Update Prow with the changed config files and send message for the update result.
func (cu *ConfigUpdater) updateProw(env config.ProwEnv) error {
	updatedFiles, err := cu.doProwUpdate(env)
    ghc := GitHubCommenter{pullrequest: cu.pullrequest, client: cu.githubclient, dryrun: cu.dryrun}
    if err == nil {
		if commenterErr := ghc.commentOnUpdateSuccess(config.ProdProwEnv, updatedFiles); commenterErr != nil {
			return fmt.Errorf("error creating comment on GitHub for the successful update: %v", commenterErr)
		}
	} else {
        if commenterErr := ghc.commentOnUpdateFailure(config.StagingProwEnv, updatedFiles, err); commenterErr != nil {
            return fmt.Errorf("error creating comment on GitHub for the failed update: %v", commenterErr)
        }
	}
	return err
}

// Decide if we need staging process by checking the changed files.
func (cu *ConfigUpdater) needsStaging() bool {
	// If any key config files for staging Prow are changed, the staging process will be needed.
	fs := collectRelevantFiles(cu.files, config.StagingProwKeyConfigPaths)
	return len(fs) != 0
}

// Start running the staging process to update and test staging Prow.
func (cu *ConfigUpdater) startStaging() error {
    // Update staging Prow.
    if err := cu.updateProw(config.StagingProwEnv); err != nil {
        return fmt.Errorf("error updating staging Prow configs: %v", err)
    }
    // Create pull request in the fork repository for the testing of staging Prow.
    fpr, err := createForkPullRequest(cu.githubclient)
    if err != nil {
        return fmt.Errorf("error creating pull request in the fork repository: %v", err)
    }
    // Create comment on the fork pull request to get it tested by staging Prow and merged.
    ghc := GitHubCommenter{client:cu.githubclient, pullrequest:fpr, dryrun:cu.dryrun}
    if err = ghc.commentToMergePullRequest(); err != nil {
        return fmt.Errorf("error creating comment on the fork pull request: %v", err)
    }
    // Wait for the fork pull request to be automatically merged by staging Prow.
    if err := waitForForkPullRequestMerged(ghc.client, *fpr.Number); err != nil {
        return fmt.Errorf("error waiting for the fork pull request to be merged: %v", err)
    }
    return nil
}

// Filter out all files that are under the given paths.
func collectRelevantFiles(files []string, paths []string) []string {
	rfs := make([]string, 0)
	for _, f := range files {
		for _, p := range paths {
			if strings.HasPrefix(f, p) {
				rfs = append(rfs, f)
			}
		}
	}
	return rfs
}

// Run the `make` command to update Prow configs.
// path is the Prow config root folder.
func (cu *ConfigUpdater) doProwUpdate(env config.ProwEnv) ([]string, error) {
	relevantFiles := make([]string, 0)
	updateCommand := ""
	switch env {
	case config.ProdProwEnv:
		relevantFiles = append(relevantFiles, collectRelevantFiles(cu.files, config.ProdProwConfigPaths)...)
		updateCommand = fmt.Sprintf(updateProwCommandTemplate, config.ProdProwConfigRoot)
	case config.StagingProwEnv:
		relevantFiles = append(relevantFiles, collectRelevantFiles(cu.files, config.StagingProwConfigPaths)...)
		updateCommand = fmt.Sprintf(updateProwCommandTemplate, config.StagingProwConfigRoot)
	default:
		return nil, fmt.Errorf("unsupported Prow environement: %q, cannot make the update", env)
	}
	if len(relevantFiles) != 0 {
		if err := helpers.Run(
			fmt.Sprintf("Updating Prow configs with command %q", updateCommand),
			func() error {
				out, err := cmd.RunCommand(updateCommand)
				log.Println(out)
				return err
			},
			cu.dryrun,
		); err != nil {
			return nil, fmt.Errorf("error updating Prow configs for %q environment: %v", env, err)
		}
	}

	// For production Prow, we also need to update Testgrid config if it's changed.
	// TODO(chizhg): this will be removed after we get rid of Testgrid config file by moving to ProwJob annotation.
	if env == config.ProdProwEnv {
		tfs := collectRelevantFiles(cu.files, []string{config.ProdTestgridConfigFile})
		if len(tfs) != 0 {
			relevantFiles = append(relevantFiles, tfs...)
			updateTestgridCommand := fmt.Sprintf(updateTestgridCommandTemplate, config.ProdProwConfigRoot)
            if err := helpers.Run(
                fmt.Sprintf("Updating Testgrid config with command %q", updateTestgridCommand),
                func() error {
                    out, err := cmd.RunCommand(updateTestgridCommand)
                    log.Println(out)
                    return err
                },
                cu.dryrun,
            ); err != nil {
                return nil, fmt.Errorf("error updating Testgrid configs for %q environment: %v", env, err)
            }
		}
	}
	return relevantFiles, nil
}


