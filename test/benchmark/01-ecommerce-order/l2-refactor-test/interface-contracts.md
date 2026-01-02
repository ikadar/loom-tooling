# Interface Contracts

Generated: 2026-01-02T21:09:41+01:00

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
| paymentMethodId | PaymentMethodId | UUID |
| type | PaymentType | CREDIT_CARD or PAYPAL |
| lastFourDigits | string | 4 digits, required for CREDIT_CARD |
| expiryMonth | integer | 1-12, required for CREDIT_CARD |
| expiryYear | integer | Required for CREDIT_CARD |
| paypalEmail | Email | Required for PAYPAL |
| isDefault | boolean |  |

### Email

| Field | Type | Constraints |
|-------|------|-------------|
| value | string | Valid email format per RFC 5322 |

### OrderStatus

| Field | Type | Constraints |
|-------|------|-------------|
| value | enum | PENDING, CONFIRMED, SHIPPED, DELIVERED, CANCELLED |

### RegistrationStatus

| Field | Type | Constraints |
|-------|------|-------------|
| value | enum | REGISTERED, UNREGISTERED |

### PaymentType

| Field | Type | Constraints |
|-------|------|-------------|
| value | enum | CREDIT_CARD, PAYPAL |

### ProductStatus

| Field | Type | Constraints |
|-------|------|-------------|
| value | enum | DRAFT, ACTIVE, INACTIVE |

### StockStatus

| Field | Type | Constraints |
|-------|------|-------------|
| value | enum | IN_STOCK, LOW_STOCK, OUT_OF_STOCK |

### SortOption

| Field | Type | Constraints |
|-------|------|-------------|
| value | enum | NAME_ASC, NAME_DESC, PRICE_ASC, PRICE_DESC, NEWEST |

### DeleteType

| Field | Type | Constraints |
|-------|------|-------------|
| value | enum | SOFT_DELETE, FULL_ERASURE |

### CartItem

| Field | Type | Constraints |
|-------|------|-------------|
| cartItemId | CartItemId | UUID |
| productId | ProductId | UUID |
| variantId | VariantId | UUID, optional |
| productName | string | Snapshot at time of add |
| unitPrice | Money | Current price |
| quantity | integer | >= 1 |
| subtotal | Money | unitPrice * quantity |

### OrderLineItem

| Field | Type | Constraints |
|-------|------|-------------|
| lineItemId | LineItemId | UUID |
| productId | ProductId | UUID |
| productName | string | Snapshot at order time |
| unitPrice | Money | Snapshot at order time |
| quantity | integer | >= 1 |
| subtotal | Money | unitPrice * quantity |

### ProductSummary

| Field | Type | Constraints |
|-------|------|-------------|
| productId | ProductId | UUID |
| name | string |  |
| price | Money |  |
| categoryId | CategoryId |  |
| primaryImageUrl | string |  |
| stockStatus | StockStatus |  |

### OrderSummary

| Field | Type | Constraints |
|-------|------|-------------|
| orderId | OrderId | UUID |
| orderNumber | string |  |
| customerId | CustomerId |  |
| status | OrderStatus |  |
| totalAmount | Money |  |
| createdAt | DateTime |  |

---

## IC-CUST-001 – Customer Service {#ic-cust-001}

**Purpose:** Manages customer registration, authentication, addresses, and payment methods

**Base URL:** `/api/v1/customers`

**Security:**
- Authentication: Bearer JWT
- Authorization: Customer can only access own profile; Admin can access all

### Operations

#### registerCustomer `POST /`

Register a new customer account with email and password

**Input:**
- `email`: Email (required)
- `password`: string (required)
- `firstName`: string (required)
- `lastName`: string (required)

**Output:**
- `customerId`: CustomerId
- `email`: Email
- `registrationStatus`: RegistrationStatus
- `createdAt`: DateTime

**Errors:**
| Code | HTTP | Message |
|------|------|----------|
| EMAIL_ALREADY_REGISTERED | 409 | Email address is already registered |
| WEAK_PASSWORD | 400 | Password must be at least 8 characters and contain at least one number |
| INVALID_EMAIL | 400 | Email format is invalid |
| MISSING_REQUIRED_FIELDS | 400 | First name, last name, email, and password are required |

**Preconditions:** [Email not already registered]

**Postconditions:** [Customer account created with REGISTERED status Verification email sent CustomerRegistered event emitted]

#### verifyEmail `POST /{customerId}/verify-email`

Verify customer email address using verification token

**Input:**
- `verificationToken`: string (required)
- `customerId`: CustomerId (required)

**Output:**
- `customerId`: CustomerId
- `emailVerified`: boolean

**Errors:**
| Code | HTTP | Message |
|------|------|----------|
| CUSTOMER_NOT_FOUND | 404 | Customer not found |
| INVALID_VERIFICATION_TOKEN | 400 | Verification token is invalid or expired |
| EMAIL_ALREADY_VERIFIED | 409 | Email is already verified |

**Preconditions:** [Customer exists Token is valid and not expired]

**Postconditions:** [Customer emailVerified set to true EmailVerified event emitted]

#### getCustomer `GET /{customerId}`

Retrieve customer profile information

**Input:**
- `customerId`: CustomerId (required)

**Output:**
- `emailVerified`: boolean
- `registrationStatus`: RegistrationStatus
- `createdAt`: DateTime
- `customerId`: CustomerId
- `email`: Email
- `firstName`: string
- `lastName`: string

**Errors:**
| Code | HTTP | Message |
|------|------|----------|
| CUSTOMER_NOT_FOUND | 404 | Customer not found |

**Preconditions:** [Customer exists]

#### addShippingAddress `POST /{customerId}/addresses`

Add a new shipping address to customer profile

**Input:**
- `recipientName`: string (required)
- `isDefault`: boolean
- `customerId`: CustomerId (required)
- `street`: string (required)
- `city`: string (required)
- `state`: string (required)
- `postalCode`: string (required)
- `country`: string (required)

**Output:**
- `addressId`: AddressId
- `street`: string
- `city`: string
- `state`: string
- `postalCode`: string
- `country`: string
- `recipientName`: string
- `isDefault`: boolean

**Errors:**
| Code | HTTP | Message |
|------|------|----------|
| CUSTOMER_NOT_FOUND | 404 | Customer not found |
| INCOMPLETE_ADDRESS | 400 | All address fields are required |
| INVALID_POSTAL_CODE | 400 | Invalid postal code format |
| INTERNATIONAL_SHIPPING_NOT_AVAILABLE | 400 | Shipping is only available within domestic country |

**Preconditions:** [Customer exists]

**Postconditions:** [Address saved to customer profile If isDefault=true, previous default unset]

#### getShippingAddresses `GET /{customerId}/addresses`

List all shipping addresses for a customer

**Input:**
- `customerId`: CustomerId (required)

**Output:**
- `addresses`: Array<ShippingAddress>

**Errors:**
| Code | HTTP | Message |
|------|------|----------|
| CUSTOMER_NOT_FOUND | 404 | Customer not found |

**Preconditions:** [Customer exists]

#### addPaymentMethod `POST /{customerId}/payment-methods`

Add a tokenized payment method to customer profile

**Input:**
- `isDefault`: boolean
- `customerId`: CustomerId (required)
- `type`: PaymentType (required)
- `paymentToken`: string (required)

**Output:**
- `paymentMethodId`: PaymentMethodId
- `type`: PaymentType
- `lastFourDigits`: string
- `isDefault`: boolean

**Errors:**
| Code | HTTP | Message |
|------|------|----------|
| CUSTOMER_NOT_FOUND | 404 | Customer not found |
| INVALID_PAYMENT_TOKEN | 400 | Payment token is invalid |
| UNSUPPORTED_CARD_TYPE | 400 | Card type not supported. Supported: Visa, Mastercard, American Express |

**Preconditions:** [Customer exists Payment token valid from gateway]

**Postconditions:** [Tokenized payment method saved If isDefault=true, previous default unset]

#### getPaymentMethods `GET /{customerId}/payment-methods`

List all saved payment methods for a customer

**Input:**
- `customerId`: CustomerId (required)

**Output:**
- `paymentMethods`: Array<PaymentMethod>

**Errors:**
| Code | HTTP | Message |
|------|------|----------|
| CUSTOMER_NOT_FOUND | 404 | Customer not found |

**Preconditions:** [Customer exists]

#### deleteCustomer `DELETE /{customerId}`

Delete or deactivate customer account (GDPR compliance)

**Input:**
- `customerId`: CustomerId (required)
- `deleteType`: DeleteType (required)

**Output:**
- `success`: boolean
- `deleteType`: DeleteType

**Errors:**
| Code | HTTP | Message |
|------|------|----------|
| CUSTOMER_NOT_FOUND | 404 | Customer not found |
| PENDING_ORDERS_EXIST | 409 | Cannot delete customer with pending orders |

**Preconditions:** [Customer exists No pending orders]

**Postconditions:** [Customer soft-deleted or data erased based on deleteType CustomerDeleted event emitted]

### Events

- **CustomerRegistered**: Emitted when a new customer account is created (payload: [customerId email createdAt])
- **EmailVerified**: Emitted when customer verifies their email address (payload: [customerId verifiedAt])
- **CustomerDeleted**: Emitted when customer account is deleted or deactivated (payload: [customerId deleteType deletedAt])

---

## IC-PROD-001 – Product Service {#ic-prod-001}

**Purpose:** Manages product catalog including creation, editing, variants, and images

**Base URL:** `/api/v1/products`

**Security:**
- Authentication: Bearer JWT
- Authorization: Read: Public for active products; Write: Admin only

### Operations

#### listProducts `GET /`

Browse products with filtering, sorting, and pagination

**Input:**
- `categoryId`: CategoryId
- `minPrice`: Money
- `maxPrice`: Money
- `inStock`: boolean
- `sortBy`: SortOption
- `page`: integer
- `pageSize`: integer

**Output:**
- `products`: Array<ProductSummary>
- `totalCount`: integer
- `page`: integer
- `pageSize`: integer
- `totalPages`: integer

**Errors:**
| Code | HTTP | Message |
|------|------|----------|
| INVALID_PRICE_RANGE | 400 | Minimum price cannot exceed maximum price |
| INVALID_PAGE | 400 | Page number must be positive |

**Postconditions:** [Returns products matching filters with pagination]

#### getProduct `GET /{productId}`

Get detailed product information including variants and stock status

**Input:**
- `productId`: ProductId (required)

**Output:**
- `productId`: ProductId
- `description`: string
- `price`: Money
- `categoryId`: CategoryId
- `variants`: Array<ProductVariant>
- `isActive`: boolean
- `stockStatus`: StockStatus
- `images`: Array<ProductImage>
- `name`: string

**Errors:**
| Code | HTTP | Message |
|------|------|----------|
| PRODUCT_NOT_FOUND | 404 | Product not found |

**Preconditions:** [Product exists and is active (for customers)]

#### createProduct `POST /`

Create a new product in draft status

**Input:**
- `name`: string (required)
- `description`: string
- `price`: Money (required)
- `categoryId`: CategoryId (required)

**Output:**
- `productId`: ProductId
- `name`: string
- `description`: string
- `price`: Money
- `categoryId`: CategoryId
- `status`: ProductStatus
- `createdAt`: DateTime
- `createdBy`: UserId

**Errors:**
| Code | HTTP | Message |
|------|------|----------|
| INVALID_PRODUCT_NAME | 400 | Product name must be between 2 and 200 characters |
| DESCRIPTION_TOO_LONG | 400 | Description cannot exceed 5000 characters |
| INVALID_PRICE | 400 | Price must be at least $0.01 |
| CATEGORY_NOT_FOUND | 404 | Category not found |

**Preconditions:** [User has admin role Category exists]

**Postconditions:** [Product created in draft status ProductCreated event emitted]

#### updateProduct `PUT /{productId}`

Update product attributes with audit logging

**Input:**
- `price`: Money
- `categoryId`: CategoryId
- `isActive`: boolean
- `productId`: ProductId (required)
- `name`: string
- `description`: string

**Output:**
- `price`: Money
- `categoryId`: CategoryId
- `isActive`: boolean
- `updatedAt`: DateTime
- `updatedBy`: UserId
- `productId`: ProductId
- `name`: string
- `description`: string

**Errors:**
| Code | HTTP | Message |
|------|------|----------|
| PRODUCT_NOT_FOUND | 404 | Product not found |
| INVALID_PRODUCT_NAME | 400 | Product name must be between 2 and 200 characters |
| DESCRIPTION_TOO_LONG | 400 | Description cannot exceed 5000 characters |
| INVALID_PRICE | 400 | Price must be at least $0.01 |
| CATEGORY_NOT_FOUND | 404 | Category not found |

**Preconditions:** [User has admin role Product exists]

**Postconditions:** [Product updated Price changes logged in audit trail ProductUpdated event emitted Carts with product reflect new price]

#### deleteProduct `DELETE /{productId}`

Soft delete product preserving order history

**Input:**
- `productId`: ProductId (required)

**Output:**
- `deletedAt`: DateTime
- `success`: boolean

**Errors:**
| Code | HTTP | Message |
|------|------|----------|
| PRODUCT_NOT_FOUND | 404 | Product not found |

**Preconditions:** [User has admin role Product exists]

**Postconditions:** [Product soft-deleted (deletedAt set) Product removed from active listings Product removed from carts ProductDeleted event emitted]

#### addProductImage `POST /{productId}/images`

Upload product image (max 10 per product)

**Input:**
- `isPrimary`: boolean
- `productId`: ProductId (required)
- `image`: binary (required)

**Output:**
- `imageId`: ImageId
- `url`: string
- `isPrimary`: boolean

**Errors:**
| Code | HTTP | Message |
|------|------|----------|
| PRODUCT_NOT_FOUND | 404 | Product not found |
| IMAGE_LIMIT_EXCEEDED | 400 | Products cannot have more than 10 images |
| INVALID_IMAGE_FORMAT | 400 | Image format not supported |

**Preconditions:** [User has admin role Product exists Product has fewer than 10 images]

**Postconditions:** [Image uploaded and associated with product If isPrimary=true, previous primary unset]

#### addProductVariant `POST /{productId}/variants`

Add a variant (size/color) to product

**Input:**
- `color`: string
- `priceAdjustment`: Money
- `productId`: ProductId (required)
- `sku`: string (required)
- `size`: string

**Output:**
- `size`: string
- `color`: string
- `priceAdjustment`: Money
- `variantId`: VariantId
- `sku`: string

**Errors:**
| Code | HTTP | Message |
|------|------|----------|
| PRODUCT_NOT_FOUND | 404 | Product not found |
| DUPLICATE_SKU | 409 | SKU already exists |
| VARIANT_ATTRIBUTE_REQUIRED | 400 | At least one variant attribute (size or color) must be specified |

**Preconditions:** [User has admin role Product exists SKU is unique]

**Postconditions:** [Variant added to product]

### Events

- **ProductCreated**: Emitted when a new product is created (payload: [productId name price categoryId createdBy createdAt])
- **ProductUpdated**: Emitted when product attributes are modified (payload: [productId changedFields previousPrice newPrice updatedBy updatedAt])
- **ProductDeleted**: Emitted when product is soft-deleted (payload: [productId deletedBy deletedAt])
- **ProductDeactivated**: Emitted when product is deactivated (payload: [productId deactivatedBy deactivatedAt])

---

## IC-CAT-001 – Category Service {#ic-cat-001}

**Purpose:** Manages product category hierarchy for organization and browsing

**Base URL:** `/api/v1/categories`

**Security:**
- Authentication: Bearer JWT for write operations
- Authorization: Read: Public; Write: Admin only

### Operations

#### listCategories `GET /`

Get category hierarchy up to 3 levels deep

**Input:**
- `parentId`: CategoryId
- `includeInactive`: boolean

**Output:**
- `categories`: Array<Category>

**Postconditions:** [Returns category tree]

#### createCategory `POST /`

Create a new category with optional parent

**Input:**
- `parentCategoryId`: CategoryId
- `name`: string (required)
- `description`: string

**Output:**
- `categoryId`: CategoryId
- `name`: string
- `description`: string
- `parentCategoryId`: CategoryId
- `depth`: integer

**Errors:**
| Code | HTTP | Message |
|------|------|----------|
| DUPLICATE_CATEGORY_NAME | 409 | Category name must be unique |
| CATEGORY_DEPTH_EXCEEDED | 400 | Categories cannot be nested more than 3 levels deep |
| PARENT_CATEGORY_NOT_FOUND | 404 | Parent category not found |

**Preconditions:** [User has admin role Name is unique Depth would not exceed 3]

**Postconditions:** [Category created CategoryCreated event emitted]

#### getCategory `GET /{categoryId}`

Get category details with products

**Input:**
- `categoryId`: CategoryId (required)

**Output:**
- `name`: string
- `description`: string
- `parentCategoryId`: CategoryId
- `childCategories`: Array<Category>
- `productCount`: integer
- `categoryId`: CategoryId

**Errors:**
| Code | HTTP | Message |
|------|------|----------|
| CATEGORY_NOT_FOUND | 404 | Category not found |

**Preconditions:** [Category exists]

### Events

- **CategoryCreated**: Emitted when a new category is created (payload: [categoryId name parentCategoryId createdAt])

---

## IC-CART-001 – Cart Service {#ic-cart-001}

**Purpose:** Manages shopping carts for customers including items, quantities, and cart merging

**Base URL:** `/api/v1/carts`

**Security:**
- Authentication: Bearer JWT or session token for guest
- Authorization: Customer can only access own cart

### Operations

#### getCart `GET /{cartId}`

Get cart contents with current prices and stock status

**Input:**
- `cartId`: CartId (required)

**Output:**
- `totalPrice`: Money
- `itemCount`: integer
- `priceChanges`: Array<PriceChange>
- `stockWarnings`: Array<StockWarning>
- `cartId`: CartId
- `customerId`: CustomerId
- `items`: Array<CartItem>

**Errors:**
| Code | HTTP | Message |
|------|------|----------|
| CART_NOT_FOUND | 404 | Cart not found |
| CART_EXPIRED | 410 | Cart has expired |

**Preconditions:** [Cart exists and not expired]

**Postconditions:** [Returns cart with current prices calculated]

#### addItem `POST /{cartId}/items`

Add product to cart or increment existing item quantity

**Input:**
- `cartId`: CartId (required)
- `productId`: ProductId (required)
- `variantId`: VariantId
- `quantity`: integer (required)

**Output:**
- `unitPrice`: Money
- `quantity`: integer
- `subtotal`: Money
- `cartItemId`: CartItemId
- `productId`: ProductId
- `variantId`: VariantId
- `productName`: string

**Errors:**
| Code | HTTP | Message |
|------|------|----------|
| CART_NOT_FOUND | 404 | Cart not found |
| PRODUCT_NOT_FOUND | 404 | Product not found |
| OUT_OF_STOCK | 409 | Product is out of stock |
| QUANTITY_EXCEEDS_STOCK | 409 | Requested quantity exceeds available stock |
| VARIANT_REQUIRED | 400 | Product has variants; variant selection required |
| VARIANT_NOT_FOUND | 404 | Variant not found |
| INVALID_QUANTITY | 400 | Quantity must be at least 1 |

**Preconditions:** [Cart exists Product exists and is active Stock available Variant selected if product has variants]

**Postconditions:** [Item added or quantity incremented Cart total recalculated ItemAddedToCart event emitted]

#### updateItemQuantity `PATCH /{cartId}/items/{cartItemId}`

Update quantity of item in cart

**Input:**
- `cartItemId`: CartItemId (required)
- `quantity`: integer (required)
- `cartId`: CartId (required)

**Output:**
- `quantity`: integer
- `subtotal`: Money
- `cartTotal`: Money
- `cartItemId`: CartItemId

**Errors:**
| Code | HTTP | Message |
|------|------|----------|
| CART_NOT_FOUND | 404 | Cart not found |
| CART_ITEM_NOT_FOUND | 404 | Cart item not found |
| QUANTITY_EXCEEDS_STOCK | 409 | Requested quantity exceeds available stock |
| INVALID_QUANTITY | 400 | Quantity must be 0 or greater |

**Preconditions:** [Cart exists Item exists in cart]

**Postconditions:** [Quantity updated (or item removed if 0) Cart total recalculated CartQuantityUpdated event emitted]

#### removeItem `DELETE /{cartId}/items/{cartItemId}`

Remove item from cart

**Input:**
- `cartId`: CartId (required)
- `cartItemId`: CartItemId (required)

**Output:**
- `success`: boolean
- `cartTotal`: Money
- `undoToken`: string

**Errors:**
| Code | HTTP | Message |
|------|------|----------|
| CART_NOT_FOUND | 404 | Cart not found |
| CART_ITEM_NOT_FOUND | 404 | Cart item not found |

**Preconditions:** [Cart exists Item exists in cart]

**Postconditions:** [Item removed Cart total recalculated ItemRemovedFromCart event emitted Undo token generated]

#### undoRemoveItem `POST /{cartId}/items/undo`

Undo recent item removal using undo token

**Input:**
- `cartId`: CartId (required)
- `undoToken`: string (required)

**Output:**
- `cartItemId`: CartItemId
- `restored`: boolean

**Errors:**
| Code | HTTP | Message |
|------|------|----------|
| CART_NOT_FOUND | 404 | Cart not found |
| INVALID_UNDO_TOKEN | 400 | Undo token is invalid or expired |

**Preconditions:** [Cart exists Undo token valid and not expired]

**Postconditions:** [Item restored to cart]

#### createGuestCart `POST /guest`

Create a session-based cart for guest user

**Input:**
- `sessionId`: string (required)

**Output:**
- `cartId`: CartId
- `expiresAt`: DateTime

**Postconditions:** [Guest cart created with 7-day expiration]

#### mergeCarts `POST /{cartId}/merge`

Merge guest cart into customer cart on login

**Input:**
- `cartId`: CartId (required)
- `guestCartId`: CartId (required)

**Output:**
- `mergedItemCount`: integer
- `stockLimitedItems`: Array<StockWarning>
- `cartId`: CartId

**Errors:**
| Code | HTTP | Message |
|------|------|----------|
| CART_NOT_FOUND | 404 | Cart not found |
| GUEST_CART_NOT_FOUND | 404 | Guest cart not found |

