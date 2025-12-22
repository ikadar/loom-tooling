---
name: loom-derive
description: Orchestrate L0 → L1 derivation using completeness-driven discovery
version: "4.1.0"
arguments:
  - name: input-file
    description: "Path to single L0 input file. Use this OR --input-dir, not both."
    required: false
  - name: input-dir
    description: "Path to directory containing L0 input files (reads all *.md files). Use this OR --input-file, not both."
    required: false
  - name: output-dir
    description: "Directory for generated L1 documents"
    required: true
  - name: decisions-file
    description: "Path to decisions.md file (default: {input-dir}/decisions.md or {input-file-dir}/decisions.md)"
    required: false
---

# Loom L0 → L1 Derivation (Orchestrator)

Orchestrates the complete L0 → L1 derivation process.

## Architecture

```
/loom-derive (this command)
    │
    ├─→ Phase 0: Read & Parse Input
    │       └─→ Read L0 files + decisions.md (if exists)
    │
    ├─→ Phase 1: Domain Discovery
    │       └─→ Extract entities, operations, relationships, UI mentions
    │
    ├─→ Phase 2: Completeness Analysis (parallel)
    │       ├─→ /loom-analyze-entities → entity ambiguities
    │       ├─→ /loom-analyze-operations → operation ambiguities
    │       └─→ /loom-analyze-ui → UI ambiguities (or request UI input)
    │
    ├─→ Phase 3: Merge & Filter Ambiguities
    │       └─→ Remove already-resolved (from decisions.md)
    │
    ├─→ Phase 4: Structured Interview
    │       └─→ /loom-interview → new resolutions only
    │
    ├─→ Phase 5: Derivation
    │       └─→ Generate AC + BR using ALL resolutions (existing + new)
    │
    └─→ Phase 6: Write Output
            ├─→ acceptance-criteria.md, business-rules.md
            └─→ APPEND new resolutions to decisions.md
```

## Reference Documents

- `.claude/docs/checklists/entity-checklist.md`
- `.claude/docs/checklists/operation-checklist.md`
- `.claude/docs/checklists/ui-checklist.md`

---

## Phase 0: Read & Parse Input

### Single File Mode

```bash
/loom-derive --input-file path/to/input.md --output-dir path/to/output
```

Read the specified file.

### Directory Mode

```bash
/loom-derive --input-dir path/to/input/ --output-dir path/to/output
```

1. Glob for all `*.md` files in directory
2. Read all files
3. Concatenate with source markers:

```markdown
<!-- SOURCE: file1.md -->
{content}

<!-- SOURCE: file2.md -->
{content}
```

### Validation

- If both `--input-file` and `--input-dir` → Error
- If neither → Error
- If no `.md` files found → Error

### Load Existing Decisions

Check for `decisions.md` in the input directory:

```bash
# Default location
{input-dir}/decisions.md
# Or explicitly specified
--decisions-file path/to/decisions.md
```

If `decisions.md` exists, parse it:

```yaml
existing_resolutions:
  - id: "AMB-ENT-001"
    question: "What happens to tasks when station deleted?"
    answer: "Block deletion if tasks exist"
    decided_at: "2025-12-21T10:30:00Z"
    source: "user"

  - id: "AMB-OP-005"
    question: "Time snap granularity?"
    answer: "15 minutes"
    decided_at: "2025-12-21T10:32:00Z"
    source: "user"
```

**Output:**
```markdown
## Existing Decisions Loaded

Found `decisions.md` with **23** previous resolutions.
These will be used and not asked again.
```

If no `decisions.md` exists, continue with empty resolutions list.

---

## Phase 1: Domain Discovery

Extract from input:

```yaml
domain:
  entities:
    - name: "Station"
      mentioned_attributes: [name, category, capacity, operating_hours]
      mentioned_operations: [create, update]
      mentioned_states: []

    - name: "Job"
      mentioned_attributes: [client, deadline, paper_status, bat_status]
      mentioned_operations: [create, delete, schedule]
      mentioned_states: [late, on_time]

  operations:
    - name: "Schedule Task"
      actor: "Scheduler"
      trigger: "drag-and-drop"
      target: "Task"
      mentioned_inputs: [task_id, station_id, start_time]
      mentioned_rules: [no_overlap, precedence, snap_grid]

  relationships:
    - from: "Job"
      to: "Task"
      type: "contains"
      cardinality: "1:N"

  business_rules:
    - "No overlapping tasks on capacity-1 stations"
    - "Task sequence must be respected"

  ui_mentions:
    - "Scheduling Grid"
    - "Left Panel"
    - "Right Panel"
    - "drag-and-drop"
```

**Output:** Present discovery summary to user.

---

## Phase 2: Completeness Analysis

Run analysis commands (conceptually in parallel):

### 2.1 Entity Analysis

Apply `/loom-analyze-entities` logic:
- Use `entity-checklist.md` reference
- Check every entity against full checklist
- Generate entity ambiguities

### 2.2 Operation Analysis

Apply `/loom-analyze-operations` logic:
- Use `operation-checklist.md` reference
- Check every operation against full checklist
- Generate operation ambiguities

### 2.3 UI Analysis

Apply `/loom-analyze-ui` logic:
- First check if UI spec exists
- If UI mentioned but no spec → **STOP and request UI input**
- If UI spec exists → analyze with `ui-checklist.md`
- Generate UI ambiguities

**If UI input missing:**

```markdown
## ⚠️ UI/UX Specification Required

UI components are mentioned but not specified:
- Scheduling Grid
- Left Panel (job list)
- Right Panel (late jobs)
- Drag-and-drop interactions

**Options:**
1. Provide UI/UX specification file
2. Answer UI questions interactively (many questions)
3. Mark UI as out of scope (will skip UI analysis)

Which option?
```

---

## Phase 3: Merge & Filter Ambiguities

### 3.1 Combine All Ambiguities

```yaml
all_ambiguities:
  total: 87

  by_severity:
    critical: 23
    important: 41
    minor: 23

  by_source:
    entities: 36
    operations: 29
    ui: 22
```

### 3.2 Filter Already-Resolved

Compare against `existing_resolutions` from decisions.md:

```yaml
# Match by question similarity, not just ID
# (IDs may change between runs, questions are stable)

filtered_ambiguities:
  total: 64  # 87 - 23 already resolved

  already_resolved: 23
  new_to_ask: 64

  by_severity:
    critical: 18  # was 23, 5 already resolved
    important: 31 # was 41, 10 already resolved
    minor: 15     # was 23, 8 already resolved
```

**Output:**
```markdown
## Ambiguity Summary

| Category | Found | Already Resolved | To Ask |
|----------|-------|------------------|--------|
| Entities | 36 | 12 | 24 |
| Operations | 29 | 8 | 21 |
| UI | 22 | 3 | 19 |
| **Total** | **87** | **23** | **64** |

Using 23 decisions from `decisions.md`.
```

---

## Phase 4: Structured Interview

Apply `/loom-interview` logic:

1. Prioritize by severity (critical first)
2. Group by source/entity for context
3. Batch questions (4-6 per round)
4. Record all answers
5. Handle follow-up questions
6. Bulk confirm minor defaults

**Loop until:**
- Zero critical ambiguities remaining
- Zero important ambiguities remaining
- All minor have resolution or confirmed default

**Output:** Full interview record.

---

## Phase 5: Derivation

With all ambiguities resolved, generate:

### 5.1 Acceptance Criteria

For each user story/operation:

```markdown
### AC-{DOMAIN}-{NUM} – {Title}

**Given** [precondition]
**When** [action]
**Then** [outcome]

**Resolved Ambiguities:**
- AMB-XXX: {answer}

**Error Cases:**
- {condition} → {behavior} (from AMB-YYY)

**Traceability:**
- Input: {source_file} § "{section}"
- Interview: Round {N}, AMB-XXX
```

### 5.2 Business Rules

For each constraint/rule:

```markdown
### BR-{DOMAIN}-{NUM} – {Title}

**Rule:** [Statement]

**Invariant:** [Formal condition using MUST/MUST NOT]

**Enforcement:**
- Check point: [where]
- Violation: [behavior]
- Error: `{ERROR_CODE}`

**Source:**
- Input: {source}
- Interview: AMB-XXX
```

