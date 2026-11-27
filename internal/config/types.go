package config

import "github.com/aott33/go-plc/internal/source"

// Config is the parsed and validated application configuration.
type Config struct {
	LogLevel  string
	Sources   []Source
	Variables []Variable
}

// Source is one configured data source. Config holds the protocol-specific
// configuration returned by the source registry.
type Source struct {
	Name   string
	Type   string
	Config source.SourceConfig
}

// Variable is a named value read from a source.
type Variable struct {
	Name     string   `yaml:"name"`
	Source   string   `yaml:"source"`
	DataType string   `yaml:"dataType"`
	Tags     []string `yaml:"tags"`
}
