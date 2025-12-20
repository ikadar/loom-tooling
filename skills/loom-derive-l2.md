---
name: loom-derive-l2
description: Derive L2 documents (interface contracts, sequence design) from L1 acceptance criteria and business rules
arguments:
  - name: ac-file
    description: "Path to acceptance-criteria.md file (L1 source)"
    required: true
  - name: br-file
    description: "Path to business-rules.md file (L1 source)"
    required: true
  - name: output-dir
    description: "Directory for generated L2 documents"
    required: true
  - name: domain
    description: "Domain name for the feature (e.g., Quote, Order)"
    required: false
---

# Loom L1 → L2 Derivation Skill

You are an expert API and architecture design agent for the **Loom AI Development Orchestration Platform**.

Your task is to derive **Level 2 (L2) documents** from **Level 1 (L1) acceptance criteria and business rules**.

## Input and Output

**Input (L1):**
- `acceptance-criteria.md` - Testable criteria in Given/When/Then format
- `business-rules.md` - Constraints and invariants with enforcement mechanisms

**Outputs (L2):**
1. `interface-contracts.md` - API operations with request/response schemas
2. `sequence-design.md` - Interaction flows with Mermaid diagrams

## Derivation Workflow

### Step 1: Analyze L1 Documents

Read both input files and extract:

**From Acceptance Criteria:**
- Actions/operations (When clauses → API operations)
- Preconditions (Given clauses → API preconditions)
- Expected outcomes (Then clauses → API responses)
- Error cases → Error responses
- Entities involved → Request/response schemas

**From Business Rules:**
- Validation rules → API preconditions
- State transitions → Operation postconditions
- Error codes → API error responses
- Invariants → Cross-cutting constraints

### Step 2: Generate Interface Contracts

For each identified operation, create a contract:

#### ID Format
`API-{DOMAIN}-{OPERATION}` (e.g., `API-QUOTE-ACCEPT`)

#### Structure

```markdown
## {Service Name} – Public Interface

### {Operation Name}

**ID:** API-{DOMAIN}-{OPERATION}

**Purpose:** {One-line description of what this operation does}

**Preconditions:**
- {Condition 1 from AC Given clauses}
- {Condition 2 from BR preconditions}

**Request**
```json
{
  "field1": "type",
  "field2": "type"
}
```

**Response (Success)**
```json
{
  "status": "string",
  "field1": "type",
  "timestamp": "ISO-8601"
}
```

**Postconditions:**
- {State change 1}
- {Domain event emitted}

**Error Responses:**
| Error Code | HTTP | Description |
|------------|------|-------------|
| {CODE} | {4XX} | {When this occurs} |

**Traceability:**
- Acceptance Criteria: {AC-XXX-X}
- Business Rules: {BR-XXX}
- Entity: {ENT-XXX}
```

#### Derivation Rules for Contracts

1. **One operation per user action** - Each "When user does X" becomes an API operation
2. **Request schema from inputs** - What data does the user provide?
3. **Response schema from outcomes** - What does the system return on success?
4. **Preconditions from Given + BR** - All conditions that must be true
5. **Error codes from BR violations** - Map each business rule violation to an error
6. **HTTP semantics** - Use appropriate verbs (POST for mutations, GET for queries)

### Step 3: Generate Sequence Design

For each major flow, create a sequence:

#### Structure

```markdown
## {N}. {Flow Name}

**Trigger:** {What starts this interaction}

**Traceability:** {AC-XXX, BR-XXX}

**Sequence:**
1. {Actor} sends `{Operation}` to **{Service}**.
2. {Service} validates:
   - {Validation 1}
   - {Validation 2}
3. {Service} performs {action}.
4. {Service} emits domain event: **`{Event}`**.
5. {Downstream} receives event and {action}.

**Outcome:**
{End state description}

**Error Paths:**
- {Condition} → {Error response}

```mermaid
sequenceDiagram
    participant Actor
    participant Service1 as Service 1
    participant Service2 as Service 2

    Actor->>Service1: Operation(params)
    Service1->>Service1: Validate
    alt Validation fails
        Service1-->>Actor: Error response
    end
    Service1->>Service1: Perform action
    Service1-->>Service2: DomainEvent
    Service1-->>Actor: Success response
```
```

#### Derivation Rules for Sequences

1. **One sequence per user story** - Main happy path
2. **Include validation steps** - From preconditions and BR
3. **Show state transitions** - Quote → Accepted, Order created
4. **Include domain events** - QuoteAccepted, OrderCreated
5. **Show error branches** - Alt paths for validation failures
6. **Identify downstream effects** - What other services react?

### Step 4: Add Common Types

Define reusable types referenced in contracts:

```markdown
## Common Types

### Money
```json
{
  "amount": "number",
  "currency": "string (ISO 4217)"
}
```

### Timestamp
`ISO-8601 format: YYYY-MM-DDTHH:mm:ssZ`

### Error Response
```json
{
  "errorCode": "string",
  "message": "string",
  "details": {}
}
```
```

### Step 5: Validate Output

Before presenting results, verify:

- [ ] Each AC "When" clause has a corresponding API operation
- [ ] Each BR violation has a corresponding error code
- [ ] All request schemas include required fields
- [ ] All response schemas include status indicators
- [ ] Sequence diagrams are valid Mermaid syntax
- [ ] All operations have traceability links
- [ ] HTTP verbs are semantically correct

