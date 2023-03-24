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

package testgrid

import (
	"fmt"
	"io/ioutil"

	"sigs.k8s.io/yaml"
)

// Config is entire testgrid config
type Config struct {
	Dashboards      []Dashboard      `json:"dashboards"`
	DashboardGroups []DashboardGroup `json:"dashboard_groups"`
}

// DashboardGroup is a group of dashboards on testgrid
type DashboardGroup struct {
	// The name for the dashboard group.
	Name string `json:"name"`
	// A list of names specifying dashboards to show links to in a separate tabbed
	// bar at the top of the page for each of the given dashboards.
	DashboardNames []string `json:"dashboard_names"`
}

// Dashboard is single dashboard on testgrid
type Dashboard struct {
	Name         string          `json:"name"`
	DashboardTab []*DashboardTab `json:"dashboard_tab,omitempty"`
}

// DashboardTab is a single tab on testgrid
type DashboardTab struct {
	Name          string `json:"name"`
	TestGroupName string `json:"test_group_name"`
}

// NewConfigFromFile loads config from file
func NewConfigFromFile(fp string) (*Config, error) {
	ac := &Config{}
	contents, err := ioutil.ReadFile(fp)
	if err == nil {
		err = yaml.Unmarshal(contents, ac)
	}
	if err != nil {
		return nil, err
	}
	return ac, err
}

// GetTabRelURL finds URL relative to testgrid home URL from testgroup name
// (generally this is prow job name)
func (ac *Config) GetTabRelURL(tgName string) (string, error) {
	for _, dashboard := range ac.Dashboards {
		for _, tab := range dashboard.DashboardTab {
			if tab.TestGroupName == tgName {
				return fmt.Sprintf("%s#%s", dashboard.Name, tab.Name), nil
			}
		}
	}
	return "", fmt.Errorf("testgroup name '%s' not exist", tgName)
}
