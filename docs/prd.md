---
stepsCompleted: [1, 2, 3, 4, 6, 7, 8, 9, 10, 11]
inputDocuments:
  - 'docs/brief.md'
documentCounts:
  briefs: 1
  research: 0
  brainstorming: 0
  projectDocs: 0
workflowType: 'prd'
lastStep: 11
project_name: 'go-plc'
user_name: 'Andy'
date: '2025-12-13'
---

# Product Requirements Document - go-plc

**Author:** Andy
**Date:** 2025-12-13

## Executive Summary

**go-plc** is a fast, simple soft PLC built in Go that enables automation engineers to develop and test PLC programs using modern software development tools and workflows. Designed for rapid iteration cycles, comprehensive industrial connectivity, and AI-assisted development - making PLC programming accessible to modern developers without requiring proprietary IDEs or vendor-specific languages.

**Primary Context:** Boot.dev capstone project with ~50 hours to MVP, built with community feedback and designed for real-world use beyond academic demonstration.

**Target Users:**
- **Primary:** Automation engineers frustrated with slow iteration cycles and proprietary tooling who want to use modern development workflows (Git, IDEs, CI/CD)
- **Secondary:** IoT/edge developers who know modern languages and need industrial protocol support without learning IEC 61131-3

### What Makes This Special

**The Key Insight:** Instead of asking "how do we make old tools work with new workflows," go-plc asks "what would a PLC look like if we designed it today with modern tools?"

**Core Differentiators:**

1. **Unapologetically modern** - Embraces a real programming language (Go) with all its tooling (Git, IDEs, testing frameworks, CI/CD) instead of trying to preserve or modernize legacy IEC 61131-3 languages

2. **Performance-first architecture** - Native Go compilation means predictable performance, real-time capability, and single binary deployment

3. **Built for modern developers** - If you know Go, you can program PLCs. No proprietary IDEs, no vendor-specific languages to learn

4. **Documentation as product** - Comprehensive Docusaurus site that makes industrial protocols and PLC concepts accessible to modern developers (not just industrial automation veterans)

5. **AI-native development** - Simple text-based configuration (YAML), code that AI models can actually read and generate (unlike vendor-specific graphical languages), potential MCP integration for AI-assisted development

6. **Open and vendor-neutral** - All standard protocols (Modbus, OPC UA, Zenoh, GraphQL), no proprietary lock-in, true portability across Linux/Windows

**Inspiration:** Inspired by James Joy's Tentacle PLC and the philosophy of abandoning IEC 61131-3 in favor of modern programming languages and tooling.

## Project Classification

**Technical Type:** IoT/Embedded (Industrial edge device with real-time control and industrial protocols)
**Domain:** General (Industrial automation without specialized regulatory requirements)
**Complexity:** Medium (Real-time requirements, multiple protocols, performance constraints)
**Project Context:** Greenfield - new project

This is an industrial IoT/edge application that must handle real-time control loop execution, multiple industrial communication protocols (Modbus TCP, OPC UA, Sparkplug B), and provide modern developer interfaces (GraphQL API, WebUI). Performance is critical with task execution overhead <50µs and API response times <10ms. The project prioritizes developer experience, comprehensive documentation, and AI-assisted development capabilities.

## Success Criteria

### User Success

**The "Worth It" Moment:**
Users know Go, read the documentation, and get their first PLC program running quickly. When they hit hurdles, comprehensive documentation and accessible community support (Discord/GitHub Issues) provide solutions. They accomplish their control tasks fast and easily without fighting the tool.

**Measurable User Success:**
- **Fast onboarding:** <30 minutes from "I want to test a control idea" to "I have it running" (first-time setup)
- **Rapid iteration:** <1 minute cycles for edit Go task → test on hardware
- **Familiar tooling:** Users leverage existing development tools (VS Code, Git, CI/CD) without proprietary IDEs
- **Protocol integration:** Automation engineers integrate industrial protocols (Modbus, OPC UA, Sparkplug B) without deep protocol expertise
- **Native language:** IoT/edge developers write PLC programs in Go rather than learning IEC 61131-3

**Emotional Success Moments:**
- **Relief:** "I don't need to learn another proprietary IDE"
- **Delight:** "It just compiled to a single binary and deployed"
- **Empowerment:** "I can version control my PLC code with Git!"

### Business Success

**Capstone Success (50 hours):**
- ✅ Tank battery reference implementation completed (migrated from CompactLogix)
- ✅ Professional documentation (Docusaurus) enables others to replicate
- ✅ Demonstrates advanced Go skills: concurrency, protocols, APIs, real-time constraints
- ✅ AI-assisted migration workflow documented (CompactLogix → go-plc)

**The Ultimate Test:**
**"Would I actually use this in a real control project?"**

If the developer would choose go-plc over CompactLogix for commissioning a real tank battery, the project succeeds.

**3 Months Post-Capstone:**
- At least one other person successfully deployed go-plc for a real or test project
- Community feedback incorporated from Discord/GitHub during development
- go-plc version of tank battery is simpler/clearer than Rockwell version

**12 Months:**
- Developer still using it / would recommend it for control projects
- Small community of users exists
- Project demonstrates deep understanding of industrial control systems

