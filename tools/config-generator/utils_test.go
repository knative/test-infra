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
	"fmt"
	"testing"

	"github.com/google/go-cmp/cmp"
	"gopkg.in/yaml.v2"
)

func TestGetString(t *testing.T) {
	var in interface{} = "abcdefg"
	out := getString(in)
	if diff := cmp.Diff(out, "abcdefg"); diff != "" {
		t.Fatalf("Unexpected output (-got +want):\n%s", diff)
	}
	if logFatalCalls != 0 {
		t.Fatalf("logFatal was called for %v", in)
	}

	out = getString(42)
	if logFatalCalls != 1 {
		t.Fatalf("logFatal was not called for %v", in)
	}
}

func TestGetInt(t *testing.T) {
	var in interface{} = 123
	out := getInt(in)
	if logFatalCalls != 0 {
		t.Fatalf("logFatal was called for %v", in)
	}
	if out != 123 {
		t.Fatalf("Expected 123, got %v", out)
	}

	getInt("abc")
	if logFatalCalls == 0 {
		t.Fatalf("Expected logFatal to be called")
	}
}

func TestGetBool(t *testing.T) {
	var in interface{} = true
	out := getBool(in)
	if logFatalCalls != 0 {
		t.Fatalf("logFatal was called for %v", in)
	}
	if !out {
		t.Fatalf("Expected true, got %v", out)
	}

	getBool(123)
	if logFatalCalls == 0 {
		t.Fatalf("Expected logFatal to be called")
	}
}

func TestGetInterfaceArray(t *testing.T) {
	in1 := []interface{}{"foo", "bar", "baz"}
	out1 := getInterfaceArray(in1)
	if fmt.Sprint(in1) != fmt.Sprint(out1) {
		t.Fatalf("Did not get same interface slice back.")
	}
	if logFatalCalls != 0 {
		t.Fatalf("Interface slice caused logFatal call")
	}

	in2 := []string{"foo", "bar", "baz"}
	getInterfaceArray(in2)
	if logFatalCalls != 1 {
		t.Fatalf("Non interface slice should have caused logFatal call")
	}
}

func TestGetStringArray(t *testing.T) {
	in := []interface{}{"foo", "bar", "baz"}
	out := getStringArray(in)
	if logFatalCalls != 0 {
		t.Fatalf("Input %v should not have caused logFatal call.", in)
	}
	if fmt.Sprint(out) != fmt.Sprint(in) {
		t.Fatalf("Expected input %v and output %v to have identical string output.", in, out)
	}
}

func TestGetMapSlice(t *testing.T) {
	var in interface{} = yaml.MapSlice{
		yaml.MapItem{Key: "abc", Value: 123},
		yaml.MapItem{Key: "def", Value: 456},
	}
	out := getMapSlice(in)
	if logFatalCalls != 0 {
		t.Fatalf("Input %v should not have caused logFatal call.", in)
	}
	if fmt.Sprint(out) != fmt.Sprint(in) {
		t.Fatalf("Expected input %v and output %v to have identical string output.", in, out)
	}
}

func TestAppendIfUnique(t *testing.T) {
	arr := []string{"foo", "bar"}
	arr = appendIfUnique(arr, "foo")
	if len(arr) != 2 {
		t.Fatalf("Expected length 2 but was %v", len(arr))
	}
	arr = appendIfUnique(arr, "baz")
	if arr[2] != "baz" {
		t.Fatalf("Expected 'baz' to be appended but wasn't.")
	}
}

func TestIsNum(t *testing.T) {
	nums := []string{"-123456.789", "-123", "0", "0.0", ".0", "123", "123456.789"}
	for _, n := range nums {
		if !isNum(n) {
			t.Fatalf("Input %v should be a num, but wasn't.", n)
		}
	}
	notNums := []string{"", ".", "abc", "123 "}
	for _, n := range notNums {
		if isNum(n) {
			t.Fatalf("Input %v should not be a num, but was identified as one.", n)
		}
	}
}

func TestQuote(t *testing.T) {
	tests := []struct {
		in           string
		expectQuotes bool
	}{
		{"foo bar baz", true},
		{"", true},
		{"\"foo bar\"", false},
		{"'foo bar'", false},
		{"123", false},
		{"abc:def", true}, // Not recognized as a key value pair without space after colon
		{"abc: def", false},
		{"abc:", false},
	}
	for _, test := range tests {
		out := quote(test.in)
		quoted := "\"" + test.in + "\""
		if test.expectQuotes && out != "\""+test.in+"\"" {
			t.Fatalf("Expected %v, got %v", quoted, out)
		}
		if !test.expectQuotes && test.in != out {
			t.Fatalf("Expected %v, got %v", test.in, out)
		}
	}
}

