# Build variables
BINARY_NAME=page-insight-tool
VERSION=1.0.0
BUILD_DATE=$(shell date +%Y-%m-%dT%H:%M:%S%z)
GIT_COMMIT=$(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")

# Go build flags
LDFLAGS=-ldflags "-X main.version=$(VERSION) -X main.buildDate=$(BUILD_DATE) -X main.gitCommit=$(GIT_COMMIT)"

# Default target
all: build

# Build the application
build:
	@echo "Building $(BINARY_NAME)..."
	go build $(LDFLAGS) -o bin/$(BINARY_NAME) ./cmd/

# Run the application
run: build
	@echo "Running $(BINARY_NAME)..."
	./bin/$(BINARY_NAME)

# Run with hot reload (requires air installed)
dev:
	@echo "Starting development server with hot reload..."
	air

# Run tests
test:
	@echo "Running tests..."
	go test -v ./...

# Run tests with coverage
test-coverage:
	@echo "Running tests with coverage..."
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out

# Clean build artifacts
clean:
	@echo "Cleaning..."
	rm -rf bin/ coverage.out

# Format code
fmt:
	@echo "Formatting code..."
	gofmt -w .

# Vet code
vet:
	@echo "Vetting code..."
	go vet ./...

# Lint code
lint:
	@echo "Linting code..."
	golangci-lint run

# Check code quality (fmt + vet + lint)
check: fmt vet lint

# Install dependencies
deps:
	@echo "Installing dependencies..."
	go mod download
	go mod verify

# Install development tools
tools:
	@echo "Installing development tools..."
	go install github.com/cosmtrek/air@latest
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

# Docker build
docker-build:
	@echo "Building Docker image..."
	docker build -t $(BINARY_NAME):$(VERSION) .

# Docker run
docker-run: docker-build
	@echo "Running Docker container..."
	docker run -p 8080:8080 $(BINARY_NAME):$(VERSION)

# Show help
help:
	@echo "Available targets:"
	@echo "  build         - Build the application"
	@echo "  run           - Build and run the application"
	@echo "  dev           - Run with hot reload (requires air)"
	@echo "  test          - Run tests"
	@echo "  test-coverage - Run tests with coverage report"
	@echo "  clean         - Clean build artifacts"
	@echo "  fmt           - Format code"
	@echo "  vet           - Vet code"
	@echo "  lint          - Lint code"
	@echo "  check         - Run all code quality checks"
	@echo "  deps          - Install dependencies"
	@echo "  tools         - Install development tools"
	@echo "  docker-build  - Build Docker image"
	@echo "  docker-run    - Build and run Docker container"
	@echo "  help          - Show this help"