**Preconditions:** [Both carts exist]

**Postconditions:** [Items merged with quantity combination Quantities capped at available stock Guest cart deleted]

### Events

- **ItemAddedToCart**: Emitted when item is added to cart (payload: [cartId cartItemId productId variantId quantity])
- **CartQuantityUpdated**: Emitted when cart item quantity changes (payload: [cartId cartItemId previousQuantity newQuantity])
- **ItemRemovedFromCart**: Emitted when item is removed from cart (payload: [cartId cartItemId productId])
- **CartCleared**: Emitted when cart is emptied (payload: [cartId reason])

---

## IC-ORDER-001 – Order Service {#ic-order-001}

**Purpose:** Manages order lifecycle from placement through delivery including status transitions

**Base URL:** `/api/v1/orders`

**Security:**
- Authentication: Bearer JWT
- Authorization: Customers: own orders only; Admin: all orders and status updates

### Operations

#### placeOrder `POST /`

Create order from cart with payment authorization

**Input:**
- `customerId`: CustomerId (required)
- `cartId`: CartId (required)
- `shippingAddressId`: AddressId (required)
- `paymentMethodId`: PaymentMethodId (required)

**Output:**
- `subtotal`: Money
- `createdAt`: DateTime
- `orderNumber`: string
- `lineItems`: Array<OrderLineItem>
- `shippingCost`: Money
- `status`: OrderStatus
- `tax`: Money
- `totalAmount`: Money
- `orderId`: OrderId

**Errors:**
| Code | HTTP | Message |
|------|------|----------|
| REGISTRATION_REQUIRED | 403 | Only registered customers can place orders |
| EMAIL_NOT_VERIFIED | 403 | Email verification required before placing orders |
| CART_EMPTY | 400 | Cannot place order with empty cart |
| OUT_OF_STOCK | 409 | One or more items are out of stock |
| INSUFFICIENT_STOCK | 409 | Requested quantity exceeds available stock |
| BACKORDER_NOT_SUPPORTED | 409 | Backorders are not supported |
| PAYMENT_REQUIRED | 402 | Payment authorization failed |
| INCOMPLETE_ADDRESS | 400 | Shipping address is incomplete |
| INTERNATIONAL_SHIPPING_NOT_AVAILABLE | 400 | International shipping not available |

**Preconditions:** [Customer registered and email verified Cart has items All items in stock Valid shipping address Valid payment method]

**Postconditions:** [Payment authorized Inventory decremented Order created with PENDING status Cart cleared Prices snapshot captured OrderPlaced event emitted Confirmation email queued]

#### getOrder `GET /{orderId}`

Get order details including line items and status

**Input:**
- `orderId`: OrderId (required)

**Output:**
- `lineItems`: Array<OrderLineItem>
- `totalAmount`: Money
- `status`: OrderStatus
- `subtotal`: Money
- `shippingCost`: Money
- `createdAt`: DateTime
- `orderId`: OrderId
- `orderNumber`: string
- `shippingAddress`: ShippingAddress
- `tax`: Money
- `trackingNumber`: string
- `updatedAt`: DateTime
- `customerId`: CustomerId

**Errors:**
| Code | HTTP | Message |
|------|------|----------|
| ORDER_NOT_FOUND | 404 | Order not found |

**Preconditions:** [Order exists User authorized to view]

#### listOrders `GET /`

List orders with filtering and pagination

**Input:**
- `page`: integer
- `pageSize`: integer
- `customerId`: CustomerId
- `status`: OrderStatus
- `fromDate`: DateTime
- `toDate`: DateTime

**Output:**
- `orders`: Array<OrderSummary>
- `totalCount`: integer
- `page`: integer
- `pageSize`: integer

**Preconditions:** [User authorized]

#### confirmOrder `POST /{orderId}/confirm`

Confirm pending order

**Input:**
- `orderId`: OrderId (required)

**Output:**
- `confirmedAt`: DateTime
- `orderId`: OrderId
- `status`: OrderStatus

**Errors:**
| Code | HTTP | Message |
|------|------|----------|
| ORDER_NOT_FOUND | 404 | Order not found |
| INVALID_STATUS_TRANSITION | 409 | Order must be in PENDING status to confirm |

**Preconditions:** [Order exists Status is PENDING]

**Postconditions:** [Status changed to CONFIRMED OrderConfirmed event emitted Status change logged]

#### shipOrder `POST /{orderId}/ship`

Mark order as shipped with optional tracking number

**Input:**
- `orderId`: OrderId (required)
- `trackingNumber`: string

**Output:**
- `trackingNumber`: string
- `shippedAt`: DateTime
- `orderId`: OrderId
- `status`: OrderStatus

**Errors:**
| Code | HTTP | Message |
|------|------|----------|
| ORDER_NOT_FOUND | 404 | Order not found |
| INVALID_STATUS_TRANSITION | 409 | Order must be in CONFIRMED status to ship |

**Preconditions:** [Order exists Status is CONFIRMED]

**Postconditions:** [Status changed to SHIPPED OrderShipped event emitted Customer notification sent]

#### deliverOrder `POST /{orderId}/deliver`

Mark order as delivered

**Input:**
- `orderId`: OrderId (required)

**Output:**
- `status`: OrderStatus
- `deliveredAt`: DateTime
- `orderId`: OrderId

**Errors:**
| Code | HTTP | Message |
|------|------|----------|
| ORDER_NOT_FOUND | 404 | Order not found |
| INVALID_STATUS_TRANSITION | 409 | Order must be in SHIPPED status to deliver |

**Preconditions:** [Order exists Status is SHIPPED]

**Postconditions:** [Status changed to DELIVERED OrderDelivered event emitted]

#### cancelOrder `POST /{orderId}/cancel`

Cancel order with inventory restoration and refund

**Input:**
- `orderId`: OrderId (required)
- `reason`: string

**Output:**
- `refundInitiated`: boolean
- `cancelledAt`: DateTime
- `orderId`: OrderId
- `status`: OrderStatus

**Errors:**
| Code | HTTP | Message |
|------|------|----------|
| ORDER_NOT_FOUND | 404 | Order not found |
| CANCELLATION_NOT_ALLOWED | 409 | Order cannot be cancelled after shipping |
| ORDER_ALREADY_CANCELLED | 409 | Order is already cancelled |

**Preconditions:** [Order exists Status is PENDING or CONFIRMED]

**Postconditions:** [Status changed to CANCELLED Inventory restored Refund initiated OrderCancelled event emitted Cancellation email sent]

#### searchOrders `GET /search`

Search orders by order number

**Input:**
- `orderNumber`: string (required)

**Output:**
- `orders`: Array<OrderSummary>

**Preconditions:** [User has admin role]

### Events

- **OrderPlaced**: Emitted when order is created (payload: [orderId orderNumber customerId totalAmount createdAt])
- **OrderConfirmed**: Emitted when order is confirmed (payload: [orderId confirmedBy confirmedAt])
- **OrderShipped**: Emitted when order is shipped (payload: [orderId trackingNumber shippedAt])
- **OrderDelivered**: Emitted when order is delivered (payload: [orderId deliveredAt])
- **OrderCancelled**: Emitted when order is cancelled (payload: [orderId reason cancelledBy cancelledAt])

---

## IC-INV-001 – Inventory Service {#ic-inv-001}

**Purpose:** Tracks and manages stock levels, reservations, and inventory adjustments

**Base URL:** `/api/v1/inventory`

**Security:**
- Authentication: Bearer JWT
- Authorization: Read: Internal services; Write: Admin only

### Operations

#### getInventory `GET /{productId}`

Get current inventory levels for a product

**Input:**
- `productId`: ProductId (required)
- `variantId`: VariantId

**Output:**
- `variantId`: VariantId
- `quantityOnHand`: integer
- `reservedQuantity`: integer
- `availableQuantity`: integer
- `inventoryId`: InventoryId
- `productId`: ProductId

**Errors:**
| Code | HTTP | Message |
|------|------|----------|
| PRODUCT_NOT_FOUND | 404 | Product not found |
| INVENTORY_NOT_FOUND | 404 | Inventory record not found |

**Preconditions:** [Product exists]

#### setStock `PUT /{productId}`

Set absolute stock level with audit logging

**Input:**
- `productId`: ProductId (required)
- `variantId`: VariantId
- `quantity`: integer (required)
- `reason`: string

**Output:**
- `inventoryId`: InventoryId
- `previousQuantity`: integer
- `newQuantity`: integer
- `availableQuantity`: integer

**Errors:**
| Code | HTTP | Message |
|------|------|----------|
| PRODUCT_NOT_FOUND | 404 | Product not found |
| NEGATIVE_INVENTORY | 400 | Inventory cannot be set below zero |

**Preconditions:** [Product exists User has admin role]

**Postconditions:** [Stock level updated Audit trail recorded Low stock alert if threshold crossed]

#### adjustStock `POST /{productId}/adjust`

Adjust stock by delta (positive or negative)

**Input:**
- `delta`: integer (required)
- `reason`: string (required)
- `productId`: ProductId (required)
- `variantId`: VariantId

**Output:**
- `newQuantity`: integer
- `availableQuantity`: integer
- `inventoryId`: InventoryId
- `previousQuantity`: integer

**Errors:**
| Code | HTTP | Message |
|------|------|----------|
| PRODUCT_NOT_FOUND | 404 | Product not found |
| NEGATIVE_INVENTORY | 400 | Adjustment would result in negative inventory |

**Preconditions:** [Product exists User has admin role Result would be >= 0]

**Postconditions:** [Stock adjusted Audit trail recorded with reason Low stock alert if threshold crossed]

#### reserveStock `POST /{productId}/reserve`

Reserve stock for pending order

**Input:**
- `orderId`: OrderId (required)
- `productId`: ProductId (required)
- `variantId`: VariantId
- `quantity`: integer (required)

**Output:**
- `reservationId`: ReservationId
- `productId`: ProductId
- `quantity`: integer
- `availableQuantity`: integer

**Errors:**
| Code | HTTP | Message |
|------|------|----------|
| PRODUCT_NOT_FOUND | 404 | Product not found |
| INSUFFICIENT_STOCK | 409 | Not enough stock to reserve |

**Preconditions:** [Product exists Available quantity >= requested]

**Postconditions:** [Reserved quantity increased InventoryReserved event emitted]

#### releaseReservation `DELETE /{productId}/reserve/{reservationId}`

Release reserved stock (e.g., order cancelled)

**Input:**
- `productId`: ProductId (required)
- `reservationId`: ReservationId (required)

**Output:**
- `released`: boolean
- `availableQuantity`: integer

**Errors:**
| Code | HTTP | Message |
|------|------|----------|
| PRODUCT_NOT_FOUND | 404 | Product not found |
| RESERVATION_NOT_FOUND | 404 | Reservation not found |

**Preconditions:** [Reservation exists]

**Postconditions:** [Reserved quantity decreased InventoryReleased event emitted]

#### deductStock `POST /{productId}/deduct`

Deduct stock from on-hand (order shipped)

**Input:**
- `productId`: ProductId (required)
- `variantId`: VariantId
- `quantity`: integer (required)
- `orderId`: OrderId (required)

**Output:**
- `previousQuantity`: integer
- `newQuantity`: integer

**Errors:**
| Code | HTTP | Message |
|------|------|----------|
| PRODUCT_NOT_FOUND | 404 | Product not found |
| INSUFFICIENT_STOCK | 409 | Not enough stock to deduct |

**Preconditions:** [Product exists Quantity on hand >= requested]

**Postconditions:** [Quantity on hand decreased InventoryDeducted event emitted]

#### restockInventory `POST /{productId}/restock`

Add stock to inventory

**Input:**
- `variantId`: VariantId
- `quantity`: integer (required)
- `reason`: string
- `productId`: ProductId (required)

**Output:**
- `newQuantity`: integer
- `availableQuantity`: integer
- `previousQuantity`: integer

**Errors:**
| Code | HTTP | Message |
|------|------|----------|
| PRODUCT_NOT_FOUND | 404 | Product not found |
| INVALID_QUANTITY | 400 | Quantity must be positive |

**Preconditions:** [Product exists Quantity > 0]

**Postconditions:** [Quantity on hand increased InventoryRestocked event emitted Audit trail recorded]

#### bulkImport `POST /bulk-import`

Import inventory levels from CSV

**Input:**
- `format`: string (required)
- `file`: binary (required)

**Output:**
- `imported`: integer
- `failed`: integer
- `errors`: Array<ImportError>

**Errors:**
| Code | HTTP | Message |
|------|------|----------|
| INVALID_FILE_FORMAT | 400 | File format not supported |
| VALIDATION_ERRORS | 400 | One or more rows have validation errors |

**Preconditions:** [User has admin role File is valid CSV]

**Postconditions:** [Valid rows imported Failed rows returned with errors]

### Events

- **InventoryReserved**: Emitted when stock is reserved for an order (payload: [inventoryId productId variantId quantity orderId])
- **InventoryReleased**: Emitted when reservation is released (payload: [inventoryId productId variantId quantity orderId])
- **InventoryDeducted**: Emitted when stock is deducted (shipped) (payload: [inventoryId productId variantId quantity orderId])
- **InventoryRestocked**: Emitted when stock is added (payload: [inventoryId productId variantId quantity reason])
- **OutOfStock**: Emitted when available quantity reaches zero (payload: [inventoryId productId variantId])
- **LowStockAlert**: Emitted when stock falls below threshold (payload: [inventoryId productId variantId availableQuantity threshold])

---

