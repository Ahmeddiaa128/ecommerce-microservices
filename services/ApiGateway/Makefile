.PHONY: run build test clean docker-up docker-down help

# Variables
SERVICE_NAME=api-gateway
BINARY_NAME=api-gateway
DOCKER_IMAGE=ecommerce/$(SERVICE_NAME)
MAIN_PATH=./cmd/main.go

## help: Display this help message
help:
	@echo "Available targets:"
	@echo "  run          - Run the API Gateway locally"
	@echo "  build        - Build the API Gateway binary"
	@echo "  test         - Run tests"
	@echo "  clean        - Clean build artifacts"
	@echo "  docker-build - Build Docker image"
	@echo "  docker-up    - Start service with Docker Compose"
	@echo "  docker-down  - Stop service with Docker Compose"
	@echo "  lint         - Run linters"
	@echo "  fmt          - Format code"

## run: Run the API Gateway locally
run:
	@echo "Running $(SERVICE_NAME)..."
	go run $(MAIN_PATH)

## build: Build the API Gateway binary
build:
	@echo "Building $(SERVICE_NAME)..."
	go build -o bin/$(BINARY_NAME) $(MAIN_PATH)
	@echo "Binary created at bin/$(BINARY_NAME)"

## test: Run tests
test:
	@echo "Running tests..."
	go test -v -race -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

## clean: Clean build artifacts
clean:
	@echo "Cleaning..."
	rm -rf bin/
	rm -rf tmp/
	rm -f coverage.out coverage.html
	@echo "Clean complete"

## docker-build: Build Docker image
docker-build:
	@echo "Building Docker image..."
	cd ../../.. && docker build -t $(DOCKER_IMAGE):latest -f services/ApiGateway/docker/Dockerfile .

## docker-up: Start service with Docker Compose
docker-up:
	@echo "Starting $(SERVICE_NAME) with Docker Compose..."
	cd docker && docker-compose up -d
	@echo "Service started"

## docker-down: Stop service with Docker Compose
docker-down:
	@echo "Stopping $(SERVICE_NAME)..."
	cd docker && docker-compose down
	@echo "Service stopped"

## docker-logs: View Docker logs
docker-logs:
	cd docker && docker-compose logs -f

## lint: Run linters
lint:
	@echo "Running linters..."
	golangci-lint run ./...

## fmt: Format code
fmt:
	@echo "Formatting code..."
	go fmt ./...
	goimports -w .

## deps: Download dependencies
deps:
	@echo "Downloading dependencies..."
	go mod download
	go mod verify

## tidy: Tidy dependencies
tidy:
	@echo "Tidying dependencies..."
	go mod tidy

## dev: Run with hot reload (requires air)
dev:
	@echo "Starting development mode with hot reload..."
	air

.DEFAULT_GOAL := help
