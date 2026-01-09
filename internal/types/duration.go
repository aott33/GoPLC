// Package types provides common type definitions used across the application.
package types

import (
	"time"

	"github.com/goccy/go-yaml"
	"github.com/goccy/go-yaml/ast"
)

// Duration wraps time.Duration to enable YAML unmarshaling of duration strings.
// Supports standard Go duration formats: "300ms", "1.5h", "2h45m", etc.
type Duration time.Duration

// UnmarshalYAML parses duration strings from YAML.
func (d *Duration) UnmarshalYAML(node ast.Node) error {
	var s string
	if err := yaml.NodeToValue(node, &s); err != nil {
		return err
	}

	duration, err := time.ParseDuration(s)
	if err != nil {
		return err
	}

	*d = Duration(duration)
	return nil
}

// MarshalYAML converts Duration back to a duration string for YAML output.
func (d Duration) MarshalYAML() ([]byte, error) {
	return []byte(time.Duration(d).String()), nil
}

// Duration returns the underlying time.Duration value.
func (d Duration) Duration() time.Duration {
	return time.Duration(d)
}
