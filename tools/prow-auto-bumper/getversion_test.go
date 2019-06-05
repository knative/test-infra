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
	"fmt"
	"log"
	"reflect"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/google/go-github/github"
	"github.com/knative/test-infra/shared/ghutil"
	"github.com/knative/test-infra/shared/ghutil/fakeghutil"
)

var (
	fakeOrg  = "fakeorg"
	fakeRepo = "fakerepo"
)

func createPullRequest(t *testing.T, fgc *fakeghutil.FakeGithubClient) *github.PullRequest {
	PR, err := fgc.CreatePullRequest(fakeOrg, fakeRepo, "user:head", "base", "title", "body")
	if nil != err {
		t.Fatalf("Create PR in %s/%s, want: no error, got: '%v'", fakeOrg, fakeRepo, err)
	}
	return PR
}

func TestGetIndex(t *testing.T) {
	datas := []struct {
		images map[string][]versions
		image  string
		tag    string
		i      int
	}{
		{
			map[string][]versions{},
			"o",
			"a-b-c",
			0,
		},
		{
			map[string][]versions{
				"o": []versions{{"", "", ""}},
			},
			"o",
			"a-b-c",
			1,
		},
		{
			map[string][]versions{
				"o": []versions{{"", "", "b"}},
			},
			"o",
			"a-b-c",
			1,
		},
		{
			map[string][]versions{
				"o": []versions{{"", "", "c"}},
			},
			"o",
			"a-b-c",
			0,
		},
		{
			map[string][]versions{
				"p": []versions{{"", "", "c"}},
			},
			"o",
			"a-b-c",
			0,
		},
	}

	for _, data := range datas {
		msg := fmt.Sprintf("getIndex with PRInfo '%v', image '%s', tag '%s'",
			data.images, data.image, data.tag)
		pv := PRVersions{
			images: data.images,
		}

		if got := pv.getIndex(data.image, data.tag); got != data.i {
			t.Fatalf("%s, want: %d, got %d", msg, data.i, got)
		}
	}
}

func TestDeconstructTag(t *testing.T) {
	datas := []struct {
		tag                 string
		datecommit, variant string
	}{
		{
			"",
			"", "",
		},
		{
			"a",
			"a", "",
		},
		{
			"a-b",
			"a-b", "",
		},
		{
			"a-b-c",
			"a-b", "c",
		},
		{
			"a-b-c-d",
			"a-b", "c-d",
		},
	}

	for _, data := range datas {
		if datecommit, variant := deconstructTag(data.tag); datecommit != data.datecommit ||
			variant != data.variant {
			log.Fatalf("deconstruct tag '%s', want: '%s, %s', got: '%s, %s'",
				data.tag, data.datecommit, data.variant, datecommit, variant)
		}
	}
}

func TestGetDominantVersion(t *testing.T) {
	datas := []struct {
		images           map[string][]versions
		dominantVersions versions
	}{
		{
			map[string][]versions{},
			versions{"", "", ""},
		},
		{
			map[string][]versions{
				"o": []versions{
					{"a-b", "h-i", ""},
				},
			},
			versions{"a-b", "h-i", ""},
		},
		{
			map[string][]versions{
				"o": []versions{
					{"a-b", "h-i", ""},
					{"c-d-x", "j-k-x", "x"},
				},
				"p": []versions{
					{"a-b-y", "h-i-y", "y"},
				},
			},
			versions{"a-b", "h-i", ""},
		},
	}

	for _, data := range datas {
		pv := PRVersions{
			images: data.images,
		}
		if dominantVersions := pv.getDominantVersions(); dominantVersions != data.dominantVersions {
			log.Fatalf("get dominant versions for '%v', want: '%v', got: '%v'", data.images, data.dominantVersions, dominantVersions)
		}
	}
}

func TestParseChangelist(t *testing.T) {
	datas := []struct {
		patches []string
		images  map[string][]versions
	}{
		{
			[]string{},
			map[string][]versions{},
		},
		{
			[]string{
				"- image: gcr.io/k8s-foofoo/bar:va-b",
				"+ image: gcr.io/k8s-foofoo/bar:vh-i",
			},
			map[string][]versions{
				"gcr.io/k8s-foofoo/bar": []versions{
					{"va-b", "vh-i", ""},
				},
			},
		},
		{
			[]string{
				"- image: gcr.io/k8s-foofoo/bar:va-b",
				"+ image: gcr.io/k8s-foofoo/bar:vh-i",
				"- image: gcr.io/k8s-foofoo/bar:va-b-x",
				"+ image: gcr.io/k8s-foofoo/bar:vh-i-x",
			},
			map[string][]versions{
				"gcr.io/k8s-foofoo/bar": []versions{
					{"va-b", "vh-i", ""},
					{"va-b-x", "vh-i-x", "x"},
				},
			},
		},
		{
			[]string{
				`- image: gcr.io/k8s-foofoo/bar:va-b
				+ image: gcr.io/k8s-foofoo/bar:vh-i
				- image: gcr.io/k8s-foofoo/bar:va-b-x
				+ image: gcr.io/k8s-foofoo/bar:vh-i-x`,
			},
			map[string][]versions{
				"gcr.io/k8s-foofoo/bar": []versions{
					{"va-b", "vh-i", ""},
					{"va-b-x", "vh-i-x", "x"},
				},
			},
		},
		{
			[]string{
				"- image: gcr.io/k8s-foofoo/bar:va-b",
				"+ image: gcr.io/k8s-foofoo/bar:vh-i",
				"- image: gcr.io/k8s-barbar/baz:vc-d",
				"+ image: gcr.io/k8s-barbar/baz:vj-k",
			},
			map[string][]versions{
				"gcr.io/k8s-foofoo/bar": []versions{
					{"va-b", "vh-i", ""},
				},
				"gcr.io/k8s-barbar/baz": []versions{
					{"vc-d", "vj-k", ""},
				},
			},
		},
	}

	for _, data := range datas {
		fgc := fakeghutil.NewFakeGithubClient()
		fcw := &GHClientWrapper{fgc}
		PR := createPullRequest(t, fgc)
		pv := PRVersions{
			images: make(map[string][]versions),
			PR:     PR,
		}
		for i, patch := range data.patches {
			SHA := strconv.Itoa(i)
			filename := fmt.Sprintf("file_%d", i)
			fgc.AddCommitToPullRequest(fakeOrg, fakeRepo, *pv.PR.Number, SHA)
			fgc.AddFileToCommit(fakeOrg, fakeRepo, SHA, filename, patch)
		}

		pv.parseChangelist(fcw)
		if eq := reflect.DeepEqual(pv.images, data.images); !eq {
			t.Fatalf("parsing PR with changes '%v', want: '%v', got: '%v'",
				data.patches, data.images, pv.images)
		}
	}
}

