package sender

import (
	"context"
	"fmt"
	cloudevents "github.com/cloudevents/sdk-go"
	log "github.com/sirupsen/logrus"
	"knative.dev/test-infra/tools/wathola/internal/config"
	"knative.dev/test-infra/tools/wathola/internal/ensure"
	"knative.dev/test-infra/tools/wathola/internal/event"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var senderConfig = &config.Instance.Sender

type sender struct {
	counter int
	active  bool
}

func (s *sender) SendContinually() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		for sig := range c {
			// sig is a ^C or term, handle it
			log.Infof("Received: %v, closing", sig.String())
			s.active = false
			s.sendFinished()
		}
	}()
	for s.active {
		err := s.sendStep()
		if err != nil {
			log.Warnf("Could not send step event, retry in %v", senderConfig.Cooldown)
			time.Sleep(senderConfig.Cooldown)
		} else {
			time.Sleep(senderConfig.Interval)
		}
	}
}

// NewCloudEvent creates a new cloud event
func NewCloudEvent(data interface{}, typ string) cloudevents.Event {
	e := cloudevents.NewEvent()
	e.SetDataContentType("application/json")
	e.SetDataContentEncoding(cloudevents.Base64)
	e.SetType(typ)
	host, err := os.Hostname()
	ensure.NoError(err)
	e.SetSource(fmt.Sprintf("knative://%s/wathola/sender", host))
	e.SetID(NewEventID())
	e.SetTime(time.Now())
	err = e.SetData(data)
	ensure.NoError(err)
	ensure.NoError(e.Validate())
	return e
}

// SendEvent will send cloud event to given url
func SendEvent(e cloudevents.Event, url string) error {
	ht, err := cloudevents.NewHTTPTransport(
		cloudevents.WithTarget(url),
		cloudevents.WithEncoding(cloudevents.HTTPBinaryV02),
	)
	ensure.NoError(err)
	c, err := cloudevents.NewClient(ht)
	ensure.NoError(err)
	ctx := context.Background()
	_, _, err = c.Send(ctx, e)
	return err
}

func (s *sender) sendStep() error {
	step := event.Step{Number: s.counter + 1}
	ce := NewCloudEvent(step, event.StepType)
	url := senderConfig.Address
	log.Infof("Sending step event #%v to %s", step.Number, url)
	err := SendEvent(ce, url)
	if err != nil {
		return err
	}
	s.counter++
	return nil
}

func (s *sender) sendFinished() {
	if s.counter == 0 {
		return
	}
	finished := event.Finished{Count: s.counter}
	url := senderConfig.Address
	ce := NewCloudEvent(finished, event.FinishedType)
	log.Infof("Sending finished event (count: %v) to %s", finished.Count, url)
	ensure.NoError(SendEvent(ce, url))
}
