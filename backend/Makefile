# Habarchy backend
GO       ?= go
GOBIN    ?= $(shell $(GO) env GOPATH)/bin
SQLC     ?= $(GOBIN)/sqlc
GOOSE    ?= $(GOBIN)/goose
LINT     ?= golangci-lint
BIN      := bin
DB_URL   ?= $(HABARCHY_DATABASE_URL)
PKGS     := ./...

.PHONY: all build run run-worker run-scheduler test test-unit test-integration lint fmt tidy sqlc openapi migrate migrate-down migrate-status tools genkey clean

all: lint test build

build:
	$(GO) build -trimpath -ldflags="-s -w" -o $(BIN)/api ./cmd/api
	$(GO) build -trimpath -ldflags="-s -w" -o $(BIN)/worker ./cmd/worker
	$(GO) build -trimpath -ldflags="-s -w" -o $(BIN)/scheduler ./cmd/scheduler

run:            ## run the API with .env
	set -a && . ./.env && set +a && $(GO) run ./cmd/api

run-worker:
	set -a && . ./.env && set +a && $(GO) run ./cmd/worker

run-scheduler:
	set -a && . ./.env && set +a && $(GO) run ./cmd/scheduler

test:           ## unit + integration (integration skipped without HABARCHY_TEST_DATABASE_URL / HABARCHY_TEST_REDIS_URL)
	$(GO) test -race -count=1 $(PKGS)

test-unit:
	HABARCHY_TEST_DATABASE_URL= HABARCHY_TEST_REDIS_URL= $(GO) test -race -count=1 $(PKGS)

test-integration:
	@test -n "$(HABARCHY_TEST_DATABASE_URL)" || (echo "set HABARCHY_TEST_DATABASE_URL" && exit 1)
	$(GO) test -race -count=1 -run 'Flow|Migrations|WithTx|Idempotency|TokenBucket|CountersAndOTP' $(PKGS)

lint:
	$(LINT) run ./...

fmt:
	gofmt -s -w .
	$(GO) vet $(PKGS)

tidy:
	$(GO) mod tidy

sqlc:           ## regenerate internal/adapters/postgres/sqlcgen
	$(SQLC) generate
	$(SQLC) vet

migrate:        ## apply migrations (uses HABARCHY_DATABASE_URL)
	$(GOOSE) -dir migrations postgres "$(DB_URL)" up

migrate-down:
	$(GOOSE) -dir migrations postgres "$(DB_URL)" down

migrate-status:
	$(GOOSE) -dir migrations postgres "$(DB_URL)" status

tools:          ## install sqlc, goose, golangci-lint
	$(GO) install github.com/sqlc-dev/sqlc/cmd/sqlc@v1.29.0
	$(GO) install github.com/pressly/goose/v3/cmd/goose@v3.24.3

openapi:        ## validate the spec and copy it to docs/openapi.yaml
	$(GO) run ./cmd/api -env >/dev/null 2>&1 || true
	cp api/openapi.yaml ../docs/openapi.yaml

genkey:         ## print a fresh HABARCHY_MASTER_KEY
	$(GO) run ./cmd/api -genkey

clean:
	rm -rf $(BIN)
