package templates

import (
	"github.com/rdeusser/vulcan/internal/scaffold"
)

var _ scaffold.Template = &Makefile{}

type Makefile struct {
	scaffold.TemplateMixin
	scaffold.ModulePathMixin
	scaffold.ProjectNameMixin
	scaffold.ProtobufMixin
}

func (t *Makefile) GetIfExistsAction() scaffold.IfExistsAction {
	return t.IfExistsAction
}

func (t *Makefile) SetTemplateDefaults() error {
	if t.Path == "" {
		t.Path = "Makefile"
	}

	t.TemplateBody = makefileTemplate
	t.IfExistsAction = scaffold.Overwrite

	return nil
}

const makefileTemplate = `# Setup variables for the package.
NAME := {{ .ProjectName }}
PKG := {{ .ModulePath }}
BUILD_PATH := $(PKG)/cmd/$(NAME)
VERSION := $(shell grep -oE "[0-9]+[.][0-9]+[.][0-9]+" version/version.go)

SEMVER := patch

OLDPWD := $(PWD)
export OLDPWD

FILES_TO_FMT ?= $(shell find . -path ./vendor -prune -o -name '*.go' -print)

DOCKER_IMAGE_REPO ?= change-me

GOBIN		   ?= $(firstword $(subst :, ,${GOPATH}))/bin
GO111MODULE	   ?= on
export GO111MODULE

# Dependencies

GOIMPORTS_VERSION             ?= master
GOIMPORTS                     ?= $(GOBIN)/goimports

REVIVE_VERSION                ?= v1.2.1
REVIVE                        ?= $(GOBIN)/revive

{{- if .ProtobufSupport -}}
BUF_VERSION		      ?= v0.20.5
BUF			      ?= $(GOBIN)/buf
PROTOC_GEN_BUF_CHECK_BREAKING ?= $(GOBIN)/protoc-gen-buf-check-breaking
PROTOC_GEN_BUF_CHECK_LINT     ?= $(GOBIN)/protoc-gen-buf-check-lint
{{- end }}

.DEFAULT_GOAL := help

define install_go_bin_version
	@go install $(1)@$(2)
endef

.PHONY: help
help: ## Display this help text.
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n\nAvailable targets:\n"} /^[\/0-9a-zA-Z_-]+:.*?##/ { printf "  \x1b[32;01m%-20s\x1b[0m %s\n", $$1, $$2 }' $(MAKEFILE_LIST)

.PHONY: tidy
tidy: $(GOIMPORTS) ## Formats Go code including imports and cleans up noise.
	@echo "==> Formatting code"
	@$(GOIMPORTS) -local $(PKG) -w $(FILES_TO_FMT)
	@echo "==> Cleaning up noise"
	@find . -type f \( -name "*.md" -o -name "*.go" \) | SED_BIN="$(SED)" xargs scripts/cleanup-noise.sh
	@echo "==> Running 'go mod tidy'"
	@go mod tidy

.PHONY: generate
generate: ## Generate code.
	@echo "==> Generating code"
	@go generate ./...
{{- if .ProtobufSupport }}
	@buf generate
{{- end }}

.PHONY: lint
lint: $(REVIVE) ## Run lint tools.
	@echo "==> Running linting tools"
	@revive -config revive.toml ./...
{{- if .ProtobufSupport }}
	@buf lint
{{- end }}

.PHONY: test
test: ## Runs all {{ .ProjectName }}'s unit tests. This excludes tests in ./test/e2e.
	@echo "==> Running unit tests (without /test/e2e)"
	@go test -v -coverprofile=coverage.out $(shell go list ./... | grep -v /test/e2e);

.PHONY: test/e2e
test/e2e: ## Runs all {{ .ProjectName }}'s e2e tests from test/e2e.
	@echo "==> Running e2e tests"
	@go test -v -tags=e2e -coverprofile=coverage.out ./test/e2e/...

.PHONY: build
build: ## Build {{ .ProjectName }}.
	@echo "==> Building {{ .ProjectName }}"
	@-CGO_ENABLED=0 \
		go build \
		-o bin/{{ .ProjectName }} \
		$(BUILD_PATH)

.PHONY: install
install: build ## Install {{ .ProjectName }}.
	@echo "==> Installing {{ .ProjectName }}"
	@mv ./bin/{{ .ProjectName }} $(GOBIN)/{{ .ProjectName }}

.PHONY: docker-build
docker-build: ## Build docker image.
	@echo "==> Building docker image"
	@docker build -t $(DOCKER_IMAGE_REPO):$(VERSION) --build-arg GOPROXY=$(GOPROXY) .

.PHONY: bump-version
bump-version: ## Bump the version in the version file. Set SEMVER to [ patch (default) | major | minor ].
	@./scripts/bump-version.sh $(SEMVER)

.PHONY: tag
tag: ## Create and push a new git tag (creates tag using version/version.go file).
	@./scripts/tag.sh

$(GOIMPORTS):
	$(call install_go_bin_version,golang.org/x/tools/cmd/goimports,$(GOIMPORTS_VERSION))

$(REVIVE):
	$(call install_go_bin_version,github.com/mgechev/revive,$(REVIVE_VERSION))

{{- if .ProtobufSupport }}
$(BUF):
	$(call install_go_bin_version,github.com/bufbuild/buf/cmd/buf,$(BUF_VERSION))
	$(call install_go_bin_version,github.com/bufbuild/buf/cmd/protoc-gen-buf-check-breaking,$(BUF_VERSION))
	$(call install_go_bin_version,github.com/bufbuild/buf/cmd/protoc-gen-buf-check-lint,$(BUF_VERSION))
{{- end -}}`
