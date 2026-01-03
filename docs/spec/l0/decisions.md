---
title: "Loom CLI Design Decisions"
generated: 2025-01-03T14:30:00Z
status: draft
level: L0
---

# Loom CLI Design Decisions

## Overview

This document records design decisions that fill gaps and resolve ambiguities in the foundational L0 documents (User Stories, Domain Vocabulary, NFRs). These decisions are the output of the interview phase and serve as input for L1 derivation.

**Format:** Each decision is documented as the result of a structured interview question, with the decision rationale and alternatives considered.

**Traceability:** Decisions inform derivation of [domain-model.md](../l1/domain-model.md), [bounded-context-map.md](../l1/bounded-context-map.md), [business-rules.md](../l1/business-rules.md), and [acceptance-criteria.md](../l1/acceptance-criteria.md).

---

## Architectural Decisions

### DEC-001: Four Bounded Contexts

**Question:** How should the domain be partitioned into bounded contexts?

**Decision:** Four bounded contexts: Derivation, Interview, Validation, CLI.

**Rationale:**
- **Derivation Context**: Core domain - document generation pipeline
- **Interview Context**: Separate because it has distinct lifecycle and can be skipped
- **Validation Context**: Separate because it's stateless and can run independently
- **CLI Context**: Orchestrator that coordinates other contexts

**Alternatives considered:**
- 3 contexts (merge CLI into Derivation) - rejected because CLI has distinct responsibilities
- 2 contexts (merge Interview and Validation into Derivation) - rejected for SRP violation
- 5+ contexts (separate each command) - rejected as over-engineering

**Source:** ai

**Affects:** l1-bounded-context-map.md, l1-domain-model.md

---

### DEC-002: Interview as Separate Context

**Question:** Should Interview be part of the Derivation Context or separate?

**Decision:** Interview is a separate bounded context (BC-INT).

**Rationale:**
- Interview has its own lifecycle (can be skipped entirely with `--skip-interview`)
- Interview state persists independently
- Interview can be conducted incrementally across multiple CLI invocations
- Clean separation enables future UI-based interview

**Alternatives considered:**
- Merge into Derivation - rejected because interview is optional
- Merge into CLI - rejected because interview has domain logic

**Source:** ai

**Affects:** l1-bounded-context-map.md

---

### DEC-003: Anti-Corruption Layer for Claude API

**Question:** How should the system integrate with the external Claude API?

**Decision:** Use Anti-Corruption Layer (ACL) pattern in Derivation Context.

**Rationale:**
- Claude API returns raw JSON that needs translation to domain objects
- API response format may change independently of domain
- ACL isolates domain from external API changes
- Enables retry logic and error handling at boundary

**Alternatives considered:**
- Direct API calls from domain - rejected to avoid coupling
- Separate Integration Context - rejected as over-engineering for single external system

**Source:** ai

**Affects:** l1-bounded-context-map.md

---

## Domain Model Decisions

### DEC-004: Phase Names

**Question:** What should the derivation phases be named?

**Decision:** Five phases: `analyze`, `interview`, `derive-l1`, `derive-l2`, `derive-l3`

**Rationale:**
- `analyze` - clear verb describing L0 analysis
- `interview` - established term for Q&A sessions
- `derive-l1/l2/l3` - consistent naming with level suffix
- Verb forms indicate actions (not nouns like "analysis")

**Alternatives considered:**
- `parse`, `question`, `generate-l1/l2/l3` - rejected for inconsistency
- Numbered phases (phase-1, phase-2) - rejected as less descriptive
- Combined phases (derive covers all levels) - rejected for granularity needs

**Source:** ai

**Affects:** l1-domain-model.md, l1-acceptance-criteria.md

---

### DEC-005: Phase Status Values

**Question:** What status values should a Phase have?

**Decision:** Four statuses: `pending`, `running`, `completed`, `failed`

**Rationale:**
- `pending` - not yet started
- `running` - currently executing
- `completed` - finished successfully
- `failed` - terminated with error

**Alternatives considered:**
- Add `skipped` - rejected because skipping is handled by not creating the phase
- Add `paused` - rejected as not needed (resume recreates from checkpoint)
- Add `cancelled` - rejected as same outcome as failed

**Source:** ai

**Affects:** l1-domain-model.md

---

### DEC-006: Answer Source Values

**Question:** How should decision sources be categorized?

**Decision:** Four sources: `user`, `default`, `existing`, `user_accepted_suggested`

**Rationale:**
- `user` - human explicitly provided the answer
- `default` - AI suggested default (used with `--skip-interview`)
- `existing` - loaded from previous decisions file
- `user_accepted_suggested` - user accepted AI suggestion without modification
- Enables audit trail and distinguishes human vs AI decisions

**Alternatives considered:**
- Binary (human/auto) - rejected as doesn't capture reuse from file
- More granular (user-interactive, user-batch, ai-default, ai-confident) - rejected as over-engineering

**Source:** user

**Affects:** l1-domain-model.md, l1-business-rules.md

---

### DEC-007: Reference Types

**Question:** What types of cross-document references should be supported?

**Decision:** Three types: `traces-to`, `implements`, `tests`

**Rationale:**
- `traces-to` - general traceability (e.g., AC traces to US)
- `implements` - technical implementation (e.g., API implements AC)
- `tests` - test coverage (e.g., TC tests AC)
- Semantic distinction enables targeted impact analysis

**Alternatives considered:**
- Single generic "references" type - rejected as loses semantic meaning
- More types (validates, depends-on, extends) - deferred for future

**Source:** ai

**Affects:** l1-domain-model.md

---

## Interface Decisions

### DEC-008: Exit Code 100 for Pending Interview

**Question:** What exit code should indicate "interview has pending questions"?

**Decision:** Exit code 100

**Rationale:**
- 0 = success/complete
- 1 = error
- 100 = special state (not error, not complete)
- High number avoids collision with common Unix exit codes
- Easy to check in scripts: `if [ $? -eq 100 ]; then`

**Alternatives considered:**
- Exit code 2 - rejected as commonly means "misuse of shell command"
- Exit code 3 - could work but 100 is more distinctive
- Negative exit code - rejected as not portable

**Source:** ai

**Affects:** l1-acceptance-criteria.md, interface-contracts.md

---

### DEC-009: JSON Output for Analysis

**Question:** What format should the analyze command output?

**Decision:** JSON to stdout

**Rationale:**
- JSON is machine-readable and widely supported
- stdout enables piping to other tools
- Structured format captures entities, operations, relationships, questions
- Progress/status messages go to stderr to keep stdout clean

**Alternatives considered:**
- Markdown output - rejected as harder to parse programmatically
- YAML output - rejected as JSON is more universal
- File output - rejected as stdout is more flexible

**Source:** ai

**Affects:** l1-acceptance-criteria.md, interface-contracts.md

---

### DEC-010: State File for Interview

**Question:** How should interview state be persisted?

**Decision:** JSON state file specified by `--state` flag

**Rationale:**
- Enables incremental interview across multiple CLI invocations
- File can be inspected/edited by humans if needed
- JSON format consistent with analysis output
- Explicit path gives user control over location

**Alternatives considered:**
- In-memory only - rejected as doesn't support incremental use
- Database - rejected as overkill for CLI tool
- Hidden file in output dir - rejected as less explicit

**Source:** ai

**Affects:** l1-acceptance-criteria.md, interface-contracts.md

---

### DEC-011: Answer JSON Format

**Question:** What format should interview answers use?

**Decision:** JSON object with `question_id`, `answer`, `source` fields

**Rationale:**
- `question_id` - identifies which question is being answered
- `answer` - the actual response text
- `source` - attribution (user, ai, existing)
- Single object per answer enables atomic operations

**Alternatives considered:**
- Positional arguments - rejected as error-prone
- Interactive prompt - rejected as conflicts with automation
- Multiple answers per call - supported via `--answers` array

