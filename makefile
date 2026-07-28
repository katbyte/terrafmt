GIT_COMMIT=$(shell git describe --always --long --dirty)
GIT_VERSION=$(shell git describe --tags --dirty 2>/dev/null | sed 's/-\([0-9]*\)-g/+\1@g/' || echo dev)
GOLANGCI_LINT_VERSION?=v2.12.2
TEST_TIMEOUT?=15m

default: fmt build

all: fmt build

tools:
	@echo "==> installing required tooling..."
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | \
		sh -s -- -b $(shell go env GOPATH)/bin ${GOLANGCI_LINT_VERSION}

fmt:
	@echo "==> Fixing source code with gofmt..."
	find . -name '*.go' | grep -v vendor | xargs gofmt -s -w
	@echo "==> Fixing source code with gofumpt..."
	find . -name '*.go' | grep -v vendor | xargs gofumpt -w
	@echo "==> Fixing imports with golangci-lint (goimports)..."
	golangci-lint fmt -E goimports ./...

goimports:
	@echo "==> Fixing imports with golangci-lint (goimports)..."
	golangci-lint fmt -E goimports ./...

test: build
	go test -race ./... -timeout ${TEST_TIMEOUT}

build:
	@echo "==> building..."
	go build -ldflags "-X github.com/katbyte/terrafmt/lib/version.GitCommit=${GIT_COMMIT} -X github.com/katbyte/terrafmt/lib/version.Version=${GIT_VERSION}"

lint:
	@echo "==> Checking source code against linters..."
	golangci-lint run ./...

lint-fix:
	@echo "==> Checking source code against linters (applying autofixes)..."
	golangci-lint run --fix ./...

depscheck:
	@echo "==> Checking source code with go mod tidy..."
	@go mod tidy
	@git diff --exit-code -- go.mod go.sum || \
		(echo; echo "Unexpected difference in go.mod/go.sum files. Run 'go mod tidy' command or revert any go.mod/go.sum changes and commit."; exit 1)
	@echo "==> Checking source code with go mod vendor..."
	@go mod vendor
	@git diff --compact-summary --exit-code -- vendor || \
		(echo; echo "Unexpected difference in vendor/ directory. Run 'go mod vendor' command or revert any go.mod/go.sum/vendor changes and commit."; exit 1)

install:
	@echo "==> installing..."
	go install -ldflags "-X github.com/katbyte/terrafmt/lib/version.GitCommit=${GIT_COMMIT} -X github.com/katbyte/terrafmt/lib/version.Version=${GIT_VERSION}" .

check-against-providers:
	@echo "==> Checking against real provider repos (golden vs main + idempotency)..."
	./scripts/check-against-providers.sh

check-all: build test lint depscheck

.PHONY: fmt goimports build lint lint-fix depscheck check-against-providers check-all install tools
