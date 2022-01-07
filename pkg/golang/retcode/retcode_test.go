package retcode_test

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

import (
	"fmt"
	"io"
	"testing"

	"knative.dev/test-infra/pkg/golang/retcode"
)

func TestCalc(t *testing.T) {
	cases := testCases()
	for i := range cases {
		tt := cases[i]
		t.Run(tt.name, func(t *testing.T) {
			if got := retcode.Calc(tt.err); got != tt.want {
				t.Errorf("Calc() = %v, want %v", got, tt.want)
			}
		})
	}
}

func testCases() []testCase {
	return []testCase{{
		name: "nil",
		err:  nil,
		want: 0,
	}, {
		name: "io.ErrUnexpectedEOF",
		err:  io.ErrUnexpectedEOF,
		want: 201,
	}, {
		name: "error of wrap caused by 12345",
		err:  fmt.Errorf("%w: 12345", io.ErrUnexpectedEOF),
		want: 160,
	}}
}

type testCase struct {
	name string
	err  error
	want int
}
