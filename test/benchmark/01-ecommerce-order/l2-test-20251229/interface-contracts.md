# Interface Contracts

Generated: 2025-12-29T15:49:16+01:00

---

## Shared Types

### Money

| Field | Type | Constraints |
|-------|------|-------------|
| amount | decimal | >= 0, precision 2 |
| currency | string | ISO 4217, default USD |

### Address

| Field | Type | Constraints |
|-------|------|-------------|
| addressId | UUID |  |
| street | string | 1-200 chars |
| city | string | 1-100 chars |
| state | string | required |
| postalCode | string | valid format |
| country | string | ISO 3166-1 |
| recipientName | string | 1-200 chars |

### Email

| Field | Type | Constraints |
|-------|------|-------------|
| value | string | RFC 5322 email format |

### Quantity

| Field | Type | Constraints |
|-------|------|-------------|
| value | integer | > 0 |

### SKU

| Field | Type | Constraints |
|-------|------|-------------|
| value | string | alphanumeric, unique |

### OrderStatus

| Field | Type | Constraints |
|-------|------|-------------|
| value | enum | pending | confirmed | shipped | delivered | cancelled |

### PaymentMethod

| Field | Type | Constraints |
|-------|------|-------------|
| type | enum | credit_card | paypal |
| token | string | payment gateway token |
| lastFourDigits | string | for credit card display |
| cardBrand | enum | visa | mastercard | amex |
| paypalEmail | string | for PayPal identification |

### Pagination

| Field | Type | Constraints |
|-------|------|-------------|
| page | integer | >= 1 |
| pageSize | integer | 1-100 |
| totalItems | integer | >= 0 |
| totalPages | integer | >= 0 |

### CartItem

| Field | Type | Constraints |
|-------|------|-------------|
| itemId | UUID |  |
| productId | UUID |  |
| variantId | UUID | nullable |
| productName | string |  |
| variantDescription | string | nullable |
| unitPrice | Money |  |
| quantity | integer | > 0 |
| subtotal | Money |  |
| isInStock | boolean |  |
| availableQuantity | integer |  |

### OrderLineItem

| Field | Type | Constraints |
|-------|------|-------------|
| lineItemId | UUID |  |
| productId | UUID |  |
| variantId | UUID | nullable |
| productName | string | snapshot |
| variantDescription | string | nullable, snapshot |
| unitPrice | Money | snapshot at order time |
| quantity | integer | > 0 |
| subtotal | Money |  |

### ProductSummary

| Field | Type | Constraints |
|-------|------|-------------|
| productId | UUID |  |
| name | string |  |
| price | Money |  |
| primaryImageUrl | string |  |
| categoryName | string |  |
| isInStock | boolean |  |
| hasVariants | boolean |  |

### ProductVariant

| Field | Type | Constraints |
|-------|------|-------------|
| variantId | UUID |  |
| sku | SKU |  |
| size | string | nullable |
| color | string | nullable |
| priceAdjustment | Money |  |
| finalPrice | Money | base + adjustment |
| availableQuantity | integer |  |
| isInStock | boolean |  |

### ProductImage

| Field | Type | Constraints |
|-------|------|-------------|
| imageId | UUID |  |
| url | string |  |
| isPrimary | boolean |  |
| displayOrder | integer |  |

### OrderSummary

| Field | Type | Constraints |
|-------|------|-------------|
| orderId | UUID |  |
| orderNumber | string |  |
| status | OrderStatus |  |
| totalAmount | Money |  |
| itemCount | integer |  |
| createdAt | DateTime |  |

### AdminOrderSummary

| Field | Type | Constraints |
|-------|------|-------------|
| orderId | UUID |  |
| orderNumber | string |  |
| status | OrderStatus |  |
| customerName | string |  |
| customerEmail | string |  |
| totalAmount | Money |  |
| itemCount | integer |  |
| createdAt | DateTime |  |

### Category

| Field | Type | Constraints |
|-------|------|-------------|
| categoryId | UUID |  |
| name | string | unique |
| description | string | nullable |
| parentCategoryId | UUID | nullable |
| depth | integer | 1-3 |
| children | array | nested categories |

