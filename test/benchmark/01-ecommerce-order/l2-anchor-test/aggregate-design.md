# Aggregate Design

Generated: 2025-12-31T13:24:23+01:00

---

## AGG-CUSTOMER-001 – Customer {#agg-customer-001}

**Purpose:** Manages customer identity, authentication credentials, and account lifecycle

### Aggregate Root: Customer

**Identity:** CustomerId (UUID)

**Attributes:**
| Name | Type | Mutable |
|------|------|----------|
| email | Email | true |
| passwordHash | PasswordHash | true |
| firstName | string | true |
| lastName | string | true |
| status | CustomerStatus | true |
| createdAt | DateTime | false |
| verifiedAt | DateTime? | true |

### Invariants

- **INV-CUST-001**: Password must be at least 8 characters and contain at least one number
  - Enforcement: Validated on password creation and update via Password value object
- **INV-CUST-002**: Email address must be unique across all customers
  - Enforcement: Enforced via unique constraint and domain service check before registration
- **INV-CUST-003**: Email must be valid format
  - Enforcement: Validated via Email value object on creation and update
- **INV-CUST-004**: Status transitions must follow valid state machine (unverified→verified→suspended, unverified→suspended)
  - Enforcement: Guard clauses in status transition methods
- **INV-CUST-005**: Customer must have non-empty first and last name
  - Enforcement: Validated on creation and profile update

### Value Objects

[Email PasswordHash CustomerStatus PersonName]

### Behaviors

| Command | Pre | Post | Emits |
|---------|-----|------|-------|
| RegisterCustomer | [email is unique (checked via domain service)] | [status == 'unverified' customer persisted] | CustomerRegistered |
| VerifyCustomer | [status == 'unverified'] | [status == 'verified' verifiedAt is set] | CustomerVerified |
| SuspendCustomer | [status in ['unverified', 'verified']] | [status == 'suspended'] | CustomerSuspended |
| ReactivateCustomer | [status == 'suspended'] | [status == 'verified'] | CustomerReactivated |
| ChangeCustomerEmail | [status != 'suspended' new email is unique] | [email updated status == 'unverified' (re-verification required)] | CustomerEmailChanged |
| ChangeCustomerPassword | [status != 'suspended' new password meets requirements] | [passwordHash updated] | CustomerPasswordChanged |
| UpdateCustomerProfile | [status != 'suspended'] | [firstName and/or lastName updated] | CustomerProfileUpdated |

### Events

- **CustomerRegistered**: [customerId email firstName lastName registeredAt]
- **CustomerVerified**: [customerId verifiedAt]
- **CustomerSuspended**: [customerId suspendedAt reason]
- **CustomerReactivated**: [customerId reactivatedAt]
- **CustomerEmailChanged**: [customerId oldEmail newEmail changedAt]
- **CustomerPasswordChanged**: [customerId changedAt]
- **CustomerProfileUpdated**: [customerId firstName lastName updatedAt]

### Repository: CustomerRepository

- Load Strategy: Load complete aggregate (single entity)
- Concurrency: Optimistic locking via version field

---

