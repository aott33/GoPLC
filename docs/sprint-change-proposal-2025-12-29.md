# Sprint Change Proposal: Zenoh Protocol Substitution

**Date:** 2025-12-29
**Status:** Approved (User: Andy)
**Triggered By:** Strategic technology decision before Phase 2 implementation
**Change Scope:** Minor - Phase 2 protocol substitution (no MVP impact)

---

## 1. Issue Summary

### Problem Statement

Before implementing Phase 2 protocol expansion, a strategic technology evaluation identified that **Zenoh protocol** offers significant advantages over the originally planned **MQTT Sparkplug B** implementation:

1. **Superior Performance:** Community-demonstrated capability of 900+ channels @ 100Hz on Raspberry Pi CM4 edge devices, compared to standard MQTT performance limitations

2. **Enhanced Feature Set:** Zenoh provides unified pub/sub + geo-distributed storage + query/reply patterns, whereas Sparkplug B only offers pub/sub messaging

3. **Deployment Simplification:** Zenoh supports peer-to-peer communication without requiring MQTT broker infrastructure

4. **Open-Source Contribution Opportunity:** Contributing to zenoh-go bindings aligns with project goals for building public portfolio and supporting open-source ecosystem

5. **Future-Proof Technology:** Zenoh's unified data model (data in motion + data at rest + data in compute) provides architectural foundation for advanced use cases

### Discovery Context

This is a **pre-implementation decision** - neither Sparkplug B nor Zenoh has been implemented yet. Both protocols are planned for Phase 2 (Post-MVP) development. This is the optimal time to make technology substitutions with zero sunk cost.

### Evidence

**Supporting Evidence:**

1. **No Implementation Exists:** Both protocols require Go library development work (go-sparkplug vs zenoh-go contribution) - effort is comparable

2. **Performance Proof:** Community evidence from industrial automation practitioners:
   - 900 channels @ 100Hz on Raspberry Pi CM4
   - Network-wide key queries returning in <few ms
   - Real-time GUI with minimal CPU impact
   - Concurrent datalogging process with no performance degradation

3. **SCADA Coverage Maintained:** OPC UA (already in MVP Epic 5) provides Ignition integration path, so Zenoh substitution doesn't eliminate SCADA connectivity

4. **Library Maturity:** Both require development effort:
   - Sparkplug B: Would use paho.mqtt.golang + custom protobuf encoding (NBIRTH/NDATA/NDEATH)
   - Zenoh: Would contribute to zenoh-go bindings (active Rust reference implementation available)

5. **Ecosystem Momentum:** Growing adoption in industrial/robotics sectors, active Eclipse Foundation project

**Trade-offs Accepted:**
- ❌ Lose native Ignition Sparkplug B integration
- ✅ Keep OPC UA → Ignition integration (already in MVP)
- ✅ Gain superior performance and feature set
- ✅ Gain open-source contribution opportunity for portfolio building

---

## 2. Impact Analysis

### Epic Impact

| Epic | Impact | Notes |
|------|--------|-------|
| Epic 1 (Foundation) | None | No changes |
| Epic 2 (Modbus I/O) | None | No changes |
| Epic 3 (Task Runtime) | None | No changes |
| Epic 4 (WebUI) | None | No changes |
| Epic 5 (OPC UA) | None | No changes |
| Epic 6 (Deployment) | None | No changes |
| **Phase 2 - Protocol Expansion** | **Modified** | Replace Sparkplug B with Zenoh |

**Phase 2 Epic Changes:**
- **REMOVE:** "Sparkplug B - MQTT publisher implementation using paho.mqtt.golang with Sparkplug B protobuf encoding (NBIRTH/NDATA/NDEATH messages) for Ignition cloud and MQTT-based SCADA integration"
- **ADD:** "Zenoh Protocol - Pub/sub implementation using zenoh-go bindings for high-performance SCADA integration (900+ channels @ 100Hz capability), geo-distributed storage, and query/reply patterns"

### Artifact Impact

| Artifact | Changes Required | Status |
|----------|------------------|--------|
| PRD | Update Growth Features section | ✅ Complete |
| Architecture Doc | Update protocol integration section, directory structure | ✅ Complete |
| Epics | No changes (protocol details not specified) | ✅ Complete |
| README | Update SCADA integration description | ✅ Complete |
| UX Specification | None - protocol is internal implementation | ✅ N/A |

### MVP Impact

✅ **ZERO MVP IMPACT**

All MVP scope (Epics 1-6) remains unchanged. This change affects only Phase 2 (Post-MVP) work that has not yet been implemented.

---

## 3. Recommended Approach

### Selected Path: Direct Adjustment (Pre-Implementation)

**Viability:** ✅ **APPROVED**

This is the lowest-risk approach because:

1. ✅ **Zero Sunk Cost** - No Sparkplug B code exists to abandon
2. ✅ **Pre-Implementation Timing** - Decision made before any Phase 2 development begins
3. ✅ **Comparable Effort** - Both protocols require Go library development work
4. ✅ **Superior Technology** - Zenoh offers measurably better performance and features
5. ✅ **Portfolio Alignment** - Contributing to zenoh-go achieves open-source contribution goals
6. ✅ **Maintains Timeline** - Remains in Phase 2 slot as originally planned

