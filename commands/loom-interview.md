---
name: loom-interview
description: Run structured interview to resolve ambiguities
version: "1.1.0"
arguments:
  - name: ambiguities
    description: "Merged ambiguity list from analysis commands (YAML or inline)"
    required: true
  - name: existing-resolutions
    description: "Previously resolved ambiguities from decisions.md (to skip)"
    required: false
  - name: batch-size
    description: "Questions per round (default: 4)"
    required: false
---

# Structured Interview

Resolve ambiguities through interactive questioning.

## Input

### Ambiguity List

Merged ambiguity list in format:

```yaml
ambiguities:
  - id: "AMB-ENT-001"
    source: "entities"  # entities, operations, ui
    entity: "Station"
    category: "deletion"
    aspect: "cascade"
    question: "What happens to tasks when station deleted?"
    severity: "critical"
    options:  # optional
      - "Block deletion"
      - "Cascade delete tasks"
      - "Reassign tasks"
    suggested_default: null

  - id: "AMB-OP-001"
    source: "operations"
    operation: "Schedule Task"
    category: "conflict"
    aspect: "overlap_behavior"
    question: "What happens when task overlaps existing?"
    severity: "critical"
    options:
      - "Block"
      - "Push down"
      - "Swap"
```

### Existing Resolutions (optional)

If `--existing-resolutions` is provided, these questions will be skipped:

```yaml
existing_resolutions:
  - id: "AMB-ENT-001"
    question: "What happens to tasks when station deleted?"
    answer: "Block deletion if tasks exist"
    decided_at: "2025-12-21T10:30:00Z"
    source: "user"
```

## Process

### Step 0: Filter Already-Resolved

Before starting interview:

1. Compare each ambiguity question against existing resolutions
2. Match by **question similarity** (not ID, since IDs may change)
3. Remove matches from interview queue
4. Report what was skipped

```markdown
## Loading Previous Decisions

Found **23** existing resolutions in decisions.md.
Matched **23** ambiguities - these will not be asked again.

**Remaining to resolve:** 64 ambiguities
```

### Step 1: Prioritize & Group

1. Sort by severity: critical → important → minor
2. Group by source/entity/operation for context coherence
3. Calculate total rounds needed

```markdown
## Interview Plan

**Total Ambiguities:** 87
**Batch Size:** 4

| Priority | Count | Rounds |
|----------|-------|--------|
| Critical | 23 | 6 |
| Important | 41 | 11 |
| Minor | 23 | (bulk confirm) |

**Estimated Rounds:** 17 + 1 bulk confirmation
```

### Step 2: Run Interview Rounds

For each batch, use AskUserQuestion tool:

```markdown
## Structured Interview: Round {N} of {total}

**Progress:** {resolved} / {total} ambiguities resolved
**This round:** {batch_size} questions ({source} - {group})

---

### Q1: {question} [AMB-XXX-001]

{context if helpful}

Options:
a) {option1}
b) {option2}
c) {option3}
d) Other: ___

---

### Q2: {question} [AMB-XXX-002]

...
```

### Step 3: Record Answers

After each round:

```yaml
resolutions:
  - id: "AMB-ENT-001"
    question: "What happens to tasks when station deleted?"
    answer: "Block deletion if tasks exist"
    answer_code: "block"
    source: "user"
    round: 1

  - id: "AMB-OP-001"
    question: "What happens when task overlaps existing?"
    answer: "Insert and push existing tasks down"
    answer_code: "push_down"
    source: "user"
    round: 1
```

### Step 4: Check for Follow-up Questions

Some answers may trigger new questions:

```markdown
## Follow-up Questions Detected

Your answer to AMB-ENT-001 ("Block deletion if tasks exist") raises:

**AMB-ENT-001-F1:** Can an admin override the block and force delete?
a) No, never allow force delete
b) Yes, admin can force delete
c) Yes, with confirmation dialog

**AMB-ENT-001-F2:** Should the UI show which tasks are blocking deletion?
a) Yes, list blocking tasks
b) No, just show "has scheduled tasks"
```

### Step 5: Minor Ambiguity Bulk Confirmation

For minor ambiguities with suggested defaults:

```markdown
## Minor Ambiguities - Confirm Defaults

I'll use these defaults unless you object. Reply "OK" to accept all, or list IDs to change.

| ID | Question | Suggested Default |
|----|----------|-------------------|
| AMB-ENT-050 | Max station name length | 100 characters |
| AMB-ENT-051 | Max job title length | 200 characters |
| AMB-OP-050 | Operation timeout | 30 seconds |
| AMB-OP-051 | Audit log retention | 90 days |
| AMB-UI-050 | Animation duration | 200ms |
| ... | ... | ... |

**Total:** 23 minor ambiguities with defaults
```

### Step 6: Completion Summary

```markdown
## Interview Complete

### Ambiguity Resolution

| Status | Count |
|--------|-------|
| Already resolved (from decisions.md) | 23 |
| Resolved this session | 64 |
| **Total resolved** | **87** |

### This Session

| Source | From User | From Default | Total |
|--------|-----------|--------------|-------|
| Entities | 20 | 4 | 24 |
| Operations | 18 | 3 | 21 |
| UI/UX | 14 | 5 | 19 |
| **Total** | **52** | **12** | **64** |

**Rounds Conducted:** 12 + 1 bulk

---

### New Resolutions (to append to decisions.md)

{List of 64 new resolutions from this session}

### All Resolutions (for derivation)

{Combined list of 87 resolutions: 23 existing + 64 new}
```

## Output Format

```yaml
interview_record:
  conducted_at: "2025-12-21T10:30:00Z"

  ambiguities:
    total_found: 87
    already_resolved: 23  # from decisions.md
    asked_this_session: 64

  rounds_conducted: 12

  summary:
    from_user: 52
    from_defaults: 12
    follow_ups_generated: 3

  # Only NEW resolutions from this session (to append to decisions.md)
  new_resolutions:
    - id: "AMB-ENT-010"
      question: "When is a job considered late?"
      answer: "When any task misses deadline by > 0 minutes"
      source: "user"
      round: 1
      decided_at: "2025-12-21T11:45:00Z"

    - id: "AMB-OP-050"
      question: "Operation timeout"
      answer: "30 seconds"
      source: "default"
      round: "bulk"
      decided_at: "2025-12-21T12:00:00Z"

  # ALL resolutions (existing + new) for use in derivation
  all_resolutions:
    - id: "AMB-ENT-001"
      question: "What happens to tasks when station deleted?"
      answer: "Block deletion if tasks exist"
      source: "existing"  # from decisions.md

    - id: "AMB-ENT-010"
      question: "When is a job considered late?"
      answer: "When any task misses deadline by > 0 minutes"
      source: "user"

    # ... all 87 resolutions
```

## Question Design Principles

1. **Provide context** - Brief explanation of why this matters
2. **Offer concrete options** - Not open-ended "how should this work?"
3. **Include "Other"** - Always allow custom answers
4. **Group logically** - Related questions together
5. **Start critical** - Most important first
6. **Explain implications** - What each option means for implementation

## Error Handling

| Situation | Action |
|-----------|--------|
| User skips question | Mark as unresolved, ask again later |
| User answers "depends" | Ask for conditions, create branching rules |
| User is unsure | Suggest most common pattern, mark as tentative |
| Contradictory answers | Flag and ask for clarification |

## Quality Criteria

Before completing:
- [ ] All critical ambiguities resolved
- [ ] All important ambiguities resolved
- [ ] All minor ambiguities have resolution or default
- [ ] No contradictions in answers
- [ ] Follow-up questions addressed
- [ ] Full resolution record generated
