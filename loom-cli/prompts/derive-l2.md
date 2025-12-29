# L1 to L2 Derivation Prompt

You are an expert test engineer and technical architect. Generate L2 documents from L1 inputs.

OUTPUT REQUIREMENT: Wrap your JSON response in ```json code blocks. Do not include any explanations before or after the code block.

## Your Task

From Acceptance Criteria (AC) and Business Rules (BR), generate:
1. **Test Cases (TC)** - executable test specifications
2. **Technical Specifications** - implementation guidance

## Input Format

You will receive:
- Acceptance Criteria (AC-XXX-NNN format)
- Business Rules (BR-XXX-NNN format)
- Domain Model context

## Test Case Generation

For EACH Acceptance Criteria, generate test cases:

### TC Format
```markdown
### TC-{AC-ID}-{NN} – {Test Name}

**Type:** {happy_path|error_case|edge_case|boundary}

**Preconditions:**
- {precondition 1}
- {precondition 2}

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| {field} | {value} | {why this value} |

**Steps:**
1. {action}
2. {action}
3. {action}

**Expected Result:**
- {assertion 1}
- {assertion 2}

**Traceability:**
- AC: {AC-ID}
- BR: {BR-IDs if applicable}
```

### Test Case Rules

1. **Coverage Requirements:**
   - Every AC MUST have at least 1 happy path TC
   - Every AC error case MUST have a TC
   - Every BR violation MUST have a negative TC

2. **Test Data:**
   - Use realistic, specific values
   - Include boundary values where applicable
   - Reference decision values from interview

3. **Naming Convention:**
   - TC-{AC-ID}-01, TC-{AC-ID}-02, etc.
   - Example: TC-AC-CART-001-01, TC-AC-CART-001-02

## Technical Specification Generation

For EACH Business Rule, generate technical specs:

### Tech Spec Format
```markdown
### TS-{BR-ID} – {Spec Name}

**Rule:** {BR statement}

**Implementation Approach:**
{How to implement this rule}

**Validation Points:**
- {where validation should occur}
- {what to validate}

**Data Requirements:**
| Field | Type | Constraints | Source |
|-------|------|-------------|--------|
| {field} | {type} | {constraints} | {BR/AC reference} |

**Error Handling:**
| Condition | Error Code | Message | HTTP Status |
|-----------|------------|---------|-------------|
| {condition} | {code} | {message} | {status} |

**Traceability:**
- BR: {BR-ID}
- Related ACs: {AC-IDs}
```

## Output Format

CRITICAL: Return ONLY raw JSON. No markdown code blocks. No explanations. No text before or after the JSON object. The response must start with { and end with }.

JSON structure:
```json
{
  "summary": {
    "test_cases_generated": N,
    "tech_specs_generated": N,
    "coverage": {
      "acs_covered": N,
      "brs_covered": N,
      "happy_path_tests": N,
      "error_tests": N,
      "edge_case_tests": N
    }
  },
  "test_cases": [
    {
      "id": "TC-AC-CART-001-01",
      "name": "Add single product to empty cart",
      "type": "happy_path",
      "ac_ref": "AC-CART-001",
      "br_refs": ["BR-STOCK-001"],
      "preconditions": ["User is logged in", "Cart is empty", "Product is in stock"],
      "test_data": [
        {"field": "product_id", "value": "PROD-001", "notes": "Valid in-stock product"},
        {"field": "quantity", "value": 1, "notes": "Minimum valid quantity"}
      ],
      "steps": [
        "Navigate to product page for PROD-001",
        "Click 'Add to Cart' button",
        "Verify cart notification appears"
      ],
      "expected_results": [
        "Cart count increases to 1",
        "Toast notification shows 'Added to cart'",
        "Cart contains PROD-001 with quantity 1"
      ]
    }
  ],
  "tech_specs": [
    {
      "id": "TS-BR-STOCK-001",
      "name": "Stock validation for cart addition",
      "br_ref": "BR-STOCK-001",
      "rule": "Products must have available stock to be added to cart",
      "implementation": "Check inventory.available_quantity > 0 before allowing add to cart",
      "validation_points": [
        "Add to cart API endpoint",
        "Cart UI before enabling button"
      ],
      "data_requirements": [
        {"field": "available_quantity", "type": "integer", "constraints": ">= 0", "source": "BR-STOCK-001"}
      ],
      "error_handling": [
        {"condition": "available_quantity == 0", "error_code": "OUT_OF_STOCK", "message": "This product is out of stock", "http_status": 400}
      ],
      "related_acs": ["AC-CART-001", "AC-CART-002"]
    }
  ]
}
```

## Quality Checklist

Before output, verify:
- [ ] Every AC has at least one test case
- [ ] Every AC error case has a test case
- [ ] Every BR has a technical specification
- [ ] All test cases have specific test data
- [ ] All traceability links are correct
- [ ] IDs follow naming conventions

---

REMINDER: Output ONLY a ```json code block with the full JSON. No explanations.

L1 INPUT (AC + BR):
