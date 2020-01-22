package event

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestStepType(t *testing.T) {
	s1 := Step{Number:1}

	assert.Equal(t, StepType, s1.Type())
}

func TestFinishedType(t *testing.T) {
	f1 := Finished{Count:441}

	assert.Equal(t, FinishedType, f1.Type())
}
