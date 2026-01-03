# Loom CLI Acceptance Criteria

## Overview

This document defines the acceptance criteria for loom-cli in Given/When/Then format. Each criterion is traceable to user stories and will be verified by test cases.

**Traceability:** Derived from [l0-loom-cli.md](l0-loom-cli.md) user stories.

---

## US-001: Analyze User Stories

### AC-ANL-001: Analyze Single File

**Given** a valid markdown file containing user stories
**And** `ANTHROPIC_API_KEY` is set
**When** I run `loom-cli analyze --input-file story.md`
**Then** JSON output is written to stdout
**And** the output contains `entities`, `operations`, `relationships`, and `questions` arrays

**Related BR:** BR-IO-001, BR-ENV-001

---

### AC-ANL-002: Analyze Directory

**Given** a directory containing multiple markdown files
**When** I run `loom-cli analyze --input-dir ./stories/`
**Then** all markdown files in the directory are analyzed
**And** combined JSON output is written to stdout

**Related BR:** BR-IO-002

---

### AC-ANL-003: Analyze with Existing Decisions

**Given** a valid markdown file and an existing decisions.md file
**When** I run `loom-cli analyze --input-file story.md --decisions decisions.md`
**Then** existing decisions are incorporated into the analysis
**And** fewer questions are generated for already-decided topics

---

### AC-ANL-004: Analyze Invalid Input

**Given** a non-existent input file
**When** I run `loom-cli analyze --input-file missing.md`
**Then** exit code is 1
**And** error message indicates file not found

**Related BR:** BR-IO-001

---

## US-002: Conduct Structured Interview

### AC-INT-001: Initialize Interview

**Given** an analysis JSON file
**When** I run `loom-cli interview --init analysis.json --state state.json`
**Then** a new interview state file is created
**And** the first question is output as JSON
**And** exit code is 100

**Related BR:** BR-INT-001, BR-INT-003

---

### AC-INT-002: View Current Question

**Given** an interview state file with pending questions
**When** I run `loom-cli interview --state state.json`
**Then** the current question is output as JSON
**And** exit code is 100

**Related BR:** BR-INT-001

---

### AC-INT-003: Answer Question

**Given** an interview state file with pending questions
**When** I run `loom-cli interview --state state.json --answer '{"question_id":"Q1", "answer":"Yes", "source":"user"}'`
**Then** the answer is recorded in the state file
**And** source is recorded as "user"
**And** the next question is output (or exit code 0 if complete)

**Related BR:** BR-INT-002

---

### AC-INT-004: Complete Interview

**Given** an interview state file with one remaining question
**When** I answer the last question
**Then** exit code is 0
**And** no question JSON is output

**Related BR:** BR-INT-003

---

### AC-INT-005: Grouped Question Mode

**Given** an analysis JSON file
**When** I run `loom-cli interview --init analysis.json --state state.json --grouped`
**Then** all questions are output at once
**And** exit code is 100

**Related BR:** BR-INT-001 (override)

---

### AC-INT-006: Batch Answers

**Given** an interview state file with multiple pending questions
**When** I run `loom-cli interview --state state.json --answers '[{"question_id":"Q1","answer":"A1","source":"user"},{"question_id":"Q2","answer":"A2","source":"user"}]'`
**Then** all provided answers are recorded
**And** remaining questions (if any) are indicated by exit code

---

## US-003: Derive Strategic Design (L1)

### AC-DRV-001: Derive L1 from Analysis

**Given** an analysis JSON file (or completed interview state)
**When** I run `loom-cli derive --output-dir ./l1 --analysis-file analysis.json`
**Then** the following files are created in ./l1:
- domain-model.md
- bounded-context-map.md
- acceptance-criteria.md
- business-rules.md
- decisions.md
**And** all files have YAML frontmatter with level: L1

**Related BR:** BR-DOC-001, BR-DRV-001

---

### AC-DRV-002: Derive with Vocabulary

**Given** an analysis file and a domain vocabulary file
**When** I run `loom-cli derive --output-dir ./l1 --analysis-file analysis.json --vocabulary vocab.md`
**Then** domain model uses vocabulary terms
**And** entity names match vocabulary definitions

---

### AC-DRV-003: Derive with NFR

