package main_test

import (
	"bytes"
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/wavesoftware/go-commandline"
	main "knative.dev/test-infra/tools/modscope"
)

func TestMainFn(t *testing.T) {
	var buf bytes.Buffer
	retcode := math.MinInt16
	main.RunMain(
		commandline.WithArgs("current"),
		commandline.WithOutput(&buf),
		commandline.WithExit(func(code int) {
			retcode = code
		}),
	)

	assert.Equal(t, math.MinInt16, retcode)
	assert.Equal(t, "knative.dev/test-infra\n", buf.String())
}
