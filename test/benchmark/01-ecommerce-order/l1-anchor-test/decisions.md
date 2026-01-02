
## Decisions from 2025-12-31

- **AMB-ENT-001: Product**
  - Q: What is the unique identifier format for products?
  - A: UUID
  - Decided: 2025-12-29 11:23 by user

- **AMB-ENT-002: Product**
  - Q: What is the data type and constraints for product price?
  - A: Decimal with 2 decimal places, minimum 0.01
  - Decided: 2025-12-29 11:26 by auto

- **AMB-ENT-003: Product**
  - Q: What is the maximum length for product name and description?
  - A: Name: 200 chars, Description: 5000 chars
  - Decided: 2025-12-29 11:26 by auto

- **AMB-ENT-004: Product**
  - Q: Should products support multiple images? How many?
  - A: Yes, up to 10 images with one primary
  - Decided: 2025-12-29 11:26 by auto

- **AMB-ENT-005: Product**
  - Q: Should removed products be soft-deleted or permanently deleted?
  - A: Soft delete to preserve order history references
  - Decided: 2025-12-29 11:26 by auto

- **AMB-ENT-006: Product**
  - Q: Can a product belong to multiple categories?
  - A: Single primary category with optional secondary categories
  - Decided: 2025-12-29 11:26 by auto

- **AMB-ENT-007: Product**
  - Q: Is product change history/audit trail required?
  - A: Track price changes and who modified
  - Decided: 2025-12-29 11:26 by auto

- **AMB-ENT-008: ProductVariant**
  - Q: Does each variant have its own price, or inherit from parent product?
  - A: Variants can have price override (optional, defaults to parent)
  - Decided: 2025-12-29 11:26 by auto

- **AMB-ENT-009: ProductVariant**
  - Q: Does each variant have its own inventory/stock level?
  - A: Yes, inventory tracked per variant
  - Decided: 2025-12-29 11:26 by auto

- **AMB-ENT-010: ProductVariant**
  - Q: What variant types are supported beyond size and color?
  - A: Configurable variant types per product
  - Decided: 2025-12-29 11:26 by auto

- **AMB-ENT-011: ProductVariant**
  - Q: Does each variant have its own SKU/identifier?
  - A: Yes, unique SKU per variant
  - Decided: 2025-12-29 11:26 by auto

- **AMB-ENT-012: Cart**
  - Q: How long should an inactive cart persist before expiration?
  - A: 30 days for logged-in users, 7 days for guests
  - Decided: 2025-12-29 11:26 by auto

- **AMB-ENT-013: Cart**
  - Q: Should guest users (not logged in) be able to have a cart?
  - A: Yes, with cart merge on login
  - Decided: 2025-12-29 11:26 by auto

- **AMB-ENT-014: Cart**
  - Q: What happens when an item in cart goes out of stock?
  - A: Item remains with warning, blocked at checkout
  - Decided: 2025-12-29 11:26 by auto

- **AMB-ENT-015: CartItem**
  - Q: Is there a maximum quantity limit per item?
  - A: Limited by available stock
  - Decided: 2025-12-29 11:26 by auto

- **AMB-ENT-016: CartItem**
  - Q: Should cart items store the price at time of adding, or always use current price?
  - A: Always use current price, notify on price changes
  - Decided: 2025-12-29 11:26 by auto

- **AMB-ENT-017: Order**
  - Q: What is the order number format?
  - A: Year prefix + sequential number (e.g., ORD-2024-00001)
  - Decided: 2025-12-29 11:26 by auto

- **AMB-ENT-018: Order**
  - Q: What are the valid state transitions for orders?
  - A: pending→confirmed→shipped→delivered, with cancelled possible from pending/confirmed
  - Decided: 2025-12-29 11:26 by auto

- **AMB-ENT-019: Order**
  - Q: Should orders track item prices at time of purchase?
  - A: Yes, snapshot prices at order time
  - Decided: 2025-12-29 11:26 by auto

- **AMB-ENT-020: Order**
  - Q: Is order modification allowed after placement (before shipping)?
  - A: No modification, must cancel and reorder
  - Decided: 2025-12-29 11:26 by auto

- **AMB-ENT-021: Order**
  - Q: Are partial shipments supported (splitting an order)?
  - A: No, single shipment per order
  - Decided: 2025-12-29 11:26 by auto

- **AMB-ENT-022: Order**
  - Q: Is order history/audit trail required for status changes?
  - A: Yes, track all status changes with timestamps
  - Decided: 2025-12-29 11:26 by auto

- **AMB-ENT-023: Order**
  - Q: What additional order attributes are needed (notes, gift message, etc.)?
  - A: Optional order notes field
  - Decided: 2025-12-29 11:26 by auto

- **AMB-ENT-024: Customer**
  - Q: What customer profile attributes are required?
  - A: Email, password, first name, last name, phone (optional)
  - Decided: 2025-12-29 11:26 by auto

- **AMB-ENT-025: Customer**
  - Q: What is the unique identifier for customers?
  - A: Email address as unique identifier
  - Decided: 2025-12-29 11:26 by auto

- **AMB-ENT-026: Customer**
  - Q: What are the password requirements?
  - A: Minimum 8 chars with at least one number
  - Decided: 2025-12-29 11:26 by auto

- **AMB-ENT-027: Customer**
  - Q: Can customers save multiple shipping addresses?
  - A: Yes, with one default address
  - Decided: 2025-12-29 11:26 by auto

- **AMB-ENT-028: Customer**
  - Q: Can customers save multiple payment methods?
  - A: Yes, with one default payment method
  - Decided: 2025-12-29 11:26 by auto

- **AMB-ENT-029: Customer**
  - Q: Should customer accounts be soft-deleted or allow permanent deletion (GDPR)?
  - A: Support both soft delete and full data erasure for GDPR
  - Decided: 2025-12-29 11:26 by auto

- **AMB-ENT-030: Customer**
  - Q: Is email verification required for registration?
  - A: Yes, email verification required before placing orders
  - Decided: 2025-12-29 11:26 by auto

- **AMB-ENT-031: ShippingAddress**
  - Q: What fields are required for a shipping address?
  - A: Street, city, state/province, postal code, country, recipient name
  - Decided: 2025-12-29 11:26 by auto

- **AMB-ENT-032: ShippingAddress**
  - Q: Which countries/regions are supported for shipping?
  - A: Domestic only (single country)
  - Decided: 2025-12-29 11:26 by auto

- **AMB-ENT-033: ShippingAddress**
  - Q: Should address validation be performed?
  - A: Basic format validation, no external service
  - Decided: 2025-12-29 11:26 by auto

- **AMB-ENT-034: PaymentMethod**
  - Q: What credit card types are accepted?
  - A: Visa, Mastercard, American Express
  - Decided: 2025-12-29 11:26 by auto

- **AMB-ENT-035: PaymentMethod**
  - Q: What payment gateway/processor will be used?
  - A: Stripe for credit cards, PayPal integration
  - Decided: 2025-12-29 11:26 by auto

- **AMB-ENT-036: PaymentMethod**
  - Q: Should payment details be stored for future purchases?
  - A: Store tokenized payment methods via payment gateway
  - Decided: 2025-12-29 11:26 by auto

- **AMB-ENT-037: Category**
  - Q: Should categories support hierarchical nesting (subcategories)?
  - A: Yes, support 2-3 levels of nesting
  - Decided: 2025-12-29 11:26 by auto

- **AMB-ENT-038: Category**
  - Q: What attributes does a category have?
  - A: Name, description, image, display order
  - Decided: 2025-12-29 11:26 by auto

- **AMB-ENT-039: Category**
  - Q: Can categories be deactivated/hidden without deletion?
  - A: Yes, support active/inactive status
  - Decided: 2025-12-29 11:26 by auto

- **AMB-ENT-040: Inventory**
  - Q: Should inventory support reserved stock (items in carts/pending orders)?
  - A: Yes, track available vs reserved stock
  - Decided: 2025-12-29 11:26 by auto

- **AMB-ENT-041: Inventory**
  - Q: Should low stock alerts be supported?
  - A: Yes, configurable threshold per product
  - Decided: 2025-12-29 11:26 by auto

- **AMB-ENT-042: Inventory**
  - Q: Is inventory history/audit trail needed?
  - A: Yes, track all stock changes with reason
  - Decided: 2025-12-29 11:26 by auto

