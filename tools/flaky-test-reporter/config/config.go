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

// config.go contains configurations for flaky tests reporting

package config

import (
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

const (
	defaultConfig = "config.yaml"
)

// Config contains all job configs for flaky tests reporting
type Config struct {
	JobConfigs []JobConfig `yaml:"jobConfigs"`
}

// JobConfig is initial configuration for a given repo, defines which job to scan
type JobConfig struct {
	Name            string         `yaml:"name"` // name of job to analyze
	Repo            string         `yaml:"repo"` // repository to test job on
	Type            string         `yaml:"type"`
	GithubIssueRepo string         `yaml:"githubIssueRepo,omitempty"`
	SlackChannels   []SlackChannel `yaml:"slackChannels,omitempty"`
}

// SlackChannel contains Slack channels info
type SlackChannel struct {
	Name     string `yaml:"name"`
	Identity string `yaml:"identity"`
}

// NewConfig parses config from configFile, default to "config.yaml" if configFile is empty string
func NewConfig(configFile string) (*Config, error) {
	if "" == configFile {
		configFile = defaultConfig
	}
	contents, err := ioutil.ReadFile(configFile)
	if nil != err {
		return nil, err
	}
	config := &Config{}
	if err = yaml.Unmarshal(contents, config); nil != err {
		return nil, err
	}
	return config, nil
}
