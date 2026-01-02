# Event & Message Design

Generated: 2026-01-02T19:24:29+01:00

---

## Domain Events

### EVT-CUST-001: CustomerRegistered

**Purpose:** Signals new customer account created

**Trigger:** Customer submits registration form

**Aggregate:** Customer | **Version:** 1.0

**Payload:**
- `customerId`: UUID
- `email`: Email
- `status`: CustomerStatus
- `registeredAt`: DateTime

**Invariants Reflected:** [Email is unique Password meets strength requirements]

**Consumers:** [EmailService CartService]

---

### EVT-CART-001: ItemAddedToCart

**Purpose:** Signals item was added to shopping cart

**Trigger:** Customer adds product to cart

**Aggregate:** Cart | **Version:** 1.0

**Payload:**
- `cartId`: UUID
- `customerId`: UUID
- `productId`: UUID
- `variantId`: UUID
- `quantity`: Integer
- `unitPrice`: Money
- `addedAt`: DateTime

**Invariants Reflected:** [Product in stock Quantity available]

**Consumers:** [AnalyticsService]

---

### EVT-CART-002: CartQuantityUpdated

**Purpose:** Signals cart item quantity changed

**Trigger:** Customer changes quantity in cart

**Aggregate:** Cart | **Version:** 1.0

**Payload:**
- `cartId`: UUID
- `productId`: UUID
- `previousQuantity`: Integer
- `newQuantity`: Integer
- `updatedAt`: DateTime

**Invariants Reflected:** [New quantity available in stock]

**Consumers:** [AnalyticsService]

---

### EVT-CART-003: ItemRemovedFromCart

**Purpose:** Signals item was removed from cart

**Trigger:** Customer removes item from cart

**Aggregate:** Cart | **Version:** 1.0

**Payload:**
- `cartId`: UUID
- `productId`: UUID
- `variantId`: UUID
- `removedAt`: DateTime

**Invariants Reflected:** [Item existed in cart]

**Consumers:** [AnalyticsService]

---

### EVT-CART-004: CartCleared

**Purpose:** Signals all items removed from cart

**Trigger:** Order placed or cart expired

**Aggregate:** Cart | **Version:** 1.0

**Payload:**
- `cartId`: UUID
- `customerId`: UUID
- `clearedAt`: DateTime

**Invariants Reflected:** [Cart was not empty]

**Consumers:** [AnalyticsService]

---

### EVT-INV-001: InventoryReserved

**Purpose:** Signals inventory reserved for order

**Trigger:** Customer places order

**Aggregate:** Inventory | **Version:** 1.0

**Payload:**
- `inventoryId`: UUID
- `productId`: UUID
- `orderId`: UUID
- `reservedQuantity`: Integer
- `reservedAt`: DateTime

**Invariants Reflected:** [Sufficient stock available]

**Consumers:** [OrderService NotificationService]

---

### EVT-ORD-001: OrderPlaced

**Purpose:** Signals new order was created

**Trigger:** Customer completes checkout

**Aggregate:** Order | **Version:** 1.0

**Payload:**
- `orderId`: UUID
- `customerId`: UUID
- `lineItems`: array
- `totalAmount`: Money
- `shippingAddress`: Address
- `paymentAuthorizationId`: String
- `placedAt`: DateTime

**Invariants Reflected:** [Cart not empty Payment authorized Inventory reserved]

**Consumers:** [NotificationService InventoryService AnalyticsService]

---

### EVT-ORD-002: OrderConfirmed

**Purpose:** Signals order payment captured and confirmed

**Trigger:** Payment captured successfully

**Aggregate:** Order | **Version:** 1.0

**Payload:**
- `orderId`: UUID
- `confirmedAt`: DateTime

**Invariants Reflected:** [Order was in PENDING status Payment captured]

**Consumers:** [NotificationService WarehouseService]

---

### EVT-ORD-003: OrderShipped

**Purpose:** Signals order was shipped to customer

**Trigger:** Warehouse marks order as shipped

**Aggregate:** Order | **Version:** 1.0

**Payload:**
- `orderId`: UUID
- `trackingNumber`: String
- `carrier`: String
- `shippedAt`: DateTime

**Invariants Reflected:** [Order was in CONFIRMED status]

**Consumers:** [NotificationService InventoryService]

---

### EVT-ORD-004: OrderDelivered

**Purpose:** Signals order was delivered to customer

**Trigger:** Carrier confirms delivery

**Aggregate:** Order | **Version:** 1.0

**Payload:**
- `orderId`: UUID
- `deliveredAt`: DateTime

**Invariants Reflected:** [Order was in SHIPPED status]

**Consumers:** [NotificationService AnalyticsService]

---

### EVT-ORD-005: OrderCancelled

**Purpose:** Signals order was cancelled

**Trigger:** Customer requests cancellation

**Aggregate:** Order | **Version:** 1.0

**Payload:**
- `orderId`: UUID
- `reason`: String
- `cancelledAt`: DateTime

**Invariants Reflected:** [Order not yet shipped]

**Consumers:** [NotificationService InventoryService PaymentService]

---

### EVT-INV-002: InventoryDeducted

**Purpose:** Signals inventory was deducted after shipping

**Trigger:** Order shipped

**Aggregate:** Inventory | **Version:** 1.0

**Payload:**
- `inventoryId`: UUID
- `productId`: UUID
- `orderId`: UUID
- `deductedQuantity`: Integer
- `deductedAt`: DateTime

**Invariants Reflected:** [Inventory was reserved]

**Consumers:** [AnalyticsService]

---

### EVT-INV-003: InventoryReleased

**Purpose:** Signals reserved inventory was released

**Trigger:** Order cancelled

**Aggregate:** Inventory | **Version:** 1.0

**Payload:**
- `inventoryId`: UUID
- `productId`: UUID
- `orderId`: UUID
- `releasedQuantity`: Integer
- `releasedAt`: DateTime

**Invariants Reflected:** [Inventory was reserved for order]

**Consumers:** [AnalyticsService]

---

### EVT-INV-004: InventoryRestocked

**Purpose:** Signals inventory was restocked

**Trigger:** Warehouse receives shipment

**Aggregate:** Inventory | **Version:** 1.0

**Payload:**
- `inventoryId`: UUID
- `productId`: UUID
- `addedQuantity`: Integer
- `newTotalQuantity`: Integer
- `restockedAt`: DateTime

**Invariants Reflected:** [Quantity is positive]

**Consumers:** [NotificationService AnalyticsService]

---

### EVT-INV-005: OutOfStock

**Purpose:** Signals product is out of stock

**Trigger:** Available quantity reaches zero

**Aggregate:** Inventory | **Version:** 1.0

**Payload:**
- `inventoryId`: UUID
- `productId`: UUID
- `occurredAt`: DateTime

**Invariants Reflected:** [Available quantity is zero]

**Consumers:** [NotificationService ProductService]

---

### EVT-PROD-001: ProductCreated

**Purpose:** Signals new product added to catalog

**Trigger:** Admin creates product

**Aggregate:** Product | **Version:** 1.0

**Payload:**
- `productId`: UUID
- `name`: String
- `price`: Money
- `categoryId`: UUID
- `createdAt`: DateTime

**Invariants Reflected:** [Price >= $0.01 Name length valid Category exists]

**Consumers:** [InventoryService SearchService]

---

### EVT-PROD-002: ProductDeactivated

**Purpose:** Signals product was deactivated

**Trigger:** Admin deactivates product

**Aggregate:** Product | **Version:** 1.0

**Payload:**
- `productId`: UUID
- `deactivatedAt`: DateTime

**Invariants Reflected:** [Product existed and was active]

**Consumers:** [SearchService CartService]

---

### EVT-PROD-003: ProductUpdated

**Purpose:** Signals product details were updated

**Trigger:** Admin updates product

**Aggregate:** Product | **Version:** 1.0

**Payload:**
- `productId`: UUID
- `updatedFields`: array
- `updatedAt`: DateTime

**Invariants Reflected:** [Price >= $0.01 if changed]

**Consumers:** [SearchService AnalyticsService]

---

### EVT-CAT-001: CategoryCreated

**Purpose:** Signals new category added to catalog

**Trigger:** Admin creates category

**Aggregate:** Category | **Version:** 1.0

**Payload:**
- `categoryId`: UUID
- `name`: String
- `parentId`: UUID
- `createdAt`: DateTime

**Invariants Reflected:** [Name is unique Hierarchy depth <= 3]

**Consumers:** [SearchService]

---

## Commands

### CMD-CUST-001: RegisterCustomer

**Intent:** Create new customer account

**Expected Outcome:** Customer created with REGISTERED status

**Required Data:**
- `email`: Email
- `password`: String
- `name`: String

**Failure Conditions:** [Email already registered Password too weak]

---

### CMD-CART-001: AddItemToCart

**Intent:** Add product to shopping cart

**Expected Outcome:** Item added to cart with recalculated total

**Required Data:**
- `customerId`: UUID
- `productId`: UUID
- `variantId`: UUID
- `quantity`: Integer

**Failure Conditions:** [Out of stock Quantity exceeds stock Variant required]

---

### CMD-CART-002: UpdateCartQuantity

**Intent:** Change quantity of item in cart

**Expected Outcome:** Cart quantity updated with recalculated total

**Required Data:**
- `cartId`: UUID
- `productId`: UUID
- `newQuantity`: Integer

**Failure Conditions:** [Quantity exceeds stock Item not in cart]

---

### CMD-CART-003: RemoveItemFromCart

**Intent:** Remove product from shopping cart

**Expected Outcome:** Item removed from cart

**Required Data:**
- `cartId`: UUID
- `productId`: UUID

**Failure Conditions:** [Item not in cart]

---

### CMD-ORD-001: PlaceOrder

**Intent:** Complete checkout and create order

**Expected Outcome:** Order created in PENDING status

**Required Data:**
- `customerId`: UUID
- `shippingAddress`: Address
- `paymentMethod`: PaymentMethod

**Failure Conditions:** [Cart empty Insufficient stock Payment failed Invalid address]

---

### CMD-ORD-002: ConfirmOrder

**Intent:** Confirm order after payment capture

**Expected Outcome:** Order status changed to CONFIRMED

**Required Data:**
- `orderId`: UUID

**Failure Conditions:** [Payment capture failed Invalid status transition]

---

### CMD-ORD-003: ShipOrder

**Intent:** Mark order as shipped with tracking

**Expected Outcome:** Order status changed to SHIPPED

**Required Data:**
- `orderId`: UUID
- `trackingNumber`: String
- `carrier`: String

**Failure Conditions:** [Invalid status transition Inventory deduction failed]

---

### CMD-ORD-004: DeliverOrder

**Intent:** Mark order as delivered

**Expected Outcome:** Order status changed to DELIVERED

**Required Data:**
- `orderId`: UUID
- `trackingNumber`: String

**Failure Conditions:** [Order not found Invalid status transition]

---

### CMD-ORD-005: CancelOrder

**Intent:** Cancel order and process refund

**Expected Outcome:** Order cancelled with refund initiated

**Required Data:**
- `orderId`: UUID
- `reason`: String

**Failure Conditions:** [Order already shipped Refund processing failed]

---

### CMD-PROD-001: CreateProduct

**Intent:** Add new product to catalog

**Expected Outcome:** Product created and active

**Required Data:**
- `name`: String
- `description`: String
- `price`: Money
- `categoryId`: UUID
- `initialStock`: Integer

**Failure Conditions:** [Invalid price Invalid name Category not found]

---

### CMD-PROD-002: DeactivateProduct

**Intent:** Soft delete product from catalog

**Expected Outcome:** Product deactivated

**Required Data:**
- `productId`: UUID

---

### CMD-PROD-003: UpdateProductPrice

**Intent:** Update product pricing

**Expected Outcome:** Product price updated

**Required Data:**
- `productId`: UUID
- `newPrice`: Money

**Failure Conditions:** [Invalid price Product not found]

---

### CMD-PROD-004: AddProductVariant

**Intent:** Add size/color variant to product

**Expected Outcome:** Variant added with inventory

**Required Data:**
- `productId`: UUID
- `sku`: String
- `size`: String
- `color`: String
- `priceModifier`: Money

**Failure Conditions:** [Duplicate SKU Product not found No attribute specified]

---

### CMD-INV-001: RestockInventory

**Intent:** Add stock to product inventory

**Expected Outcome:** Inventory increased

**Required Data:**
- `productId`: UUID
- `quantity`: Integer

**Failure Conditions:** [Invalid quantity Product not found]

---

### CMD-INV-002: ReserveInventory

**Intent:** Reserve stock for order

**Expected Outcome:** Inventory reserved for order

**Required Data:**
- `productId`: UUID
- `orderId`: UUID
- `quantity`: Integer

**Failure Conditions:** [Insufficient stock]

---

### CMD-CAT-001: CreateCategory

**Intent:** Add new product category

**Expected Outcome:** Category created in hierarchy

**Required Data:**
- `name`: String
- `parentId`: UUID

**Failure Conditions:** [Duplicate name Depth exceeded Parent not found]

---

## Integration Events

### INT-ORD-001: OrderPlacedNotification

**Purpose:** Notify external systems of new order

**Source:** OrderService

**Consumers:** [EmailService AnalyticsService]

**Payload:** [orderId customerId totalAmount lineItems]

---

### INT-ORD-002: OrderConfirmedNotification

**Purpose:** Notify systems order is confirmed

**Source:** OrderService

**Consumers:** [EmailService WarehouseService]

**Payload:** [orderId customerId shippingAddress]

---

### INT-ORD-003: OrderShippedNotification

**Purpose:** Notify customer of shipment

**Source:** OrderService

**Consumers:** [EmailService SMSService]

**Payload:** [orderId customerId trackingNumber carrier]

---

### INT-ORD-004: OrderDeliveredNotification

**Purpose:** Notify customer of delivery

**Source:** OrderService

**Consumers:** [EmailService]

**Payload:** [orderId customerId]

---

### INT-ORD-005: OrderCancelledNotification

**Purpose:** Notify customer of cancellation

**Source:** OrderService

**Consumers:** [EmailService PaymentService]

**Payload:** [orderId customerId refundAmount]

---

### INT-CUST-001: CustomerRegistrationNotification

**Purpose:** Send verification email to customer

**Source:** CustomerService

**Consumers:** [EmailService]

**Payload:** [customerId email verificationToken]

---

### INT-INV-001: OutOfStockAlert

**Purpose:** Alert operations of out-of-stock product

**Source:** InventoryService

**Consumers:** [EmailService SlackService]

**Payload:** [productId productName]

---

