# Implementation Skeletons

Generated: 2025-12-29T12:45:56+01:00

---

## ProductService (service)

**Dependencies:** [InventoryService CategoryService CartService AuditService]

**Related Specs:** [TS-BR-PRICE-001 TS-BR-PROD-001 TS-BR-PROD-002 TS-BR-PROD-003 TS-BR-AUDIT-003]

### Functions

#### `listProducts`

```typescript
listProducts(filters: ProductFilters, pagination: Pagination): PaginatedResult<Product>
```

List products with filtering, sorting, and pagination

**Implementation Steps:**
1. Build query with category, price range, and availability filters
2. Apply sorting (name, price_asc, price_desc, newest)
3. Execute paginated query excluding soft-deleted products
4. Join with inventory for stock status
5. Return products with pagination metadata

**Error Cases:** [INVALID_PRICE_RANGE]

#### `getProduct`

```typescript
getProduct(productId: string): Product
```

Get product details with stock status and variants

**Implementation Steps:**
1. Query product by ID excluding soft-deleted
2. Load product variants with inventory
3. Calculate stock availability status
4. Return product with variants and stock info

**Error Cases:** [PRODUCT_NOT_FOUND]

#### `createProduct`

```typescript
createProduct(data: CreateProductRequest, adminId: string): Product
```

Create new product with draft status

**Implementation Steps:**
1. Validate name length (2-200 chars)
2. Validate price >= $0.01
3. Validate category exists
4. Check for duplicate name (warning only)
5. Create product with draft status
6. Log creator and timestamp
7. Return created product

**Error Cases:** [INVALID_NAME NAME_TOO_LONG INVALID_PRICE CATEGORY_REQUIRED]

#### `updateProduct`

```typescript
updateProduct(productId: string, data: UpdateProductRequest, adminId: string): Product
```

Update product attributes with audit trail

**Implementation Steps:**
1. Load existing product
2. Validate updated fields
3. If price changed, log to price_history
4. Update product with last_modified_by/at
5. Return updated product

**Error Cases:** [PRODUCT_NOT_FOUND INVALID_NAME INVALID_PRICE]

#### `deleteProduct`

```typescript
deleteProduct(productId: string, adminId: string): DeletionSummary
```

Soft delete product and remove from carts

**Implementation Steps:**
1. Check if product has order history
2. Set deleted_at timestamp
3. Remove product from all active carts
4. Return impact summary (affected carts count)

**Error Cases:** [PRODUCT_NOT_FOUND]

---

## CartService (service)

**Dependencies:** [ProductService InventoryService]

**Related Specs:** [TS-BR-STOCK-001 TS-BR-CART-001 TS-BR-CART-002 TS-BR-PRICE-002 TS-BR-VAR-002 TS-BR-SHIP-001]

### Functions

#### `getCart`

```typescript
getCart(customerId: string | null, sessionId: string): Cart
```

Get cart with current prices and change notifications

**Implementation Steps:**
1. Load cart by customer ID or session ID
2. Join with products for current prices
3. Compare added_at_price with current_price
4. Flag items with price changes
5. Check stock status for each item
6. Return cart with notifications

#### `addItem`

```typescript
addItem(cartId: string, productId: string, variantId: string | null, quantity: number): CartItem
```

Add item to cart with stock and variant validation

**Implementation Steps:**
1. Validate product exists and is active
2. If product has variants, require variant_id
3. Check available stock
4. If quantity > stock, limit to stock and set flag
5. Check for existing cart item (same product/variant)
6. If exists, increment quantity (respecting stock limit)
7. If not exists, create new cart item
8. Store added_at_price
9. Update cart last_activity timestamp
10. Return cart item with quantity_limited flag

**Error Cases:** [PRODUCT_NOT_FOUND OUT_OF_STOCK VARIANT_REQUIRED INVALID_QUANTITY]

#### `updateItem`

```typescript
updateItem(cartItemId: string, quantity: number): CartItem
```

Update cart item quantity with validation

**Implementation Steps:**
1. Validate quantity >= 0
2. If quantity == 0, remove item
3. Check available stock
4. If quantity > stock, limit to stock
5. Update cart item quantity
6. Update cart last_activity
7. Return updated item

**Error Cases:** [CART_ITEM_NOT_FOUND INVALID_QUANTITY QUANTITY_EXCEEDS_STOCK]

#### `removeItem`

```typescript
removeItem(cartItemId: string): UndoToken
```

Remove item from cart with undo capability

**Implementation Steps:**
1. Store item details for undo
2. Delete cart item
3. Generate undo token with TTL
4. Update cart last_activity
5. Return undo token

**Error Cases:** [CART_ITEM_NOT_FOUND]

#### `calculateTotals`

```typescript
calculateTotals(cartId: string): CartTotals
```

Calculate cart totals with shipping and tax

**Implementation Steps:**
1. Sum item prices (current_price * quantity)
2. Apply any discounts
3. Calculate subtotal (after discounts)
4. If subtotal >= $50, shipping = $0, else shipping = $5.99
5. Calculate tax based on shipping address
6. Return totals breakdown

#### `mergeOnLogin`

```typescript
mergeOnLogin(guestSessionId: string, customerId: string): Cart
```

Merge guest cart with customer cart on login

**Implementation Steps:**
1. Load guest cart by session ID
2. Load customer cart by customer ID
3. For each guest item, check if exists in customer cart
4. If duplicate, combine quantities (respecting stock limit)
5. If new item, add to customer cart
6. Delete guest cart
7. Return merged customer cart

#### `cleanupExpiredCarts`

```typescript
cleanupExpiredCarts(): CleanupResult
```

Background job to remove expired carts

**Implementation Steps:**
1. Query guest carts with last_activity > 7 days ago
2. Query registered carts with last_activity > 30 days ago
3. Delete expired carts
4. Return count of deleted carts

---

## OrderService (service)

**Dependencies:** [CartService InventoryService PaymentService CustomerService EmailService AuditService]

**Related Specs:** [TS-BR-AUTH-001 TS-BR-STOCK-002 TS-BR-CANCEL-001 TS-BR-ORDER-001 TS-BR-ORDER-002 TS-BR-ORDER-003 TS-BR-ORDER-004 TS-BR-PAY-002 TS-BR-INV-001 TS-BR-INV-002 TS-BR-AUDIT-001]

### Functions

#### `initiateCheckout`

```typescript
initiateCheckout(customerId: string, shippingAddressId: string, paymentMethodId: string): CheckoutSession
```

Validate and initiate checkout process

**Implementation Steps:**
1. Validate customer is registered
2. Validate email is verified
3. Load cart items
4. Validate all items in stock (SELECT FOR UPDATE)
5. Validate shipping address exists and is domestic
6. Validate payment method exists
7. Create checkout session
8. Return session with totals

**Error Cases:** [REGISTRATION_REQUIRED EMAIL_NOT_VERIFIED INSUFFICIENT_STOCK ADDRESS_VALIDATION_ERROR]

#### `createOrder`

```typescript
createOrder(checkoutSessionId: string): Order
```

Create order with payment authorization and inventory decrement

**Implementation Steps:**
1. Load checkout session
2. Re-validate stock with row locking
3. Create payment intent
4. Authorize payment
5. If payment fails, throw error
6. Generate order number (ORD-YYYY-NNNNN)
7. Create order with status 'pending'
8. Copy cart items to order_items with current prices
9. Decrement inventory for all items
10. Clear cart
11. Queue confirmation email
12. Log order creation
13. Return order

**Error Cases:** [INSUFFICIENT_STOCK PAYMENT_FAILED ADDRESS_VALIDATION_ERROR]

#### `cancelOrder`

```typescript
cancelOrder(orderId: string, customerId: string, reason: string | null): Order
```

Cancel order and restore inventory

**Implementation Steps:**
1. Load order and validate ownership
2. Validate status is 'pending' or 'confirmed'
3. Update status to 'cancelled'
4. Restore inventory for all order items
5. Process refund to original payment method
6. Queue cancellation email
7. Log status change with reason
8. Return updated order

**Error Cases:** [ORDER_NOT_FOUND CANCELLATION_NOT_ALLOWED]

#### `getCustomerOrders`

```typescript
getCustomerOrders(customerId: string, pagination: Pagination): PaginatedResult<Order>
```

Get paginated order history for customer

**Implementation Steps:**
1. Query orders by customer ID
2. Include order items, status, tracking
3. Sort by created_at descending
4. Apply pagination
5. Return orders with metadata

#### `updateStatus`

```typescript
updateStatus(orderId: string, newStatus: string, adminId: string, trackingNumber: string | null): Order
```

Admin update order status with validation

**Implementation Steps:**
1. Load order
2. Validate status transition is allowed
3. Update status
4. If shipped, require and store tracking number
5. Log to order_status_history
6. Queue customer notification email
7. Return updated order

**Error Cases:** [ORDER_NOT_FOUND INVALID_STATUS_TRANSITION]

---

## CustomerService (service)

**Dependencies:** [EmailService PaymentProvider]

**Related Specs:** [TS-BR-CUST-001 TS-BR-CUST-002 TS-BR-CUST-003 TS-BR-SHIP-002 TS-BR-PAY-001]

### Functions

#### `register`

```typescript
register(data: RegistrationRequest): Customer
```

Register new customer with email verification

**Implementation Steps:**
1. Validate email format
2. Check email uniqueness (case-insensitive)
3. Validate password length >= 8
4. Validate password contains number
5. Hash password
6. Create customer with email_verified = false
7. Generate verification token
8. Queue verification email
9. Return customer (without password)

**Error Cases:** [INVALID_EMAIL DUPLICATE_EMAIL PASSWORD_TOO_SHORT PASSWORD_MISSING_NUMBER]

#### `verifyEmail`

```typescript
verifyEmail(token: string): Customer
```

Verify customer email with token

**Implementation Steps:**
1. Look up verification token
2. Check token not expired (24h TTL)
3. Set email_verified = true
4. Delete verification token
5. Return updated customer

**Error Cases:** [INVALID_TOKEN TOKEN_EXPIRED]

#### `addAddress`

```typescript
addAddress(customerId: string, data: AddressRequest): Address
```

Add shipping address with domestic validation

**Implementation Steps:**
1. Validate all required fields present
2. Validate country is in supported list (US only)
3. Create address record
4. If is_default, update other addresses
5. Return address

**Error Cases:** [INVALID_ADDRESS UNSUPPORTED_COUNTRY]

#### `addPaymentMethod`

```typescript
addPaymentMethod(customerId: string, data: PaymentMethodRequest): PaymentMethod
```

Add payment method via Stripe/PayPal

**Implementation Steps:**
1. Validate payment token with provider
2. Check card type is supported (Visa, MC, Amex, PayPal)
3. Store tokenized payment method
4. If is_default, update other methods
5. Return payment method (masked)

**Error Cases:** [INVALID_CARD UNSUPPORTED_CARD_TYPE]

#### `requestDataErasure`

```typescript
requestDataErasure(customerId: string): void
```

GDPR data erasure request

**Implementation Steps:**
1. Verify customer identity
2. Anonymize personal data in customer record
3. Anonymize order history (preserve for business records)
4. Delete addresses and payment methods
5. Log erasure request

---

## InventoryService (service)

**Dependencies:** [AuditService AlertService]

**Related Specs:** [TS-BR-STOCK-001 TS-BR-STOCK-002 TS-BR-STOCK-003 TS-BR-INV-001 TS-BR-INV-002 TS-BR-AUDIT-002]

### Functions

#### `checkAvailability`

```typescript
checkAvailability(variantId: string): InventoryStatus
```

Check stock availability for variant

**Implementation Steps:**
1. Query inventory by variant ID
2. Return available_quantity and stock status

**Error Cases:** [VARIANT_NOT_FOUND]

#### `reserveStock`

```typescript
reserveStock(items: OrderItem[]): void
```

Decrement inventory for order placement

**Implementation Steps:**
1. Begin transaction
2. For each item, SELECT FOR UPDATE inventory row
3. Validate quantity <= available_quantity
4. Decrement available_quantity
5. Log to inventory_history
6. Commit transaction

**Error Cases:** [INSUFFICIENT_STOCK]

#### `restoreStock`

```typescript
restoreStock(items: OrderItem[], reason: string): void
```

Restore inventory on order cancellation

**Implementation Steps:**
1. Begin transaction
2. For each item, increment available_quantity
3. Log to inventory_history with reason
4. Commit transaction

#### `updateInventory`

```typescript
updateInventory(variantId: string, quantity: number | null, delta: number | null, reason: string, adminId: string): Inventory
```

Admin inventory update with audit trail

**Implementation Steps:**
1. Load current inventory
2. Calculate new quantity (absolute or delta)
3. Validate new quantity >= 0
4. Update inventory
5. Log to inventory_history (before/after, reason, admin)
6. Check low stock alert threshold
7. If below threshold, trigger alert
8. Return updated inventory

**Error Cases:** [VARIANT_NOT_FOUND INVALID_STOCK_LEVEL]

#### `bulkUpdate`

```typescript
bulkUpdate(csvFile: File, adminId: string): BulkUpdateResult
```

Bulk inventory update from CSV

**Implementation Steps:**
1. Parse CSV file
2. For each row, validate SKU exists
3. For each row, validate quantity is numeric
4. Apply valid updates with audit trail
5. Collect errors for invalid rows
6. Return summary (success count, error details)

**Error Cases:** [INVALID_CSV_FORMAT]

---

## VariantService (service)

