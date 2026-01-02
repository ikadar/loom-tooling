# Initial Data Model

Generated: 2026-01-02T18:27:24+01:00

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

**Purpose:** Stores registered users who can browse products, maintain carts, and place orders

**Columns:**

| Name | Type | Constraints | Default |
|------|------|-------------|----------|
| id | UUID | PRIMARY KEY |  |
| email | VARCHAR(255) | NOT NULL UNIQUE |  |
| password_hash | VARCHAR(255) | NOT NULL |  |
| registration_status | registration_status_enum | NOT NULL | 'REGISTERED' |
| version | INTEGER | NOT NULL | 1 |
| created_at | TIMESTAMP WITH TIME ZONE | NOT NULL | NOW() |
| updated_at | TIMESTAMP WITH TIME ZONE | NOT NULL | NOW() |

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
| total_price_amount | DECIMAL(12,2) | NOT NULL | 0.00 |
| total_price_currency | CHAR(3) | NOT NULL | 'USD' |
| version | INTEGER | NOT NULL | 1 |
| created_at | TIMESTAMP WITH TIME ZONE | NOT NULL | NOW() |
| updated_at | TIMESTAMP WITH TIME ZONE | NOT NULL | NOW() |

**Primary Key:** [id]

**Indexes:**
- `idx_carts_customer`: [customer_id]

**Foreign Keys:**
- [customer_id] → customers(id) (ON DELETE CASCADE)

**Constraints:**
- `chk_cart_total_non_negative`: total_price_amount >= 0

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
| unit_price_currency | CHAR(3) | NOT NULL | 'USD' |
| quantity | INTEGER | NOT NULL |  |
| subtotal_amount | DECIMAL(12,2) | NOT NULL |  |
| subtotal_currency | CHAR(3) | NOT NULL | 'USD' |
| created_at | TIMESTAMP WITH TIME ZONE | NOT NULL | NOW() |
| updated_at | TIMESTAMP WITH TIME ZONE | NOT NULL | NOW() |

**Primary Key:** [id]

**Indexes:**
- `idx_cart_items_cart`: [cart_id]
- `idx_cart_items_cart_product`: [cart_id product_id]

**Foreign Keys:**
- [cart_id] → carts(id) (ON DELETE CASCADE)
- [product_id] → products(id) (ON DELETE RESTRICT)

**Constraints:**
- `chk_cart_item_quantity`: quantity >= 1
- `chk_cart_item_price_positive`: unit_price_amount > 0

---

### TBL-ORDER-001 – orders {#tbl-order-001}

**Aggregate:** Order

**Purpose:** Represents a customer's purchase transaction with shipping and payment details

**Columns:**

| Name | Type | Constraints | Default |
|------|------|-------------|----------|
| id | UUID | PRIMARY KEY |  |
| customer_id | UUID | NOT NULL |  |
| status | order_status_enum | NOT NULL | 'PENDING' |
| shipping_street | VARCHAR(200) | NOT NULL |  |
| shipping_city | VARCHAR(100) | NOT NULL |  |
| shipping_state | VARCHAR(100) | NOT NULL |  |
| shipping_postal_code | VARCHAR(20) | NOT NULL |  |
| shipping_country | CHAR(2) | NOT NULL |  |
| payment_type | payment_type_enum | NOT NULL |  |
| payment_last_four_digits | CHAR(4) |  |  |
| payment_expiry_month | SMALLINT |  |  |
| payment_expiry_year | SMALLINT |  |  |
| payment_paypal_email | VARCHAR(255) |  |  |
| subtotal_amount | DECIMAL(12,2) | NOT NULL |  |
| subtotal_currency | CHAR(3) | NOT NULL | 'USD' |
| shipping_cost_amount | DECIMAL(12,2) | NOT NULL |  |
| shipping_cost_currency | CHAR(3) | NOT NULL | 'USD' |
| total_amount | DECIMAL(12,2) | NOT NULL |  |
| total_currency | CHAR(3) | NOT NULL | 'USD' |
| tracking_number | VARCHAR(100) |  |  |
| cancellation_reason | VARCHAR(500) |  |  |
| version | INTEGER | NOT NULL | 1 |
| created_at | TIMESTAMP WITH TIME ZONE | NOT NULL | NOW() |
| updated_at | TIMESTAMP WITH TIME ZONE | NOT NULL | NOW() |

**Primary Key:** [id]

**Indexes:**
- `idx_orders_customer`: [customer_id]
- `idx_orders_status`: [status]
- `idx_orders_created_at`: [created_at]

**Foreign Keys:**
- [customer_id] → customers(id) (ON DELETE RESTRICT)

**Constraints:**
- `chk_order_subtotal_non_negative`: subtotal_amount >= 0
- `chk_order_shipping_non_negative`: shipping_cost_amount >= 0
- `chk_order_total_non_negative`: total_amount >= 0
- `chk_order_expiry_month`: payment_expiry_month IS NULL OR (payment_expiry_month >= 1 AND payment_expiry_month <= 12)
- `chk_order_payment_credit_card`: payment_type != 'CREDIT_CARD' OR (payment_last_four_digits IS NOT NULL AND payment_expiry_month IS NOT NULL AND payment_expiry_year IS NOT NULL)
- `chk_order_payment_paypal`: payment_type != 'PAYPAL' OR payment_paypal_email IS NOT NULL

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
| unit_price_currency | CHAR(3) | NOT NULL | 'USD' |
| quantity | INTEGER | NOT NULL |  |
| subtotal_amount | DECIMAL(12,2) | NOT NULL |  |
| subtotal_currency | CHAR(3) | NOT NULL | 'USD' |
| created_at | TIMESTAMP WITH TIME ZONE | NOT NULL | NOW() |
| updated_at | TIMESTAMP WITH TIME ZONE | NOT NULL | NOW() |

**Primary Key:** [id]

**Indexes:**
- `idx_order_line_items_order`: [order_id]

**Foreign Keys:**
- [order_id] → orders(id) (ON DELETE CASCADE)
- [product_id] → products(id) (ON DELETE RESTRICT)

**Constraints:**
- `chk_line_item_quantity`: quantity >= 1
- `chk_line_item_price_positive`: unit_price_amount > 0

---

### TBL-PRODUCT-001 – products {#tbl-product-001}

**Aggregate:** Product

**Purpose:** Represents an item available for sale in the store

**Columns:**

| Name | Type | Constraints | Default |
|------|------|-------------|----------|
| id | UUID | PRIMARY KEY |  |
| name | VARCHAR(200) | NOT NULL |  |
| description | TEXT |  |  |
| price_amount | DECIMAL(12,2) | NOT NULL |  |
| price_currency | CHAR(3) | NOT NULL | 'USD' |
| category_id | UUID | NOT NULL |  |
| is_active | BOOLEAN | NOT NULL | true |
| version | INTEGER | NOT NULL | 1 |
| created_at | TIMESTAMP WITH TIME ZONE | NOT NULL | NOW() |
| updated_at | TIMESTAMP WITH TIME ZONE | NOT NULL | NOW() |

**Primary Key:** [id]

**Indexes:**
- `idx_products_category`: [category_id]
- `idx_products_active`: [is_active]
- `idx_products_name`: [name]

**Foreign Keys:**
- [category_id] → categories(id) (ON DELETE RESTRICT)

**Constraints:**
- `chk_product_price_positive`: price_amount > 0
- `chk_product_name_length`: LENGTH(name) >= 1 AND LENGTH(name) <= 200
- `chk_product_description_length`: description IS NULL OR LENGTH(description) <= 5000

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
| size | VARCHAR(50) |  |  |
| color | VARCHAR(50) |  |  |
| price_adjustment_amount | DECIMAL(12,2) |  | 0.00 |
| price_adjustment_currency | CHAR(3) |  | 'USD' |
| created_at | TIMESTAMP WITH TIME ZONE | NOT NULL | NOW() |
| updated_at | TIMESTAMP WITH TIME ZONE | NOT NULL | NOW() |

**Primary Key:** [id]

**Indexes:**
- `idx_product_variants_product`: [product_id]
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
| description | TEXT |  |  |
| parent_category_id | UUID |  |  |
| version | INTEGER | NOT NULL | 1 |
| created_at | TIMESTAMP WITH TIME ZONE | NOT NULL | NOW() |
| updated_at | TIMESTAMP WITH TIME ZONE | NOT NULL | NOW() |

**Primary Key:** [id]

**Indexes:**
- `idx_categories_parent`: [parent_category_id]
- `idx_categories_name`: [name]

**Foreign Keys:**
- [parent_category_id] → categories(id) (ON DELETE SET NULL)

**Constraints:**
- `chk_category_not_self_parent`: parent_category_id IS NULL OR parent_category_id != id
- `chk_category_name_length`: LENGTH(name) >= 1 AND LENGTH(name) <= 100

---

### TBL-INVENTORY-001 – inventories {#tbl-inventory-001}

**Aggregate:** Inventory

**Purpose:** Tracks stock levels for products and manages availability

**Columns:**

| Name | Type | Constraints | Default |
|------|------|-------------|----------|
| id | UUID | PRIMARY KEY |  |
| product_id | UUID | NOT NULL UNIQUE |  |
| quantity_on_hand | INTEGER | NOT NULL | 0 |
| reserved_quantity | INTEGER | NOT NULL | 0 |
| version | INTEGER | NOT NULL | 1 |
| created_at | TIMESTAMP WITH TIME ZONE | NOT NULL | NOW() |
| updated_at | TIMESTAMP WITH TIME ZONE | NOT NULL | NOW() |

**Primary Key:** [id]

**Indexes:**
- `idx_inventories_product`: [product_id]

**Foreign Keys:**
- [product_id] → products(id) (ON DELETE CASCADE)

**Constraints:**
- `chk_inventory_quantity_non_negative`: quantity_on_hand >= 0
- `chk_inventory_reserved_non_negative`: reserved_quantity >= 0
- `chk_inventory_reserved_not_exceed`: reserved_quantity <= quantity_on_hand

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
        TIMESTAMP WITH TIME ZONE created_at
        TIMESTAMP WITH TIME ZONE updated_at
    }
    carts {
        UUID id
        UUID customer_id
        DECIMAL(12,2) total_price_amount
        CHAR(3) total_price_currency
        INTEGER version
        TIMESTAMP WITH TIME ZONE created_at
        TIMESTAMP WITH TIME ZONE updated_at
    }
    cart_items {
        UUID id
        UUID cart_id
        UUID product_id
        VARCHAR(200) product_name
        DECIMAL(12,2) unit_price_amount
        CHAR(3) unit_price_currency
        INTEGER quantity
        DECIMAL(12,2) subtotal_amount
        CHAR(3) subtotal_currency
        TIMESTAMP WITH TIME ZONE created_at
        TIMESTAMP WITH TIME ZONE updated_at
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
        VARCHAR(500) cancellation_reason
        INTEGER version
        TIMESTAMP WITH TIME ZONE created_at
        TIMESTAMP WITH TIME ZONE updated_at
    }
    order_line_items {
        UUID id
        UUID order_id
        UUID product_id
        VARCHAR(200) product_name
        DECIMAL(12,2) unit_price_amount
        CHAR(3) unit_price_currency
        INTEGER quantity
        DECIMAL(12,2) subtotal_amount
        CHAR(3) subtotal_currency
        TIMESTAMP WITH TIME ZONE created_at
        TIMESTAMP WITH TIME ZONE updated_at
    }
    products {
        UUID id
        VARCHAR(200) name
        TEXT description
        DECIMAL(12,2) price_amount
        CHAR(3) price_currency
        UUID category_id
        BOOLEAN is_active
        INTEGER version
        TIMESTAMP WITH TIME ZONE created_at
        TIMESTAMP WITH TIME ZONE updated_at
    }
    product_variants {
        UUID id
        UUID product_id
        VARCHAR(100) sku
        VARCHAR(50) size
        VARCHAR(50) color
        DECIMAL(12,2) price_adjustment_amount
        CHAR(3) price_adjustment_currency
        TIMESTAMP WITH TIME ZONE created_at
        TIMESTAMP WITH TIME ZONE updated_at
    }
    categories {
        UUID id
        VARCHAR(100) name
        TEXT description
        UUID parent_category_id
        INTEGER version
        TIMESTAMP WITH TIME ZONE created_at
        TIMESTAMP WITH TIME ZONE updated_at
    }
    inventories {
        UUID id
        UUID product_id
        INTEGER quantity_on_hand
        INTEGER reserved_quantity
        INTEGER version
        TIMESTAMP WITH TIME ZONE created_at
        TIMESTAMP WITH TIME ZONE updated_at
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
    products ||--o{ inventories : "FK"
```
