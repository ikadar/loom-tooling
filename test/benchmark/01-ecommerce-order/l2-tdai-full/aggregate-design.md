# Aggregate Design

Generated: 2025-12-31T12:00:34+01:00

---

## AGG-CUSTOMER-001 – Customer

**Purpose:** Manages customer identity, authentication, and serves as the anchor for shopping cart and order history

### Aggregate Root: Customer

**Identity:** CustomerId (UUID)

**Attributes:**
| Name | Type | Mutable |
|------|------|----------|
| email | Email | false |
| passwordHash | string | true |
| registrationStatus | RegistrationStatus | true |
| emailVerified | boolean | true |
| createdAt | DateTime | false |

### Invariants

- **INV-CUST-001**: Customer must be registered to place orders
  - Enforcement: Guard clause in placeOrder() validates registrationStatus='registered'
- **INV-CUST-002**: Email must be unique across all customers
  - Enforcement: Database unique constraint on email field
- **INV-CUST-003**: Password must be at least 8 characters with at least one digit
  - Enforcement: Validation in register() and changePassword() operations
- **INV-CUST-004**: Customer always has exactly one cart
  - Enforcement: Cart created automatically during customer registration

### Value Objects

[Email RegistrationStatus]

### Behaviors

| Command | Pre | Post | Emits |
|---------|-----|------|-------|
| RegisterCustomer | [email not already registered] | [registrationStatus == 'registered' cart created] | CustomerRegistered |
| VerifyCustomerEmail | [registrationStatus == 'registered' emailVerified == false] | [emailVerified == true] | CustomerEmailVerified |

### Events

- **CustomerRegistered**: [customerId email registeredAt]
- **CustomerEmailVerified**: [customerId verifiedAt]

### Repository: CustomerRepository

- Load Strategy: Load customer without cart (cart is separate aggregate)
- Concurrency: Optimistic locking via version field

---

## AGG-CART-001 – Cart

**Purpose:** Manages shopping cart state including item selection, quantity adjustments, and price calculations before checkout

### Aggregate Root: Cart

**Identity:** CartId (UUID)

**Attributes:**
| Name | Type | Mutable |
|------|------|----------|
| customerId | CustomerId | false |
| totalPrice | Money | true |
| lastActivityDate | DateTime | true |
| createdAt | DateTime | false |

### Invariants

- **INV-CART-001**: Total price must equal sum of all cart item subtotals
  - Enforcement: Recalculated after any item modification
- **INV-CART-002**: Each product can appear only once per cart
  - Enforcement: addItem() updates quantity if product exists instead of creating duplicate
- **INV-CART-003**: Only in-stock products can be added
  - Enforcement: Stock availability check via Inventory service before addItem()
- **INV-CART-004**: Cart item quantity cannot exceed available stock
  - Enforcement: Quantity validation against Inventory.availableQuantity
- **INV-CART-005**: Cart item quantity must be at least 1
  - Enforcement: Setting quantity to 0 triggers automatic item removal
- **INV-CART-006**: Products with variants require specific variant selection
  - Enforcement: Validation in addItem() checks if product has variants

### Child Entities

#### CartItem

**Identity:** CartItemId (UUID)

**Purpose:** Represents a product with quantity in the shopping cart with price snapshot

### Value Objects

[Money Quantity]

### Behaviors

| Command | Pre | Post | Emits |
|---------|-----|------|-------|
| AddItemToCart | [product is active product is in stock quantity > 0 variant selected if required] | [item added or quantity updated totalPrice recalculated] | ItemAddedToCart |
| UpdateCartItemQuantity | [product exists in cart quantity > 0 quantity <= availableStock] | [quantity updated subtotal recalculated totalPrice recalculated] | CartItemQuantityUpdated |
| RemoveItemFromCart | [product exists in cart] | [item removed totalPrice recalculated] | ItemRemovedFromCart |
| ClearCart | [] | [all items removed totalPrice == 0] | CartCleared |

### Events

- **ItemAddedToCart**: [cartId productId variantId productName unitPrice quantity]
- **CartItemQuantityUpdated**: [cartId productId oldQuantity newQuantity]
- **ItemRemovedFromCart**: [cartId productId]
- **CartCleared**: [cartId]

### Repository: CartRepository

- Load Strategy: Load complete cart with all items
- Concurrency: Optimistic locking via version field

---

## AGG-ORDER-001 – Order

**Purpose:** Manages the complete order lifecycle from placement through delivery, including payment, shipping, and status transitions

### Aggregate Root: Order

**Identity:** OrderId (UUID)

**Attributes:**
| Name | Type | Mutable |
|------|------|----------|
| customerId | CustomerId | false |
| status | OrderStatus | true |
| shippingAddress | Address | false |
| paymentMethod | PaymentMethod | false |
| subtotal | Money | false |
| shippingCost | Money | false |
| totalAmount | Money | false |
| createdAt | DateTime | false |

### Invariants

- **INV-ORDER-001**: Order must have at least one line item
  - Enforcement: Validated on creation
- **INV-ORDER-002**: Total amount must equal subtotal plus shipping cost
  - Enforcement: Calculated during order creation
- **INV-ORDER-003**: Shipping is free when subtotal exceeds $50
  - Enforcement: Shipping cost calculation sets to $0 when subtotal > $50
- **INV-ORDER-004**: Cannot modify order after status is 'shipped' or 'delivered'
  - Enforcement: Guard clause in all modification methods
- **INV-ORDER-005**: Can only cancel before status is 'shipped'
  - Enforcement: Precondition check in cancel() method
- **INV-ORDER-006**: Status transitions must follow valid state machine
  - Enforcement: State pattern: pending→confirmed→shipped→delivered, or pending/confirmed→cancelled
- **INV-ORDER-007**: Line item prices are immutable snapshots from order time
  - Enforcement: OrderLineItem attributes set at creation, no setters provided

### Child Entities

#### OrderLineItem

**Identity:** LineItemId (UUID)

**Purpose:** Immutable snapshot of a product at the time of order with quantity and price

### Value Objects

[Money Address PaymentMethod OrderStatus Quantity]

### Behaviors

| Command | Pre | Post | Emits |
|---------|-----|------|-------|
| PlaceOrder | [customer is registered cart has items all items in stock payment authorized shipping address valid domestic shipping only] | [order created with status 'pending' inventory reserved cart cleared] | OrderPlaced |
| ConfirmOrder | [status == 'pending'] | [status == 'confirmed'] | OrderConfirmed |
| ShipOrder | [status == 'confirmed'] | [status == 'shipped'] | OrderShipped |
| DeliverOrder | [status == 'shipped'] | [status == 'delivered'] | OrderDelivered |
| CancelOrder | [status in ['pending', 'confirmed']] | [status == 'cancelled' inventory restored refund initiated] | OrderCancelled |

### Events

- **OrderPlaced**: [orderId customerId lineItems totalAmount shippingAddress]
- **OrderConfirmed**: [orderId confirmedAt]
- **OrderShipped**: [orderId shippedAt trackingNumber]
- **OrderDelivered**: [orderId deliveredAt]
- **OrderCancelled**: [orderId cancelledAt reason]

### Repository: OrderRepository

- Load Strategy: Load complete order with all line items
- Concurrency: Optimistic locking via version field

---

## AGG-PRODUCT-001 – Product

**Purpose:** Manages product catalog information including descriptions, pricing, variants, and category assignment

### Aggregate Root: Product

**Identity:** ProductId (UUID)

**Attributes:**
| Name | Type | Mutable |
|------|------|----------|
| name | string | true |
| description | string | true |
| price | Money | true |
| categoryId | CategoryId | true |
| isActive | boolean | true |
| createdAt | DateTime | false |

### Invariants

- **INV-PROD-001**: Price must be greater than or equal to $0.01
  - Enforcement: Validation during creation and update
- **INV-PROD-002**: Product name must be between 2 and 200 characters
  - Enforcement: Validation during creation and update
- **INV-PROD-003**: Product must belong to a category
  - Enforcement: Required categoryId on creation
- **INV-PROD-004**: Inactive products cannot be added to cart
  - Enforcement: isActive check in cart service
- **INV-PROD-005**: Each variant SKU must be unique
  - Enforcement: Database unique constraint on SKU
- **INV-PROD-006**: Variant combination (size, color) must be unique within product
  - Enforcement: Validation during addVariant()
- **INV-PROD-007**: Products with order history must be soft-deleted only
  - Enforcement: Delete operation checks OrderLineItem references

### Child Entities

#### ProductVariant

**Identity:** VariantId (UUID)

**Purpose:** Represents a specific variation of a product with unique SKU

### Value Objects

[Money SKU]

### Behaviors

| Command | Pre | Post | Emits |
|---------|-----|------|-------|
| CreateProduct | [name is not empty price > 0 category exists] | [product created isActive == true] | ProductCreated |
| UpdateProduct | [product exists] | [attributes updated] | ProductUpdated |
| DeactivateProduct | [isActive == true] | [isActive == false] | ProductDeactivated |
| AddProductVariant | [variant SKU is unique variant combination is unique within product] | [variant added to product] | ProductVariantAdded |

### Events

- **ProductCreated**: [productId name price categoryId]
- **ProductUpdated**: [productId changedFields]
- **ProductDeactivated**: [productId deactivatedAt]
- **ProductVariantAdded**: [productId variantId sku size color]

### Repository: ProductRepository

- Load Strategy: Load product with all variants
- Concurrency: Optimistic locking via version field

---

## AGG-CATEGORY-001 – Category

**Purpose:** Organizes products into a hierarchical structure for browsing and navigation

### Aggregate Root: Category

**Identity:** CategoryId (UUID)

**Attributes:**
| Name | Type | Mutable |
|------|------|----------|
| name | string | true |
| description | string | true |
| parentCategoryId | CategoryId? | true |
| depth | integer | true |

### Invariants

- **INV-CAT-001**: Category name must be unique
  - Enforcement: Database unique constraint on name
- **INV-CAT-002**: Category cannot be its own parent
  - Enforcement: Validation during creation and update
- **INV-CAT-003**: Category hierarchy is limited to 3 levels deep
  - Enforcement: Depth calculation and validation during creation

### Behaviors

| Command | Pre | Post | Emits |
|---------|-----|------|-------|
| CreateCategory | [name is unique depth <= 3 if parent specified] | [category created] | CategoryCreated |
| RenameCategory | [new name is unique] | [name updated] | CategoryRenamed |
| MoveCategoryToParent | [new parent exists would not exceed depth 3 not creating circular reference] | [parentCategoryId updated depth recalculated] | CategoryMoved |

### Events

- **CategoryCreated**: [categoryId name parentCategoryId]
- **CategoryRenamed**: [categoryId oldName newName]
- **CategoryMoved**: [categoryId oldParentId newParentId]

### Repository: CategoryRepository

- Load Strategy: Load single category (children loaded separately)
- Concurrency: Optimistic locking via version field

---

## AGG-INVENTORY-001 – Inventory

**Purpose:** Tracks stock levels, reservations, and availability for products with full audit trail

### Aggregate Root: Inventory

**Identity:** InventoryId (UUID)

**Attributes:**
| Name | Type | Mutable |
|------|------|----------|
| productId | ProductId | false |
| stockLevel | integer | true |
| reservedQuantity | integer | true |
| availableQuantity | integer | false |

### Invariants

- **INV-INV-001**: Stock level cannot be negative
  - Enforcement: Validation in all stock adjustment operations
- **INV-INV-002**: Reserved quantity cannot exceed stock level
  - Enforcement: Validation in reserve() operation
- **INV-INV-003**: Available quantity must equal stock level minus reserved
  - Enforcement: Calculated property, recalculated on any change
- **INV-INV-004**: All inventory changes must be logged
  - Enforcement: Audit log entry created for every modification

### Behaviors

| Command | Pre | Post | Emits |
|---------|-----|------|-------|
| ReserveInventory | [quantity <= availableQuantity] | [reservedQuantity increased by quantity audit log created] | InventoryReserved |
| ReleaseInventory | [quantity <= reservedQuantity] | [reservedQuantity decreased by quantity audit log created] | InventoryReleased |
| DeductInventory | [quantity <= reservedQuantity] | [stockLevel decreased reservedQuantity decreased audit log created] | InventoryDeducted |
| RestockInventory | [quantity > 0] | [stockLevel increased by quantity audit log created] | InventoryRestocked |

### Events

- **InventoryReserved**: [inventoryId productId quantity orderId]
- **InventoryReleased**: [inventoryId productId quantity reason]
- **InventoryDeducted**: [inventoryId productId quantity orderId]
- **InventoryRestocked**: [inventoryId productId quantity previousLevel newLevel]
- **OutOfStock**: [inventoryId productId]

### Repository: InventoryRepository

- Load Strategy: Load single inventory record
- Concurrency: Pessimistic locking for stock operations to prevent overselling

---

