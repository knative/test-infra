/*
Copyright 2022 The Knative Authors

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

// pullrequest.go creates git commits and Pull Requests

package pkg

import (
	"fmt"
	"log"
	"strings"

	"github.com/google/go-github/v32/github"

	"knative.dev/test-infra/pkg/ghutil"
	"knative.dev/test-infra/pkg/git"
	"knative.dev/test-infra/pkg/helpers"
)

// CreateOrUpdatePR creates a new PR if it does not exist, or updates the PR if it
// already exists.
func CreateOrUpdatePR(gc ghutil.GithubOperations, gi git.Info, label string, dryrun bool) error {
	const matchTitle = "[Auto] Update prow jobs for release branches"
	commitMsg := matchTitle
	title := commitMsg
	body := generatePRBody()
	hasUpdates, err := git.MakeCommit(gi, commitMsg, dryrun)
	if err != nil {
		return fmt.Errorf("failed git commit: %w", err)
	}
	if !hasUpdates {
		log.Print("There is nothing committed, skip PR")
		return nil
	}

	var existPR *github.PullRequest
	existPR, err = getExistingPR(gc, gi, matchTitle)
	if err != nil {
		return fmt.Errorf("failed querying existing pullrequests: %w", err)
	}
	if existPR != nil {
		log.Printf("Found open PR %d", *existPR.Number)
		if err := helpers.Run(
			fmt.Sprintf("Updating PR %d, title: %q, body: %q", *existPR.Number, title, body),
			func() error {
				if _, err := gc.EditPullRequest(gi.Org, gi.Repo, *existPR.Number, title, body); err != nil {
					return fmt.Errorf("failed updating pullrequest: %w", err)
				}
				return nil
			},
			dryrun,
		); err != nil {
			return err
		}
	} else {
		if err := helpers.Run(
			fmt.Sprintf("Creating PR, title: %q, body: %q", title, body),
			func() error {
				existPR, err = gc.CreatePullRequest(gi.Org, gi.Repo, gi.GetHeadRef(), gi.Base, title, body)
				if err != nil {
					return fmt.Errorf("failed creating pullrequest: %w", err)
				}
				return nil
			},
			dryrun,
		); err != nil {
			return err
		}
	}

	if label != "" {
		if err := helpers.Run(
			fmt.Sprintf("Ensure label %q exists for PR", label),
			func() error {
				err = gc.EnsureLabelForPullRequest(gi.Org, gi.Repo, *existPR.Number, label)
				if err != nil {
					return fmt.Errorf("failed ensuring label %q exists: %w", label, err)
				}
				return nil
			},
			dryrun,
		); err != nil {
			return err
		}
	}

	return nil
}

func generatePRBody() string {
	body := "PR created for syncing Prow jobs for release branches changes\n\n"
	assignment := "/cc knative/productivity-wg-leads\n/assign knative/productivity-wg-leads\n\n"
	producedBy := "Produced by: release-jobs-syncer"

	return body + assignment + producedBy
}

// get existing open PR not merged yet
func getExistingPR(gc ghutil.GithubOperations, gi git.Info, matchTitle string) (*github.PullRequest, error) {
	var res *github.PullRequest
	PRs, err := gc.ListPullRequests(gi.Org, gi.Repo, gi.GetHeadRef(), gi.Base)
	if err == nil {
		for _, PR := range PRs {
			if string(ghutil.PullRequestOpenState) == *PR.State && strings.Contains(*PR.Title, matchTitle) {
				res = PR
				break
			}
		}
	}
	return res, err
}
