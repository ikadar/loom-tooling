# Technical Specifications

Generated: 2025-12-29T12:39:52+01:00

---

## TS-BR-AUTH-001 – Registration and email verification for orders

**Rule:** Only registered customers with verified email addresses can place orders

**Implementation Approach:**
Validate customer status and email_verified flag at checkout initiation and order submission endpoints. Implement middleware/guard that checks authentication token claims for registration status and email verification.

**Validation Points:**
- POST /api/checkout/initiate
- POST /api/orders
- Checkout UI component mount

**Data Requirements:**

| Field | Type | Constraints | Source |
|-------|------|-------------|--------|
| customer_id | uuid | NOT NULL, valid customer reference | Authentication token |
| customer_status | enum | MUST be 'registered' | customers.status |
| email_verified | boolean | MUST be true | customers.email_verified |

**Error Handling:**

| Condition | Error Code | Message | HTTP Status |
|-----------|------------|---------|-------------|
| customer_status != 'registered' | REGISTRATION_REQUIRED | You must create an account to place orders | 403 |
| email_verified == false | EMAIL_NOT_VERIFIED | Please verify your email address before placing orders | 403 |

**Traceability:**
- BR: BR-AUTH-001
- Related ACs: [AC-ORDER-001 AC-CUST-001]

---

## TS-BR-STOCK-001 – Stock validation for cart addition

**Rule:** Products must have available stock to be added to cart

**Implementation Approach:**
Query inventory.available_quantity before cart addition. For products with variants, check variant-specific inventory. Disable add-to-cart button in UI when stock is 0. Re-validate at checkout.

**Validation Points:**
- POST /api/cart/items
- PUT /api/cart/items/{id}
- Product detail page load
- POST /api/checkout/initiate

**Data Requirements:**
| Field | Type | Constraints | Source |
|-------|------|-------------|--------|
| product_id | uuid | NOT NULL, valid product reference | Request body |
| variant_id | uuid | NULL or valid variant reference | Request body |
| available_quantity | integer | >= 0 | inventory.available_quantity |

**Error Handling:**
| Condition | Error Code | Message | HTTP Status |
|-----------|------------|---------|-------------|
| available_quantity == 0 | OUT_OF_STOCK | This product is currently out of stock | 400 |

**Traceability:**
- BR: BR-STOCK-001
- Related ACs: [AC-CART-001 AC-PROD-002]

---

## TS-BR-STOCK-002 – No backorders validation

**Rule:** Orders cannot be placed for quantities exceeding available stock

**Implementation Approach:**
Strict stock validation at checkout. Compare each cart item quantity against current inventory. Cart can temporarily hold excess quantity but checkout is blocked. Use SELECT FOR UPDATE to prevent race conditions during checkout.

**Validation Points:**
- POST /api/checkout/initiate
- POST /api/orders

**Data Requirements:**
| Field | Type | Constraints | Source |
|-------|------|-------------|--------|
| cart_item_quantity | integer | >= 1 | cart_items.quantity |
| available_quantity | integer | >= 0 | inventory.available_quantity |
| product_variant_id | uuid | NOT NULL | cart_items.variant_id |

**Error Handling:**
| Condition | Error Code | Message | HTTP Status |
|-----------|------------|---------|-------------|
| cart_item_quantity > available_quantity | INSUFFICIENT_STOCK | Requested quantity exceeds available stock for {product_name}. Available: {available_quantity} | 400 |

**Traceability:**
- BR: BR-STOCK-002
- Related ACs: [AC-ORDER-001 AC-CART-008]

---

## TS-BR-STOCK-003 – Non-negative inventory enforcement

**Rule:** Stock levels must never be negative

**Implementation Approach:**
Database CHECK constraint on inventory.quantity >= 0. Application-level validation before any inventory adjustment. Use database transaction with row-level locking for concurrent operations.

**Validation Points:**
- PUT /api/admin/inventory/{id}
- POST /api/admin/inventory/bulk
- Order placement transaction
- Order cancellation transaction

**Data Requirements:**
| Field | Type | Constraints | Source |
|-------|------|-------------|--------|
| quantity | integer | >= 0, CHECK constraint | inventory.quantity |
| adjustment_delta | integer | result must be >= 0 | Request body |

