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
    participant D4 as Derive L4
    participant API as Claude API
    participant FS as File System

    U->>C: loom-cli cascade --input-file story.md --output-dir ./specs
    C->>FS: Create l1/, l2/, l3/, l4/ directories
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

    rect rgb(200, 240, 220)
        Note over C,API: Phase 6: Derive L4
        C->>D4: runCascadeDeriveL4()
        D4->>FS: Load loom.config.yaml
        D4->>FS: Read l2/*.md and l3/*.md files
        D4->>API: Architecture prompt
        D4->>API: Patterns prompt
        D4->>API: Coding standards prompt
        D4->>API: Project structure prompt
        D4->>API: Testing strategy prompt
        API-->>D4: L4 content (collected)
        D4->>FS: Write l4/*.md files
        D4-->>C: Phase complete
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

## SEQ-L4-001: L4 Derivation

**Traces to:** IC-DRV-004, US-010

**Actors:**
- Derive L4 Command
- Claude API
- File System
- Config Loader

```mermaid
sequenceDiagram
    participant U as User
    participant D as Derive L4
    participant CFG as Config
    participant API as Claude API
    participant FS as File System

    U->>D: loom-cli derive-l4 --input-dir ./specs --output-dir ./specs/l4
    D->>CFG: Load loom.config.yaml
    CFG-->>D: Config (language, methodology, etc.)
    D->>FS: Read L2/*.md and L3/*.md files

    rect rgb(200, 220, 240)
        Note over D,API: Phase 1: Architecture
        D->>D: Build context (aggregates, sequences, config)
        D->>API: derive-l4-architecture prompt
        API-->>D: architecture.md JSON
        D->>D: Format to markdown
        D->>FS: Write l4/architecture.md
    end

    rect rgb(220, 240, 200)
        Note over D,API: Phase 2: Patterns
        D->>D: Build context (tech-specs, skeletons, architecture)
        D->>API: derive-l4-patterns prompt
        API-->>D: patterns.md JSON
        D->>D: Format to markdown
        D->>FS: Write l4/patterns.md
    end

    rect rgb(240, 220, 200)
        Note over D,API: Phase 3: Coding Standards
        D->>D: Build context (config, language idioms)
        D->>API: derive-l4-coding-standards prompt
        API-->>D: coding-standards.md JSON
        D->>D: Format to markdown
        D->>FS: Write l4/coding-standards.md
    end

    rect rgb(240, 200, 220)
        Note over D,API: Phase 4: Project Structure
        D->>D: Build context (service-boundaries, architecture)
        D->>API: derive-l4-project-structure prompt
        API-->>D: project-structure.md JSON
        D->>D: Format to markdown
        D->>FS: Write l4/project-structure.md
    end

    rect rgb(220, 200, 240)
        Note over D,API: Phase 5: Testing Strategy
        D->>D: Build context (test-cases, config.methodology)
        D->>API: derive-l4-testing-strategy prompt
        API-->>D: testing-strategy.md JSON
        D->>D: Format to markdown
        D->>FS: Write l4/testing-strategy.md
    end

    D-->>U: Exit 0, print summary
```

---

## SEQ-GEN-001: Code Generation

**Traces to:** IC-GEN-001, US-011

**Actors:**
- Generate Command
- Claude API
- File System
- Validator
- Test Runner

```mermaid
sequenceDiagram
    participant U as User
    participant G as Generate
    participant V as Validator
    participant API as Claude API
    participant FS as File System
    participant T as TestRunner

    U->>G: loom-cli generate --input-dir ./specs --output-dir ./src

    rect rgb(200, 220, 240)
        Note over G,V: Phase 1: Pre-generation Validation
        G->>V: Validate L4 completeness
        V->>FS: Check all L4 files exist
        V->>V: Check source references
        V-->>G: Validation result
        alt Validation failed
            G-->>U: Exit 1, print errors
        end
    end

    rect rgb(220, 240, 200)
        Note over G,API: Phase 2: Code Generation
        G->>FS: Read L4/project-structure.md
        G->>FS: Create directory structure

        loop For each component
            G->>FS: Read L4/*.md context
            G->>API: Generate component code
            API-->>G: Code with traceability comments
            G->>FS: Write source file
        end
    end

    rect rgb(240, 220, 200)
        Note over G,API: Phase 3: Test Generation
        G->>FS: Read L3/test-cases.md
        G->>FS: Read L4/testing-strategy.md

        loop For each test case mapping
            G->>API: Generate test code
            API-->>G: Test code
            G->>FS: Write test file
        end
    end

    rect rgb(240, 200, 220)
        Note over G,T: Phase 4: Post-generation Validation
        G->>T: Run generated tests
        T-->>G: Test results

        alt Tests failed
            G-->>U: Exit 1, print failures
        else Tests passed
            G->>T: Check coverage
            T-->>G: Coverage report

            alt Coverage below threshold
                G-->>U: Exit 1, coverage warning
            end
        end
    end

    G-->>U: Exit 0, print summary
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
