package event

import (
	"github.com/stretchr/testify/assert"
	"knative.dev/test-infra/tools/wathola/internal/config"
	"os"
	"testing"
	"time"
)

func TestProperEventsPropagation(t *testing.T) {
	// given
	errors := NewErrorStore()
	stepsStore := NewStepsStore(errors)
	finishedStore := NewFinishedStore(stepsStore, errors)

	// when
	stepsStore.RegisterStep(&Step{Number: 1})
	stepsStore.RegisterStep(&Step{Number: 3})
	stepsStore.RegisterStep(&Step{Number: 2})
	finishedStore.RegisterFinished(&Finished{Count: 3})

	// then
	assert.Empty(t, errors.thrown)
}

func TestMissingAndDoubleEvent(t *testing.T) {
	// given
	errors := NewErrorStore()
	stepsStore := NewStepsStore(errors)
	finishedStore := NewFinishedStore(stepsStore, errors)

	// when
	stepsStore.RegisterStep(&Step{Number: 1})
	stepsStore.RegisterStep(&Step{Number: 2})
	stepsStore.RegisterStep(&Step{Number: 2})
	finishedStore.RegisterFinished(&Finished{Count: 3})

	// then
	assert.NotEmpty(t, errors.thrown)
}

func TestDoubleFinished(t *testing.T) {
	// given
	errors := NewErrorStore()
	stepsStore := NewStepsStore(errors)
	finishedStore := NewFinishedStore(stepsStore, errors)

	// when
	stepsStore.RegisterStep(&Step{Number: 1})
	stepsStore.RegisterStep(&Step{Number: 2})
	finishedStore.RegisterFinished(&Finished{Count: 2})
	finishedStore.RegisterFinished(&Finished{Count: 2})

	// then
	assert.NotEmpty(t, errors.thrown)
}

func TestMain(m *testing.M) {
	config.Instance.Receiver.Teardown.Duration = 20 * time.Millisecond
	exitcode := m.Run()
	os.Exit(exitcode)
}