**Given** an analysis file and a non-functional requirements file
**When** I run `loom-cli derive --output-dir ./l1 --analysis-file analysis.json --nfr nfr.md`
**Then** business rules include NFR-derived rules
**And** NFRs are traceable in the output

---

### AC-DRV-004: Derive Creates Output Directory

**Given** an analysis file
**And** output directory does not exist
**When** I run `loom-cli derive --output-dir ./new-l1 --analysis-file analysis.json`
**Then** the output directory is created
**And** L1 files are written successfully

**Related BR:** BR-IO-003

---

## US-004: Derive Tactical Design (L2)

### AC-L2-001: Derive L2 from L1

**Given** a directory containing L1 documents (acceptance-criteria.md, business-rules.md, domain-model.md)
**When** I run `loom-cli derive-l2 --input-dir ./l1 --output-dir ./l2`
**Then** the following files are created in ./l2:
- tech-specs.md
- interface-contracts.md
- aggregate-design.md
- sequence-design.md
- initial-data-model.md
**And** all files have YAML frontmatter with level: L2

**Related BR:** BR-DOC-001, BR-DRV-001

---

### AC-L2-002: Parallel Phase Execution

**Given** L1 documents
**When** I run `loom-cli derive-l2 --input-dir ./l1 --output-dir ./l2`
**Then** independent phases execute in parallel
**And** maximum 3 phases run concurrently
**And** total execution time is less than sequential execution

**Related BR:** BR-DRV-002

---

### AC-L2-003: Resume from Checkpoint

**Given** a previously interrupted L2 derivation with checkpoint
**When** I run `loom-cli derive-l2 --input-dir ./l1 --output-dir ./l2 --resume`
**Then** completed phases are skipped
**And** derivation continues from last incomplete phase

**Related BR:** BR-DRV-003

---

### AC-L2-004: Interactive Approval

**Given** L1 documents
**When** I run `loom-cli derive-l2 --input-dir ./l1 --output-dir ./l2 --interactive`
**Then** each generated file is previewed before writing
**And** I can approve, skip, or abort
**And** skipped files are not written

**Related BR:** BR-DRV-004

---

## US-005: Derive Operational Design (L3)

### AC-L3-001: Derive L3 from L2

**Given** a directory containing L2 documents
**When** I run `loom-cli derive-l3 --input-dir ./l2 --output-dir ./l3`
**Then** the following files are created in ./l3:
- test-cases.md
- openapi.json
- implementation-skeletons.md
- feature-tickets.md
- service-boundaries.md
- event-message-design.md
- dependency-graph.md

**Related BR:** BR-DRV-001

---

### AC-L3-002: TDAI Test Generation

**Given** L2 documents with acceptance criteria references
**When** I derive L3 test cases
**Then** each AC has at least one positive test
**And** each AC has at least one hallucination prevention test
**And** negative test ratio is at least 20%

**Related BR:** BR-TDAI-001, BR-TDAI-002

---

### AC-L3-003: OpenAPI JSON Format

**Given** L2 interface contracts
**When** I derive L3 openapi.json
**Then** the output is valid OpenAPI 3.0 JSON
**And** all operations from interface-contracts.md are included

---

## US-006: Validate Generated Documents

### AC-VAL-001: Validate All Levels

**Given** a directory with L1, L2, and L3 documents
**When** I run `loom-cli validate --input-dir ./specs`
**Then** all 10 validation rules are checked
**And** results show pass/fail for each rule
**And** violations list specific document IDs and issues

**Related BR:** BR-VAL-001

---

### AC-VAL-002: Validate Specific Level

**Given** a directory with documents
**When** I run `loom-cli validate --input-dir ./specs --level L2`
**Then** only L2-applicable rules are checked (V001-V004, V006, V007, V010)
**And** L3-only rules (V005, V008, V009) are skipped

**Related BR:** BR-VAL-002

---

### AC-VAL-003: JSON Output

**Given** a directory with documents containing violations
**When** I run `loom-cli validate --input-dir ./specs --json`
**Then** output is valid JSON
**And** structure includes `passed`, `rules`, and `violations` fields

**Related BR:** BR-VAL-003

---

### AC-VAL-004: Bidirectional Link Check

**Given** AC-ORD-001 references BR-ORD-001
**And** BR-ORD-001 does not reference back to AC-ORD-001
**When** I run validation
**Then** V004 fails
**And** violation identifies the missing back-reference

**Related BR:** BR-DOC-003

---

## US-007: Sync Bidirectional Links

### AC-SYN-001: Fix Missing Back-References

**Given** AC-ORD-001 references BR-ORD-001
**And** BR-ORD-001 does not reference back to AC-ORD-001
**When** I run `loom-cli sync-links --input-dir ./specs`
**Then** BR-ORD-001 is updated to include reference to AC-ORD-001
**And** V004 validation now passes

**Related BR:** BR-DOC-003

---

### AC-SYN-002: Dry Run Mode

**Given** documents with missing back-references
**When** I run `loom-cli sync-links --input-dir ./specs --dry-run`
**Then** no files are modified
**And** output shows what changes would be made

---

## US-008: Run Full Cascade Derivation

### AC-CAS-001: Full Cascade with Skip Interview

**Given** an L0 user story file
**When** I run `loom-cli cascade --input-file story.md --output-dir ./specs --skip-interview`
**Then** directories l1/, l2/, l3/ are created under ./specs
**And** all L1, L2, L3 documents are generated
**And** interview phase is skipped using AI defaults

**Related BR:** BR-INT-004, BR-DRV-001

---

### AC-CAS-002: Cascade with Interview

**Given** an L0 user story file
**When** I run `loom-cli cascade --input-file story.md --output-dir ./specs`
**Then** analyze phase runs first
**Then** interview phase presents questions
**And** after interview completes, derivation continues to L3

---

### AC-CAS-003: Cascade with Existing Decisions

**Given** an L0 user story file and existing decisions.md
**When** I run `loom-cli cascade --input-file story.md --output-dir ./specs --decisions decisions.md`
**Then** existing decisions are used
**And** interview is skipped or shortened based on coverage

---

### AC-CAS-004: Cascade Resume

**Given** an interrupted cascade with .cascade-state.json
**When** I run `loom-cli cascade --input-file story.md --output-dir ./specs --resume`
**Then** completed phases are skipped
**And** derivation resumes from last incomplete phase

**Related BR:** BR-DRV-003

---

### AC-CAS-005: Cascade Re-derive from Level

**Given** a completed cascade
**When** I run `loom-cli cascade --input-file story.md --output-dir ./specs --from l2`
**Then** L1 is kept unchanged
**And** L2 and L3 are re-derived

---

### AC-CAS-006: Cascade Output Structure

**Given** a cascade derivation
**When** derivation completes
**Then** output directory structure is:
```
specs/
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
│   └── [L3 documents]
└── .cascade-state.json
```

---

## US-009: Get Help and Version Information

### AC-HLP-001: Help Command

**Given** loom-cli is installed
**When** I run `loom-cli help` or `loom-cli --help` or `loom-cli -h`
**Then** usage information is displayed
**And** all commands are listed with descriptions
**And** examples are provided

---

### AC-HLP-002: Version Command

**Given** loom-cli is installed
**When** I run `loom-cli version`
**Then** version string is displayed (format: `loom-cli vX.Y.Z`)

---

### AC-HLP-003: No Arguments

**Given** loom-cli is installed
**When** I run `loom-cli` with no arguments
**Then** usage information is displayed (same as help)

---

## Related Documents

| Level | Document | Description |
|-------|----------|-------------|
| L0 | [l0-domain-vocabulary.md](l0-domain-vocabulary.md) | Domain Vocabulary |
| L0 | [l0-loom-cli.md](l0-loom-cli.md) | User Stories (source for this document) |
| L0 | [l0-nfr.md](l0-nfr.md) | Non-Functional Requirements |
| L1 | [l1-domain-model.md](l1-domain-model.md) | Domain Model |
| L1 | [l1-bounded-context-map.md](l1-bounded-context-map.md) | Bounded Context Map |
| L1 | [l1-business-rules.md](l1-business-rules.md) | Business Rules |
| L1 | This document | Acceptance Criteria |
| L1 | [l1-decisions.md](l1-decisions.md) | Design Decisions |
| L2 | [l2-cli-interface.md](l2-cli-interface.md) | CLI Interface Contract |
