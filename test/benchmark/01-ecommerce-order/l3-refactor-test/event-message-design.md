# Event & Message Design

Generated: 2026-01-02T21:37:12+01:00

---

## Domain Events

### EVT-CUST-001: CustomerRegistered

**Purpose:** Records that a new customer has completed registration

**Trigger:** Customer.register() completed successfully

**Aggregate:** Customer | **Version:** 1.0

**Payload:**
- `customerId`: CustomerId
- `email`: Email
- `registrationStatus`: RegistrationStatus
- `registeredAt`: DateTime

**Invariants Reflected:** [Customer must have valid email Customer status set to REGISTERED]

**Consumers:** [NotificationService CartService AnalyticsService]

---

### EVT-CART-001: ItemAddedToCart

**Purpose:** Records that a product was added to customer's cart

**Trigger:** Cart.addItem() called with valid product

**Aggregate:** Cart | **Version:** 1.0

**Payload:**
- `cartId`: CartId
- `customerId`: CustomerId
- `productId`: ProductId
- `productName`: string
- `unitPrice`: Money
- `quantity`: integer
- `newTotalPrice`: Money

**Invariants Reflected:** [Product must be in stock Each product appears only once in cart]

**Consumers:** [AnalyticsService RecommendationService]

---

### EVT-CART-002: CartQuantityUpdated

**Purpose:** Records that item quantity in cart was changed

**Trigger:** Cart.updateQuantity() called

**Aggregate:** Cart | **Version:** 1.0

**Payload:**
- `cartId`: CartId
- `productId`: ProductId
- `previousQuantity`: integer
- `newQuantity`: integer
- `newTotalPrice`: Money

**Invariants Reflected:** [Quantity must be at least 1 Stock must be available]

**Consumers:** [AnalyticsService]

---

### EVT-CART-003: ItemRemovedFromCart

**Purpose:** Records that an item was removed from cart

**Trigger:** Cart.removeItem() called

**Aggregate:** Cart | **Version:** 1.0

**Payload:**
- `cartId`: CartId
- `productId`: ProductId
- `removedQuantity`: integer
- `newTotalPrice`: Money

**Invariants Reflected:** [Total price recalculated after removal]

**Consumers:** [AnalyticsService RecommendationService]

---

### EVT-CART-004: CartCleared

**Purpose:** Records that all items were removed from cart

**Trigger:** Cart.clear() called or order placed

**Aggregate:** Cart | **Version:** 1.0

**Payload:**
- `cartId`: CartId
- `customerId`: CustomerId
- `reason`: string

**Invariants Reflected:** [Cart items empty after clear Total price zero after clear]

**Consumers:** [AnalyticsService]

---

### EVT-ORDER-001: OrderPlaced

**Purpose:** Records that customer placed a new order

**Trigger:** Order created from cart checkout

**Aggregate:** Order | **Version:** 1.0

**Payload:**
- `orderId`: OrderId
- `customerId`: CustomerId
- `lineItems`: List<OrderLineItem>
- `shippingAddress`: ShippingAddress
- `paymentMethod`: PaymentMethod
- `subtotal`: Money
- `shippingCost`: Money
- `totalAmount`: Money
- `createdAt`: DateTime

**Invariants Reflected:** [Order has at least one line item Shipping free when subtotal > $50]

**Consumers:** [InventoryService PaymentService NotificationService AnalyticsService]

---

### EVT-ORDER-002: OrderConfirmed

**Purpose:** Records that order was confirmed after payment capture

**Trigger:** Order.confirm() called after payment success

**Aggregate:** Order | **Version:** 1.0

**Payload:**
- `orderId`: OrderId
- `customerId`: CustomerId
- `confirmedAt`: DateTime

**Invariants Reflected:** [Status transition from PENDING to CONFIRMED]

**Consumers:** [FulfillmentService NotificationService AnalyticsService]

---

### EVT-ORDER-003: OrderShipped

**Purpose:** Records that order was shipped with tracking info

**Trigger:** Order.ship() called by fulfillment

**Aggregate:** Order | **Version:** 1.0

**Payload:**
- `orderId`: OrderId
- `customerId`: CustomerId
- `trackingNumber`: string
- `carrier`: string
- `shippedAt`: DateTime

**Invariants Reflected:** [Status transition from CONFIRMED to SHIPPED]

**Consumers:** [NotificationService AnalyticsService]

---

### EVT-ORDER-004: OrderDelivered

**Purpose:** Records that order was delivered to customer

**Trigger:** Order.deliver() called after carrier confirmation

**Aggregate:** Order | **Version:** 1.0

**Payload:**
- `orderId`: OrderId
- `customerId`: CustomerId
- `deliveredAt`: DateTime

**Invariants Reflected:** [Status transition from SHIPPED to DELIVERED]

**Consumers:** [NotificationService AnalyticsService ReviewService]

---

### EVT-ORDER-005: OrderCancelled

**Purpose:** Records that order was cancelled before shipment

**Trigger:** Order.cancel() called by customer or system

**Aggregate:** Order | **Version:** 1.0

**Payload:**
- `orderId`: OrderId
- `customerId`: CustomerId
- `reason`: string
- `cancelledAt`: DateTime

**Invariants Reflected:** [Can only cancel before SHIPPED status Inventory must be restored]

**Consumers:** [InventoryService PaymentService NotificationService AnalyticsService]

---

### EVT-PROD-001: ProductCreated

**Purpose:** Records that a new product was added to catalog

**Trigger:** Product.create() called by admin

**Aggregate:** Product | **Version:** 1.0

**Payload:**
- `productId`: ProductId
- `name`: string
- `description`: string
- `price`: Money
- `categoryId`: CategoryId
- `isActive`: boolean

**Invariants Reflected:** [Price must be greater than zero Product must belong to category]

**Consumers:** [InventoryService SearchService AnalyticsService]

---

### EVT-PROD-002: ProductUpdated

**Purpose:** Records that product details were modified

**Trigger:** Product.update() called

**Aggregate:** Product | **Version:** 1.0

**Payload:**
- `productId`: ProductId
- `changedFields`: Map<string,any>
- `updatedAt`: DateTime

**Invariants Reflected:** [Price remains greater than zero if changed]

**Consumers:** [SearchService CacheService]

---

### EVT-PROD-003: ProductDeactivated

**Purpose:** Records that product was removed from active catalog

**Trigger:** Product.deactivate() called

**Aggregate:** Product | **Version:** 1.0

**Payload:**
- `productId`: ProductId
- `deactivatedAt`: DateTime

**Invariants Reflected:** [Inactive products cannot be added to carts]

**Consumers:** [SearchService CartService AnalyticsService]

---

### EVT-CAT-001: CategoryCreated

**Purpose:** Records that a new product category was created

**Trigger:** Category.create() called

**Aggregate:** Category | **Version:** 1.0

**Payload:**
- `categoryId`: CategoryId
- `name`: string
- `parentCategoryId`: CategoryId

**Invariants Reflected:** [Category name must be unique Cannot be its own parent]

**Consumers:** [SearchService NavigationService]

---

### EVT-INV-001: InventoryReserved

**Purpose:** Records that stock was reserved for an order

**Trigger:** Inventory.reserve() called during checkout

**Aggregate:** Inventory | **Version:** 1.0

**Payload:**
- `inventoryId`: InventoryId
- `productId`: ProductId
- `orderId`: OrderId
- `quantity`: integer
- `remainingAvailable`: integer

**Invariants Reflected:** [Available quantity cannot be negative Reserved cannot exceed on-hand]

**Consumers:** [OrderService AnalyticsService]

---

### EVT-INV-002: InventoryReleased

**Purpose:** Records that reserved stock was released back

**Trigger:** Inventory.release() called on order cancel

**Aggregate:** Inventory | **Version:** 1.0

**Payload:**
- `inventoryId`: InventoryId
- `productId`: ProductId
- `orderId`: OrderId
- `quantity`: integer
- `reason`: string

**Invariants Reflected:** [Reserved quantity decreased appropriately]

**Consumers:** [OrderService AnalyticsService]

---

### EVT-INV-003: InventoryDeducted

**Purpose:** Records that stock was physically removed after shipment

**Trigger:** Inventory.deduct() called when order ships

**Aggregate:** Inventory | **Version:** 1.0

**Payload:**
- `inventoryId`: InventoryId
- `productId`: ProductId
- `orderId`: OrderId
- `quantity`: integer
- `remainingOnHand`: integer

**Invariants Reflected:** [Quantity on hand decreased Reserved decreased simultaneously]

**Consumers:** [AnalyticsService ReorderService]

---

### EVT-INV-004: InventoryRestocked

**Purpose:** Records that new stock was added to inventory

**Trigger:** Inventory.restock() called by admin

**Aggregate:** Inventory | **Version:** 1.0

**Payload:**
- `inventoryId`: InventoryId
- `productId`: ProductId
- `quantity`: integer
- `newQuantityOnHand`: integer
- `restockedBy`: UserId

**Invariants Reflected:** [Quantity must be positive]

**Consumers:** [AnalyticsService SearchService]

---

### EVT-INV-005: OutOfStockDetected

**Purpose:** Records that product available quantity reached zero

**Trigger:** Available quantity becomes zero

**Aggregate:** Inventory | **Version:** 1.0

**Payload:**
- `inventoryId`: InventoryId
- `productId`: ProductId
- `productName`: string
- `detectedAt`: DateTime

**Invariants Reflected:** [Available quantity is zero]

**Consumers:** [NotificationService ProductService AnalyticsService]

---

## Commands

### CMD-CUST-001: RegisterCustomer

**Intent:** Create new customer account with email and password

**Expected Outcome:** Customer created with REGISTERED status, cart initialized

**Required Data:**
- `email`: Email
- `password`: Password

**Failure Conditions:** [Email already registered Password does not meet requirements Invalid email format]

---

### CMD-CART-001: AddItemToCart

**Intent:** Add a product to customer's shopping cart

**Expected Outcome:** Item added to cart with updated total price

**Required Data:**
- `customerId`: CustomerId
- `productId`: ProductId
- `quantity`: integer
- `variantId`: VariantId

**Failure Conditions:** [Product out of stock Quantity exceeds available stock Product is inactive Variant required but not selected]

---

### CMD-CART-002: UpdateCartQuantity

**Intent:** Change quantity of item already in cart

**Expected Outcome:** Cart item quantity updated, total recalculated

**Required Data:**
- `cartId`: CartId
- `productId`: ProductId
- `newQuantity`: integer

**Failure Conditions:** [Item not in cart Quantity exceeds available stock Quantity less than 1]

---

### CMD-CART-003: RemoveItemFromCart

**Intent:** Remove an item from customer's cart

**Expected Outcome:** Item removed from cart, total recalculated

**Required Data:**
- `cartId`: CartId
- `productId`: ProductId

**Failure Conditions:** [Item not in cart]

---

### CMD-CART-004: ClearCart

**Intent:** Remove all items from cart

**Expected Outcome:** Cart emptied, total set to zero

**Required Data:**
- `cartId`: CartId

---

### CMD-ORDER-001: PlaceOrder

**Intent:** Convert cart into a new order with payment

**Expected Outcome:** Order created, payment authorized, inventory reserved, cart cleared

**Required Data:**
- `customerId`: CustomerId
- `shippingAddress`: ShippingAddress
- `paymentMethod`: PaymentMethod

**Failure Conditions:** [Customer not registered Cart is empty Insufficient stock Payment authorization failed International shipping not available]

---

### CMD-ORDER-002: ConfirmOrder

**Intent:** Confirm order after payment capture

**Expected Outcome:** Order status changed to CONFIRMED

**Required Data:**
- `orderId`: OrderId

**Failure Conditions:** [Order not in PENDING status Order not found]

---

### CMD-ORDER-003: ShipOrder

**Intent:** Mark order as shipped with tracking info

**Expected Outcome:** Order shipped, inventory deducted, customer notified

**Required Data:**
- `orderId`: OrderId
- `trackingNumber`: string
- `carrier`: string

**Failure Conditions:** [Order not in CONFIRMED status Invalid tracking number]

---

### CMD-ORDER-004: DeliverOrder

**Intent:** Mark order as delivered

**Expected Outcome:** Order status changed to DELIVERED

**Required Data:**
- `orderId`: OrderId
- `trackingNumber`: string

**Failure Conditions:** [Order not in SHIPPED status Order not found for tracking]

---

### CMD-ORDER-005: CancelOrder

**Intent:** Cancel order and restore inventory

**Expected Outcome:** Order cancelled, inventory released, refund initiated

**Required Data:**
- `orderId`: OrderId
- `reason`: string

**Failure Conditions:** [Order already shipped Order already delivered Order already cancelled]

---

### CMD-PROD-001: CreateProduct

**Intent:** Add new product to catalog

**Expected Outcome:** Product created, inventory initialized

**Required Data:**
- `name`: string
- `description`: string
- `price`: Money
- `categoryId`: CategoryId
- `initialStock`: integer

**Failure Conditions:** [Price below minimum Name invalid length Category not found]

---

### CMD-PROD-002: UpdateProduct

**Intent:** Modify product details

**Expected Outcome:** Product details updated

**Required Data:**
- `productId`: ProductId
- `name`: string
- `description`: string
- `price`: Money
- `categoryId`: CategoryId

**Failure Conditions:** [Product not found Price below minimum Invalid name length]

---

### CMD-PROD-003: DeactivateProduct

**Intent:** Remove product from active catalog

**Expected Outcome:** Product marked inactive, hidden from catalog

**Required Data:**
- `productId`: ProductId

**Failure Conditions:** [Product not found Product already inactive]

---

### CMD-INV-001: RestockInventory

**Intent:** Add stock to product inventory

**Expected Outcome:** Inventory quantity increased, audit logged

**Required Data:**
- `productId`: ProductId
- `quantity`: integer
- `restockedBy`: UserId

**Failure Conditions:** [Product not found Invalid quantity]

---

### CMD-INV-002: ReserveInventory

**Intent:** Reserve stock for pending order

**Expected Outcome:** Stock reserved, available quantity decreased

**Required Data:**
- `productId`: ProductId
- `quantity`: integer
- `orderId`: OrderId

**Failure Conditions:** [Insufficient available stock]

---

### CMD-INV-003: ReleaseInventory

**Intent:** Release reserved stock back to available

**Expected Outcome:** Reserved quantity decreased, available increased

**Required Data:**
- `productId`: ProductId
- `quantity`: integer
- `orderId`: OrderId
- `reason`: string

**Failure Conditions:** [Insufficient reserved quantity]

---

### CMD-CAT-001: CreateCategory

**Intent:** Create new product category

**Expected Outcome:** Category created in catalog hierarchy

**Required Data:**
- `name`: string
- `description`: string
- `parentCategoryId`: CategoryId

**Failure Conditions:** [Name already exists Parent category not found Self-reference]

---

## Integration Events

### INT-ORD-001: OrderPlacedNotification

**Purpose:** Notify external systems that a new order was placed

**Source:** OrderService

**Consumers:** [InventoryService PaymentService NotificationService AnalyticsService]

**Payload:** [orderId customerId totalAmount itemCount]

---

### INT-ORD-002: OrderShippedNotification

**Purpose:** Notify customer and external systems of shipment

**Source:** OrderService

**Consumers:** [NotificationService CarrierIntegration AnalyticsService]

**Payload:** [orderId customerId trackingNumber carrier]

---

### INT-ORD-003: OrderCancelledNotification

**Purpose:** Notify systems to reverse order-related operations

**Source:** OrderService

**Consumers:** [InventoryService PaymentService NotificationService]

**Payload:** [orderId customerId reason]

---

### INT-INV-001: StockLevelChanged

**Purpose:** Notify systems of inventory changes for sync

**Source:** InventoryService

**Consumers:** [SearchService AnalyticsService ReorderService]

**Payload:** [productId availableQuantity changeType]

---

### INT-INV-002: OutOfStockAlert

**Purpose:** Alert admin and systems when product runs out

**Source:** InventoryService

**Consumers:** [NotificationService ProductService PurchasingService]

**Payload:** [productId productName]

---

### INT-PAY-001: PaymentCaptured

**Purpose:** Notify order service that payment was captured

**Source:** PaymentService

**Consumers:** [OrderService AnalyticsService]

**Payload:** [orderId paymentId amount]

---

### INT-PAY-002: RefundInitiated

**Purpose:** Notify systems that refund process started

**Source:** PaymentService

**Consumers:** [OrderService NotificationService AnalyticsService]

**Payload:** [orderId refundId amount]

---

### INT-CUST-001: CustomerRegisteredNotification

**Purpose:** Notify external systems of new customer

**Source:** CustomerService

**Consumers:** [NotificationService MarketingService AnalyticsService]

**Payload:** [customerId email]

---