**Source:** ai

**Affects:** l1-acceptance-criteria.md, interface-contracts.md

---

## Business Rule Decisions

### DEC-012: Maximum 3 Concurrent API Calls

**Question:** How many parallel API calls should be allowed during L2 derivation?

**Decision:** Maximum 3 concurrent calls

**Rationale:**
- Claude API has rate limits
- 3 concurrent balances speed vs rate limit risk
- L2 has 5 independent phases, 3 concurrent = ~2 rounds
- Higher concurrency risks 429 errors

**Alternatives considered:**
- 1 (sequential) - rejected as too slow
- 5 (all parallel) - rejected as likely to hit rate limits
- Configurable - deferred for future (could add `--concurrency` flag)

**Source:** ai

**Affects:** l1-business-rules.md (BR-DRV-002)

---

### DEC-013: Default Overwrite Behavior

**Question:** Should existing output files be overwritten by default?

**Decision:** Yes, overwrite without confirmation (unless `--interactive` mode)

**Rationale:**
- Derivation is meant to be repeatable
- Users re-run derivation to regenerate from updated input
- Prompting every time would be annoying for common workflow
- Interactive mode exists for careful review

**Alternatives considered:**
- Never overwrite - rejected as breaks re-derivation workflow
- Always prompt - rejected as too verbose
- Backup before overwrite - deferred for future consideration

**Source:** ai

**Affects:** l1-business-rules.md (BR-IO-004)

---

### DEC-014: One Question at a Time Default

**Question:** Should interview present all questions at once or one at a time?

**Decision:** One question at a time (default), grouped mode available via `--grouped`

**Rationale:**
- One at a time focuses attention
- Prevents overwhelming user with many questions
- Each answer can inform subsequent questions
- Grouped mode available for experienced users

**Alternatives considered:**
- All at once default - rejected as overwhelming
- Configurable batch size - rejected as over-engineering
- Interactive prompt - rejected as conflicts with CLI nature

**Source:** user

**Affects:** l1-business-rules.md (BR-INT-001)

---

### DEC-015: YAML Frontmatter Fields

**Question:** What fields should be included in document frontmatter?

**Decision:** Four required fields: `title`, `generated`, `status`, `level`

**Rationale:**
- `title` - human-readable document name
- `generated` - timestamp for tracking freshness
- `status` - workflow state (draft/review/approved)
- `level` - derivation level (L1/L2/L3)
- Minimal set that enables validation and tooling

**Alternatives considered:**
- More fields (author, version, hash) - deferred as not immediately needed
- Fewer fields - rejected as loses important metadata
- No frontmatter - rejected as loses machine-readability

**Source:** ai

**Affects:** l1-business-rules.md (BR-DOC-001)

---

## ID Pattern Decisions

### DEC-016: ID Format Pattern

**Question:** What format should document element IDs follow?

**Decision:** `{PREFIX}-{CONTEXT}-{NUMBER}` format

**Rationale:**
- PREFIX identifies document type (AC, BR, TC, etc.)
- CONTEXT identifies domain area (ORD, CUST, etc.)
- NUMBER provides uniqueness within type+context
- Human-readable and sortable

**Alternatives considered:**
- UUIDs - rejected as not human-readable
- Sequential numbers only - rejected as no type/context info
- Hierarchical (AC/ORD/001) - slashes problematic in URLs/anchors

**Source:** ai

**Affects:** l1-business-rules.md (BR-DOC-002)

---

### DEC-017: ID Number Format

**Question:** How should the NUMBER portion of IDs be formatted?

**Decision:** 3 digits, zero-padded (e.g., 001, 042, 999)

**Rationale:**
- Zero-padding ensures consistent sorting
- 3 digits supports up to 999 items per type+context
- Consistent width improves readability in lists
- Standard convention in technical documentation

**Alternatives considered:**
- No padding (1, 42, 999) - rejected as sorts incorrectly
- 2 digits - rejected as may be insufficient
- 4 digits - rejected as overkill for most projects

