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

package config

import (
	"fmt"
	"testing"

	"github.com/google/go-cmp/cmp"
)

// errorMismatch compares two errors by their error string
func errorMismatch(got error, want error) string {
	exp := "<no error>"
	act := "<no error>"
	if want != nil {
		exp = want.Error()
	}
	if got != nil {
		act = got.Error()
	}
	if exp == act {
		return ""
	}
	return fmt.Sprintf("got: '%v'\nwant: '%v'", act, exp)
}

func TestParseRepoConfig(t *testing.T) {
	datas := []struct {
		yaml string
		exp  RepoConfig
		err  error
	}{
		{ // empty YAML
			"",
			RepoConfig{
				PresubmitTasks:  map[string]PresubmitTask{},
				PeriodicTasks:   map[string]PeriodicTask{},
				PostsubmitTasks: map[string]PostsubmitTask{},
			},
			nil,
		}, { // bad YAML
			"<>",
			RepoConfig{},
			fmt.Errorf("yaml: unmarshal errors:\n  line 1: cannot unmarshal !!str `<>` into map[string]interface {}"),
		}, { // unknown defaults section
			`defaults:
  foobar:
    optional: true`,
			RepoConfig{
				Defaults: Defaults{
					Settings:     DefaultSettings{},
					ProwPlugins:  map[string]ProwPlugin{},
					CodeCoverage: CodeCoverageSettings{},
				},
				PresubmitTasks:  map[string]PresubmitTask{},
				PeriodicTasks:   map[string]PeriodicTask{},
				PostsubmitTasks: map[string]PostsubmitTask{},
			},
			nil,
		}, { // bad plugin setting
			`defaults:
  prow-plugins:
    plugin1: 0.0`,
			RepoConfig{},
			fmt.Errorf("valid setting for plugin \"plugin1\" is one of [disabled|no|false|<configuration>]"),
		}, { // bad plugin string setting
			`defaults:
  prow-plugins:
    plugin1: foo`,
			RepoConfig{},
			fmt.Errorf("valid setting for plugin \"plugin1\" is one of [disabled|no|false|<configuration>]"),
		}, { // presubmit task definition with unrecognized key
			`presubmit-tasks:
  foo-tests:
    thisisanunrezognizedkey: "foo"`,
			RepoConfig{
				PresubmitTasks:  map[string]PresubmitTask{"foo-tests": {}},
				PeriodicTasks:   map[string]PeriodicTask{},
				PostsubmitTasks: map[string]PostsubmitTask{},
			},
			nil,
		}, { // periodic task definition with unrecognized key
			`periodic-tasks:
  foo-tests:
    thisisanunrezognizedkey: "foo"`,
			RepoConfig{
				PresubmitTasks:  map[string]PresubmitTask{},
				PeriodicTasks:   map[string]PeriodicTask{"foo-tests": {}},
				PostsubmitTasks: map[string]PostsubmitTask{},
			},
			nil,
		}, { // postsubmit task definition with unrecognized key
			`postsubmit-tasks:
  foo-tests:
    thisisanunrezognizedkey: "foo"`,
			RepoConfig{
				PresubmitTasks:  map[string]PresubmitTask{},
				PeriodicTasks:   map[string]PeriodicTask{},
				PostsubmitTasks: map[string]PostsubmitTask{"foo-tests": {}},
			},
			nil,
		}, { // full YAML
			`defaults:
  settings:
    merge-method: "squash"
    optional: true
    parse-go-test-output: false
    report-to-github: false
    enable-docker: true
    presubmit-filter: "(?!\\.bak)$"
    image: "gcr.io/foo/bar/special-image:dev"
    image-args:
      - "--extra-verbose"
    env:
      - name: "MY_ENV_FOO"
        value: "foo1"
    resources:
      requests:
        cpu: 7
        memory: "41Gi"
  prow-plugins:
    plugin1: disabled
    plugin2: enabled
    plugin3:
      foo: bar
  code-coverage:
    threshold: 80%
    fails-presubmit: true
    image: "gcr.io/foo/bar/special-image:dev"
    image-args:
      - "--ignore-tests"
    env:
      - name: "MY_ENV_BAR"
        value: "bar"
    resources:
      requests:
        cpu: 1
        memory: "1Gi"
presubmit-tasks:
  build-tests: "test/presubmit-tests.sh --build-tests"
  unit-tests: "go test ./..."
  custom-tests:
    command: "test/presubmit-tests.sh --xlarge-tests"
    presubmit-filter: "(?!\\.bat)$"
    optional: true
    image: "gcr.io/foo/bar/special-image2:dev"
    image-args:
    - "--ram 16GB"
    env:
      - name: "MY_ENV_FOB"
        value: "fob"
    resources:
      requests:
        cpu: 2
        memory: "2Gi"
postsubmit-tasks:
  all-tests: "go test"
  custom-all-tests:
    command: "test/presubmit-tests.sh --all"
    image: "gcr.io/foo/bar/image"
    image-args:
      - "--ram 6GB"
    env:
      - name: "MY_ENV_BOB"
        value: "bob"
    resources:
      requests:
        cpu: 9
        memory: "99Gi"
periodic-tasks:
  large-ci-tests:
    command: "test/presubmit-tests.sh --xlarge-tests"
    periodicity: "daily 8AM"
    optional: true
    parse-go-test-output: false
    enable-docker: true
    report-to-github: false
    image: "gcr.io/foo/bar/special-image3:dev"
    image-args:
      - "--ram 1GB"
    env:
      - name: "MY_ENV_FOX"
        value: "fox"
    resources:
      requests:
        cpu: 3
        memory: "3Gi"
  nightly-release:
    command: "hack/release.py --nightly"
    periodicity: "daily 1PM"
    env:
      - name: "MY_ENV_FOZ"
        value: "foz"
  auto-release:
    command: "hack/release.py --publish"
    periodicity: "every 3h"
    minor-release-periodicity: "Tuesdays 3AM"
    minor-release-arg: "--dot-release"
    image: "gcr.io/foo/bar/special-image4:dev"
    image-args:
    - "--ram 4GB"
    env:
      - name: "MY_ENV_FOY"
        value: "foy"
    resources:
      requests:
        cpu: 4
        memory: "4Gi"`,
			RepoConfig{
				Defaults: Defaults{
					Settings: DefaultSettings{
						MergeMethod:       "squash",
						Optional:          true,
						ParseGoTestOutput: false,
						ReportToGitHub:    false,
						EnableDocker:      true,
						Image:             "gcr.io/foo/bar/special-image:dev",
						ImageArgs:         []string{"--extra-verbose"},
						Env:               []EnvVar{{Name: "MY_ENV_FOO", Value: "foo1"}},
						Resources:         ResourcesConfig{Requests: ResourcesRequest{CPU: 7, Memory: "41Gi"}},
					},
					ProwPlugins: map[string]ProwPlugin{
						"plugin1": {Enabled: false},
						"plugin2": {Enabled: true},
						"plugin3": {
							Enabled: true, Parameters: map[interface{}]interface{}{string("foo"): string("bar")},
						},
					},
					CodeCoverage: CodeCoverageSettings{
						Threshold:      "80%",
						FailsPresubmit: true,
						Image:          "gcr.io/foo/bar/special-image:dev",
						ImageArgs:      []string{"--ignore-tests"},
						Env:            []EnvVar{{Name: "MY_ENV_BAR", Value: "bar"}},
						Resources:      ResourcesConfig{Requests: ResourcesRequest{CPU: 1, Memory: "1Gi"}},
					},
				},
				PresubmitTasks: map[string]PresubmitTask{
					"build-tests": {TaskParameters: Task{Command: "test/presubmit-tests.sh --build-tests"}},
					"custom-tests": {
						TaskParameters: Task{
							Command:   "test/presubmit-tests.sh --xlarge-tests",
							Image:     "gcr.io/foo/bar/special-image2:dev",
							ImageArgs: []string{"--ram 16GB"},
							Env:       []EnvVar{{Name: "MY_ENV_FOB", Value: "fob"}},
							Resources: ResourcesConfig{Requests: ResourcesRequest{CPU: 2, Memory: "2Gi"}},
						},
						PresubmitFilter: `(?!\.bat)$`,
						Optional:        true,
					},
					"unit-tests": {TaskParameters: Task{Command: "go test ./..."}},
				},
				PostsubmitTasks: map[string]PostsubmitTask{
					"all-tests": {TaskParameters: Task{Command: "go test"}},
					"custom-all-tests": {TaskParameters: Task{
						Command:   "test/presubmit-tests.sh --all",
						Image:     "gcr.io/foo/bar/image",
						ImageArgs: []string{"--ram 6GB"},
						Env:       []EnvVar{{Name: "MY_ENV_BOB", Value: "bob"}},
						Resources: ResourcesConfig{Requests: ResourcesRequest{CPU: 9, Memory: "99Gi"}},
					}},
				},
				PeriodicTasks: map[string]PeriodicTask{
					"auto-release": {
						TaskParameters: Task{
							Command:   "hack/release.py --publish",
							Image:     "gcr.io/foo/bar/special-image4:dev",
							ImageArgs: []string{"--ram 4GB"},
							Env:       []EnvVar{{Name: "MY_ENV_FOY", Value: "foy"}},
							Resources: ResourcesConfig{Requests: ResourcesRequest{CPU: 4, Memory: "4Gi"}},
						},
						Periodicity:             "every 3h",
						MinorReleasePeriodicity: "Tuesdays 3AM",
						MinorReleaseArg:         "--dot-release",
					},
					"large-ci-tests": {
						TaskParameters: Task{
							Command:      "test/presubmit-tests.sh --xlarge-tests",
							EnableDocker: true,
							Image:        "gcr.io/foo/bar/special-image3:dev",
							ImageArgs:    []string{"--ram 1GB"},
							Env:          []EnvVar{{Name: "MY_ENV_FOX", Value: "fox"}},
							Resources:    ResourcesConfig{Requests: ResourcesRequest{CPU: 3, Memory: "3Gi"}},
						},
						Periodicity: "daily 8AM",
					},
					"nightly-release": {
						TaskParameters: Task{
							Command: "hack/release.py --nightly",
							Env:     []EnvVar{{Name: "MY_ENV_FOZ", Value: "foz"}},
						},
						Periodicity: "daily 1PM",
					},
				},
			},
			nil,
		},
	}

	for _, data := range datas {
		r, err := ParseRepoConfig([]byte(data.yaml))
		errMsg := fmt.Sprintf("Parsing:\n%q\n", data.yaml)
		if m := errorMismatch(err, data.err); m != "" {
			t.Errorf("%s%s", errMsg, m)
		}
		if data.err != nil {
			continue
		}
		if dif := cmp.Diff(data.exp, r); dif != "" {
			t.Errorf("%sgot(+) is different from wanted(-)\n%v", errMsg, dif)
		}
	}
}

