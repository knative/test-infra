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

package cli_test

import (
	"bytes"
	"errors"
	"path"
	"reflect"
	"sort"
	"strings"
	"testing"

	"github.com/wavesoftware/go-commandline"
	"knative.dev/test-infra/tools/go-ls-tags/cli"
	"knative.dev/test-infra/tools/go-ls-tags/tags"
	"knative.dev/test-infra/tools/go-ls-tags/test"
)

func TestApp(t *testing.T) {
	t.Parallel()
	tcs := []testCase{{
		name:      "help",
		args:      []string{"--help"},
		fragments: []string{"Usage:\n  go-ls-tags"},
	}, {
		name: "not existing option",
		args: []string{"--not-existing-opt"},
		err:  errors.New("unknown flag: --not-existing-opt"),
	}, {
		name: "no options",
		args: []string{},
		tags: []string{"test_tag_v1", "test_tag_v2"},
	}, {
		name: "using absolute ignorefile",
		args: []string{
			"--ignore-file",
			path.Join(test.Rootdir(), "tools", "go-ls-tags", tags.DefaultIgnoreFile),
		},
		tags: []string{"test_tag_v2"},
	}, {
		name: "using relative ignorefile",
		args: []string{
			"--ignore-file",
			path.Join("tools", "go-ls-tags", tags.DefaultIgnoreFile),
		},
		tags: []string{"test_tag_v2"},
	}}
	for _, tc := range tcs {
		tc := tc
		t.Run(tc.name, tc.test)
	}
}

type testCase struct {
	name      string
	args      []string
	err       error
	fragments []string
	tags      []string
}

func (tc *testCase) test(t *testing.T) {
	t.Parallel()
	te := execute(tc.args)

	if tc.err != nil {
		if !(errors.Is(tc.err, te.err) || tc.err.Error() == te.err.Error()) {
			t.Errorf("want error: %#v, got: %#v", tc.err, te.err)
		}
	} else {
		if te.err != nil {
			t.Errorf("want no error, got: %#v", te.err)
		}
	}
	text := te.out.String()
	for _, fragment := range tc.fragments {
		if !strings.Contains(text, fragment) {
			t.Errorf("fragment %#v not found in %#v", fragment, text)
		}
	}
	if len(tc.tags) > 0 {
		want := tc.tags
		got := strings.Split(strings.Trim(text, "\n"), "\n")
		sort.Strings(want)
		sort.Strings(got)
		if !reflect.DeepEqual(got, want) {
			t.Errorf("mismatch\nwant %q\n got %q", want, got)
		}
	}
}

func execute(args []string) testExecution {
	te := testExecution{}
	test.WithDirectory(test.Rootdir(), func() {
		err := commandline.New(cli.App{}).Execute(
			commandline.WithArgs(args...),
			commandline.WithOutput(&te.out),
		)
		te.err = err
	})
	return te
}

type testExecution struct {
	out bytes.Buffer
	err error
}
