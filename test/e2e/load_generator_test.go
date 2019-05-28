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

func loadTest(t *testing.T, factors bool, profiler bool) {
	opts := loadgenerator.GeneratorOptions{
		URL:            "http://www.google.com",
		Duration:       10 * time.Second,
		BaseQPS:        10,
		NumThreads:     1,
		NumConnections: 1,
		RequestTimeout: 10 * time.Second,
	}

	if factors {
		opts.LoadFactors = []float64{1, 2, 4}
	}

	if profiler {
		opts.FileNamePrefix = t.Name()
	}

	res, err := opts.RunLoadTest()
	if err != nil {
		t.Fatalf("Error performing load test: %v", err)
	}

	if factors && len(res.Result) != 3 {
		t.Fatalf("got:%d, want: 3", len(res.Result))
	}
}
func TestFullLoad(t *testing.T) {
	loadTest(t, false, false)
}

func TestStepLoad(t *testing.T) {
	loadTest(t, true, false)
}

func TestProfiler(t *testing.T) {
	loadTest(t, false, true)
}
