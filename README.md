<p align="center">
    <img src="assets/logo.svg" alt="GoPLC" height="150">
</p>

<p align="center">
    <img src="https://img.shields.io/badge/status-development-blue" alt="Status">
</p>

# GoPLC

A soft PLC written in Go that replaces proprietary industrial programming languages with modern development tools.

## Motivation

Industrial automation is stuck in the past. Proprietary IDEs, vendor lock-in, no real version control, and languages that don't transfer to any other domain. If you're going to learn something new anyway, why not use a real programming language with modern tooling? GoPLC lets you write control logic in native Go using VS Code, Git, and CI/CD pipelines instead of being locked into Ladder Logic or Structured Text. It's time to bring software development practices from the last two decades into industrial automation.

## Quick Start

Coming soon. Installation and usage instructions will be added as the project develops.

## Usage

### What is a Soft PLC?

A traditional PLC is a specialized industrial computer that runs control logic for manufacturing equipment. A soft PLC does the same thing but runs on standard hardware—the control software without the proprietary box.

### Architecture Overview

GoPLC communicates with industrial devices via standard protocols and exposes data through multiple interfaces:

**Input/Output:**
- Modbus TCP/RTU for sensors, actuators, and I/O modules
- Automatic reconnection with exponential backoff

**SCADA Integration:**
- OPC UA server for traditional industrial SCADA systems (Ignition, Kepware, etc.)
- GraphQL API for modern web/mobile integrations
- Optional Zenoh protocol for high-performance pub/sub (Phase 2)

**Monitoring:**
- Built-in WebUI for real-time monitoring and control
- System health, connection status, variable values, task execution

Variables are defined once in YAML and automatically exposed to all protocols. The runtime handles scheduling, scaling, and protocol translation.

## Development Guidelines

**Project Structure:**
- Registry pattern for extensible protocol support
- Protocol implementations in `internal/source/<protocol>/`
- Each protocol self-registers via `init()` functions

**Error Management:**
- Use structured error format: `[Component] - [Description] (context: value, ...)`
  - Example: `[config] - Unknown source type 'opcua' (source: plc1)`
- Each package contains an `errors.go` file with helper functions
- Use `panic()` only for programmer errors; return `error` for runtime errors

**Testing:**
- Unit tests use table-driven test pattern
- Test coverage required for all error paths

## Contributing

Contributions are welcome. If you find a bug or have a feature request, please open an issue. If you want to submit code, please open a pull request with a clear description of your changes.

## License

O'Saasy License Agreement - MIT do-whatever-you-want license, but with the commercial rights for SaaS reserved for the copyright holder.
