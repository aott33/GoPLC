# Sprint Change Proposal: Story 1.2 Architecture Enhancement

**Date:** 2025-12-20
**Status:** Proposed
**Triggered By:** Story 1.2 - Configuration Schema & Loading (pre-development review)
**Change Scope:** Minor - Direct implementation within current story

---

## 1. Issue Summary

### Problem Statement

Before beginning development on Story 1.2 (Configuration Schema & Loading), an architectural review identified two improvements that should be incorporated:

1. **Unmaintained YAML Library:** The currently specified `gopkg.in/yaml.v3` library is now archived and unmaintained. Continuing with this dependency creates technical debt from day one.

2. **Tight Coupling in Source Configuration:** The current story design places `ModbusTCPConfig` and `ModbusRTUConfig` structs directly in `internal/config/`. This creates tight coupling where every new protocol (OPC-UA client, MQTT, EtherNet/IP, REST, PROFINET) would require modifying core config code.

### Discovery Context

This was identified during pre-implementation review of Story 1.2, before any code was written. This is the optimal time to make architectural adjustments - zero rework cost.

### Evidence

1. **YAML Library Status:** `gopkg.in/yaml.v3` repository is archived. The `github.com/goccy/go-yaml` library is actively maintained, passes 60+ more YAML specification test cases, and provides source-annotated error messages with line/column information - critical for debugging PLC configuration files.

2. **Industry Precedent:** Production systems like Ignition (device drivers) and Prometheus (service discovery) use registry patterns for protocol extensibility. This is a proven approach for systems that need to support multiple protocols.

---

## 2. Impact Analysis

### Epic Impact

| Epic | Impact | Notes |
|------|--------|-------|
| Epic 1 (Foundation) | Neutral | Story 1.2 scope unchanged, implementation approach improved |
| Epic 2 (Modbus) | Positive | Modbus config already in dedicated package |
| Epic 5 (OPC UA) | Positive | Pattern established for adding OPC UA source |
| Future Protocols | Positive | Zero-touch addition of new protocols |

### Artifact Impact

| Artifact | Changes Required |
|----------|------------------|
| Story 1.2 | Update implementation guide with new structure |
| Architecture Doc | Update directory structure and component boundaries |
| PRD | None - user-facing behavior unchanged |
| UX Specification | None - internal implementation detail |

### Acceptance Criteria Impact

**No changes to Story 1.2 acceptance criteria.** All three ACs remain exactly as written:
- AC1: YAML parsing with type-discriminated sources
- AC2: Validation with clear error messages
- AC3: Successful configuration loading

The registry pattern is simply a better way to achieve these same acceptance criteria.

---

## 3. Recommended Approach

### Selected Path: Direct Adjustment

Modify Story 1.2 implementation guide before development begins. This is the lowest-risk, highest-value option because:

1. **Zero rework** - No existing code to modify
2. **Optimal timing** - Changes are easiest before implementation
3. **Long-term benefit** - Every future protocol addition is simplified
4. **Maintains schedule** - Adds minimal preparation time, saves time on future stories

### Change 1: YAML Library Switch

| Aspect | Before | After |
|--------|--------|-------|
| Library | `gopkg.in/yaml.v3` | `github.com/goccy/go-yaml` |
| Status | Archived | Actively maintained |
| Error messages | Basic | Line/column annotated |
| YAML compliance | Standard | 60+ more spec tests passing |

**Developer Action:** Use `go get github.com/goccy/go-yaml` instead of yaml.v3. The API is largely compatible - primary difference is import path.

### Change 2: Registry Pattern for Protocol Sources

**Current Story 1.2 Structure (from story document):**
```
internal/config/
├── doc.go
├── types.go        # Config, Source, ModbusTCPConfig, ModbusRTUConfig, Variable...
├── loader.go       # Load() with custom UnmarshalYAML for sources
└── config_test.go
```

**Proposed Structure:**
```
internal/
├── config/
│   ├── config.go       # Config struct, Variable, Scale, Register types
│   ├── loader.go       # Load() - delegates source parsing to registry
│   └── config_test.go
├── source/
│   ├── source.go       # Source interface, SourceConfig interface
│   ├── registry.go     # Register(), ParseConfig(), RegisteredTypes()
│   └── modbus/
│       ├── tcp.go      # TCPConfig struct + init() registration
│       └── rtu.go      # RTUConfig struct + init() registration
```

