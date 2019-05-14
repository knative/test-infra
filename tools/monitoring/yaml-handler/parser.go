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

// The yamlHandler is responsible for fetching, parsing config yaml file. It also allows user to
// retrieve a particular record from the yaml.

package yamlHandler

import (
	"fmt"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
)

type AlertConditions struct {
	JobNameRegex string `yaml:"job-name-regex"`
	Occurrences  int
	JobsAffected string `yaml:"jobs-affected"`
	PrsAffected  string `yaml:"prs-affected"`
	Period       int
}

type PatternSpec struct {
	ErrorPattern string `yaml:"error-pattern"`
	Hint         string
	Alerts       []AlertConditions
}

type YamlFile struct {
	Spec []PatternSpec `yaml:"spec"`
}

type Output struct {
	Hint         string
	Occurrences  int
	JobsAffected string
	PrsAffected  string
	Period       int
}

//GetSpec gets the spec for a particular error pattern and a matching job name pattern
func GetSpec(f YamlFile, pattern, jobName string) (output Output, noMatchError error) {
	noMatchError = errors.New(fmt.Sprintf("No spec found for pattern[%s] and jobName[%s]", pattern, jobName))
	for _, patternSpec := range f.Spec {
		if pattern == patternSpec.ErrorPattern {
			noMatchError = errors.New(fmt.Sprintf("Spec found for pattern[%s], but no match for job name[%s]", pattern, jobName))
			output.Hint = patternSpec.Hint
			for _, alertCondition := range patternSpec.Alerts {
				matched, err := regexp.MatchString(alertCondition.JobNameRegex, jobName)
				if err != nil {
					continue
				}
				if matched {
					noMatchError = nil
					output.JobsAffected = alertCondition.JobsAffected
					output.Occurrences = alertCondition.Occurrences
					output.PrsAffected = alertCondition.PrsAffected
					output.Period = alertCondition.Period
					break
				}
			}
			break
		}
	}
	return output, noMatchError
}

func getFileBytes(url string) []byte {
	res, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	content, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		log.Fatal(err)
	}
	return content
}

//ParseYaml reads the yaml text and convert it to the YamlFile struct defined
func ParseYaml(url string) YamlFile {
	content := getFileBytes(url)
	file := YamlFile{}

	if err := yaml.Unmarshal(content, &file); err != nil {
		log.Fatalf("Cannot parse config %q: %v", url, err)
	}
	return file
}

//CollectErrorPatterns collects and returns all error patterns in the yaml file
func CollectErrorPatterns(f YamlFile) []string {
	var patterns []string
	for _, patternSpec := range f.Spec {
		patterns = append(patterns, patternSpec.ErrorPattern)
	}
	return patterns
}
