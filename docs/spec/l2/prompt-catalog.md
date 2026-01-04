---
title: "Loom CLI Prompt Catalog"
generated: 2025-01-03T14:30:00Z
status: draft
level: L2
---

# Loom CLI Prompt Catalog

## Overview

This document specifies all prompts used by loom-cli for AI-driven derivation. Prompts are embedded in the CLI binary using Go's `//go:embed` directive.

**Traceability:** Implements [bounded-context-map.md](../l1/bounded-context-map.md) ACL layer for Claude API integration.

**Full Prompt Content:** See [prompts/](prompts/) folder for complete prompt files.

| Prompt File | Used By | Purpose |
|-------------|---------|---------|
| [domain-discovery.md](prompts/domain-discovery.md) | analyze | Extract domain model from L0 |
| [entity-analysis.md](prompts/entity-analysis.md) | analyze | Enhance entity details |
| [operation-analysis.md](prompts/operation-analysis.md) | analyze | Enhance operation details |
| [interview.md](prompts/interview.md) | interview | Resolve ambiguities |
| [derivation.md](prompts/derivation.md) | derive | Generate AC and BR |
| [derive-domain-model.md](prompts/derive-domain-model.md) | derive | Generate L1 domain model |
| [derive-bounded-context.md](prompts/derive-bounded-context.md) | derive | Generate L1 bounded context map |
| [derive-l2.md](prompts/derive-l2.md) | derive-l2 | Combined L2 derivation |
| [derive-tech-specs.md](prompts/derive-tech-specs.md) | derive-l2 | Generate tech specs |
| [derive-interface-contracts.md](prompts/derive-interface-contracts.md) | derive-l2 | Generate interface contracts |
| [derive-aggregate-design.md](prompts/derive-aggregate-design.md) | derive-l2 | Generate aggregate design |
| [derive-sequence-design.md](prompts/derive-sequence-design.md) | derive-l2 | Generate sequence diagrams |
| [derive-data-model.md](prompts/derive-data-model.md) | derive-l2 | Generate data model |
| [derive-test-cases.md](prompts/derive-test-cases.md) | derive-l3 | Generate TDAI test cases |
| [derive-l3.md](prompts/derive-l3.md) | derive-l3 | Combined L3 derivation |
| [derive-l3-api.md](prompts/derive-l3-api.md) | derive-l3 | Generate OpenAPI spec |
| [derive-l3-skeletons.md](prompts/derive-l3-skeletons.md) | derive-l3 | Generate implementation skeletons |
| [derive-feature-tickets.md](prompts/derive-feature-tickets.md) | derive-l3 | Generate feature tickets |
| [derive-service-boundaries.md](prompts/derive-service-boundaries.md) | derive-l3 | Generate service boundaries |
| [derive-event-design.md](prompts/derive-event-design.md) | derive-l3 | Generate event/message design |
| [derive-dependency-graph.md](prompts/derive-dependency-graph.md) | derive-l3 | Generate dependency graph |

---

## Prompt Architecture

### Embedding Strategy

**File:** `prompts/prompts.go`

```go
//go:embed *.md
var promptFS embed.FS

var (
    DomainDiscovery       string // domain-discovery.md
    EntityAnalysis        string // entity-analysis.md
    OperationAnalysis     string // operation-analysis.md
    Derivation            string // derivation.md
    Interview             string // interview.md
    DeriveL2              string // derive-l2.md
    DeriveTechSpecs       string // derive-tech-specs.md
    DeriveInterfaceContracts string // derive-interface-contracts.md
    DeriveAggregateDesign string // derive-aggregate-design.md
    DeriveSequenceDesign  string // derive-sequence-design.md
    DeriveDataModel       string // derive-initial-data-model.md
    DeriveTestCases       string // derive-test-cases.md
    DeriveDomainModel     string // derive-domain-model.md
    DeriveBoundedContext  string // derive-bounded-context.md
    DeriveL3              string // derive-l3.md
    DeriveL3API           string // derive-l3-api.md
    DeriveL3Skeletons     string // derive-l3-skeletons.md
    DeriveFeatureTickets  string // derive-feature-tickets.md
    DeriveServiceBoundaries string // derive-service-boundaries.md
    DeriveEventDesign     string // derive-event-design.md
    DeriveDependencyGraph string // derive-dependency-graph.md

    // L4 Derivation prompts
    DeriveL4Architecture     string // derive-l4-architecture.md
    DeriveL4Patterns         string // derive-l4-patterns.md
    DeriveL4CodingStandards  string // derive-l4-coding-standards.md
    DeriveL4ProjectStructure string // derive-l4-project-structure.md
    DeriveL4TestingStrategy  string // derive-l4-testing-strategy.md
)
```

### Context Injection Pattern

All prompts use XML-style context injection:

```markdown
<context>
</context>
```

**Injection Method:**
```go
func buildPrompt(template string, documents ...string) string {
    context := strings.Join(documents, "\n\n---\n\n")
    return strings.Replace(template, "</context>", context + "\n</context>", 1)
}
```

**Rationale:** Anthropic best practices recommend placing context at the end of prompts for optimal attention.

---

## Embedded Checklists

**Note:** Checklists are embedded directly within prompts, not separate files.

| Checklist | Location | Purpose |
|-----------|----------|---------|
| Entity Completeness Checklist | [entity-analysis.md](prompts/entity-analysis.md) | Systematic entity attribute/state/operation analysis |
| Operation Completeness Checklist | [operation-analysis.md](prompts/operation-analysis.md) | Systematic operation input/output/rule analysis |
| Quality Checklist (Entity) | [entity-analysis.md](prompts/entity-analysis.md) | Pre-output validation for entity analysis |
| Quality Checklist (Operation) | [operation-analysis.md](prompts/operation-analysis.md) | Pre-output validation for operation analysis |
| Quality Checklist (Interview) | [interview.md](prompts/interview.md) | Pre-output validation for interview questions |

### Checklist Sections

**Entity Completeness Checklist** sections (A-H):
- A: Identity & Lifecycle
- B: Attributes & Types
- C: State Machine
- D: Operations
- E: Validation Rules
- F: Relationships
- G: Events
- H: Non-Functional Aspects

**Operation Completeness Checklist** sections (A-H):
- A: Trigger & Actor
- B: Input Validation
- C: Business Rules
- D: Processing Logic
- E: State Transitions
- F: Output & Response
- G: Error Handling
- H: Events & Side Effects

---

## Prompt Structure Standard

All prompts follow this structure:

```markdown
<role>
Expert persona with specific experience areas
Priority list (what to optimize for)
Approach description (systematic, analytical, etc.)
</role>

<task>
Clear objective statement
What to produce
</task>

<thinking_process>
Step-by-step analysis instructions
What to consider before generating output
</thinking_process>

<instructions>
Detailed requirements
Format specifications
Coverage requirements
</instructions>

<output_format>
CRITICAL REQUIREMENTS for JSON output
JSON Schema with examples
Field descriptions
</output_format>

<examples>
Named examples with description
Input → Analysis → Output pattern
Multiple complexity levels
</examples>

<self_review>
Completeness checks
Consistency checks
Format checks
Fix instructions
</self_review>

<critical_output_format>
Final reminder: PURE JSON ONLY
Start/end requirements
</critical_output_format>

<context>
</context>
```

---

## Analyze Phase Prompts

### PRM-ANL-001: Domain Discovery

**File:** `prompts/domain-discovery.md`

**Used by:** `analyze` command

**Input:** L0 user stories

**Output:** JSON with entities, operations, relationships, aggregates

**Purpose:** Extract domain model from L0 specification using DDD decision points.

**Key Features:**
- Decision Point Catalogs: EVO (Entity/Value Object), AGG (Aggregate), REF (Reference)
- Confidence levels: high, medium, low
- Interview question generation for unclear classifications

**Output Schema:**
```json
{
  "entities": [{
    "name": "string",
    "classification": "entity|value_object|unknown",
    "confidence": "high|medium|low",
    "decision_points": {
      "EVO-1": {"answer": "yes|no|unclear", "evidence": "string"}
    },
    "needs_interview": true,
    "interview_questions": ["EVO-1", "EVO-4"],
    "mentioned_attributes": ["string"],
    "mentioned_operations": ["string"],
    "mentioned_states": ["string"]
  }],
  "operations": [{
    "name": "string",
    "actor": "string",
    "trigger": "string",
    "target": "string",
    "mentioned_inputs": ["string"],
    "mentioned_rules": ["string"]
  }],
  "relationships": [{
    "from": "string",
    "to": "string",
    "type": "contains|references|belongs_to|many_to_many",
    "cardinality": "1:1|1:N|N:1|N:M",
    "confidence": "high|medium|low"
  }],
  "aggregates": [{
    "name": "string",
    "root": "string",
    "contains": ["string"],
    "confidence": "high|medium|low"
  }],
  "business_rules": ["string"],
  "interview_summary": {
    "concepts_needing_interview": ["string"],
    "total_questions": 5,
    "by_category": {"EVO": 3, "AGG": 2, "REF": 0}
  }
}
```

**Decision Points:**

| Category | ID | Decision Point | Criteria |
|----------|-----|----------------|----------|
| Entity vs VO | EVO-1 | Independent identity | Does it need to be tracked independently? → Entity |
| Entity vs VO | EVO-2 | Lifecycle independence | Can it exist without parent? → Entity |
| Entity vs VO | EVO-3 | Mutability | Need to modify while keeping identity? → Entity |
| Entity vs VO | EVO-4 | External references | Referenced from outside aggregate? → Entity |
| Entity vs VO | EVO-5 | Value equality | Equal if all attributes match? → Value Object |
| Aggregate | AGG-1 | Transactional boundary | Must be modified together atomically? → Same aggregate |
| Aggregate | AGG-2 | Consistency boundary | Must be consistent immediately? → Same aggregate |
| Aggregate | AGG-3 | Independent lifecycle | Can be created/deleted independently? → Separate aggregate |
| Aggregate | AGG-4 | Access pattern | Need to load without loading parent? → Separate aggregate |
| Reference | REF-1 | Data needs | Need full data or just ID? → Embed vs Reference |
| Reference | REF-2 | Freshness | Must always be current? → Reference |
| Reference | REF-3 | Coupling | Changes to target affect source? → Reference |

---

### PRM-ANL-002: Entity Analysis

**File:** `prompts/entity-analysis.md`

**Used by:** `analyze` command (secondary pass)

**Input:** Entities from domain discovery

**Output:** Enhanced entity details with attributes and states

---

### PRM-ANL-003: Operation Analysis

**File:** `prompts/operation-analysis.md`

**Used by:** `analyze` command (secondary pass)

**Input:** Operations from domain discovery

**Output:** Enhanced operation details with inputs, outputs, rules

---

## Interview Phase Prompts

### PRM-INT-001: Interview

**File:** `prompts/interview.md`

**Used by:** `interview` command

**Input:** Ambiguities from analysis

**Output:** Question presentation and answer processing

**Purpose:** Conduct structured interview to resolve domain model ambiguities.

**Question Categories:**
1. EVO - Entity/Value Object Classification
2. AGG - Aggregate Boundary Decisions
3. REF - Reference Type Decisions
4. AMB - General Ambiguities

**Question Format:**
```markdown
---
**[Category: Subject]** (Decision Point: {ID})

Q: {Clear question about the decision}

Options:
a) {option1} → implies {classification/outcome}
b) {option2} → implies {classification/outcome}
c) {option3}
d) Other (please specify)

Suggested default: {suggestion with rationale}
---
```

**Decision Logic:**
```
EVO Classification:
IF (EVO-1 OR EVO-2 OR EVO-3 OR EVO-4) = Yes AND EVO-5 = No:
  → ENTITY
ELSE IF EVO-5 = Yes AND (EVO-1 AND EVO-2 AND EVO-3 AND EVO-4) = No:
  → VALUE OBJECT

AGG Classification:
IF (AGG-1 OR AGG-2) = Yes:
  → SAME AGGREGATE
IF (AGG-3 OR AGG-4) = Yes:
  → SEPARATE AGGREGATE
```

---

## L1 Derivation Prompts

### PRM-DRV-001: Derivation

**File:** `prompts/derivation.md`

**Used by:** `derive` command

**Input:** Domain model + resolved decisions

**Output:** Acceptance Criteria and Business Rules

**AC Format:**
```markdown
### AC-{DOMAIN}-{NUM} – {Title}
**Given** [precondition]
**When** [action with specific inputs]
**Then** [observable outcome]

**Error Cases:**
- {condition} → {behavior}

**Traceability:**
- Source: {source reference}
- Decisions: {decision IDs used}
```

**BR Format:**
```markdown
### BR-{DOMAIN}-{NUM} – {Title}
**Rule:** [Clear statement of the constraint]
**Invariant:** [Formal condition using MUST/MUST NOT]
**Enforcement:** [Where and how it's enforced]
**Violation:** [What happens on violation]
```

---

### PRM-DRV-002: Derive Domain Model

**File:** `prompts/derive-domain-model.md`

**Used by:** `derive` command

**Output:** L1 domain-model.md content

---

### PRM-DRV-003: Derive Bounded Context

**File:** `prompts/derive-bounded-context.md`

**Used by:** `derive` command

**Output:** L1 bounded-context-map.md content

---

## L2 Derivation Prompts

### PRM-L2-001: Derive L2 (Combined)

**File:** `prompts/derive-l2.md`

**Used by:** `derive-l2` command (fallback/combined mode)

**Input:** L1 documents (AC, BR, domain model)

**Output:** Test cases and technical specifications combined

**Role:** Senior QA Architect and Technical Analyst

**Priorities:**
1. Complete Coverage - every AC and BR has derived artifacts
2. Testability - concrete, executable test cases
3. Traceability - clear links between requirements and tests
4. Precision - specific values, not vague descriptions

---

### PRM-L2-002: Derive Tech Specs

**File:** `prompts/derive-tech-specs.md`

**Used by:** `derive-l2` command (phase L2-1)

**Input:** L1 business rules, acceptance criteria

**Output:** Technical specifications

---

### PRM-L2-003: Derive Interface Contracts

**File:** `prompts/derive-interface-contracts.md`

**Used by:** `derive-l2` command (phase L2-2)

**Input:** L1 domain model, bounded context map

**Output:** Service interface contracts

---

### PRM-L2-004: Derive Aggregate Design

**File:** `prompts/derive-aggregate-design.md`

**Used by:** `derive-l2` command (phase L2-3)

**Input:** L1 domain model

**Output:** Aggregate design with Go struct definitions

---

### PRM-L2-005: Derive Sequence Design

**File:** `prompts/derive-sequence-design.md`

**Used by:** `derive-l2` command (phase L2-4)

**Input:** L1 acceptance criteria, business rules

**Output:** Sequence diagrams in Mermaid format

---

### PRM-L2-006: Derive Data Model

**File:** `prompts/derive-initial-data-model.md`

**Used by:** `derive-l2` command (phase L2-5)

**Input:** L2 aggregate design

**Output:** Initial data model with JSON schemas

---

## L3 Derivation Prompts

### PRM-L3-001: Derive Test Cases

**File:** `prompts/derive-test-cases.md`

**Used by:** `derive-l3` command (phase L3-1)

**Input:** L1 acceptance criteria

**Output:** TDAI test cases

**Role:** Principal Test Engineer with 15+ years TDD/BDD/TDAI experience

**TDAI Methodology:**

| Category | ID Suffix | Minimum per AC | Purpose |
|----------|-----------|----------------|---------|
| Positive | P01, P02 | 2 | Verify expected behavior works |
| Negative | N01, N02 | 2 | Verify error handling |
| Boundary | B01 | 1 | Test edge values |
| Hallucination | H01 | 1 | Verify NOT behaviors |

