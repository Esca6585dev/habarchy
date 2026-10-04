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

# ---- production stack (deploy/docker-compose.yml) ----------------------------
COMPOSE := docker compose -f deploy/docker-compose.yml --env-file deploy/.env
.PHONY: up down logs seed create-admin ps sdk-test

deploy/.env:
	cp deploy/.env.example deploy/.env && echo ">> edit deploy/.env (secrets, HABARCHY_PUBLIC_URL)"

up: deploy/.env          ## build + start nginx, web, api, worker, scheduler, postgres, redis, asynqmon
	$(COMPOSE) up -d --build

down:
	$(COMPOSE) down

logs:
	$(COMPOSE) logs -f --tail=200 api worker scheduler web

ps:
	$(COMPOSE) ps

seed:                    ## demo project, sample templates (tk/ru/en), test API key
	$(COMPOSE) run --rm api api -seed

create-admin:            ## admin user from HABARCHY_ADMIN_EMAIL / HABARCHY_ADMIN_PASSWORD in deploy/.env
	$(COMPOSE) run --rm api api -create-admin

sdk-test:                ## run every SDK test suite
	cd sdk/go && go test ./...
	cd sdk/ts && npm install --silent && npm test
	cd sdk/dart && dart pub get && dart test
	cd sdk/php && php tests/run.php