- **AMB-ENT-043: Inventory**
  - Q: Should backorders be supported when out of stock?
  - A: No, strict stock enforcement
  - Decided: 2025-12-29 11:26 by auto

- **AMB-ENT-044: Order**
  - Q: What happens to inventory when an order is cancelled?
  - A: Automatically restore stock on cancellation
  - Decided: 2025-12-29 11:26 by auto

- **AMB-ENT-045: Order**
  - Q: Is a refund/return process needed?
  - A: Out of scope for MVP, add later
  - Decided: 2025-12-29 11:26 by auto

- **AMB-OP-001: Browse Products**
  - Q: What filtering and search capabilities are required for browsing products?
  - A: Filter by category, price range, and availability; basic keyword search
  - Decided: 2025-12-29 11:26 by auto

- **AMB-OP-002: Browse Products**
  - Q: What sorting options are available for product listings?
  - A: Sort by name, price (asc/desc), newest first
  - Decided: 2025-12-29 11:26 by auto

- **AMB-OP-003: Browse Products**
  - Q: How should product listing be paginated?
  - A: 20 products per page with pagination controls
  - Decided: 2025-12-29 11:26 by auto

- **AMB-OP-004: Browse Products**
  - Q: Should out-of-stock products be visible in listings?
  - A: Show with 'out of stock' indicator, not purchasable
  - Decided: 2025-12-29 11:26 by auto

- **AMB-OP-005: Add to Cart**
  - Q: What happens when adding an item already in the cart?
  - A: Increment quantity of existing cart item
  - Decided: 2025-12-29 11:26 by auto

- **AMB-OP-006: Add to Cart**
  - Q: What is the minimum and maximum quantity that can be added?
  - A: Minimum 1, maximum limited by available stock
  - Decided: 2025-12-29 11:26 by auto

- **AMB-OP-007: Add to Cart**
  - Q: Can guest (non-authenticated) users add items to cart?
  - A: Yes, with session-based cart
  - Decided: 2025-12-29 11:26 by auto

- **AMB-OP-008: Add to Cart**
  - Q: How should the system handle concurrent stock claims (race condition)?
  - A: Allow add, validate stock at checkout
  - Decided: 2025-12-29 11:26 by auto

- **AMB-OP-009: Add to Cart**
  - Q: What feedback should the user receive after adding to cart?
  - A: Toast notification with view cart option
  - Decided: 2025-12-29 11:26 by auto

- **AMB-OP-010: Add to Cart**
  - Q: If product has variants, how is variant selection handled?
  - A: Variant must be selected before add to cart is enabled
  - Decided: 2025-12-29 11:26 by auto

- **AMB-OP-011: Update Cart Quantity**
  - Q: What happens when quantity is updated to exceed available stock?
  - A: Limit to available stock and show notification
  - Decided: 2025-12-29 11:26 by auto

- **AMB-OP-012: Update Cart Quantity**
  - Q: What happens when quantity is set to zero?
  - A: Treat as remove from cart
  - Decided: 2025-12-29 11:26 by auto

- **AMB-OP-013: Update Cart Quantity**
  - Q: Is this operation idempotent with the same quantity value?
  - A: Yes, setting same quantity is a no-op
  - Decided: 2025-12-29 11:26 by auto

- **AMB-OP-014: Remove from Cart**
  - Q: Should removal require confirmation?
  - A: No confirmation, but provide undo option
  - Decided: 2025-12-29 11:26 by auto

- **AMB-OP-015: Remove from Cart**
  - Q: What happens when removing the last item from cart?
  - A: Show empty cart state with continue shopping link
  - Decided: 2025-12-29 11:26 by auto

- **AMB-OP-016: Place Order**
  - Q: What payment validation occurs before order is placed?
  - A: Full payment authorization before creating order
  - Decided: 2025-12-29 11:26 by auto

- **AMB-OP-017: Place Order**
  - Q: What happens if payment fails?
  - A: Show error, keep cart intact, allow retry
  - Decided: 2025-12-29 11:26 by auto

- **AMB-OP-018: Place Order**
  - Q: What happens if stock becomes unavailable during checkout?
  - A: Show error for affected items, allow order modification
  - Decided: 2025-12-29 11:26 by auto

