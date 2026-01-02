# Implementation Skeletons

Generated: 2026-01-02T19:24:29+01:00

---

## ProductService (service)

**Dependencies:** [CategoryService AuditLogService CartService StorageService]

**Related Specs:** [TS-BR-PROD-001 TS-BR-PROD-002 TS-BR-PROD-003 TS-BR-PROD-004 TS-BR-AUDIT-003]

### Functions

#### `create`

```typescript
create(data: CreateProductDTO): Product
```

Create new product with validation

**Implementation Steps:**
1. Validate name length (2-200 chars)
2. Validate price >= $0.01
3. Validate category exists
4. Check for duplicate name (warn only)
5. Generate UUID
6. Set status to draft
7. Set created_by and timestamps
8. Persist product

**Error Cases:** [INVALID_PRODUCT_NAME INVALID_PRICE CATEGORY_NOT_FOUND]

#### `update`

```typescript
update(productId: UUID, data: UpdateProductDTO): Product
```

Update product with price audit trail

**Implementation Steps:**
1. Fetch existing product
2. Validate name length if provided
3. Validate price >= $0.01 if provided
4. Detect price change for audit
5. Update fields
6. Set modified_by and timestamps
7. Log price change to audit trail if changed
8. Persist product

**Error Cases:** [PRODUCT_NOT_FOUND INVALID_PRODUCT_NAME INVALID_PRICE]

#### `delete`

```typescript
delete(productId: UUID): void
```

Soft delete product, remove from active carts

**Implementation Steps:**
1. Check if product has order history
2. If has orders, set deleted_at timestamp (soft delete)
3. If no orders, perform hard delete
4. Remove product from all active carts
5. Notify affected cart owners

**Error Cases:** [PRODUCT_NOT_FOUND]

#### `addImage`

```typescript
addImage(productId: UUID, image: File, isPrimary: boolean): ProductImage
```

Upload product image with limit check

**Implementation Steps:**
1. Count existing images
2. Reject if count >= 10
3. Upload image to storage
4. If isPrimary, unset other primary flags
5. Create ProductImage record

**Error Cases:** [IMAGE_LIMIT_EXCEEDED PRODUCT_NOT_FOUND]

---

## VariantService (service)

**Dependencies:** [ProductService InventoryService]

**Related Specs:** [TS-BR-VAR-001]

### Functions

#### `create`

```typescript
create(productId: UUID, data: CreateVariantDTO): ProductVariant
```

Create product variant with unique SKU

**Implementation Steps:**
1. Validate product exists
2. Check SKU uniqueness
3. Create variant with options
4. Initialize inventory record with 0 stock

**Error Cases:** [PRODUCT_NOT_FOUND DUPLICATE_SKU]

---

## CartService (service)

**Dependencies:** [ProductService InventoryService]

**Related Specs:** [TS-BR-CART-001 TS-BR-CART-002 TS-BR-CART-003 TS-BR-STOCK-001 TS-BR-VAR-002]

### Functions

#### `getCart`

```typescript
getCart(sessionOrCustomerId: string): Cart
```

Get cart with current product prices

**Implementation Steps:**
1. Find cart by session or customer ID
2. Check cart expiration (7 days guest, 30 days authenticated)
3. If expired, return CART_EXPIRED error
4. For each item, fetch current price from Product/Variant
5. Detect and flag price changes since item was added
6. Calculate totals with current prices

**Error Cases:** [CART_EXPIRED]

#### `addItem`

```typescript
addItem(cartId: UUID, productId: UUID, variantId: UUID | null, quantity: number): CartItem
```

Add item to cart with stock validation

**Implementation Steps:**
1. Validate product exists and not deleted
2. If product has variants, require variantId
3. Check available stock
4. If stock == 0, return OUT_OF_STOCK error
5. Check for existing cart item (same product/variant)
6. If exists, increment quantity (cap at stock)
7. If new, create cart item
8. Store price_at_add for change detection
9. Update cart last_activity_at
10. Return actual quantity added (may be capped)

**Error Cases:** [OUT_OF_STOCK VARIANT_REQUIRED PRODUCT_NOT_FOUND QUANTITY_EXCEEDS_STOCK]

#### `updateQuantity`

```typescript
updateQuantity(itemId: UUID, quantity: number): CartItem
```

Update cart item quantity with stock limit

**Implementation Steps:**
1. Find cart item
2. If quantity == 0, remove item
3. Check available stock
4. Cap quantity at available stock
5. Update quantity
6. Update cart last_activity_at
7. Notify if quantity was capped

**Error Cases:** [CART_ITEM_NOT_FOUND QUANTITY_EXCEEDS_STOCK]

#### `removeItem`

```typescript
removeItem(itemId: UUID): UndoToken
```

Remove item immediately with undo option

**Implementation Steps:**
1. Find and remove cart item
2. Store removed item in temporary undo buffer
3. Generate undo token with TTL
4. Update cart last_activity_at

**Error Cases:** [CART_ITEM_NOT_FOUND]

#### `mergeGuestCart`

```typescript
mergeGuestCart(customerId: UUID, guestSessionId: string): Cart
```

Merge guest cart into authenticated user cart

**Implementation Steps:**
1. Find guest cart by session
2. Find or create customer cart
3. For each guest item:
4.   - Check if product/variant exists in customer cart
5.   - If exists, combine quantities (cap at stock)
6.   - If new, add to customer cart
7. Delete guest cart
8. Notify if any quantities were capped

**Error Cases:** [GUEST_CART_NOT_FOUND]

---

## CustomerService (service)

**Dependencies:** [EmailService]

**Related Specs:** [TS-BR-CUST-001 TS-BR-CUST-002 TS-BR-CUST-003]

### Functions

#### `register`

```typescript
register(data: RegisterDTO): Customer
```

Register new customer with email verification

**Implementation Steps:**
1. Validate email format
2. Check email uniqueness
3. Validate password (8+ chars, at least one digit)
4. Hash password
5. Create customer with emailVerified=false
6. Generate verification token
7. Queue verification email

**Error Cases:** [EMAIL_ALREADY_REGISTERED WEAK_PASSWORD]

#### `verifyEmail`

```typescript
verifyEmail(customerId: UUID, token: string): void
```

Verify customer email address

**Implementation Steps:**
1. Find customer
2. Validate token matches and not expired
3. Set emailVerified=true
4. Clear verification token

**Error Cases:** [INVALID_TOKEN TOKEN_EXPIRED]

#### `changePassword`

```typescript
changePassword(customerId: UUID, currentPassword: string, newPassword: string): void
```

Change customer password with validation

**Implementation Steps:**
1. Find customer
2. Verify current password
3. Validate new password strength
4. Hash new password
5. Update password

**Error Cases:** [INVALID_CURRENT_PASSWORD WEAK_PASSWORD]

#### `deleteAccount`

```typescript
deleteAccount(customerId: UUID, deletionType: 'soft_delete' | 'full_erasure', confirmed: boolean): void
```

Delete customer account with pending order check

**Implementation Steps:**
1. Verify confirmation flag is true
2. Check for pending/processing orders
3. If has incomplete orders, reject deletion
4. If soft_delete, set deactivated status
5. If full_erasure, anonymize all personal data

**Error Cases:** [CONFIRMATION_REQUIRED HAS_PENDING_ORDERS]

---

## AddressService (service)

**Dependencies:** [ConfigService]

**Related Specs:** [TS-BR-ADDR-001 TS-BR-SHIP-002]

### Functions

#### `create`

```typescript
create(customerId: UUID, data: CreateAddressDTO): ShippingAddress
```

Add shipping address with validation

**Implementation Steps:**
1. Validate all required fields present
2. Validate country is domestic
3. Validate postal code format for country
4. If is_default, unset other defaults
5. Create address

**Error Cases:** [INCOMPLETE_ADDRESS INTERNATIONAL_SHIPPING_NOT_AVAILABLE INVALID_POSTAL_CODE]

#### `validate`

```typescript
validate(address: ShippingAddress): ValidationResult
```

Validate address fields

**Implementation Steps:**
1. Check recipient_name not empty
2. Check street not empty
3. Check city not empty
4. Check state not empty
5. Check postal_code not empty
6. Check country is valid ISO code
7. Check country matches domestic config

**Error Cases:** [INCOMPLETE_ADDRESS]

---

## OrderService (service)

**Dependencies:** [CartService CustomerService InventoryService PaymentService ShippingService TaxService EmailService AuditLogService]

**Related Specs:** [TS-BR-AUTH-001 TS-BR-STOCK-001 TS-BR-INV-002 TS-BR-INV-003 TS-BR-ORDER-001 TS-BR-ORDER-002 TS-BR-PRICE-002 TS-BR-PAY-001 TS-BR-PAY-002 TS-BR-CANCEL-001 TS-BR-AUDIT-001]

### Functions

#### `create`

```typescript
create(customerId: UUID, addressId: UUID, paymentMethodId: UUID): Order
```

Create order with full validation flow

**Implementation Steps:**
1. Validate customer is registered
2. Validate customer email is verified
3. Get cart and validate not empty
4. Validate all items have sufficient stock
5. Authorize payment
6. If payment fails, return PAYMENT_REQUIRED error
7. Create order with status=pending
8. For each cart item:
9.   - Create OrderItem with captured price
10.   - Decrement inventory atomically
11. Calculate shipping (free if subtotal >= $50)
12. Calculate tax based on shipping state
13. Clear cart
14. Queue confirmation email
15. Log order creation to audit

**Error Cases:** [REGISTRATION_REQUIRED EMAIL_NOT_VERIFIED EMPTY_CART INSUFFICIENT_STOCK BACKORDER_NOT_SUPPORTED PAYMENT_REQUIRED]

#### `cancel`

```typescript
cancel(orderId: UUID, customerId: UUID): void
```

Cancel order with inventory restoration and refund

**Implementation Steps:**
1. Find order
2. Verify customer owns order
3. Validate status is pending or confirmed
4. Set status to cancelled
5. For each order item, restore inventory
6. Initiate refund via payment gateway
7. Queue cancellation email
8. Log cancellation to audit

**Error Cases:** [ORDER_NOT_FOUND CANCELLATION_NOT_ALLOWED UNAUTHORIZED]

#### `updateStatus`

```typescript
updateStatus(orderId: UUID, newStatus: OrderStatus, trackingNumber?: string): Order
```

Update order status following state machine

**Implementation Steps:**
1. Find order
2. Validate transition is allowed by state machine
3. Update status
4. If shipped, optionally set tracking number
5. Log status change to audit
6. Queue notification email to customer

**Error Cases:** [ORDER_NOT_FOUND INVALID_STATUS_TRANSITION]

---

## OrderStateMachine (utility)

**Related Specs:** [TS-BR-ORDER-001]

### Functions

#### `validateTransition`

```typescript
validateTransition(currentStatus: OrderStatus, newStatus: OrderStatus): boolean
```

Check if status transition is valid

**Implementation Steps:**
1. Define allowed transitions map:
2.   - pending → [confirmed, cancelled]
3.   - confirmed → [shipped, cancelled]
4.   - shipped → [delivered]
5.   - delivered → []
6.   - cancelled → []
7. Check if newStatus is in allowed transitions for currentStatus

**Error Cases:** [INVALID_STATUS_TRANSITION]

---

## InventoryService (service)

**Dependencies:** [AuditLogService AlertService]

**Related Specs:** [TS-BR-INV-001 TS-BR-INV-002 TS-BR-INV-003 TS-BR-AUDIT-002]

### Functions

#### `adjust`

```typescript
adjust(variantId: UUID, adjustmentType: 'absolute' | 'delta', value: number, reason?: string): Inventory
```

Adjust inventory with audit trail

**Implementation Steps:**
1. Find inventory record
2. Calculate new stock level
3. If result would be negative, reject
4. Check low stock threshold
5. Update available_stock
6. Log adjustment to audit with before/after values
7. If crossed low stock threshold, trigger alert

**Error Cases:** [VARIANT_NOT_FOUND INSUFFICIENT_STOCK]

#### `decrement`

```typescript
decrement(variantId: UUID, quantity: number): void
```

Atomically decrement inventory for order

**Implementation Steps:**
1. Execute atomic UPDATE with WHERE available_stock >= quantity
2. If no rows affected, throw INSUFFICIENT_STOCK

**Error Cases:** [INSUFFICIENT_STOCK]

#### `restore`

```typescript
restore(variantId: UUID, quantity: number): void
```

Restore inventory on order cancellation

**Implementation Steps:**
1. Increment available_stock by quantity
2. Log restoration to audit

**Error Cases:** [VARIANT_NOT_FOUND]

#### `importFromCSV`

```typescript
importFromCSV(file: File): ImportResult
```

Bulk import inventory from CSV

**Implementation Steps:**
1. Parse CSV file
2. For each row:
3.   - Validate variant exists
4.   - Validate quantity is non-negative number
5.   - Apply adjustment
6.   - Track success/failure per row
7. Return results with row-by-row errors

**Error Cases:** [INVALID_FILE_FORMAT]

---

## ShippingService (service)

**Related Specs:** [TS-BR-SHIP-001]

### Functions

#### `calculateCost`

```typescript
calculateCost(subtotal: number): number
```

Calculate shipping cost with free shipping threshold

**Implementation Steps:**
1. If subtotal >= 50.00, return 0
2. Otherwise return 5.99

---

## TaxService (service)

**Dependencies:** [TaxRateRepository]

**Related Specs:** [TS-BR-SHIP-001]

### Functions

#### `calculateTax`

```typescript
calculateTax(subtotal: number, shippingAddress: ShippingAddress): number
```

Calculate tax based on shipping state

**Implementation Steps:**
1. Validate state is provided
2. Look up tax rate for state
3. Calculate tax = subtotal * rate
4. Return tax amount (0 for no-tax states)

**Error Cases:** [STATE_REQUIRED INVALID_STATE]

---

## PaymentService (service)

**Dependencies:** [PaymentGateway (Stripe)]

**Related Specs:** [TS-BR-PAY-001 TS-BR-PAY-002]

### Functions

#### `authorize`

```typescript
authorize(paymentMethodId: UUID, amount: number): PaymentAuthorization
```

Authorize payment before order creation

**Implementation Steps:**
1. Get payment method (tokenized card or PayPal)
2. Call payment gateway to authorize amount
3. Return authorization ID on success

**Error Cases:** [PAYMENT_FAILED CARD_DECLINED]

#### `refund`

```typescript
refund(authorizationId: string, amount: number): RefundResult
```

Initiate refund for cancelled order

**Implementation Steps:**
1. Call payment gateway to refund
2. Track refund status (async completion)

**Error Cases:** [REFUND_FAILED]

---

## CheckoutService (service)

**Dependencies:** [CustomerService CartService AddressService ShippingService TaxService PaymentService OrderService]

**Related Specs:** [TS-BR-AUTH-001 TS-BR-STOCK-001 TS-BR-CUST-003]

### Functions

#### `placeOrder`

```typescript
placeOrder(customerId: UUID, shippingAddressId: UUID, paymentMethodId: UUID): Order
```

Orchestrate complete checkout flow

**Implementation Steps:**
1. Validate customer registered and email verified
2. Get and validate cart (not empty, items in stock)
3. Get shipping address
4. Calculate subtotal, shipping, tax
5. Authorize payment for total
6. Create order via OrderService
7. Return order with confirmation

**Error Cases:** [REGISTRATION_REQUIRED EMAIL_NOT_VERIFIED EMPTY_CART INSUFFICIENT_STOCK PAYMENT_REQUIRED]

#### `validateCart`

```typescript
validateCart(cart: Cart): ValidationResult
```

Validate all cart items for checkout

**Implementation Steps:**
1. Check cart not empty
2. For each item, check stock availability
3. Collect any out-of-stock items
4. Return validation result with issues

**Error Cases:** [EMPTY_CART ITEMS_OUT_OF_STOCK]

---

## CategoryService (service)

**Related Specs:** [TS-BR-CAT-001]

### Functions

#### `create`

```typescript
create(name: string, parentId: UUID | null, isActive: boolean): Category
```

Create category with depth validation

**Implementation Steps:**
1. If parentId provided, fetch parent
2. Calculate depth = parent.depth + 1
3. If depth > 3, reject
4. Create category

**Error Cases:** [PARENT_NOT_FOUND CATEGORY_DEPTH_EXCEEDED]

#### `getHierarchy`

```typescript
getHierarchy(): CategoryTree
```

Get full category tree for display

**Implementation Steps:**
1. Fetch all active categories
2. Build tree structure by parent relationships
3. Return nested tree

---

## AuditLogService (service)

**Related Specs:** [TS-BR-AUDIT-001 TS-BR-AUDIT-002 TS-BR-AUDIT-003]

### Functions

#### `log`

```typescript
log(entityType: string, entityId: UUID, action: string, beforeValue: object | null, afterValue: object, userId: UUID, reason?: string): void
```

Log audit event

**Implementation Steps:**
1. Create audit log entry with timestamp
2. Store before/after as JSONB
3. Persist to audit_log table

---

## EmailService (service)

**Dependencies:** [EmailQueue]

**Related Specs:** [TS-BR-PAY-002]

### Functions

#### `queueOrderConfirmation`

```typescript
queueOrderConfirmation(order: Order): void
```

Queue order confirmation email asynchronously

**Implementation Steps:**
1. Build email with order details
2. Add to email queue
3. Do not block on send

#### `queueCancellationConfirmation`

```typescript
queueCancellationConfirmation(order: Order): void
```

Queue cancellation confirmation email

**Implementation Steps:**
1. Build cancellation email
2. Add to email queue

#### `queueStatusNotification`

```typescript
queueStatusNotification(order: Order, newStatus: OrderStatus): void
```

Queue order status update notification

**Implementation Steps:**
1. Build status update email
2. Add to email queue

---

## CartExpirationJob (scheduled_job)

**Dependencies:** [CartRepository]

**Related Specs:** [TS-BR-CART-003]

### Functions

#### `run`

```typescript
run(): void
```

Clean up expired carts daily

**Implementation Steps:**
1. Find carts where:
2.   - Guest carts with last_activity_at + 7 days < now
3.   - Authenticated carts with last_activity_at + 30 days < now
4. Delete expired carts

---

