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

package main

import (
	"knative.dev/test-infra/pkg/ghutil"
)

const (
	org  = "knative"
	repo = "test-infra"
	// PRHead is branch name where the changes occur
	PRHead = "releasebranch"
	// PRBase is the branch name where PR targets
	PRBase = "main"

	// Paths
	repoPath           = "src/knative.dev/test-infra"
	coreConfigPath     = "config/prow/core/config.yaml"
	jobConfigPath      = "config/prow/jobs/config.yaml"
	pluginPath         = "config/prow/core/plugins.yaml"
	testgridConfigPath = "config/prow/testgrid/testgrid.yaml"
	templateConfigPath = "config/prow/config_knative.yaml"

	configGenPath = "tools/config-generator"

	configGenScript = "hack/generate-configs.sh"

	oncallAddress = "https://storage.googleapis.com/knative-infra-oncall/oncall.json"
)

// GHClientWrapper handles methods for github issues
type GHClientWrapper struct {
	ghutil.GithubOperations
}