**Source:** ai

**Affects:** l1-business-rules.md (BR-DOC-002)

---

### DEC-018: Context Abbreviation Length

**Question:** How long should context abbreviations be?

**Decision:** 3-4 uppercase letters (e.g., ORD, CUST, RFQ)

**Rationale:**
- Short enough to keep IDs readable
- Long enough to be meaningful
- Uppercase for visual distinction
- Consistent with common abbreviation practices

**Alternatives considered:**
- 2 letters - rejected as too ambiguous
- Full words - rejected as makes IDs too long
- Mixed case - rejected for consistency

**Source:** ai

**Affects:** l1-business-rules.md (BR-DOC-002)

---

## L1 Document ID Prefixes

### DEC-019: Entity ID Prefix

**Question:** What prefix should entity IDs use?

**Decision:** `ENT-{CONTEXT}` format (e.g., ENT-ORDER, ENT-CUST)

**Rationale:**
- `ENT` is more explicit than single letter `E`
- Context suffix without number allows flexible entity naming
- Matches validation regex in implementation

**Alternatives considered:**
- `E-XXX-NNN` - rejected, single letter less readable
- `ENTITY-XXX` - rejected, too long

**Source:** ai

**Affects:** l1-domain-model.md, validate.go

---

### DEC-020: Bounded Context ID Prefix

**Question:** What prefix should bounded context IDs use?

**Decision:** `BC-{CONTEXT}` format (e.g., BC-DRV, BC-INT, BC-VAL, BC-CLI)

**Rationale:**
- `BC` clearly indicates Bounded Context
- Short context suffix (3-4 chars) keeps IDs readable
- No number needed as contexts are few

**Source:** ai

**Affects:** l1-bounded-context-map.md

---

## L2 Document ID Prefixes

### DEC-021: L2 Document ID Prefixes

**Question:** What prefixes should L2 document elements use?

**Decision:** The following prefixes:
- `TS-XXX-NNN` - Tech Specs
- `IC-XXX-NNN` - Interface Contracts
- `OP-XXX-NNN` - Operations (within Interface Contracts)
- `AGG-XXX-NNN` - Aggregates
- `SEQ-XXX-NNN` - Sequences
- `TBL-XXX-NNN` - Tables (Data Model)

**Rationale:**
- Short, meaningful prefixes
- Consistent with L1 pattern
- Each document type has unique prefix

**Source:** ai

**Affects:** l1-domain-model.md (Document Type Catalog), interface-contracts.md

---

## L3 Document ID Prefixes

### DEC-022: L3 Document ID Prefixes

**Question:** What prefixes should L3 document elements use?

**Decision:** The following prefixes:
- `TC-AC-XXX-NNN-{TYPE}{NN}` - Test Cases (TYPE: P=positive, N=negative, B=boundary, H=hallucination)
- `EVT-XXX-NNN` - Domain Events
- `CMD-XXX-NNN` - Commands
- `INT-XXX-NNN` - Integration Events
- `SVC-{CONTEXT}` - Services
- `SKEL-XXX-NNN` - Implementation Skeletons
- `FDT-NNN` - Feature Definition Tickets
- `DEP-XXX-NNN` - Dependencies

**Rationale:**
- Test case ID encodes source AC and test type
- Event/Command distinction supports CQRS patterns
- Service IDs use context name without number

**Source:** ai

**Affects:** l1-domain-model.md (Document Type Catalog), validate.go

---

## CLI Flag Decisions

### DEC-023: Analyze Command Flags

**Question:** What flags should the analyze command support?

**Decision:**
- `--input-file <path>` - Single input file
- `--input-dir <path>` - Directory of input files
- `--decisions <path>` - Existing decisions file to filter resolved ambiguities

**Rationale:**
- Consistent with other commands (input-file/input-dir pattern)
- Decisions flag enables incremental analysis

**Source:** ai

**Affects:** interface-contracts.md

---

### DEC-024: Interview Command Flags