**Effort Estimate:** Medium (similar to Sparkplug B implementation, different library)

**Risk Level:** Medium
- **Risks:** zenoh-go library maturity, learning curve, contribution coordination with upstream project
- **Mitigations:**
  - Active Zenoh community and responsive maintainers
  - Mature Rust reference implementation for guidance
  - OPC UA provides fallback SCADA integration path
  - Can defer Zenoh to later phase if library blockers discovered

### Technology Comparison

| Aspect | MQTT Sparkplug B | Zenoh | Impact |
|--------|------------------|-------|--------|
| Go Library | paho.mqtt.golang + custom encoding | zenoh-go (contribute) | Development effort shifts to different library |
| Broker Required | Yes (MQTT broker infrastructure) | No (peer-to-peer capable) | Simpler deployment architecture |
| Message Pattern | NBIRTH/NDATA/NDEATH (pub/sub only) | Pub/sub + storage + query/reply | More flexible integration patterns |
| Performance | Standard MQTT (100s of channels typical) | 900 channels @ 100Hz proven | Better scalability for high-throughput scenarios |
| SCADA Integration | Native Ignition Sparkplug B support | OPC UA path maintained | No loss of Ignition connectivity |
| Ecosystem | Mature MQTT, wide SCADA adoption | Growing industrial/robotics adoption | Trade stability for innovation |

---

## 4. Document Updates Summary

### PRD Changes (✅ Complete)

**Sections Updated:**
1. **Core Differentiators** - Changed protocol list from "Modbus, OPC UA, Sparkplug B, GraphQL" to "Modbus, OPC UA, Zenoh, GraphQL"
2. **Growth Features (Phase 2)** - Replaced Sparkplug B description with Zenoh description highlighting performance and features
3. **User Journey 1 (Marcus)** - Updated SCADA integration from "Sparkplug B and OPC UA" to "OPC UA"
4. **User Journey 4 (David)** - Updated commissioning from "Sparkplug B settings" to "OPC UA server settings"
5. **User Journey 5 (Andy)** - Updated library research from "paho.mqtt.golang for Sparkplug B" to "zenoh-go contribution"
6. **Journey Requirements Summary** - Updated protocol integration list
7. **Innovation Section** - Updated YAML single-definition principle to include Zenoh
8. **Risk Mitigation** - Updated protocol complexity risk section

**Key Changes:**
- All Sparkplug B references replaced with Zenoh
- Performance characteristics documented (900+ channels @ 100Hz)
- Open-source contribution angle added
- Peer-to-peer capability highlighted
- OPC UA maintained as primary SCADA integration path

### Architecture Changes (✅ Complete)

**Sections Updated:**
1. **External Library Dependencies** - Replaced paho.mqtt.golang with zenoh-go
2. **Project Structure** - Changed `internal/sparkplug/` to `internal/zenoh/`
3. **Directory Tree** - Updated package structure with zenoh session/publisher/queryable files
4. **Docusaurus Structure** - Changed `sparkplug.md` to `zenoh.md`
5. **Component Diagram** - Updated data flow from "Sparkplug MQTT" to "Zenoh Phase 2"
6. **Data Flow Description** - Changed "Sparkplug B publishes NDATA" to "Zenoh publishes variable changes"
7. **API Boundaries Table** - Updated protocol boundary from MQTT to Zenoh peer-to-peer
8. **External Integrations Table** - Changed from MQTT broker to Zenoh routers/peers/storage
9. **Documentation Categories** - Updated protocol docs list
10. **Gap Analysis** - Replaced Sparkplug B timing with Zenoh integration patterns
11. **Future Enhancements** - Added Zenoh storage and query/reply opportunities
12. **Initialization Commands** - Updated mkdir command with zenoh directory

**Key Changes:**
- Package renamed from `internal/sparkplug/` to `internal/zenoh/`
- File structure reflects Zenoh architecture (session, publisher, queryable)
- Integration patterns updated for peer-to-peer model
- Documentation roadmap updated

### README Changes (✅ Complete)

**Section Updated:**
- **SCADA Integration** - Replaced "Sparkplug B over MQTT for modern IIoT architectures" with "Optional Zenoh protocol for high-performance pub/sub (Phase 2)"
- Emphasized OPC UA as primary SCADA integration
- Added Ignition/Kepware examples to OPC UA description

### Epics Changes (✅ Complete)

**Result:** No changes required - epics document doesn't reference protocol specifics

---

## 5. PRD MVP Impact

**Status:** ✅ **NO MVP IMPACT**

**MVP Scope (Unchanged):**
- Epic 1: Foundation (Configuration & Variable Store)
- Epic 2: Modbus I/O Integration
- Epic 3: Task Runtime & API
- Epic 4: WebUI Development
- Epic 5: OPC UA Server
- Epic 6: Deployment & Documentation

**Phase 2 Scope (Modified):**
- ❌ **REMOVED:** Sparkplug B implementation
- ✅ **ADDED:** Zenoh protocol integration

