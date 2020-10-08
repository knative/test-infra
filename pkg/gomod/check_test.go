/*
Copyright 2020 The Knative Authors

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

package gomod

import (
	"errors"
	"testing"
)

// It is not worth testing Check, the method wires other the results of well
// tested dependencies and would require mocking a http go-import server and
// a git server.

func TestError(t *testing.T) {
	tests := map[string]struct {
		err             error
		isDependencyErr bool
	}{
		"true, empty": {
			err:             &Error{},
			isDependencyErr: true,
		},
		"true, filled": {
			err: &Error{
				Module:       "foo",
				Dependencies: []string{"bar", "baz"},
			},
			isDependencyErr: true,
		},
		"false": {
			err:             errors.New("not a dep error"),
			isDependencyErr: false,
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			if errors.Is(tt.err, DependencyErr) != tt.isDependencyErr {
				t.Errorf("expected errors.Is(err, DependencyErr) to be %t", tt.isDependencyErr)
			}
		})
	}
}

func TestError_Error(t *testing.T) {
	tests := map[string]struct {
		err  error
		want string
	}{
		"empty": {
			err:  &Error{},
			want: " failed because of the following dependencies []",
		},
		"filled": {
			err: &Error{
				Module:       "foo",
				Dependencies: []string{"bar", "baz"},
			},
			want: "foo failed because of the following dependencies [bar, baz]",
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			if got := tt.err.Error(); tt.want != got {
				t.Errorf("err.Error() = %v, want %v", got, tt.want)
			}
		})
	}
}
