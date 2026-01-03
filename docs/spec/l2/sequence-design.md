---
title: "Loom CLI Sequence Design"
generated: 2025-01-03T14:30:00Z
status: draft
level: L2
---

# Loom CLI Sequence Design

## Overview

This document defines the sequence designs for loom-cli operations, showing the flow of control between components.

**Traceability:** Derived from [acceptance-criteria.md](../l1/acceptance-criteria.md) and [business-rules.md](../l1/business-rules.md).

---

## SEQ-CAS-001: Full Cascade Derivation

**Traces to:** AC-CAS-001, US-008

**Actors:**
- User (CLI caller)
- Cascade Command
- Phase Executors
- Claude API
- File System

```mermaid
sequenceDiagram
    participant U as User
    participant C as Cascade
    participant A as Analyze
    participant I as Interview
    participant D1 as Derive L1
    participant D2 as Derive L2
    participant D3 as Derive L3
    participant API as Claude API
    participant FS as File System

    U->>C: loom-cli cascade --input-file story.md --output-dir ./specs
    C->>FS: Create l1/, l2/, l3/ directories
    C->>C: Load or create CascadeState

    rect rgb(200, 220, 240)
        Note over C,API: Phase 1: Analyze
        C->>A: runCascadeAnalyze()
        A->>FS: Read input file
        A->>API: Domain discovery prompt
        API-->>A: Entities, operations, relationships
        A->>API: Completeness analysis prompt
        API-->>A: Ambiguities/questions
        A->>FS: Write .analysis.json
        A-->>C: Phase complete
    end

    rect rgb(220, 240, 200)
        Note over C,API: Phase 2: Interview (if not skipped)
        C->>I: runCascadeInterview()
        I->>FS: Read .analysis.json
        I->>FS: Write .interview-state.json
        I-->>C: Phase complete
    end

    rect rgb(240, 220, 200)
        Note over C,API: Phase 3: Derive L1
        C->>D1: runCascadeDeriveL1()
        D1->>FS: Read .interview-state.json
        D1->>API: Domain model prompt
        API-->>D1: domain-model.md content
        D1->>API: Business rules prompt
        API-->>D1: business-rules.md content
        D1->>API: Acceptance criteria prompt
        API-->>D1: acceptance-criteria.md content
        D1->>FS: Write l1/*.md files
        D1-->>C: Phase complete
    end

    rect rgb(240, 200, 220)
        Note over C,API: Phase 4: Derive L2
        C->>D2: runCascadeDeriveL2()
        D2->>FS: Read l1/*.md files
        D2->>API: Tech specs prompt (parallel)
        D2->>API: Interface contracts prompt (parallel)
        D2->>API: Aggregate design prompt (parallel)
        API-->>D2: L2 content (collected)
        D2->>FS: Write l2/*.md files
        D2-->>C: Phase complete
    end

    rect rgb(220, 200, 240)
        Note over C,API: Phase 5: Derive L3
        C->>D3: runCascadeDeriveL3()
        D3->>FS: Read l2/*.md files
        D3->>API: Test cases prompt (parallel)
        D3->>API: OpenAPI prompt (parallel)
        D3->>API: Implementation skeletons prompt (parallel)
        API-->>D3: L3 content (collected)
        D3->>FS: Write l3/*.md files
        D3-->>C: Phase complete
    end

    C->>FS: Update .cascade-state.json (completed)
    C-->>U: Exit 0, print summary
```

---

## SEQ-CAS-002: Cascade Resume

**Traces to:** AC-CAS-004, BR-DRV-003

```mermaid
sequenceDiagram
    participant U as User
    participant C as Cascade
    participant FS as File System

    U->>C: loom-cli cascade --output-dir ./specs --resume
    C->>FS: Load .cascade-state.json

    alt State file exists
        C->>C: Check phase statuses
        Note over C: Skip completed phases

        loop For each pending/failed phase
            C->>C: Execute phase
            C->>FS: Update .cascade-state.json
        end

        C-->>U: Exit 0, resume successful
    else State file not found
        C-->>U: Exit 1, "No state file found"
    end
```

