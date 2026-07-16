package modbus

import (
	"time"

	"github.com/aott33/go-plc/internal/source"
	"github.com/goccy/go-yaml/ast"
)

func init() {
	source.Register("modbus-tcp", parseTCPConfig)
}

type TCPConfig struct {
	name          string        // from parent source
	Host          string        `yaml:"host"`
	Port          int           `yaml:"port"`
	UnitID        uint8         `yaml:"unitId"`
	Timeout       time.Duration `yaml:"timeout"`
	PollInterval  time.Duration `yaml:"pollInterval"`
	RetryInterval time.Duration `yaml:"retryInterval"`
	ByteOrder     string        `yaml:"byteOrder"`
	WordOrder     string        `yaml:"wordOrder"`
}

const (
	DefaultPort          int           = 502
	DefaultUnitID        uint8         = 1
	DefaultTimout        time.Duration = 1 * time.Second
	DefaultPollInterval  time.Duration = 1000 * time.Millisecond
	DefaultRetryInterval time.Duration = 10 * time.Second
	DefaultByteOrder     string        = "big-endian"
	DefaultWordOrder     string        = "high-word-first"
)

// SourceType returns the protocol identifier.
func (cfg *TCPConfig) SourceType() string {
	return "modbus-tcp"
}

// SourceName returns the source name.
func (cfg *TCPConfig) SourceName() string {
	return cfg.name
}

// Validate checks that the TCPConfig has the correct values
func (cfg *TCPConfig) Validate() error {
	return nil
}

// parseTCPConfig is the factory function registered with the source registry.
func parseTCPConfig(name string, configNode ast.Node) (source.SourceConfig, error) {
	var cfg TCPConfig

	return &cfg, nil
}
