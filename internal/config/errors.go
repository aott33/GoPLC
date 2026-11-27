package config

import "fmt"

// errReadFile returns an error when the config file cannot be read.
func errReadFile(path string, err error) error {
	return fmt.Errorf("[config] - Cannot read config file (path: %s): %w", path, err)
}

// errParseYAML returns an error when the config file is not valid YAML.
func errParseYAML(path string, err error) error {
	return fmt.Errorf("[config] - Invalid YAML (path: %s): %w", path, err)
}

// errMissingSourceField returns an error when a source entry lacks a required field.
func errMissingSourceField(index int, field string) error {
	return fmt.Errorf("[config] - Missing required field '%s' (source index: %d)", field, index)
}

// errDuplicateSource returns an error when two sources share a name.
func errDuplicateSource(name string) error {
	return fmt.Errorf("[config] - Duplicate source name (source: %s)", name)
}

// errInvalidLogLevel returns an error when logLevel is not a known level.
func errInvalidLogLevel(level string) error {
	return fmt.Errorf("[config] - Invalid value for 'logLevel' (value: %s, allowed: debug, info, warn, error)", level)
}
