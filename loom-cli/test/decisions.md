
## Decisions from 2025-12-23

- **AMB-ENT-001: Quote**
  - Q: What additional states can a Quote have beyond 'Pending' and 'Accepted'? (e.g., Draft, Expired, Rejected, Cancelled, Converted to Order)
  - A: Quote states: Draft, Pending, Accepted, Rejected, Expired, Cancelled
  - Decided: 2025-12-23 21:41 by default

- **AMB-ENT-002: Quote**
  - Q: What happens to a Quote when it expires? Is the status automatically updated? Can an expired quote be re-activated?
  - A: System automatically marks quotes as Expired after expiration date; expired quotes cannot be accepted but can be duplicated into a new quote
  - Decided: 2025-12-23 21:41 by default

- **AMB-ENT-003: Quote**
  - Q: Can a Quote be modified after creation? If so, by whom and under what conditions? Does modification reset the expiration date?
  - A: Only sales representative can modify quotes in Draft or Pending status; modification resets expiration date and requires re-sending to customer
  - Decided: 2025-12-23 21:41 by default

- **AMB-ENT-004: Quote**
  - Q: Can a Quote be deleted? If so, by whom and under what conditions? Is it soft delete or hard delete?
  - A: Soft delete only; sales rep or admin can delete quotes in Draft or Pending status; Accepted quotes cannot be deleted
  - Decided: 2025-12-23 21:41 by default

- **AMB-ENT-005: Quote**
  - Q: What is the format and generation rule for the unique link? How long should it remain valid?
  - A: UUID-based token appended to base URL; link validity tied to quote expiration date
  - Decided: 2025-12-23 21:41 by default

- **AMB-ENT-006: Quote**
  - Q: Is audit trail/change history required for Quote modifications? What changes should be tracked?
  - A: Yes, track all changes including: status changes, line item modifications, price changes, who made changes and when
  - Decided: 2025-12-23 21:41 by default

- **AMB-ENT-007: Quote**
  - Q: How is the total value calculated? Is it stored or computed? What precision/currency handling is required?
  - A: Computed from line items; stored for performance; 2 decimal precision; single currency per quote
  - Decided: 2025-12-23 21:41 by default

- **AMB-ENT-008: Quote**
  - Q: What discount types are supported? (percentage, fixed amount, per-line-item, whole-quote)
  - A: Support both percentage and fixed amount discounts; can be applied per line item or to entire quote total
  - Decided: 2025-12-23 21:41 by default

- **AMB-ENT-009: Quote**
  - Q: What tax calculation rules apply? Are taxes calculated automatically or entered manually?
  - A: Taxes calculated automatically based on customer location and product tax categories
  - Decided: 2025-12-23 21:41 by default

- **AMB-ENT-010: Quote**
  - Q: Is there a minimum or maximum value constraint for a Quote?
  - A: Minimum quote value of $1; no maximum but quotes over $100,000 require manager approval
  - Decided: 2025-12-23 21:41 by default

- **AMB-ENT-011: Quote**
  - Q: Can the default 30-day expiration be customized? By whom? What are the min/max limits?
  - A: Sales rep can customize expiration between 7-90 days; default is 30 days
  - Decided: 2025-12-23 21:41 by default

- **AMB-ENT-012: Customer**
  - Q: What uniquely identifies a Customer? Email only? Or is there a separate customer ID?
  - A: Customers have unique system-generated ID; email is unique but can be updated; supports multiple contacts per customer organization
  - Decided: 2025-12-23 21:41 by default

- **AMB-ENT-013: Customer**
  - Q: What additional attributes does a Customer have? (name, company, address, phone, etc.)
  - A: Required: name, email, company name. Optional: phone, billing address, shipping address, tax ID
  - Decided: 2025-12-23 21:41 by default

- **AMB-ENT-014: Customer**
  - Q: Does the Customer entity represent a person or a company? Can a company have multiple contacts?
  - A: Customer represents a company/organization; separate Contact entity for individuals; one Customer can have multiple Contacts
  - Decided: 2025-12-23 21:41 by default

- **AMB-ENT-015: Customer**
  - Q: What states can a Customer have? (active, inactive, blocked, etc.)
  - A: Customer states: Active, Inactive, Blocked; only Active customers can receive and accept quotes
  - Decided: 2025-12-23 21:41 by default

- **AMB-ENT-016: Customer**
  - Q: Can Customer data be deleted? What happens to associated quotes and orders?
  - A: Soft delete only; customer marked inactive; historical quotes and orders retained for compliance; PII can be anonymized on request
  - Decided: 2025-12-23 21:41 by default

- **AMB-ENT-017: Customer**
  - Q: Is the IP address stored on Customer permanent or only captured during quote acceptance?
  - A: IP address is captured per-action (quote acceptance) and stored in QuoteAcceptance, not as a Customer attribute
  - Decided: 2025-12-23 21:41 by default

- **AMB-ENT-018: LineItem**
  - Q: Does LineItem have its own unique identifier? Or is it identified by Quote + sequence number?
  - A: LineItem has unique system ID; also has sequence number within quote for ordering
  - Decided: 2025-12-23 21:41 by default

- **AMB-ENT-019: LineItem**
  - Q: What additional attributes does LineItem need? (description, unit price vs total price, discount, tax, notes)
  - A: LineItem attributes: product reference, description, quantity, unit price, line discount, line total, tax amount, optional notes
  - Decided: 2025-12-23 21:41 by default

- **AMB-ENT-020: LineItem**
  - Q: What validation rules apply to quantity? (minimum, maximum, decimal allowed, units)
  - A: Quantity must be positive; minimum 1; decimals allowed for certain product types; maximum 99999
  - Decided: 2025-12-23 21:41 by default

- **AMB-ENT-021: LineItem**
  - Q: Is the price on LineItem locked at quote creation, or does it reference current product price?
  - A: Price is captured and locked at time of adding to quote; does not change if product price changes
  - Decided: 2025-12-23 21:41 by default

- **AMB-ENT-022: LineItem**
  - Q: Can LineItems be added/removed/modified after quote creation? Under what conditions?
  - A: LineItems can be modified while quote is in Draft or Pending status; changes require quote to be re-sent to customer
  - Decided: 2025-12-23 21:41 by default

- **AMB-ENT-023: Order**
  - Q: What is the relationship between Quote and Order? One-to-one? Can multiple orders be created from one quote?
  - A: One-to-one relationship; one accepted quote creates exactly one order; order references source quote
  - Decided: 2025-12-23 21:41 by default

- **AMB-ENT-024: Order**
  - Q: What attributes does Order inherit from Quote vs. have independently?
  - A: Order copies all quote data (customer, line items, prices, totals) at creation; adds order-specific fields (order number, order date, fulfillment status)
  - Decided: 2025-12-23 21:41 by default

- **AMB-ENT-025: Order**
  - Q: What states can an Order have? What is the order fulfillment workflow?
  - A: Order states: Created, Confirmed, Processing, Shipped, Delivered, Completed, Cancelled
  - Decided: 2025-12-23 21:41 by default

- **AMB-ENT-026: Order**
  - Q: What is the format for Order ID/number? Is it related to Quote reference number?
  - A: Order number format: ORD-XXXXX (sequential); includes reference to source quote number
  - Decided: 2025-12-23 21:41 by default

- **AMB-ENT-027: SalesRepresentative**
  - Q: What attributes does SalesRepresentative have? (name, email, employee ID, team, manager)
  - A: SalesRep attributes: employee ID, name, email, team/department, manager reference, active status
  - Decided: 2025-12-23 21:41 by default

- **AMB-ENT-028: SalesRepresentative**
  - Q: What permissions/authorization does a SalesRepresentative have? Can they see all quotes or only their own?
  - A: Sales reps see own quotes by default; team leads see team quotes; managers see all quotes in their org
  - Decided: 2025-12-23 21:41 by default

- **AMB-ENT-029: SalesRepresentative**
  - Q: What happens to quotes when a SalesRepresentative leaves? Reassignment rules?
  - A: When sales rep is deactivated, pending quotes are automatically reassigned to their manager; accepted quotes retain original rep for historical purposes
  - Decided: 2025-12-23 21:41 by default

- **AMB-ENT-030: SalesRepresentative**
  - Q: Is SalesRepresentative a separate entity or part of a general User/Employee system?
  - A: SalesRepresentative is a role within a general User system; users can have multiple roles
  - Decided: 2025-12-23 21:41 by default

- **AMB-ENT-031: QuoteAcceptance**
  - Q: Is QuoteAcceptance a separate entity or just attributes on Quote? What additional data should it capture?
  - A: Separate entity for audit purposes; captures: acceptance timestamp, IP address, user agent, geolocation (optional), acceptance method
  - Decided: 2025-12-23 21:41 by default

- **AMB-ENT-032: QuoteAcceptance**
  - Q: Should QuoteAcceptance capture digital signature or legal consent acknowledgment?
  - A: Capture explicit consent checkbox acknowledgment and terms version accepted; digital signature optional based on quote value
  - Decided: 2025-12-23 21:41 by default

- **AMB-ENT-033: QuoteAcceptance**
  - Q: Can a quote acceptance be revoked or disputed? What is the process?
  - A: Acceptance cannot be revoked by customer; disputes handled through support workflow; order can be cancelled separately
  - Decided: 2025-12-23 21:41 by default

- **AMB-ENT-034: Email**
  - Q: Is Email a separate entity that should be tracked, or just a transient notification?
  - A: Email is tracked entity for audit; stores: type, recipient, sent timestamp, delivery status, related quote/order
  - Decided: 2025-12-23 21:41 by default

- **AMB-ENT-035: Email**
  - Q: What happens if email delivery fails? Retry logic? Alternative notification?
  - A: Automatic retry 3 times over 24 hours; notify sales rep of failure; provide manual resend option
  - Decided: 2025-12-23 21:41 by default

- **AMB-ENT-036: Email**
  - Q: Should email open/click tracking be implemented for the quote link?
  - A: Track email opens and link clicks for sales analytics; data visible to sales rep on quote detail
  - Decided: 2025-12-23 21:41 by default

- **AMB-ENT-037: Product**
  - Q: Product is referenced by LineItem but not defined. What attributes does Product have?
  - A: Product entity with: SKU (identifier), name, description, unit price, tax category, active status, category
  - Decided: 2025-12-23 21:41 by default

- **AMB-ENT-038: Product**
  - Q: Can quotes include products that are discontinued or out of stock?
  - A: Cannot add discontinued products to new quotes; existing quotes with discontinued products show warning but can still be accepted
  - Decided: 2025-12-23 21:41 by default

- **AMB-OP-001: SendQuoteEmail**
  - Q: Should the email be sent immediately upon quote creation, or can it be scheduled or manually triggered?
  - A: Email is sent automatically when quote status changes to 'Pending' (ready to send); sales rep can also manually trigger resend
  - Decided: 2025-12-23 21:41 by default

- **AMB-OP-002: SendQuoteEmail**
  - Q: What should happen if the customer email address is invalid or bounces?
  - A: Validate email format before sending; on bounce, mark quote with delivery failure flag and notify sales rep; allow email address correction and resend
  - Decided: 2025-12-23 21:41 by default

- **AMB-OP-003: SendQuoteEmail**
  - Q: Can the same quote email be sent multiple times? What are the limits?
  - A: Allow resend up to 5 times; each resend generates new unique link and invalidates previous; minimum 1 hour between resends
  - Decided: 2025-12-23 21:41 by default

- **AMB-OP-004: SendQuoteEmail**
  - Q: What email template/content should be used? Is it customizable per sales rep or customer?
  - A: Standard template with company branding; sales rep can add personalized message; template includes quote summary, expiration date, and accept link
  - Decided: 2025-12-23 21:41 by default

- **AMB-OP-005: SendQuoteEmail**
  - Q: Should email delivery status be tracked and what is the expected delivery SLA?
  - A: Track delivery status (sent, delivered, bounced, opened); target delivery within 5 minutes; alert if not delivered within 1 hour
  - Decided: 2025-12-23 21:41 by default

- **AMB-OP-006: SendQuoteEmail**
  - Q: What happens if the email service is unavailable? Should quote creation fail or proceed?
  - A: Quote creation succeeds; email queued for retry; sales rep notified of delay; auto-retry every 15 minutes for up to 24 hours
  - Decided: 2025-12-23 21:41 by default

- **AMB-OP-007: ViewQuoteDetails**
  - Q: What specific quote details should be displayed to the customer?
  - A: Display: quote reference, date, expiration, line items with descriptions/quantities/prices, subtotal, taxes, discounts, total, terms and conditions
  - Decided: 2025-12-23 21:41 by default

- **AMB-OP-008: ViewQuoteDetails**
  - Q: What should be displayed if the quote has expired when the customer clicks the link?
  - A: Show friendly expiration message with quote reference; provide option to request new quote; display contact information for sales rep
  - Decided: 2025-12-23 21:41 by default

- **AMB-OP-009: ViewQuoteDetails**
  - Q: Should quote views be tracked? Is there a limit on how many times a quote can be viewed?
  - A: Track all views with timestamp and IP; no view limit; view analytics visible to sales rep
  - Decided: 2025-12-23 21:41 by default

- **AMB-OP-010: ViewQuoteDetails**
  - Q: Is authentication required to view quote details, or is the unique link sufficient?
  - A: Unique link is sufficient for viewing; no login required; link acts as bearer token
  - Decided: 2025-12-23 21:41 by default

- **AMB-OP-011: ViewQuoteDetails**
  - Q: What is the expected page load time for viewing quote details?
  - A: Page should load within 2 seconds; quote data cached for performance
  - Decided: 2025-12-23 21:41 by default

- **AMB-OP-012: ViewQuoteDetails**
  - Q: Can the customer download or print the quote from the view page?
  - A: Yes, provide PDF download and print-friendly view options
  - Decided: 2025-12-23 21:41 by default

- **AMB-OP-013: AcceptQuote**
  - Q: What confirmation or acknowledgment must the customer provide before acceptance is processed?
  - A: Customer must check 'I agree to terms and conditions' checkbox and click 'Accept Quote' button; terms link opens in new tab
  - Decided: 2025-12-23 21:41 by default

- **AMB-OP-014: AcceptQuote**
  - Q: What happens if the customer clicks accept but the quote was just modified by the sales rep?
  - A: Check quote version/modification timestamp; if changed since page load, show error and refresh quote details; require re-acceptance of updated quote
  - Decided: 2025-12-23 21:41 by default

- **AMB-OP-015: AcceptQuote**
  - Q: What happens if acceptance is attempted multiple times (double-click, retry, etc.)?
  - A: Idempotent operation; second attempt returns success with original acceptance details; no duplicate records created
  - Decided: 2025-12-23 21:41 by default

- **AMB-OP-016: AcceptQuote**
  - Q: What additional data should be captured at acceptance time for legal/audit purposes?
  - A: Capture: IP address, user agent, timestamp, browser fingerprint (optional), terms version accepted, quote version hash
  - Decided: 2025-12-23 21:41 by default

- **AMB-OP-017: AcceptQuote**
  - Q: What error message should be shown if acceptance fails due to technical issues?
  - A: Show friendly error message; preserve entered data; provide retry option; show support contact if retry fails
  - Decided: 2025-12-23 21:41 by default

- **AMB-OP-018: AcceptQuote**
  - Q: Is there a timeout for the acceptance page session? What happens on timeout?
  - A: 30-minute session timeout; on timeout, require re-click of original link; quote data refreshed on reload
  - Decided: 2025-12-23 21:41 by default

- **AMB-OP-019: AcceptQuote**
  - Q: Should sales rep be notified immediately when quote is accepted?
  - A: Yes, send immediate notification via email and in-app notification to assigned sales rep
  - Decided: 2025-12-23 21:41 by default

- **AMB-OP-020: RecordAcceptance**
  - Q: Where should acceptance data be stored? Same database as quote or separate audit log?
  - A: Store in both: QuoteAcceptance table linked to Quote, plus immutable audit log for compliance
  - Decided: 2025-12-23 21:41 by default

- **AMB-OP-021: RecordAcceptance**
  - Q: What happens if recording acceptance fails after customer clicked accept?
  - A: Wrap in transaction; if record fails, rollback status change; show error to customer; auto-retry; alert operations team
  - Decided: 2025-12-23 21:41 by default

- **AMB-OP-022: RecordAcceptance**
  - Q: Should the acceptance record be immutable once created?
  - A: Yes, acceptance records are immutable; corrections made via new records with references to original
  - Decided: 2025-12-23 21:41 by default

