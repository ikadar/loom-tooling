---
title: "Loom CLI Domain Model"
generated: 2025-01-03T14:30:00Z
status: draft
level: L1
---

# Loom CLI Domain Model

## Overview

This document defines the domain model for loom-cli, identifying the core entities, value objects, and their relationships.

**Traceability:** Derived from [loom-cli.md](../l0/loom-cli.md) user stories.

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

### ENT-PIPELINE: DerivationPipeline

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

### ENT-PHASE: Phase

A single step in the derivation pipeline.

| Attribute | Type | Description |
|-----------|------|-------------|
| name | PhaseName | Phase identifier (analyze, interview, derive-l1, derive-l2, derive-l3) |
| status | PhaseStatus | pending, running, completed, failed |
| timestamp | datetime | Last status change |
| checkpoint | JSON | Saved state for resume |

**Related:** US-004, US-008

### ENT-DOCUMENT: Document

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

### ENT-INTERVIEW: Interview

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

### ENT-AMBIGUITY: Ambiguity (Question)

An ambiguity requiring clarification during interview.

| Attribute | Type | Description |
|-----------|------|-------------|
| id | string | Unique identifier (AMB-{CAT}-{SUBJECT}-{NNN}) |
| category | string | entity, operation, ui |
| subject | string | Entity or operation name |
| question | string | Question text |
| severity | Severity | critical, important, minor |
| suggestedAnswer | string | AI's recommended answer (optional) |
| options | string[] | Suggested answers (optional) |
| checklistItem | string | Template for L1 output |
| dependsOn | SkipCondition[] | Conditions to auto-skip this question |

**Related:** US-002, DEC-L1-003, DEC-L1-007, DEC-L1-008

### ENT-DECISION: Decision

A recorded answer to an ambiguity question.

| Attribute | Type | Description |
|-----------|------|-------------|
| id | string | Question/ambiguity ID |
| question | string | The question text |
| answer | string | The decision made |
| decidedAt | datetime | When decision was made |
| source | AnswerSource | user, default, existing, user_accepted_suggested |
| category | string | entity, operation, ui |
| subject | string | Related entity/operation name |

**Related:** US-002, DEC-L1-004

### ENT-VALRUN: ValidationRun

A validation execution against a document set.

| Attribute | Type | Description |
|-----------|------|-------------|
| inputDir | string | Directory being validated |
| level | ValidationLevel | L1, L2, L3, or ALL |
| results | ValidationResult[] | Individual rule results |
| passed | boolean | Overall pass/fail |

**Related:** US-006

### ENT-VALRESULT: ValidationResult

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
| domain-model.md | ENT-XXX | Entities and relationships |
| bounded-context-map.md | BC-XXX | Context boundaries |
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
| test-cases.md | TC-AC-XXX-NNN-{T}{NN} | Test specs (T: P/N/B/H) |
| openapi.json | - | API specification |
| implementation-skeletons.md | SKEL-XXX-NNN | Code templates |
| feature-tickets.md | FDT-NNN | Development tasks |
| service-boundaries.md | SVC-XXX | Microservice definitions |
| event-message-design.md | EVT-XXX-NNN, CMD-XXX-NNN, INT-XXX-NNN | Events and commands |
| dependency-graph.md | DEP-XXX-NNN | Service dependencies |

---

## Related Documents

| Level | Document | Description |
|-------|----------|-------------|
| L0 | [domain-vocabulary.md](../l0/domain-vocabulary.md) | Domain Vocabulary (source for this document) |
| L0 | [loom-cli.md](../l0/loom-cli.md) | User Stories |
| L0 | [nfr.md](../l0/nfr.md) | Non-Functional Requirements |
| L1 | This document | Domain Model |
| L1 | [bounded-context-map.md](bounded-context-map.md) | Bounded Context Map |
| L1 | [business-rules.md](business-rules.md) | Business Rules |
| L1 | [acceptance-criteria.md](acceptance-criteria.md) | Acceptance Criteria |
| L0 | [decisions.md](../l0/decisions.md) | Design Decisions (interview output) |
| L2 | [interface-contracts.md](../l2/interface-contracts.md) | CLI Interface Contract |
