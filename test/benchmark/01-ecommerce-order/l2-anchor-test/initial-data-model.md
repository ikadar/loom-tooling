# Initial Data Model

Generated: 2025-12-31T13:24:23+01:00

---

## Enumerations

### registration_status

Values: `[unverified verified suspended]`

### order_status

Values: `[pending confirmed shipped delivered cancelled]`

### payment_type

Values: `[credit_card paypal]`

### product_status

Values: `[draft active archived]`

### inventory_action

Values: `[set adjust reserve release deduct restock]`

---

## Tables

### TBL-CUSTOMER-001 – customers {#tbl-customer-001}

**Aggregate:** Customer

**Purpose:** Stores customer account information and authentication details

**Columns:**

| Name | Type | Constraints | Default |
|------|------|-------------|----------|
| id | UUID | NOT NULL | gen_random_uuid() |
| email | VARCHAR(255) | NOT NULL UNIQUE |  |
| password_hash | VARCHAR(255) | NOT NULL |  |
| registration_status | VARCHAR(20) | NOT NULL | 'unverified' |
| email_verified | BOOLEAN | NOT NULL | false |
| first_name | VARCHAR(100) | NULL |  |
| last_name | VARCHAR(100) | NULL |  |
| version | INTEGER | NOT NULL | 1 |
| created_at | TIMESTAMP | NOT NULL | NOW() |
| updated_at | TIMESTAMP | NOT NULL | NOW() |

**Primary Key:** [id]

**Indexes:**
- `idx_customers_email`: [email]
- `idx_customers_registration_status`: [registration_status]

**Constraints:**
- `chk_customer_registration_status`: registration_status IN ('unverified', 'verified', 'suspended')

---

### TBL-CUSTOMER-002 – shipping_addresses {#tbl-customer-002}

**Aggregate:** Customer

**Purpose:** Stores customer saved shipping addresses (child of customers)

**Columns:**

| Name | Type | Constraints | Default |
|------|------|-------------|----------|
| id | UUID | NOT NULL | gen_random_uuid() |
| customer_id | UUID | NOT NULL |  |
| label | VARCHAR(50) | NULL |  |
| street | VARCHAR(200) | NOT NULL |  |
| city | VARCHAR(100) | NOT NULL |  |
| state | VARCHAR(50) | NOT NULL |  |
| postal_code | VARCHAR(20) | NOT NULL |  |
| country | CHAR(2) | NOT NULL |  |
| is_default | BOOLEAN | NOT NULL | false |
| created_at | TIMESTAMP | NOT NULL | NOW() |
| updated_at | TIMESTAMP | NOT NULL | NOW() |

**Primary Key:** [id]

**Indexes:**
- `idx_shipping_addresses_customer`: [customer_id]

**Foreign Keys:**
- [customer_id] → customers(id) (ON DELETE CASCADE)

---

### TBL-CART-001 – carts {#tbl-cart-001}

**Aggregate:** Cart

**Purpose:** Stores shopping cart header information

**Columns:**

| Name | Type | Constraints | Default |
|------|------|-------------|----------|
| id | UUID | NOT NULL | gen_random_uuid() |
| customer_id | UUID | NULL |  |
| is_guest | BOOLEAN | NOT NULL | false |
| total_price_amount | DECIMAL(12,2) | NOT NULL | 0.00 |
| total_price_currency | CHAR(3) | NOT NULL | 'USD' |
| last_activity_date | TIMESTAMP | NOT NULL | NOW() |
| version | INTEGER | NOT NULL | 1 |
| created_at | TIMESTAMP | NOT NULL | NOW() |
| updated_at | TIMESTAMP | NOT NULL | NOW() |

**Primary Key:** [id]

**Indexes:**
- `idx_carts_customer`: [customer_id]
- `idx_carts_last_activity`: [last_activity_date]

**Foreign Keys:**
- [customer_id] → customers(id) (ON DELETE SET NULL)

**Constraints:**
- `chk_cart_total_positive`: total_price_amount >= 0

---

### TBL-CART-002 – cart_items {#tbl-cart-002}

**Aggregate:** Cart

**Purpose:** Stores cart line items (child of carts)

**Columns:**

| Name | Type | Constraints | Default |
|------|------|-------------|----------|
| id | UUID | NOT NULL | gen_random_uuid() |
| cart_id | UUID | NOT NULL |  |
| product_id | UUID | NOT NULL |  |
| variant_id | UUID | NULL |  |
| product_name | VARCHAR(200) | NOT NULL |  |
| unit_price_amount | DECIMAL(12,2) | NOT NULL |  |
| unit_price_currency | CHAR(3) | NOT NULL | 'USD' |
| quantity | INTEGER | NOT NULL |  |
| subtotal_amount | DECIMAL(12,2) | NOT NULL |  |
| subtotal_currency | CHAR(3) | NOT NULL | 'USD' |
| created_at | TIMESTAMP | NOT NULL | NOW() |
| updated_at | TIMESTAMP | NOT NULL | NOW() |

**Primary Key:** [id]

**Indexes:**
- `idx_cart_items_cart`: [cart_id]
- `idx_cart_items_product`: [product_id]
- `uq_cart_items_product_variant`: [cart_id product_id variant_id]

**Foreign Keys:**
- [cart_id] → carts(id) (ON DELETE CASCADE)
- [product_id] → products(id) (ON DELETE RESTRICT)
- [variant_id] → product_variants(id) (ON DELETE RESTRICT)

**Constraints:**
- `chk_cart_item_quantity`: quantity > 0
- `chk_cart_item_price`: unit_price_amount >= 0 AND subtotal_amount >= 0

---

### TBL-ORDER-001 – orders {#tbl-order-001}

**Aggregate:** Order

**Purpose:** Stores order header information

**Columns:**

| Name | Type | Constraints | Default |
|------|------|-------------|----------|
| id | UUID | NOT NULL | gen_random_uuid() |
| order_number | VARCHAR(20) | NOT NULL UNIQUE |  |
| customer_id | UUID | NOT NULL |  |
| status | VARCHAR(20) | NOT NULL | 'pending' |
| shipping_street | VARCHAR(200) | NOT NULL |  |
| shipping_city | VARCHAR(100) | NOT NULL |  |
| shipping_state | VARCHAR(50) | NOT NULL |  |
| shipping_postal_code | VARCHAR(20) | NOT NULL |  |
| shipping_country | CHAR(2) | NOT NULL |  |
| payment_type | VARCHAR(20) | NOT NULL |  |
| payment_last_four_digits | CHAR(4) | NULL |  |
| payment_paypal_email | VARCHAR(255) | NULL |  |
| payment_transaction_id | VARCHAR(100) | NULL |  |
| subtotal_amount | DECIMAL(12,2) | NOT NULL |  |
| subtotal_currency | CHAR(3) | NOT NULL | 'USD' |
| shipping_cost_amount | DECIMAL(12,2) | NOT NULL |  |
| shipping_cost_currency | CHAR(3) | NOT NULL | 'USD' |
| tax_amount | DECIMAL(12,2) | NOT NULL | 0.00 |
| tax_currency | CHAR(3) | NOT NULL | 'USD' |
| total_amount | DECIMAL(12,2) | NOT NULL |  |
| total_currency | CHAR(3) | NOT NULL | 'USD' |
| tracking_number | VARCHAR(100) | NULL |  |
| cancellation_reason | VARCHAR(500) | NULL |  |
| version | INTEGER | NOT NULL | 1 |
| created_at | TIMESTAMP | NOT NULL | NOW() |
| updated_at | TIMESTAMP | NOT NULL | NOW() |

**Primary Key:** [id]

**Indexes:**
- `idx_orders_order_number`: [order_number]
- `idx_orders_customer`: [customer_id]
- `idx_orders_status`: [status]
- `idx_orders_created`: [created_at]

**Foreign Keys:**
- [customer_id] → customers(id) (ON DELETE RESTRICT)

**Constraints:**
- `chk_order_status`: status IN ('pending', 'confirmed', 'shipped', 'delivered', 'cancelled')
- `chk_order_amounts`: total_amount >= 0 AND subtotal_amount >= 0 AND shipping_cost_amount >= 0 AND tax_amount >= 0
- `chk_order_payment_type`: payment_type IN ('credit_card', 'paypal')

---

### TBL-ORDER-002 – order_line_items {#tbl-order-002}

**Aggregate:** Order

**Purpose:** Stores order line items (child of orders)

**Columns:**

| Name | Type | Constraints | Default |
|------|------|-------------|----------|
| id | UUID | NOT NULL | gen_random_uuid() |
| order_id | UUID | NOT NULL |  |
| product_id | UUID | NOT NULL |  |
| variant_id | UUID | NULL |  |
| product_name | VARCHAR(200) | NOT NULL |  |
| sku | VARCHAR(50) | NULL |  |
| unit_price_amount | DECIMAL(12,2) | NOT NULL |  |
| unit_price_currency | CHAR(3) | NOT NULL | 'USD' |
| quantity | INTEGER | NOT NULL |  |
| subtotal_amount | DECIMAL(12,2) | NOT NULL |  |
| subtotal_currency | CHAR(3) | NOT NULL | 'USD' |
| created_at | TIMESTAMP | NOT NULL | NOW() |

**Primary Key:** [id]

**Indexes:**
- `idx_order_line_items_order`: [order_id]
- `idx_order_line_items_product`: [product_id]

**Foreign Keys:**
- [order_id] → orders(id) (ON DELETE CASCADE)

**Constraints:**
- `chk_line_item_quantity`: quantity > 0
- `chk_line_item_price`: unit_price_amount >= 0 AND subtotal_amount >= 0

---

### TBL-PRODUCT-001 – products {#tbl-product-001}

**Aggregate:** Product

**Purpose:** Stores product catalog information

**Columns:**

| Name | Type | Constraints | Default |
|------|------|-------------|----------|
| id | UUID | NOT NULL | gen_random_uuid() |
| name | VARCHAR(200) | NOT NULL |  |
| description | TEXT | NULL |  |
| price_amount | DECIMAL(12,2) | NOT NULL |  |
| price_currency | CHAR(3) | NOT NULL | 'USD' |
| category_id | UUID | NOT NULL |  |
| is_active | BOOLEAN | NOT NULL | true |
| is_deleted | BOOLEAN | NOT NULL | false |
| status | VARCHAR(20) | NOT NULL | 'draft' |
| version | INTEGER | NOT NULL | 1 |
| created_by | UUID | NULL |  |
| updated_by | UUID | NULL |  |
| created_at | TIMESTAMP | NOT NULL | NOW() |
| updated_at | TIMESTAMP | NOT NULL | NOW() |

**Primary Key:** [id]

**Indexes:**
- `idx_products_category`: [category_id]
- `idx_products_status`: [status]
- `idx_products_active`: [is_active is_deleted]
- `idx_products_name`: [name]

**Foreign Keys:**
- [category_id] → categories(id) (ON DELETE RESTRICT)

**Constraints:**
- `chk_product_price`: price_amount >= 0.01
- `chk_product_name_length`: LENGTH(name) >= 2 AND LENGTH(name) <= 200
- `chk_product_status`: status IN ('draft', 'active', 'archived')

---

### TBL-PRODUCT-002 – product_variants {#tbl-product-002}

**Aggregate:** Product

**Purpose:** Stores product variant information (child of products)

**Columns:**

| Name | Type | Constraints | Default |
|------|------|-------------|----------|
| id | UUID | NOT NULL | gen_random_uuid() |
| product_id | UUID | NOT NULL |  |
| sku | VARCHAR(50) | NOT NULL UNIQUE |  |
| size | VARCHAR(20) | NULL |  |
| color | VARCHAR(50) | NULL |  |
| price_adjustment_amount | DECIMAL(12,2) | NOT NULL | 0.00 |
| price_adjustment_currency | CHAR(3) | NOT NULL | 'USD' |
| created_at | TIMESTAMP | NOT NULL | NOW() |
| updated_at | TIMESTAMP | NOT NULL | NOW() |

**Primary Key:** [id]

**Indexes:**
- `idx_product_variants_product`: [product_id]
- `idx_product_variants_sku`: [sku]
- `uq_product_variant_combo`: [product_id size color]

