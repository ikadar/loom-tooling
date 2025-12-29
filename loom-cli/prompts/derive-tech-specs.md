# Technical Specification Generation Prompt

You are an expert technical architect. Generate Technical Specifications from Business Rules.

OUTPUT REQUIREMENT: Wrap your JSON response in ```json code blocks. No explanations.

## Your Task

From Business Rules (BR), generate Technical Specifications for implementation.

## Tech Spec Structure

For EACH Business Rule, generate a specification covering:
- Implementation approach
- Validation points
- Data requirements
- Error handling

### TS Fields
- id: TS-{BR-ID} (e.g., TS-BR-STOCK-001)
- name: Descriptive name
- br_ref: The BR ID being specified
- rule: The business rule statement
- implementation: How to implement
- validation_points: Where to validate
- data_requirements: Field/type/constraints/source
- error_handling: Condition/error_code/message/http_status
- related_acs: Related AC IDs

## Output Format

```json
{
  "tech_specs": [
    {
      "id": "TS-BR-STOCK-001",
      "name": "Stock validation for cart",
      "br_ref": "BR-STOCK-001",
      "rule": "Products must have stock to be added",
      "implementation": "Check inventory before add to cart",
      "validation_points": ["Add to cart API", "Cart UI"],
      "data_requirements": [
        {"field": "quantity", "type": "integer", "constraints": ">= 0", "source": "BR-STOCK-001"}
      ],
      "error_handling": [
        {"condition": "quantity == 0", "error_code": "OUT_OF_STOCK", "message": "Out of stock", "http_status": 400}
      ],
      "related_acs": ["AC-CART-001"]
    }
  ],
  "summary": {
    "total": 5
  }
}
```

---

REMINDER: Output ONLY a ```json code block. No explanations.

BUSINESS RULES INPUT:
