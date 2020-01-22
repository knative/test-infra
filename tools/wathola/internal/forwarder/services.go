package forwarder

import (
	"context"
	"github.com/cloudevents/sdk-go"
	log "github.com/sirupsen/logrus"
	"knative.dev/test-infra/tools/wathola/internal/client"
	"knative.dev/test-infra/tools/wathola/internal/config"
	"knative.dev/test-infra/tools/wathola/internal/sender"
	"time"
)

var lastProgressReport = time.Now()

// New creates new forwarder
func New() Forwarder {
	config.ReadIfPresent()
	f := &forwarder{
		count: 0,
	}
	return f
}

// Stop will stop running forwarder if there is one
func Stop() {
	if cancel != nil {
		log.Info("stopping forwarder")
		cancel()
		cancel = nil
	}
}

var cancel context.CancelFunc

func (f *forwarder) Forward() {
	port := config.Instance.Forwarder.Port
	client.Receive(port, &cancel, f.forwardEvent)
}

func (f *forwarder) forwardEvent(e cloudevents.Event) {
	target := config.Instance.Forwarder.Target
	log.Tracef("Forwarding event %v to %v", e.ID(), target)
	err := sender.SendEvent(e, target)
	if err != nil {
		log.Error(err)
	}
	f.count++
	f.reportProgress()
}

func (f *forwarder) reportProgress() {
	if lastProgressReport.Add(config.Instance.Receiver.Progress.Duration).Before(time.Now()) {
		lastProgressReport = time.Now()
		log.Infof("forwarded %v events", f.count)
	}
}

type forwarder struct {
	count int
}
