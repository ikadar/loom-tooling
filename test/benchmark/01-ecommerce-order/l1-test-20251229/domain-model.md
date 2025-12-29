# Domain Model

Generated: 2025-12-29T13:37:22+01:00

**Domain:** E-Commerce Order System

Domain model for an online store supporting product browsing, shopping cart, order placement, and order tracking

---

## Entities

### ENT-CUSTOMER-001 – Customer

**Type:** aggregate_root

**Purpose:** Represents a registered user who can browse products, manage a cart, and place orders

**Attributes:**

| Name | Type | Constraints |
|------|------|-------------|
| customerId | CustomerId (UUID) | Required, immutable after creation |
| email | Email | Required, unique, valid email format |
| registrationStatus | RegistrationStatus (enum) | Required, values: registered, unregistered |
| cart | Cart | One per customer |

**Invariants:**
- Customer must be registered to place orders
- Email must be unique across all customers
- Customer always has exactly one cart

**Operations:**
- `register(email: Email, password: Password): void`
  - Pre: [Email not already registered]
  - Post: [Status set to 'registered' CustomerRegistered event emitted]
- `placeOrder(shippingAddress: Address, paymentMethod: PaymentMethod): Order`
  - Pre: [Customer is registered Cart has items All cart items in stock]
  - Post: [Order created Cart cleared OrderPlaced event emitted]

**Events:**
- `CustomerRegistered` (trigger: register() called successfully)

**Relationships:**
- contains ENT-CART-001 (one-to-one)
- contains ENT-ORDER-001 (one-to-many)

---

### ENT-CART-001 – Cart

**Type:** aggregate_root

**Purpose:** Holds products a customer intends to purchase before checkout

**Attributes:**

| Name | Type | Constraints |
|------|------|-------------|
| cartId | CartId (UUID) | Required, immutable after creation |
| customerId | CustomerId | Required, references owning customer |
| items | List<CartItem> | Can be empty |
| totalPrice | Money | Calculated, >= 0 |

**Invariants:**
- Total price must equal sum of all cart item subtotals
- Each product can appear only once (quantity adjusted instead)
- Only in-stock products can be added

**Operations:**
- `addItem(product: Product, quantity: Quantity): void`
  - Pre: [Product is in stock Quantity > 0]
  - Post: [Item added or quantity updated Total recalculated ItemAddedToCart event emitted]
- `updateQuantity(productId: ProductId, quantity: Quantity): void`
  - Pre: [Product exists in cart Quantity > 0 Requested quantity in stock]
  - Post: [Quantity updated Total recalculated CartItemQuantityUpdated event emitted]
- `removeItem(productId: ProductId): void`
  - Pre: [Product exists in cart]
  - Post: [Item removed Total recalculated ItemRemovedFromCart event emitted]
- `clear(): void`
  - Post: [All items removed Total set to zero CartCleared event emitted]

**Events:**
- `ItemAddedToCart` (trigger: addItem() called)
- `CartItemQuantityUpdated` (trigger: updateQuantity() called)
- `ItemRemovedFromCart` (trigger: removeItem() called)
- `CartCleared` (trigger: clear() called or order placed)

**Relationships:**
- contains ENT-CARTITEM-001 (one-to-many)
- belongs_to ENT-CUSTOMER-001 (many-to-one)

---

### ENT-CARTITEM-001 – CartItem

**Type:** entity

**Purpose:** Represents a product and its quantity within a shopping cart

**Attributes:**
| Name | Type | Constraints |
|------|------|-------------|
| cartItemId | CartItemId (UUID) | Required, immutable |
| productId | ProductId | Required, references Product |
| productName | string | Required, snapshot at time of adding |
| unitPrice | Money | Required, snapshot at time of adding |
| quantity | Quantity | Required, > 0 |
| subtotal | Money | Calculated: unitPrice * quantity |

**Invariants:**
- Quantity must be greater than zero
- Subtotal must equal unitPrice * quantity

**Operations:**
- `updateQuantity(newQuantity: Quantity): void`
  - Pre: [newQuantity > 0]
  - Post: [Quantity updated Subtotal recalculated]

**Relationships:**
- belongs_to ENT-CART-001 (many-to-one)
- references ENT-PRODUCT-001 (many-to-one)

---

### ENT-ORDER-001 – Order

**Type:** aggregate_root

**Purpose:** Represents a customer's purchase transaction with shipping and payment details

**Attributes:**

| Name | Type | Constraints |
|------|------|-------------|
| orderId | OrderId (UUID) | Required, immutable after creation |
| customerId | CustomerId | Required, references Customer |
| status | OrderStatus (enum) | Required, values: pending, confirmed, shipped, delivered, cancelled |
| lineItems | List<OrderLineItem> | Required, at least one item |
| shippingAddress | Address | Required |
| paymentMethod | PaymentMethod | Required |
| subtotal | Money | Calculated, sum of line items |
| shippingCost | Money | Calculated based on subtotal (free if > $50) |
| totalAmount | Money | Calculated: subtotal + shippingCost |
| createdAt | DateTime | Required, immutable |

**Invariants:**
- Order must have at least one line item
- Total amount must equal subtotal plus shipping cost
- Shipping is free when subtotal exceeds $50
- Cannot modify order after status is 'shipped' or 'delivered'
- Can only cancel before status is 'shipped'
- Status transitions: pending → confirmed → shipped → delivered (or cancelled from pending/confirmed)

**Operations:**
- `confirm(): void`
  - Pre: [Status is 'pending']
  - Post: [Status set to 'confirmed' OrderConfirmed event emitted]
- `ship(): void`
  - Pre: [Status is 'confirmed']
  - Post: [Status set to 'shipped' OrderShipped event emitted]
- `deliver(): void`
  - Pre: [Status is 'shipped']
  - Post: [Status set to 'delivered' OrderDelivered event emitted]
- `cancel(): void`
  - Pre: [Status is 'pending' or 'confirmed']
  - Post: [Status set to 'cancelled' OrderCancelled event emitted Inventory restored]

**Events:**
- `OrderPlaced` (trigger: Order created from cart)
- `OrderConfirmed` (trigger: confirm() called)
- `OrderShipped` (trigger: ship() called)
- `OrderDelivered` (trigger: deliver() called)
- `OrderCancelled` (trigger: cancel() called)

**Relationships:**
- belongs_to ENT-CUSTOMER-001 (many-to-one)
- contains ENT-ORDERLINEITEM-001 (one-to-many)

---

### ENT-ORDERLINEITEM-001 – OrderLineItem

**Type:** entity

**Purpose:** Immutable snapshot of a product at the time of order, with quantity and price

**Attributes:**

| Name | Type | Constraints |
|------|------|-------------|
| lineItemId | LineItemId (UUID) | Required, immutable |
| productId | ProductId | Required, reference to original product |
| productName | string | Required, snapshot |
| unitPrice | Money | Required, snapshot at order time |
| quantity | Quantity | Required, > 0 |
| subtotal | Money | Calculated: unitPrice * quantity |

**Invariants:**
- All attributes are immutable after creation
- Quantity must be greater than zero
- Subtotal must equal unitPrice * quantity

**Relationships:**
- belongs_to ENT-ORDER-001 (many-to-one)
- references ENT-PRODUCT-001 (many-to-one)

---

### ENT-PRODUCT-001 – Product

**Type:** aggregate_root

**Purpose:** Represents a purchasable item in the catalog with its details and variants

**Attributes:**

| Name | Type | Constraints |
|------|------|-------------|
| productId | ProductId (UUID) | Required, immutable |
| name | string | Required, 1-200 characters |
| description | string | Optional, max 2000 characters |
| price | Money | Required, > 0 |
| categoryId | CategoryId | Required, references Category |
| variants | List<ProductVariant> | Optional |
| isActive | boolean | Required, default true |

**Invariants:**
- Price must be greater than zero
- Product must belong to a category
- Inactive products cannot be added to cart

**Operations:**
- `create(name: string, description: string, price: Money, categoryId: CategoryId): Product`
  - Pre: [Name is not empty Price > 0 Category exists]
  - Post: [Product created ProductCreated event emitted]
- `update(name?: string, description?: string, price?: Money, categoryId?: CategoryId): void`
  - Pre: [Product exists]
  - Post: [Attributes updated ProductUpdated event emitted]
- `deactivate(): void`
  - Pre: [Product is active]
  - Post: [isActive set to false ProductDeactivated event emitted]
- `addVariant(variant: ProductVariant): void`
  - Pre: [Variant combination is unique]
  - Post: [Variant added ProductVariantAdded event emitted]

**Events:**
- `ProductCreated` (trigger: create() called)
- `ProductUpdated` (trigger: update() called)
- `ProductDeactivated` (trigger: deactivate() called)
- `ProductVariantAdded` (trigger: addVariant() called)

**Relationships:**
- belongs_to ENT-CATEGORY-001 (many-to-one)
- contains ENT-PRODUCTVARIANT-001 (one-to-many)
- references ENT-INVENTORY-001 (one-to-one)

---

### ENT-PRODUCTVARIANT-001 – ProductVariant

**Type:** entity

**Purpose:** Represents a specific variation of a product (e.g., size, color combination)

**Attributes:**

| Name | Type | Constraints |
|------|------|-------------|
| variantId | VariantId (UUID) | Required, immutable |
| size | string | Optional |
| color | string | Optional |
| sku | SKU | Required, unique |
| priceAdjustment | Money | Optional, can be positive or negative |

**Invariants:**
- SKU must be unique across all variants
- At least one variant attribute (size or color) must be specified
- Variant combination must be unique within product

**Relationships:**
- belongs_to ENT-PRODUCT-001 (many-to-one)

---

### ENT-CATEGORY-001 – Category

**Type:** aggregate_root

**Purpose:** Organizes products into logical groups for browsing

**Attributes:**

| Name | Type | Constraints |
|------|------|-------------|
| categoryId | CategoryId (UUID) | Required, immutable |
| name | string | Required, unique, 1-100 characters |
| description | string | Optional |
| parentCategoryId | CategoryId | Optional, for hierarchical categories |

**Invariants:**
- Category name must be unique
- Cannot be its own parent (no circular references)

**Operations:**
- `create(name: string, description?: string, parentId?: CategoryId): Category`
  - Pre: [Name is unique]
  - Post: [Category created CategoryCreated event emitted]

**Events:**
- `CategoryCreated` (trigger: create() called)

**Relationships:**
- contains ENT-PRODUCT-001 (one-to-many)

---

### ENT-INVENTORY-001 – Inventory

**Type:** aggregate_root

**Purpose:** Tracks stock levels for products and manages availability

**Attributes:**

| Name | Type | Constraints |
|------|------|-------------|
| inventoryId | InventoryId (UUID) | Required, immutable |
| productId | ProductId | Required, unique, references Product |
| stockLevel | integer | Required, >= 0 |
| reservedQuantity | integer | Required, >= 0, <= stockLevel |
| availableQuantity | integer | Calculated: stockLevel - reservedQuantity |

**Invariants:**
- Stock level cannot be negative
- Reserved quantity cannot exceed stock level
- Available quantity must equal stock level minus reserved

**Operations:**
- `reserve(quantity: integer): void`
  - Pre: [quantity <= availableQuantity]
  - Post: [reservedQuantity increased InventoryReserved event emitted]
- `release(quantity: integer): void`
  - Pre: [quantity <= reservedQuantity]
  - Post: [reservedQuantity decreased InventoryReleased event emitted]
- `deduct(quantity: integer): void`
  - Pre: [quantity <= reservedQuantity]
  - Post: [stockLevel and reservedQuantity decreased InventoryDeducted event emitted]
- `restock(quantity: integer): void`
  - Pre: [quantity > 0]
  - Post: [stockLevel increased InventoryRestocked event emitted]

**Events:**
- `InventoryReserved` (trigger: reserve() called)
- `InventoryReleased` (trigger: release() called)
- `InventoryDeducted` (trigger: deduct() called)
- `InventoryRestocked` (trigger: restock() called)
- `OutOfStock` (trigger: availableQuantity reaches 0)

**Relationships:**
- belongs_to ENT-PRODUCT-001 (one-to-one)

---

## Value Objects

### VO-MONEY-001 – Money

**Purpose:** Represents a monetary amount with currency for price calculations

**Attributes:**
- `amount` (decimal): >= 0, precision 2
- `currency` (string): ISO 4217 code (default: USD)

**Operations:** [add subtract multiply equals greaterThan lessThan]

---

### VO-ADDRESS-001 – Address

**Purpose:** Represents a shipping address for order delivery

**Attributes:**
- `street` (string): Required, 1-200 characters
- `city` (string): Required, 1-100 characters
- `state` (string): Required
- `postalCode` (string): Required, valid format
- `country` (string): Required, ISO 3166-1

**Operations:** [equals format]

---

### VO-PAYMENTMETHOD-001 – PaymentMethod

**Purpose:** Represents the payment method chosen for an order

**Attributes:**
- `type` (PaymentType (enum)): Required, values: credit_card, paypal
- `lastFourDigits` (string): Optional, for credit card display
- `paypalEmail` (string): Optional, for PayPal identification

**Operations:** [equals displayName]

---

### VO-QUANTITY-001 – Quantity

**Purpose:** Represents a positive quantity of items

**Attributes:**
- `value` (integer): > 0

**Operations:** [add subtract equals greaterThan]

---

### VO-EMAIL-001 – Email

**Purpose:** Represents a validated email address

**Attributes:**
- `value` (string): Required, valid email format (RFC 5322)

**Operations:** [equals domain]

---

### VO-SKU-001 – SKU

**Purpose:** Stock Keeping Unit identifier for products and variants

**Attributes:**
- `value` (string): Required, alphanumeric, unique

**Operations:** [equals]

---

### VO-ORDERSTATUS-001 – OrderStatus

**Purpose:** Represents the current state of an order in its lifecycle

**Attributes:**
- `value` (enum): Values: pending, confirmed, shipped, delivered, cancelled

**Operations:** [canTransitionTo equals]

---

## Summary

- Aggregate Roots: 5
- Entities: 4
- Value Objects: 7
- Total Operations: 24
- Total Events: 18
