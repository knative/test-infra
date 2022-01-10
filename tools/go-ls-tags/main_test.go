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
	"log"
	"math"
	"os"
	"path"
	"reflect"
	"runtime"
	"sort"
	"strings"
	"testing"

	"knative.dev/test-infra/tools/go-ls-tags/cli"
)

func TestMainFunc(t *testing.T) {
	for _, tc := range testCases() {
		tc := tc
		t.Run(tc.name(), tc.test)
	}
}

func (tc *testCase) test(t *testing.T) {
	te := executeMain(tc.args)

	if te.retcode != tc.retcode {
		t.Errorf("want exit code: %d, got: %v", tc.retcode, te.retcode)
	}
	if te.errOut.String() != "" {
		t.Error("Standard error output isn't empty:\n", te.errOut.String())
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

func (tc *testCase) name() string {
	name := strings.Join(tc.args, " ")
	if name == "" {
		name = "<empty>"
	}
	return name
}

func testCases() []testCase {
	return []testCase{{
		args:      []string{"--help"},
		fragments: []string{"Usage:\n  go-ls-tags"},
	}, {
		args:    []string{"--not-existing-opt"},
		retcode: 176,
	}, {
		args: []string{},
		tags: []string{"test_tag_v2"},
	}}
}

func rootdir() string {
	_, curfile, _, _ := runtime.Caller(0) //nolint:dogsled
	return path.Dir(path.Dir(path.Dir(curfile)))
}

type testExecution struct {
	out, errOut bytes.Buffer
	retcode     int
}

type testCase struct {
	args      []string
	retcode   int
	fragments []string
	tags      []string
}

func executeMain(args []string) testExecution {
	te := testExecution{
		retcode: math.MinInt64,
	}
	opts = []cli.ExecuteOption{func(ctx *cli.ExecuteContext) {
		ctx.OsExitFunc = func(code int) {
			te.retcode = code
		}
		ctx.Args = args
		ctx.Out = &te.out
		ctx.ErrOut = &te.errOut
	}}
	defer func() {
		opts = nil
	}()
	withDirectory(rootdir(), func() {
		main()
	})

	return te
}

func withDirectory(dir string, fn func()) {
	wd, err := os.Getwd()
	if err != nil {
		log.Panic(err)
	}
	defer func() {
		_ = os.Chdir(wd)
	}()
	err = os.Chdir(dir)
	if err != nil {
		log.Panic(err)
	}
	fn()
}
