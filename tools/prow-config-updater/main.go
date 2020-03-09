package main

import (
    "flag"
    "log"

    "knative.dev/pkg/test/ghutil"
)

func main() {
    githubTokenFile := flag.String("github-token-file", "", "Token file for Github authentication")
    dryrun := flag.Bool("dry-run", false, "dry run switch")
    flag.Parse()

    gc, err := ghutil.NewGithubClient(*githubTokenFile)
    if err != nil {
        log.Fatalf("Failed creating github client: %v", err)
    }

    pr, err := getLatestPullRequest(gc)
    if err != nil {
        log.Fatalf("Failed getting the latest PR number: %v", err)
    }
    fs, err := getChangedFiles(gc, *pr.Number)
    if err != nil {
        log.Fatalf("Failed getting changed files in PR %q: %v", *pr.Number, err)
    }

    if err := runProwConfigUpdate(gc, pr, fs, *dryrun); err != nil {
        log.Fatalf("Failed updating Prow configs: %v", err)
    }
}
