.PHONY: help build build-all run clean test test-verbose test-coverage mocks sqlc-generate fmt lint tidy docker-up docker-down migrate-up migrate-down migrate-create migrate-force

ifeq ($(OS),Windows_NT)
    # Windows (CMD)
    RM = del /Q /F
    RMRF = rmdir /S /Q
    # On Windows, we ignore errors by prefixing with '-' instead of '|| true'
    IGNORE_ERR = -
else
    # Linux/macOS (Bash)
    RM = rm -f
    RMRF = rm -rf
    IGNORE_ERR = 
endif

# Variables
BINARY_NAME=travel-dream-board-api
BINARY_PATH=bin/$(BINARY_NAME)
MAIN_PATH=./cmd/api/main.go
COVERAGE_FILE=coverage.out

# Database migration variables
MIGRATIONS_PATH=sql/migrations
DATABASE_URL?=postgres://postgres:password@localhost:5432/travel-dream-board?sslmode=disable

# Default target
help:
	@echo "Available targets:"
	@echo "  make build              - Build the service binary"
	@echo "  make build-all          - Build FE bundle then BE binary"
	@echo "  make run                - Run the service locally"
	@echo "  make clean              - Clean build artifacts and generated files"
	@echo ""
	@echo "  make test               - Run all tests (regenerates mocks first)"
	@echo "  make test-verbose       - Run tests with verbose output"
	@echo "  make test-coverage      - Run tests with coverage report"
	@echo "  make mocks              - Generate mocks using mockery"
	@echo ""
	@echo "  make sqlc-generate      - Generate Go code from SQL queries"
	@echo "  make fmt                - Format code"
	@echo "  make lint               - Run golangci-lint"
	@echo "  make tidy               - Tidy Go modules"
	@echo ""
	@echo "  make docker-build       - Build the Docker image"
	@echo "  make docker-up          - Build and run the service container"
	@echo "  make docker-down        - Stop and remove the service container"
	@echo ""
	@echo "  make migrate-up         - Run database migrations"
	@echo "  make migrate-down       - Rollback last migration"
	@echo "  make migrate-create     - Create new migration (use NAME=migration_name)"
	@echo "  make migrate-force      - Force migration version (use VERSION=N)"
	@echo ""
	@echo "  make deploy             - Deploys a new version of the service"

# Build the binary
build:
	@echo "Building $(BINARY_NAME)..."
	@go build -o $(BINARY_PATH) $(MAIN_PATH)
	@echo "Build complete: $(BINARY_PATH)"

# Build FE bundle then BE binary (use this for prod-like local runs).
build-all:
	@echo "Building frontend bundle..."
	@cd ../frontend && npm run build
	@echo "Building $(BINARY_NAME)..."
	@go build -o $(BINARY_PATH) $(MAIN_PATH)
	@echo "Build complete: $(BINARY_PATH)"

# Run the service
run:
	@echo "Running $(BINARY_NAME)..."
	@go run $(MAIN_PATH)

# Clean build artifacts
clean:
	@echo Cleaning build artifacts...
	$(IGNORE_ERR)$(RMRF) bin
	$(IGNORE_ERR)$(RMRF) mocks
ifneq ($(COVERAGE_FILE),)
	$(IGNORE_ERR)$(RM) $(COVERAGE_FILE)
endif
	@echo Clean complete

# Generate mocks using mockery
mocks:
	@echo "Generating mocks with mockery..."
	@mockery
	@echo "Mocks generated"

# Run tests (regenerate mocks first)
test: mocks
	@echo "Running tests..."
	@go test ./...

# Run tests with verbose output
test-verbose: mocks
	@echo "Running tests with verbose output..."
	@go test -v ./...

# Run tests with coverage
ifeq ($(OS),Windows_NT)
test-coverage: mocks
	@echo "Running tests with coverage..."
	@go test -coverprofile=$(COVERAGE_FILE) ./internal/...
	@echo "Filtering coverage report..."
	@go run scripts/filter_coverage.go $(COVERAGE_FILE) "internal/db" "internal/testutil" "internal/logger" "internal/server" "mocks"
	@go tool cover -html=$(COVERAGE_FILE)
	@echo "Coverage report generated: $(COVERAGE_FILE)"
else
test-coverage: mocks
	@echo "Running tests with coverage..."
	@if docker --version | grep -qi podman; then \
		echo "Podman detected. Disabling Ryuk and enabling manual cleanup..."; \
		export TESTCONTAINERS_RYUK_DISABLED=true; \
		trap 'echo "Cleaning up test containers..."; docker ps -a --filter "label=org.testcontainers=true" -q | xargs -r docker rm -f' EXIT; \
		go test -coverprofile=$(COVERAGE_FILE) ./internal/...; \
	else \
		go test -coverprofile=$(COVERAGE_FILE) ./internal/...; \
	fi
	@echo "Filtering coverage report..."
	@go run scripts/filter_coverage.go $(COVERAGE_FILE) "internal/db" "internal/testutil" "internal/logger" "internal/server" "mocks"
	@go tool cover -html=$(COVERAGE_FILE)
	@echo "Coverage report generated: $(COVERAGE_FILE)"
endif

# Generate sqlc code
sqlc-generate:
	@echo "Generating sqlc code..."
	@sqlc generate
	@echo "sqlc generation complete"

# Format code
fmt:
	@echo "Formatting code..."
	@go fmt ./...
	@echo "Format complete"

# Run linter
lint:
	@echo "Running linter..."
	@golangci-lint run
	@echo "Linting complete"

# Tidy dependencies
tidy:
	@echo "Tidying dependencies..."
	@go mod tidy
	@go mod verify
	@echo "Tidy complete"

# Docker operations
docker-build:
	@echo "Building Docker image..."
	@docker build --build-arg VITE_BASE_URL=/ -t $(BINARY_NAME) .
	@echo "Docker image built: $(BINARY_NAME)"

ifeq ($(OS),Windows_NT)
docker-up: docker-build
	@echo "Starting container (Windows)..."
	@docker run -d \
		--name $(BINARY_NAME) \
		-p 8080:8080 \
		--env-file .env \
		--add-host=host.docker.internal:host-gateway \
		$(BINARY_NAME)
	@echo "Container started. Note: Ensure DATABASE_URL uses 'host.docker.internal' instead of 'localhost'"
else
docker-up: docker-build
	@echo "Starting container (Linux/Mac)..."
	@docker run -d \
		--name $(BINARY_NAME) \
		--network host \
		--env-file .env \
		$(BINARY_NAME)
	@echo "Container started"
endif

docker-down:
	@echo "Stopping container..."
	-@docker stop $(BINARY_NAME)
	-@docker rm $(BINARY_NAME)
	@echo "Container stopped"

# Database migration operations
migrate-up:
	@echo "Running migrations..."
	@migrate -path $(MIGRATIONS_PATH) -database "$(DATABASE_URL)" up
	@echo "Migrations complete"

migrate-down:
	@echo "Rolling back last migration..."
	@migrate -path $(MIGRATIONS_PATH) -database "$(DATABASE_URL)" down 1
	@echo "Rollback complete"

migrate-create:
ifndef NAME
	$(error NAME is required. Usage: make migrate-create NAME=migration_name)
endif
	@echo Creating migration: $(NAME)
	@migrate create -ext sql -dir $(MIGRATIONS_PATH) -seq $(NAME)
	@echo Migration files created in $(MIGRATIONS_PATH)

migrate-force:
	@if [ -z "$(VERSION)" ]; then \
		echo "Error: VERSION is required. Usage: make migrate-force VERSION=N"; \
		exit 1; \
	fi
	@echo "Forcing migration version to $(VERSION)..."
	@migrate -path $(MIGRATIONS_PATH) -database "$(DATABASE_URL)" force $(VERSION)
	@echo "Migration version forced"

deploy:
	@kubectl rollout restart deployment travel-dream-board-service -n dev-nl
