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
	"os"
	"regexp"
	"sort"
	"strings"
)

func groupByJob(found [][]string) {
	var msgs []string
	for _, elems := range found {
		msgs = append(msgs, strings.Join(elems, ","))
	}
	log.Printf("\n\n%s", strings.Join(msgs, "\n"))
}

func groupByMatch(found [][]string) {
	outArr := make(map[string][][]string)
	for _, l := range found {
		if _, ok := outArr[l[0]]; !ok {
			outArr[l[0]] = make([][]string, 0, 0)
		}
		outArr[l[0]] = append(outArr[l[0]], l)
	}
	for _, sl := range outArr {
		sort.Slice(sl, func(i, j int) bool {
			return sl[i][1] > sl[j][1]
		})
		var msgs []string
		for _, elems := range sl {
			msgs = append(msgs, strings.Join(elems, ","))
		}
		log.Printf("\n\n%s", strings.Join(msgs, "\n"))
	}
}

func main() {
	serviceAccount := flag.String("service-account", os.Getenv("GOOGLE_APPLICATION_CREDENTIALS"), "JSON key file for GCS service account")
	repoNames := flag.String("repo", "test-infra", "repo to be analyzed, comma separated")
	startDate := flag.String("start-date", "2017-01-01", "cut off date to be analyzed")
	parseRegex := flag.String("parser", "", "regex string used for parsing")
	jobFilter := flag.String("jobs", "", "jobs to be analyzed, comma separated")
	prOnly := flag.Bool("pr-only", false, "supplied if just want to analyze PR jobs")
	ciOnly := flag.Bool("ci-only", false, "supplied if just want to analyze CI jobs")
	groupBy := flag.String("groupby", "job(default)", "output groupby, supports: match(group by matches)")
	flag.Parse()

	if "" == *parseRegex {
		log.Fatal("--parser must be provided")
	}

	c, _ := NewParser(*serviceAccount)
	c.logParser = func(s string) string {
		return regexp.MustCompile(*parseRegex).FindString(s)
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
		if !*prOnly {
			log.Println("\tProcessing postsubmit jobs")
			c.feedPostsubmitJobsFromRepo(repo)
		}
		if !*ciOnly {
			log.Println("\tProcessing presubmit jobs")
			c.feedPresubmitJobsFromRepo(repo)
		}
	}
	c.wait()
	log.Printf("Processed %d builds, and found %d matches", len(c.processed), len(c.found))
	switch *groupBy {
	case "job":
		groupByJob(c.found)
	case "match":
		groupByMatch(c.found)
	default:
		log.Printf("--groupby doesn't support %s, fallback to default", *groupBy)
		groupByJob(c.found)
	}

}
