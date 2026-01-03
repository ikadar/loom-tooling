# Loom CLI - AI-Assisted Documentation Derivation Tool

## Overview

A command-line tool that uses AI to derive structured software documentation from user stories. The tool follows a multi-level derivation process:

- **L0**: User stories and requirements (input)
- **L1**: Strategic Design (Domain Model, Acceptance Criteria, Business Rules)
- **L2**: Tactical Design (Tech Specs, Contracts, Aggregates, Sequences)
- **L3**: Operational Design (Test Cases, OpenAPI, Skeletons, Events)

## User Stories

### US-001: Analyze User Stories
As a software architect, I want to analyze user story documents so that I can discover the domain model and identify ambiguities that need clarification.

The analysis should:
- Extract entities, operations, and relationships from the text
- Identify missing information or unclear requirements
- Produce a structured JSON output for further processing

### US-002: Conduct Structured Interview
As a software architect, I want to answer questions about ambiguities so that decisions are recorded and can inform the derivation process.

The interview should:
- Present one question at a time
- Record answers with source attribution (user, AI, existing)
- Allow skipping questions
- Support grouped question mode for efficiency (`--grouped`)
- Support batch answers for automation (`--answers` JSON array)

### US-003: Derive Strategic Design (L1)
As a software architect, I want to derive L1 documents from analyzed input so that I have a consistent domain model and requirements specification.

L1 derivation supports:
- Optional domain vocabulary file (`--vocabulary`) to enhance domain model accuracy
- Optional non-functional requirements file (`--nfr`) to include NFRs in business rules

L1 outputs include:
- Domain Model (entities, value objects, relationships)
- Bounded Context Map (context boundaries, integrations)
- Acceptance Criteria (Given/When/Then format)
- Business Rules (invariants, validations)
- Decisions (recorded from interview)

### US-004: Derive Tactical Design (L2)
As a software developer, I want to derive L2 documents from L1 so that I have technical specifications for implementation.

L2 derivation supports:
- Parallel execution of independent phases (max 3 concurrent)
- Checkpoint system for resuming interrupted derivations (`--resume`)
- Interactive approval mode (`--interactive`)

L2 outputs include:
- Tech Specs (implementation details per business rule)
- Interface Contracts (API operations, DTOs, errors)
- Aggregate Design (DDD aggregates, invariants, behaviors)
- Sequence Design (interaction flows, participants)
- Initial Data Model (tables, relationships, indexes)

### US-005: Derive Operational Design (L3)
As a QA engineer, I want to derive L3 documents from L2 so that I have test cases and implementation guides.

L3 outputs include:
- Test Cases (TDAI methodology: positive, negative, boundary, hallucination prevention)
- OpenAPI Specification (JSON format)
- Implementation Skeletons (service structure, function signatures)
- Feature Tickets (development tasks with acceptance criteria)
- Service Boundaries (microservice definitions)
- Event/Message Design (domain events, commands, integration events)
- Dependency Graph (service dependencies visualization)

### US-006: Validate Generated Documents
As a software architect, I want to validate derived documents so that I can ensure consistency and completeness.

Validation rules:
- V001: Every document has IDs
- V002: IDs follow expected patterns
- V003: All references point to existing IDs
- V004: Bidirectional links are consistent (use `sync-links` to fix)
- V005: Every AC has at least 1 test case
- V006: Every Entity has an aggregate
- V007: Every Service has an interface contract
- V008: Negative test ratio >= 20%
- V009: Every AC has hallucination prevention test
- V010: No duplicate IDs

Validation supports:
- Level-specific validation (`--level L1|L2|L3|ALL`)
- JSON output for CI/CD integration (`--json`)

### US-007: Sync Bidirectional Links
As a software architect, I want to automatically fix missing bidirectional references so that traceability is maintained across documents.

When document A references document B, document B should reference back to A. The tool should detect and fix missing back-references.

Sync-links supports:
- Dry-run mode to preview changes (`--dry-run`)

### US-008: Run Full Cascade Derivation
As a software architect, I want to run the entire derivation pipeline with a single command so that I can quickly generate all documentation levels.

The cascade should:
- Run analyze, interview (optional), derive, derive-l2, derive-l3 sequentially
- Create l1/, l2/, l3/ output directories
- Support resume from interrupted state
- Support re-derivation from a specific level

### US-009: Get Help and Version Information
As a user, I want to see help documentation and version information so that I understand how to use the CLI.

The CLI should:
- Display help with `--help`, `-h`, or `help` command
- Show version with `version` command
- Provide usage examples for each command

## Business Rules

- Input files must be markdown format
- All derived documents must have YAML frontmatter (title, generated, status, level)
- IDs must follow naming conventions (AC-XXX-NNN, BR-XXX-NNN, etc.)
- Traceability must be bidirectional (if A→B then B→A)
- Test cases must include hallucination prevention tests (TDAI)
- Negative test ratio must be at least 20%

## Non-Functional Requirements

- CLI should provide clear progress feedback during long operations
- Parallel execution where possible (L2 phases run concurrently)
- Checkpoint support for resuming interrupted derivations
- JSON output option for programmatic integration
- Interactive approval mode for reviewing generated content

## Technical Context

- Written in Go
- Uses Claude API for AI-powered derivation
- Markdown output with YAML frontmatter
- OpenAPI 3.0 JSON for API specs
- Mermaid diagrams for visualizations

## Related Documents

| Level | Document | Description |
|-------|----------|-------------|
| L0 | [l0-domain-vocabulary.md](l0-domain-vocabulary.md) | Domain Vocabulary (concepts, definitions) |
| L0 | This document | User Stories |
| L0 | [l0-nfr.md](l0-nfr.md) | Non-Functional Requirements |
| L1 | [l1-domain-model.md](l1-domain-model.md) | Domain Model (entities, relationships) |
| L1 | [l1-bounded-context-map.md](l1-bounded-context-map.md) | Bounded Context Map |
| L1 | [l1-business-rules.md](l1-business-rules.md) | Business Rules (invariants, constraints) |
| L1 | [l1-acceptance-criteria.md](l1-acceptance-criteria.md) | Acceptance Criteria (Given/When/Then) |
| L2 | [l2-cli-interface.md](l2-cli-interface.md) | CLI Interface Contract (commands, options, exit codes) |
