// go-plc is a soft PLC runtime for industrial automation.
// It provides Modbus, OPC UA, and GraphQL interfaces for monitoring
// and controlling industrial processes.
package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"io"
	"log/slog"
	"os"

	"github.com/aott33/go-plc/internal/config"

	// Source packages register their types with the source registry in init().
	_ "github.com/aott33/go-plc/internal/source/modbus"
)

// Version information - will be set by build flags in the future
var (
	version = "0.1.0-dev"
)

// DefaultConfigPath is used when -config is not given.
const DefaultConfigPath = "configs/example.yaml"

// options holds the parsed command-line flags.
type options struct {
	configPath string
	logLevel   string // empty means use the level from the config file
}

// parseFlags parses command-line arguments (without the program name).
func parseFlags(args []string, output io.Writer) (options, error) {
	var opts options

	fs := flag.NewFlagSet("go-plc", flag.ContinueOnError)
	fs.SetOutput(output)
	fs.StringVar(&opts.configPath, "config", DefaultConfigPath, "path to the YAML configuration file")
	fs.StringVar(&opts.logLevel, "log-level", "", "log level: debug, info, warn or error (overrides logLevel in the config file)")

	if err := fs.Parse(args); err != nil {
		return options{}, err
	}
	if fs.NArg() > 0 {
		return options{}, fmt.Errorf("unexpected arguments: %v", fs.Args())
	}
	if opts.logLevel != "" {
		if _, err := parseLevel(opts.logLevel); err != nil {
			return options{}, err
		}
	}

	return opts, nil
}

// parseLevel converts a level name to a slog.Level.
func parseLevel(name string) (slog.Level, error) {
	switch name {
	case "debug":
		return slog.LevelDebug, nil
	case "info":
		return slog.LevelInfo, nil
	case "warn":
		return slog.LevelWarn, nil
	case "error":
		return slog.LevelError, nil
	default:
		return 0, fmt.Errorf("invalid log level %q (allowed: debug, info, warn, error)", name)
	}
}

func main() {
	os.Exit(realMain(os.Args[1:]))
}

// realMain runs the program and returns the process exit code.
func realMain(args []string) int {
	opts, err := parseFlags(args, os.Stderr)
	if errors.Is(err, flag.ErrHelp) {
		return 0
	}
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return 2
	}

	// Set up structured logging using slog (standard library since Go 1.21).
	// We use JSON format for production compatibility with log aggregators.
	// The level starts at info (or the -log-level flag) and is updated from
	// the config file once it is loaded, unless the flag was given.
	var level slog.LevelVar
	if opts.logLevel != "" {
		lvl, _ := parseLevel(opts.logLevel) // validated in parseFlags
		level.Set(lvl)
	}
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: &level})))

	slog.Info("GoPLC starting", "version", version, "config", opts.configPath)

	cfg, err := config.Load(opts.configPath)
	if err != nil {
		slog.Error("Configuration failed to load", "error", err)
		return 1
	}

	if opts.logLevel == "" {
		lvl, err := parseLevel(cfg.LogLevel)
		if err != nil {
			slog.Error("Configuration failed to load", "error", err)
			return 1
		}
		level.Set(lvl)
	}

	slog.Info("Configuration loaded successfully", "sources", len(cfg.Sources), "variables", len(cfg.Variables))
	for _, src := range cfg.Sources {
		slog.Info("Source configured", "name", src.Name, "type", src.Type, "config", src.Config)
	}

	// TODO(maintainer): signal handling and graceful shutdown.
	ctx := context.Background()

	if err := run(ctx, cfg); err != nil {
		slog.Error("Runtime stopped with error", "error", err)
		return 1
	}

	return 0
}