- **AMB-OP-023: RecordAcceptance**
  - Q: What is the retention period for acceptance records?
  - A: Retain indefinitely for legal compliance; minimum 7 years per standard business record retention
  - Decided: 2025-12-23 21:41 by default

- **AMB-OP-024: SendConfirmationEmail**
  - Q: What information should be included in the confirmation email?
  - A: Include: order confirmation number, quote reference, accepted items summary, total amount, expected next steps, contact information
  - Decided: 2025-12-23 21:41 by default

- **AMB-OP-025: SendConfirmationEmail**
  - Q: Should confirmation email include payment instructions or next steps?
  - A: Yes, include payment terms, payment methods, and timeline; link to payment portal if applicable
  - Decided: 2025-12-23 21:41 by default

- **AMB-OP-026: SendConfirmationEmail**
  - Q: Should the sales rep receive a copy of the confirmation email?
  - A: Yes, BCC sales rep on confirmation email; also send separate internal notification
  - Decided: 2025-12-23 21:41 by default

- **AMB-OP-027: SendConfirmationEmail**
  - Q: What happens if confirmation email fails to send?
  - A: Queue for retry; do not rollback acceptance; notify sales rep of delivery failure; allow manual resend
  - Decided: 2025-12-23 21:41 by default

- **AMB-OP-028: UpdateQuoteStatus**
  - Q: What other status transitions are valid besides Pending → Accepted?
  - A: Valid transitions: Draft→Pending, Pending→Accepted, Pending→Rejected, Pending→Expired, Pending→Cancelled, any→Cancelled (by admin)
  - Decided: 2025-12-23 21:41 by default

- **AMB-OP-029: UpdateQuoteStatus**
  - Q: Who can change quote status and which transitions can each role perform?
  - A: Customer: Pending→Accepted/Rejected; Sales Rep: Draft→Pending, Pending→Cancelled; System: Pending→Expired; Admin: any transition
  - Decided: 2025-12-23 21:41 by default

- **AMB-OP-030: UpdateQuoteStatus**
  - Q: Should status changes trigger notifications to relevant parties?
  - A: Yes, notify: customer on status change to Accepted/Rejected/Expired; sales rep on all status changes; configurable notification preferences
  - Decided: 2025-12-23 21:41 by default

- **AMB-OP-031: UpdateQuoteStatus**
  - Q: Is status change history maintained? How long?
  - A: Yes, maintain full status change history with timestamp, actor, and reason; retain for life of quote record
  - Decided: 2025-12-23 21:41 by default

- **AMB-OP-032: TriggerOrderCreation**
  - Q: Should order creation be synchronous or asynchronous to quote acceptance?
  - A: Asynchronous via event/message queue; quote acceptance completes immediately; order creation handled by separate process
  - Decided: 2025-12-23 21:41 by default

- **AMB-OP-033: TriggerOrderCreation**
  - Q: What happens if order creation fails? Should quote acceptance be rolled back?
  - A: Do not rollback quote acceptance; mark quote with 'order creation pending' flag; auto-retry order creation; alert operations team after 3 failures
  - Decided: 2025-12-23 21:41 by default

- **AMB-OP-034: TriggerOrderCreation**
  - Q: What data validation should occur before creating the order?
  - A: Validate: all products still active, customer status active, credit check passed (if applicable), inventory available (if tracked)
  - Decided: 2025-12-23 21:41 by default

- **AMB-OP-035: TriggerOrderCreation**
  - Q: Should order creation trigger any integrations with external systems (ERP, fulfillment, etc.)?
  - A: Yes, publish order created event to integration bus; downstream systems (ERP, fulfillment, billing) subscribe to event
  - Decided: 2025-12-23 21:41 by default

- **AMB-OP-036: CreateQuote**
  - Q: What is the minimum required information to create a quote?
  - A: Required: customer reference, at least one line item (product, quantity, price); Optional: custom expiration, notes, discount
  - Decided: 2025-12-23 21:41 by default

- **AMB-OP-037: CreateQuote**
  - Q: Can a quote be created in Draft status before sending, or is it immediately Pending?
  - A: Quote created in Draft status; explicit action required to change to Pending and send to customer
  - Decided: 2025-12-23 21:41 by default

- **AMB-OP-038: CreateQuote**
  - Q: What validation is performed on prices? Can sales rep override product prices?
  - A: Prices default to product catalog price; sales rep can override within configured discount limits (e.g., max 20%); larger discounts require manager approval
  - Decided: 2025-12-23 21:41 by default

- **AMB-OP-039: CreateQuote**
  - Q: Can a quote be duplicated from an existing quote or order?
  - A: Yes, support duplicate from: existing quote (any status), previous order; new quote gets new reference number and fresh expiration
  - Decided: 2025-12-23 21:41 by default

- **AMB-OP-040: CreateQuote**
  - Q: What happens if a product in the quote is discontinued during quote creation?
  - A: Show warning but allow quote creation with discontinued product; flag quote for review; block if product is hard-deleted
  - Decided: 2025-12-23 21:41 by default

- **AMB-OP-041: CreateQuote**
  - Q: Is there a maximum number of line items per quote?
  - A: Soft limit of 100 line items with warning; hard limit of 500; larger quotes require splitting
  - Decided: 2025-12-23 21:41 by default

- **AMB-OP-042: CreateQuote**
  - Q: Should quote creation be audited and what details should be captured?
  - A: Yes, audit log captures: creation timestamp, sales rep, customer, initial line items, source (manual/duplicate/API)
  - Decided: 2025-12-23 21:41 by default

- **AMB-OP-043: CreateQuote**
  - Q: Can quotes be created via API or integration, or only through UI?
  - A: Both UI and API supported; API requires authentication and rate limiting; same validation rules apply
  - Decided: 2025-12-23 21:41 by default


## Decisions from 2025-12-23

- **AMB-ENT-002: Quote**
  - Q: What happens when a Quote expires? Is it automatic status change or manual? Can an expired quote be extended?
  - A: System automatically changes status to 'Expired' at midnight on expiration date. Sales rep can extend by creating a new quote.
  - Decided: 2025-12-23 22:00 by default

- **AMB-ENT-003: Quote**
  - Q: Can a Quote be modified after creation? If so, by whom and in which states?
  - A: Quotes can only be modified in Draft state by the creating sales representative or their manager.
  - Decided: 2025-12-23 22:00 by default

- **AMB-ENT-004: Quote**
  - Q: Can a Quote be deleted? If so, by whom and under what conditions?
  - A: Quotes are soft-deleted (marked as cancelled) and can only be cancelled by sales rep or manager while in Draft or Pending state.
  - Decided: 2025-12-23 22:00 by default

- **AMB-ENT-005: Quote**
  - Q: What is the format and generation rule for the unique link? How long should it be valid?
  - A: UUID-based token, valid until quote expiration date. Format: https://domain.com/quotes/{uuid-token}
  - Decided: 2025-12-23 22:00 by default

- **AMB-ENT-006: Quote**
  - Q: Is the 30-day default expiration configurable? At what level (system, per customer, per quote)?
  - A: Default is 30 days at system level, but sales rep can override per quote within allowed range (7-90 days).
  - Decided: 2025-12-23 22:00 by default

- **AMB-ENT-007: Quote**
  - Q: How is total value calculated? What tax rules apply? Are discounts percentage or fixed amount?
  - A: Total = sum(line items) - discounts + taxes. Discounts can be percentage or fixed. Tax calculated based on customer location.
  - Decided: 2025-12-23 22:00 by default

- **AMB-ENT-008: Quote**
  - Q: Is version history needed for quotes? Should we track all changes made before sending to customer?
  - A: Yes, maintain version history with who changed what and when. Each send to customer creates a new version.
  - Decided: 2025-12-23 22:00 by default

- **AMB-ENT-009: Quote**
  - Q: What additional attributes are needed? (e.g., created_at, created_by, currency, payment_terms, validity_period, notes)
  - A: Add: created_at, created_by, updated_at, currency, payment_terms, internal_notes, customer_notes
  - Decided: 2025-12-23 22:00 by default

- **AMB-ENT-010: Customer**
  - Q: What is the unique identifier for Customer? Is this an existing entity in the system or created during quote process?
  - A: Customer is an existing entity with system-generated ID. Quote must reference existing customer.
  - Decided: 2025-12-23 22:00 by default

- **AMB-ENT-011: Customer**
  - Q: What additional customer attributes are needed for quoting? (e.g., name, company, address, phone, billing address)
  - A: Required: name, email, company. Optional: phone, billing_address, shipping_address, tax_id
  - Decided: 2025-12-23 22:00 by default

