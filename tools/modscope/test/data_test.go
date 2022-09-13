package test_test

import (
	"io/fs"
	"testing"

	"github.com/stretchr/testify/assert"
	"knative.dev/test-infra/tools/modscope/test"
)

func TestExampleFS(t *testing.T) {
	t.Parallel()
	tfs := test.ExampleFS()
	bytes, err := fs.ReadFile(tfs, "code/foo/go.mod")
	assert.NoError(t, err)
	assert.Equal(t, "module foo\n\ngo 1.18\n", string(bytes))
}
