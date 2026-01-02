# Initial Data Model

Generated: 2026-01-02T21:09:41+01:00

---

## Enumerations

### registration_status_enum

Values: `[REGISTERED UNREGISTERED]`

### order_status_enum

Values: `[PENDING CONFIRMED SHIPPED DELIVERED CANCELLED]`

### payment_type_enum

Values: `[CREDIT_CARD PAYPAL]`

---

## Tables

### TBL-CUSTOMER-001 – customers {#tbl-customer-001}

**Aggregate:** Customer

**Purpose:** Stores registered customers who can browse products and place orders

**Columns:**

| Name | Type | Constraints | Default |
|------|------|-------------|----------|
| id | UUID | PRIMARY KEY |  |
| email | VARCHAR(255) | NOT NULL UNIQUE |  |
| password_hash | VARCHAR(255) | NOT NULL |  |
| registration_status | registration_status_enum | NOT NULL DEFAULT 'REGISTERED' |  |
| version | INTEGER | NOT NULL DEFAULT 1 |  |
| created_at | TIMESTAMP | NOT NULL DEFAULT NOW() |  |
| updated_at | TIMESTAMP | NOT NULL DEFAULT NOW() |  |

**Primary Key:** [id]

**Indexes:**
- `idx_customers_email`: [email]

---

### TBL-CART-001 – carts {#tbl-cart-001}

**Aggregate:** Cart

**Purpose:** Holds products a customer intends to purchase before checkout

**Columns:**

| Name | Type | Constraints | Default |
|------|------|-------------|----------|
| id | UUID | PRIMARY KEY |  |
| customer_id | UUID | NOT NULL UNIQUE |  |
| total_price_amount | DECIMAL(12,2) | NOT NULL DEFAULT 0 |  |
| total_price_currency | CHAR(3) | NOT NULL DEFAULT 'USD' |  |
| version | INTEGER | NOT NULL DEFAULT 1 |  |
| created_at | TIMESTAMP | NOT NULL DEFAULT NOW() |  |
| updated_at | TIMESTAMP | NOT NULL DEFAULT NOW() |  |

**Primary Key:** [id]

**Indexes:**
- `idx_carts_customer_id`: [customer_id]

**Foreign Keys:**
- [customer_id] → customers(id) (ON DELETE CASCADE)

**Constraints:**
- `chk_cart_total_price_non_negative`: total_price_amount >= 0

---

### TBL-CARTITEM-001 – cart_items {#tbl-cartitem-001}

**Aggregate:** Cart

**Purpose:** Represents a product and quantity within a shopping cart

**Columns:**

| Name | Type | Constraints | Default |
|------|------|-------------|----------|
| id | UUID | PRIMARY KEY |  |
| cart_id | UUID | NOT NULL |  |
| product_id | UUID | NOT NULL |  |
| product_name | VARCHAR(200) | NOT NULL |  |
| unit_price_amount | DECIMAL(12,2) | NOT NULL |  |
| unit_price_currency | CHAR(3) | NOT NULL DEFAULT 'USD' |  |
| quantity | INTEGER | NOT NULL |  |
| created_at | TIMESTAMP | NOT NULL DEFAULT NOW() |  |
| updated_at | TIMESTAMP | NOT NULL DEFAULT NOW() |  |

**Primary Key:** [id]

**Indexes:**
- `idx_cart_items_cart_id`: [cart_id]
- `idx_cart_items_cart_product`: [cart_id product_id]

**Foreign Keys:**
- [cart_id] → carts(id) (ON DELETE CASCADE)
- [product_id] → products(id) (ON DELETE RESTRICT)

**Constraints:**
- `chk_cart_item_quantity_positive`: quantity >= 1
- `chk_cart_item_price_positive`: unit_price_amount >= 0

---

### TBL-ORDER-001 – orders {#tbl-order-001}

**Aggregate:** Order

**Purpose:** Represents a customer's purchase transaction with shipping and payment details

**Columns:**

| Name | Type | Constraints | Default |
|------|------|-------------|----------|
| id | UUID | PRIMARY KEY |  |
| customer_id | UUID | NOT NULL |  |
| status | order_status_enum | NOT NULL DEFAULT 'PENDING' |  |
| shipping_street | VARCHAR(200) | NOT NULL |  |
| shipping_city | VARCHAR(100) | NOT NULL |  |
| shipping_state | VARCHAR(100) | NOT NULL |  |
| shipping_postal_code | VARCHAR(20) | NOT NULL |  |
| shipping_country | CHAR(2) | NOT NULL |  |
| payment_type | payment_type_enum | NOT NULL |  |
| payment_last_four_digits | CHAR(4) | NULL |  |
| payment_expiry_month | SMALLINT | NULL |  |
| payment_expiry_year | SMALLINT | NULL |  |
| payment_paypal_email | VARCHAR(255) | NULL |  |
| subtotal_amount | DECIMAL(12,2) | NOT NULL |  |
| subtotal_currency | CHAR(3) | NOT NULL DEFAULT 'USD' |  |
| shipping_cost_amount | DECIMAL(12,2) | NOT NULL |  |
| shipping_cost_currency | CHAR(3) | NOT NULL DEFAULT 'USD' |  |
| total_amount | DECIMAL(12,2) | NOT NULL |  |
| total_currency | CHAR(3) | NOT NULL DEFAULT 'USD' |  |
| tracking_number | VARCHAR(100) | NULL |  |
| cancellation_reason | TEXT | NULL |  |
| version | INTEGER | NOT NULL DEFAULT 1 |  |
| created_at | TIMESTAMP | NOT NULL DEFAULT NOW() |  |
| updated_at | TIMESTAMP | NOT NULL DEFAULT NOW() |  |

**Primary Key:** [id]

**Indexes:**
- `idx_orders_customer_id`: [customer_id]
- `idx_orders_status`: [status]
- `idx_orders_created_at`: [created_at]

**Foreign Keys:**
- [customer_id] → customers(id) (ON DELETE RESTRICT)

**Constraints:**
- `chk_order_subtotal_non_negative`: subtotal_amount >= 0
- `chk_order_shipping_non_negative`: shipping_cost_amount >= 0
- `chk_order_total_non_negative`: total_amount >= 0
- `chk_payment_expiry_month`: payment_expiry_month IS NULL OR (payment_expiry_month >= 1 AND payment_expiry_month <= 12)

---

### TBL-ORDERLINEITEM-001 – order_line_items {#tbl-orderlineitem-001}

**Aggregate:** Order

**Purpose:** Immutable snapshot of a product at time of order with quantity and pricing

**Columns:**

| Name | Type | Constraints | Default |
|------|------|-------------|----------|
| id | UUID | PRIMARY KEY |  |
| order_id | UUID | NOT NULL |  |
| product_id | UUID | NOT NULL |  |
| product_name | VARCHAR(200) | NOT NULL |  |
| unit_price_amount | DECIMAL(12,2) | NOT NULL |  |
| unit_price_currency | CHAR(3) | NOT NULL DEFAULT 'USD' |  |
| quantity | INTEGER | NOT NULL |  |
| created_at | TIMESTAMP | NOT NULL DEFAULT NOW() |  |
| updated_at | TIMESTAMP | NOT NULL DEFAULT NOW() |  |

**Primary Key:** [id]

**Indexes:**
- `idx_order_line_items_order_id`: [order_id]

**Foreign Keys:**
- [order_id] → orders(id) (ON DELETE CASCADE)
- [product_id] → products(id) (ON DELETE RESTRICT)

**Constraints:**
- `chk_line_item_quantity_positive`: quantity >= 1
- `chk_line_item_price_positive`: unit_price_amount >= 0

---

### TBL-PRODUCT-001 – products {#tbl-product-001}

**Aggregate:** Product

**Purpose:** Represents an item available for sale in the store

**Columns:**

| Name | Type | Constraints | Default |
|------|------|-------------|----------|
| id | UUID | PRIMARY KEY |  |
| name | VARCHAR(200) | NOT NULL |  |
| description | VARCHAR(5000) | NULL |  |
| price_amount | DECIMAL(12,2) | NOT NULL |  |
| price_currency | CHAR(3) | NOT NULL DEFAULT 'USD' |  |
| category_id | UUID | NOT NULL |  |
| is_active | BOOLEAN | NOT NULL DEFAULT TRUE |  |
| version | INTEGER | NOT NULL DEFAULT 1 |  |
| created_at | TIMESTAMP | NOT NULL DEFAULT NOW() |  |
| updated_at | TIMESTAMP | NOT NULL DEFAULT NOW() |  |

**Primary Key:** [id]

**Indexes:**
- `idx_products_category_id`: [category_id]
- `idx_products_is_active`: [is_active]
- `idx_products_name`: [name]

**Foreign Keys:**
- [category_id] → categories(id) (ON DELETE RESTRICT)

**Constraints:**
- `chk_product_price_positive`: price_amount > 0
- `chk_product_name_length`: LENGTH(name) >= 1

---

### TBL-PRODUCTVARIANT-001 – product_variants {#tbl-productvariant-001}

**Aggregate:** Product

**Purpose:** Represents a specific variation of a product (e.g., size, color combination)

**Columns:**

| Name | Type | Constraints | Default |
|------|------|-------------|----------|
| id | UUID | PRIMARY KEY |  |
| product_id | UUID | NOT NULL |  |
| sku | VARCHAR(100) | NOT NULL UNIQUE |  |
| size | VARCHAR(50) | NULL |  |
| color | VARCHAR(50) | NULL |  |
| price_adjustment_amount | DECIMAL(12,2) | NULL DEFAULT 0 |  |
| price_adjustment_currency | CHAR(3) | NOT NULL DEFAULT 'USD' |  |
| created_at | TIMESTAMP | NOT NULL DEFAULT NOW() |  |
| updated_at | TIMESTAMP | NOT NULL DEFAULT NOW() |  |

**Primary Key:** [id]

**Indexes:**
- `idx_product_variants_product_id`: [product_id]
- `idx_product_variants_sku`: [sku]

**Foreign Keys:**
- [product_id] → products(id) (ON DELETE CASCADE)

