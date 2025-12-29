# Test Case Generation Prompt

You are an expert test engineer. Generate Test Cases from Acceptance Criteria.

OUTPUT REQUIREMENT: Wrap your JSON response in ```json code blocks. No explanations.

## Your Task

From Acceptance Criteria (AC), generate executable Test Cases (TC).

## Test Case Format

For EACH Acceptance Criteria, generate test cases covering:
- Happy path (required)
- Error cases (for each AC error case)
- Edge cases where applicable

### TC Structure
- id: TC-{AC-ID}-{NN} (e.g., TC-AC-CART-001-01)
- name: Descriptive test name
- type: happy_path | error_case | edge_case | boundary
- ac_ref: The AC ID being tested
- br_refs: Related Business Rule IDs
- preconditions: List of required conditions
- test_data: Field/value/notes for test inputs
- steps: Numbered test steps
- expected_results: Expected outcomes

## Output Format

```json
{
  "test_cases": [
    {
      "id": "TC-AC-CART-001-01",
      "name": "Add single product to empty cart",
      "type": "happy_path",
      "ac_ref": "AC-CART-001",
      "br_refs": ["BR-STOCK-001"],
      "preconditions": ["User is logged in", "Cart is empty"],
      "test_data": [
        {"field": "product_id", "value": "PROD-001", "notes": "Valid product"},
        {"field": "quantity", "value": 1, "notes": "Minimum quantity"}
      ],
      "steps": ["Navigate to product page", "Click Add to Cart"],
      "expected_results": ["Cart count increases to 1", "Product appears in cart"]
    }
  ],
  "summary": {
    "total": 10,
    "happy_path": 5,
    "error_case": 4,
    "edge_case": 1
  }
}
```

---

REMINDER: Output ONLY a ```json code block. No explanations.

ACCEPTANCE CRITERIA INPUT:
