# Acceptance Criteria

Generated: 2025-12-23T21:44:26+01:00

---

## AC-QUOTE-001 – Customer views quote via unique link

**Given** A customer has received an email with a unique quote link AND the quote is in 'Pending' status AND the quote has not expired
**When** The customer clicks the unique link in the email
**Then** The system displays the quote details page showing: quote reference number, creation date, expiration date, all line items with descriptions/quantities/prices, subtotal, applicable taxes, discounts, total amount, and terms and conditions

**Error Cases:**
- Invalid/malformed link token → Display 'Quote not found' error page with contact information
- Quote has expired → Display friendly expiration message with quote reference, option to request new quote, and sales rep contact info
- Quote status is not Pending (already Accepted/Rejected/Cancelled) → Display appropriate status message with quote details in read-only mode
- Quote has been soft-deleted → Display 'Quote no longer available' message

**Traceability:**
- Source: US-QUOTE-003
- Decision: AMB-ENT-005
- Decision: AMB-OP-007
- Decision: AMB-OP-008
- Decision: AMB-OP-010

---

## AC-QUOTE-002 – Quote page loads within performance requirements

**Given** A valid quote link exists for a Pending, non-expired quote
**When** The customer navigates to the quote view page
**Then** The page fully loads within 2 seconds with all quote data displayed, using cached data for performance optimization

**Error Cases:**
- Page load exceeds 2 seconds → Log performance warning for monitoring
- Cache miss → Fall back to database query without failing the request

**Traceability:**
- Source: US-QUOTE-003
- Decision: AMB-OP-011

---

## AC-QUOTE-003 – Customer can download quote as PDF

**Given** A customer is viewing a valid quote details page
**When** The customer clicks the 'Download PDF' button
**Then** The system generates and downloads a PDF document containing all quote details with company branding

**Error Cases:**
- PDF generation fails → Display error message with retry option

**Traceability:**
- Source: US-QUOTE-003
- Decision: AMB-OP-012

---

## AC-QUOTE-004 – Quote views are tracked for analytics

**Given** A customer accesses a quote via unique link
**When** The quote details page is loaded
**Then** The system records a view event with timestamp, IP address, and user agent; this data is visible to the assigned sales rep on the quote detail page

**Error Cases:**
- View tracking fails → Continue displaying quote (non-blocking), log error for investigation

**Traceability:**
- Source: US-QUOTE-003
- Decision: AMB-OP-009

---

## AC-QUOTE-005 – Customer accepts quote with terms acknowledgment

**Given** A customer is viewing a valid, non-expired quote in 'Pending' status AND the customer has checked the 'I agree to terms and conditions' checkbox
**When** The customer clicks the 'Accept Quote' button
**Then** The system: (1) validates quote is still Pending and not expired, (2) checks quote version hasn't changed, (3) records acceptance with timestamp/IP/user agent/terms version, (4) updates quote status to 'Accepted', (5) displays acceptance confirmation message, (6) sends confirmation email to customer, (7) notifies sales rep via email and in-app notification, (8) triggers asynchronous order creation

**Error Cases:**
- Terms checkbox not checked → Disable Accept button, show 'Please accept terms and conditions' message
- Quote expired between page load and acceptance → Show 'Quote has expired' error, offer to request new quote
- Quote modified by sales rep since page load → Show 'Quote has been updated' error, refresh page with new quote details
- Quote already accepted (concurrent acceptance) → Return success with original acceptance details (idempotent)
- Technical failure during acceptance → Show friendly error, preserve form state, provide retry button, show support contact after 2 retries
- Session timeout (30 minutes) → Require re-click of original link, refresh quote data

**Traceability:**
- Source: US-QUOTE-003
- Decision: AMB-OP-013
- Decision: AMB-OP-014
- Decision: AMB-OP-015
- Decision: AMB-OP-016
- Decision: AMB-OP-017
- Decision: AMB-OP-018
- Decision: AMB-OP-019

---

## AC-QUOTE-006 – Acceptance record is persisted with audit trail

**Given** A customer has clicked Accept Quote with valid preconditions
**When** The acceptance is being processed
**Then** The system creates an immutable QuoteAcceptance record containing: acceptance timestamp, customer IP address, user agent, terms version accepted, quote version hash; AND writes to both the QuoteAcceptance table and the immutable audit log

**Error Cases:**
- Database write fails → Rollback entire transaction including status change, show error to customer, auto-retry, alert operations team
- Audit log write fails → Rollback entire transaction, alert operations team immediately

**Traceability:**
- Source: US-QUOTE-003
- Decision: AMB-OP-020
- Decision: AMB-OP-021
- Decision: AMB-OP-022
- Decision: AMB-ENT-031

---

## AC-QUOTE-007 – Confirmation email sent after acceptance

**Given** A quote has been successfully accepted AND the QuoteAcceptance record has been created
**When** The acceptance transaction completes
**Then** The system sends a confirmation email to the customer containing: order confirmation number, quote reference, accepted items summary, total amount, payment terms and methods, expected next steps, and contact information; AND BCCs the assigned sales rep

