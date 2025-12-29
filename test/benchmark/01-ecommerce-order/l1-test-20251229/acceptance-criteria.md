# Acceptance Criteria

Generated: 2025-12-29T13:37:22+01:00

---

## AC-CUST-001 – Customer Registration

**Given** a visitor with email 'john@example.com' that is not registered
**When** they submit registration with email 'john@example.com', password 'SecurePass1', first name 'John', and last name 'Doe'
**Then** a customer account is created with status 'unverified', and a verification email is sent to 'john@example.com'

**Error Cases:**
- Email already registered → DUPLICATE_EMAIL error, suggest login or password reset
- Password less than 8 characters → INVALID_PASSWORD error with requirements message
- Password without number → INVALID_PASSWORD error with requirements message
- Invalid email format → INVALID_EMAIL error

**Traceability:**
- Source: Customer entity
- Source: register operation
- Decision: AMB-ENT-024
- Decision: AMB-ENT-025
- Decision: AMB-ENT-026
- Decision: AMB-ENT-030

---

## AC-CUST-002 – Email Verification Required for Orders

**Given** a registered customer with unverified email
**When** they attempt to place an order
**Then** the order is rejected with EMAIL_NOT_VERIFIED error and prompt to verify email

**Error Cases:**
- Verification link expired → VERIFICATION_EXPIRED error with resend option

**Traceability:**
- Source: Place Order operation
- Decision: AMB-ENT-030

---

## AC-CUST-003 – Save Multiple Shipping Addresses

**Given** a registered customer with one existing address
**When** they add a new shipping address with street '456 Oak Ave', city 'Boston', state 'MA', postal code '02101', country 'US', recipient 'John Doe'
**Then** the address is saved to their account, and they can select it for future orders

**Error Cases:**
- Missing required field → INVALID_ADDRESS error specifying missing field
- Invalid postal code format → INVALID_POSTAL_CODE error

**Traceability:**
- Source: ShippingAddress entity
- Source: Customer entity
- Decision: AMB-ENT-027
- Decision: AMB-ENT-031
- Decision: AMB-ENT-032
- Decision: AMB-ENT-033

---

## AC-CUST-004 – GDPR Data Erasure

**Given** a registered customer requesting account deletion
**When** they confirm full data erasure request
**Then** all personal data is permanently deleted, order history is anonymized, and confirmation is sent

**Error Cases:**
- Pending orders exist → PENDING_ORDERS_EXIST error, must wait for completion

**Traceability:**
- Source: Customer entity
- Decision: AMB-ENT-029

---

## AC-PROD-001 – Browse Products with Filters

**Given** a product catalog with products in multiple categories and price ranges
**When** a customer filters by category 'Electronics', price range $50-$200, and availability 'in stock'
**Then** only products matching all criteria are displayed, sorted by newest first, with 20 products per page

**Error Cases:**
- No products match filters → Empty result with 'No products found' message and clear filters option

**Traceability:**
- Source: Browse Products operation
- Decision: AMB-OP-001
- Decision: AMB-OP-002
- Decision: AMB-OP-003

---

## AC-PROD-002 – View Out-of-Stock Products

**Given** a product 'Widget Pro' with zero available stock
**When** a customer views the product listing
**Then** the product is displayed with 'Out of Stock' indicator and 'Add to Cart' button is disabled

**Traceability:**
- Source: Browse Products operation
- Decision: AMB-OP-004

---

## AC-PROD-003 – Product with Variants Display

**Given** a product 'T-Shirt' with variants for size (S, M, L) and color (Red, Blue)
**When** a customer views the product detail page
**Then** variant options are displayed, 'Add to Cart' is disabled until a variant is selected, and price/stock updates based on selection

**Error Cases:**
- Selected variant out of stock → Show 'Out of Stock' for that variant, allow different selection

**Traceability:**
- Source: ProductVariant entity
- Source: Add to Cart operation
- Decision: AMB-ENT-008
- Decision: AMB-ENT-009
- Decision: AMB-ENT-010
- Decision: AMB-ENT-011
- Decision: AMB-OP-010

---

## AC-CART-001 – Add Product to Cart

**Given** a product 'Widget' with 10 units in stock and price $25.00
**When** a customer adds 2 units to their cart
**Then** cart shows 'Widget' with quantity 2, unit price $25.00, subtotal $50.00, and toast notification appears with 'View Cart' option

**Error Cases:**
- Product out of stock → OUT_OF_STOCK error, cannot add
- Requested quantity exceeds stock → INSUFFICIENT_STOCK error, show available quantity

**Traceability:**
- Source: Add to Cart operation
- Source: Cart entity
- Decision: AMB-OP-005
- Decision: AMB-OP-006
- Decision: AMB-OP-009

---

## AC-CART-002 – Add Existing Item to Cart

**Given** a cart containing 'Widget' with quantity 2
**When** the customer adds 3 more units of 'Widget'
**Then** cart shows 'Widget' with quantity 5 (not a separate line item)

**Error Cases:**
- Total quantity would exceed stock → Limit to available stock with notification

**Traceability:**
- Source: Add to Cart operation
- Decision: AMB-OP-005
- Decision: AMB-OP-006

---

## AC-CART-003 – Guest User Cart

**Given** a visitor not logged in
**When** they add products to cart and later log in
**Then** guest cart items are merged with any existing account cart items

**Error Cases:**
- Merged quantity exceeds stock → Limit to stock with notification

**Traceability:**
- Source: Cart entity
- Source: Add to Cart operation
- Decision: AMB-ENT-012
- Decision: AMB-ENT-013
- Decision: AMB-OP-007

---

## AC-CART-004 – Update Cart Item Quantity

**Given** a cart with 'Widget' quantity 3 and 10 units in stock
**When** customer changes quantity to 5
**Then** quantity updates to 5, subtotal recalculates, cart total updates

**Error Cases:**
- Quantity set to 0 → Item removed from cart
- Quantity exceeds stock (e.g., 15) → Limited to 10 with notification 'Only 10 available'

**Traceability:**
- Source: Update Cart Quantity operation
- Decision: AMB-OP-011
- Decision: AMB-OP-012
- Decision: AMB-OP-013

---

## AC-CART-005 – Remove Item from Cart

**Given** a cart with 'Widget' and 'Gadget'
**When** customer removes 'Widget'
**Then** 'Widget' is removed immediately, 'Undo' option appears for 5 seconds, cart total recalculates

**Traceability:**
- Source: Remove from Cart operation
- Decision: AMB-OP-014

---

## AC-CART-006 – Remove Last Item from Cart

**Given** a cart with only 'Widget'
**When** customer removes 'Widget'
**Then** empty cart state is displayed with 'Continue Shopping' link

**Traceability:**
- Source: Remove from Cart operation
- Decision: AMB-OP-015

---

## AC-CART-007 – Cart Item Goes Out of Stock

**Given** a cart with 'Widget' quantity 2
**When** the product goes out of stock (inventory becomes 0)
**Then** item remains in cart with 'Out of Stock' warning, checkout is blocked until resolved

**Traceability:**
- Source: Cart entity
- Decision: AMB-ENT-014

---

## AC-CART-008 – Cart Price Update Notification

**Given** a cart with 'Widget' at $25.00
**When** admin changes 'Widget' price to $30.00
**Then** on next cart view, price shows $30.00 with notification 'Price updated from $25.00 to $30.00'

**Traceability:**
- Source: CartItem entity
- Decision: AMB-ENT-016

---

## AC-CART-009 – Cart Expiration

**Given** a logged-in customer cart inactive for 30 days
**When** expiration period passes
**Then** cart items are cleared, customer notified on next login

**Traceability:**
- Source: Cart entity
- Decision: AMB-ENT-012

---

## AC-ORDER-001 – Place Order with Free Shipping

**Given** a verified customer with cart subtotal $75.00 (after any discounts)
**When** they place an order with valid shipping address and credit card
**Then** order is created with status 'pending', shipping cost $0.00, payment is authorized, inventory is decremented, cart is cleared, confirmation email is queued

**Error Cases:**
- Payment authorization fails → PAYMENT_FAILED error, cart preserved, allow retry
- Item out of stock at checkout → STOCK_UNAVAILABLE error for affected items, allow modification
- Customer not verified → EMAIL_NOT_VERIFIED error

**Traceability:**
- Source: Place Order operation
- Source: Order entity
- Decision: AMB-OP-016
- Decision: AMB-OP-017
- Decision: AMB-OP-018
- Decision: AMB-OP-019
- Decision: AMB-OP-020

---

## AC-ORDER-002 – Place Order with Paid Shipping

**Given** a verified customer with cart subtotal $35.00
**When** they place an order
**Then** order is created with shipping cost $5.99, total $40.99 plus applicable tax

**Traceability:**
- Source: Place Order operation
- Decision: AMB-OP-020
- Decision: AMB-OP-021

---

## AC-ORDER-003 – Order Number Generation

**Given** the current year is 2024 and last order was ORD-2024-00042
**When** a new order is placed
**Then** order number ORD-2024-00043 is assigned

**Traceability:**
- Source: Order entity
- Decision: AMB-ENT-017

---

## AC-ORDER-004 – Order Price Snapshot

**Given** cart with 'Widget' at $25.00
**When** order is placed
**Then** order line item stores price $25.00, even if product price later changes

**Traceability:**
- Source: Order entity
- Decision: AMB-ENT-019

---

## AC-ORDER-005 – Tax Calculation

**Given** order with subtotal $100.00, shipping to California (7.25% tax rate)
**When** order totals are calculated
**Then** tax is $7.25, calculated on subtotal before shipping

**Traceability:**
- Source: Place Order operation
- Decision: AMB-OP-022

---

## AC-ORDER-006 – Track Order Status

**Given** a customer with orders ORD-2024-00001 (delivered) and ORD-2024-00042 (shipped)
**When** they view their order history
**Then** all orders displayed with pagination, showing order number, date, items summary, total, and current status

**Traceability:**
- Source: Track Order Status operation
- Decision: AMB-OP-025
- Decision: AMB-OP-026

---

## AC-ORDER-007 – View Order with Tracking

**Given** order ORD-2024-00042 with status 'shipped' and tracking number '1Z999AA10123456784'
**When** customer views order details
**Then** tracking number is displayed with link to carrier tracking page

**Traceability:**
- Source: Track Order Status operation
- Decision: AMB-OP-027

---

## AC-ORDER-008 – Cancel Pending Order

**Given** order ORD-2024-00042 with status 'pending'
**When** customer requests cancellation with reason 'Changed mind'
**Then** status changes to 'cancelled', inventory is restored, refund is initiated to original payment method, cancellation email is sent

**Traceability:**
- Source: Cancel Order operation
- Decision: AMB-OP-028
- Decision: AMB-OP-029
- Decision: AMB-OP-030
- Decision: AMB-OP-031
- Decision: AMB-OP-032

---

## AC-ORDER-009 – Cancel Confirmed Order

**Given** order ORD-2024-00042 with status 'confirmed'
**When** customer requests cancellation
**Then** status changes to 'cancelled', inventory is restored, refund is initiated

**Traceability:**
- Source: Cancel Order operation
- Decision: AMB-OP-028

---

## AC-ORDER-010 – Cannot Cancel Shipped Order

**Given** order ORD-2024-00042 with status 'shipped'
**When** customer attempts cancellation
**Then** cancellation is rejected with ORDER_ALREADY_SHIPPED error

**Traceability:**
- Source: Cancel Order operation
- Source: Order state transitions
- Decision: AMB-OP-028
- Decision: AMB-ENT-018

---

## AC-ORDER-011 – Order Confirmation Email

