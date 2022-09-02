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

// actions.go provides generic functions related to Actions

package ghutil

import (
	"github.com/google/go-github/v32/github"
)

// ListWorkflows lists workflows for a given org/repo.
func (gc *GithubClient) ListWorkflows(org, repo string) ([]*github.Workflow, error) {
	listOptions := &github.ListOptions{}
	genericList, err := gc.depaginate(
		"listing workflows",
		maxRetryCount,
		listOptions,
		func() ([]interface{}, *github.Response, error) {
			workflows, resp, err := gc.Client.Actions.ListWorkflows(ctx, org, repo, listOptions)
			var interfaceList []interface{}
			if nil == err {
				for _, workflow := range workflows.Workflows {
					interfaceList = append(interfaceList, workflow)
				}
			}
			return interfaceList, resp, err
		},
	)
	res := make([]*github.Workflow, len(genericList))
	for i, elem := range genericList {
		res[i] = elem.(*github.Workflow)
	}
	return res, err
}
