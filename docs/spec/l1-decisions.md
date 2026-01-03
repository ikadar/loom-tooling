# Loom CLI Design Decisions

## Overview

This document records the design decisions made during L1 derivation that are not explicitly specified in the L0 foundational documents. These decisions fill gaps, resolve ambiguities, and make architectural choices necessary for implementation.

**Format:** Each decision is documented as if it were the result of a structured interview question, with the decision rationale and alternatives considered.

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

**Decision:** Three sources: `user`, `ai`, `existing`

**Rationale:**
- `user` - human explicitly provided the answer
- `ai` - AI suggested default (used with `--skip-interview`)
- `existing` - loaded from previous decisions file
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

**Affects:** l1-acceptance-criteria.md, l2-cli-interface.md

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

**Affects:** l1-acceptance-criteria.md, l2-cli-interface.md

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

**Affects:** l1-acceptance-criteria.md, l2-cli-interface.md

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

**Affects:** l1-acceptance-criteria.md, l2-cli-interface.md

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

**Statistics:**
- Total decisions: 18
- Source: user: 2 (11%)
- Source: ai: 16 (89%)

---

## Related Documents

| Level | Document | Description |
|-------|----------|-------------|
| L0 | [l0-domain-vocabulary.md](l0-domain-vocabulary.md) | Domain Vocabulary |
| L0 | [l0-loom-cli.md](l0-loom-cli.md) | User Stories |
| L0 | [l0-nfr.md](l0-nfr.md) | Non-Functional Requirements |
| L1 | [l1-domain-model.md](l1-domain-model.md) | Domain Model |
| L1 | [l1-bounded-context-map.md](l1-bounded-context-map.md) | Bounded Context Map |
| L1 | [l1-business-rules.md](l1-business-rules.md) | Business Rules |
| L1 | [l1-acceptance-criteria.md](l1-acceptance-criteria.md) | Acceptance Criteria |
| L1 | This document | Design Decisions |
| L2 | [l2-cli-interface.md](l2-cli-interface.md) | CLI Interface Contract |
