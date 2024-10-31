SHELL := /usr/bin/env bash -o pipefail
.SHELLFLAGS = -ec

## Host info
GOOS ?= $(shell go env GOOS)
GOARCH ?= $(shell go env GOARCH)


## Plugin version. Bump for each release
VERSION ?= 1.2.0-dev
# Git Commit ID
# GIT_COMMIT_NO := $(shell git rev-parse HEAD 2> /dev/null || true)
# GIT_COMMIT_ID := $(if $(shell git status --porcelain --untracked-files=no),$(GIT_COMMIT_NO)-dirty,$(GIT_COMMIT_NO))

export PLUGIN_VERSION=v$(VERSION)
# export PLUGIN_GIT_COMMIT_ID=$(GIT_COMMIT_ID)

## Tool versions
GO_LICENSE_VERSION ?= 1.39.0
GOLANGCI_LINT_VERSION ?= 1.61.0
GINKGO_VERSION ?= 2.21.0 # Using ginkgo v2

## Other flags
SKIP_TESTS ?= false

##@ General

.PHONY: help
help: ## Display this help.
	@awk 'BEGIN {FS = ":.*##"; printf "Makefile for ephemeral-containers-plugin project\n\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_0-9-]+:.*?##/ { printf "  \033[36m%-30s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

.PHONY: fmt
fmt: add-license ## Run go fmt against source files.
	go fmt ./...

.PHONY: vet
vet: ## Run go vet against source files.
	go vet ./...

.PHONY: lint
lint: golangci-lint ## Run lint checks with golangci-lint
	$(GOLANGCI_LINT) run ./...

.PHONY: lint-fix
lint-fix: golangci-lint ## Apply lint fixes with golangci-lint
	$(GOLANGCI_LINT) run --fix ./...

.PHONY: shellcheck
shellcheck: ## Run shell check against scripts
	shellcheck ./scripts/*.sh

.PHONY: add-license
add-license: go-license ## Add license header to source files.
	$(GO_LICENSE) --config go-license.yaml $(shell find ./ -name "*.go")

.PHONY: add-license
check-license: go-license ## Check license header to source files.
	$(GO_LICENSE) --verify --config go-license.yaml $(shell find ./ -name "*.go")

##@ Tests

.PHONY: test-unit
test-unit: vet fmt ginkgo ## Run unit tests.
ifneq ($(SKIP_TESTS), true)
	$(GINKGO) -v -output-dir=. -cover -coverpkg=./... -r -coverprofile cover.out  ./pkg ./cmd
endif

## E2E Tests
TEST_BIN := $(shell pwd)/testbin

.PHONY: test-e2e
test-e2e: build e2e-setup ginkgo ## Run e2e tests.
ifneq ($(SKIP_TESTS), true)
	PATH="$(BUILD_DIR):$(TEST_BIN):$${PATH}" KUBECONFIG="$(TEST_BIN)/kubeconfig" $(GINKGO) --timeout=10m -v -output-dir=. ./e2e
endif

.PHONY: e2e-setup
e2e-setup: ## Setting up environment for e2e tests.
ifneq ($(SKIP_TESTS), true)
	./scripts/e2e_setup.sh
endif

.PHONY: e2e-teardown
e2e-teardown: ## Tearing environment for e2e tests.
ifneq ($(SKIP_TESTS), true)
	./scripts/e2e_cleanup.sh
endif

##@ Install tools

LOCAL_BIN ?= $(shell pwd)/bin
.PHONY: local-bin
local-bin: $(LOCAL_BIN) ## Location to install tools.
$(LOCAL_BIN):
	mkdir -p $(LOCAL_BIN)

GO_LICENSE ?= $(LOCAL_BIN)/go-license
PHONY: go-license
go-license: $(GO_LICENSE) ## Install go-license.
$(GO_LICENSE): local-bin
	test -s $(GO_LICENSE) || GOBIN=$(LOCAL_BIN) go install github.com/palantir/go-license@v$(GO_LICENSE_VERSION)

GOLANGCI_LINT ?= $(LOCAL_BIN)/golangci-lint
.PHONY: golangci-lint
golangci-lint: $(GOLANGCI_LINT) # Install golangci-lint.
$(GOLANGCI_LINT): local-bin
	test -s $(GOLANGCI_LINT) || GOBIN=$(LOCAL_BIN) go install github.com/golangci/golangci-lint/cmd/golangci-lint@v$(GOLANGCI_LINT_VERSION)


GINKGO ?= $(LOCAL_BIN)/ginkgo
.PHONY: ginkgo
ginkgo: $(GINKGO) ## Install ginkgo.
$(GINKGO): $(LOCAL_BIN)
	test -s $(GINKGO) || GOBIN=$(LOCAL_BIN) go install github.com/onsi/ginkgo/v2/ginkgo@v$(GINKGO_VERSION)

##@ Build

.PHONY: generate
generate: ## Generate go codes
	go generate ./...

BUILD_DIR ?= $(shell pwd)/build

.PHONY: build
build: generate ## Build ephemeral-containers-plugin binary (i.e. must have kubectl- prefix).
	mkdir -p $(BUILD_DIR)
	CGO_ENABLED=0 GOOS=$(GOOS) GOARCH=$(GOARCH) go build \
		-o $(BUILD_DIR)/kubectl-ephemeral_containers \
		main.go

.PHONY: clean
clean: ## Clean built binaries.
	- rm -rf $(BUILD_DIR)

# Default to $HOME/bin
INSTALL_DIR ?= $(HOME)/bin

.PHONY: install
install: build ## Build and install the plugin to PATH (Default to $HOME/bin).
ifeq (,$(findstring $(INSTALL_DIR),$(PATH)))
	$(error $(INSTALL_DIR) is not on $$PATH)
else
	@if [ -s "$(INSTALL_DIR)/kubectl-ephemeral_containers" ]; then \
		echo "$(INSTALL_DIR)/kubectl-ephemeral_containers exists. Please remove it with \"make uninstall\"."; \
		exit 1; \
	fi; \
	install -m 755 $(BUILD_DIR)/kubectl-ephemeral_containers $(INSTALL_DIR)/kubectl-ephemeral_containers
	@echo Installed to $(BUILD_DIR)/kubectl-ephemeral_containers
endif


.PHONY: uninstall
uninstall: ## Unistall the plugin from PATH.
ifeq (,$(findstring $(INSTALL_DIR),$(PATH)))
	$(error $(INSTALL_DIR) is not on $$PATH)
else
	- test -s $(INSTALL_DIR)/kubectl-ephemeral_containers && rm $(INSTALL_DIR)/kubectl-ephemeral_containers
	@echo Removed $(BUILD_DIR)/kubectl-ephemeral_containers
endif
