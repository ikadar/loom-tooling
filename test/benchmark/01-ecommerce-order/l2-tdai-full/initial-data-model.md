# Initial Data Model

Generated: 2025-12-31T12:00:34+01:00

---

## Enumerations

### registration_status

Values: `[registered unregistered]`

### order_status

Values: `[pending confirmed shipped delivered cancelled]`

### payment_type

Values: `[credit_card paypal]`

---

## Tables

### TBL-CUSTOMER-001 – customers

**Aggregate:** Customer

**Purpose:** Stores registered customer information

**Columns:**

| Name | Type | Constraints | Default |
|------|------|-------------|----------|
| id | UUID | NOT NULL | gen_random_uuid() |
| email | VARCHAR(255) | NOT NULL UNIQUE |  |
| registration_status | VARCHAR(20) | NOT NULL | 'registered' |
| version | INTEGER | NOT NULL | 1 |
| created_at | TIMESTAMP | NOT NULL | NOW() |
| updated_at | TIMESTAMP | NOT NULL | NOW() |

**Primary Key:** [id]

**Indexes:**
- `idx_customers_email`: [email]

**Constraints:**
- `chk_customer_registration_status`: registration_status IN ('registered', 'unregistered')

---

### TBL-CART-001 – carts

**Aggregate:** Cart

**Purpose:** Stores shopping cart header information

**Columns:**

| Name | Type | Constraints | Default |
|------|------|-------------|----------|
| id | UUID | NOT NULL | gen_random_uuid() |
| customer_id | UUID | NOT NULL UNIQUE |  |
| total_price_amount | DECIMAL(10,2) | NOT NULL | 0 |
| total_price_currency | CHAR(3) | NOT NULL | 'USD' |
| version | INTEGER | NOT NULL | 1 |
| created_at | TIMESTAMP | NOT NULL | NOW() |
| updated_at | TIMESTAMP | NOT NULL | NOW() |

**Primary Key:** [id]

**Indexes:**
- `idx_carts_customer`: [customer_id]

**Foreign Keys:**
- [customer_id] → customers(id) (ON DELETE CASCADE)

**Constraints:**
- `chk_cart_total_price`: total_price_amount >= 0

---

### TBL-CART-002 – cart_items

**Aggregate:** Cart

**Purpose:** Stores items within shopping carts

**Columns:**

| Name | Type | Constraints | Default |
|------|------|-------------|----------|
| id | UUID | NOT NULL | gen_random_uuid() |
| cart_id | UUID | NOT NULL |  |
| product_id | UUID | NOT NULL |  |
| product_name | VARCHAR(200) | NOT NULL |  |
| unit_price_amount | DECIMAL(10,2) | NOT NULL |  |
| unit_price_currency | CHAR(3) | NOT NULL | 'USD' |
| quantity | INTEGER | NOT NULL |  |
| subtotal_amount | DECIMAL(10,2) | NOT NULL |  |
| subtotal_currency | CHAR(3) | NOT NULL | 'USD' |
| created_at | TIMESTAMP | NOT NULL | NOW() |
| updated_at | TIMESTAMP | NOT NULL | NOW() |

**Primary Key:** [id]

**Indexes:**
- `idx_cart_items_cart`: [cart_id]
- `idx_cart_items_cart_product`: [cart_id product_id]

**Foreign Keys:**
- [cart_id] → carts(id) (ON DELETE CASCADE)
- [product_id] → products(id) (ON DELETE RESTRICT)

**Constraints:**
- `chk_cart_item_quantity`: quantity > 0
- `chk_cart_item_prices`: unit_price_amount >= 0 AND subtotal_amount >= 0

---

### TBL-ORDER-001 – orders

**Aggregate:** Order

**Purpose:** Stores order header information with shipping and payment details

**Columns:**

| Name | Type | Constraints | Default |
|------|------|-------------|----------|
| id | UUID | NOT NULL | gen_random_uuid() |
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
| subtotal_amount | DECIMAL(10,2) | NOT NULL |  |
| subtotal_currency | CHAR(3) | NOT NULL | 'USD' |
| shipping_cost_amount | DECIMAL(10,2) | NOT NULL |  |
| shipping_cost_currency | CHAR(3) | NOT NULL | 'USD' |
| total_amount | DECIMAL(10,2) | NOT NULL |  |
| total_currency | CHAR(3) | NOT NULL | 'USD' |
| version | INTEGER | NOT NULL | 1 |
| created_at | TIMESTAMP | NOT NULL | NOW() |
| updated_at | TIMESTAMP | NOT NULL | NOW() |

**Primary Key:** [id]

**Indexes:**
- `idx_orders_customer`: [customer_id]
- `idx_orders_status`: [status]
- `idx_orders_created`: [created_at]

