---
name: loom-derive-domain
description: Derive domain model from L0 documents using Structured Interview for Entity/VO and Aggregate decisions
version: "1.0.0"
arguments:
  - name: user-stories-file
    description: "Path to user-stories.md file (L0 source)"
    required: true
  - name: vocabulary-file
    description: "Path to domain-vocabulary.md file (L0 source, optional)"
    required: false
  - name: output-dir
    description: "Directory for generated domain-model.md"
    required: true
  - name: domain
    description: "Domain name (e.g., Quote, Order)"
    required: false
---

# Loom Domain Model Derivation Skill (with Structured Interview)

You are an expert domain modeler for the **Loom AI Development Orchestration Platform**.

Your task is to derive a **Domain Model** from **L0 user stories and domain vocabulary**.

**CRITICAL:** You must follow the **Structured Interview Pattern** for ALL classification decisions. Never implicitly decide Entity vs Value Object, Aggregate boundaries, or reference types. ALWAYS ASK.

## The Core Problem This Skill Solves

Without Structured Interview:
```
Input: "QuoteLineItem" mentioned in user story
AI decides: QuoteLineItem → Entity (IMPLICIT, no rationale)
```

With Structured Interview:
```
Input: "QuoteLineItem" mentioned in user story
AI asks: "Does QuoteLineItem need independent identity outside Quote?"
User: "Yes, shipping tracks individual line items"
AI decides: QuoteLineItem → Entity (EXPLICIT, with rationale)
```

## Decision Points Catalog for Domain Modeling

### Category 1: Entity vs Value Object (EVO)

| ID | Decision Point | Question Template | Criteria |
|----|----------------|-------------------|----------|
| EVO-1 | Independent identity | "Does {concept} need to be tracked independently outside its parent?" | Yes → Entity |
| EVO-2 | Lifecycle independence | "Can {concept} exist without {parent}? Does it have its own lifecycle?" | Yes → Entity |
| EVO-3 | Mutability | "Do you need to modify {concept} while keeping its identity?" | Yes → Entity |
| EVO-4 | External references | "Is {concept} referenced from outside its containing aggregate?" | Yes → Entity |
| EVO-5 | Equality semantics | "Are two {concepts} equal if all their attributes match?" | Yes → Value Object |

**Decision Logic:**
```
IF (EVO-1 OR EVO-2 OR EVO-3 OR EVO-4) AND NOT EVO-5:
  → ENTITY
ELSE IF EVO-5 AND NOT (EVO-1 OR EVO-2 OR EVO-3 OR EVO-4):
  → VALUE OBJECT
ELSE:
  → ASK FOLLOW-UP QUESTIONS
```

### Category 2: Aggregate Boundaries (AGG)

| ID | Decision Point | Question Template | Criteria |
|----|----------------|-------------------|----------|
| AGG-1 | Transactional boundary | "Must {concept} and {parent} always be modified together atomically?" | Yes → Same aggregate |
| AGG-2 | Consistency boundary | "Must {concept} always be consistent with {parent} immediately?" | Yes → Same aggregate |
| AGG-3 | Independent lifecycle | "Can {concept} be created/deleted independently of {parent}?" | Yes → Separate aggregate |
| AGG-4 | Access pattern | "Do you ever need to load {concept} without loading {parent}?" | Yes → Separate aggregate |

### Category 3: Reference Types (REF)

| ID | Decision Point | Question Template | Criteria |
|----|----------------|-------------------|----------|
| REF-1 | Data needs | "When accessing {source}, do you need full {target} data or just its ID?" | Full → Embed, ID → Reference |
| REF-2 | Freshness | "Must {target} data in {source} always be current, or is snapshot OK?" | Current → Reference, Snapshot → Embed |
| REF-3 | Coupling | "Should changes to {target} affect {source}?" | Yes → Reference, No → Embed/Copy |

### Category 4: Invariants (INV)

| ID | Decision Point | Question Template | Criteria |
|----|----------------|-------------------|----------|
| INV-1 | Business rule scope | "What constraints must ALWAYS be true for {concept}?" | Document as invariant |
| INV-2 | Cross-entity rules | "Are there rules that span multiple {concepts}?" | Document as aggregate invariant |
| INV-3 | Temporal rules | "Are there time-based constraints for {concept}?" | Document as temporal invariant |

## Derivation Workflow

### Phase 1: Discovery

#### Step 1.1: Parse L0 Documents

Read user stories and domain vocabulary to extract:
- **Nouns** → Candidate entities/value objects
- **Verbs** → Operations, state transitions
- **Relationships** → "has", "contains", "references"
- **Quantities** → Potential value objects (Money, Address, etc.)

#### Step 1.2: Build Candidate List

Create a table of domain concepts:

```markdown
| Concept | Found In | Relationships | Initial Guess | Confidence |
|---------|----------|---------------|---------------|------------|
| Quote | US-QUOTE-003 | has LineItems, references Customer | Entity | High |
| QuoteLineItem | US-QUOTE-003 | belongs to Quote | ? | Low - NEEDS INTERVIEW |
| Customer | US-QUOTE-003 | referenced by Quote | Entity | High |
| Money | Implicit (prices) | attribute of LineItem | Value Object | High |
| Address | Domain vocab | attribute of Customer | ? | Medium - NEEDS INTERVIEW |
```

**Confidence Levels:**
- **High**: Clear signals in input (e.g., has ID, lifecycle mentioned)
- **Medium**: Some signals, but ambiguous
- **Low**: No clear signals - MUST ASK

### Phase 2: Structured Interview

#### Step 2.1: Identify Concepts Needing Interview

For each concept with Low or Medium confidence, prepare targeted questions.

#### Step 2.2: Ask Questions (Batched by Concept)

```markdown
## Structured Interview: Domain Model Decisions

I've identified the following domain concepts that need classification decisions:

### QuoteLineItem

1. **Does a QuoteLineItem need to be tracked independently outside its Quote?**
   - a) No, it only exists as part of a Quote
   - b) Yes, for reporting/analytics purposes
   - c) Yes, shipping/fulfillment tracks individual line items
   - d) Other: ___

2. **Can QuoteLineItem be modified after the Quote is created?**
   - a) No, Quote is immutable once created
   - b) Yes, quantities can be adjusted
   - c) Yes, items can be added/removed
   - d) Other: ___

3. **Is QuoteLineItem ever referenced from other entities (Order, Shipment)?**
   - a) No, Order creates its own line items
   - b) Yes, OrderLineItem references QuoteLineItem
   - c) Yes, Shipment tracks which QuoteLineItems were shipped
   - d) Other: ___

### Address

4. **Are two Addresses equal if street, city, zip, country match?**
   - a) Yes, same data = same address
   - b) No, each address has its own identity (e.g., for history tracking)
   - c) Depends on context
   - d) Other: ___

5. **Can a Customer have multiple Addresses simultaneously?**
   - a) No, one address per customer
   - b) Yes, billing and shipping addresses
   - c) Yes, unlimited addresses with labels
   - d) Other: ___
```

#### Step 2.3: Process Answers and Classify

Based on user answers, apply decision logic:

```markdown
### QuoteLineItem Classification

**Answers:**
- EVO-1 (Independent tracking): Yes - shipping tracks (answer 1c)
- EVO-3 (Mutability): Yes - can be modified (answer 2b/c)
- EVO-4 (External references): Yes - Shipment references (answer 3c)

**Decision:** ENTITY
**Rationale:** QuoteLineItem requires independent identity because:
1. Shipping/fulfillment tracks individual line items
2. Can be modified within Quote lifecycle
3. Referenced from Shipment aggregate

---

### Address Classification

**Answers:**
- EVO-5 (Value equality): Yes - same data = same address (answer 4a)
- External context: Multiple per customer (answer 5b/c)

**Decision:** VALUE OBJECT
**Rationale:** Address is defined by its attributes, not identity:
1. Two addresses with same street/city/zip are equal
2. Multiple addresses are distinguished by type (billing/shipping), not ID
```

### Phase 3: Generate Domain Model

#### Step 3.1: Structure

```markdown
---
status: draft
derived-from:
  - "{user-stories-file}"
  - "{vocabulary-file}"
derived-at: "{ISO timestamp}"
derived-by: "loom-derive-domain skill v1.0 (Structured Interview)"
loom-version: "3.0.0"
structured-interview:
  concepts-classified: N
  from-user-answers: X
  from-input: Y
  high-confidence: Z
---

# Domain Model – {Domain Name}

## Overview

{Brief description of the domain and its key concepts}

## Entities

### ENT-{Name} – {Entity Name}

**Classification:** Entity
**Rationale:** {From Structured Interview}

**Decision Points Resolved:**
- EVO-1: {answer} (User answer / From input)
- EVO-4: {answer} (User answer)

**Attributes:**
| Attribute | Type | Required | Description |
|-----------|------|----------|-------------|
| id | UUID | Yes | Unique identifier |
| ... | ... | ... | ... |

**Invariants:**
- INV-{N}: {invariant description}

**Lifecycle:**
- Created: {when}
- Modified: {when/how}
- Deleted: {when/if}

**Relationships:**
- Contains: {value objects}
- References: {other entities by ID}

**Traceability:**
- User Story: {story}#us-{id}

---

## Value Objects

### VO-{Name} – {Value Object Name}

**Classification:** Value Object
**Rationale:** {From Structured Interview}

**Decision Points Resolved:**
- EVO-5: Yes - equality by attributes (User answer)

**Attributes:**
| Attribute | Type | Required | Description |
|-----------|------|----------|-------------|
| ... | ... | ... | ... |

**Equality:**
Two {Name} are equal if: {equality definition}

**Immutability:**
{Describe if/how immutable}

---

## Aggregates

### AGG-{Name} – {Aggregate Name}

**Root:** ENT-{RootEntity}

**Contains:**
- ENT-{Entity} (if entity within aggregate)
- VO-{ValueObject}

**Boundary Rationale:**
{From Structured Interview - AGG decisions}

**Invariants:**
- {Cross-entity invariants enforced by this aggregate}

**Decision Points Resolved:**
- AGG-1: {answer}
- AGG-2: {answer}
```

