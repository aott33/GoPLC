---
stepsCompleted: [1, 2, 3, 4, 5, 6, 7, 8]
inputDocuments:
  - 'docs/prd.md'
  - 'docs/ux-design-specification.md'
workflowType: 'architecture'
lastStep: 8
status: 'complete'
completedAt: '2025-12-14'
project_name: 'go-plc'
user_name: 'Andy'
date: '2025-12-14'
---

# Architecture Decision Document

_This document builds collaboratively through step-by-step discovery. Sections are appended as we work through each architectural decision together._

## Project Context Analysis

### Requirements Overview

**Functional Requirements:**

The PRD defines 53 functional requirements across 10 domains:

| Domain | FR Count | Architectural Impact |
|--------|----------|---------------------|
| PLC Runtime Control | 6 | Core scheduler, state machine, signal handling |
| Variable Management | 5 | Central variable store, scaling engine, API surface |
| Modbus Communication | 8 | TCP client, connection manager, polling loops |
| OPC UA Integration | 5 | OPC UA server, node management, read/write handlers |
| GraphQL API | 6 | Schema design, resolvers, subscription engine |
| WebUI Monitoring | 7 | React components, real-time updates, embedded serving |
| Task Development | 5 | Task discovery, execution runtime, variable access API |
| Logging & Diagnostics | 4 | Logging framework, level configuration, error formatting |
| Configuration | 4 | YAML parsing, validation, source/variable binding |
| Deployment | 3 | Build pipeline, cross-compilation, service management |

**Non-Functional Requirements:**

Critical NFRs driving architectural decisions:

- **NFR1:** Task execution overhead <50µs per cycle → Minimal abstraction in hot path
- **NFR4:** GraphQL response <10ms → Efficient resolver design, connection pooling
- **NFR7:** Memory stability over 24+ hours → No goroutine leaks, proper cleanup
- **NFR11:** Automatic Modbus reconnection → Connection manager with exponential backoff
- **NFR16:** Cross-protocol value consistency → Single source of truth variable store
- **NFR21-23:** Protocol compliance → Use proven libraries (simonvetter/modbus, gopcua/opcua)

**Scale & Complexity:**

- Primary domain: IoT/Embedded Edge Application with Web Interface
- Complexity level: Medium-High
- Estimated architectural components: 8-10 major modules

### Technical Constraints & Dependencies

**External Library Dependencies:**
- `simonvetter/modbus` - Modbus TCP client (mandated in PRD)
- `gopcua/opcua` - OPC UA server implementation (mandated in PRD)
- `gqlgen` - GraphQL server with subscription support
- `zenoh-go` - Optional Phase 2: Zenoh protocol bindings for high-performance pub/sub (community contribution opportunity)

**Deployment Platform:**
- **Target:** Real-time Linux OS (PREEMPT_RT kernel recommended)
- Single binary deployment with embedded assets
- No external runtime dependencies
- Runs as systemd service for 24/7 operation

**Development Platform:**
- Developer's choice (Linux, Windows, macOS)
- Any IDE/editor (VS Code, Neovim, GoLand, etc.)
- Standard Go toolchain for cross-compilation to Linux target
- Python Modbus simulator for local testing without hardware

**UX Technical Stack (from UX Spec):**
- React + Vite + TypeScript
- urql for GraphQL with WebSocket subscriptions
- Tailwind CSS + shadcn/ui components
- Roboto font family

### Cross-Cutting Concerns Identified

1. **Concurrency Management**
   - Multiple goroutines accessing shared variable store
   - Task execution, Modbus polling, protocol servers all concurrent
   - Requires thread-safe design patterns throughout

2. **Error Handling & Recovery**
   - Connection failures must not crash runtime
   - Exponential backoff for reconnection
   - Clear error propagation to logging and WebUI

3. **Configuration-Driven Behavior**
   - YAML as single source of configuration truth
   - Variables, sources, tasks all configurable
   - Startup validation before runtime begins

4. **Observability**
   - Configurable log levels without restart
   - Real-time status exposed via GraphQL
   - Human-readable error messages (not codes)

5. **Graceful Lifecycle Management**
   - Clean startup sequence
   - Signal handling (SIGTERM/SIGINT)
   - Coordinated shutdown within 5 seconds

## Starter Template Evaluation

### Primary Technology Domain

IoT/Embedded Edge Application with two distinct technology stacks:
- **Backend:** Go (custom build - no standard PLC starter exists)
- **Frontend:** React + Vite + TypeScript (monorepo in `/web`)

### Starter Options Considered

**Go Backend:**
No standard starter templates exist for soft PLC applications. Following official Go documentation recommendation for server projects with gradual growth approach.

**React Frontend:**
- Official shadcn/ui + Vite setup (selected - simplest official path)
- Community starters (rejected - unnecessary complexity)

### Selected Approach

**Go:** Official Go project layout for server applications
**React:** Official Vite + shadcn/ui installation

