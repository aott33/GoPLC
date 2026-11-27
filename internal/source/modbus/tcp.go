// Package modbus provides the configuration for Modbus TCP sources and
// registers it with the source registry.
package modbus

import (
	"errors"
	"time"

	"github.com/aott33/go-plc/internal/source"
	"github.com/aott33/go-plc/internal/types"
	"github.com/goccy/go-yaml"
	"github.com/goccy/go-yaml/ast"
)

// TypeName is the source type used in configuration files.
const TypeName = "modbus-tcp"

func init() {
	source.Register(TypeName, parseTCPConfig)
}

// TCPConfig is the configuration for a Modbus TCP source.
type TCPConfig struct {
	name          string         // from parent source
	Host          string         `yaml:"host"`
	Port          int            `yaml:"port"`
	UnitID        uint8          `yaml:"unitId"`
	Timeout       types.Duration `yaml:"timeout"`
	PollInterval  types.Duration `yaml:"pollInterval"`
	RetryInterval types.Duration `yaml:"retryInterval"`
	ByteOrder     string         `yaml:"byteOrder"`
	WordOrder     string         `yaml:"wordOrder"`
}

// Defaults applied to fields that are absent from the config block.
const (
	DefaultPort          int           = 502
	DefaultUnitID        uint8         = 1
	DefaultTimeout       time.Duration = 1 * time.Second
	DefaultPollInterval  time.Duration = 1000 * time.Millisecond
	DefaultRetryInterval time.Duration = 10 * time.Second
	DefaultByteOrder     string        = ByteOrderBigEndian
	DefaultWordOrder     string        = WordOrderHighFirst
)

// Allowed values for ByteOrder and WordOrder.
const (
	ByteOrderBigEndian    = "big-endian"
	ByteOrderLittleEndian = "little-endian"
	WordOrderHighFirst    = "high-word-first"
	WordOrderLowFirst     = "low-word-first"
)

// SourceType returns the protocol identifier.
func (cfg *TCPConfig) SourceType() string {
	return TypeName
}

// SourceName returns the source name.
func (cfg *TCPConfig) SourceName() string {
	return cfg.name
}

// Validate checks that the TCPConfig has usable values.
// It returns all problems found, joined with errors.Join.
func (cfg *TCPConfig) Validate() error {
	var errs []error

	if cfg.Host == "" {
		errs = append(errs, errMissingField(cfg.name, "host"))
	}
	if cfg.Port < 1 || cfg.Port > 65535 {
		errs = append(errs, errInvalidField(cfg.name, "port", cfg.Port, "1 to 65535"))
	}
	if cfg.Timeout <= 0 {
		errs = append(errs, errInvalidField(cfg.name, "timeout", cfg.Timeout.Duration(), "greater than 0"))
	}
	if cfg.PollInterval <= 0 {
		errs = append(errs, errInvalidField(cfg.name, "pollInterval", cfg.PollInterval.Duration(), "greater than 0"))
	}
	if cfg.RetryInterval <= 0 {
		errs = append(errs, errInvalidField(cfg.name, "retryInterval", cfg.RetryInterval.Duration(), "greater than 0"))
	}
	if cfg.ByteOrder != ByteOrderBigEndian && cfg.ByteOrder != ByteOrderLittleEndian {
		errs = append(errs, errInvalidField(cfg.name, "byteOrder", cfg.ByteOrder,
			ByteOrderBigEndian+", "+ByteOrderLittleEndian))
	}
	if cfg.WordOrder != WordOrderHighFirst && cfg.WordOrder != WordOrderLowFirst {
		errs = append(errs, errInvalidField(cfg.name, "wordOrder", cfg.WordOrder,
			WordOrderHighFirst+", "+WordOrderLowFirst))
	}

	return errors.Join(errs...)
}

// defaultTCPConfig returns a TCPConfig with every default applied.
func defaultTCPConfig(name string) TCPConfig {
	return TCPConfig{
		name:          name,
		Port:          DefaultPort,
		UnitID:        DefaultUnitID,
		Timeout:       types.Duration(DefaultTimeout),
		PollInterval:  types.Duration(DefaultPollInterval),
		RetryInterval: types.Duration(DefaultRetryInterval),
		ByteOrder:     DefaultByteOrder,
		WordOrder:     DefaultWordOrder,
	}
}

// parseTCPConfig is the factory function registered with the source registry.
// It starts from the defaults and decodes the config block over them, so a
// field that is present in YAML (even as zero) overrides its default.
// Unknown fields are rejected to catch typos. It does not call Validate.
func parseTCPConfig(name string, configNode ast.Node) (source.SourceConfig, error) {
	cfg := defaultTCPConfig(name)

	if configNode == nil || configNode.Type() == ast.NullType {
		return &cfg, nil
	}

	if err := yaml.NodeToValue(configNode, &cfg, yaml.Strict()); err != nil {
		return nil, errDecode(name, err)
	}

	return &cfg, nil
}