### Technical Success

**Core Success Pillars:**

**1. Ease of Development**
- Time to implement tank battery competitive with CompactLogix
- Reduced lines of code compared to Rockwell implementation
- Clear, maintainable Go code structure

**2. Ease of Integration**
- SCADA connection setup: <30 minutes from documentation
- Works with Ignition via OPC UA and Sparkplug B out of the box
- No custom driver development needed

**3. Reliability (24/7 Operation)**
- Modbus connection failures handled gracefully with automatic reconnection
- Configurable logging levels for troubleshooting without restarts
- Graceful shutdown with no data corruption
- Memory stability in long-running processes
- System health status monitoring

### Measurable Outcomes

**Performance Targets:**
- Task execution overhead: <50µs per task cycle
- GraphQL/API response time: <10ms
- Scan rate support: 100ms typical for tank battery control loops
- Code reduction: 80% less boilerplate vs. traditional PLC programming

**Validation Criteria:**
- Integration test with Python Modbus simulator
- Integration test with Ignition consuming OPC UA and/or Sparkplug B
- Performance benchmarks documented (baseline on development hardware)
- Tank battery control logic functions correctly end-to-end

**Documentation Quality:**
- Comprehensive Docusaurus site covering:
  - Installation and setup
  - Linux device configuration for real-time performance
  - Task programming guide
  - Tag/variable configuration
  - Third-party device integration (Sparkplug B, OPC UA, Modbus)
  - WebUI monitoring guide
  - Performance benchmark methodology
  - Tank battery reference implementation walkthrough

## Product Scope

### MVP - Minimum Viable Product (50 hours)

**Core PLC Runtime:**
1. **Modbus I/O** - Read/write holding registers and coils via Modbus TCP using simonvetter/modbus library (includes RTU/ASCII for future use). Python Modbus simulator for testing without physical hardware.
2. **Go-based task logic** - Native Go tasks with clean API for variable access and control logic. Tasks auto-discovered from `/tasks` folder.
3. **OPC UA server** - Using gopcua/opcua library for SCADA integration (Ignition, Kepware, etc.)
4. **GraphQL API** - gqlgen-based API with queries and subscriptions. Includes built-in GraphQL Playground for testing and development.
5. **Monitoring WebUI** - Vite + React + urql frontend connecting to GraphQL backend:
   - PLC status display with enable/disable operation control (run mode)
   - Connected device/source status (Modbus, OPC UA connection health)
   - Variable list with real-time values
   - Task list with configuration, execution status, and enable/disable control
6. **Comprehensive documentation** - Docusaurus site with setup, configuration, task programming guide, and tank battery walkthrough
7. **Tank battery reference implementation** - Real-world control project migrated from CompactLogix

**Reliability Features (MVP):**
- Modbus automatic reconnection logic with exponential backoff
- Logging framework with configurable levels
- Graceful shutdown handling
- Basic error handling and reporting

**Technical Validation:**
- Integration test suite with Python Modbus simulator
- Integration test with Ignition SCADA via OPC UA
- Performance benchmarks documented on development hardware
- Tank battery control logic verified end-to-end

**Done = Working tank battery + Docusaurus docs + <30min quickstart**

### Growth Features (Phase 2 - Post-MVP)

**Protocol Expansion:**
- **Zenoh Protocol** - High-performance pub/sub implementation using zenoh-go bindings for advanced SCADA integration (900+ channels @ 100Hz capability), geo-distributed storage, query/reply patterns, and peer-to-peer communication without requiring message broker infrastructure
- **Ethernet/IP** - Allen-Bradley ecosystem connectivity enabling:
  - VFDs (PowerFlex series)
  - Remote I/O (Point I/O, Flex I/O)
  - Other Allen-Bradley compatible devices

**Enhanced WebUI:**
- Source configuration (add, modify, enable/disable Modbus/OPC UA/Ethernet/IP sources)
- Variable configuration (add, modify, enable/disable variables)
- Configuration persistence and validation

**Enhanced Reliability:**
- Persistent state (survive restarts without data loss)
- Watchdog integration for fault detection
- Advanced fault tolerance mechanisms
- Load testing and burn-in validation

**Developer Experience:**
- Hot reload for tasks during development
- MCP server for AI-assisted PLC programming

### Vision Features (Phase 3 - Expansion)

**Industrial Protocol Gateway:**
- **PROFINET** - Siemens ecosystem connectivity (S7 PLCs, Siemens VFDs, Siemens Remote I/O)
- **EtherCAT** - High-speed motion control and servo drive integration
- **BACnet** - Building automation system integration for facility crossover applications

**Target Device Categories:**
- VFDs (Siemens, ABB, Danfoss, Yaskawa via respective protocols)
- Remote I/O (WAGO, Beckhoff, Siemens ET 200)
- HMIs (protocol-agnostic data serving)
- Servo drives and motion controllers (EtherCAT)
- Building systems (BACnet-enabled HVAC, lighting, metering)

**Protected Browser-based Task Development:**
- Disabled by default - requires explicit password configuration in YAML to enable
- Safety features: warning banners, multi-step confirmation before deployment, syntax validation, automatic backup
- Code editor (Monaco/CodeMirror) for Go task files in separate "Development" tab
- Intended for remote commissioning scenarios with appropriate guardrails

**Platform Features:**
- Remote parameter tuning via WebUI (setpoints, thresholds)
- Configuration hot-reload (YAML changes without restart)
- AI code generation for control logic
- Community-driven task and pattern library

**Long-term Vision:**
go-plc as a universal industrial protocol gateway and soft PLC platform - enabling modern developers to connect any industrial device using familiar tools (Go, Git, VS Code, CI/CD) without vendor lock-in or proprietary IDEs. Proven reliability in production deployments running 24/7 for years.

## User Journeys

### Journey 1: Marcus Chen - Breaking Free from Vendor Lock-in
**User Type: Automation Engineer**

Marcus is a controls engineer at an upstream oil & gas company who's been programming Allen-Bradley PLCs for 15 years. He's tired of waiting 20 minutes just to compile and download a simple logic change to test a new pump sequencing idea. When his manager assigns him a new site coming online - a central tank battery that needs full automation - Marcus sees an opportunity to try something different.

Late one Friday, frustrated after another slow iteration cycle on a different project, he discovers go-plc on GitHub. The promise of "sub-minute iteration cycles" and "version control your PLC code with Git" catches his attention. He decides to spend his weekend trying something new.

Saturday morning, Marcus installs Go and go-plc following the Docusaurus docs. Within 25 minutes, he has a simple control task running - reading Modbus values from a Python simulator and controlling a virtual pump. He's shocked. He makes a logic change, rebuilds, and sees it running 45 seconds later. No proprietary IDE. No licensing dongles. Just his favorite VS Code editor and a terminal.

The breakthrough comes when he implements the tank battery logic in Go. The code is clear, version-controlled in Git, and he can iterate on control strategies in minutes instead of hours. When he connects Ignition SCADA via OPC UA and sees real-time data flowing, he realizes he's found what he's been looking for. Two weeks later, the new site comes online with go-plc running the tank battery, and Marcus has a modern, maintainable control system that he can troubleshoot and update from anywhere with Git access.

**This journey reveals requirements for:**
- Clear installation and quickstart documentation
- Go task development with simple variable access API
- Modbus I/O integration with simulator support for testing
- SCADA integration via OPC UA (with optional Zenoh for advanced use cases)
- Fast build and deployment workflow (<1 minute iteration)
- Git-friendly project structure (text-based configuration)
- Tank battery reference implementation as example

### Journey 2: Sarah Rodriguez - Bridging the IT/OT Gap
**User Type: IoT/Edge Developer**

Sarah is a software developer at a tech startup building an edge analytics platform for industrial sites. She knows Python, Go, and modern web technologies, but when her product manager says they need to pull data from customer PLCs via Modbus and OPC UA, she panics. She's never touched industrial protocols and the thought of learning Ladder Logic to understand customer systems is overwhelming.

Her team considers buying expensive industrial gateway hardware, but the $5,000 per-site cost doesn't fit their startup budget. Then she finds go-plc while searching for "open source Modbus OPC UA golang." The documentation promises "industrial protocol support without deep protocol expertise."

She follows the quickstart guide and within 30 minutes has go-plc reading Modbus data from a test simulator. The Go code makes sense - it's just like any other Go application she's built. She configures a few YAML files for variables and scaling, writes a simple Go task to process the data, and suddenly she has industrial data flowing into her analytics platform via the GraphQL API.

The real win comes when a customer asks if they can integrate via OPC UA instead of the REST API. Sarah checks the go-plc docs, adds three lines to the YAML config, and OPC UA is enabled. No custom driver development. No vendor support contracts. Just configuration. Her CTO is impressed that she bridged the IT/OT gap in a weekend using familiar tools.

**This journey reveals requirements for:**
- GraphQL API for modern web/app integration
- YAML configuration documentation for non-PLC developers
- Multiple integration options (GraphQL, OPC UA, Sparkplug B)
- Industrial protocol abstraction (no deep Modbus/OPC UA expertise needed)
- Clear API documentation and examples
- Quick integration path for developers familiar with Go

### Journey 3: Jake Morrison - Night Shift Problem Solver
**User Type: Operations/Maintenance Personnel**

Jake is an operations technician working night shifts at a remote tank battery site. He's not a programmer, but he's responsible for keeping the automation running and responding when alarms go off. When the old PLC system had issues, he'd have to call Marcus (the controls engineer) at 2 AM and wait hours for remote troubleshooting.

With go-plc deployed at his site, Jake's life changes. One night, he notices tank levels aren't updating in SCADA. Instead of immediately calling for help, he opens the go-plc WebUI on his tablet. He can see that the Modbus connection to the remote I/O is disconnected - the status shows red with "Connection failed: timeout."

He checks the physical I/O rack, notices a loose network cable, reconnects it, and watches the WebUI status turn green within seconds. The tank levels start updating again. The configurable logging helped him see exactly what failed without needing to interpret cryptic PLC error codes.

Jake sends Marcus a quick text: "Fixed the Modbus issue - loose cable. WebUI made it obvious. Back to normal." Marcus, still asleep, sees it in the morning and smiles. The system's transparency and clear status monitoring just saved both of them a middle-of-the-night troubleshooting session.

