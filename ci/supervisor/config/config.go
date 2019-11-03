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

// Package config contains functions related to config files.
package config

import (
	"fmt"
	yaml "gopkg.in/yaml.v2"
)

// yamlMap is a generic type for YAML files
type yamlMap = map[interface{}]interface{}

// EnvVar holds a environment variable (name and value)
type EnvVar struct {
	Name, Value string
}

// ResourcesRequest holds resource requests for a task or job
type ResourcesRequest struct {
	CPU    int
	Memory string
}

// ResourcesConfig holds the config of resources used by a task or job
type ResourcesConfig struct {
	Requests ResourcesRequest
}

// DefaultSettings holds the default settings to be used by the supervisor
type DefaultSettings struct {
	MergeMethod       string `yaml:"merge-method"`
	Optional          bool
	ParseGoTestOutput bool   `yaml:"parse-go-test-output"`
	ReportToGitHub    bool   `yaml:"report-to-github"`
	EnableDocker      bool   `yaml:"enable-docker"`
	PresubmitFilter   string `yaml:"presumit-filter"`
	Image             string
	ImageArgs         []string `yaml:"image-args"`
	Env               []EnvVar
	Resources         ResourcesConfig
	TestGridConfig    string `yaml:"testgrid-config"`
}

// CodeCoverageSettings holds the settings for the code coverage task
type CodeCoverageSettings struct {
	Threshold      string
	FailsPresubmit bool `yaml:"fails-presubmit"`
	Image          string
	ImageArgs      []string `yaml:"image-args"`
	Env            []EnvVar
	Resources      ResourcesConfig
}

// Task holds the basic information for a task
type Task struct {
	Command           string
	ParseGoTestOutput bool `yaml:"parse-go-test-output"`
	EnableDocker      bool `yaml:"enable-docker"`
	Image             string
	ImageArgs         []string `yaml:"image-args"`
	Env               []EnvVar
	Resources         ResourcesConfig
}

// ProwPlugin holds the parameters of a Prow plugin
type ProwPlugin struct {
	Enabled    bool
	Parameters yamlMap
}

// PresubmitTask holds the parameters of a presubmit task
type PresubmitTask struct {
	TaskParameters  Task
	PresubmitFilter string `yaml:"presubmit-filter"`
	Optional        bool
	ReportToGitHub  bool `yaml:"report-to-github"`
}

// PostsubmitTask holds the parameters of a postsubmit task
type PostsubmitTask struct {
	TaskParameters Task
}

// PeriodicTask holds the parameters or a periodic task
type PeriodicTask struct {
	TaskParameters          Task
	Periodicity             string
	MinorReleasePeriodicity string `yaml:"minor-release-periodicity"`
	MinorReleaseArg         string `yaml:"minor-release-arg"`
}

// Defaults holds the default settings for a repository
type Defaults struct {
	Settings     DefaultSettings
	ProwPlugins  map[string]ProwPlugin
	CodeCoverage CodeCoverageSettings
}

// RepoConfig holds all the configuration for a repository
type RepoConfig struct {
	Defaults        Defaults
	PresubmitTasks  map[string]PresubmitTask
	PostsubmitTasks map[string]PostsubmitTask
	PeriodicTasks   map[string]PeriodicTask
}

// NightlyReleaseSettings holds the default configuration of the nightly release task
type NightlyReleaseSettings struct {
	Window string
}

// AutoReleaseSettings holds the default configuration of the auto release task
type AutoReleaseSettings struct {
	Window                  string
	Periodicity             string
	MinorReleasePeriodicity string `yaml:"minor-release-periodicity"`
	MinorReleaseArg         string `yaml:"minor-release-arg"`
}

// SupervisorDefaults all the default configuration for the repository tasks
type SupervisorDefaults struct {
	Settings       DefaultSettings
	ProwPlugins    map[string]ProwPlugin
	CodeCoverage   CodeCoverageSettings
	NightlyRelease NightlyReleaseSettings
	AutoRelease    AutoReleaseSettings
}

// SupervisorConfig holds all the configuration for supervisor
type SupervisorConfig struct {
	Defaults     SupervisorDefaults
	Repositories map[string][]string
}

// parsePlugins unmarshals a Prow plugin section into the given map
func parsePlugins(pluginMap map[string]ProwPlugin, rawData yamlMap) error {
	for p := range rawData {
		pluginName, _ := p.(string)
		plugin := ProwPlugin{Enabled: true}
		badSetting := fmt.Errorf("valid setting for plugin %q is one of [disabled|no|false|<configuration>]", pluginName)
		// Setting is a simple enable/disable string
		if v, ok := rawData[pluginName].(string); ok {
			if v == "disabled" || v == "no" || v == "false" {
				plugin.Enabled = false
			} else if v == "enabled" || v == "yes" || v == "true" {
				plugin.Enabled = true
			} else {
				return badSetting
			}
			// Otherwise it's a raw plugin setting
		} else {
			if p, ok := rawData[pluginName].(yamlMap); ok {
				plugin.Parameters = p
			} else {
				return badSetting
			}
		}
		pluginMap[pluginName] = plugin
	}
	return nil
}

// UnmarshalYAML unmarshals a Defaults struct, correctly unmarshalling the Prow plugin section
func (d *Defaults) UnmarshalYAML(unmarshal func(interface{}) error) error {
	// Expected fields
	var structuredData struct {
		Settings     DefaultSettings
		ProwPlugins  yamlMap              `yaml:"prow-plugins"`
		CodeCoverage CodeCoverageSettings `yaml:"code-coverage"`
	}
	if err := unmarshal(&structuredData); err != nil {
		return err
	}
	d.Settings = structuredData.Settings
	d.CodeCoverage = structuredData.CodeCoverage
	d.ProwPlugins = make(map[string]ProwPlugin, 0)
	return parsePlugins(d.ProwPlugins, structuredData.ProwPlugins)
}

// ParseRepoConfig unmarshals a repository config YAML
func ParseRepoConfig(content []byte) (RepoConfig, error) {
	r := RepoConfig{}
	r.PresubmitTasks = make(map[string]PresubmitTask, 0)
	r.PeriodicTasks = make(map[string]PeriodicTask, 0)
	r.PostsubmitTasks = make(map[string]PostsubmitTask, 0)

	// Extract the raw YAML
	var rawConfig map[string]interface{}
	if err := yaml.Unmarshal(content, &rawConfig); err != nil {
		return RepoConfig{}, err
	}

	// Unmarshal the defaults section
	if d, err := yaml.Marshal(rawConfig["defaults"]); err == nil {
		if err := yaml.Unmarshal(d, &r.Defaults); err != nil {
			return RepoConfig{}, err
		}
	} else {
		return RepoConfig{}, err
	}

	// Unmarshal the presubmit-tasks section
	if value, ok := rawConfig["presubmit-tasks"].(yamlMap); ok {
		for a := range value {
			taskName, _ := a.(string)
			// Setting can be simply the command to run
			if v, ok := value[taskName].(string); ok {
				r.PresubmitTasks[taskName] = PresubmitTask{TaskParameters: Task{Command: v}}
				// Otherwise it's the complete set of parameters
			} else {
				if d, err := yaml.Marshal(value[taskName]); err == nil {
					pt := PresubmitTask{}
					err = yaml.Unmarshal(d, &pt)
					err = yaml.Unmarshal(d, &pt.TaskParameters)
					r.PresubmitTasks[taskName] = pt
				} else {
					return RepoConfig{}, err
				}
			}
		}
	}

	// Unmarshal the periodic-tasks section
	if value, ok := rawConfig["periodic-tasks"].(yamlMap); ok {
		for a := range value {
			taskName, _ := a.(string)
			// Expect the complete set of parameters
			if d, err := yaml.Marshal(value[taskName]); err == nil {
				pt := PeriodicTask{}
				err = yaml.Unmarshal(d, &pt)
				err = yaml.Unmarshal(d, &pt.TaskParameters)
				r.PeriodicTasks[taskName] = pt
			} else {
				return RepoConfig{}, err
			}
		}
	}

	// Unmarshal the postsubmit-tasks section
	if value, ok := rawConfig["postsubmit-tasks"].(yamlMap); ok {
		for a := range value {
			taskName, _ := a.(string)
			// Setting can be simply the command to run
			if v, ok := value[taskName].(string); ok {
				r.PostsubmitTasks[taskName] = PostsubmitTask{TaskParameters: Task{Command: v}}
				// Otherwise it's the complete set of parameters
			} else {
				if d, err := yaml.Marshal(value[taskName]); err == nil {
					pt := PostsubmitTask{}
					err = yaml.Unmarshal(d, &pt)
					err = yaml.Unmarshal(d, &pt.TaskParameters)
					r.PostsubmitTasks[taskName] = pt
				} else {
					return RepoConfig{}, err
				}
			}
		}
	}

	return r, nil
}

// UnmarshalYAML unmarshals a SupervisorDefaults struct, correctly unmarshalling the Prow plugin section
func (sd *SupervisorDefaults) UnmarshalYAML(unmarshal func(interface{}) error) error {
	// Expected fields
	var structuredData struct {
		Settings       DefaultSettings
		ProwPlugins    yamlMap                `yaml:"prow-plugins"`
		CodeCoverage   CodeCoverageSettings   `yaml:"code-coverage"`
		NightlyRelease NightlyReleaseSettings `yaml:"nightly-release"`
		AutoRelease    AutoReleaseSettings    `yaml:"auto-release"`
	}
	if err := unmarshal(&structuredData); err != nil {
		return err
	}
	sd.Settings = structuredData.Settings
	sd.CodeCoverage = structuredData.CodeCoverage
	sd.NightlyRelease = structuredData.NightlyRelease
	sd.AutoRelease = structuredData.AutoRelease
	sd.ProwPlugins = make(map[string]ProwPlugin, 0)
	return parsePlugins(sd.ProwPlugins, structuredData.ProwPlugins)
}

// ParseSupervisorConfig unmarshals a supervisor config YAML
func ParseSupervisorConfig(content []byte) (SupervisorConfig, error) {
	s := SupervisorConfig{}
	s.Repositories = make(map[string][]string, 0)

	// Extract the raw YAML
	var rawConfig map[string]interface{}
	if err := yaml.Unmarshal(content, &rawConfig); err != nil {
		return SupervisorConfig{}, err
	}

	// Unmarshal the defaults section
	if d, err := yaml.Marshal(rawConfig["defaults"]); err == nil {
		if err := yaml.Unmarshal(d, &s.Defaults); err != nil {
			return SupervisorConfig{}, err
		}
	} else {
		return SupervisorConfig{}, err
	}

	// Unmarshal the repos section
	if d, err := yaml.Marshal(rawConfig["repos"]); err == nil {
		if err := yaml.Unmarshal(d, &s.Repositories); err != nil {
			return SupervisorConfig{}, err
		}
	} else {
		return SupervisorConfig{}, err
	}

	return s, nil
}