**Error Cases:**
- Email service unavailable → Queue email for retry, do NOT rollback acceptance, continue with success response to customer
- Email bounces → Mark delivery failure, notify sales rep for manual follow-up, allow manual resend
- Invalid customer email → Alert sales rep immediately for manual contact

**Traceability:**
- Source: US-QUOTE-003
- Decision: AMB-OP-024
- Decision: AMB-OP-025
- Decision: AMB-OP-026
- Decision: AMB-OP-027

---

## AC-QUOTE-008 – Order creation triggered asynchronously

**Given** A quote status has changed to 'Accepted'
**When** The quote acceptance transaction completes
**Then** The system publishes an event to the message queue for order creation; the order creation process runs asynchronously and creates an Order with: order number (ORD-XXXXX format), reference to source quote, copied customer/line items/prices/totals, order date, and initial status 'Created'

**Error Cases:**
- Order creation validation fails (product discontinued, customer inactive) → Mark quote with 'order creation pending' flag, alert operations team
- Order creation fails after 3 retries → Escalate to operations team, do NOT rollback quote acceptance
- Message queue unavailable → Persist event locally for later processing, alert operations team

**Traceability:**
- Source: US-QUOTE-003
- Decision: AMB-OP-032
- Decision: AMB-OP-033
- Decision: AMB-OP-034
- Decision: AMB-ENT-023
- Decision: AMB-ENT-024
- Decision: AMB-ENT-026

---

## AC-QUOTE-009 – Sales representative creates quote

**Given** A sales representative is authenticated AND has permission to create quotes AND has selected a valid, active customer
**When** The sales rep submits a new quote with at least one line item containing product, quantity, and price
**Then** The system creates a quote in 'Draft' status with: unique reference number (QT-XXXXX), default expiration of 30 days, calculated totals with taxes based on customer location, and creates an audit log entry

**Error Cases:**
- No line items provided → Show validation error 'At least one line item required'
- Customer is Inactive or Blocked → Show error 'Cannot create quote for inactive customer'
- Product is discontinued → Show warning but allow creation, flag quote for review
- Product is hard-deleted → Block line item addition, show error
- Line items exceed 100 → Show warning 'Large quote, consider splitting'
- Line items exceed 500 → Block creation, show error 'Maximum 500 line items per quote'

**Traceability:**
- Source: US-QUOTE-003
- Decision: AMB-OP-036
- Decision: AMB-OP-037
- Decision: AMB-OP-040
- Decision: AMB-OP-041
- Decision: AMB-OP-042
- Decision: AMB-ENT-015

---

## AC-QUOTE-010 – Sales representative can override prices within limits

**Given** A sales representative is creating or editing a quote line item
**When** The sales rep enters a price different from the product catalog price
**Then** The system allows the override if the discount is within configured limits (e.g., max 20%); if the discount exceeds limits, the quote is flagged for manager approval before it can be sent to customer

**Error Cases:**
- Discount exceeds configured limit without approval → Block transition to Pending status, show 'Manager approval required for discount of X%'
- Price is zero or negative → Show validation error 'Price must be greater than zero'

**Traceability:**
- Source: US-QUOTE-003
- Decision: AMB-OP-038

---

## AC-QUOTE-011 – Sales representative sends quote to customer

**Given** A quote is in 'Draft' status AND has valid customer and line items AND any required approvals are obtained
**When** The sales rep clicks 'Send to Customer'
**Then** The system: (1) changes quote status to 'Pending', (2) generates a unique UUID-based link token, (3) sends email to customer with quote summary, expiration date, and accept link using company-branded template with optional personalized message

**Error Cases:**
- Customer email format invalid → Show validation error, require correction before sending
- Email service unavailable → Queue email for retry every 15 minutes for up to 24 hours, notify sales rep of delay, quote status still changes to Pending
- Resend limit reached (5 times) → Block resend, show 'Maximum resends reached'
- Resend attempted within 1 hour of last send → Show 'Please wait before resending'

**Traceability:**
- Source: US-QUOTE-003
- Decision: AMB-OP-001
- Decision: AMB-OP-002
- Decision: AMB-OP-003
- Decision: AMB-OP-004
- Decision: AMB-OP-006
- Decision: AMB-ENT-005

---

## AC-QUOTE-012 – Quote automatically expires

**Given** A quote is in 'Pending' status AND the current date/time exceeds the quote expiration date
**When** The system expiration check runs (or customer attempts to view/accept)
**Then** The system automatically updates the quote status to 'Expired' AND sends notification to the sales rep

**Error Cases:**
- Expiration job fails → Log error, retry on next scheduled run, expired quotes blocked at view/accept time as fallback

**Traceability:**
- Source: US-QUOTE-003
- Decision: AMB-ENT-002
- Decision: AMB-OP-028

---

## AC-QUOTE-013 – Quote can be duplicated

**Given** A quote exists (any status) OR an order exists
**When** A sales rep selects 'Duplicate' on the quote or order
**Then** The system creates a new quote in 'Draft' status with: new reference number (QT-XXXXX), fresh 30-day expiration, copied line items with current product prices, assigned to the duplicating sales rep

**Error Cases:**
- Source quote/order not found → Show error 'Source not found'
- Products in source are now discontinued → Show warning, allow duplication without those products

**Traceability:**
- Source: US-QUOTE-003
- Decision: AMB-OP-039
- Decision: AMB-ENT-002

---

## AC-QUOTE-014 – Quote modification resets expiration and requires resend

**Given** A quote is in 'Draft' or 'Pending' status
**When** The assigned sales rep modifies line items, prices, discounts, or other quote details
**Then** The system: (1) saves the changes with full audit trail, (2) resets the expiration date to 30 days from modification, (3) if quote was Pending, invalidates the previous unique link AND shows prompt to resend to customer

**Error Cases:**
- Quote is in Accepted/Rejected/Expired/Cancelled status → Block modification, show 'Cannot modify quote in [status] status'
- Concurrent modification detected → Show conflict error, require refresh and re-apply changes

**Traceability:**
- Source: US-QUOTE-003
- Decision: AMB-ENT-003
- Decision: AMB-ENT-006
- Decision: AMB-ENT-022

---

## AC-QUOTE-015 – Customer rejects quote

**Given** A customer is viewing a valid, non-expired quote in 'Pending' status
**When** The customer clicks 'Reject Quote' (if provided) or explicitly declines
**Then** The system updates quote status to 'Rejected', records rejection timestamp, and notifies the assigned sales rep

**Error Cases:**
- Quote already accepted → Show 'Quote has already been accepted'
- Quote expired → Show expiration message

**Traceability:**
- Source: US-QUOTE-003
- Decision: AMB-ENT-001
- Decision: AMB-OP-028
- Decision: AMB-OP-029

---

## AC-QUOTE-016 – Sales representative cancels quote

**Given** A quote is in 'Draft' or 'Pending' status
**When** The assigned sales rep clicks 'Cancel Quote'
**Then** The system updates quote status to 'Cancelled', invalidates any active unique link, and logs the cancellation with timestamp and reason

**Error Cases:**
- Quote is Accepted → Block cancellation, show 'Cannot cancel accepted quote, cancel the order instead'
- Quote is already Cancelled/Expired → Show 'Quote is already [status]'

**Traceability:**
- Source: US-QUOTE-003
- Decision: AMB-ENT-001
- Decision: AMB-OP-028
- Decision: AMB-OP-029

---

## AC-QUOTE-017 – Quote soft deletion

**Given** A quote is in 'Draft' or 'Pending' status AND the user is the assigned sales rep or an admin
**When** The user deletes the quote
**Then** The system marks the quote as soft-deleted (not visible in normal queries) while retaining all data for audit purposes

**Error Cases:**
- Quote is Accepted → Block deletion, show 'Cannot delete accepted quote'
- User lacks permission → Show 'You do not have permission to delete this quote'

**Traceability:**
- Source: US-QUOTE-003
- Decision: AMB-ENT-004

---

## AC-QUOTE-018 – Email delivery tracking

**Given** A quote email has been sent to a customer
**When** The email service reports delivery status updates
**Then** The system tracks and stores: sent timestamp, delivery status (sent/delivered/bounced/opened), open events, and link click events; this data is visible to the sales rep on the quote detail page

**Error Cases:**
- Email not delivered within 1 hour → Alert sales rep
- Email bounced → Mark quote with delivery failure flag, notify sales rep, allow email correction and resend

**Traceability:**
- Source: US-QUOTE-003
- Decision: AMB-OP-005
- Decision: AMB-ENT-034
- Decision: AMB-ENT-035
- Decision: AMB-ENT-036

---

## AC-QUOTE-019 – Quote created via API

**Given** An authenticated API client with valid credentials and rate limit not exceeded
**When** The client sends a POST request to create a quote
**Then** The system creates the quote with the same validation rules as UI, returns the quote details with reference number, and logs the API call in the audit trail with source='API'

**Error Cases:**
- Authentication failed → Return 401 Unauthorized
- Rate limit exceeded → Return 429 Too Many Requests with retry-after header
- Validation errors → Return 400 Bad Request with detailed error messages

**Traceability:**
- Source: US-QUOTE-003
- Decision: AMB-OP-043

---

## AC-QUOTE-020 – Sales rep visibility based on role

**Given** A user is authenticated with SalesRepresentative role
**When** The user views the quotes list
**Then** The user sees only their own quotes by default; team leads see their team's quotes; managers see all quotes in their organization

**Error Cases:**
- Attempting to access quote outside visibility scope → Return 403 Forbidden

**Traceability:**
- Source: US-QUOTE-003
- Decision: AMB-ENT-028

---

## AC-QUOTE-021 – Quote reassignment on sales rep deactivation

**Given** A sales representative has pending quotes AND the sales rep is being deactivated
**When** The admin deactivates the sales rep account
**Then** The system automatically reassigns all pending quotes to the sales rep's manager; accepted quotes retain the original sales rep for historical reporting

**Error Cases:**
- Manager not defined → Escalate to admin for manual reassignment
- Manager is also inactive → Escalate to next level manager or admin

**Traceability:**
- Source: US-QUOTE-003
- Decision: AMB-ENT-029

---

