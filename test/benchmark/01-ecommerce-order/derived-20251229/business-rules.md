# Business Rules

Generated: 2025-12-29T11:31:02+01:00

---

## BR-AUTH-001 – Registration required for orders

**Rule:** Only registered customers with verified email addresses can place orders

**Invariant:** Order.customer MUST have status 'registered' AND email_verified = true

**Enforcement:** Validated at checkout initiation and order submission

**Error Code:** `REGISTRATION_REQUIRED`

**Traceability:**
- Source: Domain Model: Business Rules
- Decision: AMB-ENT-030

---

## BR-STOCK-001 – Stock requirement for cart addition

**Rule:** Products must have available stock to be added to cart

**Invariant:** Product.inventory.available_quantity MUST be > 0 for add to cart

**Enforcement:** Validated at add to cart action; stock validated again at checkout

**Error Code:** `OUT_OF_STOCK`

**Traceability:**
- Source: Domain Model: Business Rules
- Decision: AMB-OP-008
- Decision: AMB-ENT-043

---

## BR-STOCK-002 – No backorders allowed

**Rule:** Orders cannot be placed for quantities exceeding available stock

**Invariant:** CartItem.quantity MUST be <= ProductVariant.inventory.available_quantity at checkout

**Enforcement:** Strict validation at checkout; cart can temporarily hold more but blocked at order

**Error Code:** `INSUFFICIENT_STOCK`

**Traceability:**
- Source: Domain Model: Inventory
- Decision: AMB-ENT-043
- Decision: AMB-OP-018

---

## BR-STOCK-003 – Inventory cannot go negative

**Rule:** Stock levels must never be negative

**Invariant:** Inventory.quantity MUST be >= 0

**Enforcement:** Validated on all inventory operations (adjustments, order placement, cancellation)

**Error Code:** `INVALID_STOCK_LEVEL`

**Traceability:**
- Source: Domain Model: Inventory
- Decision: AMB-OP-053

---

## BR-SHIP-001 – Free shipping threshold

**Rule:** Orders with subtotal of $50 or more (before tax, after discounts) qualify for free shipping

**Invariant:** IF Order.subtotal >= 50.00 THEN Order.shipping_cost = 0

**Enforcement:** Calculated during checkout total computation

**Error Code:** `N/A`

**Traceability:**
- Source: Domain Model: Business Rules
- Decision: AMB-OP-020
- Decision: AMB-OP-021

---

## BR-CANCEL-001 – Order cancellation window

**Rule:** Orders can only be cancelled before shipping

**Invariant:** Order.status MUST be IN ('pending', 'confirmed') to allow cancellation

**Enforcement:** Validated at cancellation request; cancel button hidden for shipped/delivered orders

**Error Code:** `CANCELLATION_NOT_ALLOWED`

**Traceability:**
- Source: Domain Model: Business Rules
- Decision: AMB-OP-028
- Decision: AMB-ENT-018

---

## BR-ORDER-001 – Valid order status transitions

**Rule:** Orders must follow defined state machine: pending→confirmed→shipped→delivered, with cancelled as terminal state from pending or confirmed only

**Invariant:** Order status transitions MUST follow: {pending→confirmed, pending→cancelled, confirmed→shipped, confirmed→cancelled, shipped→delivered}

**Enforcement:** Validated at every status update; invalid transitions rejected

**Error Code:** `INVALID_STATUS_TRANSITION`

**Traceability:**
- Source: Domain Model: Order
- Decision: AMB-ENT-018
- Decision: AMB-OP-047

---

## BR-ORDER-002 – Order immutability after placement

**Rule:** Orders cannot be modified after placement; customer must cancel and reorder

**Invariant:** Order items, quantities, and addresses MUST NOT change after status != 'cart'

**Enforcement:** No edit endpoints exposed; only status updates and cancellation allowed

**Error Code:** `ORDER_MODIFICATION_NOT_ALLOWED`

**Traceability:**
- Source: Domain Model: Order
- Decision: AMB-ENT-020

---

## BR-ORDER-003 – Single shipment per order

**Rule:** Each order is shipped as a single unit; partial shipments not supported

**Invariant:** Order MUST have exactly one shipment record

**Enforcement:** System design constraint; no split shipment functionality

**Error Code:** `N/A`

**Traceability:**
- Source: Domain Model: Order
- Decision: AMB-ENT-021

---

## BR-ORDER-004 – Price snapshot at order time

**Rule:** Order items must capture the product price at the time of order placement

**Invariant:** OrderItem.unit_price MUST be set to Product.price at order creation time and MUST NOT change

**Enforcement:** Price captured during order creation; historical prices preserved

**Error Code:** `N/A`

**Traceability:**
- Source: Domain Model: Order
- Decision: AMB-ENT-019

---

## BR-PRICE-001 – Minimum product price

**Rule:** Product prices must be at least $0.01

**Invariant:** Product.price MUST be >= 0.01

**Enforcement:** Validated at product creation and update

**Error Code:** `INVALID_PRICE`

**Traceability:**
- Source: Domain Model: Product
- Decision: AMB-ENT-002

---

## BR-PRICE-002 – Cart uses current prices

**Rule:** Cart items always reflect current product prices, not historical prices

**Invariant:** CartItem displayed price MUST equal current Product.price

**Enforcement:** Price recalculated on every cart view; notifications shown for changes

**Error Code:** `N/A`

**Traceability:**
- Source: Domain Model: CartItem
- Decision: AMB-ENT-016

---

## BR-CART-001 – Cart quantity limits

**Rule:** Cart item quantity must be between 1 and available stock

**Invariant:** CartItem.quantity MUST be >= 1 AND <= available_stock

**Enforcement:** Validated at add and update operations; quantity auto-limited to stock with notification

**Error Code:** `QUANTITY_EXCEEDS_STOCK`

**Traceability:**
- Source: Domain Model: CartItem
- Decision: AMB-ENT-015
- Decision: AMB-OP-006
- Decision: AMB-OP-011

---

## BR-CART-002 – Cart expiration

**Rule:** Inactive carts expire after 30 days for logged-in users and 7 days for guests

**Invariant:** Cart MUST be purged if last_activity > (30 days for registered, 7 days for guest)

**Enforcement:** Background job cleans expired carts

**Error Code:** `N/A`

**Traceability:**
- Source: Domain Model: Cart
- Decision: AMB-ENT-012

---

## BR-CUST-001 – Unique customer email

**Rule:** Customer email addresses must be unique in the system

**Invariant:** Customer.email MUST be unique across all customers

**Enforcement:** Validated at registration; database unique constraint

**Error Code:** `DUPLICATE_EMAIL`

**Traceability:**
- Source: Domain Model: Customer
- Decision: AMB-ENT-025

---

## BR-CUST-002 – Password requirements

**Rule:** Customer passwords must be at least 8 characters with at least one number

**Invariant:** Customer.password MUST have length >= 8 AND contain at least one digit

**Enforcement:** Validated at registration and password change

**Error Code:** `PASSWORD_REQUIREMENTS_NOT_MET`

**Traceability:**
- Source: Domain Model: Customer
- Decision: AMB-ENT-026

---

## BR-CUST-003 – Email verification required

**Rule:** Customers must verify their email address before placing orders

**Invariant:** Customer.email_verified MUST be true before order placement

**Enforcement:** Validated at checkout; verification email sent on registration

**Error Code:** `EMAIL_NOT_VERIFIED`

**Traceability:**
- Source: Domain Model: Customer
- Decision: AMB-ENT-030

---

## BR-SHIP-002 – Domestic shipping only

**Rule:** Shipping is only supported within the domestic country

**Invariant:** ShippingAddress.country MUST be the supported domestic country

**Enforcement:** Validated at address creation and checkout

**Error Code:** `UNSUPPORTED_COUNTRY`

**Traceability:**
- Source: Domain Model: ShippingAddress
- Decision: AMB-ENT-032

---

## BR-PAY-001 – Supported payment types

**Rule:** Only Visa, Mastercard, American Express credit cards and PayPal are accepted

**Invariant:** PaymentMethod.type MUST be IN ('visa', 'mastercard', 'amex', 'paypal')

**Enforcement:** Validated at payment method addition and checkout

**Error Code:** `UNSUPPORTED_PAYMENT_TYPE`

**Traceability:**
- Source: Domain Model: PaymentMethod
- Decision: AMB-ENT-034
- Decision: AMB-ENT-035

---

## BR-PAY-002 – Payment authorization before order

**Rule:** Full payment must be authorized before an order is created

**Invariant:** Order creation MUST be preceded by successful payment authorization

**Enforcement:** Transactional order creation; rollback if payment fails

**Error Code:** `PAYMENT_FAILED`

**Traceability:**
- Source: Domain Model: Place Order
- Decision: AMB-OP-016

---

## BR-PROD-001 – Product soft deletion

**Rule:** Products with order history must be soft-deleted to preserve historical references

**Invariant:** Product with associated OrderItems MUST NOT be hard-deleted

**Enforcement:** Delete operation performs soft delete; hard delete blocked if orders exist

**Error Code:** `N/A`

**Traceability:**
- Source: Domain Model: Product
- Decision: AMB-ENT-005
- Decision: AMB-OP-042

---

## BR-PROD-002 – Product name constraints

**Rule:** Product names must be between 2 and 200 characters

**Invariant:** Product.name.length MUST be >= 2 AND <= 200

**Enforcement:** Validated at product creation and update

**Error Code:** `INVALID_NAME`

**Traceability:**
- Source: Domain Model: Product
- Decision: AMB-ENT-003
- Decision: AMB-OP-033

---

## BR-PROD-003 – Product requires category

**Rule:** Every product must have a primary category assigned

**Invariant:** Product.primary_category MUST NOT be null

**Enforcement:** Validated at product creation and update

**Error Code:** `CATEGORY_REQUIRED`

**Traceability:**
- Source: Domain Model: Product
- Decision: AMB-ENT-006
- Decision: AMB-OP-033

---

## BR-VAR-001 – Unique variant SKU

**Rule:** Each product variant must have a unique SKU

**Invariant:** ProductVariant.sku MUST be unique across all variants

**Enforcement:** Validated at variant creation; database unique constraint

**Error Code:** `DUPLICATE_SKU`

**Traceability:**
- Source: Domain Model: ProductVariant
- Decision: AMB-ENT-011

---

## BR-VAR-002 – Variant selection required

**Rule:** For products with variants, a specific variant must be selected before adding to cart

**Invariant:** CartItem referencing Product with variants MUST specify ProductVariant.id

**Enforcement:** Add to cart disabled until variant selected; validated at API level

**Error Code:** `VARIANT_REQUIRED`

**Traceability:**
- Source: Domain Model: Add to Cart
- Decision: AMB-OP-010

---

## BR-CAT-001 – Category nesting depth

**Rule:** Categories can only be nested up to 3 levels deep

**Invariant:** Category.depth MUST be <= 3

**Enforcement:** Validated at category creation and parent assignment

**Error Code:** `CATEGORY_DEPTH_EXCEEDED`

**Traceability:**
- Source: Domain Model: Category
- Decision: AMB-ENT-037

---

## BR-INV-001 – Inventory restoration on cancellation

**Rule:** When an order is cancelled, all reserved inventory must be automatically restored

**Invariant:** Order cancellation MUST trigger inventory.available += order_item.quantity for all items

**Enforcement:** Transactional cancellation process; inventory update in same transaction

**Error Code:** `N/A`

**Traceability:**
- Source: Domain Model: Order
- Decision: AMB-ENT-044
- Decision: AMB-OP-031

---

## BR-INV-002 – Inventory decrement on payment

**Rule:** Inventory is decremented immediately upon successful payment

**Invariant:** Successful payment MUST trigger inventory.available -= order_item.quantity for all items

**Enforcement:** Transactional order creation; inventory and payment in same transaction

**Error Code:** `N/A`

**Traceability:**
- Source: Domain Model: Place Order
- Decision: AMB-OP-019

---

## BR-AUDIT-001 – Order status audit trail

**Rule:** All order status changes must be logged with timestamp and actor

**Invariant:** Every Order.status change MUST create OrderStatusHistory record

**Enforcement:** Automatic logging in status update service

**Error Code:** `N/A`

**Traceability:**
- Source: Domain Model: Order
- Decision: AMB-ENT-022
- Decision: AMB-OP-050

---

## BR-AUDIT-002 – Inventory audit trail

**Rule:** All inventory changes must be logged with before/after values, reason, and actor

**Invariant:** Every Inventory.quantity change MUST create InventoryHistory record

**Enforcement:** Automatic logging in inventory update service

**Error Code:** `N/A`

**Traceability:**
- Source: Domain Model: Inventory
- Decision: AMB-ENT-042
- Decision: AMB-OP-055

---

## BR-AUDIT-003 – Price change audit

**Rule:** Product price changes must be tracked with who modified and when

**Invariant:** Product.price change MUST create PriceHistory record with modifier_id and timestamp

**Enforcement:** Automatic logging in product update service

**Error Code:** `N/A`

**Traceability:**
- Source: Domain Model: Product
- Decision: AMB-ENT-007

---