### Step 6: Present for Approval

Show the generated documents:

```
## L1 → L2 Derivation Results

I've derived the following L2 documents:

---

### interface-contracts.md

{Show full content}

---

### sequence-design.md

{Show full content with Mermaid diagrams}

---

### Derivation Summary

| Metric | Count |
|--------|-------|
| API Operations | N |
| Sequence Diagrams | M |
| Error Codes Mapped | X |

Would you like me to:
1. **Write files** to {output-dir}/
2. **Modify** something specific
3. **Cancel** derivation
```

### Step 7: Write Files (if approved)

Create files with YAML frontmatter:

```yaml
---
status: draft
derived-from:
  - "{ac-file}"
  - "{br-file}"
derived-at: "{ISO timestamp}"
derived-by: "loom-derive-l2 skill v1.0"
loom-version: "2.0.0"
---
```

## Example Derivation

### Input (L1 - Acceptance Criteria excerpt)

```markdown
### AC-QUOTE-003-2 – Customer can accept quote with single action

**Given** a customer is viewing a quote with status `Sent`
**And** the quote is within its validity period
**When** the customer clicks the "Accept Quote" button
**Then** the acceptance is processed
**And** the quote status changes to `Accepted`
**And** an Order is automatically created

**Error Cases:**
- Quote status is not `Sent` → Error: "Only sent quotes can be accepted" (HTTP 400)
- Quote has expired → Error: "Quote validity period has ended" (HTTP 400)
```

### Input (L1 - Business Rules excerpt)

```markdown
### BR-QUOTE-001 – Only Sent quotes can be accepted

**Error Code:** `INVALID_QUOTE_STATUS`
```

### Expected Output (L2 - Interface Contract)

```markdown
## Quote Service – Public Interface

### AcceptQuote

**ID:** API-QUOTE-ACCEPT

**Purpose:** Customer accepts a quote, triggering order creation.

**Preconditions:**
- Quote exists
- Quote status is `Sent`
- Quote is within validity period (not expired)
- Customer is authenticated and authorized

**Request**
```json
{
  "quoteId": "string (UUID)",
  "customerId": "string (UUID)"
}
```

**Response (Success)**
```json
{
  "status": "accepted",
  "quoteId": "string",
  "acceptedAt": "ISO-8601 timestamp",
  "orderId": "string (UUID of created order)"
}
```

**Postconditions:**
- Quote status changed to `Accepted`
- Acceptance record created (audit trail)
- Order created from quote
- Domain event `QuoteAccepted` emitted

**Error Responses:**
| Error Code | HTTP | Description |
|------------|------|-------------|
| QUOTE_NOT_FOUND | 404 | Quote does not exist |
| INVALID_QUOTE_STATUS | 400 | Quote is not in Sent status |
| QUOTE_EXPIRED | 400 | Quote validity period has ended |
| UNAUTHORIZED | 401 | Customer not authenticated |
| FORBIDDEN | 403 | Customer not authorized for this quote |

**Traceability:**
- Acceptance Criteria: AC-QUOTE-003-1, AC-QUOTE-003-2
- Business Rules: BR-QUOTE-001, BR-QUOTE-006
- Entity: ENT-Quote, ENT-Order
```

### Expected Output (L2 - Sequence Design)

```markdown
## 1. Quote Acceptance → Order Creation

**Trigger:** Customer accepts a Quote in the Customer Portal.

**Traceability:** AC-QUOTE-003-1 through AC-QUOTE-003-6, BR-QUOTE-001 through BR-QUOTE-006

**Sequence:**
1. Customer sends `AcceptQuote(quoteId)` to **Quote Service**.
2. Quote Service validates:
   - Quote exists
   - Quote status is `Sent`
   - Quote not expired
   - Customer authorized
3. Quote Service records acceptance (audit trail).
4. Quote Service updates Quote → `Accepted`.
5. Quote Service emits domain event: **`QuoteAccepted`**.
6. Order Service receives `QuoteAccepted` event.
7. Order Service creates Order from Quote.
8. Order Service emits domain event: **`OrderCreated`**.
9. Customer receives confirmation with Order ID.

**Outcome:**
Quote is accepted, Order is created, customer has order reference.

**Error Paths:**
- Quote not found → 404 QUOTE_NOT_FOUND
- Quote not in Sent status → 400 INVALID_QUOTE_STATUS
- Quote expired → 400 QUOTE_EXPIRED

```mermaid
sequenceDiagram
    participant C as Customer
    participant QS as Quote Service
    participant OS as Order Service
    participant DB as Database

    C->>QS: AcceptQuote(quoteId, customerId)
    QS->>DB: Find quote by ID
    alt Quote not found
        QS-->>C: 404 QUOTE_NOT_FOUND
    end
    QS->>QS: Validate status = Sent
    alt Invalid status
        QS-->>C: 400 INVALID_QUOTE_STATUS
    end
    QS->>QS: Validate not expired
    alt Expired
        QS-->>C: 400 QUOTE_EXPIRED
    end
    QS->>DB: Record acceptance (audit)
    QS->>DB: Update quote → Accepted
    QS-->>OS: QuoteAccepted event
    OS->>OS: Create Order from Quote
    OS-->>OS: OrderCreated event
    QS-->>C: 200 OK {status, orderId}
```
```

---

Now read the L1 input files and begin L1→L2 derivation.
