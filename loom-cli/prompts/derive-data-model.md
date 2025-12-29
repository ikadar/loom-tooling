# Initial Data Model Derivation Prompt

You are an expert database architect. Generate an Initial Data Model from L1 and L2 documents.

OUTPUT REQUIREMENT: Wrap your JSON response in ```json code blocks. No explanations.

## Your Task

From Domain Model and Aggregate Design, derive an Initial Data Model that specifies:
1. **Tables/Collections** - Database structure for each aggregate
2. **Fields** - Column definitions with types and constraints
3. **Indexes** - Performance optimization indexes
4. **Relationships** - Foreign keys and references

## Data Model Structure

For each table, include:
- id: Unique identifier (TBL-{NAME}-{NNN})
- name: Table name (snake_case)
- aggregate: Source aggregate
- fields: Column definitions
- primaryKey: PK definition
- indexes: Index definitions
- foreignKeys: FK relationships

## Output Format

```json
{
  "tables": [
    {
      "id": "TBL-ORDER-001",
      "name": "orders",
      "aggregate": "Order",
      "purpose": "Stores order header information",
      "fields": [
        {"name": "id", "type": "UUID", "constraints": "NOT NULL", "default": "gen_random_uuid()"},
        {"name": "customer_id", "type": "UUID", "constraints": "NOT NULL"},
        {"name": "status", "type": "VARCHAR(20)", "constraints": "NOT NULL", "default": "'pending'"},
        {"name": "subtotal_amount", "type": "DECIMAL(10,2)", "constraints": "NOT NULL"},
        {"name": "subtotal_currency", "type": "CHAR(3)", "constraints": "NOT NULL", "default": "'USD'"},
        {"name": "shipping_amount", "type": "DECIMAL(10,2)", "constraints": "NOT NULL"},
        {"name": "shipping_currency", "type": "CHAR(3)", "constraints": "NOT NULL", "default": "'USD'"},
        {"name": "total_amount", "type": "DECIMAL(10,2)", "constraints": "NOT NULL"},
        {"name": "total_currency", "type": "CHAR(3)", "constraints": "NOT NULL", "default": "'USD'"},
        {"name": "shipping_street", "type": "VARCHAR(200)", "constraints": "NOT NULL"},
        {"name": "shipping_city", "type": "VARCHAR(100)", "constraints": "NOT NULL"},
        {"name": "shipping_state", "type": "VARCHAR(50)", "constraints": "NOT NULL"},
        {"name": "shipping_postal_code", "type": "VARCHAR(20)", "constraints": "NOT NULL"},
        {"name": "shipping_country", "type": "CHAR(2)", "constraints": "NOT NULL"},
        {"name": "version", "type": "INTEGER", "constraints": "NOT NULL", "default": "1"},
        {"name": "created_at", "type": "TIMESTAMP", "constraints": "NOT NULL", "default": "NOW()"},
        {"name": "updated_at", "type": "TIMESTAMP", "constraints": "NOT NULL", "default": "NOW()"}
      ],
      "primaryKey": {"columns": ["id"]},
      "indexes": [
        {"name": "idx_orders_customer", "columns": ["customer_id"]},
        {"name": "idx_orders_status", "columns": ["status"]},
        {"name": "idx_orders_created", "columns": ["created_at"]}
      ],
      "foreignKeys": [
        {"columns": ["customer_id"], "references": "customers(id)", "onDelete": "RESTRICT"}
      ],
      "checkConstraints": [
        {"name": "chk_order_status", "expression": "status IN ('pending', 'confirmed', 'shipped', 'delivered', 'cancelled')"},
        {"name": "chk_order_amounts", "expression": "total_amount >= 0 AND subtotal_amount >= 0 AND shipping_amount >= 0"}
      ]
    },
    {
      "id": "TBL-ORDER-002",
      "name": "order_line_items",
      "aggregate": "Order",
      "purpose": "Stores order line items (child of orders)",
      "fields": [
        {"name": "id", "type": "UUID", "constraints": "NOT NULL", "default": "gen_random_uuid()"},
        {"name": "order_id", "type": "UUID", "constraints": "NOT NULL"},
        {"name": "product_id", "type": "UUID", "constraints": "NOT NULL"},
        {"name": "product_name", "type": "VARCHAR(200)", "constraints": "NOT NULL"},
        {"name": "unit_price_amount", "type": "DECIMAL(10,2)", "constraints": "NOT NULL"},
        {"name": "unit_price_currency", "type": "CHAR(3)", "constraints": "NOT NULL", "default": "'USD'"},
        {"name": "quantity", "type": "INTEGER", "constraints": "NOT NULL"},
        {"name": "subtotal_amount", "type": "DECIMAL(10,2)", "constraints": "NOT NULL"},
        {"name": "subtotal_currency", "type": "CHAR(3)", "constraints": "NOT NULL", "default": "'USD'"},
        {"name": "created_at", "type": "TIMESTAMP", "constraints": "NOT NULL", "default": "NOW()"}
      ],
      "primaryKey": {"columns": ["id"]},
      "indexes": [
        {"name": "idx_line_items_order", "columns": ["order_id"]}
      ],
      "foreignKeys": [
        {"columns": ["order_id"], "references": "orders(id)", "onDelete": "CASCADE"}
      ],
      "checkConstraints": [
        {"name": "chk_line_item_quantity", "expression": "quantity > 0"},
        {"name": "chk_line_item_price", "expression": "unit_price_amount >= 0 AND subtotal_amount >= 0"}
      ]
    }
  ],
  "enums": [
    {
      "name": "order_status",
      "values": ["pending", "confirmed", "shipped", "delivered", "cancelled"]
    },
    {
      "name": "payment_type",
      "values": ["credit_card", "paypal"]
    }
  ],
  "summary": {
    "total_tables": 8,
    "total_indexes": 15,
    "total_foreign_keys": 6,
    "total_enums": 3
  }
}
```

## Database Design Principles

Apply these principles:
1. **Aggregate = Transaction boundary** - Tables within aggregate can be JOINed
2. **Value objects embedded** - Money, Address stored as multiple columns
3. **IDs are UUIDs** - Use UUID for all primary keys
4. **Optimistic locking** - Version column for concurrency
5. **Audit columns** - created_at, updated_at on all tables
6. **Snake case** - Use snake_case for all identifiers

## Quality Checklist

Before output, verify:
- [ ] Every aggregate has corresponding table(s)
- [ ] All constraints from domain model are represented
- [ ] Indexes cover common query patterns
- [ ] Foreign keys enforce referential integrity
- [ ] Naming is consistent (snake_case)

---

REMINDER: Output ONLY a ```json code block. No explanations.

L1/L2 INPUT (Domain Model + Aggregate Design):

