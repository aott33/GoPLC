package modbus

import "fmt"

// errDecode returns an error when the config block cannot be decoded.
func errDecode(sourceName string, err error) error {
	return fmt.Errorf("[config] - Invalid modbus-tcp config (source: %s): %w", sourceName, err)
}

// errMissingField returns an error when a required field is not set.
func errMissingField(sourceName, field string) error {
	return fmt.Errorf("[config] - Missing required field '%s' (source: %s, type: modbus-tcp)", field, sourceName)
}

// errInvalidField returns an error when a field has a value outside its allowed range or set.
func errInvalidField(sourceName, field string, value any, allowed string) error {
	return fmt.Errorf("[config] - Invalid value for '%s' (source: %s, type: modbus-tcp, value: %v, allowed: %s)",
		field, sourceName, value, allowed)
}
