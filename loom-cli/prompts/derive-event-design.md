# Event & Message Design Derivation Prompt

You are an expert event-driven architect. Generate Event & Message Design from L1/L2 documents.

CRITICAL OUTPUT REQUIREMENTS:
1. Wrap response in ```json code blocks
2. NO explanations - JSON only
3. ALL string values must be SHORT (max 80 chars)
4. NO line breaks within string values

## Your Task

From Domain Model and Sequence Design, derive Event & Message definitions:
1. **Domain Events** - Facts about state changes
2. **Commands** - Intent to perform actions
3. **Integration Events** - Cross-service notifications

## Output Format

```json
{
  "domain_events": [
    {
      "id": "EVT-ORDER-001",
      "name": "OrderCreated",
      "purpose": "Signals new order was placed",
      "trigger": "Customer places order from cart",
      "aggregate": "Order",
      "payload": [
        {"field": "orderId", "type": "UUID"},
        {"field": "customerId", "type": "UUID"},
        {"field": "totalAmount", "type": "Money"},
        {"field": "lineItems", "type": "array"},
        {"field": "createdAt", "type": "DateTime"}
      ],
      "invariants_reflected": ["Order has items", "Total calculated"],
      "consumers": ["Inventory Service", "Notification Service"],
      "version": "1.0"
    }
  ],
  "commands": [
    {
      "id": "CMD-ORDER-001",
      "name": "CreateOrder",
      "intent": "Place new order from cart contents",
      "required_data": [
        {"field": "customerId", "type": "UUID"},
        {"field": "shippingAddress", "type": "Address"},
        {"field": "paymentMethod", "type": "PaymentMethod"}
      ],
      "expected_outcome": "Order created with pending status",
      "failure_conditions": ["Cart empty", "Item out of stock"]
    }
  ],
  "integration_events": [
    {
      "id": "INT-ORDER-001",
      "name": "OrderPlacedNotification",
      "purpose": "Notify external systems of new order",
      "source": "Order Service",
      "consumers": ["Email Service", "Analytics"],
      "payload": ["orderId", "customerId", "totalAmount"]
    }
  ],
  "summary": {
    "domain_events": 15,
    "commands": 10,
    "integration_events": 5
  }
}
```

REMINDER: JSON only. All strings under 80 chars.

INPUT:

