SHELL := bash
.ONESHELL:
.SHELLFLAGS := -eu -o pipefail -c
MAKEFLAGS += --warn-undefined-variables
MAKEFLAGS += --no-builtin-rules

.PHONY: help
help_spacing := 20
help: ## List all available targets with help
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | \
	awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-$(help_spacing)s\033[0m %s\n", $$1, $$2}'

.PHONY: enable-githooks
enable-githooks: ## Enable git hooks
	git config core.hooksPath .githooks

.PHONY: disable-githooks
disable-githooks: ## Disable git hooks
	git config --unset core.hooksPath

.PHONY: init
init: enable-githooks tidy codegen ## Prepare project for development

.PHONY: tidy
tidy: ## Tidying all project go modules
	go mod tidy

.PHONY: codegen
codegen: ## Run code generation
	go generate ./...

.PHONY: outdated
outdated: ## Print outdated dependencies (`go install github.com/psampaz/go-mod-outdated@latest` required)
	@go list -u -m -json all | go-mod-outdated -update -direct

.PHONY: .update
.update:
	go get -u ./...

.PHONY: update
update: .update tidy codegen ## Update go.mod dependencies

.PHONY: lint
lint: ## Run golangci-linter
	golangci-lint run

.PHONY: short-test
short-test: lint ## Run tests
	go test -short -v ./...

.PHONY: test
test: lint ## Run tests
	go test -v ./...