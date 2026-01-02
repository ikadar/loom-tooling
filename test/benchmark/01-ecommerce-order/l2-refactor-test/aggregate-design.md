# Aggregate Design

Generated: 2026-01-02T21:09:41+01:00

---

## AGG-CUSTOMER-001 – Customer {#agg-customer-001}

**Purpose:** Manages customer identity, registration status, and profile for order placement eligibility

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
| createdAt | DateTime | false |

### Invariants

- **INV-CUST-001**: Customer must have a valid email address
  - Enforcement: Email value object validates format on creation; registration rejects invalid emails
- **INV-CUST-002**: Email addresses must be unique across all customer accounts
  - Enforcement: Registration validates email uniqueness via repository lookup before account creation
- **INV-CUST-003**: Password must meet minimum security requirements
  - Enforcement: Password value object validates strength on creation and change operations
- **INV-CUST-004**: Customer must be registered and email verified to place orders
  - Enforcement: Order placement checks registration status and emailVerified flag before proceeding

### Value Objects

[CustomerId Email Password RegistrationStatus]

### Behaviors

| Command | Pre | Post | Emits |
|---------|-----|------|-------|
| RegisterCustomer | [Email not already registered Password meets strength requirements] | [Customer created with REGISTERED status emailVerified set to false] | CustomerRegistered |
| VerifyCustomerEmail | [Valid verification token Customer exists] | [emailVerified set to true] | CustomerEmailVerified |
| ChangeCustomerEmail | [New email not already registered New email is valid format] | [Email updated emailVerified reset to false] | CustomerEmailChanged |
| ChangeCustomerPassword | [New password meets strength requirements] | [Password hash updated] | CustomerPasswordChanged |

### Events

- **CustomerRegistered**: [customerId email registeredAt]
- **CustomerEmailVerified**: [customerId verifiedAt]
- **CustomerEmailChanged**: [customerId oldEmail newEmail]
- **CustomerPasswordChanged**: [customerId changedAt]

### Repository: CustomerRepository

- Load Strategy: Eager load all attributes; no child entities
- Concurrency: Optimistic locking with version field

---

## AGG-CART-001 – Cart {#agg-cart-001}

**Purpose:** Holds products a customer intends to purchase with current pricing before checkout

### Aggregate Root: Cart

**Identity:** CartId (UUID)

**Attributes:**
| Name | Type | Mutable |
|------|------|----------|
| cartId | CartId | false |
| customerId | CustomerId | false |
| totalPrice | Money | true |
| lastActivityAt | DateTime | true |

### Invariants

- **INV-CART-001**: Total price equals sum of all item subtotals
  - Enforcement: Total recalculated on every item add/update/remove operation
- **INV-CART-002**: Each product can appear only once in cart
  - Enforcement: AddItem checks for existing item and updates quantity instead of adding duplicate
- **INV-CART-003**: Cart item quantity must be at least 1
  - Enforcement: Quantity validation on add and update; remove item if quantity would be 0
- **INV-CART-004**: Cart item quantity cannot exceed available stock
  - Enforcement: Add and update operations query inventory and cap quantity at available stock
- **INV-CART-005**: Cannot add inactive or out-of-stock products
  - Enforcement: AddItem validates product is active and has available stock before adding

### Child Entities

#### CartItem

**Identity:** CartItemId (UUID)

**Purpose:** Represents a product and quantity within the shopping cart with price snapshot

### Value Objects

[CartId CartItemId Money ProductId VariantId]

### Behaviors

| Command | Pre | Post | Emits |
|---------|-----|------|-------|
| AddItemToCart | [Product is active Product/variant is in stock Quantity > 0] | [CartItem added or quantity updated Total recalculated lastActivityAt updated] | ItemAddedToCart |
| UpdateCartItemQuantity | [Item exists in cart New quantity >= 1 Stock available for new quantity] | [Quantity updated Subtotal recalculated Total recalculated] | CartQuantityUpdated |
| RemoveItemFromCart | [Item exists in cart] | [Item removed from cart Total recalculated] | ItemRemovedFromCart |
| ClearCart | [] | [All items removed Total set to zero] | CartCleared |
| RefreshCartPrices | [] | [All item prices updated to current product prices Total recalculated] | CartPricesRefreshed |

### Events

- **ItemAddedToCart**: [cartId productId variantId quantity unitPrice]
- **CartQuantityUpdated**: [cartId cartItemId oldQuantity newQuantity]
- **ItemRemovedFromCart**: [cartId productId variantId]
- **CartCleared**: [cartId itemCount]
- **CartPricesRefreshed**: [cartId priceChanges]

### Repository: CartRepository

- Load Strategy: Eager load cart with all cart items
- Concurrency: Optimistic locking with version field

---

## AGG-ORDER-001 – Order {#agg-order-001}

**Purpose:** Represents an immutable purchase transaction with line items, shipping, and payment details

### Aggregate Root: Order

**Identity:** OrderId (UUID)

**Attributes:**
| Name | Type | Mutable |
|------|------|----------|
| orderId | OrderId | false |
| customerId | CustomerId | false |
| status | OrderStatus | true |
| shippingAddress | ShippingAddress | false |
| paymentMethod | PaymentMethod | false |
| subtotal | Money | false |
| shippingCost | Money | false |
| totalAmount | Money | false |
| createdAt | DateTime | false |
| trackingNumber | string? | true |

### Invariants

- **INV-ORD-001**: Order must have at least one line item
  - Enforcement: Order creation validates lineItems is not empty; no operation allows removal of last item
- **INV-ORD-002**: Total amount equals subtotal plus shipping cost
  - Enforcement: Calculated on creation; all values are immutable after creation
- **INV-ORD-003**: Shipping is free when subtotal exceeds $50
  - Enforcement: Shipping cost calculation during order creation applies rule automatically
- **INV-ORD-004**: Orders cannot be modified after placement
  - Enforcement: No edit operations exposed; only status transitions allowed
- **INV-ORD-005**: Status transitions must follow valid state machine
  - Enforcement: Status update validates transition is allowed from current state
- **INV-ORD-006**: Can only cancel order before shipped status
  - Enforcement: Cancel operation validates current status is PENDING or CONFIRMED
- **INV-ORD-007**: Line item subtotal equals unitPrice times quantity
  - Enforcement: Calculated on line item creation; immutable thereafter
- **INV-ORD-008**: Prices captured at time of order placement
  - Enforcement: Order creation copies current prices from cart to order line items

### Child Entities

#### OrderLineItem

**Identity:** LineItemId (UUID)

**Purpose:** Immutable snapshot of a product at time of order with quantity and pricing

### Value Objects

[OrderId LineItemId OrderStatus Money ShippingAddress PaymentMethod ProductId VariantId]

### Behaviors

| Command | Pre | Post | Emits |
|---------|-----|------|-------|
| PlaceOrder | [Customer is registered Customer email is verified Cart has items All items in stock Payment authorized] | [Order created with PENDING status Line items captured with price snapshots Inventory reserved] | OrderPlaced |
| ConfirmOrder | [Status is PENDING] | [Status changed to CONFIRMED] | OrderConfirmed |
| ShipOrder | [Status is CONFIRMED Tracking number provided] | [Status changed to SHIPPED Tracking number recorded Inventory deducted] | OrderShipped |
| DeliverOrder | [Status is SHIPPED] | [Status changed to DELIVERED] | OrderDelivered |
| CancelOrder | [Status is PENDING or CONFIRMED] | [Status changed to CANCELLED Inventory released Refund initiated] | OrderCancelled |

### Events

- **OrderPlaced**: [orderId customerId lineItems totalAmount createdAt]
- **OrderConfirmed**: [orderId confirmedAt]
- **OrderShipped**: [orderId trackingNumber shippedAt]
- **OrderDelivered**: [orderId deliveredAt]
- **OrderCancelled**: [orderId reason cancelledAt]

