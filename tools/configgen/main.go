/*
Copyright 2022 The Knative Authors

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
	"flag"
	"log"

	"knative.dev/test-infra/tools/configgen/pkg"
)

var (
	prowJobsConfigInput  string
	prowJobsConfigOutput string
	allProwJobsConfig    string
	testgridConfigOutput string
)

func main() {
	flag.StringVar(&prowJobsConfigInput, "prow-jobs-config-input", "", "The input path for the prow jobs config")
	flag.StringVar(&prowJobsConfigOutput, "prow-jobs-config-output", "", "The output path for the prow jobs config")
	flag.StringVar(&allProwJobsConfig, "all-prow-jobs-config", "", "The path for all prow jobs config")
	flag.StringVar(&testgridConfigOutput, "testgrid-config-output", "", "The output path for the testgrid config")

	flag.Parse()
	if prowJobsConfigInput == "" {
		log.Fatal("--prow-jobs-config-input must be specified")
	}
	if prowJobsConfigOutput == "" {
		log.Fatal("--prow-jobs-config-output must be specified")
	}
	if allProwJobsConfig == "" {
		log.Fatal("--all-prow-jobs-config must be specified")
	}
	if testgridConfigOutput == "" {
		log.Fatal("--testgrid-config-output must be specified")
	}

	if err := pkg.GenerateProwJobsConfig(prowJobsConfigInput, prowJobsConfigOutput); err != nil {
		log.Fatalf("Error generating Prow jobs: %v", err)
	}

	if err := pkg.GenerateTestGridConfig(allProwJobsConfig, testgridConfigOutput); err != nil {
		log.Fatalf("Error generating TestGrid config: %v", err)
	}
}
