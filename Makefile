# Task Engine - Makefile

# Project variables
PROJECT_NAME := task-engine
CMD_DIR := ./cmd
BIN_DIR := ./bin

# Go
GO := go

# Database
DB_URL := postgres://postgres:password@localhost:5432/task_engine?sslmode=disable
MIGRATIONS_DIR := ./migrations

.PHONY: help build run test clean deps setup docker-up docker-down dev-full stop migrate-up migrate-down migrate-create migrate-force

.DEFAULT_GOAL := help

help: ## Show available commands
	@echo "=== $(PROJECT_NAME) - Available Commands ==="
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "  %-20s %s\n", $$1, $$2}'

setup: ## Setup development environment
	@echo "Setting up development environment..."
	@mkdir -p $(BIN_DIR)
	@if [ ! -f .env ]; then \
		echo "Creating .env file..."; \
		if [ -f .env.example ]; then \
			cp .env.example .env; \
		else \
			echo ".env.example not found"; \
		fi; \
	fi
	@echo "Installing dependencies..."
	$(GO) mod download
	$(GO) mod tidy
	@echo "Environment ready!"

deps: ## Install dependencies
	@echo "Installing dependencies..."
	$(GO) mod download
	$(GO) mod tidy

build: ## Build application
	@echo "Building application..."
	@mkdir -p $(BIN_DIR)
	$(GO) build -o $(BIN_DIR)/$(PROJECT_NAME) $(CMD_DIR)/main.go
	@echo "Built: $(BIN_DIR)/$(PROJECT_NAME)"

run: ## Run application
	@echo "Running application..."
	$(GO) run $(CMD_DIR)/main.go

test: ## Run tests
	@echo "Running tests..."
	$(GO) test -v ./...

docker-up: ## Start Docker services
	@echo "Starting Docker services..."
	docker-compose up -d
	@echo "Docker services started"

docker-down: ## Stop Docker services
	@echo "Stopping Docker services..."
	docker-compose down
	@echo "Docker services stopped"

migrate-up: ## Run database migrations up
	@echo "Running database migrations up..."
	migrate -path $(MIGRATIONS_DIR) -database "$(DB_URL)" up
	@echo "Migrations completed"

migrate-down: ## Rollback database migrations
	@echo "Rolling back database migrations..."
	migrate -path $(MIGRATIONS_DIR) -database "$(DB_URL)" down
	@echo "Migrations rollback completed"

migrate-force: ## Force migration version (use: make migrate-force VERSION=1)
	@echo "Forcing migration version to $(VERSION)..."
	migrate -path $(MIGRATIONS_DIR) -database "$(DB_URL)" force $(VERSION)
	@echo "Migration version forced to $(VERSION)"

migrate-create: ## Create new migration (use: make migrate-create NAME=migration_name)
	@echo "Creating new migration: $(NAME)..."
	migrate create -ext sql -dir $(MIGRATIONS_DIR) -seq $(NAME)
	@echo "Migration created: $(NAME)"

clean: ## Clean build artifacts
	@echo "Cleaning build artifacts..."
	$(GO) clean
	@rm -rf $(BIN_DIR)
	@echo "Cleaned"

dev-full: docker-up deps migrate-up run ## Start full development environment

stop: docker-down clean ## Stop all services and clean
