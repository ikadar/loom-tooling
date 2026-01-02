# Service Boundaries

Generated: 2026-01-02T21:37:12+01:00

---

## SVC-CUSTOMER: Customer Service

**Purpose:** Manages customer identity, registration, authentication, and profile data

**API Base:** `/api/v1/customers`

**Separation Reason:** Core identity domain with security requirements separate from transactional services

### Capabilities
- Register new customers
- Verify customer email
- Authenticate customers
- Manage customer profiles
- Change customer password

### Inputs
- **command:** RegisterCustomer
- **command:** VerifyCustomerEmail
- **command:** ChangeCustomerEmail
- **command:** ChangeCustomerPassword
- **query:** GetCustomer
- **query:** AuthenticateCustomer

### Outputs
- **event:** CustomerRegistered
- **event:** CustomerEmailVerified
- **event:** CustomerEmailChanged
- **event:** CustomerPasswordChanged
- **response:** CustomerProfile

### Owned Aggregates
- Customer

---

## SVC-CATALOG: Catalog Service

**Purpose:** Manages product information, categorization, variants, and images for browsing

**API Base:** `/api/v1/catalog`

**Separation Reason:** Read-heavy browsing workload with different scaling needs than transactional services

### Capabilities
- Create and manage products
- Organize products into categories
- Define product variants
- Manage product images
- Activate and deactivate products
- Soft-delete products with order history

### Inputs
- **command:** CreateProduct
- **command:** UpdateProduct
- **command:** DeactivateProduct
- **command:** ActivateProduct
- **command:** DeleteProduct
- **command:** AddProductVariant
- **command:** RemoveProductVariant
- **command:** AddProductImage
- **command:** SetPrimaryProductImage
- **command:** CreateCategory
- **command:** RenameCategory
- **command:** MoveCategoryToParent
- **command:** DeleteCategory
- **query:** GetProduct
- **query:** ListProducts
- **query:** GetCategory
- **query:** ListCategories

### Outputs
- **event:** ProductCreated
- **event:** ProductUpdated
- **event:** ProductDeactivated
- **event:** ProductActivated
- **event:** ProductDeleted
- **event:** ProductVariantAdded
- **event:** ProductVariantRemoved
- **event:** ProductImageAdded
- **event:** ProductPrimaryImageChanged
- **event:** CategoryCreated
- **event:** CategoryRenamed
- **event:** CategoryMoved
- **event:** CategoryDeleted
- **response:** ProductDetails
- **response:** ProductList
- **response:** CategoryTree

### Owned Aggregates
- Product
- Category

---

## SVC-SHOPPING: Shopping Cart Service

**Purpose:** Manages shopping cart lifecycle with items, quantities, and pricing before checkout

**API Base:** `/api/v1/carts`

**Separation Reason:** High-frequency temporary state with different persistence and caching patterns

### Capabilities
- Add items to cart
- Update item quantities
- Remove items from cart
- Calculate cart totals
- Refresh cart prices
- Clear cart on checkout

### Inputs
- **command:** AddItemToCart
- **command:** UpdateCartItemQuantity
- **command:** RemoveItemFromCart
- **command:** ClearCart
- **command:** RefreshCartPrices
- **query:** GetCart

### Outputs
- **event:** ItemAddedToCart
- **event:** CartQuantityUpdated
- **event:** ItemRemovedFromCart
- **event:** CartCleared
- **event:** CartPricesRefreshed
- **response:** CartDetails

### Owned Aggregates
- Cart

### Dependencies
- **Catalog Service** (sync): Fetch product details and current prices for cart items
- **Inventory Service** (sync): Check stock availability before adding or updating items
- **Customer Service** (sync): Validate customer exists for cart association

---

## SVC-ORDER: Order Service

**Purpose:** Manages complete order lifecycle from placement through fulfillment and delivery

**API Base:** `/api/v1/orders`

**Separation Reason:** Critical transactional service with audit requirements and complex state machine

### Capabilities
- Create orders from cart
- Process order confirmation
- Track order status transitions
- Handle order cancellation
- Calculate totals with shipping
- Apply free shipping rules

### Inputs
- **command:** PlaceOrder
- **command:** ConfirmOrder
- **command:** ShipOrder
- **command:** DeliverOrder
- **command:** CancelOrder
- **query:** GetOrder
- **query:** ListCustomerOrders

### Outputs
- **event:** OrderPlaced
- **event:** OrderConfirmed
- **event:** OrderShipped
- **event:** OrderDelivered
- **event:** OrderCancelled
- **response:** OrderDetails
- **response:** OrderList

### Owned Aggregates
- Order

### Dependencies
- **Customer Service** (sync): Verify customer registration and email verification status
- **Shopping Cart Service** (sync): Retrieve cart items and clear cart after order placement
- **Catalog Service** (sync): Capture product snapshots via ACL for order line items
- **Inventory Service** (async): Reserve stock on order placement, release on cancellation

---

## SVC-INVENTORY: Inventory Service

**Purpose:** Tracks stock levels, reservations, and product availability across the system

**API Base:** `/api/v1/inventory`

**Separation Reason:** Requires pessimistic locking and high consistency for stock operations

### Capabilities
- Track quantity on hand
- Reserve stock for pending orders
- Release reserved stock on cancellation
- Deduct stock on fulfillment
- Restock products
- Adjust inventory with audit trail
- Check product availability

### Inputs
- **command:** ReserveInventory
- **command:** ReleaseInventory
- **command:** DeductInventory
- **command:** RestockInventory
- **command:** AdjustInventory
- **query:** GetInventory
- **query:** CheckAvailability

### Outputs
- **event:** InventoryReserved
- **event:** InventoryReleased
- **event:** InventoryDeducted
- **event:** InventoryRestocked
- **event:** InventoryAdjusted
- **event:** OutOfStock
- **response:** InventoryStatus
- **response:** AvailabilityCheck

### Owned Aggregates
- Inventory

### Dependencies
- **Catalog Service** (sync): Shared kernel for ProductId; validates product exists for inventory records

---

