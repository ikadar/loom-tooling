# Business Rules

Generated: 2025-12-23T21:44:26+01:00

---

## BR-QUOTE-001 – Quote reference number format

**Rule:** Every quote must have a unique reference number in the format QT-XXXXX where XXXXX is a sequential number

**Invariant:** Quote.referenceNumber MUST match pattern 'QT-[0-9]{5,}' AND MUST be unique across all quotes

**Enforcement:** Generated automatically by system on quote creation; uniqueness enforced by database constraint

**Error Code:** `ERROR_DUPLICATE_REFERENCE`

**Traceability:**
- Source: US-QUOTE-003

---

## BR-QUOTE-002 – Quote expiration period limits

**Rule:** Quote expiration can be set between 7 and 90 days, with a default of 30 days

**Invariant:** Quote.expirationDays MUST be >= 7 AND <= 90

**Enforcement:** Validated on quote creation and modification; UI enforces with date picker constraints

**Error Code:** `ERROR_INVALID_EXPIRATION`

**Traceability:**
- Source: US-QUOTE-003
- Decision: AMB-ENT-011

---

## BR-QUOTE-003 – Quote minimum value

**Rule:** Total quote value must be at least $1.00

**Invariant:** Quote.totalValue MUST be >= 1.00

**Enforcement:** Validated when transitioning quote from Draft to Pending status

**Error Code:** `ERROR_MINIMUM_VALUE`

**Traceability:**
- Source: US-QUOTE-003
- Decision: AMB-ENT-010

---

## BR-QUOTE-004 – High-value quote approval

**Rule:** Quotes with total value exceeding $100,000 require manager approval before sending to customer

**Invariant:** IF Quote.totalValue > 100000 THEN Quote.managerApproval MUST be 'Approved' before status can change to 'Pending'

**Enforcement:** Enforced on status transition from Draft to Pending; workflow triggers approval request

**Error Code:** `ERROR_APPROVAL_REQUIRED`

**Traceability:**
- Source: US-QUOTE-003
- Decision: AMB-ENT-010

---

## BR-QUOTE-005 – Quote status transitions

**Rule:** Quote status can only transition through defined paths: Draft→Pending, Pending→Accepted/Rejected/Expired/Cancelled, Any→Cancelled (admin only)

**Invariant:** Quote.status transitions MUST follow the defined state machine

**Enforcement:** Enforced by status update service; invalid transitions rejected

**Error Code:** `ERROR_INVALID_TRANSITION`

**Traceability:**
- Source: US-QUOTE-003
- Decision: AMB-ENT-001
- Decision: AMB-OP-028

---

## BR-QUOTE-006 – Quote status transition authorization

**Rule:** Status transitions are role-restricted: Customer can transition Pending→Accepted/Rejected; Sales Rep can transition Draft→Pending and Pending→Cancelled; System transitions Pending→Expired; Admin can perform any transition

**Invariant:** Actor MUST have appropriate role for the requested status transition

**Enforcement:** Enforced by authorization layer before status update

**Error Code:** `ERROR_UNAUTHORIZED_TRANSITION`

**Traceability:**
- Source: US-QUOTE-003
- Decision: AMB-OP-029

---

## BR-QUOTE-007 – Quote modification restrictions

**Rule:** Quotes can only be modified when in Draft or Pending status by the assigned sales rep or admin

**Invariant:** Quote modification MUST only occur when Quote.status IN ('Draft', 'Pending') AND (user = Quote.assignedSalesRep OR user.role = 'Admin')

**Enforcement:** Enforced on all quote update operations

**Error Code:** `ERROR_QUOTE_NOT_MODIFIABLE`

**Traceability:**
- Source: US-QUOTE-003
- Decision: AMB-ENT-003

---

## BR-QUOTE-008 – Quote deletion restrictions

**Rule:** Only quotes in Draft or Pending status can be soft-deleted; Accepted quotes cannot be deleted

**Invariant:** Quote.delete() MUST only be called when Quote.status IN ('Draft', 'Pending')

**Enforcement:** Enforced on delete operation

**Error Code:** `ERROR_CANNOT_DELETE`

**Traceability:**
- Source: US-QUOTE-003
- Decision: AMB-ENT-004

---

## BR-QUOTE-009 – Line item quantity constraints

**Rule:** Line item quantity must be positive, minimum 1, maximum 99999; decimals allowed only for products configured to allow them

**Invariant:** LineItem.quantity MUST be > 0 AND <= 99999 AND (Product.allowDecimalQuantity = true OR LineItem.quantity is integer)

**Enforcement:** Validated on line item creation and modification

**Error Code:** `ERROR_INVALID_QUANTITY`

**Traceability:**
- Source: US-QUOTE-003
- Decision: AMB-ENT-020

---

## BR-QUOTE-010 – Line item price locking

**Rule:** Line item price is captured and locked at the time of adding to quote; subsequent product price changes do not affect existing quote line items

**Invariant:** LineItem.price MUST NOT change due to Product.price changes after LineItem creation

