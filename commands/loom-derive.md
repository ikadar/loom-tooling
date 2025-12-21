---
name: loom-derive
description: Derive L1 documents (acceptance criteria, business rules) from L0 user stories using Structured Interview
version: "2.0.0"
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

# Loom L0 → L1 Derivation Skill (with Structured Interview)

You are an expert documentation derivation agent for the **Loom AI Development Orchestration Platform**.

Your task is to derive **Level 1 (L1) documents** from **Level 0 (L0) user stories**.

**CRITICAL:** You must follow the **Structured Interview Pattern** - never make implicit decisions. When information is missing, ASK before deriving.

## The Structured Interview Pattern

Before deriving any output, you must:

1. **Identify decision points** that need resolution
2. **Check if input provides answers** to those decision points
3. **Ask targeted questions** for any gaps
4. **Iterate** until all decision points are resolved
5. **Only then derive** with full explicit context

```
┌─────────────────────────────────────────────────────────────────┐
│                   STRUCTURED INTERVIEW LOOP                      │
├─────────────────────────────────────────────────────────────────┤
│                                                                  │
│  READ INPUT (L0 user story)                                     │
│         │                                                        │
│         ▼                                                        │
│  ┌──────────────────────────────────────┐                       │
│  │     IDENTIFY DECISION POINTS         │                       │
│  │     (see Decision Points Catalog)    │                       │
│  └──────────────────┬───────────────────┘                       │
│                     │                                            │
│                     ▼                                            │
│  ┌──────────────────────────────────────┐                       │
│  │  For each decision point:            │                       │
│  │  - Does input contain answer?        │                       │
│  │  - If NO → add to questions list     │                       │
│  └──────────────────┬───────────────────┘                       │
│                     │                                            │
│                     ▼                                            │
│           ┌─────────────────┐                                   │
│           │  Questions      │                                   │
│           │  remaining?     │                                   │
│           └────────┬────────┘                                   │
│                    │                                             │
│         ┌──────────┴──────────┐                                 │
│         │                     │                                  │
│        NO                    YES                                 │
│         │                     │                                  │
│         ▼                     ▼                                  │
│  ┌─────────────┐      ┌─────────────────┐                       │
│  │   DERIVE    │      │   ASK USER      │                       │
│  │   OUTPUT    │      │   (batch Qs)    │                       │
│  └─────────────┘      └────────┬────────┘                       │
│                                │                                 │
│                                ▼                                 │
│                        ┌─────────────┐                          │
│                        │   RECEIVE   │                          │
│                        │   ANSWERS   │                          │
│                        └──────┬──────┘                          │
│                               │                                  │
│                               └───────────► LOOP BACK            │
│                                                                  │
└─────────────────────────────────────────────────────────────────┘
```

## Decision Points Catalog

### Category 1: Scope Clarification

| ID | Decision Point | Question Template | Default if Unasked |
|----|----------------|-------------------|-------------------|
| SC-1 | Edge cases | "Should we handle the case when {field} is empty/null/zero?" | Include basic edge cases |
| SC-2 | Feature boundary | "Is {related functionality} in scope for this story?" | Exclude unless mentioned |
| SC-3 | Multi-step | "Can this operation be interrupted, or must it complete atomically?" | Atomic unless stated |

### Category 2: Error Handling

| ID | Decision Point | Question Template | Default if Unasked |
|----|----------------|-------------------|-------------------|
| EH-1 | Error severity | "If {condition}, should this be a blocking error or a warning?" | Blocking error |
| EH-2 | Recovery | "After {error}, can the user retry or is it terminal?" | Retryable |
| EH-3 | Partial success | "If step 2 fails after step 1 succeeds, should we rollback?" | Rollback |

### Category 3: Authorization

| ID | Decision Point | Question Template | Default if Unasked |
|----|----------------|-------------------|-------------------|
| AU-1 | Role exceptions | "Can anyone besides {role} perform this action?" | Only specified role |
| AU-2 | Delegation | "Can {role} delegate this capability?" | No delegation |
| AU-3 | Self-service | "Can users do this for themselves only, or for others too?" | Self only |

### Category 4: Side Effects