**Given** order ORD-2024-00042 just placed
**When** order placement completes successfully
**Then** email queued with order number, items list, quantities, prices, total, shipping address, and estimated delivery date

**Error Cases:**
- Email service unavailable → Email queued for retry, order not blocked

**Traceability:**
- Source: Send Confirmation Email operation
- Decision: AMB-OP-056
- Decision: AMB-OP-057
- Decision: AMB-OP-058
- Decision: AMB-OP-059
- Decision: AMB-OP-060

---

## AC-ADMIN-001 – Add New Product

**Given** an admin user on the product management page
**When** they create product with name 'Super Widget', description 'A great widget', price $29.99, category 'Electronics'
**Then** product is created with status 'draft', UUID assigned, creator and timestamp logged

**Error Cases:**
- Name empty → NAME_REQUIRED error
- Name exceeds 200 chars → NAME_TOO_LONG error
- Price <= 0 → INVALID_PRICE error
- Category not selected → CATEGORY_REQUIRED error
- Duplicate name exists → Warning shown (not blocking)

**Traceability:**
- Source: Add Product operation
- Decision: AMB-ENT-001
- Decision: AMB-ENT-002
- Decision: AMB-ENT-003
- Decision: AMB-OP-033
- Decision: AMB-OP-034
- Decision: AMB-OP-035
- Decision: AMB-OP-037

---

## AC-ADMIN-002 – Add Product Images

**Given** a product 'Super Widget' in draft status
**When** admin uploads 3 images and sets one as primary
**Then** images are uploaded to cloud storage, URLs stored, primary image marked

**Error Cases:**
- More than 10 images → MAX_IMAGES_EXCEEDED error
- Invalid image format → INVALID_IMAGE_FORMAT error

**Traceability:**
- Source: Add Product operation
- Decision: AMB-ENT-004
- Decision: AMB-OP-036

---

## AC-ADMIN-003 – Edit Product Price

**Given** product 'Widget' with price $25.00, currently in customer carts
**When** admin changes price to $30.00
**Then** price updated, change logged with admin user and timestamp, carts show updated price on next view

**Traceability:**
- Source: Edit Product operation
- Decision: AMB-ENT-007
- Decision: AMB-OP-038
- Decision: AMB-OP-040

---

## AC-ADMIN-004 – Edit Product Concurrent Access

**Given** admin A and admin B both editing product 'Widget'
**When** admin A saves, then admin B saves different changes
**Then** admin B's changes overwrite admin A's, conflict warning shown to admin B

**Traceability:**
- Source: Edit Product operation
- Decision: AMB-OP-039

---

## AC-ADMIN-005 – Remove Product (Soft Delete)

**Given** product 'Widget' with existing order history
**When** admin confirms removal
**Then** product is soft-deleted (not visible to customers), order history preserved, cart items removed with notification

**Traceability:**
- Source: Remove Product operation
- Decision: AMB-ENT-005
- Decision: AMB-OP-041
- Decision: AMB-OP-042
- Decision: AMB-OP-043

---

## AC-ADMIN-006 – View All Orders with Filters

**Given** 100 orders in the system
**When** admin filters by status 'shipped', date range 'last 7 days', and searches order number 'ORD-2024'
**Then** matching orders displayed with order number, date, customer name, total, status

**Traceability:**
- Source: View All Orders operation
- Decision: AMB-OP-044
- Decision: AMB-OP-045

---

## AC-ADMIN-007 – Export Orders to CSV

**Given** filtered order list showing 25 orders
**When** admin clicks 'Export CSV'
**Then** CSV file downloads with order details for all filtered orders

**Traceability:**
- Source: View All Orders operation
- Decision: AMB-OP-046

---

## AC-ADMIN-008 – Update Order Status to Shipped

**Given** order ORD-2024-00042 with status 'confirmed'
**When** admin updates status to 'shipped' with tracking number '1Z999AA10123456784'
**Then** status changes to 'shipped', tracking number stored, status change logged, customer email sent

