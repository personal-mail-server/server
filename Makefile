COMPOSE_FILE := code/docker-compose.yml

.PHONY: up down status logs ps

up:
	docker compose -f $(COMPOSE_FILE) up --build -d

down:
	docker compose -f $(COMPOSE_FILE) down

status:
	docker compose -f $(COMPOSE_FILE) ps

ps: status

logs:
	docker compose -f $(COMPOSE_FILE) logs -f
