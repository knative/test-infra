package main

import "knative.dev/test-infra/tools/wathola/internal/forwarder"

var instance forwarder.Forwarder

func main() {
	instance = forwarder.New()
	instance.Forward()
}
