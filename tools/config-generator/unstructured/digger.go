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
	"strconv"
)

// ErrInvalidFormat when provided unstructured object has invalid format.
var ErrInvalidFormat = errors.New("invalid format")

// Digger is a function that digs in unstructured object and returns some sub
// element of that object, or an error if such sub object can't be located.
type Digger func(interface{}) (interface{}, error)

// MapKey returns a Digger that looks up the value of a key within the map.
func MapKey(key string) Digger {
	return func(un interface{}) (interface{}, error) {
		m, ok := un.(map[interface{}]interface{})
		if !ok {
			return nil, fmt.Errorf("%w: not a map: %#v", ErrInvalidFormat, un)
		}
		val, ok := m[key]
		if !ok {
			return nil, fmt.Errorf("%w: no key %#v in map: %#v",
				ErrInvalidFormat, key, un)
		}
		return val, nil
	}
}

// SliceElem returns a Digger that looks up the value of slice under provided index.
func SliceElem(idx int) Digger {
	return func(un interface{}) (interface{}, error) {
		s, ok := un.([]interface{})
		if !ok {
			return nil, fmt.Errorf("%w: not a slice: %#v", ErrInvalidFormat, un)
		}
		if idx < 0 || idx >= len(s) {
			return nil, fmt.Errorf(
				"%w: index out of range [%d] for %#v",
				ErrInvalidFormat, idx, s)
		}
		return s[idx], nil
	}
}

func toDiggers(queries []string) []Digger {
	digrs := make([]Digger, len(queries))
	for i, query := range queries {
		idx, err := strconv.Atoi(query)
		var next Digger
		if err != nil {
			next = MapKey(query)
		} else {
			next = SliceElem(idx)
		}
		digrs[i] = next
	}
	return digrs
}
