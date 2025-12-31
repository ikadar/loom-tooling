# Interface Contracts

Generated: 2025-12-31T12:00:34+01:00

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
| addressId | UUID | optional, for saved addresses |
| street | string | 1-200 chars |
| city | string | 1-100 chars |
| state | string | required |
| postalCode | string | valid format |
| country | string | ISO 3166-1 |
| recipientName | string | required for shipping |

### PaymentMethod

| Field | Type | Constraints |
|-------|------|-------------|
| paymentMethodId | UUID | optional, for saved methods |
| type | PaymentType | required |
| lastFourDigits | string | 4 digits, for credit card display |
| cardBrand | string | visa, mastercard, amex |
| paypalEmail | string | for PayPal identification |

### Quantity

| Field | Type | Constraints |
|-------|------|-------------|
| value | integer | > 0 |

### Email

| Field | Type | Constraints |
|-------|------|-------------|
| value | string | RFC 5322 valid email format |

### SKU

| Field | Type | Constraints |
|-------|------|-------------|
| value | string | alphanumeric, unique across all variants |

### OrderStatus

| Field | Type | Constraints |
|-------|------|-------------|
| value | enum | pending, confirmed, shipped, delivered, cancelled |

### RegistrationStatus

| Field | Type | Constraints |
|-------|------|-------------|
| value | enum | unverified, registered |

### PaymentType

| Field | Type | Constraints |
|-------|------|-------------|
| value | enum | credit_card, paypal |

### ProductSummary

| Field | Type | Constraints |
|-------|------|-------------|
| productId | UUID | required |
| name | string | required |
| price | Money | required |
| primaryImageUrl | string | optional |
| categoryName | string | required |
| isInStock | boolean | required |
| hasVariants | boolean | required |

### ProductVariant

| Field | Type | Constraints |
|-------|------|-------------|
| variantId | UUID | required |
| sku | SKU | required, unique |
| size | string | optional |
| color | string | optional |
| priceAdjustment | Money | optional, can be negative |
| availableQuantity | integer | >= 0 |

### ProductImage

| Field | Type | Constraints |
|-------|------|-------------|
| imageId | UUID | required |
| url | string | required, valid URL |
| isPrimary | boolean | required |
| displayOrder | integer | >= 0 |

### CartItem

| Field | Type | Constraints |
|-------|------|-------------|
| cartItemId | UUID | required |
| productId | UUID | required |
| variantId | UUID | optional |
| productName | string | required |
| unitPrice | Money | required |
| quantity | Quantity | required |
| subtotal | Money | calculated |
| isInStock | boolean | required |
| availableQuantity | integer | >= 0 |

### OrderLineItem

| Field | Type | Constraints |
|-------|------|-------------|
| lineItemId | UUID | required |
| productId | UUID | required |
| productName | string | required, snapshot |
| unitPrice | Money | required, snapshot |
| quantity | Quantity | required |
| subtotal | Money | calculated |

### OrderSummary

| Field | Type | Constraints |
|-------|------|-------------|
| orderId | UUID | required |
| orderNumber | string | required |
| status | OrderStatus | required |
| totalAmount | Money | required |
| itemCount | integer | required |
| createdAt | DateTime | required |

### AdminOrderSummary

| Field | Type | Constraints |
|-------|------|-------------|
| orderId | UUID | required |
| orderNumber | string | required |
| customerName | string | required |
| customerEmail | Email | required |
| status | OrderStatus | required |
| totalAmount | Money | required |
| createdAt | DateTime | required |

### CategoryTree

| Field | Type | Constraints |
|-------|------|-------------|
| categoryId | UUID | required |
| name | string | required |
| description | string | optional |
| depth | integer | 0-2 |
| children | List<CategoryTree> | optional |

### PriceChange

| Field | Type | Constraints |
|-------|------|-------------|
| cartItemId | UUID | required |
| productName | string | required |
| previousPrice | Money | required |
| currentPrice | Money | required |

### QuantityLimitedItem

| Field | Type | Constraints |
|-------|------|-------------|
| productId | UUID | required |
| productName | string | required |
| requestedQuantity | integer | required |
| limitedToQuantity | integer | required |

### BulkUpdateError

| Field | Type | Constraints |
|-------|------|-------------|
| rowNumber | integer | required |
| sku | string | required |
| errorCode | string | required |
| errorMessage | string | required |

### InventoryAuditEntry

| Field | Type | Constraints |
|-------|------|-------------|
| auditId | UUID | required |
| productId | UUID | required |
| previousValue | integer | required |
| newValue | integer | required |
| reason | string | required |
| userId | UUID | required |
| timestamp | DateTime | required |

### DateRange

| Field | Type | Constraints |
|-------|------|-------------|
| startDate | Date | required |
| endDate | Date | required |

---

## IC-CUSTOMER-001 – Customer Service

**Purpose:** Manages customer registration, authentication, and profile data including shipping addresses

**Base URL:** `/api/v1/customers`

**Security:**
- Authentication: Bearer JWT (except registration and verification)
- Authorization: Customer can only access own data, admin can access all

### Operations

#### registerCustomer `POST /register`

Register a new customer account with email verification

**Input:**
- `firstName`: string (required)
- `lastName`: string (required)
- `email`: Email (required)
- `password`: string (required)

**Output:**
- `customerId`: UUID
- `email`: Email
- `registrationStatus`: RegistrationStatus
- `createdAt`: DateTime

**Errors:**
| Code | HTTP | Message |
|------|------|----------|
| DUPLICATE_EMAIL | 409 | Email already registered, suggest login or password reset |
| INVALID_PASSWORD | 400 | Password must be at least 8 characters with at least one number |
| INVALID_EMAIL | 400 | Invalid email format |

**Preconditions:** [Email not already registered]

**Postconditions:** [Customer account created with status 'unverified' Verification email sent CustomerRegistered event emitted]

#### verifyEmail `POST /verify-email`

Verify customer email address using verification token

**Input:**
- `token`: string (required)

**Output:**
- `customerId`: UUID
- `emailVerified`: boolean
- `registrationStatus`: RegistrationStatus

**Errors:**
| Code | HTTP | Message |
|------|------|----------|
| INVALID_TOKEN | 400 | Invalid verification token |
| VERIFICATION_EXPIRED | 410 | Verification link expired, request new verification email |

**Preconditions:** [Valid verification token exists]

**Postconditions:** [Customer email marked as verified Registration status updated to 'registered']

#### addShippingAddress `POST /{customerId}/addresses`

Add a new shipping address to customer account

**Input:**
- `country`: string (required)
- `recipientName`: string (required)
- `customerId`: UUID (required)
- `street`: string (required)
- `city`: string (required)
- `state`: string (required)
- `postalCode`: string (required)

**Output:**
- `addressId`: UUID
- `address`: Address

**Errors:**
| Code | HTTP | Message |
|------|------|----------|
| INVALID_ADDRESS | 400 | Missing required address field |
| INVALID_POSTAL_CODE | 400 | Invalid postal code format |
| UNSUPPORTED_SHIPPING_REGION | 400 | Shipping only available within domestic country |
| CUSTOMER_NOT_FOUND | 404 | Customer does not exist |

**Preconditions:** [Customer exists and is registered]

**Postconditions:** [Address saved to customer account Address available for order selection]

#### getCustomerAddresses `GET /{customerId}/addresses`

Retrieve all shipping addresses for a customer

**Input:**
- `customerId`: UUID (required)

**Output:**
- `addresses`: List<Address>

**Errors:**
| Code | HTTP | Message |
|------|------|----------|
| CUSTOMER_NOT_FOUND | 404 | Customer does not exist |

**Preconditions:** [Customer exists]

#### requestDataErasure `DELETE /{customerId}/data`

GDPR data erasure request - permanently delete all personal data

**Input:**
- `customerId`: UUID (required)
- `confirmErasure`: boolean (required)

**Output:**
- `erasedAt`: DateTime
- `erasureConfirmed`: boolean

**Errors:**
| Code | HTTP | Message |
|------|------|----------|
| PENDING_ORDERS_EXIST | 409 | Cannot delete account with pending orders |
| CONFIRMATION_REQUIRED | 400 | Must confirm erasure request |

**Preconditions:** [Customer exists No pending orders]

**Postconditions:** [All personal data permanently deleted Order history anonymized Confirmation sent]

#### savePaymentMethod `POST /{customerId}/payment-methods`

Save a tokenized payment method for future purchases

**Input:**
- `cardBrand`: string
- `customerId`: UUID (required)
- `paymentToken`: string (required)
- `type`: PaymentType (required)
- `lastFourDigits`: string

**Output:**
- `paymentMethodId`: UUID
- `type`: PaymentType
- `displayName`: string

**Errors:**
| Code | HTTP | Message |
|------|------|----------|
| INVALID_PAYMENT_TOKEN | 400 | Invalid payment token |
| CUSTOMER_NOT_FOUND | 404 | Customer does not exist |

**Preconditions:** [Customer exists Valid payment token]

**Postconditions:** [Tokenized payment method stored Available for future checkouts]

### Events

- **CustomerRegistered**: Emitted when a new customer registers successfully (payload: [customerId email registrationStatus createdAt])
- **CustomerEmailVerified**: Emitted when customer verifies their email (payload: [customerId email verifiedAt])
- **CustomerDataErased**: Emitted when customer data is permanently deleted (GDPR) (payload: [anonymizedCustomerId erasedAt])

---

## IC-PRODUCT-001 – Product Service

**Purpose:** Manages product catalog including browsing, filtering, and product details with variants

**Base URL:** `/api/v1/products`

**Security:**
- Authentication: None for browsing, Bearer JWT for admin operations
- Authorization: Admin role required for create/update/delete operations

### Operations

#### browseProducts `GET /`

Browse and filter products with pagination

**Input:**
- `inStock`: boolean
- `sortBy`: string
- `sortOrder`: string
- `page`: integer
- `pageSize`: integer
- `categoryId`: UUID
- `minPrice`: decimal
- `maxPrice`: decimal

**Output:**
- `products`: List<ProductSummary>
- `totalCount`: integer
- `page`: integer
- `pageSize`: integer
- `totalPages`: integer

**Errors:**
| Code | HTTP | Message |
|------|------|----------|
| INVALID_CATEGORY | 400 | Category does not exist |
| INVALID_PRICE_RANGE | 400 | minPrice must be less than maxPrice |

#### getProductDetails `GET /{productId}`

Get detailed product information including variants and stock status

