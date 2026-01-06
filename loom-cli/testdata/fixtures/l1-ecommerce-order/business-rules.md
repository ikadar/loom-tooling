# Business Rules

Generated: 2025-12-31T13:11:10+01:00

---

## BR-AUTH-001 – Order placement requires registration {#br-auth-001}

**Rule:** Only registered customers with verified email addresses can place orders

**Invariant:** Order.customer MUST be a registered Customer with verified email

**Enforcement:** Checkout process validates customer registration and email verification status before allowing order submission

**Error Code:** `REGISTRATION_REQUIRED`

**Traceability:**
- Source: Business rule: Customers must be registered to place orders
- Decision: [AMB-ENT-030](decisions.md#amb-ent-030)

---

## BR-STOCK-001 – Stock availability for add to cart {#br-stock-001}

**Rule:** Products must have available inventory to be added to cart

**Invariant:** CartItem.quantity MUST NOT exceed Product/Variant available stock at checkout

**Enforcement:** Add to cart checks current stock; cart items can exist with warnings but checkout validates final availability

**Error Code:** `OUT_OF_STOCK`

**Traceability:**
- Source: Business rule: Products must be in stock to be added to cart
- Decision: [AMB-ENT-014](decisions.md#amb-ent-014)
- Decision: [AMB-OP-008](decisions.md#amb-op-008)

---

## BR-SHIP-001 – Free shipping threshold {#br-ship-001}

**Rule:** Orders with subtotal of $50 or more (before tax, after discounts) qualify for free shipping

**Invariant:** IF Order.subtotal >= 50.00 THEN Order.shippingCost MUST be 0

**Enforcement:** Shipping cost calculation during checkout; standard shipping ($5.99) applied for orders under $50

**Error Code:** `N/A`

**Traceability:**
- Source: Business rule: Orders over $50 qualify for free shipping
- Decision: [AMB-OP-020](decisions.md#amb-op-020)
- Decision: [AMB-OP-021](decisions.md#amb-op-021)

---

## BR-CANCEL-001 – Order cancellation window {#br-cancel-001}

**Rule:** Orders can only be cancelled while in pending or confirmed status (before shipping)

**Invariant:** Order cancellation MUST only be allowed when Order.status IN ('pending', 'confirmed')

**Enforcement:** Cancel order operation validates current status before allowing cancellation

**Error Code:** `CANCELLATION_NOT_ALLOWED`

**Traceability:**
- Source: Business rule: Orders can only be cancelled before shipping
- Decision: [AMB-ENT-018](decisions.md#amb-ent-018)
- Decision: [AMB-OP-028](decisions.md#amb-op-028)

---

## BR-PRICE-001 – Minimum product price {#br-price-001}

**Rule:** Product price must be at least $0.01

**Invariant:** Product.price MUST be >= 0.01 with exactly 2 decimal places

**Enforcement:** Product creation and edit validation

**Error Code:** `INVALID_PRICE`

**Traceability:**
- Source: Product entity
- Decision: [AMB-ENT-002](decisions.md#amb-ent-002)

---

## BR-PRICE-002 – Order price snapshot {#br-price-002}

**Rule:** Order items must capture prices at the time of order placement

**Invariant:** OrderItem.price MUST be set at order creation and MUST NOT change thereafter

**Enforcement:** Order placement copies current prices to order items; subsequent price changes do not affect existing orders

**Error Code:** `N/A`

**Traceability:**
- Source: Order entity
- Decision: [AMB-ENT-019](decisions.md#amb-ent-019)
- Decision: [AMB-OP-038](decisions.md#amb-op-038)

---

## BR-CART-001 – Cart price is always current {#br-cart-001}

**Rule:** Cart items always reflect current product prices, not prices at time of adding

**Invariant:** Cart.total MUST be calculated using current Product.price values

**Enforcement:** Cart total recalculates on every cart view; price change notifications shown to customer

**Error Code:** `N/A`

**Traceability:**
- Source: CartItem entity
- Decision: [AMB-ENT-016](decisions.md#amb-ent-016)

---

## BR-INV-001 – Inventory cannot go negative {#br-inv-001}

**Rule:** Inventory stock level must never be negative

**Invariant:** Inventory.availableStock MUST be >= 0

**Enforcement:** All inventory adjustments validate resulting quantity; order placement decrements inventory atomically with stock check

**Error Code:** `INSUFFICIENT_STOCK`

**Traceability:**
- Source: Inventory entity
- Decision: [AMB-ENT-043](decisions.md#amb-ent-043)
- Decision: [AMB-OP-053](decisions.md#amb-op-053)

---

## BR-INV-002 – No backorders allowed {#br-inv-002}

**Rule:** Orders cannot be placed for quantities exceeding available stock

**Invariant:** OrderItem.quantity MUST NOT exceed Inventory.availableStock at time of order placement

**Enforcement:** Checkout validates stock availability atomically; payment only authorized if stock is available

**Error Code:** `BACKORDER_NOT_SUPPORTED`

**Traceability:**
- Source: Inventory entity
- Decision: [AMB-ENT-043](decisions.md#amb-ent-043)

---

## BR-INV-003 – Inventory restoration on cancellation {#br-inv-003}

**Rule:** Cancelled orders must automatically restore inventory

**Invariant:** When Order.status changes to 'cancelled', all OrderItem quantities MUST be added back to respective Inventory.availableStock

**Enforcement:** Order cancellation operation atomically updates inventory

**Error Code:** `N/A`

**Traceability:**
- Source: Order entity
- Decision: [AMB-ENT-044](decisions.md#amb-ent-044)
- Decision: [AMB-OP-031](decisions.md#amb-op-031)

---

## BR-ORDER-001 – Valid order state transitions {#br-order-001}

**Rule:** Orders must follow defined state machine transitions

**Invariant:** Order.status transitions MUST follow: pending→confirmed→shipped→delivered, OR pending/confirmed→cancelled (terminal)

**Enforcement:** Order status update operation validates transition is allowed from current state

**Error Code:** `INVALID_STATUS_TRANSITION`

**Traceability:**
- Source: Order entity
- Decision: [AMB-ENT-018](decisions.md#amb-ent-018)
- Decision: [AMB-OP-047](decisions.md#amb-op-047)

---

## BR-ORDER-002 – Order immutability after placement {#br-order-002}

**Rule:** Orders cannot be modified after placement; customer must cancel and reorder

**Invariant:** Order items, quantities, and shipping address MUST NOT be modified after order creation

**Enforcement:** No edit operations exposed for placed orders; only status transitions allowed

**Error Code:** `ORDER_MODIFICATION_NOT_ALLOWED`

**Traceability:**
- Source: Order entity
- Decision: [AMB-ENT-020](decisions.md#amb-ent-020)

---

## BR-ORDER-003 – Single shipment per order {#br-order-003}

**Rule:** Each order ships as a single unit; partial shipments are not supported

**Invariant:** Order MUST have exactly one shipment record; all items ship together

**Enforcement:** No partial shipment functionality in system; shipped status applies to entire order

**Error Code:** `N/A`

**Traceability:**
- Source: Order entity
- Decision: [AMB-ENT-021](decisions.md#amb-ent-021)

---

## BR-CUST-001 – Unique customer email {#br-cust-001}

**Rule:** Email addresses must be unique across all customer accounts

**Invariant:** Customer.email MUST be unique in the system

**Enforcement:** Registration validates email uniqueness before account creation

**Error Code:** `EMAIL_ALREADY_REGISTERED`

**Traceability:**
- Source: Customer entity
- Decision: [AMB-ENT-025](decisions.md#amb-ent-025)

---

## BR-CUST-002 – Password requirements {#br-cust-002}

**Rule:** Customer passwords must meet minimum security requirements

**Invariant:** Customer.password MUST be at least 8 characters AND contain at least one number

**Enforcement:** Registration and password change validate password strength

**Error Code:** `WEAK_PASSWORD`

**Traceability:**
- Source: Customer entity
- Decision: [AMB-ENT-026](decisions.md#amb-ent-026)

---

## BR-CUST-003 – Email verification required for orders {#br-cust-003}

**Rule:** Customers must verify their email before placing orders

**Invariant:** Customer.emailVerified MUST be true before Order can be created

**Enforcement:** Checkout validates email verification status

**Error Code:** `EMAIL_NOT_VERIFIED`

**Traceability:**
- Source: Customer entity
- Decision: [AMB-ENT-030](decisions.md#amb-ent-030)

---

## BR-PROD-001 – Product soft delete for order history {#br-prod-001}

**Rule:** Products with order history must be soft-deleted to preserve historical data

**Invariant:** Product with associated OrderItems MUST NOT be permanently deleted; Product.deletedAt must be set instead

**Enforcement:** Delete operation checks for order history and performs soft delete

**Error Code:** `N/A`

**Traceability:**
- Source: Product entity
- Decision: [AMB-ENT-005](decisions.md#amb-ent-005)
- Decision: [AMB-OP-042](decisions.md#amb-op-042)

---

## BR-PROD-002 – Product name length constraints {#br-prod-002}

**Rule:** Product names must be between 2 and 200 characters

**Invariant:** Product.name.length MUST be >= 2 AND <= 200

**Enforcement:** Product creation and edit validation

**Error Code:** `INVALID_PRODUCT_NAME`

**Traceability:**
- Source: Product entity
- Decision: [AMB-ENT-003](decisions.md#amb-ent-003)
- Decision: [AMB-OP-033](decisions.md#amb-op-033)

---

## BR-PROD-003 – Product description length constraint {#br-prod-003}

**Rule:** Product descriptions must not exceed 5000 characters

**Invariant:** Product.description.length MUST be <= 5000

**Enforcement:** Product creation and edit validation

**Error Code:** `DESCRIPTION_TOO_LONG`

**Traceability:**
- Source: Product entity
- Decision: [AMB-ENT-003](decisions.md#amb-ent-003)

---

## BR-PROD-004 – Product image limit {#br-prod-004}

**Rule:** Products can have up to 10 images with exactly one primary image

**Invariant:** Product.images.count MUST be <= 10 AND exactly one image MUST be marked as primary

**Enforcement:** Image upload validation; primary image selection enforced

**Error Code:** `IMAGE_LIMIT_EXCEEDED`

**Traceability:**
- Source: Product entity
- Decision: [AMB-ENT-004](decisions.md#amb-ent-004)

---

## BR-VAR-001 – Variant SKU uniqueness {#br-var-001}

**Rule:** Each product variant must have a unique SKU

**Invariant:** ProductVariant.sku MUST be unique across all variants in the system

**Enforcement:** Variant creation validates SKU uniqueness

**Error Code:** `DUPLICATE_SKU`

**Traceability:**
- Source: ProductVariant entity
- Decision: [AMB-ENT-011](decisions.md#amb-ent-011)

---

## BR-VAR-002 – Variant selection required for purchase {#br-var-002}

**Rule:** Products with variants require variant selection before add to cart

**Invariant:** IF Product has variants THEN CartItem MUST reference a specific ProductVariant

**Enforcement:** Add to cart disabled until variant selected; cart item stores variant reference

**Error Code:** `VARIANT_REQUIRED`

**Traceability:**
- Source: ProductVariant entity
- Source: Add to Cart operation
- Decision: [AMB-OP-010](decisions.md#amb-op-010)

---

## BR-CART-002 – Cart quantity limited by stock {#br-cart-002}

**Rule:** Cart item quantity cannot exceed available inventory

**Invariant:** CartItem.quantity MUST be <= Inventory.availableStock for the referenced product/variant

**Enforcement:** Add to cart and quantity update operations cap quantity at available stock

**Error Code:** `QUANTITY_EXCEEDS_STOCK`

**Traceability:**
- Source: CartItem entity
- Decision: [AMB-ENT-015](decisions.md#amb-ent-015)
- Decision: [AMB-OP-006](decisions.md#amb-op-006)
- Decision: [AMB-OP-011](decisions.md#amb-op-011)

---

## BR-CART-003 – Cart expiration policy {#br-cart-003}

**Rule:** Carts expire after period of inactivity based on user type

**Invariant:** Logged-in user carts expire after 30 days of inactivity; guest carts expire after 7 days

**Enforcement:** Background job clears expired carts; cart activity timestamp updated on modifications

**Error Code:** `CART_EXPIRED`

**Traceability:**
- Source: Cart entity
- Decision: [AMB-ENT-012](decisions.md#amb-ent-012)

---

## BR-PAY-001 – Payment before order creation {#br-pay-001}

**Rule:** Payment must be fully authorized before order is created

**Invariant:** Order creation MUST only occur after successful payment authorization

**Enforcement:** Checkout flow authorizes payment first; order only created on successful authorization

**Error Code:** `PAYMENT_REQUIRED`

**Traceability:**
- Source: Place Order operation
- Decision: [AMB-OP-016](decisions.md#amb-op-016)

---

## BR-PAY-002 – Automatic refund on cancellation {#br-pay-002}

**Rule:** Cancelled orders must receive automatic refund to original payment method

**Invariant:** When Order.status changes to 'cancelled', refund MUST be initiated to Order.paymentMethod

**Enforcement:** Cancellation operation triggers refund via payment gateway

**Error Code:** `N/A`

**Traceability:**
- Source: Cancel Order operation
- Decision: [AMB-OP-030](decisions.md#amb-op-030)

---

## BR-SHIP-002 – Domestic shipping only {#br-ship-002}

**Rule:** Shipping is only available within a single country

**Invariant:** ShippingAddress.country MUST match configured domestic country

**Enforcement:** Address validation rejects international addresses

**Error Code:** `INTERNATIONAL_SHIPPING_NOT_AVAILABLE`

**Traceability:**
- Source: ShippingAddress entity
- Decision: [AMB-ENT-032](decisions.md#amb-ent-032)

---

## BR-ADDR-001 – Required shipping address fields {#br-addr-001}

**Rule:** Shipping addresses must include all required fields

**Invariant:** ShippingAddress MUST have: street, city, state/province, postalCode, country, recipientName

**Enforcement:** Address form and API validation

**Error Code:** `INCOMPLETE_ADDRESS`

**Traceability:**
- Source: ShippingAddress entity
- Decision: [AMB-ENT-031](decisions.md#amb-ent-031)

---

## BR-CAT-001 – Category hierarchy depth limit {#br-cat-001}

**Rule:** Categories can be nested up to 3 levels deep

**Invariant:** Category.depth MUST be <= 3

**Enforcement:** Category creation validates parent hierarchy depth

**Error Code:** `CATEGORY_DEPTH_EXCEEDED`

**Traceability:**
- Source: Category entity
- Decision: [AMB-ENT-037](decisions.md#amb-ent-037)

---

## BR-AUDIT-001 – Order event logging {#br-audit-001}

**Rule:** All order events must be logged for audit purposes

**Invariant:** Order status changes, creation, and cancellation MUST be logged with timestamp, user, and before/after values

**Enforcement:** Order operations write to audit log

**Error Code:** `N/A`

**Traceability:**
- Source: Order entity
- Source: Update Order Status operation
- Decision: [AMB-ENT-022](decisions.md#amb-ent-022)
- Decision: [AMB-OP-024](decisions.md#amb-op-024)
- Decision: [AMB-OP-050](decisions.md#amb-op-050)

---

## BR-AUDIT-002 – Inventory change logging {#br-audit-002}

**Rule:** All inventory changes must be logged with full audit trail

**Invariant:** Inventory adjustments MUST log before/after values, reason, and operator

**Enforcement:** Inventory operations write to audit log

**Error Code:** `N/A`

**Traceability:**
- Source: Inventory entity
- Decision: [AMB-ENT-042](decisions.md#amb-ent-042)
- Decision: [AMB-OP-055](decisions.md#amb-op-055)

---

## BR-AUDIT-003 – Price change tracking {#br-audit-003}

**Rule:** Product price changes must be tracked in audit trail

**Invariant:** Product.price changes MUST log previous price, new price, modifier, and timestamp

**Enforcement:** Product edit operation logs price changes

**Error Code:** `N/A`

**Traceability:**
- Source: Product entity
- Decision: [AMB-ENT-007](decisions.md#amb-ent-007)

---

