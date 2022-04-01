# Setup name variables for the package/tool.
NAME := vulcan
PKG := github.com/rdeusser/$(NAME)
BUILD_PATH := $(PKG)/cmd/$(NAME)
GIT_COMMIT := $(PKG)/version
VERSION := $(shell grep -oE "[0-9]+[.][0-9]+[.][0-9]+" version/version.go)

SEMVER := patch

OLDPWD := $(PWD)
export OLDPWD

OUT_DIR := $(PWD)/bin

FILES_TO_FMT ?= $(shell find . -path ./vendor -prune -o -name '*.go' -print)

GOBIN		   ?= $(firstword $(subst :, ,${GOPATH}))/bin
GO111MODULE	   ?= on
export GO111MODULE

GOIMPORTS_VERSION	      ?= master
GOIMPORTS		      ?= $(GOBIN)/goimports

GOLANGCILINT_VERSION	      ?= v1.30.0
GOLANGCILINT		      ?= $(GOBIN)/golangci-lint

BUF_VERSION		      ?= v0.20.5
BUF			      ?= $(GOBIN)/buf
PROTOC_GEN_BUF_CHECK_BREAKING ?= $(GOBIN)/protoc-gen-buf-check-breaking
PROTOC_GEN_BUF_CHECK_LINT     ?= $(GOBIN)/protoc-gen-buf-check-lint

.DEFAULT_GOAL := help

define fetch_go_bin_version
	@cd /tmp
	@go get $(1)@$(2)
	@cd -
endef

.PHONY: help
help: ## Display this help text.
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n\nTargets:\n"} /^[0-9a-zA-Z_-]+:.*?##/ { printf "    \033[36m%-12s\033[0m %s\n", $$1, $$2 }' $(MAKEFILE_LIST)

.PHONY: tidy
tidy: $(GOIMPORTS) ## Formats Go code including imports and cleans up noise.
	@echo ">> formatting code"
	@$(GOIMPORTS) -local github.com/rdeusser/vulcan -w $(FILES_TO_FMT)
	@echo ">> cleaning up noise"
	@find . -type f \( -name "*.md" -o -name "*.go" \) | SED_BIN="$(SED)" xargs scripts/cleanup-noise.sh
	@echo ">> running 'go mod tidy'"
	@go mod tidy

.PHONY: lint
lint: $(GOLANGCILINT) ## Run various static analysis tools against our code.
	@echo ">> linting all of the Go files"
	@$(GOLANGCILINT) run

.PHONY: generate
generate: $(BUF) ## Generates Protobuf and/or Go code.
	@echo ">> generating code"
	@(BUF) check lint

.PHONY: test
test: ## Runs all vulcan's unit tests. This excludes tests in ./test/e2e.
	@echo ">> running unit tests (without /test/e2e)"
	@go test -coverprofile=coverage.out $(shell go list ./... | grep -v /test/e2e);

.PHONY: test/e2e
test/e2e: generate tidy # Runs all vulcan's e2e tests from test/e2e.
	@echo ">> running e2e tests"
	@go test -v -tags=e2e ./test/e2e/... -coverprofile cover.out

.PHONY: build
build: ## Build vulcan.
	@echo ">> building vulcan"
	@-CGO_ENABLED=0 \
		go build \
		-o $(OUT_DIR)/vulcan \
		$(BUILD_PATH)

.PHONY: install
install: build ## Build and install vulcan.
	@echo ">> installing vulcan"
	 mv $(OUT_DIR)/vulcan $(GOBIN)

.PHONY: bump-version
bump-version: ## Bump the version in the version file. Set SEMVER to [ patch (default) | major | minor ].
	@./scripts/bump-version.sh $(SEMVER)

.PHONY: tag
tag: ## Create and push a new git tag (creates tag using version/version.go file).
	@./scripts/tag.sh

$(GOIMPORTS):
	$(call fetch_go_bin_version,golang.org/x/tools/cmd/goimports,$(GOIMPORTS_VERSION))

$(GOLANGCILINT):
	$(call fetch_go_bin_version,github.com/golangci/golangci-lint/cmd/golangci-lint,$(GOLANGCILINT_VERSION))

$(BUF):
	$(call fetch_go_bin_version,github.com/bufbuild/buf/cmd/buf,$(BUF_VERSION))
	$(call fetch_go_bin_version,github.com/bufbuild/buf/cmd/protoc-gen-buf-check-breaking,$(BUF_VERSION))
	$(call fetch_go_bin_version,github.com/bufbuild/buf/cmd/protoc-gen-buf-check-lint,$(BUF_VERSION))
