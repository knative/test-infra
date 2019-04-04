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

// loadgenerator.go provides a wrapper on fortio load generator.

package loadgenerator

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"fortio.org/fortio/fhttp"
	"fortio.org/fortio/periodic"
	"github.com/knative/pkg/test/helpers"
	"github.com/knative/test-infra/shared/prow"
)

const (
	p50      = 50.0
	p90      = 90.0
	p99      = 99.0
	duration = 1 * time.Minute
	qps      = 10
	jsonExt  = ".json"
)

// GeneratorOptions provides knobs to run the perf test
type GeneratorOptions struct {
	// Duration is the total time to generate the load.
	Duration time.Duration
	// NumThreads is the number of threads generating the load.
	NumThreads int
	// NumConnections controls the number of idle connections Transport keeps.
	NumConnections int
	// URL is the endpoint of the target service.
	URL string
	// Domain is the domain name of the target service.
	Domain string
	// RequestTimeout is the maximum time waiting for a response.
	RequestTimeout time.Duration
	// QPS is the total QPS generated by all threads.
	QPS float64
	// AllowInitialErrors specifies whether initial errors will abort the load.
	AllowInitialErrors bool
}

// GeneratorResults contains the results of running the per test
type GeneratorResults struct {
	Result *fhttp.HTTPRunnerResults
}

// addDefaults adds default values to non mandatory params
func (g *GeneratorOptions) addDefaults() {
	if g.RequestTimeout == 0 {
		g.RequestTimeout = duration
	}

	if g.QPS == 0 {
		g.QPS = qps
	}
}

// CreateRunnerOptions sets up the fortio client with the knobs needed to run the load test
func (g *GeneratorOptions) CreateRunnerOptions(resolvableDomain bool) *fhttp.HTTPRunnerOptions {
	g.addDefaults()
	o := fhttp.NewHTTPOptions(g.URL)

	o.NumConnections = g.NumConnections
	o.HTTPReqTimeOut = g.RequestTimeout

	// If the url does not contains a resolvable domain, we need to add the domain as a header
	if !resolvableDomain {
		o.AddAndValidateExtraHeader(fmt.Sprintf("Host: %s", g.Domain))
	}

	return &fhttp.HTTPRunnerOptions{
		RunnerOptions: periodic.RunnerOptions{
			Duration:    g.Duration,
			NumThreads:  g.NumThreads,
			Percentiles: []float64{p50, p90, p99},
			QPS:         g.QPS,
		},
		HTTPOptions:        *o,
		AllowInitialErrors: g.AllowInitialErrors,
	}
}

// RunLoadTest runs the load test with fortio and returns the response
func (g *GeneratorOptions) RunLoadTest(resolvableDomain bool) (*GeneratorResults, error) {
	r, err := fhttp.RunHTTPTest(g.CreateRunnerOptions(resolvableDomain))
	return &GeneratorResults{Result: r}, err
}

// SaveJSON saves the results as Json in the artifacts directory
func (gr *GeneratorResults) SaveJSON(testName string) error {
	dir := prow.GetLocalArtifactsDir()
	if err := helpers.CreateDir(dir); err != nil {
		return err
	}

	outputFile := dir + "/" + testName + jsonExt
	log.Printf("Storing json output in %s", outputFile)
	f, err := os.OpenFile(outputFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}

	defer f.Close()
	json, err := json.Marshal(gr)
	if err != nil {
		return err
	}
	if _, err = f.Write(json); err != nil {
		return err
	}

	return nil
}
