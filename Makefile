COMPOSE_FILE := code/docker-compose.yml
BASE_BRANCH ?= main
COMPOSE_NETWORK := code_default
MIGRATE_DATABASE_URL := postgres://postgres:postgres@db:5432/mail_server?sslmode=disable
DB_CONTAINER_NAME := personal-mail-db

.PHONY: up down status logs ps push-trigger migrate-up migrate-down

up:
	docker compose -f $(COMPOSE_FILE) up --build -d

down:
	docker compose -f $(COMPOSE_FILE) down

status:
	docker compose -f $(COMPOSE_FILE) ps

ps: status

logs:
	docker compose -f $(COMPOSE_FILE) logs -f

push-trigger:
	go -C code run ./cmd/push-trigger --base $(BASE_BRANCH)

migrate-up:
	docker compose -f $(COMPOSE_FILE) up -d db
	until [ "$$(docker inspect -f '{{.State.Health.Status}}' $(DB_CONTAINER_NAME) 2>/dev/null)" = "healthy" ]; do sleep 1; done
	docker compose -f $(COMPOSE_FILE) build backend
	docker run --rm --entrypoint /app/migrate --network $(COMPOSE_NETWORK) -e DATABASE_URL="$(MIGRATE_DATABASE_URL)" code-backend:latest --direction up

migrate-down:
	docker compose -f $(COMPOSE_FILE) up -d db
	until [ "$$(docker inspect -f '{{.State.Health.Status}}' $(DB_CONTAINER_NAME) 2>/dev/null)" = "healthy" ]; do sleep 1; done
	docker compose -f $(COMPOSE_FILE) build backend
	docker run --rm --entrypoint /app/migrate --network $(COMPOSE_NETWORK) -e DATABASE_URL="$(MIGRATE_DATABASE_URL)" code-backend:latest --direction down --steps $${STEPS:-1}
