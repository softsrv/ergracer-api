# ErgrAcer API Makefile
# =====================

# Variables
APP_NAME := ergracer-api
CONTAINER_NAME := ergracer-postgres
BIN_DIR := bin
GOOS_LINUX := linux
GOARCH_AMD64 := amd64

# Colors for output
GREEN := \033[32m
YELLOW := \033[33m
BLUE := \033[34m
RESET := \033[0m

.PHONY: help fresh api build build-linux clean test lint deps postgres-start postgres-stop postgres-logs

# Default target
help:
	@echo "$(BLUE)ErgrAcer API - Available Commands$(RESET)"
	@echo "=================================="
	@echo ""
	@echo "$(GREEN)Development:$(RESET)"
	@echo "  fresh        - Remove postgres container, setup fresh DB, and start API with hot reload"
	@echo "  api          - Start API with hot reload (air) without touching postgres"
	@echo ""
	@echo "$(GREEN)Building:$(RESET)"
	@echo "  build        - Build binary for current platform to ./bin/"
	@echo "  build-linux  - Build binary for Linux x86_64 to ./bin/"
	@echo ""
	@echo "$(GREEN)Containers:$(RESET)"
	@echo "  postgres-start - Start postgres container only"
	@echo "  postgres-stop  - Stop postgres container"
	@echo "  postgres-logs  - Show postgres container logs"
	@echo "  api-container  - Start API in container"
	@echo "  api-logs       - Show API container logs"
	@echo ""
	@echo "$(GREEN)Utilities:$(RESET)"
	@echo "  test         - Run all tests"
	@echo "  lint         - Run linter (if golangci-lint is available)"
	@echo "  deps         - Download and tidy dependencies"
	@echo "  clean        - Clean build artifacts and stop containers"

# Fresh setup: remove postgres, setup fresh, start API
fresh:
	@echo "$(YELLOW)🔄 Setting up fresh environment...$(RESET)"
	@echo "$(BLUE)Step 1: Stopping and removing existing postgres container...$(RESET)"
	-docker stop $(CONTAINER_NAME) 2>/dev/null || true
	-docker rm $(CONTAINER_NAME) 2>/dev/null || true
	@echo "$(BLUE)Step 2: Setting up fresh PostgreSQL container...$(RESET)"
	@./setup-container.sh db
	@echo "$(BLUE)Step 3: Starting API with hot reload...$(RESET)"
	@echo "$(GREEN)✅ Fresh environment ready! Starting API...$(RESET)"
	@air

# Start API only with hot reload
api:
	@echo "$(YELLOW)🚀 Starting API with hot reload...$(RESET)"
	@air

# Build for current platform
build:
	@echo "$(YELLOW)🔨 Building $(APP_NAME) for current platform...$(RESET)"
	@mkdir -p $(BIN_DIR)
	@go build -ldflags="-w -s" -o $(BIN_DIR)/$(APP_NAME) .
	@echo "$(GREEN)✅ Build complete: $(BIN_DIR)/$(APP_NAME)$(RESET)"

# Build for Linux x86_64
build-linux:
	@echo "$(YELLOW)🔨 Building $(APP_NAME) for Linux x86_64...$(RESET)"
	@mkdir -p $(BIN_DIR)
	@GOOS=$(GOOS_LINUX) GOARCH=$(GOARCH_AMD64) go build -ldflags="-w -s" -o $(BIN_DIR)/$(APP_NAME)-linux-amd64 .
	@echo "$(GREEN)✅ Linux build complete: $(BIN_DIR)/$(APP_NAME)-linux-amd64$(RESET)"

# Start postgres container only
postgres-start:
	@echo "$(YELLOW)🐘 Starting PostgreSQL container...$(RESET)"
	@./setup-container.sh db

# Stop postgres container
postgres-stop:
	@echo "$(YELLOW)🛑 Stopping PostgreSQL container...$(RESET)"
	@docker stop $(CONTAINER_NAME) 2>/dev/null || echo "Container not running"
	@echo "$(GREEN)✅ PostgreSQL container stopped$(RESET)"

# Show postgres logs
postgres-logs:
	@echo "$(BLUE)📋 PostgreSQL container logs:$(RESET)"
	@docker logs $(CONTAINER_NAME) -f

# Start API container
api-container:
	@echo "$(YELLOW)🚀 Starting API container...$(RESET)"
	@./setup-container.sh api

# Show API container logs
api-logs:
	@echo "$(BLUE)📋 API container logs:$(RESET)"
	@docker logs ergracer-api -f

# Run tests
test:
	@echo "$(YELLOW)🧪 Running tests...$(RESET)"
	@go test -v ./...
	@echo "$(GREEN)✅ Tests complete$(RESET)"

# Run linter if available
lint:
	@echo "$(YELLOW)🔍 Running linter...$(RESET)"
	@if command -v golangci-lint >/dev/null 2>&1; then \
		golangci-lint run; \
		echo "$(GREEN)✅ Linting complete$(RESET)"; \
	else \
		echo "$(YELLOW)⚠️  golangci-lint not found. Install with: go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest$(RESET)"; \
	fi

# Download and tidy dependencies
deps:
	@echo "$(YELLOW)📦 Downloading and tidying dependencies...$(RESET)"
	@go mod download
	@go mod tidy
	@echo "$(GREEN)✅ Dependencies updated$(RESET)"

# Clean build artifacts and stop containers
clean:
	@echo "$(YELLOW)🧹 Cleaning up...$(RESET)"
	@rm -rf $(BIN_DIR)
	@rm -rf tmp
	@-docker stop $(CONTAINER_NAME) 2>/dev/null || true
	@echo "$(GREEN)✅ Cleanup complete$(RESET)"

# Ensure setup script is executable
$(shell chmod +x setup-container.sh 2>/dev/null || true)