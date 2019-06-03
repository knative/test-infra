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

// prow-auto-bumper finds stable Prow components version used by k8s,
// and creates PRs updating them in knative/test-infra

package main

import (
	"flag"
	"fmt"
	"log"
	"regexp"
	"time"

	"github.com/google/go-github/github"
	"github.com/knative/test-infra/shared/ghutil"
)

const (
	// Git info for k8s Prow auto bumper PRs
	srcOrg  = "kubernetes"
	srcRepo = "test-infra"
	// srcPRHead is the head branch of k8s auto version bump PRs
	// TODO(chaodaiG): using head branch querying is less ideal than using
	// label `area/prow/bump`, which is not supported by Github API yet. Move
	// to filter using this label once it's supported
	srcPRHead = "autobump"
	// srcPRBase is the base branch of k8s auto version bump PRs
	srcPRBase = "master"
	// srcPRUserID is the user from which PR was created
	srcPRUserID = "k8s-ci-robot"

	// Git info for target repo that Prow version bump PR targets
	org  = "knative"
	repo = "test-infra"
	// PRHead is branch name where the changes occur
	PRHead = "autobump"
	// PRBase is the branch name where PR targets
	PRBase = "master"

	// Index for regex matching groups
	imageImagePart = 1
	imageTagPart   = 2
	// Max difference away from target date
	maxDelta = 2 * 24 // 2 days
	// Safe duration is the smallest amount of hours a version stayed
	safeDuration = 12 // 12 hours
	maxRetry     = 3

	oncallAddress = "https://storage.googleapis.com/knative-infra-oncall/oncall.json"
)

var (
	fileFilters = []*regexp.Regexp{regexp.MustCompile(`\.yaml$`)}
	// matching            gcr.io /k8s-(prow|testimage)/(tide|kubekin-e2e|.*)    :vYYYYMMDD-HASH-VARIANT
	imagePattern     = `\b(gcr\.io/k8s[a-z0-9-]{5,29}/[a-zA-Z0-9][a-zA-Z0-9_.-]+):(v[a-zA-Z0-9_.-]+)\b`
	imageRegexp      = regexp.MustCompile(imagePattern)
	imageLinePattern = fmt.Sprintf(`\s+[a-z]+:\s+"?'?%s"?'?`, imagePattern)
	// matching   "-    image: gcr.io /k8s-(prow|testimage)/(tide|kubekin-e2e|.*)    :vYYYYMMDD-HASH-VARIANT"
	imageMinusRegexp = regexp.MustCompile(fmt.Sprintf(`\-%s`, imageLinePattern))
	// matching   "+    image: gcr.io /k8s-(prow|testimage)/(tide|kubekin-e2e|.*)    :vYYYYMMDD-HASH-VARIANT"
	imagePlusRegexp = regexp.MustCompile(fmt.Sprintf(`\+%s`, imageLinePattern))
	// Preferred time for candidate PR creation date
	targetTime = time.Now().Add(-time.Hour * 7 * 24) // 7 days
)

// GHClientWrapper handles methods for github issues
type GHClientWrapper struct {
	ghutil.GithubOperations
}

type gitInfo struct {
	org      string
	repo     string
	head     string // PR head branch
	base     string // PR base branch
	userID   string // Github User ID of PR creator
	userName string // User display name for Git commit
	email    string // User email address for Git commit
}

// HeadRef is in the form of "user:head", i.e. "github_user:branch_foo"
func (gi *gitInfo) getHeadRef() string {
	return fmt.Sprintf("%s:%s", gi.userID, gi.head)
}

// versions holds the version change for an image
// oldVersion and newVersion are both in the format of "vYYYYMMDD-HASH-VARIANT"
type versions struct {
	oldVersion string
	newVersion string
	variant    string
}

// PRVersions contains PR and version changes in it
type PRVersions struct {
	images map[string][]versions // map of image name: versions struct
	// The way k8s updates versions doesn't guarantee the same version tag across all images,
	// dominantVersions is the version that appears most times
	dominantVersions *versions
	PR               *github.PullRequest
}

func main() {
	githubAccount := flag.String("github-account", "", "Token file for Github authentication")
	gitUserID := flag.String("git-userid", "", "The github ID of user for hosting fork, i.e. Github ID of bot")
	gitUserName := flag.String("git-username", "", "The username to use on the git commit. Requires --git-email")
	gitEmail := flag.String("git-email", "", "The email to use on the git commit. Requires --git-username")
	dryrun := flag.Bool("dry-run", false, "dry run switch")
	flag.Parse()

	if nil != dryrun && true == *dryrun {
		log.Printf("running in [dry run mode]")
	}

	gc, err := ghutil.NewGithubClient(*githubAccount)
	if nil != err {
		log.Fatalf("cannot authenticate to github: %v", err)
	}

	srcGI := gitInfo{
		org:    srcOrg,
		repo:   srcRepo,
		head:   srcPRHead,
		base:   srcPRBase,
		userID: srcPRUserID,
	}

	targetGI := gitInfo{
		org:      org,
		repo:     repo,
		head:     PRHead,
		base:     PRBase,
		userID:   *gitUserID,
		userName: *gitUserName,
		email:    *gitEmail,
	}

	gcw := &GHClientWrapper{gc}
	bestVersion, err := retryGetBestVersion(gcw, srcGI)
	if nil != err {
		log.Fatalf("cannot get best version from %s/%s: '%v'", srcGI.org, srcGI.repo, err)
	}
	log.Println("Found version to update. Old Version: '%s', New Version: '%s'",
		bestVersion.dominantVersions.oldVersion, bestVersion.dominantVersions.newVersion)

	errMsgs, err := updateAllFiles(bestVersion, fileFilters, imageRegexp, *dryrun)
	if nil != err {
		log.Fatalf("failed updating files: '%v'", err)
	}

	if err = createOrUpdatePR(gcw, bestVersion, targetGI, errMsgs, *dryrun); nil != err {
		log.Fatalf("failed creating pullrequest: '%v'", err)
	}
}
