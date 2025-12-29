# Domain Model Derivation Prompt

You are an expert domain-driven design architect. Generate a comprehensive Domain Model from L0 inputs.

OUTPUT REQUIREMENT: Wrap your JSON response in ```json code blocks. No explanations.

## Your Task

From User Stories and Domain Vocabulary, derive:
1. **Entities** - Core domain objects with identity
2. **Value Objects** - Immutable objects without identity
3. **Aggregate Roots** - Consistency boundaries
4. **Relationships** - How entities relate to each other

## Entity Structure

For each entity, include:
- id: Unique identifier (ENT-{DOMAIN}-{NNN})
- name: Entity name (PascalCase)
- type: "aggregate_root" | "entity" | "value_object"
- purpose: Why this entity exists
- attributes: List of fields with type and constraints
- invariants: Business rules that must always hold
- operations: Methods/behaviors
- events: Domain events emitted
- relationships: Links to other entities

## Output Format

```json
{
  "domain_model": {
    "name": "E-Commerce Domain",
    "description": "Domain model for online shopping platform",
    "bounded_contexts": ["Orders", "Catalog", "Customers"]
  },
  "entities": [
    {
      "id": "ENT-ORDER-001",
      "name": "Order",
      "type": "aggregate_root",
      "purpose": "Represents a customer purchase transaction",
      "attributes": [
        {
          "name": "orderId",
          "type": "OrderId (UUID)",
          "constraints": "Required, immutable after creation"
        },
        {
          "name": "status",
          "type": "OrderStatus (enum)",
          "constraints": "Required, valid transitions only"
        },
        {
          "name": "totalAmount",
          "type": "Money",
          "constraints": ">= 0"
        }
      ],
      "invariants": [
        "Order total must equal sum of line items",
        "Cannot modify after status is 'Shipped'",
        "Must have at least one line item"
      ],
      "operations": [
        {
          "name": "addItem",
          "signature": "addItem(product: Product, quantity: number): void",
          "preconditions": ["Status is 'Draft'", "Product is available"],
          "postconditions": ["Line item added", "Total recalculated"]
        },
        {
          "name": "submit",
          "signature": "submit(): void",
          "preconditions": ["Has line items", "Status is 'Draft'"],
          "postconditions": ["Status changes to 'Submitted'", "OrderSubmitted event emitted"]
        }
      ],
      "events": [
        {
          "name": "OrderSubmitted",
          "trigger": "submit() called",
          "payload": ["orderId", "customerId", "totalAmount", "timestamp"]
        }
      ],
      "relationships": [
        {
          "target": "ENT-CUSTOMER-001",
          "type": "belongs_to",
          "cardinality": "many-to-one"
        },
        {
          "target": "ENT-LINEITEM-001",
          "type": "contains",
          "cardinality": "one-to-many"
        }
      ]
    }
  ],
  "value_objects": [
    {
      "id": "VO-MONEY-001",
      "name": "Money",
      "purpose": "Represents monetary value with currency",
      "attributes": [
        {"name": "amount", "type": "decimal", "constraints": ">= 0"},
        {"name": "currency", "type": "string", "constraints": "ISO 4217 code"}
      ],
      "operations": ["add", "subtract", "multiply", "equals"]
    }
  ],
  "summary": {
    "aggregate_roots": 3,
    "entities": 5,
    "value_objects": 4,
    "total_operations": 15,
    "total_events": 8
  }
}
```

## Quality Checklist

Before output, verify:
- [ ] Every user story maps to at least one entity
- [ ] Every entity has clear invariants
- [ ] Aggregate boundaries are explicit
- [ ] All relationships have cardinality
- [ ] Value objects are truly immutable
- [ ] Events follow past-tense naming (OrderCreated, not CreateOrder)

---

REMINDER: Output ONLY a ```json code block. No explanations.

L0 INPUT (User Stories + Domain Vocabulary):
