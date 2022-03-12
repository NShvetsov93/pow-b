DEV_PROJECT := pow-b
DEV_COMPOSE_FILE := scripts/docker-compose.yml

# -----------------------------------------------------------------------------
.PHONY: app
app:
	${INFO} "Removing old instance app (if any)"
	@docker-compose -p $(DEV_PROJECT) -f $(DEV_COMPOSE_FILE) stop go
	@docker-compose -p $(DEV_PROJECT) -f $(DEV_COMPOSE_FILE) rm -f go
	${INFO} "Removing old instance redis (if any)"
	@docker-compose -p $(DEV_PROJECT) -f $(DEV_COMPOSE_FILE) stop redis
	@docker-compose -p $(DEV_PROJECT) -f $(DEV_COMPOSE_FILE) rm -f redis
	${INFO} "Running redis"
	@docker-compose -p $(DEV_PROJECT) -f $(DEV_COMPOSE_FILE) up -d redis
	${INFO} "Running backend application"
	@docker-compose -p $(DEV_PROJECT) -f $(DEV_COMPOSE_FILE) up -d go
