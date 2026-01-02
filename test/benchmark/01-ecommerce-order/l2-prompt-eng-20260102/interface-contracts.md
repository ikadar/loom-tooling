# Interface Contracts

Generated: 2026-01-02T18:27:24+01:00

---

## Shared Types

### Money

| Field | Type | Constraints |
|-------|------|-------------|
| amount | decimal | Precision 2 decimal places |
| currency | string | ISO 4217 code, default USD |

### ShippingAddress

| Field | Type | Constraints |
|-------|------|-------------|
| street | string | Required, 1-200 characters |
| city | string | Required, 1-100 characters |
| state | string | Required, 2-100 characters |
| postalCode | string | Required, valid postal format |
| country | string | Required, ISO 3166 code |
| recipientName | string | Required |

### PaymentMethod

| Field | Type | Constraints |
|-------|------|-------------|
| type | PaymentType | Required, CREDIT_CARD or PAYPAL |
| lastFourDigits | string | Required for CREDIT_CARD, 4 digits |
| expiryMonth | integer | Required for CREDIT_CARD, 1-12 |
| expiryYear | integer | Required for CREDIT_CARD |
| paypalEmail | Email | Required for PAYPAL |

### Email

| Field | Type | Constraints |
|-------|------|-------------|
| value | string | Valid email format per RFC 5322 |

### OrderStatus

| Field | Type | Constraints |
|-------|------|-------------|
| value | enum | PENDING, CONFIRMED, SHIPPED, DELIVERED, CANCELLED |

### CartItem

| Field | Type | Constraints |
|-------|------|-------------|
| cartItemId | CartItemId (UUID) | Required, immutable |
| productId | ProductId (UUID) | Required |
| variantId | VariantId (UUID) | Optional |
| productName | string | Snapshot at time of add |
| unitPrice | Money | Current price |
| quantity | integer | Required, >= 1 |
| subtotal | Money | Calculated: unitPrice * quantity |

### OrderLineItem

| Field | Type | Constraints |
|-------|------|-------------|
| lineItemId | LineItemId (UUID) | Required, immutable |
| productId | ProductId (UUID) | Required |
| productName | string | Snapshot at order time |
| unitPrice | Money | Snapshot at order time, immutable |
| quantity | integer | Required, >= 1 |
| subtotal | Money | Calculated: unitPrice * quantity |

### ProductVariant

| Field | Type | Constraints |
|-------|------|-------------|
| variantId | VariantId (UUID) | Required, immutable |
| sku | string | Required, unique |
| size | string | Optional |
| color | string | Optional |
| priceAdjustment | Money | Optional, can be positive or negative |

---

## IC-PRODUCT-001 – Product Service {#ic-product-001}

**Purpose:** Manages product catalog, categories, and product variants for the e-commerce platform

**Base URL:** `/api/v1/products`

**Security:**
- Authentication: Bearer JWT
- Authorization: Admin role required for create/update/delete operations; public read for active products

### Operations

#### listProducts `GET /`

Browse products with filtering, sorting, and pagination

**Input:**
- `page`: integer
- `pageSize`: integer
- `categoryId`: CategoryId (UUID)
- `minPrice`: Money
- `maxPrice`: Money
- `inStock`: boolean
- `sortBy`: enum(name, price_asc, price_desc, newest)

**Output:**
- `products`: Product[]
- `totalCount`: integer
- `page`: integer
- `pageSize`: integer
- `hasMore`: boolean

**Errors:**
| Code | HTTP | Message |
|------|------|----------|
| INVALID_PRICE_RANGE | 400 | Minimum price cannot exceed maximum price |

**Postconditions:** [Products matching filters are returned sorted and paginated]

#### getProduct `GET /{productId}`

Get detailed product information including variants

**Input:**
- `productId`: ProductId (UUID) (required)

**Output:**
- `categoryId`: CategoryId
- `availableQuantity`: integer
- `name`: string
- `variants`: ProductVariant[]
- `productId`: ProductId (UUID)
- `images`: ProductImage[]
- `inStock`: boolean
- `description`: string
- `isActive`: boolean
- `price`: Money

**Errors:**
| Code | HTTP | Message |
|------|------|----------|
| PRODUCT_NOT_FOUND | 404 | Product with specified ID does not exist |

**Preconditions:** [Product must exist and be active for customer view]

**Postconditions:** [Product details with variant stock status returned]

#### createProduct `POST /`

Create a new product in draft status

**Input:**
- `name`: string (required)
- `description`: string
- `price`: Money (required)
- `categoryId`: CategoryId (UUID) (required)

**Output:**
- `productId`: ProductId (UUID)
- `name`: string
- `status`: enum(draft, active)
- `createdAt`: DateTime
- `createdBy`: UserId

**Errors:**
| Code | HTTP | Message |
|------|------|----------|
| INVALID_PRICE | 400 | Price must be at least $0.01 |
| INVALID_PRODUCT_NAME | 400 | Product name must be between 2 and 200 characters |
| DESCRIPTION_TOO_LONG | 400 | Product description must not exceed 5000 characters |
| CATEGORY_NOT_FOUND | 404 | Specified category does not exist |

**Preconditions:** [User must be authenticated as admin Category must exist]

**Postconditions:** [Product created in draft status ProductCreated event emitted]

#### updateProduct `PUT /{productId}`

Update existing product attributes

**Input:**
- `productId`: ProductId (UUID) (required)
- `name`: string
- `description`: string
- `price`: Money
- `categoryId`: CategoryId (UUID)

**Output:**
- `categoryId`: CategoryId
- `updatedAt`: DateTime
- `updatedBy`: UserId
- `productId`: ProductId (UUID)
- `name`: string
- `description`: string
- `price`: Money

**Errors:**
| Code | HTTP | Message |
|------|------|----------|
| PRODUCT_NOT_FOUND | 404 | Product with specified ID does not exist |
| INVALID_PRICE | 400 | Price must be at least $0.01 |
| INVALID_PRODUCT_NAME | 400 | Product name must be between 2 and 200 characters |
| DESCRIPTION_TOO_LONG | 400 | Product description must not exceed 5000 characters |
| CATEGORY_NOT_FOUND | 404 | Specified category does not exist |

**Preconditions:** [User must be authenticated as admin Product must exist]

**Postconditions:** [Product updated Price changes logged in audit trail ProductUpdated event emitted Carts with this product reflect new price]

#### deleteProduct `DELETE /{productId}`

Soft delete a product (preserved for order history)

**Input:**
- `productId`: ProductId (UUID) (required)

**Output:**
- `success`: boolean
- `deletedAt`: DateTime

**Errors:**
| Code | HTTP | Message |
|------|------|----------|
| PRODUCT_NOT_FOUND | 404 | Product with specified ID does not exist |

**Preconditions:** [User must be authenticated as admin Product must exist]

**Postconditions:** [Product soft-deleted Removed from carts No longer appears in listings ProductDeactivated event emitted]

#### addProductImages `POST /{productId}/images`

Add images to a product (max 10 with one primary)

**Input:**
- `productId`: ProductId (UUID) (required)
- `images`: ImageUpload[] (required)
- `primaryImageIndex`: integer

**Output:**
- `primaryImage`: ProductImage
- `images`: ProductImage[]

**Errors:**
| Code | HTTP | Message |
|------|------|----------|
| PRODUCT_NOT_FOUND | 404 | Product with specified ID does not exist |
| IMAGE_LIMIT_EXCEEDED | 400 | Products cannot have more than 10 images |

**Preconditions:** [User must be authenticated as admin Product must exist Total images after add must not exceed 10]

**Postconditions:** [Images added to product Exactly one image marked as primary]

#### addProductVariant `POST /{productId}/variants`

Add a variant to a product with unique SKU

**Input:**
- `size`: string
- `color`: string
- `priceAdjustment`: Money
- `productId`: ProductId (UUID) (required)
- `sku`: string (required)

**Output:**
- `size`: string
- `color`: string
- `priceAdjustment`: Money
- `variantId`: VariantId (UUID)
- `sku`: string

**Errors:**
| Code | HTTP | Message |
|------|------|----------|
| PRODUCT_NOT_FOUND | 404 | Product with specified ID does not exist |
| DUPLICATE_SKU | 409 | A variant with this SKU already exists |
| VARIANT_ATTRIBUTES_REQUIRED | 400 | At least one variant attribute (size or color) must be specified |

**Preconditions:** [User must be authenticated as admin Product must exist SKU must be unique]

**Postconditions:** [Variant added to product]

### Events

- **ProductCreated**: Emitted when a new product is created (payload: [productId name price categoryId createdAt createdBy])
- **ProductUpdated**: Emitted when product attributes are modified (payload: [productId changedFields previousPrice newPrice updatedAt updatedBy])
- **ProductDeactivated**: Emitted when a product is soft-deleted (payload: [productId deletedAt deletedBy])

---

## IC-CATEGORY-001 – Category Service {#ic-category-001}

**Purpose:** Organizes products into hierarchical browsable groupings

**Base URL:** `/api/v1/categories`

**Security:**
- Authentication: Bearer JWT
- Authorization: Admin role required for create/update/delete; public read for active categories

### Operations

#### listCategories `GET /`

List all categories in hierarchical structure

**Input:**
- `parentId`: CategoryId (UUID)
- `includeInactive`: boolean

**Output:**
- `categories`: Category[]

**Postconditions:** [Categories returned in hierarchical structure up to 3 levels]

#### getCategory `GET /{categoryId}`

Get category details with products

**Input:**
- `categoryId`: CategoryId (UUID) (required)

**Output:**
- `parentCategoryId`: CategoryId (UUID)
- `subcategories`: Category[]
- `productCount`: integer
- `categoryId`: CategoryId (UUID)
- `name`: string
- `description`: string

**Errors:**
| Code | HTTP | Message |
|------|------|----------|
| CATEGORY_NOT_FOUND | 404 | Category with specified ID does not exist |

**Preconditions:** [Category must exist]

**Postconditions:** [Category details with subcategories returned]

#### createCategory `POST /`

Create a new category

**Input:**
- `name`: string (required)
- `description`: string
- `parentCategoryId`: CategoryId (UUID)

**Output:**
- `categoryId`: CategoryId (UUID)
- `name`: string
- `depth`: integer

**Errors:**
| Code | HTTP | Message |
|------|------|----------|
| DUPLICATE_CATEGORY_NAME | 409 | A category with this name already exists |
| PARENT_CATEGORY_NOT_FOUND | 404 | Parent category does not exist |
| CATEGORY_DEPTH_EXCEEDED | 400 | Categories cannot be nested more than 3 levels deep |

**Preconditions:** [User must be authenticated as admin Category name must be unique If parent specified, parent must exist Resulting depth must not exceed 3]

**Postconditions:** [Category created CategoryCreated event emitted]

### Events

- **CategoryCreated**: Emitted when a new category is created (payload: [categoryId name parentCategoryId depth createdAt])

---

## IC-CART-001 – Cart Service {#ic-cart-001}

**Purpose:** Manages shopping carts for customers, handling item additions, updates, and cart lifecycle

**Base URL:** `/api/v1/carts`

**Security:**
- Authentication: Bearer JWT or session token for guest carts
- Authorization: User can only access their own cart

### Operations

#### getCart `GET /`

Get current user's cart with items and totals

**Output:**
- `customerId`: CustomerId (UUID)
- `items`: CartItem[]
- `totalPrice`: Money
- `itemCount`: integer
- `priceChanges`: PriceChange[]
- `stockWarnings`: StockWarning[]
- `cartId`: CartId (UUID)

**Postconditions:** [Cart returned with current prices Price changes since last view indicated Out-of-stock items flagged]

#### addToCart `POST /items`

Add a product to the cart

