package source

import (
	"fmt"
	"strings"
)

// errUnknownSourceType returns an error when a source type is not registered in the factory registry.
// The registered types are listed so the user can see what is available.
func errUnknownSourceType(typeName, sourceName string, registered []string) error {
	return fmt.Errorf("[config] - Unknown source type '%s' (source: %s, registered: %s)",
		typeName, sourceName, strings.Join(registered, ", "))
}

// errDuplicateRegistration returns a string when a source type is registered multiple times.
// This indicates a programming error where two packages attempt to register the same type name.
func errDuplicateRegistration(typeName string) string {
	return fmt.Sprintf("[source] - Source type '%s' already registered", typeName)
}