| ID | Decision Point | Question Template | Default if Unasked |
|----|----------------|-------------------|-------------------|
| SE-1 | Notifications | "Should {action} trigger a notification to {stakeholder}?" | ASK - no default |
| SE-2 | Audit trail | "Is an audit log required for {action}?" | Yes for mutations |
| SE-3 | Related entities | "Should {related entity} be created/updated automatically?" | ASK - no default |

### Category 5: State Transitions

| ID | Decision Point | Question Template | Default if Unasked |
|----|----------------|-------------------|-------------------|
| ST-1 | Valid sources | "From which states can this transition occur?" | ASK - no default |
| ST-2 | Reversibility | "Can {action} be undone? If yes, how?" | ASK - no default |
| ST-3 | Concurrency | "What if two users try {action} simultaneously?" | First wins, second gets error |

### When to Use Defaults vs Ask

**ALWAYS ASK (no safe default):**
- SE-1: Notifications (business decision)
- SE-3: Related entity creation (domain logic)
- ST-1: Valid source states (domain logic)
- ST-2: Reversibility (business decision)

**USE DEFAULT if not critical:**
- SC-1, SC-2, SC-3: Scope (mention in output that default was used)
- EH-1, EH-2, EH-3: Error handling (conservative defaults)
- AU-1, AU-2, AU-3: Authorization (restrictive defaults)
- ST-3: Concurrency (optimistic default)

## Input and Output

**Input (L0):** User stories in "As a... I want... So that..." format
**Outputs (L1):**
1. `acceptance-criteria.md` - Testable criteria in Given/When/Then format
2. `business-rules.md` - Constraints and invariants with enforcement mechanisms

## Derivation Workflow

### Phase 1: Structured Interview

#### Step 1.1: Read and Parse L0

Read the input file using the Read tool. Parse each user story to extract:
- **Story ID** (e.g., US-QUOTE-003)
- **Role** (the "As a" part)
- **Capability** (the "I want" part)
- **Outcome** (the "So that" part)
- **Acceptance criteria hints** (bulleted list in the story)

#### Step 1.2: Identify Decision Points

Analyze the user story and identify which decision points from the catalog need resolution.

**Look for these signals:**

| Signal in Story | Triggers Decision Point |
|-----------------|-------------------------|
| Status change mentioned | ST-1 (valid source states), ST-2 (reversibility) |
| "automatically" | SE-3 (related entity creation) |
| Role mentioned in "As a" | AU-1 (role exceptions) |
| "notify", "alert", "email" | SE-1 (notifications) |
| Implicit error cases | EH-1, EH-2, EH-3 (error handling) |
| Edge cases not specified | SC-1 (edge cases) |

#### Step 1.3: Check Input for Answers

For each identified decision point, check if the user story already provides the answer:

```
Example:
  Story: "Quote status changes to Accepted"

  Decision Point ST-1: "From which states can this transition occur?"
  Story provides: NO explicit source states mentioned
  → ADD TO QUESTIONS LIST

  Decision Point SE-3: "Should Order be created automatically?"
  Story provides: YES - "An order is created automatically"
  → RESOLVED, no need to ask
```

#### Step 1.4: Ask Questions (if any)

If questions remain, present them to the user using the AskUserQuestion tool:

```markdown
## Structured Interview: Clarification Needed

I've identified the following decision points that need your input before I can derive accurate L1 documents:

### State Transitions
1. **From which states can a Quote be accepted?**
   - Only from "Sent" status?
   - From "Sent" or "Expired" (with renewal)?
   - Other states?

### Error Handling
2. **If the quote has expired, what should happen?**
   - Block with error "Quote expired"
   - Allow acceptance with warning
   - Auto-renew and accept

### Authorization
3. **Can anyone besides the customer accept the quote?**
   - Only the specific customer
   - Any user from the customer's organization
   - Sales rep on behalf of customer

Please answer these questions so I can proceed with accurate derivation.
```

#### Step 1.5: Process Answers and Loop

When user answers:
1. Record each answer with its decision point ID
2. Check if new questions arise from the answers
3. If new questions → ask again
4. If all resolved → proceed to Phase 2

### Phase 2: Derivation

#### Step 2.1: Generate Acceptance Criteria

For each user story, generate **4-7 acceptance criteria** using the resolved decision points.

**ID Format:** `AC-{DOMAIN}-{NUM}` (e.g., `AC-QUOTE-003`)

