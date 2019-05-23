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

import (
	"io/ioutil"
	"istio.io/fortio/log"
	"testing"
)

// TODO: after actual config yaml created, replace the path to point to that file
const yamlPath = "sample.yaml"

func getConfigYaml() []byte {
	text, err := ioutil.ReadFile(yamlPath)
	if err != nil {
		log.Fatalf("cannot read file [%s]: %v", yamlPath, err)
		return nil
	}
	return text
}

func Test_validate(t *testing.T) {
	type args struct {
		text []byte
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			// this test is to test the validity of the actual config yaml.
			name: "actual config yaml",
			args: args{
				text: getConfigYaml(),
			},
			wantErr: false,
		},
		{
			name: "valid config yaml",
			args: args{
				text: []byte(`spec:
  - error-pattern: 'Something went wrong: starting e2e cluster: error creating cluster'
    hint: 'Check gcp status'
    alerts:
      - job-name-regex: 'pull.*'
        occurrences: 2
        jobs-affected: 1
        prs-affected: 2
        period: 1440 # 1440 minutes = 24 hours
      - job-name-regex: '.*'
        occurrences: 5
        jobs-affected: 2
        prs-affected: 1 # for non-pull jobs, we don't care about the number of prs affected, so we set the number to 1, which will basically make this particular condition always true
        period: 1440

  - error-pattern: 'sample*error2'
    hint: 'hint_for_pattern_2'
    alerts:
      - job-name-regex: 'pull.*'
        occurrences: 20
        jobs-affected: 10
        prs-affected: 20
        period: 60 # 1440 minutes = 24 hours
      - job-name-regex: '.*'
        occurrences: 50
        jobs-affected: 20
        prs-affected: 10 # for non-pull jobs, we don't care about the number of prs affected, so we set the number to 1, which will basically make this particular condition always true
        period: 60`),
			},
			wantErr: false,
		},
		{
			name: "bad yaml - bad error pattern",
			args: args{
				text: []byte(`spec:
  - error-pattern: 'Something went wrong: starting e2e cluster: error creating cluster'
    hint: 'Check gcp status'
    alerts:
      - job-name-regex: '^pull.*'
        occurrences: 2
        jobs-affected: 1
        prs-affected: 2
        period: 1440 # 1440 minutes = 24 hours
      - job-name-regex: '.*'
        occurrences: 5
        jobs-affected: 2
        prs-affected: 1 # for non-pull jobs, we don't care about the number of prs affected, so we set the number to 1, which will basically make this particular condition always true
        period: 1440

  - error-pattern: '[3'
    hint: 'hint_for_pattern_2'
    alerts:
      - job-name-regex: '^pull.*'
        occurrences: 20
        jobs-affected: 10
        prs-affected: 20
        period: 60 # 1440 minutes = 24 hours
      - job-name-regex: '.*'
        occurrences: 50
        jobs-affected: 20
        prs-affected: 10 # for non-pull jobs, we don't care about the number of prs affected, so we set the number to 1, which will basically make this particular condition always true
        period: 60`),
			},
			wantErr: true,
		},
		{
			name: "bad yaml - bad job name pattern",
			args: args{
				text: []byte(`spec:
  - error-pattern: 'Something went wrong: starting e2e cluster: error creating cluster'
    hint: 'Check gcp status'
    alerts:
      - job-name-regex: '^pull.*'
        occurrences: 2
        jobs-affected: 1
        prs-affected: 2
        period: 1440 # 1440 minutes = 24 hours
      - job-name-regex: '.*'
        occurrences: 5
        jobs-affected: 2
        prs-affected: 1 # for non-pull jobs, we don't care about the number of prs affected, so we set the number to 1, which will basically make this particular condition always true
        period: 1440

  - error-pattern: '3'
    hint: 'hint_for_pattern_2'
    alerts:
      - job-name-regex: '^pull.*'
        occurrences: 20
        jobs-affected: 10
        prs-affected: 20
        period: 60 # 1440 minutes = 24 hours
      - job-name-regex: '[5'
        occurrences: 50
        jobs-affected: 20
        prs-affected: 10 # for non-pull jobs, we don't care about the number of prs affected, so we set the number to 1, which will basically make this particular condition always true
        period: 60`),
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := validate(tt.args.text); (err != nil) != tt.wantErr {
				t.Errorf("validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