### SavedPaymentMethod

| Field | Type | Constraints |
|-------|------|-------------|
| paymentMethodId | UUID |  |
| type | enum | credit_card | paypal |
| lastFourDigits | string | for credit card |
| cardBrand | enum | visa | mastercard | amex |
| paypalEmail | string | for PayPal |
| isDefault | boolean |  |
| expiresAt | string | MM/YY for cards |

### InventoryAuditEntry

| Field | Type | Constraints |
|-------|------|-------------|
| auditLogId | UUID |  |
| productId | UUID |  |
| previousValue | integer |  |
| newValue | integer |  |
| adjustment | integer |  |
| reason | string |  |
| userId | UUID |  |
| userName | string |  |
| timestamp | DateTime |  |

### PriceChange

| Field | Type | Constraints |
|-------|------|-------------|
| itemId | UUID |  |
| productName | string |  |
| oldPrice | Money |  |
| newPrice | Money |  |

### MergeNotification

| Field | Type | Constraints |
|-------|------|-------------|
| productName | string |  |
| message | string |  |
| type | enum | quantity_limited | item_merged |

### BulkUpdateError

| Field | Type | Constraints |
|-------|------|-------------|
| rowNumber | integer |  |
| sku | string |  |
| errorCode | string |  |
| errorMessage | string |  |

---

## IC-CUSTOMER-001 – Customer Service

**Purpose:** Manages customer registration, authentication, addresses, and account lifecycle

**Base URL:** `/api/v1/customers`

**Security:**
- Authentication: Bearer JWT (except registration)
- Authorization: Customer can only access own data

### Operations

#### registerCustomer `POST /register`

Register a new customer account

**Input:**
- `email`: Email (required)
- `password`: string (required)
- `firstName`: string (required)
- `lastName`: string (required)

**Output:**
- `customerId`: UUID
- `email`: Email
- `status`: string
- `createdAt`: DateTime

**Errors:**
| Code | HTTP | Message |
|------|------|----------|
| DUPLICATE_EMAIL | 409 | Email already registered. Please login or reset password |
| INVALID_PASSWORD | 400 | Password must be at least 8 characters with at least one number |
| INVALID_EMAIL | 400 | Invalid email format |

**Preconditions:** [Email not already registered]

**Postconditions:** [Customer account created with status 'unverified' Verification email sent CustomerRegistered event emitted]

#### verifyEmail `POST /verify-email`

Verify customer email address with token

**Input:**
- `token`: string (required)

**Output:**
- `customerId`: UUID
- `emailVerified`: boolean

**Errors:**
| Code | HTTP | Message |
|------|------|----------|
| INVALID_TOKEN | 400 | Invalid verification token |
| VERIFICATION_EXPIRED | 410 | Verification link expired. Please request a new one |

**Preconditions:** [Valid verification token exists]

**Postconditions:** [Customer emailVerified set to true]

#### addShippingAddress `POST /{customerId}/addresses`

Add a new shipping address to customer account

**Input:**
- `state`: string (required)
- `postalCode`: string (required)
- `country`: string (required)
- `recipientName`: string (required)
- `isDefault`: boolean
- `street`: string (required)
- `city`: string (required)

**Output:**
- `city`: string
- `state`: string
- `postalCode`: string
- `country`: string
- `recipientName`: string
- `isDefault`: boolean
- `addressId`: UUID
- `street`: string

**Errors:**
| Code | HTTP | Message |
|------|------|----------|
| INVALID_ADDRESS | 400 | Missing required field: {fieldName} |
| INVALID_POSTAL_CODE | 400 | Invalid postal code format |
| UNSUPPORTED_SHIPPING_REGION | 400 | Shipping not available to this country |
| CUSTOMER_NOT_FOUND | 404 | Customer not found |

**Preconditions:** [Customer exists Country is domestic]

**Postconditions:** [Address saved to customer account]

#### getCustomerAddresses `GET /{customerId}/addresses`

Get all shipping addresses for a customer

**Output:**
- `addresses`: array

**Errors:**
| Code | HTTP | Message |
|------|------|----------|
| CUSTOMER_NOT_FOUND | 404 | Customer not found |

**Preconditions:** [Customer exists]

#### requestDataErasure `POST /{customerId}/gdpr/erase`

Request GDPR data erasure for customer account

**Input:**
- `confirmDeletion`: boolean (required)

**Output:**
- `status`: string
- `erasedAt`: DateTime

**Errors:**
| Code | HTTP | Message |
|------|------|----------|
| PENDING_ORDERS_EXIST | 409 | Cannot delete account with pending orders |
| CUSTOMER_NOT_FOUND | 404 | Customer not found |

**Preconditions:** [Customer exists No pending orders]

**Postconditions:** [Personal data permanently deleted Order history anonymized Confirmation email sent]

### Events

- **CustomerRegistered**: Emitted when a new customer registers (payload: [customerId email createdAt])
- **CustomerEmailVerified**: Emitted when customer verifies their email (payload: [customerId email verifiedAt])
- **CustomerDataErased**: Emitted when customer data is erased per GDPR (payload: [anonymizedId erasedAt])

---

## IC-PRODUCT-001 – Product Service

**Purpose:** Manages product catalog, variants, and browsing functionality

**Base URL:** `/api/v1/products`

**Security:**
- Authentication: Bearer JWT (admin operations only)
- Authorization: Admin role required for create/update/delete

### Operations

#### listProducts `GET /`

Browse products with filters, sorting, and pagination

**Input:**
- `pageSize`: integer
- `categoryId`: UUID
- `minPrice`: Money
- `maxPrice`: Money
- `availability`: enum
- `sortBy`: enum
- `page`: integer

**Output:**
- `products`: array
- `pagination`: Pagination

**Errors:**
| Code | HTTP | Message |
|------|------|----------|
| INVALID_CATEGORY | 400 | Category not found |

#### getProduct `GET /{productId}`

Get product details including variants and stock status

**Output:**
- `name`: string
- `price`: Money
- `images`: array
- `hasVariants`: boolean
- `productId`: UUID
- `description`: string
- `categoryId`: UUID
- `categoryName`: string
- `variants`: array
- `availableQuantity`: integer
- `isInStock`: boolean

**Errors:**
| Code | HTTP | Message |
|------|------|----------|
| PRODUCT_NOT_FOUND | 404 | Product not found |

**Preconditions:** [Product exists and is active]

#### createProduct `POST /`

Create a new product (admin only)

**Input:**
- `name`: string (required)
- `description`: string
- `price`: Money (required)
- `categoryId`: UUID (required)

**Output:**
- `status`: string
- `createdAt`: DateTime
- `createdBy`: UUID
- `productId`: UUID
- `name`: string

**Errors:**
| Code | HTTP | Message |
|------|------|----------|
| NAME_REQUIRED | 400 | Product name is required |
| NAME_TOO_LONG | 400 | Product name cannot exceed 200 characters |
| INVALID_PRICE | 400 | Price must be at least $0.01 |
| CATEGORY_REQUIRED | 400 | Category is required |
| CATEGORY_NOT_FOUND | 404 | Category not found |

**Preconditions:** [Admin authenticated Category exists]

**Postconditions:** [Product created with status 'draft' ProductCreated event emitted]

#### updateProduct `PUT /{productId}`

Update product details (admin only)

**Input:**
- `name`: string
- `description`: string
- `price`: Money
- `categoryId`: UUID

**Output:**
- `productId`: UUID
- `name`: string
- `price`: Money
- `updatedAt`: DateTime
- `updatedBy`: UUID

**Errors:**
| Code | HTTP | Message |
|------|------|----------|
| PRODUCT_NOT_FOUND | 404 | Product not found |
| NAME_TOO_LONG | 400 | Product name cannot exceed 200 characters |
| INVALID_PRICE | 400 | Price must be at least $0.01 |

**Preconditions:** [Admin authenticated Product exists]

**Postconditions:** [Product updated Change logged ProductUpdated event emitted]

#### deleteProduct `DELETE /{productId}`

Soft delete a product (admin only)

**Output:**
- `productId`: UUID
- `status`: string
- `deletedAt`: DateTime

**Errors:**
| Code | HTTP | Message |
|------|------|----------|
| PRODUCT_NOT_FOUND | 404 | Product not found |

**Preconditions:** [Admin authenticated Product exists]

**Postconditions:** [Product soft-deleted Cart items with product removed ProductDeactivated event emitted]

#### uploadProductImages `POST /{productId}/images`

Upload images for a product (admin only)

**Input:**
- `images`: array (required)
- `primaryIndex`: integer

**Output:**
- `primaryImageUrl`: string
- `images`: array

**Errors:**
| Code | HTTP | Message |
|------|------|----------|
| PRODUCT_NOT_FOUND | 404 | Product not found |
| MAX_IMAGES_EXCEEDED | 400 | Maximum 10 images allowed |
| INVALID_IMAGE_FORMAT | 400 | Invalid image format. Allowed: JPG, PNG, WebP |

**Preconditions:** [Admin authenticated Product exists]

**Postconditions:** [Images uploaded to cloud storage URLs stored on product]

#### addProductVariant `POST /{productId}/variants`

Add a variant to a product (admin only)

**Input:**
- `size`: string
- `color`: string
- `priceAdjustment`: Money
- `sku`: SKU (required)

**Output:**
- `variantId`: UUID
- `sku`: SKU
- `size`: string
- `color`: string
- `priceAdjustment`: Money

**Errors:**
| Code | HTTP | Message |
|------|------|----------|
| PRODUCT_NOT_FOUND | 404 | Product not found |
| DUPLICATE_SKU | 409 | SKU already exists |
| VARIANT_REQUIRED | 400 | At least size or color must be specified |

**Preconditions:** [Admin authenticated Product exists SKU is unique]

**Postconditions:** [Variant added to product ProductVariantAdded event emitted]

### Events

- **ProductCreated**: Emitted when a new product is created (payload: [productId name price categoryId createdBy])
- **ProductUpdated**: Emitted when a product is updated (payload: [productId changedFields updatedBy])
- **ProductDeactivated**: Emitted when a product is soft-deleted (payload: [productId deletedBy])
- **ProductVariantAdded**: Emitted when a variant is added to a product (payload: [productId variantId sku])

---

## IC-CART-001 – Cart Service

**Purpose:** Manages shopping cart operations including add, update, remove items

**Base URL:** `/api/v1/cart`

**Security:**
- Authentication: Bearer JWT or guest session token
- Authorization: Customer can only access own cart

### Operations

#### getCart `GET /`

Get current customer's cart with all items

**Output:**
- `hasOutOfStockItems`: boolean
- `priceChanges`: array
- `cartId`: UUID
- `items`: array
- `subtotal`: Money
- `itemCount`: integer

**Errors:**
| Code | HTTP | Message |
|------|------|----------|
| CART_EXPIRED | 410 | Your cart has expired |

**Postconditions:** [Price changes detected and returned if any]

#### addToCart `POST /items`

Add a product to the cart

**Input:**
- `productId`: UUID (required)
- `variantId`: UUID
- `quantity`: Quantity (required)

**Output:**
- `cartId`: UUID
- `item`: CartItem
- `subtotal`: Money
- `quantityLimited`: boolean
- `availableQuantity`: integer

**Errors:**
| Code | HTTP | Message |
|------|------|----------|
| PRODUCT_NOT_FOUND | 404 | Product not found |
| OUT_OF_STOCK | 400 | Product is out of stock |
| INSUFFICIENT_STOCK | 400 | Only {available} units available |
| VARIANT_REQUIRED | 400 | Please select a variant |

**Preconditions:** [Product exists and is active Product in stock Variant selected if product has variants]

**Postconditions:** [Item added or quantity updated Cart total recalculated ItemAddedToCart event emitted]

#### updateCartItemQuantity `PUT /items/{itemId}`

Update quantity of an item in the cart

**Input:**
- `quantity`: Quantity (required)

**Output:**
- `itemRemoved`: boolean
- `cartId`: UUID
- `item`: CartItem
- `subtotal`: Money
- `quantityLimited`: boolean
- `availableQuantity`: integer

**Errors:**
| Code | HTTP | Message |
|------|------|----------|
| ITEM_NOT_FOUND | 404 | Item not in cart |
| INSUFFICIENT_STOCK | 400 | Only {available} units available |

**Preconditions:** [Item exists in cart]

**Postconditions:** [Quantity updated or item removed if 0 Cart total recalculated CartItemQuantityUpdated event emitted]

#### removeFromCart `DELETE /items/{itemId}`

Remove an item from the cart

**Output:**
- `cartId`: UUID
- `subtotal`: Money
- `itemCount`: integer
- `isEmpty`: boolean
- `undoToken`: string

**Errors:**
| Code | HTTP | Message |
|------|------|----------|
| ITEM_NOT_FOUND | 404 | Item not in cart |

**Preconditions:** [Item exists in cart]

**Postconditions:** [Item removed from cart Cart total recalculated ItemRemovedFromCart event emitted]

#### undoRemove `POST /items/undo`

Undo recent item removal

**Input:**
- `undoToken`: string (required)

**Output:**
- `cartId`: UUID
- `item`: CartItem
- `subtotal`: Money

**Errors:**
| Code | HTTP | Message |
|------|------|----------|
| UNDO_EXPIRED | 410 | Undo window has expired |
| INVALID_TOKEN | 400 | Invalid undo token |

**Preconditions:** [Undo token valid and not expired]

**Postconditions:** [Item restored to cart]

#### mergeGuestCart `POST /merge`

Merge guest cart into authenticated customer's cart

**Input:**
- `guestCartId`: UUID (required)

**Output:**
- `cartId`: UUID
- `items`: array
- `subtotal`: Money
- `mergeNotifications`: array

**Errors:**
| Code | HTTP | Message |
|------|------|----------|
| GUEST_CART_NOT_FOUND | 404 | Guest cart not found |

**Preconditions:** [Customer authenticated Guest cart exists]

**Postconditions:** [Guest cart items merged Quantities limited to stock if exceeded]

### Events

- **ItemAddedToCart**: Emitted when an item is added to cart (payload: [cartId productId variantId quantity])
- **CartItemQuantityUpdated**: Emitted when cart item quantity changes (payload: [cartId itemId oldQuantity newQuantity])
- **ItemRemovedFromCart**: Emitted when an item is removed from cart (payload: [cartId productId variantId])
- **CartCleared**: Emitted when cart is cleared (order placed or expired) (payload: [cartId reason])

---

## IC-ORDER-001 – Order Service

**Purpose:** Manages order lifecycle from placement through delivery or cancellation

**Base URL:** `/api/v1/orders`

**Security:**
- Authentication: Bearer JWT
- Authorization: Customer can only access own orders, admin for all operations

### Operations

#### createOrder `POST /`

Create a new order from cart contents

**Input:**
- `savePaymentMethod`: boolean
- `shippingAddressId`: UUID (required)
- `paymentMethod`: PaymentMethod (required)

**Output:**
- `orderId`: UUID
- `orderNumber`: string
- `status`: OrderStatus
- `subtotal`: Money
- `shippingCost`: Money
- `tax`: Money
- `createdAt`: DateTime
- `lineItems`: array
- `totalAmount`: Money

**Errors:**
| Code | HTTP | Message |
|------|------|----------|
| CART_EMPTY | 400 | Cart has no items |
| EMAIL_NOT_VERIFIED | 403 | Please verify your email before ordering |
| STOCK_UNAVAILABLE | 409 | Item {productName} is out of stock |
| PAYMENT_FAILED | 402 | Payment authorization failed |
| PAYMENT_DECLINED | 402 | Card declined: {reason} |
| INVALID_CARD | 400 | Invalid card number |
| UNSUPPORTED_CARD_TYPE | 400 | Card type not accepted. Use Visa, Mastercard, or Amex |
| INVALID_ADDRESS | 400 | Invalid shipping address |
| ADDRESS_NOT_FOUND | 404 | Shipping address not found |

