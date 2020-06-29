package main

import (
	"reflect"
	"testing"

	"knative.dev/test-infra/pkg/junit"
)

func Test_filterOutParentTests(t *testing.T) {
	tests := []struct {
		name          string
		originalCases []junit.TestCase
		want          []junit.TestCase
	}{
		{
			name: "no parent cases",
			originalCases: []junit.TestCase{{
				Name: "aaa",
			}, {
				Name: "bbb",
			}, {
				Name: "ccc",
			}},
			want: []junit.TestCase{{
				Name: "aaa",
			}, {
				Name: "bbb",
			}, {
				Name: "ccc",
			}},
		},
		{
			name: "one parent case",
			originalCases: []junit.TestCase{{
				Name: "aaa",
			}, {
				Name: "bbb",
			}, {
				Name: "bbb/ccc",
			}},
			want: []junit.TestCase{{
				Name: "aaa",
			}, {
				Name: "bbb/ccc",
			}},
		},
		{
			name: "two nested cases",
			originalCases: []junit.TestCase{{
				Name: "aaa",
			}, {
				Name: "bbb",
			}, {
				Name: "bbb/ccc",
			}, {
				Name: "bbb/ddd",
			}, {
				Name: "bbb/ddd/fff",
			}},
			want: []junit.TestCase{{
				Name: "aaa",
			}, {
				Name: "bbb/ccc",
			}, {
				Name: "bbb/ddd/fff",
			}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := filterOutParentTests(tt.originalCases); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("filterOutParentTests() = %v, want %v", got, tt.want)
			}
		})
	}
}
