GIT_COMMIT=$(shell git describe --always --long --dirty)
GOLANGCI_LINT_VERSION=v1.29.0

default: fmt build

all: fmt imports build

tools:
	@echo "==> installing required tooling..."
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | \
		sh -s -- -b $(shell go env GOPATH)/bin ${GOLANGCI_LINT_VERSION}

fmt:
	@echo "==> Fixing source code with gofmt..."
	find . -name '*.go' | grep -v vendor | xargs gofmt -s -w

imports:
	@echo "==> Fixing imports code with goimports..."
	goimports -w .

test: build
	go test ./...

build:
	@echo "==> building..."
	go build -ldflags "-X github.com/katbyte/terrafmt/lib/version.GitCommit=${GIT_COMMIT}"

lint:
	@echo "==> Checking source code against linters..."
	golangci-lint run ./...

install:
	@echo "==> installing..."
	go install -ldflags "-X github.com/katbyte/terrafmt/lib/version.GitCommit=${GIT_COMMIT}" .

.PHONY: fmt imports build lint install tools