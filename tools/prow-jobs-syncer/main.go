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

// prow-jobs-syncer fetches release branches,
// and creates PRs updating them in knative/test-infra

package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path"
	"strings"

	"knative.dev/test-infra/pkg/cmd"
	"knative.dev/test-infra/pkg/ghutil"

	"knative.dev/test-infra/pkg/git"
)

func main() {
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

	gopath := os.Getenv("GOPATH")

	configgenArgs := []string{
		"--prow-jobs-config-output",
		path.Join(gopath, repoPath, jobConfigPath),
		"--testgrid-config-output",
		path.Join(gopath, repoPath, testgridConfigPath),
		"--upgrade-release-branches",
		"--github-token-path",
		*githubAccount,
		path.Join(gopath, repoPath, templateConfigPath),
	}

	configgenFullPath := path.Join(gopath, repoPath, configGenPath)

	log.Print(cmd.RunCommand(fmt.Sprintf("go run %s %s",
		configgenFullPath, strings.Join(configgenArgs, " "))))

	// The code gen above updates the template file, which might not be
	// sufficient for generating all prow/testgrid configs, rerun config gen
	// script to make everything up-to-date
	log.Print(cmd.RunCommand(configGenScript))

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

	gcw := &GHClientWrapper{gc}
	if err = createOrUpdatePR(gcw, targetGI, *label, *dryrun); err != nil {
		log.Fatalf("failed creating pullrequest: '%v'", err)
	}
}
