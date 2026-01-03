---
title: "Loom CLI L2 Design Decisions"
generated: 2025-01-03T14:30:00Z
status: draft
level: L2
---

# Loom CLI L2 Design Decisions

## Overview

This document records design decisions that fill gaps between L2 (Tactical Design) and L3 (Operational Design) documents. These decisions emerge during L3 derivation when implementation details require choices not specified in L2.

**Format:** Each decision documents the gap, choice made, and rationale.

**Traceability:** Decisions inform derivation of L3 documents (test-cases.md, openapi.json, implementation-skeletons.md, etc.).

---

## Test Case Decisions

### DEC-L2-001: Test Category Encoding

**Question:** What single-character codes should represent test categories?

**Decision:** Four codes: `P` (positive), `N` (negative), `B` (boundary), `H` (hallucination)

**Rationale:**
- Single character for compact ID pattern
- P/N are intuitive opposites
- B for boundary (edge cases)
- H for hallucination prevention (TDAI methodology)
- Enables pattern: `TC-AC-{DOMAIN}-{NUM}-{TYPE}{SEQ}`

**Source:** ai

**Affects:** test-cases.md, internal/generator/testcases.go

---

### DEC-L2-002: Test Generation Batch Size

**Question:** How many Acceptance Criteria should be processed per AI generation batch?

**Decision:** 5 ACs per batch

**Rationale:**
- Balance between context window usage and parallelism
- Allows sufficient context for test generation
- Prevents timeout on large AC sets
- Configurable via `ChunkSize` field

**Source:** ai

**Affects:** test-cases.md, internal/generator/testcases.go

---

### DEC-L2-003: Test Coverage Metrics Structure

**Question:** What metrics should test coverage statistics include?

**Decision:** Nested structure with `Total`, `ByCategory` (counts), and `Coverage` (ratios and flags)

**Rationale:**
- `ByCategory.Positive/Negative/Boundary/Hallucination` - raw counts per category
- `Coverage.ACsCovered` - number of ACs with tests
- `Coverage.PositiveRatio/NegativeRatio` - float ratios for V008 validation
- `Coverage.HasHallucinationTests` - boolean for V009 validation
- Enables programmatic validation

**Source:** ai

**Affects:** test-cases.md, internal/generator/testcases.go

---

## API Specification Decisions

### DEC-L2-004: OpenAPI Version

**Question:** Which OpenAPI specification version should be generated?

**Decision:** OpenAPI 3.0.3

**Rationale:**
- 3.0.x is widely supported by tooling
- 3.0.3 is stable point release
- Supports JSON Schema draft-07 for schemas
- Compatible with most code generators

**Source:** ai

**Affects:** l3-openapi-spec.md, cmd/derive_l3.go

---

### DEC-L2-005: RESTful Endpoint Patterns

**Question:** What URL patterns should generated endpoints follow?

**Decision:** Standard RESTful patterns:
- `GET /{resource}` - List
- `POST /{resource}` - Create
- `GET /{resource}/{id}` - Get one
- `PUT /{resource}/{id}` - Update
- `DELETE /{resource}/{id}` - Delete
- `POST /{resource}/{id}/{action}` - Action

**Rationale:**
- Industry standard REST conventions
- Predictable for API consumers
- Matches DDD aggregate patterns

**Source:** ai

**Affects:** l3-openapi-spec.md

---

## Implementation Skeleton Decisions

### DEC-L2-006: Service Type Enumeration

**Question:** What types should implementation skeletons support?

**Decision:** Four types: `service`, `repository`, `controller`, `handler`

**Rationale:**
- `service` - Domain service with business logic
- `repository` - Data access layer
- `controller` - API endpoint handlers
- `handler` - Event/message handlers
- Covers typical DDD/Clean Architecture layers

**Source:** ai

**Affects:** l3-implementation-skeletons.md, cmd/derive_l3.go

---

### DEC-L2-007: Domain Code Extraction

**Question:** How should domain codes be extracted from service names for ID patterns?

**Decision:** Take first 4 characters of name (after removing common suffixes), uppercase

**Rationale:**
- Remove suffixes: "Service", "Controller", "Repository"
- Take first 4 chars: "ProductService" → "PROD"
- Uppercase for consistency
- Short enough for readable IDs

