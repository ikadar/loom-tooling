---
name: loom-analyze-operations
description: Analyze operations for completeness and generate ambiguity list
version: "1.0.0"
arguments:
  - name: input
    description: "Parsed domain discovery output (operations) OR path to input file"
    required: true
  - name: output-format
    description: "Output format: 'ambiguities' (default) or 'full-report'"
    required: false
---

# Operation Completeness Analysis

Analyze operations for completeness using the operation checklist.

## Reference

Use checklist from: `.claude/docs/checklists/operation-checklist.md`

## Input

Either:
1. Parsed domain discovery (list of operations with known details)
2. Path to L0 input file (will extract operations first)

## Process

### Step 1: Operation Extraction (if needed)

If input is a file path, extract operations:

```yaml
operations:
  - name: "Schedule Task"
    actor: "Scheduler"
    trigger: "drag-and-drop"
    target: "Task"
    mentioned_inputs: [task_id, station_id, start_time]
    mentioned_outputs: [assignment]
    mentioned_errors: [conflict, precedence_violation]
    mentioned_rules: [no_overlap, snap_to_grid]
```

### Step 2: Apply Checklist

For EACH operation, go through EVERY section of the operation checklist:

1. **A. Basic Definition** - Actor, trigger, target, preconditions, postconditions, atomicity, idempotency
2. **B. Input Definition** - Every input parameter fully specified
3. **C. Output Definition** - Success response structure
4. **D. Error Handling** - Every possible error with status, code, message, recovery
5. **E. Side Effects** - Audit, notifications, entity changes, events
6. **F. Performance & Limits** - Timeout, rate limit, batch size
7. **G. Edge Cases** - Auto-generate boundary scenarios
8. **H. Transaction Boundaries** - Rollback, compensation

### Step 3: Generate Ambiguity List

For each `?` or missing item, create an ambiguity:

```yaml
ambiguities:
  - id: "AMB-OP-001"
    operation: "Schedule Task"
    category: "error_handling"
    aspect: "conflict.behavior"
    question: "When task overlaps existing task, should it block, push, or swap?"
    severity: "critical"
    options:
      - "Block drop with error"
      - "Insert and push existing down"
      - "Swap positions"
      - "Allow with conflict warning"

  - id: "AMB-OP-002"
    operation: "Schedule Task"
    category: "input"
    aspect: "start_time.granularity"
    question: "What is the time snap granularity for scheduling?"
    severity: "critical"
    options:
      - "1 minute"
      - "5 minutes"
      - "15 minutes"
      - "30 minutes"
```

### Step 4: Severity Classification

| Severity | Criteria |
|----------|----------|
| **critical** | Core operation behavior undefined |
| **important** | Error handling or edge case undefined |
| **minor** | Performance limits or minor details |

**Critical examples:**
- Main operation behavior
- Conflict resolution
- Required inputs
- Success criteria

**Important examples:**
- Specific error messages
- Side effect details
- Recovery actions
- Concurrency handling

**Minor examples:**
- Timeout values
- Rate limits
- Logging details

## Output Format

### Ambiguities Only (default)

```markdown
## Operation Analysis: Ambiguities Found

**Operations Analyzed:** 8
**Total Ambiguities:** 34

### Critical (12)

| ID | Operation | Aspect | Question |
|----|-----------|--------|----------|
| AMB-OP-001 | Schedule Task | conflict | What happens on overlap? |
| ... |

### Important (15)

| ID | Operation | Aspect | Question |
|----|-----------|--------|----------|
| ... |

### Minor (7)

| ID | Operation | Aspect | Question | Suggested Default |
|----|-----------|--------|----------|-------------------|
| ... |
```

### Full Report

Includes the ambiguity list PLUS full checklist results showing what IS defined.

## Auto-Generated Edge Cases

For each operation, automatically generate questions for:

| Edge Case | Question Template |
|-----------|-------------------|
| All optional omitted | "What if {operation} called with only required inputs?" |
| Min boundary | "What if {input} is at minimum value?" |
| Max boundary | "What if {input} is at maximum value?" |
| Below min | "What if {input} is below minimum?" |
| Above max | "What if {input} is above maximum?" |
| Empty string | "What if {string_input} is empty?" |
| Null value | "What if {input} is null?" |
| Target not found | "What if target {entity} doesn't exist?" |
| Target just deleted | "What if {entity} deleted between check and {operation}?" |
| Rapid repeat | "What if {operation} called twice rapidly?" |
| During loading | "What if {operation} triggered while previous still loading?" |
| Network failure | "What if network fails during {operation}?" |
| Concurrent same target | "What if two users {operation} same {entity} simultaneously?" |

## Quality Criteria

Before outputting, verify:
- [ ] Every operation has been analyzed
- [ ] Every input fully specified (type, required, validation)
- [ ] Every possible error has handling defined
- [ ] Every side effect identified
- [ ] Edge cases generated for each operation
- [ ] Severities correctly assigned
- [ ] No duplicate ambiguity IDs
