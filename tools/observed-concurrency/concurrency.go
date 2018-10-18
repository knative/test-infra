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

// concurrency.go parses the latest log file of the perf test run and outputs the latency observed.

package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/knative/test-infra/tools/gcs"
	"github.com/knative/test-infra/tools/testgrid"
)

const (
	logDir      = "logs/ci-knative-serving-performance"
	buildFile   = "build-log.txt"
	perfLatency = "perf_latency"
	scaleTo5    = "ScaleTo5(ms)"
)

func getRelevantLogs(fields []string) *string {
	// I1018 01:13:41.537]     5: 2002 ms
	if len(fields) == 5 && fields[2] == "5:" {
		return &fields[3]
	}
	return nil
}

func createTestgridXML(latency float32, dir string) {
	tp := []testgrid.TestProperty{testgrid.TestProperty{Name: perfLatency, Value: latency}}
	tc := []testgrid.TestCase{testgrid.TestCase{ClassName: perfLatency, Name: scaleTo5, Properties: testgrid.TestProperties{Property: tp}}}
	ts := testgrid.TestSuite{TestCases: tc}

	if err := testgrid.CreateXMLOutput(ts, dir); err != nil {
		log.Fatalf("Cannot create the xml output file: %v", err)
	}
}

func main() {

	artifactsDir := flag.String("artifacts-dir", "./artifacts", "Directory to store the generated XML file")
	serviceAccount := flag.String("service-account", os.Getenv("GOOGLE_APPLICATION_CREDENTIALS"), "JSON key file for service account to use")
	flag.Parse()

	// Read the latest-build.txt file to get the latest build number
	ctx := context.Background()
	num, err := gcs.GetLatestBuildNumber(ctx, logDir, *serviceAccount)
	if err != nil {
		log.Fatalf("Cannot get latest build number: %v", err)
	}

	// Get observed latency from the log file
	latencyLogs := gcs.ParseLog(ctx, fmt.Sprintf("%s/%d/%s", logDir, num, buildFile), getRelevantLogs)
	if len(latencyLogs) == 0 {
		log.Fatalf("No relevant log statement found for latency")
	}

	if latency, err := strconv.ParseFloat(latencyLogs[0], 32); err != nil {
		log.Fatalf("Cannot get latency: %v", err)
	} else {
		// Write the testgrid xml to artifacts
		createTestgridXML(float32(latency), *artifactsDir)
	}
}
