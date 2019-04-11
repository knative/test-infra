// +build e2e

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

package e2e

import (
	"testing"
	"time"

	"github.com/knative/test-infra/shared/loadgenerator"
)

func loadTest(t *testing.T, numSteps int) {
	opts := loadgenerator.GeneratorOptions{
		URL:            "http://www.google.com",
		Duration:       10 * time.Second,
		QPS:            10,
		NumThreads:     1,
		NumConnections: 1,
		RequestTimeout: 10 * time.Second,
		NumSteps:       numSteps,
	}

	res, err := opts.RunLoadTest(true)
	if err != nil {
		t.Fatalf("Error performing load test: %v", err)
	}

	if len(res.Result) != numSteps {
		t.Logf("got:%d, want: %d", len(res.Result), numSteps)
		for _, i := range res.Result {
			if i == nil {
				t.Fatal("nil")
			}
			t.Log(i.DurationHistogram.Percentiles[0].Value)
			t.Log(i.DurationHistogram.Percentiles[0].Percentile)
		}
	}
}
func TestLoadGeneratorFullLoad(t *testing.T) {
	loadTest(t, 1)
}

func TestLoadGeneratorStepLoad(t *testing.T) {
	loadTest(t, 3)
}
