# Event & Message Design

Generated: 2026-01-02T20:15:12+01:00

---

## Domain Events

### EVT-CUS-001: CustomerRegistered

**Purpose:** New customer account created with pending email verification

**Trigger:** Customer submits valid registration form

**Aggregate:** Customer | **Version:** 1.0

**Payload:**
- `customerId`: UUID
- `email`: string
- `registeredAt`: timestamp

**Invariants Reflected:** [Email must be unique Password meets strength requirements]

**Consumers:** [EmailService CartService AnalyticsService]

---

### EVT-CUS-002: CustomerEmailVerified

**Purpose:** Customer email address confirmed as valid

**Trigger:** Customer clicks verification link

**Aggregate:** Customer | **Version:** 1.0

**Payload:**
- `customerId`: UUID
- `verifiedAt`: timestamp

**Invariants Reflected:** [Email verified before ordering]

**Consumers:** [NotificationService]

---

### EVT-CART-001: ItemAddedToCart

**Purpose:** Product added to customer shopping cart

**Trigger:** Customer adds item with available stock

**Aggregate:** Cart | **Version:** 1.0

**Payload:**
- `cartId`: UUID
- `customerId`: UUID
- `productId`: UUID
- `variantId`: UUID?
- `quantity`: int
- `unitPrice`: decimal
- `addedAt`: timestamp

**Invariants Reflected:** [Product must be active Stock must be available]

**Consumers:** [AnalyticsService]

---

### EVT-CART-002: CartQuantityUpdated

**Purpose:** Quantity of cart item changed

**Trigger:** Customer updates quantity in cart view

**Aggregate:** Cart | **Version:** 1.0

**Payload:**
- `cartId`: UUID
- `productId`: UUID
- `variantId`: UUID?
- `previousQuantity`: int
- `newQuantity`: int
- `updatedAt`: timestamp

**Invariants Reflected:** [Quantity must not exceed available stock]

**Consumers:** [AnalyticsService]

---

### EVT-CART-003: ItemRemovedFromCart

**Purpose:** Product removed from shopping cart

**Trigger:** Customer removes item from cart

**Aggregate:** Cart | **Version:** 1.0

**Payload:**
- `cartId`: UUID
- `productId`: UUID
- `variantId`: UUID?
- `removedAt`: timestamp

**Consumers:** [AnalyticsService]

---

### EVT-CART-004: CartCleared

**Purpose:** All items removed from cart after order placement

**Trigger:** Order successfully placed

**Aggregate:** Cart | **Version:** 1.0

**Payload:**
- `cartId`: UUID
- `customerId`: UUID
- `clearedAt`: timestamp
- `orderId`: UUID

---

### EVT-ORD-001: OrderPlaced

**Purpose:** New order created with reserved inventory and authorized payment

**Trigger:** Customer completes checkout

**Aggregate:** Order | **Version:** 1.0

**Payload:**
- `orderId`: UUID
- `customerId`: UUID
- `lineItems`: LineItem[]
- `shippingAddress`: Address
- `subtotal`: decimal
- `shippingCost`: decimal
- `total`: decimal
- `paymentAuthorizationId`: string
- `placedAt`: timestamp

**Invariants Reflected:** [Customer verified Stock available Payment authorized]

**Consumers:** [InventoryService NotificationService AnalyticsService]

---

### EVT-ORD-002: OrderConfirmed

**Purpose:** Order payment captured and ready for fulfillment

**Trigger:** Payment captured successfully

**Aggregate:** Order | **Version:** 1.0

**Payload:**
- `orderId`: UUID
- `confirmedAt`: timestamp

**Invariants Reflected:** [Payment must be captured]

**Consumers:** [NotificationService FulfillmentService]

---

### EVT-ORD-003: OrderShipped

**Purpose:** Order dispatched with tracking information

**Trigger:** Warehouse marks order as shipped

**Aggregate:** Order | **Version:** 1.0

**Payload:**
- `orderId`: UUID
- `trackingNumber`: string
- `carrier`: string
- `shippedAt`: timestamp

**Invariants Reflected:** [Order must be confirmed]

**Consumers:** [NotificationService AnalyticsService]

---

### EVT-ORD-004: OrderDelivered

**Purpose:** Order received by customer

**Trigger:** Carrier confirms delivery

**Aggregate:** Order | **Version:** 1.0

**Payload:**
- `orderId`: UUID
- `deliveredAt`: timestamp

**Invariants Reflected:** [Order must be shipped]

**Consumers:** [NotificationService AnalyticsService]

---

### EVT-ORD-005: OrderCancelled

**Purpose:** Order cancelled before shipping

**Trigger:** Customer requests cancellation

**Aggregate:** Order | **Version:** 1.0

**Payload:**
- `orderId`: UUID
- `reason`: string
- `cancelledAt`: timestamp
- `refundId`: string

**Invariants Reflected:** [Order must not be shipped]

**Consumers:** [InventoryService NotificationService PaymentService]

---

### EVT-PROD-001: ProductCreated

**Purpose:** New product added to catalog

**Trigger:** Admin creates product

**Aggregate:** Product | **Version:** 1.0

**Payload:**
- `productId`: UUID
- `name`: string
- `description`: string
- `price`: decimal
- `categoryId`: UUID
- `createdAt`: timestamp

**Invariants Reflected:** [Price >= $0.01 Category exists]

**Consumers:** [InventoryService SearchService]

---

### EVT-PROD-002: ProductUpdated

**Purpose:** Product details or price changed

**Trigger:** Admin updates product

**Aggregate:** Product | **Version:** 1.0

**Payload:**
- `productId`: UUID
- `changedFields`: string[]
- `previousPrice`: decimal?
- `newPrice`: decimal?
- `updatedAt`: timestamp

**Invariants Reflected:** [Price >= $0.01]

**Consumers:** [SearchService AnalyticsService]

---

### EVT-PROD-003: ProductDeactivated

**Purpose:** Product removed from active catalog

**Trigger:** Admin deactivates product

**Aggregate:** Product | **Version:** 1.0

**Payload:**
- `productId`: UUID
- `deactivatedAt`: timestamp

**Consumers:** [SearchService CartService]

---

### EVT-PROD-004: ProductVariantAdded

**Purpose:** New size/color variant added to product

**Trigger:** Admin adds variant

**Aggregate:** Product | **Version:** 1.0

**Payload:**
- `productId`: UUID
- `variantId`: UUID
- `sku`: string
- `size`: string?
- `color`: string?
- `priceAdjustment`: decimal
- `addedAt`: timestamp

**Invariants Reflected:** [SKU must be unique At least size or color required]

**Consumers:** [InventoryService SearchService]

---

### EVT-INV-001: InventoryReserved

**Purpose:** Stock reserved for pending order

**Trigger:** Order placed successfully

**Aggregate:** Inventory | **Version:** 1.0

**Payload:**
- `inventoryId`: UUID
- `productId`: UUID
- `variantId`: UUID?
- `orderId`: UUID
- `quantity`: int
- `reservedAt`: timestamp

**Invariants Reflected:** [Available quantity >= reserved quantity]

**Consumers:** [OrderService]

---

### EVT-INV-002: InventoryReleased

**Purpose:** Reserved stock returned to available

**Trigger:** Order cancelled

**Aggregate:** Inventory | **Version:** 1.0

**Payload:**
- `inventoryId`: UUID
- `productId`: UUID
- `variantId`: UUID?
- `orderId`: UUID
- `quantity`: int
- `reason`: string
- `releasedAt`: timestamp

**Consumers:** [OrderService]

---

### EVT-INV-003: InventoryDeducted

**Purpose:** Stock permanently reduced after shipping

**Trigger:** Order shipped

**Aggregate:** Inventory | **Version:** 1.0

**Payload:**
- `inventoryId`: UUID
- `productId`: UUID
- `variantId`: UUID?
- `orderId`: UUID
- `quantity`: int
- `deductedAt`: timestamp

**Consumers:** [AnalyticsService]

---

### EVT-INV-004: InventoryRestocked

**Purpose:** Stock added from supplier delivery

**Trigger:** Warehouse receives shipment

**Aggregate:** Inventory | **Version:** 1.0

**Payload:**
- `inventoryId`: UUID
- `productId`: UUID
- `variantId`: UUID?
- `quantity`: int
- `previousQuantity`: int
- `newQuantity`: int
- `restockedAt`: timestamp

**Invariants Reflected:** [Quantity must be positive]

**Consumers:** [AnalyticsService]

---

### EVT-INV-005: OutOfStockDetected

**Purpose:** Product available quantity reached zero

**Trigger:** Inventory deducted to zero

**Aggregate:** Inventory | **Version:** 1.0

**Payload:**
- `inventoryId`: UUID
- `productId`: UUID
- `variantId`: UUID?
- `detectedAt`: timestamp

**Consumers:** [NotificationService OperationsService]

---

### EVT-CAT-001: CategoryCreated

**Purpose:** New product category added to hierarchy

**Trigger:** Admin creates category

**Aggregate:** Category | **Version:** 1.0

**Payload:**
- `categoryId`: UUID
- `name`: string
- `parentId`: UUID?
- `createdAt`: timestamp

**Invariants Reflected:** [Name unique Hierarchy depth <= 3]

**Consumers:** [SearchService]

---

## Commands

### CMD-CUS-001: RegisterCustomer

**Intent:** Create new customer account

**Expected Outcome:** Customer created with REGISTERED status, verification email sent

**Required Data:**
- `email`: string
- `password`: string
- `name`: string

**Failure Conditions:** [Email already registered Password too weak]

---

### CMD-CUS-002: VerifyEmail

**Intent:** Confirm customer email address

**Expected Outcome:** Customer email verified, full access granted

**Required Data:**
- `customerId`: UUID
- `verificationToken`: string

**Failure Conditions:** [Invalid token Token expired]

---

### CMD-CART-001: AddItemToCart

**Intent:** Add product to shopping cart

**Expected Outcome:** Item added to cart, total recalculated

**Required Data:**
- `customerId`: UUID
- `productId`: UUID
- `variantId`: UUID?
- `quantity`: int

**Failure Conditions:** [Product out of stock Quantity exceeds stock Product inactive Variant required]

---

### CMD-CART-002: UpdateCartQuantity

**Intent:** Change quantity of cart item

**Expected Outcome:** Quantity updated, cart total recalculated

**Required Data:**
- `cartId`: UUID
- `productId`: UUID
- `variantId`: UUID?
- `newQuantity`: int

**Failure Conditions:** [Quantity exceeds stock Item not in cart]

---

### CMD-CART-003: RemoveItemFromCart

**Intent:** Remove product from cart

**Expected Outcome:** Item removed, cart total recalculated

**Required Data:**
- `cartId`: UUID
- `productId`: UUID
- `variantId`: UUID?

**Failure Conditions:** [Item not in cart]

---

### CMD-ORD-001: PlaceOrder

**Intent:** Complete checkout and create order

**Expected Outcome:** Order created in PENDING status, inventory reserved, payment authorized

**Required Data:**
- `customerId`: UUID
- `cartId`: UUID
- `shippingAddress`: Address
- `paymentMethod`: PaymentMethod

**Failure Conditions:** [Not registered Email not verified Cart empty Invalid address Stock insufficient Payment failed]

---

### CMD-ORD-002: ConfirmOrder

**Intent:** Capture payment and confirm order for fulfillment

**Expected Outcome:** Order status changed to CONFIRMED

**Required Data:**
- `orderId`: UUID

**Failure Conditions:** [Payment capture failed Order not in PENDING status]

---

### CMD-ORD-003: ShipOrder

**Intent:** Mark order as shipped with tracking

**Expected Outcome:** Order shipped, inventory deducted, customer notified

**Required Data:**
- `orderId`: UUID
- `trackingNumber`: string
- `carrier`: string

**Failure Conditions:** [Order not CONFIRMED Inventory deduction failed]

---

### CMD-ORD-004: DeliverOrder

**Intent:** Mark order as delivered to customer

**Expected Outcome:** Order status changed to DELIVERED

**Required Data:**
- `orderId`: UUID
- `trackingNumber`: string

**Failure Conditions:** [Order not found Order not SHIPPED]

---

### CMD-ORD-005: CancelOrder

**Intent:** Cancel order and process refund

**Expected Outcome:** Order cancelled, inventory released, refund initiated

**Required Data:**
- `orderId`: UUID
- `reason`: string

**Failure Conditions:** [Order already shipped Order already delivered]

---

### CMD-PROD-001: CreateProduct

**Intent:** Add new product to catalog

**Expected Outcome:** Product created and active, inventory initialized

**Required Data:**
- `name`: string
- `description`: string
- `price`: decimal
- `categoryId`: UUID
- `initialStock`: int

**Failure Conditions:** [Price < $0.01 Invalid name Description too long Category not found]

---

### CMD-PROD-002: UpdateProductPrice

**Intent:** Change product price

**Expected Outcome:** Price updated, audit trail recorded

**Required Data:**
- `productId`: UUID
- `newPrice`: decimal

**Failure Conditions:** [Price < $0.01 Product not found]

---

### CMD-PROD-003: DeactivateProduct

**Intent:** Remove product from active catalog

**Expected Outcome:** Product deactivated, no longer visible

**Required Data:**
- `productId`: UUID

---

### CMD-PROD-004: AddProductVariant

**Intent:** Add size/color variant to product

**Expected Outcome:** Variant added with inventory

**Required Data:**
- `productId`: UUID
- `sku`: string
- `size`: string?
- `color`: string?
- `priceAdjustment`: decimal
- `initialStock`: int

**Failure Conditions:** [SKU exists Product not found No attribute specified]

---

### CMD-INV-001: RestockInventory

**Intent:** Add stock from supplier

**Expected Outcome:** Inventory increased, audit trail recorded

**Required Data:**
- `productId`: UUID
- `variantId`: UUID?
- `quantity`: int

**Failure Conditions:** [Invalid quantity Product not found]

---

### CMD-CAT-001: CreateCategory

**Intent:** Add new product category

**Expected Outcome:** Category created in hierarchy

**Required Data:**
- `name`: string
- `parentId`: UUID?

**Failure Conditions:** [Name exists Depth exceeded Parent not found]

---

## Integration Events

### INT-ORD-001: OrderPlacedNotification

**Purpose:** Notify external systems of new order for fulfillment coordination

**Source:** OrderService

**Consumers:** [WarehouseSystem ShippingService AnalyticsPlatform]

**Payload:** [orderId customerId itemCount total shippingAddress]

---

### INT-ORD-002: OrderShippedNotification

**Purpose:** Notify customer and tracking systems of shipment

**Source:** OrderService

**Consumers:** [EmailService SMSService TrackingPortal]

**Payload:** [orderId customerId trackingNumber carrier estimatedDelivery]

---

### INT-ORD-003: OrderCancelledNotification

**Purpose:** Notify payment and inventory systems of cancellation

**Source:** OrderService

**Consumers:** [PaymentGateway WarehouseSystem AnalyticsPlatform]

**Payload:** [orderId customerId reason refundAmount]

---

### INT-INV-001: LowStockAlert

**Purpose:** Alert operations team when stock is critically low

**Source:** InventoryService

**Consumers:** [OperationsDashboard EmailService ProcurementSystem]

**Payload:** [productId productName currentStock reorderPoint]

---

### INT-INV-002: OutOfStockAlert

**Purpose:** Urgent notification when product becomes unavailable

**Source:** InventoryService

**Consumers:** [OperationsDashboard EmailService MarketingSystem]

**Payload:** [productId productName lastAvailableAt]

---

### INT-CUS-001: CustomerRegisteredNotification

**Purpose:** Trigger welcome email and CRM sync

**Source:** CustomerService

**Consumers:** [EmailService CRMSystem MarketingPlatform]

**Payload:** [customerId email name registeredAt]

---

### INT-PAY-001: PaymentProcessedNotification

**Purpose:** Notify accounting and reporting systems of payment

**Source:** PaymentService

**Consumers:** [AccountingSystem ReportingPlatform]

**Payload:** [paymentId orderId amount method processedAt]

---

### INT-PAY-002: RefundProcessedNotification

**Purpose:** Notify accounting of refund for reconciliation

**Source:** PaymentService

**Consumers:** [AccountingSystem CustomerService]

**Payload:** [refundId orderId amount processedAt]

---