**Error Handling:**
| Condition | Error Code | Message | HTTP Status |
|-----------|------------|---------|-------------|
| resulting_quantity < 0 | INVALID_STOCK_LEVEL | Stock level cannot be negative. Minimum is 0. | 400 |

**Traceability:**
- BR: BR-STOCK-003
- Related ACs: [AC-ADMIN-007]

---

## TS-BR-SHIP-001 – Free shipping threshold calculation

**Rule:** Orders with subtotal of $50 or more (before tax, after discounts) qualify for free shipping

**Implementation Approach:**
Calculate shipping cost during checkout totals computation. Subtotal = sum(item_price * quantity) - discounts. If subtotal >= 50.00 then shipping_cost = 0, else shipping_cost = 5.99. Display threshold progress in cart UI.

**Validation Points:**
- GET /api/cart/totals
- POST /api/checkout/initiate
- Cart summary component

**Data Requirements:**
| Field | Type | Constraints | Source |
|-------|------|-------------|--------|
| subtotal | decimal(10,2) | >= 0 | Calculated from cart items |
| discount_total | decimal(10,2) | >= 0 | Applied discounts |
| shipping_cost | decimal(10,2) | 0.00 or 5.99 | Calculated |
| free_shipping_threshold | decimal(10,2) | = 50.00 | Configuration |

**Traceability:**
- BR: BR-SHIP-001
- Related ACs: [AC-ORDER-002 AC-ORDER-003]

---

## TS-BR-CANCEL-001 – Order cancellation window validation

**Rule:** Orders can only be cancelled before shipping

**Implementation Approach:**
Check order status before processing cancellation request. Only allow cancellation if status IN ('pending', 'confirmed'). Hide cancel button in UI for shipped/delivered orders. Return appropriate error for invalid states.

**Validation Points:**
- POST /api/orders/{id}/cancel
- Order detail page load

**Data Requirements:**
| Field | Type | Constraints | Source |
|-------|------|-------------|--------|
| order_id | uuid | NOT NULL, valid order reference | URL parameter |
| order_status | enum | MUST be IN ('pending', 'confirmed') | orders.status |
| customer_id | uuid | MUST match authenticated user | orders.customer_id |

**Error Handling:**
| Condition | Error Code | Message | HTTP Status |
|-----------|------------|---------|-------------|
| order_status NOT IN ('pending', 'confirmed') | CANCELLATION_NOT_ALLOWED | Orders that have been shipped or delivered cannot be cancelled. Please contact support for assistance. | 400 |

**Traceability:**
- BR: BR-CANCEL-001
- Related ACs: [AC-ORDER-006]

---

## TS-BR-ORDER-001 – Order status state machine

**Rule:** Orders must follow defined state machine: pending→confirmed→shipped→delivered, with cancelled as terminal state from pending or confirmed only

**Implementation Approach:**
Implement state machine pattern for order status transitions. Define allowed transitions map: {pending: [confirmed, cancelled], confirmed: [shipped, cancelled], shipped: [delivered]}. Validate every status update against transition rules.

**Validation Points:**
- PUT /api/admin/orders/{id}/status
- Order status update service

**Data Requirements:**
| Field | Type | Constraints | Source |
|-------|------|-------------|--------|
| current_status | enum | IN ('pending', 'confirmed', 'shipped', 'delivered', 'cancelled') | orders.status |
| new_status | enum | Must be valid transition from current_status | Request body |
| allowed_transitions | map | Predefined state machine rules | Application configuration |

**Error Handling:**
| Condition | Error Code | Message | HTTP Status |
|-----------|------------|---------|-------------|
| new_status not in allowed_transitions[current_status] | INVALID_STATUS_TRANSITION | Cannot transition order from '{current_status}' to '{new_status}' | 400 |

**Traceability:**
- BR: BR-ORDER-001
- Related ACs: [AC-ADMIN-005 AC-ADMIN-006]

---

## TS-BR-ORDER-002 – Order immutability enforcement

**Rule:** Orders cannot be modified after placement; customer must cancel and reorder

**Implementation Approach:**
Do not expose PUT/PATCH endpoints for order items, quantities, or addresses. Only allow status updates via dedicated endpoint. Order modification requests return error with guidance to cancel and reorder.

**Validation Points:**
- API route definitions
- Order service layer

**Data Requirements:**
| Field | Type | Constraints | Source |
|-------|------|-------------|--------|
| order_id | uuid | NOT NULL | URL parameter |
| order_status | enum | Any non-cart status | orders.status |