**Foreign Keys:**
- [customer_id] → customers(id) (ON DELETE RESTRICT)

**Constraints:**
- `chk_order_status`: status IN ('pending', 'confirmed', 'shipped', 'delivered', 'cancelled')
- `chk_order_amounts`: total_amount >= 0 AND subtotal_amount >= 0 AND shipping_cost_amount >= 0
- `chk_order_payment_type`: payment_type IN ('credit_card', 'paypal')

---

### TBL-ORDER-002 – order_line_items

**Aggregate:** Order

**Purpose:** Stores immutable order line items (snapshot at order time)

**Columns:**

| Name | Type | Constraints | Default |
|------|------|-------------|----------|
| id | UUID | NOT NULL | gen_random_uuid() |
| order_id | UUID | NOT NULL |  |
| product_id | UUID | NOT NULL |  |
| product_name | VARCHAR(200) | NOT NULL |  |
| unit_price_amount | DECIMAL(10,2) | NOT NULL |  |
| unit_price_currency | CHAR(3) | NOT NULL | 'USD' |
| quantity | INTEGER | NOT NULL |  |
| subtotal_amount | DECIMAL(10,2) | NOT NULL |  |
| subtotal_currency | CHAR(3) | NOT NULL | 'USD' |
| created_at | TIMESTAMP | NOT NULL | NOW() |

**Primary Key:** [id]

**Indexes:**
- `idx_order_line_items_order`: [order_id]

**Foreign Keys:**
- [order_id] → orders(id) (ON DELETE CASCADE)

**Constraints:**
- `chk_order_line_item_quantity`: quantity > 0
- `chk_order_line_item_prices`: unit_price_amount >= 0 AND subtotal_amount >= 0

---

### TBL-PRODUCT-001 – products

**Aggregate:** Product

**Purpose:** Stores product catalog information

**Columns:**

| Name | Type | Constraints | Default |
|------|------|-------------|----------|
| id | UUID | NOT NULL | gen_random_uuid() |
| name | VARCHAR(200) | NOT NULL |  |
| description | VARCHAR(2000) | NULL |  |
| price_amount | DECIMAL(10,2) | NOT NULL |  |
| price_currency | CHAR(3) | NOT NULL | 'USD' |
| category_id | UUID | NOT NULL |  |
| is_active | BOOLEAN | NOT NULL | true |
| version | INTEGER | NOT NULL | 1 |
| created_at | TIMESTAMP | NOT NULL | NOW() |
| updated_at | TIMESTAMP | NOT NULL | NOW() |

**Primary Key:** [id]

**Indexes:**
- `idx_products_category`: [category_id]
- `idx_products_active`: [is_active]
- `idx_products_name`: [name]

**Foreign Keys:**
- [category_id] → categories(id) (ON DELETE RESTRICT)

**Constraints:**
- `chk_product_price`: price_amount > 0
- `chk_product_name_length`: LENGTH(name) >= 1 AND LENGTH(name) <= 200

---

### TBL-PRODUCT-002 – product_variants

**Aggregate:** Product

**Purpose:** Stores product variations (size, color combinations)

**Columns:**

| Name | Type | Constraints | Default |
|------|------|-------------|----------|
| id | UUID | NOT NULL | gen_random_uuid() |
| product_id | UUID | NOT NULL |  |
| size | VARCHAR(50) | NULL |  |
| color | VARCHAR(50) | NULL |  |
| sku | VARCHAR(100) | NOT NULL UNIQUE |  |
| price_adjustment_amount | DECIMAL(10,2) | NULL |  |
| price_adjustment_currency | CHAR(3) | NULL |  |
| created_at | TIMESTAMP | NOT NULL | NOW() |
| updated_at | TIMESTAMP | NOT NULL | NOW() |

**Primary Key:** [id]

**Indexes:**
- `idx_product_variants_product`: [product_id]
- `idx_product_variants_sku`: [sku]
- `idx_product_variants_unique_combo`: [product_id size color]

**Foreign Keys:**
- [product_id] → products(id) (ON DELETE CASCADE)

**Constraints:**
- `chk_variant_has_attribute`: size IS NOT NULL OR color IS NOT NULL

---

### TBL-CATEGORY-001 – categories

**Aggregate:** Category

**Purpose:** Stores product category hierarchy

**Columns:**

| Name | Type | Constraints | Default |
|------|------|-------------|----------|
| id | UUID | NOT NULL | gen_random_uuid() |
| name | VARCHAR(100) | NOT NULL UNIQUE |  |
| description | VARCHAR(500) | NULL |  |
| parent_category_id | UUID | NULL |  |
| version | INTEGER | NOT NULL | 1 |
| created_at | TIMESTAMP | NOT NULL | NOW() |
| updated_at | TIMESTAMP | NOT NULL | NOW() |

