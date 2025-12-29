# Bounded Context Map Derivation Prompt

You are an expert domain-driven design architect. Generate a Bounded Context Map from the Domain Model.

OUTPUT REQUIREMENT: Wrap your JSON response in ```json code blocks. No explanations.

## Your Task

From the Domain Model, derive:
1. **Bounded Contexts** - Logical boundaries grouping related entities
2. **Context Relationships** - How contexts communicate
3. **Integration Patterns** - Shared Kernel, Customer-Supplier, Conformist, etc.

## Context Relationship Types

- **Shared Kernel**: Two contexts share a subset of domain model
- **Customer-Supplier**: Upstream context serves downstream context
- **Conformist**: Downstream conforms to upstream's model
- **Anti-corruption Layer**: Translation layer between contexts
- **Open Host Service**: Published API for multiple consumers
- **Published Language**: Shared interchange format

## Output Format

IMPORTANT: Do NOT include the context_map_diagram field. Output ONLY the bounded_contexts, context_relationships and summary fields.

```json
{
  "bounded_contexts": [
    {
      "id": "BC-ORDER",
      "name": "Order Management",
      "purpose": "Handles order lifecycle from cart to fulfillment",
      "core_entities": ["Order", "LineItem", "OrderStatus"],
      "aggregates": ["Order"],
      "capabilities": [
        "Create and manage orders",
        "Track order status",
        "Calculate totals and taxes"
      ],
      "ubiquitous_language": {
        "Order": "A customer's purchase request",
        "LineItem": "A single product entry in an order",
        "Fulfillment": "Process of preparing and shipping an order"
      }
    }
  ],
  "context_relationships": [
    {
      "upstream": "BC-CATALOG",
      "downstream": "BC-ORDER",
      "relationship_type": "customer_supplier",
      "description": "Order context consumes product info from Catalog",
      "integration_pattern": "Open Host Service",
      "shared_data": ["productId", "productName", "price"]
    }
  ],
  "summary": {
    "total_contexts": 4,
    "total_relationships": 5,
    "integration_patterns_used": ["Open Host Service", "Domain Events", "Shared Kernel"]
  }
}
```

## Context Identification Rules

1. **High Cohesion**: Entities within a context are tightly related
2. **Loose Coupling**: Contexts communicate through well-defined interfaces
3. **Single Team Ownership**: One team should own one context
4. **Linguistic Boundary**: Terms have consistent meaning within context
5. **Transaction Boundary**: Strong consistency within, eventual consistency between

## Quality Checklist

Before output, verify:
- [ ] Every aggregate belongs to exactly one context
- [ ] No circular dependencies between contexts
- [ ] Relationship types are explicit
- [ ] Ubiquitous language is defined per context

---

REMINDER: Output ONLY a ```json code block. No explanations.

DOMAIN MODEL INPUT:
