package main

import (
	"testing"
)

func TestFileListCommentString(t *testing.T) {
	tests := []struct {
		name     string
		files    []string
		expected string
	}{{
		name:     "empty list",
		files:    make([]string, 0),
		expected: "",
	}, {
		name:     "list with one string",
		files:    []string{"file1.yaml"},
		expected: "- `file1.yaml`",
	}, {
		name:     "list with multiple strings",
		files:    []string{"file1.yaml", "file2.yaml", "file3.yaml"},
		expected: "- `file1.yaml`\n- `file2.yaml`\n- `file3.yaml`",
	},
	}

	for _, test := range tests {
		res := fileListCommentString(test.files)
		if res != test.expected {
			t.Fatalf("expect: %q, actual: %q", test.expected, res)
		}
	}
}
