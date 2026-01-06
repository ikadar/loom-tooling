# Acceptance Criteria

Generated: 2025-12-31T13:11:10+01:00

---

## AC-PROD-001 – Browse products with filtering {#ac-prod-001}

**Given** a customer is on the product listing page
**When** they apply filters for category, price range, or availability
**Then** only products matching all selected filters are displayed, sorted by the selected option (name, price asc/desc, or newest), with 20 products per page

**Error Cases:**
- No products match filters → Show empty state with 'No products found' message and clear filters option
- Invalid price range (min > max) → Show validation error

**Traceability:**
- Source: Browse Products operation
- Source: Product entity
- Decision: [AMB-OP-001](decisions.md#amb-op-001)
- Decision: [AMB-OP-002](decisions.md#amb-op-002)
- Decision: [AMB-OP-003](decisions.md#amb-op-003)
- Decision: [AMB-OP-004](decisions.md#amb-op-004)

---

## AC-PROD-002 – View out-of-stock products {#ac-prod-002}

**Given** a product has zero available inventory
**When** a customer views the product listing or product detail page
**Then** the product is displayed with an 'Out of Stock' indicator and the add-to-cart button is disabled

**Traceability:**
- Source: Browse Products operation
- Decision: [AMB-OP-004](decisions.md#amb-op-004)
- Decision: [AMB-ENT-009](decisions.md#amb-ent-009)

---

## AC-CART-001 – Add product to cart {#ac-cart-001}

**Given** a customer is viewing an in-stock product with selected variant (if applicable)
**When** they click 'Add to Cart' with quantity 1-N (where N ≤ available stock)
**Then** the item is added to their cart, a toast notification appears with 'View Cart' option, and the cart icon updates to show the new item count

**Error Cases:**
- Product out of stock → Add to cart button disabled, show 'Out of Stock'
- Quantity exceeds available stock → Limit quantity to available stock and show notification
- Variant not selected → Add to cart button disabled until variant selected

**Traceability:**
- Source: Add to Cart operation
- Source: Cart entity
- Decision: [AMB-OP-005](decisions.md#amb-op-005)
- Decision: [AMB-OP-006](decisions.md#amb-op-006)
- Decision: [AMB-OP-007](decisions.md#amb-op-007)
- Decision: [AMB-OP-008](decisions.md#amb-op-008)
- Decision: [AMB-OP-009](decisions.md#amb-op-009)
- Decision: [AMB-OP-010](decisions.md#amb-op-010)

---

## AC-CART-002 – Add existing item to cart {#ac-cart-002}

**Given** a customer has Product A (quantity 2) in their cart
**When** they add Product A again with quantity 1
**Then** the existing cart item quantity is incremented to 3 (not a new line item created)

**Error Cases:**
- Combined quantity exceeds stock → Limit to available stock and show notification

**Traceability:**
- Source: Add to Cart operation
- Decision: [AMB-OP-005](decisions.md#amb-op-005)
- Decision: [AMB-ENT-015](decisions.md#amb-ent-015)

---

## AC-CART-003 – Guest user cart {#ac-cart-003}

**Given** a user is browsing without being logged in
**When** they add items to cart
**Then** a session-based cart is created that persists for 7 days of inactivity

**Error Cases:**
- Session expires after 7 days → Cart is cleared, user sees empty cart

**Traceability:**
- Source: Cart entity
- Source: Add to Cart operation
- Decision: [AMB-ENT-012](decisions.md#amb-ent-012)
- Decision: [AMB-ENT-013](decisions.md#amb-ent-013)
- Decision: [AMB-OP-007](decisions.md#amb-op-007)

---

## AC-CART-004 – Cart merge on login {#ac-cart-004}

**Given** a guest user has items in their session cart AND has a registered account with existing cart items
**When** they log in
**Then** the guest cart items are merged with their account cart, combining quantities for duplicate products

**Error Cases:**
- Merged quantity exceeds stock → Limit to available stock and show notification

**Traceability:**
- Source: Cart entity
- Decision: [AMB-ENT-013](decisions.md#amb-ent-013)

---

## AC-CART-005 – Update cart quantity {#ac-cart-005}

**Given** a customer has an item in their cart
**When** they change the quantity to a new value
**Then** the cart item quantity is updated and the cart total recalculates

**Error Cases:**
- Quantity set to 0 → Item is removed from cart
- Quantity exceeds available stock → Quantity limited to available stock with notification
- Quantity set to same value → No change (idempotent)

**Traceability:**
- Source: Update Cart Quantity operation
- Decision: [AMB-OP-011](decisions.md#amb-op-011)
- Decision: [AMB-OP-012](decisions.md#amb-op-012)
- Decision: [AMB-OP-013](decisions.md#amb-op-013)

---

## AC-CART-006 – Remove item from cart {#ac-cart-006}

**Given** a customer has items in their cart
**When** they click remove on an item
**Then** the item is immediately removed (no confirmation) and an 'Undo' option is shown briefly

**Error Cases:**
- Last item removed → Show empty cart state with 'Continue Shopping' link

**Traceability:**
- Source: Remove from Cart operation
- Decision: [AMB-OP-014](decisions.md#amb-op-014)
- Decision: [AMB-OP-015](decisions.md#amb-op-015)

---

## AC-CART-007 – Cart item goes out of stock {#ac-cart-007}

**Given** a customer has an item in their cart
**When** that product's inventory drops to zero (purchased by others)
**Then** the item remains in cart with a warning indicator, but checkout is blocked for that item

**Traceability:**
- Source: Cart entity
- Decision: [AMB-ENT-014](decisions.md#amb-ent-014)

---

## AC-CART-008 – Cart reflects current prices {#ac-cart-008}

**Given** a customer has items in their cart
**When** a product's price changes
**Then** the cart displays the current price and notifies the customer of the price change

**Traceability:**
- Source: CartItem entity
- Decision: [AMB-ENT-016](decisions.md#amb-ent-016)

---

## AC-ORDER-001 – Place order as registered customer {#ac-order-001}

**Given** a registered and email-verified customer has items in cart, a valid shipping address, and a valid payment method
**When** they submit the order
**Then** payment is authorized, inventory is decremented, order is created with status 'pending', cart is cleared, order confirmation page is shown, and confirmation email is queued

**Error Cases:**
- Customer not registered → Redirect to registration/login
- Email not verified → Show verification required message
- Payment authorization fails → Show error, keep cart intact, allow retry
- Item out of stock during checkout → Show error for affected items, allow cart modification
- Shipping address invalid → Show validation errors

**Traceability:**
- Source: Place Order operation
- Source: Order entity
- Decision: [AMB-OP-016](decisions.md#amb-op-016)
- Decision: [AMB-OP-017](decisions.md#amb-op-017)
- Decision: [AMB-OP-018](decisions.md#amb-op-018)
- Decision: [AMB-OP-019](decisions.md#amb-op-019)
- Decision: [AMB-OP-024](decisions.md#amb-op-024)
- Decision: [AMB-ENT-030](decisions.md#amb-ent-030)

---

## AC-ORDER-002 – Free shipping threshold {#ac-order-002}

**Given** a customer is checking out
**When** their subtotal (after discounts, before tax) is $50 or more
**Then** free shipping is automatically applied

**Error Cases:**
- Subtotal below $50 → Standard shipping fee of $5.99 is applied

**Traceability:**
- Source: Place Order operation
- Decision: [AMB-OP-020](decisions.md#amb-op-020)
- Decision: [AMB-OP-021](decisions.md#amb-op-021)

---

## AC-ORDER-003 – Tax calculation {#ac-order-003}

**Given** a customer is checking out with a shipping address
**When** the order total is calculated
**Then** applicable US state tax is calculated based on the shipping address state

**Traceability:**
- Source: Place Order operation
- Decision: [AMB-OP-022](decisions.md#amb-op-022)

---

## AC-ORDER-004 – Order confirmation email {#ac-order-004}

**Given** an order has been successfully placed
**When** the order is created
**Then** an email is queued asynchronously containing order number, items, quantities, prices, total, shipping address, and estimated delivery

**Error Cases:**
- Email service fails → Queue for retry, log failure, order creation not blocked

**Traceability:**
- Source: Send Confirmation Email operation
- Decision: [AMB-OP-056](decisions.md#amb-op-056)
- Decision: [AMB-OP-057](decisions.md#amb-op-057)
- Decision: [AMB-OP-058](decisions.md#amb-op-058)
- Decision: [AMB-OP-060](decisions.md#amb-op-060)

---

## AC-ORDER-005 – Track order status {#ac-order-005}

**Given** a registered customer has placed orders
**When** they view their order history
**Then** all their orders are displayed with pagination, showing order number, date, items, quantities, prices, status, shipping address, and tracking number (if available)

**Error Cases:**
- No orders exist → Show empty state with 'Start Shopping' link

**Traceability:**
- Source: Track Order Status operation
- Decision: [AMB-OP-025](decisions.md#amb-op-025)
- Decision: [AMB-OP-026](decisions.md#amb-op-026)
- Decision: [AMB-OP-027](decisions.md#amb-op-027)

---

## AC-ORDER-006 – Cancel order {#ac-order-006}

**Given** a customer has an order in 'pending' or 'confirmed' status
**When** they request cancellation
**Then** the order status changes to 'cancelled', inventory is automatically restored, refund is initiated to original payment method, and cancellation confirmation email is sent

**Error Cases:**
- Order already shipped → Cancellation not allowed, show error message
- Order already delivered → Cancellation not allowed, show error message
- Order already cancelled → Show already cancelled message

**Traceability:**
- Source: Cancel Order operation
- Decision: [AMB-OP-028](decisions.md#amb-op-028)
- Decision: [AMB-OP-029](decisions.md#amb-op-029)
- Decision: [AMB-OP-030](decisions.md#amb-op-030)
- Decision: [AMB-OP-031](decisions.md#amb-op-031)
- Decision: [AMB-OP-032](decisions.md#amb-op-032)
- Decision: [AMB-ENT-044](decisions.md#amb-ent-044)

---

## AC-ORDER-007 – Order prices snapshot {#ac-order-007}

**Given** a customer places an order
**When** the order is created
**Then** all item prices are captured at the time of order and remain unchanged even if product prices change later

**Traceability:**
- Source: Order entity
- Decision: [AMB-ENT-019](decisions.md#amb-ent-019)

---

## AC-CUST-001 – Customer registration {#ac-cust-001}

**Given** a new user wants to create an account
**When** they submit registration with email, password (min 8 chars with at least one number), first name, and last name
**Then** the account is created, verification email is sent, and user is informed to verify email before placing orders

**Error Cases:**
- Email already registered → Show 'Email already in use' error
- Password too weak → Show password requirements
- Required fields missing → Show validation errors

**Traceability:**
- Source: Customer entity
- Decision: [AMB-ENT-024](decisions.md#amb-ent-024)
- Decision: [AMB-ENT-025](decisions.md#amb-ent-025)
- Decision: [AMB-ENT-026](decisions.md#amb-ent-026)
- Decision: [AMB-ENT-030](decisions.md#amb-ent-030)

---

## AC-CUST-002 – Customer saved addresses {#ac-cust-002}

**Given** a registered customer is managing their profile
**When** they add a shipping address with street, city, state/province, postal code, country, and recipient name
**Then** the address is saved to their account, with option to set as default

**Error Cases:**
- Required fields missing → Show validation errors
- Invalid postal code format → Show format error

**Traceability:**
- Source: Customer entity
- Source: ShippingAddress entity
- Decision: [AMB-ENT-027](decisions.md#amb-ent-027)
- Decision: [AMB-ENT-031](decisions.md#amb-ent-031)
- Decision: [AMB-ENT-032](decisions.md#amb-ent-032)
- Decision: [AMB-ENT-033](decisions.md#amb-ent-033)

---

## AC-CUST-003 – Customer saved payment methods {#ac-cust-003}

**Given** a registered customer is managing their profile
**When** they add a payment method (Visa, Mastercard, or American Express via Stripe, or PayPal)
**Then** the tokenized payment method is saved via the payment gateway, with option to set as default

**Error Cases:**
- Invalid card → Show card validation error from payment gateway
- Unsupported card type → Show supported card types

**Traceability:**
- Source: Customer entity
- Source: PaymentMethod entity
- Decision: [AMB-ENT-028](decisions.md#amb-ent-028)
- Decision: [AMB-ENT-034](decisions.md#amb-ent-034)
- Decision: [AMB-ENT-035](decisions.md#amb-ent-035)
- Decision: [AMB-ENT-036](decisions.md#amb-ent-036)

---

## AC-CUST-004 – Customer data deletion (GDPR) {#ac-cust-004}

**Given** a registered customer requests account deletion
**When** they confirm the deletion request
**Then** the system supports either soft delete (account deactivated) or full data erasure for GDPR compliance based on customer preference

**Error Cases:**
- Customer has pending orders → Block deletion until orders completed

**Traceability:**
- Source: Customer entity
- Decision: [AMB-ENT-029](decisions.md#amb-ent-029)

---

## AC-ADMIN-001 – Add new product {#ac-admin-001}

**Given** an admin user is on the product management interface
**When** they submit a new product with name (2-200 chars), description (up to 5000 chars), price (≥ $0.01), and category
**Then** the product is created in 'draft' status with a UUID, creator and timestamp are logged, and admin can add up to 10 images with one marked as primary

**Error Cases:**
- Name too short/long → Show length validation error
- Price below minimum → Show 'Price must be at least $0.01' error
- Category not selected → Show 'Category is required' error
- Duplicate name exists → Show warning (not blocking)

**Traceability:**
- Source: Add Product operation
- Source: Product entity
- Decision: [AMB-ENT-001](decisions.md#amb-ent-001)
- Decision: [AMB-ENT-002](decisions.md#amb-ent-002)
- Decision: [AMB-ENT-003](decisions.md#amb-ent-003)
- Decision: [AMB-ENT-004](decisions.md#amb-ent-004)
- Decision: [AMB-OP-033](decisions.md#amb-op-033)
- Decision: [AMB-OP-034](decisions.md#amb-op-034)
- Decision: [AMB-OP-035](decisions.md#amb-op-035)
- Decision: [AMB-OP-036](decisions.md#amb-op-036)
- Decision: [AMB-OP-037](decisions.md#amb-op-037)

---

## AC-ADMIN-002 – Edit product {#ac-admin-002}

**Given** an admin is editing an existing product
**When** they modify product attributes and save
**Then** the product is updated, last modified by/date is tracked, price changes are logged in audit trail, and carts with this product reflect the new price

**Error Cases:**
- Concurrent edit by another admin → Last write wins with conflict warning
- Validation errors → Same as Add Product

**Traceability:**
- Source: Edit Product operation
- Decision: [AMB-ENT-007](decisions.md#amb-ent-007)
- Decision: [AMB-OP-038](decisions.md#amb-op-038)
- Decision: [AMB-OP-039](decisions.md#amb-op-039)
- Decision: [AMB-OP-040](decisions.md#amb-op-040)

---

## AC-ADMIN-003 – Remove product {#ac-admin-003}

**Given** an admin wants to remove a product
**When** they click delete and confirm
**Then** the product is soft-deleted (preserved for order history), removed from any carts (with notification on next cart view), and no longer appears in product listings

**Error Cases:**
- Confirmation cancelled → No action taken

**Traceability:**
- Source: Remove Product operation
- Decision: [AMB-ENT-005](decisions.md#amb-ent-005)
- Decision: [AMB-OP-041](decisions.md#amb-op-041)
- Decision: [AMB-OP-042](decisions.md#amb-op-042)
- Decision: [AMB-OP-043](decisions.md#amb-op-043)

---

## AC-ADMIN-004 – View all orders {#ac-admin-004}

**Given** an admin is on the order management interface
**When** they view the order list
**Then** orders are displayed with order number, date, customer, total, and status, with options to filter by status/date range/customer and search by order number

**Error Cases:**
- No orders match filters → Show empty state

**Traceability:**
- Source: View All Orders operation
- Decision: [AMB-OP-044](decisions.md#amb-op-044)
- Decision: [AMB-OP-045](decisions.md#amb-op-045)
- Decision: [AMB-OP-046](decisions.md#amb-op-046)

---

## AC-ADMIN-005 – Update order status {#ac-admin-005}

**Given** an admin is viewing an order
**When** they change the status following valid transitions (pending→confirmed→shipped→delivered)
**Then** the status is updated, change is logged with admin user and timestamp, customer receives email notification (for confirmed/shipped/delivered)

**Error Cases:**
- Invalid status transition → Show error with allowed transitions
- Setting to shipped without tracking number → Allow (tracking is optional)

**Traceability:**
- Source: Update Order Status operation
- Decision: [AMB-ENT-018](decisions.md#amb-ent-018)
- Decision: [AMB-OP-047](decisions.md#amb-op-047)
- Decision: [AMB-OP-048](decisions.md#amb-op-048)
- Decision: [AMB-OP-049](decisions.md#amb-op-049)
- Decision: [AMB-OP-050](decisions.md#amb-op-050)

---

## AC-ADMIN-006 – Manage inventory {#ac-admin-006}

**Given** an admin is managing product inventory
**When** they adjust stock level (set absolute value or delta adjustment)
**Then** the inventory is updated (tracked per variant), full audit trail is recorded with before/after values and optional reason, and low stock alerts are triggered if threshold is crossed

**Error Cases:**
- Resulting quantity would be negative → Block with 'Inventory cannot go below zero' error
- Bulk CSV import with invalid data → Show row-by-row validation errors

**Traceability:**
- Source: Manage Inventory operation
- Source: Inventory entity
- Decision: [AMB-ENT-009](decisions.md#amb-ent-009)
- Decision: [AMB-ENT-040](decisions.md#amb-ent-040)
- Decision: [AMB-ENT-041](decisions.md#amb-ent-041)
- Decision: [AMB-ENT-042](decisions.md#amb-ent-042)
- Decision: [AMB-ENT-043](decisions.md#amb-ent-043)
- Decision: [AMB-OP-051](decisions.md#amb-op-051)
- Decision: [AMB-OP-052](decisions.md#amb-op-052)
- Decision: [AMB-OP-053](decisions.md#amb-op-053)
- Decision: [AMB-OP-054](decisions.md#amb-op-054)
- Decision: [AMB-OP-055](decisions.md#amb-op-055)

---

## AC-VAR-001 – Product variant selection {#ac-var-001}

**Given** a product has multiple variants (size, color, or other configurable types)
**When** a customer views the product
**Then** all variant options are displayed, each variant shows its own stock status, and price may differ if variant has price override

**Error Cases:**
- Selected variant out of stock → Show out of stock for that variant, allow selecting others

**Traceability:**
- Source: ProductVariant entity
- Decision: [AMB-ENT-008](decisions.md#amb-ent-008)
- Decision: [AMB-ENT-009](decisions.md#amb-ent-009)
- Decision: [AMB-ENT-010](decisions.md#amb-ent-010)
- Decision: [AMB-ENT-011](decisions.md#amb-ent-011)

---

## AC-CAT-001 – Category hierarchy {#ac-cat-001}

**Given** products are organized in categories
**When** a customer browses by category
**Then** categories are displayed in up to 3 levels of hierarchy, with products showing their primary category and optional secondary categories

**Error Cases:**
- Category is inactive → Category and its products not shown in customer view

**Traceability:**
- Source: Category entity
- Decision: [AMB-ENT-006](decisions.md#amb-ent-006)
- Decision: [AMB-ENT-037](decisions.md#amb-ent-037)
- Decision: [AMB-ENT-038](decisions.md#amb-ent-038)
- Decision: [AMB-ENT-039](decisions.md#amb-ent-039)

---

