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

package main

import (
	"sort"
	"strconv"
	"strings"

	"gopkg.in/yaml.v2"
)

// getString casts the given interface (expected string) as string.
// An array of length 1 is also considered a single string.
func getString(s interface{}) string {
	if _, ok := s.([]interface{}); ok {
		values := getStringArray(s)
		if len(values) == 1 {
			return values[0]
		}
		logFatalf("Entry %v is not a string or string array of size 1", s)
	}
	if str, ok := s.(string); ok {
		return str
	}
	logFatalf("Entry %v is not a string", s)
	return ""
}

// getInt casts the given interface (expected int) as int.
func getInt(s interface{}) int {
	if value, ok := s.(int); ok {
		return value
	}
	logFatalf("Entry %v is not an integer", s)
	return 0
}

// getBool casts the given interface (expected bool) as bool.
func getBool(s interface{}) bool {
	if value, ok := s.(bool); ok {
		return value
	}
	logFatalf("Entry %v is not a boolean", s)
	return false
}

// getInterfaceArray casts the given interface (expected interface array) as interface array.
func getInterfaceArray(s interface{}) []interface{} {
	if interfaceArray, ok := s.([]interface{}); ok {
		return interfaceArray
	}
	logFatalf("Entry %v is not an interface array", s)
	return nil
}

// getStringArray casts the given interface (expected string array) as string array.
func getStringArray(s interface{}) []string {
	interfaceArray := getInterfaceArray(s)
	strArray := make([]string, len(interfaceArray))
	for i := range interfaceArray {
		strArray[i] = getString(interfaceArray[i])
	}
	return strArray
}

// getMapSlice casts the given interface (expected MapSlice) as MapSlice.
func getMapSlice(m interface{}) yaml.MapSlice {
	if mm, ok := m.(yaml.MapSlice); ok {
		return mm
	}
	logFatalf("Entry %v is not a yaml.MapSlice", m)
	return nil
}

// appendIfUnique appends an element to an array of strings, unless it's already present.
func appendIfUnique(a1 []string, e2 string) []string {
	var res []string
	res = append(res, a1...)
	for _, e1 := range a1 {
		if e1 == e2 {
			return res
		}
	}
	return append(res, e2)
}

// isNum checks if the given string is a valid number
func isNum(s string) bool {
	_, err := strconv.ParseFloat(s, 64)
	return err == nil
}

// quote returns the given string quoted if it's not a number, or not a key/value pair, or already quoted.
func quote(s string) string {
	if isNum(s) {
		return s
	}
	if strings.HasPrefix(s, "'") || strings.HasPrefix(s, "\"") || strings.Contains(s, ": ") || strings.HasSuffix(s, ":") {
		return s
	}
	return "\"" + s + "\""
}

// indentBase is a helper function which returns the given array indented.
func indentBase(indentation int, prefix string, indentFirstLine bool, array []string) string {
	s := ""
	if len(array) == 0 {
		return s
	}
	indent := strings.Repeat(" ", indentation)
	for i := 0; i < len(array); i++ {
		if i > 0 || indentFirstLine {
			s += indent
		}
		s += prefix + quote(array[i]) + "\n"
	}
	return s
}

// indentArray returns the given array indented, prefixed by "-".
func indentArray(indentation int, array []string) string {
	return indentBase(indentation, "- ", false, array)
}

// indentKeys returns the given array of key/value pairs indented.
func indentKeys(indentation int, array []string) string {
	return indentBase(indentation, "", false, array)
}

// indentSectionBase is a helper function which returns the given array of key/value pairs indented inside a section.
func indentSectionBase(indentation int, title string, prefix string, array []string) string {
	keys := indentBase(indentation, prefix, true, array)
	if keys == "" {
		return keys
	}
	return title + ":\n" + keys
}

// indentArraySection returns the given array indented inside a section.
func indentArraySection(indentation int, title string, array []string) string {
	return indentSectionBase(indentation, title, "- ", array)
}

// indentSection returns the given array of key/value pairs indented inside a section.
func indentSection(indentation int, title string, array []string) string {
	return indentSectionBase(indentation, title, "", array)
}

// indentMap returns the given map indented, with each key/value separated by ": "
func indentMap(indentation int, mp map[string]string) string {
	// Extract map keys to keep order consistent.
	keys := make([]string, 0, len(mp))
	for key := range mp {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	arr := make([]string, len(mp))
	for i := 0; i < len(mp); i++ {
		arr[i] = keys[i] + ": " + quote(mp[keys[i]])
	}
	return indentBase(indentation, "", false, arr)
}

// strExists checks if the given string exists in the array
func strExists(arr []string, str string) bool {
	for _, s := range arr {
		if str == s {
			return true
		}
	}
	return false
}
