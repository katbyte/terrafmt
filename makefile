GIT_COMMIT=$(shell git describe --always --long --dirty)
GIT_VERSION=$(shell git describe --tags --dirty 2>/dev/null | sed 's/-\([0-9]*\)-g/+\1@g/' || echo dev)
GOLANGCI_LINT_VERSION?=v2.12.2
TEST_TIMEOUT?=15m

default: fmt build

all: fmt build

help: ## Show this help
	@awk 'BEGIN {FS = ":.*##"; printf "Usage: make \033[36m<target>\033[0m\n"} /^[a-zA-Z0-9_-]+:.*?##/ { printf "  \033[36m%-24s\033[0m%s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) }' $(MAKEFILE_LIST)

##@ Build
build: ## Compile terrafmt with version info from git
	@echo "==> building..."
	go build -ldflags "-X github.com/katbyte/terrafmt/lib/version.GitCommit=${GIT_COMMIT} -X github.com/katbyte/terrafmt/lib/version.Version=${GIT_VERSION}"

install: ## Install terrafmt into GOPATH/bin with version info from git
	@echo "==> installing..."
	go install -ldflags "-X github.com/katbyte/terrafmt/lib/version.GitCommit=${GIT_COMMIT} -X github.com/katbyte/terrafmt/lib/version.Version=${GIT_VERSION}" .

tools: ## Install the tools required for development (golangci-lint)
	@echo "==> installing required tooling..."
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | \
		sh -s -- -b $(shell go env GOPATH)/bin ${GOLANGCI_LINT_VERSION}

##@ Formatting
fmt: ## Fix Go formatting (gofmt, gofumpt, goimports)
	@echo "==> Fixing source code with gofmt..."
	find . -name '*.go' | grep -v vendor | xargs gofmt -s -w
	@echo "==> Fixing source code with gofumpt..."
	find . -name '*.go' | grep -v vendor | xargs gofumpt -w
	@echo "==> Fixing imports with golangci-lint (goimports)..."
	golangci-lint fmt -E goimports ./...

goimports: ## Fix imports with golangci-lint (goimports)
	@echo "==> Fixing imports with golangci-lint (goimports)..."
	golangci-lint fmt -E goimports ./...

##@ Linting & Dependencies
lint: ## Check source code with the golangci linters
	@echo "==> Checking source code against linters..."
	golangci-lint run ./...

lint-fix: ## Fix source code with all golangci linters
	@echo "==> Checking source code against linters (applying autofixes)..."
	golangci-lint run --fix ./...

depscheck: ## Check that go.mod/go.sum and vendor/ are in sync
	@echo "==> Checking source code with go mod tidy..."
	@go mod tidy
	@git diff --exit-code -- go.mod go.sum || \
		(echo; echo "Unexpected difference in go.mod/go.sum files. Run 'go mod tidy' command or revert any go.mod/go.sum changes and commit."; exit 1)
	@echo "==> Checking source code with go mod vendor..."
	@go mod vendor
	@git diff --compact-summary --exit-code -- vendor || \
		(echo; echo "Unexpected difference in vendor/ directory. Run 'go mod vendor' command or revert any go.mod/go.sum/vendor changes and commit."; exit 1)

##@ Testing
test: build ## Run the unit tests (with -race)
	go test -race ./... -timeout ${TEST_TIMEOUT}

check-against-providers: ## Check formatting against real provider repos (golden vs main + idempotency)
	@echo "==> Checking against real provider repos (golden vs main + idempotency)..."
	./scripts/check-against-providers.sh

check-all: build test lint depscheck ## Run build + test + lint + depscheck

.PHONY: default all help fmt goimports build lint lint-fix depscheck check-against-providers check-all install tools
