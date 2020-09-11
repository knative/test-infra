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
	"os"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestMain(m *testing.M) {
	ResetOutput() // Redirect output prior to each test.
	os.Exit(m.Run())
}
func TestOutputConfig(t *testing.T) {
	outputConfig("")
	if diff := cmp.Diff(GetOutput(), ""); diff != "" {
		t.Errorf("Incorrect output for empty string: (-got +want)\n%s", diff)
	}

	outputConfig(" \t\n")
	if diff := cmp.Diff(GetOutput(), ""); diff != "" {
		t.Errorf("Incorrect output for whitespace string: (-got +want)\n%s", diff)
	}
	if emittedOutput {
		t.Fatal("emittedOutput was incorrectly set")
	}

	inputLine := "some-key: some-value"
	outputConfig(inputLine)
	if diff := cmp.Diff(GetOutput(), inputLine+"\n"); diff != "" {
		t.Errorf("Incorrect output for whitespace string: (-got +want)\n%s", diff)
	}
	if !emittedOutput {
		t.Fatal("emittedOutput should have been set, but wasn't")
	}
}
