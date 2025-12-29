# Initial Data Model

Generated: 2025-12-29T15:49:16+01:00

---

## Enumerations

### order_status

Values: `[pending confirmed shipped delivered cancelled]`

### payment_type

Values: `[credit_card paypal]`

### card_brand

Values: `[visa mastercard amex]`

### registration_status

Values: `[unverified registered]`

### product_status

Values: `[draft active inactive]`

### inventory_operation_type

Values: `[set adjust reserve release deduct restock]`

### email_type

Values: `[order_confirmation order_shipped order_cancelled verification password_reset]`

### email_status

Values: `[pending processing sent failed]`

---

## Tables

### TBL-CUSTOMER-001 – customers

**Aggregate:** Customer

**Purpose:** Stores customer account information

**Columns:**

| Name | Type | Constraints | Default |
|------|------|-------------|----------|
| id | UUID | NOT NULL | gen_random_uuid() |
| email | VARCHAR(255) | NOT NULL UNIQUE |  |
| password_hash | VARCHAR(255) | NOT NULL |  |
| first_name | VARCHAR(100) | NOT NULL |  |
| last_name | VARCHAR(100) | NOT NULL |  |
| registration_status | VARCHAR(20) | NOT NULL | 'unverified' |
| email_verified | BOOLEAN | NOT NULL | false |
| verification_token | VARCHAR(255) |  |  |
| verification_token_expires_at | TIMESTAMP |  |  |
| version | INTEGER | NOT NULL | 1 |
| created_at | TIMESTAMP | NOT NULL | NOW() |
| updated_at | TIMESTAMP | NOT NULL | NOW() |

**Primary Key:** [id]

**Indexes:**
- `idx_customers_email`: [email]
- `idx_customers_verification_token`: [verification_token]

**Constraints:**
- `chk_customer_registration_status`: registration_status IN ('unverified', 'registered')

---

### TBL-CUSTOMER-002 – customer_shipping_addresses

**Aggregate:** Customer

**Purpose:** Stores customer saved shipping addresses

**Columns:**

| Name | Type | Constraints | Default |
|------|------|-------------|----------|
| id | UUID | NOT NULL | gen_random_uuid() |
| customer_id | UUID | NOT NULL |  |
| street | VARCHAR(200) | NOT NULL |  |
| city | VARCHAR(100) | NOT NULL |  |
| state | VARCHAR(100) | NOT NULL |  |
| postal_code | VARCHAR(20) | NOT NULL |  |
| country | CHAR(2) | NOT NULL |  |
| recipient_name | VARCHAR(200) | NOT NULL |  |
| is_default | BOOLEAN | NOT NULL | false |
| created_at | TIMESTAMP | NOT NULL | NOW() |
| updated_at | TIMESTAMP | NOT NULL | NOW() |

**Primary Key:** [id]

**Indexes:**
- `idx_shipping_addresses_customer`: [customer_id]

**Foreign Keys:**
- [customer_id] → customers(id) (ON DELETE CASCADE)

---

### TBL-CART-001 – carts

**Aggregate:** Cart

**Purpose:** Stores shopping cart header information

**Columns:**

| Name | Type | Constraints | Default |
|------|------|-------------|----------|
| id | UUID | NOT NULL | gen_random_uuid() |
| customer_id | UUID |  |  |
| is_guest | BOOLEAN | NOT NULL | false |
| guest_session_id | VARCHAR(255) |  |  |
| total_price_amount | DECIMAL(10,2) | NOT NULL | 0.00 |
| total_price_currency | CHAR(3) | NOT NULL | 'USD' |
| last_activity_date | TIMESTAMP | NOT NULL | NOW() |
| version | INTEGER | NOT NULL | 1 |
| created_at | TIMESTAMP | NOT NULL | NOW() |
| updated_at | TIMESTAMP | NOT NULL | NOW() |

**Primary Key:** [id]

**Indexes:**
- `idx_carts_customer`: [customer_id]
- `idx_carts_guest_session`: [guest_session_id]
- `idx_carts_last_activity`: [last_activity_date]

**Foreign Keys:**
- [customer_id] → customers(id) (ON DELETE CASCADE)

**Constraints:**
- `chk_cart_total_price`: total_price_amount >= 0

---

### TBL-CART-002 – cart_items

**Aggregate:** Cart

**Purpose:** Stores individual items in shopping carts

**Columns:**

| Name | Type | Constraints | Default |
|------|------|-------------|----------|
| id | UUID | NOT NULL | gen_random_uuid() |
| cart_id | UUID | NOT NULL |  |
| product_id | UUID | NOT NULL |  |
| variant_id | UUID |  |  |
| product_name | VARCHAR(200) | NOT NULL |  |
| variant_description | VARCHAR(200) |  |  |
| unit_price_amount | DECIMAL(10,2) | NOT NULL |  |
| unit_price_currency | CHAR(3) | NOT NULL | 'USD' |
| quantity | INTEGER | NOT NULL |  |
| subtotal_amount | DECIMAL(10,2) | NOT NULL |  |
| subtotal_currency | CHAR(3) | NOT NULL | 'USD' |
| added_at | TIMESTAMP | NOT NULL | NOW() |
| created_at | TIMESTAMP | NOT NULL | NOW() |
| updated_at | TIMESTAMP | NOT NULL | NOW() |

**Primary Key:** [id]

**Indexes:**
- `idx_cart_items_cart`: [cart_id]
- `idx_cart_items_product`: [product_id]
- `idx_cart_items_cart_product`: [cart_id product_id variant_id]

**Foreign Keys:**
- [cart_id] → carts(id) (ON DELETE CASCADE)
- [product_id] → products(id) (ON DELETE RESTRICT)
- [variant_id] → product_variants(id) (ON DELETE RESTRICT)

**Constraints:**
- `chk_cart_item_quantity`: quantity > 0
- `chk_cart_item_price`: unit_price_amount >= 0 AND subtotal_amount >= 0

---

### TBL-ORDER-001 – orders

**Aggregate:** Order

**Purpose:** Stores order header information

**Columns:**

| Name | Type | Constraints | Default |
|------|------|-------------|----------|
| id | UUID | NOT NULL | gen_random_uuid() |
| order_number | VARCHAR(20) | NOT NULL UNIQUE |  |
| customer_id | UUID | NOT NULL |  |
| status | VARCHAR(20) | NOT NULL | 'pending' |
| subtotal_amount | DECIMAL(10,2) | NOT NULL |  |
| subtotal_currency | CHAR(3) | NOT NULL | 'USD' |
| shipping_cost_amount | DECIMAL(10,2) | NOT NULL |  |
| shipping_cost_currency | CHAR(3) | NOT NULL | 'USD' |
| tax_amount | DECIMAL(10,2) | NOT NULL | 0.00 |
| tax_currency | CHAR(3) | NOT NULL | 'USD' |
| total_amount | DECIMAL(10,2) | NOT NULL |  |
| total_currency | CHAR(3) | NOT NULL | 'USD' |
| shipping_street | VARCHAR(200) | NOT NULL |  |
| shipping_city | VARCHAR(100) | NOT NULL |  |
| shipping_state | VARCHAR(100) | NOT NULL |  |
| shipping_postal_code | VARCHAR(20) | NOT NULL |  |
| shipping_country | CHAR(2) | NOT NULL |  |
| shipping_recipient_name | VARCHAR(200) | NOT NULL |  |
| payment_type | VARCHAR(20) | NOT NULL |  |
| payment_last_four_digits | CHAR(4) |  |  |
| payment_card_brand | VARCHAR(20) |  |  |
| payment_paypal_email | VARCHAR(255) |  |  |
| payment_transaction_id | VARCHAR(255) |  |  |
| tracking_number | VARCHAR(100) |  |  |
| cancellation_reason | VARCHAR(500) |  |  |
| version | INTEGER | NOT NULL | 1 |
| created_at | TIMESTAMP | NOT NULL | NOW() |
| updated_at | TIMESTAMP | NOT NULL | NOW() |

**Primary Key:** [id]

**Indexes:**
- `idx_orders_order_number`: [order_number]
- `idx_orders_customer`: [customer_id]
- `idx_orders_status`: [status]
- `idx_orders_created`: [created_at]
- `idx_orders_customer_status`: [customer_id status]

**Foreign Keys:**
- [customer_id] → customers(id) (ON DELETE RESTRICT)

**Constraints:**
- `chk_order_status`: status IN ('pending', 'confirmed', 'shipped', 'delivered', 'cancelled')
- `chk_order_amounts`: total_amount >= 0 AND subtotal_amount >= 0 AND shipping_cost_amount >= 0 AND tax_amount >= 0
- `chk_order_payment_type`: payment_type IN ('credit_card', 'paypal')

---

### TBL-ORDER-002 – order_line_items

**Aggregate:** Order

**Purpose:** Stores order line items (immutable snapshots)

**Columns:**

| Name | Type | Constraints | Default |
|------|------|-------------|----------|
| id | UUID | NOT NULL | gen_random_uuid() |
| order_id | UUID | NOT NULL |  |
| product_id | UUID | NOT NULL |  |
| variant_id | UUID |  |  |
| product_name | VARCHAR(200) | NOT NULL |  |
| variant_description | VARCHAR(200) |  |  |
| unit_price_amount | DECIMAL(10,2) | NOT NULL |  |
| unit_price_currency | CHAR(3) | NOT NULL | 'USD' |
| quantity | INTEGER | NOT NULL |  |
| subtotal_amount | DECIMAL(10,2) | NOT NULL |  |
| subtotal_currency | CHAR(3) | NOT NULL | 'USD' |
| created_at | TIMESTAMP | NOT NULL | NOW() |

**Primary Key:** [id]

**Indexes:**
- `idx_order_line_items_order`: [order_id]
- `idx_order_line_items_product`: [product_id]

**Foreign Keys:**
- [order_id] → orders(id) (ON DELETE CASCADE)
- [product_id] → products(id) (ON DELETE RESTRICT)

**Constraints:**
- `chk_line_item_quantity`: quantity > 0
- `chk_line_item_price`: unit_price_amount >= 0 AND subtotal_amount >= 0

---

### TBL-ORDER-003 – order_number_sequence

**Aggregate:** Order

**Purpose:** Tracks order number sequence per year

**Columns:**

| Name | Type | Constraints | Default |
|------|------|-------------|----------|
| year | INTEGER | NOT NULL |  |
| last_sequence | INTEGER | NOT NULL | 0 |

**Primary Key:** [year]

**Constraints:**
- `chk_sequence_positive`: last_sequence >= 0

---

### TBL-PRODUCT-001 – products

**Aggregate:** Product

**Purpose:** Stores product catalog information

**Columns:**

| Name | Type | Constraints | Default |
|------|------|-------------|----------|
| id | UUID | NOT NULL | gen_random_uuid() |
| name | VARCHAR(200) | NOT NULL |  |
| description | TEXT |  |  |
| price_amount | DECIMAL(10,2) | NOT NULL |  |
| price_currency | CHAR(3) | NOT NULL | 'USD' |
| category_id | UUID | NOT NULL |  |
| is_active | BOOLEAN | NOT NULL | true |
| is_deleted | BOOLEAN | NOT NULL | false |
| status | VARCHAR(20) | NOT NULL | 'draft' |
| created_by | UUID |  |  |
| updated_by | UUID |  |  |
| version | INTEGER | NOT NULL | 1 |
| created_at | TIMESTAMP | NOT NULL | NOW() |
| updated_at | TIMESTAMP | NOT NULL | NOW() |

**Primary Key:** [id]

**Indexes:**
- `idx_products_category`: [category_id]
- `idx_products_active`: [is_active is_deleted]
- `idx_products_name`: [name]
- `idx_products_created`: [created_at]

**Foreign Keys:**
- [category_id] → categories(id) (ON DELETE RESTRICT)

**Constraints:**
- `chk_product_price`: price_amount >= 0.01
- `chk_product_name_length`: LENGTH(name) >= 2 AND LENGTH(name) <= 200
- `chk_product_status`: status IN ('draft', 'active', 'inactive')

---

### TBL-PRODUCT-002 – product_variants

**Aggregate:** Product

**Purpose:** Stores product variant information (size, color, etc.)

**Columns:**

| Name | Type | Constraints | Default |
|------|------|-------------|----------|
| id | UUID | NOT NULL | gen_random_uuid() |
| product_id | UUID | NOT NULL |  |
| sku | VARCHAR(50) | NOT NULL UNIQUE |  |
| size | VARCHAR(50) |  |  |
| color | VARCHAR(50) |  |  |
| price_adjustment_amount | DECIMAL(10,2) | NOT NULL | 0.00 |
| price_adjustment_currency | CHAR(3) | NOT NULL | 'USD' |
| created_at | TIMESTAMP | NOT NULL | NOW() |
| updated_at | TIMESTAMP | NOT NULL | NOW() |

**Primary Key:** [id]

**Indexes:**
- `idx_product_variants_product`: [product_id]
- `idx_product_variants_sku`: [sku]
- `idx_product_variants_combination`: [product_id size color]

**Foreign Keys:**
- [product_id] → products(id) (ON DELETE CASCADE)

**Constraints:**
- `chk_variant_has_attribute`: size IS NOT NULL OR color IS NOT NULL

---

### TBL-PRODUCT-003 – product_images

**Aggregate:** Product

**Purpose:** Stores product image URLs and metadata

**Columns:**

| Name | Type | Constraints | Default |
|------|------|-------------|----------|
| id | UUID | NOT NULL | gen_random_uuid() |
| product_id | UUID | NOT NULL |  |
| url | VARCHAR(500) | NOT NULL |  |
| is_primary | BOOLEAN | NOT NULL | false |
| display_order | INTEGER | NOT NULL | 0 |
| created_at | TIMESTAMP | NOT NULL | NOW() |

**Primary Key:** [id]

**Indexes:**
- `idx_product_images_product`: [product_id]
- `idx_product_images_primary`: [product_id is_primary]

**Foreign Keys:**
- [product_id] → products(id) (ON DELETE CASCADE)

---

### TBL-CATEGORY-001 – categories

**Aggregate:** Category

**Purpose:** Stores product category hierarchy

**Columns:**

| Name | Type | Constraints | Default |
|------|------|-------------|----------|
| id | UUID | NOT NULL | gen_random_uuid() |
| name | VARCHAR(100) | NOT NULL UNIQUE |  |
| description | TEXT |  |  |
| parent_category_id | UUID |  |  |
| depth | INTEGER | NOT NULL | 1 |
| display_order | INTEGER | NOT NULL | 0 |
| is_active | BOOLEAN | NOT NULL | true |
| version | INTEGER | NOT NULL | 1 |
| created_at | TIMESTAMP | NOT NULL | NOW() |
| updated_at | TIMESTAMP | NOT NULL | NOW() |

**Primary Key:** [id]

**Indexes:**
- `idx_categories_name`: [name]
- `idx_categories_parent`: [parent_category_id]
- `idx_categories_active`: [is_active]

**Foreign Keys:**
- [parent_category_id] → categories(id) (ON DELETE RESTRICT)

**Constraints:**
- `chk_category_depth`: depth >= 1 AND depth <= 3
- `chk_category_not_self_parent`: parent_category_id IS NULL OR parent_category_id != id

---

### TBL-INVENTORY-001 – inventory

**Aggregate:** Inventory

**Purpose:** Tracks product stock levels and reservations

**Columns:**

| Name | Type | Constraints | Default |
|------|------|-------------|----------|
| id | UUID | NOT NULL | gen_random_uuid() |
| product_id | UUID | NOT NULL UNIQUE |  |
| stock_level | INTEGER | NOT NULL | 0 |
| reserved_quantity | INTEGER | NOT NULL | 0 |
| low_stock_threshold | INTEGER | NOT NULL | 10 |
| version | INTEGER | NOT NULL | 1 |
| created_at | TIMESTAMP | NOT NULL | NOW() |
| updated_at | TIMESTAMP | NOT NULL | NOW() |

**Primary Key:** [id]

**Indexes:**
- `idx_inventory_product`: [product_id]
- `idx_inventory_low_stock`: [stock_level low_stock_threshold]

**Foreign Keys:**
- [product_id] → products(id) (ON DELETE CASCADE)

**Constraints:**
- `chk_inventory_stock_level`: stock_level >= 0
- `chk_inventory_reserved`: reserved_quantity >= 0 AND reserved_quantity <= stock_level
- `chk_inventory_threshold`: low_stock_threshold >= 0

---

### TBL-INVENTORY-002 – inventory_audit_log

**Aggregate:** Inventory

**Purpose:** Tracks all inventory changes for audit purposes

**Columns:**

| Name | Type | Constraints | Default |
|------|------|-------------|----------|
| id | UUID | NOT NULL | gen_random_uuid() |
| inventory_id | UUID | NOT NULL |  |
| product_id | UUID | NOT NULL |  |
| previous_value | INTEGER | NOT NULL |  |
| new_value | INTEGER | NOT NULL |  |
| adjustment | INTEGER | NOT NULL |  |
| reason | VARCHAR(255) | NOT NULL |  |
| operation_type | VARCHAR(50) | NOT NULL |  |
| user_id | UUID |  |  |
| order_id | UUID |  |  |
| created_at | TIMESTAMP | NOT NULL | NOW() |

**Primary Key:** [id]

**Indexes:**
- `idx_inventory_audit_inventory`: [inventory_id]
- `idx_inventory_audit_product`: [product_id]
- `idx_inventory_audit_created`: [created_at]
- `idx_inventory_audit_user`: [user_id]

**Foreign Keys:**
- [inventory_id] → inventory(id) (ON DELETE CASCADE)
- [product_id] → products(id) (ON DELETE CASCADE)

**Constraints:**
- `chk_audit_operation_type`: operation_type IN ('set', 'adjust', 'reserve', 'release', 'deduct', 'restock')

---

### TBL-PAYMENT-001 – saved_payment_methods

**Aggregate:** Customer

**Purpose:** Stores customer saved payment methods (tokenized)

**Columns:**

| Name | Type | Constraints | Default |
|------|------|-------------|----------|
| id | UUID | NOT NULL | gen_random_uuid() |
| customer_id | UUID | NOT NULL |  |
| payment_type | VARCHAR(20) | NOT NULL |  |
| token | VARCHAR(255) | NOT NULL |  |
| last_four_digits | CHAR(4) |  |  |
| card_brand | VARCHAR(20) |  |  |
| paypal_email | VARCHAR(255) |  |  |
| expires_at | VARCHAR(7) |  |  |
| is_default | BOOLEAN | NOT NULL | false |
| created_at | TIMESTAMP | NOT NULL | NOW() |
| updated_at | TIMESTAMP | NOT NULL | NOW() |

**Primary Key:** [id]

**Indexes:**
- `idx_saved_payments_customer`: [customer_id]
- `idx_saved_payments_default`: [customer_id is_default]

**Foreign Keys:**
- [customer_id] → customers(id) (ON DELETE CASCADE)

**Constraints:**
- `chk_payment_method_type`: payment_type IN ('credit_card', 'paypal')
- `chk_payment_card_brand`: card_brand IS NULL OR card_brand IN ('visa', 'mastercard', 'amex')

---

### TBL-EMAIL-001 – email_queue

**Aggregate:** System

**Purpose:** Queue for outbound emails with deduplication

**Columns:**

| Name | Type | Constraints | Default |
|------|------|-------------|----------|
| id | UUID | NOT NULL | gen_random_uuid() |
| email_type | VARCHAR(50) | NOT NULL |  |
| recipient_email | VARCHAR(255) | NOT NULL |  |
| subject | VARCHAR(255) | NOT NULL |  |
| reference_type | VARCHAR(50) | NOT NULL |  |
| reference_id | UUID | NOT NULL |  |
| payload | JSONB | NOT NULL |  |
| status | VARCHAR(20) | NOT NULL | 'pending' |
| retry_count | INTEGER | NOT NULL | 0 |
| last_error | TEXT |  |  |
| sent_at | TIMESTAMP |  |  |
| created_at | TIMESTAMP | NOT NULL | NOW() |
| updated_at | TIMESTAMP | NOT NULL | NOW() |

**Primary Key:** [id]

**Indexes:**
- `idx_email_queue_status`: [status]
- `idx_email_queue_reference`: [reference_type reference_id]
- `idx_email_queue_dedup`: [email_type reference_type reference_id]

**Constraints:**
- `chk_email_status`: status IN ('pending', 'processing', 'sent', 'failed')
- `chk_email_type`: email_type IN ('order_confirmation', 'order_shipped', 'order_cancelled', 'verification', 'password_reset')

---

## Entity-Relationship Diagram

```mermaid
erDiagram
    customers {
        UUID id
        VARCHAR(255) email
        VARCHAR(255) password_hash
        VARCHAR(100) first_name
        VARCHAR(100) last_name
        VARCHAR(20) registration_status
        BOOLEAN email_verified
        VARCHAR(255) verification_token
        TIMESTAMP verification_token_expires_at
        INTEGER version
        TIMESTAMP created_at
        TIMESTAMP updated_at
    }
    customer_shipping_addresses {
        UUID id
        UUID customer_id
        VARCHAR(200) street
        VARCHAR(100) city
        VARCHAR(100) state
        VARCHAR(20) postal_code
        CHAR(2) country
        VARCHAR(200) recipient_name
        BOOLEAN is_default
        TIMESTAMP created_at
        TIMESTAMP updated_at
    }
    carts {
        UUID id
        UUID customer_id
        BOOLEAN is_guest
        VARCHAR(255) guest_session_id
        DECIMAL(10,2) total_price_amount
        CHAR(3) total_price_currency
        TIMESTAMP last_activity_date
        INTEGER version
        TIMESTAMP created_at
        TIMESTAMP updated_at
    }
    cart_items {
        UUID id
        UUID cart_id
        UUID product_id
        UUID variant_id
        VARCHAR(200) product_name
        VARCHAR(200) variant_description
        DECIMAL(10,2) unit_price_amount
        CHAR(3) unit_price_currency
        INTEGER quantity
        DECIMAL(10,2) subtotal_amount
        CHAR(3) subtotal_currency
        TIMESTAMP added_at
        TIMESTAMP created_at
        TIMESTAMP updated_at
    }
    orders {
        UUID id
        VARCHAR(20) order_number
        UUID customer_id
        VARCHAR(20) status
        DECIMAL(10,2) subtotal_amount
        CHAR(3) subtotal_currency
        DECIMAL(10,2) shipping_cost_amount
        CHAR(3) shipping_cost_currency
        DECIMAL(10,2) tax_amount
        CHAR(3) tax_currency
        DECIMAL(10,2) total_amount
        CHAR(3) total_currency
        VARCHAR(200) shipping_street
        VARCHAR(100) shipping_city
        VARCHAR(100) shipping_state
        VARCHAR(20) shipping_postal_code
        CHAR(2) shipping_country
        VARCHAR(200) shipping_recipient_name
        VARCHAR(20) payment_type
        CHAR(4) payment_last_four_digits
        VARCHAR(20) payment_card_brand
        VARCHAR(255) payment_paypal_email
        VARCHAR(255) payment_transaction_id
        VARCHAR(100) tracking_number
        VARCHAR(500) cancellation_reason
        INTEGER version
        TIMESTAMP created_at
        TIMESTAMP updated_at
    }
    order_line_items {
        UUID id
        UUID order_id
        UUID product_id
        UUID variant_id
        VARCHAR(200) product_name
        VARCHAR(200) variant_description
        DECIMAL(10,2) unit_price_amount
        CHAR(3) unit_price_currency
        INTEGER quantity
        DECIMAL(10,2) subtotal_amount
        CHAR(3) subtotal_currency
        TIMESTAMP created_at
    }
    order_number_sequence {
        INTEGER year
        INTEGER last_sequence
    }
    products {
        UUID id
        VARCHAR(200) name
        TEXT description
        DECIMAL(10,2) price_amount
        CHAR(3) price_currency
        UUID category_id
        BOOLEAN is_active
        BOOLEAN is_deleted
        VARCHAR(20) status
        UUID created_by
        UUID updated_by
        INTEGER version
        TIMESTAMP created_at
        TIMESTAMP updated_at
    }
    product_variants {
        UUID id
        UUID product_id
        VARCHAR(50) sku
        VARCHAR(50) size
        VARCHAR(50) color
        DECIMAL(10,2) price_adjustment_amount
        CHAR(3) price_adjustment_currency
        TIMESTAMP created_at
        TIMESTAMP updated_at
    }
    product_images {
        UUID id
        UUID product_id
        VARCHAR(500) url
        BOOLEAN is_primary
        INTEGER display_order
        TIMESTAMP created_at
    }
    categories {
        UUID id
        VARCHAR(100) name
        TEXT description
        UUID parent_category_id
        INTEGER depth
        INTEGER display_order
        BOOLEAN is_active
        INTEGER version
        TIMESTAMP created_at
        TIMESTAMP updated_at
    }
    inventory {
        UUID id
        UUID product_id
        INTEGER stock_level
        INTEGER reserved_quantity
        INTEGER low_stock_threshold
        INTEGER version
        TIMESTAMP created_at
        TIMESTAMP updated_at
    }
    inventory_audit_log {
        UUID id
        UUID inventory_id
        UUID product_id
        INTEGER previous_value
        INTEGER new_value
        INTEGER adjustment
        VARCHAR(255) reason
        VARCHAR(50) operation_type
        UUID user_id
        UUID order_id
        TIMESTAMP created_at
    }
    saved_payment_methods {
        UUID id
        UUID customer_id
        VARCHAR(20) payment_type
        VARCHAR(255) token
        CHAR(4) last_four_digits
        VARCHAR(20) card_brand
        VARCHAR(255) paypal_email
        VARCHAR(7) expires_at
        BOOLEAN is_default
        TIMESTAMP created_at
        TIMESTAMP updated_at
    }
    email_queue {
        UUID id
        VARCHAR(50) email_type
        VARCHAR(255) recipient_email
        VARCHAR(255) subject
        VARCHAR(50) reference_type
        UUID reference_id
        JSONB payload
        VARCHAR(20) status
        INTEGER retry_count
        TEXT last_error
        TIMESTAMP sent_at
        TIMESTAMP created_at
        TIMESTAMP updated_at
    }
    customers ||--o{ customer_shipping_addresses : "FK"
    customers ||--o{ carts : "FK"
    carts ||--o{ cart_items : "FK"
    products ||--o{ cart_items : "FK"
    product_variants ||--o{ cart_items : "FK"
    customers ||--o{ orders : "FK"
    orders ||--o{ order_line_items : "FK"
    products ||--o{ order_line_items : "FK"
    categories ||--o{ products : "FK"
    products ||--o{ product_variants : "FK"
    products ||--o{ product_images : "FK"
    categories ||--o{ categories : "FK"
    products ||--o{ inventory : "FK"
    inventory ||--o{ inventory_audit_log : "FK"
    products ||--o{ inventory_audit_log : "FK"
    customers ||--o{ saved_payment_methods : "FK"
```