**Error Cases:**
- Invalid status transition (e.g., pending→shipped) → INVALID_STATUS_TRANSITION error

**Traceability:**
- Source: Update Order Status operation
- Decision: AMB-OP-047
- Decision: AMB-OP-048
- Decision: AMB-OP-049
- Decision: AMB-OP-050

---

## AC-ADMIN-009 – Manage Inventory - Set Quantity

**Given** product 'Widget' with current stock 50
**When** admin sets stock to 100 with reason 'Shipment received'
**Then** stock updated to 100, audit log created with before (50), after (100), reason, admin user, timestamp

**Error Cases:**
- Negative quantity → INVALID_QUANTITY error

**Traceability:**
- Source: Manage Inventory operation
- Decision: AMB-OP-051
- Decision: AMB-OP-052
- Decision: AMB-OP-053
- Decision: AMB-OP-055

---

## AC-ADMIN-010 – Manage Inventory - Adjust Delta

**Given** product 'Widget' with current stock 50
**When** admin adjusts by -5 with reason 'Damaged goods'
**Then** stock updated to 45, audit log created

**Error Cases:**
- Adjustment would result in negative stock → INSUFFICIENT_STOCK error

**Traceability:**
- Source: Manage Inventory operation
- Decision: AMB-OP-051
- Decision: AMB-OP-053

---

## AC-ADMIN-011 – Bulk Inventory Update via CSV

**Given** CSV file with 50 product SKUs and new quantities
**When** admin imports the CSV
**Then** all valid rows are updated, errors reported for invalid rows, summary shown

**Error Cases:**
- SKU not found → Row skipped with error
- Invalid quantity → Row skipped with error

**Traceability:**
- Source: Manage Inventory operation
- Decision: AMB-OP-054

---

## AC-ADMIN-012 – Low Stock Alert

**Given** product 'Widget' with low stock threshold 10 and current stock 15
**When** stock drops to 8 (via order or adjustment)
**Then** low stock alert is triggered for admin notification

**Traceability:**
- Source: Inventory entity
- Decision: AMB-ENT-041

---

## AC-ADMIN-013 – Create Category Hierarchy

**Given** parent category 'Electronics' exists
**When** admin creates subcategory 'Smartphones' under 'Electronics'
**Then** subcategory created with parent reference, display order configurable

**Error Cases:**
- Exceeds 3 levels of nesting → MAX_NESTING_EXCEEDED error

**Traceability:**
- Source: Category entity
- Decision: AMB-ENT-037
- Decision: AMB-ENT-038

---

## AC-PAY-001 – Pay with Credit Card

**Given** customer at checkout with Visa card ending 4242
**When** they submit payment
**Then** payment authorized via Stripe, tokenized reference stored, order created

**Error Cases:**
- Card declined → PAYMENT_DECLINED error with reason
- Invalid card number → INVALID_CARD error
- Unsupported card type (e.g., Discover) → UNSUPPORTED_CARD_TYPE error

**Traceability:**
- Source: PaymentMethod entity
- Source: Place Order operation
- Decision: AMB-ENT-034
- Decision: AMB-ENT-035
- Decision: AMB-ENT-036

---

## AC-PAY-002 – Pay with PayPal

**Given** customer at checkout selecting PayPal
**When** they complete PayPal authorization
**Then** payment authorized, PayPal reference stored, order created

**Error Cases:**
- PayPal authorization cancelled → PAYMENT_CANCELLED, return to checkout

**Traceability:**
- Source: PaymentMethod entity
- Decision: AMB-ENT-035

---

## AC-PAY-003 – Save Payment Method for Future

**Given** customer completes order with new credit card
**When** they select 'Save for future purchases'
**Then** tokenized payment method stored, available for future checkouts

**Traceability:**
- Source: Customer entity
- Source: PaymentMethod entity
- Decision: AMB-ENT-028
- Decision: AMB-ENT-036

---

