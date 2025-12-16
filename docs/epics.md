---
stepsCompleted: [1, 2, 3, 4]
inputDocuments:
  - 'docs/prd.md'
  - 'docs/architecture.md'
  - 'docs/ux-design-specification.md'
project_name: 'go-plc'
user_name: 'Andy'
date: '2025-12-15'
---

# go-plc - Epic Breakdown

## Overview

This document provides the complete epic and story breakdown for go-plc, decomposing the requirements from the PRD, UX Design, and Architecture into implementable stories for human developers.

## Requirements Inventory

### Functional Requirements

**PLC Runtime Control (FR1-FR6):**
- FR1: Operators can enable/disable the PLC runtime (run mode control)
- FR2: System can execute Go-based tasks at configurable scan rates
- FR3: System can auto-discover tasks from the `/tasks` folder
- FR4: Operators can enable/disable individual tasks at runtime
- FR5: System can gracefully shutdown without data corruption
- FR6: System can start with a YAML configuration file

**Variable Management (FR7-FR11):**
- FR7: Developers can define variables in YAML with source binding and scaling
- FR8: System can expose variables to all protocols from a single definition
- FR9: Tasks can read and write variable values through a clean API
- FR10: Operators can view all variables and their current values
- FR11: System can apply linear scaling to variable values (raw to engineering units)

**Modbus Communication (FR12-FR19):**
- FR12: System can connect to Modbus TCP devices as a client
- FR13: System can read holding registers from Modbus devices
- FR14: System can write holding registers to Modbus devices
- FR15: System can read coils from Modbus devices
- FR16: System can write coils to Modbus devices
- FR17: System can automatically reconnect on Modbus connection failure
- FR18: Operators can view Modbus connection status (connected/disconnected/error)
- FR19: Developers can configure multiple Modbus sources with independent polling intervals

**OPC UA Integration (FR20-FR24):**
- FR20: System can expose variables as an OPC UA server
- FR21: SCADA systems can connect to go-plc via OPC UA
- FR22: SCADA systems can read variable values via OPC UA
- FR23: SCADA systems can write variable values via OPC UA
- FR24: Operators can view OPC UA server status

**GraphQL API (FR25-FR30):**
- FR25: External applications can query current variable values via GraphQL
- FR26: External applications can subscribe to variable value changes via GraphQL
- FR27: External applications can query PLC status via GraphQL
- FR28: External applications can query task status via GraphQL
- FR29: External applications can query source/connection status via GraphQL
- FR30: Developers can test GraphQL queries using built-in GraphQL Playground

**WebUI Monitoring (FR31-FR37):**
- FR31: Operators can view PLC runtime status in WebUI
- FR32: Operators can enable/disable PLC runtime from WebUI
- FR33: Operators can view all source/device connection status in WebUI
- FR34: Operators can view variable list with real-time values in WebUI
- FR35: Operators can view task list with configuration and status in WebUI
- FR36: Operators can enable/disable individual tasks from WebUI
- FR37: WebUI can display real-time updates without page refresh

**Task Development (FR38-FR42):**
- FR38: Developers can write control logic in native Go
- FR39: Tasks can access variables through a simple API (no verbose patterns)
- FR40: Developers can specify task scan rate in task configuration
- FR41: Developers can rebuild and deploy tasks in under 1 minute
- FR42: System can report task execution errors with clear messages

**Logging & Diagnostics (FR43-FR46):**
- FR43: System can log events at configurable levels (debug, info, warn, error)
- FR44: Operators can view human-readable error messages (no cryptic codes)
- FR45: System can log connection state changes
- FR46: System can log task execution errors

**Configuration (FR47-FR50):**
- FR47: Developers can define sources (Modbus devices) in YAML
- FR48: Developers can define variables with source bindings in YAML
- FR49: System can validate configuration on startup with clear error messages
- FR50: System can report configuration errors before starting runtime

**Deployment (FR51-FR53):**
- FR51: System can compile to a single binary with embedded WebUI
- FR52: System can run on Linux and Windows platforms
- FR53: System can run as a systemd service on Linux

### Non-Functional Requirements

**Performance (NFR1-NFR9):**
- NFR1: Task execution overhead must be <50µs per cycle
- NFR2: Task scheduler must support scan rates from 10ms to 10s
- NFR3: Variable read/write operations must complete within task cycle budget
- NFR4: GraphQL query response time must be <10ms for variable reads
- NFR5: GraphQL subscriptions must deliver updates within 100ms of value change
- NFR6: OPC UA read operations must complete within 50ms
- NFR7: Memory usage must remain stable during 24+ hour operation (no memory leaks)
- NFR8: CPU usage must remain <50% on target hardware during normal operation
- NFR9: Single binary size must be <50MB (including embedded WebUI)

**Reliability (NFR10-NFR17):**
- NFR10: System must support 24/7 continuous operation
- NFR11: System must recover from Modbus connection failures without operator intervention
- NFR12: System must complete graceful shutdown within 5 seconds
- NFR13: All errors must be logged with human-readable messages
- NFR14: Connection failures must trigger automatic reconnection with exponential backoff
- NFR15: Configuration errors must be reported at startup before runtime begins
- NFR16: Variable values must remain consistent across all protocols (Modbus, OPC UA, GraphQL)
- NFR17: No data corruption on graceful shutdown

**Security (NFR18-NFR20):**
- NFR18: MVP assumes deployment in trusted, firewalled industrial network
- NFR19: No authentication required for MVP (documented limitation)
- NFR20: All network services bind to configurable interfaces (not hardcoded to 0.0.0.0)

**Integration (NFR21-NFR26):**
- NFR21: Modbus TCP implementation must comply with Modbus Application Protocol Specification
- NFR22: OPC UA server must be compatible with standard OPC UA clients (Ignition, Kepware)
- NFR23: GraphQL API must follow GraphQL specification for queries and subscriptions
- NFR24: System must successfully integrate with Ignition SCADA via OPC UA
- NFR25: System must work with Python pymodbus simulator for testing
- NFR26: WebUI must function in modern browsers (Chrome, Firefox, Edge - latest 2 versions)

**Maintainability (NFR27-NFR32):**
- NFR27: Code must follow standard Go formatting (gofmt)
- NFR28: Code must pass go vet with no warnings
- NFR29: Public APIs must have documentation comments
- NFR30: System must compile to single binary for Linux (amd64, arm64) and Windows (amd64)
- NFR31: Configuration changes must not require recompilation
- NFR32: Logs must support configurable output levels without restart

### Additional Requirements

**From Architecture - Starter Template:**
- Project uses official Go project layout with `internal/` packages
- Frontend uses Vite + React + TypeScript + shadcn/ui
- Monorepo structure with `/web` for frontend
- Docker support for development environment

**From Architecture - Technical Decisions:**
- Variable Store: `sync.RWMutex` + `map[string]*Variable` for thread-safe access
- Type-discriminated YAML configuration schema (sources with protocol-specific config blocks)
- GraphQL filtered subscriptions pattern
- Logging via Go 1.21+ `log/slog` (standard library)
- Structured errors: `[Component] - [Description] (context: [details])`
- Channel-based change notifications for GraphQL subscription broadcast

**From Architecture - Implementation Patterns:**
- `camelCase` for JSON, GraphQL fields, YAML config
- `PascalCase` for Go exported identifiers
- Co-located tests (`*_test.go` next to source)
- Frontend organized by type (`components/`, `hooks/`, `lib/`)

**From UX Design - WebUI Requirements:**
- React + Vite + TypeScript with urql for GraphQL
- Tailwind CSS + shadcn/ui components
- Roboto font family
- Dark/light theme with system preference detection
- ISA-101 High Performance HMI alignment

**From UX Design - Component Requirements:**
- InfoBar: Full-width top bar with alerts toggle, PLC status, live indicator
- AlertsPanel: Collapsible panel showing error/warning details
- Sidebar: Collapsible navigation (56px collapsed, 200px expanded)
- Panel: Reusable card container for content sections
- DataTable: Consistent table styling for Sources, Tasks, Variables
- SourceTag: Colored pill showing source name with connection status
- StatusDot: Small colored indicator for status