func TestParseSupervisorConfig(t *testing.T) {
	datas := []struct {
		yaml string
		exp  SupervisorConfig
		err  error
	}{
		{ // empty YAML
			"",
			SupervisorConfig{},
			nil,
		}, { // bad YAML
			"<>",
			SupervisorConfig{},
			fmt.Errorf("yaml: unmarshal errors:\n  line 1: cannot unmarshal !!str `<>` into map[string]interface {}"),
		}, { // unknown defaults section
			`defaults:
  foobar:
    optional: true`,
			SupervisorConfig{
				Defaults: SupervisorDefaults{
					Settings:     DefaultSettings{},
					ProwPlugins:  map[string]ProwPlugin{},
					CodeCoverage: CodeCoverageSettings{},
				},
			},
			nil,
		}, { // bad plugin setting
			`defaults:
  prow-plugins:
    plugin1: 0.0`,
			SupervisorConfig{},
			fmt.Errorf("valid setting for plugin \"plugin1\" is one of [disabled|no|false|<configuration>]"),
		}, { // bad plugin string setting
			`defaults:
  prow-plugins:
    plugin1: foo`,
			SupervisorConfig{},
			fmt.Errorf("valid setting for plugin \"plugin1\" is one of [disabled|no|false|<configuration>]"),
		}, { // full YAML
			`defaults:
  settings:
    merge-method: "squash"
    optional: true
    parse-go-test-output: false
    report-to-github: false
    enable-docker: false
    presubmit-filter: "(?!\\.bak)$"
    testgrid-config: "gs://foo/bar"
    image: "gcr.io/foo/bar/special-image:dev"
    image-args:
      - "--extra-verbose"
    env:
      - name: "MY_ENV_FOO"
        value: "foo"
    resources:
      requests:
        cpu: 1
        memory: "1Gi"
  code-coverage:
    threshold: 80%
    fails-presubmit: true
    image: "gcr.io/foo/bar/special-image1:dev"
    image-args:
    - "--ignore-tests"
    env:
      - name: "MY_ENV_FOB"
        value: "fob"
    resources:
      requests:
        cpu: 2
        memory: "2Gi"
  prow-plugins:
    plugin1: disabled
    plugin2: enabled
    plugin3:
      foo: bar
  nightly-release:
    window: "1AM-3AM PDT"
  auto-release:
    periodicity: "every 2h"
    minor-release-periodicity: "Tuesdays 4AM"
    minor-release-arg: "--dot-release"
repos:
  org1:
  - repo1
  - repo2
`,
			SupervisorConfig{
				Defaults: SupervisorDefaults{
					Settings: DefaultSettings{
						MergeMethod:       "squash",
						Optional:          true,
						ParseGoTestOutput: false,
						ReportToGitHub:    false,
						EnableDocker:      false,
						PresubmitFilter:   "",
						Image:             "gcr.io/foo/bar/special-image:dev",
						ImageArgs:         []string{"--extra-verbose"},
						Env:               []EnvVar{{Name: "MY_ENV_FOO", Value: "foo"}},
						Resources:         ResourcesConfig{Requests: ResourcesRequest{CPU: 1, Memory: "1Gi"}},
						TestGridConfig:    "gs://foo/bar",
					},
					ProwPlugins: map[string]ProwPlugin{
						"plugin1": {},
						"plugin2": {Enabled: true},
						"plugin3": {Enabled: true, Parameters: map[interface{}]interface{}{string("foo"): string("bar")}},
					},
					CodeCoverage: CodeCoverageSettings{
						Threshold:      "80%",
						FailsPresubmit: true,
						Image:          "gcr.io/foo/bar/special-image1:dev",
						ImageArgs:      []string{"--ignore-tests"},
						Env:            []EnvVar{{Name: "MY_ENV_FOB", Value: "fob"}},
						Resources:      ResourcesConfig{Requests: ResourcesRequest{CPU: 2, Memory: "2Gi"}},
					},
					NightlyRelease: NightlyReleaseSettings{Window: "1AM-3AM PDT"},
					AutoRelease: AutoReleaseSettings{
						Window:                  "",
						Periodicity:             "every 2h",
						MinorReleasePeriodicity: "Tuesdays 4AM",
						MinorReleaseArg:         "--dot-release",
					},
				},
				Repositories: map[string][]string{"org1": {"repo1", "repo2"}},
			},
			nil,
		},
	}

	for _, data := range datas {
		r, err := ParseSupervisorConfig([]byte(data.yaml))
		errMsg := fmt.Sprintf("Parsing:\n%q\n", data.yaml)
		if m := errorMismatch(err, data.err); m != "" {
			t.Errorf("%s%s", errMsg, m)
		}
		if data.err != nil {
			continue
		}
		if dif := cmp.Diff(data.exp, r); dif != "" {
			t.Errorf("%sgot(+) is different from wanted(-)\n%v", errMsg, dif)
		}
	}
}
