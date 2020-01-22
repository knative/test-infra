package client

import (
	"context"
	cloudevents "github.com/cloudevents/sdk-go"
	cloudeventshttp "github.com/cloudevents/sdk-go/pkg/cloudevents/transport/http"
	log "github.com/sirupsen/logrus"
	"knative.dev/test-infra/tools/wathola/internal/config"
	"knative.dev/test-infra/tools/wathola/internal/ensure"
	nethttp "net/http"
	"strings"
)

// ReceiveEvent represents a function that receive event
type ReceiveEvent func(e cloudevents.Event)

// Receive events and push then to passed fn
func Receive(
	port int,
	cancel *context.CancelFunc,
	receiveEvent ReceiveEvent,
	middlewares ... cloudeventshttp.Middleware) {
	portOpt := cloudevents.WithPort(port)
	opts := make([]cloudeventshttp.Option, 0)
	opts = append(opts, portOpt)
	if config.Instance.Readiness.Enabled {
		readyOpt := cloudevents.WithMiddleware(readinessMiddleware)
		opts = append(opts, readyOpt)
	}
	for _, m := range middlewares {
		opt := cloudevents.WithMiddleware(m)
		opts = append(opts, opt)
	}
	http, err := cloudevents.NewHTTPTransport(opts...)
	if err != nil {
		log.Fatalf("failed to create http transport, %v", err)
	}
	c, err := cloudevents.NewClient(http)
	if err != nil {
		log.Fatalf("failed to create client, %v", err)
	}
	ctx, ccancel := context.WithCancel(context.Background())
	cancel = &ccancel
	log.Infof("listening for events on port %v", port)
	err = c.StartReceiver(ctx, receiveEvent)
	if err != nil {
		log.Fatal(err)
	}
}

func readinessMiddleware(next nethttp.Handler) nethttp.Handler {
	log.Debugf("Using readiness probe: %v", config.Instance.Readiness.URI)
	return &readinessProbe{
		next: next,
	}
}

type readinessProbe struct {
	next nethttp.Handler
}

func (r readinessProbe) ServeHTTP(rw nethttp.ResponseWriter, req *nethttp.Request) {
	if req.RequestURI == config.Instance.Readiness.URI {
		rw.WriteHeader(config.Instance.Readiness.Status)
		_, err := rw.Write([]byte(config.Instance.Readiness.Message))
		ensure.NoError(err)
		log.Debugf("Received ready check. Headers: %v", headersOf(req))
	} else {
		r.next.ServeHTTP(rw, req)
	}
}

func headersOf(req *nethttp.Request) string {
	var b strings.Builder
	ensure.NoError(req.Header.Write(&b))
	headers := b.String()
	return strings.ReplaceAll(headers, "\r\n", "; ")
}
