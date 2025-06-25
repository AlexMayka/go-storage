# Go Storage Makefile
.PHONY: help build run stop clean logs test

# Default target
help: ## Show this help message
	@echo "Go Storage Project Commands:"
	@echo ""
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

# Docker commands
build: ## Build all Docker images
	docker-compose build

up: ## Start all services in detached mode
	docker-compose up -d

down: ## Stop and remove all containers
	docker-compose down

restart: ## Restart all services
	docker-compose restart

stop: ## Stop all services without removing containers
	docker-compose stop

start: ## Start existing containers
	docker-compose start

logs: ## Show logs from all services
	docker-compose logs -f

logs-app: ## Show logs from app service only
	docker-compose logs -f app

logs-db: ## Show logs from database service only
	docker-compose logs -f db

logs-minio: ## Show logs from MinIO service only
	docker-compose logs -f minio

logs-migrate: ## Show logs from migration service only
	docker-compose logs migrate

# Development commands
dev: ## Start services for development (rebuild app)
	docker-compose up -d db minio
	docker-compose build app
	docker-compose up app

dev-down: ## Stop development environment
	docker-compose down

# Database commands
db-migrate: ## Run database migrations manually
	docker-compose run --rm migrate

db-migrate-down: ## Rollback last migration
	docker-compose run --rm migrate goose -dir /migrations postgres "postgres://admin:admin@db:5432/storage?sslmode=disable" down

db-migrate-status: ## Show migration status
	docker-compose run --rm migrate -dir /migrations postgres "postgres://admin:admin@db:5432/storage?sslmode=disable" status

db-migrate-reset: ## Reset database (dangerous!)
	docker-compose run --rm migrate goose -dir /migrations postgres "postgres://admin:admin@db:5432/storage?sslmode=disable" reset

db-shell: ## Access PostgreSQL shell
	docker-compose exec db psql -U admin -d storage

# MinIO commands
minio-console: ## Open MinIO console URL
	@echo "MinIO Console: http://localhost:9001"
	@echo "Username: admin"
	@echo "Password: secret123"

# Application commands
app-shell: ## Access application container shell
	docker-compose exec app sh

# Cleanup commands
clean: ## Remove all containers, networks, and volumes
	docker-compose down -v --remove-orphans
	docker system prune -f

clean-all: ## Remove everything including images
	docker-compose down -v --remove-orphans --rmi all
	docker system prune -a -f

# Testing
test: ## Run tests locally
	go test -v ./...

test-coverage: ## Run tests with coverage
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out

# Documentation
swagger: ## Generate/update Swagger documentation
	swag init -g cmd/api/main.go -o cmd/api/docs

# Build local binary
build-local: ## Build application locally
	go build -o bin/go-storage ./cmd/api

run-local: ## Run application locally (requires local DB and MinIO)
	./bin/go-storage

# Status
status: ## Show status of all services
	docker-compose ps

health: ## Check health of all services
	@echo "Checking service health..."
	@docker-compose exec app curl -f -s http://localhost:8080/swagger/index.html > /dev/null && echo "✓ App is healthy" || echo "✗ App is unhealthy"
	@docker-compose exec db pg_isready -U admin -d storage && echo "✓ Database is healthy" || echo "✗ Database is unhealthy"
	@docker-compose exec minio curl -f http://localhost:9000/minio/health/live && echo "✓ MinIO is healthy" || echo "✗ MinIO is unhealthy"

# Environment
env-copy: ## Copy example environment file
	cp .env.example .env
	@echo "Created .env file from .env.example"
	@echo "Please edit .env file with your settings"

# Quick start
quick-start: env-copy build up ## Quick start: copy env, build and start all services
	@echo ""
	@echo "🚀 Go Storage is starting up..."
	@echo ""
	@echo "Services:"
	@echo "  📱 App:          http://localhost:8080"
	@echo "  📚 Swagger:      http://localhost:8080/swagger/index.html"
	@echo "  🗄️  MinIO Console: http://localhost:9001 (admin/secret123)"
	@echo "  🐘 PostgreSQL:   localhost:5432 (admin/admin/storage)"
	@echo ""
	@echo "Use 'make logs' to see logs from all services"
	@echo "Use 'make health' to check service health"