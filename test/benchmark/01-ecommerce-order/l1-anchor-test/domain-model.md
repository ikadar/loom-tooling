# Domain Model

Generated: 2025-12-31T13:11:10+01:00

**Domain:** E-Commerce Order System

Domain model for an online store supporting product browsing, shopping cart, order placement, and order management

---

## Entities

### ENT-CUSTOMER-001 – Customer {#ent-customer-001}

**Type:** aggregate_root

**Purpose:** Represents a registered user who can browse products, maintain a cart, and place orders

**Attributes:**
| Name | Type | Constraints |
|------|------|-------------|
| customerId | CustomerId (UUID) | Required, immutable after creation |
| email | Email | Required, unique, valid email format |
| registrationStatus | RegistrationStatus (enum) | Required, values: REGISTERED, UNREGISTERED |
| cart | Cart | Required, one per customer |

**Invariants:**
- Customer must have a valid email address
- Customer must be registered to place orders
- Each customer has exactly one active cart

**Operations:**
- `register(email: Email, password: Password): void`
  - Pre: [Email not already registered]
  - Post: [Status set to REGISTERED CustomerRegistered event emitted]
- `placeOrder(shippingAddress: ShippingAddress, paymentMethod: PaymentMethod): Order`
  - Pre: [Status is REGISTERED Cart has items All cart items in stock]
  - Post: [Order created Cart emptied OrderPlaced event emitted]

**Events:**
- `CustomerRegistered` (trigger: register() completed)

**Relationships:**
- contains ENT-CART-001 (one-to-one)
- has_many ENT-ORDER-001 (one-to-many)

---

### ENT-CART-001 – Cart {#ent-cart-001}

**Type:** aggregate_root

**Purpose:** Holds products a customer intends to purchase before checkout

**Attributes:**
| Name | Type | Constraints |
|------|------|-------------|
| cartId | CartId (UUID) | Required, immutable after creation |
| customerId | CustomerId | Required, references owning customer |
| items | List<CartItem> | Can be empty |
| totalPrice | Money | >= 0, calculated from items |

**Invariants:**
- Total price equals sum of all item subtotals
- Each product can appear only once (quantities aggregated)
- Cannot add out-of-stock products

**Operations:**
- `addItem(product: Product, quantity: number): void`
  - Pre: [Product is in stock Quantity > 0]
  - Post: [CartItem added or quantity updated Total recalculated ItemAddedToCart event emitted]
- `updateQuantity(productId: ProductId, quantity: number): void`
  - Pre: [Item exists in cart Quantity > 0 Stock available]
  - Post: [Quantity updated Total recalculated CartQuantityUpdated event emitted]
- `removeItem(productId: ProductId): void`
  - Pre: [Item exists in cart]
  - Post: [Item removed Total recalculated ItemRemovedFromCart event emitted]
- `clear(): void`
  - Post: [All items removed Total set to zero CartCleared event emitted]

**Events:**
- `ItemAddedToCart` (trigger: addItem() called)
- `CartQuantityUpdated` (trigger: updateQuantity() called)
- `ItemRemovedFromCart` (trigger: removeItem() called)
- `CartCleared` (trigger: clear() called or order placed)

**Relationships:**
- contains ENT-CARTITEM-001 (one-to-many)
- belongs_to ENT-CUSTOMER-001 (many-to-one)

---

### ENT-CARTITEM-001 – CartItem {#ent-cartitem-001}

**Type:** entity

**Purpose:** Represents a product and quantity within a shopping cart

**Attributes:**
| Name | Type | Constraints |
|------|------|-------------|
| cartItemId | CartItemId (UUID) | Required, immutable after creation |
| productId | ProductId | Required, references product |
| productName | string | Required, snapshot at time of add |
| unitPrice | Money | Required, snapshot at time of add |
| quantity | integer | Required, >= 1 |
| subtotal | Money | Calculated: unitPrice * quantity |

**Invariants:**
- Quantity must be at least 1
- Subtotal equals unitPrice multiplied by quantity

**Operations:**
- `updateQuantity(newQuantity: number): void`
  - Pre: [newQuantity >= 1]
  - Post: [Quantity updated Subtotal recalculated]

**Relationships:**
- belongs_to ENT-CART-001 (many-to-one)
- references ENT-PRODUCT-001 (many-to-one)

---

### ENT-ORDER-001 – Order {#ent-order-001}

**Type:** aggregate_root

**Purpose:** Represents a customer's purchase transaction with line items, shipping, and payment details

**Attributes:**
| Name | Type | Constraints |
|------|------|-------------|
| orderId | OrderId (UUID) | Required, immutable after creation |
| customerId | CustomerId | Required, references ordering customer |
| status | OrderStatus (enum) | Required, values: PENDING, CONFIRMED, SHIPPED, DELIVERED, CANCELLED |
| lineItems | List<OrderLineItem> | Required, at least one item |
| shippingAddress | ShippingAddress | Required |
| paymentMethod | PaymentMethod | Required |
| subtotal | Money | >= 0, sum of line items |
| shippingCost | Money | >= 0, free if subtotal > $50 |
| totalAmount | Money | >= 0, subtotal + shippingCost |
| createdAt | DateTime | Required, immutable |

**Invariants:**
- Order must have at least one line item
- Total amount equals subtotal plus shipping cost
- Shipping is free when subtotal exceeds $50
- Cannot modify order after status is SHIPPED or DELIVERED
- Can only cancel order before SHIPPED status
- Status transitions must follow valid state machine

**Operations:**
- `confirm(): void`
  - Pre: [Status is PENDING]
  - Post: [Status changes to CONFIRMED OrderConfirmed event emitted]
- `ship(trackingNumber: string): void`
  - Pre: [Status is CONFIRMED]
  - Post: [Status changes to SHIPPED OrderShipped event emitted]
- `deliver(): void`
  - Pre: [Status is SHIPPED]
  - Post: [Status changes to DELIVERED OrderDelivered event emitted]
- `cancel(reason: string): void`
  - Pre: [Status is PENDING or CONFIRMED]
  - Post: [Status changes to CANCELLED OrderCancelled event emitted Inventory restored]
- `updateStatus(newStatus: OrderStatus): void`
  - Pre: [Transition is valid per state machine]
  - Post: [Status updated Appropriate event emitted]

**Events:**
- `OrderPlaced` (trigger: Order created from cart)
- `OrderConfirmed` (trigger: confirm() called)
- `OrderShipped` (trigger: ship() called)
- `OrderDelivered` (trigger: deliver() called)
- `OrderCancelled` (trigger: cancel() called)

**Relationships:**
- belongs_to ENT-CUSTOMER-001 (many-to-one)
- contains ENT-ORDERLINEITEM-001 (one-to-many)
- contains VO-SHIPPINGADDRESS-001 (one-to-one)
- contains VO-PAYMENTMETHOD-001 (one-to-one)

---

### ENT-ORDERLINEITEM-001 – OrderLineItem {#ent-orderlineitem-001}

**Type:** entity

**Purpose:** Immutable snapshot of a product at time of order with quantity and pricing

**Attributes:**
| Name | Type | Constraints |
|------|------|-------------|
| lineItemId | LineItemId (UUID) | Required, immutable |
| productId | ProductId | Required, reference to original product |
| productName | string | Required, snapshot at order time |
| unitPrice | Money | Required, snapshot at order time |
| quantity | integer | Required, >= 1 |
| subtotal | Money | Calculated: unitPrice * quantity |

**Invariants:**
- All attributes are immutable after creation
- Quantity must be at least 1
- Subtotal equals unitPrice multiplied by quantity

**Relationships:**
- belongs_to ENT-ORDER-001 (many-to-one)
- references ENT-PRODUCT-001 (many-to-one)

---

### ENT-PRODUCT-001 – Product {#ent-product-001}

**Type:** aggregate_root

**Purpose:** Represents an item available for sale in the store

**Attributes:**
| Name | Type | Constraints |
|------|------|-------------|
| productId | ProductId (UUID) | Required, immutable after creation |
| name | string | Required, 1-200 characters |
| description | string | Optional, max 5000 characters |
| price | Money | Required, > 0 |
| categoryId | CategoryId | Required, references category |
| variants | List<ProductVariant> | Optional, for products with size/color options |
| isActive | boolean | Required, determines visibility |

**Invariants:**
- Price must be greater than zero
- Product must belong to a category
- Inactive products cannot be added to carts

**Operations:**
- `create(name: string, description: string, price: Money, categoryId: CategoryId): Product`
  - Pre: [Price > 0 Category exists]
  - Post: [Product created as active ProductCreated event emitted]
- `update(name?: string, description?: string, price?: Money, categoryId?: CategoryId): void`
  - Pre: [If price provided, price > 0]
  - Post: [Product updated ProductUpdated event emitted]
- `deactivate(): void`
  - Pre: [Product is active]
  - Post: [isActive set to false ProductDeactivated event emitted]
- `addVariant(variant: ProductVariant): void`
  - Pre: [Variant combination is unique]
  - Post: [Variant added to product]

**Events:**
- `ProductCreated` (trigger: create() called)
- `ProductUpdated` (trigger: update() called)
- `ProductDeactivated` (trigger: deactivate() called)

**Relationships:**
- belongs_to ENT-CATEGORY-001 (many-to-one)
- contains ENT-PRODUCTVARIANT-001 (one-to-many)
- has_one ENT-INVENTORY-001 (one-to-one)

---

### ENT-PRODUCTVARIANT-001 – ProductVariant {#ent-productvariant-001}

**Type:** entity

**Purpose:** Represents a specific variation of a product (e.g., size, color combination)

**Attributes:**
| Name | Type | Constraints |
|------|------|-------------|
| variantId | VariantId (UUID) | Required, immutable |
| sku | string | Required, unique across all variants |
| size | string | Optional |
| color | string | Optional |
| priceAdjustment | Money | Optional, can be positive or negative |

**Invariants:**
- SKU must be unique
- At least one variant attribute (size or color) must be specified

**Relationships:**
- belongs_to ENT-PRODUCT-001 (many-to-one)

---

### ENT-CATEGORY-001 – Category {#ent-category-001}

**Type:** aggregate_root

**Purpose:** Organizes products into browsable groupings

**Attributes:**
| Name | Type | Constraints |
|------|------|-------------|
| categoryId | CategoryId (UUID) | Required, immutable |
| name | string | Required, unique, 1-100 characters |
| description | string | Optional |
| parentCategoryId | CategoryId | Optional, for nested categories |

**Invariants:**
- Category name must be unique
- Cannot be its own parent

**Operations:**
- `create(name: string, description?: string, parentId?: CategoryId): Category`
  - Pre: [Name is unique]
  - Post: [Category created CategoryCreated event emitted]

**Events:**
- `CategoryCreated` (trigger: create() called)

**Relationships:**
- has_many ENT-PRODUCT-001 (one-to-many)
- parent_of ENT-CATEGORY-001 (one-to-many)

---

### ENT-INVENTORY-001 – Inventory {#ent-inventory-001}

**Type:** aggregate_root

**Purpose:** Tracks stock levels for products and manages availability

**Attributes:**
| Name | Type | Constraints |
|------|------|-------------|
| inventoryId | InventoryId (UUID) | Required, immutable |
| productId | ProductId | Required, unique, references product |
| quantityOnHand | integer | Required, >= 0 |
| reservedQuantity | integer | Required, >= 0 |
| availableQuantity | integer | Calculated: quantityOnHand - reservedQuantity |

**Invariants:**
- Available quantity cannot be negative
- Reserved quantity cannot exceed quantity on hand
- Each product has exactly one inventory record

**Operations:**
- `reserve(quantity: number): void`
  - Pre: [Available quantity >= requested quantity]
  - Post: [Reserved quantity increased InventoryReserved event emitted]
- `release(quantity: number): void`
  - Pre: [Reserved quantity >= requested quantity]
  - Post: [Reserved quantity decreased InventoryReleased event emitted]
- `deduct(quantity: number): void`
  - Pre: [Quantity on hand >= requested quantity]
  - Post: [Quantity on hand decreased InventoryDeducted event emitted]
- `restock(quantity: number): void`
  - Pre: [Quantity > 0]
  - Post: [Quantity on hand increased InventoryRestocked event emitted]
- `isInStock(quantity: number): boolean`
  - Post: [Returns true if available quantity >= requested]

**Events:**
- `InventoryReserved` (trigger: reserve() called)
- `InventoryReleased` (trigger: release() called)
- `InventoryDeducted` (trigger: deduct() called)
- `InventoryRestocked` (trigger: restock() called)
- `OutOfStock` (trigger: Available quantity reaches 0)

**Relationships:**
- tracks ENT-PRODUCT-001 (one-to-one)

---

## Value Objects

### VO-MONEY-001 – Money {#vo-money-001}

**Purpose:** Represents a monetary value with currency for precise financial calculations

**Attributes:**
- `amount` (decimal): Required, precision 2 decimal places
- `currency` (string): Required, ISO 4217 code (default: USD)

**Operations:** [add subtract multiply divide equals greaterThan lessThan]

---

### VO-SHIPPINGADDRESS-001 – ShippingAddress {#vo-shippingaddress-001}

**Purpose:** Represents a complete shipping destination for order delivery

**Attributes:**
- `street` (string): Required, 1-200 characters
- `city` (string): Required, 1-100 characters
- `state` (string): Required, 2-100 characters
- `postalCode` (string): Required, valid postal format
- `country` (string): Required, ISO 3166 code

**Operations:** [equals format validate]

---

### VO-PAYMENTMETHOD-001 – PaymentMethod {#vo-paymentmethod-001}

**Purpose:** Represents the payment instrument used for an order

**Attributes:**
- `type` (PaymentType (enum)): Required, values: CREDIT_CARD, PAYPAL
- `lastFourDigits` (string): Required for CREDIT_CARD, 4 digits
- `expiryMonth` (integer): Required for CREDIT_CARD, 1-12
- `expiryYear` (integer): Required for CREDIT_CARD
- `paypalEmail` (Email): Required for PAYPAL

**Operations:** [equals isExpired mask]

---

### VO-EMAIL-001 – Email {#vo-email-001}

**Purpose:** Represents a validated email address

**Attributes:**
- `value` (string): Required, valid email format per RFC 5322

**Operations:** [equals getDomain validate]

---

### VO-ORDERID-001 – OrderId {#vo-orderid-001}

**Purpose:** Strongly typed identifier for orders

**Attributes:**
- `value` (UUID): Required, immutable

**Operations:** [equals toString]

---

### VO-PRODUCTID-001 – ProductId {#vo-productid-001}

**Purpose:** Strongly typed identifier for products

**Attributes:**
- `value` (UUID): Required, immutable

**Operations:** [equals toString]

---

### VO-CUSTOMERID-001 – CustomerId {#vo-customerid-001}

**Purpose:** Strongly typed identifier for customers

**Attributes:**
- `value` (UUID): Required, immutable

**Operations:** [equals toString]

---

## Summary

- Aggregate Roots: 6
- Entities: 10
- Value Objects: 7
- Total Operations: 28
- Total Events: 18