- **AMB-OP-019: Place Order**
  - Q: When is inventory decremented for ordered items?
  - A: Immediately upon successful payment
  - Decided: 2025-12-29 11:26 by auto

- **AMB-OP-020: Place Order**
  - Q: How is the free shipping threshold ($50) calculated?
  - A: Based on subtotal before tax, after discounts
  - Decided: 2025-12-29 11:26 by auto

- **AMB-OP-021: Place Order**
  - Q: What shipping options and costs are available below $50?
  - A: Standard shipping at flat rate (e.g., $5.99)
  - Decided: 2025-12-29 11:26 by auto

- **AMB-OP-022: Place Order**
  - Q: Is tax calculation required? How is it determined?
  - A: Calculate based on shipping address (US state tax)
  - Decided: 2025-12-29 11:26 by auto

- **AMB-OP-023: Place Order**
  - Q: Are discount codes/coupons supported?
  - A: Out of scope for MVP
  - Decided: 2025-12-29 11:26 by auto

- **AMB-OP-024: Place Order**
  - Q: Should order placement be logged for audit?
  - A: Yes, log all order events with timestamp
  - Decided: 2025-12-29 11:26 by auto

- **AMB-OP-025: Track Order Status**
  - Q: Can customers view all their past orders or only recent ones?
  - A: All orders with pagination
  - Decided: 2025-12-29 11:26 by auto

- **AMB-OP-026: Track Order Status**
  - Q: What order details are visible to customers?
  - A: Items, quantities, prices, status, shipping address, tracking number
  - Decided: 2025-12-29 11:26 by auto

- **AMB-OP-027: Track Order Status**
  - Q: Is shipment tracking integration required?
  - A: Store tracking number, link to carrier website
  - Decided: 2025-12-29 11:26 by auto

- **AMB-OP-028: Cancel Order**
  - Q: What statuses allow order cancellation?
  - A: Only 'pending' and 'confirmed' statuses
  - Decided: 2025-12-29 11:26 by auto

- **AMB-OP-029: Cancel Order**
  - Q: Is a cancellation reason required?
  - A: Optional with predefined options
  - Decided: 2025-12-29 11:26 by auto

- **AMB-OP-030: Cancel Order**
  - Q: How is payment refunded on cancellation?
  - A: Automatic refund to original payment method
  - Decided: 2025-12-29 11:26 by auto

- **AMB-OP-031: Cancel Order**
  - Q: Is inventory automatically restored on cancellation?
  - A: Yes, automatically restore stock
  - Decided: 2025-12-29 11:26 by auto

- **AMB-OP-032: Cancel Order**
  - Q: Should cancellation send a confirmation email?
  - A: Yes, email with cancellation details and refund info
  - Decided: 2025-12-29 11:26 by auto

- **AMB-OP-033: Add Product**
  - Q: What validations are required for product data?
  - A: Name required (2-200 chars), price > 0, category required
  - Decided: 2025-12-29 11:26 by auto

- **AMB-OP-034: Add Product**
  - Q: Should product names be unique?
  - A: No, but warn if duplicate name exists
  - Decided: 2025-12-29 11:26 by auto

- **AMB-OP-035: Add Product**
  - Q: What is the initial state/status of a new product?
  - A: Draft status, requires explicit publish
  - Decided: 2025-12-29 11:26 by auto

- **AMB-OP-036: Add Product**
  - Q: How are product images uploaded and stored?
  - A: Upload to cloud storage (S3), store URL reference
  - Decided: 2025-12-29 11:26 by auto

- **AMB-OP-037: Add Product**
  - Q: Should product creation be audited?
  - A: Yes, log creator and timestamp
  - Decided: 2025-12-29 11:26 by auto

- **AMB-OP-038: Edit Product**
  - Q: Can product price be changed if it's in active orders or carts?
  - A: Yes, existing orders keep original price, carts get updated price
  - Decided: 2025-12-29 11:26 by auto

- **AMB-OP-039: Edit Product**
  - Q: Is concurrent editing by multiple admins handled?
  - A: Last write wins with conflict warning
  - Decided: 2025-12-29 11:26 by auto

- **AMB-OP-040: Edit Product**
  - Q: Should product edit history be maintained?
  - A: Track last modified by/date, not full history
  - Decided: 2025-12-29 11:26 by auto

- **AMB-OP-041: Remove Product**
  - Q: What happens to cart items referencing a removed product?
  - A: Remove from carts with notification on next cart view
  - Decided: 2025-12-29 11:26 by auto

- **AMB-OP-042: Remove Product**
  - Q: Can products with order history be removed?
  - A: Soft delete only to preserve order history
  - Decided: 2025-12-29 11:26 by auto

- **AMB-OP-043: Remove Product**
  - Q: Should product removal require confirmation?
  - A: Yes, confirmation with impact summary
  - Decided: 2025-12-29 11:26 by auto

- **AMB-OP-044: View All Orders**
  - Q: What filtering and search options are available for admin order list?
  - A: Filter by status, date range, customer; search by order number
  - Decided: 2025-12-29 11:26 by auto

- **AMB-OP-045: View All Orders**
  - Q: What order information is shown in the admin list view?
  - A: Order number, date, customer, total, status
  - Decided: 2025-12-29 11:26 by auto

- **AMB-OP-046: View All Orders**
  - Q: Can orders be exported (CSV, PDF)?
  - A: CSV export for reporting
  - Decided: 2025-12-29 11:26 by auto

- **AMB-OP-047: Update Order Status**
  - Q: What are the allowed status transitions?
  - A: pending→confirmed→shipped→delivered (with cancelled as terminal from pending/confirmed)
  - Decided: 2025-12-29 11:26 by auto

- **AMB-OP-048: Update Order Status**
  - Q: Should status changes trigger customer notifications?
  - A: Email notification on confirmed, shipped, delivered
  - Decided: 2025-12-29 11:26 by auto

- **AMB-OP-049: Update Order Status**
  - Q: Can tracking number be added when marking as shipped?
  - A: Yes, optional tracking number field when setting shipped
  - Decided: 2025-12-29 11:26 by auto

- **AMB-OP-050: Update Order Status**
  - Q: Should status change history be recorded?
  - A: Yes, log all changes with admin user and timestamp
  - Decided: 2025-12-29 11:26 by auto

- **AMB-OP-051: Manage Inventory**
  - Q: What inventory adjustment types are supported?
  - A: Set absolute quantity, or adjust by delta (+/-)
  - Decided: 2025-12-29 11:26 by auto

- **AMB-OP-052: Manage Inventory**
  - Q: Is a reason required for inventory adjustments?
  - A: Optional reason with predefined options
  - Decided: 2025-12-29 11:26 by auto

- **AMB-OP-053: Manage Inventory**
  - Q: Can inventory go negative?
  - A: No, minimum is zero
  - Decided: 2025-12-29 11:26 by auto

- **AMB-OP-054: Manage Inventory**
  - Q: Should bulk inventory updates be supported?
  - A: Yes, CSV import for bulk updates
  - Decided: 2025-12-29 11:26 by auto

- **AMB-OP-055: Manage Inventory**
  - Q: Should inventory changes be audited?
  - A: Yes, full audit trail with before/after values
  - Decided: 2025-12-29 11:26 by auto

- **AMB-OP-056: Send Confirmation Email**
  - Q: What information is included in the confirmation email?
  - A: Order number, items, total, shipping address, estimated delivery
  - Decided: 2025-12-29 11:26 by auto

- **AMB-OP-057: Send Confirmation Email**
  - Q: What happens if email sending fails?
  - A: Queue for retry, log failure, don't block order
  - Decided: 2025-12-29 11:26 by auto

- **AMB-OP-058: Send Confirmation Email**
  - Q: Is this operation idempotent (can it be safely retried)?
  - A: Yes with deduplication to prevent duplicate emails
  - Decided: 2025-12-29 11:26 by auto

- **AMB-OP-059: Send Confirmation Email**
  - Q: What email service/provider will be used?
  - A: Transactional email service (SendGrid, AWS SES)
  - Decided: 2025-12-29 11:26 by auto

- **AMB-OP-060: Send Confirmation Email**
  - Q: Should email sending be synchronous or asynchronous?
  - A: Asynchronous via message queue
  - Decided: 2025-12-29 11:26 by auto

