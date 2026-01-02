# TDAI Test Cases

Generated: 2026-01-02T11:02:21+01:00

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
2. Enter email, password, first name, last name
3. Submit registration form

**Expected Result:**
- Account created with status 'unverified'
- Verification email sent to john@example.com

**Traceability:**
- AC: [AC-CUST-001](../l1/acceptance-criteria.md#ac-cust-001)
- BR: [BR-CUST-001](../l1/business-rules.md#br-cust-001)

---

### TC-AC-CUST-001-P02 – Valid registration with different valid email {#tc-ac-cust-001-p02}

**Preconditions:**
- Email jane.doe@company.org not registered

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| email | jane.doe@company.org | Valid email |
| password | MyPass123 | Valid password |
| firstName | Jane | Required field |
| lastName | Smith | Required field |

**Steps:**
1. Navigate to registration page
2. Enter all required fields
3. Submit registration form

**Expected Result:**
- Account created with status 'unverified'
- Verification email sent to jane.doe@company.org

**Traceability:**
- AC: [AC-CUST-001](../l1/acceptance-criteria.md#ac-cust-001)
- BR: [BR-CUST-001](../l1/business-rules.md#br-cust-001)

---

### TC-AC-CUST-001-P03 – Valid registration with complex password {#tc-ac-cust-001-p03}

**Preconditions:**
- Email test.user@domain.com not registered

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| email | test.user@domain.com | Valid email |
| password | C0mpl3x!Pass#2024 | Strong password |
| firstName | Test | Required field |
| lastName | User | Required field |

**Steps:**
1. Navigate to registration page
2. Enter all required fields with complex password
3. Submit registration form

**Expected Result:**
- Account created with status 'unverified'
- Verification email sent

**Traceability:**
- AC: [AC-CUST-001](../l1/acceptance-criteria.md#ac-cust-001)
- BR: [BR-CUST-001](../l1/business-rules.md#br-cust-001)

---

## Negative Tests (Error Cases)

### TC-AC-CUST-001-N01 – Reject duplicate email {#tc-ac-cust-001-n01}

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
2. Enter existing email with valid data
3. Submit registration form

**Expected Result:**
- DUPLICATE_EMAIL error returned
- Suggest login or password reset

**Traceability:**
- AC: [AC-CUST-001](../l1/acceptance-criteria.md#ac-cust-001)
- BR: [BR-CUST-001](../l1/business-rules.md#br-cust-001)

---

### TC-AC-CUST-001-N02 – Reject password less than 8 characters {#tc-ac-cust-001-n02}

**Preconditions:**
- Email new@example.com not registered

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| email | new@example.com | Valid email |
| password | Pass1 | Only 5 chars |
| firstName | John | Required field |
| lastName | Doe | Required field |

**Steps:**
1. Navigate to registration page
2. Enter password with less than 8 characters
3. Submit registration form

**Expected Result:**
- INVALID_PASSWORD error returned
- Requirements message displayed

**Traceability:**
- AC: [AC-CUST-001](../l1/acceptance-criteria.md#ac-cust-001)
- BR: [BR-CUST-001](../l1/business-rules.md#br-cust-001)

---

### TC-AC-CUST-001-N03 – Reject password without number {#tc-ac-cust-001-n03}

**Preconditions:**
- Email new2@example.com not registered

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| email | new2@example.com | Valid email |
| password | SecurePass | No number |
| firstName | John | Required field |
| lastName | Doe | Required field |

**Steps:**
1. Navigate to registration page
2. Enter password without any number
3. Submit registration form

**Expected Result:**
- INVALID_PASSWORD error returned
- Requirements message displayed

**Traceability:**
- AC: [AC-CUST-001](../l1/acceptance-criteria.md#ac-cust-001)
- BR: [BR-CUST-001](../l1/business-rules.md#br-cust-001)

---

### TC-AC-CUST-001-N04 – Reject invalid email format {#tc-ac-cust-001-n04}

**Preconditions:**
- User on registration page

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| email | invalid-email | Missing @ and domain |
| password | SecurePass1 | Valid password |
| firstName | John | Required field |
| lastName | Doe | Required field |

**Steps:**
1. Navigate to registration page
2. Enter invalid email format
3. Submit registration form

**Expected Result:**
- INVALID_EMAIL error returned

**Traceability:**
- AC: [AC-CUST-001](../l1/acceptance-criteria.md#ac-cust-001)
- BR: [BR-CUST-001](../l1/business-rules.md#br-cust-001)

---

### TC-AC-CUST-001-N05 – Reject email without domain {#tc-ac-cust-001-n05}

**Preconditions:**
- User on registration page

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| email | john@ | Missing domain |
| password | SecurePass1 | Valid password |
| firstName | John | Required field |
| lastName | Doe | Required field |

**Steps:**
1. Navigate to registration page
2. Enter email without domain
3. Submit registration form

**Expected Result:**
- INVALID_EMAIL error returned

**Traceability:**
- AC: [AC-CUST-001](../l1/acceptance-criteria.md#ac-cust-001)
- BR: [BR-CUST-001](../l1/business-rules.md#br-cust-001)

---

## Boundary Tests

### TC-AC-CUST-001-B01 – Password at minimum length (8 chars) {#tc-ac-cust-001-b01}

**Preconditions:**
- Email boundary1@example.com not registered

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| email | boundary1@example.com | Valid |
| password | Secure1! | Exactly 8 chars |
| firstName | Test | Required field |
| lastName | User | Required field |

**Steps:**
1. Navigate to registration page
2. Enter password with exactly 8 characters
3. Submit registration form

**Expected Result:**
- Registration succeeds
- Account created with status 'unverified'

**Traceability:**
- AC: [AC-CUST-001](../l1/acceptance-criteria.md#ac-cust-001)
- BR: [BR-CUST-001](../l1/business-rules.md#br-cust-001)

---

### TC-AC-CUST-001-B02 – Password at 7 chars (below minimum) {#tc-ac-cust-001-b02}

**Preconditions:**
- Email boundary2@example.com not registered

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| email | boundary2@example.com | Valid |
| password | Secur1! | Only 7 chars |
| firstName | Test | Required field |
| lastName | User | Required field |

**Steps:**
1. Navigate to registration page
2. Enter password with 7 characters
3. Submit registration form

**Expected Result:**
- INVALID_PASSWORD error returned
- Requirements message displayed

**Traceability:**
- AC: [AC-CUST-001](../l1/acceptance-criteria.md#ac-cust-001)
- BR: [BR-CUST-001](../l1/business-rules.md#br-cust-001)

---

## Hallucination Prevention Tests

### TC-AC-CUST-001-H01 – Phone number should NOT be required {#tc-ac-cust-001-h01}

**⚠️ Should NOT:** Require phone number (not specified in AC)

**Preconditions:**
- Email hal1@example.com not registered

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| email | hal1@example.com | Valid email |
| password | SecurePass1 | Valid password |
| firstName | John | Required field |
| lastName | Doe | Required field |
| phone |  | Empty - not in AC |

**Steps:**
1. Navigate to registration page
2. Fill only email, password, firstName, lastName
3. Submit without phone number

**Expected Result:**
- Registration succeeds without phone

**Traceability:**
- AC: [AC-CUST-001](../l1/acceptance-criteria.md#ac-cust-001)

---

### TC-AC-CUST-001-H02 – Address should NOT be required {#tc-ac-cust-001-h02}

**⚠️ Should NOT:** Require address fields (not specified in AC)

**Preconditions:**
- Email hal2@example.com not registered

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| email | hal2@example.com | Valid email |
| password | SecurePass1 | Valid password |
| firstName | John | Required field |
| lastName | Doe | Required field |
| address |  | Empty - not in AC |

**Steps:**
1. Navigate to registration page
2. Fill only required fields per AC
3. Submit without address

**Expected Result:**
- Registration succeeds without address

**Traceability:**
- AC: [AC-CUST-001](../l1/acceptance-criteria.md#ac-cust-001)

---

### TC-AC-CUST-001-H03 – Account should NOT be auto-verified {#tc-ac-cust-001-h03}

**⚠️ Should NOT:** Auto-verify account (AC states 'unverified')

**Preconditions:**
- Email hal3@example.com not registered

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| email | hal3@example.com | Valid email |
| password | SecurePass1 | Valid password |
| firstName | John | Required field |
| lastName | Doe | Required field |

**Steps:**
1. Complete valid registration
2. Check account status immediately after

**Expected Result:**
- Account status is 'unverified'
- Not 'verified' or 'active'

**Traceability:**
- AC: [AC-CUST-001](../l1/acceptance-criteria.md#ac-cust-001)

---

