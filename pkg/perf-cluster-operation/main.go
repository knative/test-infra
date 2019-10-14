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

package main

import (
	"flag"

	"log"
)

var (
	gcpProjectName      string
	repoName            string
	benchmarkRootFolder string
)

func main() {
	flag.StringVar(&gcpProjectName, "gcp-project", "", "name of the GCP project for cluster operations")
	flag.StringVar(&repoName, "repository", "", "name of the repository")
	flag.StringVar(&benchmarkRootFolder, "benchmark-root", "", "root folder of the benchmarks")
	flag.Parse()
	args := flag.Args()
	if len(args) != 1 {
		log.Fatal("only one operation name can be provided as the arg")
	}

	client, err := newClient()
	if err != nil {
		log.Fatalf("failed setting up GKE client, cannot proceed: %v", err)
	}
	operation := args[0]
	switch operation {
	case "recreate":
		if err := client.recreateClusters(gcpProjectName, repoName, benchmarkRootFolder); err != nil {
			log.Fatalf("Failed recreating clusters for repo %q: %v", repoName, err)
		} else {
			log.Printf("Done with recreating clusters for repo %q", repoName)
		}
	case "reconcile":
		if err := client.reconcileClusters(gcpProjectName, repoName, benchmarkRootFolder); err != nil {
			log.Fatalf("Failed reconciling clusters for repo %q: %v", repoName, err)
		} else {
			log.Printf("Done with reconciling clusters for repo %q", repoName)
		}
	default:
		log.Fatalf("operation %q is not supported", operation)
	}
}
