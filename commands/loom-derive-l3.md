---
name: loom-derive-l3
description: Generate L3 test cases from L2 interface contracts using TDAI and Structured Interview
version: "2.0.0"
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

# Loom L2 → L3 TDAI Test Generation Skill (v2.0 - Structured Interview)

You are an expert test-driven development agent for the **Loom AI Development Orchestration Platform**.

Your task is to generate **Level 3 (L3) test cases** from **Level 2 (L2) interface contracts** using **TDAI (Test-Driven AI Development)** principles.

## ⚠️ CRITICAL: Structured Interview Required

**BEFORE generating ANY test code, you MUST conduct a Structured Interview to resolve test strategy decisions.**

The Structured Interview prevents implicit decisions by:
1. Identifying decision points in test generation
2. Checking if input documents answer them
3. Asking the user for unresolved decisions
4. Only then generating tests with explicit choices

---

## Decision Points Catalog for L3 Test Generation

### Category 1: Test Strategy (TST)

| ID | Decision Point | Question Template | Default if Unasked |
|----|----------------|-------------------|-------------------|
| TST-1 | Pyramid ratio | "Standard 70:20:10 (unit:integration:E2E) or custom ratio?" | 70:20:10 |
| TST-2 | Test style | "BDD style (describe/it) or traditional (test functions)?" | BDD |
| TST-3 | Async handling | "How to handle async operations in tests? (await/callbacks/done)" | await |
| TST-4 | Coverage target | "What code coverage target? (80%/90%/100%)" | 80% |

### Category 2: Mock Strategy (MOC)

| ID | Decision Point | Question Template | Default if Unasked |
|----|----------------|-------------------|-------------------|
| MOC-1 | External services | "Mock all external services or use test doubles for some?" | ASK - no default |
| MOC-2 | Database | "Mock DB, use in-memory DB, or real test DB?" | ASK - no default |
| MOC-3 | Message queue | "Mock message queue or use test container?" | Mock |
| MOC-4 | Time/dates | "Mock time for date-dependent tests?" | Yes, mock time |

### Category 3: Test Data (TDA)

| ID | Decision Point | Question Template | Default if Unasked |
|----|----------------|-------------------|-------------------|
| TDA-1 | Data creation | "Fixtures (static JSON) or factories (builder functions)?" | ASK - no default |
| TDA-2 | Shared data | "Shared test data across suites or isolated per test?" | Isolated |
| TDA-3 | Edge cases | "Auto-generate edge cases or manually specify?" | Auto-generate |
| TDA-4 | Cleanup | "Clean up test data after each test or after suite?" | After each |

### Category 4: Coverage Priority (COV)

| ID | Decision Point | Question Template | Default if Unasked |
|----|----------------|-------------------|-------------------|
| COV-1 | Critical paths | "Which operations need E2E tests? (list or all)" | ASK - no default |
| COV-2 | Error scenarios | "Test all error codes or only critical ones?" | All |
| COV-3 | Performance | "Include performance/load test stubs?" | No |
| COV-4 | Security | "Include security test cases (auth, injection)?" | Yes |

### Category 5: Environment (ENV)

| ID | Decision Point | Question Template | Default if Unasked |
|----|----------------|-------------------|-------------------|
| ENV-1 | CI execution | "Tests run in CI? Need specific setup?" | ASK if unclear |
| ENV-2 | Parallelization | "Run tests in parallel or sequential?" | Parallel |
| ENV-3 | Test containers | "Use Docker/Testcontainers for dependencies?" | ASK - no default |
| ENV-4 | Snapshot testing | "Use snapshot tests for response structures?" | No |

---

## Phase 1: Structured Interview

### Step 1.1: Analyze Inputs for Answers

Read the L2 interface contracts and extract:
- Test framework already specified → TST decisions
- Mock patterns mentioned → MOC decisions
- Test data examples → TDA decisions
- Critical paths identified → COV decisions

### Step 1.2: Identify Gaps

For each decision point category, check:
- Is it answered in input? → Mark "from-input"
- Is it a safe default? → Mark "default-applied"
- Must ask user? → Add to interview questions

### Step 1.3: Conduct Interview

Present questions grouped by category:

```
## Test Strategy Interview

I need to resolve some decisions before generating tests.

**Answered from input:**
- TST-2: BDD style (jest with describe/it in examples)
- MOC-3: Mock message queue (event-driven architecture)

**Questions for you:**

### Mock Strategy
1. **MOC-1 (External Services):** The interface contracts show calls to OrderService and NotificationService. How should tests handle these?
   a) Mock all external services (faster, isolated)
   b) Use test doubles with real behavior verification
   c) Integration tests with real services in Docker

2. **MOC-2 (Database):** Quote acceptance updates database. How to handle in tests?
   a) Mock repository layer
   b) In-memory SQLite/H2
   c) Real PostgreSQL in Docker (testcontainers)

### Test Data
3. **TDA-1 (Data Creation):** How to create test quotes and customers?
   a) Static fixtures (JSON files)
   b) Factory functions (QuoteFactory.create())
   c) Builder pattern (QuoteBuilder().withStatus('Sent').build())

### Coverage Priority
4. **COV-1 (Critical Paths):** Which flows need full E2E tests?
   a) All API operations
   b) Only Quote Acceptance flow (critical business path)
   c) Acceptance + Reversal (both state changes)

Please answer with numbers (e.g., "1a, 2c, 3b, 4c") or provide custom answers.
```

### Step 1.4: Record Decisions

After receiving answers, document all decisions:

```markdown
## Decisions Resolved

| ID | Decision | Answer | Source |
|----|----------|--------|--------|
| TST-2 | Test style | BDD (describe/it) | Input |
| MOC-1 | External services | Mock all | User (1a) |
| MOC-2 | Database | Testcontainers | User (2c) |
| TDA-1 | Data creation | Builder pattern | User (3c) |
| COV-1 | Critical paths | Accept + Reversal | User (4c) |
```

---

## Phase 2: TDAI Test Generation (After Interview)

**Only proceed to this phase after all decision points are resolved.**

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

Create with YAML frontmatter including Structured Interview metadata:

```yaml
---
status: draft
derived-from:
  - "{contracts-file}"
  - "{ac-file}"
  - "{br-file}"
derived-at: "{ISO timestamp}"
derived-by: "loom-derive-l3 skill v2.0 (Structured Interview)"
loom-version: "3.0.0"
tdai-version: "1.0"
structured-interview:
  decision-points-resolved: {N}
  from-user-answers: {X}
  from-input: {Y}
  defaults-applied: {Z}
  test-strategy:
    pyramid-ratio: "70:20:10"
    mock-strategy: "mock-all"
    data-creation: "builder-pattern"
    critical-paths:
      - "quote-acceptance"
      - "quote-reversal"
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

## Structured Interview Impact for L3

### Without Structured Interview (Implicit Decisions)

```typescript
// Test file generated with implicit assumptions:

it('should accept quote and create order', async () => {
  // IMPLICIT: Using fixtures (could be wrong for team)
  const quote = testFixtures.validQuote;

  // IMPLICIT: Mocking all services (maybe integration needed)
  jest.mock('../services/orderService');

  // IMPLICIT: Only happy path (missing critical error tests)
  const result = await quoteService.accept(quote.id);
  expect(result.status).toBe('accepted');
});

// Problems:
// - Team prefers builder pattern → fixtures cause confusion
// - OrderService mock hides integration bugs
// - No E2E test for critical business flow
// - 95% positive tests, 5% negative (should be 80%/20%)
```

### With Structured Interview (Explicit Decisions)

```typescript
// Test file generated with explicit choices from interview:

// TDA-1: Builder pattern (User choice: 3c)
const quoteBuilder = new QuoteBuilder();

// MOC-2: Testcontainers (User choice: 2c)
beforeAll(async () => {
  await PostgresContainer.start();
});

// COV-1: E2E for critical paths (User choice: 4c)
describe('E2E: Quote Acceptance Flow', () => {
  it('should accept quote and trigger order creation', async () => {
    const quote = quoteBuilder
      .withStatus('Sent')
      .withValidUntil(tomorrow)
      .build();

    const response = await api.post(`/quotes/${quote.id}/accept`);

    expect(response.status).toBe(200);
    // Integration: verify order actually created
    const order = await orderRepository.findByQuoteId(quote.id);
    expect(order).toBeDefined();
  });
});

// Benefits:
// - Builder pattern matches team conventions
// - Real DB catches integration issues
// - E2E tests cover critical business flows
// - Balanced test categories (explicit from interview)
```

### Test Strategy Traceability

Every test case includes decision point references:

```typescript
/**
 * @traceability AC-QUOTE-003-2, BR-QUOTE-001
 * @decision MOC-2: Real DB via testcontainers
 * @decision TDA-1: Builder pattern for test data
 */
it('should reject expired quote', async () => {
  // Test implementation...
});
```

---

## Phase 3: Validation

Before presenting tests, validate:

1. **Interview Complete:** All "ASK - no default" decisions resolved
2. **Test Balance:** ≥20% negative tests, ~70:20:10 pyramid
3. **Traceability:** Every test links to AC/BR/API
4. **Decision Documentation:** All choices recorded in frontmatter

---

Now conduct the Structured Interview, then generate comprehensive TDAI test cases.
