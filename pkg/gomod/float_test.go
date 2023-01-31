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

package gomod

import (
	"testing"

	"github.com/stretchr/testify/require"
	"knative.dev/test-infra/pkg/git"
)

// TestFloat - This is an integration test, it will make a call out to the internet.
func TestFloat(t *testing.T) {
	tests := map[string]struct {
		gomod         string
		release       string
		moduleRelease string
		domain        string
		rule          git.RulesetType
		want          map[string]git.RefType
		wantRef       map[string]string // it is a bad idea to test with anything other than `release`.
	}{
		"demo1, v0.15, knative.dev, any rule": {
			gomod:         "./testdata/gomod.float1",
			release:       "v0.15",
			moduleRelease: "v0.15",
			domain:        "knative.dev",
			rule:          git.AnyRule,
			want: map[string]git.RefType{
				"knative.dev/pkg":      git.ReleaseBranchRef,
				"knative.dev/eventing": git.ReleaseRef,
			},
		},
		"demo1, v0.15, knative.dev, release rule": {
			gomod:         "./testdata/gomod.float1",
			release:       "v0.15",
			moduleRelease: "v0.15",
			domain:        "knative.dev",
			rule:          git.ReleaseRule,
			want: map[string]git.RefType{
				"knative.dev/eventing": git.ReleaseRef,
			},
		},
		"demo1, v0.15, knative.dev, release branch rule": {
			gomod:         "./testdata/gomod.float1",
			release:       "v0.15",
			moduleRelease: "v0.15",
			domain:        "knative.dev",
			rule:          git.ReleaseBranchRule,
			want: map[string]git.RefType{
				"knative.dev/pkg":      git.ReleaseBranchRef,
				"knative.dev/eventing": git.ReleaseBranchRef,
			},
		},
		"demo1, v99.99, knative.dev, release branch or release rule": {
			gomod:         "./testdata/gomod.float1",
			release:       "v99.99",
			moduleRelease: "v99.99",
			domain:        "knative.dev",
			rule:          git.ReleaseOrReleaseBranchRule,
			want:          map[string]git.RefType{},
		},
		"demo1, v0.16, knative.dev, any rule": {
			gomod:         "./testdata/gomod.float1",
			release:       "v0.16",
			moduleRelease: "v0.16",
			domain:        "knative.dev",
			rule:          git.AnyRule,
			want: map[string]git.RefType{
				"knative.dev/pkg":      git.ReleaseBranchRef,
				"knative.dev/eventing": git.ReleaseRef,
			},
		},
		"demo1, v0.16, k8s.io, any rule": {
			gomod:         "./testdata/gomod.float1",
			release:       "v0.16",
			moduleRelease: "v0.16",
			domain:        "knative.dev",
			rule:          git.AnyRule,
			want: map[string]git.RefType{
				"knative.dev/pkg":      git.ReleaseBranchRef,
				"knative.dev/eventing": git.ReleaseRef,
			},
		},
		"demo1, v0.16, v0.15 mods, k8s.io, any rule": {
			gomod:         "./testdata/gomod.float1",
			release:       "v0.16",
			moduleRelease: "v0.15",
			domain:        "knative.dev",
			rule:          git.AnyRule,
			want: map[string]git.RefType{
				"knative.dev/pkg":      git.ReleaseBranchRef,
				"knative.dev/eventing": git.ReleaseRef,
			},
			wantRef: map[string]string{
				"knative.dev/pkg": "release-0.16",
			},
		},
		"demo1, v99.99, knative.dev, any rule": {
			gomod:         "./testdata/gomod.float1",
			release:       "v99.99",
			moduleRelease: "v99.99",
			domain:        "knative.dev",
			rule:          git.AnyRule,
			want: map[string]git.RefType{
				"knative.dev/pkg":      git.BranchRef,
				"knative.dev/eventing": git.BranchRef,
			},
		},
		"demo1, v0.15, v99.99 mods, knative.dev, any rule": {
			gomod:         "./testdata/gomod.float1",
			release:       "v0.15",
			moduleRelease: "v99.99",
			domain:        "knative.dev",
			rule:          git.AnyRule,
			want: map[string]git.RefType{
				"knative.dev/pkg":      git.ReleaseBranchRef,
				"knative.dev/eventing": git.ReleaseBranchRef,
			},
			wantRef: map[string]string{
				"knative.dev/pkg":      "release-0.15",
				"knative.dev/eventing": "release-0.15",
			},
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			selector, err := DefaultSelector(tt.domain)
			require.NoError(t, err)
			deps, err := Float(tt.gomod, tt.release, tt.moduleRelease, selector, tt.rule)
			if err != nil {
				t.Fatal(err)
			}
			for _, dep := range deps {
				module, ref, got := git.ParseRef(dep)
				if want, ok := tt.want[module]; ok {
					if got != want {
						t.Errorf("[ref type] Float() %s; got %q, want: %q", module, got, want)
					}
				} else {
					t.Error("untested float dep: ", dep)
				}
				// Optional test on returned ref.
				if want, ok := tt.wantRef[module]; ok {
					if got := ref; got != want {
						t.Errorf("[ref value] Float() %s; got %q, want: %q", module, got, want)
					}
				}
			}
		})
	}
}

// TestFloatUnhappy - This is an integration test, it will make a call out to the internet.
func TestFloatUnhappy(t *testing.T) {
	tests := map[string]struct {
		gomod   string
		release string
		domain  string
		rule    git.RulesetType
	}{
		"bad go mod file": {
			gomod:   "./testdata/bad.example",
			release: "v0.15",
			domain:  "knative.dev",
			rule:    git.AnyRule,
		},
		"bad version": {
			gomod:   "./testdata/gomod.float1",
			release: "jupiter/is/angry",
			domain:  "knative.dev",
			rule:    git.AnyRule,
		},
		"bad go module": {
			gomod:   "./testdata/gomod.float1",
			release: "v0.15",
			domain:  "does-not-exist.nope",
			rule:    git.AnyRule,
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			selector, err := DefaultSelector(tt.domain)
			require.NoError(t, err)
			_, err = Float(tt.gomod, tt.release, tt.release, selector, tt.rule)
			if err == nil {
				t.Error("Expected an error")
			}
		})
	}
}
