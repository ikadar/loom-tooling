# Loom CLI Domain Model

## Overview

This document defines the domain model for loom-cli, identifying the core entities, value objects, and their relationships.

**Traceability:** Derived from [l0-loom-cli.md](l0-loom-cli.md) user stories.

---

## Bounded Contexts

### 1. Derivation Context

Core domain for document derivation pipeline.

**Entities:**
- DerivationPipeline
- Phase
- Document

**Value Objects:**
- DerivationLevel
- PhaseStatus

### 2. Interview Context

Handles ambiguity resolution through structured questioning.

**Entities:**
- Interview
- Question
- Decision

**Value Objects:**
- AnswerSource
- QuestionStatus

### 3. Validation Context

Ensures document consistency and completeness.

**Entities:**
- ValidationRun
- ValidationResult

**Value Objects:**
- ValidationRule
- ValidationLevel

---

## Entities

### E-DRV-001: DerivationPipeline

The orchestrator for multi-level document derivation.

| Attribute | Type | Description |
|-----------|------|-------------|
| id | string | Unique pipeline identifier |
| inputFiles | InputFile[] | Source L0 documents |
| outputDir | string | Base output directory |
| currentPhase | Phase | Currently executing phase |
| state | PipelineState | Current execution state |

**Behaviors:**
- Start derivation from L0 input
- Resume from checkpoint
- Re-derive from specific level

**Related:** US-008

### E-DRV-002: Phase

A single step in the derivation pipeline.

| Attribute | Type | Description |
|-----------|------|-------------|
| name | PhaseName | Phase identifier (analyze, interview, derive-l1, derive-l2, derive-l3) |
| status | PhaseStatus | pending, running, completed, failed |
| timestamp | datetime | Last status change |
| checkpoint | JSON | Saved state for resume |

**Related:** US-004, US-008

### E-DOC-001: Document

A generated specification document at any derivation level.

| Attribute | Type | Description |
|-----------|------|-------------|
| id | DocumentId | Unique document identifier |
| level | DerivationLevel | L0, L1, L2, or L3 |
| title | string | Document title |
| content | markdown | Document body |
| frontmatter | Frontmatter | YAML metadata |
| references | Reference[] | Links to other documents |

**Invariants:**
- Must have YAML frontmatter
- All IDs must follow naming conventions
- References must be bidirectional

**Related:** US-003, US-004, US-005

### E-INT-001: Interview

A structured interview session for resolving ambiguities.

| Attribute | Type | Description |
|-----------|------|-------------|
| id | string | Interview session identifier |
| questions | Question[] | Questions to be answered |
| currentIndex | int | Current question position |
| decisions | Decision[] | Recorded decisions |

**Behaviors:**
- Initialize from analysis
- Present next question
- Record answer
- Mark complete when all answered

**Related:** US-002

### E-INT-002: Question

An ambiguity requiring clarification.

| Attribute | Type | Description |
|-----------|------|-------------|
| id | QuestionId | Unique question identifier (Q1, Q2, ...) |
| text | string | Question text |
| context | string | Background information |
| options | string[] | Suggested answers (optional) |
| status | QuestionStatus | pending, answered, skipped |

**Related:** US-002

### E-INT-003: Decision

A recorded answer to an ambiguity question.

| Attribute | Type | Description |
|-----------|------|-------------|
| questionId | QuestionId | Reference to question |
| answer | string | The decision made |
| source | AnswerSource | user, ai, existing |
| timestamp | datetime | When decision was made |

**Related:** US-002

### E-VAL-001: ValidationRun

A validation execution against a document set.

| Attribute | Type | Description |
|-----------|------|-------------|
| inputDir | string | Directory being validated |
| level | ValidationLevel | L1, L2, L3, or ALL |
| results | ValidationResult[] | Individual rule results |
| passed | boolean | Overall pass/fail |

**Related:** US-006

### E-VAL-002: ValidationResult

Result of a single validation rule.

| Attribute | Type | Description |
|-----------|------|-------------|
| rule | ValidationRule | The rule checked |
| passed | boolean | Whether rule passed |
| violations | Violation[] | List of violations if failed |

**Related:** US-006

---

## Value Objects

### DerivationLevel

```
enum DerivationLevel {
  L0  // User Stories (input)
  L1  // Strategic Design
  L2  // Tactical Design
  L3  // Operational Design
}
```

### PhaseStatus

```
enum PhaseStatus {
  pending
  running
  completed
  failed
}
```

### PhaseName

```
enum PhaseName {
  analyze
  interview
  derive-l1
  derive-l2
  derive-l3
}
```

### AnswerSource

```
enum AnswerSource {
  user      // Human provided answer
  ai        // AI suggested default
  existing  // From previous decisions file
}
```

### QuestionStatus

```
enum QuestionStatus {
  pending
  answered
  skipped
}
```

### ValidationLevel

```
enum ValidationLevel {
  L1
  L2
  L3
  ALL
}
```

### ValidationRule

```
enum ValidationRule {
  V001  // Every document has IDs
  V002  // IDs follow expected patterns
  V003  // All references point to existing IDs
  V004  // Bidirectional links are consistent
  V005  // Every AC has at least 1 test case
  V006  // Every Entity has an aggregate
  V007  // Every Service has an interface contract
  V008  // Negative test ratio >= 20%
  V009  // Every AC has hallucination prevention test
  V010  // No duplicate IDs
}
```

### Frontmatter

```
value object Frontmatter {
  title: string
  generated: datetime
  status: "draft" | "review" | "approved"
  level: DerivationLevel
}
```

### Reference

```
value object Reference {
  targetId: string      // e.g., "AC-ORD-001"
  targetFile: string    // e.g., "acceptance-criteria.md"
  type: "traces-to" | "implements" | "tests"
}
```

### InputFile

```
value object InputFile {
  path: string
  format: "markdown"
  content: string
}
```

### Violation

```
value object Violation {
  documentId: string
  message: string
  location: string
}
```

---

## Relationships

```
DerivationPipeline
    │
    ├── 1:N ──► Phase
    │             │
    │             └── generates ──► Document (per level)
    │
    └── 0:1 ──► Interview
                  │
                  ├── 1:N ──► Question
                  │
                  └── 1:N ──► Decision

ValidationRun
    │
    └── 1:N ──► ValidationResult
                  │
                  └── 0:N ──► Violation

Document
    │
    └── N:M ──► Document (via Reference)
```

---

## Document Type Catalog

### L1 Documents

| Document | ID Pattern | Description |
|----------|------------|-------------|
| domain-model.md | E-XXX-NNN | Entities and relationships |
| bounded-context-map.md | BC-XXX-NNN | Context boundaries |
| acceptance-criteria.md | AC-XXX-NNN | Given/When/Then criteria |
| business-rules.md | BR-XXX-NNN | Business invariants |
| decisions.md | DEC-NNN | Recorded decisions |

### L2 Documents

| Document | ID Pattern | Description |
|----------|------------|-------------|
| tech-specs.md | TS-XXX-NNN | Implementation specs |
| interface-contracts.md | IC-XXX-NNN, OP-XXX-NNN | API contracts |
| aggregate-design.md | AGG-XXX-NNN | DDD aggregates |
| sequence-design.md | SEQ-XXX-NNN | Interaction flows |
| initial-data-model.md | TBL-XXX-NNN | Database schema |

### L3 Documents

| Document | ID Pattern | Description |
|----------|------------|-------------|
| test-cases.md | TC-XXX-NNN | Test specifications |
| openapi.json | - | API specification |
| implementation-skeletons.md | SK-XXX-NNN | Code templates |
| feature-tickets.md | FT-XXX-NNN | Development tasks |
| service-boundaries.md | SVC-XXX-NNN | Microservice definitions |
| event-message-design.md | EVT-XXX-NNN | Domain events |
| dependency-graph.md | - | Service dependencies |

---

## Related Documents

| Level | Document | Description |
|-------|----------|-------------|
| L0 | [l0-domain-vocabulary.md](l0-domain-vocabulary.md) | Domain Vocabulary (source for this document) |
| L0 | [l0-loom-cli.md](l0-loom-cli.md) | User Stories |
| L0 | [l0-nfr.md](l0-nfr.md) | Non-Functional Requirements |
| L1 | This document | Domain Model |
| L1 | [l1-business-rules.md](l1-business-rules.md) | Business Rules |
| L1 | [l1-acceptance-criteria.md](l1-acceptance-criteria.md) | Acceptance Criteria |
| L2 | [l2-cli-interface.md](l2-cli-interface.md) | CLI Interface Contract |