---

## SEQ-INT-001: Interview Flow

**Traces to:** AC-INT-001 through AC-INT-004, BR-INT-001

```mermaid
sequenceDiagram
    participant U as User
    participant I as Interview
    participant FS as File System

    U->>I: loom-cli interview --init analysis.json --state state.json
    I->>FS: Read analysis.json
    I->>I: Extract questions from ambiguities
    I->>FS: Create state.json with questions
    I-->>U: Output first question JSON, Exit 100

    loop Until all questions answered
        U->>I: loom-cli interview --state state.json --answer '{...}'
        I->>FS: Read state.json
        I->>I: Record decision
        I->>I: Advance to next question
        I->>FS: Update state.json

        alt More questions
            I-->>U: Output next question JSON, Exit 100
        else All answered
            I-->>U: Exit 0
        end
    end
```

---

## SEQ-INT-002: Grouped Interview Mode

**Traces to:** AC-INT-005, BR-INT-001 (override)

```mermaid
sequenceDiagram
    participant U as User
    participant I as Interview
    participant FS as File System

    U->>I: loom-cli interview --init analysis.json --state state.json --grouped
    I->>FS: Read analysis.json
    I->>I: Extract all questions
    I->>FS: Create state.json
    I-->>U: Output ALL questions JSON, Exit 100

    U->>I: loom-cli interview --state state.json --answers '[{...}, {...}]'
    I->>FS: Read state.json
    I->>I: Record all decisions
    I->>FS: Update state.json
    I-->>U: Exit 0
```

---

## SEQ-VAL-001: Document Validation

**Traces to:** AC-VAL-001, BR-VAL-001

```mermaid
sequenceDiagram
    participant U as User
    participant V as Validate
    participant FS as File System

    U->>V: loom-cli validate --input-dir ./specs --level ALL

    rect rgb(200, 220, 240)
        Note over V,FS: Phase 1: Collect IDs
        V->>FS: Glob *.md files
        loop For each file
            V->>FS: Read file
            V->>V: Extract IDs with regex
            V->>V: Build ID registry
        end
    end

    rect rgb(220, 240, 200)
        Note over V: Phase 2: Structural Validation
        V->>V: V001: Check documents have IDs
        V->>V: V002: Validate ID patterns
        V->>V: V010: Check for duplicates
    end

    rect rgb(240, 220, 200)
        Note over V: Phase 3: Traceability Validation
        V->>V: V003: Validate references exist
        V->>V: V004: Check bidirectional links
    end

    rect rgb(240, 200, 220)
        Note over V: Phase 4: Completeness Validation
        V->>V: V005: AC has test cases
        V->>V: V006: Entity has aggregate
        V->>V: V007: Service has contract
    end

    rect rgb(220, 200, 240)
        Note over V: Phase 5: TDAI Validation
        V->>V: V008: Negative test ratio >= 20%
        V->>V: V009: Hallucination prevention tests
    end

    V->>V: Calculate summary

    alt --json flag
        V-->>U: Output JSON result
    else text output
        V-->>U: Output formatted text
    end

    alt Validation passed
        V-->>U: Exit 0
    else Validation failed
        V-->>U: Exit 1
    end
```

---

## SEQ-DRV-001: L2 Parallel Derivation

**Traces to:** AC-L2-002, BR-DRV-002

```mermaid
sequenceDiagram
    participant D as Derive L2
    participant PE as ParallelExecutor
    participant API as Claude API
    participant FS as File System

    D->>FS: Read L1 documents
    D->>PE: Create executor (max 3 concurrent)

    par Tech Specs (slot 1)
        PE->>API: Tech specs prompt
        API-->>PE: tech-specs.md content
    and Interface Contracts (slot 2)
        PE->>API: Interface contracts prompt
        API-->>PE: interface-contracts.md content
    and Aggregate Design (slot 3)
        PE->>API: Aggregate design prompt
        API-->>PE: aggregate-design.md content
    end

    Note over PE: Wait for slot, then continue

    par Sequence Design (slot 1)
        PE->>API: Sequence design prompt
        API-->>PE: sequence-design.md content
    and Data Model (slot 2)
        PE->>API: Data model prompt
        API-->>PE: initial-data-model.md content
    end

    PE-->>D: All results collected

    loop For each result
        D->>FS: Write l2/*.md file
    end

    D-->>D: Phase complete
```

---

## SEQ-INT-003: Interactive Approval Flow

**Traces to:** AC-L2-004, BR-DRV-004, DEC-031

```mermaid
sequenceDiagram
    participant U as User
    participant D as Derive
    participant API as Claude API
    participant FS as File System

    D->>API: Generate content
    API-->>D: Content

    loop For each generated file
        D->>U: Render preview (first 20 lines)
        D->>U: Prompt: [A]pprove / [E]dit / [R]egenerate / [S]kip / [Q]uit

        alt Approve (A)
            U->>D: 'a'
            D->>FS: Write file
        else Edit (E)
            U->>D: 'e'
            D->>FS: Write to temp file
            D->>U: Open $EDITOR
            U->>D: Save and close
            D->>FS: Read edited content
            D->>FS: Write final file
        else Regenerate (R)
            U->>D: 'r'
            D->>API: Regenerate prompt
            API-->>D: New content
            Note over D: Loop back to preview
        else Skip (S)
            U->>D: 's'
            Note over D: Continue to next file
        else Quit (Q)
            U->>D: 'q'
            D->>U: Confirm quit?
            alt Confirmed
                D-->>U: Exit 0
            else Not confirmed
                Note over D: Loop back to prompt
            end
        end
    end
```

---

## SEQ-SYN-001: Sync Links

**Traces to:** AC-SYN-001, BR-DOC-003

```mermaid
sequenceDiagram
    participant U as User
    participant S as SyncLinks
    participant FS as File System

    U->>S: loom-cli sync-links --input-dir ./specs

    S->>FS: Read all *.md files
    S->>S: Parse document IDs
    S->>S: Parse references

    loop For each reference A → B
        S->>S: Check if B → A exists
        alt Missing back-reference
            S->>S: Add B → A to fix list
        end
    end

    alt --dry-run
        S-->>U: Print proposed changes
    else Apply changes
        loop For each fix
            S->>FS: Read file
            S->>S: Add reference to Traceability section
            S->>FS: Write updated file
        end
        S-->>U: Print applied changes
    end

    S-->>U: Exit 0
```

---

## Related Documents

| Level | Document | Description |
|-------|----------|-------------|
| L0 | [domain-vocabulary.md](../l0/domain-vocabulary.md) | Domain Vocabulary |
| L0 | [loom-cli.md](../l0/loom-cli.md) | User Stories |
| L0 | [nfr.md](../l0/nfr.md) | Non-Functional Requirements |
| L0 | [decisions.md](../l0/decisions.md) | Design Decisions (L0→L1) |
| L1 | [domain-model.md](../l1/domain-model.md) | Domain Model |
| L1 | [decisions.md](../l1/decisions.md) | Design Decisions (L1→L2) |
| L2 | [decisions.md](decisions.md) | Design Decisions (L2→L3) |
| L2 | [prompt-catalog.md](prompt-catalog.md) | Prompt Catalog |
| L1 | [bounded-context-map.md](../l1/bounded-context-map.md) | Bounded Context Map |
| L1 | [business-rules.md](../l1/business-rules.md) | Business Rules |
| L1 | [acceptance-criteria.md](../l1/acceptance-criteria.md) | Acceptance Criteria (source) |
| L2 | [tech-specs.md](tech-specs.md) | Technical Specifications |
| L2 | [interface-contracts.md](interface-contracts.md) | CLI Interface Contract |
| L2 | [aggregate-design.md](aggregate-design.md) | Aggregate Design |
| L2 | This document | Sequence Design |
| L2 | [initial-data-model.md](initial-data-model.md) | Data Model |