**Input:**
- `productId`: UUID (required)

**Output:**
- `images`: List<ProductImage>
- `name`: string
- `description`: string
- `availableQuantity`: integer
- `isActive`: boolean
- `productId`: UUID
- `price`: Money
- `categoryId`: UUID
- `categoryName`: string
- `variants`: List<ProductVariant>

**Errors:**
| Code | HTTP | Message |
|------|------|----------|
| PRODUCT_NOT_FOUND | 404 | Product does not exist |

#### createProduct `POST /`

Create a new product in draft status (admin only)

**Input:**
- `name`: string (required)
- `description`: string
- `price`: Money (required)
- `categoryId`: UUID (required)

**Output:**
- `createdBy`: UUID
- `productId`: UUID
- `name`: string
- `status`: string
- `createdAt`: DateTime

**Errors:**
| Code | HTTP | Message |
|------|------|----------|
| NAME_REQUIRED | 400 | Product name is required |
| NAME_TOO_LONG | 400 | Product name exceeds 200 characters |
| INVALID_PRICE | 400 | Price must be at least $0.01 |
| CATEGORY_REQUIRED | 400 | Category must be selected |
| CATEGORY_NOT_FOUND | 404 | Category does not exist |

**Preconditions:** [Admin user authenticated Category exists]

**Postconditions:** [Product created with status 'draft' UUID assigned Creator and timestamp logged ProductCreated event emitted]

#### updateProduct `PUT /{productId}`

Update product details (admin only)

**Input:**
- `description`: string
- `price`: Money
- `categoryId`: UUID
- `productId`: UUID (required)
- `name`: string

**Output:**
- `updatedBy`: UUID
- `productId`: UUID
- `name`: string
- `price`: Money
- `updatedAt`: DateTime

**Errors:**
| Code | HTTP | Message |
|------|------|----------|
| PRODUCT_NOT_FOUND | 404 | Product does not exist |
| NAME_TOO_LONG | 400 | Product name exceeds 200 characters |
| INVALID_PRICE | 400 | Price must be at least $0.01 |
| CATEGORY_NOT_FOUND | 404 | Category does not exist |

**Preconditions:** [Admin user authenticated Product exists]

**Postconditions:** [Product attributes updated Change logged with admin user and timestamp ProductUpdated event emitted Carts show updated price on next view]

#### deleteProduct `DELETE /{productId}`

Soft-delete a product (admin only)

**Input:**
- `productId`: UUID (required)

**Output:**
- `productId`: UUID
- `deletedAt`: DateTime
- `softDeleted`: boolean

**Errors:**
| Code | HTTP | Message |
|------|------|----------|
| PRODUCT_NOT_FOUND | 404 | Product does not exist |

**Preconditions:** [Admin user authenticated Product exists]

**Postconditions:** [Product soft-deleted (not visible to customers) Order history preserved Cart items removed with notification ProductDeactivated event emitted]

#### uploadProductImages `POST /{productId}/images`

Upload product images (admin only)

**Input:**
- `images`: List<MultipartFile> (required)
- `primaryImageIndex`: integer
- `productId`: UUID (required)

**Output:**
- `primaryImageUrl`: string
- `images`: List<ProductImage>

**Errors:**
| Code | HTTP | Message |
|------|------|----------|
| PRODUCT_NOT_FOUND | 404 | Product does not exist |
| MAX_IMAGES_EXCEEDED | 400 | Maximum 10 images allowed per product |
| INVALID_IMAGE_FORMAT | 400 | Invalid image format, supported: JPG, PNG, WebP |

**Preconditions:** [Admin user authenticated Product exists]

**Postconditions:** [Images uploaded to cloud storage URLs stored Primary image marked]

#### addProductVariant `POST /{productId}/variants`

Add a variant to a product (admin only)

**Input:**
- `productId`: UUID (required)
- `sku`: SKU (required)
- `size`: string
- `color`: string
- `priceAdjustment`: Money

**Output:**
- `variantId`: UUID
- `sku`: SKU
- `size`: string
- `color`: string
- `priceAdjustment`: Money

**Errors:**
| Code | HTTP | Message |
|------|------|----------|
| PRODUCT_NOT_FOUND | 404 | Product does not exist |
| DUPLICATE_SKU | 409 | SKU already exists |
| VARIANT_ATTRIBUTE_REQUIRED | 400 | At least size or color must be specified |
| DUPLICATE_VARIANT | 409 | Variant combination already exists |

**Preconditions:** [Admin user authenticated Product exists SKU is unique]

**Postconditions:** [Variant added to product ProductVariantAdded event emitted]

### Events

- **ProductCreated**: Emitted when a new product is created (payload: [productId name price categoryId createdBy])
- **ProductUpdated**: Emitted when product details are updated (payload: [productId updatedFields updatedBy previousPrice newPrice])
- **ProductDeactivated**: Emitted when a product is soft-deleted (payload: [productId deletedBy deletedAt])
- **ProductVariantAdded**: Emitted when a variant is added to a product (payload: [productId variantId sku])

---

## IC-CART-001 – Cart Service

**Purpose:** Manages shopping cart operations including adding, updating, and removing items

**Base URL:** `/api/v1/carts`

**Security:**
- Authentication: Bearer JWT for logged-in users, session token for guest carts
- Authorization: User can only access own cart

### Operations

#### getCart `GET /{cartId}`

Get current cart contents with calculated totals

**Input:**
- `cartId`: UUID (required)

**Output:**
- `hasOutOfStockItems`: boolean
- `priceChanges`: List<PriceChange>
- `cartId`: UUID
- `customerId`: UUID
- `items`: List<CartItem>
- `totalPrice`: Money
- `itemCount`: integer

**Errors:**
| Code | HTTP | Message |
|------|------|----------|
| CART_NOT_FOUND | 404 | Cart does not exist |
| CART_EXPIRED | 410 | Cart has expired due to inactivity |

**Preconditions:** [Cart exists]

#### addItemToCart `POST /{cartId}/items`

Add a product to the cart or increase quantity if already present

**Input:**
- `quantity`: Quantity (required)
- `cartId`: UUID (required)
- `productId`: UUID (required)
- `variantId`: UUID

**Output:**
- `cartItemId`: UUID
- `productId`: UUID
- `productName`: string
- `unitPrice`: Money
- `quantity`: Quantity
- `subtotal`: Money
- `cartTotal`: Money
- `cartId`: UUID

**Errors:**
| Code | HTTP | Message |
|------|------|----------|
| OUT_OF_STOCK | 409 | Product is out of stock |
| INSUFFICIENT_STOCK | 409 | Requested quantity exceeds available stock |
| PRODUCT_NOT_FOUND | 404 | Product does not exist |
| PRODUCT_INACTIVE | 400 | Product is not available for purchase |
| VARIANT_REQUIRED | 400 | Product has variants, specific variant must be selected |
| QUANTITY_LIMITED | 200 | Quantity limited to available stock |

**Preconditions:** [Product exists and is active Product is in stock Quantity > 0 Variant selected if product has variants]

**Postconditions:** [Item added or quantity updated Total recalculated ItemAddedToCart event emitted]

#### updateCartItemQuantity `PUT /{cartId}/items/{cartItemId}`

Update quantity of an item in the cart

**Input:**
- `cartId`: UUID (required)
- `cartItemId`: UUID (required)
- `quantity`: Quantity (required)

**Output:**
- `quantity`: Quantity
- `subtotal`: Money
- `cartTotal`: Money
- `quantityLimited`: boolean
- `availableQuantity`: integer
- `cartItemId`: UUID

**Errors:**
| Code | HTTP | Message |
|------|------|----------|
| CART_NOT_FOUND | 404 | Cart does not exist |
| CART_ITEM_NOT_FOUND | 404 | Item not found in cart |
| QUANTITY_LIMITED | 200 | Quantity limited to available stock |

**Preconditions:** [Cart exists Item exists in cart]

**Postconditions:** [Quantity updated (or item removed if quantity is 0) Total recalculated CartItemQuantityUpdated event emitted]

#### removeCartItem `DELETE /{cartId}/items/{cartItemId}`

Remove an item from the cart

**Input:**
- `cartId`: UUID (required)
- `cartItemId`: UUID (required)

**Output:**
- `cartTotal`: Money
- `itemCount`: integer
- `undoToken`: string
- `undoExpiresAt`: DateTime
- `cartId`: UUID
- `removedItemId`: UUID

**Errors:**
| Code | HTTP | Message |
|------|------|----------|
| CART_NOT_FOUND | 404 | Cart does not exist |
| CART_ITEM_NOT_FOUND | 404 | Item not found in cart |

**Preconditions:** [Cart exists Item exists in cart]

**Postconditions:** [Item removed from cart Total recalculated Undo option available for 5 seconds ItemRemovedFromCart event emitted]

#### undoRemoveCartItem `POST /{cartId}/items/undo`

Undo the last item removal within 5 seconds

**Input:**
- `cartId`: UUID (required)
- `undoToken`: string (required)

**Output:**
- `cartItemId`: UUID
- `restored`: boolean
- `cartTotal`: Money

**Errors:**
| Code | HTTP | Message |
|------|------|----------|
| UNDO_EXPIRED | 410 | Undo window has expired |
| INVALID_UNDO_TOKEN | 400 | Invalid undo token |

**Preconditions:** [Valid undo token Within 5 second window]

**Postconditions:** [Item restored to cart Total recalculated]

#### mergeGuestCart `POST /{cartId}/merge`

Merge guest cart with authenticated user's cart on login

**Input:**
- `cartId`: UUID (required)
- `guestCartId`: UUID (required)

**Output:**
- `quantityLimitedItems`: List<QuantityLimitedItem>
- `cartTotal`: Money
- `cartId`: UUID
- `mergedItems`: integer

**Errors:**
| Code | HTTP | Message |
|------|------|----------|
| CART_NOT_FOUND | 404 | Cart does not exist |
| GUEST_CART_NOT_FOUND | 404 | Guest cart does not exist |

**Preconditions:** [Both carts exist]

**Postconditions:** [Guest cart items merged into user cart Quantities limited by stock if needed Guest cart cleared]

#### clearCart `DELETE /{cartId}`

Clear all items from the cart

**Input:**
- `cartId`: UUID (required)

**Output:**
- `cleared`: boolean
- `cartId`: UUID

**Errors:**
| Code | HTTP | Message |
|------|------|----------|
| CART_NOT_FOUND | 404 | Cart does not exist |