**Preconditions:** [Customer registered and email verified Cart has items All items in stock Valid shipping address Valid payment method]

**Postconditions:** [Order created with status 'pending' Payment authorized Inventory decremented Cart cleared Confirmation email queued OrderPlaced event emitted]

#### getOrder `GET /{orderId}`

Get order details including tracking information

**Output:**
- `shippingAddress`: Address
- `subtotal`: Money
- `tax`: Money
- `totalAmount`: Money
- `shippingCost`: Money
- `trackingNumber`: string
- `trackingUrl`: string
- `createdAt`: DateTime
- `updatedAt`: DateTime
- `orderId`: UUID
- `orderNumber`: string
- `status`: OrderStatus
- `lineItems`: array

**Errors:**
| Code | HTTP | Message |
|------|------|----------|
| ORDER_NOT_FOUND | 404 | Order not found |

**Preconditions:** [Order exists Customer owns order or is admin]

#### listCustomerOrders `GET /`

List customer's orders with pagination

**Input:**
- `status`: OrderStatus
- `page`: integer
- `pageSize`: integer

**Output:**
- `orders`: array
- `pagination`: Pagination

**Preconditions:** [Customer authenticated]

#### cancelOrder `POST /{orderId}/cancel`

Cancel an order (customer or admin)

**Input:**
- `reason`: string (required)

**Output:**
- `status`: OrderStatus
- `refundStatus`: string
- `cancelledAt`: DateTime
- `orderId`: UUID

**Errors:**
| Code | HTTP | Message |
|------|------|----------|
| ORDER_NOT_FOUND | 404 | Order not found |
| ORDER_ALREADY_SHIPPED | 400 | Cannot cancel shipped order |
| ORDER_NOT_MODIFIABLE | 400 | Order cannot be modified |
| REFUND_FAILED | 500 | Refund processing failed |

**Preconditions:** [Order exists Status is pending or confirmed]

**Postconditions:** [Status set to cancelled Inventory restored Refund initiated Cancellation email sent OrderCancelled event emitted]

#### listAllOrders `GET /admin`

List all orders with filters (admin only)

**Input:**
- `status`: OrderStatus
- `startDate`: Date
- `endDate`: Date
- `search`: string
- `page`: integer
- `pageSize`: integer

**Output:**
- `pagination`: Pagination
- `orders`: array

**Preconditions:** [Admin authenticated]

#### exportOrders `GET /admin/export`

Export filtered orders to CSV (admin only)

**Input:**
- `status`: OrderStatus
- `startDate`: Date
- `endDate`: Date
- `search`: string

**Output:**
- `contentType`: string
- `filename`: string
- `data`: binary

**Preconditions:** [Admin authenticated]

#### updateOrderStatus `PUT /{orderId}/status`

Update order status (admin only)

**Input:**
- `status`: OrderStatus (required)
- `trackingNumber`: string

**Output:**
- `updatedAt`: DateTime
- `orderId`: UUID
- `status`: OrderStatus
- `trackingNumber`: string

**Errors:**
| Code | HTTP | Message |
|------|------|----------|
| ORDER_NOT_FOUND | 404 | Order not found |
| INVALID_STATUS_TRANSITION | 400 | Cannot transition from {current} to {new} |

**Preconditions:** [Admin authenticated Order exists Valid status transition]

**Postconditions:** [Status updated Change logged Customer email sent if shipped OrderConfirmed/OrderShipped/OrderDelivered event emitted]

### Events

- **OrderPlaced**: Emitted when a new order is placed (payload: [orderId orderNumber customerId totalAmount lineItems])
- **OrderConfirmed**: Emitted when order is confirmed (payload: [orderId orderNumber confirmedAt])
- **OrderShipped**: Emitted when order is shipped (payload: [orderId orderNumber trackingNumber shippedAt])
- **OrderDelivered**: Emitted when order is delivered (payload: [orderId orderNumber deliveredAt])
- **OrderCancelled**: Emitted when order is cancelled (payload: [orderId orderNumber reason cancelledAt])

---

## IC-INVENTORY-001 – Inventory Service

**Purpose:** Manages product stock levels, reservations, and audit trails

**Base URL:** `/api/v1/inventory`

**Security:**
- Authentication: Bearer JWT
- Authorization: Admin role required for all operations

### Operations

#### getInventory `GET /{productId}`

Get inventory status for a product

**Output:**
- `stockLevel`: integer
- `reservedQuantity`: integer
- `availableQuantity`: integer
- `lowStockThreshold`: integer
- `isLowStock`: boolean
- `productId`: UUID

**Errors:**
| Code | HTTP | Message |
|------|------|----------|
| PRODUCT_NOT_FOUND | 404 | Product not found |

**Preconditions:** [Admin authenticated Product exists]

#### setInventory `PUT /{productId}`

Set absolute inventory quantity (admin only)

**Input:**
- `quantity`: integer (required)
- `reason`: string (required)

**Output:**
- `productId`: UUID
- `previousQuantity`: integer
- `newQuantity`: integer
- `auditLogId`: UUID

**Errors:**
| Code | HTTP | Message |
|------|------|----------|
| PRODUCT_NOT_FOUND | 404 | Product not found |
| INVALID_QUANTITY | 400 | Quantity cannot be negative |

**Preconditions:** [Admin authenticated Product exists]

**Postconditions:** [Stock level updated Audit log created Low stock alert if below threshold]

#### adjustInventory `POST /{productId}/adjust`

Adjust inventory by delta (admin only)

**Input:**
- `adjustment`: integer (required)
- `reason`: string (required)

**Output:**
- `newQuantity`: integer
- `auditLogId`: UUID
- `productId`: UUID
- `previousQuantity`: integer
- `adjustment`: integer

**Errors:**
| Code | HTTP | Message |
|------|------|----------|
| PRODUCT_NOT_FOUND | 404 | Product not found |
| INSUFFICIENT_STOCK | 400 | Adjustment would result in negative stock |

**Preconditions:** [Admin authenticated Product exists Result >= 0]

**Postconditions:** [Stock level adjusted Audit log created Low stock alert if below threshold]

#### bulkUpdateInventory `POST /bulk`

Bulk update inventory from CSV (admin only)

**Input:**
- `file`: File (required)

**Output:**
- `totalRows`: integer
- `successCount`: integer
- `errorCount`: integer
- `errors`: array

**Errors:**
| Code | HTTP | Message |
|------|------|----------|
| INVALID_FILE_FORMAT | 400 | Invalid CSV format |

**Preconditions:** [Admin authenticated Valid CSV file]

**Postconditions:** [Valid rows updated Audit logs created Errors reported]

#### getAuditLog `GET /{productId}/audit`

Get inventory audit log for a product (admin only)

**Input:**
- `startDate`: Date
- `endDate`: Date
- `page`: integer
- `pageSize`: integer

**Output:**
- `productId`: UUID
- `auditEntries`: array
- `pagination`: Pagination

**Errors:**
| Code | HTTP | Message |
|------|------|----------|
| PRODUCT_NOT_FOUND | 404 | Product not found |

**Preconditions:** [Admin authenticated Product exists]

### Events

- **InventoryReserved**: Emitted when inventory is reserved for checkout (payload: [productId quantity orderId])
- **InventoryReleased**: Emitted when reserved inventory is released (payload: [productId quantity reason])
- **InventoryDeducted**: Emitted when inventory is deducted after order (payload: [productId quantity orderId])
- **InventoryRestocked**: Emitted when inventory is restocked (payload: [productId quantity reason userId])
- **LowStockAlert**: Emitted when stock falls below threshold (payload: [productId currentStock threshold])
- **OutOfStock**: Emitted when product becomes out of stock (payload: [productId])

---

## IC-CATEGORY-001 – Category Service

**Purpose:** Manages product category hierarchy

**Base URL:** `/api/v1/categories`

