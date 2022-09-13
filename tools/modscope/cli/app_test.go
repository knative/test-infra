package cli_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"knative.dev/test-infra/tools/modscope/cli"
)

func TestApp(t *testing.T) {
	t.Parallel()
	app := cli.App{}
	c := app.Command()

	assert.Equal(t, "modscope", c.Use)
}
