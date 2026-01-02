# Acceptance Criteria

Generated: 2025-12-29T13:37:22+01:00

---

## AC-CUST-001 – Customer Registration

**Given** a visitor with email 'john@example.com' that is not registered
**When** they submit registration with email 'john@example.com', password 'SecurePass1', first name 'John', and last name 'Doe'
**Then** a customer account is created with status 'unverified', and a verification email is sent to 'john@example.com'

**Error Cases:**
- Email already registered → DUPLICATE_EMAIL error, suggest login or password reset
- Password less than 8 characters → INVALID_PASSWORD error with requirements message
- Password without number → INVALID_PASSWORD error with requirements message
- Invalid email format → INVALID_EMAIL error

**Traceability:**
- Source: Customer entity
- Source: register operation
- Decision: AMB-ENT-024
- Decision: AMB-ENT-025
- Decision: AMB-ENT-026
- Decision: AMB-ENT-030

---

## AC-CUST-002 – Email Verification Required for Orders

**Given** a registered customer with unverified email
**When** they attempt to place an order
**Then** the order is rejected with EMAIL_NOT_VERIFIED error and prompt to verify email

**Error Cases:**
- Verification link expired → VERIFICATION_EXPIRED error with resend option

**Traceability:**
- Source: Place Order operation
- Decision: AMB-ENT-030

---

## AC-CUST-003 – Save Multiple Shipping Addresses

**Given** a registered customer with one existing address
**When** they add a new shipping address with street '456 Oak Ave', city 'Boston', state 'MA', postal code '02101', country 'US', recipient 'John Doe'
**Then** the address is saved to their account, and they can select it for future orders

**Error Cases:**
- Missing required field → INVALID_ADDRESS error specifying missing field
- Invalid postal code format → INVALID_POSTAL_CODE error

**Traceability:**
- Source: ShippingAddress entity
- Source: Customer entity
- Decision: AMB-ENT-027
- Decision: AMB-ENT-031
- Decision: AMB-ENT-032
- Decision: AMB-ENT-033

---

## AC-CUST-004 – GDPR Data Erasure

**Given** a registered customer requesting account deletion
**When** they confirm full data erasure request
**Then** all personal data is permanently deleted, order history is anonymized, and confirmation is sent

**Error Cases:**
- Pending orders exist → PENDING_ORDERS_EXIST error, must wait for completion

**Traceability:**
- Source: Customer entity
- Decision: AMB-ENT-029

---

## AC-PROD-001 – Browse Products with Filters

**Given** a product catalog with products in multiple categories and price ranges
**When** a customer filters by category 'Electronics', price range $50-$200, and availability 'in stock'
**Then** only products matching all criteria are displayed, sorted by newest first, with 20 products per page

**Error Cases:**
- No products match filters → Empty result with 'No products found' message and clear filters option

**Traceability:**
- Source: Browse Products operation
- Decision: AMB-OP-001
- Decision: AMB-OP-002
- Decision: AMB-OP-003

---

## AC-PROD-002 – View Out-of-Stock Products

**Given** a product 'Widget Pro' with zero available stock
**When** a customer views the product listing
**Then** the product is displayed with 'Out of Stock' indicator and 'Add to Cart' button is disabled