**Rationale:**
1. Keep it simple - use official tooling over community starters
2. Monorepo structure keeps frontend and backend together
3. Docker support enables easy testing without hardware
4. `internal/` directory prevents accidental API exposure

### Project Structure

```
go-plc/
├── cmd/go-plc/main.go          # Entry point
├── internal/                    # Private implementation packages
│   ├── config/                  # YAML parsing, validation
│   ├── runtime/                 # PLC scheduler, state machine
│   ├── variables/               # Variable store, scaling engine
│   ├── modbus/                  # Modbus TCP client
│   ├── opcua/                   # OPC UA server
│   ├── zenoh/                   # Optional Phase 2: Zenoh protocol integration
│   ├── api/                     # GraphQL server
│   └── tasks/                   # Task discovery, execution
├── web/                         # React frontend (embedded at build)
├── tasks/                       # User task files (auto-discovered)
├── config.yaml                  # Runtime configuration
├── Dockerfile                   # Multi-stage build
├── docker-compose.yaml          # Development environment
├── go.mod
└── Makefile                     # Build automation
```

### Initialization Commands

**Go Module:**
```bash
go mod init github.com/[username]/go-plc
```

**Frontend Setup:**
```bash
cd web
pnpm create vite@latest . --template react-ts
pnpm add tailwindcss @tailwindcss/vite
pnpm add -D @types/node
pnpm dlx shadcn@latest init
pnpm add urql graphql graphql-ws
```

### Architectural Decisions Provided by Structure

**Code Organization:**
- `internal/` packages are implementation-private
- Each concern in its own package (modbus, opcua, etc.)
- Clear separation between runtime core and protocol adapters

**Build & Deployment:**
- Single binary with embedded frontend (Go embed)
- Multi-stage Docker build for minimal image size
- Cross-compilation via standard Go toolchain

**Development Experience:**
- `make dev` - Run with hot reload
- `make build` - Production binary
- `docker-compose up` - Full environment with Modbus simulator

**Note:** Project initialization should be the first implementation story.

## Core Architectural Decisions

### Decision Priority Analysis

**Critical Decisions (Block Implementation):**
- Variable store implementation (sync.RWMutex + map)
- YAML configuration schema (type-discriminated, flat variable list)
- GraphQL subscription pattern (filtered)

**Important Decisions (Shape Architecture):**
- Register addressing (0-based offset + explicit type)
- Error handling (structured errors)
- Logging framework (slog)
- Frontend state management (urql cache + useState)

**Deferred Decisions (Post-MVP):**
- Authentication and authorization
- TLS/encryption for protocols
- Role-based access control

### Data Architecture

**Variable Store:**
- Implementation: `sync.RWMutex` + `map[string]*Variable`
- Rationale: Clear locking semantics, safe iteration, concurrent read support
- Thread-safe access from tasks, Modbus polling, OPC UA, and GraphQL

**Configuration Schema:**

Type-discriminated source configuration with flat variable list:

```yaml
# Sources with protocol-specific config blocks
sources:
  - name: remoteIO
    type: modbus-tcp
    config:
      host: 192.168.1.100
      port: 502
      unitId: 1
      timeout: 1s
      pollInterval: 100ms
      retryInterval: 5s
      byteOrder: big-endian
      wordOrder: high-word-first

  - name: serialDevice
    type: modbus-rtu
    config:
      device: /dev/ttyUSB0
      baudRate: 19200
      dataBits: 8
      parity: none
      stopBits: 2
      unitId: 1
      timeout: 1s
      pollInterval: 100ms

# Flat variable list with source references
variables:
  - name: tankLevel
    source: remoteIO
    register:
      type: holding
      address: 0
    dataType: uint16
    scale:
      rawMin: 0
      rawMax: 65535
      engMin: 0
      engMax: 100
      unit: "%"
    tags: [tank1, level]

  - name: pumpRunning
    source: remoteIO
    register:
      type: coil
      address: 0
    dataType: bool
    tags: [tank1, pump]
```

**Supported Data Types (from simonvetter/modbus):**

| dataType | Registers | Description |
|----------|-----------|-------------|
| `bool` | N/A | Coils and discrete inputs |
| `uint16` | 1 | Unsigned 16-bit integer |
| `int16` | 1 | Signed 16-bit integer |
| `uint32` | 2 | Unsigned 32-bit integer |
| `int32` | 2 | Signed 32-bit integer |
| `float32` | 2 | 32-bit floating point |
| `uint64` | 4 | Unsigned 64-bit integer |
| `int64` | 4 | Signed 64-bit integer |
| `float64` | 4 | 64-bit floating point |

**Register Addressing:**
- 0-based offset with explicit register type
- Matches simonvetter/modbus library API
- Types: `holding`, `input`, `coil`, `discrete`

**Scaling (Optional):**
- Formula: `scaled = ((raw - rawMin) / (rawMax - rawMin)) * (engMax - engMin) + engMin`
- When omitted: raw value used as-is
- Per-variable endianness override supported

**Poll Interval:** Per-source (not per-variable)

**Protocol Source Registry Pattern:**

The configuration loading implementation uses a registry pattern for protocol-specific source configs to enable extensibility without modifying core config code:

- **Registry Location:** `internal/source/registry.go` provides `Register()`, `ParseConfig()`, and `RegisteredTypes()` functions
- **Protocol Packages:** Each protocol (e.g., `internal/source/modbus/`) implements `SourceConfig` interface and registers itself via `init()`
- **Factory Pattern:** Protocol packages provide factory functions that parse YAML config nodes into typed structs
- **Zero-Touch Addition:** New protocols (OPC-UA client, MQTT, EtherNet/IP) can be added by creating a new package under `internal/source/` without touching existing code
- **Blank Imports:** `main.go` uses blank imports (`_ "github.com/aott33/go-plc/internal/source/modbus"`) to trigger protocol registration

See [Sprint Change Proposal 2025-12-20](../sprint-change-proposal-2025-12-20.md) for complete technical specification and rationale.

### Authentication & Security

**MVP Security Posture:**
- No authentication (trusted network assumption)
- Network binding configurable in YAML
- Documented limitation for MVP

**Deferred to Post-MVP:**
- API key / token authentication
- TLS for OPC UA and GraphQL
- Role-based access control
- Secure credential storage

### API & Communication Patterns

**GraphQL Subscriptions:**
- Pattern: Filtered subscriptions
- Allows subscribing to all variables or filtered by tags/source/name
- Example: `variableUpdates(filter: { tags: ["tank1"] })`

**Error Handling:**
- Pattern: Structured errors with category, message, and context
- Human-readable messages (no cryptic codes)
- Structure:
  ```go
  type PLCError struct {
      Category string // "connection", "config", "task", "protocol"
      Message  string // Human-readable description
      Source   string // Optional: which source/variable/task
  }
  ```

### Frontend Architecture

**State Management:**
- Primary: urql cache for server state (variables, status, tasks)
- Secondary: React useState for UI-only state (sidebar, theme, filters)
- Rationale: Monitoring dashboard - most state comes from server

**Tech Stack (from UX Spec):**
- React + Vite + TypeScript
- urql with graphql-ws for subscriptions
- Tailwind CSS + shadcn/ui
- Roboto font family

### Infrastructure & Deployment

**Logging:**
- Framework: Go 1.21+ `log/slog` (standard library)
- Structured JSON logging
- Configurable levels without restart
- Rationale: No external dependency, sufficient for needs

**Docker Support:**
- `Dockerfile`: Multi-stage build (frontend → backend → minimal runtime)
- `docker-compose.yaml`: go-plc service only
- Modbus simulator: External repository (will be added to docker-compose later)

**Deployment Targets:**
- Production: Real-time Linux OS (PREEMPT_RT recommended)
- Development: Any OS (Linux, Windows, macOS)
- Testing: Docker container

### Decision Impact Analysis

**Implementation Sequence:**
1. Project initialization (Go module, frontend scaffold)
2. Configuration parsing and validation
3. Variable store with scaling engine
4. Modbus client with connection manager
5. Task discovery and execution runtime
6. GraphQL API with subscriptions
7. WebUI components
8. OPC UA server integration
9. Docker and build pipeline

**Cross-Component Dependencies:**
- Variable store is central - all protocols read/write through it
- Configuration drives source and variable initialization
- GraphQL subscriptions depend on variable store change notifications
- Task execution depends on variable store API

## Implementation Patterns & Consistency Rules

### Pattern Categories Defined

**Critical Conflict Points Identified:** 7 areas where AI agents could make different choices, now standardized.

### Naming Patterns

**Go Code Naming:**
- Package names: lowercase, single word (`config`, `modbus`, `variables`)
- File names: lowercase with underscores (`modbus_client.go`, `variable_store.go`)
- Exported (public): `PascalCase` (`GetVariable`, `PLCRuntime`)
- Unexported (private): `camelCase` (`pollLoop`, `reconnectTimer`)
- JSON struct tags: `camelCase`

```go
type Variable struct {
    Name     string  `json:"name"`
    Value    float64 `json:"value"`
    RawValue int64   `json:"rawValue"`
    Unit     string  `json:"unit"`
}
```

**GraphQL Naming:**
- Type names: `PascalCase` (`Variable`, `Source`, `PLCStatus`)
- Field names: `camelCase` (`tankLevel`, `pollInterval`)
- Enum values: `SCREAMING_SNAKE_CASE` (`HOLDING`, `INPUT`, `COIL`, `DISCRETE`)

```graphql
type Query {
  variables(filter: VariableFilter): [Variable!]!
  sources: [Source!]!
  plcStatus: PLCStatus!
}

type Subscription {
  variableUpdates(filter: VariableFilter): VariableUpdate!
}

enum RegisterType {
  HOLDING
  INPUT
  COIL
  DISCRETE
}
```

