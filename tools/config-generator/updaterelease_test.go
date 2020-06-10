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
	"errors"
	"io/ioutil"
	"os"
	"syscall"
	"testing"

	"github.com/google/go-cmp/cmp"
	"gopkg.in/yaml.v2"
)

var (
	// errUnwrappable: Some errors not wrappable
	errUnwrappable = errors.New("unwrappable")
)

func TestUpgradeReleaseBranchesTemplate(t *testing.T) {
	tests := []struct {
		name      string
		fileExist bool
		in        string
		want      string
		wantErr   error
	}{
		{
			"Change",
			true,
			`periodics:
  repo1:
  - branch-ci: true
    release: "0.5"`,
			`periodics:
  repo1:
  - branch-ci: true
    release: "0.5"
  - branch-ci: true
    release: "0.15"
`,
			nil,
		}, {
			"No_op",
			true,
			`periodics:
  repo1:
  - branch-ci: true
    release: "0.15"`,
			`periodics:
  repo1:
  - branch-ci: true
    release: "0.15"
`,
			nil,
		}, {
			"File_not_exit",
			false,
			`doesnt matter`,
			`doesnt matter`,
			syscall.ENOENT, // os.PathError is not usable, use syscall instead
		}, {
			"Not marshallable",
			true,
			`doesnt matter`,
			`doesnt matter`,
			errUnwrappable,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var fn string
			fn = "file_not_exist"
			if tt.fileExist {
				fi, err := ioutil.TempFile(os.TempDir(), "TestUpgradeReleaseBranchesTemplate")
				if err == nil {
					fn = fi.Name()
					err = ioutil.WriteFile(fi.Name(), []byte(tt.in), 0644)
				}
				if err != nil {
					t.Fatalf("Failed creating temp file: %v", err)
				}
				t.Logf("Temp file created at %q", fi.Name())
			}
			err := upgradeReleaseBranchesTemplate(fn)
			if !errors.Is(err, tt.wantErr) && (err != nil && tt.wantErr != errUnwrappable) {
				t.Fatalf("Error not expected. Want: '%v', got: '%v'", tt.wantErr, err)
			}
			if !tt.fileExist {
				return
			}
			gotBytes, err := ioutil.ReadFile(fn)
			if !errors.Is(err, tt.wantErr) && (err != nil && tt.wantErr != errUnwrappable) {
				t.Fatalf("Error not expected. Want: '%v', got: '%v'", tt.wantErr, err)
			}
			got := string(gotBytes)
			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Fatalf("Mismatch, got(+), want(-): \n%s", diff)
			}
		})
	}
}

func TestGetReposMap(t *testing.T) {
	const (
		latest = "0.6"
	)
	tests := []struct {
		name string
		in   string
		want string
	}{
		{
			"Simple_update_case",
			`repo1:
- branch-ci: true
  release: "0.5"`,
			`repo1:
- branch-ci: true
  release: "0.5"
- branch-ci: true
  release: "0.6"
`,
		}, {
			"Simple_update_case2",
			`repo1:
- branch-ci: true
  release: "0.1"`,
			`repo1:
- branch-ci: true
  release: "0.1"
- branch-ci: true
  release: "0.6"
`,
		}, {
			"Simple_update_case3",
			`repo1:
- branch-ci: true
  release: "0.1"
- branch-ci: true
  release: "0.3"`,
			`repo1:
- branch-ci: true
  release: "0.1"
- branch-ci: true
  release: "0.3"
- branch-ci: true
  release: "0.6"
`,
		}, {
			"Simple_update_case4",
			`repo1:
- dot-release: true
  release: "0.5"`,
			`repo1:
- dot-release: true
  release: "0.5"
- dot-release: true
  release: "0.6"
`,
		}, {
			"Delete_old_branches",
			`repo1:
- branch-ci: true
  release: "0.2"
- branch-ci: true
  release: "0.3"
- branch-ci: true
  release: "0.4"
- branch-ci: true
  release: "0.5"`,
			`repo1:
- branch-ci: true
  release: "0.3"
- branch-ci: true
  release: "0.4"
- branch-ci: true
  release: "0.5"
- branch-ci: true
  release: "0.6"
`,
		}, {
			"No_op",
			`repo1:
- branch-ci: true
  release: "0.3"
- branch-ci: true
  release: "0.4"
- branch-ci: true
  release: "0.5"
- branch-ci: true
  release: "0.6"`,
			`repo1:
- branch-ci: true
  release: "0.3"
- branch-ci: true
  release: "0.4"
- branch-ci: true
  release: "0.5"
- branch-ci: true
  release: "0.6"
`,
		}, {
			"No_delete",
			`repo1:
- branch-ci: true
  release: "0.2"
- branch-ci: true
  release: "0.3"
- branch-ci: true
  release: "0.4"
- branch-ci: true
  release: "0.5"
- branch-ci: true
  release: "0.6"`,
			`repo1:
- branch-ci: true
  release: "0.2"
- branch-ci: true
  release: "0.3"
- branch-ci: true
  release: "0.4"
- branch-ci: true
  release: "0.5"
- branch-ci: true
  release: "0.6"
`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			inStruct := yaml.MapSlice{}
			if err := yaml.Unmarshal([]byte(tt.in), &inStruct); err != nil {
				t.Fatalf("Failed unmarshal %q: %v", tt.in, err)
			}
			gotStruct := getReposMap(inStruct, latest)
			gotBytes, err := yaml.Marshal(gotStruct)
			if err != nil {
				t.Fatalf("Failed marshal: %v", err)
			}
			got := string(gotBytes)
			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Fatalf("Mismatch, got(+), want(-): \n%s", diff)
			}
		})
	}
}
