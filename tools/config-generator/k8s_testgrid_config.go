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
	"sort"
)

const (
	k8sTestgridTempl = "k8s_testgrid.yaml"
)

type k8sTestgridData struct {
	NamedDashboards          []string
	KnativeDashboards        []string
	KnativeSandboxDashboards []string
	GoogleDashboards         []string
}

func generateK8sTestgrid(knativeDashboards, sandboxDashboards, googleDashboards []string) {
	namedDashboardsSet := make(map[string]struct{})
	for _, dashboard := range knativeDashboards {
		namedDashboardsSet["name: "+dashboard] = struct{}{}
	}
	for _, dashboard := range sandboxDashboards {
		namedDashboardsSet["name: "+dashboard] = struct{}{}
	}
	for _, dashboard := range googleDashboards {
		namedDashboardsSet["name: "+dashboard] = struct{}{}
	}
	namedDashboards := stringSetToSlice(namedDashboardsSet)
	sort.Strings(knativeDashboards)
	sort.Strings(sandboxDashboards)
	sort.Strings(googleDashboards)
	sort.Strings(namedDashboards)
	data := k8sTestgridData{
		KnativeDashboards:        knativeDashboards,
		KnativeSandboxDashboards: sandboxDashboards,
		GoogleDashboards:         googleDashboards,
		NamedDashboards:          namedDashboards,
	}
	executeTemplate("k8s testgrid", readTemplate(k8sTestgridTempl), data)
}
