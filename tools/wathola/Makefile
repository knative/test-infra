PROJECT_DIR   = $(shell readlink -f .)
BUILD_DIR     = "$(PROJECT_DIR)/build"
SENDER_DIR    = "$(PROJECT_DIR)/cmd/sender"
SENDER_BIN    = "$(BUILD_DIR)/wathola-sender"
FORWARDER_DIR = "$(PROJECT_DIR)/cmd/forwarder"
FORWARDER_BIN = "$(BUILD_DIR)/wathola-forwarder"
RECEIVER_DIR  = "$(PROJECT_DIR)/cmd/receiver"
RECEIVER_BIN  = "$(BUILD_DIR)/wathola-receiver"

GO           ?= go
RICHGO       ?= rich$(GO)

.PHONY: default
default: binaries

.PHONY: builddeps
builddeps:
	@GO111MODULE=off $(GO) get github.com/kyoh86/richgo
	@GO111MODULE=off $(GO) get github.com/mgechev/revive

.PHONY: builddir
builddir:
	@mkdir -p build

.PHONY: clean
clean: builddeps
	@echo "🛁 Cleaning"
	@rm -frv $(BUILD_DIR)

.PHONY: check
check: builddeps
	@echo "🛂 Checking"
	revive -config revive.toml -formatter stylish ./...

.PHONY: test
test: builddir check
	@echo "✔️ Testing"
	$(RICHGO) test -v -covermode=count -coverprofile=build/coverage.out ./...

.PHONY: sender
sender: builddir test
	@echo "🔨 Building sender"
	CGO_ENABLED=0 GOOS=linux $(RICHGO) build -o $(SENDER_BIN) $(SENDER_DIR)

.PHONY: forwarder
forwarder: builddir test
	@echo "🔨 Building forwarder"
	CGO_ENABLED=0 GOOS=linux $(RICHGO) build -o $(FORWARDER_BIN) $(FORWARDER_DIR)

.PHONY: receiver
receiver: builddir test
	@echo "🔨 Building receiver"
	CGO_ENABLED=0 GOOS=linux $(RICHGO) build -o $(RECEIVER_BIN) $(RECEIVER_DIR)

.PHONY: binaries
binaries: sender forwarder receiver