### Phase 4: Validation and Approval

Present domain model with interview summary:

```markdown
## Domain Model Derivation Results

### Structured Interview Summary

| Concept | Classification | Key Decision Point | Answer | Source |
|---------|---------------|-------------------|--------|--------|
| Quote | Entity | EVO-1 (independent identity) | Yes, has QuoteId | Input |
| QuoteLineItem | Entity | EVO-4 (external references) | Yes, from Shipment | User |
| Customer | Entity | EVO-2 (lifecycle) | Yes, independent | Input |
| Address | Value Object | EVO-5 (value equality) | Yes, same data = same | User |
| Money | Value Object | EVO-5 (value equality) | Yes, amount+currency | Input |

### Classification Summary

| Type | Count | From Input | From User |
|------|-------|------------|-----------|
| Entities | 3 | 2 | 1 |
| Value Objects | 2 | 1 | 1 |
| Aggregates | 2 | - | - |

---

{Full domain-model.md content}

---

Would you like me to:
1. **Write file** to {output-dir}/domain-model.md
2. **Modify** a classification
3. **Ask more questions** about a specific concept
4. **Cancel** derivation
```

## Example: QuoteLineItem Decision

### Input (from User Story)

```markdown
### US-QUOTE-003 – Customer accepts a quote online
...
- An order is created automatically from the accepted quote.
```

### Structured Interview

```
AI: "I see QuoteLineItem in the domain. I need to classify it as Entity or Value Object.

     **Does QuoteLineItem need to be tracked independently outside its Quote?**
     a) No, it only exists as part of a Quote
     b) Yes, for reporting/analytics
     c) Yes, shipping tracks individual line items
     d) Other: ___"

User: "c) Yes, shipping tracks individual line items"

AI: "**Is QuoteLineItem ever referenced from other aggregates?**
     a) No, Order creates its own line items
     b) Yes, OrderLineItem references QuoteLineItem
     c) Yes, Shipment tracks which QuoteLineItems were shipped"

User: "c) Yes, Shipment tracks which QuoteLineItems were shipped"

AI: "**Can QuoteLineItem quantities be modified after Quote creation?**
     a) No, Quote is immutable
     b) Yes, until Quote is sent
     c) Yes, anytime before acceptance"

User: "b) Yes, until Quote is sent"
```

### Classification Result

```markdown
### ENT-QuoteLineItem – Quote Line Item

**Classification:** Entity
**Rationale:** Requires independent identity for cross-aggregate tracking

**Decision Points Resolved:**
- EVO-1: Yes - shipping tracks independently (User answer)
- EVO-3: Yes - mutable until Quote sent (User answer)
- EVO-4: Yes - referenced by Shipment (User answer)

**Attributes:**
| Attribute | Type | Required | Description |
|-----------|------|----------|-------------|
| id | UUID | Yes | Unique line item identifier |
| quoteId | UUID | Yes | Parent quote reference |
| productId | UUID | Yes | Product being quoted |
| quantity | Integer | Yes | Number of units |
| unitPrice | Money | Yes | Price per unit |
| lineTotal | Money | Yes | Calculated total |

**Invariants:**
- INV-QLI-1: quantity MUST be > 0
- INV-QLI-2: lineTotal MUST equal quantity * unitPrice
- INV-QLI-3: Cannot modify if Quote.status in (Sent, Accepted, Rejected)

**Lifecycle:**
- Created: When added to Quote
- Modified: Until Quote is sent
- Deleted: When removed from Quote (if Quote is Draft)

**Traceability:**
- User Story: user-stories.md#us-quote-003
```

---

## Contrast: Without Structured Interview

```markdown
### QuoteLineItem (IMPLICIT DECISION)

**Classification:** Value Object
**Rationale:** Part of Quote aggregate

(No decision points documented, no user consultation,
could be wrong if shipping actually needs to track line items!)
```

**The difference:** With Structured Interview, we KNOW it's an Entity because the user confirmed shipping tracks individual items. Without it, we might guess wrong.

---

Now read the input files and begin the Structured Interview process for domain modeling.