**Security:**
- Authentication: Bearer JWT (admin operations only)
- Authorization: Admin role required for create/update/delete

### Operations

#### listCategories `GET /`

List all categories in hierarchical structure

**Input:**
- `flat`: boolean

**Output:**
- `categories`: array

#### createCategory `POST /`

Create a new category (admin only)

**Input:**
- `name`: string (required)
- `description`: string
- `parentCategoryId`: UUID
- `displayOrder`: integer

**Output:**
- `categoryId`: UUID
- `name`: string
- `parentCategoryId`: UUID
- `depth`: integer

**Errors:**
| Code | HTTP | Message |
|------|------|----------|
| DUPLICATE_CATEGORY_NAME | 409 | Category name already exists |
| PARENT_NOT_FOUND | 404 | Parent category not found |
| MAX_NESTING_EXCEEDED | 400 | Category nesting limited to 3 levels |

**Preconditions:** [Admin authenticated Name is unique Nesting depth <= 3]

**Postconditions:** [Category created CategoryCreated event emitted]

### Events

- **CategoryCreated**: Emitted when a new category is created (payload: [categoryId name parentCategoryId])

---

## IC-PAYMENT-001 – Payment Service

**Purpose:** Handles payment processing, authorization, and saved payment methods

**Base URL:** `/api/v1/payments`

**Security:**
- Authentication: Bearer JWT
- Authorization: Customer can only access own payment methods

### Operations

#### getSavedPaymentMethods `GET /methods`

Get customer's saved payment methods

**Output:**
- `paymentMethods`: array

**Preconditions:** [Customer authenticated]

#### savePaymentMethod `POST /methods`

Save a new payment method for future use

**Input:**
- `type`: enum (required)
- `token`: string (required)
- `isDefault`: boolean

**Output:**
- `lastFourDigits`: string
- `isDefault`: boolean
- `paymentMethodId`: UUID
- `type`: string

**Errors:**
| Code | HTTP | Message |
|------|------|----------|
| INVALID_TOKEN | 400 | Invalid payment token |
| UNSUPPORTED_CARD_TYPE | 400 | Card type not accepted |

**Preconditions:** [Customer authenticated Valid payment token]

**Postconditions:** [Payment method tokenized and saved]

#### deletePaymentMethod `DELETE /methods/{paymentMethodId}`

Delete a saved payment method

**Output:**
- `deleted`: boolean

**Errors:**
| Code | HTTP | Message |
|------|------|----------|
| PAYMENT_METHOD_NOT_FOUND | 404 | Payment method not found |

**Preconditions:** [Customer authenticated Payment method exists]

**Postconditions:** [Payment method deleted]

#### initiatePayPalAuth `POST /paypal/initiate`

Initiate PayPal authorization flow

**Input:**
- `amount`: Money (required)
- `returnUrl`: string (required)
- `cancelUrl`: string (required)

**Output:**
- `authorizationUrl`: string
- `paypalOrderId`: string

**Errors:**
| Code | HTTP | Message |
|------|------|----------|
| PAYPAL_ERROR | 500 | PayPal service error |

**Preconditions:** [Customer authenticated]

**Postconditions:** [PayPal authorization initiated]

#### completePayPalAuth `POST /paypal/complete`

Complete PayPal authorization after redirect

**Input:**
- `paypalOrderId`: string (required)

**Output:**
- `status`: string
- `paymentToken`: string

**Errors:**
| Code | HTTP | Message |
|------|------|----------|
| PAYMENT_CANCELLED | 400 | PayPal authorization cancelled |
| PAYPAL_ERROR | 500 | PayPal service error |

**Preconditions:** [Customer authenticated PayPal authorization pending]

**Postconditions:** [PayPal authorization completed or cancelled]

### Events

- **PaymentAuthorized**: Emitted when payment is authorized (payload: [orderId amount paymentMethod authorizationId])
- **PaymentCaptured**: Emitted when payment is captured (payload: [orderId amount captureId])
- **RefundInitiated**: Emitted when refund is initiated (payload: [orderId amount refundId])

---

