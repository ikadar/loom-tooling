# Technical Specifications

Generated: 2025-12-29T15:49:16+01:00

---

## TS-BR-CUST-001 – Order placement authorization

**Rule:** Only registered and email-verified customers can place orders

**Implementation Approach:**
Check customer registration status and email verification flag before order creation in order service

**Validation Points:**
- Order creation API
- Checkout flow initiation

**Data Requirements:**
| Field | Type | Constraints | Source |
|-------|------|-------------|--------|
| customerId | UUID | NOT NULL, valid reference | BR-CUST-001 |
| registrationStatus | enum | = 'registered' | BR-CUST-001 |
| emailVerified | boolean | = true | BR-CUST-001 |

**Error Handling:**
| Condition | Error Code | Message | HTTP Status |
|-----------|------------|---------|-------------|
| registrationStatus != 'registered' | UNAUTHORIZED_ORDER | Registration required to place orders | 403 |
| emailVerified != true | EMAIL_NOT_VERIFIED | Please verify your email before ordering | 403 |

**Traceability:**
- BR: BR-CUST-001
- Related ACs: [AC-CUST-002 AC-ORDER-001]

---

## TS-BR-CUST-002 – Unique email constraint

**Rule:** Each customer account must have a unique email address

**Implementation Approach:**
Add UNIQUE constraint on email column, check before insert in registration service

**Validation Points:**
- Customer registration API
- Email update API

**Data Requirements:**
| Field | Type | Constraints | Source |
|-------|------|-------------|--------|
| email | varchar(255) | UNIQUE, NOT NULL, valid email format | BR-CUST-002 |

**Error Handling:**
| Condition | Error Code | Message | HTTP Status |
|-----------|------------|---------|-------------|
| email already exists in database | DUPLICATE_EMAIL | Email already registered. Please login or reset password | 409 |

**Traceability:**
- BR: BR-CUST-002
- Related ACs: [AC-CUST-001]

---

## TS-BR-CUST-003 – Password validation rules

**Rule:** Customer passwords must be at least 8 characters with at least one number

**Implementation Approach:**
Regex validation pattern: ^(?=.*[0-9]).{8,}$ applied before password hashing

**Validation Points:**
- Registration API
- Password change API
- Password reset API

**Data Requirements:**
| Field | Type | Constraints | Source |
|-------|------|-------------|--------|
| password | varchar(255) | length >= 8, contains digit | BR-CUST-003 |

**Error Handling:**
| Condition | Error Code | Message | HTTP Status |
|-----------|------------|---------|-------------|
| password.length < 8 | INVALID_PASSWORD | Password must be at least 8 characters | 400 |
| !password.match(/[0-9]/) | INVALID_PASSWORD | Password must contain at least one number | 400 |

**Traceability:**
- BR: BR-CUST-003
- Related ACs: [AC-CUST-001]

---

## TS-BR-PROD-001 – Product price validation

**Rule:** Product prices must be positive with minimum $0.01

**Implementation Approach:**
Decimal validation with 2 decimal places, minimum value check in product service

**Validation Points:**
- Product creation API
- Product edit API

**Data Requirements:**
| Field | Type | Constraints | Source |
|-------|------|-------------|--------|
| price | decimal(10,2) | >= 0.01 | BR-PROD-001 |

**Error Handling:**
| Condition | Error Code | Message | HTTP Status |
|-----------|------------|---------|-------------|
| price < 0.01 | INVALID_PRICE | Price must be at least $0.01 | 400 |
| price has more than 2 decimals | INVALID_PRICE | Price must have exactly 2 decimal places | 400 |

**Traceability:**
- BR: BR-PROD-001
- Related ACs: [AC-ADMIN-001]

---

## TS-BR-PROD-002 – Product name length validation

**Rule:** Product names must be between 2 and 200 characters

**Implementation Approach:**
String length validation in product service before database insert/update

**Validation Points:**
- Product creation API
- Product edit API

**Data Requirements:**
| Field | Type | Constraints | Source |
|-------|------|-------------|--------|
| name | varchar(200) | 2 <= length <= 200, NOT NULL | BR-PROD-002 |

**Error Handling:**
| Condition | Error Code | Message | HTTP Status |
|-----------|------------|---------|-------------|
| name.length < 2 || name is empty | NAME_REQUIRED | Product name is required | 400 |
| name.length > 200 | NAME_TOO_LONG | Product name cannot exceed 200 characters | 400 |

**Traceability:**
- BR: BR-PROD-002
- Related ACs: [AC-ADMIN-001]

---

## TS-BR-PROD-003 – Product soft delete implementation

**Rule:** Products with order history must be soft-deleted, not permanently deleted

**Implementation Approach:**
Check OrderLineItem references before delete; set isDeleted=true and deletedAt timestamp instead of physical delete

**Validation Points:**
- Product delete API

**Data Requirements:**
| Field | Type | Constraints | Source |
|-------|------|-------------|--------|
| isDeleted | boolean | default false | BR-PROD-003 |
| deletedAt | timestamp | nullable | BR-PROD-003 |

**Traceability:**
- BR: BR-PROD-003
- Related ACs: [AC-ADMIN-005]

---

## TS-BR-PROD-004 – Variant SKU uniqueness

**Rule:** Each product variant must have a unique SKU

**Implementation Approach:**
Add UNIQUE constraint on sku column in product_variants table

**Validation Points:**
- Variant creation API
- Variant edit API
- Bulk import

**Data Requirements:**
| Field | Type | Constraints | Source |
|-------|------|-------------|--------|
| sku | varchar(50) | UNIQUE, NOT NULL | BR-PROD-004 |

**Error Handling:**
| Condition | Error Code | Message | HTTP Status |
|-----------|------------|---------|-------------|
| sku already exists | DUPLICATE_SKU | SKU already exists | 409 |

**Traceability:**
- BR: BR-PROD-004
- Related ACs: [AC-PROD-003]

---

## TS-BR-CART-001 – Stock availability for cart add

**Rule:** Products must be in stock to be added to cart

**Implementation Approach:**
Query inventory.availableQuantity before cart item insert; reject if zero

**Validation Points:**
- Add to cart API
- Cart UI add button

**Data Requirements:**
| Field | Type | Constraints | Source |
|-------|------|-------------|--------|
| productId | UUID | valid product reference | BR-CART-001 |
| availableQuantity | integer | > 0 for add | BR-CART-001 |

**Error Handling:**
| Condition | Error Code | Message | HTTP Status |
|-----------|------------|---------|-------------|
| availableQuantity == 0 | OUT_OF_STOCK | Product is out of stock | 400 |

**Traceability:**
- BR: BR-CART-001
- Related ACs: [AC-CART-001]

---

## TS-BR-CART-002 – Cart quantity stock limit

**Rule:** Cart item quantity cannot exceed available stock

**Implementation Approach:**
Cap requested quantity to availableQuantity; return notification if capped

**Validation Points:**
- Add to cart API
- Update quantity API

**Data Requirements:**
| Field | Type | Constraints | Source |
|-------|------|-------------|--------|
| quantity | integer | <= availableQuantity | BR-CART-002 |
| availableQuantity | integer | from inventory table | BR-CART-002 |

**Error Handling:**
| Condition | Error Code | Message | HTTP Status |
|-----------|------------|---------|-------------|
| requestedQty > availableQuantity | QUANTITY_LIMITED | Only {available} available | 200 |

**Traceability:**
- BR: BR-CART-002
- Related ACs: [AC-CART-001 AC-CART-002 AC-CART-004]

---

## TS-BR-CART-003 – Positive cart quantity enforcement

**Rule:** Cart item quantities must be at least 1

**Implementation Approach:**
If quantity set to 0 or less, remove item from cart instead of rejecting

**Validation Points:**
- Update quantity API

**Data Requirements:**
| Field | Type | Constraints | Source |
|-------|------|-------------|--------|
| quantity | integer | >= 1 or triggers removal | BR-CART-003 |

**Traceability:**
- BR: BR-CART-003
- Related ACs: [AC-CART-004]

---

## TS-BR-CART-004 – Cart expiration mechanism

**Rule:** Inactive carts expire after 30 days (logged-in) or 7 days (guest)

**Implementation Approach:**
Background job runs daily to clear expired carts; lazy check on cart access

**Validation Points:**
- Cart access API
- Background scheduler

**Data Requirements:**
| Field | Type | Constraints | Source |
|-------|------|-------------|--------|
| lastActivityDate | timestamp | NOT NULL, updated on cart ops | BR-CART-004 |
| isGuest | boolean | determines expiry threshold | BR-CART-004 |

**Error Handling:**
| Condition | Error Code | Message | HTTP Status |
|-----------|------------|---------|-------------|
| cart expired on access | CART_EXPIRED | Your cart has expired | 410 |

**Traceability:**
- BR: BR-CART-004
- Related ACs: [AC-CART-009]

---

## TS-BR-CART-005 – Variant selection requirement

**Rule:** For products with variants, a specific variant must be selected before adding to cart

**Implementation Approach:**
Check if product.hasVariants; require variantId in add-to-cart if true

**Validation Points:**
- Add to cart API
- Product detail UI

**Data Requirements:**
| Field | Type | Constraints | Source |
|-------|------|-------------|--------|
| variantId | UUID | required if product has variants | BR-CART-005 |
| hasVariants | boolean | from product | BR-CART-005 |

**Error Handling:**
| Condition | Error Code | Message | HTTP Status |
|-----------|------------|---------|-------------|
| hasVariants && !variantId | VARIANT_REQUIRED | Please select a variant | 400 |

**Traceability:**
- BR: BR-CART-005
- Related ACs: [AC-PROD-003 AC-CART-001]

---

## TS-BR-ORDER-001 – Free shipping calculation

**Rule:** Orders over $50 qualify for free shipping

**Implementation Approach:**
If subtotal > 50.00, set shippingCost = 0.00; else apply standard rate

**Validation Points:**
- Checkout totals calculation
- Order creation

**Data Requirements:**
| Field | Type | Constraints | Source |
|-------|------|-------------|--------|
| subtotal | decimal(10,2) | sum of line items | BR-ORDER-001 |
| shippingCost | decimal(10,2) | 0.00 if subtotal > 50 | BR-ORDER-001 |
| freeShippingThreshold | decimal | config value = 50.00 | BR-ORDER-001 |

**Traceability:**
- BR: BR-ORDER-001
- Related ACs: [AC-ORDER-001 AC-ORDER-002]

---

## TS-BR-ORDER-002 – Order cancellation status check

**Rule:** Orders can only be cancelled before shipping

**Implementation Approach:**
Verify order.status in ('pending', 'confirmed') before allowing cancel

**Validation Points:**
- Cancel order API

**Data Requirements:**
| Field | Type | Constraints | Source |
|-------|------|-------------|--------|
| status | enum | pending|confirmed for cancel | BR-ORDER-002 |

**Error Handling:**
| Condition | Error Code | Message | HTTP Status |
|-----------|------------|---------|-------------|
| status in ('shipped', 'delivered', 'cancelled') | ORDER_ALREADY_SHIPPED | Cannot cancel shipped order | 400 |

**Traceability:**
- BR: BR-ORDER-002
- Related ACs: [AC-ORDER-008 AC-ORDER-009 AC-ORDER-010]

---

## TS-BR-ORDER-003 – Order status state machine

**Rule:** Order status can only transition through valid states

**Implementation Approach:**
State machine with allowed transitions map; validate before update

**Validation Points:**
- Update order status API

**Data Requirements:**
| Field | Type | Constraints | Source |
|-------|------|-------------|--------|
| currentStatus | enum | valid order status | BR-ORDER-003 |
| newStatus | enum | valid transition from current | BR-ORDER-003 |

**Error Handling:**
| Condition | Error Code | Message | HTTP Status |
|-----------|------------|---------|-------------|
| invalid transition | INVALID_STATUS_TRANSITION | Cannot transition from {current} to {new} | 400 |

**Traceability:**
- BR: BR-ORDER-003
- Related ACs: [AC-ADMIN-008]

---

## TS-BR-ORDER-004 – Order immutability after shipping

**Rule:** Orders cannot be modified after shipping

**Implementation Approach:**
Check status before any modification operation; reject if shipped/delivered/cancelled

**Validation Points:**
- All order modification endpoints

**Data Requirements:**
| Field | Type | Constraints | Source |
|-------|------|-------------|--------|
| status | enum | must not be shipped/delivered/cancelled | BR-ORDER-004 |

**Error Handling:**
| Condition | Error Code | Message | HTTP Status |
|-----------|------------|---------|-------------|
| status in ('shipped', 'delivered', 'cancelled') | ORDER_NOT_MODIFIABLE | Order cannot be modified | 400 |

**Traceability:**
- BR: BR-ORDER-004
- Related ACs: [AC-ORDER-008 AC-ORDER-010]

---

## TS-BR-ORDER-005 – Minimum line item validation

**Rule:** Orders must contain at least one line item

**Implementation Approach:**
Validate lineItems array length >= 1 before order creation

**Validation Points:**
- Place order API

**Data Requirements:**
| Field | Type | Constraints | Source |
|-------|------|-------------|--------|
| lineItems | array | length >= 1 | BR-ORDER-005 |

**Error Handling:**
| Condition | Error Code | Message | HTTP Status |
|-----------|------------|---------|-------------|
| lineItems.length == 0 | EMPTY_ORDER | Order must contain at least one item | 400 |

**Traceability:**
- BR: BR-ORDER-005
- Related ACs: [AC-ORDER-001]

---

## TS-BR-ORDER-006 – Price snapshot implementation

**Rule:** Order line items must capture the price at order placement, not current price

**Implementation Approach:**
Copy product.price to orderLineItem.unitPrice at order creation; store as immutable

**Validation Points:**
- Place order API

**Data Requirements:**
| Field | Type | Constraints | Source |
|-------|------|-------------|--------|
| unitPrice | decimal(10,2) | copied from product, immutable | BR-ORDER-006 |
| productPriceAtOrder | decimal(10,2) | snapshot value | BR-ORDER-006 |

**Traceability:**
- BR: BR-ORDER-006
- Related ACs: [AC-ORDER-004]

---

## TS-BR-INV-001 – Non-negative inventory enforcement

**Rule:** Inventory stock level cannot go below zero

**Implementation Approach:**
Database CHECK constraint; validate before adjustment operations

**Validation Points:**
- Inventory adjustment API
- Order placement stock decrement

**Data Requirements:**
| Field | Type | Constraints | Source |
|-------|------|-------------|--------|
| stockLevel | integer | >= 0, CHECK constraint | BR-INV-001 |

**Error Handling:**
| Condition | Error Code | Message | HTTP Status |
|-----------|------------|---------|-------------|
| adjustment would result in negative | INSUFFICIENT_STOCK | Insufficient stock available | 400 |

**Traceability:**
- BR: BR-INV-001
- Related ACs: [AC-ADMIN-009 AC-ADMIN-010]

---

## TS-BR-INV-002 – Stock validation at checkout

**Rule:** Products cannot be sold when out of stock (no backorders)

**Implementation Approach:**
Validate all cart items have sufficient stock atomically before order creation

**Validation Points:**
- Place order API

**Data Requirements:**
| Field | Type | Constraints | Source |
|-------|------|-------------|--------|
| requestedQuantity | integer | from cart item | BR-INV-002 |
| availableQuantity | integer | >= requestedQuantity | BR-INV-002 |

**Error Handling:**
| Condition | Error Code | Message | HTTP Status |
|-----------|------------|---------|-------------|
| availableQuantity < requestedQuantity | STOCK_UNAVAILABLE | {product} has insufficient stock | 400 |

**Traceability:**
- BR: BR-INV-002
- Related ACs: [AC-ORDER-001]

---

## TS-BR-INV-003 – Inventory restoration on cancel

**Rule:** Cancelled orders must restore inventory automatically

**Implementation Approach:**
In cancel operation, iterate line items and increment stock by quantity

**Validation Points:**
- Cancel order API

**Data Requirements:**
| Field | Type | Constraints | Source |
|-------|------|-------------|--------|
| lineItems | array | order line items with quantities | BR-INV-003 |
| stockLevel | integer | incremented by cancelled qty | BR-INV-003 |

**Traceability:**
- BR: BR-INV-003
- Related ACs: [AC-ORDER-008 AC-ORDER-009]

---

## TS-BR-INV-004 – Inventory audit logging

**Rule:** All inventory changes must be logged with before/after values and reason

**Implementation Approach:**
Create audit log entry in same transaction as inventory change

**Validation Points:**
- All inventory modification operations

**Data Requirements:**
| Field | Type | Constraints | Source |
|-------|------|-------------|--------|
| previousValue | integer | stock before change | BR-INV-004 |
| newValue | integer | stock after change | BR-INV-004 |
| reason | varchar(255) | NOT NULL | BR-INV-004 |
| userId | UUID | user who made change | BR-INV-004 |
| timestamp | timestamp | auto-generated | BR-INV-004 |

**Traceability:**
- BR: BR-INV-004
- Related ACs: [AC-ADMIN-009 AC-ADMIN-010]

---

## TS-BR-PAY-001 – Payment authorization flow

**Rule:** Payment must be authorized before order is created

**Implementation Approach:**
Call payment gateway authorization; only create order if authorization succeeds

**Validation Points:**
- Place order API - before order insert

**Data Requirements:**
| Field | Type | Constraints | Source |
|-------|------|-------------|--------|
| paymentMethodId | UUID | valid payment method | BR-PAY-001 |
| amount | decimal(10,2) | order total | BR-PAY-001 |
| authorizationId | varchar(100) | from gateway response | BR-PAY-001 |

**Error Handling:**
| Condition | Error Code | Message | HTTP Status |
|-----------|------------|---------|-------------|
| authorization fails | PAYMENT_FAILED | Payment authorization failed | 402 |

**Traceability:**
- BR: BR-PAY-001
- Related ACs: [AC-ORDER-001 AC-PAY-001]

---

## TS-BR-PAY-002 – Credit card type validation

**Rule:** Only Visa, Mastercard, and American Express credit cards are accepted

**Implementation Approach:**
Detect card type from BIN; reject if not in allowed list

**Validation Points:**
- Payment method creation
- Checkout payment

**Data Requirements:**
| Field | Type | Constraints | Source |
|-------|------|-------------|--------|
| cardType | enum | visa|mastercard|amex | BR-PAY-002 |
| cardNumber | varchar(19) | for BIN detection | BR-PAY-002 |

**Error Handling:**
| Condition | Error Code | Message | HTTP Status |
|-----------|------------|---------|-------------|
| cardType not in allowed list | UNSUPPORTED_CARD_TYPE | Card type not accepted. Use Visa, Mastercard, or Amex | 400 |

**Traceability:**
- BR: BR-PAY-002
- Related ACs: [AC-PAY-001]

---

## TS-BR-PAY-003 – Automatic refund on cancellation

**Rule:** Cancelled orders must be refunded to the original payment method

**Implementation Approach:**
Call payment gateway refund API with original transaction reference

**Validation Points:**
- Cancel order API - after status change

**Data Requirements:**
| Field | Type | Constraints | Source |
|-------|------|-------------|--------|
| originalTransactionId | varchar(100) | from order payment | BR-PAY-003 |
| refundAmount | decimal(10,2) | order total | BR-PAY-003 |

**Error Handling:**
| Condition | Error Code | Message | HTTP Status |
|-----------|------------|---------|-------------|
| refund API fails | REFUND_FAILED | Refund processing failed | 500 |

**Traceability:**
- BR: BR-PAY-003
- Related ACs: [AC-ORDER-008]

---

## TS-BR-SHIP-001 – Domestic shipping validation

**Rule:** Shipping is only available within the domestic country

**Implementation Approach:**
Compare address country to configured domestic country code

**Validation Points:**
- Address save API
- Checkout shipping step

**Data Requirements:**
| Field | Type | Constraints | Source |
|-------|------|-------------|--------|
| country | varchar(2) | ISO country code | BR-SHIP-001 |
| domesticCountry | varchar(2) | config value | BR-SHIP-001 |

**Error Handling:**
| Condition | Error Code | Message | HTTP Status |
|-----------|------------|---------|-------------|
| country != domesticCountry | UNSUPPORTED_SHIPPING_REGION | Shipping not available to this country | 400 |

**Traceability:**
- BR: BR-SHIP-001
- Related ACs: [AC-CUST-003]

---

## TS-BR-SHIP-002 – Required address fields validation

**Rule:** Shipping addresses must include all required fields

**Implementation Approach:**
Validate all required fields are non-empty before save

**Validation Points:**
- Address save API
- Checkout shipping step

**Data Requirements:**
| Field | Type | Constraints | Source |
|-------|------|-------------|--------|
| street | varchar(255) | NOT NULL, NOT EMPTY | BR-SHIP-002 |
| city | varchar(100) | NOT NULL, NOT EMPTY | BR-SHIP-002 |
| state | varchar(100) | NOT NULL, NOT EMPTY | BR-SHIP-002 |
| postalCode | varchar(20) | NOT NULL, NOT EMPTY | BR-SHIP-002 |
| country | varchar(2) | NOT NULL, NOT EMPTY | BR-SHIP-002 |
| recipientName | varchar(200) | NOT NULL, NOT EMPTY | BR-SHIP-002 |

**Error Handling:**
| Condition | Error Code | Message | HTTP Status |
|-----------|------------|---------|-------------|
| any required field empty | INVALID_ADDRESS | Missing required field: {fieldName} | 400 |

**Traceability:**
- BR: BR-SHIP-002
- Related ACs: [AC-CUST-003]

---

## TS-BR-CAT-001 – Category nesting depth limit

**Rule:** Category hierarchy is limited to 3 levels deep

**Implementation Approach:**
Calculate depth by traversing parent chain; reject if would exceed 3

**Validation Points:**
- Category creation API
- Category move API

**Data Requirements:**
| Field | Type | Constraints | Source |
|-------|------|-------------|--------|
| parentId | UUID | nullable, valid category ref | BR-CAT-001 |
| depth | integer | <= 3 | BR-CAT-001 |

**Error Handling:**
| Condition | Error Code | Message | HTTP Status |
|-----------|------------|---------|-------------|
| calculated depth > 3 | MAX_NESTING_EXCEEDED | Category nesting limited to 3 levels | 400 |

**Traceability:**
- BR: BR-CAT-001
- Related ACs: [AC-ADMIN-013]

---

## TS-BR-EMAIL-001 – Async email delivery

**Rule:** Email delivery failures must not block order completion

**Implementation Approach:**
Publish email event to message queue; process async with retry logic

**Validation Points:**
- Order creation - after commit

**Data Requirements:**
| Field | Type | Constraints | Source |
|-------|------|-------------|--------|
| orderId | UUID | for email content | BR-EMAIL-001 |
| customerEmail | varchar(255) | recipient | BR-EMAIL-001 |
| emailType | enum | order_confirmation | BR-EMAIL-001 |

**Traceability:**
- BR: BR-EMAIL-001
- Related ACs: [AC-ORDER-011]

---

## TS-BR-EMAIL-002 – Email deduplication

**Rule:** Duplicate email sends must be prevented

**Implementation Approach:**
Check email_sent table for orderId+emailType before sending; mark sent atomically

**Validation Points:**
- Email processor - before send

**Data Requirements:**
| Field | Type | Constraints | Source |
|-------|------|-------------|--------|
| orderId | UUID | unique with emailType | BR-EMAIL-002 |
| emailType | enum | part of unique key | BR-EMAIL-002 |
| sentAt | timestamp | when email was sent | BR-EMAIL-002 |

**Traceability:**
- BR: BR-EMAIL-002
- Related ACs: [AC-ORDER-011]

---

