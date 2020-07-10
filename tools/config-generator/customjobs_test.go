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
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"gopkg.in/yaml.v2"
	"k8s.io/apimachinery/pkg/util/sets"
)

var (
	defaultTemplateConfigPath = "../../config/prod/prow/jobs/custom"
)

type customJobStruct struct {
	Presubmits  map[string][]singleCustomJob `yaml:"presubmits,omitempty"`
	Postsubmits map[string][]singleCustomJob `yaml:"postsubmits,omitempty"`
	Periodics   []singleCustomJob            `yaml:"periodics,omitempty"`
}

type singleCustomJob struct {
	Name string `yaml:"name"`
}

func TestEnsureCustomJob(t *testing.T) {
	validJobs := sets.NewString()
	filepath.Walk(defaultTemplateConfigPath, func(path string, info os.FileInfo, err error) error {
		if strings.HasSuffix(path, ".yaml") {
			content, err := ioutil.ReadFile(path)
			if err != nil {
				t.Fatalf("Failed reading template file %q: %v", path, err)
			}

			allCustomJobs := customJobStruct{}
			if err = yaml.Unmarshal(content, &allCustomJobs); err != nil {
				t.Fatalf("Failed unmarshalling: %v", err)
			}
			for _, sjs := range allCustomJobs.Presubmits {
				for _, sj := range sjs {
					validJobs.Insert(sj.Name)
				}
			}
			for _, sjs := range allCustomJobs.Postsubmits {
				for _, sj := range sjs {
					validJobs.Insert(sj.Name)
				}
			}
			for _, sj := range allCustomJobs.Periodics {
				validJobs.Insert(sj.Name)
			}
		}
		return nil
	})

	for _, job := range customJobnames {
		if !validJobs.Has(job) {
			t.Fatalf("Job %q doesn't exist in %q", job, defaultTemplateConfigPath)
		}
	}
}
