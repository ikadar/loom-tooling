# Derivation Prompt

You are an expert requirements engineer. Based on the domain model and resolved decisions, generate:

1. ACCEPTANCE CRITERIA in Given/When/Then format
2. BUSINESS RULES with enforcement mechanisms

## Format for Acceptance Criteria

```markdown
### AC-{DOMAIN}-{NUM} – {Title}
**Given** [precondition]
**When** [action with specific inputs]
**Then** [observable outcome]

**Error Cases:**
- {condition} → {behavior}

**Traceability:**
- Source: {source reference}
- Decisions: {decision IDs used}
```

## Format for Business Rules

```markdown
### BR-{DOMAIN}-{NUM} – {Title}
**Rule:** [Clear statement of the constraint]
**Invariant:** [Formal condition using MUST/MUST NOT]
**Enforcement:** [Where and how it's enforced]
**Violation:** [What happens on violation]

**Traceability:**
- Source: {source reference}
- Decisions: {decision IDs used}
```

## Coverage Requirements

Generate comprehensive ACs and BRs covering:
- All happy paths
- All error cases from decisions
- All business rules and constraints
- All state transitions
- All authorization requirements

## Output Format

Return as JSON:
```json
{
  "acceptance_criteria": [
    {
      "id": "AC-XXX-001",
      "title": "Title",
      "given": "precondition",
      "when": "action",
      "then": "outcome",
      "error_cases": ["case1", "case2"],
      "source_refs": ["source1"],
      "decision_refs": ["AMB-001"]
    }
  ],
  "business_rules": [
    {
      "id": "BR-XXX-001",
      "title": "Title",
      "rule": "statement",
      "invariant": "MUST condition",
      "enforcement": "how enforced",
      "error_code": "ERROR_CODE",
      "source_refs": ["source1"],
      "decision_refs": ["AMB-001"]
    }
  ]
}
```

---

DOMAIN MODEL:
