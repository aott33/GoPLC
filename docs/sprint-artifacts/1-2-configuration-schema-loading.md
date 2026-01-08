# Story 1.2: Configuration Schema & Loading

**Status:** ready-for-dev
**Epic:** 1 - Project Foundation & Core Runtime
**Estimated Effort:** Medium (3-4 hours)
**Priority:** Critical Path - Blocks Stories 1.3-1.8

---

## Overview

This story implements the YAML configuration parsing and validation system for go-plc. You'll create the configuration structs that define sources (Modbus devices) and variables, implement a loader that parses and validates the configuration, and integrate it into the main application startup sequence.

**Why This Matters:**
Configuration is the foundation of go-plc. Every other component (variable store, Modbus client, OPC UA server, GraphQL API) depends on this configuration. Getting the schema right now prevents painful migrations later. The type-discriminated source pattern enables future protocol additions without breaking changes.

---

## Context from Previous Work

**Module Path:** `github.com/aott33/go-plc` (defined in go.mod)

**Current main.go Structure:**
The entry point (`cmd/go-plc/main.go`) already has slog configured with JSON output. Your config loading will be called early in the startup sequence - after flag parsing but before variable store initialization. The existing main.go comments outline where config loading fits:
```
// parse command line flags
// load and validate yaml file    <-- THIS STORY
// initialize variable store
// connect to sources
// ...
```

**Development Environment Notes from Story 1.1:**
- Go version is 1.25 (verify with `go version`)
- Windows users: Makefile requires bash - use WSL, Git Bash, or VSCode integrated terminal
- `go mod tidy` should be run after adding yaml.v3 dependency

---

## User Story

> **As a** developer,
> **I want** to define sources and variables in a YAML configuration file that is parsed and validated on startup,
> **So that** I can configure the PLC without code changes and catch configuration errors early.

---

## Acceptance Criteria

### AC1: YAML Configuration Parsing with Type-Discriminated Sources

**Given** a YAML configuration file with the architecture-mandated schema
**When** the application starts with `-config path/to/config.yaml`
**Then** sources are parsed with protocol-specific config blocks (modbus-tcp, modbus-rtu)
**And** variables are parsed with source references, register definitions, and optional scaling

### AC2: Configuration Validation with Clear Error Messages

**Given** a YAML file with invalid configuration (missing required field, invalid type, unknown source reference)
**When** the application attempts to load the config
**Then** a clear, human-readable error message is logged in format: `[config] - [description] (context: [details])`
**And** the application exits with non-zero status before starting the runtime

### AC3: Successful Configuration Loading

**Given** a valid configuration
**When** config loading completes
**Then** all sources and variables are accessible via typed Go structs
**And** INFO level log confirms "Configuration loaded successfully"

---

## Sprint Change Proposal

**IMPORTANT:** This story's implementation approach has been enhanced by [Sprint Change Proposal 2025-12-20](../sprint-change-proposal-2025-12-20.md).

**Key Changes:**
1. **YAML Library:** Use `github.com/goccy/go-yaml` instead of `gopkg.in/yaml.v3`
2. **Architecture Pattern:** Implement registry pattern for protocol-specific source configs
3. **Package Structure:** Protocol configs live in `internal/source/modbus/` instead of `internal/config/`

**Implementation Guide:** Follow Section 4 of the sprint change proposal for complete technical specifications, code examples, and file structure.

**Why This Matters:** The registry pattern eliminates tight coupling and enables future protocol additions (OPC-UA, MQTT, EtherNet/IP) without modifying core config code.

---

## Technical Requirements

### Architecture Compliance

You MUST follow these patterns from the architecture document AND the sprint change proposal:

**Naming Conventions:**
- Package names: `config`, `source` (lowercase, single word)
- File names: `config.go`, `loader.go`, `config_test.go`, `registry.go`, `tcp.go`, `rtu.go`
- Exported types: `PascalCase` (e.g., `Config`, `Source`, `Variable`)
- Unexported: `camelCase` (e.g., `validateSources`, `loadFile`)
- JSON/YAML tags: `camelCase` (e.g., `json:"pollInterval"` NOT `json:"poll_interval"`)

**Error Message Format:**
```
[Component] - [Human description] (context: [details])
```
Examples:
- `[config] - Missing required field 'host' (source: remoteIO, type: modbus-tcp)`
- `[config] - Unknown source reference (variable: tankLevel, source: unknownDevice)`

**Logging:**
- Use `log/slog` exclusively (NOT `log`, `logrus`, `zap`, or `fmt.Println`)
- JSON structured format
- Appropriate log levels: INFO for success, ERROR for failures

### YAML Configuration Schema

The architecture mandates this exact schema structure. Study it carefully - your structs must match this exactly:

```yaml
# Application-level settings
logLevel: info  # debug, info, warn, error

# Sources with type-discriminated config blocks
sources:
  - name: remoteIO              # Unique identifier
    type: modbus-tcp            # Protocol type
    config:                     # Protocol-specific block
      host: 192.168.1.100
      port: 502
      unitId: 1
      timeout: 1s
      pollInterval: 100ms
      retryInterval: 5s
      byteOrder: big-endian     # big-endian or little-endian
      wordOrder: high-word-first # high-word-first or low-word-first

  - name: serialDevice
    type: modbus-rtu
    config:
      device: /dev/ttyUSB0
      baudRate: 19200
      dataBits: 8
      parity: none              # none, even, odd
      stopBits: 2
      unitId: 1
      timeout: 1s
      pollInterval: 100ms

# Flat variable list with source references
variables:
  - name: tankLevel
    source: remoteIO            # References source by name
    register:
      type: holding             # holding, input, coil, discrete
      address: 0                # 0-based offset
    dataType: uint16            # See supported types below
    scale:                      # Optional - omit for raw values
      rawMin: 0
      rawMax: 65535
      engMin: 0
      engMax: 100
      unit: "%"
    tags: [tank1, level]        # Optional - for filtering/grouping

  - name: pumpRunning
    source: remoteIO
    register:
      type: coil
      address: 0
    dataType: bool
    tags: [tank1, pump]
```

### Supported Data Types

From the simonvetter/modbus library - you'll need these for validation:

| dataType | Registers | Description |
|----------|-----------|-------------|
| `bool` | N/A | Coils and discrete inputs only |
| `uint16` | 1 | Unsigned 16-bit integer |
| `int16` | 1 | Signed 16-bit integer |
| `uint32` | 2 | Unsigned 32-bit integer |
| `int32` | 2 | Signed 32-bit integer |
| `float32` | 2 | 32-bit floating point |
| `uint64` | 4 | Unsigned 64-bit integer |
| `int64` | 4 | Signed 64-bit integer |
| `float64` | 4 | 64-bit floating point |

### Register Types

| type | Function Codes | Use Case |
|------|----------------|----------|
| `holding` | FC03/FC06/FC16 | Read/write registers |
| `input` | FC04 | Read-only registers |
| `coil` | FC01/FC05/FC15 | Read/write single bits |
| `discrete` | FC02 | Read-only single bits |

---

## Implementation Guide

