#!/bin/bash

# Script to run Task Engine migrations
# Usage: ./scripts/run-migrations.sh [up|down|force|create|status|help]

set -e

MIGRATIONS_DIR="./migrations"
DB_URL="postgres://postgres:password@localhost:5432/task_engine?sslmode=disable"

check_postgres() {
    echo "[INFO] Checking PostgreSQL connection..."
    if ! pg_isready -h localhost -p 5432 -U postgres > /dev/null 2>&1; then
        echo "[ERROR] PostgreSQL is not running or unreachable"
        echo "Start PostgreSQL with: make docker-up"
        exit 1
    fi
    echo "[INFO] PostgreSQL is running"
}

check_migrate() {
    if ! command -v migrate &> /dev/null; then
        echo "[ERROR] 'migrate' tool not found"
        echo "Install it with: go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest"
        exit 1
    fi
}

run_migrations_up() {
    echo "[INFO] Running migrations UP..."
    migrate -path "$MIGRATIONS_DIR" -database "$DB_URL" up
    echo "[INFO] Migrations UP completed"
}

run_migrations_down() {
    echo "[INFO] Running migrations DOWN (rollback)..."
    migrate -path "$MIGRATIONS_DIR" -database "$DB_URL" down
    echo "[INFO] Migrations DOWN completed"
}

force_migration() {
    if [ -z "$VERSION" ]; then
        echo "[ERROR] Version not specified. Use: VERSION=1 ./scripts/run-migrations.sh force"
        exit 1
    fi
    echo "[INFO] Forcing migration version $VERSION..."
    migrate -path "$MIGRATIONS_DIR" -database "$DB_URL" force "$VERSION"
    echo "[INFO] Migration forced to version $VERSION"
}

create_migration() {
    if [ -z "$NAME" ]; then
        echo "[ERROR] Migration name not specified. Use: NAME=name ./scripts/run-migrations.sh create"
        exit 1
    fi
    echo "[INFO] Creating new migration: $NAME..."
    migrate create -ext sql -dir "$MIGRATIONS_DIR" -seq "$NAME"
    echo "[INFO] Migration '$NAME' created"
}

show_status() {
    echo "[INFO] Migration status:"
    migrate -path "$MIGRATIONS_DIR" -database "$DB_URL" version
}

show_help() {
    echo "Usage: $0 [command]"
    echo "Commands:"
    echo "  up       - Apply all pending migrations"
    echo "  down     - Rollback all migrations"
    echo "  force    - Force specific migration version (use VERSION=N)"
    echo "  create   - Create a new migration (use NAME=name)"
    echo "  status   - Show migration status"
    echo "  help     - Show this help message"
}

main() {
    check_postgres
    check_migrate

    case "${1:-help}" in
        up) run_migrations_up ;;
        down) run_migrations_down ;;
        force) force_migration ;;
        create) create_migration ;;
        status) show_status ;;
        help|*) show_help ;;
    esac
}

main "$@"