**Input:**
- `quantity`: integer (required)
- `productId`: ProductId (UUID) (required)
- `variantId`: VariantId (UUID)

**Output:**
- `itemCount`: integer
- `cartId`: CartId (UUID)
- `cartItem`: CartItem
- `totalPrice`: Money

**Errors:**
| Code | HTTP | Message |
|------|------|----------|
| PRODUCT_NOT_FOUND | 404 | Product with specified ID does not exist |
| OUT_OF_STOCK | 409 | Product is out of stock |
| QUANTITY_EXCEEDS_STOCK | 409 | Requested quantity exceeds available stock |
| VARIANT_REQUIRED | 400 | Product has variants; variant selection is required |
| VARIANT_NOT_FOUND | 404 | Specified variant does not exist |
| PRODUCT_INACTIVE | 409 | Product is not available for purchase |

**Preconditions:** [Product must exist and be active If product has variants, variant must be specified Quantity must be greater than 0 Available stock must be sufficient]

**Postconditions:** [Item added to cart or quantity updated if already exists Cart total recalculated ItemAddedToCart event emitted]

#### updateCartItemQuantity `PATCH /items/{cartItemId}`

Update quantity of a cart item

**Input:**
- `cartItemId`: CartItemId (UUID) (required)
- `quantity`: integer (required)

**Output:**
- `totalPrice`: Money
- `cartItem`: CartItem

**Errors:**
| Code | HTTP | Message |
|------|------|----------|
| CART_ITEM_NOT_FOUND | 404 | Cart item with specified ID does not exist |
| QUANTITY_EXCEEDS_STOCK | 409 | Requested quantity exceeds available stock |
| INVALID_QUANTITY | 400 | Quantity must be at least 0 |

**Preconditions:** [Cart item must exist Quantity must be >= 0]

**Postconditions:** [Quantity updated (or item removed if quantity is 0) Cart total recalculated CartQuantityUpdated event emitted]

#### removeFromCart `DELETE /items/{cartItemId}`

Remove an item from the cart

**Input:**
- `cartItemId`: CartItemId (UUID) (required)

**Output:**
- `success`: boolean
- `undoToken`: string
- `undoExpiresAt`: DateTime

**Errors:**
| Code | HTTP | Message |
|------|------|----------|
| CART_ITEM_NOT_FOUND | 404 | Cart item with specified ID does not exist |

**Preconditions:** [Cart item must exist]

**Postconditions:** [Item removed from cart Cart total recalculated Undo option available briefly ItemRemovedFromCart event emitted]

#### undoRemoveFromCart `POST /items/undo`

Restore a recently removed cart item

**Input:**
- `undoToken`: string (required)

**Output:**
- `totalPrice`: Money
- `cartItem`: CartItem

**Errors:**
| Code | HTTP | Message |
|------|------|----------|
| UNDO_TOKEN_EXPIRED | 410 | Undo window has expired |
| UNDO_TOKEN_INVALID | 400 | Invalid undo token |

**Preconditions:** [Undo token must be valid and not expired]

**Postconditions:** [Item restored to cart Cart total recalculated]

#### clearCart `DELETE /`

Remove all items from the cart

**Output:**
- `success`: boolean

**Postconditions:** [All items removed from cart CartCleared event emitted]

#### mergeGuestCart `POST /merge`

Merge guest cart into authenticated user's cart on login

**Input:**
- `guestCartId`: CartId (UUID) (required)

**Output:**
- `mergeConflicts`: MergeConflict[]
- `cartId`: CartId (UUID)
- `items`: CartItem[]
- `totalPrice`: Money

**Errors:**
| Code | HTTP | Message |
|------|------|----------|
| GUEST_CART_NOT_FOUND | 404 | Guest cart with specified ID does not exist |
| GUEST_CART_EXPIRED | 410 | Guest cart has expired |

**Preconditions:** [User must be authenticated Guest cart must exist]

**Postconditions:** [Guest cart items merged into user cart Quantities combined for duplicate products Stock limits enforced with notifications]

### Events

- **ItemAddedToCart**: Emitted when an item is added to cart (payload: [cartId customerId productId variantId quantity unitPrice])
- **CartQuantityUpdated**: Emitted when cart item quantity changes (payload: [cartId cartItemId previousQuantity newQuantity])
- **ItemRemovedFromCart**: Emitted when an item is removed from cart (payload: [cartId cartItemId productId])
- **CartCleared**: Emitted when cart is emptied (payload: [cartId customerId reason])

---

## IC-ORDER-001 – Order Service {#ic-order-001}

**Purpose:** Handles order placement, status management, and order history for customers and administrators

**Base URL:** `/api/v1/orders`

**Security:**
- Authentication: Bearer JWT
- Authorization: Customers can view/cancel own orders; Admin role required for status updates

### Operations

#### placeOrder `POST /`

Create a new order from the customer's cart

**Input:**
- `shippingAddress`: ShippingAddress
- `paymentMethodId`: PaymentMethodId (UUID)
- `paymentMethod`: PaymentMethod
- `shippingAddressId`: AddressId (UUID)

**Output:**
- `shippingCost`: Money
- `orderNumber`: string
- `totalAmount`: Money
- `orderId`: OrderId (UUID)
- `status`: OrderStatus
- `tax`: Money
- `shippingAddress`: ShippingAddress
- `lineItems`: OrderLineItem[]
- `subtotal`: Money
- `createdAt`: DateTime

**Errors:**
| Code | HTTP | Message |
|------|------|----------|
| REGISTRATION_REQUIRED | 403 | Only registered customers can place orders |
| EMAIL_NOT_VERIFIED | 403 | Email verification required before placing orders |
| CART_EMPTY | 400 | Cannot place order with empty cart |
| OUT_OF_STOCK | 409 | One or more items are out of stock |
| INSUFFICIENT_STOCK | 409 | Requested quantity exceeds available stock |
| BACKORDER_NOT_SUPPORTED | 409 | Backorders are not allowed |
| PAYMENT_REQUIRED | 402 | Payment authorization failed |
| INCOMPLETE_ADDRESS | 400 | Shipping address is missing required fields |
| INTERNATIONAL_SHIPPING_NOT_AVAILABLE | 400 | Shipping is only available domestically |

**Preconditions:** [Customer must be registered with verified email Cart must not be empty All items must be in stock Shipping address must be valid and domestic Payment must be authorized]

**Postconditions:** [Order created with status PENDING Inventory decremented Cart cleared Prices captured at order time OrderPlaced event emitted Confirmation email queued]

#### getOrder `GET /{orderId}`

Get order details by ID

**Input:**
- `orderId`: OrderId (UUID) (required)

**Output:**
- `status`: OrderStatus
- `subtotal`: Money
- `tax`: Money
- `statusHistory`: StatusChange[]
- `orderId`: OrderId (UUID)
- `orderNumber`: string
- `shippingAddress`: ShippingAddress
- `trackingNumber`: string
- `createdAt`: DateTime
- `lineItems`: OrderLineItem[]
- `shippingCost`: Money
- `totalAmount`: Money

**Errors:**
| Code | HTTP | Message |
|------|------|----------|
| ORDER_NOT_FOUND | 404 | Order with specified ID does not exist |
| ACCESS_DENIED | 403 | You do not have permission to view this order |

**Preconditions:** [Order must exist User must be order owner or admin]

**Postconditions:** [Order details returned]

#### listOrders `GET /`

List orders for current customer or all orders for admin

**Input:**
- `pageSize`: integer
- `status`: OrderStatus
- `startDate`: Date
- `endDate`: Date
- `customerId`: CustomerId (UUID)
- `orderNumber`: string
- `page`: integer

**Output:**
- `totalCount`: integer
- `page`: integer
- `pageSize`: integer
- `orders`: OrderSummary[]

**Preconditions:** [User must be authenticated]

**Postconditions:** [Orders returned with pagination]

#### cancelOrder `POST /{orderId}/cancel`

Cancel an order before it ships

