# Technical Specifications

Generated: 2026-01-02T21:09:41+01:00

---

## TS-BR-AUTH-001 – Customer registration and email verification for orders {#ts-br-auth-001}

**Rule:** Only registered customers with verified email addresses can place orders

**Implementation Approach:**
Guard clause in checkout service validates Customer exists and emailVerified=true before order creation

**Validation Points:**
- POST /orders API endpoint
- CheckoutService.placeOrder()
- OrderAggregate.create()

**Data Requirements:**
| Field | Type | Constraints | Source |
|-------|------|-------------|--------|
| customerId | UUID | NOT NULL, FK to customers | Session/Auth context |
| emailVerified | boolean | must be true | Customer aggregate |

**Error Handling:**
| Condition | Error Code | Message | HTTP Status |
|-----------|------------|---------|-------------|
| Customer not found or not authenticated | REGISTRATION_REQUIRED | You must be registered to place an order | 401 |
| emailVerified = false | EMAIL_NOT_VERIFIED | Please verify your email before placing an order | 403 |

**Traceability:**
- BR: [BR-AUTH-001](../l1/business-rules.md#br-auth-001)
- Related ACs: [AC-AUTH-001](../l1/acceptance-criteria.md#ac-auth-001)

---

## TS-BR-STOCK-001 – Stock availability validation for cart additions {#ts-br-stock-001}

**Rule:** Products must have available inventory to be added to cart

**Implementation Approach:**
Check Inventory.availableStock > 0 on add-to-cart; allow cart items with stock warnings but block checkout if unavailable

**Validation Points:**
- POST /cart/items API
- CartService.addItem()
- CheckoutService.validateCart()

**Data Requirements:**
| Field | Type | Constraints | Source |
|-------|------|-------------|--------|
| availableStock | integer | >= 0 | Inventory aggregate |
| productId | UUID | NOT NULL | Request payload |
| variantId | UUID | NULL if no variants | Request payload |

**Error Handling:**
| Condition | Error Code | Message | HTTP Status |
|-----------|------------|---------|-------------|
| availableStock = 0 | OUT_OF_STOCK | This product is currently out of stock | 409 |

**Traceability:**
- BR: [BR-STOCK-001](../l1/business-rules.md#br-stock-001)
- Related ACs: [AC-STOCK-001](../l1/acceptance-criteria.md#ac-stock-001)

---

## TS-BR-SHIP-001 – Free shipping threshold calculation {#ts-br-ship-001}

**Rule:** Orders with subtotal of $50 or more qualify for free shipping

**Implementation Approach:**
ShippingCalculator checks order.subtotal >= 50.00; if true, set shippingCost = 0; else shippingCost = 5.99

**Validation Points:**
- ShippingService.calculateShipping()
- CheckoutService.calculateTotals()
- Cart summary display

**Data Requirements:**
| Field | Type | Constraints | Source |
|-------|------|-------------|--------|
| subtotal | decimal(10,2) | >= 0 | Calculated from cart/order items after discounts |
| shippingCost | decimal(10,2) | 0.00 or 5.99 | Calculated field |

**Traceability:**
- BR: [BR-SHIP-001](../l1/business-rules.md#br-ship-001)
- Related ACs: [AC-SHIP-001](../l1/acceptance-criteria.md#ac-ship-001)

---

## TS-BR-CANCEL-001 – Order cancellation status validation {#ts-br-cancel-001}

**Rule:** Orders can only be cancelled while in pending or confirmed status

**Implementation Approach:**
Guard clause in cancel operation checks order.status IN ('pending', 'confirmed') before allowing transition to 'cancelled'

**Validation Points:**
- POST /orders/{id}/cancel API
- OrderService.cancelOrder()
- Order.cancel()

**Data Requirements:**
| Field | Type | Constraints | Source |
|-------|------|-------------|--------|
| status | OrderStatus enum | pending|confirmed|shipped|delivered|cancelled | Order aggregate |

**Error Handling:**
| Condition | Error Code | Message | HTTP Status |
|-----------|------------|---------|-------------|
| status NOT IN ('pending', 'confirmed') | CANCELLATION_NOT_ALLOWED | This order cannot be cancelled as it has already been shipped | 409 |

**Traceability:**
- BR: [BR-CANCEL-001](../l1/business-rules.md#br-cancel-001)
- Related ACs: [AC-CANCEL-001](../l1/acceptance-criteria.md#ac-cancel-001)

---

## TS-BR-PRICE-001 – Minimum product price enforcement {#ts-br-price-001}

**Rule:** Product price must be at least $0.01

**Implementation Approach:**
Validate price >= 0.01 with exactly 2 decimal places on product create/update

**Validation Points:**
- POST /products API
- PUT /products/{id} API
- ProductService.createProduct()
- ProductService.updateProduct()
- DB CHECK constraint

**Data Requirements:**
| Field | Type | Constraints | Source |
|-------|------|-------------|--------|
| price | decimal(10,2) | >= 0.01 | Request payload |

**Error Handling:**
| Condition | Error Code | Message | HTTP Status |
|-----------|------------|---------|-------------|
| price < 0.01 | INVALID_PRICE | Product price must be at least $0.01 | 400 |
| price has more than 2 decimal places | INVALID_PRICE | Price must have exactly 2 decimal places | 400 |

**Traceability:**
- BR: [BR-PRICE-001](../l1/business-rules.md#br-price-001)
- Related ACs: [AC-PRICE-001](../l1/acceptance-criteria.md#ac-price-001)

---

## TS-BR-PRICE-002 – Order price snapshot on placement {#ts-br-price-002}

**Rule:** Order items must capture prices at the time of order placement

**Implementation Approach:**
Copy current Product.price to OrderItem.unitPrice at order creation; OrderItem.price is immutable after creation

**Validation Points:**
- OrderService.placeOrder()
- OrderItem.create()
- DB: no UPDATE on orderItems.unitPrice

**Data Requirements:**
| Field | Type | Constraints | Source |
|-------|------|-------------|--------|
| unitPrice | decimal(10,2) | immutable after creation | Product.price at order time |

**Traceability:**
- BR: [BR-PRICE-002](../l1/business-rules.md#br-price-002)
- Related ACs: [AC-PRICE-002](../l1/acceptance-criteria.md#ac-price-002)

---

## TS-BR-CART-001 – Cart real-time price calculation {#ts-br-cart-001}

**Rule:** Cart items always reflect current product prices

**Implementation Approach:**
Cart.total calculated by joining to current Product.price on every cart retrieval; show price change notifications if price differs from when added

**Validation Points:**
- GET /cart API
- CartService.getCart()
- Cart view rendering

**Data Requirements:**
| Field | Type | Constraints | Source |
|-------|------|-------------|--------|
| priceAtAdd | decimal(10,2) | optional, for change detection | Snapshot at add time |
| currentPrice | decimal(10,2) | fetched from Product | Product aggregate |

**Traceability:**
- BR: [BR-CART-001](../l1/business-rules.md#br-cart-001)
- Related ACs: [AC-CART-001](../l1/acceptance-criteria.md#ac-cart-001)

---

## TS-BR-INV-001 – Non-negative inventory constraint {#ts-br-inv-001}

**Rule:** Inventory stock level must never be negative

**Implementation Approach:**
DB CHECK constraint on availableStock >= 0; all inventory adjustments use atomic decrement with stock check

**Validation Points:**
- InventoryService.adjustStock()
- OrderService.placeOrder()
- DB CHECK(availableStock >= 0)

**Data Requirements:**
| Field | Type | Constraints | Source |
|-------|------|-------------|--------|
| availableStock | integer | >= 0, NOT NULL | Inventory aggregate |

**Error Handling:**
| Condition | Error Code | Message | HTTP Status |
|-----------|------------|---------|-------------|
| Adjustment would result in negative stock | INSUFFICIENT_STOCK | Insufficient stock to complete this operation | 409 |

**Traceability:**
- BR: [BR-INV-001](../l1/business-rules.md#br-inv-001)
- Related ACs: [AC-INV-001](../l1/acceptance-criteria.md#ac-inv-001)

---

## TS-BR-INV-002 – Backorder prevention {#ts-br-inv-002}

**Rule:** Orders cannot be placed for quantities exceeding available stock

**Implementation Approach:**
Atomic check-and-decrement at checkout: SELECT FOR UPDATE inventory, verify each item.quantity <= availableStock, then decrement

**Validation Points:**
- CheckoutService.placeOrder()
- InventoryService.reserveStock()
- DB transaction with row locks

**Data Requirements:**
| Field | Type | Constraints | Source |
|-------|------|-------------|--------|
| requestedQuantity | integer | > 0 | OrderItem |
| availableStock | integer | >= 0 | Inventory aggregate |

**Error Handling:**
| Condition | Error Code | Message | HTTP Status |
|-----------|------------|---------|-------------|
| requestedQuantity > availableStock | BACKORDER_NOT_SUPPORTED | Cannot order more items than currently in stock | 409 |

**Traceability:**
- BR: [BR-INV-002](../l1/business-rules.md#br-inv-002)
- Related ACs: [AC-INV-002](../l1/acceptance-criteria.md#ac-inv-002)

---

## TS-BR-INV-003 – Inventory restoration on order cancellation {#ts-br-inv-003}

**Rule:** Cancelled orders must automatically restore inventory

**Implementation Approach:**
Order cancellation triggers InventoryService.restoreStock() for each OrderItem in same transaction

**Validation Points:**
- OrderService.cancelOrder()
- InventoryService.restoreStock()
- DB transaction atomicity

**Data Requirements:**
| Field | Type | Constraints | Source |
|-------|------|-------------|--------|
| orderItems | OrderItem[] | quantity and productId/variantId | Order aggregate |

**Traceability:**
- BR: [BR-INV-003](../l1/business-rules.md#br-inv-003)
- Related ACs: [AC-INV-003](../l1/acceptance-criteria.md#ac-inv-003)

---

## TS-BR-ORDER-001 – Order state machine transitions {#ts-br-order-001}

**Rule:** Orders must follow defined state machine transitions

**Implementation Approach:**
State machine in Order aggregate with allowed transitions map; validate transition before applying

**Validation Points:**
- Order.transitionTo()
- OrderService.updateStatus()
- PUT /orders/{id}/status API

**Data Requirements:**
| Field | Type | Constraints | Source |
|-------|------|-------------|--------|
| status | OrderStatus enum | pending|confirmed|shipped|delivered|cancelled | Order aggregate |
| allowedTransitions | Map<Status, Status[]> | pending→[confirmed,cancelled], confirmed→[shipped,cancelled], shipped→[delivered] | Domain logic |

**Error Handling:**
| Condition | Error Code | Message | HTTP Status |
|-----------|------------|---------|-------------|
| Requested transition not in allowed transitions | INVALID_STATUS_TRANSITION | Cannot transition order from {current} to {requested} | 409 |

**Traceability:**
- BR: [BR-ORDER-001](../l1/business-rules.md#br-order-001)
- Related ACs: [AC-ORDER-001](../l1/acceptance-criteria.md#ac-order-001)

---

## TS-BR-ORDER-002 – Order immutability after placement {#ts-br-order-002}

**Rule:** Orders cannot be modified after placement

**Implementation Approach:**
No PUT/PATCH endpoints for order items, quantities, or shipping address; only status transitions allowed

**Validation Points:**
- API layer: no modification endpoints
- OrderService: no edit methods
- Order aggregate: immutable fields

**Data Requirements:**
| Field | Type | Constraints | Source |
|-------|------|-------------|--------|
| orderItems | OrderItem[] | immutable after creation | Order aggregate |
| shippingAddress | ShippingAddress | immutable after creation | Order aggregate |

**Error Handling:**
| Condition | Error Code | Message | HTTP Status |
|-----------|------------|---------|-------------|
| Attempt to modify placed order | ORDER_MODIFICATION_NOT_ALLOWED | Orders cannot be modified after placement. Please cancel and place a new order. | 409 |

**Traceability:**
- BR: [BR-ORDER-002](../l1/business-rules.md#br-order-002)
- Related ACs: [AC-ORDER-002](../l1/acceptance-criteria.md#ac-order-002)

---

## TS-BR-ORDER-003 – Single shipment per order constraint {#ts-br-order-003}

**Rule:** Each order ships as a single unit

**Implementation Approach:**
Order has one-to-one relationship with Shipment; shipped status applies to entire order atomically

**Validation Points:**
- Order aggregate design
- ShipmentService.createShipment()
- DB: unique constraint on order_id in shipments

**Data Requirements:**
| Field | Type | Constraints | Source |
|-------|------|-------------|--------|
| shipmentId | UUID | one per order | Shipment aggregate |

**Traceability:**
- BR: [BR-ORDER-003](../l1/business-rules.md#br-order-003)
- Related ACs: [AC-ORDER-003](../l1/acceptance-criteria.md#ac-order-003)

---

## TS-BR-CUST-001 – Unique customer email constraint {#ts-br-cust-001}

**Rule:** Email addresses must be unique across all customer accounts

**Implementation Approach:**
DB unique index on customers.email; check uniqueness before insert

**Validation Points:**
- POST /customers (registration) API
- CustomerService.register()
- DB UNIQUE constraint on email

**Data Requirements:**
| Field | Type | Constraints | Source |
|-------|------|-------------|--------|
| email | varchar(255) | NOT NULL, UNIQUE, valid email format | Request payload |

**Error Handling:**
| Condition | Error Code | Message | HTTP Status |
|-----------|------------|---------|-------------|
| Email already exists in database | EMAIL_ALREADY_REGISTERED | An account with this email already exists | 409 |

**Traceability:**
- BR: [BR-CUST-001](../l1/business-rules.md#br-cust-001)
- Related ACs: [AC-CUST-001](../l1/acceptance-criteria.md#ac-cust-001)

---

## TS-BR-CUST-002 – Password strength requirements {#ts-br-cust-002}

**Rule:** Customer passwords must meet minimum security requirements

**Implementation Approach:**
Validate password with regex: ^(?=.*\d).{8,}$ on registration and password change

**Validation Points:**
- POST /customers (registration) API
- PUT /customers/{id}/password API
- CustomerService.setPassword()

**Data Requirements:**
| Field | Type | Constraints | Source |
|-------|------|-------------|--------|
| password | string | min 8 chars, at least one digit | Request payload |

**Error Handling:**
| Condition | Error Code | Message | HTTP Status |
|-----------|------------|---------|-------------|
| Password does not meet requirements | WEAK_PASSWORD | Password must be at least 8 characters and contain at least one number | 400 |

**Traceability:**
- BR: [BR-CUST-002](../l1/business-rules.md#br-cust-002)
- Related ACs: [AC-CUST-002](../l1/acceptance-criteria.md#ac-cust-002)

---

## TS-BR-CUST-003 – Email verification requirement for checkout {#ts-br-cust-003}

**Rule:** Customers must verify their email before placing orders

**Implementation Approach:**
Check customer.emailVerified = true in checkout flow before allowing order placement

**Validation Points:**
- POST /orders API
- CheckoutService.placeOrder()
- Pre-checkout validation

**Data Requirements:**
| Field | Type | Constraints | Source |
|-------|------|-------------|--------|
| emailVerified | boolean | must be true for checkout | Customer aggregate |

**Error Handling:**
| Condition | Error Code | Message | HTTP Status |
|-----------|------------|---------|-------------|
| emailVerified = false | EMAIL_NOT_VERIFIED | Please verify your email address before placing an order | 403 |

**Traceability:**
- BR: [BR-CUST-003](../l1/business-rules.md#br-cust-003)
- Related ACs: [AC-CUST-003](../l1/acceptance-criteria.md#ac-cust-003)

---

## TS-BR-PROD-001 – Product soft delete for order history preservation {#ts-br-prod-001}

**Rule:** Products with order history must be soft-deleted

**Implementation Approach:**
Delete operation checks for existing OrderItems; if found, set deletedAt timestamp instead of hard delete

**Validation Points:**
- DELETE /products/{id} API
- ProductService.deleteProduct()
- Query filters exclude deletedAt IS NOT NULL

**Data Requirements:**
| Field | Type | Constraints | Source |
|-------|------|-------------|--------|
| deletedAt | timestamp | NULL if active | Product aggregate |
| hasOrderHistory | boolean | derived from OrderItems | OrderItem table |

**Traceability:**
- BR: [BR-PROD-001](../l1/business-rules.md#br-prod-001)
- Related ACs: [AC-PROD-001](../l1/acceptance-criteria.md#ac-prod-001)

---

## TS-BR-PROD-002 – Product name length validation {#ts-br-prod-002}

**Rule:** Product names must be between 2 and 200 characters

**Implementation Approach:**
Validate name.length >= 2 AND <= 200 on product create/update

**Validation Points:**
- POST /products API
- PUT /products/{id} API
- ProductService.createProduct()
- ProductService.updateProduct()

**Data Requirements:**
| Field | Type | Constraints | Source |
|-------|------|-------------|--------|
| name | varchar(200) | 2 <= length <= 200, NOT NULL | Request payload |

**Error Handling:**
| Condition | Error Code | Message | HTTP Status |
|-----------|------------|---------|-------------|
| name.length < 2 OR name.length > 200 | INVALID_PRODUCT_NAME | Product name must be between 2 and 200 characters | 400 |

**Traceability:**
- BR: [BR-PROD-002](../l1/business-rules.md#br-prod-002)
- Related ACs: [AC-PROD-002](../l1/acceptance-criteria.md#ac-prod-002)

---

## TS-BR-PROD-003 – Product description length constraint {#ts-br-prod-003}

**Rule:** Product descriptions must not exceed 5000 characters

**Implementation Approach:**
Validate description.length <= 5000 on product create/update

**Validation Points:**
- POST /products API
- PUT /products/{id} API
- ProductService.createProduct()
- ProductService.updateProduct()

**Data Requirements:**
| Field | Type | Constraints | Source |
|-------|------|-------------|--------|
| description | text | length <= 5000, nullable | Request payload |

**Error Handling:**
| Condition | Error Code | Message | HTTP Status |
|-----------|------------|---------|-------------|
| description.length > 5000 | DESCRIPTION_TOO_LONG | Product description cannot exceed 5000 characters | 400 |

**Traceability:**
- BR: [BR-PROD-003](../l1/business-rules.md#br-prod-003)
- Related ACs: [AC-PROD-003](../l1/acceptance-criteria.md#ac-prod-003)

---

## TS-BR-PROD-004 – Product image limit and primary image enforcement {#ts-br-prod-004}

**Rule:** Products can have up to 10 images with exactly one primary

**Implementation Approach:**
Validate images.count <= 10 and exactly one isPrimary=true on image add/update

**Validation Points:**
- POST /products/{id}/images API
- PUT /products/{id}/images/{imageId} API
- ProductImageService.addImage()
- ProductImageService.setPrimary()

**Data Requirements:**
| Field | Type | Constraints | Source |
|-------|------|-------------|--------|
| images | ProductImage[] | count <= 10 | Product aggregate |
| isPrimary | boolean | exactly one true per product | ProductImage |

**Error Handling:**
| Condition | Error Code | Message | HTTP Status |
|-----------|------------|---------|-------------|
| images.count > 10 | IMAGE_LIMIT_EXCEEDED | Products cannot have more than 10 images | 400 |
| No primary image set | PRIMARY_IMAGE_REQUIRED | Product must have exactly one primary image | 400 |

**Traceability:**
- BR: [BR-PROD-004](../l1/business-rules.md#br-prod-004)
- Related ACs: [AC-PROD-004](../l1/acceptance-criteria.md#ac-prod-004)

---

## TS-BR-VAR-001 – Variant SKU uniqueness constraint {#ts-br-var-001}

**Rule:** Each product variant must have a unique SKU

**Implementation Approach:**
DB unique index on product_variants.sku; validate uniqueness before insert

**Validation Points:**
- POST /products/{id}/variants API
- VariantService.createVariant()
- DB UNIQUE constraint on sku

**Data Requirements:**
| Field | Type | Constraints | Source |
|-------|------|-------------|--------|
| sku | varchar(100) | NOT NULL, UNIQUE | Request payload |

**Error Handling:**
| Condition | Error Code | Message | HTTP Status |
|-----------|------------|---------|-------------|
| SKU already exists | DUPLICATE_SKU | A product variant with this SKU already exists | 409 |

**Traceability:**
- BR: [BR-VAR-001](../l1/business-rules.md#br-var-001)
- Related ACs: [AC-VAR-001](../l1/acceptance-criteria.md#ac-var-001)

---

## TS-BR-VAR-002 – Variant selection required for products with variants {#ts-br-var-002}

**Rule:** Products with variants require variant selection before add to cart

**Implementation Approach:**
Check if product.hasVariants; if true, require variantId in add-to-cart request

**Validation Points:**
- POST /cart/items API
- CartService.addItem()
- UI: disable add button until variant selected

**Data Requirements:**
| Field | Type | Constraints | Source |
|-------|------|-------------|--------|
| hasVariants | boolean | derived from variants.count > 0 | Product aggregate |
| variantId | UUID | required if hasVariants | Request payload |

**Error Handling:**
| Condition | Error Code | Message | HTTP Status |
|-----------|------------|---------|-------------|
| Product has variants but variantId not provided | VARIANT_REQUIRED | Please select a product variant before adding to cart | 400 |

**Traceability:**
- BR: [BR-VAR-002](../l1/business-rules.md#br-var-002)
- Related ACs: [AC-VAR-002](../l1/acceptance-criteria.md#ac-var-002)

---

## TS-BR-CART-002 – Cart quantity stock limit enforcement {#ts-br-cart-002}

**Rule:** Cart item quantity cannot exceed available inventory

**Implementation Approach:**
Cap quantity at availableStock on add/update; show notification if requested quantity exceeds stock

**Validation Points:**
- POST /cart/items API
- PUT /cart/items/{id} API
- CartService.addItem()
- CartService.updateQuantity()

**Data Requirements:**
| Field | Type | Constraints | Source |
|-------|------|-------------|--------|
| quantity | integer | > 0, <= availableStock | Request payload |
| availableStock | integer | >= 0 | Inventory aggregate |

**Error Handling:**
| Condition | Error Code | Message | HTTP Status |
|-----------|------------|---------|-------------|
| quantity > availableStock | QUANTITY_EXCEEDS_STOCK | Only {available} items available. Quantity adjusted. | 409 |

**Traceability:**
- BR: [BR-CART-002](../l1/business-rules.md#br-cart-002)
- Related ACs: [AC-CART-002](../l1/acceptance-criteria.md#ac-cart-002)

---

## TS-BR-CART-003 – Cart expiration policy by user type {#ts-br-cart-003}

**Rule:** Carts expire after period of inactivity based on user type

**Implementation Approach:**
Background job runs daily to delete carts where lastModified + expirationPeriod < now(); update lastModified on cart operations

**Validation Points:**
- CartCleanupJob (scheduled)
- CartService.onCartModification()
- GET /cart API (check expiration)

**Data Requirements:**
| Field | Type | Constraints | Source |
|-------|------|-------------|--------|
| lastModified | timestamp | NOT NULL, updated on changes | Cart aggregate |
| isGuest | boolean | determines expiration period | Cart aggregate |

**Error Handling:**
| Condition | Error Code | Message | HTTP Status |
|-----------|------------|---------|-------------|
| Cart has expired | CART_EXPIRED | Your cart has expired. Please add items again. | 410 |

**Traceability:**
- BR: [BR-CART-003](../l1/business-rules.md#br-cart-003)
- Related ACs: [AC-CART-003](../l1/acceptance-criteria.md#ac-cart-003)

---

## TS-BR-PAY-001 – Payment authorization before order creation {#ts-br-pay-001}

**Rule:** Payment must be fully authorized before order is created

**Implementation Approach:**
Checkout flow: 1) validate cart, 2) authorize payment, 3) if authorized, create order atomically with inventory decrement

**Validation Points:**
- CheckoutService.placeOrder()
- PaymentService.authorize()
- Transaction boundary

**Data Requirements:**
| Field | Type | Constraints | Source |
|-------|------|-------------|--------|
| paymentAuthorizationId | string | from payment gateway | Payment gateway response |
| paymentStatus | PaymentStatus enum | must be 'authorized' | Payment service |

**Error Handling:**
| Condition | Error Code | Message | HTTP Status |
|-----------|------------|---------|-------------|
| Payment authorization failed | PAYMENT_REQUIRED | Payment could not be authorized. Please try again. | 402 |

**Traceability:**
- BR: [BR-PAY-001](../l1/business-rules.md#br-pay-001)
- Related ACs: [AC-PAY-001](../l1/acceptance-criteria.md#ac-pay-001)

---

## TS-BR-PAY-002 – Automatic refund on order cancellation {#ts-br-pay-002}

**Rule:** Cancelled orders must receive automatic refund

**Implementation Approach:**
Order cancellation triggers PaymentService.refund(order.paymentId) in same transaction; store refund status

**Validation Points:**
- OrderService.cancelOrder()
- PaymentService.refund()
- Refund status tracking

**Data Requirements:**
| Field | Type | Constraints | Source |
|-------|------|-------------|--------|
| paymentId | string | original payment reference | Order aggregate |
| refundStatus | RefundStatus enum | pending|completed|failed | Payment service |

**Traceability:**
- BR: [BR-PAY-002](../l1/business-rules.md#br-pay-002)
- Related ACs: [AC-PAY-002](../l1/acceptance-criteria.md#ac-pay-002)

---

## TS-BR-SHIP-002 – Domestic shipping only restriction {#ts-br-ship-002}

**Rule:** Shipping is only available within a single country

**Implementation Approach:**
Validate address.country equals config.domesticCountry; reject international addresses

**Validation Points:**
- POST /orders API
- AddressValidator.validate()
- Checkout address form

**Data Requirements:**
| Field | Type | Constraints | Source |
|-------|------|-------------|--------|
| country | varchar(2) | ISO 3166-1 alpha-2, must match domestic | Request payload |
| domesticCountry | varchar(2) | configured value | Application config |

**Error Handling:**
| Condition | Error Code | Message | HTTP Status |
|-----------|------------|---------|-------------|
| country != domesticCountry | INTERNATIONAL_SHIPPING_NOT_AVAILABLE | We currently only ship within {country} | 400 |

**Traceability:**
- BR: [BR-SHIP-002](../l1/business-rules.md#br-ship-002)
- Related ACs: [AC-SHIP-002](../l1/acceptance-criteria.md#ac-ship-002)

---

## TS-BR-ADDR-001 – Required shipping address fields validation {#ts-br-addr-001}

**Rule:** Shipping addresses must include all required fields

**Implementation Approach:**
Validate all required fields are present and non-empty in address payload

**Validation Points:**
- POST /orders API
- PUT /customers/{id}/addresses API
- AddressValidator.validate()
- Checkout form validation

**Data Requirements:**
| Field | Type | Constraints | Source |
|-------|------|-------------|--------|
| street | varchar(255) | NOT NULL, non-empty | Request payload |
| city | varchar(100) | NOT NULL, non-empty | Request payload |
| state | varchar(100) | NOT NULL, non-empty | Request payload |
| postalCode | varchar(20) | NOT NULL, non-empty | Request payload |
| country | varchar(2) | NOT NULL, ISO 3166-1 | Request payload |
| recipientName | varchar(200) | NOT NULL, non-empty | Request payload |

**Error Handling:**
| Condition | Error Code | Message | HTTP Status |
|-----------|------------|---------|-------------|
| Any required field missing or empty | INCOMPLETE_ADDRESS | Please provide all required address fields | 400 |

**Traceability:**
- BR: [BR-ADDR-001](../l1/business-rules.md#br-addr-001)
- Related ACs: [AC-ADDR-001](../l1/acceptance-criteria.md#ac-addr-001)

---

## TS-BR-CAT-001 – Category hierarchy depth limit {#ts-br-cat-001}

**Rule:** Categories can be nested up to 3 levels deep

**Implementation Approach:**
Calculate depth by traversing parent chain; reject if depth > 3 on category create/update with parent

**Validation Points:**
- POST /categories API
- PUT /categories/{id} API
- CategoryService.createCategory()
- CategoryService.updateCategory()

**Data Requirements:**
| Field | Type | Constraints | Source |
|-------|------|-------------|--------|
| parentId | UUID | nullable, FK to categories | Request payload |
| depth | integer | 1 <= depth <= 3 | Calculated from hierarchy |

**Error Handling:**
| Condition | Error Code | Message | HTTP Status |
|-----------|------------|---------|-------------|
| Resulting depth > 3 | CATEGORY_DEPTH_EXCEEDED | Categories cannot be nested more than 3 levels deep | 400 |

**Traceability:**
- BR: [BR-CAT-001](../l1/business-rules.md#br-cat-001)
- Related ACs: [AC-CAT-001](../l1/acceptance-criteria.md#ac-cat-001)

---

## TS-BR-AUDIT-001 – Order event audit logging {#ts-br-audit-001}

**Rule:** All order events must be logged for audit purposes

**Implementation Approach:**
Emit OrderEvent to audit log on create, status change, cancel; include timestamp, userId, previousStatus, newStatus

**Validation Points:**
- OrderService.placeOrder()
- OrderService.updateStatus()
- OrderService.cancelOrder()
- AuditLogService.logOrderEvent()

**Data Requirements:**
| Field | Type | Constraints | Source |
|-------|------|-------------|--------|
| eventType | OrderEventType enum | created|status_changed|cancelled | Domain event |
| timestamp | timestamp | NOT NULL, UTC | System clock |
| userId | UUID | NOT NULL | Auth context |
| previousValue | json | nullable | Before state |
| newValue | json | nullable | After state |

**Traceability:**
- BR: [BR-AUDIT-001](../l1/business-rules.md#br-audit-001)
- Related ACs: [AC-AUDIT-001](../l1/acceptance-criteria.md#ac-audit-001)

---

## TS-BR-AUDIT-002 – Inventory change audit logging {#ts-br-audit-002}

**Rule:** All inventory changes must be logged with full audit trail

**Implementation Approach:**
Emit InventoryEvent on every stock adjustment; include productId, previousStock, newStock, reason, operatorId

**Validation Points:**
- InventoryService.adjustStock()
- InventoryService.reserveStock()
- InventoryService.restoreStock()
- AuditLogService.logInventoryEvent()

**Data Requirements:**
| Field | Type | Constraints | Source |
|-------|------|-------------|--------|
| productId | UUID | NOT NULL | Inventory aggregate |
| previousStock | integer | >= 0 | Before state |
| newStock | integer | >= 0 | After state |
| reason | AdjustmentReason enum | order_placed|order_cancelled|manual_adjustment|... | Operation context |
| operatorId | UUID | NOT NULL | Auth context |

**Traceability:**
- BR: [BR-AUDIT-002](../l1/business-rules.md#br-audit-002)
- Related ACs: [AC-AUDIT-002](../l1/acceptance-criteria.md#ac-audit-002)

---

## TS-BR-AUDIT-003 – Price change audit tracking {#ts-br-audit-003}

**Rule:** Product price changes must be tracked in audit trail

**Implementation Approach:**
Compare price before/after on product update; if changed, emit PriceChangeEvent to audit log

**Validation Points:**
- ProductService.updateProduct()
- AuditLogService.logPriceChange()

**Data Requirements:**
| Field | Type | Constraints | Source |
|-------|------|-------------|--------|
| productId | UUID | NOT NULL | Product aggregate |
| previousPrice | decimal(10,2) | >= 0.01 | Before state |
| newPrice | decimal(10,2) | >= 0.01 | After state |
| modifierId | UUID | NOT NULL | Auth context |
| timestamp | timestamp | NOT NULL, UTC | System clock |

**Traceability:**
- BR: [BR-AUDIT-003](../l1/business-rules.md#br-audit-003)
- Related ACs: [AC-AUDIT-003](../l1/acceptance-criteria.md#ac-audit-003)

---

