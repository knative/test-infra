package config

import "path/filepath"

type ProwEnv string

const (
	ProdProwEnv    ProwEnv = "prow"
	StagingProwEnv ProwEnv = "prow-staging"
)

const (
	ProwBotName = "knative-prow-robot"
	ForkOrgName = ProwBotName
	OrgName     = "knative"
	RepoName    = "test-infra"

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
	ProdTestgridConfigFile = filepath.Join(ProdProwConfigRoot, "testgrid/testgrid.yaml")
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
