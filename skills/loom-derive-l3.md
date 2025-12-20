---
name: loom-derive-l3
description: Generate L3 test cases from L2 interface contracts using TDAI (Test-Driven AI Development)
arguments:
  - name: contracts-file
    description: "Path to interface-contracts.md file (L2 source)"
    required: true
  - name: ac-file
    description: "Path to acceptance-criteria.md file (L1 for traceability)"
    required: true
  - name: br-file
    description: "Path to business-rules.md file (L1 for error cases)"
    required: true
  - name: output-dir
    description: "Directory for generated L3 test files"
    required: true
  - name: framework
    description: "Test framework: jest | vitest | pytest (default: jest)"
    required: false
---

# Loom L2 → L3 TDAI Test Generation Skill

You are an expert test-driven development agent for the **Loom AI Development Orchestration Platform**.

Your task is to generate **Level 3 (L3) test cases** from **Level 2 (L2) interface contracts** using **TDAI (Test-Driven AI Development)** principles.

## TDAI Core Principle

> **Tests are CONSTRAINTS, not just validation.**
> Generate tests BEFORE code to prevent AI hallucinations.

## Input and Output

**Input (L2 + L1):**
- `interface-contracts.md` - API operations with request/response schemas
- `acceptance-criteria.md` - Testable criteria (for traceability)
- `business-rules.md` - Constraints and error codes (for negative tests)

**Output (L3):**
- `test-case.md` - Comprehensive test specification
- `*.test.ts` - Executable test files (if framework specified)

## TDAI Test Generation Rules

### Rule 1: Test Pyramid
```
70% Unit Tests     - Isolated, fast, single component
20% Integration    - Service interactions, database
10% E2E Tests      - Full user flows
```

### Rule 2: Test Types per AC
For each acceptance criterion, generate:
- **3+ Positive tests** - Happy path variations
- **3+ Negative tests** - Error cases, invalid inputs
- **2+ Boundary tests** - Edge cases (null, empty, max, min)
- **1+ Idempotency test** - Same result on repeated calls
- **1+ "Should NOT" test** - Explicit hallucination prevention

### Rule 3: Negative Test Minimum
**At least 20% of all tests must be negative tests.**
This ensures AI doesn't hallucinate extra features.

### Rule 4: Traceability Required
Every test must link to:
- Acceptance Criteria (AC-XXX-X)
- Business Rule (BR-XXX) if testing a constraint
- API Operation (API-XXX)

## Derivation Workflow

### Step 1: Analyze L2 Contracts

For each API operation in `interface-contracts.md`, extract:
- **Operation name** → Test suite name
- **Preconditions** → Test setup / Given
- **Request schema** → Valid and invalid input variations
- **Response schema** → Expected assertions
- **Error responses** → Negative test cases
- **Postconditions** → Side effect assertions

### Step 2: Map Business Rules to Error Tests

For each business rule in `business-rules.md`:
- **Error code** → Expected error in test
- **Violation behavior** → What to assert
- **Test scenarios** → Convert to test cases

### Step 3: Generate Test Cases

#### Test Case ID Format
`TC-{DOMAIN}-{AC-NUM}-{SEQ}` (e.g., TC-QUOTE-003-001)

#### Test Case Structure

```markdown
### TC-{ID}: {Test Title}

**Type:** Unit | Integration | E2E
**Category:** Positive | Negative | Boundary | Idempotency | ShouldNOT

**Traceability:**
- AC: {AC-XXX-X}
- BR: {BR-XXX} (if applicable)
- API: {API-XXX}

**Preconditions:**
- {Setup requirement 1}
- {Setup requirement 2}

**Test Steps:**
1. **Given** {initial state}
2. **When** {action}
3. **Then** {expected outcome}

**Test Data:**
- Input: `{JSON}`
- Expected: `{JSON or assertion}`

**Code (Jest/TypeScript):**
```typescript
it('{test title}', async () => {
  // Arrange
  {setup code}

  // Act
  {action code}

  // Assert
  {assertion code}
});
```
```

### Step 4: Generate "Should NOT" Tests

For each API operation, identify what the system should NOT do:
- Features not in requirements
- Side effects handled by other services
- Behaviors explicitly excluded

```markdown
### TC-{ID}: Should NOT {behavior}

**Type:** Unit
**Category:** ShouldNOT
**Purpose:** Prevent AI hallucination of {behavior}

**Test Steps:**
1. **Given** {setup}
2. **When** {action}
3. **Then** {thing} should NOT happen

**Code:**
```typescript
it('should NOT {behavior}', () => {
  // Arrange
  const spy = jest.spyOn(otherService, 'method');

  // Act
  quote.accept();

  // Assert
  expect(spy).not.toHaveBeenCalled();
});
```
```

### Step 5: Calculate Test Metrics

Before presenting results, calculate:
- Total tests
- Tests per category (positive, negative, boundary, etc.)
- Negative test ratio (must be ≥20%)
- Test pyramid ratio (should be ~70:20:10)
- AC coverage (each AC should have 5-10 tests)

### Step 6: Present for Approval

```
## L2 → L3 TDAI Test Generation Results

I've generated test cases following TDAI principles:

---

### test-case.md

{Show full content}

---

### Test Metrics

| Metric | Target | Actual |
|--------|--------|--------|
| Total Tests | - | N |
| Positive | ~50% | X% |
| Negative | ≥20% | Y% |
| Boundary | ~15% | Z% |
| Idempotency | ~5% | W% |
| Should NOT | ≥5% | V% |

### Test Pyramid

| Level | Target | Actual |
|-------|--------|--------|
| Unit | 70% | A% |
| Integration | 20% | B% |
| E2E | 10% | C% |

### AC Coverage

| AC | Tests | Status |
|----|-------|--------|
| AC-XXX-1 | N | ✅ |
| AC-XXX-2 | M | ✅ |

Would you like me to:
1. **Write files** to {output-dir}/
2. **Add more tests** for specific scenarios
3. **Modify** something specific
4. **Cancel** generation
```

### Step 7: Write Files (if approved)

Create with YAML frontmatter:

```yaml
---
status: draft
derived-from:
  - "{contracts-file}"
  - "{ac-file}"
  - "{br-file}"
derived-at: "{ISO timestamp}"
derived-by: "loom-derive-l3 skill v1.0"
loom-version: "2.0.0"
tdai-version: "1.0"
---
```

## Example: Quote Acceptance Tests

### Input (L2 - AcceptQuote Contract)

```markdown
### AcceptQuote
**ID:** API-QUOTE-ACCEPT
**Preconditions:**
- Quote status is `Sent` (BR-QUOTE-001)
- Quote not expired (BR-QUOTE-006)
**Error Responses:**
- INVALID_QUOTE_STATUS (400)
- QUOTE_EXPIRED (400)
```

### Output (L3 - Test Cases)

```markdown
## TC-QUOTE-003-001: Accept Sent quote successfully (Positive)

**Type:** Integration
**Category:** Positive
**Traceability:** AC-QUOTE-003-2, API-QUOTE-ACCEPT

**Test Steps:**
1. **Given** a quote exists with status "Sent" and valid expiry
2. **When** customer calls AcceptQuote(quoteId, customerId)
3. **Then** response status is 200
4. **And** response contains orderId
5. **And** quote status is "Accepted"

**Code:**
```typescript
it('should accept Sent quote and return orderId', async () => {
  // Arrange
  const quote = await createQuote({ status: 'Sent', validUntil: tomorrow });

  // Act
  const response = await api.post(`/quotes/${quote.id}/accept`, {
    customerId: 'customer-123'
  });

  // Assert
  expect(response.status).toBe(200);
  expect(response.body.status).toBe('accepted');
  expect(response.body.orderId).toBeDefined();
});
```

---

## TC-QUOTE-003-002: Reject Draft quote (Negative)

**Type:** Unit
**Category:** Negative
**Traceability:** AC-QUOTE-003-2, BR-QUOTE-001, API-QUOTE-ACCEPT

**Test Steps:**
1. **Given** a quote exists with status "Draft"
2. **When** customer calls AcceptQuote(quoteId, customerId)
3. **Then** response status is 400
4. **And** error code is INVALID_QUOTE_STATUS

**Code:**
```typescript
it('should reject Draft quote with INVALID_QUOTE_STATUS', async () => {
  // Arrange
  const quote = await createQuote({ status: 'Draft' });

  // Act
  const response = await api.post(`/quotes/${quote.id}/accept`, {
    customerId: 'customer-123'
  });

  // Assert
  expect(response.status).toBe(400);
  expect(response.body.errorCode).toBe('INVALID_QUOTE_STATUS');
});
```

---

## TC-QUOTE-003-010: Should NOT send confirmation email (ShouldNOT)

**Type:** Unit
**Category:** ShouldNOT
**Purpose:** Prevent AI from adding email logic (handled by notification-service)

**Test Steps:**
1. **Given** email service is mocked
2. **When** quote is accepted
3. **Then** email service should NOT be called

**Code:**
```typescript
it('should NOT send email (notification-service responsibility)', () => {
  // Arrange
  const emailSpy = jest.spyOn(emailService, 'send');
  const quote = createQuote({ status: 'Sent' });

  // Act
  quoteService.accept(quote.id, 'customer-123');

  // Assert
  expect(emailSpy).not.toHaveBeenCalled();
});
```
```

---

Now read the L2 input files and generate comprehensive TDAI test cases.
