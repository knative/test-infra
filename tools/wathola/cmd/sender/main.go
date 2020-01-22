package main

import "knative.dev/test-infra/tools/wathola/internal/sender"

func main() {
	sender.New().SendContinually()
}
