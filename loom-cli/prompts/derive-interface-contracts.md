# Interface Contracts Derivation Prompt

You are an expert API designer. Generate Interface Contracts from L1 documents.

OUTPUT REQUIREMENT: Wrap your JSON response in ```json code blocks. No explanations.

## Your Task

From Domain Model, Business Rules, and Acceptance Criteria, derive Interface Contracts that define:
1. **Service APIs** - HTTP endpoints with full request/response schemas
2. **Operations** - Each operation with inputs, outputs, errors, pre/postconditions
3. **Shared Types** - Reusable data types and value objects

## Interface Contract Structure

For each service interface, include:
- id: Unique identifier (IC-{SERVICE}-{NNN})
- serviceName: Name of the service
- purpose: Why this interface exists
- operations: List of API operations
- sharedTypes: Reusable type definitions
- securityRequirements: Authentication/authorization rules

## Output Format

```json
{
  "interface_contracts": [
    {
      "id": "IC-ORDER-001",
      "serviceName": "Order Service",
      "purpose": "Manages order lifecycle from creation to delivery",
      "baseUrl": "/api/v1/orders",
      "operations": [
        {
          "id": "OP-ORDER-001",
          "name": "createOrder",
          "method": "POST",
          "path": "/",
          "description": "Create a new order from cart contents",
          "inputSchema": {
            "customerId": {"type": "UUID", "required": true},
            "shippingAddress": {"type": "Address", "required": true},
            "paymentMethod": {"type": "PaymentMethod", "required": true}
          },
          "outputSchema": {
            "orderId": {"type": "UUID"},
            "status": {"type": "OrderStatus"},
            "totalAmount": {"type": "Money"},
            "createdAt": {"type": "DateTime"}
          },
          "errors": [
            {"code": "CART_EMPTY", "httpStatus": 400, "message": "Cart has no items"},
            {"code": "INSUFFICIENT_STOCK", "httpStatus": 409, "message": "One or more items out of stock"},
            {"code": "CUSTOMER_NOT_FOUND", "httpStatus": 404, "message": "Customer does not exist"}
          ],
          "preconditions": ["Customer must be registered", "Cart must have items"],
          "postconditions": ["Order created with status 'pending'", "Cart cleared", "Inventory reserved"],
          "relatedACs": ["AC-ORDER-001"],
          "relatedBRs": ["BR-ORDER-001", "BR-ORDER-002"]
        }
      ],
      "events": [
        {
          "name": "OrderCreated",
          "description": "Emitted when a new order is created",
          "payload": ["orderId", "customerId", "totalAmount", "lineItems"]
        }
      ],
      "securityRequirements": {
        "authentication": "Bearer JWT",
        "authorization": "Customer must own the cart"
      }
    }
  ],
  "shared_types": [
    {
      "name": "Money",
      "fields": [
        {"name": "amount", "type": "decimal", "constraints": ">= 0, precision 2"},
        {"name": "currency", "type": "string", "constraints": "ISO 4217, default USD"}
      ]
    },
    {
      "name": "Address",
      "fields": [
        {"name": "street", "type": "string", "constraints": "1-200 chars"},
        {"name": "city", "type": "string", "constraints": "1-100 chars"},
        {"name": "state", "type": "string", "constraints": "required"},
        {"name": "postalCode", "type": "string", "constraints": "valid format"},
        {"name": "country", "type": "string", "constraints": "ISO 3166-1"}
      ]
    }
  ],
  "summary": {
    "total_contracts": 5,
    "total_operations": 15,
    "total_events": 8,
    "total_shared_types": 6
  }
}
```

## Quality Checklist

Before output, verify:
- [ ] Every AC maps to at least one operation
- [ ] Every BR is referenced in validation or error handling
- [ ] All operations have clear input/output schemas
- [ ] Error codes are consistent and meaningful
- [ ] Pre/postconditions are explicit

---

REMINDER: Output ONLY a ```json code block. No explanations.

L1 INPUT (Domain Model + Business Rules + Acceptance Criteria):

