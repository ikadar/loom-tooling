# Technical Specifications

Generated: 2025-12-31T12:00:34+01:00

---

## TS-BR-CUST-001 – Customer registration validation for orders

**Rule:** Only registered and email-verified customers can place orders

**Implementation Approach:**
Validate customer registrationStatus and emailVerified fields in order service before order creation

**Validation Points:**
- Place Order API
- Checkout flow

**Data Requirements:**
| Field | Type | Constraints | Source |
|-------|------|-------------|--------|
| registrationStatus | enum | must equal 'registered' | Customer entity |
| emailVerified | boolean | must be true | Customer entity |

**Error Handling:**
| Condition | Error Code | Message | HTTP Status |
|-----------|------------|---------|-------------|
| registrationStatus != 'registered' OR emailVerified == false | UNAUTHORIZED_ORDER | Only registered and email-verified customers can place orders | 403 |

**Traceability:**
- BR: BR-CUST-001

---

## TS-BR-CUST-002 – Unique email constraint for customers

**Rule:** Each customer account must have a unique email address

**Implementation Approach:**
Database unique constraint on email column, application-level check during registration

**Validation Points:**
- Registration API
- Email update API

**Data Requirements:**
| Field | Type | Constraints | Source |
|-------|------|-------------|--------|
| email | string | unique, valid email format | Customer entity |

**Error Handling:**
| Condition | Error Code | Message | HTTP Status |
|-----------|------------|---------|-------------|
| email already exists in database | DUPLICATE_EMAIL | An account with this email address already exists | 409 |

**Traceability:**
- BR: BR-CUST-002

---

## TS-BR-CUST-003 – Password strength validation

**Rule:** Customer passwords must be at least 8 characters with at least one number

**Implementation Approach:**
Regex validation: ^(?=.*[0-9]).{8,}$ applied during registration and password change

**Validation Points:**
- Registration API
- Password change API
- Password reset API

**Data Requirements:**
| Field | Type | Constraints | Source |
|-------|------|-------------|--------|
| password | string | length >= 8, contains at least one digit [0-9] | Customer entity |

**Error Handling:**
| Condition | Error Code | Message | HTTP Status |
|-----------|------------|---------|-------------|
| password.length < 8 OR !password.matches(/[0-9]/) | INVALID_PASSWORD | Password must be at least 8 characters and contain at least one number | 400 |

**Traceability:**
- BR: BR-CUST-003

---

## TS-BR-PROD-001 – Product price validation

**Rule:** Product prices must be positive with minimum $0.01

**Implementation Approach:**
Validate price >= 0.01 and has exactly 2 decimal places during product creation and update

**Validation Points:**
- Create Product API
- Update Product API
- Admin UI

**Data Requirements:**
| Field | Type | Constraints | Source |
|-------|------|-------------|--------|
| price | decimal(10,2) | >= 0.01, exactly 2 decimal places | Product entity |

**Error Handling:**
| Condition | Error Code | Message | HTTP Status |
|-----------|------------|---------|-------------|
| price < 0.01 OR decimal places != 2 | INVALID_PRICE | Price must be at least $0.01 with exactly 2 decimal places | 400 |

**Traceability:**
- BR: BR-PROD-001

---

## TS-BR-PROD-002 – Product name length validation

**Rule:** Product names must be between 2 and 200 characters

**Implementation Approach:**
String length validation on product name field during creation and update

**Validation Points:**
- Create Product API
- Update Product API
- Admin UI

**Data Requirements:**
| Field | Type | Constraints | Source |
|-------|------|-------------|--------|
| name | string | 2 <= length <= 200 | Product entity |

**Error Handling:**
| Condition | Error Code | Message | HTTP Status |
|-----------|------------|---------|-------------|
| name.length < 2 OR name.length > 200 | INVALID_NAME | Product name must be between 2 and 200 characters | 400 |

**Traceability:**
- BR: BR-PROD-002

---

## TS-BR-PROD-003 – Product soft delete with order references

**Rule:** Products with order history must be soft-deleted, not permanently deleted

**Implementation Approach:**
Check for OrderLineItem references before delete; if exists, set deletedAt timestamp instead of hard delete

**Validation Points:**
- Delete Product API

**Data Requirements:**
| Field | Type | Constraints | Source |
|-------|------|-------------|--------|
| deletedAt | timestamp | nullable, set on soft delete | Product entity |
| orderLineItems | relation | check for existence | OrderLineItem entity |