**Coverage Requirements:**
- Negative ratio >= 30% (stricter than V008's 20%)
- Every AC has at least 6 tests (2P + 2N + 1B + 1H)
- Every test has `source_quote` from AC

**Self-Review Checklist:**
- [ ] Every AC has at least 6 tests
- [ ] All four categories represented per AC
- [ ] All IDs unique and properly formatted
- [ ] All strings under 60 characters
- [ ] Negative ratio >= 30%
- [ ] Every test has meaningful source_quote
- [ ] JSON valid (no trailing commas)

---

### PRM-L3-002: Derive API Specification

**File:** `prompts/derive-l3-api.md`

**Used by:** `derive-l3` command (phase L3-2a)

**Input:** L2 tech specs, interface contracts

**Output:** OpenAPI 3.0.3 specification

---

### PRM-L3-003: Derive Implementation Skeletons

**File:** `prompts/derive-l3-skeletons.md`

**Used by:** `derive-l3` command (phase L3-2b)

**Input:** L2 tech specs, aggregate design

**Output:** Service implementation skeletons with TypeScript signatures

---

### PRM-L3-004: Derive Feature Tickets

**File:** `prompts/derive-feature-tickets.md`

**Used by:** `derive-l3` command (phase L3-3)

**Input:** L1 acceptance criteria, L2 tech specs

**Output:** Feature definition tickets

---

### PRM-L3-005: Derive Service Boundaries

**File:** `prompts/derive-service-boundaries.md`

**Used by:** `derive-l3` command (phase L3-4)

**Input:** L1 bounded context map, L2 aggregate design

**Output:** Microservice boundaries

**Role:** Microservice Architect with 12+ years experience

**Design Principles:**
1. Single Responsibility - one service, one purpose
2. Loose Coupling - minimal inter-service dependencies
3. High Cohesion - related capabilities together
4. Autonomous - services can evolve independently

**Analysis Process:**
1. Capability Identification - operations, commands, events
2. Boundary Definition - aggregates, API surface, data
3. Dependency Analysis - sync vs async, data deps
4. Contract Specification - inputs, outputs, API base

---

### PRM-L3-006: Derive Event Design

**File:** `prompts/derive-event-design.md`

**Used by:** `derive-l3` command (phase L3-5)

**Input:** L1 domain model, L2 sequence design

**Output:** Domain events, commands, integration events

---

### PRM-L3-007: Derive Dependency Graph

**File:** `prompts/derive-dependency-graph.md`

**Used by:** `derive-l3` command (phase L3-6)

**Input:** L3 service boundaries

**Output:** Service dependency graph with Mermaid diagram

---

## Prompt Quality Standards

### JSON Output Requirements

All derivation prompts enforce:

```markdown
<critical_output_format>
YOUR RESPONSE MUST BE PURE JSON ONLY.
- Start with { character immediately
- End with } character
- No text before the JSON
- No text after the JSON
- No markdown code blocks
- No explanations or summaries
</critical_output_format>
```

### String Length Constraints

| Context | Max Length | Rationale |
|---------|------------|-----------|
| Test case names | 60 chars | Fits in standard terminal width |
| Service descriptions | 80 chars | Single line readability |
| Error messages | 100 chars | User-friendly display |

### Self-Review Pattern

Every prompt includes self-review instructions:

```markdown
<self_review>
After generating output, verify these criteria. If any fail, fix before outputting:

COMPLETENESS CHECK:
- [specific completeness criteria]

CONSISTENCY CHECK:
- [specific consistency criteria]

FORMAT CHECK:
- Is JSON valid (no trailing commas)?
- Does output start with { character?

If issues found, fix before outputting.
</self_review>
```

---

## Prompt Versioning

| Prompt | Version | Last Updated | Breaking Changes |
|--------|---------|--------------|------------------|
| domain-discovery.md | 1.0 | 2024-01 | Initial |
| derive-test-cases.md | 1.1 | 2024-01 | Added source_quote requirement |
| derive-service-boundaries.md | 1.0 | 2024-01 | Initial |

**Version Policy:**
- Major version: Breaking schema changes
- Minor version: New optional fields or improved instructions

---

## Related Documents

| Level | Document | Description |
|-------|----------|-------------|
| L0 | [domain-vocabulary.md](../l0/domain-vocabulary.md) | Domain Vocabulary |
| L0 | [loom-cli.md](../l0/loom-cli.md) | User Stories |
| L0 | [decisions.md](../l0/decisions.md) | Design Decisions (L0→L1) |
| L1 | [bounded-context-map.md](../l1/bounded-context-map.md) | Bounded Context Map |
| L1 | [decisions.md](../l1/decisions.md) | Design Decisions (L1→L2) |
| L2 | [decisions.md](decisions.md) | Design Decisions (L2→L3) |
| L2 | [tech-specs.md](tech-specs.md) | Technical Specifications (Claude integration) |
| L2 | This document | Prompt Catalog |
