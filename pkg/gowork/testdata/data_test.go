package testdata_test

import (
	"io/fs"
	"testing"

	"github.com/stretchr/testify/assert"
	"knative.dev/test-infra/pkg/gowork/testdata"
)

func TestExampleFS(t *testing.T) {
	t.Parallel()
	tfs := testdata.ExampleFS()
	bytes, err := fs.ReadFile(tfs, "code/foo/go.mod")
	assert.NoError(t, err)
	assert.Equal(t, "module foo\n\ngo 1.18\n", string(bytes))
}
