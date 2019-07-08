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

// github_commenter.go finds the relevant pull requests for the failed jobs that
// triggered the retryer and posts comments to it, either retrying the test or
// telling the contributors why we cannot retry.

package main

import (
	"log"

	"github.com/knative/test-infra/shared/ghutil"
)

// GithubClient wraps the ghutil Github client
type GithubClient struct {
	*ghutil.GithubClient
}

func NewGithubClient(githubAccount string) (*GithubClient, error) {
	ghc, err := ghutil.NewGithubClient(githubAccount)
	log.Printf("temporary - otherwise compiler will yell about ghc being an unused var: %s\n", ghc)
	return &GithubClient{ghc}, err
}