**Preconditions:** [Cart exists]

**Postconditions:** [All items removed Total set to zero CartCleared event emitted]

### Events

- **ItemAddedToCart**: Emitted when an item is added to the cart (payload: [cartId cartItemId productId quantity unitPrice])
- **CartItemQuantityUpdated**: Emitted when cart item quantity is changed (payload: [cartId cartItemId previousQuantity newQuantity])
- **ItemRemovedFromCart**: Emitted when an item is removed from the cart (payload: [cartId cartItemId productId])
- **CartCleared**: Emitted when cart is cleared (manually or on order placement) (payload: [cartId customerId reason])
- **CartExpired**: Emitted when cart expires due to inactivity (payload: [cartId customerId expiredAt])

---

## IC-ORDER-001 – Order Service

**Purpose:** Manages order lifecycle from creation through delivery including status tracking and cancellation

**Base URL:** `/api/v1/orders`

**Security:**
- Authentication: Bearer JWT required for all operations
- Authorization: Customer can only access own orders, admin operations require admin role

### Operations

#### placeOrder `POST /`

Create a new order from cart contents with payment authorization

**Input:**
- `paymentMethodId`: UUID
- `paymentMethod`: PaymentMethod
- `savePaymentMethod`: boolean
- `customerId`: UUID (required)
- `cartId`: UUID (required)
- `shippingAddressId`: UUID
- `shippingAddress`: Address

**Output:**
- `orderId`: UUID
- `status`: OrderStatus
- `lineItems`: List<OrderLineItem>
- `subtotal`: Money
- `tax`: Money
- `totalAmount`: Money
- `shippingAddress`: Address
- `orderNumber`: string
- `shippingCost`: Money
- `estimatedDelivery`: DateRange
- `createdAt`: DateTime

**Errors:**
| Code | HTTP | Message |
|------|------|----------|
| EMAIL_NOT_VERIFIED | 403 | Email must be verified to place orders |
| CART_EMPTY | 400 | Cart has no items |
| STOCK_UNAVAILABLE | 409 | One or more items out of stock |
| PAYMENT_FAILED | 402 | Payment authorization failed |
| PAYMENT_DECLINED | 402 | Card was declined |
| INVALID_CARD | 400 | Invalid card number |
| UNSUPPORTED_CARD_TYPE | 400 | Only Visa, Mastercard, and American Express accepted |
| INVALID_ADDRESS | 400 | Invalid shipping address |
| UNSUPPORTED_SHIPPING_REGION | 400 | Shipping only available within domestic country |
| CUSTOMER_NOT_FOUND | 404 | Customer does not exist |

**Preconditions:** [Customer is registered and email verified Cart has items All cart items in stock Valid shipping address Valid payment method]

**Postconditions:** [Order created with status 'pending' Payment authorized Inventory reserved Cart cleared Confirmation email queued OrderPlaced event emitted]

#### getOrder `GET /{orderId}`

Get order details including tracking information

**Input:**
- `orderId`: UUID (required)

**Output:**
- `trackingUrl`: string
- `updatedAt`: DateTime
- `orderId`: UUID
- `lineItems`: List<OrderLineItem>
- `tax`: Money
- `shippingAddress`: Address
- `trackingNumber`: string
- `createdAt`: DateTime
- `orderNumber`: string
- `status`: OrderStatus
- `subtotal`: Money
- `shippingCost`: Money
- `totalAmount`: Money

**Errors:**
| Code | HTTP | Message |
|------|------|----------|
| ORDER_NOT_FOUND | 404 | Order does not exist |
| ACCESS_DENIED | 403 | Not authorized to view this order |

**Preconditions:** [Order exists]

#### getCustomerOrders `GET /`

Get paginated order history for a customer

**Input:**
- `customerId`: UUID (required)
- `status`: OrderStatus
- `page`: integer
- `pageSize`: integer

**Output:**
- `orders`: List<OrderSummary>
- `totalCount`: integer
- `page`: integer
- `pageSize`: integer
- `totalPages`: integer

**Errors:**
| Code | HTTP | Message |
|------|------|----------|
| CUSTOMER_NOT_FOUND | 404 | Customer does not exist |

**Preconditions:** [Customer exists]

#### cancelOrder `POST /{orderId}/cancel`

Cancel an order and initiate refund

**Input:**
- `orderId`: UUID (required)
- `reason`: string (required)

**Output:**
- `status`: OrderStatus
- `refundStatus`: string
- `refundAmount`: Money
- `cancelledAt`: DateTime
- `orderId`: UUID

**Errors:**
| Code | HTTP | Message |
|------|------|----------|
| ORDER_NOT_FOUND | 404 | Order does not exist |
| ORDER_ALREADY_SHIPPED | 409 | Cannot cancel order that has been shipped |
| ORDER_NOT_MODIFIABLE | 409 | Order cannot be modified in current status |
| REFUND_FAILED | 500 | Refund initiation failed |

**Preconditions:** [Order exists Order status is 'pending' or 'confirmed']

**Postconditions:** [Status changed to 'cancelled' Inventory restored Refund initiated to original payment method Cancellation email sent OrderCancelled event emitted]

#### confirmOrder `POST /{orderId}/confirm`

Confirm a pending order (admin only)

**Input:**
- `orderId`: UUID (required)

