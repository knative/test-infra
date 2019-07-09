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

package performance

import (
	"fmt"

	"github.com/knative/test-infra/shared/junit"
)

const (
	// Latency of the performance test, it's a property name used by testgrid
	perfLatency = "perf_latency"
	// Error rate of the performance test
	perfErrorRate = "perf_error_rate"
)

// CreatePerfLatencyTestCase creates a perf latency test case with the provided name and value
func CreatePerfLatencyTestCase(metricValue float32, metricName, testName string) junit.TestCase {
	tp := []junit.TestProperty{{Name: perfLatency, Value: fmt.Sprintf("%f", metricValue)}}
	tc := junit.TestCase{
		ClassName:  testName,
		Name:       fmt.Sprintf("%s/%s", testName, metricName),
		Properties: junit.TestProperties{Properties: tp}}

	return tc
}

// CreatePerfErrorRateTestCase creates a perf error rate test case with the provided value
func CreatePerfErrorRateTestCase(metricValue float32, testName string) junit.TestCase {
	tp := []junit.TestProperty{{Name: perfErrorRate, Value: fmt.Sprintf("%f", metricValue)}}
	tc := junit.TestCase{
		ClassName:  testName,
		Name:       fmt.Sprintf("%s/error_rate", testName),
		Properties: junit.TestProperties{Properties: tp}}

	return tc
}
