GIT_COMMIT=$(shell git describe --always --long --dirty)

default: fmt build

all: fmt imports build

fmt:
	@echo "==> Fixing source code with gofmt..."
	find . -name '*.go' | grep -v vendor | xargs gofmt -s -w

imports:
	@echo "==> Fixing imports code with goimports..."
	goimports -w .

build:
	@echo "==> building..."
	cd cmd/terrafmt && go build -ldflags "-X github.com/katbyte/terrafmt/version.GitCommit=${GIT_COMMIT}" . && mv terrafmt ../../

install:
	@echo "==> installing..."
	cd cmd/terrafmt && go install -ldflags "-X github.com/katbyte/terrafmt/version.GitCommit=${GIT_COMMIT}" .

.PHONY: fmt imports build