package main

import "knative.dev/test-infra/tools/wathola/internal/receiver"

var instance receiver.Receiver

func main() {
	instance = receiver.New()
	instance.Receive()
}
