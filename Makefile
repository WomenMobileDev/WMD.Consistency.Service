.PHONY: dev dev-docker test build clean docker-prod docker-stop help

all: build

dev: 
	@echo "Starting development server with live reloading..."
	@air -c .air.toml

dev-docker:
	@echo "Starting development environment with Docker and live reloading..."
	@docker compose -f docker-compose.dev.yml up --build
test:
	@echo "Running tests..."
	@go test -v ./...

build:
	@echo "Building application..."
	@go build -o bin/main ./cmd/main.go
clean:
	@echo "Cleaning build artifacts..."
	@rm -rf bin tmp

docker-prod:
	@echo "Starting production environment with Docker..."
	@docker compose up --build -d

docker-stop:
	@echo "Stopping Docker containers..."
	@docker compose down

help:
	@echo "Available commands:"
	@echo "  make dev         - Run development server with Air (local)"
	@echo "  make dev-docker  - Run development server with Docker and live reloading"
	@echo "  make test        - Run tests"
	@echo "  make build       - Build the application"
	@echo "  make clean       - Clean build artifacts"
	@echo "  make docker-prod - Run production Docker environment"
	@echo "  make docker-stop - Stop Docker containers"
