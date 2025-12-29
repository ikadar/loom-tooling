# Implementation Skeletons

Generated: 2025-12-29T16:22:00+01:00

---

## CustomerService (service)

**Dependencies:** [EmailService CustomerRepository]

**Related Specs:** [TS-BR-CUST-001 TS-BR-CUST-002 TS-BR-CUST-003 TS-BR-SHIP-001 TS-BR-SHIP-002]

### Functions

#### `register`

```typescript
register(email: string, password: string, firstName: string, lastName: string): Customer
```

Register new customer with email uniqueness check

**Implementation Steps:**
1. Validate email format
2. Check email uniqueness in database
3. Validate password meets requirements (8+ chars, 1 number)
4. Hash password
5. Create customer with unverified status
6. Queue verification email

**Error Cases:** [DUPLICATE_EMAIL INVALID_PASSWORD]

#### `addAddress`

```typescript
addAddress(customerId: UUID, address: AddressInput): Address
```

Add shipping address with validation

**Implementation Steps:**
1. Validate all required fields present
2. Validate postal code format
3. Validate country is domestic
4. Save address

**Error Cases:** [INVALID_ADDRESS INVALID_POSTAL_CODE UNSUPPORTED_SHIPPING_REGION]

#### `requestDataErasure`

```typescript
requestDataErasure(customerId: UUID): void
```

GDPR data erasure with pending order check

**Implementation Steps:**
1. Check for pending orders
2. Delete personal data
3. Anonymize order history
4. Send confirmation email

**Error Cases:** [PENDING_ORDERS_EXIST]

#### `verifyOrderEligibility`

```typescript
verifyOrderEligibility(customerId: UUID): boolean
```

Check if customer can place orders

**Implementation Steps:**
1. Check registration status
2. Check email verification status

**Error Cases:** [UNAUTHORIZED_ORDER EMAIL_NOT_VERIFIED]

---

## ProductService (service)

**Dependencies:** [ProductRepository StorageService InventoryService]

**Related Specs:** [TS-BR-PROD-001 TS-BR-PROD-002 TS-BR-PROD-003 TS-BR-PROD-004]

### Functions

#### `listProducts`

```typescript
listProducts(filters: ProductFilters, page: number, limit: number): ProductList
```

Filter and paginate products

**Implementation Steps:**
1. Build query from filters
2. Apply pagination
3. Return products with total count

#### `createProduct`

```typescript
createProduct(input: ProductInput): Product
```

Create new product as draft

**Implementation Steps:**
1. Validate name length (2-200 chars)
2. Validate price (>= 0.01)
3. Create product with draft status
4. Generate UUID

**Error Cases:** [INVALID_PRICE NAME_REQUIRED NAME_TOO_LONG]

#### `updateProduct`

```typescript
updateProduct(productId: UUID, input: ProductInput): Product
```

Update product with conflict detection

**Implementation Steps:**
1. Fetch current product
2. Validate input fields
3. Check for concurrent edits
4. Update product
5. Log change

**Error Cases:** [INVALID_PRICE NAME_REQUIRED NAME_TOO_LONG]

#### `deleteProduct`

```typescript
deleteProduct(productId: UUID): void
```

Soft delete product if has order history

**Implementation Steps:**
1. Check for order line item references
2. If has references: soft delete (set isDeleted, deletedAt)
3. If no references: allow hard delete

#### `uploadImages`

```typescript
uploadImages(productId: UUID, images: File[], primaryIndex: number): ImageList
```

Upload product images with limit check

**Implementation Steps:**
1. Check current image count
2. Validate total <= 10
3. Upload to storage
4. Set primary image

**Error Cases:** [MAX_IMAGES_EXCEEDED]

---

## CartService (service)

**Dependencies:** [CartRepository InventoryService ProductService]

**Related Specs:** [TS-BR-CART-001 TS-BR-CART-002 TS-BR-CART-003 TS-BR-CART-004 TS-BR-CART-005]

### Functions

#### `addItem`

```typescript
addItem(cartId: UUID, productId: UUID, variantId: UUID | null, quantity: number): CartItem
```

Add item to cart with stock validation

**Implementation Steps:**
1. Check if product has variants and variant required
2. Validate stock availability (must be > 0)
3. Check for existing cart item
4. If exists: increase quantity
5. If new: create cart item
6. Cap quantity to available stock
7. Update cart lastActivityDate

**Error Cases:** [OUT_OF_STOCK VARIANT_REQUIRED QUANTITY_LIMITED]

#### `updateItemQuantity`

```typescript
updateItemQuantity(cartId: UUID, itemId: UUID, quantity: number): CartItem | null
```

Update cart item quantity or remove if zero

**Implementation Steps:**
1. If quantity <= 0: remove item
2. Check available stock
3. Cap to available if needed
4. Update quantity
5. Recalculate totals

**Error Cases:** [QUANTITY_LIMITED]

#### `removeItem`

```typescript
removeItem(cartId: UUID, itemId: UUID): void
```

Remove item from cart

**Implementation Steps:**
1. Delete cart item
2. Return undo token (5 second window)

#### `mergeCarts`

```typescript
mergeCarts(guestCartId: UUID, customerCartId: UUID): Cart
```

Merge guest cart into customer cart on login

**Implementation Steps:**
1. Get both carts
2. For each guest item: add to customer cart or increase quantity
3. Apply stock limits
4. Delete guest cart

#### `validateCart`

```typescript
validateCart(cartId: UUID): CartValidationResult
```

Validate cart items for checkout

**Implementation Steps:**
1. Check each item's current stock
2. Check for price changes
3. Mark invalid items
4. Return validation result

**Error Cases:** [OUT_OF_STOCK]

#### `expireInactiveCarts`

```typescript
expireInactiveCarts(): number
```

Background job to expire old carts

**Implementation Steps:**
1. Find carts where lastActivityDate > threshold
2. 30 days for logged-in, 7 days for guest
3. Delete expired carts
4. Return count deleted

---

## OrderService (service)

**Dependencies:** [CustomerService CartService InventoryService PaymentService EmailService OrderRepository]

**Related Specs:** [TS-BR-CUST-001 TS-BR-ORDER-001 TS-BR-ORDER-002 TS-BR-ORDER-003 TS-BR-ORDER-004 TS-BR-ORDER-005 TS-BR-ORDER-006 TS-BR-INV-002 TS-BR-INV-003]

### Functions

#### `placeOrder`

```typescript
placeOrder(customerId: UUID, shippingAddressId: UUID, paymentMethodId: UUID): Order
```

Create order with all validations

**Implementation Steps:**
1. Verify customer order eligibility
2. Get and validate cart (non-empty, stock available)
3. Calculate totals (subtotal, shipping, tax)
4. Authorize payment
5. Create order with price snapshots
6. Generate order number
7. Decrement inventory
8. Clear cart
9. Queue confirmation email

**Error Cases:** [EMAIL_NOT_VERIFIED EMPTY_ORDER STOCK_UNAVAILABLE PAYMENT_FAILED]

#### `cancelOrder`

```typescript
cancelOrder(orderId: UUID, reason: string): Order
```

Cancel order with refund and inventory restore

**Implementation Steps:**
1. Verify order status is cancellable (pending or confirmed)
2. Update status to cancelled
3. Restore inventory for all line items
4. Initiate refund
5. Send cancellation email

**Error Cases:** [ORDER_ALREADY_SHIPPED ORDER_NOT_MODIFIABLE]

#### `updateStatus`

```typescript
updateStatus(orderId: UUID, newStatus: OrderStatus, trackingNumber: string | null): Order
```

Update order status with transition validation

**Implementation Steps:**
1. Fetch order
2. Validate status transition using state machine
3. Update status
4. If shipped: store tracking number
5. Send status notification email

**Error Cases:** [INVALID_STATUS_TRANSITION]

#### `getOrderHistory`

```typescript
getOrderHistory(customerId: UUID, page: number, limit: number): OrderList
```

Get paginated order history

**Implementation Steps:**
1. Query orders by customer
2. Sort by date descending
3. Apply pagination
4. Include status and basic details

#### `calculateTotals`

```typescript
calculateTotals(cart: Cart, shippingAddress: Address): OrderTotals
```

Calculate order totals including shipping and tax

**Implementation Steps:**
1. Sum line item prices for subtotal
2. Apply free shipping if subtotal > 50
3. Calculate tax based on shipping address
4. Return totals breakdown

#### `generateOrderNumber`

```typescript
generateOrderNumber(): string
```

Generate sequential order number

**Implementation Steps:**
1. Get current year
2. Get last order number for year
3. Increment sequence
4. Format as ORD-YYYY-NNNNN

---

## InventoryService (service)

**Dependencies:** [InventoryRepository InventoryAuditRepository NotificationService]

**Related Specs:** [TS-BR-INV-001 TS-BR-INV-002 TS-BR-INV-003 TS-BR-INV-004]

### Functions

#### `getAvailableQuantity`

```typescript
getAvailableQuantity(productId: UUID, variantId: UUID | null): number
```

Get current available stock

**Implementation Steps:**
1. Query inventory table
2. Return stockLevel

#### `setQuantity`

```typescript
setQuantity(productId: UUID, quantity: number, reason: string, userId: UUID): void
```

Set absolute inventory quantity

**Implementation Steps:**
1. Get current stock level
2. Validate quantity >= 0
3. Update stock level
4. Create audit log entry

#### `adjustQuantity`

```typescript
adjustQuantity(productId: UUID, adjustment: number, reason: string, userId: UUID): void
```

Adjust inventory by delta amount

**Implementation Steps:**
1. Get current stock level
2. Calculate new level
3. Validate new level >= 0
4. Update stock level
5. Create audit log entry
6. Check low stock alert threshold

**Error Cases:** [INSUFFICIENT_STOCK]

#### `decrementForOrder`

```typescript
decrementForOrder(lineItems: OrderLineItem[]): void
```

Atomically decrement stock for order

**Implementation Steps:**
1. Start transaction
2. For each line item: decrement by quantity
3. Verify no negative stock
4. Commit transaction

**Error Cases:** [INSUFFICIENT_STOCK]

#### `restoreForCancellation`

```typescript
restoreForCancellation(orderId: UUID): void
```

Restore inventory when order cancelled

**Implementation Steps:**
1. Get order line items
2. For each item: increment stock by quantity
3. Create audit log entries

#### `bulkUpdate`

```typescript
bulkUpdate(updates: InventoryUpdate[]): BulkUpdateResult
```

Process CSV bulk update

**Implementation Steps:**
1. Parse CSV file
2. Validate SKUs exist
3. Apply valid updates
4. Collect errors for invalid rows
5. Return summary

#### `checkLowStockAlerts`

```typescript
checkLowStockAlerts(): LowStockAlert[]
```

Check products below threshold

**Implementation Steps:**
1. Query products where stock < threshold
2. Generate alerts
3. Send notifications

---

## PaymentService (service)

**Dependencies:** [PaymentGateway PaymentMethodRepository]

**Related Specs:** [TS-BR-PAY-001 TS-BR-PAY-002 TS-BR-PAY-003]

### Functions

#### `authorizePayment`

```typescript
authorizePayment(paymentMethodId: UUID, amount: number): PaymentAuthorization
```

Authorize payment before order creation

**Implementation Steps:**
1. Get payment method details
2. Call payment gateway authorization API
3. Store authorization ID
4. Return authorization result

**Error Cases:** [PAYMENT_FAILED PAYMENT_DECLINED]

#### `capturePayment`

```typescript
capturePayment(authorizationId: string, amount: number): PaymentCapture
```

Capture authorized payment

**Implementation Steps:**
1. Call payment gateway capture API
2. Update payment status
3. Return capture result

**Error Cases:** [CAPTURE_FAILED]

#### `refundPayment`

```typescript
refundPayment(orderId: UUID): Refund
```

Refund payment for cancelled order

**Implementation Steps:**
1. Get original transaction ID
2. Call payment gateway refund API
3. Store refund record
4. Return refund result

**Error Cases:** [REFUND_FAILED]

#### `savePaymentMethod`

```typescript
savePaymentMethod(customerId: UUID, cardDetails: CardDetails, saveForFuture: boolean): PaymentMethod
```

Tokenize and save payment method

**Implementation Steps:**
1. Detect card type from BIN
2. Validate card type accepted (Visa, MC, Amex)
3. Tokenize with payment gateway
4. If saveForFuture: store token
5. Return payment method

**Error Cases:** [UNSUPPORTED_CARD_TYPE]

#### `initiatePayPal`

```typescript
initiatePayPal(amount: number, returnUrl: string, cancelUrl: string): PayPalSession
```

Start PayPal authorization flow

**Implementation Steps:**
1. Create PayPal order
2. Get approval URL
3. Return redirect URL

#### `handlePayPalCallback`

```typescript
handlePayPalCallback(token: string, payerId: string | null): PayPalResult
```

Process PayPal callback

**Implementation Steps:**
1. If payerId present: capture payment
2. If cancelled: mark as cancelled
3. Return result

**Error Cases:** [PAYMENT_CANCELLED]

---

## CategoryService (service)

**Dependencies:** [CategoryRepository]

**Related Specs:** [TS-BR-CAT-001]

### Functions

#### `createCategory`

```typescript
createCategory(name: string, parentId: UUID | null): Category
```

Create category with nesting validation

**Implementation Steps:**
1. If parentId: calculate depth by traversing parent chain
2. Validate depth <= 3
3. Create category

**Error Cases:** [MAX_NESTING_EXCEEDED]

#### `getCategoryTree`

```typescript
getCategoryTree(): CategoryNode[]
```

Get full category hierarchy

**Implementation Steps:**
1. Query all categories
2. Build tree structure
3. Return root nodes with children

---

## EmailService (service)

**Dependencies:** [MessageQueue EmailProvider EmailSentRepository]

**Related Specs:** [TS-BR-EMAIL-001 TS-BR-EMAIL-002]

### Functions

#### `queueOrderConfirmation`

```typescript
queueOrderConfirmation(orderId: UUID): void
```

Queue order confirmation email async

**Implementation Steps:**
1. Create email event
2. Publish to message queue
3. Return immediately (non-blocking)

#### `processEmailQueue`

```typescript
processEmailQueue(): void
```

Background worker to process email queue

**Implementation Steps:**
1. Consume message from queue
2. Check deduplication table
3. If not sent: send email
4. Mark as sent in deduplication table
5. Retry on failure

#### `sendEmail`

```typescript
sendEmail(type: EmailType, recipient: string, data: object): void
```

Send email via provider

**Implementation Steps:**
1. Render email template
2. Call email provider API
3. Log result

---

## AdminOrderService (service)

**Dependencies:** [OrderRepository]

**Related Specs:** [TS-BR-ORDER-003]

### Functions

#### `listOrders`

```typescript
listOrders(filters: OrderFilters, page: number, limit: number): OrderList
```

List orders with admin filters

**Implementation Steps:**
1. Build query from filters (status, date range)
2. Apply pagination
3. Return orders with totals

#### `exportOrdersCsv`

```typescript
exportOrdersCsv(filters: OrderFilters): Buffer
```

Export filtered orders to CSV

**Implementation Steps:**
1. Query orders matching filters
2. Format as CSV
3. Return file buffer

---