---

## 4. Technical Specification for Human Developer

This section provides the complete technical details needed to implement the registry pattern.

### 4.1 Key Go Concepts

**Go Interfaces:**
An interface defines a contract - a set of method signatures. Any struct that implements all the methods automatically satisfies the interface. No explicit "implements" keyword needed.

```go
// Any struct with these methods satisfies Source interface
type Source interface {
    Connect(ctx context.Context) error
    Close() error
    // ... other methods
}
```

**Registry Pattern:**
A map from string keys to factory functions. Protocols register themselves at startup, and the config loader looks up the factory by type name.

```go
var registry = make(map[string]SourceFactory)

func Register(typeName string, factory SourceFactory) {
    registry[typeName] = factory
}
```

**Blank Imports:**
The `_ "package/path"` syntax imports a package solely for its `init()` side effects (registration). The underscore discards the package reference since we don't call it directly.

```go
import (
    _ "github.com/aott33/go-plc/internal/source/modbus" // triggers init()
)
```

**Deferred YAML Parsing:**
Parse the source `type` field first, look up the factory in the registry, then let the factory parse the `config` block into the correct typed struct.

### 4.2 Interface Definitions

**File:** `internal/source/source.go`

```go
package source

import (
    "context"

    "github.com/goccy/go-yaml"
)

// Source represents a data source that can be connected and polled.
// This interface will be fully implemented in Epic 2 (Modbus I/O).
// For Story 1.2, we only need the config-related methods.
type Source interface {
    // Name returns the configured source name
    Name() string

    // Type returns the protocol type (e.g., "modbus-tcp")
    Type() string

    // Future methods (Epic 2):
    // Connect(ctx context.Context) error
    // Close() error
    // Poll(ctx context.Context) error
}

// SourceConfig represents protocol-specific configuration.
// Each protocol implements this interface.
type SourceConfig interface {
    // Validate checks the configuration and returns errors
    Validate() error

    // SourceType returns the protocol type string
    SourceType() string

    // SourceName returns the configured name
    SourceName() string
}

// SourceFactory creates a SourceConfig from YAML node data.
// The factory is responsible for parsing protocol-specific fields.
type SourceFactory func(name string, configNode yaml.Node) (SourceConfig, error)
```

### 4.3 Registry Implementation

**File:** `internal/source/registry.go`

```go
package source

import (
    "fmt"
    "sync"

    "github.com/goccy/go-yaml"
)

var (
    registry = make(map[string]SourceFactory)
    mu       sync.RWMutex
)

// Register adds a source factory to the registry.
// Called from init() functions in protocol packages.
func Register(typeName string, factory SourceFactory) {
    mu.Lock()
    defer mu.Unlock()

    if _, exists := registry[typeName]; exists {
        panic(fmt.Sprintf("source type already registered: %s", typeName))
    }
    registry[typeName] = factory
}

// ParseConfig parses a source configuration using the registered factory.
// Returns an error if the source type is not registered.
func ParseConfig(typeName, name string, configNode yaml.Node) (SourceConfig, error) {
    mu.RLock()
    factory, exists := registry[typeName]
    mu.RUnlock()

    if !exists {
        return nil, fmt.Errorf("[config] - Unknown source type '%s' (source: %s)", typeName, name)
    }

    return factory(name, configNode)
}

// RegisteredTypes returns a list of all registered source types.
// Useful for validation error messages.
func RegisteredTypes() []string {
    mu.RLock()
    defer mu.RUnlock()

    types := make([]string, 0, len(registry))
    for t := range registry {
        types = append(types, t)
    }
    return types
}
```

### 4.4 Modbus TCP Implementation

**File:** `internal/source/modbus/tcp.go`

