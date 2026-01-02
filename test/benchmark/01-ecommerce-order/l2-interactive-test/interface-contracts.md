# Interface Contracts

Generated: 2026-01-02T11:02:21+01:00

---

## Shared Types

### Email

| Field | Type | Constraints |
|-------|------|-------------|
| value | string | valid email format per RFC 5322 |

### CustomerStatus

| Field | Type | Constraints |
|-------|------|-------------|
| value | enum | unverified | verified | suspended |

---

## IC-CUSTOMER-001 – Customer Service {#ic-customer-001}

**Purpose:** Manages customer registration and account lifecycle

**Base URL:** `/api/v1/customers`

**Security:**
- Authentication: None (public endpoint)
- Authorization: None

### Operations

#### registerCustomer `POST /register`

Register a new customer account

**Input:**
- `email`: Email (required)
- `password`: string (required)
- `firstName`: string (required)
- `lastName`: string (required)

**Output:**
- `customerId`: UUID
- `email`: Email
- `firstName`: string
- `lastName`: string
- `status`: CustomerStatus
- `createdAt`: DateTime

**Errors:**
| Code | HTTP | Message |
|------|------|----------|
| DUPLICATE_EMAIL | 409 | Email already registered. Please login or reset your password. |
| INVALID_PASSWORD | 400 | Password must be at least 8 characters and contain at least one number. |
| INVALID_EMAIL | 400 | Invalid email format. |

**Preconditions:** [Email must not be registered to existing account]

**Postconditions:** [Customer account created with status 'unverified' Verification email sent to customer]

### Events

- **CustomerRegistered**: Emitted when a new customer account is created (payload: [customerId email firstName lastName status])
- **VerificationEmailSent**: Emitted when a verification email is dispatched (payload: [customerId email sentAt])

---

