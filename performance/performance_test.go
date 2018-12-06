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

package performance_test

import (
	"testing"
	"time"

	"github.com/knative/test-infra/performance"
)

const (
	testTime = 1 * time.Minute
	testNum  = 5
	testUrl  = "http://example.com"
	testQPS  = 100.0
)

func getPerfOptions() *performance.PerfOptions {
	return &performance.PerfOptions{
		Duration:       testTime,
		NumThreads:     testNum,
		NumConnections: testNum,
		URL:            testUrl,
		Domain:         testUrl,
		RequestTimeout: testTime,
		QPS:            testQPS,
	}
}

func TestRunLoadTest(t *testing.T) {
	o := getPerfOptions()
	opts := o.CreateRunnerOptions(false)

	if opts.RunnerOptions.Duration != testTime {
		t.Fatalf("Duration is %v. Expected %v", opts.RunnerOptions.Duration, testTime)
	}

	if opts.RunnerOptions.NumThreads != testNum {
		t.Fatalf("Number of threads is %d. Expected %d", opts.RunnerOptions.NumThreads, testNum)
	}

	if opts.RunnerOptions.QPS != testQPS {
		t.Fatalf("QPS is %f. Expected %f", opts.RunnerOptions.QPS, testQPS)
	}

	if opts.HTTPOptions.NumConnections != testNum {
		t.Fatalf("Number of connections is %d. Expected %d", opts.HTTPOptions.NumConnections, testNum)
	}

	if opts.HTTPOptions.HTTPReqTimeOut != testTime {
		t.Fatalf("Request timeout is %v. Expected %v", opts.HTTPOptions.HTTPReqTimeOut, testTime)
	}

	if opts.HTTPOptions.URL != testUrl {
		t.Fatalf("Url is %s. Expected %s", opts.HTTPOptions.URL, testUrl)
	}
}
