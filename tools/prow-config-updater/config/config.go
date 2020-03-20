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

import (
	"fmt"
	"log"
	"path/filepath"

	"knative.dev/pkg/test/cmd"
	"knative.dev/pkg/test/helpers"
)

type ProwEnv string

const (
	ProdProwEnv    ProwEnv = "prow"
	StagingProwEnv ProwEnv = "prow-staging"
)

const (
	OrgName  = "knative"
	RepoName = "test-infra"
	PRHead   = "autoupdate"
	PRBase   = "master"
	// the label that needs to be applied on the PR to get it automatically merged
	AutoMergeLabel = "auto-merge"

	PullBotName          = "pull[bot]"
	PullEndpointTemplate = "https://pull.git.ci/process/%s/%s"

	configPath         = "config"
	configTemplatePath = "tools/config-generator/templates"

	// Prow config subfolder names
	core    = "core"
	jobs    = "jobs"
	cluster = "cluster"
)

var (
	// Prow config root paths.
	prodProwConfigRoot    = filepath.Join(configPath, string(ProdProwEnv))
	stagingProwConfigRoot = filepath.Join(configPath, string(StagingProwEnv))

	// Commands that generate and update Prow configs.
	// These are commands for both staging and production Prow.
	generateProwConfigFilesCommand = "./hack/generate-configs.sh"
	updateProwCommandTemplate      = "make -C %s update-prow-cluster"
	updateProdProwCommand          = fmt.Sprintf(updateProwCommandTemplate, prodProwConfigRoot)
	updateStagingProwCommand       = fmt.Sprintf(updateProwCommandTemplate, stagingProwConfigRoot)
	// This command is only used for production prow in this tool.
	updateTestgridCommand = fmt.Sprintf("make -C %s update-testgrid-config", prodProwConfigRoot)

	// Config paths that need to be handled by prow-config-updater if files under them are changed.

	ProdProwConfigPaths = []string{
		filepath.Join(prodProwConfigRoot, core),
		filepath.Join(prodProwConfigRoot, jobs),
		filepath.Join(prodProwConfigRoot, cluster),
	}
	StagingProwConfigPaths = []string{
		filepath.Join(stagingProwConfigRoot, core),
		filepath.Join(stagingProwConfigRoot, jobs),
		filepath.Join(stagingProwConfigRoot, cluster),
	}
	ProdTestgridConfigPath = filepath.Join(prodProwConfigRoot, "testgrid")

	// Config paths that need to be gated and tested by prow-config-updater.
	ProdProwKeyConfigPaths = []string{
		filepath.Join(prodProwConfigRoot, cluster),
		filepath.Join(configTemplatePath, string(ProdProwEnv)),
	}
	StagingProwKeyConfigPaths = []string{
		filepath.Join(stagingProwConfigRoot, cluster),
		filepath.Join(configTemplatePath, string(StagingProwEnv)),
	}
)

// UpdateProw will update Prow with the existing Prow config files.
func UpdateProw(env ProwEnv, dryrun bool) error {
	updateCommand := ""
	switch env {
	case ProdProwEnv:
		updateCommand = updateProdProwCommand
	case StagingProwEnv:
		updateCommand = updateStagingProwCommand
	default:
		return fmt.Errorf("unsupported Prow environement: %q, cannot make the update", env)
	}

	return helpers.Run(
		fmt.Sprintf("Updating Prow configs with command %q", updateCommand),
		func() error {
			out, err := cmd.RunCommand(updateCommand)
			log.Println(out)
			return err
		},
		dryrun,
	)
}

// UpdateTestgrid will update testgrid with the existing testgrid config file.
func UpdateTestgrid(env ProwEnv, dryrun bool) error {
	if env != ProdProwEnv {
		log.Printf("no testgrid config needs to be updated for %q Prow environment", env)
		return nil
	}

	return helpers.Run(
		fmt.Sprintf("Updating Testgrid config with command %q", updateTestgridCommand),
		func() error {
			out, err := cmd.RunCommand(updateTestgridCommand)
			log.Println(out)
			return err
		},
		dryrun,
	)
}

// GenerateProwConfigFiles will run the config generator command to generate new Prow config files.
func GenerateProwConfigFiles() error {
	_, err := cmd.RunCommand(generateProwConfigFilesCommand)
	return err
}
