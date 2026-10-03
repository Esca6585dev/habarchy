# Habarchy monorepo. Per-app targets live in each app directory.
.PHONY: dev-deps dev-deps-down backend-test backend-lint backend-build test lint

dev-deps:            ## start Postgres + Redis for local development
	docker compose -f deploy/docker-compose.dev.yml up -d

dev-deps-down:
	docker compose -f deploy/docker-compose.dev.yml down

backend-test:
	$(MAKE) -C backend test

backend-lint:
	$(MAKE) -C backend lint

backend-build:
	$(MAKE) -C backend build

test: backend-test
lint: backend-lint
