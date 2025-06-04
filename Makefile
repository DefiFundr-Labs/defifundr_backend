# Observability targets
validate-observability: ## Validate observability setup
	@echo "🔍 Validating observability setup..."
	@echo "📊 Checking metrics endpoint..."
	@if curl -f -s http://localhost:8080/metrics > /dev/null; then \
		echo "✅ Metrics endpoint is accessible"; \
	else \
		echo "❌ Metrics endpoint is not accessible - make sure the server is running"; \
		exit 1; \
	fi
	@echo "�� Checking health endpoint..."
	@if curl -f -s http://localhost:8080/health > /dev/null; then \
		echo "✅ Health endpoint is accessible"; \
	else \
		echo "❌ Health endpoint is not accessible"; \
		exit 1; \
	fi
	@echo "📈 Validating Prometheus metrics format..."
	@if curl -s http://localhost:8080/metrics | grep -q "^# HELP"; then \
		echo "✅ Prometheus metrics format is valid"; \
	else \
		echo "❌ Prometheus metrics format is invalid"; \
		exit 1; \
	fi
	@echo "🎯 All observability checks passed!"
  # Load .env file if it exists
ifneq (,$(wildcard .env))
    include .env
    export
endif

DB_URL = ${DB_SOURCE}
DB_NAME ?= ${DB_NAME}
DB_USER ?= ${DB_USER}
DB_PASSWORD ?= ${DB_PASSWORD}
DB_PORT ?= 5433

# Docker commands
docker-up:
	docker-compose up -d

docker-down:
	docker-compose down

docker-logs:
	docker-compose logs -f

docker-ps:
	docker-compose ps

docker-build:
	docker-compose build

docker-restart:
	docker-compose restart

# Database commands
postgres:
	docker run --name $(DB_NAME) -p $(DB_PORT):5432 -e POSTGRES_USER=$(DB_USER) -e POSTGRES_PASSWORD=$(DB_PASSWORD) -d postgres:15-alpine

createdb:
	docker exec -it $(DB_NAME) createdb --username=$(DB_USER) --owner=$(DB_USER) $(DB_NAME)

dockerlogs:
	docker logs defi

dropdb:
	docker exec -it defi dropdb $(DB_NAME) 

# Migration commands (using goose)
migrate-create:
	@read -p "Enter migration name: " name; \
	goose -dir db/migrations create $${name} sql

migrate-up:
	goose -dir db/migrations postgres "$(DB_URL)" up

migrate-up-one:
	goose -dir db/migrations postgres "$(DB_URL)" up-by-one

migrate-down:
	goose -dir db/migrations postgres "$(DB_URL)" down

migrate-down-one:
	goose -dir db/migrations postgres "$(DB_URL)" down-by-one

migrate-status:
	goose -dir db/migrations postgres "$(DB_URL)" status

migrate-reset:
	goose -dir db/migrations postgres "$(DB_URL)" reset 

# Documentation commands
# Database documentation
DB_DOCS_DIR=cmd/api/docs/db_diagram
DB_DBML_FILE=$(DB_DOCS_DIR)/db.dbml

.PHONY: db_docs db_diagram

# Create directory if it doesn't exist
$(DB_DOCS_DIR):
	mkdir -p $(DB_DOCS_DIR)

# Generate DBML file from the dbdiagram source
db_dbml: $(DB_DOCS_DIR)
	cp dbdiagram.txt $(DB_DBML_FILE)

db_docs: 
	@echo "Installing dbdocs CLI..."
	npm install -g dbdocs
	@echo "Building database documentation..."
	dbdocs build $(DB_DBML_FILE)

db_schema:
	dbml2sql --postgress -o cmd/api/docs/db_diagram/schema.sql cmd/api/docs/db_diagram/db.dbml

# Development commands
sqlc:
	sqlc generate

test:
	go test -v -cover ./...

server:
	go run cmd/api/main.go

air:
	air

seed:
	go run cmd/seed/main.go

# Mock generation
mock:
	@mkdir -p db/mock
	mockgen -package sqlc -destination db/sqlc/querier.go -source db/sqlc/querier.go

# Linting
lint:
	golangci-lint run ./...

# Build
build:
	go build -o bin/api cmd/api/main.go

# Clean
clean:
	rm -rf bin/
	rm -rf smart-contracts/build/

# Install tools
install-tools:
	go install github.com/air-verse/air@latest
	go install github.com/pressly/goose/v3/cmd/goose@latest
	go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest
	go install github.com/golang/mock/mockgen@latest
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

# Swagger
swagger:
	swag init -g cmd/api/main.go -o cmd/api/docs

# Help command
help:
	@echo "DefiFundr - Blockchain Payroll System"
	@echo ""
	@echo "Database Commands:"
	@echo "  postgres            - Run PostgreSQL container"
	@echo "  createdb            - Create database"
	@echo "  dropdb              - Drop database"
	@echo "  dockerlogs          - View PostgreSQL container logs"
	@echo ""
	@echo "Migration Commands (Goose):"
	@echo "  migrate-create      - Create a new migration file"
	@echo "  migrate-up          - Run all pending migrations"
	@echo "  migrate-up-one      - Run the next pending migration"
	@echo "  migrate-down        - Revert the last migration"
	@echo "  migrate-down-one    - Revert the last migration"
	@echo "  migrate-status      - Show migration status"
	@echo "  migrate-reset       - Revert all migrations"
	@echo ""
	@echo "Docker Commands:"
	@echo "  docker-up           - Start all Docker containers"
	@echo "  docker-down         - Stop and remove all Docker containers"
	@echo "  docker-logs         - View logs from all Docker containers"
	@echo "  docker-ps           - List running Docker containers"
	@echo "  docker-build        - Rebuild Docker images"
	@echo "  docker-restart      - Restart all Docker containers"
	@echo ""
	@echo "Development Commands:"
	@echo "  sqlc                - Generate SQL code with sqlc"
	@echo "  mock                - Generate mock code for testing"
	@echo "  server              - Run the API server"
	@echo "  air                 - Run the server with hot reload"
	@echo "  test                - Run tests"
	@echo "  lint                - Run linters"
	@echo "  build               - Build the application"
	@echo "  clean               - Clean build artifacts"
	@echo "  seed                - Seed the database with test data"
	@echo ""
	@echo "Smart Contract Commands:"
	@echo "  gencontract         - Generate Go bindings for smart contracts"
	@echo ""
	@echo "Documentation Commands:"
	@echo "  db_docs             - Generate DB documentation with dbdocs"
	@echo "  db_schema           - Generate SQL schema from DBML"
	@echo ""
	@echo "Setup Commands:"
	@echo "  install-tools       - Install development tools"
	@echo ""
	@echo "Swagger Commands:"
	@echo "  swagger             - Generate Swagger documentation"
	@echo ""

.PHONY: postgres createdb dockerlogs dropdb migrate-create migrate-up migrate-up-one migrate-down migrate-down-one migrate-status migrate-reset migrateup migrateup1 migratedown migratedown1 db_docs db_schema sqlc test server air mock gencontract docker-up docker-down docker-logs docker-ps docker-build docker-restart help seed lint build clean install-tools
