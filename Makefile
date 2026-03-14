COMPOSE_FILE := code/docker-compose.yml
BASE_BRANCH ?= main

.PHONY: up down status logs ps push-trigger

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
	go run ./code/cmd/push-trigger --base $(BASE_BRANCH)
