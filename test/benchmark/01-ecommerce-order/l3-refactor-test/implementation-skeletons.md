# Implementation Skeletons

Generated: 2026-01-02T21:37:12+01:00

---

## SKEL-CUST-001: CustomerService (service)

**Dependencies:** [EmailService PasswordService]

**Related Specs:** [TS-BR-CUST-001 TS-BR-CUST-002 TS-BR-CUST-003 TS-BR-AUTH-001]

### Functions

#### `register`

```typescript
register(data: { email: string; password: string; name: string }): Promise<Customer>
```

Register a new customer with email and password

**Implementation Steps:**
1. Validate email format using regex
2. Validate password meets requirements (min 8 chars, at least one digit)
3. Check email uniqueness in database
4. Hash password using secure algorithm
5. Create customer record with emailVerified=false
6. Generate email verification token
7. Send verification email via EmailService
8. Return created customer

**Error Cases:** [WEAK_PASSWORD EMAIL_ALREADY_REGISTERED]

#### `verifyEmail`

```typescript
verifyEmail(token: string): Promise<Customer>
```

Verify customer email address using verification token

**Implementation Steps:**
1. Validate token exists and not expired
2. Find customer by verification token
3. Set emailVerified=true
4. Clear verification token
5. Return updated customer

#### `setPassword`

```typescript
setPassword(customerId: UUID, newPassword: string): Promise<void>
```

Update customer password with validation

**Implementation Steps:**
1. Validate password meets requirements (min 8 chars, at least one digit)
2. Hash new password
3. Update customer password in database

**Error Cases:** [WEAK_PASSWORD]

#### `getById`

```typescript
getById(customerId: UUID): Promise<Customer | null>
```

Retrieve customer by ID

**Implementation Steps:**
1. Query customer by ID
2. Return customer or null if not found

---

## SKEL-PROD-002: ProductService (service)

**Dependencies:** [AuditLogService]

**Related Specs:** [TS-BR-PRICE-001 TS-BR-PROD-001 TS-BR-PROD-002 TS-BR-PROD-003 TS-BR-AUDIT-003]

### Functions

#### `createProduct`

```typescript
createProduct(data: { name: string; description?: string; price: number }): Promise<Product>
```

Create a new product with validation

**Implementation Steps:**
1. Validate name length (2-200 characters)
2. Validate description length if provided (max 5000 characters)
3. Validate price >= 0.01 with exactly 2 decimal places
4. Create product record
5. Return created product

**Error Cases:** [INVALID_PRODUCT_NAME DESCRIPTION_TOO_LONG INVALID_PRICE]

#### `updateProduct`

```typescript
updateProduct(productId: UUID, data: { name?: string; description?: string; price?: number }): Promise<Product>
```

Update existing product with validation and price audit

**Implementation Steps:**
1. Fetch existing product
2. Validate name length if provided (2-200 characters)
3. Validate description length if provided (max 5000 characters)
4. Validate price if provided (>= 0.01, exactly 2 decimal places)
5. If price changed, emit PriceChangeEvent to AuditLogService
6. Update product record
7. Return updated product

**Error Cases:** [INVALID_PRODUCT_NAME DESCRIPTION_TOO_LONG INVALID_PRICE]

#### `deleteProduct`

```typescript
deleteProduct(productId: UUID): Promise<void>
```

Delete or soft-delete product based on order history

**Implementation Steps:**
1. Check if product has existing OrderItems
2. If has order history, set deletedAt timestamp (soft delete)
3. If no order history, perform hard delete
4. Return success

#### `getProduct`

```typescript
getProduct(productId: UUID): Promise<Product | null>
```

Retrieve product by ID, excluding soft-deleted

**Implementation Steps:**
1. Query product by ID where deletedAt IS NULL
2. Return product or null

---

## SKEL-PROD-003: ProductImageService (service)

**Related Specs:** [TS-BR-PROD-004]

### Functions

#### `addImage`

```typescript
addImage(productId: UUID, imageData: { url: string; isPrimary: boolean }): Promise<ProductImage>
```

Add image to product with limit enforcement

**Implementation Steps:**
1. Fetch current image count for product
2. Validate count < 10
3. If isPrimary=true, unset previous primary image
4. If this is first image, set isPrimary=true automatically
5. Create image record
6. Return created image

**Error Cases:** [IMAGE_LIMIT_EXCEEDED]

#### `setPrimary`

```typescript
setPrimary(productId: UUID, imageId: UUID): Promise<void>
```

Set an image as the primary image for product

**Implementation Steps:**
1. Verify image belongs to product
2. Unset current primary image
3. Set specified image as primary
4. Return success

#### `removeImage`

```typescript
removeImage(productId: UUID, imageId: UUID): Promise<void>
```

Remove image from product with primary image validation

**Implementation Steps:**
1. Verify image belongs to product
2. Delete image record
3. If deleted image was primary, validate remaining images have a primary
4. If no primary and images remain, set first remaining as primary

**Error Cases:** [PRIMARY_IMAGE_REQUIRED]

---

## SKEL-VARI-004: VariantService (service)

**Related Specs:** [TS-BR-VAR-001]

### Functions

#### `createVariant`

```typescript
createVariant(productId: UUID, data: { sku: string; attributes: object; price?: number }): Promise<ProductVariant>
```

Create a new product variant with unique SKU

**Implementation Steps:**
1. Validate SKU uniqueness across all variants
2. Validate price if provided (>= 0.01)
3. Create variant record linked to product
4. Return created variant

**Error Cases:** [DUPLICATE_SKU INVALID_PRICE]

#### `getVariantsBySku`

```typescript
getVariantsBySku(sku: string): Promise<ProductVariant | null>
```

Find variant by SKU

**Implementation Steps:**
1. Query variant by SKU
2. Return variant or null

---

## SKEL-INVE-005: InventoryService (service)

**Dependencies:** [AuditLogService]

**Related Specs:** [TS-BR-INV-001 TS-BR-INV-002 TS-BR-INV-003 TS-BR-AUDIT-002]

### Functions

#### `adjustStock`

```typescript
adjustStock(productId: UUID, variantId: UUID | null, adjustment: number, reason: AdjustmentReason): Promise<Inventory>
```

Adjust inventory stock level with audit logging

**Implementation Steps:**
1. Fetch current inventory with row lock (SELECT FOR UPDATE)
2. Calculate new stock level
3. Validate new stock >= 0
4. Update inventory record
5. Emit InventoryEvent to AuditLogService with previousStock, newStock, reason
6. Return updated inventory

**Error Cases:** [INSUFFICIENT_STOCK]

#### `reserveStock`

```typescript
reserveStock(items: Array<{ productId: UUID; variantId?: UUID; quantity: number }>): Promise<void>
```

Atomically reserve stock for order placement

**Implementation Steps:**
1. Begin transaction
2. For each item, SELECT FOR UPDATE inventory row
3. Validate each item.quantity <= availableStock
4. Decrement availableStock for each item
5. Emit InventoryEvent for each item with reason=order_placed
6. Commit transaction

**Error Cases:** [BACKORDER_NOT_SUPPORTED INSUFFICIENT_STOCK]

#### `restoreStock`

```typescript
restoreStock(items: Array<{ productId: UUID; variantId?: UUID; quantity: number }>): Promise<void>
```

Restore stock for cancelled order

**Implementation Steps:**
1. Begin transaction
2. For each item, increment availableStock
3. Emit InventoryEvent for each item with reason=order_cancelled
4. Commit transaction

#### `getAvailableStock`

```typescript
getAvailableStock(productId: UUID, variantId?: UUID): Promise<number>
```

Get current available stock for product/variant

**Implementation Steps:**
1. Query inventory by productId and variantId
2. Return availableStock or 0 if not found

---

## SKEL-CART-006: CartService (service)

**Dependencies:** [InventoryService ProductService]

**Related Specs:** [TS-BR-STOCK-001 TS-BR-VAR-002 TS-BR-CART-001 TS-BR-CART-002 TS-BR-CART-003]

### Functions

#### `addItem`

```typescript
addItem(cartId: UUID, data: { productId: UUID; variantId?: UUID; quantity: number }): Promise<CartItem>
```

Add item to cart with stock validation

**Implementation Steps:**
1. Fetch product and check if it has variants
2. If product has variants, validate variantId is provided
3. Check availableStock > 0
4. Validate quantity <= availableStock, cap if exceeds
5. Snapshot current price as priceAtAdd
6. Create or update cart item
7. Update cart lastModified timestamp
8. Return cart item (with notification if quantity was capped)

