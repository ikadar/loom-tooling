---
name: loom-analyze-entities
description: Analyze entities for completeness and generate ambiguity list
version: "1.0.0"
arguments:
  - name: input
    description: "Parsed domain discovery output (entities, relationships) OR path to input file"
    required: true
  - name: output-format
    description: "Output format: 'ambiguities' (default) or 'full-report'"
    required: false
---

# Entity Completeness Analysis

Analyze entities for completeness using the entity checklist.

## Reference

Use checklist from: `.claude/docs/checklists/entity-checklist.md`

## Input

Either:
1. Parsed domain discovery (list of entities with their known attributes)
2. Path to L0 input file (will extract entities first)

## Process

### Step 1: Entity Extraction (if needed)

If input is a file path, extract entities:

```yaml
entities:
  - name: "EntityName"
    mentioned_attributes:
      - name: "attr1"
        type: "inferred or ?"
        details: "what we know"
    mentioned_operations: [create, update, delete]
    mentioned_states: [state1, state2]
    mentioned_relationships:
      - target: "OtherEntity"
        type: "inferred"
```

### Step 2: Apply Checklist

For EACH entity, go through EVERY section of the entity checklist:

1. **A. Attribute Definition** - For each mentioned attribute + common expected attributes
2. **B. Enum Values** - For each status/enum field
3. **C. Lifecycle** - Creation, Update, Deletion
4. **D. Relationships** - For each relationship
5. **E. Constraints** - Business rules on this entity
6. **F. Concurrent Access** - Multi-user scenarios
7. **G. History & Audit** - Change tracking
8. **H. Edge Cases** - Auto-generate boundary questions

### Step 3: Generate Ambiguity List

For each `?` or missing item, create an ambiguity:

```yaml
ambiguities:
  - id: "AMB-ENT-001"
    entity: "Station"
    category: "attribute"
    aspect: "name.unique"
    question: "Must station name be unique?"
    severity: "important"  # critical, important, minor
    suggested_default: null  # or a sensible default

  - id: "AMB-ENT-002"
    entity: "Station"
    category: "lifecycle"
    aspect: "deletion.cascade"
    question: "What happens to scheduled tasks when station is deleted?"
    severity: "critical"
    suggested_default: null
```

### Step 4: Severity Classification

| Severity | Criteria |
|----------|----------|
| **critical** | Blocks implementation - can't write code without knowing |
| **important** | Affects behavior - code would work but might be wrong |
| **minor** | Has sensible default - can propose and confirm |

**Critical examples:**
- Data types for attributes
- Required vs optional
- Deletion cascade behavior
- State transition rules

**Important examples:**
- Validation rules
- Default values
- Audit requirements
- Concurrent access behavior

**Minor examples:**
- Field length limits
- Specific error messages
- History retention period

## Output Format

### Ambiguities Only (default)

```markdown
## Entity Analysis: Ambiguities Found

**Entities Analyzed:** 5
**Total Ambiguities:** 47

### Critical (15)

| ID | Entity | Aspect | Question |
|----|--------|--------|----------|
| AMB-ENT-001 | Station | deletion | What happens to tasks when station deleted? |
| ... |

### Important (22)

| ID | Entity | Aspect | Question |
|----|--------|--------|----------|
| ... |

### Minor (10)

| ID | Entity | Aspect | Question | Suggested Default |
|----|--------|--------|----------|-------------------|
| ... |
```

### Full Report

Includes the ambiguity list PLUS full checklist results showing what IS defined.

## Quality Criteria

Before outputting, verify:
- [ ] Every entity has been analyzed
- [ ] Every attribute has type/required/default checked
- [ ] Every enum has values/transitions checked
- [ ] Every relationship has cascade behavior checked
- [ ] Every lifecycle phase checked
- [ ] Edge cases generated for each entity
- [ ] Severities correctly assigned
- [ ] No duplicate ambiguity IDs
