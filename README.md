<img src="assets/logo-blue.svg" alt="GoPLC" width="150">

_My first attempt at the logo, looks more like a seal. Will update when I have time_

---

![Status](https://img.shields.io/badge/status-development-blue)
<!--
![Status](https://img.shields.io/badge/status-planning-yellow)
![Status](https://img.shields.io/badge/status-release%20v1.0.0-green)
-->

A soft PLC written in Go that replaces proprietary industrial programming languages with modern development tools.

## What Problem Does This Solve?

Industrial automation tools are decades behind modern software development. Proprietary IDEs, vendor lock-in, slow iteration cycles, no version control. If you have to learn vendor-specific languages anyway, you might as well learn a real programming language that transfers to other domains.

GoPLC lets you write control logic in native Go using VS Code, Git, and CI/CD pipelines instead of Ladder Logic or Structured Text.

## What is a Soft PLC?

A traditional PLC is a specialized industrial computer that runs control logic for manufacturing equipment. A soft PLC does the same thing but runs on standard hardware—the control software without the proprietary box.

## Architecture Overview

GoPLC communicates with industrial devices via standard protocols and exposes data through multiple interfaces:

**Input/Output:**
- Modbus TCP/RTU for sensors, actuators, and I/O modules
- Automatic reconnection with exponential backoff

**SCADA Integration:**
- OPC UA server for traditional industrial SCADA systems
- Sparkplug B over MQTT for modern IIoT architectures
- GraphQL API for custom integrations

**Monitoring:**
- Built-in WebUI for real-time monitoring and control
- System health, connection status, variable values, task execution

Variables are defined once in YAML and automatically exposed to all protocols. The runtime handles scheduling, scaling, and protocol translation.

## Inspiration

This project is inspired by [Tentacle PLC](https://joyautomation.com/software/packages/tentacle), which proved that modern programming languages can replace IEC 61131-3 standards. GoPLC implements that philosophy in Go.

## Development Methodology

This project uses a hybrid AI-assisted workflow:

| Phase | Who |
|-------|-----|
| Planning & Documentation | AI-assisted (BMAD framework) |
| PRD, UX Design, Architecture | AI-assisted |
| Epics & Stories | AI-assisted (written for human developers) |
| **Development** | **Human (me)** |
| Code Review & QA | AI-assisted |

The goal: leverage AI where it excels (structured documentation, planning, review) while keeping implementation skills sharp through hands-on development.

## Current Status

**Planning stage.** PRD, UX design specifications, and architecture documentation are complete. Development begins soon.

## License

MIT License - use and modify with attribution.
