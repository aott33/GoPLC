// Package source provides a registry-based system for protocol-specific data sources.
// This package defines the core interfaces and factory pattern for extensible protocol support.
//
// The registry pattern allows protocols (Modbus TCP/RTU, OPC-UA, etc.) to self-register
// at startup via init() functions, enabling zero-touch addition of new protocols without
// modifying core configuration code.
package source

import "github.com/goccy/go-yaml/ast"

// Source represents a data source that can be connected and polled.
type Source interface {
	// Name returns the user-defined source name from configuration
	Name() string

	// Type returns the protocol identifier (e.g., "modbus-tcp", "modbus-rtu", "opcua-client")
	Type() string
}

// SourceConfig represents protocol-specific configuration.
// Each protocol implements this interface.
type SourceConfig interface {
	// Validate checks required fields, applies defaults, and validates field values.
	Validate() error

	// SourceType returns the protocol type string (e.g., "modbus-tcp").
	// Must match the type name used during factory registration.
	SourceType() string

	// SourceName returns the user-defined source name from configuration.
	SourceName() string
}

// SourceFactory creates a SourceConfig from YAML node data.
// The factory is responsible for parsing protocol-specific fields.
type SourceFactory func(name string, configNode ast.Node) (SourceConfig, error)
