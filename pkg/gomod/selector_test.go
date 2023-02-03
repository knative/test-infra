package gomod_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"knative.dev/test-infra/pkg/gomod"
)

func TestDefaultSelector(t *testing.T) {
	sel, err := gomod.DefaultSelector("knative.dev")
	require.NoError(t, err)
	mods := []string{
		"knative.dev/test-infra",
		"knative.dev/serving",
		"knative.dev/eventing",
		"knative.dev/pkg",
		"github.com/blang/semver/v4",
		"github.com/google/go-cmp",
		"go.uber.org/atomic",
	}
	selected := make([]string, 0, len(mods))
	for _, mod := range mods {
		if sel(mod) {
			selected = append(selected, mod)
		}
	}
	assert.Equal(t, []string{
		"knative.dev/serving",
		"knative.dev/eventing",
		"knative.dev/pkg",
	}, selected)
}
