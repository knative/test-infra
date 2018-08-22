package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/google/go-github/github"
)

var (
	// Info about the current PR
	repoOwner = os.Getenv("REPO_OWNER")
	repoName = os.Getenv("REPO_NAME")
	pullNumber = atoi(os.Getenv("PULL_NUMBER"), "pull number")

	// Shared useful variables
	ctx = context.Background()
	onePageList = &github.ListOptions{Page: 1}
	verbose = false
	anonymousGitHubClient *github.Client
)

// atoi is a convenience function to convert a string to integer, failing in case of error.
func atoi(str, valueName string) int {
        value, err := strconv.Atoi(str)
        if err != nil {
                log.Fatalf("Unexpected non number '%s' for %s: %v", str, valueName, err)
        }
        return value
}

// infof if a convenience wrapper around log.Infof, and does nothing unless --verbose is passed.
func infof(template string, args ...interface{}) {
	if verbose {
		log.Printf(template, args...)
	}
}

func main() {
	listChangedFilesFlag := flag.Bool("list-changed-files", false, "List the files changed by the current pull request")
	verboseFlag := flag.Bool("verbose", false, "Whether to dump extra info on output or not; intended for debugging")
        flag.Parse()

	verbose = *verboseFlag
	anonymousGitHubClient = github.NewClient(nil)
	if *listChangedFilesFlag {
		listChangedFiles()
	}
}

func listChangedFiles() {
	infof("Listing changed files for PR %d in repository %s/%s", pullNumber, repoOwner, repoName)
	files, _, err := anonymousGitHubClient.PullRequests.ListFiles(ctx, repoOwner, repoName, pullNumber, onePageList)
	if err != nil {
		log.Fatalf("Error listing files: %v", err)
	}
	for _, file := range files {
		fmt.Println(*file.Filename)
	}
}
