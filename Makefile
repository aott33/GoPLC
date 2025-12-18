# GoPLC Makefile
# --------------------------------------------
#  Build commands for GoPLC
#  Run `make help` to see available targets

# Configuration
BINARY_NAME := go-plc
BUILD_DIR := .
GO := go

# Build flags (can be extended later for version injection)
LDFLAGS :=

# Default target
.DEFAULT_GOAL := help

# ============================================================
# BUILD TARGETS
# ============================================================

.PHONY: build
build: ## Build the go-plc binary
	$(GO) build -ldflags "$(LDFLAGS)" -o $(BUILD_DIR)/$(BINARY_NAME) ./cmd/go-plc
	@echo "Built: $(BUILD_DIR)/$(BINARY_NAME)"

.PHONY: build-all
build-all: ## Build for all target platforms
	@echo "Cross-platform builds not yet configured"
	@echo "TODO: Add GOOS/GOARCH combinations in Story 6.2"

# ============================================================
# TEST TARGETS
# ============================================================

.PHONY: test
test: ## Run all tests with race detection
	$(GO) test -race -v ./...

.PHONY: test-coverage
test-coverage: ## Run tests with coverage report
	$(GO) test -race -coverprofile=coverage.out ./...
	$(GO) tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report: coverage.html"

# ============================================================
# CODE QUALITY
# ============================================================

.PHONY: lint
lint: ## Run linter (requires golangci-lint)
	@which golangci-lint > /dev/null 2>&1 || \
		(echo "golangci-lint not installed. Install with:" && \
		 echo "  go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest" && \
		 exit 1)
	golangci-lint run ./...

.PHONY: fmt
fmt: ## Format all Go code
	$(GO) fmt ./...

.PHONY: vet
vet: ## Run go vet
	$(GO) vet ./...

# ============================================================
# DEVELOPMENT
# ============================================================

.PHONY: run
run: build ## Build and run the application
	./$(BINARY_NAME)

.PHONY: dev
dev: ## Run in development mode (placeholder)
	@echo "Development mode not yet implemented"
	@echo "TODO: Add hot-reload with air or similar tool"

# ============================================================
# CLEANUP
# ============================================================

.PHONY: clean
clean: ## Remove build artifacts
	rm -f $(BUILD_DIR)/$(BINARY_NAME)
	rm -f $(BUILD_DIR)/$(BINARY_NAME).exe
	rm -f coverage.out coverage.html
	@echo "Cleaned build artifacts"

# ============================================================
# HELP
# ============================================================

.PHONY: help
help: ## Show this help message
	@echo "go-plc - Soft PLC Runtime"
	@echo ""
	@echo "Usage: make [target]"
	@echo ""
	@echo "Targets:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | \
		awk 'BEGIN {FS = ":.*?## "}; {printf "  %-15s %s\n", $$1, $$2}'
