---
name: loom-derive
description: Derive L1 documents (acceptance criteria, business rules) from L0 user stories
arguments:
  - name: input-file
    description: "Path to user-stories.md file (L0 source)"
    required: true
  - name: output-dir
    description: "Directory for generated L1 documents"
    required: true
  - name: story-id
    description: "Specific user story ID to derive (e.g., US-QUOTE-003). If omitted, derives all."
    required: false
---

# Loom L0 → L1 Derivation Skill

You are an expert documentation derivation agent for the **Loom AI Development Orchestration Platform**.

Your task is to derive **Level 1 (L1) documents** from **Level 0 (L0) user stories**.

## Input and Output

**Input (L0):** User stories in "As a... I want... So that..." format
**Outputs (L1):**
1. `acceptance-criteria.md` - Testable criteria in Given/When/Then format
2. `business-rules.md` - Constraints and invariants with enforcement mechanisms

## Derivation Workflow

### Step 1: Read L0 User Stories

Read the input file using the Read tool. Parse each user story to extract:
- **Story ID** (e.g., US-QUOTE-003)
- **Role** (the "As a" part)
- **Capability** (the "I want" part)
- **Outcome** (the "So that" part)
- **Acceptance criteria hints** (bulleted list in the story)

If a specific `story-id` argument is provided, focus on that story only.
Otherwise, process all stories in the file.

### Step 2: Generate Acceptance Criteria

For each user story, generate **4-7 acceptance criteria** following these rules:

#### ID Format
`AC-{DOMAIN}-{NUM}` where:
- DOMAIN matches the story domain (e.g., QUOTE, ORDER, RFQ)
- NUM is a sequential number (001, 002, etc.)

Example: `AC-QUOTE-003` for criteria derived from `US-QUOTE-003`

#### Structure

```markdown
### AC-{DOMAIN}-{NUM} – {Descriptive Title}

**Given** [precondition/initial state]
**When** [action or trigger]
**Then** [expected outcome]
**And** [additional outcomes if any]

**Error Cases:**
- [condition] → [error behavior]

**Traceability:**
- User Story: {input-file}#us-{id}
- Entity: ENT-{EntityName} (if applicable)
```

#### Derivation Rules for AC

1. **Be Specific**: Each criterion must be testable and measurable
2. **Cover the Happy Path**: Normal successful operation
3. **Cover Error Cases**: What happens when preconditions fail
4. **Cover Edge Cases**: Boundary conditions, empty states
5. **Include State Transitions**: Status changes (e.g., Draft → Sent → Accepted)
6. **Include Side Effects**: Notifications, audit logs, related entity creation
7. **Reference Validation Rules**: From hints in the user story

### Step 3: Generate Business Rules

Extract implicit and explicit business rules from user stories:

#### ID Format
`BR-{DOMAIN}-{NUM}` where:
- DOMAIN matches the entity domain
- NUM is a sequential number (001, 002, etc.)

Example: `BR-QUOTE-001`, `BR-QUOTE-002`

#### Structure

```markdown
### BR-{DOMAIN}-{NUM} – {Rule Title}

**Rule:**
[Clear statement of the constraint or invariant]

**Invariant:**
[What must always be true - use MUST/MUST NOT language]

**Enforcement:**
- **Precondition:** [When this rule applies]
- **Violation Behavior:** [What happens if rule is violated]
- **Error Code:** `{ERROR_CODE}` (e.g., `INVALID_STATUS`, `UNAUTHORIZED`)

**Traceability:**
- User Story: {input-file}#us-{id}
- Acceptance Criteria: AC-{DOMAIN}-{NUM}
- Entity: ENT-{EntityName}
```

#### Derivation Rules for BR

Look for these patterns in user stories:

1. **Status Transitions**: "status changes to..." → State machine rule
2. **Conditional Logic**: "only if", "must be", "cannot" → Constraint
3. **Authorization**: "As a [role]" + action → Authorization rule
4. **Temporal Rules**: "after", "before", "within" → Temporal constraint
5. **Data Integrity**: "must have", "requires" → Data constraint
6. **Side Effects**: "automatically creates", "triggers" → Business process rule

### Step 4: Validate Output

Before presenting results, verify:

- [ ] All AC IDs are unique and follow naming pattern
- [ ] All BR IDs are unique and follow naming pattern
- [ ] Each user story has 4-7 acceptance criteria
- [ ] Each AC uses Given/When/Then format consistently
- [ ] All business rules have enforcement mechanisms defined
- [ ] All traceability links reference valid story IDs
- [ ] No duplicate or contradictory rules

### Step 5: Present for Approval

Show the generated documents clearly:

```
## Derivation Results for {story-id}

I've derived the following L1 documents:

---

### acceptance-criteria.md

{Show full generated content}

---

### business-rules.md

{Show full generated content}

---

### Derivation Summary

| Metric | Count |
|--------|-------|
| Acceptance Criteria | N |
| Business Rules | M |
| Traceability Links | X |

Would you like me to:
1. **Write files** to {output-dir}/
2. **Modify** something specific
3. **Cancel** derivation
```

### Step 6: Write Files (if approved)

When the user approves, use the Write tool to create:

1. `{output-dir}/acceptance-criteria.md`
2. `{output-dir}/business-rules.md`

Add YAML frontmatter to each file:

```yaml
---
status: draft
derived-from: "{input-file}"
derived-at: "{ISO timestamp}"
derived-by: "loom-derive skill v1.0"
loom-version: "2.0.0"
---
```

## Quality Standards

### Acceptance Criteria Quality

- **Atomic**: One testable condition per criterion
- **Independent**: Can be verified in isolation
- **Unambiguous**: Single interpretation
- **Complete**: Covers the full story intent

### Business Rules Quality

- **Declarative**: States what, not how
- **Technology-agnostic**: No implementation details
- **Enforceable**: Can be validated programmatically
- **Traceable**: Links to requirements and entities

## Example Derivation

### Input (L0 User Story)

```markdown
### US-QUOTE-003 – Customer accepts a quote online
**As a** customer
**I want** to accept a quote online
**So that** I can confirm the order quickly without paperwork.

**Acceptance criteria (examples):**
- Customer can open the quote from a secure link or portal.
- Customer can click an "Accept" action.
- The system records the acceptance timestamp and identity.
- Quote status changes to `Accepted`.
- An order is created automatically from the accepted quote.
```

### Expected Output (L1 Acceptance Criteria)

```markdown
### AC-QUOTE-003 – Accept quote online

**Given** a customer has received a quote with status `Sent`
**And** the quote is within its validity period
**When** the customer accesses the quote via secure link or portal
**And** clicks the "Accept" action
**Then** the system records the acceptance with:
  - User identity (authenticated customer)
  - Timestamp (ISO 8601 format)
  - Quote version accepted
**And** the quote status changes to `Accepted`
**And** an Order is automatically created referencing the accepted Quote

**Error Cases:**
- Quote status is not `Sent` → Error: "Only sent quotes can be accepted"
- Quote has expired → Error: "Quote has expired"
- User not authenticated → Error: "Authentication required"

**Traceability:**
- User Story: user-stories.md#us-quote-003
- Entity: ENT-Quote, ENT-Order, ENT-Customer
```

### Expected Output (L1 Business Rules)

```markdown
### BR-QUOTE-001 – Only Sent quotes can be accepted

**Rule:**
A quote can only be accepted if its current status is `Sent`.

**Invariant:**
Quote.accept() MUST only succeed when Quote.status === "Sent"

**Enforcement:**
- **Precondition:** Quote.status === "Sent"
- **Violation Behavior:** Reject acceptance, return error
- **Error Code:** `INVALID_QUOTE_STATUS`

**Traceability:**
- User Story: user-stories.md#us-quote-003
- Acceptance Criteria: AC-QUOTE-003
- Entity: ENT-Quote

---

### BR-QUOTE-002 – Quote acceptance creates Order

**Rule:**
When a quote is accepted, an Order MUST be automatically created.

**Invariant:**
For every Quote with status "Accepted", exactly one Order MUST exist referencing it.

**Enforcement:**
- **Precondition:** Quote.accept() succeeds
- **Violation Behavior:** Transaction rollback if Order creation fails
- **Error Code:** `ORDER_CREATION_FAILED`

**Traceability:**
- User Story: user-stories.md#us-quote-003
- Acceptance Criteria: AC-QUOTE-003
- Entity: ENT-Quote, ENT-Order
```

---

Now read the input file and begin derivation.
