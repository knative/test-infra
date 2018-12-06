TAG = $(shell date -u +v%Y%m%d)-$(shell git describe --tags --always --dirty)

# Build and push specific variables.
REGISTRY ?= gcr.io
PROJECT  ?= knative-tests
PUSH     ?= docker push

# Making dev version of coverage-image
coverage-image-dev:
	@echo "Making dev version..."
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o coverage
	chmod +x coverage
	docker build -t "$(REGISTRY)/$(PROJECT)/test-infra/coverage-dev:$(TAG)" -t "$(REGISTRY)/$(PROJECT)/test-infra/coverage-dev:latest-dev" .
	$(PUSH) "$(REGISTRY)/$(PROJECT)/test-infra/coverage-dev:$(TAG)"
	$(PUSH) "$(REGISTRY)/$(PROJECT)/test-infra/coverage-dev:latest-dev"

# Making production version of coverage-image
production:
	@echo "Making production version..."
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o coverage
	chmod +x coverage
	docker build -t "$(REGISTRY)/$(PROJECT)/test-infra/coverage:$(TAG)" -t "$(REGISTRY)/$(PROJECT)/test-infra/coverage:latest" .
	$(PUSH) "$(REGISTRY)/$(PROJECT)/test-infra/coverage:$(TAG)"
	$(PUSH) "$(REGISTRY)/$(PROJECT)/test-infra/coverage:latest"
