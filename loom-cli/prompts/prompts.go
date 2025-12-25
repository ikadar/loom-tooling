package prompts

// DomainDiscovery is the prompt for Phase 1
const DomainDiscovery = `You are a domain analysis expert. Extract the domain model from the following L0 specification.

Return a JSON object with this structure:
{
  "entities": [
    {
      "name": "EntityName",
      "mentioned_attributes": ["attr1", "attr2"],
      "mentioned_operations": ["op1", "op2"],
      "mentioned_states": ["state1", "state2"]
    }
  ],
  "operations": [
    {
      "name": "OperationName",
      "actor": "ActorName",
      "trigger": "how it's triggered",
      "target": "TargetEntity",
      "mentioned_inputs": ["input1", "input2"],
      "mentioned_rules": ["rule1", "rule2"]
    }
  ],
  "relationships": [
    {
      "from": "Entity1",
      "to": "Entity2",
      "type": "contains|references|extends",
      "cardinality": "1:1|1:N|N:M"
    }
  ],
  "business_rules": ["Rule statement 1", "Rule statement 2"],
  "ui_mentions": ["UI component 1", "UI element 2"]
}

Be thorough - extract EVERY entity, operation, relationship, rule, and UI mention from the input.

L0 INPUT:
`

// EntityAnalysis is the prompt for analyzing entities
const EntityAnalysis = `You are a requirements completeness analyst. Analyze each entity against this checklist and identify ambiguities.

ENTITY COMPLETENESS CHECKLIST:
1. IDENTITY: Is there a unique identifier? What format?
2. ATTRIBUTES: Are all attributes listed? Types? Required vs optional?
3. LIFECYCLE: Creation, modification, deletion rules? Who can do each?
4. STATES: What states can it have? State transitions?
5. VALIDATION: What validation rules apply to each attribute?
6. RELATIONSHIPS: Are all relationships to other entities clear? Cardinality?
7. HISTORY: Is change tracking needed? Audit trail?
8. SOFT DELETE: Is deletion permanent or soft?
9. DEFAULTS: What are default values for optional attributes?
10. CONSTRAINTS: Any business constraints not yet covered?

For each ambiguity found, return JSON:
{
  "ambiguities": [
    {
      "id": "AMB-ENT-001",
      "category": "entity",
      "subject": "EntityName",
      "question": "Clear question about what's missing or unclear",
      "severity": "critical|important|minor",
      "suggested_answer": "A reasonable default if applicable",
      "options": ["Option A", "Option B", "Option C"],
      "checklist_item": "Which checklist item this relates to"
    }
  ]
}

Severity guide:
- critical: Blocks implementation, must be answered
- important: Affects multiple features, should be answered
- minor: Nice to have, can use default

ENTITIES TO ANALYZE:
`

// OperationAnalysis is the prompt for analyzing operations
const OperationAnalysis = `You are a requirements completeness analyst. Analyze each operation against this checklist and identify ambiguities.

OPERATION COMPLETENESS CHECKLIST:
1. INPUTS: Are all inputs defined? Types? Required vs optional?
2. OUTPUTS: What does the operation return? Success/failure indicators?
3. PRECONDITIONS: What must be true before execution?
4. POSTCONDITIONS: What is guaranteed after execution?
5. AUTHORIZATION: Who can perform this? Role/permission requirements?
6. VALIDATION: What input validation is needed?
7. ERROR HANDLING: What errors can occur? How to handle each?
8. SIDE EFFECTS: What other state changes occur?
9. IDEMPOTENCY: Can it be safely retried?
10. CONCURRENCY: What happens with concurrent executions?
11. PERFORMANCE: Any timeout or performance requirements?
12. AUDIT: Should this be logged/audited?

For each ambiguity found, return JSON:
{
  "ambiguities": [
    {
      "id": "AMB-OP-001",
      "category": "operation",
      "subject": "OperationName",
      "question": "Clear question about what's missing or unclear",
      "severity": "critical|important|minor",
      "suggested_answer": "A reasonable default if applicable",
      "options": ["Option A", "Option B", "Option C"],
      "checklist_item": "Which checklist item this relates to"
    }
  ]
}

OPERATIONS TO ANALYZE:
`

// DerivationPrompt is the prompt for Phase 5 - generating AC and BR
const DerivationPrompt = `You are an expert requirements engineer. Based on the domain model and resolved decisions, generate:

1. ACCEPTANCE CRITERIA in Given/When/Then format
2. BUSINESS RULES with enforcement mechanisms

FORMAT FOR ACCEPTANCE CRITERIA:
### AC-{DOMAIN}-{NUM} – {Title}
**Given** [precondition]
**When** [action with specific inputs]
**Then** [observable outcome]

**Error Cases:**
- {condition} → {behavior}

**Traceability:**
- Source: {source reference}
- Decisions: {decision IDs used}

FORMAT FOR BUSINESS RULES:
### BR-{DOMAIN}-{NUM} – {Title}
**Rule:** [Clear statement of the constraint]
**Invariant:** [Formal condition using MUST/MUST NOT]
**Enforcement:** [Where and how it's enforced]
**Violation:** [What happens on violation]

**Traceability:**
- Source: {source reference}
- Decisions: {decision IDs used}

Generate comprehensive ACs and BRs covering:
- All happy paths
- All error cases from decisions
- All business rules and constraints
- All state transitions
- All authorization requirements

Return as JSON:
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

DOMAIN MODEL:
`

// InterviewPrompt is the prompt for structuring interview questions
const InterviewPrompt = `You are conducting a structured interview to resolve specification ambiguities.

Present questions clearly and concisely. Group related questions together.

For each ambiguity:
1. Provide context (what entity/operation it relates to)
2. Ask the specific question
3. Offer options if available
4. Suggest a default if applicable

Format output as:
---
**[Category: Subject]**

Q: {question}

Options:
a) {option1}
b) {option2}
c) {option3}
d) Other (please specify)

Suggested default: {suggestion}
---

AMBIGUITIES TO ASK:
`
