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
	"os"
	"log"
	"strings"
	"regexp"
)

func main() {
	serviceAccount := flag.String("service-account", os.Getenv("GOOGLE_APPLICATION_CREDENTIALS"), "JSON key file for GCS service account")
	repoNames := flag.String("repo", "test-infra", "repo to be analyzed, comma separated")
	startDate := flag.String("start-date", "2019-02-22", "cut off date to be analyzed")
	parseRegex := flag.String("parser", "", "regex string used for parsing")
	jobFilter := flag.String("jobs", "", "jobs to be analyzed, comma separated")
	prOnly := flag.Bool("pr-only", false, "supplied if just want to analyze PR jobs")
	ciOnly := flag.Bool("ci-only", false, "supplied if just want to analyze CI jobs")
	flag.Parse()

	if "" == *parseRegex {
		log.Fatal("--parser must be provided")
	}

	c, _ := NewParser(*serviceAccount)
	c.logParser = func(s string) bool {
		return regexp.MustCompile(*parseRegex).MatchString(s)
	}
	c.CleanupOnInterrupt()
	defer c.cleanup()

	c.setStartDate(*startDate)
	for _, j := range strings.Split(*jobFilter, ",") {
		if "" != j {
			c.jobFilter = append(c.jobFilter, j)
		}
	}

	for _, repo := range strings.Split(*repoNames, ",") {
		log.Printf("Repo: '%s'", repo)
		if ! *prOnly {
			log.Println("\tProcessing postsubmit jobs")
			c.feedPostsubmitJobsFromRepo(repo)
		}
		if ! *ciOnly {
			log.Println("\tProcessing presubmit jobs")
			c.feedPresubmitJobsFromRepo(repo)
		}
	}
	c.wait()
	log.Printf("Processed %d builds, and found %d matches", len(c.processed), len(c.found))
	for _, l := range c.found {
		log.Printf("\t%s", l)
	}
}