### Repository: OrderRepository

- Load Strategy: Eager load order with all line items
- Concurrency: Optimistic locking with version field for status transitions

---

## AGG-PRODUCT-001 – Product {#agg-product-001}

**Purpose:** Represents an item available for sale with variants, pricing, and categorization

### Aggregate Root: Product

**Identity:** ProductId (UUID)

**Attributes:**
| Name | Type | Mutable |
|------|------|----------|
| productId | ProductId | false |
| name | string | true |
| description | string? | true |
| price | Money | true |
| categoryId | CategoryId | true |
| isActive | boolean | true |
| deletedAt | DateTime? | true |
| createdAt | DateTime | false |

### Invariants

- **INV-PROD-001**: Price must be at least $0.01
  - Enforcement: Money value object validates minimum amount on creation and update
- **INV-PROD-002**: Product must belong to a category
  - Enforcement: Product creation requires valid categoryId; category existence checked
- **INV-PROD-003**: Product name must be 2-200 characters
  - Enforcement: Name validation on creation and update operations
- **INV-PROD-004**: Product description must not exceed 5000 characters
  - Enforcement: Description validation on creation and update operations
- **INV-PROD-005**: Variant SKU must be unique across all variants
  - Enforcement: Variant creation validates SKU uniqueness via repository
- **INV-PROD-006**: Products with order history must be soft-deleted
  - Enforcement: Delete operation checks for order history and performs soft delete if found
- **INV-PROD-007**: Products can have up to 10 images with one primary
  - Enforcement: Image upload validates count limit; primary image selection enforced

### Child Entities

#### ProductVariant

**Identity:** VariantId (UUID)

**Purpose:** Represents a specific variation of a product (size, color combination)

#### ProductImage

**Identity:** ImageId (UUID)

**Purpose:** Product image with primary flag

### Value Objects

[ProductId VariantId ImageId CategoryId Money Sku]

### Behaviors

| Command | Pre | Post | Emits |
|---------|-----|------|-------|
| CreateProduct | [Price > 0 Category exists Name length valid] | [Product created as active] | ProductCreated |
| UpdateProduct | [If price provided, price > 0 If name provided, length valid] | [Product attributes updated Price change logged if applicable] | ProductUpdated |
| DeactivateProduct | [Product is active] | [isActive set to false] | ProductDeactivated |
| ActivateProduct | [Product is inactive Product not deleted] | [isActive set to true] | ProductActivated |
| DeleteProduct | [] | [If has orders: deletedAt set (soft delete); else: permanently removed] | ProductDeleted |
| AddProductVariant | [SKU is unique At least size or color specified] | [Variant added to product] | ProductVariantAdded |
| RemoveProductVariant | [Variant exists No pending orders with this variant] | [Variant removed] | ProductVariantRemoved |
| AddProductImage | [Image count < 10] | [Image added If first image, set as primary] | ProductImageAdded |
| SetPrimaryProductImage | [Image exists on product] | [Specified image set as primary Previous primary cleared] | ProductPrimaryImageChanged |

### Events

- **ProductCreated**: [productId name price categoryId]
- **ProductUpdated**: [productId changedFields]
- **ProductDeactivated**: [productId deactivatedAt]
- **ProductActivated**: [productId activatedAt]
- **ProductDeleted**: [productId deletedAt isSoftDelete]
- **ProductVariantAdded**: [productId variantId sku]
- **ProductVariantRemoved**: [productId variantId]
- **ProductImageAdded**: [productId imageId url]
- **ProductPrimaryImageChanged**: [productId imageId]

### Repository: ProductRepository

- Load Strategy: Eager load product with variants and images
- Concurrency: Optimistic locking with version field

---

## AGG-CATEGORY-001 – Category {#agg-category-001}

**Purpose:** Organizes products into browsable hierarchical groupings

### Aggregate Root: Category

**Identity:** CategoryId (UUID)

**Attributes:**
| Name | Type | Mutable |
|------|------|----------|
| categoryId | CategoryId | false |
| name | string | true |
| description | string? | true |
| parentCategoryId | CategoryId? | true |
| depth | integer | true |

### Invariants

- **INV-CAT-001**: Category name must be unique
  - Enforcement: Category creation validates name uniqueness via repository
- **INV-CAT-002**: Category cannot be its own parent
  - Enforcement: Parent assignment validates parentCategoryId != categoryId
- **INV-CAT-003**: Category hierarchy limited to 3 levels deep
  - Enforcement: Category creation validates parent hierarchy depth before assignment

### Value Objects

[CategoryId]

### Behaviors

| Command | Pre | Post | Emits |
|---------|-----|------|-------|
| CreateCategory | [Name is unique If parent specified, parent exists and depth < 3] | [Category created with calculated depth] | CategoryCreated |
| RenameCategory | [New name is unique] | [Category name updated] | CategoryRenamed |
| MoveCategoryToParent | [New parent exists Would not exceed depth limit Not creating cycle] | [Parent updated Depth recalculated for self and descendants] | CategoryMoved |
| DeleteCategory | [No products in category No child categories] | [Category removed] | CategoryDeleted |

### Events

- **CategoryCreated**: [categoryId name parentCategoryId depth]
- **CategoryRenamed**: [categoryId oldName newName]
- **CategoryMoved**: [categoryId oldParentId newParentId newDepth]
- **CategoryDeleted**: [categoryId]

### Repository: CategoryRepository

- Load Strategy: Lazy load children on demand
- Concurrency: Optimistic locking with version field

---

## AGG-INVENTORY-001 – Inventory {#agg-inventory-001}

**Purpose:** Tracks stock levels for products and manages availability with reservations

### Aggregate Root: Inventory

**Identity:** InventoryId (UUID), also keyed by ProductId

**Attributes:**
| Name | Type | Mutable |
|------|------|----------|
| inventoryId | InventoryId | false |
| productId | ProductId | false |
| quantityOnHand | integer | true |
| reservedQuantity | integer | true |
| availableQuantity | integer | true |

### Invariants

- **INV-INV-001**: Available quantity cannot be negative
  - Enforcement: All operations validate resulting availableQuantity >= 0 before committing
- **INV-INV-002**: Reserved quantity cannot exceed quantity on hand
  - Enforcement: Reserve operation validates reservedQuantity + requested <= quantityOnHand
- **INV-INV-003**: No backorders allowed
  - Enforcement: Reserve operation rejects if requested > availableQuantity
- **INV-INV-004**: Each product has exactly one inventory record
  - Enforcement: Inventory creation validates no existing record for productId

### Value Objects

[InventoryId ProductId Quantity]

### Behaviors

| Command | Pre | Post | Emits |
|---------|-----|------|-------|
| ReserveInventory | [Available quantity >= requested quantity] | [Reserved quantity increased Available quantity decreased] | InventoryReserved |
| ReleaseInventory | [Reserved quantity >= requested quantity] | [Reserved quantity decreased Available quantity increased] | InventoryReleased |
| DeductInventory | [Reserved quantity >= requested quantity] | [Quantity on hand decreased Reserved quantity decreased] | InventoryDeducted |
| RestockInventory | [Quantity > 0] | [Quantity on hand increased Available quantity increased] | InventoryRestocked |
| AdjustInventory | [Resulting quantities would be valid] | [Quantities adjusted Audit log entry created] | InventoryAdjusted |

### Events

- **InventoryReserved**: [inventoryId productId quantity orderId]
- **InventoryReleased**: [inventoryId productId quantity reason]
- **InventoryDeducted**: [inventoryId productId quantity orderId]
- **InventoryRestocked**: [inventoryId productId quantity previousQuantity]
- **InventoryAdjusted**: [inventoryId productId oldQuantity newQuantity reason operator]
- **OutOfStock**: [inventoryId productId]

### Repository: InventoryRepository

- Load Strategy: Eager load single record
- Concurrency: Pessimistic locking for reserve/deduct operations to prevent overselling

---

