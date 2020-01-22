package config

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDefaultValues(t *testing.T) {
	assert.NotNil(t, Instance)
	assert.Condition(t, func() (success bool) {
		return Instance.Receiver.Teardown.Duration.Seconds() >= 1
	})
}
