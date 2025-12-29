# Business Rules

Generated: 2025-12-29T13:37:22+01:00

---

## BR-CUST-001 – Registration Required for Orders

**Rule:** Only registered and email-verified customers can place orders

**Invariant:** Customer MUST have registrationStatus='registered' AND emailVerified=true to place an order

**Enforcement:** Order service validates customer status before order creation

**Error Code:** `UNAUTHORIZED_ORDER`

**Traceability:**
- Source: Business rules
- Source: Customer entity
- Decision: AMB-ENT-030

---

## BR-CUST-002 – Unique Customer Email

**Rule:** Each customer account must have a unique email address

**Invariant:** Customer.email MUST be unique across all customer records

**Enforcement:** Database unique constraint, checked during registration

**Error Code:** `DUPLICATE_EMAIL`

**Traceability:**
- Source: Customer entity
- Decision: AMB-ENT-025

---

## BR-CUST-003 – Password Requirements

**Rule:** Customer passwords must be at least 8 characters with at least one number

**Invariant:** Customer.password MUST have length >= 8 AND contain at least one digit [0-9]

**Enforcement:** Password validation during registration and password change

**Error Code:** `INVALID_PASSWORD`

**Traceability:**
- Source: Customer entity
- Decision: AMB-ENT-026

---

## BR-PROD-001 – Product Price Minimum

**Rule:** Product prices must be positive with minimum $0.01

**Invariant:** Product.price MUST be >= 0.01 with exactly 2 decimal places

**Enforcement:** Validation during product creation and edit

**Error Code:** `INVALID_PRICE`

**Traceability:**
- Source: Product entity
- Decision: AMB-ENT-002

---

## BR-PROD-002 – Product Name Length

**Rule:** Product names must be between 2 and 200 characters

**Invariant:** Product.name MUST have 2 <= length <= 200

**Enforcement:** Validation during product creation and edit

**Error Code:** `INVALID_NAME`

**Traceability:**
- Source: Product entity
- Decision: AMB-ENT-003
- Decision: AMB-OP-033

---

## BR-PROD-003 – Product Soft Delete

**Rule:** Products with order history must be soft-deleted, not permanently deleted

**Invariant:** Product MUST NOT be permanently deleted if OrderLineItem references exist

**Enforcement:** Delete operation checks for order references, applies soft delete flag

**Error Code:** `N/A`

**Traceability:**
- Source: Product entity
- Source: Remove Product operation
- Decision: AMB-ENT-005
- Decision: AMB-OP-042

---

## BR-PROD-004 – Unique Variant SKU

**Rule:** Each product variant must have a unique SKU

**Invariant:** ProductVariant.sku MUST be unique across all variants

**Enforcement:** Database unique constraint, checked during variant creation

**Error Code:** `DUPLICATE_SKU`

**Traceability:**
- Source: ProductVariant entity
- Decision: AMB-ENT-011

---

## BR-CART-001 – In-Stock Products Only

**Rule:** Products must be in stock to be added to cart

**Invariant:** CartItem.productId MUST reference a product with Inventory.availableQuantity > 0 at time of add

**Enforcement:** Stock check during add-to-cart operation

**Error Code:** `OUT_OF_STOCK`

**Traceability:**
- Source: Business rules
- Source: Add to Cart operation
- Decision: AMB-OP-008

---

## BR-CART-002 – Quantity Limited by Stock

**Rule:** Cart item quantity cannot exceed available stock

**Invariant:** CartItem.quantity MUST be <= Inventory.availableQuantity for referenced product

**Enforcement:** Validation during add and quantity update, automatic capping

**Error Code:** `QUANTITY_LIMITED`

**Traceability:**
- Source: CartItem entity
- Source: Update Cart Quantity operation
- Decision: AMB-ENT-015
- Decision: AMB-OP-006
- Decision: AMB-OP-011

---

## BR-CART-003 – Positive Quantity Required

**Rule:** Cart item quantities must be at least 1

**Invariant:** CartItem.quantity MUST be >= 1 (setting to 0 triggers removal)

**Enforcement:** Quantity validation, zero triggers item removal

**Error Code:** `N/A`

**Traceability:**
- Source: CartItem entity
- Decision: AMB-OP-006
- Decision: AMB-OP-012

---

## BR-CART-004 – Cart Expiration

**Rule:** Inactive carts expire after 30 days (logged-in) or 7 days (guest)

**Invariant:** Cart MUST be cleared when lastActivityDate > expiration threshold

**Enforcement:** Background job or lazy evaluation on cart access

**Error Code:** `CART_EXPIRED`

**Traceability:**
- Source: Cart entity
- Decision: AMB-ENT-012

---

## BR-CART-005 – Variant Selection Required

**Rule:** For products with variants, a specific variant must be selected before adding to cart

**Invariant:** CartItem MUST reference specific ProductVariant if Product has variants

**Enforcement:** UI disables add-to-cart until variant selected, API validates

**Error Code:** `VARIANT_REQUIRED`

**Traceability:**
- Source: Add to Cart operation
- Decision: AMB-OP-010

---

## BR-ORDER-001 – Free Shipping Threshold

**Rule:** Orders over $50 qualify for free shipping

**Invariant:** Order.shippingCost MUST be $0.00 when Order.subtotal > $50.00

**Enforcement:** Shipping cost calculation during checkout

**Error Code:** `N/A`

**Traceability:**
- Source: Business rules
- Source: Place Order operation
- Decision: AMB-OP-020
- Decision: AMB-OP-021

---

## BR-ORDER-002 – Cancellation Before Shipping Only

**Rule:** Orders can only be cancelled before shipping

**Invariant:** Order.cancel() MUST only be allowed when status IN ('pending', 'confirmed')

**Enforcement:** Status check in cancel operation

**Error Code:** `ORDER_ALREADY_SHIPPED`

**Traceability:**
- Source: Business rules
- Source: Cancel Order operation
- Decision: AMB-OP-028
- Decision: AMB-ENT-018

---

## BR-ORDER-003 – Valid Order Status Transitions

**Rule:** Order status can only transition through valid states

**Invariant:** Order.status transitions MUST follow: pending→confirmed→shipped→delivered, OR pending/confirmed→cancelled

**Enforcement:** State machine validation in order service

**Error Code:** `INVALID_STATUS_TRANSITION`

**Traceability:**
- Source: Order entity
- Source: Update Order Status operation
- Decision: AMB-ENT-018
- Decision: AMB-OP-047

---

## BR-ORDER-004 – Order Immutability After Shipping

**Rule:** Orders cannot be modified after shipping

**Invariant:** Order MUST NOT allow modifications when status IN ('shipped', 'delivered', 'cancelled')

**Enforcement:** Status check before any order modification

**Error Code:** `ORDER_NOT_MODIFIABLE`

**Traceability:**
- Source: Order entity
- Decision: AMB-ENT-020

---

## BR-ORDER-005 – Minimum One Line Item

**Rule:** Orders must contain at least one line item

**Invariant:** Order.lineItems.length MUST be >= 1

**Enforcement:** Validation during order creation

**Error Code:** `EMPTY_ORDER`

**Traceability:**
- Source: Order entity

---

## BR-ORDER-006 – Price Snapshot at Order Time

**Rule:** Order line items must capture the price at order placement, not current price

**Invariant:** OrderLineItem.unitPrice MUST be set from Product.price at order creation time and MUST NOT change

**Enforcement:** Copy price during order creation, immutable afterwards

**Error Code:** `N/A`

**Traceability:**
- Source: Order entity
- Decision: AMB-ENT-019

---

## BR-INV-001 – Non-Negative Inventory

**Rule:** Inventory stock level cannot go below zero

**Invariant:** Inventory.stockLevel MUST be >= 0

**Enforcement:** Validation during all stock adjustments

**Error Code:** `INSUFFICIENT_STOCK`

**Traceability:**
- Source: Inventory entity
- Source: Manage Inventory operation
- Decision: AMB-OP-053
- Decision: AMB-ENT-043

---

## BR-INV-002 – No Backorders

**Rule:** Products cannot be sold when out of stock (no backorders)

**Invariant:** Order MUST NOT be created with items where Inventory.availableQuantity < requested quantity

**Enforcement:** Stock validation at checkout before order creation

**Error Code:** `STOCK_UNAVAILABLE`

**Traceability:**
- Source: Inventory entity
- Decision: AMB-ENT-043
- Decision: AMB-OP-018

---

## BR-INV-003 – Inventory Restoration on Cancellation

**Rule:** Cancelled orders must restore inventory automatically

**Invariant:** Order.cancel() MUST restore Inventory.stockLevel by sum of OrderLineItem quantities

**Enforcement:** Automatic stock restoration in cancel operation

**Error Code:** `N/A`

**Traceability:**
- Source: Order entity
- Source: Inventory entity
- Decision: AMB-ENT-044
- Decision: AMB-OP-031

---

## BR-INV-004 – Inventory Audit Trail

**Rule:** All inventory changes must be logged with before/after values and reason

**Invariant:** Every Inventory change MUST create InventoryAuditLog with previousValue, newValue, reason, userId, timestamp

**Enforcement:** Audit logging in inventory service for all operations

**Error Code:** `N/A`

**Traceability:**
- Source: Inventory entity
- Decision: AMB-ENT-042
- Decision: AMB-OP-055

---

## BR-PAY-001 – Payment Authorization Before Order

**Rule:** Payment must be authorized before order is created

**Invariant:** Order MUST NOT be created until payment authorization is successful

**Enforcement:** Payment gateway authorization in checkout flow

**Error Code:** `PAYMENT_FAILED`

**Traceability:**
- Source: Place Order operation
- Decision: AMB-OP-016
- Decision: AMB-OP-017

---

## BR-PAY-002 – Supported Credit Card Types

**Rule:** Only Visa, Mastercard, and American Express credit cards are accepted

**Invariant:** PaymentMethod.cardType MUST be IN ('visa', 'mastercard', 'amex') for credit card payments

**Enforcement:** Card type validation during payment processing

**Error Code:** `UNSUPPORTED_CARD_TYPE`

**Traceability:**
- Source: PaymentMethod entity
- Decision: AMB-ENT-034

---

## BR-PAY-003 – Automatic Refund on Cancellation

**Rule:** Cancelled orders must be refunded to the original payment method

**Invariant:** Order.cancel() MUST initiate refund to Order.paymentMethod

**Enforcement:** Refund API call in cancel operation

**Error Code:** `REFUND_FAILED`

**Traceability:**
- Source: Cancel Order operation
- Decision: AMB-OP-030

---

## BR-SHIP-001 – Domestic Shipping Only

**Rule:** Shipping is only available within the domestic country

**Invariant:** ShippingAddress.country MUST be the configured domestic country code

**Enforcement:** Address validation during checkout

**Error Code:** `UNSUPPORTED_SHIPPING_REGION`

**Traceability:**
- Source: ShippingAddress entity
- Decision: AMB-ENT-032

---

## BR-SHIP-002 – Required Address Fields

**Rule:** Shipping addresses must include all required fields

**Invariant:** ShippingAddress MUST have non-empty: street, city, state, postalCode, country, recipientName

**Enforcement:** Address validation during checkout and address save

**Error Code:** `INVALID_ADDRESS`

**Traceability:**
- Source: ShippingAddress entity
- Decision: AMB-ENT-031

---

## BR-CAT-001 – Category Nesting Limit

**Rule:** Category hierarchy is limited to 3 levels deep

**Invariant:** Category.depth MUST be <= 3

**Enforcement:** Depth validation during category creation with parent

**Error Code:** `MAX_NESTING_EXCEEDED`

**Traceability:**
- Source: Category entity
- Decision: AMB-ENT-037

---

## BR-EMAIL-001 – Email Delivery Non-Blocking

**Rule:** Email delivery failures must not block order completion

**Invariant:** Order creation MUST succeed even if email delivery fails

**Enforcement:** Asynchronous email via message queue with retry

**Error Code:** `N/A`

**Traceability:**
- Source: Send Confirmation Email operation
- Decision: AMB-OP-057
- Decision: AMB-OP-060

---

## BR-EMAIL-002 – Email Idempotency

**Rule:** Duplicate email sends must be prevented

**Invariant:** Confirmation email MUST be sent exactly once per order

**Enforcement:** Deduplication check using order ID before sending

**Error Code:** `N/A`

**Traceability:**
- Source: Send Confirmation Email operation
- Decision: AMB-OP-058

---

