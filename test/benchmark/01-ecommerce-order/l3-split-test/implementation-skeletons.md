# Implementation Skeletons

Generated: 2026-01-02T20:15:12+01:00

---

## SKEL-CUST-001: CustomerService (service)

**Dependencies:** [EmailService PasswordHashService]

**Related Specs:** [TS-BR-AUTH-001 TS-BR-CUST-001 TS-BR-CUST-002 TS-BR-CUST-003]

### Functions

#### `register`

```typescript
register(data: RegisterDTO): Promise<Customer>
```

Register a new customer with email and password

**Implementation Steps:**
1. Validate email format using regex
2. Check email uniqueness in database
3. Validate password strength (min 8 chars, at least one digit)
4. Hash password using bcrypt
5. Create customer record with emailVerified=false
6. Send verification email
7. Return customer object

**Error Cases:** [EMAIL_ALREADY_REGISTERED WEAK_PASSWORD]

#### `verifyEmail`

```typescript
verifyEmail(token: string): Promise<Customer>
```

Verify customer email address using verification token

**Implementation Steps:**
1. Find customer by verification token
2. Check token validity and expiration
3. Set emailVerified=true
4. Clear verification token
5. Return updated customer

#### `changePassword`

```typescript
changePassword(customerId: UUID, currentPassword: string, newPassword: string): Promise<void>
```

Change customer password with validation

**Implementation Steps:**
1. Find customer by ID
2. Verify current password matches
3. Validate new password strength (min 8 chars, at least one digit)
4. Hash new password
5. Update customer password
6. Return success

**Error Cases:** [WEAK_PASSWORD]

---

## SKEL-PROD-002: ProductService (service)

**Dependencies:** [AuditLogService]

**Related Specs:** [TS-BR-PRICE-001 TS-BR-PROD-001 TS-BR-PROD-002 TS-BR-PROD-003 TS-BR-PROD-004 TS-BR-AUDIT-003]

### Functions

#### `create`

```typescript
create(data: CreateProductDTO): Promise<Product>
```

Create a new product with validation

**Implementation Steps:**
1. Validate name length (2-200 characters)
2. Validate description length (max 5000 characters)
3. Validate price >= 0.01
4. Create product record
5. Return created product

**Error Cases:** [INVALID_PRODUCT_NAME DESCRIPTION_TOO_LONG INVALID_PRICE]

#### `update`

```typescript
update(productId: UUID, data: UpdateProductDTO): Promise<Product>
```

Update existing product with validation and price audit

**Implementation Steps:**
1. Find product by ID
2. Validate name length if provided (2-200 characters)
3. Validate description length if provided (max 5000 characters)
4. Validate price >= 0.01 if provided
5. Check if price changed for audit logging
6. Update product record
7. If price changed, log to audit via AuditLogService
8. Return updated product

**Error Cases:** [INVALID_PRODUCT_NAME DESCRIPTION_TOO_LONG INVALID_PRICE]

#### `delete`

```typescript
delete(productId: UUID): Promise<void>
```

Soft delete product if has order history, otherwise hard delete

**Implementation Steps:**
1. Find product by ID
2. Check if product has OrderItem references
3. If has orders: set deletedAt timestamp (soft delete)
4. If no orders: perform hard delete
5. Return success

#### `addImage`

```typescript
addImage(productId: UUID, image: ImageDTO): Promise<ProductImage>
```

Add image to product with limit validation

**Implementation Steps:**
1. Find product by ID
2. Count existing images
3. Validate count < 10
4. If first image, set as primary
5. Save image record
6. Return created image

**Error Cases:** [IMAGE_LIMIT_EXCEEDED]

#### `setPrimaryImage`

```typescript
setPrimaryImage(productId: UUID, imageId: UUID): Promise<void>
```

Set specified image as primary, unsetting others

**Implementation Steps:**
1. Find product by ID
2. Find image by ID
3. Unset isPrimary on all product images
4. Set isPrimary on target image
5. Return success

**Error Cases:** [PRIMARY_IMAGE_REQUIRED]

---

## SKEL-VARI-003: VariantService (service)

**Dependencies:** [ProductService]

**Related Specs:** [TS-BR-VAR-001]

### Functions

#### `create`

```typescript
create(productId: UUID, data: CreateVariantDTO): Promise<ProductVariant>
```

Create a new product variant with unique SKU

**Implementation Steps:**
1. Find product by ID
2. Validate SKU uniqueness across all variants
3. Create variant record linked to product
4. Return created variant

**Error Cases:** [DUPLICATE_SKU]

#### `update`

```typescript
update(variantId: UUID, data: UpdateVariantDTO): Promise<ProductVariant>
```

Update existing variant with SKU uniqueness check

**Implementation Steps:**
1. Find variant by ID
2. If SKU changed, validate uniqueness
3. Update variant record
4. Return updated variant

**Error Cases:** [DUPLICATE_SKU]

---

## SKEL-INVE-004: InventoryService (service)

**Dependencies:** [AuditLogService]

**Related Specs:** [TS-BR-STOCK-001 TS-BR-INV-001 TS-BR-INV-002 TS-BR-INV-003 TS-BR-AUDIT-002]

### Functions

#### `adjust`

```typescript
adjust(productId: UUID, variantId: UUID | null, quantity: number, reason: string): Promise<Inventory>
```

Adjust inventory stock level with audit logging

**Implementation Steps:**
1. Find inventory record for product/variant
2. Calculate new stock level
3. Validate resulting quantity >= 0
4. Update availableStock atomically
5. Log adjustment to audit via AuditLogService
6. Return updated inventory

**Error Cases:** [INSUFFICIENT_STOCK]

#### `decrement`

```typescript
decrement(productId: UUID, variantId: UUID | null, quantity: number): Promise<Inventory>
```

Decrement inventory for order placement

**Implementation Steps:**
1. Find inventory record for product/variant
2. Validate availableStock >= quantity
3. Atomically decrement availableStock with DB check constraint
4. Log decrement to audit via AuditLogService
5. Return updated inventory

**Error Cases:** [INSUFFICIENT_STOCK]

#### `restore`

```typescript
restore(orderItems: OrderItem[]): Promise<void>
```

Restore inventory for cancelled order items

**Implementation Steps:**
1. For each OrderItem in list
2. Find inventory record for product/variant
3. Increment availableStock by item quantity
4. Log restoration to audit via AuditLogService
5. Return success

#### `checkAvailability`

```typescript
checkAvailability(productId: UUID, variantId: UUID | null): Promise<number>
```

Get available stock for a product/variant

**Implementation Steps:**
1. Find inventory record for product/variant
2. Return availableStock value

---

## SKEL-CART-005: CartService (service)

**Dependencies:** [InventoryService ShippingService]

**Related Specs:** [TS-BR-STOCK-001 TS-BR-CART-001 TS-BR-CART-002 TS-BR-CART-003 TS-BR-VAR-002]

### Functions

#### `getCart`

```typescript
getCart(cartId: UUID): Promise<Cart>
```

Get cart with current prices and expiration check

**Implementation Steps:**
1. Find cart by ID
2. Check if cart expired based on lastActivityAt and user type
3. If expired, throw CART_EXPIRED error
4. For each cart item, fetch current product price
5. Calculate totals using current prices
6. Flag any items with price changes since addition
7. Return cart with calculated totals

**Error Cases:** [CART_EXPIRED]

#### `addItem`

```typescript
addItem(cartId: UUID, productId: UUID, variantId: UUID | null, quantity: number): Promise<CartItem>
```

Add item to cart with stock and variant validation

**Implementation Steps:**
1. Find cart by ID
2. Find product by ID
3. If product has variants and variantId is null, throw VARIANT_REQUIRED
4. Check inventory availableStock
5. If availableStock == 0, throw OUT_OF_STOCK
6. Cap quantity at availableStock if exceeded
7. Store current price as priceAtAdd
8. Create or update cart item
9. Update cart lastActivityAt
10. Return cart item with actual quantity added

**Error Cases:** [VARIANT_REQUIRED OUT_OF_STOCK QUANTITY_EXCEEDS_STOCK]

#### `updateQuantity`

```typescript
updateQuantity(cartItemId: UUID, quantity: number): Promise<CartItem>
```

Update cart item quantity with stock validation

**Implementation Steps:**
1. Find cart item by ID
2. Check inventory availableStock
3. Cap quantity at availableStock if exceeded
4. Update cart item quantity
5. Update cart lastActivityAt
6. Return updated cart item

**Error Cases:** [QUANTITY_EXCEEDS_STOCK]

#### `removeItem`

```typescript
removeItem(cartItemId: UUID): Promise<void>
```

Remove item from cart

**Implementation Steps:**
1. Find cart item by ID
2. Delete cart item
3. Update cart lastActivityAt
4. Return success

#### `calculateTotal`

```typescript
calculateTotal(cart: Cart): Promise<CartTotals>
```

Calculate cart totals with current prices

**Implementation Steps:**
1. For each cart item, fetch current product/variant price
2. Calculate line totals (price * quantity)
3. Sum subtotal from all line items
4. Calculate shipping using ShippingService
5. Return totals object with subtotal, shipping, total

---

## SKEL-SHIP-006: ShippingService (service)

**Related Specs:** [TS-BR-SHIP-001 TS-BR-SHIP-002 TS-BR-ADDR-001]

### Functions

#### `calculateCost`

```typescript
calculateCost(subtotal: number): Promise<number>
```

Calculate shipping cost based on free shipping threshold

**Implementation Steps:**
1. Check if subtotal >= 50.00
2. If yes, return 0.00 (free shipping)
3. Otherwise, return 5.99

#### `validateAddress`

```typescript
validateAddress(address: ShippingAddressDTO): Promise<void>
```

Validate shipping address for domestic delivery

**Implementation Steps:**
1. Validate all required fields are present and non-empty
2. Validate country is ISO 3166-1 alpha-2 format
3. Check country matches DOMESTIC_COUNTRY config
4. If validation fails, throw appropriate error

**Error Cases:** [INCOMPLETE_ADDRESS INTERNATIONAL_SHIPPING_NOT_AVAILABLE]

---

## SKEL-ORDE-007: OrderService (service)

**Dependencies:** [CustomerService CartService InventoryService PaymentService AuditLogService]

**Related Specs:** [TS-BR-AUTH-001 TS-BR-PRICE-002 TS-BR-CANCEL-001 TS-BR-INV-002 TS-BR-INV-003 TS-BR-ORDER-001 TS-BR-ORDER-002 TS-BR-ORDER-003 TS-BR-AUDIT-001]

### Functions

#### `create`

```typescript
create(customerId: UUID, cartId: UUID, paymentAuthorizationId: string, shippingAddress: ShippingAddress): Promise<Order>
```

Create order from cart with price snapshot

**Implementation Steps:**
1. Find customer by ID
2. Validate customer.emailVerified == true
3. Find cart by ID
4. Validate all cart items have sufficient stock
5. For each cart item, create OrderItem with current price as unitPrice
6. Calculate line totals (unitPrice * quantity)
7. Calculate subtotal, shipping, and total
8. Decrement inventory for all items via InventoryService
9. Create order with status='pending'
10. Store paymentAuthorizationId
11. Clear cart
12. Log order creation to audit via AuditLogService
13. Return created order

**Error Cases:** [REGISTRATION_REQUIRED EMAIL_NOT_VERIFIED INSUFFICIENT_STOCK BACKORDER_NOT_SUPPORTED]

#### `updateStatus`

```typescript
updateStatus(orderId: UUID, newStatus: OrderStatus): Promise<Order>
```

Update order status following state machine rules

**Implementation Steps:**
1. Find order by ID
2. Get current status
3. Validate transition is allowed per state machine
4. Update order status
5. Log status change to audit via AuditLogService
6. Return updated order

**Error Cases:** [INVALID_STATUS_TRANSITION]

#### `cancel`

```typescript
cancel(orderId: UUID): Promise<Order>
```

Cancel order with inventory restoration and refund

**Implementation Steps:**
1. Find order by ID
2. Validate current status is 'pending' or 'confirmed'
3. Update status to 'cancelled'
4. Restore inventory for all order items via InventoryService
5. Initiate refund via PaymentService
6. Log cancellation to audit via AuditLogService
7. Return cancelled order

**Error Cases:** [CANCELLATION_NOT_ALLOWED]

#### `getById`

```typescript
getById(orderId: UUID): Promise<Order>
```

Get order by ID

**Implementation Steps:**
1. Find order by ID
2. Return order with items

---

## SKEL-CHEC-008: CheckoutService (service)

**Dependencies:** [CustomerService CartService InventoryService PaymentService ShippingService OrderService]

**Related Specs:** [TS-BR-AUTH-001 TS-BR-STOCK-001 TS-BR-INV-002 TS-BR-CUST-003 TS-BR-PAY-001]

### Functions

#### `placeOrder`

```typescript
placeOrder(customerId: UUID, cartId: UUID, paymentMethod: PaymentMethodDTO, shippingAddress: ShippingAddressDTO): Promise<Order>
```

Complete checkout flow with validation, payment, and order creation

**Implementation Steps:**
1. Validate customer exists and emailVerified == true
2. Validate cart and all items via validateCart()
3. Validate shipping address via ShippingService
4. Authorize payment via PaymentService
5. If payment fails, throw PAYMENT_REQUIRED
6. Create order via OrderService with payment authorization
7. Return created order

**Error Cases:** [REGISTRATION_REQUIRED EMAIL_NOT_VERIFIED OUT_OF_STOCK INSUFFICIENT_STOCK BACKORDER_NOT_SUPPORTED PAYMENT_REQUIRED INCOMPLETE_ADDRESS INTERNATIONAL_SHIPPING_NOT_AVAILABLE]

#### `validateCart`

```typescript
validateCart(cartId: UUID): Promise<ValidationResult>
```

Validate cart for checkout readiness

**Implementation Steps:**
1. Get cart via CartService
2. For each item, check inventory availableStock
3. If any item has availableStock == 0, flag as out of stock
4. If any item quantity > availableStock, flag as insufficient
5. Return validation result with any issues

**Error Cases:** [OUT_OF_STOCK INSUFFICIENT_STOCK]

#### `validateStock`

```typescript
validateStock(cartItems: CartItem[]): Promise<void>
```

Atomically validate all items have sufficient stock

**Implementation Steps:**
1. For each cart item
2. Get availableStock from InventoryService
3. If quantity > availableStock, throw BACKORDER_NOT_SUPPORTED
4. Return success if all items pass

**Error Cases:** [BACKORDER_NOT_SUPPORTED]

---

## SKEL-PAYM-009: PaymentService (service)

**Related Specs:** [TS-BR-PAY-001 TS-BR-PAY-002]

### Functions

#### `authorize`

```typescript
authorize(amount: number, paymentMethod: PaymentMethodDTO): Promise<PaymentAuthorization>
```

Authorize payment for order total

**Implementation Steps:**
1. Validate payment method details
2. Send authorization request to payment gateway
3. If successful, return authorization ID and status
4. If failed, throw PAYMENT_REQUIRED

**Error Cases:** [PAYMENT_REQUIRED]

#### `capture`

```typescript
capture(authorizationId: string): Promise<PaymentCapture>
```

Capture previously authorized payment

**Implementation Steps:**
1. Send capture request to payment gateway
2. Update payment status to 'captured'
3. Return capture confirmation

#### `refund`

```typescript
refund(authorizationId: string): Promise<PaymentRefund>
```

Initiate refund for cancelled order

**Implementation Steps:**
1. Send refund request to payment gateway
2. Create refund record with status 'pending'
3. Return refund confirmation

---

## SKEL-CATE-010: CategoryService (service)

**Related Specs:** [TS-BR-CAT-001]

### Functions

#### `create`

```typescript
create(data: CreateCategoryDTO): Promise<Category>
```

Create category with depth validation

**Implementation Steps:**
1. If parentId provided, find parent category
2. Calculate depth = parent.depth + 1 (or 1 if no parent)
3. Validate depth <= 3
4. Create category record
5. Return created category

**Error Cases:** [CATEGORY_DEPTH_EXCEEDED]

#### `update`

```typescript
update(categoryId: UUID, data: UpdateCategoryDTO): Promise<Category>
```

Update category details

**Implementation Steps:**
1. Find category by ID
2. Update category fields
3. Return updated category

#### `delete`

```typescript
delete(categoryId: UUID): Promise<void>
```

Delete category and handle children

**Implementation Steps:**
1. Find category by ID
2. Check for child categories
3. Delete category record
4. Return success

---

## SKEL-ADDR-011: AddressService (service)

**Related Specs:** [TS-BR-ADDR-001 TS-BR-SHIP-002]

### Functions

#### `create`

```typescript
create(customerId: UUID, data: CreateAddressDTO): Promise<ShippingAddress>
```

Create shipping address with validation

**Implementation Steps:**
1. Validate all required fields present and non-empty
2. Validate country is ISO 3166-1 alpha-2
3. Validate country matches DOMESTIC_COUNTRY
4. Create address record linked to customer
5. Return created address

**Error Cases:** [INCOMPLETE_ADDRESS INTERNATIONAL_SHIPPING_NOT_AVAILABLE]

#### `update`

```typescript
update(addressId: UUID, data: UpdateAddressDTO): Promise<ShippingAddress>
```

Update shipping address with validation

**Implementation Steps:**
1. Find address by ID
2. Validate all required fields if provided
3. Validate country if provided
4. Update address record
5. Return updated address

**Error Cases:** [INCOMPLETE_ADDRESS INTERNATIONAL_SHIPPING_NOT_AVAILABLE]

#### `validate`

```typescript
validate(address: ShippingAddressDTO): Promise<void>
```

Validate address fields and domestic shipping

**Implementation Steps:**
1. Check recipientName is non-empty
2. Check street is non-empty
3. Check city is non-empty
4. Check state is non-empty
5. Check postalCode is non-empty
6. Check country is valid ISO code
7. Check country matches DOMESTIC_COUNTRY
8. Throw appropriate error if any validation fails

**Error Cases:** [INCOMPLETE_ADDRESS INTERNATIONAL_SHIPPING_NOT_AVAILABLE]

---

## SKEL-AUDI-012: AuditLogService (service)

**Related Specs:** [TS-BR-AUDIT-001 TS-BR-AUDIT-002 TS-BR-AUDIT-003]

### Functions

#### `log`

```typescript
log(entityType: string, entityId: UUID, action: string, beforeValue: object | null, afterValue: object, userId: UUID, reason?: string): Promise<AuditLog>
```

Create audit log entry for entity changes

**Implementation Steps:**
1. Create audit log record with timestamp
2. Store entityType, entityId, action
3. Store beforeValue and afterValue as JSONB
4. Store userId and optional reason
5. Return created audit log entry

#### `getByEntity`

```typescript
getByEntity(entityType: string, entityId: UUID): Promise<AuditLog[]>
```

Get audit history for an entity

**Implementation Steps:**
1. Query audit logs by entityType and entityId
2. Order by timestamp descending
3. Return list of audit entries

---

## SKEL-CART-013: CartExpirationJob (service)

**Dependencies:** [CartService]

**Related Specs:** [TS-BR-CART-003]

### Functions

#### `run`

```typescript
run(): Promise<void>
```

Scheduled job to expire inactive carts

**Implementation Steps:**
1. Calculate cutoff timestamp for authenticated carts (now - 30 days)
2. Calculate cutoff timestamp for guest carts (now - 7 days)
3. Delete guest carts where lastActivityAt < guest cutoff
4. Delete authenticated carts where lastActivityAt < auth cutoff
5. Log number of expired carts

---

