# Technical Specifications

Generated: 2025-12-31T13:24:23+01:00

---

## TS-BR-CUST-001 – Password validation requirements {#ts-br-cust-001}

**Rule:** Password must be at least 8 characters and contain at least one number

**Implementation Approach:**
Apply regex validation on password field during registration and password change operations. Validate on both client-side (for UX) and server-side (for security).

**Validation Points:**
- Registration API endpoint
- Password change API endpoint
- Registration form UI
- Password change form UI

**Data Requirements:**
| Field | Type | Constraints | Source |
|-------|------|-------------|--------|
| password | string | minLength: 8, pattern: /.*[0-9].*/ | BR-CUST-001 |

**Error Handling:**
| Condition | Error Code | Message | HTTP Status |
|-----------|------------|---------|-------------|
| password.length < 8 | PASSWORD_TOO_SHORT | Password must be at least 8 characters long | 400 |
| !password.match(/[0-9]/) | PASSWORD_MISSING_NUMBER | Password must contain at least one number | 400 |

**Traceability:**
- BR: [BR-CUST-001](../l1/business-rules.md#br-cust-001)
- Related ACs: [AC-REG-001](../l1/acceptance-criteria.md#ac-reg-001), [AC-PWD-001](../l1/acceptance-criteria.md#ac-pwd-001)

---

## TS-BR-CUST-002 – Email uniqueness constraint {#ts-br-cust-002}

**Rule:** Each email address can only be associated with one customer account

**Implementation Approach:**
Add unique constraint on email column in customers table. Check for existing email before creating new account. Use case-insensitive comparison (normalize to lowercase before storage).

**Validation Points:**
- Registration API endpoint
- Email change API endpoint
- Database constraint level

**Data Requirements:**
| Field | Type | Constraints | Source |
|-------|------|-------------|--------|
| email | string | unique, format: email, maxLength: 255, stored as lowercase | BR-CUST-002 |

**Error Handling:**
| Condition | Error Code | Message | HTTP Status |
|-----------|------------|---------|-------------|
| email already exists in database | EMAIL_ALREADY_EXISTS | An account with this email address already exists | 409 |
| invalid email format | INVALID_EMAIL_FORMAT | Please provide a valid email address | 400 |

**Traceability:**
- BR: [BR-CUST-002](../l1/business-rules.md#br-cust-002)
- Related ACs: [AC-REG-001](../l1/acceptance-criteria.md#ac-reg-001), [AC-EMAIL-001](../l1/acceptance-criteria.md#ac-email-001)

---