**Output:**
- `orderId`: UUID
- `status`: OrderStatus
- `confirmedAt`: DateTime

**Errors:**
| Code | HTTP | Message |
|------|------|----------|
| ORDER_NOT_FOUND | 404 | Order does not exist |
| INVALID_STATUS_TRANSITION | 409 | Order must be in 'pending' status to confirm |

**Preconditions:** [Admin user authenticated Order exists Order status is 'pending']

**Postconditions:** [Status changed to 'confirmed' OrderConfirmed event emitted]

#### shipOrder `POST /{orderId}/ship`

Mark order as shipped with tracking info (admin only)

**Input:**
- `orderId`: UUID (required)
- `trackingNumber`: string (required)
- `carrier`: string (required)

**Output:**
- `trackingUrl`: string
- `shippedAt`: DateTime
- `orderId`: UUID
- `status`: OrderStatus
- `trackingNumber`: string

**Errors:**
| Code | HTTP | Message |
|------|------|----------|
| ORDER_NOT_FOUND | 404 | Order does not exist |
| INVALID_STATUS_TRANSITION | 409 | Order must be in 'confirmed' status to ship |

**Preconditions:** [Admin user authenticated Order exists Order status is 'confirmed']

**Postconditions:** [Status changed to 'shipped' Tracking number stored Customer notified via email OrderShipped event emitted]

#### deliverOrder `POST /{orderId}/deliver`

Mark order as delivered (admin only)

**Input:**
- `orderId`: UUID (required)

**Output:**
- `deliveredAt`: DateTime
- `orderId`: UUID
- `status`: OrderStatus

**Errors:**
| Code | HTTP | Message |
|------|------|----------|
| ORDER_NOT_FOUND | 404 | Order does not exist |
| INVALID_STATUS_TRANSITION | 409 | Order must be in 'shipped' status to mark as delivered |

**Preconditions:** [Admin user authenticated Order exists Order status is 'shipped']

**Postconditions:** [Status changed to 'delivered' OrderDelivered event emitted]

#### getAllOrders `GET /admin/all`

Get all orders with filters (admin only)

**Input:**
- `pageSize`: integer
- `status`: OrderStatus
- `startDate`: Date
- `endDate`: Date
- `orderNumberSearch`: string
- `page`: integer

**Output:**
- `pageSize`: integer
- `totalPages`: integer
- `orders`: List<AdminOrderSummary>
- `totalCount`: integer
- `page`: integer

**Errors:**
| Code | HTTP | Message |
|------|------|----------|
| INVALID_DATE_RANGE | 400 | Start date must be before end date |

**Preconditions:** [Admin user authenticated]

#### exportOrdersCsv `GET /admin/export`

Export filtered orders to CSV (admin only)

**Input:**
- `status`: OrderStatus
- `startDate`: Date
- `endDate`: Date
- `orderNumberSearch`: string

**Output:**
- `file`: binary
- `filename`: string

**Errors:**
| Code | HTTP | Message |
|------|------|----------|
| EXPORT_TOO_LARGE | 400 | Too many orders to export, narrow filter criteria |

**Preconditions:** [Admin user authenticated]

**Postconditions:** [CSV file generated with all filtered orders]

### Events

- **OrderPlaced**: Emitted when a new order is created (payload: [orderId orderNumber customerId totalAmount lineItems shippingAddress])
- **OrderConfirmed**: Emitted when an order is confirmed (payload: [orderId customerId confirmedAt])
- **OrderShipped**: Emitted when an order is shipped (payload: [orderId customerId trackingNumber carrier shippedAt])
- **OrderDelivered**: Emitted when an order is delivered (payload: [orderId customerId deliveredAt])
- **OrderCancelled**: Emitted when an order is cancelled (payload: [orderId customerId reason refundAmount cancelledAt])

---

## IC-INVENTORY-001 – Inventory Service

**Purpose:** Manages product stock levels, reservations, and inventory adjustments with audit logging

**Base URL:** `/api/v1/inventory`

**Security:**
- Authentication: Bearer JWT for admin operations, internal auth for reserve/release/deduct
- Authorization: Admin role required for manual adjustments and audit access

### Operations

#### getInventory `GET /{productId}`

Get current inventory levels for a product

**Input:**
- `productId`: UUID (required)

**Output:**
- `stockLevel`: integer
- `reservedQuantity`: integer
- `availableQuantity`: integer
- `lowStockThreshold`: integer
- `isLowStock`: boolean
- `inventoryId`: UUID
- `productId`: UUID

**Errors:**
| Code | HTTP | Message |
|------|------|----------|
| PRODUCT_NOT_FOUND | 404 | Product does not exist |

**Preconditions:** [Product exists]

#### setInventoryLevel `PUT /{productId}`

Set absolute inventory level (admin only)

**Input:**
- `quantity`: integer (required)
- `reason`: string (required)
- `productId`: UUID (required)

**Output:**
- `adjustedAt`: DateTime
- `inventoryId`: UUID
- `productId`: UUID
- `previousLevel`: integer
- `newLevel`: integer
- `adjustedBy`: UUID

**Errors:**
| Code | HTTP | Message |
|------|------|----------|
| PRODUCT_NOT_FOUND | 404 | Product does not exist |
| INVALID_QUANTITY | 400 | Quantity cannot be negative |

