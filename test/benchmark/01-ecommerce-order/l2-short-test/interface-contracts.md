# Interface Contracts

Generated: 2025-12-31T10:31:00+01:00

---

## Shared Types

### Email

| Field | Type | Constraints |
|-------|------|-------------|
| value | string | RFC 5322 format, max 255 chars, unique per customer |

### CustomerStatus

| Field | Type | Constraints |
|-------|------|-------------|
| value | enum | unverified | verified | suspended |

### Password

| Field | Type | Constraints |
|-------|------|-------------|
| value | string | min 8 chars, must contain at least one number, stored as bcrypt hash |

### Customer

| Field | Type | Constraints |
|-------|------|-------------|
| id | UUID | primary key |
| email | Email | unique, required |
| firstName | string | 1-100 chars, required |
| lastName | string | 1-100 chars, required |
| status | CustomerStatus | default: unverified |
| createdAt | DateTime | auto-generated |
| verifiedAt | DateTime | nullable |

---

## IC-CUST-001 – Customer Service

**Purpose:** Manages customer registration, authentication, and account lifecycle

**Base URL:** `/api/v1/customers`

**Security:**
- Authentication: Bearer JWT (except registration endpoint)
- Authorization: Customers can only access their own data

### Operations

#### registerCustomer `POST /register`

Register a new customer account

**Input:**
- `email`: Email (required)
- `password`: string (required)
- `firstName`: string (required)
- `lastName`: string (required)

**Output:**
- `lastName`: string
- `status`: CustomerStatus
- `createdAt`: DateTime
- `customerId`: UUID
- `email`: Email
- `firstName`: string

**Errors:**
| Code | HTTP | Message |
|------|------|----------|
| DUPLICATE_EMAIL | 409 | Email already registered. Please login or reset password |
| INVALID_PASSWORD | 400 | Password must be at least 8 characters and contain at least one number |
| INVALID_EMAIL | 400 | Please provide a valid email address |
| MISSING_REQUIRED_FIELD | 400 | Required field is missing |

**Preconditions:** [Email must not be registered]

**Postconditions:** [Customer account created with status 'unverified' Verification email sent to customer]

#### verifyEmail `POST /{customerId}/verify`

Verify customer email address using verification token

**Input:**
- `customerId`: UUID (required)
- `verificationToken`: string (required)

**Output:**
- `customerId`: UUID
- `status`: CustomerStatus
- `verifiedAt`: DateTime

**Errors:**
| Code | HTTP | Message |
|------|------|----------|
| CUSTOMER_NOT_FOUND | 404 | Customer does not exist |
| INVALID_TOKEN | 400 | Verification token is invalid or expired |
| ALREADY_VERIFIED | 409 | Email already verified |

**Preconditions:** [Customer must exist Customer status must be 'unverified']

**Postconditions:** [Customer status changed to 'verified']

#### getCustomer `GET /{customerId}`

Retrieve customer details by ID

**Input:**
- `customerId`: UUID (required)

**Output:**
- `status`: CustomerStatus
- `createdAt`: DateTime
- `verifiedAt`: DateTime
- `customerId`: UUID
- `email`: Email
- `firstName`: string
- `lastName`: string

**Errors:**
| Code | HTTP | Message |
|------|------|----------|
| CUSTOMER_NOT_FOUND | 404 | Customer does not exist |
| UNAUTHORIZED | 401 | Authentication required |
| FORBIDDEN | 403 | Access denied |

**Preconditions:** [Customer must be authenticated Customer can only access own profile]

#### resendVerificationEmail `POST /{customerId}/resend-verification`

Resend verification email to customer

**Input:**
- `customerId`: UUID (required)

**Output:**
- `message`: string
- `sentAt`: DateTime

**Errors:**
| Code | HTTP | Message |
|------|------|----------|
| CUSTOMER_NOT_FOUND | 404 | Customer does not exist |
| ALREADY_VERIFIED | 409 | Email already verified |
| RATE_LIMITED | 429 | Too many requests. Please wait before trying again |

**Preconditions:** [Customer must exist Customer status must be 'unverified']

**Postconditions:** [New verification email sent]

### Events

- **CustomerRegistered**: Emitted when a new customer account is created (payload: [customerId email firstName lastName status])
- **CustomerEmailVerified**: Emitted when customer verifies their email address (payload: [customerId email verifiedAt])
- **VerificationEmailSent**: Emitted when a verification email is sent (payload: [customerId email sentAt])

---

