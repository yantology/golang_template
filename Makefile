# Go Backend Template Makefile

# Variables
BINARY_NAME=golang-template
DOCKER_IMAGE=golang-template
DOCKER_TAG=latest

# Default target
.DEFAULT_GOAL := help

# Development commands
.PHONY: dev
dev: ## Start Go application
	go run cmd/api/main.go


.PHONY: db-up
db-up: ## Start database services (PostgreSQL + Adminer)
	docker-compose up -d

.PHONY: db-down
db-down: ## Stop database services
	docker-compose down

.PHONY: db-clean
db-clean: ## Clean database (remove volumes)
	docker-compose down -v

.PHONY: db-logs
db-logs: ## Show database logs
	docker-compose logs -f postgres

# Build commands
.PHONY: build
build: ## Build the application binary
	CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o bin/$(BINARY_NAME) cmd/api/main.go

.PHONY: build-docker
build-docker: ## Build Docker image
	docker build -t $(DOCKER_IMAGE):$(DOCKER_TAG) .

.PHONY: clean
clean: ## Clean build artifacts
	rm -rf bin/
	docker system prune -f

# Database commands
.PHONY: db-connect
db-connect: ## Connect to database via psql
	docker-compose exec postgres psql -U postgres -d golang_template_dev

.PHONY: db-reset
db-reset: ## Reset database (remove all data)
	docker-compose down -v
	docker-compose up -d

# Migration commands (requires migrate CLI tool)
# Database URL can be set via DATABASE_URL environment variable
# For development, defaults to local postgres
DB_URL ?= postgres://postgres:dev_password@localhost:5432/golang_template_dev?sslmode=disable
MIGRATION_PATH ?= ./internal/data/migrations

.PHONY: migrate-up
migrate-up: ## Run database migrations up (set DATABASE_URL for production)
	migrate -path $(MIGRATION_PATH) -database "$(DB_URL)" up

.PHONY: migrate-down
migrate-down: ## Run database migrations down (set DATABASE_URL for production)
	migrate -path $(MIGRATION_PATH) -database "$(DB_URL)" down

.PHONY: migrate-create
migrate-create: ## Create new migration file (usage: make migrate-create NAME=migration_name)
	@if [ -z "$(NAME)" ]; then echo "Usage: make migrate-create NAME=migration_name"; exit 1; fi
	migrate create -ext sql -dir $(MIGRATION_PATH) -seq $(NAME)

.PHONY: migrate-status
migrate-status: ## Check migration status (set DATABASE_URL for production)
	migrate -path $(MIGRATION_PATH) -database "$(DB_URL)" version

.PHONY: migrate-force
migrate-force: ## Force migration version (usage: make migrate-force VERSION=1)
	@if [ -z "$(VERSION)" ]; then echo "Usage: make migrate-force VERSION=1"; exit 1; fi
	migrate -path $(MIGRATION_PATH) -database "$(DB_URL)" force $(VERSION)

.PHONY: migrate-up-one
migrate-up-one: ## Run one migration up (set DATABASE_URL for production)
	migrate -path $(MIGRATION_PATH) -database "$(DB_URL)" up 1

.PHONY: migrate-down-one
migrate-down-one: ## Run one migration down (set DATABASE_URL for production)
	migrate -path $(MIGRATION_PATH) -database "$(DB_URL)" down 1

# Advanced migration commands using helper script
.PHONY: migrate
migrate: ## Run migration script (usage: make migrate ARGS="up|down|status|create name")
	@./scripts/migrate.sh $(ARGS)

.PHONY: migrate-prod
migrate-prod: ## Run migration for production (requires production env vars)
	@APP_SERVER_ENV=production ./scripts/migrate.sh $(ARGS)

# Testing commands
.PHONY: test
test: ## Run all tests
	go test -v ./...

.PHONY: test-coverage
test-coverage: ## Run tests with coverage
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

.PHONY: test-integration
test-integration: ## Run integration tests
	go test -v ./tests/integration/...

.PHONY: test-e2e
test-e2e: ## Run E2E tests
	go test -v ./tests/e2e/...

.PHONY: test-unit
test-unit: ## Run unit tests only
	go test -v ./internal/...

# Code quality commands
.PHONY: fmt
fmt: ## Format Go code
	go fmt ./...

.PHONY: lint
lint: ## Run golangci-lint
	golangci-lint run

.PHONY: vet
vet: ## Run go vet
	go vet ./...

.PHONY: tidy
tidy: ## Run go mod tidy
	go mod tidy

.PHONY: verify
verify: fmt vet lint tidy test ## Run all verification steps

# Development tools
.PHONY: tools
tools: ## Install development tools
	go install github.com/golang-migrate/migrate/v4/cmd/migrate@latest
	go install github.com/vektra/mockery/v2@latest
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest


# Utility commands
.PHONY: env
env: ## Copy environment template
	cp .env.example .env
	@echo "Environment file created. Please edit .env with your settings."

.PHONY: adminer
adminer: ## Open Adminer in browser (requires xdg-open or open command)
	@command -v xdg-open >/dev/null 2>&1 && xdg-open http://localhost:8081 || \
	 command -v open >/dev/null 2>&1 && open http://localhost:8081 || \
	 echo "Please open http://localhost:8081 in your browser"

.PHONY: api
api: ## Open API in browser
	@command -v xdg-open >/dev/null 2>&1 && xdg-open http://localhost:8080/health || \
	 command -v open >/dev/null 2>&1 && open http://localhost:8080/health || \
	 echo "Please open http://localhost:8080 in your browser"


# Help command
.PHONY: help
help: ## Show this help message
	@echo "Go Backend Template - Available commands:"
	@echo ""
	@awk 'BEGIN {FS = ":.*##"; printf "\033[36m\033[0m"} /^[a-zA-Z_-]+:.*?##/ { printf "  \033[36m%-20s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)
	@echo ""
	@echo "Examples:"
	@echo "  make db-up              # Start PostgreSQL + Adminer"
	@echo "  make migrate-up         # Run database migrations"
	@echo "  make dev                # Start Go application"
	@echo "  make migrate-create NAME=add_users_table"