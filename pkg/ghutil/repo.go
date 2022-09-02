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

// repo.go provides generic functions related to Repo

package ghutil

import (
	"fmt"

	"github.com/google/go-github/v32/github"
)

// ListRepos lists repos under org
func (gc *GithubClient) ListRepos(org string) ([]string, error) {
	repoListOptions := &github.RepositoryListOptions{}
	genericList, err := gc.depaginate(
		"listing repos",
		maxRetryCount,
		&repoListOptions.ListOptions,
		func() ([]interface{}, *github.Response, error) {
			page, resp, err := gc.Client.Repositories.List(ctx, org, repoListOptions)
			var interfaceList []interface{}
			if nil == err {
				for _, repo := range page {
					interfaceList = append(interfaceList, repo)
				}
			}
			return interfaceList, resp, err
		},
	)
	res := make([]string, len(genericList))
	for i, elem := range genericList {
		res[i] = elem.(*github.Repository).GetName()
	}
	return res, err
}

// ListBranches lists branchs for given repo
func (gc *GithubClient) ListBranches(org, repo string) ([]*github.Branch, error) {
	branchListOptions := &github.BranchListOptions{}
	genericList, err := gc.depaginate(
		fmt.Sprintf("listing Pull request from org %q and base %q", org, repo),
		maxRetryCount,
		&branchListOptions.ListOptions,
		func() ([]interface{}, *github.Response, error) {
			page, resp, err := gc.Client.Repositories.ListBranches(ctx, org, repo, branchListOptions)
			var interfaceList []interface{}
			if nil == err {
				for _, PR := range page {
					interfaceList = append(interfaceList, PR)
				}
			}
			return interfaceList, resp, err
		},
	)
	res := make([]*github.Branch, len(genericList))
	for i, elem := range genericList {
		res[i] = elem.(*github.Branch)
	}
	return res, err
}
