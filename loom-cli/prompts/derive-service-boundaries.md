# Service Boundaries Derivation Prompt

You are an expert system architect. Generate Service Boundaries from L1/L2 documents.

CRITICAL OUTPUT REQUIREMENTS:
1. Wrap response in ```json code blocks
2. NO explanations - JSON only
3. ALL string values must be SHORT (max 80 chars)
4. NO line breaks within string values

## Your Task

From Bounded Context Map and Aggregate Design, derive Service Boundaries that define:
1. **Purpose** - Core responsibility
2. **Capabilities** - What the service can do
3. **Inputs/Outputs** - Commands, events, responses
4. **Ownership** - Aggregates owned
5. **Dependencies** - External services needed

## Output Format

```json
{
  "services": [
    {
      "id": "SVC-ORDER",
      "name": "Order Service",
      "purpose": "Manages complete order lifecycle",
      "capabilities": [
        "Create orders from cart",
        "Track order status",
        "Cancel orders"
      ],
      "inputs": [
        {"type": "command", "name": "CreateOrder"},
        {"type": "command", "name": "CancelOrder"}
      ],
      "outputs": [
        {"type": "event", "name": "OrderCreated"},
        {"type": "event", "name": "OrderCancelled"}
      ],
      "owned_aggregates": ["Order"],
      "dependencies": [
        {"service": "Inventory Service", "type": "sync", "reason": "Reserve stock"},
        {"service": "Cart Service", "type": "sync", "reason": "Get cart items"}
      ],
      "api_base": "/api/v1/orders",
      "separation_reason": "Order lifecycle is distinct bounded context"
    }
  ],
  "summary": {
    "total_services": 5,
    "total_dependencies": 12
  }
}
```

REMINDER: JSON only. All strings under 80 chars.

INPUT:

