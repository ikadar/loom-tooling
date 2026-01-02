# Technical Specifications

Generated: 2026-01-02T18:27:24+01:00

---

## TS-BR-AUTH-001 – Customer registration and email verification for orders {#ts-br-auth-001}

**Rule:** Only registered customers with verified email addresses can place orders

**Implementation Approach:**
Checkout service validates customer exists and has emailVerified=true before order creation

**Validation Points:**
- POST /orders API
- CheckoutService.placeOrder()
- OrderService.create()

**Data Requirements:**
| Field | Type | Constraints | Source |
|-------|------|-------------|--------|
| customer_id | UUID | NOT NULL, FK to customers | Order aggregate |
| email_verified | boolean | must be true for order placement | Customer aggregate |

**Error Handling:**
| Condition | Error Code | Message | HTTP Status |
|-----------|------------|---------|-------------|
| customer not found or not authenticated | REGISTRATION_REQUIRED | You must be registered to place an order | 401 |
| customer.emailVerified == false | EMAIL_NOT_VERIFIED | Please verify your email address before placing an order | 403 |

**Traceability:**
- BR: [BR-AUTH-001](../l1/business-rules.md#br-auth-001)

---

## TS-BR-STOCK-001 – Stock availability validation for cart additions {#ts-br-stock-001}

**Rule:** Products must have available inventory to be added to cart

**Implementation Approach:**
Check inventory.availableStock > 0 on add to cart; allow cart items with warnings but validate at checkout

**Validation Points:**
- POST /cart/items API
- CartService.addItem()
- CheckoutService.validateCart()

**Data Requirements:**
| Field | Type | Constraints | Source |
|-------|------|-------------|--------|
| available_stock | integer | >= 0 | Inventory aggregate |
| product_id | UUID | NOT NULL | CartItem entity |
| variant_id | UUID | nullable, required if product has variants | CartItem entity |

**Error Handling:**
| Condition | Error Code | Message | HTTP Status |
|-----------|------------|---------|-------------|
| availableStock == 0 | OUT_OF_STOCK | This product is currently out of stock | 409 |
| availableStock < requested at checkout | INSUFFICIENT_STOCK | Not enough stock available for checkout | 409 |

**Traceability:**
- BR: [BR-STOCK-001](../l1/business-rules.md#br-stock-001)

---

## TS-BR-SHIP-001 – Free shipping threshold calculation {#ts-br-ship-001}

**Rule:** Orders with subtotal of $50 or more qualify for free shipping

**Implementation Approach:**
ShippingCalculator checks subtotal >= 50.00; returns 0 for shipping cost if met, otherwise $5.99

**Validation Points:**
- ShippingService.calculateCost()
- Cart.getShippingEstimate()
- Order.calculateTotals()

**Data Requirements:**
| Field | Type | Constraints | Source |
|-------|------|-------------|--------|
| subtotal | decimal(10,2) | calculated from line items after discounts, before tax | Order/Cart aggregate |
| shipping_cost | decimal(10,2) | 0.00 if subtotal >= 50.00, else 5.99 | Order aggregate |

**Traceability:**
- BR: [BR-SHIP-001](../l1/business-rules.md#br-ship-001)

---

## TS-BR-CANCEL-001 – Order cancellation status validation {#ts-br-cancel-001}

**Rule:** Orders can only be cancelled while in pending or confirmed status

**Implementation Approach:**
CancelOrder operation checks current status is in allowed set before proceeding

**Validation Points:**
- POST /orders/{id}/cancel API
- OrderService.cancel()

**Data Requirements:**
| Field | Type | Constraints | Source |
|-------|------|-------------|--------|
| status | OrderStatus enum | must be 'pending' or 'confirmed' to allow cancellation | Order aggregate |

**Error Handling:**
| Condition | Error Code | Message | HTTP Status |
|-----------|------------|---------|-------------|
| status NOT IN ('pending', 'confirmed') | CANCELLATION_NOT_ALLOWED | This order cannot be cancelled as it has already been shipped | 409 |

**Traceability:**
- BR: [BR-CANCEL-001](../l1/business-rules.md#br-cancel-001)

---

## TS-BR-PRICE-001 – Minimum product price validation {#ts-br-price-001}

**Rule:** Product price must be at least $0.01

**Implementation Approach:**
Validate price >= 0.01 on product create/update; enforce 2 decimal precision

**Validation Points:**
- POST /products API
- PUT /products/{id} API
- ProductService.create()
- ProductService.update()
- DB CHECK constraint

**Data Requirements:**
| Field | Type | Constraints | Source |
|-------|------|-------------|--------|
| price | decimal(10,2) | >= 0.01 | Product aggregate |

**Error Handling:**
| Condition | Error Code | Message | HTTP Status |
|-----------|------------|---------|-------------|
| price < 0.01 | INVALID_PRICE | Product price must be at least $0.01 | 400 |

**Traceability:**
- BR: [BR-PRICE-001](../l1/business-rules.md#br-price-001)

---

## TS-BR-PRICE-002 – Order price snapshot at placement {#ts-br-price-002}

**Rule:** Order items must capture prices at the time of order placement

**Implementation Approach:**
Copy current Product/Variant price to OrderItem.unitPrice during order creation; OrderItem prices are immutable after creation

**Validation Points:**
- OrderService.create()
- OrderItem entity constructor

**Data Requirements:**
| Field | Type | Constraints | Source |
|-------|------|-------------|--------|
| unit_price | decimal(10,2) | immutable after order creation | OrderItem entity |
| quantity | integer | >= 1 | OrderItem entity |
| line_total | decimal(10,2) | calculated: unit_price * quantity | OrderItem entity |

**Traceability:**
- BR: [BR-PRICE-002](../l1/business-rules.md#br-price-002)

---

## TS-BR-CART-001 – Cart prices reflect current product prices {#ts-br-cart-001}

**Rule:** Cart items always reflect current product prices

**Implementation Approach:**
Cart.calculateTotal() fetches current prices from Product/Variant; detect and flag price changes since item was added

**Validation Points:**
- GET /cart API
- CartService.getCart()
- Cart.calculateTotal()

**Data Requirements:**
| Field | Type | Constraints | Source |
|-------|------|-------------|--------|
| product_id | UUID | NOT NULL, FK to products | CartItem entity |
| price_at_add | decimal(10,2) | stored for change detection only | CartItem entity |
| current_price | decimal(10,2) | fetched from Product at runtime | Product aggregate |

**Traceability:**
- BR: [BR-CART-001](../l1/business-rules.md#br-cart-001)

---

## TS-BR-INV-001 – Non-negative inventory constraint {#ts-br-inv-001}

**Rule:** Inventory stock level must never be negative

**Implementation Approach:**
All inventory operations validate resulting quantity >= 0; use atomic decrement with check in database

**Validation Points:**
- InventoryService.adjust()
- InventoryService.decrement()
- DB CHECK constraint on available_stock

**Data Requirements:**
| Field | Type | Constraints | Source |
|-------|------|-------------|--------|
| available_stock | integer | >= 0, CHECK constraint | Inventory aggregate |

**Error Handling:**
| Condition | Error Code | Message | HTTP Status |
|-----------|------------|---------|-------------|
| resulting stock would be < 0 | INSUFFICIENT_STOCK | Insufficient stock available | 409 |

**Traceability:**
- BR: [BR-INV-001](../l1/business-rules.md#br-inv-001)

---

## TS-BR-INV-002 – No backorder validation at checkout {#ts-br-inv-002}

**Rule:** Orders cannot be placed for quantities exceeding available stock

**Implementation Approach:**
Checkout atomically validates all line items have sufficient stock before payment authorization

**Validation Points:**
- CheckoutService.validateStock()
- OrderService.create()

**Data Requirements:**
| Field | Type | Constraints | Source |
|-------|------|-------------|--------|
| available_stock | integer | >= order quantity | Inventory aggregate |
| quantity | integer | >= 1 | OrderItem/CartItem entity |

**Error Handling:**
| Condition | Error Code | Message | HTTP Status |
|-----------|------------|---------|-------------|
| any item quantity > available stock | BACKORDER_NOT_SUPPORTED | One or more items exceed available stock. Backorders are not supported. | 409 |

**Traceability:**
- BR: [BR-INV-002](../l1/business-rules.md#br-inv-002)

---

## TS-BR-INV-003 – Inventory restoration on order cancellation {#ts-br-inv-003}

**Rule:** Cancelled orders must automatically restore inventory

**Implementation Approach:**
Order cancellation handler iterates OrderItems and increments Inventory.availableStock for each

**Validation Points:**
- OrderService.cancel()
- InventoryService.restore()

**Data Requirements:**
| Field | Type | Constraints | Source |
|-------|------|-------------|--------|
| order_items | OrderItem[] | list of items to restore | Order aggregate |
| available_stock | integer | incremented by item quantity | Inventory aggregate |

**Traceability:**
- BR: [BR-INV-003](../l1/business-rules.md#br-inv-003)

---

## TS-BR-ORDER-001 – Order state machine transition validation {#ts-br-order-001}

**Rule:** Orders must follow defined state machine transitions

**Implementation Approach:**
OrderStateMachine validates transition legality before allowing status change

**Validation Points:**
- PUT /orders/{id}/status API
- OrderService.updateStatus()
- Order.transitionTo()

**Data Requirements:**
| Field | Type | Constraints | Source |
|-------|------|-------------|--------|
| status | OrderStatus enum | pending|confirmed|shipped|delivered|cancelled | Order aggregate |
| allowed_transitions | map | pending→[confirmed,cancelled], confirmed→[shipped,cancelled], shipped→[delivered], delivered→[], cancelled→[] | Order state machine |

**Error Handling:**
| Condition | Error Code | Message | HTTP Status |
|-----------|------------|---------|-------------|
| transition not in allowed set | INVALID_STATUS_TRANSITION | Cannot transition order from {current} to {requested} | 409 |

**Traceability:**
- BR: [BR-ORDER-001](../l1/business-rules.md#br-order-001)

---

## TS-BR-ORDER-002 – Order immutability after placement {#ts-br-order-002}

**Rule:** Orders cannot be modified after placement

**Implementation Approach:**
No PUT endpoints for order items/quantities/address; only status transitions exposed

**Validation Points:**
- API layer - no edit endpoints exposed
- OrderService - no modify methods

**Data Requirements:**
| Field | Type | Constraints | Source |
|-------|------|-------------|--------|
| order_items | OrderItem[] | immutable after creation | Order aggregate |
| shipping_address | ShippingAddress | immutable after creation | Order aggregate |

**Error Handling:**
| Condition | Error Code | Message | HTTP Status |
|-----------|------------|---------|-------------|
| attempt to modify placed order | ORDER_MODIFICATION_NOT_ALLOWED | Orders cannot be modified after placement. Please cancel and create a new order. | 403 |

**Traceability:**
- BR: [BR-ORDER-002](../l1/business-rules.md#br-order-002)

---

## TS-BR-ORDER-003 – Single shipment per order constraint {#ts-br-order-003}

**Rule:** Each order ships as a single unit

**Implementation Approach:**
Order has single shipment reference; no partial shipment API or functionality

**Validation Points:**
- Order entity design
- ShipmentService.ship()

**Data Requirements:**
| Field | Type | Constraints | Source |
|-------|------|-------------|--------|
| shipment_id | UUID | nullable, single value per order | Order aggregate |
| tracking_number | string | single tracking number per order | Shipment entity |

**Traceability:**
- BR: [BR-ORDER-003](../l1/business-rules.md#br-order-003)

---

## TS-BR-CUST-001 – Unique customer email validation {#ts-br-cust-001}

**Rule:** Email addresses must be unique across all customer accounts

**Implementation Approach:**
Check email uniqueness before customer creation; DB unique constraint as backup

**Validation Points:**
- POST /customers API
- CustomerService.register()
- DB UNIQUE constraint on email

**Data Requirements:**
| Field | Type | Constraints | Source |
|-------|------|-------------|--------|
| email | varchar(255) | UNIQUE, NOT NULL, valid email format | Customer aggregate |

**Error Handling:**
| Condition | Error Code | Message | HTTP Status |
|-----------|------------|---------|-------------|
| email already exists | EMAIL_ALREADY_REGISTERED | An account with this email address already exists | 409 |

**Traceability:**
- BR: [BR-CUST-001](../l1/business-rules.md#br-cust-001)

---

## TS-BR-CUST-002 – Password strength validation {#ts-br-cust-002}

**Rule:** Customer passwords must meet minimum security requirements

**Implementation Approach:**
Validate password on registration and password change; regex: ^(?=.*[0-9]).{8,}$

**Validation Points:**
- POST /customers API
- PUT /customers/{id}/password API
- CustomerService.register()
- CustomerService.changePassword()

**Data Requirements:**
| Field | Type | Constraints | Source |
|-------|------|-------------|--------|
| password | varchar(255) | min 8 chars, at least one digit, stored as hash | Customer aggregate |

**Error Handling:**
| Condition | Error Code | Message | HTTP Status |
|-----------|------------|---------|-------------|
| password < 8 chars OR no digit | WEAK_PASSWORD | Password must be at least 8 characters and contain at least one number | 400 |

**Traceability:**
- BR: [BR-CUST-002](../l1/business-rules.md#br-cust-002)

---

## TS-BR-CUST-003 – Email verification required for checkout {#ts-br-cust-003}

**Rule:** Customers must verify their email before placing orders

**Implementation Approach:**
Check customer.emailVerified == true in checkout flow before order creation

**Validation Points:**
- CheckoutService.placeOrder()
- POST /orders API

**Data Requirements:**
| Field | Type | Constraints | Source |
|-------|------|-------------|--------|
| email_verified | boolean | default false, set true after verification | Customer aggregate |

**Error Handling:**
| Condition | Error Code | Message | HTTP Status |
|-----------|------------|---------|-------------|
| emailVerified == false | EMAIL_NOT_VERIFIED | Please verify your email address before placing an order | 403 |

**Traceability:**
- BR: [BR-CUST-003](../l1/business-rules.md#br-cust-003)

---

## TS-BR-PROD-001 – Product soft delete for order history preservation {#ts-br-prod-001}

**Rule:** Products with order history must be soft-deleted

**Implementation Approach:**
Delete operation checks for OrderItem references; if found, set deletedAt timestamp instead of hard delete

**Validation Points:**
- DELETE /products/{id} API
- ProductService.delete()

**Data Requirements:**
| Field | Type | Constraints | Source |
|-------|------|-------------|--------|
| deleted_at | timestamp | nullable, set for soft delete | Product aggregate |
| has_orders | boolean | derived from OrderItem existence | Computed |

**Traceability:**
- BR: [BR-PROD-001](../l1/business-rules.md#br-prod-001)

---

## TS-BR-PROD-002 – Product name length validation {#ts-br-prod-002}

**Rule:** Product names must be between 2 and 200 characters

**Implementation Approach:**
Validate name length on create and update operations

**Validation Points:**
- POST /products API
- PUT /products/{id} API
- ProductService.create()
- ProductService.update()

**Data Requirements:**
| Field | Type | Constraints | Source |
|-------|------|-------------|--------|
| name | varchar(200) | length >= 2 AND <= 200, NOT NULL | Product aggregate |

**Error Handling:**
| Condition | Error Code | Message | HTTP Status |
|-----------|------------|---------|-------------|
| name.length < 2 OR name.length > 200 | INVALID_PRODUCT_NAME | Product name must be between 2 and 200 characters | 400 |

**Traceability:**
- BR: [BR-PROD-002](../l1/business-rules.md#br-prod-002)

---

## TS-BR-PROD-003 – Product description length validation {#ts-br-prod-003}

**Rule:** Product descriptions must not exceed 5000 characters

**Implementation Approach:**
Validate description length on create and update operations

**Validation Points:**
- POST /products API
- PUT /products/{id} API
- ProductService.create()
- ProductService.update()

**Data Requirements:**
| Field | Type | Constraints | Source |
|-------|------|-------------|--------|
| description | text | length <= 5000, nullable | Product aggregate |

**Error Handling:**
| Condition | Error Code | Message | HTTP Status |
|-----------|------------|---------|-------------|
| description.length > 5000 | DESCRIPTION_TOO_LONG | Product description cannot exceed 5000 characters | 400 |

**Traceability:**
- BR: [BR-PROD-003](../l1/business-rules.md#br-prod-003)

---

## TS-BR-PROD-004 – Product image limit and primary image validation {#ts-br-prod-004}

**Rule:** Products can have up to 10 images with exactly one primary image

**Implementation Approach:**
Validate image count on upload; enforce single primary via toggle logic (setting one as primary unsets others)

**Validation Points:**
- POST /products/{id}/images API
- ProductService.addImage()
- ProductService.setPrimaryImage()

**Data Requirements:**
| Field | Type | Constraints | Source |
|-------|------|-------------|--------|
| images | ProductImage[] | max 10 items | Product aggregate |
| is_primary | boolean | exactly one true per product | ProductImage entity |

**Error Handling:**
| Condition | Error Code | Message | HTTP Status |
|-----------|------------|---------|-------------|
| images.count >= 10 on upload | IMAGE_LIMIT_EXCEEDED | Products cannot have more than 10 images | 400 |
| no primary image set | PRIMARY_IMAGE_REQUIRED | Product must have exactly one primary image | 400 |

**Traceability:**
- BR: [BR-PROD-004](../l1/business-rules.md#br-prod-004)

---

## TS-BR-VAR-001 – Variant SKU uniqueness validation {#ts-br-var-001}

**Rule:** Each product variant must have a unique SKU

**Implementation Approach:**
Check SKU uniqueness on variant create/update; DB unique constraint as backup

**Validation Points:**
- POST /products/{id}/variants API
- PUT /variants/{id} API
- VariantService.create()
- DB UNIQUE constraint on sku

**Data Requirements:**
| Field | Type | Constraints | Source |
|-------|------|-------------|--------|
| sku | varchar(100) | UNIQUE, NOT NULL | ProductVariant entity |

**Error Handling:**
| Condition | Error Code | Message | HTTP Status |
|-----------|------------|---------|-------------|
| SKU already exists | DUPLICATE_SKU | A variant with this SKU already exists | 409 |

**Traceability:**
- BR: [BR-VAR-001](../l1/business-rules.md#br-var-001)

---

## TS-BR-VAR-002 – Variant selection required for products with variants {#ts-br-var-002}

**Rule:** Products with variants require variant selection before add to cart

**Implementation Approach:**
Add to cart validates: if product.hasVariants then variantId is required

**Validation Points:**
- POST /cart/items API
- CartService.addItem()

**Data Requirements:**
| Field | Type | Constraints | Source |
|-------|------|-------------|--------|
| has_variants | boolean | derived from variants.count > 0 | Product aggregate |
| variant_id | UUID | required if product has variants | CartItem entity |

**Error Handling:**
| Condition | Error Code | Message | HTTP Status |
|-----------|------------|---------|-------------|
| product has variants AND variantId is null | VARIANT_REQUIRED | Please select a variant before adding to cart | 400 |

**Traceability:**
- BR: [BR-VAR-002](../l1/business-rules.md#br-var-002)

---

## TS-BR-CART-002 – Cart quantity limited by available stock {#ts-br-cart-002}

**Rule:** Cart item quantity cannot exceed available inventory

**Implementation Approach:**
Cap requested quantity at availableStock on add/update; return actual quantity added

**Validation Points:**
- POST /cart/items API
- PUT /cart/items/{id} API
- CartService.addItem()
- CartService.updateQuantity()

**Data Requirements:**
| Field | Type | Constraints | Source |
|-------|------|-------------|--------|
| quantity | integer | >= 1 AND <= availableStock | CartItem entity |
| available_stock | integer | >= 0 | Inventory aggregate |

**Error Handling:**
| Condition | Error Code | Message | HTTP Status |
|-----------|------------|---------|-------------|
| requested quantity > availableStock | QUANTITY_EXCEEDS_STOCK | Requested quantity exceeds available stock. Quantity has been adjusted. | 200 |

**Traceability:**
- BR: [BR-CART-002](../l1/business-rules.md#br-cart-002)

---

## TS-BR-CART-003 – Cart expiration by user type {#ts-br-cart-003}

**Rule:** Carts expire after period of inactivity based on user type

**Implementation Approach:**
Background job runs daily to delete carts where lastActivityAt + TTL < now; TTL = 30 days for authenticated, 7 days for guest

**Validation Points:**
- CartExpirationJob (scheduled)
- CartService.getCart() - check expiration

**Data Requirements:**
| Field | Type | Constraints | Source |
|-------|------|-------------|--------|
| last_activity_at | timestamp | updated on any cart modification | Cart aggregate |
| customer_id | UUID | nullable, null for guest carts | Cart aggregate |

**Error Handling:**
| Condition | Error Code | Message | HTTP Status |
|-----------|------------|---------|-------------|
| cart accessed after expiration | CART_EXPIRED | Your cart has expired due to inactivity | 410 |

**Traceability:**
- BR: [BR-CART-003](../l1/business-rules.md#br-cart-003)

---

## TS-BR-PAY-001 – Payment authorization before order creation {#ts-br-pay-001}

**Rule:** Payment must be fully authorized before order is created

**Implementation Approach:**
Checkout flow: 1) validate cart, 2) authorize payment, 3) create order only on successful auth

**Validation Points:**
- CheckoutService.placeOrder()
- PaymentService.authorize()

**Data Requirements:**
| Field | Type | Constraints | Source |
|-------|------|-------------|--------|
| payment_authorization_id | string | NOT NULL on order creation | Order aggregate |
| payment_status | enum | authorized|captured|failed | Payment entity |

**Error Handling:**
| Condition | Error Code | Message | HTTP Status |
|-----------|------------|---------|-------------|
| payment authorization fails | PAYMENT_REQUIRED | Payment authorization failed. Please try again or use a different payment method. | 402 |

**Traceability:**
- BR: [BR-PAY-001](../l1/business-rules.md#br-pay-001)

---

## TS-BR-PAY-002 – Automatic refund on order cancellation {#ts-br-pay-002}

**Rule:** Cancelled orders must receive automatic refund

**Implementation Approach:**
Order cancellation triggers refund via payment gateway; refund is async but must be initiated

**Validation Points:**
- OrderService.cancel()
- PaymentService.refund()

**Data Requirements:**
| Field | Type | Constraints | Source |
|-------|------|-------------|--------|
| payment_authorization_id | string | required for refund | Order aggregate |
| refund_status | enum | pending|completed|failed | Payment entity |

**Traceability:**
- BR: [BR-PAY-002](../l1/business-rules.md#br-pay-002)

---

## TS-BR-SHIP-002 – Domestic shipping only validation {#ts-br-ship-002}

**Rule:** Shipping is only available within configured domestic country

**Implementation Approach:**
Validate ShippingAddress.country matches system configured domestic country

**Validation Points:**
- POST /addresses API
- CheckoutService.validateAddress()

**Data Requirements:**
| Field | Type | Constraints | Source |
|-------|------|-------------|--------|
| country | varchar(2) | ISO 3166-1 alpha-2, must match DOMESTIC_COUNTRY config | ShippingAddress entity |
| domestic_country | string | system configuration | Application config |

**Error Handling:**
| Condition | Error Code | Message | HTTP Status |
|-----------|------------|---------|-------------|
| country != DOMESTIC_COUNTRY | INTERNATIONAL_SHIPPING_NOT_AVAILABLE | We currently only ship within {country}. International shipping is not available. | 400 |

**Traceability:**
- BR: [BR-SHIP-002](../l1/business-rules.md#br-ship-002)

---

## TS-BR-ADDR-001 – Required shipping address fields validation {#ts-br-addr-001}

**Rule:** Shipping addresses must include all required fields

**Implementation Approach:**
Validate all required fields are present and non-empty on address create/update

**Validation Points:**
- POST /addresses API
- PUT /addresses/{id} API
- AddressService.validate()

**Data Requirements:**
| Field | Type | Constraints | Source |
|-------|------|-------------|--------|
| recipient_name | varchar(100) | NOT NULL, non-empty | ShippingAddress entity |
| street | varchar(200) | NOT NULL, non-empty | ShippingAddress entity |
| city | varchar(100) | NOT NULL, non-empty | ShippingAddress entity |
| state | varchar(100) | NOT NULL, non-empty | ShippingAddress entity |
| postal_code | varchar(20) | NOT NULL, non-empty | ShippingAddress entity |
| country | varchar(2) | NOT NULL, ISO 3166-1 alpha-2 | ShippingAddress entity |

**Error Handling:**
| Condition | Error Code | Message | HTTP Status |
|-----------|------------|---------|-------------|
| any required field missing or empty | INCOMPLETE_ADDRESS | Please provide all required address fields: name, street, city, state, postal code, and country | 400 |

**Traceability:**
- BR: [BR-ADDR-001](../l1/business-rules.md#br-addr-001)

---

## TS-BR-CAT-001 – Category hierarchy depth limit {#ts-br-cat-001}

**Rule:** Categories can be nested up to 3 levels deep

**Implementation Approach:**
On category create with parent, calculate depth and reject if > 3; depth = parent.depth + 1

**Validation Points:**
- POST /categories API
- CategoryService.create()

**Data Requirements:**
| Field | Type | Constraints | Source |
|-------|------|-------------|--------|
| parent_id | UUID | nullable, FK to categories | Category entity |
| depth | integer | >= 1 AND <= 3 | Category entity |

**Error Handling:**
| Condition | Error Code | Message | HTTP Status |
|-----------|------------|---------|-------------|
| calculated depth > 3 | CATEGORY_DEPTH_EXCEEDED | Categories cannot be nested more than 3 levels deep | 400 |

**Traceability:**
- BR: [BR-CAT-001](../l1/business-rules.md#br-cat-001)

---

## TS-BR-AUDIT-001 – Order event audit logging {#ts-br-audit-001}

**Rule:** All order events must be logged for audit purposes

**Implementation Approach:**
OrderEventListener logs all order domain events to audit_log table

**Validation Points:**
- OrderService.create()
- OrderService.updateStatus()
- OrderService.cancel()
- AuditLogService.log()

**Data Requirements:**
| Field | Type | Constraints | Source |
|-------|------|-------------|--------|
| entity_type | varchar(50) | 'Order' | AuditLog entity |
| entity_id | UUID | Order ID | AuditLog entity |
| action | varchar(50) | created|status_changed|cancelled | AuditLog entity |
| before_value | jsonb | nullable, previous state | AuditLog entity |
| after_value | jsonb | new state | AuditLog entity |
| user_id | UUID | actor who made change | AuditLog entity |
| timestamp | timestamp | NOT NULL | AuditLog entity |

**Traceability:**
- BR: [BR-AUDIT-001](../l1/business-rules.md#br-audit-001)

---

## TS-BR-AUDIT-002 – Inventory change audit logging {#ts-br-audit-002}

**Rule:** All inventory changes must be logged with full audit trail

**Implementation Approach:**
InventoryEventListener logs all inventory adjustments to audit_log table

**Validation Points:**
- InventoryService.adjust()
- InventoryService.decrement()
- InventoryService.restore()
- AuditLogService.log()

**Data Requirements:**
| Field | Type | Constraints | Source |
|-------|------|-------------|--------|
| entity_type | varchar(50) | 'Inventory' | AuditLog entity |
| entity_id | UUID | Inventory ID | AuditLog entity |
| action | varchar(50) | adjusted|decremented|restored | AuditLog entity |
| before_value | jsonb | previous stock level | AuditLog entity |
| after_value | jsonb | new stock level | AuditLog entity |
| reason | varchar(200) | adjustment reason | AuditLog entity |
| user_id | UUID | operator who made change | AuditLog entity |

**Traceability:**
- BR: [BR-AUDIT-002](../l1/business-rules.md#br-audit-002)

---

## TS-BR-AUDIT-003 – Price change audit logging {#ts-br-audit-003}

**Rule:** Product price changes must be tracked in audit trail

**Implementation Approach:**
ProductEventListener detects price changes and logs to audit_log table

**Validation Points:**
- ProductService.update()
- AuditLogService.log()

**Data Requirements:**
| Field | Type | Constraints | Source |
|-------|------|-------------|--------|
| entity_type | varchar(50) | 'Product' | AuditLog entity |
| entity_id | UUID | Product ID | AuditLog entity |
| action | varchar(50) | price_changed | AuditLog entity |
| before_value | jsonb | {"price": previous_price} | AuditLog entity |
| after_value | jsonb | {"price": new_price} | AuditLog entity |
| user_id | UUID | modifier who changed price | AuditLog entity |

**Traceability:**
- BR: [BR-AUDIT-003](../l1/business-rules.md#br-audit-003)

---