**Enforcement:** Price copied from product catalog at line item creation; no automatic price sync

**Error Code:** `N/A`

**Traceability:**
- Source: US-QUOTE-003
- Decision: AMB-ENT-021

---

## BR-QUOTE-011 – Quote line item limits

**Rule:** Quotes have a soft limit of 100 line items (warning) and hard limit of 500 line items (blocked)

**Invariant:** Quote.lineItems.count MUST be <= 500

**Enforcement:** Warning shown at 100+ items; creation blocked at 500+ items

**Error Code:** `ERROR_TOO_MANY_ITEMS`

**Traceability:**
- Source: US-QUOTE-003
- Decision: AMB-OP-041

---

## BR-QUOTE-012 – Discontinued product restrictions

**Rule:** Discontinued products cannot be added to new quotes; existing quotes with discontinued products show warning but can be accepted

**Invariant:** New LineItem MUST NOT reference Product where Product.status = 'Discontinued'

**Enforcement:** Validated on line item addition; warning displayed for existing quotes during view

**Error Code:** `ERROR_PRODUCT_DISCONTINUED`

**Traceability:**
- Source: US-QUOTE-003
- Decision: AMB-ENT-038

---

## BR-QUOTE-013 – Quote resend limits

**Rule:** Quote email can be resent maximum 5 times with minimum 1 hour between resends; each resend generates new unique link and invalidates previous

**Invariant:** Quote.emailSendCount MUST be <= 5 AND timeSinceLastSend MUST be >= 1 hour

**Enforcement:** Enforced on resend operation

**Error Code:** `ERROR_RESEND_LIMIT`

**Traceability:**
- Source: US-QUOTE-003
- Decision: AMB-OP-003

---

## BR-QUOTE-014 – Unique link validity

**Rule:** Quote unique link (UUID-based token) is valid only until quote expiration date or until a new link is generated

**Invariant:** QuoteLink.token MUST be invalidated when Quote.expirationDate is passed OR new link is generated

**Enforcement:** Validated on every link access attempt

**Error Code:** `ERROR_LINK_INVALID`

**Traceability:**
- Source: US-QUOTE-003
- Decision: AMB-ENT-005

---

## BR-QUOTE-015 – Quote acceptance terms requirement

**Rule:** Quote acceptance requires explicit acknowledgment of terms and conditions via checkbox

**Invariant:** QuoteAcceptance MUST include termsAccepted = true AND termsVersion reference

**Enforcement:** Accept button disabled until checkbox checked; validated on server side

**Error Code:** `ERROR_TERMS_NOT_ACCEPTED`

**Traceability:**
- Source: US-QUOTE-003
- Decision: AMB-OP-013
- Decision: AMB-ENT-032

---

## BR-QUOTE-016 – Quote acceptance preconditions

**Rule:** Quote can only be accepted if status is Pending AND not expired AND quote version hasn't changed since page load

**Invariant:** AcceptQuote MUST only succeed when Quote.status = 'Pending' AND Quote.expirationDate > now() AND Quote.version = submittedVersion

**Enforcement:** Validated atomically during acceptance transaction

**Error Code:** `ERROR_QUOTE_NOT_ACCEPTABLE`

**Traceability:**
- Source: US-QUOTE-003
- Decision: AMB-OP-014

---

## BR-QUOTE-017 – Quote acceptance idempotency

**Rule:** Multiple acceptance attempts for the same quote return the same successful result without creating duplicate records

**Invariant:** IF Quote.status = 'Accepted' THEN AcceptQuote returns existing QuoteAcceptance without modification

**Enforcement:** Checked at start of acceptance operation; return cached result if already accepted

**Error Code:** `N/A`

**Traceability:**
- Source: US-QUOTE-003
- Decision: AMB-OP-015

---

## BR-QUOTE-018 – Acceptance record immutability

**Rule:** QuoteAcceptance records are immutable once created; corrections must be made via new records referencing the original

**Invariant:** QuoteAcceptance records MUST NOT be updated or deleted after creation

**Enforcement:** No update/delete endpoints exposed; database-level constraints

**Error Code:** `ERROR_RECORD_IMMUTABLE`

**Traceability:**
- Source: US-QUOTE-003
- Decision: AMB-OP-022

---

## BR-QUOTE-019 – Acceptance record retention

**Rule:** QuoteAcceptance records must be retained for minimum 7 years for legal compliance

**Invariant:** QuoteAcceptance records MUST NOT be deleted before 7 years from creation date

**Enforcement:** Data retention policy; no hard delete capability; archival process preserves records

**Error Code:** `N/A`

**Traceability:**
- Source: US-QUOTE-003
- Decision: AMB-OP-023

---

## BR-QUOTE-020 – Quote to order relationship

**Rule:** One accepted quote creates exactly one order; order references source quote

**Invariant:** Quote:Order relationship MUST be 1:1; Order.sourceQuoteId MUST reference valid Quote with status 'Accepted'

**Enforcement:** Enforced during order creation; database constraint on unique quote reference

**Error Code:** `ERROR_ORDER_EXISTS`

