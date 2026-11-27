package config

import (
	"errors"
	"os"

	"github.com/aott33/go-plc/internal/source"
	"github.com/goccy/go-yaml"
	"github.com/goccy/go-yaml/ast"
)

// DefaultLogLevel is used when logLevel is not set in the file.
const DefaultLogLevel = "info"

// rawConfig mirrors the YAML file. Source config blocks are kept as
// AST nodes and handed to the source registry.
type rawConfig struct {
	LogLevel  string      `yaml:"logLevel"`
	Sources   []rawSource `yaml:"sources"`
	Variables []Variable  `yaml:"variables"`
}

type rawSource struct {
	Name   string  `yaml:"name"`
	Type   string  `yaml:"type"`
	Config rawNode `yaml:"config"`
}

// rawNode captures a YAML node without decoding it.
type rawNode struct {
	node ast.Node
}

// UnmarshalYAML stores the node for later decoding.
func (r *rawNode) UnmarshalYAML(node ast.Node) error {
	r.node = node
	return nil
}

// Load reads the YAML file at path, parses each source through the source
// registry and validates it. Source packages must be imported (usually with
// a blank import in main) so their types are registered before Load is called.
//
// Source errors are collected and returned together with errors.Join.
func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, errReadFile(path, err)
	}

	return parse(path, data)
}

// parse builds a Config from YAML bytes. path is used only in error messages.
func parse(path string, data []byte) (*Config, error) {
	var raw rawConfig
	if err := yaml.Unmarshal(data, &raw); err != nil {
		return nil, errParseYAML(path, err)
	}

	cfg := &Config{
		LogLevel:  raw.LogLevel,
		Variables: raw.Variables,
	}
	if cfg.LogLevel == "" {
		cfg.LogLevel = DefaultLogLevel
	}

	var errs []error

	switch cfg.LogLevel {
	case "debug", "info", "warn", "error":
	default:
		errs = append(errs, errInvalidLogLevel(cfg.LogLevel))
	}

	seen := make(map[string]bool, len(raw.Sources))

	for i, rs := range raw.Sources {
		if rs.Name == "" {
			errs = append(errs, errMissingSourceField(i, "name"))
			continue
		}
		if rs.Type == "" {
			errs = append(errs, errMissingSourceField(i, "type"))
			continue
		}
		if seen[rs.Name] {
			errs = append(errs, errDuplicateSource(rs.Name))
			continue
		}
		seen[rs.Name] = true

		sc, err := source.ParseConfig(rs.Type, rs.Name, rs.Config.node)
		if err != nil {
			errs = append(errs, err)
			continue
		}
		if err := sc.Validate(); err != nil {
			errs = append(errs, err)
			continue
		}

		cfg.Sources = append(cfg.Sources, Source{Name: rs.Name, Type: rs.Type, Config: sc})
	}

	if err := errors.Join(errs...); err != nil {
		return nil, err
	}

	return cfg, nil
}
