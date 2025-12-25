# US-QUOTE-003 – Customer accepts a quote online

**As a** customer
**I want** to accept a quote online
**So that** I can confirm the order quickly without paperwork.

## Acceptance Notes

- Customer receives quote via email with unique link
- Quote has expiration date (default 30 days)
- Customer can view quote details before accepting
- Acceptance requires clicking "Accept Quote" button
- System should record acceptance timestamp and IP address
- Customer receives confirmation email after acceptance
- Quote status changes from "Pending" to "Accepted"
- Accepted quote triggers order creation process

## Business Context

- Quotes are created by sales representatives
- Each quote has a unique reference number (QT-XXXXX)
- Quote contains line items with products, quantities, prices
- Total value includes applicable taxes and discounts
- Customer may have multiple pending quotes
