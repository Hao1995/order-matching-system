# Variables
COMPOSE_FILE = docker-compose.yml
DOCKER_COMPOSE = docker-compose -f $(COMPOSE_FILE)

# Targets
.PHONY: build up down logs restart

# Build images without starting
build:
	$(DOCKER_COMPOSE) build

# Start services
up:
	$(DOCKER_COMPOSE) up --build -d

# Stop services
down:
	$(DOCKER_COMPOSE) down

# Show logs
logs:
	$(DOCKER_COMPOSE) logs -f

# Restart services
restart: down up