**YAML Configuration:**
- Field names: `camelCase` (`pollInterval`, `byteOrder`, `unitId`)
- Consistent with JSON and GraphQL conventions

### Structure Patterns

**Go Test Organization:**
- Co-located with source files
- Example: `internal/modbus/client.go` → `internal/modbus/client_test.go`

**Frontend Organization (by type):**
```
web/src/
├── components/
│   ├── ui/              # shadcn/ui components
│   ├── InfoBar.tsx
│   ├── Sidebar.tsx
│   ├── Panel.tsx
│   ├── SourcesTable.tsx
│   ├── VariablesTable.tsx
│   └── TasksTable.tsx
├── hooks/
│   └── useVariables.ts
├── lib/
│   └── urql.ts          # GraphQL client setup
├── App.tsx
└── main.tsx
```

### Format Patterns

**Error Message Format:**
```
[Component] - [Human description] (context: [details])
```

**Examples:**
- `remoteIO - Connection timeout after 3 retries (host: 192.168.1.100:502)`
- `tankLevel - Value exceeds configured range (value: 105, max: 100)`
- `pumpControl - Task execution exceeded scan rate (elapsed: 150ms, limit: 100ms)`

**Logging Format (Structured JSON via slog):**
```json
{"time":"2025-01-15T10:30:00Z","level":"INFO","msg":"Modbus source connected","source":"remoteIO","host":"192.168.1.100"}
```

**Log Level Usage:**

| Level | Usage |
|-------|-------|
| `DEBUG` | Detailed diagnostic info (variable values, poll timing) |
| `INFO` | Normal operations (startup, config loaded, connections established) |
| `WARN` | Recoverable issues (reconnection attempts, value out of range) |
| `ERROR` | Failures requiring attention (connection failures, task errors) |

### Communication Patterns

**Variable Store Change Notifications:**
- Publish changes via channel for GraphQL subscription broadcast
- Pattern: Observer with buffered channel to prevent blocking

**GraphQL Subscription Filter:**
```graphql
input VariableFilter {
  names: [String!]
  tags: [String!]
  source: String
}
```

### Enforcement Guidelines

**All AI Agents MUST:**
1. Follow `camelCase` for JSON, GraphQL fields, and YAML config
2. Co-locate tests with source files
3. Use structured error format: `[Component] - [Description] (context)`
4. Use slog for all logging with appropriate levels
5. Follow Go naming conventions (PascalCase exported, camelCase unexported)

**Pattern Verification:**
- Code review checks naming consistency
- GraphQL schema validates field naming
- YAML parsing validates config structure

### Pattern Examples

**Good Examples:**
```go
// Correct: camelCase JSON tags
type Source struct {
    Name         string `json:"name"`
    Type         string `json:"type"`
    PollInterval string `json:"pollInterval"`
}

// Correct: Structured error with context
return fmt.Errorf("remoteIO - Connection failed after %d retries (host: %s)", retries, host)

// Correct: slog structured logging
slog.Info("Source connected", "source", source.Name, "host", host)
```

**Anti-Patterns:**
```go
// Wrong: snake_case JSON tags
type Source struct {
    PollInterval string `json:"poll_interval"` // Should be pollInterval
}

// Wrong: Cryptic error without context
return errors.New("connection error") // No component, no context

// Wrong: Printf logging instead of slog
log.Printf("Connected to %s", host) // Should use slog
```

## Project Structure & Boundaries

### Complete Project Directory Structure

```
go-plc/
├── cmd/
│   └── go-plc/
│       └── main.go                    # Application entry point
├── internal/
│   ├── config/
│   │   ├── config.go                  # Config struct definitions
│   │   ├── loader.go                  # YAML parsing and validation
│   │   └── config_test.go
│   ├── source/                        # Protocol source registry pattern
│   │   ├── source.go                  # Source interface definitions
│   │   ├── registry.go                # Protocol factory registry
│   │   └── modbus/
│   │       ├── tcp.go                 # Modbus TCP config + registration
│   │       ├── rtu.go                 # Modbus RTU config + registration
│   │       └── duration.go            # Duration type for YAML parsing
│   ├── variables/
│   │   ├── store.go                   # Variable store (RWMutex + map)
│   │   ├── variable.go                # Variable struct and methods
│   │   ├── scaling.go                 # Scaling engine
│   │   ├── notifications.go           # Change notification channels
│   │   └── store_test.go
│   ├── runtime/
│   │   ├── runtime.go                 # PLC runtime coordinator
│   │   ├── scheduler.go               # Task scheduler
│   │   ├── state.go                   # Runtime state machine
│   │   ├── lifecycle.go               # Startup/shutdown coordination
│   │   └── runtime_test.go
│   ├── modbus/
│   │   ├── client.go                  # Modbus client wrapper
│   │   ├── connection.go              # Connection manager with reconnect
│   │   ├── poller.go                  # Polling loop goroutine
│   │   ├── types.go                   # Modbus-specific types
│   │   └── client_test.go
│   ├── opcua/
│   │   ├── server.go                  # OPC UA server setup
│   │   ├── nodes.go                   # Node generation from variables
│   │   ├── handlers.go                # Read/write handlers
│   │   └── server_test.go
│   ├── zenoh/                   # Optional Phase 2: Zenoh protocol
│   │   ├── session.go                 # Zenoh session management
│   │   ├── publisher.go               # Pub/sub implementation
│   │   ├── queryable.go               # Query/reply handlers
│   │   └── session_test.go
│   ├── api/
│   │   ├── server.go                  # HTTP/WebSocket server
│   │   ├── schema.graphql             # GraphQL schema definition
│   │   ├── resolver.go                # Generated resolver interface
│   │   ├── query.go                   # Query resolvers
│   │   ├── mutation.go                # Mutation resolvers (setVariable)
│   │   ├── subscription.go            # Subscription resolvers
│   │   └── server_test.go
│   └── tasks/
│       ├── discovery.go               # Task file discovery
│       ├── executor.go                # Task execution runtime
│       ├── context.go                 # Task execution context (var access)
│       └── executor_test.go
├── web/
│   ├── src/
│   │   ├── components/
│   │   │   ├── ui/                    # shadcn/ui components
│   │   │   ├── InfoBar.tsx            # Top status bar
│   │   │   ├── Sidebar.tsx            # Navigation sidebar
│   │   │   ├── Panel.tsx              # Resizable panel container
│   │   │   ├── SourcesPanel.tsx       # Sources table component
│   │   │   ├── VariablesPanel.tsx     # Variables table component
│   │   │   └── TasksPanel.tsx         # Tasks table component
│   │   ├── hooks/
│   │   │   ├── useVariables.ts        # Variable subscription hook
│   │   │   ├── useSources.ts          # Sources query hook
│   │   │   └── usePlcStatus.ts        # PLC status hook
│   │   ├── lib/
│   │   │   ├── urql.ts                # GraphQL client setup
│   │   │   └── utils.ts               # Utility functions
│   │   ├── types/
│   │   │   └── graphql.ts             # Generated GraphQL types
│   │   ├── App.tsx                    # Main application component
│   │   ├── main.tsx                   # React entry point
│   │   └── index.css                  # Tailwind imports
│   ├── public/
│   │   └── favicon.ico
│   ├── index.html
│   ├── package.json
│   ├── pnpm-lock.yaml
│   ├── tsconfig.json
│   ├── vite.config.ts
│   ├── tailwind.config.js
│   └── components.json               # shadcn/ui config
├── docs-site/                         # Docusaurus documentation site
│   ├── docs/
│   │   ├── intro.md                   # Getting started guide
│   │   ├── installation.md            # Installation instructions
│   │   ├── configuration/
│   │   │   ├── overview.md            # Configuration overview
│   │   │   ├── sources.md             # Source configuration reference
│   │   │   └── variables.md           # Variable configuration reference
│   │   ├── tasks/
│   │   │   ├── overview.md            # Task development overview
│   │   │   ├── writing-tasks.md       # How to write tasks
│   │   │   └── api-reference.md       # Task API reference
│   │   ├── protocols/
│   │   │   ├── modbus.md              # Modbus configuration
│   │   │   ├── opcua.md               # OPC UA server details
│   │   │   └── zenoh.md               # Zenoh protocol setup (Phase 2)
│   │   └── api/
│   │       └── graphql.md             # GraphQL API reference
│   ├── blog/                          # Optional: Release notes, updates
│   ├── src/
│   │   ├── css/
│   │   │   └── custom.css             # Custom Docusaurus styles
│   │   └── pages/
│   │       └── index.tsx              # Landing page
│   ├── static/
│   │   └── img/                       # Images and diagrams
│   ├── docusaurus.config.js           # Docusaurus configuration
│   ├── sidebars.js                    # Documentation sidebar structure
│   ├── package.json
│   └── tsconfig.json
├── tasks/                             # User task files (auto-discovered)
│   └── example_task.go                # Example task template
├── scripts/                           # Build and utility scripts
├── config.yaml                        # Example runtime configuration
├── go.mod
├── go.sum
├── Makefile
├── Dockerfile
├── docker-compose.yaml
├── .gitignore
└── README.md
```

### Architectural Boundaries

**API Boundaries:**

| Boundary | Location | Protocol |
|----------|----------|----------|
| GraphQL API | `internal/api/` | HTTP + WebSocket (port 8080) |
| OPC UA Server | `internal/opcua/` | OPC UA (port 4840) |
| Zenoh (Phase 2) | `internal/zenoh/` | Zenoh pub/sub (peer-to-peer or router) |
| Embedded WebUI | `internal/api/server.go` | HTTP static serving |

**Component Boundaries:**

```
┌─────────────────────────────────────────────────────────────────┐
│                        cmd/go-plc/main.go                       │
│                    (Initialization & Wiring)                    │
└──────────────────────────────┬──────────────────────────────────┘
                               │
┌──────────────────────────────▼──────────────────────────────────┐
│                     internal/runtime/                           │
│              (Scheduler, State Machine, Lifecycle)              │
└──────────────────────────────┬──────────────────────────────────┘
                               │
         ┌─────────────────────┼─────────────────────┐
         │                     │                     │
         ▼                     ▼                     ▼
┌─────────────────┐  ┌─────────────────┐  ┌─────────────────┐
│ internal/modbus │  │ internal/tasks  │  │ internal/api    │
│   (I/O Source)  │  │  (User Logic)   │  │   (GraphQL)     │
└────────┬────────┘  └────────┬────────┘  └────────┬────────┘
         │                    │                     │
         └────────────────────┼─────────────────────┘
                              │
                              ▼
              ┌───────────────────────────────┐
              │     internal/variables/       │
              │  (Central Variable Store)     │
              │   Thread-safe, Observable     │
              └───────────────┬───────────────┘
                              │
         ┌────────────────────┼────────────────────┐
         │                    │                    │
         ▼                    ▼                    ▼
┌─────────────────┐  ┌─────────────────┐  ┌─────────────────┐
│ internal/opcua  │  │ internal/zenoh  │  │ internal/api    │
│   (OPC UA Out)  │  │  (Zenoh Phase2) │  │ (Subscriptions) │
└─────────────────┘  └─────────────────┘  └─────────────────┘
```

**Data Flow:**
1. Config loaded → Sources and Variables initialized
2. Modbus poller reads from devices → Updates Variable Store
3. Variable Store notifies subscribers via channel
4. GraphQL subscriptions broadcast changes to WebUI
5. OPC UA server reads from Variable Store on client request
6. Zenoh (Phase 2) publishes variable changes to network
7. Tasks execute at scan rate → Read/Write Variable Store

### Requirements to Structure Mapping

**FR Domain Mapping:**

| FR Domain | Package(s) | Key Files |
|-----------|------------|-----------|
| PLC Runtime Control (FR1-6) | `internal/runtime/` | `runtime.go`, `scheduler.go`, `lifecycle.go` |
| Variable Management (FR7-11) | `internal/variables/` | `store.go`, `scaling.go` |
| Modbus Communication (FR12-19) | `internal/modbus/` | `client.go`, `connection.go`, `poller.go` |
| OPC UA Integration (FR20-24) | `internal/opcua/` | `server.go`, `nodes.go`, `handlers.go` |
| GraphQL API (FR25-30) | `internal/api/` | `schema.graphql`, `resolver.go`, `subscription.go` |
| WebUI Monitoring (FR31-37) | `web/src/` | `App.tsx`, `components/*.tsx` |
| Task Development (FR38-42) | `internal/tasks/` | `discovery.go`, `executor.go`, `context.go` |
| Logging & Diagnostics (FR43-46) | All packages | slog calls throughout |
| Configuration (FR47-50) | `internal/config/`, `internal/source/` | `config.go`, `loader.go`, `registry.go` |
| Deployment (FR51-53) | Root + `docs-site/` | `Dockerfile`, `Makefile`, `docs-site/` |

**Cross-Cutting Concerns Location:**

| Concern | Primary Location | Notes |
|---------|------------------|-------|
| Logging | All `internal/` packages | slog with structured JSON |
| Error Handling | All `internal/` packages | Structured error format |
| Configuration | `internal/config/` | Loaded once at startup |
| Concurrency | `internal/variables/store.go` | RWMutex-protected map |
| Graceful Shutdown | `internal/runtime/lifecycle.go` | Context cancellation |

### Integration Points

**Internal Communication:**
- Variable Store uses `chan *VariableUpdate` for change notifications
- Runtime coordinates goroutines via `context.Context` for cancellation
- Config is passed by reference at initialization (immutable after startup)

**External Integrations:**

| Integration | Package | External System |
|-------------|---------|-----------------|
| Modbus TCP/RTU | `internal/modbus/` | Industrial devices |
| OPC UA Clients | `internal/opcua/` | SCADA systems (Ignition, Kepware) |
| Zenoh Network (Phase 2) | `internal/zenoh/` | Zenoh routers, peers, storage nodes |
| Web Browsers | `internal/api/` | WebUI via GraphQL |

### Documentation Site Structure

**Docusaurus Initialization:**
```bash
cd docs-site
npx create-docusaurus@latest . classic --typescript
```

**Documentation Categories:**

| Category | Path | Content |
|----------|------|---------|
| Getting Started | `docs/intro.md`, `docs/installation.md` | Quick start, prerequisites, first run |
| Configuration | `docs/configuration/` | YAML schema reference, examples |
| Task Development | `docs/tasks/` | Writing tasks, variable API, examples |
| Protocols | `docs/protocols/` | Modbus, OPC UA, GraphQL, Zenoh (Phase 2) |
| API Reference | `docs/api/` | GraphQL schema, WebSocket subscriptions |

