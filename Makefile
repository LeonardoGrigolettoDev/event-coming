.PHONY: build run run-worker test test-coverage migrate-up migrate-down docker-up docker-down swagger lint tidy clean

# Build targets
build:
	@echo "Building binaries..."
	@go build -o bin/api cmd/api/main.go
	@go build -o bin/worker cmd/worker/main.go
	@echo "Build complete!"

# Run targets
run:
	@echo "Running API server..."
	@go run cmd/api/main.go

run-worker:
	@echo "Running workers..."
	@go run cmd/worker/main.go

# Test targets
test:
	@echo "Running tests..."
	@go test -v ./...

test-coverage:
	@echo "Running tests with coverage..."
	@go test -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

# Database migration targets
migrate-up:
	@echo "Running migrations..."
	@migrate -path migrations -database "$(DATABASE_URL)" up

migrate-down:
	@echo "Rolling back migrations..."
	@migrate -path migrations -database "$(DATABASE_URL)" down

migrate-create:
	@echo "Creating new migration: $(name)"
	@migrate create -ext sql -dir migrations -seq $(name)

# Docker targets
docker-up:
	@echo "Starting Docker containers..."
	@docker-compose up -d

docker-down:
	@echo "Stopping Docker containers..."
	@docker-compose down

docker-logs:
	@docker-compose logs -f

docker-build:
	@echo "Building Docker images..."
	@docker-compose build

# Development targets
swagger:
	@echo "Generating Swagger documentation..."
	@swag init -g cmd/api/main.go -o docs/swagger

lint:
	@echo "Running linter..."
	@golangci-lint run

tidy:
	@echo "Tidying modules..."
	@go mod tidy

clean:
	@echo "Cleaning build artifacts..."
	@rm -rf bin/
	@rm -f coverage.out coverage.html
	@echo "Clean complete!"

# Install tools
install-tools:
	@echo "Installing development tools..."
	@go install github.com/swaggo/swag/cmd/swag@latest
	@go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	@go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
	@echo "Tools installed!"

# Help
help:
	@echo "Available targets:"
	@echo "  build           - Build API and worker binaries"
	@echo "  run             - Run API server"
	@echo "  run-worker      - Run workers"
	@echo "  test            - Run tests"
	@echo "  test-coverage   - Run tests with coverage report"
	@echo "  migrate-up      - Run database migrations"
	@echo "  migrate-down    - Rollback database migrations"
	@echo "  migrate-create  - Create new migration (usage: make migrate-create name=migration_name)"
	@echo "  docker-up       - Start Docker containers"
	@echo "  docker-down     - Stop Docker containers"
	@echo "  docker-logs     - View Docker logs"
	@echo "  docker-build    - Build Docker images"
	@echo "  swagger         - Generate Swagger documentation"
	@echo "  lint            - Run linter"
	@echo "  tidy            - Tidy Go modules"
	@echo "  clean           - Clean build artifacts"
	@echo "  install-tools   - Install development tools"
