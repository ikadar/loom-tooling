# Aggregate Design Derivation Prompt

You are an expert DDD architect. Generate detailed Aggregate Designs from L1 documents.

OUTPUT REQUIREMENT: Wrap your JSON response in ```json code blocks. No explanations.

## Your Task

From Domain Model and Business Rules, derive detailed Aggregate Designs that specify:
1. **Aggregate boundaries** - What belongs inside each aggregate
2. **Invariants** - Business rules that must always hold
3. **Behaviors** - Commands and state transitions
4. **Events** - Domain events emitted
5. **Repository patterns** - How aggregates are persisted

## Aggregate Structure

For each aggregate, include:
- id: Unique identifier (AGG-{NAME}-{NNN})
- name: Aggregate root name
- purpose: Why this aggregate exists
- invariants: Business rules that must always hold
- entities: Child entities within the boundary
- valueObjects: Immutable value types used
- behaviors: Commands the aggregate handles
- events: Domain events emitted
- repository: Persistence patterns

## Output Format

CRITICAL: "valueObjects" must be a simple string array (e.g., `["Money", "Address"]`), NOT objects with fields.

```json
{
  "aggregates": [
    {
      "id": "AGG-ORDER-001",
      "name": "Order",
      "purpose": "Manages the complete order lifecycle from placement through delivery",
      "invariants": [
        {
          "id": "INV-ORDER-001",
          "rule": "Order must have at least one line item",
          "enforcement": "Validated on creation and after any item removal"
        },
        {
          "id": "INV-ORDER-002",
          "rule": "Total amount must equal sum of line items plus shipping",
          "enforcement": "Recalculated after any item change"
        },
        {
          "id": "INV-ORDER-003",
          "rule": "Cannot modify order after status is 'shipped'",
          "enforcement": "Guard clause in all modification methods"
        },
        {
          "id": "INV-ORDER-004",
          "rule": "Status transitions must follow valid state machine",
          "enforcement": "State pattern with allowed transitions"
        }
      ],
      "root": {
        "entity": "Order",
        "identity": "OrderId (UUID)",
        "attributes": [
          {"name": "status", "type": "OrderStatus", "mutable": true},
          {"name": "totalAmount", "type": "Money", "mutable": true},
          {"name": "shippingCost", "type": "Money", "mutable": false},
          {"name": "createdAt", "type": "DateTime", "mutable": false}
        ]
      },
      "entities": [
        {
          "name": "OrderLineItem",
          "identity": "LineItemId",
          "purpose": "Immutable snapshot of product at order time",
          "attributes": [
            {"name": "productId", "type": "ProductId"},
            {"name": "productName", "type": "string"},
            {"name": "unitPrice", "type": "Money"},
            {"name": "quantity", "type": "Quantity"},
            {"name": "subtotal", "type": "Money"}
          ]
        }
      ],
      "valueObjects": ["Money", "Address", "OrderStatus", "Quantity"],
      "behaviors": [
        {
          "name": "confirm",
          "command": "ConfirmOrder",
          "preconditions": ["status == 'pending'"],
          "postconditions": ["status == 'confirmed'"],
          "emits": "OrderConfirmed"
        },
        {
          "name": "ship",
          "command": "ShipOrder",
          "preconditions": ["status == 'confirmed'"],
          "postconditions": ["status == 'shipped'"],
          "emits": "OrderShipped"
        },
        {
          "name": "deliver",
          "command": "DeliverOrder",
          "preconditions": ["status == 'shipped'"],
          "postconditions": ["status == 'delivered'"],
          "emits": "OrderDelivered"
        },
        {
          "name": "cancel",
          "command": "CancelOrder",
          "preconditions": ["status in ['pending', 'confirmed']"],
          "postconditions": ["status == 'cancelled'"],
          "emits": "OrderCancelled"
        }
      ],
      "events": [
        {"name": "OrderCreated", "payload": ["orderId", "customerId", "lineItems", "totalAmount"]},
        {"name": "OrderConfirmed", "payload": ["orderId", "confirmedAt"]},
        {"name": "OrderShipped", "payload": ["orderId", "shippedAt", "trackingNumber"]},
        {"name": "OrderDelivered", "payload": ["orderId", "deliveredAt"]},
        {"name": "OrderCancelled", "payload": ["orderId", "cancelledAt", "reason"]}
      ],
      "repository": {
        "name": "OrderRepository",
        "methods": [
          {"name": "findById", "params": "OrderId", "returns": "Order?"},
          {"name": "findByCustomerId", "params": "CustomerId", "returns": "List<Order>"},
          {"name": "save", "params": "Order", "returns": "void"}
        ],
        "loadStrategy": "Load complete aggregate with all line items",
        "concurrency": "Optimistic locking via version field"
      },
      "externalReferences": [
        {"aggregate": "Customer", "via": "customerId", "type": "reference"},
        {"aggregate": "Product", "via": "productId in line items", "type": "snapshot"}
      ]
    }
  ],
  "summary": {
    "total_aggregates": 5,
    "total_invariants": 18,
    "total_behaviors": 25,
    "total_events": 20
  }
}
```

## Design Principles

Apply these DDD principles:
1. **Small aggregates** - Prefer multiple small aggregates over one large one
2. **Reference by ID** - Reference other aggregates by ID only, never by object
3. **Single transaction** - One aggregate per transaction
4. **Strong invariants** - All business rules enforced in the aggregate root

## Quality Checklist

Before output, verify:
- [ ] Each aggregate has clear boundaries
- [ ] All invariants are explicitly stated
- [ ] State transitions have pre/postconditions
- [ ] Events follow past-tense naming
- [ ] External references are by ID only

---

REMINDER: Output ONLY a ```json code block. No explanations.

L1 INPUT (Domain Model + Business Rules):

