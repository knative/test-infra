/*
Copyright 2019 The Knative Authors

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    https://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package githubUtil

import (
	"os"
	"testing"
)

func TestFilePathProfileToGithub(t *testing.T) {
	t.Run("repo=knative.dev", func(t *testing.T) {
		input := "knative.dev/test-infra/pkg/ab/cde"
		expectedOutput := "pkg/ab/cde"
		actualOutput := FilePathProfileToGithub(input)
		if actualOutput != expectedOutput {
			t.Fatalf("input=%s; expected output=%s; actual output=%s", input, expectedOutput,
				actualOutput)
		}
	})
	t.Run("repo=github.com/{repo}", func(t *testing.T) {
		input := "github.com/myRepoOwner/myRepoName/pkg/ab/cde"
		expectedOutput := "pkg/ab/cde"
		repoRoot := "/d1/d2/d3/gopath/src/github.com/myRepoOwner/myRepoName"
		getRepoRoot = func() (string, error) {
			return repoRoot, nil
		}
		gopath := os.Getenv("GOPATH")
		os.Setenv("GOPATH", "/d1/d2/d3/gopath")
		defer os.Setenv("GOPATH", gopath)
		actualOutput := FilePathProfileToGithub(input)

		if actualOutput != expectedOutput {
			t.Fatalf("input=%s; expected output=%s; actual output=%s", input, expectedOutput,
				actualOutput)
		}
	})
}
