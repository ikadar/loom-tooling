# Service Boundaries

Generated: 2025-12-29T16:22:00+01:00

---

## SVC-CUSTOMER: Customer Service

**Purpose:** Manages customer registration, authentication, and addresses

**API Base:** `/api/v1/customers`

**Separation Reason:** Customer identity is core bounded context with auth concerns

### Capabilities
- Register new customers
- Verify customer email
- Manage shipping addresses
- Authenticate customers
- Process GDPR data erasure

### Inputs
- **command:** RegisterCustomer
- **command:** VerifyCustomerEmail
- **command:** AddShippingAddress
- **command:** UpdateShippingAddress
- **command:** RemoveShippingAddress
- **command:** RequestGDPRErasure

### Outputs
- **event:** CustomerRegistered
- **event:** CustomerEmailVerified
- **event:** ShippingAddressAdded
- **event:** ShippingAddressUpdated
- **event:** ShippingAddressRemoved
- **event:** CustomerDataErased

### Owned Aggregates
- Customer

### Dependencies
- **Email Service** (async): Send verification emails
- **Order Service** (sync): Check pending orders for GDPR

---

## SVC-CATALOG: Product Catalog Service

**Purpose:** Manages product information, variants, categories, and images

**API Base:** `/api/v1/products`

**Separation Reason:** Product catalog is distinct domain with admin-focused operations

### Capabilities
- Create and update products
- Manage product variants
- Upload product images
- Organize category hierarchy
- Activate/deactivate products
- Soft delete products

### Inputs
- **command:** CreateProduct
- **command:** UpdateProduct
- **command:** ActivateProduct
- **command:** DeactivateProduct
- **command:** DeleteProduct
- **command:** AddProductVariant
- **command:** UploadProductImages
- **command:** CreateCategory
- **command:** MoveCategory

### Outputs
- **event:** ProductCreated
- **event:** ProductUpdated
- **event:** ProductActivated
- **event:** ProductDeactivated
- **event:** ProductDeleted
- **event:** ProductVariantAdded
- **event:** CategoryCreated
- **event:** CategoryMoved

### Owned Aggregates
- Product
- Category

### Dependencies
- **Storage Service** (sync): Upload product images
- **Order Service** (sync): Check order history for soft delete

---

## SVC-INVENTORY: Inventory Service

**Purpose:** Tracks stock levels, reservations, and provides audit trail

**API Base:** `/api/v1/inventory`

**Separation Reason:** Inventory has distinct consistency requirements and audit needs

### Capabilities
- Track stock levels per product
- Reserve inventory for orders
- Release reserved inventory
- Deduct inventory on fulfillment
- Adjust inventory with audit
- Bulk update via CSV
- Trigger low stock alerts

### Inputs
- **command:** SetInventoryQuantity
- **command:** AdjustInventory
- **command:** ReserveInventory
- **command:** ReleaseInventory
- **command:** DeductInventory
- **command:** RestockInventory
- **command:** BulkUpdateInventory
- **event:** OrderPlaced
- **event:** OrderCancelled

### Outputs
- **event:** InventoryQuantitySet
- **event:** InventoryAdjusted
- **event:** InventoryReserved
- **event:** InventoryReleased
- **event:** InventoryDeducted
- **event:** LowStockAlert
- **event:** OutOfStock
- **query:** GetAvailableQuantity

### Owned Aggregates
- Inventory

### Dependencies
- **Notification Service** (async): Send low stock alerts

---

## SVC-CART: Shopping Cart Service

**Purpose:** Manages shopping cart operations and checkout preparation

**API Base:** `/api/v1/cart`

**Separation Reason:** Cart is transient state with different lifecycle than orders

### Capabilities
- Add items to cart
- Update item quantities
- Remove items from cart
- Merge guest cart on login
- Validate cart for checkout
- Sync prices with catalog
- Expire inactive carts

### Inputs
- **command:** AddItemToCart
- **command:** UpdateCartItemQuantity
- **command:** RemoveItemFromCart
- **command:** ClearCart
- **command:** MergeGuestCart
- **command:** RefreshCartPrices
- **event:** ProductDeleted
- **event:** ProductUpdated

### Outputs
- **event:** ItemAddedToCart
- **event:** CartItemQuantityUpdated
- **event:** ItemRemovedFromCart
- **event:** CartCleared
- **event:** CartMerged
- **event:** CartExpired
- **query:** GetCart

### Owned Aggregates
- Cart

### Dependencies
- **Product Catalog Service** (sync): Get product details and prices
- **Inventory Service** (sync): Check stock availability
- **Customer Service** (sync): Associate cart with customer

---

## SVC-ORDER: Order Service

**Purpose:** Manages complete order lifecycle from placement to delivery

**API Base:** `/api/v1/orders`

**Separation Reason:** Order lifecycle is core business process with strict invariants

### Capabilities
- Place orders from cart
- Track order status
- Update order status
- Cancel orders
- Calculate shipping and tax
- Generate order numbers
- View order history

### Inputs
- **command:** PlaceOrder
- **command:** ConfirmOrder
- **command:** ShipOrder
- **command:** DeliverOrder
- **command:** CancelOrder

### Outputs
- **event:** OrderPlaced
- **event:** OrderConfirmed
- **event:** OrderShipped
- **event:** OrderDelivered
- **event:** OrderCancelled
- **query:** GetOrder
- **query:** GetOrderHistory

### Owned Aggregates
- Order

### Dependencies
- **Customer Service** (sync): Verify customer eligibility
- **Shopping Cart Service** (sync): Get cart items for order
- **Inventory Service** (sync): Reserve and deduct stock
- **Payment Service** (sync): Authorize and capture payment
- **Email Service** (async): Send order notifications

---

## SVC-PAYMENT: Payment Service

**Purpose:** Handles payment authorization, capture, and refunds

**API Base:** `/api/v1/payments`

**Separation Reason:** Payment processing requires PCI compliance isolation

### Capabilities
- Authorize credit card payments
- Process PayPal payments
- Capture authorized payments
- Process refunds
- Save payment methods
- Validate card types

### Inputs
- **command:** AuthorizePayment
- **command:** CapturePayment
- **command:** RefundPayment
- **command:** SavePaymentMethod
- **command:** InitiatePayPal
- **event:** OrderCancelled

### Outputs
- **event:** PaymentAuthorized
- **event:** PaymentCaptured
- **event:** PaymentRefunded
- **event:** PaymentFailed
- **event:** PaymentMethodSaved

### Dependencies
- **Stripe Gateway** (sync): Process credit card payments
- **PayPal Gateway** (sync): Process PayPal payments

---

## SVC-EMAIL: Email Service

**Purpose:** Handles async email delivery with deduplication

**API Base:** `/api/v1/emails`

**Separation Reason:** Email delivery is async infrastructure concern

### Capabilities
- Queue confirmation emails
- Send verification emails
- Send status update emails
- Retry failed emails
- Prevent duplicate sends

### Inputs
- **event:** CustomerRegistered
- **event:** OrderPlaced
- **event:** OrderShipped
- **event:** OrderCancelled
- **command:** SendEmail

### Outputs
- **event:** EmailSent
- **event:** EmailFailed

### Dependencies
- **Message Queue** (async): Queue emails for processing
- **Email Provider** (sync): Deliver emails via SMTP/API

---

## SVC-ADMIN: Admin Service

**Purpose:** Provides admin-specific operations and reporting

**API Base:** `/api/v1/admin`

**Separation Reason:** Admin operations have different auth and audit requirements

### Capabilities
- Filter and view all orders
- Export orders to CSV
- Manage order fulfillment
- View audit logs

### Inputs
- **query:** ListOrders
- **query:** ExportOrdersCsv
- **command:** UpdateOrderStatus

### Outputs
- **query:** FilteredOrderList
- **file:** OrdersCsvExport

### Dependencies
- **Order Service** (sync): Query and update orders
- **Inventory Service** (sync): View and update inventory
- **Product Catalog Service** (sync): Manage products

---