- **AMB-ENT-012: Customer**
  - Q: Can a customer have multiple email addresses? Which one receives the quote?
  - A: Customer can have multiple emails. Quote is sent to primary email unless sales rep specifies otherwise.
  - Decided: 2025-12-23 22:00 by default

- **AMB-ENT-013: Customer**
  - Q: Does the IP address belong to Customer entity or only to QuoteAcceptance? Should we track customer login IPs?
  - A: IP address belongs to QuoteAcceptance only. Do not track general customer IPs.
  - Decided: 2025-12-23 22:00 by default

- **AMB-ENT-014: LineItem**
  - Q: What is the unique identifier for LineItem? Is it composite (quote_id + sequence) or independent?
  - A: System-generated UUID per line item, with sequence number for display ordering.
  - Decided: 2025-12-23 22:00 by default

- **AMB-ENT-015: LineItem**
  - Q: Is price the unit price or total price? What about line-level discounts?
  - A: Store unit_price and quantity. Calculate line_total. Support optional line-level discount (percentage or amount).
  - Decided: 2025-12-23 22:00 by default

- **AMB-ENT-016: LineItem**
  - Q: What additional attributes are needed? (e.g., description, SKU, tax_rate, notes)
  - A: Add: description (can override product default), sku, tax_rate, notes, sequence_number
  - Decided: 2025-12-23 22:00 by default

- **AMB-ENT-017: LineItem**
  - Q: What validation rules apply to quantity? Minimum? Maximum? Decimal allowed?
  - A: Quantity must be > 0, maximum 99999, decimals allowed to 2 places for unit-based products.
  - Decided: 2025-12-23 22:00 by default

- **AMB-ENT-018: LineItem**
  - Q: Can line items be added/removed/modified after quote is sent (Pending status)?
  - A: No modifications allowed after quote enters Pending status. Must cancel and create new quote.
  - Decided: 2025-12-23 22:00 by default

- **AMB-ENT-019: Order**
  - Q: What attributes should Order have? Is it a separate entity or reference to external system?
  - A: Order is a separate entity: order_id, order_number, quote_reference, status, created_at, customer_id
  - Decided: 2025-12-23 22:00 by default

- **AMB-ENT-020: Order**
  - Q: What is the relationship between Quote and Order? One-to-one or can multiple orders come from one quote?
  - A: One-to-one: each accepted quote creates exactly one order. Quote stores order_id reference.
  - Decided: 2025-12-23 22:00 by default

- **AMB-ENT-021: Order**
  - Q: What triggers should exist after order creation? (e.g., notify warehouse, send to ERP, notify sales rep)
  - A: After order creation: notify sales rep, send to fulfillment system, update customer record.
  - Decided: 2025-12-23 22:00 by default

- **AMB-ENT-022: SalesRepresentative**
  - Q: What is the unique identifier for SalesRepresentative? Internal user ID or separate entity?
  - A: SalesRepresentative is a role/reference to internal User entity with user_id as identifier.
  - Decided: 2025-12-23 22:00 by default

- **AMB-ENT-023: SalesRepresentative**
  - Q: What attributes are needed for SalesRepresentative? (e.g., name, email, department, manager, territory)
  - A: Reference User entity. Additional: territory, manager_id, quota, commission_rate
  - Decided: 2025-12-23 22:00 by default

- **AMB-ENT-024: SalesRepresentative**
  - Q: Can quotes be reassigned to different sales representatives? Under what conditions?
  - A: Yes, quotes can be reassigned by manager. Track original creator and current owner separately.
  - Decided: 2025-12-23 22:00 by default

- **AMB-ENT-025: QuoteAcceptance**
  - Q: Is QuoteAcceptance a separate entity or embedded in Quote? Should it track additional legal/compliance data?
  - A: Separate entity for audit purposes. Add: user_agent, acceptance_method, terms_version_accepted
  - Decided: 2025-12-23 22:00 by default

- **AMB-ENT-026: QuoteAcceptance**
  - Q: What additional data should be captured at acceptance? (e.g., browser info, geolocation, terms version)
  - A: Capture: IP address, user_agent, timestamp, terms_version, acceptance_token (for verification)
  - Decided: 2025-12-23 22:00 by default

- **AMB-ENT-027: QuoteAcceptance**
  - Q: How long should acceptance records be retained? Any compliance requirements?
  - A: Retain indefinitely for legal compliance. Never delete, even if quote is archived.
  - Decided: 2025-12-23 22:00 by default

- **AMB-ENT-028: Email**
  - Q: Is Email a separate entity or just an operation? Should we track all emails sent related to a quote?
  - A: Separate EmailLog entity to track all quote-related emails: type, sent_at, recipient, status, quote_id
  - Decided: 2025-12-23 22:00 by default

- **AMB-ENT-029: Email**
  - Q: What email types need to be supported? (e.g., quote sent, reminder, expiring soon, accepted, rejected)
  - A: Required: quote_sent, quote_accepted. Optional: expiration_reminder, quote_expired, quote_viewed
  - Decided: 2025-12-23 22:00 by default

- **AMB-ENT-030: Email**
  - Q: Should email delivery status be tracked? (e.g., sent, delivered, opened, bounced)
  - A: Track: sent, delivered, bounced, opened (if using tracking pixel). Store in EmailLog.
  - Decided: 2025-12-23 22:00 by default

- **AMB-ENT-031: Quote**
  - Q: What currency formats are supported? Single currency or multi-currency quotes?
  - A: Support multiple currencies. Quote has single currency. Exchange rates locked at quote creation.
  - Decided: 2025-12-23 22:00 by default

- **AMB-ENT-032: Quote**
  - Q: Is there a minimum or maximum quote value? Any approval workflows for high-value quotes?
  - A: No minimum. Quotes over $10,000 require manager approval before sending to customer.
  - Decided: 2025-12-23 22:00 by default

- **AMB-ENT-033: Customer**
  - Q: Can customers reject quotes explicitly, or only let them expire? Should rejection reason be captured?
  - A: Add explicit reject option with optional reason. Quote status changes to 'Rejected'. Notify sales rep.
  - Decided: 2025-12-23 22:00 by default

- **AMB-ENT-034: Quote**
  - Q: Can a customer request modifications to a quote before accepting? How is this handled?
  - A: Customer can request changes via reply email or form. Sales rep creates new version of quote.
  - Decided: 2025-12-23 22:00 by default

- **AMB-ENT-035: LineItem**
  - Q: What is the relationship between LineItem and Product? Must line items reference existing products?
  - A: Line items should reference Product entity but allow description/price override. Support custom line items without product reference.
  - Decided: 2025-12-23 22:00 by default

- **AMB-OP-001: SendQuoteEmail**
  - Q: Should the email be sent automatically when a quote is created, or manually triggered by the sales representative?
  - A: Sales rep explicitly clicks 'Send to Customer' button after reviewing the quote.
  - Decided: 2025-12-23 22:00 by default

- **AMB-OP-002: SendQuoteEmail**
  - Q: What happens if the email fails to send? Retry logic? Notification to sales rep?
  - A: Retry 3 times with exponential backoff. If all fail, notify sales rep and mark quote as 'Email Failed'.
  - Decided: 2025-12-23 22:00 by default

- **AMB-OP-003: SendQuoteEmail**
  - Q: Can a quote email be resent? Under what conditions? Is there a limit?
  - A: Yes, can be resent while quote is Pending. Maximum 5 resends. Each resend logged.
  - Decided: 2025-12-23 22:00 by default

- **AMB-OP-004: SendQuoteEmail**
  - Q: Should the email include a PDF attachment of the quote, or only the link?
  - A: Include both: link for acceptance and PDF attachment for offline review.
  - Decided: 2025-12-23 22:00 by default

- **AMB-OP-005: SendQuoteEmail**
  - Q: Should the sales representative be CC'd or BCC'd on the quote email?
  - A: BCC sales rep by default, configurable in user preferences.
  - Decided: 2025-12-23 22:00 by default

- **AMB-OP-006: SendQuoteEmail**
  - Q: Is this operation audited? What details should be logged?
  - A: Yes, log: timestamp, quote_id, recipient, sender, email_status, link_generated
  - Decided: 2025-12-23 22:00 by default