**Question:** What flags should the interview command support?

**Decision:**
- `--init <path>` - Initialize from analysis JSON
- `--state <path>` - Interview state file path
- `--answer <json>` - Single answer JSON
- `--answers <json>` - Batch answers JSON array
- `--grouped`, `-g` - Show all questions at once

**Rationale:**
- `--init` separates initialization from continuation
- `--answers` (plural) enables batch mode for automation
- `-g` short flag for common grouped mode

**Source:** ai

**Affects:** interface-contracts.md

---

### DEC-025: Derive Command Flags

**Question:** What flags should the derive (L1) command support?

**Decision:**
- `--output-dir <path>` - Output directory (required)
- `--analysis-file <path>` - Analysis JSON or interview state
- `--decisions <path>` - Existing decisions to append
- `--vocabulary <path>` - Domain vocabulary for accuracy
- `--nfr <path>` - NFR file for business rules

**Rationale:**
- `--analysis-file` name indicates it accepts analysis output
- Optional vocabulary/nfr enhance derivation quality

**Source:** ai

**Affects:** interface-contracts.md

---

### DEC-026: Cascade Re-derive Flag

**Question:** How should cascade support re-derivation from a specific level?

**Decision:** `--from <level>` flag with values: l1, l2, l3

**Rationale:**
- Intuitive "from this level" semantics
- Lowercase level names match directory structure
- Enables partial re-derivation without full restart

**Source:** ai

**Affects:** interface-contracts.md

---

### DEC-027: Short Flags

**Question:** Which flags should have short versions?

**Decision:**
- `-i` for `--interactive`
- `-g` for `--grouped`
- `-r` for `--resume`

**Rationale:**
- Only most commonly used flags get short versions
- Avoids short flag collision
- Matches Unix convention

**Source:** ai

**Affects:** interface-contracts.md

---

## Validation Decisions

### DEC-028: Validation Rule Level Mapping

**Question:** Which validation rules apply to which levels?

**Decision:**
| Rule | Level |
|------|-------|
| V001-V004, V010 | ALL |
| V006, V007 | L2+ |
| V005, V008, V009 | L3 |

**Rationale:**
- Structural rules (V001-V004, V010) apply universally
- Entity/Service completeness (V006-V007) requires L2 docs
- Test coverage (V005, V008, V009) requires L3 docs

**Source:** ai

**Affects:** l1-business-rules.md, validate.go

---

## Document Format Decisions

### DEC-029: Frontmatter Status Values

**Question:** What status values should documents have?

**Decision:** Three values: `draft`, `review`, `approved`

**Rationale:**
- `draft` - initial generated state
- `review` - under human review
- `approved` - ready for use
- Simple workflow, extensible later

**Source:** ai

**Affects:** l1-business-rules.md (BR-DOC-001)

---

### DEC-030: Version String Format

**Question:** What format should the version command output?

**Decision:** `loom-cli vX.Y.Z` format (e.g., `loom-cli v0.3.0`)

**Rationale:**
- Includes tool name for clarity in logs
- `v` prefix is common convention
- Semantic versioning (X.Y.Z)

**Source:** ai

**Affects:** interface-contracts.md, root.go

---

## Interactive Mode Decisions

### DEC-031: Interactive Mode Actions

**Question:** What actions should interactive mode support?

**Decision:** Five actions:
- `[A]pprove` - Accept and write file
- `[E]dit` - Open in $EDITOR
- `[R]egenerate` - Generate again with AI
- `[S]kip` - Skip this file
- `[Q]uit` - Exit derivation

**Rationale:**
- Approve is default (Enter key)
- Edit enables manual refinement
- Regenerate retries AI generation
- Skip allows partial derivation
- Quit provides clean exit

**Alternatives considered:**
- Only approve/skip/abort - rejected as Edit and Regenerate are valuable
- Numbered options - rejected, letter keys are faster

**Source:** ai

**Affects:** l1-acceptance-criteria.md, interactive.go

---

## DDD Pattern Decisions