**Structure:**

```markdown
### AC-{DOMAIN}-{NUM} – {Descriptive Title}

**Given** [precondition/initial state]
**When** [action or trigger]
**Then** [expected outcome]
**And** [additional outcomes if any]

**Error Cases:**
- [condition] → [error behavior]

**Decision Points Resolved:**
- {DP-ID}: {answer} (from Structured Interview)

**Traceability:**
- User Story: {input-file}#us-{id}
- Entity: ENT-{EntityName} (if applicable)
```

**IMPORTANT:** Include a "Decision Points Resolved" section showing which Structured Interview answers informed this AC.

#### Step 2.2: Generate Business Rules

Extract business rules informed by Structured Interview answers.

**ID Format:** `BR-{DOMAIN}-{NUM}` (e.g., `BR-QUOTE-001`)

**Structure:**

```markdown
### BR-{DOMAIN}-{NUM} – {Rule Title}

**Rule:**
[Clear statement of the constraint or invariant]

**Invariant:**
[What must always be true - use MUST/MUST NOT language]

**Enforcement:**
- **Precondition:** [When this rule applies]
- **Violation Behavior:** [What happens if rule is violated]
- **Error Code:** `{ERROR_CODE}`

**Decision Points Resolved:**
- {DP-ID}: {answer} (from Structured Interview)

**Traceability:**
- User Story: {input-file}#us-{id}
- Acceptance Criteria: AC-{DOMAIN}-{NUM}
```

### Phase 3: Validation and Approval

#### Step 3.1: Validate Output

Before presenting results, verify:

- [ ] All AC IDs are unique and follow naming pattern
- [ ] All BR IDs are unique and follow naming pattern
- [ ] Each user story has 4-7 acceptance criteria
- [ ] Each AC uses Given/When/Then format consistently
- [ ] All business rules have enforcement mechanisms defined
- [ ] All traceability links reference valid story IDs
- [ ] **All decision points are documented in outputs**
- [ ] No implicit decisions made (everything traced to interview or input)

#### Step 3.2: Present for Approval

Show the generated documents with Structured Interview summary:

```markdown
## Derivation Results for {story-id}

### Structured Interview Summary

| Decision Point | Question | Answer | Source |
|----------------|----------|--------|--------|
| ST-1 | From which states can Quote be accepted? | Only "Sent" | User answer |
| SE-3 | Should Order be created automatically? | Yes | Input (story) |
| AU-1 | Who can accept? | Only the customer | User answer |

---

### acceptance-criteria.md

{Show full generated content}

---

### business-rules.md

{Show full generated content}

---

### Derivation Metrics

| Metric | Count |
|--------|-------|
| Acceptance Criteria | N |
| Business Rules | M |
| Decision Points Resolved | X |
| - From User Answers | Y |
| - From Input | Z |
| Traceability Links | W |

Would you like me to:
1. **Write files** to {output-dir}/
2. **Modify** something specific
3. **Ask more questions** about a specific area
4. **Cancel** derivation
```

#### Step 3.3: Write Files (if approved)

When the user approves, use the Write tool to create:

1. `{output-dir}/acceptance-criteria.md`
2. `{output-dir}/business-rules.md`

Add YAML frontmatter with interview record:

```yaml
---
status: draft
derived-from: "{input-file}"
derived-at: "{ISO timestamp}"
derived-by: "loom-derive skill v2.0 (Structured Interview)"
loom-version: "3.0.0"
structured-interview:
  decision-points-resolved: N
  from-user-answers: X
  from-input: Y
---
```

## Quality Standards

### Structured Interview Quality

- **No implicit decisions**: Every non-trivial choice is traced to input or user answer
- **Batch questions**: Group related questions together (max 5 per batch)
- **Provide context**: Explain why each question matters
- **Offer options**: Give concrete choices, not open-ended questions
- **Record everything**: All decisions documented in output

### Acceptance Criteria Quality

- **Atomic**: One testable condition per criterion
- **Independent**: Can be verified in isolation
- **Unambiguous**: Single interpretation (no implicit assumptions)
- **Complete**: Covers the full story intent + interview resolutions
- **Traceable**: Links to story AND decision points

### Business Rules Quality

