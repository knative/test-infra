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
	"fmt"
	"log"
	"math"
	"strings"
	"time"

	"github.com/knative/test-infra/shared/ghutil"
)

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
	if nil != pv.dominantVersions {
		return *pv.dominantVersions
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

	pv.dominantVersions = &versions{
		oldVersion: getDominantKey(cOld),
		newVersion: getDominantKey(cNew),
	}

	return *pv.dominantVersions
}

// parse changelist, find all version changes, and store them in image name: versions map
func (pv *PRVersions) parseChangelist(gcw *GHClientWrapper, gi gitInfo) error {
	fs, err := gcw.ListFiles(gi.org, gi.repo, *pv.PR.Number)
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
func getBestVersion(gcw *GHClientWrapper, gi gitInfo) (*PRVersions, error) {
	visited := make(map[string]PRVersions)
	var bestPv *PRVersions
	var overallErr error
	var bestDelta float64 = maxDelta + 1
	PRs, err := gcw.ListPullRequests(gi.org, gi.repo, gi.getHeadRef(), gi.base)
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
		if err := pv.parseChangelist(gcw, gi); nil != err {
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

func retryGetBestVersion(gcw *GHClientWrapper, gi gitInfo) (*PRVersions, error) {
	var bestPv *PRVersions
	var overallErr error
	// retry if there is github related error
	for retryCount := 0; nil == overallErr && retryCount < maxRetry; retryCount++ {
		bestPv, overallErr = getBestVersion(gcw, gi)
		if nil != overallErr {
			log.Println(overallErr)
			if maxRetry-1 != retryCount {
				log.Printf("Retry #%d", retryCount+1)
			}
		}
	}
	return bestPv, overallErr
}