func TestIndentBase(t *testing.T) {
	tests := []struct {
		input           []string
		indentation     int
		prefix          string
		indentFirstLine bool
		expected        string
	}{
		{
			input:           []string{"foo", "bar", "baz"},
			indentation:     2,
			prefix:          "",
			indentFirstLine: false,
			expected:        fmt.Sprintf("%q\n  %q\n  %q\n", "foo", "bar", "baz"),
		},
		{
			input:           []string{"foo", "bar", "baz"},
			indentation:     0,
			prefix:          "",
			indentFirstLine: false,
			expected:        fmt.Sprintf("%q\n%q\n%q\n", "foo", "bar", "baz"),
		},
		{
			input:           []string{"foo", "bar", "baz"},
			indentation:     2,
			prefix:          "",
			indentFirstLine: true,
			expected:        fmt.Sprintf("  %q\n  %q\n  %q\n", "foo", "bar", "baz"),
		},
		{
			input:           []string{"foo", "bar", "baz"},
			indentation:     2,
			prefix:          "__",
			indentFirstLine: false,
			expected:        fmt.Sprintf("__%q\n  __%q\n  __%q\n", "foo", "bar", "baz"),
		},
	}
	for _, test := range tests {
		out := indentBase(
			test.indentation,
			test.prefix,
			test.indentFirstLine,
			test.input)
		if diff := cmp.Diff(out, test.expected); diff != "" {
			t.Fatalf("Unexpected output (-got +want):\n%s", diff)
		}
	}
}

func TestIndentArray(t *testing.T) {
	input := []string{"'foo'", "42", "key: value", "bar"}
	indentation := 2
	expected := "- 'foo'\n  - 42\n  - key: value\n  - \"bar\"\n"

	if diff := cmp.Diff(indentArray(indentation, input), expected); diff != "" {
		t.Fatalf("Unexpected output (-got +want):\n%s", diff)
	}
}

func TestIndentKeys(t *testing.T) {
	input := []string{"abc: def", "foo: bar"}
	indentation := 2
	expected := "abc: def\n  foo: bar\n"

	if diff := cmp.Diff(indentKeys(indentation, input), expected); diff != "" {
		t.Fatalf("Unexpected output (-got +want):\n%s", diff)
	}
}

func TestIndentSectionBase(t *testing.T) {
	indentation := 2
	title := "foo"
	prefix := "__"
	input := []string{"abc: def", "bar", "42"}
	expected := "foo:\n  __abc: def\n  __\"bar\"\n  __42\n"

	out := indentSectionBase(indentation, title, prefix, input)
	if diff := cmp.Diff(out, expected); diff != "" {
		t.Fatalf("Unexpected output (-got +want):\n%s", diff)
	}

	out = indentSectionBase(indentation, title, prefix, []string{})
	if diff := cmp.Diff(out, ""); diff != "" {
		t.Fatalf("Unexpected output (-got +want):\n%s", diff)
	}
}

func TestIndentArraySection(t *testing.T) {
	indentation := 2
	title := "foo"
	input := []string{"abc: def", "bar", "42"}
	expected := "foo:\n  - abc: def\n  - \"bar\"\n  - 42\n"

	out := indentArraySection(indentation, title, input)
	if diff := cmp.Diff(out, expected); diff != "" {
		t.Fatalf("Unexpected output (-got +want):\n%s", diff)
	}

	out = indentArraySection(indentation, title, []string{})
	if diff := cmp.Diff(out, ""); diff != "" {
		t.Fatalf("Unexpected output (-got +want):\n%s", diff)
	}
}

func TestIndentSection(t *testing.T) {
	indentation := 2
	title := "foo"
	input := []string{"abc: def", "bar: baz", "magic_num: 42"}
	expected := "foo:\n  abc: def\n  bar: baz\n  magic_num: 42\n"

	out := indentSection(indentation, title, input)
	if diff := cmp.Diff(out, expected); diff != "" {
		t.Fatalf("Unexpected output (-got +want):\n%s", diff)
	}

	out = indentSection(indentation, title, []string{})
	if diff := cmp.Diff(out, ""); diff != "" {
		t.Fatalf("Unexpected output (-got +want):\n%s", diff)
	}
}

func TestIndentMap(t *testing.T) {
	indentation := 2
	input := map[string]string{
		"foo": "bar",
		"abc": "def",
		"num": "42",
	}
	expected := "abc: \"def\"\n  foo: \"bar\"\n  num: 42\n"

	out := indentMap(indentation, input)
	if diff := cmp.Diff(out, expected); diff != "" {
		t.Fatalf("Unexpected output (-got +want):\n%s", diff)
	}
}

func TestStrExists(t *testing.T) {
	sArray := []string{"foo", "bar", "baz"}

	if strExists(sArray, "abc") {
		t.Fatalf("String abc should not exist in %v", sArray)
	}

	if !strExists(sArray, "bar") {
		t.Fatalf("String bar should exist in %v", sArray)
	}
}
