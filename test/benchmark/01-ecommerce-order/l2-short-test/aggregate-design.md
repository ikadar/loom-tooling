# Aggregate Design

Generated: 2025-12-31T10:31:00+01:00

---

## AGG-CUST-001 – Customer

**Purpose:** Manages customer identity, authentication credentials, and account lifecycle

### Aggregate Root: Customer

**Identity:** CustomerId (UUID)

**Attributes:**
| Name | Type | Mutable |
|------|------|----------|
| email | Email | false |
| passwordHash | PasswordHash | true |
| firstName | string | true |
| lastName | string | true |
| status | CustomerStatus | true |
| createdAt | DateTime | false |
| verifiedAt | DateTime | true |

### Invariants

- **INV-CUST-001**: Email must be unique across all customers
  - Enforcement: Checked via domain service before registration; unique constraint in database
- **INV-CUST-002**: Password must be at least 8 characters
  - Enforcement: Validated in Password value object constructor
- **INV-CUST-003**: Password must contain at least one number
  - Enforcement: Validated in Password value object constructor
- **INV-CUST-004**: Email must be valid format
  - Enforcement: Validated in Email value object constructor
- **INV-CUST-005**: Status transitions must follow valid state machine
  - Enforcement: Guard clauses in status change methods
- **INV-CUST-006**: First name and last name are required
  - Enforcement: Validated on creation and update

### Value Objects

[Email Password PasswordHash CustomerStatus FullName]

### Behaviors

| Command | Pre | Post | Emits |
|---------|-----|------|-------|
| RegisterCustomer | [email not already registered] | [customer created with status 'unverified' verification token generated] | CustomerRegistered |
| VerifyCustomerEmail | [status == 'unverified' valid verification token] | [status == 'verified' verifiedAt set] | CustomerEmailVerified |
| ChangeCustomerPassword | [status != 'suspended' new password meets requirements] | [passwordHash updated] | CustomerPasswordChanged |
| UpdateCustomerProfile | [status != 'suspended'] | [firstName and/or lastName updated] | CustomerProfileUpdated |
| SuspendCustomer | [status != 'suspended'] | [status == 'suspended'] | CustomerSuspended |
| ReactivateCustomer | [status == 'suspended'] | [status == 'verified'] | CustomerReactivated |

### Events

- **CustomerRegistered**: [customerId email firstName lastName status]
- **CustomerEmailVerified**: [customerId email verifiedAt]
- **CustomerPasswordChanged**: [customerId changedAt]
- **CustomerProfileUpdated**: [customerId firstName lastName updatedAt]
- **CustomerSuspended**: [customerId suspendedAt reason]
- **CustomerReactivated**: [customerId reactivatedAt]

### Repository: CustomerRepository

- Load Strategy: Load complete aggregate (single entity, no children)
- Concurrency: Optimistic locking via version field

---

