# TDAI Test Cases

Generated: 2025-12-31T10:31:00+01:00

**Methodology:** Test-Driven AI Development (TDAI)

## Summary

| Category | Count | Ratio |
|----------|-------|-------|
| Positive | 3 | 21.0% |
| Negative | 5 | 36.0% |
| Boundary | 2 | - |
| Hallucination Prevention | 4 | - |
| **Total** | **14** | - |

**Coverage:** 1 ACs covered
**Hallucination Prevention:** ✓ Enabled

---

## Positive Tests (Happy Path)

### TC-AC-CUST-001-P01 – Valid registration with all required fields

**Preconditions:**
- Email john@example.com not registered

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| email | john@example.com | Valid email |
| password | SecurePass1 | Valid password |
| firstName | John | Required |
| lastName | Doe | Required |

**Steps:**
1. Navigate to registration page
2. Enter all required fields
3. Submit registration form

**Expected Result:**
- Customer account created with status 'unverified'
- Verification email sent to john@example.com

**Traceability:**
- AC: AC-CUST-001
- BR: [BR-CUST-001]

---

### TC-AC-CUST-001-P02 – Valid registration with different valid email

**Preconditions:**
- Email jane.doe@company.org not registered

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| email | jane.doe@company.org | Valid |
| password | MyPass123 | Valid password |
| firstName | Jane | Required |
| lastName | Doe | Required |

**Steps:**
1. Navigate to registration page
2. Enter all required fields
3. Submit registration form

**Expected Result:**
- Customer account created with status 'unverified'
- Verification email sent to jane.doe@company.org

**Traceability:**
- AC: AC-CUST-001
- BR: [BR-CUST-001]

---

### TC-AC-CUST-001-P03 – Valid registration with complex password

**Preconditions:**
- Email test@example.com not registered

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| email | test@example.com | Valid |
| password | C0mpl3x!Pass#2024 | Complex |
| firstName | Test | Required |
| lastName | User | Required |

**Steps:**
1. Navigate to registration page
2. Enter all required fields with complex password
3. Submit registration form

**Expected Result:**
- Customer account created with status 'unverified'
- Verification email sent

**Traceability:**
- AC: AC-CUST-001
- BR: [BR-CUST-001]

---

## Negative Tests (Error Cases)

### TC-AC-CUST-001-N01 – Reject duplicate email address

**Preconditions:**
- Email exists@example.com already registered

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| email | exists@example.com | Duplicate |
| password | SecurePass1 | Valid |
| firstName | John | Required |
| lastName | Doe | Required |

**Steps:**
1. Navigate to registration page
2. Enter existing email with valid other fields
3. Submit registration form

**Expected Result:**
- DUPLICATE_EMAIL error returned
- Suggestion to login or reset password displayed

**Traceability:**
- AC: AC-CUST-001
- BR: [BR-CUST-001]

---

### TC-AC-CUST-001-N02 – Reject password less than 8 characters

**Preconditions:**
- Email new@example.com not registered

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| email | new@example.com | Valid |
| password | Pass1 | Only 5 chars |
| firstName | John | Required |
| lastName | Doe | Required |

**Steps:**
1. Navigate to registration page
2. Enter password with less than 8 characters
3. Submit registration form

**Expected Result:**
- INVALID_PASSWORD error returned
- Password requirements message displayed

**Traceability:**
- AC: AC-CUST-001
- BR: [BR-CUST-001]

---

### TC-AC-CUST-001-N03 – Reject password without number

**Preconditions:**
- Email new@example.com not registered

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| email | new@example.com | Valid |
| password | SecurePass | No number |
| firstName | John | Required |
| lastName | Doe | Required |

**Steps:**
1. Navigate to registration page
2. Enter password without any number
3. Submit registration form

**Expected Result:**
- INVALID_PASSWORD error returned
- Password requirements message displayed

**Traceability:**
- AC: AC-CUST-001
- BR: [BR-CUST-001]

---

### TC-AC-CUST-001-N04 – Reject invalid email format

**Preconditions:**
- User on registration page

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| email | invalid-email | No @ symbol |
| password | SecurePass1 | Valid |
| firstName | John | Required |
| lastName | Doe | Required |

**Steps:**
1. Navigate to registration page
2. Enter email without @ symbol
3. Submit registration form

**Expected Result:**
- INVALID_EMAIL error returned

**Traceability:**
- AC: AC-CUST-001
- BR: [BR-CUST-001]

---

### TC-AC-CUST-001-N05 – Reject email without domain

**Preconditions:**
- User on registration page

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| email | john@ | No domain |
| password | SecurePass1 | Valid |
| firstName | John | Required |
| lastName | Doe | Required |

**Steps:**
1. Navigate to registration page
2. Enter email without domain part
3. Submit registration form

**Expected Result:**
- INVALID_EMAIL error returned

**Traceability:**
- AC: AC-CUST-001
- BR: [BR-CUST-001]

---

## Boundary Tests

### TC-AC-CUST-001-B01 – Password at minimum length (8 characters)

**Preconditions:**
- Email boundary@example.com not registered

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| email | boundary@example.com | Valid |
| password | Pass123! | Exactly 8 chars |
| firstName | John | Required |
| lastName | Doe | Required |

**Steps:**
1. Navigate to registration page
2. Enter password with exactly 8 characters
3. Submit registration form

**Expected Result:**
- Registration succeeds
- Customer account created with status 'unverified'

**Traceability:**
- AC: AC-CUST-001
- BR: [BR-CUST-001]

---

### TC-AC-CUST-001-B02 – Password at 7 characters (below minimum)

**Preconditions:**
- Email boundary@example.com not registered

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| email | boundary@example.com | Valid |
| password | Pass12! | Only 7 chars |
| firstName | John | Required |
| lastName | Doe | Required |

**Steps:**
1. Navigate to registration page
2. Enter password with 7 characters
3. Submit registration form

**Expected Result:**
- INVALID_PASSWORD error returned
- Password requirements message displayed

**Traceability:**
- AC: AC-CUST-001
- BR: [BR-CUST-001]

---

## Hallucination Prevention Tests

### TC-AC-CUST-001-H01 – Middle name should NOT be required

**⚠️ Should NOT:** Require middle name (not in AC)

**Preconditions:**
- Email hal@example.com not registered

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| email | hal@example.com | Valid |
| password | SecurePass1 | Valid |
| firstName | John | Required |
| lastName | Doe | Required |
| middleName |  | Not provided |

**Steps:**
1. Navigate to registration page
2. Fill form without middle name
3. Submit registration form

**Expected Result:**
- Registration succeeds without middle name

**Traceability:**
- AC: AC-CUST-001

---

### TC-AC-CUST-001-H02 – Phone number should NOT be required

**⚠️ Should NOT:** Require phone number (not in AC)

**Preconditions:**
- Email hal2@example.com not registered

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| email | hal2@example.com | Valid |
| password | SecurePass1 | Valid |
| firstName | John | Required |
| lastName | Doe | Required |
| phone |  | Not provided |

**Steps:**
1. Navigate to registration page
2. Fill form without phone number
3. Submit registration form

**Expected Result:**
- Registration succeeds without phone

**Traceability:**
- AC: AC-CUST-001

---

### TC-AC-CUST-001-H03 – Address should NOT be required for registration

**⚠️ Should NOT:** Require address for registration (not in AC)

**Preconditions:**
- Email hal3@example.com not registered

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| email | hal3@example.com | Valid |
| password | SecurePass1 | Valid |
| firstName | John | Required |
| lastName | Doe | Required |
| address |  | Not provided |

**Steps:**
1. Navigate to registration page
2. Fill form without address
3. Submit registration form

**Expected Result:**
- Registration succeeds without address

**Traceability:**
- AC: AC-CUST-001

---

### TC-AC-CUST-001-H04 – Account should NOT be auto-verified

**⚠️ Should NOT:** Auto-verify account (AC specifies 'unverified' status)

**Preconditions:**
- Email hal4@example.com not registered

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| email | hal4@example.com | Valid |
| password | SecurePass1 | Valid |
| firstName | John | Required |
| lastName | Doe | Required |

**Steps:**
1. Navigate to registration page
2. Complete registration with valid data
3. Check account status immediately

**Expected Result:**
- Account status is 'unverified'
- Account is NOT automatically verified

**Traceability:**
- AC: AC-CUST-001

---

