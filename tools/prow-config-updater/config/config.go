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

package config

import "path/filepath"

type ProwEnv string

const (
	ProdProwEnv    ProwEnv = "prow"
	StagingProwEnv ProwEnv = "prow-staging"
)

const (
	OrgName        = "knative"
	RepoName       = "test-infra"
	PRHead         = "autoupdate"
	PRBase         = "master"
	// the label that needs to be applied on the PR to get it automatically merged
	AutoMergeLabel = "auto-merge"

	configPath         = "config"
	configTemplatePath = "config-generator/templates"
)

var (
	ProdProwConfigRoot    = filepath.Join(configPath, string(ProdProwEnv))
	StagingProwConfigRoot = filepath.Join(configPath, string(StagingProwEnv))

	// Config paths that need to be handled by prow-config-updater if files under them are changed.
	ProdProwConfigPaths = []string{
		filepath.Join(ProdProwConfigRoot, "core"),
		filepath.Join(ProdProwConfigRoot, "jobs"),
		filepath.Join(ProdProwConfigRoot, "cluster"),
	}
	StagingProwConfigPaths = []string{
		filepath.Join(StagingProwConfigRoot, "core"),
		filepath.Join(StagingProwConfigRoot, "jobs"),
		filepath.Join(StagingProwConfigRoot, "cluster"),
	}
	ProdTestgridConfigPath = filepath.Join(ProdProwConfigRoot, "testgrid")
)

// Config paths that need to be gated and tested by prow-config-updater.
var (
	ProdProwKeyConfigPaths = []string{
		filepath.Join(ProdProwConfigRoot, "cluster"),
		filepath.Join(configTemplatePath, string(ProdProwEnv)),
	}
	StagingProwKeyConfigPaths = []string{
		filepath.Join(StagingProwConfigRoot, "cluster"),
		filepath.Join(configTemplatePath, string(StagingProwEnv)),
	}
)
