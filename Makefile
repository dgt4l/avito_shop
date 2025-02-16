# Environment
ENV_FILE := .env.dev
include $(ENV_FILE)

# Go Configuration
GO_MOD         := github.com/dgt4l/avito_shop
APP_NAME       := avito_shop
FULL_SRC_PATH  := $(GO_MOD)/cmd/$(APP_NAME)
BIN_DIR        := bin/
BUILD_FLAGS    := -ldflags="-s -w"
COVER_FLAGS    := -cover
COVER_DIR      := build/
BUILD_CMD      := go build $(BUILD_FLAGS) -o $(BIN_DIR)/$(APP_NAME) $(FULL_SRC_PATH)

# Tool Configurations
GCI_CONFIG_PATH    := .golangci.yml
D_COMPOSE_YML_PATH := deploy/docker-compose.yml
DOCKERFILE_PATH    := deploy/Dockerfile
D_COMPOSE_CMD      := docker compose

# Migration Configuration
MIGRATION_DIR           := migrations/
MIGRATION_COMMAND_SETUP := migrate -path $(MIGRATION_DIR) -database "$(DB_DRIVER)://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=disable"

# Phony Targets
.PHONY: help all build run coverage clean test lint hardlint fmt vet mod sec docker-up docker-down docker-build docker-buildup docker-restart migration-up migration-down migration-fix check-deps

all: run

build: check-deps ## Build binary
	@mkdir -p $(BIN_DIR)
	@$(BUILD_CMD)

run: build ## Build and run the application
	@./$(BIN_DIR)/$(APP_NAME)

coverage: BUILD_FLAGS+=$(COVER_FLAGS) ## Build and run with coverage
coverage: export GOCOVERDIR=$(COVER_DIR)
coverage: build run
	@mkdir -p $(COVER_DIR)
	@go tool covdata percent -i=$(COVER_DIR)

clean: ## Remove binaries and coverage files
	@rm -rf $(BIN_DIR) $(COVER_DIR)
	@go clean

test: ## Run tests
	@go test -v ./...

lint: ## Run linters
	@golangci-lint -v run --config $(GCI_CONFIG_PATH)

mod: ## Update dependencies
	@go mod tidy

docker-up: check-deps ## Start Docker containers
	@${D_COMPOSE_CMD} --env-file $(ENV_FILE) -f ${D_COMPOSE_YML_PATH} up -d

docker-down: check-deps ## Stop Docker containers
	@${D_COMPOSE_CMD} --env-file $(ENV_FILE) -f ${D_COMPOSE_YML_PATH} down

docker-build: check-deps ## Build Docker images
	@${D_COMPOSE_CMD} --env-file $(ENV_FILE) -f ${D_COMPOSE_YML_PATH} build

docker-buildup: docker-build docker-up ## Build and start Docker containers

docker-restart: docker-down docker-up ## Restart Docker containers

help: ## Show this help message
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST) | sort