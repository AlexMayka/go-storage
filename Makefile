# Go Storage Makefile
.PHONY: help build run stop clean logs test

# Default target
help:
	@echo "Go Storage Project Commands:"
	@echo ""
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

# Docker commands
build:
	docker-compose build

up:
	docker-compose up -d

down:
	docker-compose down

restart:
	docker-compose restart

stop:
	docker-compose stop

start:
	docker-compose start

logs:
	docker-compose logs -f

logs-app:
	docker-compose logs -f app

logs-db:
	docker-compose logs -f db

logs-minio:
	docker-compose logs -f minio

logs-migrate:
	docker-compose logs migrate

# Development commands
dev:
	docker-compose up -d db minio
	docker-compose build app
	docker-compose up app

dev-down:
	docker-compose down

# Database commands
db-migrate:
	docker-compose run --rm migrate

db-migrate-down:
	docker-compose run --rm migrate goose -dir /migrations postgres "postgres://admin:admin@db:5432/storage?sslmode=disable" down

db-migrate-status:
	docker-compose run --rm migrate -dir /migrations postgres "postgres://admin:admin@db:5432/storage?sslmode=disable" status

db-migrate-reset:
	docker-compose run --rm migrate goose -dir /migrations postgres "postgres://admin:admin@db:5432/storage?sslmode=disable" reset

db-shell:
	docker-compose exec db psql -U admin -d storage

# MinIO commands
minio-console:
	@echo "MinIO Console: http://localhost:9001"
	@echo "Username: admin"
	@echo "Password: secret123"

# Application commands
app-shell:
	docker-compose exec app sh

# Cleanup commands
clean:
	docker-compose down -v --remove-orphans
	docker system prune -f

clean-all:
	docker-compose down -v --remove-orphans --rmi all
	docker system prune -a -f

# Testing
test:
	go test -v ./...

test-coverage:
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out

# Documentation
swagger:
	swag init -g cmd/api/main.go -o cmd/api/docs

# Build local binary
build-local:
	go build -o bin/go-storage ./cmd/api

run-local:
	./bin/go-storage

# Status
status:
	docker-compose ps

health:
	@echo "Checking service health..."
	@docker-compose exec app curl -f -s http://localhost:8080/swagger/index.html > /dev/null && echo "✓ App is healthy" || echo "✗ App is unhealthy"
	@docker-compose exec db pg_isready -U admin -d storage && echo "✓ Database is healthy" || echo "✗ Database is unhealthy"
	@docker-compose exec minio curl -f http://localhost:9000/minio/health/live && echo "✓ MinIO is healthy" || echo "✗ MinIO is unhealthy"

# Environment
env-copy:
	cp .env.example .env
	@echo "Created .env file from .env.example"
	@echo "Please edit .env file with your settings"

# Quick start
quick-start: env-copy build up