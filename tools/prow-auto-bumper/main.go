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

// prow-auto-bumper finds stable Prow components version used by k8s,
// and creates PRs updating them in knative/test-infra

package main

import (
	"flag"
	"fmt"
	"log"
	"math"
	"regexp"
	"strings"
	"time"

	"github.com/google/go-github/github"
	"github.com/knative/test-infra/shared/ghutil"
)

const (
	org  = "kubernetes"
	repo = "test-infra"
	// PRHead is the head branch of k8s auto version bump PRs
	// TODO(chaodaiG): using head branch querying is less ideal than using
	// label `area/prow/bump`, which is not supported by Github API yet. Move
	// to filter using this label once it's supported
	PRHead = "k8s-ci-robot:autobump"
	// PRBase is the base branch of k8s auto version bump PRs
	PRBase = "master"
	// Index for regex matching groups
	imageImagePart = 1
	imageTagPart   = 2
	// Max difference away from target date
	maxDelta = 2 * 24 // 2 days
	// Safe duration is the smallest amount of hours a version stayed
	safeDuration = 12 // 12 hours
	maxRetry     = 3
)

var (
	// matching            gcr.io /k8s-(prow|testimage)/(tide|kubekin-e2e|.*)    :vYYYYMMDD-HASH-VARIANT
	imagePattern     = `\b(gcr\.io/k8s[a-z0-9-]{5,29}/[a-zA-Z0-9][a-zA-Z0-9_.-]+):(v[a-zA-Z0-9_.-]+)\b`
	imageRegexp      = regexp.MustCompile(imagePattern)
	imageLinePattern = fmt.Sprintf(`\s+[a-z]+:\s+"?'?%s"?'?`, imagePattern)
	// matching   "-    image: gcr.io /k8s-(prow|testimage)/(tide|kubekin-e2e|.*)    :vYYYYMMDD-HASH-VARIANT"
	imageMinusRegexp = regexp.MustCompile(fmt.Sprintf(`\-%s`, imageLinePattern))
	// matching   "+    image: gcr.io /k8s-(prow|testimage)/(tide|kubekin-e2e|.*)    :vYYYYMMDD-HASH-VARIANT"
	imagePlusRegexp = regexp.MustCompile(fmt.Sprintf(`\+%s`, imageLinePattern))
	// Preferred time for candidate PR creation date
	targetTime = time.Now().Add(-time.Hour * 7 * 24) // 7 days
)

// Client handles methods for github issues
type Client struct {
	client ghutil.GithubOperations
}

// versions holds the version change for an image
// oldVersion and newVersion are both in the format of "vYYYYMMDD-HASH-VARIANT"
type versions struct {
	oldVersion string
	newVersion string
	variant    string
}

// PRVersions contains PR and version changes in it
type PRVersions struct {
	images map[string][]versions // map of image name: versions struct
	// The way k8s updates versions doesn't guarantee the same version tag across all images,
	// dominantVs is the version that appears most times
	dominantVs *versions
	PR         *github.PullRequest
}

// Helper method for adding a newly discovered tag into pv
func (pv *PRVersions) getIndex(image, tag string) int {
	if _, ok := pv.images[image]; !ok {
		pv.images[image] = make([]versions, 0, 0)
	}
	_, variant := deconstructTag(tag)
	iv := -1
	for i, vs := range pv.images[image] {
		if vs.variant == variant {
			iv = i
			break
		}
	}
	if -1 == iv {
		pv.images[image] = append(pv.images[image], versions{variant: variant})
		iv = len(pv.images[image]) - 1
	}
	return iv
}

// Tags could be in the form of: v[YYYYMMDD]-[GIT_HASH](-[VARIANT_PART]),
// separate it to `v[YYYYMMDD]-[GIT_HASH]` and `[VARIANT_PART]`
func deconstructTag(in string) (string, string) {
	dateCommit := in
	var variant string
	parts := strings.Split(in, "-")
	if len(parts) > 2 {
		variant = strings.Join(parts[2:], "-")
	}
	if len(parts) > 1 {
		dateCommit = fmt.Sprintf("%s-%s", parts[0], parts[1])
	}
	return dateCommit, variant
}

// get key with highest value
func getDominantKey(m map[string]int) string {
	var res string
	for key, v := range m {
		if "" == res || v > m[res] {
			res = key
		}
	}
	return res
}

