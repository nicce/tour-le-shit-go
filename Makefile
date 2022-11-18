SHELL := /bin/bash

OS := $(shell uname)
ARCHITECTURE := $(shell uname -m | sed "s/x86_64/amd64/g")
BRANCH_NAME := $(shell git rev-parse --abbrev-ref HEAD)
CREATED := $(shell date +%Y-%m-%dT%T%z)
GIT_REPO := $(shell git config --get remote.origin.url)
GIT_TOKEN ?= $(shell cat git-token.txt)
REPO_NAME := $(shell basename ${GIT_REPO} .git)
REVISION_ID := $(shell git rev-parse HEAD)
SHORT_SHA := $(shell git rev-parse --short HEAD)
TAG_NAME := $(shell git describe --exact-match --tags 2> /dev/null)
VERSION := $(if ${TAG_NAME},${TAG_NAME},"unversioned")
VERSION_PATH := github.com/ingka-group-digital/${REPO_NAME}/internal/version

LICENSED_VERSION := 3.3.1
LICENSED_PATH := bin/licensed_$(LICENSED_VERSION)/licensed
LICENSED_URL := https://github.com/github/licensed/releases/download/$(LICENSED_VERSION)/licensed-$(LICENSED_VERSION)-$(OS)-x64.tar.gz

GOLANGCI_LINT_VERSION := v1.50.1
GOLANGCI_LINT := bin/golangci-lint_$(GOLANGCI_LINT_VERSION)/golangci-lint
GOLINTCI_URL := https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh

GOTESTSUM_VERSION := 1.8.0
GOTESTSUM_PATH := bin/gotestsum_v$(GOTESTSUM_VERSION)/gotestsum
GOTESTSUM_URL := https://github.com/gotestyourself/gotestsum/releases/download/v$(GOTESTSUM_VERSION)/gotestsum_$(GOTESTSUM_VERSION)_$(OS)_$(ARCHITECTURE).tar.gz

all: help

## clean: Clean up all build artifacts
clean:
	@echo "🚀 Cleaning up old artifacts"
	@rm -f ${REPO_NAME}

## build: Build the application artifacts. Linting can be skipped by setting env variable IGNORE_LINTING.
build: clean $(if $(IGNORE_LINTING), , lint)
	@echo "🚀 Building artifacts"
	@go build -ldflags="-s -w -X '${VERSION_PATH}.Version=${VERSION}' -X '${VERSION_PATH}.Commit=${SHORT_SHA}'" -o bin/${REPO_NAME} .

## run: Run the application
run: build
	@echo "🚀 Running binary"
	@./bin/${REPO_NAME}

## licensed: Checks license of dependencies
licensed: licensed-cache licensed-status

## licensed-cache: Builds license cache for dependencies
licensed-cache: ${LICENSED_PATH}
	@echo "🚀 Building license cache"
	@${LICENSED_PATH} cache

## licensed-info: Returns information about the current licensed version being used
licensed-info:
	@echo ${LICENSED_PATH}

## licensed-status: Builds license cache for dependencies
licensed-status: ${LICENSED_PATH}
	@echo "🚀 Checking license status"
	@${LICENSED_PATH} status

## lint: Lint the source code
lint: ${GOLANGCI_LINT}
	@echo "🚀 Linting code"
	@$(GOLANGCI_LINT) run

## lint-info: Returns information about the current linter being used
lint-info:
	@echo ${GOLANGCI_LINT}

## test: Run Go tests
test: ${GOTESTSUM_PATH}
	@echo "🚀 Running tests"
	@set -o pipefail; ${GOTESTSUM_PATH} --format testname --no-color=false | grep -v 'EMPTY'; exit $$?

## test-benchmark: Run Go benchmark tests
test-benchmark:
	@echo "🚀 Running benchmark tests"
	@go test -bench=. -benchmem ./...

## install-hooks: Install Git hooks
install-hooks:
	@echo "🚀 Installing Git hooks"
	@cp hooks/pre-push .git/hooks/pre-push

## uninstall-hooks: Uninstall Git hooks
uninstall-hooks:
	@echo "🚀 Uninstalling Git hooks"
	@rm -f .git/hooks/pre-push

## coverage: Create a test coverage report in HTML format
coverage:
	@echo "🚀 Creating coverage report in HTML format"
	@go test -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out

${GOLANGCI_LINT}:
	@echo "📦 Installing golangci-lint ${GOLANGCI_LINT_VERSION}"
	@mkdir -p $(dir ${GOLANGCI_LINT})
	@curl -sfL ${GOLINTCI_URL} | sh -s -- -b ./$(patsubst %/,%,$(dir ${GOLANGCI_LINT})) ${GOLANGCI_LINT_VERSION} > /dev/null 2>&1

${GOTESTSUM_PATH}:
	@echo "📦 Installing GoTestSum ${GOTESTSUM_VERSION}"
	@mkdir -p $(dir ${GOTESTSUM_PATH})
	@curl -sSL ${GOTESTSUM_URL} > bin/gotestsum.tar.gz
	@tar -xzf bin/gotestsum.tar.gz -C $(patsubst %/,%,$(dir ${GOTESTSUM_PATH}))
	@rm -f bin/gotestsum.tar.gz

${LICENSED_PATH}:
	@echo "📦 Installing Licensed ${LICENSED_VERSION}"
	@mkdir -p $(dir ${LICENSED_PATH})
	@curl -sSL ${LICENSED_URL} > bin/licensed.tar.gz
	@tar -xzf bin/licensed.tar.gz -C $(patsubst %/,%,$(dir ${LICENSED_PATH}))
	@rm -f bin/licensed.tar.gz

help: Makefile
	@echo
	@echo "📗 Choose a command run in "${REPO_NAME}":"
	@echo
	@sed -n 's/^##//p' $< | column -t -s ':' |  sed -e 's/^/ /'
	@echo
