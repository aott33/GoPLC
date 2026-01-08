package source

import "fmt"

// errUnknownSourceType returns an error when a source type is not registered in the factory registry.
func errUnknownSourceType(typeName, sourceName string) error {
	return fmt.Errorf("[config] - Unknown source type '%s' (source: %s)", typeName, sourceName)
}

// errDuplicateRegistration returns a string when a source type is registered multiple times.
// This indicates a programming error where two packages attempt to register the same type name.
func errDuplicateRegistration(typeName string) string {
	return fmt.Sprintf("[source] - Source type '%s' already registered", typeName)
}
