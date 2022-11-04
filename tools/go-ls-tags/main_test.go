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
	"bytes"
	"reflect"
	"strings"
	"testing"

	"github.com/wavesoftware/go-commandline"
	"knative.dev/test-infra/tools/go-ls-tags/cli"
	"knative.dev/test-infra/tools/go-ls-tags/test"
)

func TestMainFunc(t *testing.T) {
	tcs := []testCase{{
		name:      "help",
		args:      []string{"--help"},
		fragments: []string{"Usage:\n  go-ls-tags"},
	}, {
		name:    "not existing option",
		args:    []string{"--not-existing-opt"},
		retcode: 176,
	}, {
		name: "no options",
		args: []string{},
		tags: []string{"test_tag_v1", "test_tag_v2", "test_tag_v3"},
	}}
	for _, tc := range tcs {
		tc := tc
		t.Run(tc.name, tc.test)
	}
}

type testCase struct {
	name      string
	args      []string
	retcode   int
	fragments []string
	tags      []string
}

func (tc *testCase) test(t *testing.T) {
	te := executeMain(tc.args)

	if te.retcode != tc.retcode {
		t.Errorf("want exit code: %d, got: %v", tc.retcode, te.retcode)
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
		if !reflect.DeepEqual(got, want) {
			t.Errorf("mismatch\nwant %q\n got %q", want, got)
		}
	}
}

type testExecution struct {
	out     bytes.Buffer
	retcode int
}

func executeMain(args []string) testExecution {
	te := testExecution{}
	cli.Options = []commandline.Option{
		commandline.WithArgs(args...),
		commandline.WithOutput(&te.out),
		commandline.WithExit(func(code int) {
			te.retcode = code
		}),
	}
	defer func() {
		cli.Options = nil
	}()
	test.WithDirectory(test.Rootdir(), func() {
		main()
	})

	return te
}
