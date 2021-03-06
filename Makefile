.DEFAULT_GOAL := help

TEST_FLAGS ?= -race
PKG_BASE   ?= $(shell go list .)
PKGS       ?= $(shell go list ./... | grep -v /vendor/)
SOURCES     = $(shell find . -name '*.go')
BINARIES    = $(shell find cmd/ -mindepth 1 -maxdepth 1 -type d | sed -e 's/cmd\//build\//g')

.PHONY: help
help:
	@grep -E '^[a-zA-Z0-9-]+:.*?## .*$$' Makefile | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "[32m%-10s[0m %s\n", $$1, $$2}'

.PHONY: test
test: ## run tests
	go test $(TEST_FLAGS) $(PKGS)

.PHONY: vet
vet: ## run go vet
	go vet $(PKGS)

.PHONY: coverage
coverage: ## generate code coverage
	go test $(TEST_FLAGS) -covermode=atomic -coverprofile=coverage.txt $(PKGS)
	go tool cover -func=coverage.txt

.PHONY: lint
lint: ## run golangci-lint
	golangci-lint run

.PHONY: clean
clean: ## cleanup build dir
	rm -rf build/
	mkdir build/

.PHONY: build
build: $(BINARIES) ## build all binaries

build/%: $(SOURCES)
	go build -ldflags "-s -w" -o build/$* ./cmd/$*
