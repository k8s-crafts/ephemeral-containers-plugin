SHELL := /usr/bin/env bash -o pipefail
.SHELLFLAGS = -ec


##@ General

.PHONY: help
help: ## Display this help.
	@awk 'BEGIN {FS = ":.*##"; printf "Makefile for ephemeral-containers-plugin project\n\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_0-9-]+:.*?##/ { printf "  \033[36m%-30s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

.PHONY: fmt
fmt: ## Run go fmt against source files.
	go fmt ./...

.PHONY: vet
vet: ## Run go vet agaist source files.
	go vet ./...

.PHONY: add-license
add-license: go-license ## Add license header to source files.
	$(GO_LICENSE) --config license.yaml $(shell find ./ -name "*.go")

.PHONY: add-license
check-license: go-license ## Check license header to source files.
	$(GO_LICENSE) --verify --config license.yaml $(shell find ./ -name "*.go")
	
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
	test -s $(GO_LICENSE) || GOBIN=$(LOCAL_BIN) go install github.com/palantir/go-license@v1.39.0

##@ Build

.PHONY: build
build: ## Build ephemeral-containers-plugin binary (i.e. must have kubectl- prefix).
	mkdir -p build && go build -o build/kubectl-ephemeral-containers-plugin main.go