- **AMB-OP-007: ViewQuoteDetails**
  - Q: What should be displayed if the quote has expired when customer clicks the link?
  - A: Show friendly message that quote has expired with expiration date. Provide contact info for sales rep to request new quote.
  - Decided: 2025-12-23 22:00 by default

- **AMB-OP-008: ViewQuoteDetails**
  - Q: What should be displayed if the quote has already been accepted?
  - A: Show quote details with 'Already Accepted' status, acceptance date, and order reference if available.
  - Decided: 2025-12-23 22:00 by default

- **AMB-OP-009: ViewQuoteDetails**
  - Q: Should viewing the quote be tracked/logged? Should sales rep be notified when customer views?
  - A: Log all views with timestamp and IP. Notify sales rep on first view only.
  - Decided: 2025-12-23 22:00 by default

- **AMB-OP-010: ViewQuoteDetails**
  - Q: Does the customer need to authenticate in any way, or is the unique link sufficient?
  - A: Unique link is sufficient for viewing. No authentication required.
  - Decided: 2025-12-23 22:00 by default

- **AMB-OP-011: ViewQuoteDetails**
  - Q: What is the expected page load time? Any performance requirements?
  - A: Page should load within 2 seconds. Quote details should be pre-rendered or cached.
  - Decided: 2025-12-23 22:00 by default

- **AMB-OP-012: ViewQuoteDetails**
  - Q: What happens if the unique link is invalid or tampered with?
  - A: Show generic 'Quote not found' message. Log the attempt with IP for security monitoring.
  - Decided: 2025-12-23 22:00 by default

- **AMB-OP-013: AcceptQuote**
  - Q: Is any additional confirmation required before acceptance? (e.g., checkbox for terms, CAPTCHA)
  - A: Require checkbox to accept terms and conditions. No CAPTCHA for better UX.
  - Decided: 2025-12-23 22:00 by default

- **AMB-OP-014: AcceptQuote**
  - Q: What specific data should be returned/displayed after successful acceptance?
  - A: Show confirmation page with: acceptance confirmation number, quote reference, expected next steps, order number (when created).
  - Decided: 2025-12-23 22:00 by default

- **AMB-OP-015: AcceptQuote**
  - Q: What happens if customer double-clicks or submits acceptance multiple times?
  - A: Operation is idempotent. Second click shows 'Already accepted' confirmation. No duplicate orders created.
  - Decided: 2025-12-23 22:00 by default

- **AMB-OP-016: AcceptQuote**
  - Q: What happens if two people (e.g., from same company) try to accept simultaneously?
  - A: First acceptance wins. Second attempt shows 'Already accepted by [name/time]'. Use optimistic locking.
  - Decided: 2025-12-23 22:00 by default

- **AMB-OP-017: AcceptQuote**
  - Q: What if the acceptance process partially fails (e.g., acceptance recorded but email fails)?
  - A: Use saga pattern. If email fails, acceptance still valid. Retry email async. Show success to customer.
  - Decided: 2025-12-23 22:00 by default

- **AMB-OP-018: AcceptQuote**
  - Q: Should acceptance capture electronic signature or is button click sufficient?
  - A: Button click with terms checkbox is sufficient. Store acceptance metadata as proof.
  - Decided: 2025-12-23 22:00 by default

- **AMB-OP-019: AcceptQuote**
  - Q: What is the maximum acceptable response time for the acceptance operation?
  - A: User-facing acceptance should complete within 3 seconds. Order creation can be async.
  - Decided: 2025-12-23 22:00 by default

- **AMB-OP-020: RecordAcceptance**
  - Q: What additional metadata should be captured beyond timestamp and IP? (e.g., user agent, geolocation)
  - A: Capture: timestamp, IP address, user agent, browser fingerprint (optional), referrer, session ID
  - Decided: 2025-12-23 22:00 by default

- **AMB-OP-021: RecordAcceptance**
  - Q: Should this record be immutable? Can it ever be modified or deleted?
  - A: Immutable. Cannot be modified or deleted. Append-only audit trail.
  - Decided: 2025-12-23 22:00 by default

- **AMB-OP-022: RecordAcceptance**
  - Q: What if IP address cannot be determined (e.g., proxy, VPN)?
  - A: Store whatever IP is available (may be proxy). Add flag indicating if proxy detected. Do not block acceptance.
  - Decided: 2025-12-23 22:00 by default

- **AMB-OP-023: SendConfirmationEmail**
  - Q: What should the confirmation email contain? Just confirmation or detailed order info?
  - A: Include: confirmation number, quote summary, acceptance timestamp, next steps, contact info, link to view order status
  - Decided: 2025-12-23 22:00 by default

- **AMB-OP-024: SendConfirmationEmail**
  - Q: Should the sales representative also receive a notification email?
  - A: Yes, send separate notification to sales rep with acceptance details and customer info.
  - Decided: 2025-12-23 22:00 by default

- **AMB-OP-026: UpdateQuoteStatus**
  - Q: What other status transitions are valid? Define the complete state machine.
  - A: Valid transitions: Draft→Pending, Pending→Accepted, Pending→Expired, Pending→Rejected, Pending→Cancelled, Draft→Cancelled
  - Decided: 2025-12-23 22:00 by default

- **AMB-OP-027: UpdateQuoteStatus**
  - Q: Should status changes trigger notifications? To whom?
  - A: Notify: Sales rep on all changes, Customer on Accepted/Expired/Rejected, Manager on high-value quote acceptance
  - Decided: 2025-12-23 22:00 by default

- **AMB-OP-028: UpdateQuoteStatus**
  - Q: Should all status changes be logged in an audit trail?
  - A: Yes, log: old_status, new_status, changed_by, timestamp, reason (if applicable)
  - Decided: 2025-12-23 22:00 by default

- **AMB-OP-029: UpdateQuoteStatus**
  - Q: Can status be reverted? (e.g., Accepted back to Pending)
  - A: No. Status changes are forward-only. To 'undo', cancel current quote and create new one.
  - Decided: 2025-12-23 22:00 by default

- **AMB-OP-030: TriggerOrderCreation**
  - Q: Should order creation be synchronous or asynchronous?
  - A: Asynchronous. Customer sees acceptance confirmation immediately. Order created in background within 30 seconds.
  - Decided: 2025-12-23 22:00 by default

- **AMB-OP-031: TriggerOrderCreation**
  - Q: What happens if order creation fails? Should the quote acceptance be rolled back?
  - A: No rollback. Quote remains Accepted. Alert operations team. Retry order creation. Manual intervention if retries exhausted.
  - Decided: 2025-12-23 22:00 by default

- **AMB-OP-032: TriggerOrderCreation**
  - Q: What order data should be copied from the quote vs. referenced?
  - A: Copy: line items, prices, customer details, totals at acceptance time. Reference: quote_id, customer_id
  - Decided: 2025-12-23 22:00 by default

- **AMB-OP-033: TriggerOrderCreation**
  - Q: Should the quote-to-order process integrate with external systems (ERP, inventory)?
  - A: Yes, publish OrderCreated event. External systems subscribe. Don't block on external integration.
  - Decided: 2025-12-23 22:00 by default

- **AMB-OP-034: TriggerOrderCreation**
  - Q: What is the retry policy for failed order creation?
  - A: Retry 5 times with exponential backoff (1s, 2s, 4s, 8s, 16s). Then alert operations for manual handling.
  - Decided: 2025-12-23 22:00 by default

- **AMB-OP-035: CreateQuote**
  - Q: What is the complete list of required vs optional inputs for creating a quote?
  - A: Required: customer_id, at least one line item. Optional: expiration_date, notes, discount, payment_terms, custom_fields
  - Decided: 2025-12-23 22:00 by default

- **AMB-OP-036: CreateQuote**
  - Q: What validation rules apply to line items? (e.g., min items, max items, valid products)
  - A: Min 1 line item, max 100. Products must be active. Quantities > 0. Prices > 0 (or allow zero for free items).
  - Decided: 2025-12-23 22:00 by default

- **AMB-OP-037: CreateQuote**
  - Q: How is the reference number (QT-XXXXX) generated? Is it sequential, random, or formatted?
  - A: Sequential with prefix: QT-{YEAR}{5-digit-sequence}. Example: QT-2024-00001. Reset sequence annually.
  - Decided: 2025-12-23 22:00 by default

- **AMB-OP-038: CreateQuote**
  - Q: Can any sales representative create quotes for any customer, or are there territory/assignment restrictions?
  - A: Sales reps can only create quotes for customers in their assigned territory. Managers can override.
  - Decided: 2025-12-23 22:00 by default

- **AMB-OP-039: CreateQuote**
  - Q: Should quote creation be audited? What details?
  - A: Yes. Log: created_by, created_at, customer_id, total_value, line_item_count, initial_status
  - Decided: 2025-12-23 22:00 by default

- **AMB-OP-040: CreateQuote**
  - Q: Can the same quote (same customer, same items) be created multiple times? Any duplicate detection?
  - A: Allow duplicates but warn sales rep if similar quote exists for same customer within 7 days.
  - Decided: 2025-12-23 22:00 by default

- **AMB-OP-041: CreateQuote**
  - Q: What happens if product prices change after quote creation but before acceptance?
  - A: Quote prices are locked at creation time. Customer accepts at quoted price regardless of current price.
  - Decided: 2025-12-23 22:00 by default

- **AMB-OP-042: CreateQuote**
  - Q: Is there a draft/save feature, or must quotes be complete on creation?
  - A: Support draft status. Quotes can be saved incomplete. Must be complete before sending to customer.
  - Decided: 2025-12-23 22:00 by default

- **AMB-OP-043: AcceptQuote**
  - Q: Can a quote be accepted if the referenced products are no longer available?
  - A: Yes, acceptance honors the quote as-is. Inventory/availability issues handled during order fulfillment.
  - Decided: 2025-12-23 22:00 by default

- **AMB-OP-044: CreateQuote**
  - Q: What is the initial status of a newly created quote?
  - A: Draft by default. Changes to Pending when sales rep clicks 'Send to Customer'.
  - Decided: 2025-12-23 22:00 by default

- **AMB-OP-045: ViewQuoteDetails**
  - Q: Should customers be able to download/print the quote from the view page?
  - A: Yes, provide 'Download PDF' and 'Print' buttons on the quote view page.
  - Decided: 2025-12-23 22:00 by default


## Decisions from 2025-12-23

- **AMB-ENT-002: Quote**
  - Q: What happens when a Quote expires?
  - A: System automatically transitions quote to Expired state via scheduled job that runs daily. Expired quotes cannot be accepted but remain viewable.
  - Decided: 2025-12-23 22:13 by ai_suggested

- **AMB-ENT-003: Quote**
  - Q: Can a Quote be modified after creation?
  - A: Quote can be modified only in Draft or Pending state by the creating sales representative or their manager. Modifications are logged in audit trail.
  - Decided: 2025-12-23 22:13 by ai_suggested

- **AMB-ENT-004: Quote**
  - Q: Can a Quote be deleted?
  - A: Soft delete only with audit trail. Allowed for Draft/Pending quotes by sales rep (own quotes) or admin (any quote). Accepted/Expired quotes cannot be deleted.
  - Decided: 2025-12-23 22:13 by ai_suggested

- **AMB-ENT-005: Quote**
  - Q: What is the format for the unique reference number?
  - A: QT- prefix followed by 5-digit zero-padded sequential number (e.g., QT-00001). Sequence is global, not per-year.
  - Decided: 2025-12-23 22:13 by ai_suggested

- **AMB-ENT-006: Quote**
  - Q: What attributes are required vs optional when creating a Quote?
  - A: Required: customer. Empty quotes allowed for Draft state. At least one line item required before transitioning to Pending. Expiration date defaults to 30 days if not specified.
  - Decided: 2025-12-23 22:13 by ai_suggested

- **AMB-ENT-007: Quote**
  - Q: What is the data type and precision for monetary values?
  - A: Integer cents (smallest currency unit) to avoid floating point issues. Currency code stored per quote, defaulting to system currency. All calculations performed in cents.
  - Decided: 2025-12-23 22:13 by ai_suggested

- **AMB-ENT-008: Quote**
  - Q: Is an audit trail required for Quote changes?
  - A: Full field-level audit trail with timestamp, user ID, and before/after values for all changes. Stored in separate audit log table.
  - Decided: 2025-12-23 22:13 by ai_suggested

- **AMB-ENT-009: Quote**
  - Q: What is the validation rule for expiration date?
  - A: Must be future date. Maximum 90 days from creation. Default 30 days. Cannot be modified after quote is sent (Pending state).
  - Decided: 2025-12-23 22:13 by ai_suggested

- **AMB-ENT-010: Quote**
  - Q: How are discounts represented?
  - A: Support both percentage and fixed amount discounts at both line item and quote level. Line discounts applied first, then quote-level discount on subtotal.
  - Decided: 2025-12-23 22:13 by ai_suggested

- **AMB-ENT-011: Customer**
  - Q: What uniquely identifies a Customer?
  - A: System-generated customer ID (UUID) as primary key. Email address as unique secondary identifier. Both are immutable once set.
  - Decided: 2025-12-23 22:13 by ai_suggested

- **AMB-ENT-012: Customer**
  - Q: What additional attributes does a Customer have?
  - A: Full profile: name (required), company name, phone, billing address, shipping address, tax ID. Billing address required for quote acceptance.
  - Decided: 2025-12-23 22:13 by ai_suggested

- **AMB-ENT-013: Customer**
  - Q: Does a Customer need an account to accept quotes?
  - A: Guest acceptance via unique link with optional account creation. Customer details collected at acceptance time. Account creation offered but not required.
  - Decided: 2025-12-23 22:13 by ai_suggested

- **AMB-ENT-015: Customer**
  - Q: Can a Customer be deleted?
  - A: Anonymization instead of deletion for GDPR compliance. Personal data anonymized but quotes and orders retained with anonymized reference for business records.
  - Decided: 2025-12-23 22:13 by ai_suggested

- **AMB-ENT-016: Customer**
  - Q: Where is IP address stored?
  - A: IP address stored only on QuoteAcceptance record for legal/audit purposes. Not stored on Customer entity.
  - Decided: 2025-12-23 22:13 by ai_suggested

- **AMB-ENT-017: Customer**
  - Q: What uniquely identifies a LineItem?
  - A: System-generated line item ID (UUID), allowing same product to appear multiple times with different configurations or quantities.
  - Decided: 2025-12-23 22:13 by ai_suggested

- **AMB-ENT-018: LineItem**
  - Q: What are the validation rules for quantity?
  - A: Positive integer only, minimum 1, maximum 99999. No decimal quantities. Validation error if outside range.
  - Decided: 2025-12-23 22:13 by ai_suggested

- **AMB-ENT-019: LineItem**
  - Q: Is price unit or total? When is it locked?
  - A: Store both unit price and line total. Prices locked at time of adding to quote (snapshot from product catalog). Line total = unit price × quantity - line discount.
  - Decided: 2025-12-23 22:13 by ai_suggested

- **AMB-ENT-020: Customer**
  - Q: What additional attributes might a LineItem have?
  - A: Description (from product), SKU, unit price, quantity, line discount (amount or percentage), line total, optional notes field for customizations.
  - Decided: 2025-12-23 22:13 by ai_suggested

- **AMB-ENT-021: LineItem**
  - Q: Can LineItems be modified after Quote creation?
  - A: LineItems can be added, modified, or removed when Quote is in Draft or Pending state. Changes not allowed after customer views (to prevent confusion).
  - Decided: 2025-12-23 22:13 by ai_suggested

- **AMB-ENT-022: Order**
  - Q: What is the relationship between Quote and Order?
  - A: One-to-one relationship: one accepted quote creates exactly one order. Order references source quote ID.
  - Decided: 2025-12-23 22:13 by ai_suggested

- **AMB-ENT-023: Order**
  - Q: What attributes does Order have?
  - A: Full independent Order entity with: order ID, quote reference, customer snapshot, line items (copied), status, timestamps, payment info, shipping info, fulfillment details.
  - Decided: 2025-12-23 22:13 by ai_suggested

- **AMB-ENT-024: Order**
  - Q: What states does an Order go through?
  - A: States: Created, Confirmed, Processing, Shipped, Delivered, Cancelled, Refunded. State machine enforces valid transitions.
  - Decided: 2025-12-23 22:13 by ai_suggested

- **AMB-ENT-025: Order**
  - Q: What is the format of the order number?
  - A: ORD- prefix followed by 5-digit zero-padded sequential number (e.g., ORD-00001). Same pattern as quotes but separate sequence.
  - Decided: 2025-12-23 22:13 by ai_suggested

- **AMB-ENT-026: Order**
  - Q: What attributes does a SalesRepresentative have?
  - A: Integrated from HR/identity system: employee ID, name, email, phone, team/department, manager reference, active status. Extended attributes (commission, quota) in separate sales module.
  - Decided: 2025-12-23 22:13 by ai_suggested

- **AMB-ENT-027: SalesRepresentative**
  - Q: What permissions does a SalesRepresentative have?
  - A: Role-based: Sales reps manage own quotes; Team leads view/manage team quotes; Managers see department; Admins have full access. All based on RBAC system.
  - Decided: 2025-12-23 22:13 by ai_suggested

- **AMB-ENT-028: SalesRepresentative**
  - Q: What happens to quotes when a SalesRepresentative leaves?
  - A: Auto-reassign pending quotes to manager. Original creator retained in audit history. Completed quotes remain associated with original rep for reporting.
  - Decided: 2025-12-23 22:13 by ai_suggested

- **AMB-ENT-029: QuoteAcceptance**
  - Q: Is QuoteAcceptance a separate entity or embedded?
  - A: Embedded within Quote as acceptance metadata fields: accepted_at, accepted_ip, accepted_user_agent, tc_version_accepted, acceptance_method.
  - Decided: 2025-12-23 22:13 by ai_suggested

- **AMB-ENT-030: QuoteAcceptance**
  - Q: Should QuoteAcceptance capture legal acknowledgment?
  - A: Mandatory T&C acceptance checkbox with timestamp. Digital signature optional, configurable for quotes above value threshold (e.g., $10,000).
  - Decided: 2025-12-23 22:13 by ai_suggested

- **AMB-ENT-031: Quote**
  - Q: What is the format and security of the UniqueLink?
  - A: Cryptographically secure UUID v4 token. Valid until quote expiration. Multi-use for viewing, single-use for acceptance (prevents duplicate orders).
  - Decided: 2025-12-23 22:13 by ai_suggested

- **AMB-ENT-032: UniqueLink**
  - Q: Can a new UniqueLink be generated?
  - A: Sales rep can regenerate link via UI action. Previous link immediately invalidated. Action logged in audit trail. Customer notified of new link.
  - Decided: 2025-12-23 22:13 by ai_suggested

- **AMB-ENT-033: UniqueLink**
  - Q: Is UniqueLink a separate entity or attribute?
  - A: Attribute of Quote (access_token field) rather than separate entity. Simple and sufficient for one-to-one relationship.
  - Decided: 2025-12-23 22:13 by ai_suggested

- **AMB-ENT-034: Email**
  - Q: Is Email an entity that needs to be tracked?
  - A: Track as EmailLog entity: id, type (quote_sent, confirmation, reminder), recipient, timestamp, status (queued/sent/delivered/bounced/failed), quote_id, error_message.
  - Decided: 2025-12-23 22:13 by ai_suggested

- **AMB-ENT-035: Email**
  - Q: What happens if email delivery fails?
  - A: Automatic retry with exponential backoff: 3 attempts over 24 hours (immediate, 1 hour, 8 hours). Notify sales rep on final failure. Log all attempts.
  - Decided: 2025-12-23 22:13 by ai_suggested

- **AMB-ENT-036: Quote**
  - Q: Is there a Product entity?
  - A: Product entity exists with: product_id, SKU, name, description, base_price, active status, category_id. Integration point for external product/inventory systems.
  - Decided: 2025-12-23 22:13 by ai_suggested

- **AMB-ENT-037: Quote**
  - Q: Can a customer reject a quote?
  - A: Add Rejected state to quote lifecycle. Customer can reject via link with optional reason (dropdown + free text). Sales rep notified immediately.
  - Decided: 2025-12-23 22:13 by ai_suggested

- **AMB-ENT-038: Quote**
  - Q: What timestamps should Quote track?
  - A: Track: created_at, updated_at, sent_at, first_viewed_at, accepted_at, rejected_at, expired_at, cancelled_at. All nullable except created_at/updated_at.
  - Decided: 2025-12-23 22:13 by ai_suggested

- **AMB-ENT-039: Quote**
  - Q: Should the system track quote views?
  - A: Track first_viewed_at timestamp and view_count. Sufficient for sales insights. Detailed view log (with IPs) available in audit trail if needed.
  - Decided: 2025-12-23 22:13 by ai_suggested

- **AMB-ENT-040: Quote**
  - Q: Are approval workflows required?
  - A: Approval required when: discount exceeds 15% or total value exceeds $50,000. Approver is manager of quote creator. Configurable thresholds.
  - Decided: 2025-12-23 22:13 by ai_suggested

- **AMB-OP-001: SendQuoteEmail**
  - Q: Is email sent automatically or explicitly?
  - A: Sales rep explicitly triggers sending via 'Send to Customer' action. Quote transitions from Draft to Pending state upon send.
  - Decided: 2025-12-23 22:13 by ai_suggested

- **AMB-OP-002: Email**
  - Q: What happens if email bounces?
  - A: Log bounce event, update EmailLog status, notify sales rep immediately via in-app notification. Quote remains in Pending state. Sales rep can update email and resend.
  - Decided: 2025-12-23 22:13 by ai_suggested

- **AMB-OP-003: SendQuoteEmail**
  - Q: Can quote email be resent?
  - A: Yes, sales rep can resend. Maximum 5 times total. Minimum 1 hour between sends (cooldown). Same unique link used. All sends logged.
  - Decided: 2025-12-23 22:13 by ai_suggested

- **AMB-OP-004: SendQuoteEmail**
  - Q: What content should quote email include?
  - A: Summary: customer name, quote number, line item count, total value, expiration date, plus unique link to view full details. Company branding/logo in header.
  - Decided: 2025-12-23 22:13 by ai_suggested

- **AMB-OP-005: SendQuoteEmail**
  - Q: Should SendQuoteEmail be audited?
  - A: Yes. Log: timestamp, quote_id, recipient_email, sender_user_id, email_service_message_id, status. Stored in EmailLog entity.
  - Decided: 2025-12-23 22:13 by ai_suggested

- **AMB-OP-006: SendQuoteEmail**
  - Q: What is the email sending timeout?
  - A: Queue-based with background processing. 30 second timeout per attempt. Async from user action (user sees immediate 'Email queued' confirmation).
  - Decided: 2025-12-23 22:13 by ai_suggested

- **AMB-OP-009: ViewQuoteDetails**
  - Q: What information is displayed to customer?
  - A: Display: company info, quote number, line items (description, quantity, unit price, line total), subtotal, discounts, taxes, grand total, expiration date, T&C link, Accept/Reject buttons.
  - Decided: 2025-12-23 22:13 by ai_suggested

- **AMB-OP-011: ViewQuoteDetails**
  - Q: What happens if unique link is invalid?
  - A: Show generic 'This link is invalid or has expired' error page. Do not reveal whether quote exists. Log security event with IP for monitoring.
  - Decided: 2025-12-23 22:13 by ai_suggested

- **AMB-OP-012: ViewQuoteDetails**
  - Q: Is there rate limiting on quote viewing?
  - A: Rate limit: 60 requests per hour per IP address. After threshold: show CAPTCHA challenge. Log excessive attempts for security review.
  - Decided: 2025-12-23 22:13 by ai_suggested

- **AMB-OP-013: AcceptQuote**
  - Q: Does customer need to agree to T&C?
  - A: Yes, mandatory checkbox: 'I agree to the Terms and Conditions' with link to T&C document. Record T&C version accepted and checkbox timestamp.
  - Decided: 2025-12-23 22:13 by ai_suggested

- **AMB-OP-014: AcceptQuote**
  - Q: What customer information is required at acceptance?
  - A: Verify/collect: full name (required), billing address (required), shipping address (if different, optional), phone (required), PO number (optional). Pre-fill from customer record if available.
  - Decided: 2025-12-23 22:13 by ai_suggested

- **AMB-OP-015: AcceptQuote**
  - Q: Is payment part of quote acceptance?
  - A: No payment at acceptance. Quote acceptance is commitment to purchase. Payment handled separately in order processing (invoice, credit terms, or payment gateway based on customer type).
  - Decided: 2025-12-23 22:13 by ai_suggested

- **AMB-OP-016: AcceptQuote**
  - Q: What happens if order creation fails?
  - A: Show error to customer with apology message. Queue for automatic retry (3 attempts over 15 minutes). Alert operations team. Keep quote in 'Acceptance Processing' intermediate state until resolved.
  - Decided: 2025-12-23 22:13 by ai_suggested

- **AMB-OP-017: AcceptQuote**
  - Q: Can quote be accepted multiple times?
  - A: Idempotent operation. Subsequent acceptance attempts show 'Quote already accepted' confirmation page with order number and status link. No duplicate orders created.
  - Decided: 2025-12-23 22:13 by ai_suggested

- **AMB-OP-018: AcceptQuote**
  - Q: What happens with simultaneous acceptance?
  - A: Use database-level optimistic locking on quote record. First transaction to complete wins. Second sees 'Quote already accepted' message. No duplicate orders possible.
  - Decided: 2025-12-23 22:13 by ai_suggested

- **AMB-OP-019: AcceptQuote**
  - Q: Should user agent be captured?
  - A: Capture IP address and user agent string. Useful for fraud detection and customer support. Store in QuoteAcceptance metadata.
  - Decided: 2025-12-23 22:13 by ai_suggested

- **AMB-OP-020: AcceptQuote**
  - Q: Is there a confirmation step?
  - A: Two-step process: Click 'Accept Quote' → Confirmation dialog shows quote summary, T&C checkbox, customer details form → Click 'Confirm and Place Order' to finalize.
  - Decided: 2025-12-23 22:13 by ai_suggested

- **AMB-OP-021: RecordAcceptance**
  - Q: Is RecordAcceptance a separate operation?
  - A: Part of AcceptQuote database transaction. Recording acceptance metadata and updating quote status happen atomically. Both succeed or both rollback.
  - Decided: 2025-12-23 22:13 by ai_suggested

- **AMB-OP-022: RecordAcceptance**
  - Q: What metadata to record beyond timestamp and IP?
  - A: Record: acceptance_timestamp, ip_address, user_agent, session_id, tc_version_accepted, quote_version_hash (for integrity verification). All stored in Quote.acceptance_metadata JSON field.
  - Decided: 2025-12-23 22:13 by ai_suggested

- **AMB-OP-024: SendConfirmationEmail**
  - Q: What should confirmation email include?
  - A: Include: 'Thank you' header, quote/order summary, order number prominently displayed, next steps (what happens now), estimated timeline, customer service contact info, link to order status page.
  - Decided: 2025-12-23 22:13 by ai_suggested

- **AMB-OP-025: SendConfirmationEmail**
  - Q: Should sales rep receive notification?
  - A: Yes, both email and in-app notification. Email: summary of acceptance with customer details and order number. In-app: real-time notification with link to order.
  - Decided: 2025-12-23 22:13 by ai_suggested

- **AMB-OP-026: SendConfirmationEmail**
  - Q: What if confirmation email fails?
  - A: Do not rollback acceptance or order. Queue email for retry (3 attempts). Alert operations on final failure. Customer can view confirmation on web via success page. Log failure.
  - Decided: 2025-12-23 22:13 by ai_suggested

- **AMB-OP-027: SendConfirmationEmail**
  - Q: Should confirmation email include details or link?
  - A: Include summary (order number, items count, total, estimated delivery) plus prominent link to full order details in customer portal/order status page.
  - Decided: 2025-12-23 22:13 by ai_suggested

- **AMB-OP-028: UpdateQuoteStatus**
  - Q: What status transitions are valid?
  - A: Valid transitions: Draft→Pending, Pending→Accepted, Pending→Expired, Pending→Cancelled, Pending→Rejected. No backward transitions from terminal states. Admin cannot reverse Accepted status.
  - Decided: 2025-12-23 22:13 by ai_suggested

- **AMB-OP-029: UpdateQuoteStatus**
  - Q: Should status updates trigger notifications?
  - A: Notify customer on: Accepted (confirmation), Expired (opportunity to request extension), Rejected (confirmation). Notify sales rep on all status changes via in-app + email for Accepted/Rejected.
  - Decided: 2025-12-23 22:13 by ai_suggested

- **AMB-OP-031: TriggerOrderCreation**
  - Q: Is order creation synchronous or async?
  - A: Synchronous for core order creation (customer waits, sees order number). Asynchronous for downstream: inventory reservation, fulfillment queue, accounting notifications.
  - Decided: 2025-12-23 22:13 by ai_suggested

- **AMB-OP-032: TriggerOrderCreation**
  - Q: What data is copied vs referenced?
  - A: Copy (snapshot): line items with prices, customer details, addresses. Reference: quote_id, product_ids (for linking to current product data). Snapshot ensures order accuracy regardless of future changes.
  - Decided: 2025-12-23 22:13 by ai_suggested

- **AMB-OP-033: TriggerOrderCreation**
  - Q: What validation during order creation?
  - A: Validate: products still active, inventory available (warning only, not blocking), shipping address format valid, customer not blocked. Credit check deferred to payment processing.
  - Decided: 2025-12-23 22:13 by ai_suggested

- **AMB-OP-034: TriggerOrderCreation**
  - Q: What if inventory not available?
  - A: Create order with backorder status for unavailable items. Notify customer in confirmation: 'Some items on backorder, estimated restock date X'. Notify sales rep. Do not block order creation.
  - Decided: 2025-12-23 22:13 by ai_suggested

- **AMB-OP-035: TriggerOrderCreation**
  - Q: Should order creation trigger downstream processes?
  - A: Trigger via event/message queue: inventory reservation request, fulfillment queue entry, accounting/ERP notification, CRM opportunity update. All async, non-blocking.
  - Decided: 2025-12-23 22:13 by ai_suggested

- **AMB-OP-036: CreateQuote**
  - Q: What are required vs optional inputs for CreateQuote?
  - A: Required: customer (existing or new). Optional: line items (can add after), expiration date (default 30 days), notes, discount, custom terms, PO reference. Quote starts in Draft state.
  - Decided: 2025-12-23 22:13 by ai_suggested

- **AMB-OP-037: CreateQuote**
  - Q: Can quote be created for new customer?
  - A: Both allowed. UI provides: search/select existing customer OR 'Create New Customer' inline form with required fields (name, email). New customer created atomically with quote.
  - Decided: 2025-12-23 22:13 by ai_suggested

- **AMB-OP-038: CreateQuote**
  - Q: What validation rules for line items?
  - A: Validate: product must be active, quantity must be positive integer (1-99999), unit price must be within range (base price ± configurable discount limit, e.g., max 25% below base).
  - Decided: 2025-12-23 22:13 by ai_suggested

- **AMB-OP-039: CreateQuote**
  - Q: Who can create quotes?
  - A: Users with 'sales_rep' role or higher. Optionally restricted by territory (if enabled) or product line permissions. Configurable in role management.
  - Decided: 2025-12-23 22:13 by ai_suggested

- **AMB-OP-040: CreateQuote**
  - Q: Is quote creation draft first or immediately pending?
  - A: Create as Draft state. Sales rep edits as needed, then explicit 'Send to Customer' action validates completeness and transitions to Pending state.
  - Decided: 2025-12-23 22:13 by ai_suggested

- **AMB-OP-041: CreateQuote**
  - Q: Can quote be cloned?
  - A: Yes, 'Clone Quote' action available. Copies: customer (can change), line items, notes. Resets: dates (new expiration), status (Draft), quote number (new). Modification required before send.
  - Decided: 2025-12-23 22:13 by ai_suggested

- **AMB-OP-042: CreateQuote**
  - Q: What if pricing changes between creation and acceptance?
  - A: Prices locked at quote creation. Quote represents commitment to sell at stated prices until expiration. Product price changes do not affect existing quotes.
  - Decided: 2025-12-23 22:13 by ai_suggested

- **AMB-OP-043: CreateQuote**
  - Q: Should quote creation be audited?
  - A: Yes. Log: creator_user_id, created_at, customer_id, initial_total_value, line_item_count. Stored in audit log. Full quote snapshot available via audit trail.
  - Decided: 2025-12-23 22:13 by ai_suggested

- **AMB-OP-044: AcceptQuote**
  - Q: What is expected response time for accept?
  - A: Target: complete within 5 seconds. Show progress indicator after 2 seconds. Hard timeout at 30 seconds with 'Please try again' message and retry button.
  - Decided: 2025-12-23 22:13 by ai_suggested

- **AMB-OP-045: AcceptQuote**
  - Q: What is returned upon successful acceptance?
  - A: Confirmation page showing: success message, order number (prominent), order summary, next steps ('You will receive confirmation email', 'Estimated delivery: X'), link to order status page.
  - Decided: 2025-12-23 22:13 by ai_suggested

