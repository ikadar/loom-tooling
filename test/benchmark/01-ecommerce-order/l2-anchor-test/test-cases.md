# TDAI Test Cases

Generated: 2025-12-31T13:24:23+01:00

**Methodology:** Test-Driven AI Development (TDAI)

## Summary

| Category | Count | Ratio |
|----------|-------|-------|
| Positive | 3 | 23.0% |
| Negative | 5 | 38.0% |
| Boundary | 2 | - |
| Hallucination Prevention | 3 | - |
| **Total** | **13** | - |

**Coverage:** 1 ACs covered
**Hallucination Prevention:** ✓ Enabled

---

## Positive Tests (Happy Path)

### TC-AC-CUST-001-P01 – Valid registration with all required fields {#tc-ac-cust-001-p01}

**Preconditions:**
- Email john@example.com not registered

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| email | john@example.com | Valid email |
| password | SecurePass1 | Valid password |
| firstName | John | Required field |
| lastName | Doe | Required field |

**Steps:**
1. Navigate to registration page
2. Enter all required fields with valid data
3. Submit registration form

**Expected Result:**
- Customer account created with status 'unverified'
- Verification email sent to john@example.com

**Traceability:**
- AC: [AC-CUST-001](../l1/acceptance-criteria.md#ac-cust-001)
- BR: [BR-CUST-001](../l1/business-rules.md#br-cust-001)

---

### TC-AC-CUST-001-P02 – Valid registration with different email domain {#tc-ac-cust-001-p02}

**Preconditions:**
- Email jane@company.org not registered

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| email | jane@company.org | Different domain |
| password | MyPass123 | Valid password |
| firstName | Jane | Required field |
| lastName | Smith | Required field |

**Steps:**
1. Navigate to registration page
2. Enter valid data with .org email domain
3. Submit registration form

**Expected Result:**
- Customer account created with status 'unverified'
- Verification email sent to jane@company.org

**Traceability:**
- AC: [AC-CUST-001](../l1/acceptance-criteria.md#ac-cust-001)
- BR: [BR-CUST-001](../l1/business-rules.md#br-cust-001)

---

### TC-AC-CUST-001-P03 – Valid registration with complex password {#tc-ac-cust-001-p03}

**Preconditions:**
- Email test@test.com not registered

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| email | test@test.com | Valid email |
| password | C0mpl3x!Pass | Complex password |
| firstName | Test | Required field |
| lastName | User | Required field |

**Steps:**
1. Navigate to registration page
2. Enter complex password with special chars
3. Submit registration form

**Expected Result:**
- Customer account created with status 'unverified'
- Verification email sent

**Traceability:**
- AC: [AC-CUST-001](../l1/acceptance-criteria.md#ac-cust-001)
- BR: [BR-CUST-001](../l1/business-rules.md#br-cust-001)

---

## Negative Tests (Error Cases)

### TC-AC-CUST-001-N01 – Reject duplicate email registration {#tc-ac-cust-001-n01}

**Preconditions:**
- Email exists@example.com already registered

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| email | exists@example.com | Duplicate |
| password | SecurePass1 | Valid password |
| firstName | John | Required field |
| lastName | Doe | Required field |

**Steps:**
1. Navigate to registration page
2. Enter already registered email
3. Submit registration form

**Expected Result:**
- DUPLICATE_EMAIL error returned
- Suggest login or password reset
- No account created

**Traceability:**
- AC: [AC-CUST-001](../l1/acceptance-criteria.md#ac-cust-001)
- BR: [BR-CUST-001](../l1/business-rules.md#br-cust-001)

---

### TC-AC-CUST-001-N02 – Reject password shorter than 8 characters {#tc-ac-cust-001-n02}

**Preconditions:**
- Email new@example.com not registered

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| email | new@example.com | Valid email |
| password | Short1 | Only 6 chars |
| firstName | John | Required field |
| lastName | Doe | Required field |

**Steps:**
1. Navigate to registration page
2. Enter password with only 6 characters
3. Submit registration form

**Expected Result:**
- INVALID_PASSWORD error returned
- Requirements message displayed
- No account created

**Traceability:**
- AC: [AC-CUST-001](../l1/acceptance-criteria.md#ac-cust-001)
- BR: [BR-CUST-001](../l1/business-rules.md#br-cust-001)

---

### TC-AC-CUST-001-N03 – Reject password without number {#tc-ac-cust-001-n03}

**Preconditions:**
- Email new@example.com not registered

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| email | new@example.com | Valid email |
| password | NoNumberPass | No digits |
| firstName | John | Required field |
| lastName | Doe | Required field |

**Steps:**
1. Navigate to registration page
2. Enter password without any numbers
3. Submit registration form

**Expected Result:**
- INVALID_PASSWORD error returned
- Requirements message displayed
- No account created

**Traceability:**
- AC: [AC-CUST-001](../l1/acceptance-criteria.md#ac-cust-001)
- BR: [BR-CUST-001](../l1/business-rules.md#br-cust-001)

---

### TC-AC-CUST-001-N04 – Reject invalid email format {#tc-ac-cust-001-n04}

**Preconditions:**
- None

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| email | invalid-email | Missing @ and domain |
| password | SecurePass1 | Valid password |
| firstName | John | Required field |
| lastName | Doe | Required field |

**Steps:**
1. Navigate to registration page
2. Enter email without @ symbol
3. Submit registration form

**Expected Result:**
- INVALID_EMAIL error returned
- No account created

**Traceability:**
- AC: [AC-CUST-001](../l1/acceptance-criteria.md#ac-cust-001)
- BR: [BR-CUST-001](../l1/business-rules.md#br-cust-001)

---

### TC-AC-CUST-001-N05 – Reject email without domain extension {#tc-ac-cust-001-n05}

**Preconditions:**
- None

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| email | user@domain | Missing TLD |
| password | SecurePass1 | Valid password |
| firstName | John | Required field |
| lastName | Doe | Required field |

**Steps:**
1. Navigate to registration page
2. Enter email without domain extension
3. Submit registration form

**Expected Result:**
- INVALID_EMAIL error returned
- No account created

**Traceability:**
- AC: [AC-CUST-001](../l1/acceptance-criteria.md#ac-cust-001)
- BR: [BR-CUST-001](../l1/business-rules.md#br-cust-001)

---

## Boundary Tests

### TC-AC-CUST-001-B01 – Password at minimum length (8 characters) {#tc-ac-cust-001-b01}

**Preconditions:**
- Email boundary@example.com not registered

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| email | boundary@example.com | Valid |
| password | Pass1234 | Exactly 8 chars |
| firstName | Test | Required field |
| lastName | User | Required field |

**Steps:**
1. Navigate to registration page
2. Enter password with exactly 8 characters
3. Submit registration form

**Expected Result:**
- Registration succeeds
- Customer account created

**Traceability:**
- AC: [AC-CUST-001](../l1/acceptance-criteria.md#ac-cust-001)
- BR: [BR-CUST-001](../l1/business-rules.md#br-cust-001)

---

### TC-AC-CUST-001-B02 – Password below minimum (7 characters) {#tc-ac-cust-001-b02}

**Preconditions:**
- Email boundary@example.com not registered

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| email | boundary@example.com | Valid |
| password | Pass123 | Only 7 chars |
| firstName | Test | Required field |
| lastName | User | Required field |

**Steps:**
1. Navigate to registration page
2. Enter password with 7 characters
3. Submit registration form

**Expected Result:**
- INVALID_PASSWORD error returned
- No account created

**Traceability:**
- AC: [AC-CUST-001](../l1/acceptance-criteria.md#ac-cust-001)
- BR: [BR-CUST-001](../l1/business-rules.md#br-cust-001)

---

## Hallucination Prevention Tests

### TC-AC-CUST-001-H01 – Phone number should NOT be required {#tc-ac-cust-001-h01}

**⚠️ Should NOT:** Require phone number (not specified in AC)

**Preconditions:**
- Email nophone@example.com not registered

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| email | nophone@example.com | Valid |
| password | SecurePass1 | Valid password |
| firstName | John | Required field |
| lastName | Doe | Required field |
| phone |  | Empty - not in AC |

**Steps:**
1. Navigate to registration page
2. Fill all AC-specified fields
3. Leave phone field empty
4. Submit registration form

**Expected Result:**
- Registration succeeds without phone

**Traceability:**
- AC: [AC-CUST-001](../l1/acceptance-criteria.md#ac-cust-001)

---

### TC-AC-CUST-001-H02 – Address should NOT be required {#tc-ac-cust-001-h02}

**⚠️ Should NOT:** Require address fields (not specified in AC)

**Preconditions:**
- Email noaddr@example.com not registered

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| email | noaddr@example.com | Valid |
| password | SecurePass1 | Valid password |
| firstName | John | Required field |
| lastName | Doe | Required field |
| address |  | Empty - not in AC |

**Steps:**
1. Navigate to registration page
2. Fill only AC-specified fields
3. Leave address empty
4. Submit registration form

**Expected Result:**
- Registration succeeds without address

**Traceability:**
- AC: [AC-CUST-001](../l1/acceptance-criteria.md#ac-cust-001)

---

### TC-AC-CUST-001-H03 – Account should NOT be auto-verified {#tc-ac-cust-001-h03}

**⚠️ Should NOT:** Auto-verify account (AC specifies 'unverified' status)

**Preconditions:**
- Email verify@example.com not registered

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| email | verify@example.com | Valid |
| password | SecurePass1 | Valid password |
| firstName | John | Required field |
| lastName | Doe | Required field |

**Steps:**
1. Complete registration successfully
2. Check account status immediately

**Expected Result:**
- Account status is 'unverified'
- Verification email sent

**Traceability:**
- AC: [AC-CUST-001](../l1/acceptance-criteria.md#ac-cust-001)

---