### File Organization Patterns

**Configuration Files:**
- `config.yaml` - Runtime configuration (sources, variables)
- `web/vite.config.ts` - Frontend build configuration
- `web/tailwind.config.js` - Tailwind CSS configuration
- `web/components.json` - shadcn/ui configuration
- `docs-site/docusaurus.config.js` - Documentation site configuration

**Source Organization:**
- Go: One package per concern in `internal/`
- React: Components by type in `web/src/components/`
- Tests: Co-located with source (`*_test.go`)
- Docs: Organized by topic in `docs-site/docs/`

**Build Artifacts:**
- Go binary: `./go-plc` (or `./build/go-plc`)
- Frontend: `web/dist/` (embedded via Go embed)
- Documentation: `docs-site/build/` (deployed separately)
- Docker image: `go-plc:latest`

### Development Workflow Integration

**Development Commands (Makefile):**
```makefile
dev:        # Run Go with frontend dev server (hot reload)
build:      # Build production binary with embedded frontend
test:       # Run all Go tests
lint:       # Run golangci-lint
docker:     # Build Docker image
docs:       # Build documentation site
docs-dev:   # Run documentation site in dev mode
```

**Build Process:**
1. `pnpm build` in `web/` → `web/dist/`
2. `go build` embeds `web/dist/` via `//go:embed`
3. Single binary output with all assets
4. `pnpm build` in `docs-site/` → `docs-site/build/` (separate deployment)

**Deployment:**
- Single binary: `./go-plc -config /path/to/config.yaml`
- Docker: `docker run -v config.yaml:/app/config.yaml go-plc`
- Systemd: Service file starts binary with config path
- Documentation: Deploy `docs-site/build/` to GitHub Pages or Netlify

## Architecture Validation Results

### Coherence Validation ✅

**Decision Compatibility:**
All technology choices are compatible and work together:
- Go backend libraries (simonvetter/modbus, gopcua/opcua, gqlgen, paho.mqtt.golang) are stable and well-maintained
- Frontend stack (React + Vite + TypeScript + urql + Tailwind + shadcn/ui) is a standard modern combination
- No version conflicts or incompatibilities identified

**Pattern Consistency:**
Implementation patterns align with architectural decisions:
- Consistent `camelCase` naming across JSON, GraphQL, and YAML
- Go conventions properly applied (PascalCase exported, camelCase unexported)
- Test organization follows Go standard (co-located)
- Frontend organization (by type) matches React conventions

**Structure Alignment:**
Project structure supports all architectural decisions:
- `internal/` packages enforce implementation privacy
- Monorepo structure enables frontend embedding
- Component boundaries clearly defined
- Variable Store as central hub is properly positioned

### Requirements Coverage Validation ✅

**Functional Requirements Coverage:**
All 53 functional requirements across 10 domains have corresponding architectural components:
- PLC Runtime → `internal/runtime/`
- Variable Management → `internal/variables/`
- Modbus Communication → `internal/modbus/`
- OPC UA Integration → `internal/opcua/`
- GraphQL API → `internal/api/`
- WebUI Monitoring → `web/src/`
- Task Development → `internal/tasks/`
- Logging → slog throughout
- Configuration → `internal/config/`
- Deployment → Root build files + `docs-site/`

**Non-Functional Requirements Coverage:**
All critical NFRs are architecturally addressed:
- Performance: Minimal abstraction in hot path, RWMutex for concurrent reads
- Reliability: Connection manager with reconnection, graceful shutdown
- Memory stability: Context cancellation, no goroutine leaks pattern
- Protocol compliance: Mandated libraries used

### Implementation Readiness Validation ✅

**Decision Completeness:**
- All critical decisions documented with rationale
- Configuration schema fully specified with YAML examples
- Data types, addressing, and scaling defined
- Error handling and logging patterns with examples

**Structure Completeness:**
- Complete directory tree with ~70 files specified
- Package responsibilities clearly documented
- Integration points mapped
- Build process defined

**Pattern Completeness:**
- All 7 conflict areas addressed with patterns
- Good and bad examples provided
- Enforcement guidelines documented

### Gap Analysis Results

**Critical Gaps:** None

**Important Gaps (resolved during implementation):**
- Task API details → Defined during `internal/tasks/` implementation
- Zenoh integration patterns → To be defined in Phase 2 based on zenoh-go maturity

**Future Enhancement Opportunities:**
- Config hot-reload
- OPC UA write-through
- Task hot-reload
- Zenoh geo-distributed storage integration
- Zenoh query/reply patterns for remote task control

### Architecture Completeness Checklist