---

## Phase 6: Write Output

### Files Generated

1. `{output-dir}/acceptance-criteria.md`
2. `{output-dir}/business-rules.md`
3. `{input-dir}/decisions.md` - **APPEND** new resolutions (persistent)
4. `{output-dir}/interview-record.md` - Full session log (optional, for audit)

### 6.1 Update decisions.md

**CRITICAL:** Append new resolutions to `decisions.md`, preserving existing ones.

```markdown
---
# decisions.md - Loom Decision Log
# This file persists interview answers across derivation runs.
# DO NOT delete - answers will be asked again!
---

## Entity Decisions

### Station

- **AMB-ENT-001: Deletion behavior**
  - Q: What happens to tasks when station deleted?
  - A: Block deletion if tasks exist
  - Decided: 2025-12-21 by user

- **AMB-ENT-002: Name uniqueness**
  - Q: Must station name be unique?
  - A: Yes, unique within organization
  - Decided: 2025-12-21 by user

### Job

- **AMB-ENT-010: Late threshold**
  - Q: When is a job considered "late"?
  - A: When any task misses its deadline by > 0 minutes
  - Decided: 2025-12-21 by user

## Operation Decisions

### Schedule Task

- **AMB-OP-001: Overlap behavior**
  - Q: What happens when task overlaps existing?
  - A: Block with error, show conflict details
  - Decided: 2025-12-21 by user

## UI Decisions

- **AMB-UI-001: Drag feedback**
  - Q: What visual feedback during task drag?
  - A: Ghost image + valid/invalid drop zone highlighting
  - Decided: 2025-12-21 by user

## Defaults Accepted

The following minor decisions use suggested defaults:

| ID | Question | Default | Accepted |
|----|----------|---------|----------|
| AMB-ENT-050 | Max station name length | 100 chars | 2025-12-21 |
| AMB-ENT-051 | Max job title length | 200 chars | 2025-12-21 |
| AMB-OP-050 | Operation timeout | 30 seconds | 2025-12-21 |
```

### 6.2 Frontmatter for L1 Documents

```yaml
---
id: L1-AC
status: draft
derived-from:
  - "{input-file-1}"
  - "{input-file-2}"
  - "decisions.md"
derived-at: "{ISO timestamp}"
derived-by: "loom-derive v4.1.0"
completeness-analysis:
  entities-analyzed: 5
  operations-analyzed: 8
  ambiguities-found: 87
  ambiguities-resolved: 87
decisions:
  from-existing: 23    # loaded from decisions.md
  from-this-session: 64  # asked in this run
  total: 87
---
```

---

## Quality Criteria

Before writing output, verify:

- [ ] All entities analyzed with full checklist
- [ ] All operations analyzed with full checklist
- [ ] UI analyzed OR explicitly marked out of scope
- [ ] Zero critical ambiguities remaining
- [ ] Zero important ambiguities remaining
- [ ] All minor have resolution or confirmed default
- [ ] All ACs trace to input and/or interview
- [ ] All BRs trace to input and/or interview
- [ ] Interview record complete

---

## Summary Output

```markdown
## Derivation Complete

### Input
- Files: {list}
- Total lines: {N}
- Existing decisions: {M} (from decisions.md)

### Analysis
| Category | Analyzed | Ambiguities | Already Resolved | Asked |
|----------|----------|-------------|------------------|-------|
| Entities | 5 | 36 | 12 | 24 |
| Operations | 8 | 29 | 8 | 21 |
| UI | 1 screen | 22 | 3 | 19 |
| **Total** | | **87** | **23** | **64** |

### Interview (this session)
| Rounds | From User | From Defaults |
|--------|-----------|---------------|
| 12 | 52 | 12 |

### Output
| Document | Items |
|----------|-------|
| Acceptance Criteria | 34 |
| Business Rules | 12 |

### Files Written
- {output-dir}/acceptance-criteria.md
- {output-dir}/business-rules.md
- {input-dir}/decisions.md (23 existing + 64 new = 87 total)

**Expansion:** {input_lines} → {output_lines} ({ratio}x)
```

---

## Now: Execute

Begin with Phase 0: Read input files.
