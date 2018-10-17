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

// testgrid.go provides methods to perform action on testgrid.

package main

import (
	"log"

	"github.com/knative/test-infra/tools/testgrid"
)

const (
	apiCoverage    = "api_coverage"
	overallRoute   = "OverallRoute"
	overallConfig  = "OverallConfiguration"
	overallService = "OverallService"
)

func createCases(tcName string, covered map[string]int, notCovered map[string]int) []testgrid.TestCase {
	var tc []testgrid.TestCase

	var percentCovered = float32(100 * len(covered) / (len(covered) + len(notCovered)))
	tp := []testgrid.TestProperty{testgrid.TestProperty{Name: apiCoverage, Value: percentCovered}}
	tc = append(tc, testgrid.TestCase{Name: tcName, Properties: testgrid.TestProperties{Property: tp}, Fail: false})

	for key, value := range covered {
		tp := []testgrid.TestProperty{testgrid.TestProperty{Name: apiCoverage, Value: float32(value)}}
		tc = append(tc, testgrid.TestCase{Name: tcName + "/" + key, Properties: testgrid.TestProperties{Property: tp}, Fail: false})
	}

	for key, value := range notCovered {
		tp := []testgrid.TestProperty{testgrid.TestProperty{Name: apiCoverage, Value: float32(value)}}
		tc = append(tc, testgrid.TestCase{Name: tcName + "/" + key, Properties: testgrid.TestProperties{Property: tp}, Fail: true})
	}
	return tc
}

func createTestgridXML(coverage *OverallAPICoverage, artifactsDir string) {
	tc := createCases(overallRoute, coverage.RouteAPICovered, coverage.RouteAPINotCovered)
	tc = append(tc, createCases(overallConfig, coverage.ConfigurationAPICovered, coverage.ConfigurationAPINotCovered)...)
	tc = append(tc, createCases(overallService, coverage.ServiceAPICovered, coverage.ServiceAPINotCovered)...)
	ts := testgrid.TestSuite{TestCases: tc}

	if err := testgrid.CreateXMLOutput(ts, artifactsDir); err != nil {
		log.Fatalf("Cannot create the xml output file: %v", err)
	}
}
