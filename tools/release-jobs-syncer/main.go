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

// release-jobs-syncer fetches latest release branches,
// and creates PRs to update Prow jobs for them in knative/test-infra

package main

import (
	"flag"
	"log"
	"path/filepath"

	"knative.dev/test-infra/pkg/ghutil"
	"knative.dev/test-infra/pkg/git"
	"knative.dev/test-infra/tools/release-jobs-syncer/pkg"
)

const (
	org  = "knative"
	repo = "test-infra"
	// PRHead is branch name where the changes occur
	PRHead = "releasebranch"
	// PRBase is the branch name where PR targets
	PRBase = "main"
)

func main() {
	// prowJobConfigRootPath := flag.String("prow-job-config-root-path",
	//  filepath.Join(helpers.MustGetRootDir(), "prow/jobs_config"), "Root path for the Prow jobs config")
	prowJobConfigRootPath := flag.String("prow-job-config-root-path",
		filepath.Join("/home/prow/go/src/knative.dev/test-infra", "prow/jobs_config"), "Root path for the Prow jobs config")
	githubAccount := flag.String("github-account", "", "Token file for Github authentication")
	gitUserID := flag.String("git-userid", "", "The github ID of user for hosting fork, i.e. Github ID of bot")
	gitUserName := flag.String("git-username", "", "The username to use on the git commit. Requires --git-email")
	gitEmail := flag.String("git-email", "", "The email to use on the git commit. Requires --git-username")
	label := flag.String("label", "", "The label to add on the PR")
	dryrun := flag.Bool("dry-run", false, "dry run switch")
	flag.Parse()

	if *dryrun {
		log.Println("Running in [dry run mode]")
	}

	gc, err := ghutil.NewGithubClient(*githubAccount)
	if err != nil {
		log.Fatalf("cannot authenticate to github: %v", err)
	}

	targetGI := git.Info{
		Org:      org,
		Repo:     repo,
		Head:     PRHead,
		Base:     PRBase,
		UserID:   *gitUserID,
		UserName: *gitUserName,
		Email:    *gitEmail,
	}

	if err := pkg.UpdateReleaseBranchConfig(gc, *prowJobConfigRootPath); err != nil {
		log.Fatalf("error updating release branch config: %v", err)
	}
	if err = pkg.CreateOrUpdatePR(gc, targetGI, *label, *dryrun); err != nil {
		log.Fatalf("error creating pullrequest: %v", err)
	}
}