**Timeline Impact:** None - Phase 2 work begins after MVP completion

---

## 6. High-Level Action Plan

### Immediate Actions (✅ Complete)

1. ✅ Update PRD Growth Features section
2. ✅ Update Architecture document protocol integration sections
3. ✅ Update README SCADA integration description
4. ✅ Verify Epics document (no changes needed)
5. ✅ Generate this Sprint Change Proposal for record-keeping

### Phase 2 Implementation Actions (Future)

When Phase 2 begins:

1. **Research zenoh-go Library:**
   - Review current zenoh-go implementation status
   - Identify contribution opportunities
   - Assess API stability and feature coverage

2. **Design Integration Pattern:**
   - Define Zenoh session lifecycle management
   - Design variable-to-key mapping strategy
   - Plan pub/sub patterns for variable updates
   - Explore storage and query/reply use cases

3. **Implement Zenoh Integration:**
   - Create `internal/zenoh/` package structure
   - Implement session management
   - Implement publisher for variable updates
   - Implement queryable for remote variable access (optional)
   - Write integration tests

4. **Contribute to zenoh-go:**
   - Submit bug reports for issues encountered
   - Contribute bug fixes or feature enhancements
   - Participate in community discussions
   - Document industrial automation use case

5. **Update Documentation:**
   - Write Zenoh protocol guide for Docusaurus
   - Create example YAML configurations
   - Document performance characteristics
   - Provide comparison guidance (when to use OPC UA vs Zenoh)

---

## 7. Agent Handoff Plan

### Handoff Recipients

| Role | Responsibility | Status |
|------|----------------|--------|
| **Scrum Master (Bob)** | Document change proposal, update tracking | ✅ Complete |
| **Product Manager** | Approve strategic direction | ✅ Approved (User: Andy) |
| **Architect** | Update architecture documents | ✅ Complete |
| **Tech Writer** | Update protocol documentation plan | 📋 Phase 2 |
| **Developer (Andy)** | Implement Zenoh integration in Phase 2 | 📋 Phase 2 |

### Success Criteria

**Approval Criteria (✅ Met):**
- ✅ User explicitly approved Zenoh substitution
- ✅ All affected documents updated
- ✅ No MVP impact confirmed
- ✅ Change proposal documented

**Implementation Criteria (Phase 2):**
- Zenoh session successfully manages variable pub/sub
- Performance benchmarks validate 100Hz+ update rates
- Integration tests verify data flow
- Documentation guides users through Zenoh configuration
- Optional: Contribution to zenoh-go upstream accepted

---

## 8. Recommendation Summary

**Change Scope:** Minor (Phase 2 only, no MVP impact)

**Recommendation:** ✅ **ZENOH SUBSTITUTION APPROVED**

**Decision Rationale:**

1. ✅ **Zero Sunk Cost** - No existing implementation to abandon
2. ✅ **Comparable Effort** - Both require Go library development work
3. ✅ **Superior Performance** - Demonstrated 900 ch @ 100Hz capability
4. ✅ **Enhanced Features** - Storage + query/reply beyond pub/sub
5. ✅ **SCADA Coverage** - OPC UA maintains Ignition integration
6. ✅ **Portfolio Value** - Open-source contribution opportunity
7. ✅ **Deployment Simplification** - Peer-to-peer removes broker requirement
8. ✅ **Future-Proof** - Unified data model supports advanced use cases

**Key Benefits:**
- Better performance for high-throughput scenarios
- Simplified deployment (no MQTT broker)
- Richer protocol features (storage, query/reply)
- Open-source contribution builds portfolio
- Active community and modern technology

**Accepted Trade-offs:**
- Lose native Sparkplug B support (niche use case)
- Keep OPC UA for standard SCADA integration
- Accept zenoh-go library maturity risk (mitigated by Rust reference)

---

## Approval

**Proposed By:** Bob (Scrum Master Agent)
**Date:** 2025-12-29

**User Approval:** ✅ **APPROVED**

**User:** Andy
**Decision:** "2, update all documents even readme if needed"

**Notes:**

User confirmed:
- No development has begun on MQTT Sparkplug B
- Willing to contribute to zenoh-go as part of project
- Main drivers: performance improvements, future-proofing
- Acceptable trade-off: lose Ignition Sparkplug B, keep OPC UA
- Phase 2 timing is appropriate (same as Sparkplug B would have been)

---

## Change Log

| Date | Change | Author |
|------|--------|--------|
| 2025-12-29 | Initial proposal created | Bob (SM Agent) |
| 2025-12-29 | PRD updated with Zenoh substitution | Bob (SM Agent) |
| 2025-12-29 | Architecture updated with Zenoh integration | Bob (SM Agent) |
| 2025-12-29 | README updated with Zenoh reference | Bob (SM Agent) |
| 2025-12-29 | Proposal approved by user | Andy |

---

Generated with [Claude Code](https://claude.com/claude-code)
**BMAD Workflow:** `.bmad/bmm/workflows/4-implementation/correct-course/`