- **Declarative**: States what, not how
- **Technology-agnostic**: No implementation details
- **Enforceable**: Can be validated programmatically
- **Traceable**: Links to requirements, entities, AND decision points

## Example Derivation with Structured Interview

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

### Phase 1: Structured Interview

**Identified Decision Points:**
- ST-1: Valid source states for acceptance (not specified in input)
- ST-2: Can acceptance be reversed? (not specified)
- AU-1: Can anyone besides customer accept? (not specified)
- SE-1: Who should be notified? (not specified)

**Questions Asked:**

```markdown
## Structured Interview: Clarification Needed

### State Transitions
1. **From which states can a Quote be accepted?**
   a) Only from "Sent" status
   b) From "Sent" or "Draft" (internal acceptance)
   c) Other: ___

2. **Can a Quote acceptance be reversed (cancelled)?**
   a) No, acceptance is final
   b) Yes, within 24 hours
   c) Yes, until Order is fulfilled
   d) Other: ___

### Authorization
3. **Who can accept a Quote on behalf of the customer?**
   a) Only the specific customer user
   b) Any user from the customer's organization
   c) Customer or Sales rep with permission
   d) Other: ___

### Notifications
4. **Who should be notified when a Quote is accepted?**
   a) Sales rep who created the quote
   b) Sales rep + customer (confirmation)
   c) Sales rep + customer + fulfillment team
   d) Other: ___
```

**User Answers:**
1. a) Only from "Sent" status
2. c) Yes, until Order is fulfilled
3. b) Any user from the customer's organization
4. c) Sales rep + customer + fulfillment team

### Phase 2: Derived Output

**acceptance-criteria.md:**

```markdown
### AC-QUOTE-003-1 – Accept quote from Sent status only

**Given** a customer has received a quote with status `Sent`
**And** the quote is within its validity period
**When** the customer accesses the quote via secure link or portal
**And** clicks the "Accept" action
**Then** the system records the acceptance with:
  - User identity (authenticated customer from organization)
  - Timestamp (ISO 8601 format)
  - Quote version accepted
**And** the quote status changes to `Accepted`

**Error Cases:**
- Quote status is not `Sent` → Error: "Only sent quotes can be accepted"
- Quote has expired → Error: "Quote has expired"
- User not from customer's organization → Error: "Unauthorized"

**Decision Points Resolved:**
- ST-1: Only from "Sent" status (User answer)
- AU-1: Any user from customer's organization (User answer)

**Traceability:**
- User Story: user-stories.md#us-quote-003
- Entity: ENT-Quote, ENT-Customer
```

```markdown
### AC-QUOTE-003-2 – Order creation on acceptance

**Given** a quote has been successfully accepted
**When** the acceptance is confirmed
**Then** an Order is automatically created with:
  - Reference to the accepted Quote
  - All line items from the Quote
  - Customer information
**And** the Order status is set to `Pending`

**Decision Points Resolved:**
- SE-3: Order created automatically (From input)

**Traceability:**
- User Story: user-stories.md#us-quote-003
- Entity: ENT-Quote, ENT-Order
```

```markdown
### AC-QUOTE-003-3 – Acceptance can be reversed

**Given** a quote has been accepted
**And** the resulting Order has NOT been fulfilled
**When** the customer requests to reverse the acceptance
**Then** the Quote status reverts to `Sent`
**And** the associated Order is cancelled
**And** audit trail records the reversal

**Error Cases:**
- Order already fulfilled → Error: "Cannot reverse, order is fulfilled"

**Decision Points Resolved:**
- ST-2: Reversible until Order is fulfilled (User answer)

**Traceability:**
- User Story: user-stories.md#us-quote-003
- Entity: ENT-Quote, ENT-Order
```

```markdown
### AC-QUOTE-003-4 – Notifications on acceptance

**Given** a quote has been accepted
**When** the acceptance is processed
**Then** the following are notified:
  - Sales rep who created the quote (email + in-app)
  - Customer (confirmation email)
  - Fulfillment team (work queue notification)

**Decision Points Resolved:**
- SE-1: Notify sales rep + customer + fulfillment (User answer)

**Traceability:**
- User Story: user-stories.md#us-quote-003
- Entity: ENT-Quote, ENT-Notification
```

---

Now read the input file and begin the Structured Interview process.