**Input:**
- `orderId`: OrderId (UUID) (required)
- `reason`: string

**Output:**
- `status`: OrderStatus
- `cancelledAt`: DateTime
- `refundInitiated`: boolean
- `orderId`: OrderId (UUID)

**Errors:**
| Code | HTTP | Message |
|------|------|----------|
| ORDER_NOT_FOUND | 404 | Order with specified ID does not exist |
| CANCELLATION_NOT_ALLOWED | 409 | Order cannot be cancelled after shipping |
| ORDER_ALREADY_CANCELLED | 409 | Order has already been cancelled |

**Preconditions:** [Order must exist Order status must be PENDING or CONFIRMED User must be order owner or admin]

**Postconditions:** [Order status changed to CANCELLED Inventory restored Refund initiated OrderCancelled event emitted Cancellation email sent]

#### confirmOrder `POST /{orderId}/confirm`

Confirm a pending order (admin operation)

**Input:**
- `orderId`: OrderId (UUID) (required)

**Output:**
- `confirmedAt`: DateTime
- `orderId`: OrderId (UUID)
- `status`: OrderStatus

**Errors:**
| Code | HTTP | Message |
|------|------|----------|
| ORDER_NOT_FOUND | 404 | Order with specified ID does not exist |
| INVALID_STATUS_TRANSITION | 409 | Order can only be confirmed from pending status |

**Preconditions:** [User must be authenticated as admin Order must exist Order status must be PENDING]

**Postconditions:** [Order status changed to CONFIRMED OrderConfirmed event emitted Customer notified]

#### shipOrder `POST /{orderId}/ship`

Mark order as shipped with optional tracking number

**Input:**
- `orderId`: OrderId (UUID) (required)
- `trackingNumber`: string

**Output:**
- `trackingNumber`: string
- `shippedAt`: DateTime
- `orderId`: OrderId (UUID)
- `status`: OrderStatus

**Errors:**
| Code | HTTP | Message |
|------|------|----------|
| ORDER_NOT_FOUND | 404 | Order with specified ID does not exist |
| INVALID_STATUS_TRANSITION | 409 | Order can only be shipped from confirmed status |

**Preconditions:** [User must be authenticated as admin Order must exist Order status must be CONFIRMED]

**Postconditions:** [Order status changed to SHIPPED All items ship as single unit OrderShipped event emitted Customer notified]

#### deliverOrder `POST /{orderId}/deliver`

Mark order as delivered

**Input:**
- `orderId`: OrderId (UUID) (required)

**Output:**
- `status`: OrderStatus
- `deliveredAt`: DateTime
- `orderId`: OrderId (UUID)

**Errors:**
| Code | HTTP | Message |
|------|------|----------|
| ORDER_NOT_FOUND | 404 | Order with specified ID does not exist |
| INVALID_STATUS_TRANSITION | 409 | Order can only be delivered from shipped status |

**Preconditions:** [User must be authenticated as admin Order must exist Order status must be SHIPPED]

**Postconditions:** [Order status changed to DELIVERED OrderDelivered event emitted Customer notified]

### Events

- **OrderPlaced**: Emitted when a new order is created (payload: [orderId orderNumber customerId totalAmount lineItemCount createdAt])
- **OrderConfirmed**: Emitted when order is confirmed by admin (payload: [orderId confirmedAt confirmedBy])
- **OrderShipped**: Emitted when order is marked as shipped (payload: [orderId trackingNumber shippedAt shippedBy])
- **OrderDelivered**: Emitted when order is marked as delivered (payload: [orderId deliveredAt])
- **OrderCancelled**: Emitted when order is cancelled (payload: [orderId reason cancelledAt cancelledBy refundAmount])

---

## IC-CUSTOMER-001 – Customer Service {#ic-customer-001}

**Purpose:** Manages customer registration, authentication, profiles, and saved addresses/payment methods

**Base URL:** `/api/v1/customers`

**Security:**
- Authentication: Bearer JWT
- Authorization: Users can only access their own profile

### Operations

