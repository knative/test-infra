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

// data definitions that are used for the config file generation of k8s testgrid

package main

import (
	"regexp"
	"sort"

	"k8s.io/apimachinery/pkg/util/sets"
)

const (
	k8sTestgridTempl      = "k8s_testgrid.yaml"
	k8sTestgridGroupTempl = "k8s_testgrid_testgroup.yaml"
)

type k8sTestgridData struct {
	AllRepos     []string
	OrgsAndRepos map[string][]string
}

func generateK8sTestgrid(metaData TestGridMetaData) {
	// Regex expression for `knative-0.21`, `knative-sandbox-1.00`
	reReleaseBranch := regexp.MustCompile(`(knative|knative\-sandbox|google)\-[\d]+\.[\d]+`)

	allReposSet := sets.NewString("name: utilities")
	// Sort orgsAndRepos to maintain the output order
	allOrgs := []string{"maintenance", "prow-tests"}
	for org := range metaData.md {
		allOrgs = append(allOrgs, org)
	}
	sort.Strings(allOrgs)
	orgsAndRepos := map[string][]string{
		"maintenance": {"utilities"},
	}
	for org, repos := range metaData.md {
		// If org name matches release branch then this is a ungrouped
		if reReleaseBranch.MatchString(org) {
			allReposSet.Insert("name: " + org)
			continue
		}
		renamedReposForOrg := []string{}
		for repo := range repos {
			allReposSet.Insert("name: " + repo)
			if repo == "utilities" {
				continue
			}
			renamedReposForOrg = append(renamedReposForOrg, repo)
		}
		orgsAndRepos[org] = renamedReposForOrg
	}
	allRepos := allReposSet.List() // Returns in sorted order.

	executeTemplate("k8s testgrid",
		readTemplate(k8sTestgridTempl),
		struct{ AllRepos []string }{allRepos})

	for _, org := range allOrgs {
		repos := orgsAndRepos[org]
		sort.Strings(repos)
		groupName := org
		// If group name matches release branch then skip
		if reReleaseBranch.MatchString(groupName) {
			continue
		}
		executeTemplate("k8s testgrid group",
			readTemplate(k8sTestgridGroupTempl),
			struct {
				Org   string
				Repos []string
			}{groupName, repos})
	}
}
