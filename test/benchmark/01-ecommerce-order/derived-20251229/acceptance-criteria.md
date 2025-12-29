# Acceptance Criteria

Generated: 2025-12-29T11:31:02+01:00

---

## AC-PROD-001 – Browse products with filtering

**Given** A customer is on the product listing page
**When** They apply filters for category, price range, or availability
**Then** Only products matching all applied filters are displayed, sorted by the selected option (name, price asc/desc, newest), with 20 products per page

**Error Cases:**
- No products match filters → Show empty state with 'No products found' message and option to clear filters
- Invalid price range (min > max) → Show validation error

**Traceability:**
- Source: Domain Model: Browse Products
- Decision: AMB-OP-001
- Decision: AMB-OP-002
- Decision: AMB-OP-003
- Decision: AMB-OP-004

---

## AC-PROD-002 – View out-of-stock products

**Given** A product has zero available inventory
**When** A customer views the product listing or product detail page
**Then** The product is displayed with an 'out of stock' indicator and the add to cart button is disabled

**Traceability:**
- Source: Domain Model: Browse Products
- Decision: AMB-OP-004

---

## AC-CART-001 – Add product to cart

**Given** A customer is viewing an in-stock product with variant selected (if applicable)
**When** They click 'Add to Cart' with quantity between 1 and available stock
**Then** The item is added to their cart, a toast notification appears with 'View Cart' option, and the cart item count updates

**Error Cases:**
- Product out of stock → Add to cart button disabled, show 'Out of Stock' message
- Quantity exceeds available stock → Limit quantity to available stock with notification
- Variant not selected (for products with variants) → Add to cart button disabled until variant selected

**Traceability:**
- Source: Domain Model: Add to Cart
- Decision: AMB-OP-005
- Decision: AMB-OP-006
- Decision: AMB-OP-007
- Decision: AMB-OP-008
- Decision: AMB-OP-009
- Decision: AMB-OP-010

---

## AC-CART-002 – Add existing item to cart

**Given** A customer has Product A (quantity 2) in their cart
**When** They add Product A again with quantity 1
**Then** The cart item quantity is incremented to 3 (not a new line item created)

**Error Cases:**
- New total exceeds available stock → Quantity set to max available with notification

**Traceability:**
- Source: Domain Model: Add to Cart
- Decision: AMB-OP-005
- Decision: AMB-OP-006

---

## AC-CART-003 – Guest user cart

**Given** A user is not logged in
**When** They add items to cart
**Then** A session-based cart is created and persists for 7 days of inactivity

**Traceability:**
- Source: Domain Model: Cart
- Decision: AMB-OP-007
- Decision: AMB-ENT-012
- Decision: AMB-ENT-013

---

## AC-CART-004 – Cart merge on login

**Given** A guest user has items in their session cart
**When** They log in to an existing account that also has cart items
**Then** The guest cart items are merged with the existing cart, combining quantities for duplicate products

**Error Cases:**
- Merged quantity exceeds stock → Limit to available stock with notification

**Traceability:**
- Source: Domain Model: Cart
- Decision: AMB-ENT-013

---

## AC-CART-005 – Update cart item quantity

**Given** A customer has a product in their cart
**When** They update the quantity to a valid number (1 to available stock)
**Then** The cart item quantity is updated, cart total recalculates, and current product price is used

**Error Cases:**
- Quantity set to 0 → Item is removed from cart
- Quantity exceeds available stock → Quantity limited to available stock with notification
- Negative quantity → Validation error, quantity unchanged

**Traceability:**
- Source: Domain Model: Update Cart Quantity
- Decision: AMB-OP-011
- Decision: AMB-OP-012
- Decision: AMB-OP-013
- Decision: AMB-ENT-015
- Decision: AMB-ENT-016

---

## AC-CART-006 – Remove item from cart

**Given** A customer has items in their cart
**When** They click remove on a cart item
**Then** The item is immediately removed (no confirmation) and an 'Undo' option is shown temporarily

**Traceability:**
- Source: Domain Model: Remove from Cart
- Decision: AMB-OP-014

---

## AC-CART-007 – Remove last item from cart

**Given** A customer has exactly one item in their cart
**When** They remove that item
**Then** An empty cart state is displayed with a 'Continue Shopping' link

**Traceability:**
- Source: Domain Model: Remove from Cart
- Decision: AMB-OP-015

---

## AC-CART-008 – Cart item becomes out of stock

**Given** A customer has a product in their cart
**When** That product goes out of stock while in the cart
**Then** The item remains in cart with a warning indicator, but checkout is blocked until resolved

**Traceability:**
- Source: Domain Model: Cart
- Decision: AMB-ENT-014

---

## AC-CART-009 – Price change notification

**Given** A customer has a product in their cart
**When** The product price changes and they view the cart
**Then** The cart shows the current price and a notification that the price has changed

**Traceability:**
- Source: Domain Model: CartItem
- Decision: AMB-ENT-016

---

## AC-ORDER-001 – Place order as registered customer

**Given** A registered customer with verified email has items in cart, valid shipping address, and valid payment method
**When** They submit the order with subtotal $60 (before tax)
**Then** Payment is authorized, inventory is decremented, order is created with status 'pending', cart is cleared, confirmation email is queued, and order number format is ORD-YYYY-NNNNN

**Error Cases:**
- Customer not registered → REGISTRATION_REQUIRED error, redirect to registration
- Customer email not verified → EMAIL_NOT_VERIFIED error, prompt to verify
- Payment authorization fails → PAYMENT_FAILED error, cart preserved, allow retry
- Item out of stock at checkout → INSUFFICIENT_STOCK error for affected items, allow modification
- Invalid shipping address → ADDRESS_VALIDATION_ERROR, highlight invalid fields

**Traceability:**
- Source: Domain Model: Place Order
- Decision: AMB-OP-016
- Decision: AMB-OP-017
- Decision: AMB-OP-018
- Decision: AMB-OP-019
- Decision: AMB-ENT-017
- Decision: AMB-ENT-030

---

## AC-ORDER-002 – Free shipping qualification

**Given** A customer has cart items with subtotal $50 or more (before tax, after discounts)
**When** They proceed to checkout
**Then** Free shipping is automatically applied and displayed

**Traceability:**
- Source: Domain Model: Place Order
- Decision: AMB-OP-020

---

## AC-ORDER-003 – Standard shipping charge

**Given** A customer has cart items with subtotal below $50
**When** They proceed to checkout
**Then** Standard shipping fee of $5.99 is applied

**Traceability:**
- Source: Domain Model: Place Order
- Decision: AMB-OP-021

---

## AC-ORDER-004 – Tax calculation

**Given** A customer enters a US shipping address
**When** The order total is calculated
**Then** State tax is calculated based on the shipping address state

**Traceability:**
- Source: Domain Model: Place Order
- Decision: AMB-OP-022

---

## AC-ORDER-005 – Track order status

**Given** A logged-in customer has placed orders
**When** They view their order history
**Then** All orders are displayed with pagination, showing order number, items, quantities, prices, status, shipping address, and tracking number (if available)

**Traceability:**
- Source: Domain Model: Track Order Status
- Decision: AMB-OP-025
- Decision: AMB-OP-026
- Decision: AMB-OP-027

---

## AC-ORDER-006 – Cancel pending order

**Given** A customer has an order with status 'pending' or 'confirmed'
**When** They request cancellation with an optional reason
**Then** Order status changes to 'cancelled', payment is refunded to original method, inventory is restored, and cancellation email is sent

**Error Cases:**
- Order status is 'shipped' or 'delivered' → CANCELLATION_NOT_ALLOWED error, suggest contacting support

**Traceability:**
- Source: Domain Model: Cancel Order
- Decision: AMB-OP-028
- Decision: AMB-OP-029
- Decision: AMB-OP-030
- Decision: AMB-OP-031
- Decision: AMB-OP-032
- Decision: AMB-ENT-044

---

## AC-ORDER-007 – Order confirmation email

**Given** An order has been successfully placed
**When** The system processes the order
**Then** A confirmation email is queued asynchronously containing order number, items, total, shipping address, and estimated delivery

**Error Cases:**
- Email service unavailable → Queue for retry, log failure, order creation not blocked

**Traceability:**
- Source: Domain Model: Send Confirmation Email
- Decision: AMB-OP-056
- Decision: AMB-OP-057
- Decision: AMB-OP-058
- Decision: AMB-OP-059
- Decision: AMB-OP-060

---

## AC-CUST-001 – Customer registration

**Given** A guest user is on the registration page
**When** They submit valid registration data (email, password 8+ chars with number, first name, last name, optional phone)
**Then** Account is created, verification email is sent, and user cannot place orders until email is verified

**Error Cases:**
- Email already registered → DUPLICATE_EMAIL error
- Password too short → PASSWORD_TOO_SHORT error
- Password missing number → PASSWORD_MISSING_NUMBER error
- Invalid email format → INVALID_EMAIL error

**Traceability:**
- Source: Domain Model: Customer
- Decision: AMB-ENT-024
- Decision: AMB-ENT-025
- Decision: AMB-ENT-026
- Decision: AMB-ENT-030

---

## AC-CUST-002 – Manage shipping addresses

**Given** A registered customer is logged in
**When** They add a new shipping address with all required fields
**Then** The address is saved and they can set it as default; multiple addresses are supported

**Error Cases:**
- Missing required field → INVALID_ADDRESS error with field specification
- Country not supported (non-domestic) → UNSUPPORTED_COUNTRY error

**Traceability:**
- Source: Domain Model: ShippingAddress
- Decision: AMB-ENT-027
- Decision: AMB-ENT-031
- Decision: AMB-ENT-032
- Decision: AMB-ENT-033

---

## AC-CUST-003 – Manage payment methods

**Given** A registered customer is logged in
**When** They add a payment method (Visa, Mastercard, Amex, or PayPal)
**Then** The payment method is tokenized via Stripe/PayPal and stored; they can set it as default

**Error Cases:**
- Invalid card number → INVALID_CARD error
- Unsupported card type → UNSUPPORTED_CARD_TYPE error

**Traceability:**
- Source: Domain Model: PaymentMethod
- Decision: AMB-ENT-028
- Decision: AMB-ENT-034
- Decision: AMB-ENT-035
- Decision: AMB-ENT-036

---

## AC-ADMIN-001 – Add new product

**Given** An admin user is on the product management page
**When** They create a product with name (2-200 chars), price (>$0.01), category, description (up to 5000 chars), and images (up to 10)
**Then** Product is created in 'draft' status with UUID, creator and timestamp logged

**Error Cases:**
- Name too short (<2 chars) → INVALID_NAME error
- Name too long (>200 chars) → NAME_TOO_LONG error
- Price zero or negative → INVALID_PRICE error
- No category selected → CATEGORY_REQUIRED error
- Duplicate name exists → Warning displayed (not blocking)

**Traceability:**
- Source: Domain Model: Add Product
- Decision: AMB-ENT-001
- Decision: AMB-ENT-002
- Decision: AMB-ENT-003
- Decision: AMB-ENT-004
- Decision: AMB-OP-033
- Decision: AMB-OP-034
- Decision: AMB-OP-035
- Decision: AMB-OP-036
- Decision: AMB-OP-037

---

## AC-ADMIN-002 – Edit product

**Given** An admin is editing an existing product
**When** They modify product attributes and save
**Then** Changes are saved, last modified by/date is updated, price changes are logged; existing orders retain original price, carts get updated price

**Error Cases:**
- Concurrent edit by another admin → Last write wins with conflict warning displayed

**Traceability:**
- Source: Domain Model: Edit Product
- Decision: AMB-ENT-007
- Decision: AMB-OP-038
- Decision: AMB-OP-039
- Decision: AMB-OP-040

---

## AC-ADMIN-003 – Remove product

**Given** An admin wants to remove a product
**When** They confirm removal after seeing impact summary
**Then** Product is soft-deleted, removed from all active carts (users notified on next cart view), but preserved for order history

**Traceability:**
- Source: Domain Model: Remove Product
- Decision: AMB-ENT-005
- Decision: AMB-OP-041
- Decision: AMB-OP-042
- Decision: AMB-OP-043

---

## AC-ADMIN-004 – View all orders

**Given** An admin is on the order management page
**When** They apply filters (status, date range, customer) or search by order number
**Then** Matching orders are displayed showing order number, date, customer, total, status with pagination

**Traceability:**
- Source: Domain Model: View All Orders
- Decision: AMB-OP-044
- Decision: AMB-OP-045
- Decision: AMB-OP-046

---

## AC-ADMIN-005 – Update order status

**Given** An admin is viewing an order with status 'pending'
**When** They update status to 'confirmed'
**Then** Status changes, customer receives email notification, change is logged with admin user and timestamp

**Error Cases:**
- Invalid transition (e.g., pending→delivered) → INVALID_STATUS_TRANSITION error

**Traceability:**
- Source: Domain Model: Update Order Status
- Decision: AMB-ENT-018
- Decision: AMB-ENT-022
- Decision: AMB-OP-047
- Decision: AMB-OP-048
- Decision: AMB-OP-049
- Decision: AMB-OP-050

---

## AC-ADMIN-006 – Mark order as shipped

**Given** An admin is viewing a 'confirmed' order
**When** They update status to 'shipped' with optional tracking number
**Then** Status changes, tracking number is stored, customer receives email with tracking link to carrier website

**Traceability:**
- Source: Domain Model: Update Order Status
- Decision: AMB-OP-027
- Decision: AMB-OP-049

---

## AC-ADMIN-007 – Manage inventory

**Given** An admin is on the inventory management page
**When** They update stock for a product variant (absolute or delta adjustment) with optional reason
**Then** Stock is updated, audit trail records before/after values with timestamp and admin user, low stock alert triggers if threshold crossed

**Error Cases:**
- Result would be negative → INVALID_STOCK_LEVEL error, minimum is zero

**Traceability:**
- Source: Domain Model: Manage Inventory
- Decision: AMB-ENT-009
- Decision: AMB-ENT-040
- Decision: AMB-ENT-041
- Decision: AMB-ENT-042
- Decision: AMB-ENT-043
- Decision: AMB-OP-051
- Decision: AMB-OP-052
- Decision: AMB-OP-053
- Decision: AMB-OP-054
- Decision: AMB-OP-055

---

## AC-ADMIN-008 – Bulk inventory update via CSV

**Given** An admin has a CSV file with SKU and quantity columns
**When** They upload the file for bulk inventory update
**Then** All valid rows are processed, audit trail created for each change, summary shows success/failure counts

**Error Cases:**
- Invalid SKU → Row skipped with error in summary
- Invalid quantity format → Row skipped with error in summary

**Traceability:**
- Source: Domain Model: Manage Inventory
- Decision: AMB-OP-054
- Decision: AMB-OP-055

---

## AC-VAR-001 – Add product variant

**Given** An admin is editing a product
**When** They add a variant with unique SKU, variant attributes (size, color, etc.), and optional price override
**Then** Variant is created with its own inventory tracking, inheriting parent price if no override specified

**Error Cases:**
- Duplicate SKU → DUPLICATE_SKU error

**Traceability:**
- Source: Domain Model: ProductVariant
- Decision: AMB-ENT-008
- Decision: AMB-ENT-009
- Decision: AMB-ENT-010
- Decision: AMB-ENT-011

---

## AC-CAT-001 – Manage categories

**Given** An admin is on category management
**When** They create/edit a category with name, description, image, display order, and optional parent (up to 3 levels)
**Then** Category is saved with active/inactive status support

**Error Cases:**
- Nesting exceeds 3 levels → CATEGORY_DEPTH_EXCEEDED error

**Traceability:**
- Source: Domain Model: Category
- Decision: AMB-ENT-037
- Decision: AMB-ENT-038
- Decision: AMB-ENT-039

---

## AC-GDPR-001 – Customer data erasure

**Given** A customer requests full data deletion (GDPR)
**When** The request is processed
**Then** All personal data is erased, order history is anonymized (preserved for business records)

**Traceability:**
- Source: Domain Model: Customer
- Decision: AMB-ENT-029

---