**Implementation:**
```go
func extractDomainCode(name string) string {
    name = strings.TrimSuffix(name, "Service")
    name = strings.TrimSuffix(name, "Controller")
    name = strings.TrimSuffix(name, "Repository")
    if len(name) > 4 { name = name[:4] }
    return strings.ToUpper(name)
}
```

**Source:** ai

**Affects:** l3-implementation-skeletons.md, l3-dependency-graph.md, cmd/derive_l3.go

---

### DEC-L2-008: Function Signature Language

**Question:** What language should function signatures use in skeletons?

**Decision:** TypeScript-style signatures

**Rationale:**
- TypeScript is language-agnostic enough for documentation
- Type annotations are clear and readable
- Async/Promise syntax widely understood
- Example: `async createOrder(request: CreateOrderRequest): Promise<Order>`

**Source:** ai

**Affects:** l3-implementation-skeletons.md

---

## Feature Ticket Decisions

### DEC-L2-009: Feature Ticket Status Values

**Question:** What status values should feature tickets support?

**Decision:** Four statuses: `draft`, `ready`, `in_progress`, `done`

**Rationale:**
- `draft` - Initial state, under review
- `ready` - Approved for implementation
- `in_progress` - Currently being worked on
- `done` - Implementation complete
- Simple Kanban-style workflow

**Source:** ai

**Affects:** l3-feature-tickets.md, cmd/derive_l3.go

---

### DEC-L2-010: Feature Ticket Priority Values

**Question:** What priority levels should feature tickets support?

**Decision:** Four levels: `critical`, `high`, `medium`, `low`

**Rationale:**
- `critical` - Must be in MVP, blocking
- `high` - Important for launch
- `medium` - Nice to have
- `low` - Future consideration
- Standard prioritization scheme

**Source:** ai

**Affects:** l3-feature-tickets.md, cmd/derive_l3.go

---

### DEC-L2-011: Feature Ticket Complexity Values

**Question:** What complexity levels should feature tickets support?

**Decision:** Five levels with time guidance:
- `trivial` - Less than 1 day
- `simple` - 1-2 days
- `medium` - 3-5 days
- `complex` - 1-2 weeks
- `very_complex` - More than 2 weeks

**Rationale:**
- Provides rough estimation guidance
- T-shirt sizing approach (S/M/L/XL/XXL equivalent)
- Helps with sprint planning
- Time ranges are guidelines, not commitments

**Source:** ai

**Affects:** l3-feature-tickets.md, cmd/derive_l3.go

---

## Service Boundary Decisions

### DEC-L2-012: Service Dependency Types

**Question:** What types of service dependencies should be distinguished?

**Decision:** Three types: `sync`, `async`, `data`

**Rationale:**
- `sync` - Synchronous API calls (HTTP, gRPC)
- `async` - Asynchronous messaging (events, queues)
- `data` - Shared data dependency (database, cache)
- Critical for architecture decisions and failure isolation

**Source:** ai

**Affects:** l3-service-boundaries.md, cmd/derive_l3.go

---

## Event Design Decisions

### DEC-L2-013: Event Versioning Strategy

**Question:** How should domain events be versioned?

**Decision:** Explicit `version` field with semver-style string (e.g., "1.0")

**Rationale:**
- Enables schema evolution
- Consumers can handle version-specific logic
- String type allows for flexible versioning schemes
- Default initial version: "1.0"

**Source:** ai

**Affects:** l3-event-message-design.md, cmd/derive_l3.go

---

### DEC-L2-014: Event Type Taxonomy

**Question:** How should different types of events be categorized?

**Decision:** Three categories:
- **Domain Events** (EVT-) - Internal aggregate state changes
- **Commands** (CMD-) - Explicit requests to perform actions
- **Integration Events** (INT-) - Cross-boundary notifications

**Rationale:**
- Domain Events are past-tense facts (OrderCreated)
- Commands are imperative requests (CreateOrder)
- Integration Events are for external system coordination
- Clear separation of concerns

**Source:** ai

**Affects:** l3-event-message-design.md, cmd/derive_l3.go

---

## Dependency Graph Decisions

### DEC-L2-015: Mermaid Diagram Component Shapes

