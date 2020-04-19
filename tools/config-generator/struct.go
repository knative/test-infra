package main

import (
	"log"

	"github.com/imdario/mergo"
)

type allConfig struct {
	// Toplevels holds repo level configs
	Toplevels []toplevel `yaml:"toplevels,omitempty"`
	// Presubmit holds presubmit configs for each repo(key of map)
	Presubmits []repoConfig `yaml:"presubmits,omitempty"`
	// Postsubmit holds presubmit configs for each repo(key of map)
	Postsubmits []repoConfig `yaml:"postsubmits,omitempty"`
	// Periodics holds periodic configs for each repo(key of map)
	Periodics []repoConfig `yaml:"periodics,omitempty"`
}

// Go coverage is pretty special, enabling it in presubmit means go coverage
// teastlso in postsubmit, periodic
func (ac *allConfig) init() {
	// First, merging prowJob configs with defaults set with "*" or higher level
	// config with same repo name
	ac.overrideConfig()

	// Second, perform custom logic go coverage logic
	m := make(map[string]prowJob)
	for _, rc := range ac.Presubmits {
		for _, pj := range rc.Jobs {
			if pj.Type == "go-coverage" {
				m[rc.Repo] = pj
			}
		}
	}

	for _, rc := range ac.Postsubmits {
		if pj, ok := m[rc.Repo]; ok {
			rc.Jobs = append(rc.Jobs, pj)
			// ac.Postsubmits[i] = append(ac.Postsubmits[i], pj)
		}
	}

	for _, rc := range ac.Periodics {
		if pj, ok := m[rc.Repo]; ok {
			rc.Jobs = append(rc.Jobs, pj)
			// ac.Postsubmits[i] = append(ac.Postsubmits[i], pj)
		}
	}

	// for i, rc := range ac.Periodics {
	// 	if pj, ok := m[rc.Repo]; ok {
	// 		ac.Periodics[i] = append(ac.Periodics[i], pj)
	// 	}
	// }

	/* Below are the implementation of supporting top level definitions */
	// var toplevelOverride func([]toplevel, []repoConfig) []repoConfig
	// toplevelOverride = func(tls []toplevel, rcs []repoConfig) []repoConfig {
	// 	res := rcs[:]
	// 	for _, tl := range ac.Toplevels {
	// 		if tl.GoCoverage {
	// 			// Add repo if not exist
	// 			targetIndex := -1
	// 			for i, rc := range rcs {
	// 				if rc.Repo == tl.Repo {
	// 					targetIndex = i
	// 					break
	// 				}
	// 			}
	// 			if targetIndex == -1 {
	// 				res = append(res, repoConfig{Repo: tl.Repo})
	// 				targetIndex = len(res) - 1
	// 			}

	// 			// Add job if not exist
	// 			newJobs := res[targetIndex].Jobs

	// 			var hasCov bool
	// 			for _, pj := range res[targetIndex].Jobs {
	// 				if pj.Type == "go-coverage" {
	// 					hasCov = true
	// 				}
	// 			}
	// 			if !hasCov {
	// 				newJobs = append(newJobs, prowJob{Type: "go-coverage"})
	// 			}

	// 			res[targetIndex].Jobs = newJobs
	// 		}
	// 	}
	// 	return res
	// }
	// ac.Presubmits = toplevelOverride(ac.Toplevels, ac.Presubmits)
	// ac.Postsubmits = toplevelOverride(ac.Toplevels, ac.Postsubmits)
	// ac.Periodics = toplevelOverride(ac.Toplevels, ac.Periodics)
}

func (ac *allConfig) overrideConfig() {
	// overrideFunc overrides configs, overrides if a value of the struct is
	// default value.
	// There is a glitch that false bool might not be able to override true, if
	// it comes needs may need to convert these to pointers instead of relying
	// on defaults. Or, change it to string representing bools
	var overrideFunc func([]repoConfig, string, *prowJob) prowJob
	overrideFunc = func(rcs []repoConfig, repo string, pj *prowJob) prowJob {
		for _, rc := range rcs {
			if rc.Repo == repo {
				for _, job := range rc.Jobs {
					if job.Type == "*" {
						if err := mergo.Merge(pj, job); err != nil {
							log.Fatalf("failed overriding: %v", err)
						}
					}
				}
			}
		}
		for _, rc := range rcs {
			if rc.Repo == "*" {
				for _, job := range rc.Jobs {
					if job.Type == pj.Type {
						if err := mergo.Merge(pj, job); err != nil {
							log.Fatalf("failed overriding: %v", err)
						}
					}
				}
			}
		}
		for _, rc := range rcs {
			if rc.Repo == "*" {
				for _, job := range rc.Jobs {
					if job.Type == "*" {
						if err := mergo.Merge(pj, job); err != nil {
							log.Fatalf("failed overriding: %v", err)
						}
					}
				}
			}
		}
		return *pj
	}

	var validJobs func(validJobs []repoConfig) []repoConfig
	validJobs = func(rcs []repoConfig) []repoConfig {
		var res []repoConfig
		for _, rc := range rcs {
			if rc.Repo == "*" {
				continue
			}
			var jobs []prowJob
			for _, pj := range rc.Jobs {
				if pj.Type == "*" {
					continue
				}
				jobs = append(jobs, overrideFunc(rcs, rc.Repo, &pj))
			}
			rc.Jobs = jobs
			res = append(res, rc)
		}
		return res
	}

	ac.Presubmits = validJobs(ac.Presubmits)
	ac.Postsubmits = validJobs(ac.Postsubmits)
	ac.Periodics = validJobs(ac.Periodics)
}

type toplevel struct {
	Repo       string `yaml:"repo"`
	GoCoverage bool   `yaml:"go-coverage,omitempty"`
}

type repoConfig struct {
	Repo string    `yaml:"repo"`
	Jobs []prowJob `yaml:"jobs,omitempty"`
}

type prowJob struct {
	// Template is the template to use for the job
	Template string `yaml:"template,omitempty"`
	// Type is job type, i.e. build-tests, custom-test
	Type string `yaml:"type"`
	// Name is the identifier of job name, for example:
	//	Name `upgrade-tests` in serving repo means `pull-kantive-serving-upgrade-tests`
	Name string `yaml:"name,omitempty"`
	// Skipped can be used to negate global default
	Skipped bool `yaml:"skipped,omitempty"`
	// DotDev: is this repo using knative.dev alias
	DotDev    bool `yaml:"dot-dev,omitempty"`
	AlwaysRun bool `yaml:"always-run,omitempty"`
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
