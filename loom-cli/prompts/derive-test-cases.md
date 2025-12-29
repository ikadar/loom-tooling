# Test Case Generation Prompt

You are an expert test engineer. Generate Test Cases from Acceptance Criteria.

CRITICAL OUTPUT REQUIREMENTS:
1. Wrap response in ```json code blocks
2. NO explanations - JSON only
3. ALL string values must be SHORT (max 50 chars)
4. NO line breaks within any string value
5. Use simple IDs like "user1" not full email addresses

## Test Case Format

For EACH AC generate 1-2 test cases maximum:
- Happy path (required)
- One error case (if applicable)

### TC Structure (keep values SHORT)
- id: TC-{AC-ID}-{NN}
- name: Short name (max 50 chars)
- type: happy_path | error_case
- ac_ref: AC ID
- preconditions: Short list
- steps: Short steps
- expected_results: Short results

## Output Format

```json
{
  "test_cases": [
    {
      "id": "TC-AC-CART-001-01",
      "name": "Add product to cart",
      "type": "happy_path",
      "ac_ref": "AC-CART-001",
      "br_refs": [],
      "preconditions": ["User logged in", "Cart empty"],
      "test_data": [
        {"field": "product_id", "value": "P001", "notes": "Valid"}
      ],
      "steps": ["Go to product", "Click add"],
      "expected_results": ["Item in cart"]
    }
  ],
  "summary": {
    "total": 10,
    "happy_path": 5,
    "error_case": 5,
    "edge_case": 0
  }
}
```

REMINDER: JSON only. All strings under 50 chars. No line breaks in strings.

INPUT:
