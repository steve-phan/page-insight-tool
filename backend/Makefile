.PHONY: build run test test-coverage clean deps tools docker-build docker-run dev fmt vet lint check swagger swagger-serve

# Application variables
APP_NAME = page-insight-tool
BIN_DIR = bin
MAIN_GO_PATH = cmd/main.go

# Version information
VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
COMMIT ?= $(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")
BUILD_DATE ?= $(shell date -u +"%Y-%m-%dT%H:%M:%SZ")

# Build flags
LDFLAGS = -X main.version=$(VERSION) \
          -X main.buildDate=$(BUILD_DATE) \
          -X main.gitCommit=$(COMMIT)

# Go build
build:
	@echo "Building $(APP_NAME)..."
	go build -ldflags "$(LDFLAGS)" -o $(BIN_DIR)/$(APP_NAME) $(MAIN_GO_PATH)

# Run the application
run: build
	@echo "Running $(APP_NAME)..."
	$(BIN_DIR)/$(APP_NAME)

# Development mode with hot reload
dev:
	@echo "Starting development server..."
	go run $(MAIN_GO_PATH)

# Run tests
test:
	@echo "Running tests..."
	go test -v ./...

# Run tests with coverage
test-coverage:
	@echo "Running tests with coverage..."
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

# Run specific test packages
test-config:
	go test -v ./internal/config/...

test-handlers:
	go test -v ./internal/handlers/...

test-routes:
	go test -v ./internal/routes/...

test-integration:
	go test -v ./tests/...

# Format code
fmt:
	@echo "Formatting code..."
	go fmt ./...

# Vet code
vet:
	@echo "Vetting code..."
	go vet ./...

# Lint code (requires golangci-lint)
lint:
	@echo "Linting code..."
	golangci-lint run

# Run all checks
check: fmt vet test
	@echo "All checks passed!"

# Download dependencies
deps:
	@echo "Downloading Go dependencies..."
	go mod download
	go mod tidy

# Install development tools
tools:
	@echo "Installing development tools..."
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	go install github.com/stretchr/testify@latest

# Docker build
docker-build:
	@echo "Building Docker image..."
	docker build -t $(APP_NAME):$(VERSION) .

# Docker run
docker-run: docker-build
	@echo "Running Docker container..."
	docker run -p 8080:8080 $(APP_NAME):$(VERSION)

# Clean build artifacts
clean:
	@echo "Cleaning..."
	rm -rf $(BIN_DIR)/
	rm -f coverage.out coverage.html

# Generate Swagger documentation
swagger:
	@echo "Generating Swagger documentation..."
	@which swag > /dev/null || (echo "swag not found. Installing..." && go install github.com/swaggo/swag/cmd/swag@latest)
	swag init -g cmd/main.go -o ./docs/api --parseDependency --parseInternal --parseDepth 2
	@echo "Converting to OpenAPI 3.0..."
	@which swag2openapi > /dev/null || (echo "swag2openapi not found. Installing..." && go install github.com/swaggo/swag/cmd/swag2openapi@latest)
	swag2openapi ./docs/api/swagger.json -o ./docs/api/openapi.json

# Serve Swagger UI (requires running server)
swagger-serve:
	@echo "Swagger UI available at http://localhost:8080/swagger/index.html"
	@echo "Make sure the server is running first!"

# Help
help:
	@echo "Available targets:"
	@echo "  build         - Build the application"
	@echo "  run           - Build and run the application"
	@echo "  dev           - Run in development mode"
	@echo "  test          - Run all tests"
	@echo "  test-coverage - Run tests with coverage report"
	@echo "  test-config   - Run config tests only"
	@echo "  test-handlers - Run handler tests only"
	@echo "  test-routes   - Run route tests only"
	@echo "  test-integration - Run integration tests only"
	@echo "  fmt           - Format code"
	@echo "  vet           - Vet code"
	@echo "  lint          - Lint code"
	@echo "  check         - Run all checks"
	@echo "  deps          - Download dependencies"
	@echo "  tools         - Install development tools"
	@echo "  docker-build  - Build Docker image"
	@echo "  docker-run    - Build and run Docker container"
	@echo "  clean         - Clean build artifacts"
	@echo "  swagger       - Generate Swagger/OpenAPI documentation"
	@echo "  swagger-serve - Display Swagger UI URL (requires running server)"
	@echo "  help          - Show this help"
