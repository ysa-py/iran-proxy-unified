.PHONY: all build clean test install deps docker-build docker-push help

# Project variables
PROJECT_NAME := iran-proxy-ultimate
VERSION := 3.2.0
BUILD_DIR := bin
SRC_DIR := src
DOCKER_REGISTRY := ghcr.io
DOCKER_IMAGE := $(DOCKER_REGISTRY)/$(PROJECT_NAME)

# Build variables
GO := go
GOFLAGS := -v
LDFLAGS := -s -w -X main.Version=$(VERSION) -X main.BuildTime=$(shell date -u '+%Y-%m-%d_%H:%M:%S')
GOBUILD := $(GO) build $(GOFLAGS) -ldflags="$(LDFLAGS)"
GOTEST := $(GO) test -v -race -coverprofile=coverage.out

# Colors for output
CYAN := \033[0;36m
GREEN := \033[0;32m
RED := \033[0;31m
YELLOW := \033[0;33m
NC := \033[0m # No Color

##@ General

help: ## Display this help
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make $(CYAN)<target>$(NC)\n"} /^[a-zA-Z_0-9-]+:.*?##/ { printf "  $(CYAN)%-15s$(NC) %s\n", $$1, $$2 } /^##@/ { printf "\n$(YELLOW)%s$(NC)\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

##@ Development

deps: ## Download Go dependencies
	@echo "$(GREEN)Downloading dependencies...$(NC)"
	@cd $(SRC_DIR) && $(GO) mod download
	@cd $(SRC_DIR) && $(GO) mod verify
	@echo "$(GREEN)✓ Dependencies downloaded$(NC)"

fmt: ## Format Go code
	@echo "$(GREEN)Formatting code...$(NC)"
	@gofmt -s -w $(SRC_DIR)
	@echo "$(GREEN)✓ Code formatted$(NC)"

lint: ## Run linters
	@echo "$(GREEN)Running linters...$(NC)"
	@cd $(SRC_DIR) && $(GO) vet ./...
	@echo "$(GREEN)✓ Linting complete$(NC)"

test: ## Run tests
	@echo "$(GREEN)Running tests...$(NC)"
	@cd $(SRC_DIR) && $(GOTEST) ./...
	@echo "$(GREEN)✓ Tests passed$(NC)"

coverage: test ## Generate coverage report
	@echo "$(GREEN)Generating coverage report...$(NC)"
	@cd $(SRC_DIR) && $(GO) tool cover -html=coverage.out -o coverage.html
	@echo "$(GREEN)✓ Coverage report generated: src/coverage.html$(NC)"

##@ Build

build: deps ## Build the project
	@echo "$(GREEN)Building $(PROJECT_NAME) v$(VERSION)...$(NC)"
	@mkdir -p $(BUILD_DIR)
	@cd $(SRC_DIR) && $(GOBUILD) -o ../$(BUILD_DIR)/$(PROJECT_NAME) main.go main_iran.go
	@echo "$(GREEN)✓ Build complete: $(BUILD_DIR)/$(PROJECT_NAME)$(NC)"

build-all: ## Build for all platforms
	@echo "$(GREEN)Building for all platforms...$(NC)"
	@mkdir -p $(BUILD_DIR)
	@cd $(SRC_DIR) && \
		GOOS=linux GOARCH=amd64 $(GOBUILD) -o ../$(BUILD_DIR)/$(PROJECT_NAME)-linux-amd64 main.go main_iran.go && \
		GOOS=linux GOARCH=arm64 $(GOBUILD) -o ../$(BUILD_DIR)/$(PROJECT_NAME)-linux-arm64 main.go main_iran.go && \
		GOOS=darwin GOARCH=amd64 $(GOBUILD) -o ../$(BUILD_DIR)/$(PROJECT_NAME)-darwin-amd64 main.go main_iran.go && \
		GOOS=darwin GOARCH=arm64 $(GOBUILD) -o ../$(BUILD_DIR)/$(PROJECT_NAME)-darwin-arm64 main.go main_iran.go && \
		GOOS=windows GOARCH=amd64 $(GOBUILD) -o ../$(BUILD_DIR)/$(PROJECT_NAME)-windows-amd64.exe main.go main_iran.go
	@echo "$(GREEN)✓ Cross-platform build complete$(NC)"

install: build ## Install the binary
	@echo "$(GREEN)Installing $(PROJECT_NAME)...$(NC)"
	@install -m 755 $(BUILD_DIR)/$(PROJECT_NAME) /usr/local/bin/
	@echo "$(GREEN)✓ Installed to /usr/local/bin/$(PROJECT_NAME)$(NC)"

clean: ## Clean build artifacts
	@echo "$(YELLOW)Cleaning build artifacts...$(NC)"
	@rm -rf $(BUILD_DIR)
	@rm -f $(SRC_DIR)/coverage.out $(SRC_DIR)/coverage.html
	@echo "$(GREEN)✓ Clean complete$(NC)"

##@ Docker

docker-build: ## Build Docker image
	@echo "$(GREEN)Building Docker image...$(NC)"
	@docker build -t $(DOCKER_IMAGE):$(VERSION) -t $(DOCKER_IMAGE):latest .
	@echo "$(GREEN)✓ Docker image built$(NC)"

docker-push: docker-build ## Push Docker image to registry
	@echo "$(GREEN)Pushing Docker image...$(NC)"
	@docker push $(DOCKER_IMAGE):$(VERSION)
	@docker push $(DOCKER_IMAGE):latest
	@echo "$(GREEN)✓ Docker image pushed$(NC)"

docker-run: ## Run Docker container
	@echo "$(GREEN)Running Docker container...$(NC)"
	@docker run --rm -it $(DOCKER_IMAGE):latest

##@ Run

run: build ## Build and run the application
	@echo "$(GREEN)Running $(PROJECT_NAME)...$(NC)"
	@./$(BUILD_DIR)/$(PROJECT_NAME) --iran-mode --performance-mode balanced

run-iran: build ## Run with Iran-specific optimizations
	@echo "$(GREEN)Running with Iran optimizations...$(NC)"
	@./$(BUILD_DIR)/$(PROJECT_NAME) \
		--iran-mode \
		--dpi-evasion-level aggressive \
		--performance-mode balanced \
		--max-concurrent 200 \
		--timeout 15 \
		--enable-monitoring \
		--enable-self-healing

##@ Utilities

version: ## Show version
	@echo "$(CYAN)$(PROJECT_NAME) version $(VERSION)$(NC)"

check: lint test ## Run all checks
	@echo "$(GREEN)✓ All checks passed$(NC)"

release: clean build-all ## Create a release
	@echo "$(GREEN)Creating release v$(VERSION)...$(NC)"
	@mkdir -p releases/$(VERSION)
	@cp $(BUILD_DIR)/* releases/$(VERSION)/
	@cd releases && tar -czf $(VERSION).tar.gz $(VERSION)
	@echo "$(GREEN)✓ Release created: releases/$(VERSION).tar.gz$(NC)"

.DEFAULT_GOAL := help
