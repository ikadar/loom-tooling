# Aggregate Design

Generated: 2025-12-29T15:49:16+01:00

---

## AGG-CUSTOMER-001 – Customer

**Purpose:** Manages customer identity, authentication, and associated shipping addresses

### Aggregate Root: Customer

**Identity:** CustomerId (UUID)

**Attributes:**
| Name | Type | Mutable |
|------|------|----------|
| customerId | CustomerId | false |
| email | Email | true |
| passwordHash | string | true |
| registrationStatus | RegistrationStatus | true |
| emailVerified | boolean | true |
| firstName | string | true |
| lastName | string | true |
| shippingAddresses | List<ShippingAddress> | true |
| createdAt | DateTime | false |

### Invariants

- **INV-CUST-001**: Customer email must be unique across all customers
  - Enforcement: Database unique constraint, validated on registration
- **INV-CUST-002**: Customer must be registered and email verified to place orders
  - Enforcement: Status check in placeOrder operation
- **INV-CUST-003**: Password must be at least 8 characters with at least one number
  - Enforcement: Validation during registration and password change
- **INV-CUST-004**: Email must be valid RFC 5322 format
  - Enforcement: Email value object validation on creation

### Child Entities

#### ShippingAddress

**Identity:** AddressId

**Purpose:** Customer's saved shipping address for order delivery

### Value Objects

[Email Password RegistrationStatus Address]

### Behaviors

| Command | Pre | Post | Emits |
|---------|-----|------|-------|
| RegisterCustomer | [email not already registered] | [status == 'unverified' verification email queued] | CustomerRegistered |
| VerifyCustomerEmail | [valid verification token token not expired] | [emailVerified == true] | CustomerEmailVerified |
| AddShippingAddress | [customer registered address fields valid country is domestic] | [address added to shippingAddresses] | ShippingAddressAdded |
| UpdateShippingAddress | [address exists address fields valid] | [address updated] | ShippingAddressUpdated |
| RemoveShippingAddress | [address exists] | [address removed from shippingAddresses] | ShippingAddressRemoved |
| RequestGDPRErasure | [no pending orders] | [personal data erased order history anonymized] | CustomerDataErased |

### Events

- **CustomerRegistered**: [customerId email createdAt]
- **CustomerEmailVerified**: [customerId verifiedAt]
- **ShippingAddressAdded**: [customerId addressId address]
- **ShippingAddressUpdated**: [customerId addressId address]
- **ShippingAddressRemoved**: [customerId addressId]
- **CustomerDataErased**: [anonymizedId erasedAt]

### Repository: CustomerRepository

- Load Strategy: Load customer with all shipping addresses
- Concurrency: Optimistic locking via version field

---

## AGG-CART-001 – Cart

**Purpose:** Manages shopping cart items and totals before checkout

### Aggregate Root: Cart

**Identity:** CartId (UUID)

**Attributes:**
| Name | Type | Mutable |
|------|------|----------|
| cartId | CartId | false |
| customerId | CustomerId | false |
| isGuest | boolean | false |
| items | List<CartItem> | true |
| totalPrice | Money | true |
| lastActivityDate | DateTime | true |
| createdAt | DateTime | false |

### Invariants

- **INV-CART-001**: Total price must equal sum of all cart item subtotals
  - Enforcement: Recalculated after any item change
- **INV-CART-002**: Each product can appear only once (quantity adjusted instead)
  - Enforcement: Check for existing item before add, merge if exists
- **INV-CART-003**: Cart item quantity must be at least 1
  - Enforcement: Setting to 0 triggers removal
- **INV-CART-004**: Cart item quantity cannot exceed available stock
  - Enforcement: Stock check during add and update, auto-cap if exceeded
- **INV-CART-005**: Products with variants require specific variant selection
  - Enforcement: Validation on addItem when product hasVariants
- **INV-CART-006**: Inactive carts expire after threshold (30 days logged-in, 7 days guest)
  - Enforcement: Background job or lazy evaluation on cart access

### Child Entities

#### CartItem

**Identity:** CartItemId

**Purpose:** Product and quantity snapshot within shopping cart

### Value Objects

[Money Quantity]

### Behaviors

| Command | Pre | Post | Emits |
|---------|-----|------|-------|
| AddItemToCart | [product is active product is in stock quantity > 0 variant selected if required] | [item added or quantity merged totalPrice recalculated lastActivityDate updated] | ItemAddedToCart |
| UpdateCartItemQuantity | [item exists in cart quantity >= 0 quantity <= available stock] | [quantity updated or item removed if 0 totalPrice recalculated lastActivityDate updated] | CartItemQuantityUpdated |
| RemoveItemFromCart | [item exists in cart] | [item removed totalPrice recalculated lastActivityDate updated] | ItemRemovedFromCart |
| ClearCart | [] | [all items removed totalPrice == 0] | CartCleared |
| MergeGuestCart | [guest cart exists customer authenticated] | [guest cart items merged quantities limited to stock] | CartMerged |
| RefreshCartPrices | [] | [prices updated to current price changes tracked] | CartPricesRefreshed |

### Events

- **ItemAddedToCart**: [cartId productId variantId quantity unitPrice]
- **CartItemQuantityUpdated**: [cartId cartItemId oldQuantity newQuantity]
- **ItemRemovedFromCart**: [cartId productId variantId]
- **CartCleared**: [cartId reason]
- **CartMerged**: [cartId guestCartId mergedItems]
- **CartPricesRefreshed**: [cartId priceChanges]
- **CartExpired**: [cartId customerId]

### Repository: CartRepository

- Load Strategy: Load cart with all items eagerly
- Concurrency: Optimistic locking via version field

---

## AGG-ORDER-001 – Order

**Purpose:** Manages the complete order lifecycle from placement through delivery or cancellation

### Aggregate Root: Order

**Identity:** OrderId (UUID)

**Attributes:**
| Name | Type | Mutable |
|------|------|----------|
| orderId | OrderId | false |
| orderNumber | string | false |
| customerId | CustomerId | false |
| status | OrderStatus | true |
| lineItems | List<OrderLineItem> | false |
| shippingAddress | Address | false |
| paymentMethod | PaymentMethod | false |
| paymentTransactionId | string | false |
| subtotal | Money | false |
| shippingCost | Money | false |
| tax | Money | false |
| totalAmount | Money | false |
| trackingNumber | string | true |
| cancellationReason | string | true |
| createdAt | DateTime | false |
| updatedAt | DateTime | true |

### Invariants

- **INV-ORDER-001**: Order must have at least one line item
  - Enforcement: Validated on creation
- **INV-ORDER-002**: Total amount must equal subtotal plus shipping cost
  - Enforcement: Calculated on creation, immutable
- **INV-ORDER-003**: Shipping is free when subtotal exceeds $50
  - Enforcement: Shipping cost calculation on creation
- **INV-ORDER-004**: Cannot modify order after status is shipped, delivered, or cancelled
  - Enforcement: Guard clause in all modification methods
- **INV-ORDER-005**: Can only cancel before status is shipped
  - Enforcement: Status check in cancel operation
- **INV-ORDER-006**: Status transitions must follow valid state machine
  - Enforcement: State pattern: pending→confirmed→shipped→delivered OR pending/confirmed→cancelled
- **INV-ORDER-007**: Line item prices are immutable snapshots from order time
  - Enforcement: Copy price during order creation, no updates allowed
- **INV-ORDER-008**: Order number must be unique with format ORD-{YEAR}-{SEQUENCE}
  - Enforcement: Generated atomically on creation

### Child Entities

#### OrderLineItem

**Identity:** LineItemId

**Purpose:** Immutable snapshot of product at order time with quantity and price

### Value Objects

[Money Address PaymentMethod OrderStatus Quantity]

### Behaviors

| Command | Pre | Post | Emits |
|---------|-----|------|-------|
| PlaceOrder | [customer registered and verified cart has items all items in stock payment authorized valid shipping address] | [order created with status 'pending' line items snapshot created totals calculated] | OrderPlaced |
| ConfirmOrder | [status == 'pending'] | [status == 'confirmed'] | OrderConfirmed |
| ShipOrder | [status == 'confirmed' tracking number provided] | [status == 'shipped' tracking number stored] | OrderShipped |
| DeliverOrder | [status == 'shipped'] | [status == 'delivered'] | OrderDelivered |
| CancelOrder | [status in ['pending', 'confirmed']] | [status == 'cancelled' refund initiated inventory restoration triggered] | OrderCancelled |

### Events

- **OrderPlaced**: [orderId orderNumber customerId lineItems totalAmount shippingAddress]
- **OrderConfirmed**: [orderId orderNumber confirmedAt]
- **OrderShipped**: [orderId orderNumber trackingNumber shippedAt]
- **OrderDelivered**: [orderId orderNumber deliveredAt]
- **OrderCancelled**: [orderId orderNumber reason cancelledAt lineItems]

### Repository: OrderRepository

- Load Strategy: Load order with all line items eagerly
- Concurrency: Optimistic locking via version field

---

## AGG-PRODUCT-001 – Product

**Purpose:** Manages product catalog information including variants and images

### Aggregate Root: Product

**Identity:** ProductId (UUID)

**Attributes:**
| Name | Type | Mutable |
|------|------|----------|
| productId | ProductId | false |
| name | string | true |
| description | string | true |
| price | Money | true |
| categoryId | CategoryId | true |
| variants | List<ProductVariant> | true |
| images | List<ProductImage> | true |
| isActive | boolean | true |
| isDeleted | boolean | true |
| status | ProductStatus | true |
| createdAt | DateTime | false |
| createdBy | UserId | false |
| updatedAt | DateTime | true |
| updatedBy | UserId | true |

### Invariants

- **INV-PROD-001**: Product price must be at least $0.01
  - Enforcement: Validation on creation and update
- **INV-PROD-002**: Product name must be 2-200 characters
  - Enforcement: Validation on creation and update
- **INV-PROD-003**: Product must belong to a category
  - Enforcement: Required categoryId on creation
- **INV-PROD-004**: Inactive products cannot be added to cart
  - Enforcement: isActive check in cart add operation
- **INV-PROD-005**: Variant SKU must be unique across all products
  - Enforcement: Database unique constraint
- **INV-PROD-006**: Variant combination must be unique within product
  - Enforcement: Check before adding variant
- **INV-PROD-007**: Products with order history must be soft-deleted
  - Enforcement: Check for OrderLineItem references before delete
- **INV-PROD-008**: Maximum 10 images per product
  - Enforcement: Validation on image upload

### Child Entities

#### ProductVariant

**Identity:** VariantId

**Purpose:** Specific variation of product (size, color combination)

#### ProductImage

**Identity:** ImageId

**Purpose:** Product image with display order

### Value Objects

[Money SKU ProductStatus]

### Behaviors

| Command | Pre | Post | Emits |
|---------|-----|------|-------|
| CreateProduct | [name is valid price >= 0.01 category exists] | [product created with status 'draft' UUID assigned] | ProductCreated |
| UpdateProduct | [product exists fields valid] | [attributes updated change logged] | ProductUpdated |
| ActivateProduct | [product exists isActive == false] | [isActive == true status == 'active'] | ProductActivated |
| DeactivateProduct | [product exists isActive == true] | [isActive == false] | ProductDeactivated |
| DeleteProduct | [product exists] | [isDeleted == true cart items with product removed] | ProductDeleted |
| AddProductVariant | [product exists SKU unique variant combination unique] | [variant added to product] | ProductVariantAdded |
| RemoveProductVariant | [variant exists no active cart/order references] | [variant removed] | ProductVariantRemoved |
| UploadProductImages | [product exists total images <= 10 valid image format] | [images stored URLs saved primary set] | ProductImagesUploaded |

### Events

- **ProductCreated**: [productId name price categoryId createdBy]
- **ProductUpdated**: [productId changedFields updatedBy]
- **ProductActivated**: [productId activatedAt]
- **ProductDeactivated**: [productId deactivatedAt]
- **ProductDeleted**: [productId deletedAt deletedBy]
- **ProductVariantAdded**: [productId variantId sku]
- **ProductVariantRemoved**: [productId variantId]
- **ProductImagesUploaded**: [productId imageUrls primaryImageUrl]

### Repository: ProductRepository

- Load Strategy: Load product with variants and images eagerly
- Concurrency: Optimistic locking via version field

---

## AGG-INVENTORY-001 – Inventory

**Purpose:** Tracks stock levels and manages inventory reservations with full audit trail

### Aggregate Root: Inventory

**Identity:** InventoryId (UUID)

**Attributes:**
| Name | Type | Mutable |
|------|------|----------|
| inventoryId | InventoryId | false |
| productId | ProductId | false |
| stockLevel | integer | true |
| reservedQuantity | integer | true |
| lowStockThreshold | integer | true |
| updatedAt | DateTime | true |

### Invariants

- **INV-INV-001**: Stock level cannot be negative
  - Enforcement: Validation on all stock adjustments
- **INV-INV-002**: Reserved quantity cannot exceed stock level
  - Enforcement: Check in reserve operation
- **INV-INV-003**: Available quantity must equal stock level minus reserved
  - Enforcement: Calculated property
- **INV-INV-004**: All inventory changes must be logged with before/after values and reason
  - Enforcement: Audit log entry in same transaction as change
- **INV-INV-005**: Low stock threshold cannot be negative
  - Enforcement: Validation on set

### Child Entities

#### InventoryAuditLog

**Identity:** AuditLogId

**Purpose:** Immutable record of inventory changes

### Behaviors

| Command | Pre | Post | Emits |
|---------|-----|------|-------|
| SetInventoryQuantity | [quantity >= 0 reason provided] | [stockLevel set audit log created low stock alert if below threshold] | InventoryQuantitySet |
| AdjustInventory | [result >= 0 reason provided] | [stockLevel adjusted audit log created low stock alert if below threshold] | InventoryAdjusted |
| ReserveInventory | [quantity <= availableQuantity] | [reservedQuantity increased] | InventoryReserved |
| ReleaseInventory | [quantity <= reservedQuantity] | [reservedQuantity decreased] | InventoryReleased |
| DeductInventory | [quantity <= reservedQuantity] | [stockLevel decreased reservedQuantity decreased audit log created] | InventoryDeducted |
| RestockInventory | [quantity > 0 reason provided] | [stockLevel increased audit log created] | InventoryRestocked |
| SetLowStockThreshold | [threshold >= 0] | [lowStockThreshold set] | LowStockThresholdSet |

### Events

- **InventoryQuantitySet**: [inventoryId productId previousValue newValue reason userId]
- **InventoryAdjusted**: [inventoryId productId adjustment newValue reason userId]
- **InventoryReserved**: [inventoryId productId quantity orderId]
- **InventoryReleased**: [inventoryId productId quantity reason]
- **InventoryDeducted**: [inventoryId productId quantity orderId]
- **InventoryRestocked**: [inventoryId productId quantity reason userId]
- **LowStockAlert**: [inventoryId productId currentStock threshold]
- **OutOfStock**: [inventoryId productId]
- **LowStockThresholdSet**: [inventoryId productId threshold]

### Repository: InventoryRepository

- Load Strategy: Load inventory without audit log by default, load audit log on demand
- Concurrency: Pessimistic locking for reserve/release operations

---

## AGG-CATEGORY-001 – Category

**Purpose:** Organizes products into hierarchical browsing structure

### Aggregate Root: Category

**Identity:** CategoryId (UUID)

**Attributes:**
| Name | Type | Mutable |
|------|------|----------|
| categoryId | CategoryId | false |
| name | string | true |
| description | string | true |
| parentCategoryId | CategoryId | true |
| depth | integer | true |
| displayOrder | integer | true |
| isActive | boolean | true |
| createdAt | DateTime | false |

### Invariants

- **INV-CAT-001**: Category name must be unique
  - Enforcement: Database unique constraint
- **INV-CAT-002**: Category cannot be its own parent (no circular references)
  - Enforcement: Validation on parent assignment
- **INV-CAT-003**: Category hierarchy limited to 3 levels deep
  - Enforcement: Depth calculation and validation on create/move
- **INV-CAT-004**: Category name must be 1-100 characters
  - Enforcement: Validation on creation and update

### Behaviors

| Command | Pre | Post | Emits |
|---------|-----|------|-------|
| CreateCategory | [name unique name 1-100 chars parent exists if specified resulting depth <= 3] | [category created depth calculated] | CategoryCreated |
| UpdateCategory | [category exists name unique if changed] | [attributes updated] | CategoryUpdated |
| MoveCategory | [category exists new parent exists no circular reference resulting depth <= 3] | [parentCategoryId updated depth recalculated for subtree] | CategoryMoved |
| DeactivateCategory | [category exists no active products in category] | [isActive == false] | CategoryDeactivated |

### Events

- **CategoryCreated**: [categoryId name parentCategoryId depth]
- **CategoryUpdated**: [categoryId changedFields]
- **CategoryMoved**: [categoryId oldParentId newParentId newDepth]
- **CategoryDeactivated**: [categoryId deactivatedAt]

### Repository: CategoryRepository

- Load Strategy: Load category without children, load hierarchy on demand
- Concurrency: Optimistic locking via version field

---

