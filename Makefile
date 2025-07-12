# Makefile for enterprisedata-exchange

# Variables
BINARY_NAME=enterprisedata-exchange
MAIN_PATH=./cmd/server
DOCKER_COMPOSE_FILE=docker-compose.yml
MIGRATIONS_PATH=./migrations

# Default target
.PHONY: help
help: ## Show this help message
	@echo "Available commands:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-15s\033[0m %s\n", $$1, $$2}'

# Development commands
.PHONY: build
build: ## Build the application
	@echo "Building application..."
	go build -o bin/$(BINARY_NAME) $(MAIN_PATH)

.PHONY: run
run: ## Run the application
	@echo "Running application..."
	CONFIG_PATH=configs/local.yaml go run $(MAIN_PATH)

.PHONY: dev
dev: wire build ## Generate wire, build and run the application
	@echo "Running application..."
	CONFIG_PATH=configs/local.yaml go run $(MAIN_PATH)

.PHONY: clean
clean: ## Clean build artifacts
	@echo "Cleaning..."
	rm -rf bin/
	rm -rf temp/
	go clean

# Wire generation
.PHONY: wire
wire: ## Generate wire dependency injection code
	@echo "Generating wire code..."
	wire $(MAIN_PATH)

# Testing
.PHONY: test
test: ## Run tests
	@echo "Running tests..."
	go test -v ./...

.PHONY: test-coverage
test-coverage: ## Run tests with coverage
	@echo "Running tests with coverage..."
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

# Code quality
.PHONY: lint
lint: ## Run linter
	@echo "Running linter..."
	golangci-lint run

.PHONY: fmt
fmt: ## Format code
	@echo "Formatting code..."
	go fmt ./...

.PHONY: vet
vet: ## Run go vet
	@echo "Running go vet..."
	go vet ./...

# Dependencies
.PHONY: deps
deps: ## Download dependencies
	@echo "Downloading dependencies..."
	go mod download

.PHONY: deps-tidy
deps-tidy: ## Tidy dependencies
	@echo "Tidying dependencies..."
	go mod tidy

# Docker commands
.PHONY: docker-build
docker-build: ## Build Docker image
	@echo "Building Docker image..."
	docker build -t $(BINARY_NAME) .

.PHONY: docker-up
docker-up: ## Start services with docker-compose
	@echo "Starting services with docker-compose..."
	docker-compose -f $(DOCKER_COMPOSE_FILE) up -d

.PHONY: docker-down
docker-down: ## Stop services with docker-compose
	@echo "Stopping services with docker-compose..."
	docker-compose -f $(DOCKER_COMPOSE_FILE) down

.PHONY: docker-logs
docker-logs: ## Show docker-compose logs
	docker-compose -f $(DOCKER_COMPOSE_FILE) logs -f

# Database commands
.PHONY: migrate-up
migrate-up: ## Run database migrations up
	@echo "Running migrations up..."
	goose -dir $(MIGRATIONS_PATH) sqlite3 storage/app.db up

.PHONY: migrate-down
migrate-down: ## Run database migrations down
	@echo "Running migrations down..."
	goose -dir $(MIGRATIONS_PATH) sqlite3 storage/app.db down

.PHONY: migrate-create
migrate-create: ## Create new migration file (usage: make migrate-create NAME=migration_name)
	@echo "Creating migration $(NAME)..."
	goose -dir $(MIGRATIONS_PATH) create $(NAME) sql

.PHONY: migrate-status
migrate-status: ## Show migration status
	@echo "Migration status..."
	goose -dir $(MIGRATIONS_PATH) sqlite3 storage/app.db status

.PHONY: migrate-reset
migrate-reset: ## Reset all migrations
	@echo "Resetting all migrations..."
	goose -dir $(MIGRATIONS_PATH) sqlite3 storage/app.db reset

.PHONY: migrate-version
migrate-version: ## Show current migration version
	@echo "Current migration version..."
	goose -dir $(MIGRATIONS_PATH) sqlite3 storage/app.db version

.PHONY: db-create
db-create: ## Create database file if it doesn't exist
	@echo "Creating database..."
	@mkdir -p storage
	@touch storage/app.db

# Installation commands
.PHONY: install-tools
install-tools: ## Install development tools
	@echo "Installing development tools..."
	go install github.com/google/wire/cmd/wire@latest
	go install github.com/pressly/goose/v3/cmd/goose@latest
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

# All-in-one commands
.PHONY: setup
setup: install-tools deps db-create wire ## Setup development environment

.PHONY: check
check: fmt vet lint test ## Run all checks

.PHONY: full-build
full-build: clean wire fmt vet build ## Clean, generate, format, vet and build 