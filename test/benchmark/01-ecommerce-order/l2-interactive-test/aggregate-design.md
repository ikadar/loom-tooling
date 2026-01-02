# Aggregate Design

Generated: 2026-01-02T11:02:21+01:00

---

## AGG-CUSTOMER-001 – Customer {#agg-customer-001}

**Purpose:** Manages customer account lifecycle including registration, verification, and account status

### Aggregate Root: Customer

**Identity:** CustomerId (UUID)

**Attributes:**
| Name | Type | Mutable |
|------|------|----------|
| email | Email | true |
| password | Password | true |
| firstName | string | true |
| lastName | string | true |
| status | CustomerStatus | true |
| createdAt | DateTime | false |
| verifiedAt | DateTime? | true |

### Invariants

- **INV-CUST-001**: Password must be at least 8 characters and contain at least one number
  - Enforcement: Validated on password creation and change via Password value object
- **INV-CUST-002**: Email address must be unique across all customer accounts
  - Enforcement: Enforced at repository level with unique constraint, checked before save
- **INV-CUST-003**: Email must be valid format
  - Enforcement: Validated via Email value object on creation and update
- **INV-CUST-004**: Status transitions must follow valid state machine (unverified→verified→suspended, unverified→suspended)
  - Enforcement: Guard clauses in status transition methods
- **INV-CUST-005**: Customer must have first name and last name
  - Enforcement: Validated on creation, non-empty string check

### Value Objects

[Email Password CustomerStatus]

### Behaviors

| Command | Pre | Post | Emits |
|---------|-----|------|-------|
| RegisterCustomer | [email is unique (checked via repository)] | [status == 'unverified' customer persisted] | CustomerRegistered |
| VerifyCustomer | [status == 'unverified'] | [status == 'verified' verifiedAt is set] | CustomerVerified |
| SuspendCustomer | [status in ['unverified', 'verified']] | [status == 'suspended'] | CustomerSuspended |
| ChangeCustomerEmail | [status != 'suspended' new email is unique] | [email updated status == 'unverified'] | CustomerEmailChanged |
| ChangeCustomerPassword | [status != 'suspended'] | [password updated] | CustomerPasswordChanged |
| UpdateCustomerProfile | [status != 'suspended'] | [firstName and/or lastName updated] | CustomerProfileUpdated |

### Events

- **CustomerRegistered**: [customerId email firstName lastName registeredAt]
- **CustomerVerified**: [customerId verifiedAt]
- **CustomerSuspended**: [customerId suspendedAt reason]
- **CustomerEmailChanged**: [customerId oldEmail newEmail changedAt]
- **CustomerPasswordChanged**: [customerId changedAt]
- **CustomerProfileUpdated**: [customerId firstName lastName updatedAt]

### Repository: CustomerRepository

- Load Strategy: Load complete aggregate (single entity, no children)
- Concurrency: Optimistic locking via version field

---