```go
package modbus

import (
    "fmt"
    "time"

    "github.com/goccy/go-yaml"

    "github.com/aott33/go-plc/internal/source"
)

func init() {
    source.Register("modbus-tcp", parseTCPConfig)
}

// TCPConfig holds Modbus TCP connection settings.
type TCPConfig struct {
    name          string        // from parent source, not YAML
    Host          string        `yaml:"host"`
    Port          int           `yaml:"port"`
    UnitID        uint8         `yaml:"unitId"`
    Timeout       Duration      `yaml:"timeout"`
    PollInterval  Duration      `yaml:"pollInterval"`
    RetryInterval Duration      `yaml:"retryInterval"`
    ByteOrder     string        `yaml:"byteOrder"`
    WordOrder     string        `yaml:"wordOrder"`
}

// SourceType returns the protocol identifier.
func (c *TCPConfig) SourceType() string { return "modbus-tcp" }

// SourceName returns the configured source name.
func (c *TCPConfig) SourceName() string { return c.name }

// Validate checks required fields and applies defaults.
func (c *TCPConfig) Validate() error {
    if c.Host == "" {
        return fmt.Errorf("[config] - Missing required field 'host' (source: %s, type: modbus-tcp)", c.name)
    }

    // Apply defaults
    if c.Port == 0 {
        c.Port = 502
    }
    if c.UnitID == 0 {
        c.UnitID = 1
    }
    if c.Timeout == 0 {
        c.Timeout = Duration(1 * time.Second)
    }
    if c.PollInterval == 0 {
        c.PollInterval = Duration(100 * time.Millisecond)
    }
    if c.RetryInterval == 0 {
        c.RetryInterval = Duration(5 * time.Second)
    }
    if c.ByteOrder == "" {
        c.ByteOrder = "big-endian"
    }
    if c.WordOrder == "" {
        c.WordOrder = "high-word-first"
    }

    // Validate enum values
    if c.ByteOrder != "big-endian" && c.ByteOrder != "little-endian" {
        return fmt.Errorf("[config] - Invalid byteOrder '%s', must be 'big-endian' or 'little-endian' (source: %s)", c.ByteOrder, c.name)
    }
    if c.WordOrder != "high-word-first" && c.WordOrder != "low-word-first" {
        return fmt.Errorf("[config] - Invalid wordOrder '%s', must be 'high-word-first' or 'low-word-first' (source: %s)", c.WordOrder, c.name)
    }

    return nil
}

// parseTCPConfig is the factory function registered with the source registry.
func parseTCPConfig(name string, configNode yaml.Node) (source.SourceConfig, error) {
    var cfg TCPConfig
    cfg.name = name

    if err := configNode.Decode(&cfg); err != nil {
        return nil, fmt.Errorf("[config] - Failed to parse modbus-tcp config (source: %s): %w", name, err)
    }

    return &cfg, nil
}
```

### 4.5 Modbus RTU Implementation

**File:** `internal/source/modbus/rtu.go`

```go
package modbus

import (
    "fmt"
    "time"

    "github.com/goccy/go-yaml"

    "github.com/aott33/go-plc/internal/source"
)

func init() {
    source.Register("modbus-rtu", parseRTUConfig)
}

// RTUConfig holds Modbus RTU (serial) connection settings.
type RTUConfig struct {
    name         string   // from parent source, not YAML
    Device       string   `yaml:"device"`
    BaudRate     int      `yaml:"baudRate"`
    DataBits     int      `yaml:"dataBits"`
    Parity       string   `yaml:"parity"`
    StopBits     int      `yaml:"stopBits"`
    UnitID       uint8    `yaml:"unitId"`
    Timeout      Duration `yaml:"timeout"`
    PollInterval Duration `yaml:"pollInterval"`
}

// SourceType returns the protocol identifier.
func (c *RTUConfig) SourceType() string { return "modbus-rtu" }

// SourceName returns the configured source name.
func (c *RTUConfig) SourceName() string { return c.name }

// Validate checks required fields and applies defaults.
func (c *RTUConfig) Validate() error {
    if c.Device == "" {
        return fmt.Errorf("[config] - Missing required field 'device' (source: %s, type: modbus-rtu)", c.name)
    }
    if c.BaudRate == 0 {
        return fmt.Errorf("[config] - Missing required field 'baudRate' (source: %s, type: modbus-rtu)", c.name)
    }

    // Apply defaults
    if c.DataBits == 0 {
        c.DataBits = 8
    }
    if c.Parity == "" {
        c.Parity = "none"
    }
    if c.StopBits == 0 {
        c.StopBits = 1
    }
    if c.UnitID == 0 {
        c.UnitID = 1
    }
    if c.Timeout == 0 {
        c.Timeout = Duration(1 * time.Second)
    }
    if c.PollInterval == 0 {
        c.PollInterval = Duration(100 * time.Millisecond)
    }

    // Validate parity
    if c.Parity != "none" && c.Parity != "even" && c.Parity != "odd" {
        return fmt.Errorf("[config] - Invalid parity '%s', must be 'none', 'even', or 'odd' (source: %s)", c.Parity, c.name)
    }

    return nil
}

// parseRTUConfig is the factory function registered with the source registry.
func parseRTUConfig(name string, configNode yaml.Node) (source.SourceConfig, error) {
    var cfg RTUConfig
    cfg.name = name

    if err := configNode.Decode(&cfg); err != nil {
        return nil, fmt.Errorf("[config] - Failed to parse modbus-rtu config (source: %s): %w", name, err)
    }

    return &cfg, nil
}
```

### 4.6 Duration Type Helper

**File:** `internal/source/modbus/duration.go`

```go
package modbus

import (
    "time"

    "github.com/goccy/go-yaml"
)

// Duration wraps time.Duration for YAML unmarshaling.
// YAML strings like "100ms" or "5s" need custom parsing.
type Duration time.Duration

// UnmarshalYAML parses duration strings like "100ms", "5s", "1m".
func (d *Duration) UnmarshalYAML(node *yaml.Node) error {
    var s string
    if err := node.Decode(&s); err != nil {
        return err
    }

    duration, err := time.ParseDuration(s)
    if err != nil {
        return err
    }

    *d = Duration(duration)
    return nil
}

// Duration returns the underlying time.Duration value.
func (d Duration) Duration() time.Duration {
    return time.Duration(d)
}
```

### 4.7 Updated Config Loader

**File:** `internal/config/loader.go`

```go
package config

import (
    "fmt"
    "log/slog"
    "os"

    "github.com/goccy/go-yaml"

    "github.com/aott33/go-plc/internal/source"
)

// rawSource is used for initial YAML parsing before type discrimination.
type rawSource struct {
    Name   string    `yaml:"name"`
    Type   string    `yaml:"type"`
    Config yaml.Node `yaml:"config"`
}

// rawConfig is the top-level YAML structure for initial parsing.
type rawConfig struct {
    LogLevel  string      `yaml:"logLevel"`
    Sources   []rawSource `yaml:"sources"`
    Variables []Variable  `yaml:"variables"`
}

// Load reads and parses the configuration file.
func Load(path string) (*Config, error) {
    data, err := os.ReadFile(path)
    if err != nil {
        return nil, fmt.Errorf("[config] - Failed to read config file (path: %s): %w", path, err)
    }

    var raw rawConfig
    if err := yaml.Unmarshal(data, &raw); err != nil {
        return nil, fmt.Errorf("[config] - Failed to parse YAML: %w", err)
    }

    cfg := &Config{
        LogLevel:  raw.LogLevel,
        Variables: raw.Variables,
    }

    // Apply log level default
    if cfg.LogLevel == "" {
        cfg.LogLevel = "info"
    }

    // Parse sources using registry
    for _, rs := range raw.Sources {
        if rs.Name == "" {
            return nil, fmt.Errorf("[config] - Source missing required field 'name'")
        }
        if rs.Type == "" {
            return nil, fmt.Errorf("[config] - Missing required field 'type' (source: %s)", rs.Name)
        }

        srcConfig, err := source.ParseConfig(rs.Type, rs.Name, rs.Config)
        if err != nil {
            return nil, err
        }

        cfg.Sources = append(cfg.Sources, Source{
            Name:   rs.Name,
            Type:   rs.Type,
            Config: srcConfig,
        })
    }

    // Validate the complete configuration
    if err := cfg.Validate(); err != nil {
        return nil, err
    }

    slog.Info("Configuration loaded successfully",
        "sources", len(cfg.Sources),
        "variables", len(cfg.Variables),
        "logLevel", cfg.LogLevel,
    )

    return cfg, nil
}
```

