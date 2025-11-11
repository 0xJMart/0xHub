GO ?= /home/admin/go/bin/go
GOLANGCI_LINT ?= golangci-lint
DOCKER ?= docker
GOOSE ?= $(GO) run github.com/pressly/goose/v3/cmd/goose
MIGRATIONS_DIR ?= $(CURDIR)/migrations
DB_DRIVER ?= postgres
DB_DSN ?= postgres://hub:hub@localhost:5432/hub?sslmode=disable

.PHONY: all lint test generate docker-build migrate-up migrate-down migrate-status migrate-reset

all: lint test

lint:
	$(GOLANGCI_LINT) run ./...

test:
	$(GO) test ./...

generate:
	$(GO) generate ./...

docker-build:
	$(DOCKER) build -t hub-api:dev .

migrate-up:
	$(GOOSE) -dir $(MIGRATIONS_DIR) $(DB_DRIVER) $(DB_DSN) up

migrate-down:
	$(GOOSE) -dir $(MIGRATIONS_DIR) $(DB_DRIVER) $(DB_DSN) down

migrate-status:
	$(GOOSE) -dir $(MIGRATIONS_DIR) $(DB_DRIVER) $(DB_DSN) status

migrate-reset:
	$(GOOSE) -dir $(MIGRATIONS_DIR) $(DB_DRIVER) $(DB_DSN) reset


