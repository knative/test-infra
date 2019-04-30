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

package loadgenerator_test

import (
	"os"
	"testing"

	"github.com/knative/test-infra/shared/loadgenerator"
	"github.com/knative/test-infra/shared/prow"
)

func TestSaveJSON(t *testing.T) {
	res := &loadgenerator.GeneratorResults{FileNamePrefix: t.Name()}
	err := res.SaveJSON()
	if err != nil {
		t.Fatalf("Cannot save JSON: %v", err)
	}

	// Delete the test json file created
	dir := prow.GetLocalArtifactsDir()
	if err = os.Remove(dir + "/" + "TestSaveJSON.json"); err != nil {
		t.Logf("Cannot delete test file: %v", err)
	}
}
