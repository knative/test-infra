package main

import (
	"flag"
	"log"

	"knative.dev/pkg/test/ghutil"
)

func main() {
	mainGithubTokenFile := flag.String("main-github-token-file", "",
		"Token file for Github authentication, used for most of important interactions with Github")
	commentGithubTokenFile := flag.String("comment-github-token-file", "",
		"Token file for Github authentication, used for adding comments on Github")
	dryrun := flag.Bool("dry-run", false, "dry run switch")
	flag.Parse()

	mgc, err := ghutil.NewGithubClient(*mainGithubTokenFile)
	if err != nil {
		log.Fatalf("Failed creating main github client: %v", err)
	}
	cgc, err := ghutil.NewGithubClient(*commentGithubTokenFile)
	if err != nil {
		log.Fatalf("Failed creating commenter github client: %v", err)
	}

	pr, err := getLatestPullRequest(mgc)
	if err != nil {
		log.Fatalf("Failed getting the latest PR number: %v", err)
	}
	fs, err := getChangedFiles(mgc, *pr.Number)
	if err != nil {
		log.Fatalf("Failed getting changed files in PR %q: %v", *pr.Number, err)
	}

	if err := runProwConfigUpdate(mgc, cgc, pr, fs, *dryrun); err != nil {
		log.Fatalf("Failed updating Prow configs: %v", err)
	}
}
