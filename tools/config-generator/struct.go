package main

type allConfig struct {
	// Presubmit holds presubmit configs for each repo(key of map)
	Presubmits map[string][]prowJob `yaml:"presubmits,omitempty"`
	// Periodics holds periodic configs for each repo(key of map)
	Periodics map[string][]prowJob `yaml:"periodics,omitempty"`
}

type prowJob struct {
	// Type is job type, i.e. build-tests, custom-test
	Type string `yaml:"type"`
	// Name is the identifier of job name, for example:
	//	Name `upgrade-tests` in serving repo means `pull-kantive-serving-upgrade-tests`
	Name string `yaml:"name,omitempty"`
	// Skipped can be used to negate global default
	Skipped bool `yaml:"skipped,omitempty"`
	// DotDev: is this repo using knative.dev alias
	DotDev    bool `yaml:"dot-dev,omitempty"`
	AlwaysRun bool `yaml:"always_run,omitempty"`
	Optional  bool `yaml:"optional,omitempty"`
	// NeedsMonitor specifies if crier sends it's pubsub message to GCP
	NeedsMonitor bool `yaml:"needs-monitor,omitempty"`
	// Resources contains CPU and memory requests/limits for a job
	Resources resources `yaml:"resources,omitempty"`
	// Command overrides the command
	Command string `yaml:"command,omitempty"`
	// Args
	Args []string `yaml:"args,omitempty"`
	// EnvVars
	EnvVars      []string `yaml:"env-vars,omitempty"`
	Timeout      int      `yaml:"timeout,omitempty"`
	NeedsDind    bool     `yaml:"needs-dind,omitempty"`
	Performance  bool     `yaml:"performance,omitempty"`
	SkipBranches []string `yaml:"skip-branches,omitempty"`
	Branches     []string `yaml:"branches,omitempty"`
	Cron         string   `yaml:"cron,omitempty"`

	/* For release branches */
	Release string `yaml:"release,omitempty"`

	/* Coverage only */
	// GoCoverageThreshold is the threshold for coverage test
	GoCoverageThreshold int `yaml:"go-coverage-threshold,omitempty"`

	/* Temporary */
	Go113         bool     `yaml:"go113,omitempty"`
	Go114         bool     `yaml:"go114,omitempty"`
	Go112Branches []string `yaml:"go112-branches,omitempty"`
}

type resources struct {
	Requests resourcesDef `yaml:"requests,omitempty"`
	Limits   resourcesDef `yaml:"limits,omitempty"`
}

type resourcesDef struct {
	Memory string `yaml:"memory,omitempty"`
	CPU    string `yaml:"cpu,omitempty"`
}