**Traceability:**
- BR: BR-PROD-003

---

## TS-BR-PROD-004 – Unique SKU for product variants

**Rule:** Each product variant must have a unique SKU

**Implementation Approach:**
Database unique constraint on sku column, application-level check during variant creation

**Validation Points:**
- Create Variant API
- Update Variant API

**Data Requirements:**
| Field | Type | Constraints | Source |
|-------|------|-------------|--------|
| sku | string | unique across all variants | ProductVariant entity |

**Error Handling:**
| Condition | Error Code | Message | HTTP Status |
|-----------|------------|---------|-------------|
| sku already exists in database | DUPLICATE_SKU | A variant with this SKU already exists | 409 |

**Traceability:**
- BR: BR-PROD-004

---

## TS-BR-CART-001 – Stock availability check for cart

**Rule:** Products must be in stock to be added to cart

**Implementation Approach:**
Query Inventory.availableQuantity for product before adding to cart

**Validation Points:**
- Add to Cart API
- Cart UI

**Data Requirements:**
| Field | Type | Constraints | Source |
|-------|------|-------------|--------|
| availableQuantity | integer | > 0 | Inventory entity |
| productId | uuid | must reference valid product | CartItem entity |

**Error Handling:**
| Condition | Error Code | Message | HTTP Status |
|-----------|------------|---------|-------------|
| availableQuantity <= 0 | OUT_OF_STOCK | This product is currently out of stock | 400 |

**Traceability:**
- BR: BR-CART-001

---

## TS-BR-CART-002 – Cart quantity limited by stock

**Rule:** Cart item quantity cannot exceed available stock

**Implementation Approach:**
Compare requested quantity against Inventory.availableQuantity; cap or reject if exceeded

**Validation Points:**
- Add to Cart API
- Update Cart Quantity API

**Data Requirements:**
| Field | Type | Constraints | Source |
|-------|------|-------------|--------|
| quantity | integer | <= availableQuantity | CartItem entity |
| availableQuantity | integer | from Inventory | Inventory entity |

**Error Handling:**
| Condition | Error Code | Message | HTTP Status |
|-----------|------------|---------|-------------|
| quantity > availableQuantity | QUANTITY_LIMITED | Requested quantity exceeds available stock. Maximum available: {availableQuantity} | 400 |

**Traceability:**
- BR: BR-CART-002

---

## TS-BR-CART-003 – Positive cart quantity validation

**Rule:** Cart item quantities must be at least 1

**Implementation Approach:**
Validate quantity >= 1; if quantity set to 0, remove item from cart

**Validation Points:**
- Add to Cart API
- Update Cart Quantity API

**Data Requirements:**
| Field | Type | Constraints | Source |
|-------|------|-------------|--------|
| quantity | integer | >= 1 (0 triggers removal) | CartItem entity |

**Traceability:**
- BR: BR-CART-003

---

## TS-BR-CART-004 – Cart expiration handling

**Rule:** Inactive carts expire after 30 days (logged-in) or 7 days (guest)

**Implementation Approach:**
Background job or lazy evaluation on cart access checks lastActivityDate against threshold

**Validation Points:**
- Cart access API
- Background cleanup job

**Data Requirements:**
| Field | Type | Constraints | Source |
|-------|------|-------------|--------|
| lastActivityDate | timestamp | updated on cart interaction | Cart entity |
| isGuest | boolean | determines expiration threshold | Cart entity |

**Error Handling:**
| Condition | Error Code | Message | HTTP Status |
|-----------|------------|---------|-------------|
| cart expired based on lastActivityDate | CART_EXPIRED | Your cart has expired due to inactivity | 410 |

**Traceability:**
- BR: BR-CART-004

---

## TS-BR-CART-005 – Variant selection requirement

**Rule:** For products with variants, a specific variant must be selected before adding to cart

**Implementation Approach:**
Check if product has variants; if so, require variantId in add-to-cart request

**Validation Points:**
- Add to Cart API
- Product UI

**Data Requirements:**
| Field | Type | Constraints | Source |
|-------|------|-------------|--------|
| variantId | uuid | required if product has variants | CartItem entity |
| hasVariants | boolean | derived from Product.variants.length > 0 | Product entity |

**Error Handling:**
| Condition | Error Code | Message | HTTP Status |
|-----------|------------|---------|-------------|
| product has variants AND variantId not provided | VARIANT_REQUIRED | Please select a specific variant before adding to cart | 400 |

**Traceability:**
- BR: BR-CART-005

---

## TS-BR-ORDER-001 – Free shipping threshold calculation

**Rule:** Orders over $50 qualify for free shipping

**Implementation Approach:**
Calculate shipping cost as $0.00 when order subtotal exceeds $50.00

**Validation Points:**
- Checkout flow
- Order creation

**Data Requirements:**
| Field | Type | Constraints | Source |
|-------|------|-------------|--------|
| subtotal | decimal(10,2) | sum of line item totals | Order entity |
| shippingCost | decimal(10,2) | $0.00 if subtotal > $50.00 | Order entity |

**Traceability:**
- BR: BR-ORDER-001

---

## TS-BR-ORDER-002 – Order cancellation status check

**Rule:** Orders can only be cancelled before shipping

**Implementation Approach:**
Check order status is 'pending' or 'confirmed' before allowing cancellation

**Validation Points:**
- Cancel Order API

**Data Requirements:**
| Field | Type | Constraints | Source |
|-------|------|-------------|--------|
| status | enum | must be 'pending' or 'confirmed' for cancellation | Order entity |

**Error Handling:**
| Condition | Error Code | Message | HTTP Status |
|-----------|------------|---------|-------------|
| status NOT IN ('pending', 'confirmed') | ORDER_ALREADY_SHIPPED | Cannot cancel order that has already been shipped | 400 |

**Traceability:**
- BR: BR-ORDER-002

---

## TS-BR-ORDER-003 – Order status state machine

**Rule:** Order status can only transition through valid states

**Implementation Approach:**
Implement state machine validating transitions: pending→confirmed→shipped→delivered, pending/confirmed→cancelled

**Validation Points:**
- Update Order Status API

**Data Requirements:**
| Field | Type | Constraints | Source |
|-------|------|-------------|--------|
| status | enum | pending|confirmed|shipped|delivered|cancelled | Order entity |
| previousStatus | enum | current status before transition | Order entity |

**Error Handling:**
| Condition | Error Code | Message | HTTP Status |
|-----------|------------|---------|-------------|
| invalid status transition attempted | INVALID_STATUS_TRANSITION | Cannot transition order from {currentStatus} to {requestedStatus} | 400 |

**Traceability:**
- BR: BR-ORDER-003

---

## TS-BR-ORDER-004 – Order immutability after shipping

**Rule:** Orders cannot be modified after shipping

**Implementation Approach:**
Check order status before any modification; reject if shipped, delivered, or cancelled

**Validation Points:**
- Update Order API
- All order modification endpoints

**Data Requirements:**
| Field | Type | Constraints | Source |
|-------|------|-------------|--------|
| status | enum | must NOT be 'shipped', 'delivered', or 'cancelled' | Order entity |

**Error Handling:**
| Condition | Error Code | Message | HTTP Status |
|-----------|------------|---------|-------------|
| status IN ('shipped', 'delivered', 'cancelled') | ORDER_NOT_MODIFIABLE | Order cannot be modified in its current status | 400 |

**Traceability:**
- BR: BR-ORDER-004

---

## TS-BR-ORDER-005 – Minimum line item validation

**Rule:** Orders must contain at least one line item

**Implementation Approach:**
Validate lineItems array has at least one element during order creation

**Validation Points:**
- Place Order API

**Data Requirements:**
| Field | Type | Constraints | Source |
|-------|------|-------------|--------|
| lineItems | array | length >= 1 | Order entity |

**Error Handling:**
| Condition | Error Code | Message | HTTP Status |
|-----------|------------|---------|-------------|
| lineItems.length == 0 | EMPTY_ORDER | Order must contain at least one item | 400 |

**Traceability:**
- BR: BR-ORDER-005

---

## TS-BR-ORDER-006 – Price snapshot at order time

**Rule:** Order line items must capture the price at order placement, not current price

**Implementation Approach:**
Copy Product.price to OrderLineItem.unitPrice during order creation; field is immutable after

**Validation Points:**
- Place Order API

**Data Requirements:**
| Field | Type | Constraints | Source |
|-------|------|-------------|--------|
| unitPrice | decimal(10,2) | copied from Product.price at order time, immutable | OrderLineItem entity |

**Traceability:**
- BR: BR-ORDER-006

---

## TS-BR-INV-001 – Non-negative inventory validation

**Rule:** Inventory stock level cannot go below zero

**Implementation Approach:**
Validate stockLevel >= 0 during all stock adjustments; use database constraint

**Validation Points:**
- Manage Inventory API
- Order fulfillment
- Stock adjustment

**Data Requirements:**
| Field | Type | Constraints | Source |
|-------|------|-------------|--------|
| stockLevel | integer | >= 0 | Inventory entity |

**Error Handling:**
| Condition | Error Code | Message | HTTP Status |
|-----------|------------|---------|-------------|
| adjustment would result in stockLevel < 0 | INSUFFICIENT_STOCK | Insufficient stock for this operation | 400 |

**Traceability:**
- BR: BR-INV-001

---

## TS-BR-INV-002 – No backorder validation at checkout

**Rule:** Products cannot be sold when out of stock (no backorders)

**Implementation Approach:**
Validate availableQuantity >= requested quantity for all line items at checkout

**Validation Points:**
- Place Order API
- Checkout validation

**Data Requirements:**
| Field | Type | Constraints | Source |
|-------|------|-------------|--------|
| availableQuantity | integer | >= requested quantity | Inventory entity |
| quantity | integer | requested order quantity | OrderLineItem entity |

**Error Handling:**
| Condition | Error Code | Message | HTTP Status |
|-----------|------------|---------|-------------|
| availableQuantity < requested quantity | STOCK_UNAVAILABLE | Insufficient stock for {productName}. Available: {availableQuantity} | 400 |

**Traceability:**
- BR: BR-INV-002

---

## TS-BR-INV-003 – Inventory restoration on cancellation

**Rule:** Cancelled orders must restore inventory automatically

**Implementation Approach:**
On order cancellation, iterate line items and restore stockLevel by each item's quantity

**Validation Points:**
- Cancel Order API

**Data Requirements:**
| Field | Type | Constraints | Source |
|-------|------|-------------|--------|
| stockLevel | integer | incremented by cancelled quantity | Inventory entity |
| quantity | integer | from each OrderLineItem | OrderLineItem entity |

**Traceability:**
- BR: BR-INV-003

---

## TS-BR-INV-004 – Inventory audit trail logging

**Rule:** All inventory changes must be logged with before/after values and reason

**Implementation Approach:**
Create InventoryAuditLog entry for every inventory modification with previousValue, newValue, reason, userId, timestamp

**Validation Points:**
- All inventory modification operations

**Data Requirements:**
| Field | Type | Constraints | Source |
|-------|------|-------------|--------|
| previousValue | integer | stock level before change | InventoryAuditLog entity |
| newValue | integer | stock level after change | InventoryAuditLog entity |
| reason | string | description of change reason | InventoryAuditLog entity |
| userId | uuid | user who made change | InventoryAuditLog entity |
| timestamp | timestamp | time of change | InventoryAuditLog entity |

**Traceability:**
- BR: BR-INV-004

---

## TS-BR-PAY-001 – Payment authorization before order creation

**Rule:** Payment must be authorized before order is created

**Implementation Approach:**
Call payment gateway authorization API; only create order if authorization succeeds

**Validation Points:**
- Place Order API
- Checkout flow

**Data Requirements:**
| Field | Type | Constraints | Source |
|-------|------|-------------|--------|
| authorizationId | string | returned from payment gateway on success | Payment entity |
| paymentStatus | enum | must be 'authorized' | Payment entity |

**Error Handling:**
| Condition | Error Code | Message | HTTP Status |
|-----------|------------|---------|-------------|
| payment authorization fails | PAYMENT_FAILED | Payment authorization failed. Please check your payment details and try again | 402 |

**Traceability:**
- BR: BR-PAY-001

---

## TS-BR-PAY-002 – Supported credit card type validation

**Rule:** Only Visa, Mastercard, and American Express credit cards are accepted

**Implementation Approach:**
Validate card type from BIN/IIN or explicit type field against allowed list

**Validation Points:**
- Payment processing API
- Add Payment Method API

**Data Requirements:**
| Field | Type | Constraints | Source |
|-------|------|-------------|--------|
| cardType | enum | visa|mastercard|amex | PaymentMethod entity |

**Error Handling:**
| Condition | Error Code | Message | HTTP Status |
|-----------|------------|---------|-------------|
| cardType NOT IN ('visa', 'mastercard', 'amex') | UNSUPPORTED_CARD_TYPE | Only Visa, Mastercard, and American Express cards are accepted | 400 |

**Traceability:**
- BR: BR-PAY-002

---

## TS-BR-PAY-003 – Automatic refund on order cancellation

**Rule:** Cancelled orders must be refunded to the original payment method

**Implementation Approach:**
On order cancellation, initiate refund via payment gateway to original payment method

**Validation Points:**
- Cancel Order API

**Data Requirements:**
| Field | Type | Constraints | Source |
|-------|------|-------------|--------|
| paymentMethodId | uuid | original payment method used | Order entity |
| refundAmount | decimal(10,2) | order total amount | Order entity |

**Error Handling:**
| Condition | Error Code | Message | HTTP Status |
|-----------|------------|---------|-------------|
| refund API call fails | REFUND_FAILED | Unable to process refund. Please contact customer support | 500 |

**Traceability:**
- BR: BR-PAY-003

---

## TS-BR-SHIP-001 – Domestic shipping validation

**Rule:** Shipping is only available within the domestic country

**Implementation Approach:**
Validate ShippingAddress.country against configured domestic country code

**Validation Points:**
- Checkout flow
- Add Shipping Address API

**Data Requirements:**
| Field | Type | Constraints | Source |
|-------|------|-------------|--------|
| country | string | must match configured domestic country code | ShippingAddress entity |

**Error Handling:**
| Condition | Error Code | Message | HTTP Status |
|-----------|------------|---------|-------------|
| country != domesticCountryCode | UNSUPPORTED_SHIPPING_REGION | We currently only ship within {domesticCountry} | 400 |

**Traceability:**
- BR: BR-SHIP-001

---

## TS-BR-SHIP-002 – Required shipping address fields validation

**Rule:** Shipping addresses must include all required fields

**Implementation Approach:**
Validate all required fields are present and non-empty

**Validation Points:**
- Checkout flow
- Add Shipping Address API
- Update Address API

**Data Requirements:**
| Field | Type | Constraints | Source |
|-------|------|-------------|--------|
| street | string | required, non-empty | ShippingAddress entity |
| city | string | required, non-empty | ShippingAddress entity |
| state | string | required, non-empty | ShippingAddress entity |
| postalCode | string | required, non-empty | ShippingAddress entity |
| country | string | required, non-empty | ShippingAddress entity |
| recipientName | string | required, non-empty | ShippingAddress entity |

**Error Handling:**
| Condition | Error Code | Message | HTTP Status |
|-----------|------------|---------|-------------|
| any required field is missing or empty | INVALID_ADDRESS | Please provide all required address fields | 400 |

**Traceability:**
- BR: BR-SHIP-002

---

## TS-BR-CAT-001 – Category nesting depth limit

**Rule:** Category hierarchy is limited to 3 levels deep

**Implementation Approach:**
Calculate depth when creating category with parent; reject if depth > 3

**Validation Points:**
- Create Category API
- Update Category Parent API

**Data Requirements:**
| Field | Type | Constraints | Source |
|-------|------|-------------|--------|
| parentId | uuid | nullable, references parent category | Category entity |
| depth | integer | <= 3 | Category entity |

**Error Handling:**
| Condition | Error Code | Message | HTTP Status |
|-----------|------------|---------|-------------|
| calculated depth > 3 | MAX_NESTING_EXCEEDED | Category hierarchy cannot exceed 3 levels | 400 |

**Traceability:**
- BR: BR-CAT-001

---

## TS-BR-EMAIL-001 – Non-blocking email delivery

**Rule:** Email delivery failures must not block order completion

**Implementation Approach:**
Send confirmation emails asynchronously via message queue with retry mechanism

**Validation Points:**
- Place Order API (post-commit)

**Data Requirements:**
| Field | Type | Constraints | Source |
|-------|------|-------------|--------|
| orderId | uuid | reference to created order | EmailQueue entity |
| emailType | enum | order_confirmation | EmailQueue entity |
| retryCount | integer | tracks retry attempts | EmailQueue entity |

**Traceability:**
- BR: BR-EMAIL-001

---

## TS-BR-EMAIL-002 – Email idempotency check

**Rule:** Duplicate email sends must be prevented

**Implementation Approach:**
Check for existing sent email record by orderId before sending; skip if already sent

**Validation Points:**
- Email sending service

**Data Requirements:**
| Field | Type | Constraints | Source |
|-------|------|-------------|--------|
| orderId | uuid | unique per confirmation email | SentEmailLog entity |
| sentAt | timestamp | records when email was sent | SentEmailLog entity |

**Traceability:**
- BR: BR-EMAIL-002

---

