# Loom CLI Domain Vocabulary

## Overview

This document defines the core domain concepts, terms, and their relationships for loom-cli. It serves as the foundational vocabulary from which the domain model is derived.

**Level:** L0 (Foundational, Human-Provided)

---

## Core Concepts

### Derivation

- **Definition:** The process of automatically generating structured documentation from higher-level inputs using AI.
- **Notes:** Derivation flows from L0 → L1 → L2 → L3. Each level depends on previous levels.
- **Related terms:** Derivation Level, Derivation Pipeline, Phase

### Derivation Level

- **Definition:** A classification of documentation by abstraction level and generation order.
- **Allowed values:** L0, L1, L2, L3
- **Notes:**
  - L0: Foundational (human-provided) - User Stories, Vocabulary, NFRs
  - L1: Strategic Design (AI-derived) - Domain Model, AC, BR
  - L2: Tactical Design (AI-derived) - Interface Contracts, Aggregates, Sequences
  - L3: Operational Design (AI-derived) - Test Cases, OpenAPI, Skeletons
- **Related terms:** Derivation, Document

### Derivation Pipeline

- **Definition:** The orchestrated sequence of phases that transforms L0 input into L1, L2, and L3 outputs.
- **Notes:** Can be run as individual commands or as a single cascade command.
- **Related terms:** Phase, Cascade

### Phase

- **Definition:** A single step in the derivation pipeline that performs one transformation.
- **Allowed values:** analyze, interview, derive-l1, derive-l2, derive-l3
- **Notes:** Phases have status (pending, running, completed, failed) and can be checkpointed for resume.
- **Related terms:** Derivation Pipeline, Checkpoint

### Cascade

- **Definition:** A single-command execution of the full derivation pipeline (L0 → L1 → L2 → L3).
- **Notes:** Supports skip-interview mode, resume from checkpoint, and re-derivation from specific level.
- **Related terms:** Derivation Pipeline, Phase

---

## Document Types

### Document

- **Definition:** A markdown file containing structured specification content at a specific derivation level.
- **Notes:** All documents have YAML frontmatter with metadata (title, generated, status, level).
- **Related terms:** Derivation Level, Frontmatter, Reference

### User Story

- **Definition:** A high-level requirement in "As a... I want... So that..." format.
- **ID Pattern:** US-XXX-NNN (e.g., US-RFQ-001)
- **Level:** L0
- **Notes:** Human-provided input that drives all derivations.
- **Related terms:** Acceptance Criteria, Business Rule

### Acceptance Criteria

- **Definition:** A testable requirement in Given/When/Then (BDD) format derived from a user story.
- **ID Pattern:** AC-XXX-NNN (e.g., AC-RFQ-001)
- **Level:** L1
- **Notes:** Each user story should have 3-7 acceptance criteria.
- **Related terms:** User Story, Test Case

### Business Rule

- **Definition:** An invariant, constraint, or policy that the system must enforce.
- **ID Pattern:** BR-XXX-NNN (e.g., BR-ORD-001)
- **Level:** L1
- **Notes:** Derived from domain vocabulary constraints and user story requirements.
- **Related terms:** User Story, Domain Vocabulary, Validation Rule

### Domain Model

- **Definition:** A structured representation of entities, value objects, and their relationships.
- **ID Pattern:** E-XXX-NNN for entities (e.g., E-DRV-001)
- **Level:** L1
- **Notes:** Derived from domain vocabulary terms and relationships.
- **Related terms:** Entity, Value Object, Aggregate

### Interface Contract

- **Definition:** An API specification defining operations, inputs, outputs, and errors for a service.
- **ID Pattern:** IC-XXX-NNN for contracts, OP-XXX-NNN for operations
- **Level:** L2
- **Notes:** Derived from domain model and acceptance criteria.
- **Related terms:** Operation, Service

### Test Case

- **Definition:** A specification for verifying that an acceptance criterion is correctly implemented.
- **ID Pattern:** TC-XXX-NNN (e.g., TC-AC-RFQ-001-P01)
- **Level:** L3
- **Notes:** Generated using TDAI methodology (positive, negative, boundary, hallucination prevention).
- **Related terms:** Acceptance Criteria, TDAI

---

## Analysis & Interview

### Analysis

- **Definition:** The process of extracting entities, operations, relationships, and ambiguities from L0 input.
- **Notes:** Produces a JSON structure with discovered elements and questions for clarification.
- **Related terms:** Question, Entity, Operation

### Question