**Foreign Keys:**
- [product_id] → products(id) (ON DELETE CASCADE)

**Constraints:**
- `chk_variant_has_attribute`: size IS NOT NULL OR color IS NOT NULL

---

### TBL-PRODUCT-003 – product_images {#tbl-product-003}

**Aggregate:** Product

**Purpose:** Stores product image references (child of products)

**Columns:**

| Name | Type | Constraints | Default |
|------|------|-------------|----------|
| id | UUID | NOT NULL | gen_random_uuid() |
| product_id | UUID | NOT NULL |  |
| url | VARCHAR(500) | NOT NULL |  |
| alt_text | VARCHAR(200) | NULL |  |
| display_order | INTEGER | NOT NULL | 0 |
| is_primary | BOOLEAN | NOT NULL | false |
| created_at | TIMESTAMP | NOT NULL | NOW() |

**Primary Key:** [id]

**Indexes:**
- `idx_product_images_product`: [product_id]
- `idx_product_images_order`: [product_id display_order]

**Foreign Keys:**
- [product_id] → products(id) (ON DELETE CASCADE)

**Constraints:**
- `chk_image_display_order`: display_order >= 0

---

### TBL-INVENTORY-001 – inventory {#tbl-inventory-001}

**Aggregate:** Inventory

**Purpose:** Stores stock levels and reservations for products

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

### TBL-INVENTORY-002 – inventory_audit_logs {#tbl-inventory-002}

**Aggregate:** Inventory

**Purpose:** Stores immutable inventory change audit trail (child of inventory)

**Columns:**

| Name | Type | Constraints | Default |
|------|------|-------------|----------|
| id | UUID | NOT NULL | gen_random_uuid() |
| inventory_id | UUID | NOT NULL |  |
| action | VARCHAR(30) | NOT NULL |  |
| previous_value | INTEGER | NOT NULL |  |
| new_value | INTEGER | NOT NULL |  |
| quantity_changed | INTEGER | NOT NULL |  |
| reason | VARCHAR(500) | NOT NULL |  |
| order_id | UUID | NULL |  |
| user_id | UUID | NULL |  |
| created_at | TIMESTAMP | NOT NULL | NOW() |

**Primary Key:** [id]

**Indexes:**
- `idx_inventory_audit_inventory`: [inventory_id]
- `idx_inventory_audit_created`: [created_at]
- `idx_inventory_audit_order`: [order_id]

**Foreign Keys:**
- [inventory_id] → inventory(id) (ON DELETE CASCADE)

**Constraints:**
- `chk_audit_action`: action IN ('set', 'adjust', 'reserve', 'release', 'deduct', 'restock')
- `chk_audit_values`: previous_value >= 0 AND new_value >= 0

---

### TBL-CATEGORY-001 – categories {#tbl-category-001}

**Aggregate:** Category

**Purpose:** Stores product category hierarchy

**Columns:**

| Name | Type | Constraints | Default |
|------|------|-------------|----------|
| id | UUID | NOT NULL | gen_random_uuid() |
| name | VARCHAR(100) | NOT NULL UNIQUE |  |
| description | TEXT | NULL |  |
| parent_category_id | UUID | NULL |  |
| depth | INTEGER | NOT NULL | 0 |
| display_order | INTEGER | NOT NULL | 0 |
| is_active | BOOLEAN | NOT NULL | true |
| version | INTEGER | NOT NULL | 1 |
| created_at | TIMESTAMP | NOT NULL | NOW() |
| updated_at | TIMESTAMP | NOT NULL | NOW() |

**Primary Key:** [id]

**Indexes:**
- `idx_categories_name`: [name]
- `idx_categories_parent`: [parent_category_id]
- `idx_categories_depth`: [depth]
- `idx_categories_display_order`: [display_order]

**Foreign Keys:**
- [parent_category_id] → categories(id) (ON DELETE RESTRICT)

**Constraints:**
- `chk_category_name_length`: LENGTH(name) >= 1 AND LENGTH(name) <= 100
- `chk_category_depth`: depth >= 0 AND depth <= 3
- `chk_category_no_self_parent`: parent_category_id IS NULL OR parent_category_id <> id