func TestGetBestVersion(t *testing.T) {
	type PRInfo struct {
		delta      int // hours before current date
		state      ghutil.PullRequestState
		oldVersion string
		newVersion string
	}
	datas := []struct {
		PRInfos          []PRInfo
		dominantVersions *versions
	}{
		{
			[]PRInfo{
				{7 * 24, ghutil.PullRequestCloseState, "va-b", "vh-i"},
			},
			&versions{"va-b", "vh-i", ""},
		},
		{
			[]PRInfo{
				{7 * 24, ghutil.PullRequestCloseState, "va-b", "vh-i"},
				{6 * 24, ghutil.PullRequestCloseState, "vc-d", "vj-k"},
				{8 * 24, ghutil.PullRequestCloseState, "ve-f", "vl-m"},
			},
			&versions{"va-b", "vh-i", ""},
		},
		{ // Reverted
			[]PRInfo{
				{7 * 24, ghutil.PullRequestCloseState, "va-b", "vh-i"},
				{7*24 - 3, ghutil.PullRequestCloseState, "vh-i", "va-b"}, // later PR
			},
			&versions{"vh-i", "va-b", ""},
		},
		{
			[]PRInfo{
				{7 * 24, ghutil.PullRequestOpenState, "va-b", "vh-i"},
			},
			nil,
		},
		{
			[]PRInfo{
				{4 * 24, ghutil.PullRequestCloseState, "va-b", "vh-i"},
				{10 * 24, ghutil.PullRequestCloseState, "vc-d", "vj-k"},
			},
			nil,
		},
	}

	for _, data := range datas {
		fgc := fakeghutil.NewFakeGithubClient()
		fcw := &GHClientWrapper{fgc}
		dateNow := time.Now()
		for i, PI := range data.PRInfos {
			PR := createPullRequest(t, fgc)
			timeCreated := dateNow.Add(-time.Hour * time.Duration(PI.delta))
			stateStr := string(PI.state)
			PR.State = &stateStr
			PR.CreatedAt = &timeCreated
			SHA := strconv.Itoa(i)
			filename := "fakefile"
			patch := fmt.Sprintf(`
				- image: gcr.io/k8s-foofoo/bar:%s
				+ image: gcr.io/k8s-foofoo/bar:%s
			`, PI.oldVersion, PI.newVersion)
			fgc.AddCommitToPullRequest(fakeOrg, fakeRepo, *PR.Number, SHA)
			fgc.AddFileToCommit(fakeOrg, fakeRepo, SHA, filename, patch)
		}

		pv, err := getBestVersion(fcw, fakeOrg, fakeRepo, "user:head", "base")
		if nil != err {
			t.Fatalf("get best versions with PRs '%v', want: no error, got: '%v'", data.PRInfos, err)
		}
		if nil == data.dominantVersions {
			if nil != pv && nil != pv.dominantVersions {
				t.Fatalf("get best versions with PRs '%v', want: nil, got: '%v'", data.PRInfos, pv.getDominantVersions())
			}
		} else if eq := reflect.DeepEqual(*data.dominantVersions, pv.getDominantVersions()); !eq {
			t.Fatalf("get best versions with PRs '%v', want: '%v', got: '%v'", data.PRInfos, data.dominantVersions, pv.getDominantVersions())
		}
	}
}

func TestRetryGetBestVersion(t *testing.T) {
	fgc := fakeghutil.NewFakeGithubClient()
	fcw := &GHClientWrapper{fgc}

	// only the error case exercises all the code in the function
	_, err := retryGetBestVersion(fcw, fakeOrg, fakeRepo, "user:head", "base")
	if nil == err || !strings.Contains(err.Error(), "failed list pull request") {
		t.Fatalf("retry get best version with no PRs, want error: 'failed list pull request', got: '%v'", err)
	}
}