**From UX Design - Status Colors (ISA-101):**
- OK/Connected: Muted green (#4A7C59)
- Error/Disconnected: Attention red (#C45C5C)
- Warning: Amber (#D4A84B)

### FR Coverage Map

| FR | Epic | Description |
|----|------|-------------|
| FR1 | Epic 1 | Enable/disable PLC runtime |
| FR2 | Epic 1 | Execute Go tasks at scan rates |
| FR3 | Epic 1 | Auto-discover tasks from /tasks |
| FR4 | Epic 1 | Enable/disable individual tasks |
| FR5 | Epic 1 | Graceful shutdown |
| FR6 | Epic 1 | Start with YAML config |
| FR7 | Epic 1 | Define variables in YAML |
| FR8 | Epic 1 | Expose variables to all protocols |
| FR9 | Epic 1 | Task variable access API |
| FR10 | Epic 1 | View variables and values |
| FR11 | Epic 1 | Linear scaling for variables |
| FR12 | Epic 2 | Connect to Modbus TCP |
| FR13 | Epic 2 | Read holding registers |
| FR14 | Epic 2 | Write holding registers |
| FR15 | Epic 2 | Read coils |
| FR16 | Epic 2 | Write coils |
| FR17 | Epic 2 | Auto reconnect on failure |
| FR18 | Epic 2 | View Modbus connection status |
| FR19 | Epic 2 | Multiple Modbus sources |
| FR20 | Epic 5 | OPC UA server for variables |
| FR21 | Epic 5 | SCADA connect via OPC UA |
| FR22 | Epic 5 | SCADA read via OPC UA |
| FR23 | Epic 5 | SCADA write via OPC UA |
| FR24 | Epic 5 | View OPC UA server status |
| FR25 | Epic 3 | Query variables via GraphQL |
| FR26 | Epic 3 | Subscribe to variable changes |
| FR27 | Epic 3 | Query PLC status via GraphQL |
| FR28 | Epic 3 | Query task status via GraphQL |
| FR29 | Epic 3 | Query source status via GraphQL |
| FR30 | Epic 3 | GraphQL Playground |
| FR31 | Epic 4 | View PLC status in WebUI |
| FR32 | Epic 4 | Enable/disable PLC from WebUI |
| FR33 | Epic 4 | View source status in WebUI |
| FR34 | Epic 4 | View variables in WebUI |
| FR35 | Epic 4 | View tasks in WebUI |
| FR36 | Epic 4 | Enable/disable tasks from WebUI |
| FR37 | Epic 4 | Real-time WebUI updates |
| FR38 | Epic 1 | Write control logic in Go |
| FR39 | Epic 1 | Simple variable access API |
| FR40 | Epic 1 | Task scan rate config |
| FR41 | Epic 1 | Rebuild/deploy under 1 minute |
| FR42 | Epic 1 | Task execution error reporting |
| FR43 | Epic 1 | Configurable log levels |
| FR44 | Epic 1 | Human-readable error messages |
| FR45 | Epic 1 | Log connection state changes |
| FR46 | Epic 1 | Log task execution errors |
| FR47 | Epic 1 | Define sources in YAML |
| FR48 | Epic 1 | Define variables in YAML |
| FR49 | Epic 1 | Validate config on startup |
| FR50 | Epic 1 | Report config errors |
| FR51 | Epic 6 | Single binary with embedded WebUI |
| FR52 | Epic 6 | Run on Linux and Windows |
| FR53 | Epic 6 | Run as systemd service |

## Epic List

### Epic 1: Project Foundation & Core Runtime
Developers can initialize the project, configure sources and variables in YAML, and have a running PLC runtime that executes Go tasks with proper logging and error handling.

**FRs covered:** FR1, FR2, FR3, FR4, FR5, FR6, FR7, FR8, FR9, FR10, FR11, FR38, FR39, FR40, FR41, FR42, FR43, FR44, FR45, FR46, FR47, FR48, FR49, FR50

### Epic 2: Modbus I/O Integration
Automation engineers can connect to industrial Modbus TCP devices, read and write holding registers and coils, with reliable automatic reconnection when connections fail.

**FRs covered:** FR12, FR13, FR14, FR15, FR16, FR17, FR18, FR19

### Epic 3: GraphQL API & Real-Time Data
External applications can query system status, variables, and tasks via GraphQL, and subscribe to real-time value changes through WebSocket connections.

**FRs covered:** FR25, FR26, FR27, FR28, FR29, FR30

### Epic 4: WebUI Monitoring Dashboard
Operators can monitor PLC status, view source connections, watch real-time variable values, manage tasks, and troubleshoot issues through a responsive web interface.

**FRs covered:** FR31, FR32, FR33, FR34, FR35, FR36, FR37

### Epic 5: OPC UA SCADA Integration
Traditional SCADA systems (Ignition, Kepware) can connect to go-plc via OPC UA to read and write variable values for industrial integration.

**FRs covered:** FR20, FR21, FR22, FR23, FR24

### Epic 6: Production Deployment & Documentation
Developers can build single-binary deployments with embedded WebUI, run go-plc as a systemd service on Linux, and follow comprehensive Docusaurus documentation for setup and operation.

**FRs covered:** FR51, FR52, FR53

---

## Epic 1: Project Foundation & Core Runtime

Developers can initialize the project, configure sources and variables in YAML, and have a running PLC runtime that executes Go tasks with proper logging and error handling.

### Story 1.1: Project Initialization & Structure

As a **developer**,
I want **the go-plc project initialized with the correct directory structure, Go module, and build tooling**,
So that **I have a solid foundation following the architecture document to build upon**.

**Acceptance Criteria:**

**Given** a clean workspace
**When** I run the initialization commands
**Then** the following structure exists:
- `cmd/go-plc/main.go` with basic entry point
- `internal/` directories: `config/`, `variables/`, `runtime/`, `modbus/`, `opcua/`, `api/`, `tasks/`
- `web/` directory (empty placeholder for frontend)
- `tasks/` directory for user task files
- `go.mod` with module path
- `Makefile` with `build`, `test`, `lint` targets
- `.gitignore` for Go projects

**And** `go build ./...` succeeds with no errors
**And** `make build` produces a binary in the expected location

---

### Story 1.2: Configuration Schema & Loading

As a **developer**,
I want **to define sources and variables in a YAML configuration file that is parsed and validated on startup**,
So that **I can configure the PLC without code changes and catch configuration errors early**.

**Acceptance Criteria:**

**Given** a YAML configuration file with the schema from Architecture (type-discriminated sources, flat variable list)
**When** the application starts with `-config path/to/config.yaml`
**Then** sources are parsed with protocol-specific config blocks (modbus-tcp, modbus-rtu)
**And** variables are parsed with source references, register definitions, and optional scaling

**Given** a YAML file with invalid configuration (missing required field, invalid type, unknown source reference)
**When** the application attempts to load the config
**Then** a clear, human-readable error message is logged (format: `[config] - [description] (context)`)
**And** the application exits with non-zero status before starting the runtime

**Given** a valid configuration
**When** config loading completes
**Then** all sources and variables are accessible via typed Go structs
**And** INFO level log confirms "Configuration loaded successfully"

---

### Story 1.3: Variable Store Implementation

As a **developer**,
I want **a thread-safe variable store that holds all variable values with scaling support**,
So that **multiple goroutines (tasks, protocols, API) can safely read and write variables**.

**Acceptance Criteria:**

**Given** a loaded configuration with variables
**When** the variable store is initialized
**Then** all configured variables are registered in the store with initial values (zero or configured default)

**Given** a variable with scaling configuration (rawMin, rawMax, engMin, engMax)
**When** a raw value is written to the variable
**Then** the scaled (engineering) value is calculated using linear interpolation
**And** both raw and scaled values are stored and accessible

**Given** multiple goroutines reading and writing to the same variable
**When** concurrent access occurs
**Then** no data races occur (verified by `go test -race`)
**And** values remain consistent (RWMutex protects the map)

**Given** a variable store with registered variables
**When** I call `store.Get("variableName")`
**Then** I receive the variable with current value, raw value, unit, and timestamp
**And** unknown variable names return an error

---

### Story 1.4: Logging Framework

As a **developer**,
I want **structured logging using slog with configurable log levels**,
So that **I can debug issues with clear, machine-parseable logs that follow the project patterns**.

**Acceptance Criteria:**

**Given** the application starts
**When** log level is configured in YAML or environment variable
**Then** only messages at or above that level are output
**And** levels are: DEBUG, INFO, WARN, ERROR

**Given** a log message is written
**When** output is produced
**Then** it follows structured JSON format: `{"time":"...","level":"...","msg":"...","key":"value"}`
**And** includes relevant context attributes (source name, variable name, etc.)

**Given** an error condition
**When** an error is logged
**Then** the message follows format: `[Component] - [Description] (context: [details])`
**And** the error includes the category, message, and source context

---

### Story 1.5: Task Discovery & Registration

As a **developer**,
I want **the system to auto-discover Go task files from the /tasks folder and register them**,
So that **I can add new control logic by simply adding task files without modifying core code**.

**Acceptance Criteria:**

**Given** Go files in the `/tasks` directory that implement the Task interface
**When** the application starts
**Then** all valid task files are discovered and registered
**And** INFO log shows each discovered task: `Task discovered: [taskName]`

**Given** a task file with a `TaskConfig` (name, description, scanRate, enabled)
**When** the task is registered
**Then** the configuration is stored for runtime execution
**And** default scan rate is 100ms if not specified

**Given** a task file with syntax errors or missing interface implementation
**When** discovery runs
**Then** an ERROR is logged with the file name and error details
**And** the application can still start with remaining valid tasks (graceful degradation)

**Given** the task discovery completes
**When** I query registered tasks
**Then** I can see task name, description, scan rate, and enabled status for each

---

### Story 1.6: Task Execution Runtime

As a **developer**,
I want **tasks to execute at their configured scan rates with access to the variable store**,
So that **control logic runs deterministically and can read/write process values**.

**Acceptance Criteria:**

**Given** a registered task with scan rate of 100ms
**When** the runtime is started
**Then** the task's Execute() method is called every 100ms (+/- reasonable jitter)
**And** execution timing is logged at DEBUG level

**Given** a task's Execute() method
**When** the task runs
**Then** it receives a TaskContext with methods: `GetVariable(name)`, `SetVariable(name, value)`
**And** variable access uses the thread-safe variable store

**Given** a task execution that exceeds its scan rate (e.g., 150ms execution for 100ms scan)
**When** the overrun is detected
**Then** a WARN log is produced: `[taskName] - Task execution exceeded scan rate (elapsed: 150ms, limit: 100ms)`
**And** the next execution starts immediately (no backlog accumulation)

**Given** a task that throws a panic
**When** the panic occurs
**Then** the panic is recovered, ERROR logged, and the task continues on next cycle
**And** other tasks are not affected

**Given** task execution overhead measurement
**When** measured under normal conditions
**Then** overhead (scheduling, context setup) is <50µs per cycle (NFR1)

---

### Story 1.7: PLC Runtime Coordinator

As a **developer**,
I want **a runtime coordinator that manages the PLC lifecycle with start/stop control and graceful shutdown**,
So that **the system can be operated reliably and shuts down cleanly without data corruption**.

**Acceptance Criteria:**

**Given** configuration loaded, variable store initialized, and tasks discovered
**When** `runtime.Start()` is called
**Then** all enabled tasks begin execution at their scan rates
**And** the runtime state changes to "running"
**And** INFO log: `PLC runtime started`

**Given** a running PLC
**When** `runtime.Stop()` is called
**Then** all task goroutines are signaled to stop via context cancellation
**And** the runtime waits for tasks to complete current cycle (graceful)
**And** the runtime state changes to "stopped"
**And** INFO log: `PLC runtime stopped`

**Given** SIGTERM or SIGINT signal received
**When** signal handler executes
**Then** `runtime.Stop()` is called for graceful shutdown
**And** shutdown completes within 5 seconds (NFR12)
**And** exit code is 0 for clean shutdown

**Given** the runtime is stopped
**When** I query the runtime state
**Then** I can see: state (running/stopped), uptime, task count, variable count

**Given** individual task enable/disable control
**When** `runtime.EnableTask(name)` or `runtime.DisableTask(name)` is called
**Then** the task starts or stops execution accordingly
**And** state change is logged at INFO level

---

### Story 1.8: Example Task & Integration Test

As a **developer**,
I want **an example task demonstrating the task API and an integration test validating the complete runtime**,
So that **I can verify the system works end-to-end and have a reference for writing tasks**.

**Acceptance Criteria:**

**Given** the example task file in `/tasks/example_task.go`
**When** I review the code
**Then** it demonstrates:
- Task interface implementation
- TaskConfig definition (name, description, scanRate)
- GetVariable and SetVariable usage
- Simple control logic (e.g., toggle a bool, increment a counter)

**Given** a test configuration with mock variables
**When** I run `go test ./...`
**Then** integration tests verify:
- Config loads successfully
- Variable store initializes with configured variables
- Tasks are discovered and registered
- Runtime starts and tasks execute
- Variables are updated by task execution
- Graceful shutdown completes

**Given** `make test` is run
**When** all tests pass
**Then** code coverage report is generated
**And** no race conditions detected (`-race` flag)

---

## Epic 2: Modbus I/O Integration

Automation engineers can connect to industrial Modbus TCP devices, read and write holding registers and coils, with reliable automatic reconnection when connections fail.

### Story 2.1: Modbus TCP Client Setup

As an **automation engineer**,
I want **go-plc to connect to Modbus TCP devices using the configured host and port**,
So that **I can communicate with industrial I/O devices using standard protocols**.

**Acceptance Criteria:**

**Given** a source configured with type `modbus-tcp` in YAML (host, port, unitId, timeout)
**When** the Modbus client is initialized
**Then** a TCP connection is established to the specified host:port
**And** the connection uses the simonvetter/modbus library
**And** INFO log: `[sourceName] - Modbus TCP connected (host: [ip]:[port])`

**Given** an unreachable Modbus device (wrong IP, device offline)
**When** connection attempt fails
**Then** ERROR log with human-readable message: `[sourceName] - Connection failed (host: [ip]:[port], error: [details])`
**And** the source is marked as disconnected
**And** other sources continue to operate independently

**Given** byte order and word order configuration (big-endian, high-word-first)
**When** multi-register values are read/written
**Then** the library is configured with correct endianness settings

---

### Story 2.2: Register & Coil Read/Write Operations

As an **automation engineer**,
I want **to read and write holding registers, input registers, coils, and discrete inputs**,
So that **I can exchange data with sensors, actuators, and I/O modules**.

**Acceptance Criteria:**

**Given** a variable configured with register type `holding` and address `0`
**When** a read operation is performed
**Then** the holding register at address 0 is read from the Modbus device
**And** the raw value is stored in the variable store

**Given** a variable configured with register type `coil` and address `0`
**When** a read operation is performed
**Then** the coil state (true/false) is read and stored as a boolean

**Given** a variable with dataType `uint32` (2 registers)
**When** reading the value
**Then** registers at address and address+1 are read and combined per configured word order

**Given** supported data types from simonvetter/modbus
**When** variables are configured
**Then** the following types work correctly: `bool`, `uint16`, `int16`, `uint32`, `int32`, `float32`, `uint64`, `int64`, `float64`

**Given** a write operation to a holding register or coil
**When** `SetVariable` is called with a new value
**Then** the value is written to the Modbus device
**And** the variable store is updated with the new value

---

### Story 2.3: Connection Manager with Automatic Reconnection

As an **automation engineer**,
I want **Modbus connections to automatically reconnect with exponential backoff when they fail**,
So that **the system recovers from network issues without manual intervention**.

**Acceptance Criteria:**

**Given** an established Modbus connection
**When** the connection is lost (network failure, device restart)
**Then** the source state changes to "disconnected"
**And** WARN log: `[sourceName] - Connection lost, reconnecting...`

**Given** a disconnected source with retryInterval configured (default 5s)
**When** reconnection is attempted
**Then** attempts use exponential backoff: 5s, 10s, 20s, 40s... up to max 5 minutes
**And** each attempt is logged at DEBUG level

**Given** a reconnection attempt succeeds
**When** connection is restored
**Then** the source state changes to "connected"
**And** INFO log: `[sourceName] - Reconnected after [N] attempts`
**And** polling resumes automatically

**Given** the connection manager
**When** I query source status
**Then** I can see: name, type, state (connected/disconnected), lastConnected timestamp, lastError message, reconnectAttempts count

---

### Story 2.4: Modbus Polling Loop

As an **automation engineer**,
I want **Modbus sources to continuously poll device registers at configured intervals**,
So that **variable values are kept up-to-date with actual device states**.

**Acceptance Criteria:**

**Given** a source with pollInterval of 100ms
**When** the source is connected and polling
**Then** all variables bound to that source are read every 100ms
**And** read values are written to the variable store (with scaling applied)

**Given** a polling cycle
**When** variables are updated in the store
**Then** the variable's timestamp is updated to the current time
**And** change notifications are sent (for future GraphQL subscriptions)

**Given** a read error during polling (timeout, device error)
**When** the error occurs
**Then** the variable value is NOT updated (stale value retained)
**And** WARN log with error details
**And** polling continues on next cycle

**Given** multiple variables on the same source
**When** polling occurs
**Then** reads are batched efficiently where possible (contiguous registers)
**And** polling completes within reasonable time

---

### Story 2.5: Multiple Sources & Integration Test

As an **automation engineer**,
I want **to configure multiple independent Modbus sources with different polling intervals**,
So that **I can connect to multiple devices with appropriate update rates**.

**Acceptance Criteria:**

**Given** configuration with multiple Modbus sources (e.g., remoteIO at 100ms, flowMeter at 500ms)
**When** the runtime starts
**Then** each source connects independently
**And** each source polls at its configured interval
**And** sources don't block each other (independent goroutines)

**Given** one source fails while others are healthy
**When** the failed source attempts reconnection
**Then** healthy sources continue operating normally
**And** the failed source reconnects independently

**Given** a Python pymodbus simulator (or similar)
**When** running integration tests
**Then** go-plc successfully:
- Connects to the simulator
- Reads holding registers with correct values
- Writes holding registers and verifies the write
- Reads coils with correct boolean values
- Handles simulator restart with automatic reconnection

**Given** `make test-integration` (or similar)
**When** tests run against the simulator
**Then** all Modbus functionality is verified end-to-end

---

## Epic 3: GraphQL API & Real-Time Data

External applications can query system status, variables, and tasks via GraphQL, and subscribe to real-time value changes through WebSocket connections.

### Story 3.1: GraphQL Schema & Server Setup

As a **developer**,
I want **a GraphQL API server with a well-defined schema**,
So that **external applications can query go-plc data using standard GraphQL tooling**.

**Acceptance Criteria:**

**Given** the go-plc application starts
**When** the GraphQL server initializes
**Then** an HTTP server listens on configurable port (default 8080)
**And** the `/graphql` endpoint accepts GraphQL requests
**And** INFO log: `GraphQL server listening on :8080`

**Given** the GraphQL schema definition
**When** schema is generated with gqlgen
**Then** types are defined for: `Variable`, `Source`, `Task`, `PLCStatus`
**And** field names follow `camelCase` convention per Architecture
**And** enum values use `SCREAMING_SNAKE_CASE` (e.g., `HOLDING`, `COIL`)

**Given** a valid GraphQL query
**When** sent to the `/graphql` endpoint
**Then** response is returned in <10ms for variable reads (NFR4)
**And** response follows GraphQL specification format

---

### Story 3.2: Query Resolvers

As an **IoT/edge developer**,
I want **to query current variable values, PLC status, task status, and source status via GraphQL**,
So that **I can integrate go-plc data into my applications**.

**Acceptance Criteria:**

**Given** the query `{ variables { name value unit timestamp } }`
**When** executed
**Then** returns all variables with current values from the variable store
**And** includes: name, value (scaled), rawValue, unit, timestamp, tags, sourceName

**Given** the query `{ variables(filter: { tags: ["tank1"] }) { name value } }`
**When** executed
**Then** returns only variables that have the "tank1" tag
**And** filter supports: names, tags, source

**Given** the query `{ plcStatus { state uptime taskCount variableCount } }`
**When** executed
**Then** returns current PLC runtime status
**And** state is "running" or "stopped"
**And** uptime is duration since start

**Given** the query `{ tasks { name description scanRate enabled lastExecutionTime } }`
**When** executed
**Then** returns all registered tasks with their configuration and status

**Given** the query `{ sources { name type state lastConnected lastError } }`
**When** executed
**Then** returns all configured sources with connection status

---

### Story 3.3: Variable Mutations

As an **IoT/edge developer**,
I want **to write variable values through the GraphQL API**,
So that **I can control the PLC from external applications**.

**Acceptance Criteria:**

**Given** the mutation `mutation { setVariable(name: "pumpRunning", value: true) { name value } }`
**When** executed
**Then** the variable value is updated in the variable store
**And** if the variable is bound to a Modbus coil/register, the value is written to the device
**And** the mutation returns the updated variable

**Given** a mutation to set a non-existent variable
**When** executed
**Then** a GraphQL error is returned with message: `Variable not found: [name]`
**And** HTTP status remains 200 (GraphQL error handling)

**Given** a mutation with invalid value type (e.g., string for boolean)
**When** executed
**Then** a GraphQL validation error is returned

---

### Story 3.4: Real-Time Subscriptions

As an **IoT/edge developer**,
I want **to subscribe to variable value changes via GraphQL subscriptions**,
So that **my application receives real-time updates without polling**.

**Acceptance Criteria:**

**Given** the subscription `subscription { variableUpdates { name value timestamp } }`
**When** a WebSocket connection is established
**Then** the client receives updates whenever any variable changes
**And** updates are delivered within 100ms of value change (NFR5)

**Given** the subscription `subscription { variableUpdates(filter: { tags: ["tank1"] }) { name value } }`
**When** variables with "tank1" tag are updated
**Then** only those variables are sent to the subscriber
**And** variables without the tag do not trigger updates

**Given** the variable store change notification channel
**When** a variable value changes (from Modbus poll, task write, or API mutation)
**Then** the change is broadcast to all active subscriptions
**And** filtering is applied per subscriber

**Given** a WebSocket connection is closed
**When** the client disconnects
**Then** the subscription is cleaned up
**And** no goroutine leaks occur

---

### Story 3.5: GraphQL Playground & Testing

As a **developer**,
I want **a built-in GraphQL Playground for testing queries and an automated test suite**,
So that **I can explore the API interactively and ensure it works correctly**.

**Acceptance Criteria:**

**Given** the go-plc application is running
**When** I navigate to `/playground` in a browser
**Then** the GraphQL Playground UI is displayed
**And** I can write and execute queries, mutations, and subscriptions
**And** schema documentation is available in the playground

**Given** the GraphQL test suite
**When** `go test ./internal/api/...` is run
**Then** tests verify:
- All query resolvers return correct data
- Mutations update variables correctly
- Subscriptions deliver updates
- Error cases return appropriate GraphQL errors
- Response times meet NFR4 (<10ms for queries)

**Given** the playground is served
**When** in production mode (future consideration)
**Then** playground can be disabled via configuration for security

---

## Epic 4: WebUI Monitoring Dashboard

Operators can monitor PLC status, view source connections, watch real-time variable values, manage tasks, and troubleshoot issues through a responsive web interface.

### Story 4.1: Frontend Project Setup

As a **developer**,
I want **the WebUI frontend initialized with the correct tech stack and design system**,
So that **I can build UI components following the UX design specification**.

**Acceptance Criteria:**

**Given** the `web/` directory
**When** I run the frontend initialization commands
**Then** Vite + React + TypeScript project is created
**And** Tailwind CSS is configured with custom color tokens from UX spec
**And** shadcn/ui is initialized with required components (Table, Badge, Card, Button, Collapsible, Input)
**And** urql is configured with GraphQL WebSocket subscription support
**And** Roboto font is loaded

**Given** the Tailwind configuration
**When** custom colors are applied
**Then** brand colors are available: graphite (#353535), stormy-teal (#3C6E71), yale-blue (#284B63), alabaster-grey (#D9D9D9)
**And** status colors are available: status-ok (#4A7C59), status-error (#C45C5C), status-warning (#D4A84B)

**Given** `pnpm dev` is run in the `web/` directory
**When** the dev server starts
**Then** the frontend is accessible at localhost:5173 (or configured port)
**And** hot module replacement works for rapid development

---

### Story 4.2: Layout & Navigation Components

As an **operator**,
I want **a clear layout with navigation and status visibility**,
So that **I can quickly see system health and navigate between views**.

**Acceptance Criteria:**

**Given** the WebUI loads
**When** the page renders
**Then** the InfoBar is displayed at the top (full width)
**And** the Sidebar is displayed on the left (collapsed by default, 56px)
**And** the main content area fills the remaining space

**Given** the InfoBar component
**When** rendered
**Then** it shows: Alerts toggle, error badge (if errors), warning badge (if warnings), Live indicator, PLC status
**And** Live indicator pulses to confirm WebSocket connection
**And** PLC status shows "Running" (green dot) or "Stopped" (red dot)

**Given** the Sidebar component
**When** collapsed (default)
**Then** it shows logo icon and navigation icons with tooltips
**When** expanded (click logo)
**Then** it shows full navigation labels and theme toggle
**And** navigation items: Overview (active), Config (disabled/future), Tasks (disabled/future), Logs (disabled/future)

**Given** the AlertsPanel component
**When** alerts exist and toggle is clicked
**Then** panel expands below InfoBar showing alert details
**And** each alert shows: severity icon, source/variable name, message, timestamp
**And** clicking toggle again collapses the panel

---

### Story 4.3: Sources Panel

As an **operator**,
I want **to see all source connections and their status at a glance**,
So that **I can quickly identify connectivity issues**.

**Acceptance Criteria:**

**Given** the Sources panel in the Overview page
**When** rendered
**Then** it displays a table with columns: Name, Type, Status, Address, Last Poll
**And** each source shows connection status with StatusDot (green=connected, red=disconnected)
**And** panel header shows source count

**Given** a connected source
**When** displayed in the table
**Then** status shows green dot with "connected" text
**And** last poll shows time since last successful poll (e.g., "45ms ago")

**Given** a disconnected source
**When** displayed in the table
**Then** status shows red dot with "disconnected" text
**And** last error message is visible (or on hover/expand)
**And** the row has subtle visual emphasis per ISA-101 (errors stand out)

**Given** the sources query
**When** data is fetched via GraphQL
**Then** `{ sources { name type state lastConnected lastError } }` is used
**And** data updates when source status changes

---

### Story 4.4: Variables Panel

As an **operator**,
I want **to view all variables with their current values and search/filter them**,
So that **I can monitor process values and find specific variables quickly**.

**Acceptance Criteria:**

**Given** the Variables panel in the Overview page
**When** rendered
**Then** it displays a searchable table with columns: Name, Value, Unit, Source, Tags
**And** panel header shows variable count
**And** search input is at the top of the panel

**Given** variables in the table
**When** displayed
**Then** values are shown in monospace font (Roboto Mono)
**And** source is shown as SourceTag component (colored pill with connection status)
**And** tags are shown as neutral badges

**Given** the search input
**When** user types a search term
**Then** variables are filtered instantly (client-side)
**And** filter matches variable name OR any tag
**And** "No variables match '[query]'" shown if no results

**Given** real-time variable updates
**When** a variable value changes
**Then** the new value is displayed immediately (via GraphQL subscription)
**And** a brief highlight animation indicates the change (subtle flash)

---

### Story 4.5: Tasks Panel

As an **operator**,
I want **to view all tasks with their status and enable/disable them**,
So that **I can manage which control logic is running**.

**Acceptance Criteria:**

**Given** the Tasks panel in the Overview page
**When** rendered
**Then** it displays a table with columns: Name, Description, Scan Rate, Status, Last Execution, Actions
**And** panel header shows task count

**Given** a task in the table
**When** displayed
**Then** scan rate is shown as a badge (e.g., "100ms")
**And** status shows "Enabled" or "Disabled"
**And** last execution time shows duration (e.g., "12ms")

**Given** the Actions column
**When** a task is enabled
**Then** a "Disable" button is shown
**When** a task is disabled
**Then** an "Enable" button is shown

**Given** the enable/disable button is clicked
**When** the action is triggered
**Then** a GraphQL mutation is sent to toggle the task
**And** the UI updates to reflect the new state
**And** confirmation feedback is provided (button state change)

---

### Story 4.6: PLC Runtime Controls

As an **operator**,
I want **to start and stop the PLC runtime from the WebUI**,
So that **I can control the system without command-line access**.

**Acceptance Criteria:**

**Given** the InfoBar shows PLC status
**When** PLC is running
**Then** status shows green dot with "PLC Running"
**And** a "Stop" button is available (or in System panel)

**Given** the InfoBar shows PLC status
**When** PLC is stopped
**Then** status shows red dot with "PLC Stopped"
**And** a "Start" button is available

**Given** the Start button is clicked
**When** action is triggered
**Then** a GraphQL mutation `startPlc` is sent
**And** the PLC runtime starts
**And** status updates to "Running"

**Given** the Stop button is clicked
**When** action is triggered
**Then** a confirmation dialog is shown (stopping PLC is significant)
**And** upon confirmation, GraphQL mutation `stopPlc` is sent
**And** the PLC runtime stops gracefully
**And** status updates to "Stopped"

---

### Story 4.7: Real-Time Updates & Integration

As an **operator**,
I want **the WebUI to show real-time updates without manual refresh**,
So that **I always see current system state**.

**Acceptance Criteria:**

**Given** the WebUI is open
**When** a WebSocket connection is established
**Then** the Live indicator in InfoBar shows pulsing green dot
**And** GraphQL subscriptions are active

**Given** a variable value changes (from Modbus poll or task)
**When** the change occurs
**Then** the Variables panel updates within 100ms
**And** no manual refresh is needed

**Given** a source connection status changes
**When** source connects or disconnects
**Then** the Sources panel updates immediately
**And** alerts are added/removed as appropriate
**And** alert badges in InfoBar update

**Given** the WebSocket connection is lost
**When** disconnect is detected
**Then** the Live indicator shows grey/warning state
**And** reconnection is attempted automatically
**And** user is informed of connection status

**Given** the frontend build
**When** `pnpm build` is run
**Then** optimized production assets are generated in `web/dist/`
**And** assets are ready for Go embed integration (Epic 6)

---

## Epic 5: OPC UA SCADA Integration

Traditional SCADA systems (Ignition, Kepware) can connect to go-plc via OPC UA to read and write variable values for industrial integration.

### Story 5.1: OPC UA Server Setup

As a **system integrator**,
I want **go-plc to expose an OPC UA server that SCADA systems can connect to**,
So that **I can integrate with traditional industrial SCADA platforms**.

**Acceptance Criteria:**

**Given** the go-plc application starts with OPC UA enabled in config
**When** the OPC UA server initializes
**Then** the server listens on configurable port (default 4840)
**And** the server uses the gopcua/opcua library
**And** INFO log: `OPC UA server listening on :4840`

**Given** an OPC UA client (e.g., Ignition, UaExpert)
**When** connecting to `opc.tcp://[host]:4840`
**Then** the connection is established successfully
**And** the server endpoint is browsable

**Given** OPC UA server configuration in YAML
**When** settings are specified
**Then** port, application name, and application URI are configurable
**And** security mode can be configured (None for MVP per NFR19)

---

### Story 5.2: Variable Node Generation

As a **system integrator**,
I want **all go-plc variables to appear as OPC UA nodes**,
So that **SCADA systems can browse and access all process values**.

**Acceptance Criteria:**

**Given** variables defined in the go-plc configuration
**When** the OPC UA server starts
**Then** each variable is exposed as an OPC UA node
**And** nodes are organized under a "Variables" folder in the address space

**Given** an OPC UA node for a variable
**When** browsing the node
**Then** the node has:
- NodeId based on variable name
- DisplayName matching variable name
- DataType matching variable type (Boolean, Int16, Float, etc.)
- Description (if available from config)

**Given** variable tags in configuration
**When** nodes are generated
**Then** tags can be used to organize nodes into folders (optional enhancement)

**Given** new variables (if config allows runtime addition - future)
**When** variables are added
**Then** corresponding OPC UA nodes are created dynamically

---

### Story 5.3: Read/Write Handlers

As a **system integrator**,
I want **SCADA systems to read and write variable values via OPC UA**,
So that **bidirectional data exchange is possible with traditional industrial systems**.

**Acceptance Criteria:**

**Given** an OPC UA client reads a variable node
**When** the read request is received
**Then** the current value is returned from the variable store
**And** the response includes timestamp and status code (Good)
**And** read completes within 50ms (NFR6)

**Given** an OPC UA client writes to a variable node
**When** the write request is received
**Then** the value is updated in the variable store
**And** if the variable is bound to Modbus, the value is written to the device
**And** the write returns success status

**Given** an OPC UA client subscribes to a variable (monitored item)
**When** the variable value changes
**Then** the client receives a data change notification
**And** updates are delivered based on configured sampling interval

**Given** a read/write to a non-existent node
**When** the request is processed
**Then** an appropriate OPC UA error status is returned (BadNodeIdUnknown)

---

### Story 5.4: Server Status & Integration Test

As a **system integrator**,
I want **to view OPC UA server status and validate integration with Ignition SCADA**,
So that **I can verify the integration works correctly for production deployment**.

**Acceptance Criteria:**

**Given** the OPC UA server is running
**When** I query server status (via GraphQL or internal API)
**Then** I can see: server state (running/stopped), port, connected clients count, endpoint URL

**Given** the sources query in GraphQL
**When** OPC UA source type is included
**Then** `{ sources { name type state } }` includes OPC UA server status

**Given** Ignition SCADA (or UaExpert for testing)
**When** configured to connect to go-plc OPC UA server
**Then** connection is established successfully
**And** variables are browsable in the OPC UA tree
**And** values can be read and displayed in Ignition tags
**And** values can be written from Ignition and reflected in go-plc

**Given** integration test with OPC UA client
**When** `go test ./internal/opcua/...` is run
**Then** tests verify:
- Server starts and accepts connections
- Nodes are created for all variables
- Read operations return correct values
- Write operations update variable store
- Monitored items receive data change notifications

---

## Epic 6: Production Deployment & Documentation

Developers can build single-binary deployments with embedded WebUI, run go-plc as a systemd service on Linux, and deploy across Linux and Windows platforms.

### Story 6.1: Single Binary Build with Embedded WebUI

As a **system administrator**,
I want **GoPLC to compile as a single binary with the WebUI embedded**,
So that **deployment is simple with no external file dependencies**.

**Acceptance Criteria:**

**Given** the Go build process is configured with `//go:embed` directives
**When** I run `go build` for the GoPLC project
**Then** the resulting binary contains all WebUI static assets (HTML, JS, CSS)
**And** the binary serves the WebUI without requiring external files
**And** the binary size remains reasonable (<50MB per NFR17)

---

### Story 6.2: Cross-Platform Build Configuration

As a **system administrator**,
I want **GoPLC binaries available for Linux (amd64/arm64) and Windows (amd64)**,
So that **I can deploy on my target hardware platform**.

**Acceptance Criteria:**

**Given** the build system is configured for cross-compilation
**When** I build for each target platform (GOOS/GOARCH combinations)
**Then** functional binaries are produced for Linux amd64, Linux arm64, and Windows amd64
**And** each binary runs correctly on its target platform
**And** platform-specific features (file paths, signals) work correctly

---

### Story 6.3: systemd Service Configuration

As a **system administrator**,
I want **GoPLC to run as a systemd service on Linux**,
So that **it starts automatically on boot and can be managed with standard tools**.

**Acceptance Criteria:**

**Given** a systemd unit file is created for GoPLC
**When** I install and enable the service
**Then** GoPLC starts automatically on system boot
**And** the service can be controlled via `systemctl start/stop/restart/status`
**And** logs are captured by journald and accessible via `journalctl`
**And** the service restarts automatically on failure (with backoff)