### 4.8 Updated Config Types

**File:** `internal/config/config.go`

```go
package config

import (
    "fmt"

    "github.com/aott33/go-plc/internal/source"
)

// Config is the root configuration structure.
type Config struct {
    LogLevel  string     `yaml:"logLevel"`
    Sources   []Source   `yaml:"-"` // Parsed via registry, not direct YAML
    Variables []Variable `yaml:"variables"`
}

// Source holds a parsed source configuration.
type Source struct {
    Name   string
    Type   string
    Config source.SourceConfig
}

// Variable defines a process variable bound to a source.
type Variable struct {
    Name     string          `yaml:"name"`
    Source   string          `yaml:"source"`
    Register RegisterConfig  `yaml:"register"`
    DataType string          `yaml:"dataType"`
    Scale    *ScaleConfig    `yaml:"scale,omitempty"`
    Tags     []string        `yaml:"tags,omitempty"`
}

// RegisterConfig defines Modbus register addressing.
type RegisterConfig struct {
    Type    string `yaml:"type"`
    Address uint16 `yaml:"address"`
}

// ScaleConfig defines linear scaling parameters.
type ScaleConfig struct {
    RawMin float64 `yaml:"rawMin"`
    RawMax float64 `yaml:"rawMax"`
    EngMin float64 `yaml:"engMin"`
    EngMax float64 `yaml:"engMax"`
    Unit   string  `yaml:"unit"`
}

// Validate checks the configuration for errors.
// Collects all errors rather than failing on first.
func (c *Config) Validate() error {
    var errors []error

    // Validate log level
    validLevels := map[string]bool{"debug": true, "info": true, "warn": true, "error": true}
    if !validLevels[c.LogLevel] {
        errors = append(errors, fmt.Errorf("[config] - Invalid logLevel '%s', must be debug/info/warn/error", c.LogLevel))
    }

    // Validate sources (call each source's Validate method)
    sourceNames := make(map[string]bool)
    for _, src := range c.Sources {
        if sourceNames[src.Name] {
            errors = append(errors, fmt.Errorf("[config] - Duplicate source name '%s'", src.Name))
        }
        sourceNames[src.Name] = true

        if err := src.Config.Validate(); err != nil {
            errors = append(errors, err)
        }
    }

    // Validate variables
    varNames := make(map[string]bool)
    for _, v := range c.Variables {
        if v.Name == "" {
            errors = append(errors, fmt.Errorf("[config] - Variable missing required field 'name'"))
            continue
        }

        if varNames[v.Name] {
            errors = append(errors, fmt.Errorf("[config] - Duplicate variable name '%s'", v.Name))
        }
        varNames[v.Name] = true

        // Check source reference
        if !sourceNames[v.Source] {
            errors = append(errors, fmt.Errorf("[config] - Unknown source reference (variable: %s, source: %s)", v.Name, v.Source))
        }

        // Validate data type
        if err := validateDataType(v); err != nil {
            errors = append(errors, err)
        }

        // Validate register type
        if err := validateRegisterType(v); err != nil {
            errors = append(errors, err)
        }

        // Validate scaling if present
        if v.Scale != nil {
            if err := validateScale(v); err != nil {
                errors = append(errors, err)
            }
        }
    }

    if len(errors) > 0 {
        // Combine errors into a single message
        msg := fmt.Sprintf("[config] - %d validation error(s):", len(errors))
        for _, e := range errors {
            msg += "\n  - " + e.Error()
        }
        return fmt.Errorf(msg)
    }

    return nil
}

// validateDataType checks that dataType is valid.
func validateDataType(v Variable) error {
    validTypes := map[string]bool{
        "bool": true, "uint16": true, "int16": true,
        "uint32": true, "int32": true, "float32": true,
        "uint64": true, "int64": true, "float64": true,
    }

    if !validTypes[v.DataType] {
        return fmt.Errorf("[config] - Invalid dataType '%s' (variable: %s)", v.DataType, v.Name)
    }
    return nil
}

// validateRegisterType checks register type validity and compatibility.
func validateRegisterType(v Variable) error {
    validTypes := map[string]bool{
        "holding": true, "input": true, "coil": true, "discrete": true,
    }

    if !validTypes[v.Register.Type] {
        return fmt.Errorf("[config] - Invalid register type '%s' (variable: %s)", v.Register.Type, v.Name)
    }

    // Bool can only be used with coil or discrete
    if v.DataType == "bool" && v.Register.Type != "coil" && v.Register.Type != "discrete" {
        return fmt.Errorf("[config] - dataType 'bool' requires register type 'coil' or 'discrete' (variable: %s)", v.Name)
    }

    // Multi-register types should only be used with holding or input
    multiRegTypes := map[string]bool{"uint32": true, "int32": true, "float32": true, "uint64": true, "int64": true, "float64": true}
    if multiRegTypes[v.DataType] && v.Register.Type != "holding" && v.Register.Type != "input" {
        return fmt.Errorf("[config] - dataType '%s' requires register type 'holding' or 'input' (variable: %s)", v.DataType, v.Name)
    }

    return nil
}

// validateScale checks scaling configuration.
func validateScale(v Variable) error {
    s := v.Scale
    if s.RawMin >= s.RawMax {
        return fmt.Errorf("[config] - rawMin must be less than rawMax (variable: %s)", v.Name)
    }
    return nil
}
```

