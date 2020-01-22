package main

import (
	"github.com/phayes/freeport"
	"github.com/stretchr/testify/assert"
	"knative.dev/test-infra/tools/wathola/internal/config"
	"knative.dev/test-infra/tools/wathola/internal/forwarder"
	"testing"
	"time"
)

func TestForwarderMain(t *testing.T) {
	config.Instance.Forwarder.Port = freeport.GetPort()
	go main()
	defer forwarder.Stop()

	time.Sleep(time.Second)

	assert.NotNil(t, instance)
}
