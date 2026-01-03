# Loom CLI Business Rules

## Overview

This document defines the business rules that govern loom-cli behavior. These rules are invariants that must be enforced by the implementation.

**Traceability:** Derived from [l0-loom-cli.md](l0-loom-cli.md) user stories.

---

## Input/Output Rules

### BR-IO-001: Markdown Input Format

**Rule:** All input files must be in markdown format.

**Rationale:** Markdown is human-readable and supports structured content extraction.

**Enforcement:** Validate file extension (.md) and content parsing.

**Related:** US-001, US-003, US-008

---

### BR-IO-002: Input Source Specification

**Rule:** Commands accepting input must receive exactly one of: `--input-file` (single file) or `--input-dir` (directory).

**Rationale:** Prevents ambiguous input specification.

**Enforcement:** CLI argument validation; error if both or neither provided.

**Applies to:** analyze, cascade

**Related:** US-001, US-008

---

### BR-IO-003: Output Directory Creation

**Rule:** Output directories are created automatically if they do not exist.

**Rationale:** Reduces manual setup steps for users.

**Enforcement:** `os.MkdirAll()` before writing files.

**Related:** US-003, US-004, US-005, US-008

---

### BR-IO-004: Output Overwrite Behavior

**Rule:** Existing files in output directories are overwritten without confirmation (unless `--interactive` mode).

**Rationale:** Enables repeatable derivation runs.

**Enforcement:** Direct file write; interactive mode prompts for confirmation.

**Related:** US-003, US-004, US-005

---

## Document Format Rules

### BR-DOC-001: YAML Frontmatter Required

**Rule:** All generated markdown documents must include YAML frontmatter with: title, generated (timestamp), status, level.

**Rationale:** Enables document metadata extraction and validation.

**Enforcement:** Formatter functions always emit frontmatter.

**Example:**
```yaml
---
title: Acceptance Criteria
generated: 2024-01-15T10:30:00Z
status: draft
level: L1
---
```

**Related:** US-003, US-004, US-005

---

### BR-DOC-002: ID Naming Conventions

**Rule:** Document element IDs must follow the pattern `{PREFIX}-{CONTEXT}-{NUMBER}` where:
- PREFIX: Document type (AC, BR, TC, TS, IC, OP, AGG, SEQ, etc.)
- CONTEXT: Domain context abbreviation (3-4 uppercase letters)
- NUMBER: Sequential number (3 digits, zero-padded)

**Valid examples:** AC-ORD-001, BR-CUST-002, TC-AC-ORD-001-P01

**Invalid examples:** AC001, ac-ord-001, AC-ORDER-1

**Rationale:** Enables consistent cross-referencing and validation.

**Enforcement:** V002 validation rule.

**Related:** US-006

---

### BR-DOC-003: Bidirectional Traceability

**Rule:** If document A references document B, then document B must reference back to A.

**Rationale:** Ensures complete traceability and impact analysis capability.

**Enforcement:** V004 validation rule; `sync-links` command to fix.

**Example:** If AC-ORD-001 is tested by TC-AC-ORD-001-P01, then TC-AC-ORD-001-P01 must list AC-ORD-001 in its "Tests" field.

**Related:** US-006, US-007

---

## Interview Rules

### BR-INT-001: One Question at a Time

**Rule:** Default interview mode presents exactly one question at a time.

**Rationale:** Focuses user attention and ensures deliberate answers.

**Enforcement:** Interview state tracks `currentIndex`; only current question output.

**Override:** `--grouped` flag shows all questions at once.

**Related:** US-002

---

### BR-INT-002: Answer Source Attribution

**Rule:** Every decision must record its source: `user` (human input), `ai` (AI-suggested default), or `existing` (from decisions file).

**Rationale:** Enables audit trail and distinguishes human decisions from AI suggestions.

**Enforcement:** Decision struct requires source field.

**Related:** US-002

---

### BR-INT-003: Interview Completion Signal

**Rule:** Interview command returns exit code 0 only when all questions are answered; exit code 100 indicates pending questions.

**Rationale:** Enables scripted interview workflows.

**Enforcement:** Exit code based on unanswered question count.

**Related:** US-002

---

### BR-INT-004: Skip Interview Option

**Rule:** The `--skip-interview` flag bypasses interview phase, using AI-suggested defaults for all decisions.

**Rationale:** Enables fast, non-interactive derivation for automation.

**Enforcement:** Cascade command skips interview phase; derive uses AI defaults.

**Related:** US-008

---

## Derivation Rules

