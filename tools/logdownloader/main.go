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

// flaky-test-reporter collects test results from continuous flows,
// identifies flaky tests, tracking flaky tests related github issues,
// and sends slack notifications.

package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/knative/test-infra/shared/common"
	"github.com/knative/test-infra/shared/prow"
)

type leanBuild struct {
	job *prow.Job
	ID  int // Build ID
}

type logInfo struct {
	build *prow.Build
	l     string
}

type prInfo struct {
	repoName string
	ID       int
}

type jobInfo struct {
	repoName string
	ID       int
	jobName  string
	buildIDs []int
}

// IDFromPR is
func IDFromPR(pr string) int {
	ID, err := strconv.Atoi(path.Base(pr))
	if nil != err {
		return -1
	}
	return ID
}

func processPR(prChan chan prInfo, buildChan chan leanBuild, wg *sync.WaitGroup) {
	for {
		select {
		case pr := <-prChan:
			// fmt.Printf("Process PR '%v'\n", pr)
			for _, j := range prow.ListJobsFromPR(pr.repoName, pr.ID) {
				// fmt.Printf("job: %s\n", j.Name)
				for _, buildID := range prow.NewJob(j.Name, prow.PresubmitJob, pr.repoName, j.PullID).GetBuildIDs() {
					wg.Add(1)
					// fmt.Printf("build: %s %d\n", j.Name, buildID)
					buildChan <- leanBuild{job: j, ID: buildID}
				}
			}
		}
	}
}

func processBuild(buildChan chan leanBuild, logChan chan logInfo, startTimestamp int64, wg, wg2 *sync.WaitGroup) {
	for {
		select {
		case b := <-buildChan:
			// fmt.Print(b.job, b.ID)
			build := b.job.NewBuild(b.ID)
			if build.FinishTime != nil && *build.FinishTime > startTimestamp {
				wg2.Add(1)
				content, _ := build.ReadFile("build-log.txt")
				logChan <- logInfo{
					build: build,
					l:     string(content),
				}
			}
			wg.Done()
		}
	}
}

func parseFunction(l string) bool {
	return strings.Contains(l, "does not have enough resources")
}

func processLog(logChan chan logInfo, parseFunc func(string) bool, wg2 *sync.WaitGroup) {
	for {
		select {
		case li := <-logChan:
			if parseFunc(li.l) {
				fmt.Println(li.build.StoragePath)
			}
			wg2.Done()
		}
	}
}

func downloadSingle(repo string, startTimestamp int64) {
	prChan := make(chan prInfo)
	buildChan := make(chan leanBuild, 500)
	logChan := make(chan logInfo, 500)

	wg := &sync.WaitGroup{}
	wg2 := &sync.WaitGroup{}

	defer func() {
		close(prChan)
		close(buildChan)
		close(logChan)
	}()

	for i := 0; i < 100; i++ {
		go processPR(prChan, buildChan, wg)
	}
	for i := 0; i < 500; i++ {
		go processBuild(buildChan, logChan, startTimestamp, wg, wg2)
	}
	for i := 0; i < 500; i++ {
		go processLog(logChan, parseFunction, wg2)
	}

	prs := prow.ListPRs(repo)
	for _, pr := range prs {
		if ID := IDFromPR(pr); -1 != ID {
			prChan <- prInfo{
				repoName: repo,
				ID:       ID,
			}
		}
	}

	wg.Wait()
	wg2.Wait()
}

func download(repoNames string, startTimestamp int64) error {
	repos := strings.Split(repoNames, ",")
	for _, repo := range repos {
		downloadSingle(repo, startTimestamp)
	}
	return nil
}

func main() {
	serviceAccount := flag.String("service-account", os.Getenv("GOOGLE_APPLICATION_CREDENTIALS"), "JSON key file for GCS service account")
	repoNames := flag.String("repo", "serving", "repo to be downloaded")
	startDate := flag.String("date-since", "2019-02-22", "cut off date to be analyzed")
	dryrun := flag.Bool("dry-run", false, "dry run switch")
	flag.Parse()

	if nil != dryrun && true == *dryrun {
		log.Printf("running in [dry run mode]")
	}

	if err := prow.Initialize(*serviceAccount); nil != err { // Explicit authenticate with gcs Client
		log.Fatalf("Failed authenticating GCS: '%v'", err)
	}

	// Clean up local artifacts directory, this will be used later for artifacts uploads
	if err := common.CreateDir(prow.GetLocalArtifactsDir()); nil != err {
		log.Fatalf("Failed preparing local artifacts directory: %v", err)
	}

	tt, err := time.Parse("2006-01-02", *startDate)
	if nil != err {
		log.Fatalf("invalid start date string, expecing format YYYY-MM-DD: '%v'", err)
	}

	if err := download(*repoNames, tt.Unix()); nil != err {
		log.Fatalf("Failed downloading '%v'", err)
	}

	log.Println("Done")
}