#### registerCustomer `POST /register`

Create a new customer account

**Input:**
- `email`: Email (required)
- `password`: string (required)
- `firstName`: string (required)
- `lastName`: string (required)

**Output:**
- `registrationStatus`: RegistrationStatus
- `emailVerified`: boolean
- `customerId`: CustomerId (UUID)
- `email`: Email

**Errors:**
| Code | HTTP | Message |
|------|------|----------|
| EMAIL_ALREADY_REGISTERED | 409 | An account with this email already exists |
| WEAK_PASSWORD | 400 | Password must be at least 8 characters and contain at least one number |
| INVALID_EMAIL | 400 | Invalid email format |

**Preconditions:** [Email must not be registered]

**Postconditions:** [Customer account created Verification email sent CustomerRegistered event emitted]

#### verifyEmail `POST /verify-email`

Verify customer email address

**Input:**
- `token`: string (required)

**Output:**
- `success`: boolean
- `emailVerified`: boolean

**Errors:**
| Code | HTTP | Message |
|------|------|----------|
| INVALID_VERIFICATION_TOKEN | 400 | Verification token is invalid or expired |

**Preconditions:** [Token must be valid and not expired]

**Postconditions:** [Email marked as verified Customer can place orders]

#### getCustomerProfile `GET /me`

Get current customer's profile

**Output:**
- `customerId`: CustomerId (UUID)
- `email`: Email
- `firstName`: string
- `lastName`: string
- `emailVerified`: boolean
- `addresses`: ShippingAddress[]
- `paymentMethods`: PaymentMethodSummary[]

**Preconditions:** [User must be authenticated]

**Postconditions:** [Customer profile returned]

#### addShippingAddress `POST /me/addresses`

Add a shipping address to customer profile

**Input:**
- `country`: string (required)
- `recipientName`: string (required)
- `isDefault`: boolean
- `street`: string (required)
- `city`: string (required)
- `state`: string (required)
- `postalCode`: string (required)

**Output:**
- `isDefault`: boolean
- `addressId`: AddressId (UUID)
- `address`: ShippingAddress

**Errors:**
| Code | HTTP | Message |
|------|------|----------|
| INCOMPLETE_ADDRESS | 400 | All required address fields must be provided |
| INVALID_POSTAL_CODE | 400 | Invalid postal code format |
| INTERNATIONAL_SHIPPING_NOT_AVAILABLE | 400 | Shipping is only available domestically |

**Preconditions:** [User must be authenticated Address must be domestic]

**Postconditions:** [Address saved to profile If marked default, previous default cleared]

#### addPaymentMethod `POST /me/payment-methods`

Add a payment method to customer profile

**Input:**
- `isDefault`: boolean
- `type`: PaymentType (required)
- `token`: string (required)

**Output:**
- `paymentMethodId`: PaymentMethodId (UUID)
- `type`: PaymentType
- `lastFourDigits`: string
- `isDefault`: boolean

**Errors:**
| Code | HTTP | Message |
|------|------|----------|
| INVALID_PAYMENT_TOKEN | 400 | Payment token validation failed |
| UNSUPPORTED_CARD_TYPE | 400 | Only Visa, Mastercard, and American Express are supported |

**Preconditions:** [User must be authenticated Payment token must be valid from payment gateway]

**Postconditions:** [Tokenized payment method saved If marked default, previous default cleared]

#### deleteCustomer `DELETE /me`

Request account deletion (GDPR compliance)

**Input:**
- `deleteType`: enum(soft_delete, full_erasure) (required)
- `confirmation`: string (required)

**Output:**
- `success`: boolean
- `scheduledDeletionDate`: DateTime

**Errors:**
| Code | HTTP | Message |
|------|------|----------|
| PENDING_ORDERS_EXIST | 409 | Cannot delete account while orders are pending |
| INVALID_CONFIRMATION | 400 | Deletion confirmation does not match |

**Preconditions:** [User must be authenticated No pending orders]

**Postconditions:** [Account deactivated or scheduled for erasure]

