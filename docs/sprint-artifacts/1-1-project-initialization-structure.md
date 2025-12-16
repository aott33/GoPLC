# Story 1.1: Project Initialization & Structure

**Status:** ready-for-dev
**Epic:** 1 - Project Foundation & Core Runtime
**Estimated Effort:** Small (1-2 hours)
**Priority:** Critical Path - Must complete before all other stories

### Prerequisites

| Requirement | Version | Verify With |
|-------------|---------|-------------|
| Go | 1.23+ (using 1.25.0) | `go version` |
| Git | Any recent | `git --version` |
| Make | GNU Make | `make --version` |
| Bash shell | Via VSCode terminal | Available in VSCode |

---

## Overview

This story establishes the foundational Go project structure for go-plc. You'll create the directory layout, initialize the Go module, set up build tooling, and create a minimal entry point. This is a straightforward "scaffolding" story - no complex logic, just getting the house in order.

**Why This Matters:**
Every future story depends on this structure. The architecture document specifies exact directory names and organization patterns that all AI agents and human developers must follow consistently.

---

## User Story

> **As a** developer,
> **I want** the go-plc project initialized with the correct directory structure, Go module, and build tooling,
> **So that** I have a solid foundation following the architecture document to build upon.

---

## What You're Building

At the end of this story, the repository will have:

```
go-plc/
├── cmd/
│   └── go-plc/
│       └── main.go              <- You create this (entry point)
├── internal/                    <- You create these directories
│   ├── config/                  <- Will hold YAML parsing (Story 1.2)
│   ├── variables/               <- Will hold variable store (Story 1.3)
│   ├── runtime/                 <- Will hold scheduler (Story 1.7)
│   ├── modbus/                  <- Will hold Modbus client (Epic 2)
│   ├── opcua/                   <- Will hold OPC UA server (Epic 5)
│   ├── sparkplug/               <- Will hold MQTT client (future)
│   ├── api/                     <- Will hold GraphQL (Epic 3)
│   └── tasks/                   <- Will hold task executor (Story 1.5)
├── web/                         <- Frontend placeholder (Epic 4)
├── tasks/                       <- User-written task files
├── scripts/                     <- Build utilities
├── go.mod                       <- You create this
├── Makefile                     <- You create this
└── (existing files: README.md, LICENSE, docs/, etc.)
```

---

## Acceptance Criteria

### AC1: Directory Structure Created

**What to verify:**
- All `internal/` subdirectories exist and contain at least a `doc.go` file
- `web/`, `tasks/`, `scripts/` directories exist (can use `.gitkeep` placeholders)
- No extra directories are created

### AC2: Go Module Initialized

**What to verify:**
- Running `go mod init` creates `go.mod` with appropriate module path
- `go build ./...` compiles without errors

### AC3: Working Build System

**What to verify:**
- Running `make build` creates an executable at `./go-plc`
- Running `make test` executes (even if no tests exist yet)
- Running `make clean` removes build artifacts

### AC4: Minimal Working Entry Point

**What to verify:**
- `cmd/go-plc/main.go` exists with a valid `main()` function
- The binary runs and exits cleanly (doesn't crash)
- Uses `log/slog` for logging (not fmt.Println or other loggers)

### AC5: Clean Git Setup

**What to verify:**
- `.gitignore` updated with Go-specific patterns
- No binaries or build artifacts committed

---

## Step-by-Step Implementation Guide

### Step 1: Create the Directory Structure

**Time estimate: 2 minutes**

Create all required directories with a single command:

```bash
# From the project root (go-plc/) - creates everything in one shot
mkdir -p cmd/go-plc internal/{config,variables,runtime,modbus,opcua,sparkplug,api,tasks} web tasks scripts
```

**Why these specific directories?**
The `internal/` directory is a Go convention - code here can only be imported by code within the same module. This prevents external packages from depending on implementation details we might change later.

---

### Step 2: Initialize the Go Module

**Time estimate: 2 minutes**

```bash
go mod init github.com/yourusername/go-plc
```

**Replace `yourusername`** with your actual GitHub username (e.g., `andrewott`).

This creates `go.mod` which tells Go:
- What this module is called
- What Go version we're using (will be set to 1.25.0)
- What dependencies we need (none yet)

After running `go mod init`, verify the go.mod file specifies the correct Go version:
```
module github.com/yourusername/go-plc

go 1.25.0
```

---

### Step 3: Create Package Placeholder Files

**Time estimate: 10 minutes**

Each `internal/` package needs at least one `.go` file. We'll create `doc.go` files with package documentation:

**internal/config/doc.go:**
```go
// Package config handles YAML configuration parsing and validation.
// It provides the configuration structs and loading logic for sources,
// variables, and runtime settings.
package config
```

**internal/variables/doc.go:**
```go
// Package variables implements the thread-safe variable store.
// It provides storage, scaling, and change notification for all
// process variables used by the PLC runtime.
package variables
```

**internal/runtime/doc.go:**
```go
// Package runtime implements the PLC runtime coordinator.
// It manages the lifecycle, task scheduling, and graceful shutdown
// of the soft PLC system.
package runtime
```

**internal/modbus/doc.go:**
```go
// Package modbus implements the Modbus TCP client.
// It handles connections to industrial Modbus devices, polling,
// and automatic reconnection on failure.
package modbus
```

**internal/opcua/doc.go:**
```go
// Package opcua implements the OPC UA server.
// It exposes PLC variables to SCADA systems via the OPC UA protocol.
package opcua
```

**internal/sparkplug/doc.go:**
```go
// Package sparkplug implements Sparkplug B MQTT messaging.
// It enables SCADA integration via MQTT with the Sparkplug B specification.
package sparkplug
```

**internal/api/doc.go:**
```go
// Package api implements the GraphQL API server.
// It provides queries, mutations, and real-time subscriptions for
// external applications and the web UI.
package api
```

**internal/tasks/doc.go:**
```go
// Package tasks handles task discovery and execution.
// It auto-discovers user-written Go tasks and executes them at
// configured scan rates.
package tasks
```

---

### Step 4: Create the Main Entry Point

**Time estimate: 10 minutes**

This is the heart of this story - a minimal but properly structured main.go:

**cmd/go-plc/main.go:**
```go
// go-plc is a soft PLC runtime for industrial automation.
// It provides Modbus, OPC UA, and GraphQL interfaces for monitoring
// and controlling industrial processes.
package main

import (
	"log/slog"
	"os"
)

// Version information - will be set by build flags in the future
var (
	version = "0.1.0-dev"
)

func main() {
	// Set up structured logging using slog (standard library since Go 1.21).
	// We use JSON format for production compatibility with log aggregators.
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))
	slog.SetDefault(logger)

	// Log startup
	slog.Info("go-plc starting",
		"version", version,
	)

	// ============================================================
	// FUTURE INITIALIZATION SEQUENCE (Story 1.7: Runtime Coordinator)
	// ============================================================
	// The full startup sequence will be implemented in later stories:
	//
	// 1. Parse command-line flags (-config path/to/config.yaml)
	// 2. Load and validate YAML configuration (Story 1.2)
	// 3. Initialize the variable store (Story 1.3)
	// 4. Connect to Modbus sources (Epic 2)
	// 5. Start OPC UA server (Epic 5)
	// 6. Discover and register tasks (Story 1.5)
	// 7. Start task scheduler (Story 1.6)
	// 8. Start GraphQL API server (Epic 3)
	// 9. Block until shutdown signal (SIGTERM/SIGINT)
	// 10. Graceful shutdown
	// ============================================================

	slog.Info("go-plc initialized successfully - runtime not yet implemented")
}
```

**Key points about this code:**

1. **Uses `log/slog`** - This is the standard library structured logger (added in Go 1.21, available in our Go 1.25.0). The architecture mandates this - do NOT use `log`, `logrus`, `zap`, or any other logging library.

2. **JSON output** - Production systems need machine-parseable logs. JSON format works with any log aggregation system.

3. **Comments document the future** - The TODO comments serve as documentation for what comes next. Future stories will fill in these pieces.

---

### Step 5: Create the Makefile

**Time estimate: 10 minutes**

The Makefile provides standard commands for building and testing:

**Makefile:**
```makefile
# go-plc Makefile
# ================
# Standard build commands for the go-plc project.
# Run `make help` to see available targets.

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
```

**What each target does:**

| Target | Purpose |
|--------|---------|
| `make build` | Compiles the binary |
| `make test` | Runs tests with race detector |
| `make lint` | Checks code quality |
| `make fmt` | Formats code |
| `make clean` | Removes build artifacts |
| `make help` | Shows available commands |

---

### Step 6: Update .gitignore

**Time estimate: 5 minutes**

Add Go-specific patterns to the existing `.gitignore`:

**Append to .gitignore:**
```gitignore
# Go build artifacts
go-plc
go-plc.exe
*.exe
build/

# Test artifacts
coverage.out
coverage.html
*.coverprofile
*.test

# IDE and editor files (if not already present)
.idea/
*.swp
*.swo
*~
.vscode/settings.json

# OS files
.DS_Store
Thumbs.db

# Dependency directories
vendor/
```

---

### Step 7: Create Placeholder Files with README Breadcrumbs

**Time estimate: 5 minutes**

For directories that will be populated later, create README files that serve as **breadcrumbs** - signposts that help future developers (and AI agents) understand:
- What belongs in this directory
- Which story/epic will populate it
- What NOT to put here

This prevents confusion and misplaced files as the project grows.

**web/README.md:**
```markdown
# WebUI Frontend

This directory will contain the React + Vite + TypeScript frontend.

**Setup:** Epic 4, Story 4.1 (Frontend Project Setup)

**Tech Stack (from architecture):**
- React + Vite + TypeScript
- Tailwind CSS + shadcn/ui
- urql for GraphQL with WebSocket subscriptions
- Roboto font family

**Do not manually create files here** - Story 4.1 will initialize the complete frontend project.
```

**tasks/README.md:**
```markdown
# User Task Files

This directory contains user-written Go task files that are auto-discovered by the runtime.

**Setup:** Story 1.5 (Task Discovery) + Story 1.8 (Example Task)

**Task Interface:**
Tasks must implement the Task interface (defined in Story 1.5):
- `Name() string` - Task identifier
- `Config() TaskConfig` - Scan rate, enabled status
- `Execute(ctx TaskContext) error` - Control logic

**Example:** An example task will be added in Story 1.8.
```

**scripts/README.md:**
```markdown
# Build and Utility Scripts

This directory contains build automation and utility scripts.

**Current contents:** None yet

**Planned additions:**
- Cross-compilation scripts (Story 6.2)
- Deployment helpers (Story 6.3)
```

---

### Step 8: Verify Everything Works

**Time estimate: 5 minutes**

Run these commands to verify your setup:

```bash
# 1. Verify the Go module is set up correctly
go mod tidy

# 2. Verify everything compiles
go build ./...

# 3. Build the binary
make build

# 4. Run the binary (should print startup message and exit)
./go-plc

# 5. Run tests (no tests yet, but should not error)
make test

# 6. Clean up
make clean
```

**Expected output from `./go-plc`:**
```json
{"time":"2025-12-15T...","level":"INFO","msg":"go-plc starting","version":"0.1.0-dev"}
{"time":"2025-12-15T...","level":"INFO","msg":"go-plc initialized successfully - runtime not yet implemented"}
```

---

## Checklist Before Marking Complete

- [x] Go 1.25.0 installed and verified (`go version`)
- [ ] All directories created per the structure above
- [ ] `go.mod` exists with module path and `go 1.25.0`
- [ ] Each `internal/` package has a `doc.go` file
- [ ] `cmd/go-plc/main.go` compiles and runs
- [ ] `make build` produces a working binary
- [ ] `make test` runs without errors
- [ ] `.gitignore` updated with Go patterns
- [ ] README breadcrumbs in `web/`, `tasks/`, `scripts/`
- [ ] No binaries committed to git

---

## Post-Implementation Workflow

After completing development, follow these steps:

### 1. Update Story Status
Change status in `docs/sprint-artifacts/sprint-status.yaml`:
```yaml
1-1-project-initialization-structure: review  # was: ready-for-dev
```

### 2. Run Code Review
```bash
/bmad:bmm:workflows:code-review
```
This runs an **adversarial review** that:
- Validates directory structure matches architecture
- Checks naming conventions (camelCase, PascalCase rules)
- Verifies Makefile functionality
- Confirms slog usage (not fmt.Println)
- **Must find 3-10 issues** (it's designed to be thorough, not rubber-stamp)

### 3. Address Review Feedback
- Fix any issues identified
- Re-run review if needed

### 4. Mark Story Complete
Update status in `docs/sprint-artifacts/sprint-status.yaml`:
```yaml
1-1-project-initialization-structure: done  # was: review
```

### 5. Proceed to Next Story
Run create-story for Story 1.2:
```bash
/bmad:bmm:workflows:create-story
```

---

**Testing Notes:**
- Story 1.1 has no unit tests (just scaffolding)
- `make test` should run without errors (empty test suite is OK)
- TDD begins with Story 1.2 (Configuration Schema)
- Go's built-in `testing` package is our framework - no external dependencies needed

---

## What's NOT in This Story

To keep this story focused, the following are **explicitly excluded**:

| Item | Handled In |
|------|------------|
| `config.yaml` file | Story 1.2 (Configuration Schema) |
| Parsing command-line flags | Story 1.2 |
| YAML configuration loading | Story 1.2 |
| Variable store implementation | Story 1.3 |
| Logging configuration | Story 1.4 |
| Task discovery | Story 1.5 |
| Task execution runtime | Story 1.6 |
| Runtime lifecycle management | Story 1.7 |
| Example task | Story 1.8 |
| Frontend setup | Story 4.1 |
| Docker configuration | Story 6.1 |

> **Note:** Do NOT create `config.yaml` in this story. The configuration schema will be defined in Story 1.2 after the config package is implemented.

---

## Troubleshooting

### "package not found" errors
Make sure each `internal/` directory has at least one `.go` file with a valid package declaration.

### "go.mod not found"
Run `go mod init` from the project root directory.

### golangci-lint not found
This is optional for this story. Install with:
```bash
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
```

### Binary won't run on Windows
The Makefile creates `go-plc` (Linux/Mac). For Windows, use:
```bash
go build -o go-plc.exe ./cmd/go-plc
```

---

## Architecture References

This story implements the project structure defined in:

- [Architecture: Project Structure](docs/architecture.md) - Complete directory specification
- [Architecture: Naming Patterns](docs/architecture.md) - Naming conventions
- [Architecture: Implementation Patterns](docs/architecture.md) - Code patterns

---

## Dev Notes After Completion

- running go mod init create a go.mod file with go version 1.24.5
  - I ran go mod edit -go=1.25 to update the version

### Files Created

_List all files created_

### Deviations from Plan

_Note any changes made during implementation and why_

### Lessons Learned

_Any insights for future stories_