**Error Cases:** [OUT_OF_STOCK VARIANT_REQUIRED QUANTITY_EXCEEDS_STOCK]

#### `updateQuantity`

```typescript
updateQuantity(cartId: UUID, itemId: UUID, quantity: number): Promise<CartItem>
```

Update cart item quantity with stock validation

**Implementation Steps:**
1. Fetch cart item
2. Get current availableStock for product/variant
3. Validate quantity <= availableStock
4. Update cart item quantity
5. Update cart lastModified timestamp
6. Return updated cart item

**Error Cases:** [QUANTITY_EXCEEDS_STOCK]

#### `getCart`

```typescript
getCart(cartId: UUID): Promise<Cart>
```

Retrieve cart with current prices and totals

**Implementation Steps:**
1. Check cart expiration based on lastModified and user type
2. If expired, return CART_EXPIRED error
3. Fetch cart with items
4. For each item, join to current Product.price
5. Calculate total using current prices
6. Include price change notifications where priceAtAdd differs from currentPrice
7. Return cart with items and calculated totals

**Error Cases:** [CART_EXPIRED]

#### `removeItem`

```typescript
removeItem(cartId: UUID, itemId: UUID): Promise<void>
```

Remove item from cart

**Implementation Steps:**
1. Delete cart item
2. Update cart lastModified timestamp

---

## SKEL-CART-007: CartCleanupJob (service)

**Related Specs:** [TS-BR-CART-003]

### Functions

#### `cleanupExpiredCarts`

```typescript
cleanupExpiredCarts(): Promise<{ deleted: number }>
```

Scheduled job to delete expired carts

**Implementation Steps:**
1. Query carts where isGuest=true AND lastModified + 7 days < now()
2. Query carts where isGuest=false AND lastModified + 30 days < now()
3. Delete matching carts and their items
4. Return count of deleted carts

---

## SKEL-CHEC-008: CheckoutService (service)

**Dependencies:** [CustomerService CartService InventoryService PaymentService ShippingService OrderService AuditLogService AddressValidator]

**Related Specs:** [TS-BR-AUTH-001 TS-BR-CUST-003 TS-BR-INV-002 TS-BR-PRICE-002 TS-BR-PAY-001 TS-BR-ADDR-001 TS-BR-SHIP-002]

### Functions

#### `placeOrder`

```typescript
placeOrder(customerId: UUID, cartId: UUID, shippingAddress: ShippingAddress, paymentDetails: PaymentDetails): Promise<Order>
```

Complete checkout flow with all validations

**Implementation Steps:**
1. Validate customer exists and is authenticated
2. Validate customer.emailVerified = true
3. Validate cart exists and not expired
4. Validate all cart items still have available stock
5. Validate shipping address has all required fields
6. Validate address.country matches domesticCountry config
7. Calculate order totals using calculateTotals()
8. Authorize payment via PaymentService
9. If payment fails, return PAYMENT_REQUIRED error
10. Begin transaction
11. Reserve stock via InventoryService.reserveStock()
12. Create Order with status='pending'
13. Create OrderItems with unitPrice snapshot from current Product.price
14. Emit OrderEvent to AuditLogService with type=created
15. Clear cart
16. Commit transaction
17. Return created order

**Error Cases:** [REGISTRATION_REQUIRED EMAIL_NOT_VERIFIED OUT_OF_STOCK BACKORDER_NOT_SUPPORTED INCOMPLETE_ADDRESS INTERNATIONAL_SHIPPING_NOT_AVAILABLE PAYMENT_REQUIRED]

#### `validateCart`

```typescript
validateCart(cartId: UUID): Promise<ValidationResult>
```

Validate cart items for checkout readiness

**Implementation Steps:**
1. Fetch cart with items
2. For each item, check availableStock >= quantity
3. Return validation result with any stock issues

**Error Cases:** [OUT_OF_STOCK]

#### `calculateTotals`

```typescript
calculateTotals(cartId: UUID): Promise<OrderTotals>
```

Calculate order totals including shipping

