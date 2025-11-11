COMPOSE ?= docker compose
ENV_FILE ?= config/env.example

.PHONY: dev-up dev-down dev-logs ci-up ci-down format lint typecheck test

dev-up:
	$(COMPOSE) --env-file $(ENV_FILE) -f docker-compose.dev.yml up --build -d

dev-down:
	$(COMPOSE) --env-file $(ENV_FILE) -f docker-compose.dev.yml down --remove-orphans

dev-logs:
	$(COMPOSE) --env-file $(ENV_FILE) -f docker-compose.dev.yml logs -f

ci-up:
	$(COMPOSE) -f docker-compose.ci.yml up --build -d

ci-down:
	$(COMPOSE) -f docker-compose.ci.yml down --remove-orphans

format:
	npm run format

lint:
	npm run lint

typecheck:
	npm run typecheck

test:
	npm run test

