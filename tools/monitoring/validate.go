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

package main

import "fmt"

// collectPatterns collects all regexp patterns, including both error message patterns
// and job name patterns
func collectPatterns(config *Config) []string {
	var patterns []string
	for _, patternSpec := range config.Spec {
		patterns = append(patterns, patternSpec.ErrorPattern)
		for _, alertCondition := range patternSpec.Alerts {
			patterns = append(patterns, alertCondition.JobNameRegex)
		}
	}

	return patterns
}

// validate checks if the given yaml content meets our requirement for a monitoring config file
func validate(text []byte) error {
	config, err := newConfig(text)
	if err != nil {
		return err
	}
	patterns := collectPatterns(config)
	_, badPatterns := compilePatterns(patterns)
	if len(badPatterns) > 0 {
		return fmt.Errorf("bad patterns found: %v", badPatterns)
	}
	return nil
}
