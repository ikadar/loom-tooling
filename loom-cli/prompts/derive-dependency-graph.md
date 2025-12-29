# Dependency Graph Derivation Prompt

You are an expert system architect. Generate Dependency Graphs from Service Boundaries.

CRITICAL OUTPUT REQUIREMENTS:
1. Wrap response in ```json code blocks
2. NO explanations - JSON only
3. ALL string values must be SHORT (max 60 chars)
4. NO line breaks within string values

## Your Task

From Service Boundaries, derive a Dependency Graph showing:
1. **Components** - All services and external systems
2. **Dependencies** - Directed relationships with types
3. **Dependency types** - sync, async, data, external

## Output Format

```json
{
  "components": [
    {
      "id": "Order Service",
      "type": "domain_service",
      "description": "Manages order lifecycle"
    },
    {
      "id": "Payment Gateway",
      "type": "external",
      "description": "Third-party payment processor"
    }
  ],
  "dependencies": [
    {
      "from": "Order Service",
      "to": "Inventory Service",
      "type": "sync",
      "description": "Reserve stock on order"
    },
    {
      "from": "Order Service",
      "to": "Notification Service",
      "type": "async",
      "description": "OrderCreated event"
    },
    {
      "from": "Order Service",
      "to": "Payment Gateway",
      "type": "external",
      "description": "Process payments"
    }
  ],
  "summary": {
    "total_components": 8,
    "total_dependencies": 15,
    "by_type": {"sync": 6, "async": 7, "external": 2}
  }
}
```

REMINDER: JSON only. All strings under 60 chars.

INPUT:

