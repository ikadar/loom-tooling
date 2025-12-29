# Event & Message Design

Generated: 2025-12-29T16:22:00+01:00

---

## Domain Events

### EVT-CUST-001: CustomerRegistered

**Purpose:** Signals new customer account was created

**Trigger:** Customer completes registration with valid email and password

**Aggregate:** Customer | **Version:** 1.0

**Payload:**
- `customerId`: UUID
- `email`: Email
- `firstName`: string
- `lastName`: string
- `createdAt`: DateTime

**Invariants Reflected:** [Email is unique Password meets requirements]

**Consumers:** [Email Service]

---

### EVT-CUST-002: CustomerEmailVerified

**Purpose:** Signals customer email address was verified

**Trigger:** Customer clicks valid verification link

**Aggregate:** Customer | **Version:** 1.0

**Payload:**
- `customerId`: UUID
- `email`: Email
- `verifiedAt`: DateTime

**Invariants Reflected:** [Customer can now place orders]

**Consumers:** [Customer Service]

---

### EVT-CUST-003: ShippingAddressAdded

**Purpose:** Signals new shipping address was added to customer profile

**Trigger:** Customer saves new shipping address

**Aggregate:** Customer | **Version:** 1.0

**Payload:**
- `customerId`: UUID
- `addressId`: UUID
- `street`: string
- `city`: string
- `state`: string
- `postalCode`: string
- `country`: string

**Invariants Reflected:** [Address is valid Country is domestic]

**Consumers:** [Customer Service]

---

### EVT-CUST-004: CustomerDataErased

**Purpose:** Signals customer data was deleted per GDPR request

**Trigger:** Customer confirms data erasure request with no pending orders

**Aggregate:** Customer | **Version:** 1.0

**Payload:**
- `anonymizedId`: UUID
- `erasedAt`: DateTime

**Invariants Reflected:** [No pending orders existed Data anonymized]

**Consumers:** [Audit Service]

---

### EVT-CART-001: ItemAddedToCart

**Purpose:** Signals product was added to shopping cart

**Trigger:** Customer adds product to cart with valid quantity

**Aggregate:** Cart | **Version:** 1.0

**Payload:**
- `cartId`: UUID
- `customerId`: UUID
- `productId`: UUID
- `variantId`: UUID
- `quantity`: integer
- `unitPrice`: Money
- `addedAt`: DateTime

**Invariants Reflected:** [Product is in stock Quantity capped to available]

**Consumers:** [Cart Service]

---

### EVT-CART-002: CartItemQuantityUpdated

**Purpose:** Signals cart item quantity was changed

**Trigger:** Customer changes quantity in cart

**Aggregate:** Cart | **Version:** 1.0

**Payload:**
- `cartId`: UUID
- `cartItemId`: UUID
- `productId`: UUID
- `oldQuantity`: integer
- `newQuantity`: integer
- `updatedAt`: DateTime

**Invariants Reflected:** [Quantity >= 1 Quantity <= stock]

**Consumers:** [Cart Service]

---

### EVT-CART-003: ItemRemovedFromCart

**Purpose:** Signals product was removed from cart

**Trigger:** Customer removes item or sets quantity to zero

**Aggregate:** Cart | **Version:** 1.0

**Payload:**
- `cartId`: UUID
- `productId`: UUID
- `variantId`: UUID
- `removedAt`: DateTime

**Invariants Reflected:** [Item existed in cart]

**Consumers:** [Cart Service]

---

### EVT-CART-004: CartCleared

**Purpose:** Signals all items removed from cart

**Trigger:** Order placed or customer clears cart

**Aggregate:** Cart | **Version:** 1.0

**Payload:**
- `cartId`: UUID
- `customerId`: UUID
- `reason`: string
- `clearedAt`: DateTime

**Invariants Reflected:** [Cart total is zero]

**Consumers:** [Cart Service]

---

### EVT-CART-005: CartMerged

**Purpose:** Signals guest cart was merged with customer cart on login

**Trigger:** Guest user logs in with items in cart

**Aggregate:** Cart | **Version:** 1.0

**Payload:**
- `cartId`: UUID
- `guestCartId`: UUID
- `customerId`: UUID
- `mergedItems`: array
- `mergedAt`: DateTime

**Invariants Reflected:** [Quantities capped to stock]

**Consumers:** [Cart Service]

---

### EVT-CART-006: CartExpired

**Purpose:** Signals cart was cleared due to inactivity

**Trigger:** Background job detects cart exceeds inactivity threshold

**Aggregate:** Cart | **Version:** 1.0

**Payload:**
- `cartId`: UUID
- `customerId`: UUID
- `lastActivityDate`: DateTime
- `expiredAt`: DateTime

**Invariants Reflected:** [Inactivity exceeded threshold]

**Consumers:** [Notification Service]

---

### EVT-ORDER-001: OrderPlaced

**Purpose:** Signals new order was created from cart

**Trigger:** Customer places order with valid payment

**Aggregate:** Order | **Version:** 1.0

**Payload:**
- `orderId`: UUID
- `orderNumber`: string
- `customerId`: UUID
- `lineItems`: array
- `subtotal`: Money
- `shippingCost`: Money
- `tax`: Money
- `totalAmount`: Money
- `shippingAddress`: Address
- `createdAt`: DateTime

**Invariants Reflected:** [Has items Payment authorized Stock reserved]

**Consumers:** [Inventory Service Email Service Payment Service]

---

### EVT-ORDER-002: OrderConfirmed

**Purpose:** Signals order was confirmed by admin

**Trigger:** Admin confirms pending order

**Aggregate:** Order | **Version:** 1.0

**Payload:**
- `orderId`: UUID
- `orderNumber`: string
- `confirmedAt`: DateTime
- `confirmedBy`: UUID

**Invariants Reflected:** [Status was pending]

**Consumers:** [Order Service]

---

### EVT-ORDER-003: OrderShipped

**Purpose:** Signals order was shipped with tracking

**Trigger:** Admin marks order as shipped with tracking number

**Aggregate:** Order | **Version:** 1.0

**Payload:**
- `orderId`: UUID
- `orderNumber`: string
- `trackingNumber`: string
- `carrier`: string
- `shippedAt`: DateTime

**Invariants Reflected:** [Status was confirmed Tracking provided]

**Consumers:** [Email Service Inventory Service]

---

### EVT-ORDER-004: OrderDelivered

**Purpose:** Signals order was delivered to customer

**Trigger:** Admin marks order as delivered

**Aggregate:** Order | **Version:** 1.0

**Payload:**
- `orderId`: UUID
- `orderNumber`: string
- `deliveredAt`: DateTime

**Invariants Reflected:** [Status was shipped]

**Consumers:** [Email Service]

---

### EVT-ORDER-005: OrderCancelled

**Purpose:** Signals order was cancelled before shipping

**Trigger:** Customer or admin cancels pending/confirmed order

**Aggregate:** Order | **Version:** 1.0

**Payload:**
- `orderId`: UUID
- `orderNumber`: string
- `customerId`: UUID
- `reason`: string
- `lineItems`: array
- `cancelledAt`: DateTime

**Invariants Reflected:** [Status was pending or confirmed]

**Consumers:** [Inventory Service Payment Service Email Service]

---

### EVT-PROD-001: ProductCreated

**Purpose:** Signals new product was added to catalog

**Trigger:** Admin creates product with valid data

**Aggregate:** Product | **Version:** 1.0

**Payload:**
- `productId`: UUID
- `name`: string
- `price`: Money
- `categoryId`: UUID
- `status`: string
- `createdBy`: UUID
- `createdAt`: DateTime

**Invariants Reflected:** [Name valid Price >= 0.01 Category exists]

**Consumers:** [Inventory Service Catalog Service]

---

### EVT-PROD-002: ProductUpdated

**Purpose:** Signals product attributes were changed

**Trigger:** Admin updates product details

**Aggregate:** Product | **Version:** 1.0

**Payload:**
- `productId`: UUID
- `changedFields`: object
- `updatedBy`: UUID
- `updatedAt`: DateTime

**Invariants Reflected:** [Fields valid after change]

**Consumers:** [Cart Service Catalog Service]

---

### EVT-PROD-003: ProductDeleted

**Purpose:** Signals product was soft-deleted from catalog

**Trigger:** Admin deletes product

**Aggregate:** Product | **Version:** 1.0

**Payload:**
- `productId`: UUID
- `deletedBy`: UUID
- `deletedAt`: DateTime

**Invariants Reflected:** [Order history preserved via soft delete]

**Consumers:** [Cart Service Inventory Service]

---

### EVT-PROD-004: ProductActivated

**Purpose:** Signals product is now available for purchase

**Trigger:** Admin activates draft product

**Aggregate:** Product | **Version:** 1.0

**Payload:**
- `productId`: UUID
- `activatedAt`: DateTime

**Invariants Reflected:** [Product was inactive]

**Consumers:** [Catalog Service]

---

### EVT-PROD-005: ProductDeactivated

**Purpose:** Signals product is no longer available for purchase

**Trigger:** Admin deactivates product

**Aggregate:** Product | **Version:** 1.0

**Payload:**
- `productId`: UUID
- `deactivatedAt`: DateTime

**Invariants Reflected:** [Product was active]

**Consumers:** [Cart Service Catalog Service]

---

### EVT-INV-001: InventoryQuantitySet

**Purpose:** Signals inventory was set to absolute value

**Trigger:** Admin sets inventory quantity with reason

**Aggregate:** Inventory | **Version:** 1.0

**Payload:**
- `inventoryId`: UUID
- `productId`: UUID
- `previousValue`: integer
- `newValue`: integer
- `reason`: string
- `userId`: UUID
- `timestamp`: DateTime

**Invariants Reflected:** [Stock >= 0 Audit logged]

**Consumers:** [Inventory Service]

---

### EVT-INV-002: InventoryAdjusted

**Purpose:** Signals inventory was adjusted by delta

**Trigger:** Admin adjusts inventory with positive or negative delta

**Aggregate:** Inventory | **Version:** 1.0

**Payload:**
- `inventoryId`: UUID
- `productId`: UUID
- `adjustment`: integer
- `newValue`: integer
- `reason`: string
- `userId`: UUID
- `timestamp`: DateTime

**Invariants Reflected:** [Result >= 0 Audit logged]

**Consumers:** [Inventory Service]

---

### EVT-INV-003: InventoryReserved

**Purpose:** Signals stock was reserved for pending order

**Trigger:** Order placed successfully

**Aggregate:** Inventory | **Version:** 1.0

**Payload:**
- `inventoryId`: UUID
- `productId`: UUID
- `quantity`: integer
- `orderId`: UUID
- `reservedAt`: DateTime

**Invariants Reflected:** [Quantity <= available]

**Consumers:** [Order Service]

---

### EVT-INV-004: InventoryReleased

**Purpose:** Signals reserved stock was released

**Trigger:** Order cancelled or reservation expired

**Aggregate:** Inventory | **Version:** 1.0

**Payload:**
- `inventoryId`: UUID
- `productId`: UUID
- `quantity`: integer
- `reason`: string
- `releasedAt`: DateTime

**Invariants Reflected:** [Quantity <= reserved]

**Consumers:** [Order Service]

---

### EVT-INV-005: InventoryDeducted

**Purpose:** Signals stock was permanently deducted

**Trigger:** Order shipped successfully

**Aggregate:** Inventory | **Version:** 1.0

**Payload:**
- `inventoryId`: UUID
- `productId`: UUID
- `quantity`: integer
- `orderId`: UUID
- `deductedAt`: DateTime

**Invariants Reflected:** [Quantity <= reserved Audit logged]

**Consumers:** [Inventory Service]

---

### EVT-INV-006: LowStockAlert

**Purpose:** Signals product stock fell below threshold

**Trigger:** Stock level drops below configured threshold

**Aggregate:** Inventory | **Version:** 1.0

**Payload:**
- `inventoryId`: UUID
- `productId`: UUID
- `currentStock`: integer
- `threshold`: integer
- `triggeredAt`: DateTime

**Invariants Reflected:** [Stock < threshold]

**Consumers:** [Notification Service Admin Service]

---

### EVT-INV-007: OutOfStock

**Purpose:** Signals product has zero available stock

**Trigger:** Available quantity reaches zero

**Aggregate:** Inventory | **Version:** 1.0

**Payload:**
- `inventoryId`: UUID
- `productId`: UUID
- `triggeredAt`: DateTime

**Invariants Reflected:** [Available quantity == 0]

**Consumers:** [Cart Service Catalog Service Notification Service]

---

### EVT-CAT-001: CategoryCreated

**Purpose:** Signals new category was added to hierarchy

**Trigger:** Admin creates category with valid name and parent

**Aggregate:** Category | **Version:** 1.0

**Payload:**
- `categoryId`: UUID
- `name`: string
- `parentCategoryId`: UUID
- `depth`: integer
- `createdAt`: DateTime

**Invariants Reflected:** [Name unique Depth <= 3]

**Consumers:** [Catalog Service]

---

### EVT-PAY-001: PaymentAuthorized

**Purpose:** Signals payment was authorized by gateway

**Trigger:** Payment gateway approves authorization request

**Aggregate:** Payment | **Version:** 1.0

**Payload:**
- `paymentId`: UUID
- `orderId`: UUID
- `amount`: Money
- `authorizationId`: string
- `paymentMethod`: string
- `authorizedAt`: DateTime

**Invariants Reflected:** [Amount matches order total]

**Consumers:** [Order Service]

---

### EVT-PAY-002: PaymentRefunded

**Purpose:** Signals refund was processed for cancelled order

**Trigger:** Order cancelled with successful refund

**Aggregate:** Payment | **Version:** 1.0

**Payload:**
- `paymentId`: UUID
- `orderId`: UUID
- `refundAmount`: Money
- `refundId`: string
- `refundedAt`: DateTime

**Invariants Reflected:** [Original payment existed]

**Consumers:** [Order Service Email Service]

---

## Commands

### CMD-CUST-001: RegisterCustomer

**Intent:** Create new customer account with email and password

**Expected Outcome:** Customer created with unverified status, verification email queued

**Required Data:**
- `email`: Email
- `password`: string
- `firstName`: string
- `lastName`: string

**Failure Conditions:** [Email already registered Invalid password format Invalid email]

---

### CMD-CUST-002: VerifyCustomerEmail

**Intent:** Verify customer email address via token

**Expected Outcome:** Customer email verified, can now place orders

**Required Data:**
- `verificationToken`: string

**Failure Conditions:** [Invalid token Token expired]

---

### CMD-CUST-003: AddShippingAddress

**Intent:** Add new shipping address to customer profile

**Expected Outcome:** Address saved and available for orders

**Required Data:**
- `customerId`: UUID
- `street`: string
- `city`: string
- `state`: string
- `postalCode`: string
- `country`: string
- `recipientName`: string

**Failure Conditions:** [Invalid postal code Country not supported Missing field]

---

### CMD-CUST-004: RequestGDPRErasure

**Intent:** Delete all customer personal data per GDPR

**Expected Outcome:** Personal data deleted, order history anonymized

**Required Data:**
- `customerId`: UUID
- `confirmation`: boolean

**Failure Conditions:** [Pending orders exist]

---

### CMD-CART-001: AddItemToCart

**Intent:** Add product to shopping cart

**Expected Outcome:** Item added to cart, totals recalculated

**Required Data:**
- `customerId`: UUID
- `productId`: UUID
- `variantId`: UUID
- `quantity`: integer

**Failure Conditions:** [Out of stock Variant required but missing Product inactive]

---

### CMD-CART-002: UpdateCartItemQuantity

**Intent:** Change quantity of item in cart

**Expected Outcome:** Quantity updated or item removed if zero

**Required Data:**
- `cartId`: UUID
- `cartItemId`: UUID
- `quantity`: integer

**Failure Conditions:** [Item not in cart Quantity exceeds stock]

---

### CMD-CART-003: RemoveItemFromCart

**Intent:** Remove product from shopping cart

**Expected Outcome:** Item removed, totals recalculated, undo available

**Required Data:**
- `cartId`: UUID
- `cartItemId`: UUID

**Failure Conditions:** [Item not in cart]

---

### CMD-CART-004: MergeGuestCart

**Intent:** Merge guest cart into customer cart on login

**Expected Outcome:** Carts merged, quantities capped to stock

**Required Data:**
- `guestCartId`: UUID
- `customerId`: UUID

**Failure Conditions:** [Guest cart not found]

---

### CMD-ORDER-001: PlaceOrder

**Intent:** Create order from cart contents with payment

**Expected Outcome:** Order created, payment authorized, inventory reserved, cart cleared

**Required Data:**
- `customerId`: UUID
- `shippingAddressId`: UUID
- `paymentMethodId`: UUID

**Failure Conditions:** [Cart empty Out of stock Payment failed Email not verified]

---

### CMD-ORDER-002: ConfirmOrder

**Intent:** Confirm pending order for fulfillment

**Expected Outcome:** Order status changed to confirmed

**Required Data:**
- `orderId`: UUID

**Failure Conditions:** [Order not found Invalid status transition]

---

### CMD-ORDER-003: ShipOrder

**Intent:** Mark order as shipped with tracking number

**Expected Outcome:** Status shipped, inventory deducted, customer notified

**Required Data:**
- `orderId`: UUID
- `trackingNumber`: string
- `carrier`: string

**Failure Conditions:** [Order not found Invalid status transition Missing tracking]

---

### CMD-ORDER-004: CancelOrder

**Intent:** Cancel order before shipping with refund

**Expected Outcome:** Order cancelled, inventory released, refund initiated

**Required Data:**
- `orderId`: UUID
- `reason`: string

**Failure Conditions:** [Order not found Already shipped]

---

### CMD-PROD-001: CreateProduct

**Intent:** Add new product to catalog

**Expected Outcome:** Product created as draft with UUID

**Required Data:**
- `name`: string
- `description`: string
- `price`: Money
- `categoryId`: UUID

**Failure Conditions:** [Invalid name Invalid price Category not found]

---

### CMD-PROD-002: UpdateProduct

**Intent:** Update product attributes

**Expected Outcome:** Product updated, change logged

**Required Data:**
- `productId`: UUID
- `name`: string
- `price`: Money
- `description`: string

**Failure Conditions:** [Product not found Invalid data]

---

### CMD-PROD-003: DeleteProduct

**Intent:** Remove product from catalog via soft delete

**Expected Outcome:** Product soft-deleted, removed from carts

**Required Data:**
- `productId`: UUID

**Failure Conditions:** [Product not found]

---

### CMD-INV-001: SetInventoryQuantity

**Intent:** Set absolute inventory quantity with audit

**Expected Outcome:** Stock level set, audit log created

**Required Data:**
- `productId`: UUID
- `quantity`: integer
- `reason`: string

**Failure Conditions:** [Product not found Negative quantity]

---

### CMD-INV-002: AdjustInventory

**Intent:** Adjust inventory by positive or negative delta

**Expected Outcome:** Stock adjusted, audit log created

**Required Data:**
- `productId`: UUID
- `adjustment`: integer
- `reason`: string

**Failure Conditions:** [Product not found Would result in negative stock]

---

### CMD-INV-003: ReserveInventory

**Intent:** Reserve stock for pending order

**Expected Outcome:** Inventory reserved for order

**Required Data:**
- `productId`: UUID
- `quantity`: integer
- `orderId`: UUID

**Failure Conditions:** [Insufficient available stock]

---

### CMD-INV-004: ReleaseInventory

**Intent:** Release reserved inventory

**Expected Outcome:** Reserved quantity released

**Required Data:**
- `productId`: UUID
- `quantity`: integer
- `reason`: string

**Failure Conditions:** [Quantity exceeds reserved]

---

### CMD-CAT-001: CreateCategory

**Intent:** Create new product category in hierarchy

**Expected Outcome:** Category created with calculated depth

**Required Data:**
- `name`: string
- `parentCategoryId`: UUID

**Failure Conditions:** [Name not unique Parent not found Exceeds max nesting]

---

### CMD-PAY-001: AuthorizePayment

**Intent:** Authorize payment with gateway

**Expected Outcome:** Payment authorized, authorization ID stored

**Required Data:**
- `paymentMethodId`: UUID
- `amount`: Money
- `orderId`: UUID

**Failure Conditions:** [Card declined Unsupported card type Gateway error]

---

### CMD-PAY-002: RefundPayment

**Intent:** Refund payment for cancelled order

**Expected Outcome:** Refund processed to original payment method

**Required Data:**
- `orderId`: UUID
- `amount`: Money

**Failure Conditions:** [Original payment not found Refund failed]

---

### CMD-PAY-003: InitiatePayPal

**Intent:** Start PayPal authorization flow

**Expected Outcome:** PayPal redirect URL returned

**Required Data:**
- `amount`: Money
- `returnUrl`: string
- `cancelUrl`: string

**Failure Conditions:** [PayPal service error]

---

## Integration Events

### INT-ORDER-001: OrderPlacedNotification

**Purpose:** Notify external systems of new order for processing

**Source:** Order Service

**Consumers:** [Email Service Analytics Service Fulfillment System]

**Payload:** [orderId orderNumber customerId totalAmount itemCount]

---

### INT-ORDER-002: OrderShippedNotification

**Purpose:** Notify customer and systems that order shipped

**Source:** Order Service

**Consumers:** [Email Service Customer App Analytics Service]

**Payload:** [orderId orderNumber customerId trackingNumber carrier]

---

### INT-ORDER-003: OrderCancelledNotification

**Purpose:** Notify systems of order cancellation

**Source:** Order Service

**Consumers:** [Email Service Analytics Service Payment Service]

**Payload:** [orderId orderNumber customerId reason refundAmount]

---

### INT-INV-001: LowStockNotification

**Purpose:** Alert admins when product stock is low

**Source:** Inventory Service

**Consumers:** [Admin Dashboard Email Service Slack Integration]

**Payload:** [productId productName currentStock threshold]

---

### INT-INV-002: OutOfStockNotification

**Purpose:** Alert systems when product becomes unavailable

**Source:** Inventory Service

**Consumers:** [Catalog Service Cart Service Admin Dashboard]

**Payload:** [productId productName]

---

### INT-CUST-001: CustomerRegisteredNotification

**Purpose:** Notify systems of new customer registration

**Source:** Customer Service

**Consumers:** [Email Service Analytics Service CRM System]

**Payload:** [customerId email firstName]

---

### INT-PAY-001: PaymentFailedNotification

**Purpose:** Alert on payment failure for monitoring

**Source:** Payment Service

**Consumers:** [Admin Dashboard Monitoring Service]

**Payload:** [orderId customerId amount failureReason]

---

### INT-PAY-002: RefundProcessedNotification

**Purpose:** Notify customer and systems of completed refund

**Source:** Payment Service

**Consumers:** [Email Service Analytics Service]

**Payload:** [orderId customerId refundAmount refundId]

---