### Events

- **CustomerRegistered**: Emitted when a new customer registers (payload: [customerId email createdAt])

---

## IC-INVENTORY-001 – Inventory Service {#ic-inventory-001}

**Purpose:** Tracks and manages stock levels for products and variants

**Base URL:** `/api/v1/inventory`

**Security:**
- Authentication: Bearer JWT
- Authorization: Admin role required for all operations

### Operations

#### getInventory `GET /{productId}`

Get inventory levels for a product

**Input:**
- `productId`: ProductId (UUID) (required)

**Output:**
- `availableQuantity`: integer
- `inventoryId`: InventoryId (UUID)
- `productId`: ProductId (UUID)
- `quantityOnHand`: integer
- `reservedQuantity`: integer

**Errors:**
| Code | HTTP | Message |
|------|------|----------|
| PRODUCT_NOT_FOUND | 404 | Product with specified ID does not exist |

**Preconditions:** [Product must exist]

**Postconditions:** [Inventory levels returned]

#### adjustInventory `PATCH /{productId}`

Adjust stock level (set absolute or delta)

**Input:**
- `productId`: ProductId (UUID) (required)
- `adjustmentType`: enum(set, delta) (required)
- `quantity`: integer (required)
- `reason`: string

**Output:**
- `adjustedBy`: UserId
- `inventoryId`: InventoryId (UUID)
- `previousQuantity`: integer
- `newQuantity`: integer
- `adjustedAt`: DateTime

**Errors:**
| Code | HTTP | Message |
|------|------|----------|
| PRODUCT_NOT_FOUND | 404 | Product with specified ID does not exist |
| NEGATIVE_INVENTORY | 400 | Inventory cannot go below zero |
| INSUFFICIENT_STOCK | 409 | Cannot reduce inventory below reserved quantity |

**Preconditions:** [User must be authenticated as admin Product must exist Resulting quantity must be >= 0]

**Postconditions:** [Inventory adjusted Full audit trail recorded Low stock alert triggered if applicable InventoryAdjusted event emitted]

#### restockInventory `POST /{productId}/restock`

Add stock to inventory

**Input:**
- `quantity`: integer (required)
- `reason`: string
- `productId`: ProductId (UUID) (required)

**Output:**
- `previousQuantity`: integer
- `newQuantity`: integer
- `restockedAt`: DateTime
- `inventoryId`: InventoryId (UUID)

**Errors:**
| Code | HTTP | Message |
|------|------|----------|
| PRODUCT_NOT_FOUND | 404 | Product with specified ID does not exist |
| INVALID_QUANTITY | 400 | Restock quantity must be greater than zero |

**Preconditions:** [User must be authenticated as admin Product must exist Quantity must be > 0]

**Postconditions:** [Inventory increased Audit trail recorded InventoryRestocked event emitted]

#### bulkImportInventory `POST /import`

Bulk update inventory from CSV

**Input:**
- `file`: CSV file (required)

**Output:**
- `failureCount`: integer
- `errors`: ImportError[]
- `successCount`: integer

**Errors:**
| Code | HTTP | Message |
|------|------|----------|
| INVALID_CSV_FORMAT | 400 | CSV file format is invalid |
| PARTIAL_IMPORT_FAILURE | 207 | Some rows failed validation |

**Preconditions:** [User must be authenticated as admin CSV must have valid format]

**Postconditions:** [Valid rows processed Row-by-row errors returned for failures]

### Events

- **InventoryReserved**: Emitted when stock is reserved for an order (payload: [inventoryId productId quantity orderId])
- **InventoryReleased**: Emitted when reserved stock is released (payload: [inventoryId productId quantity reason])
- **InventoryDeducted**: Emitted when stock is permanently reduced (payload: [inventoryId productId quantity orderId])
- **InventoryRestocked**: Emitted when stock is added (payload: [inventoryId productId previousQuantity newQuantity reason])
- **OutOfStock**: Emitted when available quantity reaches zero (payload: [inventoryId productId])

---

