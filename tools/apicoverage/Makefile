REGISTRY ?= gcr.io
PROJECT  ?= knative-tests/test-infra
PUSH     ?= docker push
apicoverage-image:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o apicoverage
	docker build -t "$(REGISTRY)/$(PROJECT)/apicoverage:latest" .
	$(PUSH) "$(REGISTRY)/$(PROJECT)/apicoverage:latest"
