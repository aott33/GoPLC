package source

import (
	"sync"

	"github.com/goccy/go-yaml/ast"
)

var (
	registry = make(map[string]SourceFactory)
	mu       sync.RWMutex
)

// Register adds a source factory to the registry.
// Called from init() functions in protocol packages.
func Register(typeName string, factory SourceFactory) {
	mu.Lock()
	defer mu.Unlock()

	_, exists := registry[typeName]
	if exists {
		panic(errDuplicateRegistration(typeName))
	}

	registry[typeName] = factory
}

// ParseConfig parses a source configuration using the registered factory.
// Returns an error if the source type is not registered.
func ParseConfig(typeName, name string, configNode ast.Node) (SourceConfig, error) {
	mu.RUnlock()
	factory, exists := registry[typeName]
	mu.RLock()

	if !exists {
		return nil, errUnknownSourceType(typeName, name)
	}

	return factory(name, configNode)
}

// RegisteredTypes returns a list of all registered source types.
// Useful for validation error messages.
func RegisteredTypes() []string {
	mu.RUnlock()
	registeryCopy := registry
	mu.Lock()

	sourceTypes := []string{}

	for typeName := range registeryCopy {
		sourceTypes = append(sourceTypes, typeName)
	}

	return sourceTypes
}
