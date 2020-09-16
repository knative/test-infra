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
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestOutputConfig(t *testing.T) {
	output.outputConfig("")
	if diff := cmp.Diff(GetOutput(), ""); diff != "" {
		t.Errorf("Incorrect output for empty string: (-got +want)\n%s", diff)
	}

	output.outputConfig(" \t\n")
	if diff := cmp.Diff(GetOutput(), ""); diff != "" {
		t.Errorf("Incorrect output for whitespace string: (-got +want)\n%s", diff)
	}
	if output.count != 0 {
		t.Fatalf("Output count should have been 0, but was %d", output.count)
	}

	inputLine := "some-key: some-value"
	output.outputConfig(inputLine)
	if diff := cmp.Diff(GetOutput(), inputLine+"\n"); diff != "" {
		t.Errorf("Incorrect output for whitespace string: (-got +want)\n%s", diff)
	}
	if output.count != 1 {
		t.Fatalf("Output count should have been exactly 1, but was %d", output.count)
	}
}