**Error Handling:**
| Condition | Error Code | Message | HTTP Status |
|-----------|------------|---------|-------------|
| attempt to modify order items/address after placement | ORDER_MODIFICATION_NOT_ALLOWED | Orders cannot be modified after placement. Please cancel this order and create a new one. | 400 |

**Traceability:**
- BR: BR-ORDER-002
- Related ACs: [AC-ORDER-001]

---

## TS-BR-ORDER-003 – Single shipment constraint

**Rule:** Each order is shipped as a single unit; partial shipments not supported

**Implementation Approach:**
Database schema enforces one-to-one relationship between orders and shipments. No split shipment functionality in admin UI or API. All order items transition to shipped status together.

**Validation Points:**
- Database schema design
- Order fulfillment service

**Data Requirements:**
| Field | Type | Constraints | Source |
|-------|------|-------------|--------|
| order_id | uuid | UNIQUE in shipments table | orders.id |
| shipment_id | uuid | One per order | shipments.id |

**Traceability:**
- BR: BR-ORDER-003
- Related ACs: [AC-ADMIN-006]

---

## TS-BR-ORDER-004 – Price snapshot at order time

**Rule:** Order items must capture the product price at the time of order placement

**Implementation Approach:**
Copy product price to order_items.unit_price during order creation. Do not reference products.price for historical orders. Store as immutable field with no update capability.

**Validation Points:**
- Order creation service
- Order item insertion

**Data Requirements:**
| Field | Type | Constraints | Source |
|-------|------|-------------|--------|
| unit_price | decimal(10,2) | > 0, immutable after creation | products.price at order time |
| product_id | uuid | Reference to product | cart_items.product_id |
| quantity | integer | >= 1 | cart_items.quantity |

**Traceability:**
- BR: BR-ORDER-004
- Related ACs: [AC-ORDER-001 AC-ADMIN-002]

---

## TS-BR-PRICE-001 – Minimum product price validation

**Rule:** Product prices must be at least $0.01

**Implementation Approach:**
Validate price >= 0.01 at product creation and update. Database CHECK constraint for price >= 0.01. Reject zero or negative prices with validation error.

**Validation Points:**
- POST /api/admin/products
- PUT /api/admin/products/{id}

**Data Requirements:**
| Field | Type | Constraints | Source |
|-------|------|-------------|--------|
| price | decimal(10,2) | >= 0.01, CHECK constraint | Request body |

**Error Handling:**
| Condition | Error Code | Message | HTTP Status |
|-----------|------------|---------|-------------|
| price < 0.01 | INVALID_PRICE | Product price must be at least $0.01 | 400 |

**Traceability:**
- BR: BR-PRICE-001
- Related ACs: [AC-ADMIN-001]

---

## TS-BR-PRICE-002 – Cart current price display

**Rule:** Cart items always reflect current product prices, not historical prices

**Implementation Approach:**
Join cart_items with products on every cart fetch to get current prices. Compare against cart_items.added_at_price for change detection. Display notification when price differs from when item was added.

**Validation Points:**
- GET /api/cart
- Cart page component

**Data Requirements:**
| Field | Type | Constraints | Source |
|-------|------|-------------|--------|
| current_price | decimal(10,2) | >= 0.01 | products.price |
| added_at_price | decimal(10,2) | Stored at add time | cart_items.added_at_price |
| price_changed | boolean | Calculated field | current_price != added_at_price |

**Traceability:**
- BR: BR-PRICE-002
- Related ACs: [AC-CART-009]

---

## TS-BR-CART-001 – Cart quantity limits validation

**Rule:** Cart item quantity must be between 1 and available stock

**Implementation Approach:**
Validate quantity >= 1 and <= available_stock at add and update operations. Auto-limit quantity to available stock with notification rather than rejecting. UI quantity selector bounded by stock.

**Validation Points:**
- POST /api/cart/items
- PUT /api/cart/items/{id}

**Data Requirements:**
| Field | Type | Constraints | Source |
|-------|------|-------------|--------|
| quantity | integer | >= 1 | Request body |
| available_stock | integer | >= 0 | inventory.available_quantity |
| quantity_limited | boolean | Response flag | Calculated when quantity adjusted |

**Error Handling:**
| Condition | Error Code | Message | HTTP Status |
|-----------|------------|---------|-------------|
| quantity < 1 | INVALID_QUANTITY | Quantity must be at least 1 | 400 |
| quantity > available_stock | QUANTITY_EXCEEDS_STOCK | Quantity limited to available stock ({available_stock}) | 200 |

**Traceability:**
- BR: BR-CART-001
- Related ACs: [AC-CART-001 AC-CART-002 AC-CART-005]

---

## TS-BR-CART-002 – Cart expiration cleanup

**Rule:** Inactive carts expire after 30 days for logged-in users and 7 days for guests

**Implementation Approach:**
Scheduled background job runs daily to delete expired carts. Query carts where last_activity < (NOW - expiry_period). Expiry period determined by customer_id presence (NULL = guest = 7 days, NOT NULL = registered = 30 days).

**Validation Points:**
- Background job scheduler
- Cart cleanup service

**Data Requirements:**
| Field | Type | Constraints | Source |
|-------|------|-------------|--------|
| last_activity | timestamp | NOT NULL, updated on cart operations | carts.last_activity |
| customer_id | uuid | NULL for guests | carts.customer_id |
| session_id | string | NOT NULL for guests | carts.session_id |
| guest_expiry_days | integer | = 7 | Configuration |
| registered_expiry_days | integer | = 30 | Configuration |

**Traceability:**
- BR: BR-CART-002
- Related ACs: [AC-CART-003]

---

## TS-BR-CUST-001 – Unique customer email enforcement

**Rule:** Customer email addresses must be unique in the system

**Implementation Approach:**
Database UNIQUE constraint on customers.email. Case-insensitive comparison (store lowercase). Check for existence before insert to provide friendly error message.

**Validation Points:**
- POST /api/customers/register
- Database constraint

**Data Requirements:**
| Field | Type | Constraints | Source |
|-------|------|-------------|--------|
| email | varchar(255) | UNIQUE, NOT NULL, lowercase, valid email format | Request body |

**Error Handling:**
| Condition | Error Code | Message | HTTP Status |
|-----------|------------|---------|-------------|
| email already exists | DUPLICATE_EMAIL | An account with this email address already exists | 409 |

**Traceability:**
- BR: BR-CUST-001
- Related ACs: [AC-CUST-001]

---

## TS-BR-CUST-002 – Password requirements validation

**Rule:** Customer passwords must be at least 8 characters with at least one number

**Implementation Approach:**
Validate password with regex pattern: ^(?=.*\d).{8,}$. Check at registration and password change endpoints. Return specific error for each failed requirement.

**Validation Points:**
- POST /api/customers/register
- PUT /api/customers/password

**Data Requirements:**
| Field | Type | Constraints | Source |
|-------|------|-------------|--------|
| password | string | length >= 8, contains at least one digit | Request body |

**Error Handling:**
| Condition | Error Code | Message | HTTP Status |
|-----------|------------|---------|-------------|
| password.length < 8 | PASSWORD_TOO_SHORT | Password must be at least 8 characters | 400 |
| password does not contain digit | PASSWORD_MISSING_NUMBER | Password must contain at least one number | 400 |

**Traceability:**
- BR: BR-CUST-002
- Related ACs: [AC-CUST-001]

---

## TS-BR-CUST-003 – Email verification requirement

**Rule:** Customers must verify their email address before placing orders

**Implementation Approach:**
Set email_verified = false on registration. Send verification email with secure token. Verify token at confirmation endpoint and set email_verified = true. Block checkout if email_verified = false.

**Validation Points:**
- POST /api/customers/register
- GET /api/customers/verify-email
- POST /api/checkout/initiate

**Data Requirements:**
| Field | Type | Constraints | Source |
|-------|------|-------------|--------|
| email_verified | boolean | default false | customers.email_verified |
| verification_token | string | Secure random token, expires in 24h | email_verifications table |

**Error Handling:**
| Condition | Error Code | Message | HTTP Status |
|-----------|------------|---------|-------------|
| email_verified == false at checkout | EMAIL_NOT_VERIFIED | Please verify your email address to complete your order | 403 |

**Traceability:**
- BR: BR-CUST-003
- Related ACs: [AC-CUST-001 AC-ORDER-001]

---

## TS-BR-SHIP-002 – Domestic shipping only validation

**Rule:** Shipping is only supported within the domestic country

**Implementation Approach:**
Validate country code against allowed list (e.g., ['US']). Reject international addresses at address creation and checkout. Display supported countries in UI.

**Validation Points:**
- POST /api/customers/addresses
- PUT /api/customers/addresses/{id}
- POST /api/checkout/initiate

**Data Requirements:**
| Field | Type | Constraints | Source |
|-------|------|-------------|--------|
| country | varchar(2) | ISO country code, IN supported list | Request body |
| supported_countries | array | ['US'] | Configuration |

**Error Handling:**
| Condition | Error Code | Message | HTTP Status |
|-----------|------------|---------|-------------|
| country not in supported_countries | UNSUPPORTED_COUNTRY | We currently only ship within the United States | 400 |

**Traceability:**
- BR: BR-SHIP-002
- Related ACs: [AC-CUST-002]

---

## TS-BR-PAY-001 – Supported payment types validation

**Rule:** Only Visa, Mastercard, American Express credit cards and PayPal are accepted

**Implementation Approach:**
Validate card type from BIN/IIN lookup or PayPal integration. Stripe/PayPal SDK handles card type detection. Reject unsupported card types before tokenization.

**Validation Points:**
- POST /api/customers/payment-methods
- Payment method UI component

**Data Requirements:**
| Field | Type | Constraints | Source |
|-------|------|-------------|--------|
| payment_type | enum | IN ('visa', 'mastercard', 'amex', 'paypal') | Payment provider response |
| card_number | string | Validated by payment provider | Request body (not stored) |

**Error Handling:**
| Condition | Error Code | Message | HTTP Status |
|-----------|------------|---------|-------------|
| card_type not in supported types | UNSUPPORTED_CARD_TYPE | We accept Visa, Mastercard, American Express, and PayPal | 400 |
| invalid card number | INVALID_CARD | Invalid card number | 400 |

**Traceability:**
- BR: BR-PAY-001
- Related ACs: [AC-CUST-003]

---

## TS-BR-PAY-002 – Payment authorization before order

**Rule:** Full payment must be authorized before an order is created

**Implementation Approach:**
Transactional order creation: 1) Create payment intent, 2) Authorize payment, 3) If successful create order and decrement inventory, 4) If failed rollback transaction. Use database transaction with payment provider webhook confirmation.

**Validation Points:**
- POST /api/orders
- Payment service layer

**Data Requirements:**
| Field | Type | Constraints | Source |
|-------|------|-------------|--------|
| payment_intent_id | string | From payment provider | Stripe/PayPal |
| authorization_status | enum | MUST be 'succeeded' | Payment provider response |
| order_total | decimal(10,2) | > 0 | Calculated from cart |

**Error Handling:**
| Condition | Error Code | Message | HTTP Status |
|-----------|------------|---------|-------------|
| payment authorization fails | PAYMENT_FAILED | Payment could not be processed. Please try again or use a different payment method. | 402 |

**Traceability:**
- BR: BR-PAY-002
- Related ACs: [AC-ORDER-001]

---

## TS-BR-PROD-001 – Product soft deletion

**Rule:** Products with order history must be soft-deleted to preserve historical references

**Implementation Approach:**
Add deleted_at timestamp column to products. Delete operation sets deleted_at = NOW() instead of removing row. Query filters exclude deleted products by default. Hard delete blocked if order_items reference exists.

**Validation Points:**
- DELETE /api/admin/products/{id}
- Product queries

**Data Requirements:**
| Field | Type | Constraints | Source |
|-------|------|-------------|--------|
| deleted_at | timestamp | NULL for active, NOT NULL for deleted | products.deleted_at |
| has_orders | boolean | Calculated from order_items join | Existence check |

**Traceability:**
- BR: BR-PROD-001
- Related ACs: [AC-ADMIN-003]

---

## TS-BR-PROD-002 – Product name constraints validation

**Rule:** Product names must be between 2 and 200 characters

**Implementation Approach:**
Validate name length at product creation and update. Trim whitespace before validation. Database VARCHAR(200) with CHECK constraint for minimum length.

**Validation Points:**
- POST /api/admin/products
- PUT /api/admin/products/{id}

**Data Requirements:**
| Field | Type | Constraints | Source |
|-------|------|-------------|--------|
| name | varchar(200) | length >= 2 AND length <= 200, NOT NULL | Request body |