**Constraints:**
- `chk_variant_has_attribute`: size IS NOT NULL OR color IS NOT NULL

---

### TBL-CATEGORY-001 – categories {#tbl-category-001}

**Aggregate:** Category

**Purpose:** Organizes products into browsable groupings

**Columns:**

| Name | Type | Constraints | Default |
|------|------|-------------|----------|
| id | UUID | PRIMARY KEY |  |
| name | VARCHAR(100) | NOT NULL UNIQUE |  |
| description | TEXT | NULL |  |
| parent_category_id | UUID | NULL |  |
| version | INTEGER | NOT NULL DEFAULT 1 |  |
| created_at | TIMESTAMP | NOT NULL DEFAULT NOW() |  |
| updated_at | TIMESTAMP | NOT NULL DEFAULT NOW() |  |

**Primary Key:** [id]

**Indexes:**
- `idx_categories_name`: [name]
- `idx_categories_parent_id`: [parent_category_id]

**Foreign Keys:**
- [parent_category_id] → categories(id) (ON DELETE RESTRICT)

**Constraints:**
- `chk_category_not_self_parent`: id != parent_category_id
- `chk_category_name_length`: LENGTH(name) >= 1

---

### TBL-INVENTORY-001 – inventory {#tbl-inventory-001}

**Aggregate:** Inventory

**Purpose:** Tracks stock levels for products and manages availability

**Columns:**

| Name | Type | Constraints | Default |
|------|------|-------------|----------|
| id | UUID | PRIMARY KEY |  |
| product_id | UUID | NOT NULL UNIQUE |  |
| quantity_on_hand | INTEGER | NOT NULL DEFAULT 0 |  |
| reserved_quantity | INTEGER | NOT NULL DEFAULT 0 |  |
| version | INTEGER | NOT NULL DEFAULT 1 |  |
| created_at | TIMESTAMP | NOT NULL DEFAULT NOW() |  |
| updated_at | TIMESTAMP | NOT NULL DEFAULT NOW() |  |

**Primary Key:** [id]

**Indexes:**
- `idx_inventory_product_id`: [product_id]

**Foreign Keys:**
- [product_id] → products(id) (ON DELETE CASCADE)

**Constraints:**
- `chk_inventory_quantity_non_negative`: quantity_on_hand >= 0
- `chk_inventory_reserved_non_negative`: reserved_quantity >= 0
- `chk_inventory_reserved_valid`: reserved_quantity <= quantity_on_hand

---

## Entity-Relationship Diagram

```mermaid
erDiagram
    customers {
        UUID id
        VARCHAR(255) email
        VARCHAR(255) password_hash
        registration_status_enum registration_status
        INTEGER version
        TIMESTAMP created_at
        TIMESTAMP updated_at
    }
    carts {
        UUID id
        UUID customer_id
        DECIMAL(12,2) total_price_amount
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
        DECIMAL(12,2) unit_price_amount
        CHAR(3) unit_price_currency
        INTEGER quantity
        TIMESTAMP created_at
        TIMESTAMP updated_at
    }
    orders {
        UUID id
        UUID customer_id
        order_status_enum status
        VARCHAR(200) shipping_street
        VARCHAR(100) shipping_city
        VARCHAR(100) shipping_state
        VARCHAR(20) shipping_postal_code
        CHAR(2) shipping_country
        payment_type_enum payment_type
        CHAR(4) payment_last_four_digits
        SMALLINT payment_expiry_month
        SMALLINT payment_expiry_year
        VARCHAR(255) payment_paypal_email
        DECIMAL(12,2) subtotal_amount
        CHAR(3) subtotal_currency
        DECIMAL(12,2) shipping_cost_amount
        CHAR(3) shipping_cost_currency
        DECIMAL(12,2) total_amount
        CHAR(3) total_currency
        VARCHAR(100) tracking_number
        TEXT cancellation_reason
        INTEGER version
        TIMESTAMP created_at
        TIMESTAMP updated_at
    }
    order_line_items {
        UUID id
        UUID order_id
        UUID product_id
        VARCHAR(200) product_name
        DECIMAL(12,2) unit_price_amount
        CHAR(3) unit_price_currency
        INTEGER quantity
        TIMESTAMP created_at
        TIMESTAMP updated_at
    }
    products {
        UUID id
        VARCHAR(200) name
        VARCHAR(5000) description
        DECIMAL(12,2) price_amount
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
        VARCHAR(100) sku
        VARCHAR(50) size
        VARCHAR(50) color
        DECIMAL(12,2) price_adjustment_amount
        CHAR(3) price_adjustment_currency
        TIMESTAMP created_at
        TIMESTAMP updated_at
    }
    categories {
        UUID id
        VARCHAR(100) name
        TEXT description
        UUID parent_category_id
        INTEGER version
        TIMESTAMP created_at
        TIMESTAMP updated_at
    }
    inventory {
        UUID id
        UUID product_id
        INTEGER quantity_on_hand
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
    products ||--o{ order_line_items : "FK"
    categories ||--o{ products : "FK"
    products ||--o{ product_variants : "FK"
    categories ||--o{ categories : "FK"
    products ||--o{ inventory : "FK"
```