**Implementation Steps:**
1. Fetch cart with current product prices
2. Calculate subtotal from item quantities and prices
3. Apply ShippingService.calculateShipping() for shipping cost
4. Calculate grand total
5. Return totals breakdown

---

## SKEL-ORDE-009: OrderService (service)

**Dependencies:** [InventoryService PaymentService AuditLogService]

**Related Specs:** [TS-BR-CANCEL-001 TS-BR-PRICE-002 TS-BR-ORDER-001 TS-BR-ORDER-002 TS-BR-ORDER-003 TS-BR-INV-003 TS-BR-PAY-002 TS-BR-AUDIT-001]

### Functions

#### `createOrder`

```typescript
createOrder(data: CreateOrderDTO): Promise<Order>
```

Create order with immutable items and address

**Implementation Steps:**
1. Create Order aggregate with status='pending'
2. Create OrderItems with immutable unitPrice
3. Create immutable ShippingAddress snapshot
4. Create one-to-one Shipment placeholder
5. Return created order

#### `cancelOrder`

```typescript
cancelOrder(orderId: UUID): Promise<Order>
```

Cancel order with stock restoration and refund

**Implementation Steps:**
1. Fetch order
2. Validate order.status IN ('pending', 'confirmed')
3. Begin transaction
4. Update order status to 'cancelled'
5. Restore stock via InventoryService.restoreStock() for all items
6. Process refund via PaymentService.refund()
7. Emit OrderEvent to AuditLogService with type=cancelled
8. Commit transaction
9. Return updated order

**Error Cases:** [CANCELLATION_NOT_ALLOWED]

#### `updateStatus`

```typescript
updateStatus(orderId: UUID, newStatus: OrderStatus): Promise<Order>
```

Transition order to new status following state machine

**Implementation Steps:**
1. Fetch order with current status
2. Get allowed transitions for current status from state machine
3. Validate newStatus is in allowed transitions
4. Update order status
5. Emit OrderEvent to AuditLogService with type=status_changed, previousStatus, newStatus
6. Return updated order

**Error Cases:** [INVALID_STATUS_TRANSITION]

#### `getOrder`

```typescript
getOrder(orderId: UUID): Promise<Order | null>
```

Retrieve order by ID

**Implementation Steps:**
1. Query order with items and shipping address
2. Return order or null

#### `getCustomerOrders`

```typescript
getCustomerOrders(customerId: UUID): Promise<Order[]>
```

Get all orders for a customer

**Implementation Steps:**
1. Query orders by customerId
2. Return orders list

---

## SKEL-SHIP-010: ShippingService (service)

**Related Specs:** [TS-BR-SHIP-001]

### Functions

#### `calculateShipping`

```typescript
calculateShipping(subtotal: number): Promise<{ shippingCost: number; freeShipping: boolean }>
```

Calculate shipping cost based on order subtotal

**Implementation Steps:**
1. Check if subtotal >= 50.00
2. If true, return { shippingCost: 0.00, freeShipping: true }
3. Else return { shippingCost: 5.99, freeShipping: false }

---

## SKEL-SHIP-011: ShipmentService (service)

**Related Specs:** [TS-BR-ORDER-003]

### Functions

#### `createShipment`

```typescript
createShipment(orderId: UUID): Promise<Shipment>
```

Create shipment for order (one-to-one)

**Implementation Steps:**
1. Validate no existing shipment for order
2. Create shipment record linked to order
3. Return created shipment

---

## SKEL-PAYM-012: PaymentService (service)

**Related Specs:** [TS-BR-PAY-001 TS-BR-PAY-002]

### Functions

#### `authorize`

```typescript
authorize(paymentDetails: PaymentDetails, amount: number): Promise<PaymentAuthorization>
```

Authorize payment with payment gateway

**Implementation Steps:**
1. Validate payment details
2. Call payment gateway authorization API
3. If successful, return authorization with paymentAuthorizationId
4. If failed, throw PAYMENT_REQUIRED error

**Error Cases:** [PAYMENT_REQUIRED]

#### `refund`

```typescript
refund(paymentId: string): Promise<RefundResult>
```

Process refund for cancelled order

**Implementation Steps:**
1. Call payment gateway refund API
2. Track refund status (pending|completed|failed)
3. Return refund result

---

## SKEL-ADDR-013: AddressValidator (service)

**Related Specs:** [TS-BR-ADDR-001 TS-BR-SHIP-002]

### Functions

#### `validate`

```typescript
validate(address: ShippingAddress): ValidationResult
```

Validate shipping address completeness and domestic restriction

**Implementation Steps:**
1. Validate street is present and non-empty
2. Validate city is present and non-empty
3. Validate state is present and non-empty
4. Validate postalCode is present and non-empty
5. Validate country is present and valid ISO 3166-1 alpha-2
6. Validate recipientName is present and non-empty
7. Validate country matches domesticCountry config
8. Return validation result

**Error Cases:** [INCOMPLETE_ADDRESS INTERNATIONAL_SHIPPING_NOT_AVAILABLE]

---

## SKEL-CATE-014: CategoryService (service)

**Related Specs:** [TS-BR-CAT-001]

### Functions

#### `createCategory`

```typescript
createCategory(data: { name: string; parentId?: UUID }): Promise<Category>
```

Create category with hierarchy depth validation

**Implementation Steps:**
1. If parentId provided, fetch parent category
2. Calculate resulting depth by traversing parent chain
3. Validate depth <= 3
4. Create category record
5. Return created category

**Error Cases:** [CATEGORY_DEPTH_EXCEEDED]

#### `updateCategory`

```typescript
updateCategory(categoryId: UUID, data: { name?: string; parentId?: UUID }): Promise<Category>
```

Update category with hierarchy depth validation

**Implementation Steps:**
1. Fetch category
2. If parentId changed, calculate new depth
3. Validate new depth <= 3
4. Validate no circular references
5. Update category record
6. Return updated category

**Error Cases:** [CATEGORY_DEPTH_EXCEEDED]

---

## SKEL-AUDI-015: AuditLogService (service)

**Related Specs:** [TS-BR-AUDIT-001 TS-BR-AUDIT-002 TS-BR-AUDIT-003]

### Functions

#### `logOrderEvent`

```typescript
logOrderEvent(event: { orderId: UUID; eventType: OrderEventType; userId: UUID; previousValue?: object; newValue?: object }): Promise<void>
```

Log order-related events for audit trail

**Implementation Steps:**
1. Create audit log entry with timestamp (UTC)
2. Include orderId, eventType, userId
3. Include previousValue and newValue as JSON
4. Persist to audit log storage

#### `logInventoryEvent`

```typescript
logInventoryEvent(event: { productId: UUID; variantId?: UUID; previousStock: number; newStock: number; reason: AdjustmentReason; operatorId: UUID }): Promise<void>
```

Log inventory changes for audit trail

**Implementation Steps:**
1. Create audit log entry with timestamp (UTC)
2. Include productId, variantId, previousStock, newStock
3. Include reason and operatorId
4. Persist to audit log storage

#### `logPriceChange`

```typescript
logPriceChange(event: { productId: UUID; previousPrice: number; newPrice: number; modifierId: UUID }): Promise<void>
```

Log product price changes for audit trail

**Implementation Steps:**
1. Create audit log entry with timestamp (UTC)
2. Include productId, previousPrice, newPrice, modifierId
3. Persist to audit log storage

---

## SKEL-PASS-016: PasswordService (service)

**Related Specs:** [TS-BR-CUST-002]

### Functions

#### `validate`

```typescript
validate(password: string): boolean
```

Validate password meets strength requirements

**Implementation Steps:**
1. Check password length >= 8
2. Check password contains at least one digit using regex
3. Return true if all checks pass

#### `hash`

```typescript
hash(password: string): Promise<string>
```

Hash password using secure algorithm

**Implementation Steps:**
1. Generate salt
2. Hash password with salt using bcrypt or argon2
3. Return hashed password

#### `verify`

```typescript
verify(password: string, hash: string): Promise<boolean>
```

Verify password against stored hash

**Implementation Steps:**
1. Compare password with hash
2. Return true if match

---

## SKEL-EMAI-017: EmailService (service)

**Related Specs:** [TS-BR-AUTH-001]

### Functions

#### `sendVerificationEmail`

```typescript
sendVerificationEmail(email: string, token: string): Promise<void>
```

Send email verification link to customer

**Implementation Steps:**
1. Build verification URL with token
2. Render email template
3. Send email via email provider
4. Log email sent

---