- **Definition:** An ambiguity identified during analysis that requires human clarification.
- **ID Pattern:** Q1, Q2, Q3...
- **Notes:** Questions have status (pending, answered, skipped) and may have suggested options.
- **Related terms:** Analysis, Interview, Decision

### Interview

- **Definition:** A structured Q&A session to resolve ambiguities identified during analysis.
- **Notes:** Presents one question at a time (default) or grouped. Records answers with source attribution.
- **Related terms:** Question, Decision

### Decision

- **Definition:** A recorded answer to an ambiguity question with source attribution.
- **Notes:** Source can be: user (human input), ai (AI-suggested default), existing (from decisions file).
- **Related terms:** Interview, Question, Answer Source

### Answer Source

- **Definition:** The origin of a decision's answer.
- **Allowed values:** user, ai, existing
- **Notes:**
  - user: Human provided the answer interactively
  - ai: AI suggested a default (used with --skip-interview)
  - existing: Loaded from a previous decisions.md file
- **Related terms:** Decision, Interview

---

## Validation

### Validation

- **Definition:** The process of checking derived documents for consistency, completeness, and correctness.
- **Notes:** Can be scoped to specific levels (L1, L2, L3, ALL).
- **Related terms:** Validation Rule, Violation

### Validation Rule

- **Definition:** A specific check performed during validation.
- **Allowed values:** V001-V010
- **Notes:**
  - V001: Every document has IDs
  - V002: IDs follow expected patterns
  - V003: All references point to existing IDs
  - V004: Bidirectional links are consistent
  - V005: Every AC has at least 1 test case
  - V006: Every Entity has an aggregate
  - V007: Every Service has an interface contract
  - V008: Negative test ratio >= 20%
  - V009: Every AC has hallucination prevention test
  - V010: No duplicate IDs
- **Related terms:** Validation, Violation

### Violation

- **Definition:** A failed validation rule with details about the specific issue.
- **Notes:** Contains document ID, rule violated, and descriptive message.
- **Related terms:** Validation Rule

### Traceability

- **Definition:** The ability to trace any element back to its source and forward to its dependents.
- **Notes:** Implemented via bidirectional references between documents.
- **Related terms:** Reference, Sync-Links

### Reference

- **Definition:** A link from one document element to another, enabling traceability.
- **Notes:** Must be bidirectional (if A references B, B must reference A).
- **Related terms:** Traceability, Sync-Links

### Sync-Links

- **Definition:** The operation of automatically adding missing bidirectional references.
- **Notes:** Can run in dry-run mode to preview changes.
- **Related terms:** Reference, Traceability, Validation Rule V004

---

## Technical Concepts

### Checkpoint

- **Definition:** A saved state that enables resuming an interrupted derivation.
- **Notes:** Stored as JSON files (.cascade-state.json for cascade, phase-specific for L2).
- **Related terms:** Phase, Resume

### Resume

- **Definition:** Continuing a derivation from the last successful checkpoint.
- **Notes:** Skips completed phases and continues from the first incomplete one.
- **Related terms:** Checkpoint, Phase

### Frontmatter

- **Definition:** YAML metadata at the beginning of a markdown document.
- **Required fields:** title, generated, status, level
- **Notes:** All derived documents must have frontmatter.
- **Related terms:** Document

### TDAI

- **Definition:** Test-Driven AI Development - A methodology for generating comprehensive test cases.
- **Notes:** Requires positive tests, negative tests (≥20%), boundary tests, and hallucination prevention tests.
- **Related terms:** Test Case, Hallucination Prevention Test

### Hallucination Prevention Test

- **Definition:** A test case that verifies the system does NOT implement behaviors not specified in requirements.
- **Notes:** Critical for AI-generated code to prevent unintended features.
- **Related terms:** TDAI, Test Case

---

## Related Documents

| Level | Document | Description |
|-------|----------|-------------|
| L0 | This document | Domain Vocabulary |
| L0 | [l0-loom-cli.md](l0-loom-cli.md) | User Stories |
| L0 | [l0-nfr.md](l0-nfr.md) | Non-Functional Requirements |
| L1 | [l1-domain-model.md](l1-domain-model.md) | Domain Model (derived from vocabulary) |
| L1 | [l1-bounded-context-map.md](l1-bounded-context-map.md) | Bounded Context Map |
| L1 | [l1-business-rules.md](l1-business-rules.md) | Business Rules |
| L1 | [l1-acceptance-criteria.md](l1-acceptance-criteria.md) | Acceptance Criteria |
| L1 | [l1-decisions.md](l1-decisions.md) | Design Decisions |
| L2 | [l2-cli-interface.md](l2-cli-interface.md) | CLI Interface Contract |