**This journey reveals requirements for:**
- WebUI with real-time connection status monitoring
- Clear error messages (no cryptic codes)
- Visual health indicators (connection state: green/red)
- Mobile-friendly WebUI (tablet access)
- Configurable logging with human-readable output
- System diagnostics accessible to non-programmers

### Journey 4: David Park - Fast Commissioning Under Pressure
**User Type: System Integrator**

David runs a small controls integration company that gets called when oil & gas operators need automation systems commissioned quickly. He's just landed a contract to integrate three new tank battery sites with the client's existing Ignition SCADA system. The deadline is tight - 6 weeks for all three sites.

The client mentions their controls engineer (Marcus) has been piloting go-plc and wants David to deploy it across all three sites. David is skeptical - he's used to traditional PLCs and knows their quirks. But the timeline is aggressive, so he agrees to try.

David downloads go-plc and reads through the Docusaurus deployment guide. He's pleasantly surprised - there's a complete section on Linux device configuration for real-time performance, SCADA integration examples, and even a tank battery reference implementation that matches exactly what he needs to deploy.

At the first site, he installs go-plc on an industrial Linux edge device, copies the tank battery configuration YAML, adjusts the Modbus addresses for the local I/O, and starts the service. Within 2 hours, the PLC is running. He configures the OPC UA server settings, and Ignition connects immediately to start pulling real-time data.

The breakthrough moment comes when the client requests a logic change during commissioning. David opens the Go task file, makes the adjustment, rebuilds the binary, and restarts the service. Total time: 3 minutes. With a traditional PLC, this would have meant opening a proprietary IDE, connecting remotely, downloading the program, and testing - easily 30 minutes.

By week 4, all three sites are commissioned and running. David finishes two weeks early and realizes go-plc just changed his business model. He can bid more competitively and deliver faster than competitors still using traditional PLCs.

**This journey reveals requirements for:**
- Linux deployment documentation (installation, configuration, service setup)
- Real-time performance tuning guide for Linux
- SCADA integration configuration examples (OPC UA, with optional Zenoh for advanced scenarios)
- Reusable configuration templates (tank battery example)
- Single binary deployment (no complex dependencies)
- Fast rebuild and restart workflow for field changes
- Production-ready reliability and error handling

### Journey 5: Andy - Proving Modern PLC Development
**User Type: Developer/Creator (Capstone Context)**

Andy is completing his Boot.dev capstone project and wants to build something that demonstrates real-world engineering skills - not just another CRUD app that sits in a GitHub repo collecting dust. He's worked in industrial automation and knows the pain points: slow iteration, proprietary tools, vendor lock-in. He decides to build go-plc to prove there's a better way.

He has an existing CompactLogix tank battery program from a previous project - real production code with tank level monitoring, pump sequencing, and alarm logic. This becomes his validation case: if he can recreate this in go-plc and it's actually better, the project succeeds.

Andy starts by researching existing Go libraries to accelerate development. He finds `simonvetter/modbus` for Modbus TCP and `gopcua/opcua` for OPC UA server capabilities. For advanced protocol support in Phase 2, he discovers Zenoh - a next-generation pub/sub protocol with demonstrated 900+ channels @ 100Hz performance on edge devices. He plans to contribute to the `zenoh-go` bindings, building both his portfolio and the open-source ecosystem. Standing on the shoulders of these existing libraries means he can focus on the core PLC runtime architecture rather than reinventing low-level protocol handling.

He builds the project structure, integrating these libraries into a cohesive runtime. The YAML configuration system comes next - he wants users to define variables once and have them auto-expose to all protocols (Modbus, OPC UA, GraphQL, and eventually Zenoh).

The critical moment comes when he uses AI to analyze his CompactLogix program and help translate the logic into Go tasks. The AI understands the control patterns and generates clean Go code. Andy refines it, adds proper error handling, and implements the tank battery logic. When he runs it for the first time and sees the same control behavior as the CompactLogix version - but with faster iteration cycles and readable code - he knows he's onto something.

He spends significant time on documentation, creating a comprehensive Docusaurus site. He wants someone like Marcus (the automation engineer from Journey 1) to be able to adopt go-plc without Andy having to personally onboard them. Installation guides, configuration examples, the tank battery walkthrough, performance benchmarks - all documented thoroughly.

When Andy deploys the final system, integrates it with Ignition SCADA via OPC UA, and validates the performance metrics (<50µs task execution, <10ms API response), he has his answer to the ultimate question: "Would I use this in a real control project?"

Yes. Absolutely yes.

He submits his capstone with confidence - not just because it demonstrates Go skills, concurrency, real-time systems, and protocol integration - but because he's built something he'd actually recommend to other engineers.

**This journey reveals requirements for:**
- Integration with existing Go libraries (simonvetter/modbus, gopcua/opcua)
- Optional Zenoh protocol integration for Phase 2 (high-performance pub/sub with zenoh-go contribution)
- YAML configuration system with single-definition principle
- AI-friendly code structure for migration assistance
- Performance validation framework (benchmarking <50µs, <10ms targets)
- Comprehensive Docusaurus documentation site
- Tank battery reference implementation with migration walkthrough
- Testing framework (Python Modbus simulator, Ignition integration tests)

