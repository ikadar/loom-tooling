# Aggregate Design

Generated: 2026-01-02T18:27:24+01:00

---

## AGG-CUSTOMER-001 – Customer {#agg-customer-001}

**Purpose:** Manages customer registration, authentication, and profile data with order placement capability

### Aggregate Root: Customer

**Identity:** CustomerId (UUID)

**Attributes:**
| Name | Type | Mutable |
|------|------|----------|
| customerId | CustomerId | false |
| email | Email | true |
| password | HashedPassword | true |
| registrationStatus | RegistrationStatus | true |
| emailVerified | boolean | true |
| cartId | CartId | false |

### Invariants

- **INV-CUST-001**: Customer must have a valid email address
  - Enforcement: Email value object validates format on creation; constructor rejects invalid emails
- **INV-CUST-002**: Email addresses must be unique across all customer accounts
  - Enforcement: Repository checks uniqueness before persisting; database unique constraint as backup
- **INV-CUST-003**: Customer must be registered with verified email to place orders
  - Enforcement: placeOrder method checks registrationStatus is REGISTERED and emailVerified is true
- **INV-CUST-004**: Password must meet minimum security requirements
  - Enforcement: Password value object validates on creation

### Value Objects

[CustomerId Email HashedPassword RegistrationStatus]

### Behaviors

| Command | Pre | Post | Emits |
|---------|-----|------|-------|
| RegisterCustomer | [Email not already registered Password meets requirements] | [Customer created with REGISTERED status Cart created for customer] | CustomerRegistered |
| VerifyCustomerEmail | [Valid verification token] | [emailVerified set to true] | EmailVerified |
| ChangeCustomerEmail | [New email not already registered Email format valid] | [Email updated emailVerified reset to false] | CustomerEmailChanged |
| ChangeCustomerPassword | [Password meets requirements] | [Password updated] | CustomerPasswordChanged |

### Events

- **CustomerRegistered**: [customerId email registeredAt]
- **EmailVerified**: [customerId verifiedAt]
- **CustomerEmailChanged**: [customerId oldEmail newEmail]
- **CustomerPasswordChanged**: [customerId changedAt]

### Repository: CustomerRepository

- Load Strategy: Eager load all attributes; no child entities
- Concurrency: Optimistic locking with version field

---

## AGG-CART-001 – Cart {#agg-cart-001}

**Purpose:** Manages shopping cart items with real-time pricing and stock validation before checkout

### Aggregate Root: Cart

**Identity:** CartId (UUID)

**Attributes:**
| Name | Type | Mutable |
|------|------|----------|
| cartId | CartId | false |
| customerId | CustomerId | false |
| lastActivityAt | DateTime | true |

### Invariants

- **INV-CART-001**: Total price equals sum of all item subtotals
  - Enforcement: Total recalculated on every item add/update/remove operation
- **INV-CART-002**: Each product can appear only once in cart
  - Enforcement: addItem checks for existing product and updates quantity instead of adding duplicate
- **INV-CART-003**: Cart item quantity must be at least 1
  - Enforcement: CartItem constructor and updateQuantity validate quantity >= 1
- **INV-CART-004**: Cart items reflect current product prices
  - Enforcement: Cart recalculates totals by fetching current prices on cart view
- **INV-CART-005**: Cart item quantity cannot exceed available stock
  - Enforcement: Application service validates stock before cart operations

### Child Entities

#### CartItem

**Identity:** CartItemId (UUID)

**Purpose:** Represents a product and quantity within the shopping cart

### Value Objects

[CartId CartItemId ProductId VariantId Money Quantity]

### Behaviors

| Command | Pre | Post | Emits |
|---------|-----|------|-------|
| AddItemToCart | [Product exists and is active Stock available Quantity > 0 If product has variants, variant must be specified] | [CartItem added or quantity updated lastActivityAt updated] | ItemAddedToCart |
| UpdateCartItemQuantity | [Item exists in cart Quantity >= 1 Stock available for new quantity] | [Quantity updated lastActivityAt updated] | CartQuantityUpdated |
| RemoveItemFromCart | [Item exists in cart] | [Item removed lastActivityAt updated] | ItemRemovedFromCart |
| ClearCart | [] | [All items removed lastActivityAt updated] | CartCleared |

### Events

- **ItemAddedToCart**: [cartId productId variantId quantity]
- **CartQuantityUpdated**: [cartId cartItemId oldQuantity newQuantity]
- **ItemRemovedFromCart**: [cartId productId]
- **CartCleared**: [cartId clearedAt]

### Repository: CartRepository

- Load Strategy: Eager load Cart with all CartItems
- Concurrency: Optimistic locking with version field

---

## AGG-ORDER-001 – Order {#agg-order-001}

**Purpose:** Represents a completed purchase transaction with immutable line items and state machine for fulfillment

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
| cancellationReason | string? | true |

### Invariants

- **INV-ORD-001**: Order must have at least one line item
  - Enforcement: Constructor validates lineItems list is not empty
- **INV-ORD-002**: Total amount equals subtotal plus shipping cost
  - Enforcement: Calculated field; recalculated on creation only (immutable after)
- **INV-ORD-003**: Shipping is free when subtotal exceeds $50
  - Enforcement: Shipping calculation applies zero cost when subtotal >= 50
- **INV-ORD-004**: Order items are immutable after creation
  - Enforcement: No setter methods for lineItems or shippingAddress after construction
- **INV-ORD-005**: Status transitions must follow valid state machine
  - Enforcement: updateStatus method validates transition against allowed transitions map
- **INV-ORD-006**: Can only cancel before shipped status
  - Enforcement: cancel method checks current status is PENDING or CONFIRMED
- **INV-ORD-007**: Order line items capture prices at order time
  - Enforcement: OrderLineItem created with price snapshot; no price setters

### Child Entities

#### OrderLineItem

**Identity:** LineItemId (UUID)

**Purpose:** Immutable snapshot of a product at time of order with quantity and pricing

### Value Objects

[OrderId LineItemId OrderStatus Money ShippingAddress PaymentMethod CustomerId ProductId]

### Behaviors

| Command | Pre | Post | Emits |
|---------|-----|------|-------|
| PlaceOrder | [Customer is registered with verified email At least one line item All items in stock Payment authorized Valid shipping address] | [Order created with PENDING status Price snapshots captured Shipping cost calculated Inventory reserved] | OrderPlaced |
| ConfirmOrder | [Status is PENDING] | [Status changes to CONFIRMED] | OrderConfirmed |
| ShipOrder | [Status is CONFIRMED Tracking number provided] | [Status changes to SHIPPED Tracking number recorded] | OrderShipped |
| DeliverOrder | [Status is SHIPPED] | [Status changes to DELIVERED] | OrderDelivered |
| CancelOrder | [Status is PENDING or CONFIRMED] | [Status changes to CANCELLED Cancellation reason recorded] | OrderCancelled |

### Events

- **OrderPlaced**: [orderId customerId lineItems totalAmount createdAt]
- **OrderConfirmed**: [orderId confirmedAt]
- **OrderShipped**: [orderId trackingNumber shippedAt]
- **OrderDelivered**: [orderId deliveredAt]
- **OrderCancelled**: [orderId reason cancelledAt]

### Repository: OrderRepository

- Load Strategy: Eager load Order with all OrderLineItems
- Concurrency: Optimistic locking with version field

---

## AGG-PRODUCT-001 – Product {#agg-product-001}

**Purpose:** Represents a sellable item in the catalog with variants, pricing, and category association

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

### Invariants

- **INV-PROD-001**: Price must be greater than zero
  - Enforcement: Money value object validates minimum of 0.01
- **INV-PROD-002**: Product name must be between 2 and 200 characters
  - Enforcement: Constructor and update method validate name length
- **INV-PROD-003**: Product description must not exceed 5000 characters
  - Enforcement: Constructor and update method validate description length
- **INV-PROD-004**: Product must belong to a category
  - Enforcement: Constructor requires categoryId; cannot be null
- **INV-PROD-005**: Inactive products cannot be added to carts
  - Enforcement: Application service checks isActive before cart operations
- **INV-PROD-006**: Products can have up to 10 images with exactly one primary
  - Enforcement: addImage validates count; setPrimaryImage ensures exactly one primary
- **INV-PROD-007**: Variant SKU must be unique
  - Enforcement: Repository validates SKU uniqueness before save
- **INV-PROD-008**: Variant must have at least one attribute
  - Enforcement: ProductVariant constructor validates size or color is provided

### Child Entities

#### ProductVariant

**Identity:** VariantId (UUID)

**Purpose:** Represents a specific variation of a product (size, color)

#### ProductImage

**Identity:** ImageId (UUID)

**Purpose:** Product image with primary flag

### Value Objects

[ProductId VariantId ImageId CategoryId Money SKU]

### Behaviors

| Command | Pre | Post | Emits |
|---------|-----|------|-------|
| CreateProduct | [Price > 0 Name 2-200 chars Category exists] | [Product created as active] | ProductCreated |
| UpdateProduct | [If price provided, price > 0 If name provided, 2-200 chars] | [Product attributes updated] | ProductUpdated |
| DeactivateProduct | [Product is active] | [isActive set to false] | ProductDeactivated |
| DeleteProduct | [] | [deletedAt set to current timestamp] | ProductDeleted |
| AddProductVariant | [SKU is unique Size or color specified] | [Variant added to product] | VariantAdded |
| AddProductImage | [Image count < 10] | [Image added; if first image, set as primary] | ImageAdded |
| SetPrimaryProductImage | [Image exists] | [Selected image is primary; others are not] | PrimaryImageChanged |

### Events

- **ProductCreated**: [productId name price categoryId]
- **ProductUpdated**: [productId changedFields]
- **ProductDeactivated**: [productId deactivatedAt]
- **ProductDeleted**: [productId deletedAt]
- **VariantAdded**: [productId variantId sku]
- **ImageAdded**: [productId imageId url]
- **PrimaryImageChanged**: [productId imageId]

### Repository: ProductRepository

- Load Strategy: Eager load Product with variants and images
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
  - Enforcement: Repository validates uniqueness before save
- **INV-CAT-002**: Category cannot be its own parent
  - Enforcement: setParent method validates parentCategoryId != categoryId
- **INV-CAT-003**: Category hierarchy limited to 3 levels
  - Enforcement: setParent validates resulting depth <= 3

### Value Objects

[CategoryId]

### Behaviors

| Command | Pre | Post | Emits |
|---------|-----|------|-------|
| CreateCategory | [Name is unique If parent provided, resulting depth <= 3] | [Category created Depth calculated from parent] | CategoryCreated |
| UpdateCategory | [If name changed, new name is unique] | [Category attributes updated] | CategoryUpdated |
| SetCategoryParent | [Not setting self as parent Resulting depth <= 3] | [parentCategoryId updated Depth recalculated] | CategoryParentChanged |

### Events

- **CategoryCreated**: [categoryId name parentCategoryId]
- **CategoryUpdated**: [categoryId changedFields]
- **CategoryParentChanged**: [categoryId oldParentId newParentId]

### Repository: CategoryRepository

- Load Strategy: Load category without children; load children on demand
- Concurrency: Optimistic locking with version field

---

## AGG-INVENTORY-001 – Inventory {#agg-inventory-001}

**Purpose:** Tracks stock levels for products and manages availability with reservation support

### Aggregate Root: Inventory

**Identity:** InventoryId (UUID)

**Attributes:**
| Name | Type | Mutable |
|------|------|----------|
| inventoryId | InventoryId | false |
| productId | ProductId | false |
| quantityOnHand | integer | true |
| reservedQuantity | integer | true |

### Invariants

- **INV-INV-001**: Available quantity cannot be negative
  - Enforcement: All operations validate resulting availableQuantity >= 0
- **INV-INV-002**: Reserved quantity cannot exceed quantity on hand
  - Enforcement: reserve method validates reservedQuantity + requested <= quantityOnHand
- **INV-INV-003**: Each product has exactly one inventory record
  - Enforcement: Repository enforces unique constraint on productId
- **INV-INV-004**: No backorders allowed
  - Enforcement: isInStock method used to validate before order placement

### Value Objects

[InventoryId ProductId Quantity]

### Behaviors

| Command | Pre | Post | Emits |
|---------|-----|------|-------|
| ReserveInventory | [Available quantity >= requested quantity] | [Reserved quantity increased by requested amount] | InventoryReserved |
| ReleaseInventory | [Reserved quantity >= requested quantity] | [Reserved quantity decreased by requested amount] | InventoryReleased |
| DeductInventory | [Reserved quantity >= requested quantity] | [Quantity on hand decreased Reserved quantity decreased] | InventoryDeducted |
| RestockInventory | [Quantity > 0] | [Quantity on hand increased] | InventoryRestocked |
| RestoreInventory | [Quantity > 0] | [Quantity on hand increased (for cancelled orders)] | InventoryRestored |

### Events

- **InventoryReserved**: [inventoryId productId quantity newAvailable]
- **InventoryReleased**: [inventoryId productId quantity newAvailable]
- **InventoryDeducted**: [inventoryId productId quantity newOnHand]
- **InventoryRestocked**: [inventoryId productId quantity newOnHand]
- **InventoryRestored**: [inventoryId productId quantity reason]
- **OutOfStock**: [inventoryId productId]

### Repository: InventoryRepository

- Load Strategy: Simple load; no child entities
- Concurrency: Pessimistic locking for reserve/release/deduct to prevent overselling

---

