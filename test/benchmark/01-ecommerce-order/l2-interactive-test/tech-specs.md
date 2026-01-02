# Technical Specifications

Generated: 2026-01-02T11:02:21+01:00

---

## TS-BR-CUST-001 – Password validation requirements {#ts-br-cust-001}

**Rule:** Password must be at least 8 characters and contain at least one number

**Implementation Approach:**
Validate password against regex pattern on registration and password change endpoints. Pattern: ^(?=.*[0-9]).{8,}$

**Validation Points:**
- Registration API
- Password change API
- Registration UI
- Password change UI

**Data Requirements:**
| Field | Type | Constraints | Source |
|-------|------|-------------|--------|
| password | string | min_length: 8, pattern: contains at least one digit [0-9] | BR-CUST-001 |

**Error Handling:**
| Condition | Error Code | Message | HTTP Status |
|-----------|------------|---------|-------------|
| password.length < 8 | PASSWORD_TOO_SHORT | Password must be at least 8 characters | 400 |
| !password.matches(/[0-9]/) | PASSWORD_MISSING_NUMBER | Password must contain at least one number | 400 |

**Traceability:**
- BR: [BR-CUST-001](../l1/business-rules.md#br-cust-001)

---

## TS-BR-CUST-002 – Email uniqueness constraint {#ts-br-cust-002}

**Rule:** Each email address can only be associated with one customer account

**Implementation Approach:**
Enforce unique constraint on email column in customers table. Check for existing email before insert on registration endpoint.

**Validation Points:**
- Registration API
- Email update API
- Database constraint

**Data Requirements:**
| Field | Type | Constraints | Source |
|-------|------|-------------|--------|
| email | string | unique, valid email format, max_length: 255, case_insensitive comparison | BR-CUST-002 |

**Error Handling:**
| Condition | Error Code | Message | HTTP Status |
|-----------|------------|---------|-------------|
| email already exists in database | EMAIL_ALREADY_REGISTERED | This email address is already associated with an account | 409 |

**Traceability:**
- BR: [BR-CUST-002](../l1/business-rules.md#br-cust-002)

---