**Error Handling:**
| Condition | Error Code | Message | HTTP Status |
|-----------|------------|---------|-------------|
| name.length < 2 | INVALID_NAME | Product name must be at least 2 characters | 400 |
| name.length > 200 | NAME_TOO_LONG | Product name cannot exceed 200 characters | 400 |

**Traceability:**
- BR: BR-PROD-002
- Related ACs: [AC-ADMIN-001]

---

## TS-BR-PROD-003 – Product category requirement

**Rule:** Every product must have a primary category assigned

**Implementation Approach:**
Validate primary_category_id is not null at product creation and update. Foreign key constraint to categories table. UI requires category selection before save.

**Validation Points:**
- POST /api/admin/products
- PUT /api/admin/products/{id}

**Data Requirements:**
| Field | Type | Constraints | Source |
|-------|------|-------------|--------|
| primary_category_id | uuid | NOT NULL, valid category reference | Request body |

**Error Handling:**
| Condition | Error Code | Message | HTTP Status |
|-----------|------------|---------|-------------|
| primary_category_id is null | CATEGORY_REQUIRED | Please select a category for this product | 400 |

**Traceability:**
- BR: BR-PROD-003
- Related ACs: [AC-ADMIN-001]

---

## TS-BR-VAR-001 – Unique variant SKU enforcement

**Rule:** Each product variant must have a unique SKU

**Implementation Approach:**
Database UNIQUE constraint on product_variants.sku. Check for SKU existence before insert. SKU format validation (alphanumeric with hyphens).

**Validation Points:**
- POST /api/admin/products/{id}/variants
- PUT /api/admin/variants/{id}

**Data Requirements:**
| Field | Type | Constraints | Source |
|-------|------|-------------|--------|
| sku | varchar(50) | UNIQUE, NOT NULL, alphanumeric with hyphens | Request body |

**Error Handling:**
| Condition | Error Code | Message | HTTP Status |
|-----------|------------|---------|-------------|
| sku already exists | DUPLICATE_SKU | SKU '{sku}' is already in use | 409 |

**Traceability:**
- BR: BR-VAR-001
- Related ACs: [AC-VAR-001]

---

## TS-BR-VAR-002 – Variant selection requirement

**Rule:** For products with variants, a specific variant must be selected before adding to cart

**Implementation Approach:**
Check if product has variants before cart addition. If product.has_variants = true, require variant_id in request. Disable add-to-cart button until variant selected in UI.

**Validation Points:**
- POST /api/cart/items
- Product detail page component

**Data Requirements:**
| Field | Type | Constraints | Source |
|-------|------|-------------|--------|
| product_id | uuid | NOT NULL | Request body |
| variant_id | uuid | Required if product has variants | Request body |
| has_variants | boolean | Calculated from variants count | products join product_variants |

**Error Handling:**
| Condition | Error Code | Message | HTTP Status |
|-----------|------------|---------|-------------|
| has_variants == true AND variant_id is null | VARIANT_REQUIRED | Please select a size/color/option before adding to cart | 400 |

**Traceability:**
- BR: BR-VAR-002
- Related ACs: [AC-CART-001]

---

## TS-BR-CAT-001 – Category nesting depth validation

**Rule:** Categories can only be nested up to 3 levels deep

**Implementation Approach:**
Calculate depth on category creation by traversing parent chain. Store depth field for efficient queries. Reject category creation if parent.depth >= 3.

**Validation Points:**
- POST /api/admin/categories
- PUT /api/admin/categories/{id}

**Data Requirements:**
| Field | Type | Constraints | Source |
|-------|------|-------------|--------|
| parent_id | uuid | NULL for root categories | Request body |
| depth | integer | <= 3, calculated from parent chain | categories.depth |

**Error Handling:**
| Condition | Error Code | Message | HTTP Status |
|-----------|------------|---------|-------------|
| resulting depth > 3 | CATEGORY_DEPTH_EXCEEDED | Categories cannot be nested more than 3 levels deep | 400 |

**Traceability:**
- BR: BR-CAT-001
- Related ACs: [AC-CAT-001]

---

## TS-BR-INV-001 – Inventory restoration on cancellation

**Rule:** When an order is cancelled, all reserved inventory must be automatically restored

**Implementation Approach:**
In order cancellation transaction: for each order_item, UPDATE inventory SET available_quantity = available_quantity + order_item.quantity. Use database transaction to ensure atomicity with status change.

**Validation Points:**
- POST /api/orders/{id}/cancel
- Order cancellation service

**Data Requirements:**
| Field | Type | Constraints | Source |
|-------|------|-------------|--------|
| order_items | array | All items from cancelled order | order_items table |
| variant_id | uuid | NOT NULL | order_items.variant_id |
| quantity | integer | >= 1 | order_items.quantity |

**Traceability:**
- BR: BR-INV-001
- Related ACs: [AC-ORDER-006]

---

## TS-BR-INV-002 – Inventory decrement on payment

**Rule:** Inventory is decremented immediately upon successful payment

**Implementation Approach:**
In order creation transaction: after payment authorization succeeds, UPDATE inventory SET available_quantity = available_quantity - cart_item.quantity for all items. Use SELECT FOR UPDATE to prevent race conditions.

**Validation Points:**
- POST /api/orders
- Order creation service

**Data Requirements:**
| Field | Type | Constraints | Source |
|-------|------|-------------|--------|
| cart_items | array | All items from cart | cart_items table |
| variant_id | uuid | NOT NULL | cart_items.variant_id |
| quantity | integer | >= 1, <= available | cart_items.quantity |

**Traceability:**
- BR: BR-INV-002
- Related ACs: [AC-ORDER-001]

---

## TS-BR-AUDIT-001 – Order status audit trail

**Rule:** All order status changes must be logged with timestamp and actor

**Implementation Approach:**
Create order_status_history table. On every status update, INSERT record with order_id, old_status, new_status, changed_by (user_id), changed_at (timestamp). Implement via database trigger or service layer interceptor.

**Validation Points:**
- PUT /api/admin/orders/{id}/status
- Order status update service

**Data Requirements:**
| Field | Type | Constraints | Source |
|-------|------|-------------|--------|
| order_id | uuid | NOT NULL, FK to orders | orders.id |
| old_status | enum | Previous status value | orders.status before update |
| new_status | enum | New status value | Request body |
| changed_by | uuid | NOT NULL, FK to users | Authenticated user |
| changed_at | timestamp | NOT NULL, default NOW() | Database |

**Traceability:**
- BR: BR-AUDIT-001
- Related ACs: [AC-ADMIN-005]

---

## TS-BR-AUDIT-002 – Inventory audit trail

**Rule:** All inventory changes must be logged with before/after values, reason, and actor

**Implementation Approach:**
Create inventory_history table. On every inventory update, INSERT record with variant_id, quantity_before, quantity_after, reason, modified_by, modified_at. Capture reason from request or system action type.

**Validation Points:**
- PUT /api/admin/inventory/{id}
- POST /api/admin/inventory/bulk
- Order placement
- Order cancellation

**Data Requirements:**
| Field | Type | Constraints | Source |
|-------|------|-------------|--------|
| variant_id | uuid | NOT NULL, FK to product_variants | inventory.variant_id |
| quantity_before | integer | >= 0 | inventory.quantity before update |
| quantity_after | integer | >= 0 | inventory.quantity after update |
| reason | varchar(255) | NULL allowed | Request body or system action |
| modified_by | uuid | NOT NULL | Authenticated user or system |
| modified_at | timestamp | NOT NULL, default NOW() | Database |

**Traceability:**
- BR: BR-AUDIT-002
- Related ACs: [AC-ADMIN-007 AC-ADMIN-008]

---

## TS-BR-AUDIT-003 – Price change audit trail

**Rule:** Product price changes must be tracked with who modified and when

**Implementation Approach:**
Create price_history table. On product.price update, INSERT record with product_id, old_price, new_price, modified_by, modified_at. Implement via service layer or database trigger.

**Validation Points:**
- PUT /api/admin/products/{id}

**Data Requirements:**
| Field | Type | Constraints | Source |
|-------|------|-------------|--------|
| product_id | uuid | NOT NULL, FK to products | products.id |
| old_price | decimal(10,2) | >= 0.01 | products.price before update |
| new_price | decimal(10,2) | >= 0.01 | Request body |
| modified_by | uuid | NOT NULL, FK to users | Authenticated admin user |
| modified_at | timestamp | NOT NULL, default NOW() | Database |

**Traceability:**
- BR: BR-AUDIT-003
- Related ACs: [AC-ADMIN-002]

---

