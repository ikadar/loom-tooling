# Service Boundaries

Generated: 2026-01-02T19:24:29+01:00

---

## SVC-CUSTOMER: Customer Service

**Purpose:** Manages customer registration, authentication, and profile data

**API Base:** `/api/v1/customers`

**Separation Reason:** Customer identity and auth is a distinct bounded context

### Capabilities
- Register new customers
- Verify customer email
- Update customer profile
- Authenticate customers

### Inputs
- **command:** RegisterCustomer
- **command:** VerifyCustomerEmail
- **command:** ChangeCustomerEmail
- **command:** ChangeCustomerPassword

### Outputs
- **event:** CustomerRegistered
- **event:** EmailVerified
- **event:** CustomerEmailChanged
- **event:** CustomerPasswordChanged

### Owned Aggregates
- Customer

### Dependencies
- **Cart Service** (sync): Create cart on registration

---

## SVC-CART: Cart Service

**Purpose:** Manages shopping cart items with pricing and stock validation

**API Base:** `/api/v1/carts`

**Separation Reason:** Shopping cart has distinct lifecycle from orders

### Capabilities
- Add items to cart
- Update cart quantities
- Remove items from cart
- Clear cart
- Calculate cart totals

### Inputs
- **command:** AddItemToCart
- **command:** UpdateCartItemQuantity
- **command:** RemoveItemFromCart
- **command:** ClearCart

### Outputs
- **event:** ItemAddedToCart
- **event:** CartQuantityUpdated
- **event:** ItemRemovedFromCart
- **event:** CartCleared

### Owned Aggregates
- Cart

### Dependencies
- **Product Service** (sync): Get product details and prices
- **Inventory Service** (sync): Check stock availability

---

## SVC-ORDER: Order Service

**Purpose:** Manages order lifecycle from placement through fulfillment

**API Base:** `/api/v1/orders`

**Separation Reason:** Order fulfillment is a distinct bounded context with state machine

### Capabilities
- Place orders from cart
- Confirm orders
- Ship orders with tracking
- Deliver orders
- Cancel orders

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

### Owned Aggregates
- Order

### Dependencies
- **Customer Service** (sync): Verify customer registration
- **Cart Service** (sync): Get cart items for order
- **Inventory Service** (sync): Reserve stock on placement
- **Product Service** (sync): Capture price snapshots

---

## SVC-PRODUCT: Product Service

**Purpose:** Manages product catalog with variants, pricing, and images

**API Base:** `/api/v1/products`

**Separation Reason:** Product catalog is independent domain with rich entity model

### Capabilities
- Create and update products
- Manage product variants
- Manage product images
- Activate/deactivate products
- Browse product catalog

### Inputs
- **command:** CreateProduct
- **command:** UpdateProduct
- **command:** DeactivateProduct
- **command:** DeleteProduct
- **command:** AddProductVariant
- **command:** AddProductImage
- **command:** SetPrimaryProductImage

### Outputs
- **event:** ProductCreated
- **event:** ProductUpdated
- **event:** ProductDeactivated
- **event:** ProductDeleted
- **event:** VariantAdded
- **event:** ImageAdded
- **event:** PrimaryImageChanged

### Owned Aggregates
- Product
- Category

---

## SVC-CATEGORY: Category Service

**Purpose:** Organizes products into hierarchical browsable categories

**API Base:** `/api/v1/categories`

**Separation Reason:** Category taxonomy is independent of product details

### Capabilities
- Create categories
- Update category details
- Manage category hierarchy
- Browse category tree

### Inputs
- **command:** CreateCategory
- **command:** UpdateCategory
- **command:** SetCategoryParent

### Outputs
- **event:** CategoryCreated
- **event:** CategoryUpdated
- **event:** CategoryParentChanged

### Owned Aggregates
- Category

---

## SVC-INVENTORY: Inventory Service

**Purpose:** Tracks stock levels and manages reservations for products

**API Base:** `/api/v1/inventory`

**Separation Reason:** Stock management requires pessimistic locking and high concurrency

### Capabilities
- Check stock availability
- Reserve inventory for orders
- Release reserved inventory
- Deduct inventory on fulfillment
- Restock inventory

### Inputs
- **command:** ReserveInventory
- **command:** ReleaseInventory
- **command:** DeductInventory
- **command:** RestockInventory
- **command:** RestoreInventory

### Outputs
- **event:** InventoryReserved
- **event:** InventoryReleased
- **event:** InventoryDeducted
- **event:** InventoryRestocked
- **event:** InventoryRestored
- **event:** OutOfStock

### Owned Aggregates
- Inventory

### Dependencies
- **Product Service** (async): React to product creation

---