**Dependencies:** [ProductService InventoryService]

**Related Specs:** [TS-BR-VAR-001 TS-BR-VAR-002]

### Functions

#### `createVariant`

```typescript
createVariant(productId: string, data: CreateVariantRequest): Variant
```

Create product variant with unique SKU

**Implementation Steps:**
1. Validate product exists
2. Validate SKU format
3. Check SKU uniqueness
4. Create variant with attributes
5. Initialize inventory record
6. Return variant

**Error Cases:** [PRODUCT_NOT_FOUND DUPLICATE_SKU]

#### `listVariants`

```typescript
listVariants(productId: string): Variant[]
```

List all variants for a product

**Implementation Steps:**
1. Query variants by product ID
2. Include inventory for each
3. Return variants with stock info

**Error Cases:** [PRODUCT_NOT_FOUND]

---

## CategoryService (service)

**Related Specs:** [TS-BR-CAT-001]

### Functions

#### `createCategory`

```typescript
createCategory(data: CreateCategoryRequest): Category
```

Create category with depth validation

**Implementation Steps:**
1. If parent_id provided, load parent
2. Calculate depth from parent chain
3. Validate depth <= 3
4. Create category
5. Return category

**Error Cases:** [CATEGORY_DEPTH_EXCEEDED PARENT_NOT_FOUND]

#### `listCategories`

```typescript
listCategories(): CategoryTree
```

Get category tree structure

**Implementation Steps:**
1. Query all categories
2. Build hierarchical tree structure
3. Return tree with products count per category

---

## PaymentService (service)

**Dependencies:** [StripeClient PayPalClient]

**Related Specs:** [TS-BR-PAY-001 TS-BR-PAY-002]

### Functions

#### `createPaymentIntent`

```typescript
createPaymentIntent(amount: number, currency: string, paymentMethodId: string): PaymentIntent
```

Create Stripe payment intent

**Implementation Steps:**
1. Load payment method details
2. Call Stripe API to create intent
3. Return payment intent ID

**Error Cases:** [INVALID_PAYMENT_METHOD]

#### `authorizePayment`

```typescript
authorizePayment(paymentIntentId: string): AuthorizationResult
```

Authorize payment with provider

**Implementation Steps:**
1. Confirm payment intent with Stripe
2. Handle success/failure response
3. Return authorization result

**Error Cases:** [PAYMENT_FAILED]

#### `refundPayment`

```typescript
refundPayment(paymentIntentId: string, amount: number): RefundResult
```

Process refund for cancelled order

**Implementation Steps:**
1. Call Stripe refund API
2. Log refund transaction
3. Return refund result

**Error Cases:** [REFUND_FAILED]

---

## EmailService (service)

**Dependencies:** [EmailQueue TemplateEngine]

**Related Specs:** [AC-ORDER-007 AC-CUST-001]

### Functions

#### `queueVerificationEmail`

```typescript
queueVerificationEmail(customerId: string, email: string, token: string): void
```

Queue email verification message

**Implementation Steps:**
1. Build verification email template
2. Add to email queue
3. Handle queue failures gracefully

#### `queueOrderConfirmation`

```typescript
queueOrderConfirmation(orderId: string): void
```

Queue order confirmation email

**Implementation Steps:**
1. Load order details
2. Build confirmation template with items, total, shipping
3. Add to email queue
4. Log any failures without blocking order

#### `queueShippingNotification`

```typescript
queueShippingNotification(orderId: string, trackingNumber: string): void
```

Queue shipping notification with tracking

**Implementation Steps:**
1. Load order and customer
2. Build shipping template with tracking link
3. Add to email queue

---

## AuditService (service)

**Related Specs:** [TS-BR-AUDIT-001 TS-BR-AUDIT-002 TS-BR-AUDIT-003]

### Functions

#### `logOrderStatusChange`

```typescript
logOrderStatusChange(orderId: string, oldStatus: string, newStatus: string, userId: string): void
```

Log order status transition

**Implementation Steps:**
1. Insert record to order_status_history
2. Include timestamp and actor

#### `logInventoryChange`

```typescript
logInventoryChange(variantId: string, quantityBefore: number, quantityAfter: number, reason: string, userId: string): void
```

Log inventory adjustment

**Implementation Steps:**
1. Insert record to inventory_history
2. Include before/after values, reason, timestamp

#### `logPriceChange`

```typescript
logPriceChange(productId: string, oldPrice: number, newPrice: number, userId: string): void
```

Log product price change

**Implementation Steps:**
1. Insert record to price_history
2. Include old/new prices, timestamp, modifier

---