### Journey Requirements Summary

**Core PLC Runtime Capabilities:**
- Go-based task execution with simple variable access API
- YAML configuration for variables, I/O sources, and scaling
- Fast build and deployment workflow (<1 minute iteration cycles)
- Single binary deployment with embedded WebUI
- Configurable logging framework with multiple levels
- Graceful error handling and automatic reconnection logic

**Protocol Integration:**
- Modbus TCP client (using simonvetter/modbus library)
- OPC UA server (using gopcua/opcua library)
- GraphQL API with queries and subscriptions
- Optional Zenoh pub/sub (Phase 2 - using zenoh-go bindings for high-performance scenarios)

**Monitoring and Operations:**
- WebUI showing real-time status (tags, connections, task execution)
- Clear visual health indicators (connection state)
- Mobile-friendly interface for field access
- Human-readable error messages and logging
- System diagnostics accessible to non-programmers

**Documentation and Developer Experience:**
- Comprehensive Docusaurus site with:
  - Installation and quickstart (<30 minutes to first running task)
  - Linux device configuration for real-time performance
  - Task programming guide (Go code examples)
  - YAML configuration reference
  - Protocol integration guides (Modbus, OPC UA, GraphQL, optional Zenoh)
  - Tank battery reference implementation walkthrough
  - Performance benchmark methodology
  - Troubleshooting and diagnostics guide

**Testing and Validation:**
- Python Modbus simulator for development without hardware
- Integration test framework
- Performance benchmarking tools
- Ignition SCADA integration examples

## Innovation & Novel Patterns

### Detected Innovation Areas

**Core Innovation: Complete Rejection of IEC 61131-3**

go-plc follows the groundbreaking philosophy pioneered by James Joy's Tentacle PLC: industrial automation doesn't need proprietary IEC 61131-3 languages (Ladder Logic, Structured Text, FBD). Modern programming languages are superior for PLC development.

**go-plc's Contribution:**

Building on Tentacle's proven approach of using modern languages, go-plc explores whether native compiled Go can deliver even better performance and simplicity compared to interpreted languages. The goal is simple and fast PLC development using tools modern developers already know.

**Key Innovative Aspects:**

1. **Native Compiled Approach** - Pure Go compilation (no runtime engine overhead) for predictable performance, real-time capability, and single binary deployment

2. **YAML Single-Definition Principle** - Define variables once in YAML, automatically expose to all protocols (Modbus, OPC UA, GraphQL, Zenoh) - eliminating duplicate definitions across systems

3. **AI-Assisted Migration** - Leveraging AI to analyze existing PLC programs (CompactLogix, etc.) and translate logic into clean Go code, making migration from proprietary systems accessible

4. **Comprehensive Protocol Integration** - All industrial protocols in one platform (Modbus TCP, OPC UA, GraphQL, optional Zenoh) without vendor lock-in or complex gateway hardware

5. **Developer Experience as Product** - Git version control, modern IDEs, CI/CD pipelines, and comprehensive documentation make industrial automation accessible to the broader software development community

### Market Context & Competitive Landscape

**Existing Approaches:**

- **Traditional PLCs (Allen-Bradley, Siemens, etc.):** Proprietary IEC 61131-3 implementations with slow iteration cycles, expensive IDEs, vendor lock-in
- **Tentacle PLC:** Pioneered modern language approach using JavaScript, proving the concept works in production
- **go-plc Position:** Evolution of Tentacle's philosophy using native Go for optimal performance and simplicity

**Target Gap:**