func (pv *PRVersions) getDominantVersions() versions {
	if nil != pv.dominantVs {
		return *pv.dominantVs
	}

	cOld := make(map[string]int)
	cNew := make(map[string]int)
	for _, vss := range pv.images {
		for _, vs := range vss {
			normOldTag, _ := deconstructTag(vs.oldVersion)
			normNewTag, _ := deconstructTag(vs.newVersion)
			cOld[normOldTag]++
			cNew[normNewTag]++
		}
	}

	pv.dominantVs = &versions{
		oldVersion: getDominantKey(cOld),
		newVersion: getDominantKey(cNew),
	}

	return *pv.dominantVs
}

// parse changelist, find all version changes, and store them in image name: versions map
func (pv *PRVersions) parseChangelist(client *Client) error {
	fs, err := client.client.ListFiles(org, repo, *pv.PR.Number)
	if nil != err {
		return err
	}
	for _, f := range fs {
		if nil == f.Patch {
			continue
		}
		minuses := imageMinusRegexp.FindAllStringSubmatch(*f.Patch, -1)
		for _, minus := range minuses {
			iv := pv.getIndex(minus[imageImagePart], minus[imageTagPart])
			pv.images[minus[imageImagePart]][iv].oldVersion = minus[imageTagPart]
		}

		pluses := imagePlusRegexp.FindAllStringSubmatch(*f.Patch, -1)
		for _, plus := range pluses {
			iv := pv.getIndex(plus[imageImagePart], plus[imageTagPart])
			pv.images[plus[imageImagePart]][iv].newVersion = plus[imageTagPart]
		}
	}

	return nil
}

// Query all PRs from "k8s-ci-robot:autobump", find PR roughly 7 days old and was not reverted later.
// Only return error if it's github related
func getBestVersion(client *Client, org, repo, head, base string) (*PRVersions, error) {
	visited := make(map[string]PRVersions)
	var bestPv *PRVersions
	var overallErr error
	var bestDelta float64 = maxDelta + 1
	PRs, err := client.client.ListPullRequests(org, repo, head, base)
	if nil != err {
		return bestPv, fmt.Errorf("failed list pull request: '%v'", err)
	}

	for _, PR := range PRs {
		if nil == PR.State || string(ghutil.PullRequestCloseState) != *PR.State {
			continue
		}
		delta := targetTime.Sub(*PR.CreatedAt).Hours()
		if delta > maxDelta {
			break // Over 9 days old, too old
		}
		pv := PRVersions{
			images: make(map[string][]versions),
			PR:     PR,
		}
		if err := pv.parseChangelist(client); nil != err {
			overallErr = fmt.Errorf("failed listing files from PR '%d': '%v'", *PR.Number, err)
			break
		}
		vs := pv.getDominantVersions()
		if "" == vs.oldVersion || "" == vs.newVersion {
			log.Printf("Warning: found PR misses version change '%d'", *PR.Number)
			continue
		}
		visited[vs.oldVersion] = pv
		// check if too fresh here as need the data in visited
		if delta < -maxDelta { // In past 5 days, too fresh
			continue
		}
		if updatePR, ok := visited[vs.newVersion]; ok {
			if updatePR.getDominantVersions().newVersion == vs.oldVersion { // The updatePR is reverting this PR
				continue
			}
			if updatePR.PR.CreatedAt.Before(PR.CreatedAt.Add(time.Hour * safeDuration)) {
				// The update PR is within 12 hours of current PR, consider unsafe
				continue
			}
		}
		if nil == bestPv || math.Abs(delta) < math.Abs(bestDelta) {
			bestDelta = delta
			bestPv = &pv
		}
	}
	return bestPv, overallErr
}

func retryGetBestVersion(client *Client, org, repo, head, base string) (*PRVersions, error) {
	var bestPv *PRVersions
	var overallErr error
	// retry if there is github related error
	for retryCount := 0; nil == overallErr && retryCount < maxRetry; retryCount++ {
		bestPv, overallErr = getBestVersion(client, org, repo, head, base)
		if nil != overallErr {
			log.Println(overallErr)
			if maxRetry-1 != retryCount {
				log.Printf("Retry #%d", retryCount+1)
			}
		}
	}
	return bestPv, overallErr
}

func main() {
	githubAccount := flag.String("github-account", "", "Token file for Github authentication")
	flag.Parse()

	GHClient, err := ghutil.NewGithubClient(*githubAccount)
	if nil != err {
		log.Fatalf("cannot authenticate to github: %v", err)
	}
	client := &Client{GHClient}

	bestVersion, err := retryGetBestVersion(client, org, repo, PRHead, PRBase)
	if nil != err {
		log.Fatalf("cannot get best version from %s/%s: '%v'", org, repo, err)
	}

	log.Println(bestVersion.images)
	log.Println(bestVersion.dominantVs)
}
