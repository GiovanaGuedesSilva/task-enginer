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

.PHONY: help build run test clean deps setup docker-up docker-down dev-full stop migrate-up migrate-down migrate-create migrate-force migrate-status migrate-script

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

# ===================================================== MIGRATION COMMANDS

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

migrate-status: ## Show migration status
	@echo "Migration status:"
	migrate -path $(MIGRATIONS_DIR) -database "$(DB_URL)" version

# ===================================================== MIGRATION SCRIPT COMMANDS (Alternative)

migrate-script-up: ## Run migrations using script (up)
	@echo "Running migrations using script..."
	@chmod +x scripts/run-migrations.sh
	./scripts/run-migrations.sh up

migrate-script-down: ## Run migrations using script (down)
	@echo "Running migrations rollback using script..."
	@chmod +x scripts/run-migrations.sh
	./scripts/run-migrations.sh down

migrate-script-force: ## Force migration using script (use: make migrate-script-force VERSION=1)
	@echo "Forcing migration using script..."
	@chmod +x scripts/run-migrations.sh
	VERSION=$(VERSION) ./scripts/run-migrations.sh force

migrate-script-create: ## Create migration using script (use: make migrate-script-create NAME=migration_name)
	@echo "Creating migration using script..."
	@chmod +x scripts/run-migrations.sh
	NAME=$(NAME) ./scripts/run-migrations.sh create

migrate-script-status: ## Show migration status using script
	@echo "Migration status using script:"
	@chmod +x scripts/run-migrations.sh
	./scripts/run-migrations.sh status

# ===================================================== DEVELOPMENT COMMANDS

dev-full: docker-up deps migrate-up run ## Start full development environment

dev-without-migrate: docker-up deps run ## Start development environment without migrations

clean: ## Clean build artifacts
	@echo "Cleaning build artifacts..."
	$(GO) clean
	@rm -rf $(BIN_DIR)
	@echo "Cleaned"

stop: docker-down clean ## Stop all services and clean

# ===================================================== UTILITY COMMANDS

install-migrate: ## Install migrate tool
	@echo "Installing migrate tool..."
	@if ! command -v migrate &> /dev/null; then \
		echo "Installing migrate via curl..."; \
		curl -L https://github.com/golang-migrate/migrate/releases/download/v4.18.3/migrate.linux-amd64.tar.gz | tar xvz; \
		sudo mv migrate /usr/local/bin/; \
		echo "Migrate installed successfully"; \
	else \
		echo "Migrate is already installed"; \
	fi

check-db: ## Check database connection
	@echo "Checking database connection..."
	@if pg_isready -h localhost -p 5432 -U postgres > /dev/null 2>&1; then \
		echo "Database is ready ✓"; \
	else \
		echo "Database is not ready. Start with: make docker-up"; \
		exit 1; \
	fi

logs: ## Show Docker logs
	@echo "Showing Docker logs..."
	docker-compose logs -f

restart: stop dev-full ## Restart all services