The industry has been trying to "modernize" IEC 61131-3 by adding Git integration or web IDEs to legacy languages. This approach is backwards. go-plc (following Tentacle's lead) asks: "What if we just used modern languages designed for concurrent, real-time systems?"

**Differentiation:**

Not competing with Tentacle PLC - building on its proven philosophy with a different implementation approach (compiled Go vs. JavaScript) optimized for performance-critical applications and Boot.dev capstone demonstration of advanced Go skills.

### Validation Approach

**Primary Validation: CompactLogix Tank Battery Migration**

The innovation will be validated through a real-world test case:

1. **Existing Production Logic:** Tank battery control program from CompactLogix (tank level monitoring, pump sequencing, alarm logic)
2. **AI-Assisted Translation:** Use AI to analyze Rockwell logic and generate equivalent Go code
3. **Functional Equivalence:** go-plc implementation must control tank battery identically to CompactLogix version
4. **Performance Validation:** Achieve <50µs task execution overhead and <10ms API response times
5. **Integration Validation:** Successfully integrate with Ignition SCADA via OPC UA

**Success Criteria:**

Answer the question: "Would I actually use this in a real control project?"

If go-plc can replace CompactLogix for a production tank battery with:
- Faster iteration cycles (<1 minute vs. 20+ minutes)
- Clearer, more maintainable code
- Modern development workflow (Git, VS Code, CI/CD)
- Equivalent or better performance

Then the innovation succeeds.

**Validation Tools:**

- Python Modbus simulator for hardware-independent testing
- Ignition SCADA for protocol integration validation
- Performance benchmarking framework for objective measurement
- CompactLogix program as reference implementation

### Risk Mitigation

**Innovation Risk: Automation Engineer Adoption**

**Risk:** Automation engineers may resist learning Go instead of familiar IEC 61131-3 languages.

**Mitigation:**
- **Philosophy from Tentacle PLC:** Engineers already learn vendor-specific IEC 61131-3 implementations - might as well learn a real, transferable language
- **Documentation as Product:** Comprehensive Docusaurus site bridges knowledge gap for automation engineers new to Go
- **Target IoT Developers First:** Secondary user group (IoT/edge developers) already knows modern languages and needs industrial protocol support
- **Reference Implementation:** Tank battery example demonstrates complete real-world application
- **AI-Assisted Migration:** Lowers barrier to migrating existing PLC programs

**Innovation Risk: Performance Requirements**

**Risk:** Compiled Go might not meet <50µs task execution overhead target.

**Mitigation:**
- **Early Benchmarking:** Performance validation framework in MVP
- **Proven Technology:** Go designed for concurrent, real-time systems
- **Fallback:** If performance insufficient, document findings and scope lessons learned for capstone

**Innovation Risk: Protocol Complexity**

**Risk:** Implementing OPC UA, Modbus, GraphQL might exceed 50-hour timeline.

**Mitigation:**
- **Leverage Existing Libraries:** simonvetter/modbus, gopcua/opcua
- **Clear MVP Scope:** Working tank battery with Modbus + OPC UA + GraphQL proves concept
- **Phase 2 Advanced Protocols:** Zenoh integration deferred to post-MVP for advanced performance scenarios

**The Ultimate Fallback:**

Even if adoption is limited, the capstone demonstrates:
- Advanced Go skills (concurrency, protocols, real-time systems)
- Real-world engineering problem-solving
- Comprehensive documentation and testing
- Something the developer would actually use

## IoT/Embedded Specific Requirements

### Project-Type Overview

go-plc is designed as an industrial IoT/embedded application optimized for edge deployment in automation environments. Unlike traditional PLCs with custom hardware, go-plc runs on standard industrial Linux edge devices, enabling modern deployment practices while maintaining industrial reliability requirements.

**Deployment Model:**
- Primary: Industrial Linux edge devices at remote sites (tank batteries, process facilities)
- Development: Linux/Windows workstations for development and testing
- Target: Cross-platform single binary deployment

### Hardware Requirements

**Industrial Edge Devices:**
- Standard industrial PC or edge gateway hardware
- AC-powered (no special power constraints)
- Typical industrial edge device specifications sufficient
- Linux operating system (Ubuntu, Debian, or industrial Linux distributions)
- Windows support for development environments

**Minimum Specifications:**
- Modern multi-core CPU (for concurrent task execution and protocol handling)
- Sufficient RAM for in-memory variable storage (scale based on tag count)
- Network interface (Ethernet) for Modbus TCP, SCADA, and API communication
- Storage for single binary and log files

**No Special Requirements:**
- No battery/UPS integration in MVP (standard industrial site power assumed)
- No specialized industrial-hardened hardware required
- No fanless or extreme temperature considerations for MVP

### Connectivity & Protocol Architecture

**Field-Level Communication:**
- **Modbus TCP** - Primary I/O protocol for connecting to remote I/O, sensors, actuators
- **simonvetter/modbus** library provides TCP client implementation
- Automatic reconnection with exponential backoff on connection failures
- Configurable polling intervals per source (default 100ms)

**SCADA Integration:**
- **OPC UA Server** - gopcua/opcua library for standard industrial SCADA connectivity (Ignition, Kepware, etc.)
- **Optional Zenoh (Phase 2)** - zenoh-go bindings for high-performance pub/sub, geo-distributed storage, and query/reply patterns
- Variables automatically exposed to all protocols via single YAML definition

**Modern Integration:**
- **GraphQL API** - gqlgen-based API with queries and subscriptions for web/app integration
- **WebUI** - Embedded web interface for monitoring (real-time status via GraphQL subscriptions)

**Network Assumptions:**
- Ethernet-based networking (wired industrial networks)
- No wireless/cellular connectivity in MVP
- Standard TCP/IP networking stack
- Firewalled industrial network environment assumed

### Security Model

**MVP Security Approach:**

go-plc MVP prioritizes rapid development and functional validation over comprehensive security hardening. Security implementation follows an iterative approach:

**MVP Security Posture:**
- **No authentication/authorization** - Internal development tool, trusted network assumption
- **Network Security:** Assumed deployment behind industrial firewalls and VLANs
- **Protocol Security:** Standard protocol implementations without TLS/encryption
- **Access Control:** No role-based access control in MVP

**Security Assumptions:**
- Deployment in controlled industrial network environments
- Physical security of edge devices
- Network segmentation from untrusted networks
- Internal development/testing use case for MVP validation

**Post-MVP Security Roadmap:**

Cybersecurity is critical for production industrial deployments. Post-MVP security enhancements include:

- **Authentication & Authorization:**
  - API key authentication for GraphQL API
  - User authentication for WebUI access
  - Role-based access control (RBAC) for different user types

- **Protocol Security:**
  - TLS/SSL encryption for OPC UA connections
  - Secure MQTT (TLS) for Sparkplug B
  - HTTPS for WebUI and GraphQL endpoints

- **Data Security:**
  - Encrypted configuration storage
  - Audit logging for security events
  - Secure credential management

- **Network Security:**
  - Rate limiting and DDoS protection
  - Input validation and sanitization
  - Security headers and CSP for WebUI

**Security Documentation:**

Comprehensive security documentation will be provided in Docusaurus site:
- Security assumptions and limitations (MVP)
- Network architecture recommendations
- Post-MVP security hardening guide
- Compliance considerations for different deployment scenarios

### Deployment & Update Mechanisms

**MVP Deployment Model:**

Manual deployment process following standard DevOps best practices:

**Deployment Process:**
1. Build single binary for target platform (Linux/Windows)
2. Transfer binary to edge device (scp, sftp, or deployment tool)
3. Stop existing service (if running)
4. Replace binary
5. Start service with configuration file
6. Verify connectivity and operation via WebUI

**Configuration Management:**
- YAML configuration files managed separately from binary
- Version control for configuration (Git recommended)
- Configuration validation on startup with clear error messages

**Standard CI/CD Integration:**

Developers building PLC projects with go-plc should follow modern CI/CD practices:

- **Version Control:** Git for task code and configuration
- **Automated Testing:** Integration tests with Python Modbus simulator
- **Build Pipeline:** Automated builds on commit/merge
- **Deployment Automation:** Scripted deployment to edge devices
- **Rollback Strategy:** Version-tagged binaries for quick rollback

**Documentation Focus:**

Comprehensive deployment guides in Docusaurus site:

1. **Local Development Deployment**
   - Running go-plc on development workstation
   - Python Modbus simulator setup
   - Hot-reload workflow for rapid iteration

2. **Production Industrial Edge Deployment**
   - Linux service configuration (systemd)
   - Network configuration and firewall rules
   - Logging configuration and rotation
   - Monitoring and health checks

3. **CI/CD Integration Examples**
   - GitHub Actions workflow examples
   - GitLab CI pipeline templates
   - Automated testing strategies
   - Deployment script examples

4. **Version Management Strategies**
   - Semantic versioning for go-plc runtime
   - Configuration versioning approaches
   - Migration guides between versions

**No OTA (Over-The-Air) Updates in MVP:**

OTA update mechanisms are explicitly out of scope for MVP. Post-MVP growth features may include:
- Remote update capabilities with verification
- Staged rollout mechanisms
- Automatic rollback on failure detection

### Real-Time Performance Requirements

**Task Execution:**
- <50µs task execution overhead per cycle
- Deterministic scheduling for real-time control
- Goroutine-based concurrency model for parallel task execution

**Linux Real-Time Configuration:**

Docusaurus documentation will include Linux real-time tuning guide:
- PREEMPT_RT kernel configuration (optional for soft real-time)
- CPU isolation and affinity settings
- Process priority configuration
- Network stack tuning for deterministic I/O

**Performance Validation:**

Performance benchmarking framework included in MVP:
- Task execution overhead measurement
- API response time measurement
- Protocol latency tracking
- System resource utilization monitoring

### Implementation Considerations

**Single Binary Deployment:**
- Go compilation produces single executable
- WebUI assets embedded using Go embed package
- No external dependencies beyond standard Linux libraries
- Cross-compilation for Linux and Windows targets

**Concurrent Architecture:**
- Goroutines for I/O sources (Modbus polling loops)
- Goroutines for task execution (scheduled per task scan rate)
- Goroutines for protocol servers (OPC UA, GraphQL, WebUI)
- Channel-based communication between components

**Graceful Shutdown:**
- Signal handling for SIGTERM/SIGINT
- Coordinated shutdown of all goroutines
- Flush logs and close connections cleanly
- No data corruption on normal shutdown

**Logging & Diagnostics:**
- Configurable logging levels (debug, info, warn, error)
- Structured logging for machine parsing
- Human-readable error messages for operators
- Log rotation and retention policies (documented)

**Resource Management:**
- In-memory variable storage (no persistence in MVP)
- Memory-efficient protocol buffers for Sparkplug B
- Connection pooling and resource cleanup
- Memory leak detection in testing

## Functional Requirements

### PLC Runtime Control

- FR1: Operators can enable/disable the PLC runtime (run mode control)
- FR2: System can execute Go-based tasks at configurable scan rates
- FR3: System can auto-discover tasks from the `/tasks` folder
- FR4: Operators can enable/disable individual tasks at runtime
- FR5: System can gracefully shutdown without data corruption
- FR6: System can start with a YAML configuration file

### Variable Management

- FR7: Developers can define variables in YAML with source binding and scaling
- FR8: System can expose variables to all protocols from a single definition
- FR9: Tasks can read and write variable values through a clean API
- FR10: Operators can view all variables and their current values
- FR11: System can apply linear scaling to variable values (raw to engineering units)

### Modbus Communication

- FR12: System can connect to Modbus TCP devices as a client
- FR13: System can read holding registers from Modbus devices
- FR14: System can write holding registers to Modbus devices
- FR15: System can read coils from Modbus devices
- FR16: System can write coils to Modbus devices
- FR17: System can automatically reconnect on Modbus connection failure
- FR18: Operators can view Modbus connection status (connected/disconnected/error)
- FR19: Developers can configure multiple Modbus sources with independent polling intervals

### OPC UA Integration

- FR20: System can expose variables as an OPC UA server
- FR21: SCADA systems can connect to go-plc via OPC UA
- FR22: SCADA systems can read variable values via OPC UA
- FR23: SCADA systems can write variable values via OPC UA
- FR24: Operators can view OPC UA server status

### GraphQL API

- FR25: External applications can query current variable values via GraphQL
- FR26: External applications can subscribe to variable value changes via GraphQL
- FR27: External applications can query PLC status via GraphQL
- FR28: External applications can query task status via GraphQL
- FR29: External applications can query source/connection status via GraphQL
- FR30: Developers can test GraphQL queries using built-in GraphQL Playground

### WebUI Monitoring

- FR31: Operators can view PLC runtime status in WebUI
- FR32: Operators can enable/disable PLC runtime from WebUI
- FR33: Operators can view all source/device connection status in WebUI
- FR34: Operators can view variable list with real-time values in WebUI
- FR35: Operators can view task list with configuration and status in WebUI
- FR36: Operators can enable/disable individual tasks from WebUI
- FR37: WebUI can display real-time updates without page refresh

### Task Development

- FR38: Developers can write control logic in native Go
- FR39: Tasks can access variables through a simple API (no verbose patterns)
- FR40: Developers can specify task scan rate in task configuration
- FR41: Developers can rebuild and deploy tasks in under 1 minute
- FR42: System can report task execution errors with clear messages

### Logging & Diagnostics

- FR43: System can log events at configurable levels (debug, info, warn, error)
- FR44: Operators can view human-readable error messages (no cryptic codes)
- FR45: System can log connection state changes
- FR46: System can log task execution errors

### Configuration

- FR47: Developers can define sources (Modbus devices) in YAML
- FR48: Developers can define variables with source bindings in YAML
- FR49: System can validate configuration on startup with clear error messages
- FR50: System can report configuration errors before starting runtime

### Deployment

- FR51: System can compile to a single binary with embedded WebUI
- FR52: System can run on Linux and Windows platforms
- FR53: System can run as a systemd service on Linux

## Non-Functional Requirements

### Performance

**Task Execution:**
- NFR1: Task execution overhead must be <50µs per cycle
- NFR2: Task scheduler must support scan rates from 10ms to 10s
- NFR3: Variable read/write operations must complete within task cycle budget

**API Response:**
- NFR4: GraphQL query response time must be <10ms for variable reads
- NFR5: GraphQL subscriptions must deliver updates within 100ms of value change
- NFR6: OPC UA read operations must complete within 50ms

**Resource Efficiency:**
- NFR7: Memory usage must remain stable during 24+ hour operation (no memory leaks)
- NFR8: CPU usage must remain <50% on target hardware during normal operation
- NFR9: Single binary size must be <50MB (including embedded WebUI)

### Reliability

**Availability:**
- NFR10: System must support 24/7 continuous operation
- NFR11: System must recover from Modbus connection failures without operator intervention
- NFR12: System must complete graceful shutdown within 5 seconds

**Error Handling:**
- NFR13: All errors must be logged with human-readable messages
- NFR14: Connection failures must trigger automatic reconnection with exponential backoff
- NFR15: Configuration errors must be reported at startup before runtime begins

**Data Integrity:**
- NFR16: Variable values must remain consistent across all protocols (Modbus, OPC UA, GraphQL)
- NFR17: No data corruption on graceful shutdown

### Security

**MVP Security Posture (Trusted Network):**
- NFR18: MVP assumes deployment in trusted, firewalled industrial network
- NFR19: No authentication required for MVP (documented limitation)
- NFR20: All network services bind to configurable interfaces (not hardcoded to 0.0.0.0)

**Post-MVP Security (Documented for Future):**
- Authentication and authorization deferred to Phase 2
- TLS/encryption deferred to Phase 2
- Security hardening guide to be included in documentation

### Integration

**Protocol Compliance:**
- NFR21: Modbus TCP implementation must comply with Modbus Application Protocol Specification
- NFR22: OPC UA server must be compatible with standard OPC UA clients (Ignition, Kepware)
- NFR23: GraphQL API must follow GraphQL specification for queries and subscriptions

**Interoperability:**
- NFR24: System must successfully integrate with Ignition SCADA via OPC UA
- NFR25: System must work with Python pymodbus simulator for testing
- NFR26: WebUI must function in modern browsers (Chrome, Firefox, Edge - latest 2 versions)

### Maintainability

**Code Quality:**
- NFR27: Code must follow standard Go formatting (gofmt)
- NFR28: Code must pass go vet with no warnings
- NFR29: Public APIs must have documentation comments

**Deployment:**
- NFR30: System must compile to single binary for Linux (amd64, arm64) and Windows (amd64)
- NFR31: Configuration changes must not require recompilation
- NFR32: Logs must support configurable output levels without restart