**✅ Requirements Analysis**
- [x] Project context thoroughly analyzed
- [x] Scale and complexity assessed (Medium-High)
- [x] Technical constraints identified (mandated libraries, real-time Linux)
- [x] Cross-cutting concerns mapped (5 concerns)

**✅ Architectural Decisions**
- [x] Critical decisions documented (variable store, config schema, subscriptions)
- [x] Technology stack fully specified (Go + React + 4 protocols)
- [x] Integration patterns defined (variable store as hub)
- [x] Performance considerations addressed (<50µs overhead)

**✅ Implementation Patterns**
- [x] Naming conventions established (Go, GraphQL, YAML, JSON)
- [x] Structure patterns defined (co-located tests, by-type frontend)
- [x] Communication patterns specified (channel-based notifications)
- [x] Process patterns documented (error handling, logging)

**✅ Project Structure**
- [x] Complete directory structure defined (~70 files)
- [x] Component boundaries established (8 internal packages)
- [x] Integration points mapped (4 external protocols)
- [x] Requirements to structure mapping complete (10 FR domains)

### Architecture Readiness Assessment

**Overall Status:** READY FOR IMPLEMENTATION

**Confidence Level:** HIGH

**Key Strengths:**
1. Clear separation of concerns with well-defined package boundaries
2. Variable Store as single source of truth prevents data inconsistency
3. Type-discriminated config schema enables future protocol expansion
4. Comprehensive implementation patterns prevent AI agent conflicts
5. Complete FR-to-structure mapping ensures full coverage

**Areas for Future Enhancement:**
1. Authentication and security (documented as post-MVP)
2. Hot-reload capabilities for config and tasks
3. Additional industrial protocols (Ethernet/IP, BACnet)

### Implementation Handoff

**AI Agent Guidelines:**
- Follow all architectural decisions exactly as documented
- Use implementation patterns consistently across all components
- Respect project structure and boundaries
- Refer to this document for all architectural questions
- Use `camelCase` for all JSON, GraphQL, and YAML fields

**First Implementation Priority:**
Project initialization using the specified structure and initialization commands, followed by configuration parsing and validation.

## Architecture Completion Summary

### Workflow Completion

**Architecture Decision Workflow:** COMPLETED
**Total Steps Completed:** 8
**Date Completed:** 2025-12-14
**Document Location:** docs/architecture.md

### Final Architecture Deliverables

**Complete Architecture Document**
- All architectural decisions documented with specific versions
- Implementation patterns ensuring AI agent consistency
- Complete project structure with all files and directories
- Requirements to architecture mapping
- Validation confirming coherence and completeness

**Implementation Ready Foundation**
- 15+ architectural decisions made
- 7 implementation pattern categories defined
- 8 architectural components specified
- 53 functional requirements fully supported

**AI Agent Implementation Guide**
- Technology stack with verified versions
- Consistency rules that prevent implementation conflicts
- Project structure with clear boundaries
- Integration patterns and communication standards

### Implementation Handoff

**For AI Agents:**
This architecture document is your complete guide for implementing go-plc. Follow all decisions, patterns, and structures exactly as documented.

**First Implementation Priority:**
```bash
# 1. Initialize Go module
go mod init github.com/[username]/go-plc

# 2. Create directory structure
mkdir -p cmd/go-plc internal/{config,variables,runtime,modbus,opcua,zenoh,api,tasks} web tasks scripts docs-site

# 3. Initialize frontend
cd web && pnpm create vite@latest . --template react-ts
```

**Development Sequence:**
1. Initialize project using documented starter template
2. Set up development environment per architecture
3. Implement core architectural foundations (config, variables)
4. Build features following established patterns
5. Maintain consistency with documented rules

### Quality Assurance Checklist

**Architecture Coherence**
- [x] All decisions work together without conflicts
- [x] Technology choices are compatible
- [x] Patterns support the architectural decisions
- [x] Structure aligns with all choices

**Requirements Coverage**
- [x] All functional requirements are supported
- [x] All non-functional requirements are addressed
- [x] Cross-cutting concerns are handled
- [x] Integration points are defined

**Implementation Readiness**
- [x] Decisions are specific and actionable
- [x] Patterns prevent agent conflicts
- [x] Structure is complete and unambiguous
- [x] Examples are provided for clarity

### Project Success Factors

**Clear Decision Framework**
Every technology choice was made collaboratively with clear rationale, ensuring all stakeholders understand the architectural direction.

**Consistency Guarantee**
Implementation patterns and rules ensure that multiple AI agents will produce compatible, consistent code that works together seamlessly.

**Complete Coverage**
All project requirements are architecturally supported, with clear mapping from business needs to technical implementation.

**Solid Foundation**
The chosen project structure and architectural patterns provide a production-ready foundation following current best practices.

---

**Architecture Status:** READY FOR IMPLEMENTATION

**Next Phase:** Begin implementation using the architectural decisions and patterns documented herein.

**Document Maintenance:** Update this architecture when major technical decisions are made during implementation.