### BR-DRV-001: Sequential Level Derivation

**Rule:** Derivation levels must be generated in order: L0 → L1 → L2 → L3. Each level depends on the previous.

**Rationale:** Higher levels are derived from lower levels; skipping creates inconsistencies.

**Enforcement:** Commands require input from previous level.

**Related:** US-003, US-004, US-005, US-008

---

### BR-DRV-002: Parallel Phase Execution

**Rule:** Independent L2 phases (tech-specs, interface-contracts, aggregates, sequences, data-model) may execute in parallel, with maximum 3 concurrent to respect API rate limits.

**Rationale:** Reduces total derivation time while avoiding rate limiting.

**Enforcement:** ParallelExecutor with concurrency limit.

**Related:** US-004

---

### BR-DRV-003: Checkpoint for Resume

**Rule:** Long-running operations must save checkpoints that enable resuming after interruption.

**Rationale:** Prevents loss of work on failures; enables incremental processing.

**Enforcement:** `.cascade-state.json` and phase-level checkpoints.

**Related:** US-004, US-008

---

### BR-DRV-004: Interactive Approval Mode

**Rule:** When `--interactive` flag is set, each generated file must be previewed and approved before being written.

**Rationale:** Enables human review of AI-generated content.

**Enforcement:** Workflow prompts for approval; skipped files not written.

**Related:** US-004

---

## Validation Rules

### BR-VAL-001: Validation Rule Catalog

**Rule:** The following validation rules must be implemented:

| Rule | Description | Applicable Level |
|------|-------------|------------------|
| V001 | Every document has IDs | ALL |
| V002 | IDs follow expected patterns | ALL |
| V003 | All references point to existing IDs | ALL |
| V004 | Bidirectional links are consistent | ALL |
| V005 | Every AC has at least 1 test case | L3 |
| V006 | Every Entity has an aggregate | L2 |
| V007 | Every Service has an interface contract | L2 |
| V008 | Negative test ratio >= 20% | L3 |
| V009 | Every AC has hallucination prevention test | L3 |
| V010 | No duplicate IDs | ALL |

**Related:** US-006

---

### BR-VAL-002: Level-Specific Validation

**Rule:** Validation can be scoped to specific levels (L1, L2, L3) or ALL.

**Rationale:** Enables incremental validation during development.

**Enforcement:** `--level` flag filters applicable rules.

**Related:** US-006

---

### BR-VAL-003: JSON Output for CI/CD

**Rule:** Validation results can be output as JSON for programmatic consumption.

**Rationale:** Enables integration with CI/CD pipelines.

**Enforcement:** `--json` flag changes output format.

**Related:** US-006

---

## Test Quality Rules (TDAI)

### BR-TDAI-001: Hallucination Prevention Tests

**Rule:** Every acceptance criterion must have at least one hallucination prevention test case.

**Rationale:** AI-generated code may include behaviors not specified in requirements; negative tests catch this.

**Enforcement:** V009 validation rule.

**Related:** US-005, US-006

---

### BR-TDAI-002: Negative Test Ratio

**Rule:** At least 20% of test cases must be negative tests (testing what should NOT happen).

**Rationale:** Ensures thorough error handling and boundary testing.

**Enforcement:** V008 validation rule.

**Related:** US-005, US-006

---

## Environment Rules

### BR-ENV-001: API Key Required

**Rule:** The `ANTHROPIC_API_KEY` environment variable must be set for AI-powered operations.

**Rationale:** API authentication is required for Claude API calls.

**Enforcement:** Error on missing key before API call.

**Related:** US-001, US-002, US-003, US-004, US-005

---

## Related Documents

| Level | Document | Description |
|-------|----------|-------------|
| L0 | [l0-domain-vocabulary.md](l0-domain-vocabulary.md) | Domain Vocabulary |
| L0 | [l0-loom-cli.md](l0-loom-cli.md) | User Stories |
| L0 | [l0-nfr.md](l0-nfr.md) | Non-Functional Requirements (source for NFR-derived rules) |
| L1 | [l1-domain-model.md](l1-domain-model.md) | Domain Model |
| L1 | [l1-bounded-context-map.md](l1-bounded-context-map.md) | Bounded Context Map |
| L1 | This document | Business Rules |
| L1 | [l1-acceptance-criteria.md](l1-acceptance-criteria.md) | Acceptance Criteria |
| L1 | [l1-decisions.md](l1-decisions.md) | Design Decisions |
| L2 | [l2-cli-interface.md](l2-cli-interface.md) | CLI Interface Contract |
