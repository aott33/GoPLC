# Project Context: go-plc

> **Purpose:** Concise guide for AI agents implementing go-plc. Read this before writing any code.

## Project Overview

**go-plc** is a soft PLC (Programmable Logic Controller) written in Go with a React-based monitoring WebUI. It communicates with industrial devices via Modbus TCP/RTU, exposes data through OPC UA server and Sparkplug B MQTT, and provides a GraphQL API for the embedded web interface.

**Target Platform:** Real-time Linux OS (PREEMPT_RT kernel recommended)
**Development Platform:** Any OS (Linux, Windows, macOS)

## Critical Rules

### Naming Conventions (MUST FOLLOW)

| Context | Convention | Example |
|---------|------------|---------|
| Go packages | lowercase, single word | `config`, `modbus`, `variables` |
| Go files | lowercase with underscores | `modbus_client.go`, `variable_store.go` |
| Go exported | PascalCase | `GetVariable`, `PLCRuntime` |
| Go unexported | camelCase | `pollLoop`, `reconnectTimer` |
| JSON struct tags | camelCase | `json:"pollInterval"` |
| GraphQL types | PascalCase | `Variable`, `PLCStatus` |
| GraphQL fields | camelCase | `tankLevel`, `pollInterval` |
| GraphQL enums | SCREAMING_SNAKE_CASE | `HOLDING`, `INPUT`, `COIL` |
| YAML config | camelCase | `pollInterval`, `byteOrder` |

### Error Message Format (MUST FOLLOW)

```
[Component] - [Human description] (context: [details])
```

Examples:
- `remoteIO - Connection timeout after 3 retries (host: 192.168.1.100:502)`
- `tankLevel - Value exceeds configured range (value: 105, max: 100)`

### Logging (MUST FOLLOW)

Use `log/slog` (Go 1.21+ standard library) for ALL logging:

```go
// Correct
slog.Info("Source connected", "source", source.Name, "host", host)
slog.Error("Connection failed", "source", source.Name, "err", err)

// Wrong - never use log.Printf
log.Printf("Connected to %s", host)
```

| Level | Usage |
|-------|-------|
| `DEBUG` | Variable values, poll timing, detailed diagnostics |
| `INFO` | Startup, config loaded, connections established |
| `WARN` | Reconnection attempts, values out of range |
| `ERROR` | Connection failures, task errors |

## Architecture Essentials

### Variable Store is Central

All protocols read/write through the Variable Store (`internal/variables/`):

```
Modbus Poller → Variable Store ← Task Execution
                     ↓
              Change Notifications
                     ↓
    ┌────────────────┼────────────────┐
    ↓                ↓                ↓
GraphQL Subs    OPC UA Server   Sparkplug B
```

**Implementation:** `sync.RWMutex` + `map[string]*Variable`

### Package Boundaries

All implementation code goes in `internal/` - this is private and cannot be imported externally.

| Package | Responsibility |
|---------|---------------|
| `internal/config` | YAML parsing, validation, config structs |
| `internal/variables` | Variable store, scaling engine, change notifications |
| `internal/runtime` | PLC scheduler, state machine, lifecycle |
| `internal/modbus` | Modbus TCP/RTU client, connection manager, polling |
| `internal/opcua` | OPC UA server, node generation, handlers |
| `internal/sparkplug` | Sparkplug B client, NBIRTH/NDATA/NDEATH encoding |
| `internal/api` | GraphQL server, resolvers, subscriptions |
| `internal/tasks` | Task discovery, execution runtime |

### Required Libraries

These are mandated - do not substitute:

| Library | Purpose |
|---------|---------|
| `github.com/simonvetter/modbus` | Modbus TCP/RTU client |
| `github.com/gopcua/opcua` | OPC UA server |
| `github.com/eclipse/paho.mqtt.golang` | MQTT for Sparkplug B |
| `github.com/99designs/gqlgen` | GraphQL server |

### Frontend Stack

| Library | Purpose |
|---------|---------|
| React + Vite + TypeScript | Framework |
| urql + graphql-ws | GraphQL client with subscriptions |
| Tailwind CSS + shadcn/ui | Styling and components |
| Roboto | Font family |

## Configuration Schema

Type-discriminated sources with flat variable list:

```yaml
sources:
  - name: remoteIO
    type: modbus-tcp  # or modbus-rtu
    config:
      host: 192.168.1.100
      port: 502
      unitId: 1
      timeout: 1s
      pollInterval: 100ms
      retryInterval: 5s
      byteOrder: big-endian
      wordOrder: high-word-first

variables:
  - name: tankLevel
    source: remoteIO
    register:
      type: holding  # holding, input, coil, discrete
      address: 0     # 0-based
    dataType: uint16
    scale:           # Optional - omit for raw value
      rawMin: 0
      rawMax: 65535
      engMin: 0
      engMax: 100
      unit: "%"
    tags: [tank1, level]
```

### Supported Data Types

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

### Scaling Formula

When `scale` block is present:
```
scaled = ((raw - rawMin) / (rawMax - rawMin)) * (engMax - engMin) + engMin
```

## Anti-Patterns (NEVER DO)

### Wrong JSON Tags
```go
// WRONG - snake_case
type Source struct {
    PollInterval string `json:"poll_interval"`
}

// CORRECT - camelCase
type Source struct {
    PollInterval string `json:"pollInterval"`
}
```

### Wrong Error Messages
```go
// WRONG - no context
return errors.New("connection error")

// CORRECT - structured with context
return fmt.Errorf("remoteIO - Connection failed after %d retries (host: %s)", retries, host)
```

### Wrong Logging
```go
// WRONG - printf style
log.Printf("Connected to %s", host)
fmt.Println("Debug:", value)

// CORRECT - structured slog
slog.Info("Connected", "host", host)
slog.Debug("Value read", "variable", name, "value", value)
```

### Wrong Package Structure
```go
// WRONG - putting code outside internal/
package plc  // in /pkg/plc/

// CORRECT - everything in internal/
package runtime  // in /internal/runtime/
```

## Testing Patterns

- **Location:** Co-located with source files
- **Naming:** `*_test.go` suffix
- **Example:** `internal/modbus/client.go` → `internal/modbus/client_test.go`

## Build & Deployment

### Single Binary Output
Frontend is embedded via `//go:embed` directive. Final binary includes all assets.

### Docker
- `Dockerfile`: Multi-stage build (frontend → backend → minimal runtime)
- `docker-compose.yaml`: go-plc service only
- Modbus simulator: External repository, will be added to compose later

### Development Commands (Makefile)
```makefile
dev         # Run Go with frontend dev server
build       # Build production binary with embedded frontend
test        # Run all Go tests
lint        # Run golangci-lint
docker      # Build Docker image
docs        # Build Docusaurus documentation site
docs-dev    # Run documentation site in dev mode
```

## Performance Requirements

- **Task execution overhead:** <50µs per cycle
- **GraphQL response:** <10ms
- **Memory stability:** 24+ hours without leaks
- **Graceful shutdown:** Within 5 seconds

## Deferred to Post-MVP

Do NOT implement these features:
- Authentication/authorization
- TLS/encryption for protocols
- Role-based access control
- Config hot-reload
- Task hot-reload

## Quick Reference

```go
// Standard imports for this project
import (
    "log/slog"
    "sync"
    "context"

    "github.com/simonvetter/modbus"
    "github.com/gopcua/opcua"
    "github.com/eclipse/paho.mqtt.golang"
)

// Variable store access pattern
type Store struct {
    mu   sync.RWMutex
    vars map[string]*Variable
}

func (s *Store) Get(name string) (*Variable, bool) {
    s.mu.RLock()
    defer s.mu.RUnlock()
    v, ok := s.vars[name]
    return v, ok
}
```

---

**Reference:** See [architecture.md](architecture.md) for complete architectural decisions and project structure.