**Preconditions:** [Admin user authenticated Product exists]

**Postconditions:** [Stock level updated Audit log created with before/after values, reason, user, timestamp Low stock alert triggered if applicable]

#### adjustInventory `POST /{productId}/adjust`

Adjust inventory by delta amount (admin only)

**Input:**
- `reason`: string (required)
- `productId`: UUID (required)
- `delta`: integer (required)

**Output:**
- `delta`: integer
- `adjustedBy`: UUID
- `adjustedAt`: DateTime
- `inventoryId`: UUID
- `productId`: UUID
- `previousLevel`: integer
- `newLevel`: integer

**Errors:**
| Code | HTTP | Message |
|------|------|----------|
| PRODUCT_NOT_FOUND | 404 | Product does not exist |
| INSUFFICIENT_STOCK | 409 | Adjustment would result in negative stock |

**Preconditions:** [Admin user authenticated Product exists Resulting quantity >= 0]

**Postconditions:** [Stock level adjusted by delta Audit log created Low stock alert triggered if applicable]

#### bulkUpdateInventory `POST /bulk`

Bulk update inventory from CSV file (admin only)

**Input:**
- `file`: MultipartFile (required)

**Output:**
- `totalRows`: integer
- `successfulUpdates`: integer
- `failedUpdates`: integer
- `errors`: List<BulkUpdateError>

**Errors:**
| Code | HTTP | Message |
|------|------|----------|
| INVALID_CSV_FORMAT | 400 | CSV file format is invalid |
| EMPTY_FILE | 400 | CSV file is empty |

**Preconditions:** [Admin user authenticated Valid CSV file]

**Postconditions:** [Valid rows updated Audit logs created for each update Error report generated for invalid rows]

#### reserveInventory `POST /{productId}/reserve`

Reserve inventory for a pending order (internal use)

**Input:**
- `productId`: UUID (required)
- `quantity`: integer (required)
- `orderId`: UUID (required)

**Output:**
- `reservationId`: UUID
- `productId`: UUID
- `quantity`: integer
- `availableQuantity`: integer

**Errors:**
| Code | HTTP | Message |
|------|------|----------|
| PRODUCT_NOT_FOUND | 404 | Product does not exist |
| INSUFFICIENT_STOCK | 409 | Insufficient stock to reserve |

**Preconditions:** [Product exists Requested quantity <= available quantity]

**Postconditions:** [Reserved quantity increased InventoryReserved event emitted]

#### releaseInventory `POST /{productId}/release`

Release reserved inventory (e.g., on order cancellation)

**Input:**
- `productId`: UUID (required)
- `quantity`: integer (required)
- `orderId`: UUID (required)

**Output:**
- `productId`: UUID
- `releasedQuantity`: integer
- `availableQuantity`: integer

**Errors:**
| Code | HTTP | Message |
|------|------|----------|
| PRODUCT_NOT_FOUND | 404 | Product does not exist |
| INVALID_RELEASE_QUANTITY | 400 | Release quantity exceeds reserved amount |

**Preconditions:** [Product exists Quantity <= reserved quantity]

**Postconditions:** [Reserved quantity decreased InventoryReleased event emitted]

#### deductInventory `POST /{productId}/deduct`

Deduct reserved inventory on order fulfillment (internal use)

**Input:**
- `productId`: UUID (required)
- `quantity`: integer (required)
- `orderId`: UUID (required)

**Output:**
- `productId`: UUID
- `deductedQuantity`: integer
- `newStockLevel`: integer

**Errors:**
| Code | HTTP | Message |
|------|------|----------|
| PRODUCT_NOT_FOUND | 404 | Product does not exist |
| INVALID_DEDUCT_QUANTITY | 400 | Deduct quantity exceeds reserved amount |

**Preconditions:** [Product exists Quantity <= reserved quantity]

**Postconditions:** [Stock level and reserved quantity decreased InventoryDeducted event emitted OutOfStock event if stock reaches 0]

#### getAuditLog `GET /{productId}/audit`

Get inventory audit log for a product (admin only)

**Input:**
- `productId`: UUID (required)
- `startDate`: Date
- `endDate`: Date
- `page`: integer
- `pageSize`: integer

**Output:**
- `auditEntries`: List<InventoryAuditEntry>
- `totalCount`: integer
- `page`: integer
- `pageSize`: integer

**Errors:**
| Code | HTTP | Message |
|------|------|----------|
| PRODUCT_NOT_FOUND | 404 | Product does not exist |

**Preconditions:** [Admin user authenticated Product exists]

### Events

- **InventoryReserved**: Emitted when inventory is reserved for an order (payload: [productId orderId quantity availableQuantity])
- **InventoryReleased**: Emitted when reserved inventory is released (payload: [productId orderId quantity availableQuantity])
- **InventoryDeducted**: Emitted when inventory is permanently deducted (payload: [productId orderId quantity newStockLevel])
- **InventoryRestocked**: Emitted when inventory is restocked (payload: [productId quantity newStockLevel adjustedBy])
- **OutOfStock**: Emitted when product available quantity reaches zero (payload: [productId productName])
- **LowStockAlert**: Emitted when stock drops below threshold (payload: [productId productName currentStock threshold])

---

## IC-CATEGORY-001 – Category Service

**Purpose:** Manages product category hierarchy for catalog organization

**Base URL:** `/api/v1/categories`

**Security:**
- Authentication: None for reading, Bearer JWT for admin operations
- Authorization: Admin role required for create/update/delete

### Operations

#### getCategories `GET /`

Get all categories in hierarchical structure

**Input:**
- `includeChildren`: boolean
- `parentId`: UUID

**Output:**
- `categories`: List<CategoryTree>

#### createCategory `POST /`

Create a new category (admin only)

**Input:**
- `displayOrder`: integer
- `name`: string (required)
- `description`: string
- `parentCategoryId`: UUID

**Output:**
- `createdAt`: DateTime
- `categoryId`: UUID
- `name`: string
- `parentCategoryId`: UUID
- `depth`: integer

**Errors:**
| Code | HTTP | Message |
|------|------|----------|
| DUPLICATE_CATEGORY_NAME | 409 | Category name already exists |
| PARENT_NOT_FOUND | 404 | Parent category does not exist |
| MAX_NESTING_EXCEEDED | 400 | Category hierarchy cannot exceed 3 levels |

**Preconditions:** [Admin user authenticated Name is unique Nesting depth <= 3]

**Postconditions:** [Category created with parent reference CategoryCreated event emitted]

#### updateCategory `PUT /{categoryId}`

Update category details (admin only)

**Input:**
- `categoryId`: UUID (required)
- `name`: string
- `description`: string
- `displayOrder`: integer

**Output:**
- `categoryId`: UUID
- `name`: string
- `updatedAt`: DateTime

**Errors:**
| Code | HTTP | Message |
|------|------|----------|
| CATEGORY_NOT_FOUND | 404 | Category does not exist |
| DUPLICATE_CATEGORY_NAME | 409 | Category name already exists |

**Preconditions:** [Admin user authenticated Category exists]

**Postconditions:** [Category attributes updated]

#### deleteCategory `DELETE /{categoryId}`

Delete a category (admin only)

**Input:**
- `categoryId`: UUID (required)

**Output:**
- `deleted`: boolean

**Errors:**
| Code | HTTP | Message |
|------|------|----------|
| CATEGORY_NOT_FOUND | 404 | Category does not exist |
| CATEGORY_HAS_PRODUCTS | 409 | Cannot delete category with products |
| CATEGORY_HAS_CHILDREN | 409 | Cannot delete category with subcategories |

**Preconditions:** [Admin user authenticated Category exists No products in category No child categories]

**Postconditions:** [Category deleted]

### Events

- **CategoryCreated**: Emitted when a new category is created (payload: [categoryId name parentCategoryId])

---

## IC-NOTIFICATION-001 – Notification Service

**Purpose:** Handles email notifications for order confirmations, status updates, and alerts

**Base URL:** `/api/v1/notifications`

**Security:**
- Authentication: Internal service authentication
- Authorization: Only internal services can trigger notifications

### Operations

#### sendOrderConfirmation `POST /order-confirmation`

Queue order confirmation email (internal use)

**Input:**
- `orderId`: UUID (required)
- `customerId`: UUID (required)
- `email`: Email (required)

**Output:**
- `status`: string
- `queuedAt`: DateTime
- `notificationId`: UUID

**Errors:**
| Code | HTTP | Message |
|------|------|----------|
| ORDER_NOT_FOUND | 404 | Order does not exist |
| DUPLICATE_NOTIFICATION | 409 | Confirmation email already sent for this order |

**Preconditions:** [Order exists Email not already sent for this order]

**Postconditions:** [Email queued with order details Retry mechanism enabled if delivery fails]

#### sendShippingNotification `POST /shipping-update`

Queue shipping status email (internal use)

**Input:**
- `orderId`: UUID (required)
- `customerId`: UUID (required)
- `email`: Email (required)
- `trackingNumber`: string (required)
- `carrier`: string (required)

**Output:**
- `notificationId`: UUID
- `status`: string
- `queuedAt`: DateTime

**Errors:**
| Code | HTTP | Message |
|------|------|----------|
| ORDER_NOT_FOUND | 404 | Order does not exist |

**Preconditions:** [Order exists]

**Postconditions:** [Email queued with tracking information]

#### sendCancellationNotification `POST /cancellation`

Queue order cancellation email (internal use)

**Input:**
- `orderId`: UUID (required)
- `customerId`: UUID (required)
- `email`: Email (required)
- `refundAmount`: Money (required)

**Output:**
- `notificationId`: UUID
- `status`: string
- `queuedAt`: DateTime

**Errors:**
| Code | HTTP | Message |
|------|------|----------|
| ORDER_NOT_FOUND | 404 | Order does not exist |

**Preconditions:** [Order exists]

**Postconditions:** [Email queued with cancellation and refund details]

#### sendLowStockAlert `POST /low-stock-alert`

Send low stock alert to admin (internal use)

**Input:**
- `productId`: UUID (required)
- `productName`: string (required)
- `currentStock`: integer (required)
- `threshold`: integer (required)

**Output:**
- `status`: string
- `notificationId`: UUID

**Postconditions:** [Alert sent to admin notification channel]

### Events

- **NotificationSent**: Emitted when a notification is successfully sent (payload: [notificationId type recipientEmail sentAt])
- **NotificationFailed**: Emitted when a notification fails to send (payload: [notificationId type recipientEmail error retryCount])

---

