/*
Copyright 2018 The Knative Authors

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

package testgrid_test

import (
	"os"
	"testing"

	"github.com/knative/test-infra/tools/testgrid"
)

func TestGetArtifacts(t *testing.T) {
	if v := testgrid.GetArtifactsDir(); v != "./artifacts" {
		t.Fatalf("Default value is %s and not artifacts", v)
	}
}

func TestXMLOutput(t *testing.T) {
	// Create a test file
	if err := testgrid.CreateXMLOutput(testgrid.TestSuite{}, "."); err != nil {
		t.Fatalf("Error when creating xml output file: %v", err)
	}

	// Make sure we can append to the file
	if err := testgrid.CreateXMLOutput(testgrid.TestSuite{}, "."); err != nil {
		t.Fatalf("Error when creating xml output file: %v", err)
	}

	// Delete the test file created
	if err := os.Remove("./" + testgrid.Filename); err != nil {
		t.Logf("Cannot delete test file")
	}
}