### DEC-032: Context Relationship Patterns

**Question:** What DDD relationship patterns should be used?

**Decision:** Four patterns:
- **Upstream/Downstream** - Data flow direction
- **Conformist** - Downstream adopts upstream's model
- **ACL (Anti-Corruption Layer)** - Translate external model
- **Shared Kernel** - Common types across contexts

**Rationale:**
- Standard DDD strategic patterns
- Cover common integration scenarios
- ACL specifically for Claude API integration

**Source:** ai

**Affects:** l1-bounded-context-map.md

---

### DEC-033: Shared Kernel Types

**Question:** What types should be in the Shared Kernel?

**Decision:**
- `DerivationLevel` enum (L0, L1, L2, L3)
- Document ID patterns
- `Frontmatter` value object

**Rationale:**
- These types used across all bounded contexts
- Centralized definition prevents inconsistency
- Minimal shared surface area

**Source:** ai

**Affects:** l1-bounded-context-map.md, l1-domain-model.md

---

## Summary

| ID | Decision | Source |
|----|----------|--------|
| DEC-001 | Four bounded contexts | ai |
| DEC-002 | Interview as separate context | ai |
| DEC-003 | ACL for Claude API | ai |
| DEC-004 | Phase names (analyze, interview, derive-l1/l2/l3) | ai |
| DEC-005 | Phase status (pending, running, completed, failed) | ai |
| DEC-006 | Answer source (user, ai, existing) | user |
| DEC-007 | Reference types (traces-to, implements, tests) | ai |
| DEC-008 | Exit code 100 for pending interview | ai |
| DEC-009 | JSON output for analysis | ai |
| DEC-010 | State file for interview | ai |
| DEC-011 | Answer JSON format | ai |
| DEC-012 | Max 3 concurrent API calls | ai |
| DEC-013 | Default overwrite behavior | ai |
| DEC-014 | One question at a time default | user |
| DEC-015 | YAML frontmatter fields | ai |
| DEC-016 | ID format pattern | ai |
| DEC-017 | ID number format (3 digits, zero-padded) | ai |
| DEC-018 | Context abbreviation length (3-4 chars) | ai |
| DEC-019 | Entity ID prefix (ENT-XXX) | ai |
| DEC-020 | Bounded Context ID prefix (BC-XXX) | ai |
| DEC-021 | L2 document ID prefixes | ai |
| DEC-022 | L3 document ID prefixes | ai |
| DEC-023 | Analyze command flags | ai |
| DEC-024 | Interview command flags | ai |
| DEC-025 | Derive command flags | ai |
| DEC-026 | Cascade re-derive flag (--from) | ai |
| DEC-027 | Short flags (-i, -g, -r) | ai |
| DEC-028 | Validation rule level mapping | ai |
| DEC-029 | Frontmatter status values | ai |
| DEC-030 | Version string format | ai |
| DEC-031 | Interactive mode actions | ai |
| DEC-032 | Context relationship patterns | ai |
| DEC-033 | Shared Kernel types | ai |

**Statistics:**
- Total decisions: 33
- Source: user: 2 (6%)
- Source: ai: 31 (94%)

---

## Related Documents

| Level | Document | Description |
|-------|----------|-------------|
| L0 | [domain-vocabulary.md](domain-vocabulary.md) | Domain Vocabulary |
| L0 | [loom-cli.md](loom-cli.md) | User Stories |
| L0 | [nfr.md](nfr.md) | Non-Functional Requirements |
| L0 | This document | Design Decisions (interview output) |
| L1 | [domain-model.md](../l1/domain-model.md) | Domain Model |
| L1 | [bounded-context-map.md](../l1/bounded-context-map.md) | Bounded Context Map |
| L1 | [business-rules.md](../l1/business-rules.md) | Business Rules |
| L1 | [acceptance-criteria.md](../l1/acceptance-criteria.md) | Acceptance Criteria |
| L2 | [interface-contracts.md](../l2/interface-contracts.md) | CLI Interface Contract |
