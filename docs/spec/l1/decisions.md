---
title: "Loom CLI L1 Design Decisions"
generated: 2025-01-03T14:30:00Z
status: draft
level: L1
---

# Loom CLI L1 Design Decisions

## Overview

This document records design decisions that fill gaps between L1 (Strategic Design) and L2 (Tactical Design) documents. These decisions emerge during L2 derivation when implementation details require choices not specified in L1.

**Format:** Each decision documents the gap, choice made, and rationale.

**Traceability:** Decisions inform derivation of L2 documents (tech-specs.md, aggregate-design.md, sequence-design.md, initial-data-model.md).

---

## Technical Integration Decisions

### DEC-L1-001: Claude CLI Integration

**Question:** How should the CLI integrate with Claude AI?

**Decision:** Use `claude` CLI tool (Claude Code) instead of direct Anthropic API calls.

**Rationale:**
- Claude Code CLI provides session management, tool use, and context handling
- Simplifies integration - no direct HTTP API calls needed
- Enables `--resume` session functionality
- Uses `CLAUDE_CODE_MAX_OUTPUT_TOKENS` for large outputs

**Alternatives considered:**
- Direct Anthropic API calls - rejected as more complex, requires API key handling
- OpenAI-compatible endpoint - rejected for Anthropic-specific features

**Source:** ai

**Affects:** l2-tech-specs.md, internal/claude/client.go

---

### DEC-L1-002: Max Output Tokens

**Question:** What maximum output token limit should be used for AI calls?

**Decision:** 100,000 tokens via `CLAUDE_CODE_MAX_OUTPUT_TOKENS` environment variable

**Rationale:**
- Large enough for comprehensive document generation
- Set via environment to override Claude Code defaults
- Prevents truncation of L1/L2/L3 documents

**Source:** ai

**Affects:** l2-tech-specs.md, internal/claude/client.go

---

## Domain Struct Decisions

### DEC-L1-003: Severity Levels

**Question:** What severity levels should ambiguities/questions have?

**Decision:** Three levels: `critical`, `important`, `minor`

**Rationale:**
- `critical` - Blocks derivation, must be resolved
- `important` - Significant impact, should be resolved
- `minor` - Nice to have, can use AI suggestion
- More actionable than generic high/medium/low

**Source:** ai

**Affects:** l1-domain-model.md, internal/domain/types.go

---

### DEC-L1-004: Decision Source Values

**Question:** What source values should decisions have?

**Decision:** Four values: `user`, `default`, `existing`, `user_accepted_suggested`

**Rationale:**
- `user` - Human explicitly provided answer
- `default` - AI suggested default (--skip-interview)
- `existing` - Loaded from previous decisions file
- `user_accepted_suggested` - User accepted AI suggestion without modification
- Enables precise audit trail

**Source:** ai

**Affects:** l1-domain-model.md, l2-aggregate-design.md

---

### DEC-L1-005: Entity Analysis Struct Design

**Question:** What fields should the Entity analysis struct contain?

**Decision:** `Name`, `MentionedAttributes`, `MentionedOperations`, `MentionedStates`

**Rationale:**
- Captures what was mentioned vs what was fully specified
- `Mentioned*` prefix indicates extraction from text, not final design
- Enables completeness analysis (find what's mentioned but not defined)
- Distinct from L1 entity definition which has `Type` and formal structure

**Source:** ai

**Affects:** l2-aggregate-design.md, internal/domain/types.go

---

### DEC-L1-006: Operation Analysis Struct Design

**Question:** What fields should the Operation analysis struct contain?

**Decision:** `Name`, `Actor`, `Trigger`, `Target`, `MentionedInputs`, `MentionedRules`

**Rationale:**
- Captures who does what, when, to what
- `Actor` identifies the user role
- `Trigger` identifies when operation occurs
- `Target` identifies affected entity
- Enables comprehensive operation analysis

**Source:** ai

**Affects:** l2-aggregate-design.md, internal/domain/types.go

---

### DEC-L1-007: Ambiguity Struct Design

**Question:** What fields should the Ambiguity struct contain?

**Decision:** `ID`, `Category`, `Subject`, `Question`, `Severity`, `SuggestedAnswer`, `Options`, `ChecklistItem`, `DependsOn`

**Rationale:**
- `Category` - entity/operation/ui classification
- `Subject` - which entity/operation it relates to
- `SuggestedAnswer` - AI's recommended answer
- `ChecklistItem` - template for L1 output
- `DependsOn` - enables conditional skip logic

**Source:** ai

**Affects:** l2-aggregate-design.md, internal/domain/types.go

---

### DEC-L1-008: Question Dependency Skip Logic

**Question:** How should questions be skipped based on previous answers?

**Decision:** `DependsOn` field with `SkipCondition` array containing `QuestionID` and `SkipIfAnswer` patterns

**Rationale:**
- Avoids asking irrelevant follow-up questions
- Example: Skip "what happens after deletion?" if user said "cannot be deleted"
- Pattern matching with string contains for flexibility
- Multiple skip conditions can be OR'd

**Source:** ai

**Affects:** l2-aggregate-design.md, l2-sequence-design.md

---

## Data Model Decisions

### DEC-L1-009: Cascade State File Version

**Question:** What version identifier should the cascade state use?

**Decision:** String version `"1.0"`

**Rationale:**
- Semantic versioning for state file format
- Enables future format migrations
- String type allows prerelease identifiers

**Source:** ai

**Affects:** initial-data-model.md, cascade.go

---

### DEC-L1-010: Input Hash Format

**Question:** How should input file hash be computed and stored?

**Decision:** SHA256 hash, first 16 characters (hex encoded)

**Rationale:**
- SHA256 is cryptographically secure
- 16 chars = 64 bits of entropy, sufficient for collision avoidance
- Short enough for display in logs/UI
- Consistent with git short hash convention

**Source:** ai

**Affects:** initial-data-model.md, cascade.go

---

### DEC-L1-011: Interview Output Fields

**Question:** What fields should interview output include?

**Decision:** `Status`, `Question`/`Group`, `Progress`, `RemainingCount`, `SkippedCount`, `Message`

**Rationale:**
- `Status` - question/group/complete/error state machine
- `Progress` - "5/23" format for human display
- `RemainingCount`/`SkippedCount` - machine-readable counts
- `Group` - for grouped interview mode

**Source:** ai

**Affects:** initial-data-model.md, internal/domain/types.go

---

### DEC-L1-012: InterviewState Extended Fields

**Question:** What additional fields should InterviewState contain beyond questions and decisions?

**Decision:** `SessionID`, `DomainModel`, `Skipped`, `InputContent`, `Complete`

**Rationale:**
- `SessionID` - unique identifier for resume capability
- `DomainModel` - cached domain model for derivation
- `Skipped` - track which questions were auto-skipped
- `InputContent` - original L0 content for context
- `Complete` - explicit completion flag

**Source:** ai

**Affects:** initial-data-model.md, l2-aggregate-design.md

---

## Validation Decisions

### DEC-L1-013: V004 Implementation Status

**Question:** Should V004 (bidirectional link validation) be implemented in initial version?

**Decision:** Deferred - marked as "skip" in initial implementation

**Rationale:**
- Bidirectional validation is complex
- Other validation rules provide sufficient coverage
- Can be added in future iteration
- Marked as "skip" not "fail" to indicate intentional deferral

**Source:** ai

**Affects:** l1-business-rules.md, validate.go

---

### DEC-L1-014: Go Version Requirement

**Question:** What Go version should be required?

**Decision:** Go 1.21 or later

**Rationale:**
- Single binary distribution (no runtime dependencies)
- Fast startup time for CLI responsiveness
- Strong typing for reliability
- Excellent concurrency support for parallel phase execution
- 1.21+ provides modern language features

**Source:** ai

**Affects:** l2-tech-specs.md, go.mod

---

### DEC-L1-015: API Retry Configuration

**Question:** What retry configuration should be used for Claude API calls?

**Decision:** 3 attempts, 2s base delay, 30s max delay with exponential backoff

**Rationale:**
- 3 attempts balances reliability vs user wait time
- 2s base delay avoids hammering rate limits
- 30s max prevents excessive waits
- Exponential backoff (2^n) is standard practice
- Retryable: rate limit, timeout, 5xx errors
- Non-retryable: auth errors, bad request

**Source:** ai

**Affects:** l2-tech-specs.md, internal/claude/retry.go

---

### DEC-L1-016: Interactive Preview Limits

**Question:** What limits should apply to interactive mode previews?

**Decision:** Max 20 lines, 80 character width

**Rationale:**
- 20 lines fits typical terminal without scrolling
- 80 characters is standard terminal width
- Truncation indicator (...) for longer content
- Enables quick review without overwhelming

**Source:** ai

**Affects:** l2-tech-specs.md, cmd/interactive.go

---

### DEC-L1-017: Editor Selection Priority

**Question:** How should the external editor be selected for interactive edit mode?

**Decision:** Priority: $EDITOR → $VISUAL → vim → nano → vi

**Rationale:**
- $EDITOR is standard Unix convention
- $VISUAL is alternative for visual editors
- Fallback chain ensures editor is always found
- vim is most common, nano is beginner-friendly
- vi is ultimate fallback (always available)

**Source:** ai

**Affects:** l2-tech-specs.md, cmd/interactive.go

---

## Summary

| ID | Decision | Source |
|----|----------|--------|
| DEC-L1-001 | Claude CLI integration (not direct API) | ai |
| DEC-L1-002 | Max output tokens 100,000 | ai |
| DEC-L1-003 | Severity levels (critical, important, minor) | ai |
| DEC-L1-004 | Decision source values (user, default, existing, user_accepted_suggested) | ai |
| DEC-L1-005 | Entity analysis struct design | ai |
| DEC-L1-006 | Operation analysis struct design | ai |
| DEC-L1-007 | Ambiguity struct design | ai |
| DEC-L1-008 | Question dependency skip logic | ai |
| DEC-L1-009 | Cascade state file version "1.0" | ai |
| DEC-L1-010 | Input hash SHA256 first 16 chars | ai |
| DEC-L1-011 | Interview output fields | ai |
| DEC-L1-012 | InterviewState extended fields | ai |
| DEC-L1-013 | V004 deferred implementation | ai |
| DEC-L1-014 | Go version 1.21+ | ai |
| DEC-L1-015 | Retry config (3 attempts, 2s base, 30s max) | ai |
| DEC-L1-016 | Preview limits (20 lines, 80 chars) | ai |
| DEC-L1-017 | Editor priority (EDITOR→VISUAL→vim→nano→vi) | ai |

**Statistics:**
- Total decisions: 17
- Source: ai: 17 (100%)

---

## Related Documents

| Level | Document | Description |
|-------|----------|-------------|
| L0 | [domain-vocabulary.md](../l0/domain-vocabulary.md) | Domain Vocabulary |
| L0 | [loom-cli.md](../l0/loom-cli.md) | User Stories |
| L0 | [nfr.md](../l0/nfr.md) | Non-Functional Requirements |
| L0 | [decisions.md](../l0/decisions.md) | L0→L1 Design Decisions |
| L1 | [domain-model.md](domain-model.md) | Domain Model |
| L1 | [bounded-context-map.md](bounded-context-map.md) | Bounded Context Map |
| L1 | [business-rules.md](business-rules.md) | Business Rules |
| L1 | [acceptance-criteria.md](acceptance-criteria.md) | Acceptance Criteria |
| L1 | This document | L1→L2 Design Decisions |
| L2 | [tech-specs.md](../l2/tech-specs.md) | Technical Specifications |
| L2 | [interface-contracts.md](../l2/interface-contracts.md) | CLI Interface Contract |
| L2 | [aggregate-design.md](../l2/aggregate-design.md) | Aggregate Design |
| L2 | [sequence-design.md](../l2/sequence-design.md) | Sequence Design |
| L2 | [initial-data-model.md](../l2/initial-data-model.md) | Data Model |
