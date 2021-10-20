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

package unstructured

import (
	"errors"
	"fmt"
	"reflect"
)

// ErrAsserting when provided unstructured object doesn't match assertion.
var ErrAsserting = errors.New("asserting")

// Assertion is a function that verifies unstructured object, and return error
// if found any problems with its structure.
type Assertion func(interface{}) error

// Equals returns Assertion that checks if two unstructured are equal.
func Equals(want interface{}) Assertion {
	return func(got interface{}) error {
		if !reflect.DeepEqual(got, want) {
			return fmt.Errorf("%w: %#v != %#v", ErrAsserting,
				got, want)
		}
		return nil
	}
}

// EqualsStringSlice returns an Assertion that checks if unstructured slice
// equals given string slice.
func EqualsStringSlice(want []string) Assertion {
	return func(val interface{}) error {
		got, err := toStringSlice(val)
		if err != nil {
			return err
		}
		if !reflect.DeepEqual(got, want) {
			return fmt.Errorf("%w: %#v != %#v", ErrAsserting,
				got, want)
		}
		return nil
	}
}

func toStringSlice(val interface{}) ([]string, error) {
	raw, ok := val.([]interface{})
	if !ok {
		return nil, fmt.Errorf("%w: not a slice: %#v", ErrInvalidFormat, val)
	}
	strs, err := retypeSliceToStrings(raw)
	if err != nil {
		return nil, err
	}
	return strs, nil
}

func retypeSliceToStrings(in []interface{}) ([]string, error) {
	out := make([]string, len(in))
	for i, v := range in {
		var ok bool
		out[i], ok = v.(string)
		if !ok {
			return nil, fmt.Errorf("%w: not []string: %#v", ErrInvalidFormat, in)
		}
	}
	return out, nil
}
