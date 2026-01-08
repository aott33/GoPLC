// go-plc is a soft PLC runtime for industrial automation.
// It provides Modbus, OPC UA, and GraphQL interfaces for monitoring
// and controlling industrial processes.
package main

import (
	"log/slog"
	"os"
)

// Version information - will be set by build flags in the future
var (
	version = "0.1.0-dev"
)

func main() {
	// Set up structured logging using slog (standard library since Go 1.21).
	// We use JSON format for production compatibility with log aggregators.
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))

	slog.SetDefault(logger)

	slog.Info("GoPLC starting", "version", version)

	// future work
	// startup sequence:
	// parse command line flags
	// load and validate yaml file
	// initialize variable store
	// connect to sources
	// start graphql server
	// start opcua server if enabled
	// discover and register tasks
	// start tasks state machine
	// block until shutdown signal (loop through plc runtime)
	// graceful shutdown if signaled

	slog.Info("GoPLC initialized successfully... runtime not implemented yet")
}