**CRITICAL:** Follow the complete implementation guide in [Sprint Change Proposal - Section 4](../sprint-change-proposal-2025-12-20.md#4-technical-specification-for-human-developer). The sections below provide high-level context, but Section 4 of the proposal contains the authoritative code examples and structure.

### Step 1: Create the Source Registry Pattern

**Locations:**
- `internal/source/source.go` - Interface definitions
- `internal/source/registry.go` - Registry implementation
- `internal/source/modbus/tcp.go` - Modbus TCP config
- `internal/source/modbus/rtu.go` - Modbus RTU config
- `internal/source/modbus/duration.go` - Duration helper

**What to implement:**

1. **Define Source Interfaces** (`internal/source/source.go`):
   - `Source` interface (Name(), Type() methods - full interface in Epic 2)
   - `SourceConfig` interface (Validate(), SourceType(), SourceName() methods)
   - `SourceFactory` type (function signature for creating configs from YAML)

2. **Implement the Registry** (`internal/source/registry.go`):
   - Global registry map with mutex protection
   - `Register(typeName string, factory SourceFactory)` - Called from protocol init() functions
   - `ParseConfig(typeName, name string, configNode yaml.Node)` - Looks up factory and parses config
   - `RegisteredTypes()` - Returns list of registered types for error messages

3. **Implement Modbus Protocol Configs** (`internal/source/modbus/tcp.go` and `rtu.go`):
   - Each file implements the config struct, Validate() method, and factory function
   - Each file has an `init()` function that calls `source.Register()`
   - See proposal Section 4.4 and 4.5 for complete code

4. **Create Duration Helper** (`internal/source/modbus/duration.go`):
   - Wraps `time.Duration` for YAML unmarshaling
   - Parses strings like "100ms", "5s" into `time.Duration`
   - See proposal Section 4.6 for complete code

### Step 2: Create Core Configuration Types

**Location:** `internal/config/config.go`

**What to implement:**

1. **Create the root Config struct** containing:
   - `LogLevel` field (string with validation for debug/info/warn/error)
   - `Sources` slice (holds `Source` structs with `source.SourceConfig` interface)
   - `Variables` slice

2. **Create Source struct:**
   - `Name` (string)
   - `Type` (string)
   - `Config` (source.SourceConfig interface - holds protocol-specific config)

3. **Create Variable struct** with fields:
   - `Name` (string, required, unique)
   - `Source` (string, required - must reference existing source)
   - `Register` (RegisterConfig)
   - `DataType` (string - must be one of the supported types)
   - `Scale` (optional *ScaleConfig)
   - `Tags` ([]string, optional)

4. **Create RegisterConfig struct:**
   - `Type` (string: "holding", "input", "coil", "discrete")
   - `Address` (uint16 - 0-based offset)

5. **Create ScaleConfig struct:**
   - `RawMin`, `RawMax`, `EngMin`, `EngMax` (float64)
   - `Unit` (string - display unit like "%", "PSI", "degC")

6. **Implement Config.Validate() method:**
   - Validates all sources by calling their `Config.Validate()` methods
   - Validates variables (see Section 4.8 in proposal for complete logic)
   - Collects ALL errors before returning

See proposal Section 4.8 for complete code example.

### Step 3: Implement the Configuration Loader

**Location:** `internal/config/loader.go`

**What to implement:**

1. **Create rawSource and rawConfig structs** for initial YAML parsing:
   - `rawSource` has Name, Type, and `Config yaml.Node` (deferred parsing)
   - `rawConfig` has LogLevel, Sources slice, Variables slice

2. **Implement Load function**:
   - Read file with `os.ReadFile`
   - Unmarshal into `rawConfig` using `github.com/goccy/go-yaml`
   - For each rawSource, call `source.ParseConfig()` to get typed config
   - Build final `Config` struct with typed sources
   - Call `cfg.Validate()` before returning

3. **Default value application:**
   - LogLevel defaults to "info" if not specified
   - Protocol-specific defaults are applied in each protocol's Validate() method

See proposal Section 4.7 for complete implementation with error handling.

### Step 4: Validation (Already in Config.Validate())

**Location:** `internal/config/config.go` (within the Validate() method)

Validation is implemented in the `Config.Validate()` method (see Step 2). The key rules are:

1. **Source validation:** Delegated to each protocol's `SourceConfig.Validate()` method
2. **Variable validation:** Implemented in helper functions (see proposal Section 4.8)
3. **Error collection:** ALL errors are collected before returning

**Important:** Protocol-specific validation (like checking required fields) happens in each protocol package's Validate() method. The config package only validates cross-cutting concerns like duplicate names and variable-to-source references.

### Step 5: Integrate with Main Application

**Location:** `cmd/go-plc/main.go`

**What to implement:**

1. **Add blank import for protocol registration:**
   ```go
   import (
       _ "github.com/aott33/go-plc/internal/source/modbus"
   )
   ```
   This triggers the `init()` functions in tcp.go and rtu.go, registering them with the source registry.

2. **Add command-line flag and config loading:**
   - Add `-config` flag for config file path
   - Call `config.Load(configPath)`
   - On error: log and exit
   - On success: log with source/variable counts

See proposal Section 4.9 for the complete import block structure.

### Step 6: Create Example Configuration File

**Location:** `config.yaml` (project root)

The YAML configuration schema is unchanged from the original story specification. See the YAML Configuration Schema section above for the complete example. The registry pattern is an internal implementation detail - the user-facing YAML remains identical.

### Step 7: Write Unit Tests

**Locations:**
- `internal/config/config_test.go` - Config package tests
- `internal/source/modbus/tcp_test.go` - TCP config tests (optional but recommended)
- `internal/source/modbus/rtu_test.go` - RTU config tests (optional but recommended)

**Test cases to implement:**

1. **Happy path tests:**
   - Load valid config with modbus-tcp and modbus-rtu sources
   - Verify default values are applied in protocol Validate() methods
   - Verify all fields are correctly parsed

2. **Error case tests:**
   - Missing config file
   - Invalid YAML syntax
   - Unknown source type (registry returns error with registered types)
   - Variable referencing non-existent source
   - Invalid dataType
   - Bool dataType with holding register (should fail)
   - Duplicate source/variable names

3. **Registry tests:**
   - Test that unknown source type produces helpful error message
   - Optionally create a mock "test-source" type to demonstrate extensibility

Use table-driven test pattern for validation cases.

---

## File Structure After Completion

```
internal/
├── config/
│   ├── doc.go          # Package documentation (exists from Story 1.1)
│   ├── config.go       # Config, Source, Variable, Register, Scale types + Validate()
│   ├── loader.go       # Load() function, rawSource/rawConfig for parsing
│   └── config_test.go  # Unit tests
├── source/
│   ├── source.go       # Source interface, SourceConfig interface, SourceFactory type
│   ├── registry.go     # Register(), ParseConfig(), RegisteredTypes()
│   └── modbus/
│       ├── tcp.go      # TCPConfig + init() registration + parseTCPConfig factory
│       ├── rtu.go      # RTUConfig + init() registration + parseRTUConfig factory
│       └── duration.go # Duration type for YAML time parsing

cmd/go-plc/
└── main.go             # Updated with blank import for modbus package

config.yaml             # Example configuration (unchanged schema)
go.mod                  # Updated: github.com/goccy/go-yaml dependency
```

See [Sprint Change Proposal - Section 5](../sprint-change-proposal-2025-12-20.md#5-updated-file-structure) for detailed file structure.

---

## Dependencies to Add

Add the actively-maintained YAML library with better error messages:

```bash
go get github.com/goccy/go-yaml
```

**Why this library:** Actively maintained (unlike archived yaml.v3), passes 60+ more YAML spec tests, provides line/column-annotated error messages - critical for debugging PLC configuration files.

---

## Common Pitfalls to Avoid

1. **Don't forget the blank import**. Without `_ "github.com/aott33/go-plc/internal/source/modbus"` in main.go, the protocol types won't register and you'll get "unknown source type" errors.

2. **Don't modify internal/config/ for new protocols**. Future protocols (OPC-UA, MQTT) should create new packages in `internal/source/opcua/`, `internal/source/mqtt/`, etc. - NOT modify config code.

3. **Don't validate in the Load function**. Call `cfg.Validate()` at the end of Load(), but keep parsing and validation as separate concerns.

4. **Don't forget byte order defaults**. Apply defaults in the Validate() method of each protocol config.

5. **Don't exit on first error**. The Config.Validate() method collects ALL errors before returning - follow this pattern in protocol Validate() methods too.

6. **Don't skip the Duration helper**. See proposal Section 4.6 for the complete Duration type implementation that handles YAML time strings.

---

## Testing Your Implementation

### Manual Testing Steps

1. **Build and run with valid config:**
   ```bash
   make build
   ./go-plc -config config.yaml
   ```
   Expected: "Configuration loaded successfully" log message

2. **Run with missing config flag:**
   ```bash
   ./go-plc
   ```
   Expected: Error message about missing config

3. **Run with invalid config:**
   Create a `bad-config.yaml` with errors (missing host, invalid type, etc.)
   ```bash
   ./go-plc -config bad-config.yaml
   ```
   Expected: Clear error messages for each problem

### Automated Testing

```bash
go test -v ./internal/config/...
go test -race ./internal/config/...
```

---

## Definition of Done Checklist

- [x] `internal/source/source.go` created with interface definitions
- [x] `internal/source/registry.go` created with Register/ParseConfig functions
- [ ] `internal/source/modbus/tcp.go` created with TCPConfig and init() registration
- [ ] `internal/source/modbus/rtu.go` created with RTUConfig and init() registration
- [ ] `internal/source/modbus/duration.go` created for time.Duration YAML parsing
- [ ] `internal/config/config.go` created with Config, Variable, Scale, Register types
- [ ] `internal/config/loader.go` created with Load() using registry pattern
- [ ] `cmd/go-plc/main.go` updated with `-config` flag and blank modbus import
- [ ] `config.yaml` example file created in project root
- [ ] Unit tests cover happy paths and error cases (including unknown source type)
- [ ] All tests pass with `go test -race ./...`
- [ ] `go vet ./...` passes with no warnings
- [ ] Logging uses slog with correct format
- [ ] Error messages follow `[Component] - [Description] (context)` format

---

## References

- [Architecture: Configuration Schema](docs/architecture.md#data-architecture) - The authoritative schema specification
- [Architecture: Naming Patterns](docs/architecture.md#naming-patterns) - Naming conventions
- [Architecture: Error Message Format](docs/architecture.md#format-patterns) - Error formatting rules
- [yaml.v3 Documentation](https://pkg.go.dev/gopkg.in/yaml.v3) - YAML library reference
- [Previous Story 1.1](docs/sprint-artifacts/1-1-project-initialization-structure.md) - Project structure context

---

## Next Steps After This Story

Once this story is complete, Story 1.3 (Variable Store Implementation) will use these config structs to initialize the variable store. The `Variable` and `ScaleConfig` structs you create here will be passed to the variable store for registration.

---

## Dev Agent Record

### Agent Model Used
Claude Opus 4.5 (claude-opus-4-5-20251101)

### Completion Notes
- Story created by create-story workflow
- Optimized for human developer implementation (detailed steps, not code)
- All architecture patterns incorporated from docs/architecture.md
- Previous story (1.1) context included

### Context References
- docs/architecture.md - Primary technical reference
- docs/epics.md - Story requirements and acceptance criteria
- docs/sprint-artifacts/1-1-project-initialization-structure.md - Previous story context