**Primary Key:** [id]

**Indexes:**
- `idx_categories_name`: [name]
- `idx_categories_parent`: [parent_category_id]

**Foreign Keys:**
- [parent_category_id] → categories(id) (ON DELETE SET NULL)

**Constraints:**
- `chk_category_not_self_parent`: parent_category_id IS NULL OR parent_category_id != id
- `chk_category_name_length`: LENGTH(name) >= 1 AND LENGTH(name) <= 100

---

### TBL-INVENTORY-001 – inventory

**Aggregate:** Inventory

**Purpose:** Tracks stock levels and reservations for products

**Columns:**

| Name | Type | Constraints | Default |
|------|------|-------------|----------|
| id | UUID | NOT NULL | gen_random_uuid() |
| product_id | UUID | NOT NULL UNIQUE |  |
| stock_level | INTEGER | NOT NULL | 0 |
| reserved_quantity | INTEGER | NOT NULL | 0 |
| version | INTEGER | NOT NULL | 1 |
| created_at | TIMESTAMP | NOT NULL | NOW() |
| updated_at | TIMESTAMP | NOT NULL | NOW() |

**Primary Key:** [id]

**Indexes:**
- `idx_inventory_product`: [product_id]
- `idx_inventory_stock_level`: [stock_level]

**Foreign Keys:**
- [product_id] → products(id) (ON DELETE CASCADE)

**Constraints:**
- `chk_inventory_stock_level`: stock_level >= 0
- `chk_inventory_reserved_quantity`: reserved_quantity >= 0
- `chk_inventory_reserved_not_exceed_stock`: reserved_quantity <= stock_level

---

## Entity-Relationship Diagram

```mermaid
erDiagram
    customers {
        UUID id
        VARCHAR(255) email
        VARCHAR(20) registration_status
        INTEGER version
        TIMESTAMP created_at
        TIMESTAMP updated_at
    }
    carts {
        UUID id
        UUID customer_id
        DECIMAL(10,2) total_price_amount
        CHAR(3) total_price_currency
        INTEGER version
        TIMESTAMP created_at
        TIMESTAMP updated_at
    }
    cart_items {
        UUID id
        UUID cart_id
        UUID product_id
        VARCHAR(200) product_name
        DECIMAL(10,2) unit_price_amount
        CHAR(3) unit_price_currency
        INTEGER quantity
        DECIMAL(10,2) subtotal_amount
        CHAR(3) subtotal_currency
        TIMESTAMP created_at
        TIMESTAMP updated_at
    }
    orders {
        UUID id
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
        DECIMAL(10,2) subtotal_amount
        CHAR(3) subtotal_currency
        DECIMAL(10,2) shipping_cost_amount
        CHAR(3) shipping_cost_currency
        DECIMAL(10,2) total_amount
        CHAR(3) total_currency
        INTEGER version
        TIMESTAMP created_at
        TIMESTAMP updated_at
    }
    order_line_items {
        UUID id
        UUID order_id
        UUID product_id
        VARCHAR(200) product_name
        DECIMAL(10,2) unit_price_amount
        CHAR(3) unit_price_currency
        INTEGER quantity
        DECIMAL(10,2) subtotal_amount
        CHAR(3) subtotal_currency
        TIMESTAMP created_at
    }
    products {
        UUID id
        VARCHAR(200) name
        VARCHAR(2000) description
        DECIMAL(10,2) price_amount
        CHAR(3) price_currency
        UUID category_id
        BOOLEAN is_active
        INTEGER version
        TIMESTAMP created_at
        TIMESTAMP updated_at
    }
    product_variants {
        UUID id
        UUID product_id
        VARCHAR(50) size
        VARCHAR(50) color
        VARCHAR(100) sku
        DECIMAL(10,2) price_adjustment_amount
        CHAR(3) price_adjustment_currency
        TIMESTAMP created_at
        TIMESTAMP updated_at
    }
    categories {
        UUID id
        VARCHAR(100) name
        VARCHAR(500) description
        UUID parent_category_id
        INTEGER version
        TIMESTAMP created_at
        TIMESTAMP updated_at
    }
    inventory {
        UUID id
        UUID product_id
        INTEGER stock_level
        INTEGER reserved_quantity
        INTEGER version
        TIMESTAMP created_at
        TIMESTAMP updated_at
    }
    customers ||--o{ carts : "FK"
    carts ||--o{ cart_items : "FK"
    products ||--o{ cart_items : "FK"
    customers ||--o{ orders : "FK"
    orders ||--o{ order_line_items : "FK"
    categories ||--o{ products : "FK"
    products ||--o{ product_variants : "FK"
    categories ||--o{ categories : "FK"
    products ||--o{ inventory : "FK"
```
