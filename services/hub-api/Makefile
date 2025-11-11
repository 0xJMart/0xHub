GO ?= /home/admin/go/bin/go
GOLANGCI_LINT ?= golangci-lint
DOCKER ?= docker

.PHONY: all lint test generate docker-build

all: lint test

lint:
	$(GOLANGCI_LINT) run ./...

test:
	$(GO) test ./...

generate:
	$(GO) generate ./...

docker-build:
	$(DOCKER) build -t hub-api:dev .


