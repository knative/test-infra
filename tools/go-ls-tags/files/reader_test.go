package files_test

import (
	"context"
	"path"
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"knative.dev/test-infra/tools/go-ls-tags/files"
)

func TestReadlines(t *testing.T) {
	t.Parallel()
	testfile := path.Join(dir(), "testdata", "example.txt")
	lines := make([]string, 0)
	err := files.ReadLines(context.TODO(), testfile, func(line string) error {
		lines = append(lines, line)
		return nil
	})
	require.NoError(t, err)
	assert.Equal(t, []string{"foo", "bar"}, lines)
}

func dir() string {
	_, file, _, _ := runtime.Caller(0)
	return path.Dir(file)
}
