# Loom CLI Non-Functional Requirements

## Overview

This document defines the non-functional requirements (NFRs) for loom-cli. These requirements specify quality attributes that constrain how the system delivers its functionality.

**Level:** L0 (Foundational, Human-Provided)

---

## Performance

### NFR-PERF-001: Derivation Response Time

**Requirement:** Individual derivation phases should complete within reasonable time bounds.

| Phase | Target Time | Maximum Time |
|-------|-------------|--------------|
| analyze | < 30 seconds | 60 seconds |
| derive (L1) | < 60 seconds | 120 seconds |
| derive-l2 | < 120 seconds | 300 seconds |
| derive-l3 | < 120 seconds | 300 seconds |

**Rationale:** Long-running operations frustrate users and reduce productivity.

**Measurement:** Wall-clock time from command start to completion.

---

### NFR-PERF-002: Parallel Execution

**Requirement:** Independent L2 derivation phases must execute in parallel to minimize total derivation time.

**Constraint:** Maximum 3 concurrent API calls to respect rate limits.

**Rationale:** Sequential execution of 5 L2 phases would take 5x longer than necessary.

**Measurement:** Total L2 derivation time should be < 2x single phase time.

---

### NFR-PERF-003: Progress Feedback

**Requirement:** Long-running operations must provide progress feedback to stderr.

**Details:**
- Phase start/completion messages
- Current phase indicator
- Elapsed time for completed phases

**Rationale:** Users need visibility into progress during multi-minute operations.

---

## Reliability

### NFR-REL-001: Checkpoint and Resume

**Requirement:** Long-running operations must support checkpointing to enable resume after interruption.

**Details:**
- Cascade command saves state to `.cascade-state.json`
- L2 derivation saves phase checkpoints
- Resume skips completed phases

**Rationale:** Network failures, API errors, or user interruption should not lose completed work.

---

### NFR-REL-002: Graceful Error Handling

**Requirement:** All errors must be reported with actionable messages.

**Details:**
- Clear error descriptions (not stack traces)
- Suggestions for resolution where applicable
- Exit code 1 for all errors

**Rationale:** Users should understand what went wrong and how to fix it.

---

### NFR-REL-003: API Failure Recovery

**Requirement:** Transient API failures should be retried before failing.

**Details:**
- Retry up to 3 times with exponential backoff
- Log retry attempts to stderr
- Fail gracefully after retries exhausted

**Rationale:** Claude API may have temporary availability issues.

---

## Usability

### NFR-USE-001: Self-Documenting CLI

**Requirement:** The CLI must be self-documenting through help text.

**Details:**
- `--help` available on all commands
- Usage examples in help output
- Option descriptions with types and defaults

**Rationale:** Users should not need external documentation for basic usage.

---

### NFR-USE-002: Consistent Option Naming

**Requirement:** CLI options must follow consistent naming conventions.

**Details:**
- `--input-file` / `--input-dir` for input specification
- `--output-dir` for output directory
- `--interactive` / `-i` for interactive mode
- `--resume` for checkpoint resume
- `--dry-run` for preview mode

**Rationale:** Predictable option names reduce learning curve.

---

### NFR-USE-003: Non-Destructive Defaults

**Requirement:** Default behavior should not destroy existing work without confirmation.

**Details:**
- Warn if overwriting existing files (unless `--interactive`)
- Warn if checkpoint exists (suggest `--resume`)
- Never auto-delete user files

**Rationale:** Protect users from accidental data loss.

---

## Interoperability

### NFR-INT-001: JSON Output Option

**Requirement:** Commands that produce structured output must support JSON format.

**Details:**
- `analyze` outputs JSON to stdout
- `validate --json` outputs JSON results
- Interview outputs question/state as JSON

**Rationale:** Enables integration with CI/CD pipelines and other tools.

---

### NFR-INT-002: Exit Codes

**Requirement:** Exit codes must be meaningful and consistent.

| Code | Meaning |
|------|---------|
| 0 | Success / Interview complete |
| 1 | Error (general) |
| 100 | Interview: question pending |

**Rationale:** Scripts and automation need reliable exit code semantics.

---

### NFR-INT-003: Standard Streams

**Requirement:** Use standard streams consistently.

**Details:**
- stdout: Primary output (JSON, generated content)
- stderr: Progress messages, warnings, errors
- stdin: Not used (no interactive input in commands)

**Rationale:** Enables piping and redirection in shell scripts.

---

## Maintainability

### NFR-MNT-001: Document Format Consistency

**Requirement:** All generated documents must follow consistent formatting.

**Details:**
- YAML frontmatter (title, generated, status, level)
- Markdown with consistent heading structure
- ID patterns follow conventions
- Mermaid diagrams where applicable

**Rationale:** Consistent format enables validation and tooling.

---

### NFR-MNT-002: Traceability

**Requirement:** All derived elements must be traceable to their sources.

**Details:**
- L1 traces to L0 (user stories)
- L2 traces to L1 (domain model, AC, BR)
- L3 traces to L2 (interface contracts, tech specs)
- References are bidirectional

**Rationale:** Impact analysis and change tracking require traceability.

---

### NFR-MNT-003: Version Identification

**Requirement:** Generated documents must identify the tool version that created them.

**Details:**
- Frontmatter includes `loom-cli-version` field
- `loom-cli version` command available

**Rationale:** Debugging and compatibility tracking.

---

## Security

### NFR-SEC-001: API Key Handling

**Requirement:** API keys must be handled securely.

**Details:**
- Read from environment variable (ANTHROPIC_API_KEY)
- Never log or display API keys
- Never write API keys to generated documents

**Rationale:** Protect user credentials.

---

### NFR-SEC-002: File System Safety

**Requirement:** CLI must not access files outside intended scope.

**Details:**
- Only read from specified input paths
- Only write to specified output paths
- No execution of generated content

**Rationale:** Prevent unintended file system modifications.

---

## Portability

### NFR-PORT-001: Cross-Platform Support

**Requirement:** CLI must work on major operating systems.

**Supported platforms:**
- macOS (arm64, amd64)
- Linux (amd64)
- Windows (amd64) - optional

**Rationale:** Developers use various operating systems.

---

### NFR-PORT-002: Minimal Dependencies

**Requirement:** CLI should have minimal runtime dependencies.

**Details:**
- Single binary distribution (Go compiled)
- No external runtime required
- Only dependency: ANTHROPIC_API_KEY environment variable

**Rationale:** Easy installation and deployment.

---

## Testability

### NFR-TEST-001: Deterministic Output

**Requirement:** Given the same input and AI responses, output should be deterministic.

**Details:**
- Timestamps may vary (acceptable)
- Content structure and IDs should be consistent
- Sorting order should be stable

**Rationale:** Enables automated testing and comparison.

---

### NFR-TEST-002: Dry-Run Support

**Requirement:** Destructive operations should support dry-run mode.

**Details:**
- `sync-links --dry-run` shows changes without applying
- Cascade could support `--dry-run` in future

**Rationale:** Users can preview changes before committing.

---

## Related Documents

| Level | Document | Description |
|-------|----------|-------------|
| L0 | [l0-domain-vocabulary.md](l0-domain-vocabulary.md) | Domain Vocabulary |
| L0 | [l0-loom-cli.md](l0-loom-cli.md) | User Stories |
| L0 | This document | Non-Functional Requirements |
| L1 | [l1-domain-model.md](l1-domain-model.md) | Domain Model |
| L1 | [l1-bounded-context-map.md](l1-bounded-context-map.md) | Bounded Context Map |
| L1 | [l1-business-rules.md](l1-business-rules.md) | Business Rules (includes NFR-derived rules) |
| L1 | [l1-acceptance-criteria.md](l1-acceptance-criteria.md) | Acceptance Criteria |
| L2 | [l2-cli-interface.md](l2-cli-interface.md) | CLI Interface Contract |
