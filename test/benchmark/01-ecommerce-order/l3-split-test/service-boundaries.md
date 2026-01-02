# Service Boundaries

Generated: 2026-01-02T20:15:12+01:00

---

## SVC-CUSTOMER: Customer Service

**Purpose:** Manages customer registration, authentication, and profile data

**API Base:** `/api/v1/customers`

**Separation Reason:** Customer identity and auth are core domain with distinct security requirements

### Capabilities
- Register new customers
- Verify customer email
- Manage customer profile
- Authenticate customers

### Inputs
- **command:** RegisterCustomer
- **command:** VerifyCustomerEmail
- **command:** ChangeCustomerEmail
- **command:** ChangeCustomerPassword
- **query:** GetCustomer

### Outputs
- **event:** CustomerRegistered
- **event:** EmailVerified
- **event:** CustomerEmailChanged
- **event:** CustomerPasswordChanged

### Owned Aggregates
- Customer

### Dependencies
- **Cart Service** (sync): Creates cart on customer registration

---

## SVC-CART: Cart Service

**Purpose:** Manages shopping cart items with real-time pricing and stock validation

**API Base:** `/api/v1/carts`

**Separation Reason:** Cart is transient pre-order state with different persistence and caching needs

### Capabilities
- Add items to cart
- Update item quantities
- Remove items from cart
- Clear cart
- Calculate cart totals

### Inputs
- **command:** AddItemToCart
- **command:** UpdateCartItemQuantity
- **command:** RemoveItemFromCart
- **command:** ClearCart
- **query:** GetCart

### Outputs
- **event:** ItemAddedToCart
- **event:** CartQuantityUpdated
- **event:** ItemRemovedFromCart
- **event:** CartCleared

### Owned Aggregates
- Cart

### Dependencies
- **Product Service** (sync): Fetches current product prices and details
- **Inventory Service** (sync): Validates stock availability

---

## SVC-ORDER: Order Service

**Purpose:** Manages complete order lifecycle from placement through delivery

**API Base:** `/api/v1/orders`

**Separation Reason:** Order is core business transaction with complex state machine and audit needs

### Capabilities
- Place orders
- Confirm orders
- Ship orders
- Deliver orders
- Cancel orders
- Track order status

### Inputs
- **command:** PlaceOrder
- **command:** ConfirmOrder
- **command:** ShipOrder
- **command:** DeliverOrder
- **command:** CancelOrder
- **query:** GetOrder
- **query:** GetCustomerOrders

### Outputs
- **event:** OrderPlaced
- **event:** OrderConfirmed
- **event:** OrderShipped
- **event:** OrderDelivered
- **event:** OrderCancelled

### Owned Aggregates
- Order

### Dependencies
- **Customer Service** (sync): Validates customer registration status
- **Cart Service** (sync): Retrieves cart items for order creation
- **Inventory Service** (sync): Reserves stock during order placement
- **Payment Service** (sync): Authorizes payment before order creation

---

## SVC-PRODUCT: Product Service

**Purpose:** Manages product catalog with variants, pricing, and categorization

**API Base:** `/api/v1/products`

**Separation Reason:** Product catalog is read-heavy with caching needs, separate from transactions

### Capabilities
- Create and update products
- Manage product variants
- Manage product images
- Deactivate and delete products
- Browse product catalog

### Inputs
- **command:** CreateProduct
- **command:** UpdateProduct
- **command:** DeactivateProduct
- **command:** DeleteProduct
- **command:** AddProductVariant
- **command:** AddProductImage
- **command:** SetPrimaryProductImage
- **query:** GetProduct
- **query:** ListProducts

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

**Purpose:** Organizes products into browsable hierarchical groupings

**API Base:** `/api/v1/categories`

**Separation Reason:** Category hierarchy is mostly static with high read:write ratio, cacheable

### Capabilities
- Create categories
- Update categories
- Manage category hierarchy
- Browse category tree

### Inputs
- **command:** CreateCategory
- **command:** UpdateCategory
- **command:** SetCategoryParent
- **query:** GetCategory
- **query:** GetCategoryTree

### Outputs
- **event:** CategoryCreated
- **event:** CategoryUpdated
- **event:** CategoryParentChanged

### Owned Aggregates
- Category

---

## SVC-INVENTORY: Inventory Service

**Purpose:** Tracks stock levels and manages availability with reservation support

**API Base:** `/api/v1/inventory`

**Separation Reason:** Inventory requires pessimistic locking and high consistency for stock accuracy

### Capabilities
- Check stock availability
- Reserve inventory
- Release reservations
- Deduct inventory
- Restock inventory

### Inputs
- **command:** ReserveInventory
- **command:** ReleaseInventory
- **command:** DeductInventory
- **command:** RestockInventory
- **command:** RestoreInventory
- **query:** GetInventory
- **query:** CheckStock

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
- **Order Service** (async): Listens for OrderCancelled to restore stock

---