**Question:** What Mermaid shapes should represent different component types?

**Decision:**
- `service` → Rounded: `([name])`
- `domain_service` → Square: `[name]`
- `external` → Double square: `[[name]]`

**Rationale:**
- Visual distinction between component types
- Rounded for internal services (soft edges)
- Square for core domain (solid)
- Double square for external systems (boundary)

**Source:** ai

**Affects:** l3-dependency-graph.md, cmd/derive_l3.go

---

### DEC-L2-016: Mermaid Diagram Arrow Styles

**Question:** What Mermaid arrow styles should represent different dependency types?

**Decision:**
- `sync` → Solid arrow: `-->`
- `async` → Dashed arrow with label: `-.->|async|`
- `external` → Thick arrow with label: `==>|ext|`

**Rationale:**
- Visual distinction between dependency types
- Solid for synchronous (direct connection)
- Dashed for async (loose coupling)
- Thick for external (boundary crossing)

**Source:** ai

**Affects:** l3-dependency-graph.md, cmd/derive_l3.go

---

### DEC-L2-017: Mermaid ID Sanitization

**Question:** How should IDs be sanitized for Mermaid compatibility?

**Decision:** Remove all non-alphanumeric characters

**Implementation:**
```go
func sanitizeID(id string) string {
    result := ""
    for _, c := range id {
        if (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || (c >= '0' && c <= '9') {
            result += string(c)
        }
    }
    return result
}
```

**Rationale:**
- Mermaid doesn't support hyphens or special chars in node IDs
- Preserves readability while ensuring compatibility

**Source:** ai

**Affects:** l3-dependency-graph.md, cmd/derive_l3.go

---

## Summary

| ID | Decision | Source |
|----|----------|--------|
| DEC-L2-001 | Test category codes (P, N, B, H) | ai |
| DEC-L2-002 | Test generation batch size (5 ACs) | ai |
| DEC-L2-003 | Test coverage metrics structure | ai |
| DEC-L2-004 | OpenAPI version 3.0.3 | ai |
| DEC-L2-005 | RESTful endpoint patterns | ai |
| DEC-L2-006 | Service type enumeration | ai |
| DEC-L2-007 | Domain code extraction (first 4 chars) | ai |
| DEC-L2-008 | TypeScript function signatures | ai |
| DEC-L2-009 | Feature ticket status values | ai |
| DEC-L2-010 | Feature ticket priority values | ai |
| DEC-L2-011 | Feature ticket complexity values | ai |
| DEC-L2-012 | Service dependency types | ai |
| DEC-L2-013 | Event versioning strategy | ai |
| DEC-L2-014 | Event type taxonomy (EVT/CMD/INT) | ai |
| DEC-L2-015 | Mermaid component shapes | ai |
| DEC-L2-016 | Mermaid arrow styles | ai |
| DEC-L2-017 | Mermaid ID sanitization | ai |

**Statistics:**
- Total decisions: 17
- Source: ai: 17 (100%)

---

## Related Documents

| Level | Document | Description |
|-------|----------|-------------|
| L0 | [decisions.md](../l0/decisions.md) | Design Decisions (L0→L1) |
| L1 | [decisions.md](../l1/decisions.md) | Design Decisions (L1→L2) |
| L2 | This document | Design Decisions (L2→L3) |
| L2 | [tech-specs.md](tech-specs.md) | Technical Specifications |
| L2 | [interface-contracts.md](interface-contracts.md) | CLI Interface Contract |
| L2 | [aggregate-design.md](aggregate-design.md) | Aggregate Design |
| L3 | [test-cases.md](../l3/test-cases.md) | Test Cases Specification |
| L3 | [openapi-spec.md](../l3/openapi-spec.md) | API Specification |
| L3 | [implementation-skeletons.md](../l3/implementation-skeletons.md) | Implementation Skeletons |
| L3 | [feature-tickets.md](../l3/feature-tickets.md) | Feature Tickets |
| L3 | [service-boundaries.md](../l3/service-boundaries.md) | Service Boundaries |
| L3 | [event-message-design.md](../l3/event-message-design.md) | Event & Message Design |
| L3 | [dependency-graph.md](../l3/dependency-graph.md) | Dependency Graph |