---

## Entity-Relationship Diagram

```mermaid
erDiagram
    customers {
        UUID id
        VARCHAR(255) email
        VARCHAR(255) password_hash
        VARCHAR(20) registration_status
        BOOLEAN email_verified
        VARCHAR(100) first_name
        VARCHAR(100) last_name
        INTEGER version
        TIMESTAMP created_at
        TIMESTAMP updated_at
    }
    shipping_addresses {
        UUID id
        UUID customer_id
        VARCHAR(50) label
        VARCHAR(200) street
        VARCHAR(100) city
        VARCHAR(50) state
        VARCHAR(20) postal_code
        CHAR(2) country
        BOOLEAN is_default
        TIMESTAMP created_at
        TIMESTAMP updated_at
    }
    carts {
        UUID id
        UUID customer_id
        BOOLEAN is_guest
        DECIMAL(12,2) total_price_amount
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
        DECIMAL(12,2) unit_price_amount
        CHAR(3) unit_price_currency
        INTEGER quantity
        DECIMAL(12,2) subtotal_amount
        CHAR(3) subtotal_currency
        TIMESTAMP created_at
        TIMESTAMP updated_at
    }
    orders {
        UUID id
        VARCHAR(20) order_number
        UUID customer_id
        VARCHAR(20) status
        VARCHAR(200) shipping_street
        VARCHAR(100) shipping_city
        VARCHAR(50) shipping_state
        VARCHAR(20) shipping_postal_code
        CHAR(2) shipping_country
        VARCHAR(20) payment_type
        CHAR(4) payment_last_four_digits
        VARCHAR(255) payment_paypal_email
        VARCHAR(100) payment_transaction_id
        DECIMAL(12,2) subtotal_amount
        CHAR(3) subtotal_currency
        DECIMAL(12,2) shipping_cost_amount
        CHAR(3) shipping_cost_currency
        DECIMAL(12,2) tax_amount
        CHAR(3) tax_currency
        DECIMAL(12,2) total_amount
        CHAR(3) total_currency
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
        VARCHAR(50) sku
        DECIMAL(12,2) unit_price_amount
        CHAR(3) unit_price_currency
        INTEGER quantity
        DECIMAL(12,2) subtotal_amount
        CHAR(3) subtotal_currency
        TIMESTAMP created_at
    }
    products {
        UUID id
        VARCHAR(200) name
        TEXT description
        DECIMAL(12,2) price_amount
        CHAR(3) price_currency
        UUID category_id
        BOOLEAN is_active
        BOOLEAN is_deleted
        VARCHAR(20) status
        INTEGER version
        UUID created_by
        UUID updated_by
        TIMESTAMP created_at
        TIMESTAMP updated_at
    }
    product_variants {
        UUID id
        UUID product_id
        VARCHAR(50) sku
        VARCHAR(20) size
        VARCHAR(50) color
        DECIMAL(12,2) price_adjustment_amount
        CHAR(3) price_adjustment_currency
        TIMESTAMP created_at
        TIMESTAMP updated_at
    }
    product_images {
        UUID id
        UUID product_id
        VARCHAR(500) url
        VARCHAR(200) alt_text
        INTEGER display_order
        BOOLEAN is_primary
        TIMESTAMP created_at
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
    inventory_audit_logs {
        UUID id
        UUID inventory_id
        VARCHAR(30) action
        INTEGER previous_value
        INTEGER new_value
        INTEGER quantity_changed
        VARCHAR(500) reason
        UUID order_id
        UUID user_id
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
    customers ||--o{ shipping_addresses : "FK"
    customers ||--o{ carts : "FK"
    carts ||--o{ cart_items : "FK"
    products ||--o{ cart_items : "FK"
    product_variants ||--o{ cart_items : "FK"
    customers ||--o{ orders : "FK"
    orders ||--o{ order_line_items : "FK"
    categories ||--o{ products : "FK"
    products ||--o{ product_variants : "FK"
    products ||--o{ product_images : "FK"
    products ||--o{ inventory : "FK"
    inventory ||--o{ inventory_audit_logs : "FK"
    categories ||--o{ categories : "FK"
```