### 4.9 Main.go Blank Import

**File:** `cmd/go-plc/main.go` (add import)

```go
import (
    "flag"
    "log/slog"
    "os"

    "github.com/aott33/go-plc/internal/config"

    // Register source types - triggers init() for self-registration
    _ "github.com/aott33/go-plc/internal/source/modbus"
)
```

---

## 5. Updated File Structure

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
go.mod                  # Updated: github.com/goccy/go-yaml instead of yaml.v3
```

---

## 6. Adding Future Protocols

This section demonstrates why the registry pattern pays off. To add a new protocol (e.g., OPC-UA client source):

**Step 1:** Create protocol package

```
internal/source/opcua/
└── client.go
```

**Step 2:** Implement interfaces and register

```go
package opcua

import (
    "github.com/aott33/go-plc/internal/source"
)

func init() {
    source.Register("opcua-client", parseClientConfig)
}

type ClientConfig struct {
    name     string
    Endpoint string `yaml:"endpoint"`
    // ... other fields
}

func (c *ClientConfig) SourceType() string { return "opcua-client" }
func (c *ClientConfig) SourceName() string { return c.name }
func (c *ClientConfig) Validate() error { /* ... */ }

func parseClientConfig(name string, node yaml.Node) (source.SourceConfig, error) {
    // ...
}
```

**Step 3:** Add blank import to main.go

```go
import (
    _ "github.com/aott33/go-plc/internal/source/modbus"
    _ "github.com/aott33/go-plc/internal/source/opcua"  // NEW
)
```

**That's it.** No changes to `internal/config/` or `internal/source/registry.go`.

---

## 7. Definition of Done (Updated)

Story 1.2 completion checklist with registry pattern:

- [ ] `internal/source/source.go` created with interface definitions
- [ ] `internal/source/registry.go` created with Register/ParseConfig functions
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

## 8. Implementation Handoff

### Change Scope Classification: Minor

This change can be implemented directly by the development team within the existing Story 1.2 work.

### Handoff Recipients

| Role | Responsibility |
|------|----------------|
| Human Developer | Implement Story 1.2 using registry pattern per this specification |
| Scrum Master | Update Story 1.2 document to reference this proposal |

### Next Steps

1. **Update Story 1.2 document** to reference this Sprint Change Proposal
2. **Begin development** following the technical specification in Section 4
3. **No architecture document update needed yet** - can be updated after Story 1.2 is complete

### Success Criteria

- Story 1.2 acceptance criteria all pass (unchanged)
- Registry pattern correctly parses modbus-tcp and modbus-rtu sources
- Adding a mock "test-source" type in tests demonstrates extensibility
- Error messages include line/column from goccy/go-yaml where applicable

---

## Approval

**Proposed By:** Bob (Scrum Master Agent)
**Date:** 2025-12-20

**User Approval:** [ ] Approved / [ ] Rejected / [ ] Needs Revision

**Notes:**
_Space for user comments or conditions_

---

Generated with [Claude Code](https://claude.com/claude-code)
