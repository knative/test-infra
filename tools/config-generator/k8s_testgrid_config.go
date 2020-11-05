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

import "sort"

const (
	k8sTestgridTempl      = "k8s_testgrid.yaml"
	k8sTestgridGroupTempl = "k8s_testgrid_testgroup.yaml"
)

type k8sTestgridData struct {
	AllRepos     []string
	OrgsAndRepos map[string][]string
}

func generateK8sTestgrid(orgsAndRepos map[string][]string) {
	allReposSet := make(map[string]struct{})
	for _, repos := range orgsAndRepos {
		for _, repo := range repos {
			allReposSet["name: "+repo] = struct{}{}
		}
	}
	allRepos := stringSetToSlice(allReposSet)
	sort.Strings(allRepos)

	executeTemplate("k8s testgrid",
		readTemplate(k8sTestgridTempl),
		struct{ AllRepos []string }{allRepos})

	for org, repos := range orgsAndRepos {
		sort.Strings(repos)
		executeTemplate("k8s testgrid group",
			readTemplate(k8sTestgridGroupTempl),
			struct {
				Org   string
				Repos []string
			}{org, repos})
	}
}
