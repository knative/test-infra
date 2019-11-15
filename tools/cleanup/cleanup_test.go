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
	"github.com/google/go-cmp/cmp"
	"github.com/pkg/errors"
	"strings"
	"testing"
	"time"
)

var (
	oneYearAgo = time.Now().Add(-24 * 365 * time.Hour)
)

// errorMismatch compares two errors by their error string
func errorMismatch(got error, want error) string {
	exp := "<no error>"
	act := "<no error>"
	if want != nil {
		exp = want.Error()
	}
	if got != nil {
		act = got.Error()
	}
	if strings.Contains(act, exp) {
		return ""
	}
	return fmt.Sprintf("got: '%v'\nwant: '%v'", act, exp)
}

func TestSelectProjects(t *testing.T) {
	datas := []struct {
		projectFlag  string
		yamlFileFlag string
		regexFlag    string
		exp          []string
		err          error
	}{
		{ // Project provided.
			"foo",
			"",
			"",
			[]string{"foo"},
			nil,
		},
		{ // File provided.
			"",
			"testdata/resources.yaml",
			"knative-boskos-.*",
			[]string{"knative-boskos-01", "knative-boskos-02"},
			nil,
		},
		{ // Bad file provided.
			"",
			"/foobar_resources.yamlfoo",
			"",
			[]string{},
			errors.New("no such file or directory"),
		},
		{ // Empty file provided.
			"",
			"testdata/empty.yaml",
			".*",
			[]string{},
			errors.New("no project found"),
		},
		{ // Bad regex provided.
			"",
			"testdata/resources.yaml",
			"--->}][{<---",
			[]string{},
			errors.New("invalid character class range"),
		},
		{ // Unmatching regex provided.
			"",
			"testdata/resources.yaml",
			"foobar-[0-9]",
			[]string{},
			errors.New("no project found"),
		},
	}
	for _, data := range datas {
		r, err := selectProjects(data.projectFlag, data.yamlFileFlag, data.regexFlag)
		errMsg := fmt.Sprintf("Select projects for %q/%q/%q: ", data.projectFlag, data.yamlFileFlag, data.regexFlag)
		if m := errorMismatch(err, data.err); m != "" {
			t.Errorf("%s%s", errMsg, m)
		}
		if data.err != nil {
			continue
		}
		if dif := cmp.Diff(data.exp, r); dif != "" {
			t.Errorf("%sgot(+) is different from wanted(-)\n%v", errMsg, dif)
		}
	}
}

func TestShowStats(t *testing.T) {
	showStats(0, []string{})
	showStats(1, []string{"foo"})
}

func TestDeleteResources(t *testing.T) {
	datas := []struct {
		projects []string
		before   time.Time
		fn       resourceDeleter
		expCount int
		expErr   []string
	}{
		{ // No projects.
			[]string{},
			oneYearAgo,
			func(string, time.Time) (int, error) {
				t.Error("delete function should not be called")
				return 0, nil
			},
			0,
			[]string{},
		},
		{ // With less than 10 projects.
			[]string{"p1", "p2"},
			time.Now(),
			func(string, time.Time) (int, error) {
				return 1, nil
			},
			2,
			[]string{},
		},
		{ // With more than 10 projects.
			[]string{"p1", "p2", "p3", "p4", "p5", "p6", "p7", "p8", "p9", "pA", "pB", "pC", "pD"},
			time.Now(),
			func(string, time.Time) (int, error) {
				return 1, nil
			},
			13,
			[]string{},
		},
		{ // With projects, but errors.
			[]string{"p1e", "p2e", "p3e"},
			time.Now(),
			func(p string, t time.Time) (int, error) {
				return 1, errors.New("error")
			},
			-1, // There are errors, returned count might vary, ignore.
			[]string{"error"},
		},
	}

	for _, data := range datas {
		c, err := deleteResources(data.projects, data.before, data.fn)
		errMsg := fmt.Sprintf("Delete projects %q before %v: ", data.projects, data.before)
		if c != data.expCount && data.expCount > -1 {
			t.Errorf("%sgot %d, wanted %d", errMsg, c, data.expCount)
		}
		if dif := cmp.Diff(data.expErr, err); dif != "" {
			t.Errorf("%sgot(+) is different from wanted(-)\n%v", errMsg, dif)
		}
	}
}

/*
  TODO(adrcunha): Test the other functions.
  TODO(adrcunha): Test default flag values.

  deleteImage (requires mocking remote)
  - bad reference
  - error deleting image
  - delete image

  deleteClusters (requires mocking gkeClient)
  - bad project
  - error listing clusters
  - bad timestamp
  - delete cluster (dry run)
  - delete cluster
  - error deleting cluster

  deleteImages (requires mocking name, google)
  - bad project
  - bad registry
  - error walking down registry
  - no images to delete
  - error deleting tag
  - error deleting image

  cleanup (requires mocking common, gke)
  - bad flags
  - error creating GKE client
  - error activating account
  - skipping deleting images
  - skipping deleting clusters
  - deleting images and clusters
*/
