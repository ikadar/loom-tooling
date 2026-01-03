---
title: "Loom CLI Interface Contract"
generated: 2025-01-03T14:30:00Z
status: draft
level: L2
---

# Loom CLI Interface Contract

## Overview

This document defines the complete CLI interface specification for loom-cli. It serves as the L2 (Tactical Design) interface contract, derived from L0 user stories.

**Traceability:** Commands trace to user stories (L0) and implement acceptance criteria (L1).

## Related Documents

| Level | Document | Description |
|-------|----------|-------------|
| L0 | [domain-vocabulary.md](../l0/domain-vocabulary.md) | Domain Vocabulary |
| L0 | [loom-cli.md](../l0/loom-cli.md) | User Stories |
| L0 | [nfr.md](../l0/nfr.md) | Non-Functional Requirements |
| L1 | [domain-model.md](../l1/domain-model.md) | Domain Model |
| L1 | [bounded-context-map.md](../l1/bounded-context-map.md) | Bounded Context Map |
| L1 | [business-rules.md](../l1/business-rules.md) | Business Rules |
| L1 | [acceptance-criteria.md](../l1/acceptance-criteria.md) | Acceptance Criteria (source for this document) |
| L0 | [decisions.md](../l0/decisions.md) | Design Decisions (L0→L1) |
| L1 | [decisions.md](../l1/decisions.md) | Design Decisions (L1→L2) |
| L2 | [decisions.md](decisions.md) | Design Decisions (L2→L3) |
| L2 | [prompt-catalog.md](prompt-catalog.md) | Prompt Catalog |
| L2 | This document | CLI Interface Contract |

## Global Behavior

### Exit Codes

| Code | Meaning |
|------|---------|
| 0 | Success |
| 1 | General error |
| 100 | Interview: question available (interview command only) |

### Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `ANTHROPIC_API_KEY` | Claude API key | Required |

---

## Commands

### IC-ANL-001: analyze

**Traces to:** US-001

Analyzes user story documents to discover domain model and identify ambiguities.

```
loom-cli analyze [options]
```

**Options:**

| Option | Type | Required | Description |
|--------|------|----------|-------------|
| `--input-file <path>` | string | * | Single input file (markdown) |
| `--input-dir <path>` | string | * | Directory containing input files |
| `--decisions <path>` | string | - | Path to existing decisions.md |

\* Either `--input-file` or `--input-dir` required

**Output:**
- Writes JSON to stdout - redirect to file: `loom-cli analyze ... > analysis.json`
- Contains: entities, operations, relationships, and questions

---

### IC-INT-001: interview

**Traces to:** US-002

Conducts structured interview to resolve ambiguities from analysis.

```
loom-cli interview [options]
```

**Options:**

| Option | Type | Required | Description |
|--------|------|----------|-------------|
| `--init <path>` | string | - | Initialize interview from analysis JSON file |
| `--state <path>` | string | Yes | Path to interview state file |
| `--answer <json>` | string | - | Answer as JSON: `{"question_id":"...", "answer":"...", "source":"user"}` |
| `--grouped`, `-g` | flag | - | Show all questions at once (grouped mode) |
| `--answers <json>` | string | - | Batch answers as JSON array |

**Workflow:**
1. `--init analysis.json --state state.json` - Start new interview from analysis
2. `--state state.json` - View current question (outputs JSON)
3. `--state state.json --answer '{"question_id":"Q1", "answer":"...", "source":"user"}'` - Answer question
4. Repeat until exit code 0

**Exit Codes:**
- `0` - Interview complete, no more questions
- `1` - Error
- `100` - Question available (output contains question JSON)

---

### IC-DRV-001: derive

**Traces to:** US-003

Derives L1 (Strategic Design) documents from analysis.

```
loom-cli derive [options]
```

**Options:**

| Option | Type | Required | Description |
|--------|------|----------|-------------|
| `--output-dir <path>` | string | Yes | Output directory for L1 documents |
| `--analysis-file <path>` | string | Yes | Path to analysis JSON or interview state file |
| `--decisions <path>` | string | - | Path to decisions.md (to append new decisions) |
| `--vocabulary <path>` | string | - | Domain vocabulary file for enhanced accuracy |
| `--nfr <path>` | string | - | Non-functional requirements file |

**Output:**
- `{output-dir}/domain-model.md`
- `{output-dir}/bounded-context-map.md`
- `{output-dir}/acceptance-criteria.md`
- `{output-dir}/business-rules.md`
- `{output-dir}/decisions.md`

---

### IC-DRV-002: derive-l2

**Traces to:** US-004

Derives L2 (Tactical Design) documents from L1.

```
loom-cli derive-l2 [options]
```

**Options:**

| Option | Type | Required | Description |
|--------|------|----------|-------------|
| `--input-dir <path>` | string | Yes | Directory containing L1 documents |
| `--output-dir <path>` | string | Yes | Output directory for L2 documents |
| `--interactive`, `-i` | flag | - | Interactive approval mode |

**Execution:**
- Phases run in parallel (max 3 concurrent) where independent
- Resume capability available via `cascade --resume`

**Output:**
- `{output-dir}/tech-specs.md`
- `{output-dir}/interface-contracts.md`
- `{output-dir}/aggregate-design.md`
- `{output-dir}/sequence-design.md`
- `{output-dir}/initial-data-model.md`

---

### IC-DRV-003: derive-l3

**Traces to:** US-005

Derives L3 (Operational Design) documents from L2.

```
loom-cli derive-l3 [options]
```

**Options:**

| Option | Type | Required | Description |
|--------|------|----------|-------------|
| `--input-dir <path>` | string | Yes | Directory containing L2 documents |
| `--output-dir <path>` | string | Yes | Output directory for L3 documents |

**Output:**
- `{output-dir}/test-cases.md`
- `{output-dir}/openapi.json`
- `{output-dir}/implementation-skeletons.md`
- `{output-dir}/feature-tickets.md`
- `{output-dir}/service-boundaries.md`
- `{output-dir}/event-message-design.md`
- `{output-dir}/dependency-graph.md`
- `{output-dir}/l3-output.json` (combined JSON output)

---

### IC-VAL-001: validate

**Traces to:** US-006

Validates derived documents for consistency and completeness.

```
loom-cli validate [options]
```

**Options:**

| Option | Type | Required | Description |
|--------|------|----------|-------------|
| `--input-dir <path>` | string | Yes | Directory containing documents to validate |
| `--level <level>` | string | - | Validation level: L1, L2, L3, or ALL (default: ALL) |
| `--json` | flag | - | Output results as JSON for CI/CD integration |

**Validation Rules:**

| Rule | Description | Level |
|------|-------------|-------|
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

---

### IC-SYN-001: sync-links

**Traces to:** US-007

Automatically fixes missing bidirectional references.

```
loom-cli sync-links [options]
```

**Options:**

| Option | Type | Required | Description |
|--------|------|----------|-------------|
| `--input-dir <path>` | string | Yes | Directory containing documents |
| `--dry-run` | flag | - | Preview changes without modifying files |

---

### IC-CAS-001: cascade

**Traces to:** US-008

Runs the entire derivation pipeline with a single command.

```
loom-cli cascade [options]
```

**Options:**

| Option | Type | Required | Description |
|--------|------|----------|-------------|
| `--input-file <path>` | string | * | Single L0 input file |
| `--input-dir <path>` | string | * | Directory containing L0 input files |
| `--output-dir <path>` | string | Yes | Base output directory |
| `--skip-interview` | flag | - | Skip interview, use AI-suggested defaults |
| `--decisions <path>` | string | - | Use existing decisions file |
| `--interactive`, `-i` | flag | - | Interactive approval at each level |
| `--resume`, `-r` | flag | - | Resume from interrupted state |
| `--from <level>` | string | - | Re-derive from specific level (l1, l2, l3) |

\* Either `--input-file` or `--input-dir` required

**Output Structure:**
```
{output-dir}/
├── l1/
│   ├── domain-model.md
│   ├── bounded-context-map.md
│   ├── acceptance-criteria.md
│   ├── business-rules.md
│   └── decisions.md
├── l2/
│   ├── tech-specs.md
│   ├── interface-contracts.md
│   ├── aggregate-design.md
│   ├── sequence-design.md
│   └── initial-data-model.md
├── l3/
│   ├── test-cases.md
│   ├── openapi.json
│   ├── implementation-skeletons.md
│   ├── feature-tickets.md
│   ├── service-boundaries.md
│   ├── event-message-design.md
│   └── dependency-graph.md
└── .cascade-state.json
```

**State File:** `.cascade-state.json`

Tracks progress for resume capability:
```json
{
  "version": "1.0",
  "inputHash": "sha256...",
  "phases": {
    "analyze": { "status": "completed", "timestamp": "..." },
    "interview": { "status": "completed", "timestamp": "..." },
    "derive-l1": { "status": "completed", "timestamp": "..." },
    "derive-l2": { "status": "completed", "timestamp": "..." },
    "derive-l3": { "status": "completed", "timestamp": "..." }
  }
}
```

Phase status values: `pending`, `running`, `completed`, `failed`

---

### IC-HLP-001: help

**Traces to:** US-009

Displays help documentation.

```
loom-cli help
loom-cli --help
loom-cli -h
```

---

### IC-VER-001: version

**Traces to:** US-009

Displays version information.

```
loom-cli version
```

**Output:** `loom-cli vX.Y.Z` (e.g., `loom-cli v0.3.0`)

---

## Common Patterns

### Input Specification

Commands that accept input support two mutually exclusive options:

| Pattern | Description |
|---------|-------------|
| `--input-file <path>` | Single markdown file |
| `--input-dir <path>` | Directory of markdown files |

### Output Directory

All commands that produce output require `--output-dir`:
- Directory is created if it doesn't exist
- Existing files are overwritten

### YAML Frontmatter

All generated markdown documents include YAML frontmatter:

```yaml
---
title: Document Title
generated: 2024-01-15T10:30:00Z
status: draft
level: L1
---
```
