# Environment
ENV_FILE := deploy/.env.dev
include $(ENV_FILE)

# Go Configuration
GO_MOD         := github.com/dgt4l/avito_shop
APP_NAME       := avito_shop
FULL_SRC_PATH  := $(GO_MOD)/cmd/$(APP_NAME)
BIN_DIR        := bin/
BUILD_FLAGS    := -ldflags="-s -w"
BUILD_CMD      := go build $(BUILD_FLAGS) -o $(BIN_DIR)/$(APP_NAME) $(FULL_SRC_PATH)
COV_FILE	   := cov.out

# Tool Configurations
GCI_CONFIG_PATH    := .golangci.yml
D_COMPOSE_YML_PATH := deploy/docker-compose.yml
DOCKERFILE_PATH    := deploy/Dockerfile
D_COMPOSE_CMD      := docker compose

# Migration Configuration
MIGRATION_DIR           := migrations/
MIGRATION_COMMAND_SETUP := migrate -path $(MIGRATION_DIR) -database "$(DB_DRIVER)://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=disable"

# Phony Targets
.PHONY: help all build run coverage clean test lint mod docker-up docker-down docker-buildup docker-restart

all: run

build:  ## Build binary
	@mkdir -p $(BIN_DIR)
	@$(BUILD_CMD)

run: build ## Build and run the application
	@./$(BIN_DIR)/$(APP_NAME)

coverage: ## Run tests with coverage
	@go test -coverprofile=${COV_FILE} ./internal/... && go tool cover -func=${COV_FILE}

clean: ## Remove binaries
	@rm -rf $(BIN_DIR) 
	@rm ${COV_FILE}
	@go clean

test: ## Run tests
	@go test ./internal/...

lint: ## Run linters
	@golangci-lint -v run --config $(GCI_CONFIG_PATH)

mod: ## Update dependencies
	@go mod tidy

docker-up: ## Start Docker containers
	@${D_COMPOSE_CMD} --env-file $(ENV_FILE) -f ${D_COMPOSE_YML_PATH} up -d

docker-down: ## Stop Docker containers
	@${D_COMPOSE_CMD} --env-file $(ENV_FILE) -f ${D_COMPOSE_YML_PATH} down

docker-buildup: ## Build and Start Docker Containers
	@${D_COMPOSE_CMD} --env-file $(ENV_FILE) -f ${D_COMPOSE_YML_PATH} up --build -d

docker-restart: docker-down docker-up ## Restart Docker containers

help: ## Show this help message
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST) | sort