**Traceability:**
- Source: US-QUOTE-003
- Decision: AMB-ENT-023

---

## BR-QUOTE-021 – Order number format

**Rule:** Every order must have a unique reference number in the format ORD-XXXXX where XXXXX is a sequential number

**Invariant:** Order.orderNumber MUST match pattern 'ORD-[0-9]{5,}' AND MUST be unique across all orders

**Enforcement:** Generated automatically by system on order creation; uniqueness enforced by database constraint

**Error Code:** `ERROR_DUPLICATE_ORDER`

**Traceability:**
- Source: US-QUOTE-003
- Decision: AMB-ENT-026

---

## BR-QUOTE-022 – Customer quote eligibility

**Rule:** Only customers with Active status can receive and accept quotes

**Invariant:** Quote.customer.status MUST be 'Active' for SendQuoteEmail and AcceptQuote operations

**Enforcement:** Validated when transitioning quote to Pending and during acceptance

**Error Code:** `ERROR_CUSTOMER_INACTIVE`

**Traceability:**
- Source: US-QUOTE-003
- Decision: AMB-ENT-015

---

## BR-QUOTE-023 – Price discount limits

**Rule:** Sales representatives can override product prices within configured discount limits (default max 20%); larger discounts require manager approval

**Invariant:** IF LineItem.discount > configuredMaxDiscount THEN Quote.managerApproval MUST be 'Approved' before status can change to 'Pending'

**Enforcement:** Validated on quote submission for sending; workflow triggers approval for excessive discounts

**Error Code:** `ERROR_DISCOUNT_LIMIT`

**Traceability:**
- Source: US-QUOTE-003
- Decision: AMB-OP-038

---

## BR-QUOTE-024 – Quote value calculation

**Rule:** Quote total value is computed from line items and stored with 2 decimal precision in a single currency per quote

**Invariant:** Quote.totalValue MUST equal SUM(LineItem.lineTotal) - Quote.discountAmount + Quote.taxAmount with 2 decimal precision

**Enforcement:** Calculated on every line item change; stored value recalculated on quote read if stale

**Error Code:** `N/A`

**Traceability:**
- Source: US-QUOTE-003
- Decision: AMB-ENT-007

---

## BR-QUOTE-025 – Quote requires at least one line item

**Rule:** A quote must have at least one line item to be sent to customer

**Invariant:** Quote.lineItems.count MUST be >= 1 when transitioning from Draft to Pending

**Enforcement:** Validated on status transition

**Error Code:** `ERROR_NO_LINE_ITEMS`

**Traceability:**
- Source: US-QUOTE-003
- Decision: AMB-OP-036

---

## BR-QUOTE-026 – Acceptance page session timeout

**Rule:** Quote acceptance page has a 30-minute session timeout; expired sessions require re-navigation via original link

**Invariant:** Acceptance submission MUST occur within 30 minutes of page load

**Enforcement:** Client-side timeout warning; server-side session validation

**Error Code:** `ERROR_SESSION_EXPIRED`

**Traceability:**
- Source: US-QUOTE-003
- Decision: AMB-OP-018

---

## BR-QUOTE-027 – Email delivery retry policy

**Rule:** Failed quote emails are automatically retried 3 times over 24 hours; sales rep is notified of persistent failures

**Invariant:** Email.retryCount MUST be <= 3 AND retries MUST be spread over 24-hour period

**Enforcement:** Email service retry logic; notification service for failure alerts

**Error Code:** `N/A`

**Traceability:**
- Source: US-QUOTE-003
- Decision: AMB-ENT-035

---

## BR-QUOTE-028 – Audit trail completeness

**Rule:** All quote changes must be tracked in audit log including: status changes, line item modifications, price changes, actor, and timestamp

**Invariant:** Every Quote modification MUST create corresponding AuditLog entry with changeType, actor, timestamp, and before/after values

**Enforcement:** Audit interceptor on all quote update operations

**Error Code:** `N/A`

**Traceability:**
- Source: US-QUOTE-003
- Decision: AMB-ENT-006

---

## BR-QUOTE-029 – Quote acceptance cannot be revoked

**Rule:** Once a quote is accepted by the customer, the acceptance cannot be revoked; disputes are handled through support workflow

**Invariant:** QuoteAcceptance record MUST NOT be deleted or marked as revoked

**Enforcement:** No revocation endpoint; order cancellation handled separately

**Error Code:** `ERROR_CANNOT_REVOKE`

**Traceability:**
- Source: US-QUOTE-003
- Decision: AMB-ENT-033

---

## BR-QUOTE-030 – Tax calculation based on customer location

**Rule:** Taxes are calculated automatically based on customer location (billing address) and product tax categories

**Invariant:** Quote.taxAmount MUST be calculated using Customer.billingAddress jurisdiction rules AND Product.taxCategory for each line item

**Enforcement:** Automatic recalculation on customer selection and line item changes

**Error Code:** `ERROR_TAX_CALCULATION`

**Traceability:**
- Source: US-QUOTE-003
- Decision: AMB-ENT-009

---

