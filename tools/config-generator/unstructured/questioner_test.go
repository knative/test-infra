/*
Copyright 2021 The Knative Authors

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

package unstructured_test

import (
	"errors"
	"fmt"
	"testing"

	"gopkg.in/yaml.v2"
	"knative.dev/test-infra/tools/config-generator/unstructured"
)

func TestQuestioner(t *testing.T) {
	tests := []struct {
		query  string
		assert unstructured.Assertion
		want   error
	}{{
		query:  "foo.bar.0.fizz",
		assert: unstructured.EqualsStringSlice([]string{"alpha", "beta", "gamma"}),
	}, {
		query:  "foo.bar.0.fizz",
		assert: unstructured.EqualsStringSlice([]string{"alpha", "beta"}),
		want:   unstructured.ErrAsserting,
	}, {
		query:  "foo.bar.0.bazz",
		assert: unstructured.Equals(true),
	}, {
		query:  "foo.bar.0.bazz",
		assert: unstructured.Equals("yellow"),
		want:   unstructured.ErrAsserting,
	}, {
		query:  "foo.bar.0.bazz",
		assert: unstructured.EqualsStringSlice([]string{"alpha", "beta"}),
		want:   unstructured.ErrInvalidFormat,
	}, {
		query: "bla.dada",
		want:  unstructured.ErrInvalidFormat,
	}, {
		query: "foo.bar.42.bazz",
		want:  unstructured.ErrInvalidFormat,
	}, {
		query: "foo.42.bazz",
		want:  unstructured.ErrInvalidFormat,
	}, {
		query: "foo.bla",
		want:  unstructured.ErrInvalidFormat,
	}, {
		query: "foo.bar.bla",
		want:  unstructured.ErrInvalidFormat,
	}, {
		query:  "foo.bar",
		assert: unstructured.EqualsStringSlice([]string{"alpha", "beta", "gamma"}),
		want:   unstructured.ErrInvalidFormat,
	}}
	for i, tc := range tests {
		tc := tc
		t.Run(fmt.Sprintf("%d-%s", i, tc.query), func(t *testing.T) {
			err := testQuestionerQuery(t, tc.query, tc.assert)
			checkErr(t, err, tc.want)
		})
	}
}

func exampleUnstructured(tb testing.TB) interface{} {
	tb.Helper()
	un := make(map[interface{}]interface{})
	err := yaml.Unmarshal([]byte(`---
foo:
  bar:
  - fizz:
    - alpha
    - beta
    - gamma
    bazz: true
`), &un)
	if err != nil {
		tb.Fatal(err)
	}
	return un
}

func testQuestionerQuery(tb testing.TB, query string, assert unstructured.Assertion) error {
	tb.Helper()
	questioner := unstructured.NewQuestioner(exampleUnstructured(tb))
	val, err := questioner.Query(query)
	if err != nil {
		return err
	}
	if assert != nil {
		return assert(val)
	}
	return nil
}

func checkErr(tb testing.TB, got, want error) {
	tb.Helper()
	if want == nil {
		if got != nil {
			tb.Fatal(got)
		}
	}
	if !errors.Is(got, want) {
		tb.Fatalf("got: %#v, want: %#v", got, want)
	}
}
