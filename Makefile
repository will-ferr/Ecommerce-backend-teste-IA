.PHONY: test test-coverage build run clean docker-build docker-run docker-stop

# Go commands
GO := go
GOOS := $(shell go env GOOS)
GOARCH := $(shell go env GOARCH)

# Application info
APP_NAME := smart-choice
VERSION := $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
BUILD_TIME := $(shell date -u +"%Y-%m-%dT%H:%M:%SZ")
LDFLAGS := -ldflags "-X main.appVersion=$(VERSION) -X main.buildTime=$(BUILD_TIME)"

# Directories
BUILD_DIR := build
COVERAGE_DIR := coverage

# Docker
DOCKER_IMAGE := $(APP_NAME):$(VERSION)
DOCKER_REGISTRY := your-registry.com
DOCKER_TAG := latest

# Test
TEST_TIMEOUT := 30s
TEST_VERBOSE := -v

.PHONY: help
help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-20s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

.PHONY: deps
deps: ## Download dependencies
	$(GO) mod download
	$(GO) mod tidy

.PHONY: build
build: ## Build the application
	@echo "Building $(APP_NAME) version $(VERSION)..."
	@mkdir -p $(BUILD_DIR)
	$(GO) build $(LDFLAGS) -o $(BUILD_DIR)/$(APP_NAME) .

.PHONY: run
run: ## Run the application
	$(GO) run .

.PHONY: test
test: ## Run tests
	$(GO) test $(TEST_VERBOSE) -timeout $(TEST_TIMEOUT) ./...

.PHONY: test-coverage
test-coverage: ## Run tests with coverage
	@mkdir -p $(COVERAGE_DIR)
	$(GO) test $(TEST_VERBOSE) -timeout $(TEST_TIMEOUT) -coverprofile=$(COVERAGE_DIR)/coverage.out ./...
	$(GO) tool cover -html=$(COVERAGE_DIR)/coverage.out -o $(COVERAGE_DIR)/coverage.html
	@echo "Coverage report generated: $(COVERAGE_DIR)/coverage.html"

.PHONY: test-unit
test-unit: ## Run unit tests only
	$(GO) test $(TEST_VERBOSE) -timeout $(TEST_TIMEOUT) -short ./...

.PHONY: test-integration
test-integration: ## Run integration tests
	$(GO) test $(TEST_VERBOSE) -timeout $(TEST_TIMEOUT) -tags=integration ./...

.PHONY: lint
lint: ## Run linter
	@echo "Running golangci-lint..."
	@if command -v golangci-lint >/dev/null 2>&1; then \
		golangci-lint run; \
	else \
		echo "golangci-lint not installed. Install with: curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b \$$(go env GOPATH)/bin v1.54.2"; \
	fi

.PHONY: fmt
fmt: ## Format code
	$(GO) fmt ./...

.PHONY: vet
vet: ## Run go vet
	$(GO) vet ./...

.PHONY: clean
clean: ## Clean build artifacts
	@echo "Cleaning..."
	@rm -rf $(BUILD_DIR)
	@rm -rf $(COVERAGE_DIR)
	@$(GO) clean -cache -modcache -testcache

.PHONY: docker-build
docker-build: ## Build Docker image
	@echo "Building Docker image $(DOCKER_IMAGE)..."
	docker build -t $(DOCKER_IMAGE) .
	docker tag $(DOCKER_IMAGE) $(DOCKER_REGISTRY)/$(APP_NAME):$(DOCKER_TAG)

.PHONY: docker-run
docker-run: ## Run Docker container
	@echo "Running Docker container..."
	docker-compose up -d

.PHONY: docker-stop
docker-stop: ## Stop Docker containers
	@echo "Stopping Docker containers..."
	docker-compose down

.PHONY: docker-logs
docker-logs: ## Show Docker logs
	docker-compose logs -f

.PHONY: migrate-up
migrate-up: ## Run database migrations
	@echo "Running database migrations..."
	$(GO) run . migrate up

.PHONY: migrate-down
migrate-down: ## Rollback database migrations
	@echo "Rolling back database migrations..."
	$(GO) run . migrate down

.PHONY: seed
seed: ## Seed database with test data
	@echo "Seeding database..."
	$(GO) run . seed

.PHONY: benchmark
benchmark: ## Run benchmarks
	$(GO) test -bench=. -benchmem ./...

.PHONY: profile
profile: ## Run CPU profiling
	$(GO) test -cpuprofile=cpu.prof -bench=. ./...
	$(GO) tool pprof cpu.prof

.PHONY: security-scan
security-scan: ## Run security scan
	@echo "Running security scan..."
	@if command -v gosec >/dev/null 2>&1; then \
		gosec ./...; \
	else \
		echo "gosec not installed. Install with: go install github.com/securecodewarrior/gosec/v2/cmd/gosec@latest"; \
	fi

.PHONY: deps-update
deps-update: ## Update dependencies
	$(GO) get -u ./...
	$(GO) mod tidy

.PHONY: generate
generate: ## Generate code
	$(GO) generate ./...

.PHONY: mod-verify
mod-verify: ## Verify dependencies
	$(GO) mod verify

.PHONY: ci
ci: fmt vet lint test security-scan ## Run all CI checks

.PHONY: pre-commit
pre-commit: fmt vet test ## Run pre-commit checks

.PHONY: release
release: clean test build ## Create release build
	@echo "Release build created: $(BUILD_DIR)/$(APP_NAME)"
