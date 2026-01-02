# TDAI Test Cases

Generated: 2026-01-02T18:27:24+01:00

**Methodology:** Test-Driven AI Development (TDAI)

## Summary

| Category | Count | Ratio |
|----------|-------|-------|
| Positive | 89 | 31.1% |
| Negative | 77 | 26.9% |
| Boundary | 62 | - |
| Hallucination Prevention | 58 | - |
| **Total** | **286** | - |

**Coverage:** 29 ACs covered
**Hallucination Prevention:** ✓ Enabled

---

## Positive Tests (Happy Path)

### TC-AC-PROD-001-P01 – Filter products by single category {#tc-ac-prod-001-p01}

**Preconditions:**
- Products exist in multiple categories
- User is on product listing page

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| category | Electronics | Valid category with products |

**Steps:**
1. Navigate to product listing
2. Select category filter 'Electronics'
3. Apply filters

**Expected Result:**
- Only Electronics products displayed
- Other categories hidden

**Traceability:**
- AC: [AC-PROD-001](../l1/acceptance-criteria.md#ac-prod-001)
- BR: [BR-PROD](../l1/business-rules.md#br-prod)

---

### TC-AC-PROD-001-P02 – Filter products by price range and availability {#tc-ac-prod-001-p02}

**Preconditions:**
- Products exist with various prices
- Some products in/out of stock

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| priceMin | 10 | Lower bound |
| priceMax | 50 | Upper bound |
| availability | inStock | Only available items |

**Steps:**
1. Set price range 10-50
2. Select 'In Stock' filter
3. Apply filters

**Expected Result:**
- Only in-stock products $10-$50 shown
- Out of stock items hidden

**Traceability:**
- AC: [AC-PROD-001](../l1/acceptance-criteria.md#ac-prod-001)
- BR: [BR-PROD](../l1/business-rules.md#br-prod)

---

### TC-AC-PROD-001-P03 – Sort products by price descending {#tc-ac-prod-001-p03}

**Preconditions:**
- Multiple products with different prices exist

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| sortBy | price_desc | Price high to low |

**Steps:**
1. Navigate to product listing
2. Select sort by 'Price: High to Low'

**Expected Result:**
- Products ordered by price descending
- Highest price first

**Traceability:**
- AC: [AC-PROD-001](../l1/acceptance-criteria.md#ac-prod-001)
- BR: [BR-PROD](../l1/business-rules.md#br-prod)

---

### TC-AC-PROD-001-P04 – Display 20 products per page {#tc-ac-prod-001-p04}

**Preconditions:**
- More than 20 products exist matching criteria

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| totalProducts | 50 | Enough for pagination |

**Steps:**
1. Navigate to product listing with 50+ products
2. Count displayed products

**Expected Result:**
- Exactly 20 products shown
- Pagination controls visible

**Traceability:**
- AC: [AC-PROD-001](../l1/acceptance-criteria.md#ac-prod-001)
- BR: [BR-PROD](../l1/business-rules.md#br-prod)

---

### TC-AC-PROD-002-P01 – Display out of stock indicator on product listing {#tc-ac-prod-002-p01}

**Preconditions:**
- Product exists with zero inventory

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| inventory | 0 | Zero stock |

**Steps:**
1. Navigate to product listing
2. Locate out-of-stock product

**Expected Result:**
- 'Out of Stock' indicator visible
- Product still displayed

**Traceability:**
- AC: [AC-PROD-002](../l1/acceptance-criteria.md#ac-prod-002)
- BR: [BR-PROD](../l1/business-rules.md#br-prod)

---

### TC-AC-PROD-002-P02 – Display out of stock indicator on product detail page {#tc-ac-prod-002-p02}

**Preconditions:**
- Product exists with zero inventory

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| productId | PROD-001 | Out of stock product |

**Steps:**
1. Navigate to product detail page for out-of-stock item

**Expected Result:**
- 'Out of Stock' indicator visible on detail page

**Traceability:**
- AC: [AC-PROD-002](../l1/acceptance-criteria.md#ac-prod-002)
- BR: [BR-PROD](../l1/business-rules.md#br-prod)

---

### TC-AC-PROD-002-P03 – Add to cart button disabled for out of stock product {#tc-ac-prod-002-p03}

**Preconditions:**
- Product has zero inventory

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| inventory | 0 | No stock available |

**Steps:**
1. View out-of-stock product detail
2. Examine add-to-cart button

**Expected Result:**
- Add to cart button is disabled
- Button not clickable

**Traceability:**
- AC: [AC-PROD-002](../l1/acceptance-criteria.md#ac-prod-002)
- BR: [BR-PROD](../l1/business-rules.md#br-prod)

---

### TC-AC-CART-001-P01 – Add single quantity of in-stock product to cart {#tc-ac-cart-001-p01}

**Preconditions:**
- Product in stock
- Customer on product page

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| productId | PROD-001 | In-stock product |
| quantity | 1 | Single item |

**Steps:**
1. View in-stock product
2. Set quantity to 1
3. Click Add to Cart

**Expected Result:**
- Item added to cart
- Toast notification appears

**Traceability:**
- AC: [AC-CART-001](../l1/acceptance-criteria.md#ac-cart-001)
- BR: [BR-CART](../l1/business-rules.md#br-cart)

---

### TC-AC-CART-001-P02 – Toast notification with View Cart option appears {#tc-ac-cart-001-p02}

**Preconditions:**
- Product added to cart successfully

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| productId | PROD-002 | Valid product |

**Steps:**
1. Add product to cart
2. Observe notification

**Expected Result:**
- Toast notification visible
- 'View Cart' link in toast

**Traceability:**
- AC: [AC-CART-001](../l1/acceptance-criteria.md#ac-cart-001)
- BR: [BR-CART](../l1/business-rules.md#br-cart)

---

### TC-AC-CART-001-P03 – Cart icon updates with new item count {#tc-ac-cart-001-p03}

**Preconditions:**
- Cart initially empty or has known count

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| initialCount | 0 | Empty cart |
| quantity | 3 | Adding 3 items |

**Steps:**
1. Note cart icon count
2. Add 3 items
3. Check cart icon

**Expected Result:**
- Cart icon shows updated count of 3

**Traceability:**
- AC: [AC-CART-001](../l1/acceptance-criteria.md#ac-cart-001)
- BR: [BR-CART](../l1/business-rules.md#br-cart)

---

### TC-AC-CART-001-P04 – Add product with selected variant to cart {#tc-ac-cart-001-p04}

**Preconditions:**
- Product has variants
- Variant selected

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| productId | PROD-VAR | Product with variants |
| variant | Size: Large | Selected variant |

**Steps:**
1. Select variant
2. Click Add to Cart

**Expected Result:**
- Correct variant added to cart
- Toast shown

**Traceability:**
- AC: [AC-CART-001](../l1/acceptance-criteria.md#ac-cart-001)
- BR: [BR-CART](../l1/business-rules.md#br-cart)

---

### TC-AC-CART-002-P01 – Increment quantity when adding existing product {#tc-ac-cart-002-p01}

**Preconditions:**
- Product A in cart with quantity 2

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| productId | PROD-A | Already in cart |
| existingQty | 2 | Current quantity |
| addQty | 1 | Adding 1 more |

**Steps:**
1. Add Product A again with quantity 1

**Expected Result:**
- Cart shows Product A with quantity 3

**Traceability:**
- AC: [AC-CART-002](../l1/acceptance-criteria.md#ac-cart-002)
- BR: [BR-CART](../l1/business-rules.md#br-cart)

---

### TC-AC-CART-002-P02 – No new line item created for existing product {#tc-ac-cart-002-p02}

**Preconditions:**
- Cart has 1 line item for Product A

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| initialLineItems | 1 | One line item |
| productId | PROD-A | Same product |

**Steps:**
1. Count line items
2. Add Product A again
3. Count line items

**Expected Result:**
- Still only 1 line item
- Quantity increased

**Traceability:**
- AC: [AC-CART-002](../l1/acceptance-criteria.md#ac-cart-002)
- BR: [BR-CART](../l1/business-rules.md#br-cart)

---

### TC-AC-CART-002-P03 – Multiple additions increment correctly {#tc-ac-cart-002-p03}

**Preconditions:**
- Product in cart with quantity 1
- Stock is 10

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| existingQty | 1 | Starting quantity |
| additions | 3 | Three separate adds of 1 |

**Steps:**
1. Add product 1
2. Add same product again
3. Add again
4. Add again

**Expected Result:**
- Final quantity is 4
- Single line item

**Traceability:**
- AC: [AC-CART-002](../l1/acceptance-criteria.md#ac-cart-002)
- BR: [BR-CART](../l1/business-rules.md#br-cart)

---

### TC-AC-CART-003-P01 – Guest user can add items to cart without login {#tc-ac-cart-003-p01}

**Preconditions:**
- User not logged in
- Product available

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| userState | guest | Not authenticated |

**Steps:**
1. Browse as guest
2. Add product to cart

**Expected Result:**
- Product added successfully
- Cart created for session

**Traceability:**
- AC: [AC-CART-003](../l1/acceptance-criteria.md#ac-cart-003)
- BR: [BR-CART](../l1/business-rules.md#br-cart)

---

### TC-AC-CART-003-P02 – Session-based cart persists across page navigation {#tc-ac-cart-003-p02}

**Preconditions:**
- Guest has items in cart

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| cartItems | 2 | Items in guest cart |

**Steps:**
1. Add items as guest
2. Navigate to other pages
3. Return to cart

**Expected Result:**
- Cart still contains items
- No data lost

**Traceability:**
- AC: [AC-CART-003](../l1/acceptance-criteria.md#ac-cart-003)
- BR: [BR-CART](../l1/business-rules.md#br-cart)

---

### TC-AC-CART-003-P03 – Guest cart persists after browser close within 7 days {#tc-ac-cart-003-p03}

**Preconditions:**
- Guest added items to cart

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| daysSinceActivity | 3 | Within 7 day window |

**Steps:**
1. Add items as guest
2. Close browser
3. Return after 3 days

**Expected Result:**
- Cart still contains items

**Traceability:**
- AC: [AC-CART-003](../l1/acceptance-criteria.md#ac-cart-003)
- BR: [BR-CART](../l1/business-rules.md#br-cart)

---

### TC-AC-CART-004-P01 – Merge guest cart with empty account cart {#tc-ac-cart-004-p01}

**Preconditions:**
- Guest has 2 items in session cart
- Account cart is empty

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| guestCartItems | [SKU-A:2, SKU-B:1] | Guest items |
| accountCartItems | [] | Empty account cart |

**Steps:**
1. Add items to guest cart
2. Login with registered account
3. Verify merged cart

**Expected Result:**
- Account cart contains SKU-A:2 and SKU-B:1

**Traceability:**
- AC: [AC-CART-004](../l1/acceptance-criteria.md#ac-cart-004)
- BR: [AMB-ENT-013](../l1/business-rules.md#amb-ent-013)

---

### TC-AC-CART-004-P02 – Merge carts combining quantities for duplicate products {#tc-ac-cart-004-p02}

**Preconditions:**
- Guest has SKU-A:2
- Account has SKU-A:3

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| guestCartItems | [SKU-A:2] | Guest quantity |
| accountCartItems | [SKU-A:3] | Account quantity |

**Steps:**
1. Add SKU-A to guest cart
2. Login
3. Verify quantities combined

**Expected Result:**
- Account cart contains SKU-A:5 (2+3)

**Traceability:**
- AC: [AC-CART-004](../l1/acceptance-criteria.md#ac-cart-004)
- BR: [AMB-ENT-013](../l1/business-rules.md#amb-ent-013)

---

### TC-AC-CART-004-P03 – Merge carts with mix of unique and duplicate items {#tc-ac-cart-004-p03}

**Preconditions:**
- Guest has SKU-A:1, SKU-B:2
- Account has SKU-A:1, SKU-C:1

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| guestCartItems | [SKU-A:1, SKU-B:2] | Guest items |
| accountCartItems | [SKU-A:1, SKU-C:1] | Account items |

**Steps:**
1. Setup guest and account carts
2. Login
3. Verify merged result

**Expected Result:**
- Cart contains SKU-A:2, SKU-B:2, SKU-C:1

**Traceability:**
- AC: [AC-CART-004](../l1/acceptance-criteria.md#ac-cart-004)
- BR: [AMB-ENT-013](../l1/business-rules.md#amb-ent-013)

---

### TC-AC-CART-005-P01 – Update quantity and recalculate cart total {#tc-ac-cart-005-p01}

**Preconditions:**
- Cart has SKU-A qty 2 at $10 each

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| oldQuantity | 2 | Current qty |
| newQuantity | 5 | New qty |
| unitPrice | 10.00 | Price per item |

**Steps:**
1. View cart
2. Change quantity to 5
3. Verify total

**Expected Result:**
- Quantity shows 5
- Total recalculated to $50

**Traceability:**
- AC: [AC-CART-005](../l1/acceptance-criteria.md#ac-cart-005)
- BR: [AMB-OP-011](../l1/business-rules.md#amb-op-011), [AMB-OP-012](../l1/business-rules.md#amb-op-012), [AMB-OP-013](../l1/business-rules.md#amb-op-013)

---

### TC-AC-CART-005-P02 – Decrease quantity successfully {#tc-ac-cart-005-p02}

**Preconditions:**
- Cart has SKU-A qty 5

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| oldQuantity | 5 | Current |
| newQuantity | 2 | Decreased |

**Steps:**
1. Change quantity from 5 to 2
2. Verify update

**Expected Result:**
- Quantity updated to 2
- Total recalculated

**Traceability:**
- AC: [AC-CART-005](../l1/acceptance-criteria.md#ac-cart-005)
- BR: [AMB-OP-011](../l1/business-rules.md#amb-op-011)

---

### TC-AC-CART-006-P01 – Remove item immediately without confirmation {#tc-ac-cart-006-p01}

**Preconditions:**
- Cart has SKU-A and SKU-B

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| itemToRemove | SKU-A | Target item |

**Steps:**
1. Click remove on SKU-A
2. Verify immediate removal

**Expected Result:**
- SKU-A removed immediately
- No confirmation dialog shown

**Traceability:**
- AC: [AC-CART-006](../l1/acceptance-criteria.md#ac-cart-006)
- BR: [AMB-OP-014](../l1/business-rules.md#amb-op-014), [AMB-OP-015](../l1/business-rules.md#amb-op-015)

---

### TC-AC-CART-006-P02 – Undo option shown briefly after removal {#tc-ac-cart-006-p02}

**Preconditions:**
- Cart has items

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| itemToRemove | SKU-A | Any item |

**Steps:**
1. Remove item
2. Verify undo option appears

**Expected Result:**
- Undo option displayed
- Option disappears after brief period

**Traceability:**
- AC: [AC-CART-006](../l1/acceptance-criteria.md#ac-cart-006)
- BR: [AMB-OP-015](../l1/business-rules.md#amb-op-015)

---

### TC-AC-CART-006-P03 – Undo restores removed item {#tc-ac-cart-006-p03}

**Preconditions:**
- Item just removed
- Undo option visible

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| removedItem | SKU-A:2 | Item to restore |

**Steps:**
1. Remove item
2. Click Undo
3. Verify restoration

**Expected Result:**
- Item restored to cart
- Original quantity preserved

**Traceability:**
- AC: [AC-CART-006](../l1/acceptance-criteria.md#ac-cart-006)
- BR: [AMB-OP-015](../l1/business-rules.md#amb-op-015)

---

### TC-AC-CART-007-P01 – Out of stock item remains in cart with warning {#tc-ac-cart-007-p01}

**Preconditions:**
- Cart has SKU-A
- SKU-A stock drops to 0

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| item | SKU-A | Out of stock item |
| stock | 0 | No inventory |

**Steps:**
1. Add item to cart
2. Stock drops to 0
3. View cart

**Expected Result:**
- Item still in cart
- Warning indicator shown

**Traceability:**
- AC: [AC-CART-007](../l1/acceptance-criteria.md#ac-cart-007)
- BR: [AMB-ENT-014](../l1/business-rules.md#amb-ent-014)

---

### TC-AC-CART-007-P02 – Warning indicator visible on out of stock item {#tc-ac-cart-007-p02}

**Preconditions:**
- Item in cart has 0 stock

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| stockLevel | 0 | Zero stock |

**Steps:**
1. View cart with out of stock item

**Expected Result:**
- Visual warning indicator displayed on item

**Traceability:**
- AC: [AC-CART-007](../l1/acceptance-criteria.md#ac-cart-007)
- BR: [AMB-ENT-014](../l1/business-rules.md#amb-ent-014)

---

### TC-AC-CART-008-P01 – Cart displays current price after price increase {#tc-ac-cart-008-p01}

**Preconditions:**
- Cart has SKU-A at $10
- Price changes to $12

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| oldPrice | 10.00 | Original |
| newPrice | 12.00 | Updated |

**Steps:**
1. Add item at $10
2. Price changes to $12
3. View cart

**Expected Result:**
- Cart shows $12 price
- Total reflects new price

**Traceability:**
- AC: [AC-CART-008](../l1/acceptance-criteria.md#ac-cart-008)
- BR: [AMB-ENT-016](../l1/business-rules.md#amb-ent-016)

---

### TC-AC-CART-008-P02 – Cart displays current price after price decrease {#tc-ac-cart-008-p02}

**Preconditions:**
- Cart has SKU-A at $10
- Price changes to $8

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| oldPrice | 10.00 | Original |
| newPrice | 8.00 | Reduced |

**Steps:**
1. Add item at $10
2. Price drops to $8
3. View cart

**Expected Result:**
- Cart shows $8 price
- Customer benefits from reduction

**Traceability:**
- AC: [AC-CART-008](../l1/acceptance-criteria.md#ac-cart-008)
- BR: [AMB-ENT-016](../l1/business-rules.md#amb-ent-016)

---

### TC-AC-CART-008-P03 – Customer notified of price change {#tc-ac-cart-008-p03}

**Preconditions:**
- Item price changed since added to cart

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| priceChange | +2.00 | Price increased |

**Steps:**
1. View cart after price change

**Expected Result:**
- Notification about price change shown to customer

**Traceability:**
- AC: [AC-CART-008](../l1/acceptance-criteria.md#ac-cart-008)
- BR: [AMB-ENT-016](../l1/business-rules.md#amb-ent-016)

---

### TC-AC-ORDER-001-P01 – Successfully place order with all prerequisites met {#tc-ac-order-001-p01}

**Preconditions:**
- Customer is registered
- Email is verified
- Cart has items
- Valid shipping address exists
- Valid payment method exists

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| cartItems | 2 valid products | Multiple items to verify all decremented |
| paymentMethod | valid credit card | Pre-authorized payment method |

**Steps:**
1. Login as verified customer
2. Add items to cart
3. Proceed to checkout
4. Select shipping address
5. Select payment method
6. Submit order

**Expected Result:**
- Payment is authorized
- Inventory decremented for all items
- Order created with status pending
- Cart is cleared
- Confirmation page shown

**Traceability:**
- AC: [AC-ORDER-001](../l1/acceptance-criteria.md#ac-order-001)
- BR: [BR-ORDER](../l1/business-rules.md#br-order)

---

### TC-AC-ORDER-001-P02 – Confirmation email queued after successful order {#tc-ac-order-001-p02}

**Preconditions:**
- Customer is registered and verified
- Cart has items
- Valid address and payment

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| email | customer@example.com | Verified email address |

**Steps:**
1. Complete checkout process
2. Submit order successfully
3. Check email queue

**Expected Result:**
- Order confirmation email is added to queue
- Email contains order details

**Traceability:**
- AC: [AC-ORDER-001](../l1/acceptance-criteria.md#ac-order-001)
- BR: [BR-ORDER](../l1/business-rules.md#br-order)

---

### TC-AC-ORDER-001-P03 – Cart cleared after successful order placement {#tc-ac-order-001-p03}

**Preconditions:**
- Verified customer with items in cart

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| cartItems | 3 products | Multiple items to verify all cleared |

**Steps:**
1. Add multiple items to cart
2. Complete checkout
3. View cart after order

**Expected Result:**
- Cart shows 0 items
- Cart total is $0.00

**Traceability:**
- AC: [AC-ORDER-001](../l1/acceptance-criteria.md#ac-order-001)
- BR: [BR-ORDER](../l1/business-rules.md#br-order)

---

### TC-AC-ORDER-002-P01 – Free shipping applied at $50 subtotal {#tc-ac-order-002-p01}

**Preconditions:**
- Customer at checkout

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| subtotal | $50.00 | Exact threshold amount |

**Steps:**
1. Add items totaling $50
2. Proceed to checkout
3. View shipping options

**Expected Result:**
- Free shipping is automatically applied
- Shipping cost shows $0.00

**Traceability:**
- AC: [AC-ORDER-002](../l1/acceptance-criteria.md#ac-order-002)
- BR: [BR-SHIPPING](../l1/business-rules.md#br-shipping)

---

### TC-AC-ORDER-002-P02 – Free shipping applied at $75 subtotal (above threshold) {#tc-ac-order-002-p02}

**Preconditions:**
- Customer at checkout

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| subtotal | $75.00 | Well above threshold |

**Steps:**
1. Add items totaling $75
2. Proceed to checkout

**Expected Result:**
- Free shipping automatically applied

**Traceability:**
- AC: [AC-ORDER-002](../l1/acceptance-criteria.md#ac-order-002)
- BR: [BR-SHIPPING](../l1/business-rules.md#br-shipping)

---

### TC-AC-ORDER-002-P03 – Free shipping based on post-discount subtotal {#tc-ac-order-002-p03}

**Preconditions:**
- Customer has discount code

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| originalSubtotal | $60.00 | Pre-discount total |
| discount | $5.00 | Applied discount |
| finalSubtotal | $55.00 | Still above $50 |

**Steps:**
1. Add $60 worth of items
2. Apply $5 discount
3. Check shipping

**Expected Result:**
- Free shipping applied based on $55 subtotal

**Traceability:**
- AC: [AC-ORDER-002](../l1/acceptance-criteria.md#ac-order-002)
- BR: [BR-SHIPPING](../l1/business-rules.md#br-shipping)

---

### TC-AC-ORDER-003-P01 – Tax calculated based on shipping state {#tc-ac-order-003-p01}

**Preconditions:**
- Customer at checkout with shipping address

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| shippingState | CA | California tax rate |
| subtotal | $100.00 | Easy calculation base |

**Steps:**
1. Add items to cart
2. Enter California shipping address
3. View order total

**Expected Result:**
- California state tax rate applied
- Tax amount shown in order total

**Traceability:**
- AC: [AC-ORDER-003](../l1/acceptance-criteria.md#ac-order-003)
- BR: [BR-TAX](../l1/business-rules.md#br-tax)

---

### TC-AC-ORDER-003-P02 – Tax recalculated when shipping address changes {#tc-ac-order-003-p02}

**Preconditions:**
- Customer has multiple saved addresses

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| initialState | NY | Higher tax state |
| newState | OR | No sales tax state |

**Steps:**
1. Select NY address
2. Note tax amount
3. Change to OR address
4. Note new tax

**Expected Result:**
- Tax recalculated for new state
- OR shows $0 tax

**Traceability:**
- AC: [AC-ORDER-003](../l1/acceptance-criteria.md#ac-order-003)
- BR: [BR-TAX](../l1/business-rules.md#br-tax)

---

### TC-AC-ORDER-004-P01 – Confirmation email queued with all required fields {#tc-ac-order-004-p01}

**Preconditions:**
- Order successfully placed

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| orderNumber | ORD-12345 | Generated order ID |
| itemCount | 3 | Multiple items ordered |

**Steps:**
1. Place order successfully
2. Check email queue
3. Verify email content

**Expected Result:**
- Email queued with order number
- Contains all items with quantities and prices
- Shows total and shipping address
- Includes estimated delivery

**Traceability:**
- AC: [AC-ORDER-004](../l1/acceptance-criteria.md#ac-order-004)
- BR: [BR-EMAIL](../l1/business-rules.md#br-email)

---

### TC-AC-ORDER-004-P02 – Email queued asynchronously without blocking order {#tc-ac-order-004-p02}

**Preconditions:**
- Order being placed

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| emailDelay | simulated | Slow email service |

**Steps:**
1. Place order
2. Measure order creation time
3. Verify email in queue

**Expected Result:**
- Order completes without waiting for email
- Email appears in async queue

**Traceability:**
- AC: [AC-ORDER-004](../l1/acceptance-criteria.md#ac-order-004)
- BR: [BR-EMAIL](../l1/business-rules.md#br-email)

---

### TC-AC-ORDER-005-P01 – View order history with all order details {#tc-ac-order-005-p01}

**Preconditions:**
- Customer has placed orders

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| orderCount | 5 | Multiple orders to display |

**Steps:**
1. Login as customer
2. Navigate to order history
3. View order list

**Expected Result:**
- All orders displayed
- Each shows number, date, items, status
- Shipping address visible
- Tracking number shown if available

**Traceability:**
- AC: [AC-ORDER-005](../l1/acceptance-criteria.md#ac-order-005)
- BR: [BR-TRACKING](../l1/business-rules.md#br-tracking)

---

### TC-AC-ORDER-005-P02 – Order history displays with pagination {#tc-ac-order-005-p02}

**Preconditions:**
- Customer has many orders

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| orderCount | 50 | Exceeds single page |

**Steps:**
1. Login as customer with 50 orders
2. View order history
3. Navigate pages

**Expected Result:**
- Orders split across pages
- Pagination controls visible
- Can navigate between pages

**Traceability:**
- AC: [AC-ORDER-005](../l1/acceptance-criteria.md#ac-order-005)
- BR: [BR-TRACKING](../l1/business-rules.md#br-tracking)

---

### TC-AC-ORDER-005-P03 – Tracking number shown for shipped orders {#tc-ac-order-005-p03}

**Preconditions:**
- Order has been shipped with tracking

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| trackingNumber | 1Z999AA10123456784 | UPS tracking |

**Steps:**
1. View shipped order in history

**Expected Result:**
- Tracking number displayed
- Tracking number is clickable/usable

**Traceability:**
- AC: [AC-ORDER-005](../l1/acceptance-criteria.md#ac-order-005)
- BR: [BR-TRACKING](../l1/business-rules.md#br-tracking)

---

### TC-AC-ORDER-006-P01 – Cancel pending order successfully {#tc-ac-order-006-p01}

**Preconditions:**
- Customer has existing order
- Order status is 'pending'

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| orderId | ORD-12345 | Valid pending order |
| orderStatus | pending | Cancellable state |

**Steps:**
1. Navigate to order details
2. Click cancel order
3. Confirm cancellation

**Expected Result:**
- Order status changes to 'cancelled'
- Success message displayed

**Traceability:**
- AC: [AC-ORDER-006](../l1/acceptance-criteria.md#ac-order-006)
- BR: [AMB-OP-028](../l1/business-rules.md#amb-op-028)

---

### TC-AC-ORDER-006-P02 – Cancel confirmed order successfully {#tc-ac-order-006-p02}

**Preconditions:**
- Customer has existing order
- Order status is 'confirmed'

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| orderId | ORD-12346 | Valid confirmed order |
| orderStatus | confirmed | Cancellable state |

**Steps:**
1. Navigate to order details
2. Click cancel order
3. Confirm cancellation

**Expected Result:**
- Order status changes to 'cancelled'

**Traceability:**
- AC: [AC-ORDER-006](../l1/acceptance-criteria.md#ac-order-006)
- BR: [AMB-OP-028](../l1/business-rules.md#amb-op-028)

---

### TC-AC-ORDER-006-P03 – Inventory restored on cancellation {#tc-ac-order-006-p03}

**Preconditions:**
- Order exists with 3 items
- Order status is 'pending'

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| productId | PROD-001 | Product in order |
| quantity | 3 | Ordered quantity |

**Steps:**
1. Note current inventory level
2. Cancel order
3. Check inventory level

**Expected Result:**
- Inventory increased by ordered quantity

**Traceability:**
- AC: [AC-ORDER-006](../l1/acceptance-criteria.md#ac-order-006)
- BR: [AMB-OP-029](../l1/business-rules.md#amb-op-029)

---

### TC-AC-ORDER-006-P04 – Refund initiated to original payment method {#tc-ac-order-006-p04}

**Preconditions:**
- Order paid via credit card
- Order status is 'pending'

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| paymentMethod | Visa ending 4242 | Original payment |
| orderTotal | 99.99 | Amount to refund |

**Steps:**
1. Cancel order
2. Check payment gateway logs

**Expected Result:**
- Refund request sent to original payment method

**Traceability:**
- AC: [AC-ORDER-006](../l1/acceptance-criteria.md#ac-order-006)
- BR: [AMB-OP-030](../l1/business-rules.md#amb-op-030)

---

### TC-AC-ORDER-006-P05 – Cancellation confirmation email sent {#tc-ac-order-006-p05}

**Preconditions:**
- Customer has valid email
- Order is cancellable

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| customerEmail | test@example.com | Customer email |

**Steps:**
1. Cancel order
2. Check email queue/logs

**Expected Result:**
- Cancellation confirmation email queued for delivery

**Traceability:**
- AC: [AC-ORDER-006](../l1/acceptance-criteria.md#ac-order-006)
- BR: [AMB-OP-031](../l1/business-rules.md#amb-op-031)

---

### TC-AC-ORDER-007-P01 – Item prices captured at order creation {#tc-ac-order-007-p01}

**Preconditions:**
- Product exists with price $29.99
- Customer has items in cart

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| productPrice | 29.99 | Current price |
| quantity | 2 | Ordered quantity |

**Steps:**
1. Place order
2. Check order details

**Expected Result:**
- Order shows captured price $29.99 per item

**Traceability:**
- AC: [AC-ORDER-007](../l1/acceptance-criteria.md#ac-order-007)
- BR: [AMB-ENT-019](../l1/business-rules.md#amb-ent-019)

---

### TC-AC-ORDER-007-P02 – Order prices unchanged after product price increase {#tc-ac-order-007-p02}

**Preconditions:**
- Order placed at $29.99
- Admin increases price to $39.99

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| originalPrice | 29.99 | Order time price |
| newPrice | 39.99 | Updated catalog price |

**Steps:**
1. Place order at $29.99
2. Admin updates price to $39.99
3. View order

**Expected Result:**
- Order still shows $29.99
- Product catalog shows $39.99

**Traceability:**
- AC: [AC-ORDER-007](../l1/acceptance-criteria.md#ac-order-007)
- BR: [AMB-ENT-019](../l1/business-rules.md#amb-ent-019)

---

### TC-AC-ORDER-007-P03 – Order prices unchanged after product price decrease {#tc-ac-order-007-p03}

**Preconditions:**
- Order placed at $29.99
- Admin decreases price to $19.99

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| originalPrice | 29.99 | Order time price |
| newPrice | 19.99 | Reduced catalog price |

**Steps:**
1. Place order at $29.99
2. Admin updates price to $19.99
3. View order

**Expected Result:**
- Order still shows $29.99

**Traceability:**
- AC: [AC-ORDER-007](../l1/acceptance-criteria.md#ac-order-007)
- BR: [AMB-ENT-019](../l1/business-rules.md#amb-ent-019)

---

### TC-AC-CUST-001-P01 – Register with valid email and compliant password {#tc-ac-cust-001-p01}

**Preconditions:**
- Email not previously registered

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| email | newuser@example.com | Valid email |
| password | SecurePass1 | 8+ chars with number |
| firstName | John | Required field |
| lastName | Doe | Required field |

**Steps:**
1. Navigate to registration
2. Fill all required fields
3. Submit form

**Expected Result:**
- Account created
- Success message displayed

**Traceability:**
- AC: [AC-CUST-001](../l1/acceptance-criteria.md#ac-cust-001)
- BR: [AMB-ENT-024](../l1/business-rules.md#amb-ent-024), [AMB-ENT-025](../l1/business-rules.md#amb-ent-025)

---

### TC-AC-CUST-001-P02 – Verification email sent after registration {#tc-ac-cust-001-p02}

**Preconditions:**
- Valid registration data

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| email | verify@example.com | Email to verify |

**Steps:**
1. Complete registration
2. Check email queue

**Expected Result:**
- Verification email queued
- Email contains verification link

**Traceability:**
- AC: [AC-CUST-001](../l1/acceptance-criteria.md#ac-cust-001)
- BR: [AMB-ENT-026](../l1/business-rules.md#amb-ent-026)

---

### TC-AC-CUST-001-P03 – User informed to verify email before ordering {#tc-ac-cust-001-p03}

**Preconditions:**
- User just registered

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| verificationStatus | pending | Not yet verified |

**Steps:**
1. Complete registration
2. View confirmation page

**Expected Result:**
- Message displayed about email verification requirement

**Traceability:**
- AC: [AC-CUST-001](../l1/acceptance-criteria.md#ac-cust-001)
- BR: [AMB-ENT-030](../l1/business-rules.md#amb-ent-030)

---

### TC-AC-CUST-002-P01 – Add shipping address with all required fields {#tc-ac-cust-002-p01}

**Preconditions:**
- Customer is logged in
- Profile management page open

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| street | 123 Main St | Valid street |
| city | New York | Valid city |
| state | NY | Valid state |
| postalCode | 10001 | Valid postal |
| country | USA | Valid country |
| recipientName | John Doe | Required |

**Steps:**
1. Navigate to addresses
2. Click add address
3. Fill all fields
4. Save

**Expected Result:**
- Address saved to account
- Address appears in list

**Traceability:**
- AC: [AC-CUST-002](../l1/acceptance-criteria.md#ac-cust-002)
- BR: [AMB-ENT-027](../l1/business-rules.md#amb-ent-027), [AMB-ENT-031](../l1/business-rules.md#amb-ent-031)

---

### TC-AC-CUST-002-P02 – Set address as default {#tc-ac-cust-002-p02}

**Preconditions:**
- Customer has at least one address

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| isDefault | true | Set as default |

**Steps:**
1. Add new address
2. Check 'set as default'
3. Save

**Expected Result:**
- Address marked as default
- Previous default unmarked

**Traceability:**
- AC: [AC-CUST-002](../l1/acceptance-criteria.md#ac-cust-002)
- BR: [AMB-ENT-032](../l1/business-rules.md#amb-ent-032)

---

### TC-AC-CUST-002-P03 – Add multiple addresses to account {#tc-ac-cust-002-p03}

**Preconditions:**
- Customer has one existing address

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| addressCount | 2 | Second address |

**Steps:**
1. Add second address
2. Save

**Expected Result:**
- Both addresses visible in list

**Traceability:**
- AC: [AC-CUST-002](../l1/acceptance-criteria.md#ac-cust-002)
- BR: [AMB-ENT-027](../l1/business-rules.md#amb-ent-027)

---

### TC-AC-CUST-003-P01 – Add Visa card via Stripe tokenization {#tc-ac-cust-003-p01}

**Preconditions:**
- Customer logged in
- Stripe integration active

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| cardType | Visa | Supported type |
| cardNumber | 4242424242424242 | Stripe test card |

**Steps:**
1. Navigate to payment methods
2. Add card
3. Enter Visa details
4. Save

**Expected Result:**
- Card tokenized via Stripe
- Saved to account

**Traceability:**
- AC: [AC-CUST-003](../l1/acceptance-criteria.md#ac-cust-003)
- BR: [AMB-ENT-028](../l1/business-rules.md#amb-ent-028), [AMB-ENT-034](../l1/business-rules.md#amb-ent-034)

---

### TC-AC-CUST-003-P02 – Add PayPal payment method {#tc-ac-cust-003-p02}

**Preconditions:**
- Customer logged in
- PayPal integration active

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| paymentType | PayPal | PayPal option |

**Steps:**
1. Select PayPal
2. Complete PayPal authorization
3. Save

**Expected Result:**
- PayPal linked to account

**Traceability:**
- AC: [AC-CUST-003](../l1/acceptance-criteria.md#ac-cust-003)
- BR: [AMB-ENT-028](../l1/business-rules.md#amb-ent-028)

---

### TC-AC-CUST-003-P03 – Set payment method as default {#tc-ac-cust-003-p03}

**Preconditions:**
- Customer has existing payment method

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| isDefault | true | Set as default |

**Steps:**
1. Add new payment method
2. Check 'set as default'
3. Save

**Expected Result:**
- New method is default
- Previous default unmarked

**Traceability:**
- AC: [AC-CUST-003](../l1/acceptance-criteria.md#ac-cust-003)
- BR: [AMB-ENT-036](../l1/business-rules.md#amb-ent-036)

---

### TC-AC-CUST-003-P04 – Add Mastercard via Stripe {#tc-ac-cust-003-p04}

**Preconditions:**
- Customer logged in

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| cardType | Mastercard | Supported type |
| cardNumber | 5555555555554444 | MC test card |

**Steps:**
1. Add Mastercard details
2. Save

**Expected Result:**
- Mastercard tokenized and saved

**Traceability:**
- AC: [AC-CUST-003](../l1/acceptance-criteria.md#ac-cust-003)
- BR: [AMB-ENT-034](../l1/business-rules.md#amb-ent-034)

---

### TC-AC-CUST-004-P01 – Soft delete customer account on preference selection {#tc-ac-cust-004-p01}

**Preconditions:**
- Customer is registered
- Customer has no pending orders

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| customerId | cust-123 | Valid customer |
| deletionType | soft_delete | Deactivation preference |

**Steps:**
1. Login as customer
2. Navigate to account settings
3. Request account deletion
4. Select soft delete option
5. Confirm deletion

**Expected Result:**
- Account status changed to deactivated
- Customer cannot login
- Data preserved in system

**Traceability:**
- AC: [AC-CUST-004](../l1/acceptance-criteria.md#ac-cust-004)
- BR: [BR-CUST](../l1/business-rules.md#br-cust)

---

### TC-AC-CUST-004-P02 – Full GDPR erasure on customer preference {#tc-ac-cust-004-p02}

**Preconditions:**
- Customer is registered
- No pending orders exist

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| customerId | cust-456 | Valid customer |
| deletionType | full_erasure | GDPR erasure preference |

**Steps:**
1. Login as customer
2. Navigate to account settings
3. Request account deletion
4. Select full erasure option
5. Confirm deletion

**Expected Result:**
- All personal data permanently deleted
- Anonymized records for order history
- Confirmation sent

**Traceability:**
- AC: [AC-CUST-004](../l1/acceptance-criteria.md#ac-cust-004)
- BR: [BR-CUST](../l1/business-rules.md#br-cust)

---

### TC-AC-CUST-004-P03 – Deletion confirmation required before processing {#tc-ac-cust-004-p03}

**Preconditions:**
- Customer logged in
- No pending orders

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| customerId | cust-789 | Valid customer |

**Steps:**
1. Request account deletion
2. System displays confirmation dialog
3. Customer confirms

**Expected Result:**
- Deletion only proceeds after explicit confirmation

**Traceability:**
- AC: [AC-CUST-004](../l1/acceptance-criteria.md#ac-cust-004)
- BR: [BR-CUST](../l1/business-rules.md#br-cust)

---

### TC-AC-ADMIN-001-P01 – Create product with all valid required fields {#tc-ac-admin-001-p01}

**Preconditions:**
- Admin logged in
- On product management interface

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| name | Test Product | Valid 12 chars |
| description | Product description | Valid length |
| price | 19.99 | Above minimum |
| category | Electronics | Valid category |

**Steps:**
1. Navigate to add product
2. Enter valid name
3. Enter description
4. Set price
5. Select category
6. Submit

**Expected Result:**
- Product created with UUID
- Status is draft
- Creator logged
- Timestamp recorded

**Traceability:**
- AC: [AC-ADMIN-001](../l1/acceptance-criteria.md#ac-admin-001)
- BR: [BR-ADMIN](../l1/business-rules.md#br-admin), [BR-PROD](../l1/business-rules.md#br-prod)

---

### TC-AC-ADMIN-001-P02 – Create product with images and primary image set {#tc-ac-admin-001-p02}

**Preconditions:**
- Admin logged in

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| name | Product With Images | Valid name |
| images | 5 images | Within limit |
| primaryImage | image_3 | Third image as primary |

**Steps:**
1. Create product
2. Upload 5 images
3. Mark third as primary
4. Submit

**Expected Result:**
- Product created
- 5 images attached
- Primary image correctly set

**Traceability:**
- AC: [AC-ADMIN-001](../l1/acceptance-criteria.md#ac-admin-001)
- BR: [BR-ADMIN](../l1/business-rules.md#br-admin), [BR-PROD](../l1/business-rules.md#br-prod)

---

### TC-AC-ADMIN-001-P03 – Duplicate name shows warning but allows creation {#tc-ac-admin-001-p03}

**Preconditions:**
- Product with same name exists

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| name | Existing Product | Duplicate name |
| price | 29.99 | Valid price |

**Steps:**
1. Enter duplicate product name
2. Submit product

**Expected Result:**
- Warning displayed
- Product created successfully
- Both products exist

**Traceability:**
- AC: [AC-ADMIN-001](../l1/acceptance-criteria.md#ac-admin-001)
- BR: [BR-ADMIN](../l1/business-rules.md#br-admin), [BR-PROD](../l1/business-rules.md#br-prod)

---

### TC-AC-ADMIN-002-P01 – Update product attributes successfully {#tc-ac-admin-002-p01}

**Preconditions:**
- Product exists
- Admin logged in

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| productId | prod-123 | Existing product |
| newPrice | 29.99 | Updated price |

**Steps:**
1. Open product for editing
2. Modify price
3. Save changes

**Expected Result:**
- Product updated
- Modified by admin recorded
- Timestamp updated

**Traceability:**
- AC: [AC-ADMIN-002](../l1/acceptance-criteria.md#ac-admin-002)
- BR: [BR-ADMIN](../l1/business-rules.md#br-admin), [BR-PROD](../l1/business-rules.md#br-prod)

---

### TC-AC-ADMIN-002-P02 – Price change logged in audit trail {#tc-ac-admin-002-p02}

**Preconditions:**
- Product exists with price $19.99

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| oldPrice | 19.99 | Original price |
| newPrice | 24.99 | New price |

**Steps:**
1. Edit product
2. Change price from $19.99 to $24.99
3. Save

**Expected Result:**
- Audit trail shows old and new price
- Change timestamp recorded

**Traceability:**
- AC: [AC-ADMIN-002](../l1/acceptance-criteria.md#ac-admin-002)
- BR: [BR-ADMIN](../l1/business-rules.md#br-admin), [BR-PROD](../l1/business-rules.md#br-prod)

---

### TC-AC-ADMIN-002-P03 – Cart prices reflect updated product price {#tc-ac-admin-002-p03}

**Preconditions:**
- Product in customer cart
- Admin updates price

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| oldPrice | 19.99 | Price when added to cart |
| newPrice | 24.99 | Updated price |

**Steps:**
1. Customer adds product to cart
2. Admin updates price
3. Customer views cart

**Expected Result:**
- Cart shows new price $24.99
- Cart total recalculated

**Traceability:**
- AC: [AC-ADMIN-002](../l1/acceptance-criteria.md#ac-admin-002)
- BR: [BR-ADMIN](../l1/business-rules.md#br-admin), [BR-PROD](../l1/business-rules.md#br-prod)

---

### TC-AC-ADMIN-003-P01 – Soft delete product preserves order history {#tc-ac-admin-003-p01}

**Preconditions:**
- Product exists
- Product has order history

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| productId | prod-456 | Product with orders |

**Steps:**
1. Select product to delete
2. Click delete
3. Confirm deletion

**Expected Result:**
- Product soft-deleted
- Order history intact
- Product data preserved

**Traceability:**
- AC: [AC-ADMIN-003](../l1/acceptance-criteria.md#ac-admin-003)
- BR: [BR-ADMIN](../l1/business-rules.md#br-admin), [BR-PROD](../l1/business-rules.md#br-prod)

---

### TC-AC-ADMIN-003-P02 – Deleted product removed from customer carts {#tc-ac-admin-003-p02}

**Preconditions:**
- Product in customer cart
- Admin deletes product

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| productId | prod-789 | Product in carts |

**Steps:**
1. Delete product with items in carts
2. Customer views cart

**Expected Result:**
- Product removed from cart
- Notification shown to customer

**Traceability:**
- AC: [AC-ADMIN-003](../l1/acceptance-criteria.md#ac-admin-003)
- BR: [BR-ADMIN](../l1/business-rules.md#br-admin), [BR-PROD](../l1/business-rules.md#br-prod)

---

### TC-AC-ADMIN-003-P03 – Deleted product no longer in listings {#tc-ac-admin-003-p03}

**Preconditions:**
- Product exists in catalog

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| productId | prod-101 | Visible product |

**Steps:**
1. Delete product
2. Browse product catalog
3. Search for product

**Expected Result:**
- Product not in listings
- Product not in search results

**Traceability:**
- AC: [AC-ADMIN-003](../l1/acceptance-criteria.md#ac-admin-003)
- BR: [BR-ADMIN](../l1/business-rules.md#br-admin), [BR-PROD](../l1/business-rules.md#br-prod)

---

### TC-AC-ADMIN-004-P01 – Display orders with all required fields {#tc-ac-admin-004-p01}

**Preconditions:**
- Orders exist in system
- Admin logged in

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| orderCount | 10 | Multiple orders |

**Steps:**
1. Navigate to order management
2. View order list

**Expected Result:**
- Each order shows number
- Shows date
- Shows customer
- Shows total
- Shows status

**Traceability:**
- AC: [AC-ADMIN-004](../l1/acceptance-criteria.md#ac-admin-004)
- BR: [BR-ADMIN](../l1/business-rules.md#br-admin), [BR-ORDER](../l1/business-rules.md#br-order)

---

### TC-AC-ADMIN-004-P02 – Filter orders by status {#tc-ac-admin-004-p02}

**Preconditions:**
- Orders with various statuses exist

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| statusFilter | pending | Filter value |

**Steps:**
1. View order list
2. Apply status filter for pending

**Expected Result:**
- Only pending orders displayed
- Other statuses hidden

**Traceability:**
- AC: [AC-ADMIN-004](../l1/acceptance-criteria.md#ac-admin-004)
- BR: [BR-ADMIN](../l1/business-rules.md#br-admin), [BR-ORDER](../l1/business-rules.md#br-order)

---

### TC-AC-ADMIN-004-P03 – Filter orders by date range {#tc-ac-admin-004-p03}

**Preconditions:**
- Orders from multiple dates exist

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| startDate | 2024-01-01 | Range start |
| endDate | 2024-01-31 | Range end |

**Steps:**
1. Set date range filter
2. Apply filter

**Expected Result:**
- Only orders within date range shown

**Traceability:**
- AC: [AC-ADMIN-004](../l1/acceptance-criteria.md#ac-admin-004)
- BR: [BR-ADMIN](../l1/business-rules.md#br-admin), [BR-ORDER](../l1/business-rules.md#br-order)

---

### TC-AC-ADMIN-004-P04 – Search orders by order number {#tc-ac-admin-004-p04}

**Preconditions:**
- Order ORD-12345 exists

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| searchTerm | ORD-12345 | Exact order number |

**Steps:**
1. Enter order number in search
2. Execute search

**Expected Result:**
- Matching order displayed
- Quick access to order details

**Traceability:**
- AC: [AC-ADMIN-004](../l1/acceptance-criteria.md#ac-admin-004)
- BR: [BR-ADMIN](../l1/business-rules.md#br-admin), [BR-ORDER](../l1/business-rules.md#br-order)

---

### TC-AC-ADMIN-005-P01 – Update order from pending to confirmed successfully {#tc-ac-admin-005-p01}

**Preconditions:**
- Admin is authenticated
- Order exists with pending status

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| orderId | ORD-12345 | Valid order ID |
| currentStatus | pending | Starting status |
| newStatus | confirmed | Valid next status |

**Steps:**
1. Navigate to order details
2. Select confirmed status
3. Submit change

**Expected Result:**
- Status updated to confirmed
- Change logged with admin and timestamp
- Customer receives email notification

**Traceability:**
- AC: [AC-ADMIN-005](../l1/acceptance-criteria.md#ac-admin-005)
- BR: [AMB-ENT-018](../l1/business-rules.md#amb-ent-018), [AMB-OP-047](../l1/business-rules.md#amb-op-047)

---

### TC-AC-ADMIN-005-P02 – Update order from shipped to delivered with email {#tc-ac-admin-005-p02}

**Preconditions:**
- Admin is authenticated
- Order exists with shipped status

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| orderId | ORD-12346 | Valid order ID |
| currentStatus | shipped | Starting status |
| newStatus | delivered | Final status |

**Steps:**
1. Navigate to order details
2. Select delivered status
3. Submit change

**Expected Result:**
- Status updated to delivered
- Audit log entry created
- Delivery notification email sent to customer

**Traceability:**
- AC: [AC-ADMIN-005](../l1/acceptance-criteria.md#ac-admin-005)
- BR: [AMB-OP-048](../l1/business-rules.md#amb-op-048), [AMB-OP-049](../l1/business-rules.md#amb-op-049)

---

### TC-AC-ADMIN-005-P03 – Set shipped status without tracking number {#tc-ac-admin-005-p03}

**Preconditions:**
- Admin is authenticated
- Order exists with confirmed status

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| orderId | ORD-12347 | Valid order ID |
| newStatus | shipped | Target status |
| trackingNumber |  | Optional field left empty |

**Steps:**
1. Navigate to order
2. Select shipped status
3. Leave tracking empty
4. Submit

**Expected Result:**
- Status updated to shipped successfully
- No validation error for missing tracking

**Traceability:**
- AC: [AC-ADMIN-005](../l1/acceptance-criteria.md#ac-admin-005)
- BR: [AMB-OP-050](../l1/business-rules.md#amb-op-050)

---

### TC-AC-ADMIN-006-P01 – Set absolute stock value successfully {#tc-ac-admin-006-p01}

**Preconditions:**
- Admin authenticated
- Product variant exists with stock 50

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| variantId | VAR-001 | Valid variant |
| adjustmentType | absolute | Set exact value |
| newValue | 100 | New stock level |

**Steps:**
1. Navigate to inventory
2. Select variant
3. Set absolute value 100
4. Save

**Expected Result:**
- Stock updated to 100
- Audit trail shows before:50 after:100

**Traceability:**
- AC: [AC-ADMIN-006](../l1/acceptance-criteria.md#ac-admin-006)
- BR: [AMB-ENT-009](../l1/business-rules.md#amb-ent-009), [AMB-OP-051](../l1/business-rules.md#amb-op-051)

---

### TC-AC-ADMIN-006-P02 – Apply delta adjustment to increase stock {#tc-ac-admin-006-p02}

**Preconditions:**
- Admin authenticated
- Variant has stock of 30

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| variantId | VAR-002 | Valid variant |
| adjustmentType | delta | Relative change |
| delta | +20 | Increase by 20 |
| reason | Restock delivery | Optional reason |

**Steps:**
1. Select variant
2. Choose delta adjustment
3. Enter +20
4. Add reason
5. Save

**Expected Result:**
- Stock updated to 50
- Audit shows reason and before/after values

**Traceability:**
- AC: [AC-ADMIN-006](../l1/acceptance-criteria.md#ac-admin-006)
- BR: [AMB-OP-052](../l1/business-rules.md#amb-op-052)

---

### TC-AC-ADMIN-006-P03 – Trigger low stock alert when threshold crossed {#tc-ac-admin-006-p03}

**Preconditions:**
- Variant has stock 15
- Low stock threshold is 10

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| currentStock | 15 | Above threshold |
| delta | -10 | Reduces to 5 |
| threshold | 10 | Alert trigger point |

**Steps:**
1. Reduce stock by 10
2. Save adjustment

**Expected Result:**
- Stock reduced to 5
- Low stock alert triggered

**Traceability:**
- AC: [AC-ADMIN-006](../l1/acceptance-criteria.md#ac-admin-006)
- BR: [AMB-ENT-040](../l1/business-rules.md#amb-ent-040), [AMB-OP-053](../l1/business-rules.md#amb-op-053)

---

### TC-AC-VAR-001-P01 – Display all variant options for product with size/color {#tc-ac-var-001-p01}

**Preconditions:**
- Product exists with size (S,M,L) and color (Red,Blue) variants

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| productId | PROD-SHIRT | Product with variants |
| sizes | S,M,L | Size options |
| colors | Red,Blue | Color options |

**Steps:**
1. Navigate to product page
2. View variant selection area

**Expected Result:**
- All size options displayed
- All color options displayed

**Traceability:**
- AC: [AC-VAR-001](../l1/acceptance-criteria.md#ac-var-001)
- BR: [AMB-ENT-008](../l1/business-rules.md#amb-ent-008)

---

### TC-AC-VAR-001-P02 – Show individual stock status per variant {#tc-ac-var-001-p02}

**Preconditions:**
- Product has variants with different stock levels

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| variant1 | S-Red | In stock (qty 10) |
| variant2 | M-Blue | Out of stock (qty 0) |

**Steps:**
1. View product
2. Check each variant stock indicator

**Expected Result:**
- S-Red shows in stock
- M-Blue shows out of stock

**Traceability:**
- AC: [AC-VAR-001](../l1/acceptance-criteria.md#ac-var-001)
- BR: [AMB-ENT-009](../l1/business-rules.md#amb-ent-009)

---

### TC-AC-VAR-001-P03 – Display variant-specific price override {#tc-ac-var-001-p03}

**Preconditions:**
- Product base price $20
- XL variant has +$5 override

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| basePrice | 20.00 | Product base price |
| xlOverride | 25.00 | XL variant price |

**Steps:**
1. View product
2. Select XL variant
3. Check displayed price

**Expected Result:**
- XL variant shows $25.00 price

**Traceability:**
- AC: [AC-VAR-001](../l1/acceptance-criteria.md#ac-var-001)
- BR: [AMB-ENT-010](../l1/business-rules.md#amb-ent-010)

---

### TC-AC-CAT-001-P01 – Display three-level category hierarchy {#tc-ac-cat-001-p01}

**Preconditions:**
- Categories exist: Electronics > Phones > Smartphones

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| level1 | Electronics | Top level |
| level2 | Phones | Second level |
| level3 | Smartphones | Third level |

**Steps:**
1. Navigate to category browser
2. Expand Electronics
3. Expand Phones

**Expected Result:**
- All three levels visible
- Proper parent-child relationships shown

**Traceability:**
- AC: [AC-CAT-001](../l1/acceptance-criteria.md#ac-cat-001)
- BR: [AMB-ENT-006](../l1/business-rules.md#amb-ent-006), [AMB-ENT-037](../l1/business-rules.md#amb-ent-037)

---

### TC-AC-CAT-001-P02 – Show product with primary and secondary categories {#tc-ac-cat-001-p02}

**Preconditions:**
- Product in primary: Phones, secondary: Accessories

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| productId | PROD-CASE | Phone case product |
| primary | Accessories | Main category |
| secondary | Phones | Additional category |

**Steps:**
1. View product details
2. Check category assignments

**Expected Result:**
- Primary category displayed
- Secondary categories shown

**Traceability:**
- AC: [AC-CAT-001](../l1/acceptance-criteria.md#ac-cat-001)
- BR: [AMB-ENT-038](../l1/business-rules.md#amb-ent-038), [AMB-ENT-039](../l1/business-rules.md#amb-ent-039)

---

### TC-AC-CAT-001-P03 – Browse products by selecting subcategory {#tc-ac-cat-001-p03}

**Preconditions:**
- Smartphones category has 5 products

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| category | Smartphones | Leaf category |
| productCount | 5 | Expected products |

**Steps:**
1. Navigate to Smartphones category
2. View product listing

**Expected Result:**
- 5 products displayed
- Products belong to selected category

**Traceability:**
- AC: [AC-CAT-001](../l1/acceptance-criteria.md#ac-cat-001)
- BR: [AMB-ENT-006](../l1/business-rules.md#amb-ent-006)

---

## Negative Tests (Error Cases)

### TC-AC-PROD-001-N01 – Show empty state when no products match filters {#tc-ac-prod-001-n01}

**Preconditions:**
- Products exist but none match filter combination

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| category | NonExistent | No matching products |

**Steps:**
1. Apply filters that match no products
2. Observe result

**Expected Result:**
- 'No products found' message displayed
- Clear filters option shown

**Traceability:**
- AC: [AC-PROD-001](../l1/acceptance-criteria.md#ac-prod-001)
- BR: [BR-PROD](../l1/business-rules.md#br-prod)

---

### TC-AC-PROD-001-N02 – Reject invalid price range where min exceeds max {#tc-ac-prod-001-n02}

**Preconditions:**
- User is on product listing page

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| priceMin | 100 | Higher than max |
| priceMax | 50 | Lower than min |

**Steps:**
1. Enter min price 100
2. Enter max price 50
3. Apply filters

**Expected Result:**
- Validation error displayed
- Filters not applied

**Traceability:**
- AC: [AC-PROD-001](../l1/acceptance-criteria.md#ac-prod-001)
- BR: [BR-PROD](../l1/business-rules.md#br-prod)

---

### TC-AC-PROD-001-N03 – Reject negative price values in filter {#tc-ac-prod-001-n03}

**Preconditions:**
- User is on product listing page

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| priceMin | -10 | Invalid negative price |

**Steps:**
1. Enter negative price in min field
2. Attempt to apply filters

**Expected Result:**
- Validation error for invalid price
- Filter not applied

**Traceability:**
- AC: [AC-PROD-001](../l1/acceptance-criteria.md#ac-prod-001)
- BR: [BR-PROD](../l1/business-rules.md#br-prod)

---

### TC-AC-PROD-002-N01 – Cannot add out of stock product via direct API call {#tc-ac-prod-002-n01}

**Preconditions:**
- Product has zero inventory
- User attempts API bypass

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| productId | PROD-OOS | Out of stock item |
| quantity | 1 | Attempted quantity |

**Steps:**
1. Send direct add-to-cart API request for out-of-stock item

**Expected Result:**
- Request rejected
- Appropriate error returned

**Traceability:**
- AC: [AC-PROD-002](../l1/acceptance-criteria.md#ac-prod-002)
- BR: [BR-PROD](../l1/business-rules.md#br-prod)

---

### TC-AC-PROD-002-N02 – Out of stock indicator persists after page refresh {#tc-ac-prod-002-n02}

**Preconditions:**
- Product is out of stock
- User views product

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| inventory | 0 | Remains zero |

**Steps:**
1. View out-of-stock product
2. Refresh page
3. Check indicator

**Expected Result:**
- Out of Stock indicator still visible
- Button still disabled

**Traceability:**
- AC: [AC-PROD-002](../l1/acceptance-criteria.md#ac-prod-002)
- BR: [BR-PROD](../l1/business-rules.md#br-prod)

---

### TC-AC-CART-001-N01 – Add to cart disabled for out of stock product {#tc-ac-cart-001-n01}

**Preconditions:**
- Product has zero inventory

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| inventory | 0 | Out of stock |

**Steps:**
1. View out-of-stock product
2. Check Add to Cart button

**Expected Result:**
- Button disabled
- 'Out of Stock' message shown

**Traceability:**
- AC: [AC-CART-001](../l1/acceptance-criteria.md#ac-cart-001)
- BR: [BR-CART](../l1/business-rules.md#br-cart)

---

### TC-AC-CART-001-N02 – Quantity limited when exceeding available stock {#tc-ac-cart-001-n02}

**Preconditions:**
- Product has limited stock

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| inventory | 5 | Limited stock |
| requestedQty | 10 | Exceeds stock |

**Steps:**
1. Set quantity to 10
2. Click Add to Cart

**Expected Result:**
- Quantity limited to 5
- Notification about limit shown

**Traceability:**
- AC: [AC-CART-001](../l1/acceptance-criteria.md#ac-cart-001)
- BR: [BR-CART](../l1/business-rules.md#br-cart)

---

### TC-AC-CART-001-N03 – Add to cart disabled until variant selected {#tc-ac-cart-001-n03}

**Preconditions:**
- Product requires variant selection

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| hasVariants | true | Variants required |
| selectedVariant | null | None selected |

**Steps:**
1. View product with variants
2. Do not select variant
3. Check button

**Expected Result:**
- Add to Cart button disabled

**Traceability:**
- AC: [AC-CART-001](../l1/acceptance-criteria.md#ac-cart-001)
- BR: [BR-CART](../l1/business-rules.md#br-cart)

---

### TC-AC-CART-002-N01 – Limit combined quantity to available stock {#tc-ac-cart-002-n01}

**Preconditions:**
- Product A in cart qty 8
- Stock is 10

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| existingQty | 8 | Already in cart |
| addQty | 5 | Would exceed stock |
| stock | 10 | Max available |

**Steps:**
1. Try to add 5 more of Product A

**Expected Result:**
- Quantity limited to 10
- Notification shown about limit

**Traceability:**
- AC: [AC-CART-002](../l1/acceptance-criteria.md#ac-cart-002)
- BR: [BR-CART](../l1/business-rules.md#br-cart)

---

### TC-AC-CART-002-N02 – Cannot increment beyond stock even with multiple attempts {#tc-ac-cart-002-n02}

**Preconditions:**
- Cart at max stock for product

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| existingQty | 10 | At max stock |
| stock | 10 | Max available |

**Steps:**
1. Try to add 1 more when at max
2. Try again

**Expected Result:**
- Quantity stays at 10
- Notification each attempt

**Traceability:**
- AC: [AC-CART-002](../l1/acceptance-criteria.md#ac-cart-002)
- BR: [BR-CART](../l1/business-rules.md#br-cart)

---

### TC-AC-CART-003-N01 – Cart cleared after 7 days of inactivity {#tc-ac-cart-003-n01}

**Preconditions:**
- Guest cart exists
- 7 days passed

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| daysSinceActivity | 8 | Beyond 7 day limit |

**Steps:**
1. Create guest cart
2. Wait 8 days
3. Return to site

**Expected Result:**
- Cart is empty
- User sees empty cart

**Traceability:**
- AC: [AC-CART-003](../l1/acceptance-criteria.md#ac-cart-003)
- BR: [BR-CART](../l1/business-rules.md#br-cart)

---

### TC-AC-CART-003-N02 – Guest cannot access another guest session cart {#tc-ac-cart-003-n02}

**Preconditions:**
- Two separate guest sessions exist

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| session1 | ABC123 | First guest |
| session2 | XYZ789 | Second guest |

**Steps:**
1. Guest 1 adds items
2. Guest 2 opens site
3. Check Guest 2 cart

**Expected Result:**
- Guest 2 has empty cart
- Sessions isolated

**Traceability:**
- AC: [AC-CART-003](../l1/acceptance-criteria.md#ac-cart-003)
- BR: [BR-CART](../l1/business-rules.md#br-cart)

---

### TC-AC-CART-004-N01 – Merged quantity exceeds stock - limit to available {#tc-ac-cart-004-n01}

**Preconditions:**
- Guest has SKU-A:8
- Account has SKU-A:7
- Stock is 10

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| guestQty | 8 | Guest quantity |
| accountQty | 7 | Account quantity |
| availableStock | 10 | Max available |

**Steps:**
1. Setup carts totaling 15
2. Login
3. Verify quantity limited

**Expected Result:**
- Cart contains SKU-A:10
- Notification shown about limit

**Traceability:**
- AC: [AC-CART-004](../l1/acceptance-criteria.md#ac-cart-004)
- BR: [AMB-ENT-013](../l1/business-rules.md#amb-ent-013)

---

### TC-AC-CART-004-N02 – Show notification when merge quantity is limited {#tc-ac-cart-004-n02}

**Preconditions:**
- Combined quantity would exceed stock

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| combinedQty | 15 | Exceeds stock |
| stock | 10 | Available |

**Steps:**
1. Trigger merge that exceeds stock
2. Verify notification

**Expected Result:**
- User notification displayed about stock limit

**Traceability:**
- AC: [AC-CART-004](../l1/acceptance-criteria.md#ac-cart-004)
- BR: [AMB-ENT-013](../l1/business-rules.md#amb-ent-013)

---

### TC-AC-CART-005-N01 – Setting quantity to 0 removes item from cart {#tc-ac-cart-005-n01}

**Preconditions:**
- Cart has SKU-A qty 2

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| newQuantity | 0 | Zero removes item |

**Steps:**
1. Set quantity to 0
2. Verify cart

**Expected Result:**
- Item removed from cart
- Cart total excludes item

**Traceability:**
- AC: [AC-CART-005](../l1/acceptance-criteria.md#ac-cart-005)
- BR: [AMB-OP-011](../l1/business-rules.md#amb-op-011)

---

### TC-AC-CART-005-N02 – Quantity exceeding stock limited with notification {#tc-ac-cart-005-n02}

**Preconditions:**
- Cart has SKU-A
- Stock is 10

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| requestedQty | 15 | Over stock |
| stock | 10 | Available |

**Steps:**
1. Try to set quantity to 15
2. Verify behavior

**Expected Result:**
- Quantity set to 10
- Notification shown

**Traceability:**
- AC: [AC-CART-005](../l1/acceptance-criteria.md#ac-cart-005)
- BR: [AMB-OP-012](../l1/business-rules.md#amb-op-012)

---

### TC-AC-CART-005-N03 – Setting same quantity is idempotent - no change {#tc-ac-cart-005-n03}

**Preconditions:**
- Cart has SKU-A qty 3

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| currentQty | 3 | Current |
| newQty | 3 | Same value |

**Steps:**
1. Set quantity to current value 3
2. Verify no change

**Expected Result:**
- Quantity remains 3
- No unnecessary updates triggered

**Traceability:**
- AC: [AC-CART-005](../l1/acceptance-criteria.md#ac-cart-005)
- BR: [AMB-OP-013](../l1/business-rules.md#amb-op-013)

---

### TC-AC-CART-006-N01 – Last item removed shows empty cart state {#tc-ac-cart-006-n01}

**Preconditions:**
- Cart has only one item

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| cartItems | [SKU-A:1] | Single item |

**Steps:**
1. Remove the only item
2. Verify empty cart state

**Expected Result:**
- Empty cart state displayed

**Traceability:**
- AC: [AC-CART-006](../l1/acceptance-criteria.md#ac-cart-006)
- BR: [AMB-OP-014](../l1/business-rules.md#amb-op-014)

---

### TC-AC-CART-006-N02 – Empty cart shows Continue Shopping link {#tc-ac-cart-006-n02}

**Preconditions:**
- Cart is empty after removal

**Steps:**
1. Remove last item
2. Verify Continue Shopping link

**Expected Result:**
- Continue Shopping link visible
- Link navigates to shop

**Traceability:**
- AC: [AC-CART-006](../l1/acceptance-criteria.md#ac-cart-006)
- BR: [AMB-OP-014](../l1/business-rules.md#amb-op-014)

---

### TC-AC-CART-007-N01 – Checkout blocked for out of stock item {#tc-ac-cart-007-n01}

**Preconditions:**
- Cart has out of stock item

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| item | SKU-A | Out of stock |

**Steps:**
1. Try to proceed to checkout with out of stock item

**Expected Result:**
- Checkout blocked
- User informed about blocked item

**Traceability:**
- AC: [AC-CART-007](../l1/acceptance-criteria.md#ac-cart-007)
- BR: [AMB-ENT-014](../l1/business-rules.md#amb-ent-014)

---

### TC-AC-CART-007-N02 – Cannot increase quantity of out of stock item {#tc-ac-cart-007-n02}

**Preconditions:**
- Item has 0 stock
- Item qty in cart is 2

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| currentQty | 2 | Cart quantity |
| stock | 0 | No stock |

**Steps:**
1. Try to increase quantity of out of stock item

**Expected Result:**
- Quantity increase rejected or limited

**Traceability:**
- AC: [AC-CART-007](../l1/acceptance-criteria.md#ac-cart-007)
- BR: [AMB-ENT-014](../l1/business-rules.md#amb-ent-014)

---

### TC-AC-CART-008-N01 – Old price not used for cart total {#tc-ac-cart-008-n01}

**Preconditions:**
- Price changed from $10 to $15

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| qty | 2 | Quantity |
| newPrice | 15.00 | Current price |

**Steps:**
1. Check cart total calculation

**Expected Result:**
- Total is $30 (2x$15)
- Not $20 (old price)

**Traceability:**
- AC: [AC-CART-008](../l1/acceptance-criteria.md#ac-cart-008)
- BR: [AMB-ENT-016](../l1/business-rules.md#amb-ent-016)

---

### TC-AC-CART-008-N02 – Price change notification distinguishes increase vs decrease {#tc-ac-cart-008-n02}

**Preconditions:**
- Price changed

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| changeType | increase | Direction matters |

**Steps:**
1. View notification for price increase

**Expected Result:**
- Notification clearly indicates price went up

**Traceability:**
- AC: [AC-CART-008](../l1/acceptance-criteria.md#ac-cart-008)
- BR: [AMB-ENT-016](../l1/business-rules.md#amb-ent-016)

---

### TC-AC-ORDER-001-N01 – Reject order from unregistered customer {#tc-ac-order-001-n01}

**Preconditions:**
- User is not logged in
- Guest session with cart items

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| sessionType | guest | Anonymous user |

**Steps:**
1. Add items to cart as guest
2. Attempt to proceed to checkout

**Expected Result:**
- User is redirected to registration/login page
- Cart items are preserved

**Traceability:**
- AC: [AC-ORDER-001](../l1/acceptance-criteria.md#ac-order-001)
- BR: [BR-ORDER](../l1/business-rules.md#br-order)

---

### TC-AC-ORDER-001-N02 – Reject order when email not verified {#tc-ac-order-001-n02}

**Preconditions:**
- Customer registered but email not verified

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| emailVerified | false | Unverified email status |

**Steps:**
1. Login as unverified customer
2. Add items to cart
3. Attempt checkout

**Expected Result:**
- Verification required message is displayed
- Order is not created

**Traceability:**
- AC: [AC-ORDER-001](../l1/acceptance-criteria.md#ac-order-001)
- BR: [BR-ORDER](../l1/business-rules.md#br-order)

---

### TC-AC-ORDER-001-N03 – Handle payment authorization failure {#tc-ac-order-001-n03}

**Preconditions:**
- Verified customer
- Cart has items
- Payment will fail

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| cardNumber | 4000000000000002 | Test card that declines |

**Steps:**
1. Proceed to checkout
2. Enter declining payment method
3. Submit order

**Expected Result:**
- Payment error message shown
- Cart remains intact
- Retry option available

**Traceability:**
- AC: [AC-ORDER-001](../l1/acceptance-criteria.md#ac-order-001)
- BR: [BR-ORDER](../l1/business-rules.md#br-order)

---

### TC-AC-ORDER-001-N04 – Handle out of stock during checkout {#tc-ac-order-001-n04}

**Preconditions:**
- Item becomes out of stock between cart add and checkout

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| productStock | 0 | Stock depleted by other order |

**Steps:**
1. Add item to cart
2. Another user buys last stock
3. Submit order

**Expected Result:**
- Error shown for out of stock items
- Cart modification allowed

**Traceability:**
- AC: [AC-ORDER-001](../l1/acceptance-criteria.md#ac-order-001)
- BR: [BR-ORDER](../l1/business-rules.md#br-order)

---

### TC-AC-ORDER-001-N05 – Reject invalid shipping address {#tc-ac-order-001-n05}

**Preconditions:**
- Verified customer with items in cart

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| zipCode | invalid | Malformed postal code |
| state |  | Missing required state |

**Steps:**
1. Proceed to checkout
2. Enter invalid shipping address
3. Submit order

**Expected Result:**
- Validation errors displayed for invalid fields
- Order not created

**Traceability:**
- AC: [AC-ORDER-001](../l1/acceptance-criteria.md#ac-order-001)
- BR: [BR-ORDER](../l1/business-rules.md#br-order)

---

### TC-AC-ORDER-002-N01 – Standard shipping fee for subtotal below $50 {#tc-ac-order-002-n01}

**Preconditions:**
- Customer at checkout

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| subtotal | $30.00 | Below threshold |

**Steps:**
1. Add items totaling $30
2. Proceed to checkout

**Expected Result:**
- Standard shipping fee of $5.99 is applied

**Traceability:**
- AC: [AC-ORDER-002](../l1/acceptance-criteria.md#ac-order-002)
- BR: [BR-SHIPPING](../l1/business-rules.md#br-shipping)

---

### TC-AC-ORDER-002-N02 – No free shipping when discount drops subtotal below $50 {#tc-ac-order-002-n02}

**Preconditions:**
- Discount causes subtotal to drop below threshold

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| originalSubtotal | $55.00 | Above threshold |
| discount | $10.00 | Drops below threshold |
| finalSubtotal | $45.00 | Now below $50 |

**Steps:**
1. Add $55 worth of items
2. Apply $10 discount
3. Check shipping

**Expected Result:**
- $5.99 shipping fee applied
- Free shipping not available

**Traceability:**
- AC: [AC-ORDER-002](../l1/acceptance-criteria.md#ac-order-002)
- BR: [BR-SHIPPING](../l1/business-rules.md#br-shipping)

---

### TC-AC-ORDER-003-N01 – Tax calculation fails with incomplete address {#tc-ac-order-003-n01}

**Preconditions:**
- Shipping address missing state

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| state |  | Missing state field |

**Steps:**
1. Enter address without state
2. Attempt to calculate total

**Expected Result:**
- Error requesting state selection
- Tax cannot be calculated

**Traceability:**
- AC: [AC-ORDER-003](../l1/acceptance-criteria.md#ac-order-003)
- BR: [BR-TAX](../l1/business-rules.md#br-tax)

---

### TC-AC-ORDER-003-N02 – Invalid state code rejected {#tc-ac-order-003-n02}

**Preconditions:**
- Customer entering shipping address

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| state | XX | Invalid state code |

**Steps:**
1. Enter address with invalid state code
2. Submit address

**Expected Result:**
- Validation error for invalid state
- Tax not calculated

**Traceability:**
- AC: [AC-ORDER-003](../l1/acceptance-criteria.md#ac-order-003)
- BR: [BR-TAX](../l1/business-rules.md#br-tax)

---

### TC-AC-ORDER-004-N01 – Order creation continues when email service fails {#tc-ac-order-004-n01}

**Preconditions:**
- Email service is down

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| emailServiceStatus | unavailable | Simulated outage |

**Steps:**
1. Disable email service
2. Place order
3. Check order status

**Expected Result:**
- Order is created successfully
- Email queued for retry
- Failure is logged

**Traceability:**
- AC: [AC-ORDER-004](../l1/acceptance-criteria.md#ac-order-004)
- BR: [BR-EMAIL](../l1/business-rules.md#br-email)

---

### TC-AC-ORDER-004-N02 – Failed email logged with error details {#tc-ac-order-004-n02}

**Preconditions:**
- Email service returns error

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| errorType | SMTP_TIMEOUT | Specific error type |

**Steps:**
1. Trigger email failure
2. Check system logs

**Expected Result:**
- Error logged with timestamp
- Error details captured
- Order ID in log

**Traceability:**
- AC: [AC-ORDER-004](../l1/acceptance-criteria.md#ac-order-004)
- BR: [BR-EMAIL](../l1/business-rules.md#br-email)

---

### TC-AC-ORDER-005-N01 – Empty state shown when no orders exist {#tc-ac-order-005-n01}

**Preconditions:**
- New customer with no order history

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| orderCount | 0 | No orders placed |

**Steps:**
1. Login as new customer
2. Navigate to order history

**Expected Result:**
- Empty state message displayed
- Start Shopping link visible and functional

**Traceability:**
- AC: [AC-ORDER-005](../l1/acceptance-criteria.md#ac-order-005)
- BR: [BR-TRACKING](../l1/business-rules.md#br-tracking)

---

### TC-AC-ORDER-005-N02 – Unauthenticated user cannot view order history {#tc-ac-order-005-n02}

**Preconditions:**
- User not logged in

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| sessionType | guest | No authentication |

**Steps:**
1. Navigate to order history without login

**Expected Result:**
- Redirected to login page
- Order history not accessible

**Traceability:**
- AC: [AC-ORDER-005](../l1/acceptance-criteria.md#ac-order-005)
- BR: [BR-TRACKING](../l1/business-rules.md#br-tracking)

---

### TC-AC-ORDER-006-N01 – Reject cancellation for shipped order {#tc-ac-order-006-n01}

**Preconditions:**
- Order exists
- Order status is 'shipped'

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| orderStatus | shipped | Non-cancellable state |

**Steps:**
1. Navigate to shipped order
2. Attempt to cancel

**Expected Result:**
- Error message: cancellation not allowed
- Order status unchanged

**Traceability:**
- AC: [AC-ORDER-006](../l1/acceptance-criteria.md#ac-order-006)
- BR: [AMB-OP-032](../l1/business-rules.md#amb-op-032)

---

### TC-AC-ORDER-006-N02 – Reject cancellation for delivered order {#tc-ac-order-006-n02}

**Preconditions:**
- Order exists
- Order status is 'delivered'

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| orderStatus | delivered | Non-cancellable state |

**Steps:**
1. Navigate to delivered order
2. Attempt to cancel

**Expected Result:**
- Error message: cancellation not allowed

**Traceability:**
- AC: [AC-ORDER-006](../l1/acceptance-criteria.md#ac-order-006)
- BR: [AMB-OP-032](../l1/business-rules.md#amb-op-032)

---

### TC-AC-ORDER-006-N03 – Show message for already cancelled order {#tc-ac-order-006-n03}

**Preconditions:**
- Order exists
- Order status is 'cancelled'

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| orderStatus | cancelled | Already cancelled |

**Steps:**
1. Navigate to cancelled order
2. Check cancel option

**Expected Result:**
- Message: order already cancelled
- No cancel action available

**Traceability:**
- AC: [AC-ORDER-006](../l1/acceptance-criteria.md#ac-order-006)
- BR: [AMB-OP-032](../l1/business-rules.md#amb-op-032)

---

### TC-AC-ORDER-006-N04 – Reject cancellation by non-owner customer {#tc-ac-order-006-n04}

**Preconditions:**
- Order belongs to Customer A
- Customer B is logged in

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| orderOwner | customer-A | Order owner |
| requestingUser | customer-B | Different customer |

**Steps:**
1. Customer B attempts to access order
2. Attempt cancellation

**Expected Result:**
- Access denied or order not found

**Traceability:**
- AC: [AC-ORDER-006](../l1/acceptance-criteria.md#ac-order-006)
- BR: [AMB-OP-028](../l1/business-rules.md#amb-op-028)

---

### TC-AC-ORDER-007-N01 – Price snapshot includes discounted prices {#tc-ac-order-007-n01}

**Preconditions:**
- Product on sale at $19.99 (original $29.99)

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| salePrice | 19.99 | Active sale price |
| originalPrice | 29.99 | Regular price |

**Steps:**
1. Place order during sale
2. Sale ends
3. View order

**Expected Result:**
- Order shows sale price $19.99, not original $29.99

**Traceability:**
- AC: [AC-ORDER-007](../l1/acceptance-criteria.md#ac-order-007)
- BR: [AMB-ENT-019](../l1/business-rules.md#amb-ent-019)

---

### TC-AC-ORDER-007-N02 – Reject manual price edit on existing order {#tc-ac-order-007-n02}

**Preconditions:**
- Order exists with captured prices

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| capturedPrice | 29.99 | Original order price |

**Steps:**
1. Attempt to edit order item price via API
2. Check order

**Expected Result:**
- Price modification rejected or not available

**Traceability:**
- AC: [AC-ORDER-007](../l1/acceptance-criteria.md#ac-order-007)
- BR: [AMB-ENT-019](../l1/business-rules.md#amb-ent-019)

---

### TC-AC-CUST-001-N01 – Reject registration with already used email {#tc-ac-cust-001-n01}

**Preconditions:**
- Email already exists in system

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| email | existing@example.com | Already registered |

**Steps:**
1. Attempt registration with existing email

**Expected Result:**
- Error: 'Email already in use'
- Account not created

**Traceability:**
- AC: [AC-CUST-001](../l1/acceptance-criteria.md#ac-cust-001)
- BR: [AMB-ENT-024](../l1/business-rules.md#amb-ent-024)

---

### TC-AC-CUST-001-N02 – Reject password without number {#tc-ac-cust-001-n02}

**Preconditions:**
- Registration form open

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| password | NoNumbers | Missing number requirement |

**Steps:**
1. Enter password without number
2. Submit form

**Expected Result:**
- Error showing password requirements

**Traceability:**
- AC: [AC-CUST-001](../l1/acceptance-criteria.md#ac-cust-001)
- BR: [AMB-ENT-025](../l1/business-rules.md#amb-ent-025)

---

### TC-AC-CUST-001-N03 – Reject password shorter than 8 characters {#tc-ac-cust-001-n03}

**Preconditions:**
- Registration form open

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| password | Short1 | Only 6 characters |

**Steps:**
1. Enter short password
2. Submit form

**Expected Result:**
- Error: password must be at least 8 characters

**Traceability:**
- AC: [AC-CUST-001](../l1/acceptance-criteria.md#ac-cust-001)
- BR: [AMB-ENT-025](../l1/business-rules.md#amb-ent-025)

---

### TC-AC-CUST-001-N04 – Reject registration with missing required fields {#tc-ac-cust-001-n04}

**Preconditions:**
- Registration form open

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| email | test@example.com | Provided |
| firstName |  | Missing required field |

**Steps:**
1. Submit form with missing first name

**Expected Result:**
- Validation error for missing first name

**Traceability:**
- AC: [AC-CUST-001](../l1/acceptance-criteria.md#ac-cust-001)
- BR: [AMB-ENT-024](../l1/business-rules.md#amb-ent-024)

---

### TC-AC-CUST-001-N05 – Reject invalid email format {#tc-ac-cust-001-n05}

**Preconditions:**
- Registration form open

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| email | notanemail | Invalid format |

**Steps:**
1. Enter invalid email format
2. Submit form

**Expected Result:**
- Validation error for email format

**Traceability:**
- AC: [AC-CUST-001](../l1/acceptance-criteria.md#ac-cust-001)
- BR: [AMB-ENT-024](../l1/business-rules.md#amb-ent-024)

---

### TC-AC-CUST-002-N01 – Reject address with missing required fields {#tc-ac-cust-002-n01}

**Preconditions:**
- Add address form open

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| city |  | Missing required field |

**Steps:**
1. Submit address without city

**Expected Result:**
- Validation error for missing city

**Traceability:**
- AC: [AC-CUST-002](../l1/acceptance-criteria.md#ac-cust-002)
- BR: [AMB-ENT-031](../l1/business-rules.md#amb-ent-031)

---

### TC-AC-CUST-002-N02 – Reject invalid postal code format {#tc-ac-cust-002-n02}

**Preconditions:**
- Add address form open
- Country is USA

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| postalCode | ABCDE | Invalid US postal |

**Steps:**
1. Enter invalid postal code
2. Submit form

**Expected Result:**
- Error: invalid postal code format

**Traceability:**
- AC: [AC-CUST-002](../l1/acceptance-criteria.md#ac-cust-002)
- BR: [AMB-ENT-033](../l1/business-rules.md#amb-ent-033)

---

### TC-AC-CUST-002-N03 – Reject address without recipient name {#tc-ac-cust-002-n03}

**Preconditions:**
- Add address form open

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| recipientName |  | Missing recipient |

**Steps:**
1. Submit address without recipient name

**Expected Result:**
- Validation error for missing recipient name

**Traceability:**
- AC: [AC-CUST-002](../l1/acceptance-criteria.md#ac-cust-002)
- BR: [AMB-ENT-031](../l1/business-rules.md#amb-ent-031)

---

### TC-AC-CUST-003-N01 – Reject invalid card number {#tc-ac-cust-003-n01}

**Preconditions:**
- Add payment form open

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| cardNumber | 1234567890123456 | Invalid card |

**Steps:**
1. Enter invalid card number
2. Attempt to save

**Expected Result:**
- Error from payment gateway displayed

**Traceability:**
- AC: [AC-CUST-003](../l1/acceptance-criteria.md#ac-cust-003)
- BR: [AMB-ENT-035](../l1/business-rules.md#amb-ent-035)

---

### TC-AC-CUST-003-N02 – Reject unsupported card type (Discover) {#tc-ac-cust-003-n02}

**Preconditions:**
- Add payment form open

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| cardType | Discover | Not in supported list |

**Steps:**
1. Attempt to add Discover card

**Expected Result:**
- Error: show supported card types (Visa, MC, Amex)

**Traceability:**
- AC: [AC-CUST-003](../l1/acceptance-criteria.md#ac-cust-003)
- BR: [AMB-ENT-034](../l1/business-rules.md#amb-ent-034)

---

### TC-AC-CUST-003-N03 – Reject expired card {#tc-ac-cust-003-n03}

**Preconditions:**
- Add payment form open

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| expiryDate | 01/2020 | Expired date |

**Steps:**
1. Enter expired card details
2. Attempt to save

**Expected Result:**
- Error: card expired

**Traceability:**
- AC: [AC-CUST-003](../l1/acceptance-criteria.md#ac-cust-003)
- BR: [AMB-ENT-035](../l1/business-rules.md#amb-ent-035)

---

### TC-AC-CUST-004-N01 – Block deletion when customer has pending orders {#tc-ac-cust-004-n01}

**Preconditions:**
- Customer is registered
- Customer has pending order

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| customerId | cust-111 | Customer with pending order |
| orderStatus | pending | Blocking condition |

**Steps:**
1. Login as customer with pending order
2. Request account deletion

**Expected Result:**
- Deletion blocked
- Error message displayed
- Order completion required

**Traceability:**
- AC: [AC-CUST-004](../l1/acceptance-criteria.md#ac-cust-004)
- BR: [BR-CUST](../l1/business-rules.md#br-cust)

---

### TC-AC-CUST-004-N02 – Block deletion when orders are in processing state {#tc-ac-cust-004-n02}

**Preconditions:**
- Customer has order in processing

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| customerId | cust-222 | Customer with processing order |
| orderStatus | processing | Incomplete order |

**Steps:**
1. Login as customer
2. Request deletion with processing order

**Expected Result:**
- Deletion blocked
- Shows list of incomplete orders

**Traceability:**
- AC: [AC-CUST-004](../l1/acceptance-criteria.md#ac-cust-004)
- BR: [BR-CUST](../l1/business-rules.md#br-cust)

---

### TC-AC-CUST-004-N03 – Reject deletion without confirmation {#tc-ac-cust-004-n03}

**Preconditions:**
- Customer is registered

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| customerId | cust-333 | Valid customer |

**Steps:**
1. Request account deletion
2. Cancel at confirmation dialog

**Expected Result:**
- Account remains active
- No data modified

**Traceability:**
- AC: [AC-CUST-004](../l1/acceptance-criteria.md#ac-cust-004)
- BR: [BR-CUST](../l1/business-rules.md#br-cust)

---

### TC-AC-ADMIN-001-N01 – Reject product with name too short {#tc-ac-admin-001-n01}

**Preconditions:**
- Admin on add product form

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| name | A | 1 char - below minimum 2 |

**Steps:**
1. Enter single character name
2. Submit form

**Expected Result:**
- Validation error shown
- Product not created

**Traceability:**
- AC: [AC-ADMIN-001](../l1/acceptance-criteria.md#ac-admin-001)
- BR: [BR-ADMIN](../l1/business-rules.md#br-admin), [BR-PROD](../l1/business-rules.md#br-prod)

---

### TC-AC-ADMIN-001-N02 – Reject product with name too long {#tc-ac-admin-001-n02}

**Preconditions:**
- Admin on add product form

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| name | 201 character string | Exceeds 200 char limit |

**Steps:**
1. Enter 201 character name
2. Submit form

**Expected Result:**
- Length validation error shown
- Product not created

**Traceability:**
- AC: [AC-ADMIN-001](../l1/acceptance-criteria.md#ac-admin-001)
- BR: [BR-ADMIN](../l1/business-rules.md#br-admin), [BR-PROD](../l1/business-rules.md#br-prod)

---

### TC-AC-ADMIN-001-N03 – Reject product with price below minimum {#tc-ac-admin-001-n03}

**Preconditions:**
- Admin on add product form

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| price | 0.00 | Below $0.01 minimum |

**Steps:**
1. Enter valid name
2. Set price to $0.00
3. Submit

**Expected Result:**
- Error: Price must be at least $0.01
- Product not created

**Traceability:**
- AC: [AC-ADMIN-001](../l1/acceptance-criteria.md#ac-admin-001)
- BR: [BR-ADMIN](../l1/business-rules.md#br-admin), [BR-PROD](../l1/business-rules.md#br-prod)

---

### TC-AC-ADMIN-001-N04 – Reject product without category {#tc-ac-admin-001-n04}

**Preconditions:**
- Admin on add product form

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| category | null | No category selected |

**Steps:**
1. Enter valid name and price
2. Leave category empty
3. Submit

**Expected Result:**
- Error: Category is required
- Product not created

**Traceability:**
- AC: [AC-ADMIN-001](../l1/acceptance-criteria.md#ac-admin-001)
- BR: [BR-ADMIN](../l1/business-rules.md#br-admin), [BR-PROD](../l1/business-rules.md#br-prod)

---

### TC-AC-ADMIN-002-N01 – Concurrent edit shows conflict warning {#tc-ac-admin-002-n01}

**Preconditions:**
- Two admins editing same product

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| admin1 | admin-A | First editor |
| admin2 | admin-B | Second editor |

**Steps:**
1. Admin A opens product
2. Admin B opens same product
3. Admin A saves
4. Admin B saves

**Expected Result:**
- Admin B sees conflict warning
- Last write persisted

**Traceability:**
- AC: [AC-ADMIN-002](../l1/acceptance-criteria.md#ac-admin-002)
- BR: [BR-ADMIN](../l1/business-rules.md#br-admin), [BR-PROD](../l1/business-rules.md#br-prod)

---

### TC-AC-ADMIN-002-N02 – Reject edit with invalid name length {#tc-ac-admin-002-n02}

**Preconditions:**
- Product exists

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| name | X | Too short - 1 char |

**Steps:**
1. Edit product
2. Change name to single character
3. Save

**Expected Result:**
- Validation error shown
- Changes not saved

**Traceability:**
- AC: [AC-ADMIN-002](../l1/acceptance-criteria.md#ac-admin-002)
- BR: [BR-ADMIN](../l1/business-rules.md#br-admin), [BR-PROD](../l1/business-rules.md#br-prod)

---

### TC-AC-ADMIN-002-N03 – Reject edit with price below minimum {#tc-ac-admin-002-n03}

**Preconditions:**
- Product exists

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| price | 0.00 | Below $0.01 minimum |

**Steps:**
1. Edit product
2. Set price to $0.00
3. Save

**Expected Result:**
- Price validation error
- Original price retained

**Traceability:**
- AC: [AC-ADMIN-002](../l1/acceptance-criteria.md#ac-admin-002)
- BR: [BR-ADMIN](../l1/business-rules.md#br-admin), [BR-PROD](../l1/business-rules.md#br-prod)

---

### TC-AC-ADMIN-003-N01 – Cancel deletion takes no action {#tc-ac-admin-003-n01}

**Preconditions:**
- Product exists

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| productId | prod-202 | Product to not delete |

**Steps:**
1. Click delete on product
2. Cancel confirmation dialog

**Expected Result:**
- Product still exists
- Product visible in listings

**Traceability:**
- AC: [AC-ADMIN-003](../l1/acceptance-criteria.md#ac-admin-003)
- BR: [BR-ADMIN](../l1/business-rules.md#br-admin), [BR-PROD](../l1/business-rules.md#br-prod)

---

### TC-AC-ADMIN-003-N02 – Delete requires confirmation dialog {#tc-ac-admin-003-n02}

**Preconditions:**
- Product exists

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| productId | prod-303 | Product to delete |

**Steps:**
1. Click delete button
2. Observe confirmation requirement

**Expected Result:**
- Confirmation dialog appears
- No deletion without confirm

**Traceability:**
- AC: [AC-ADMIN-003](../l1/acceptance-criteria.md#ac-admin-003)
- BR: [BR-ADMIN](../l1/business-rules.md#br-admin), [BR-PROD](../l1/business-rules.md#br-prod)

---

### TC-AC-ADMIN-004-N01 – Show empty state when no orders match filter {#tc-ac-admin-004-n01}

**Preconditions:**
- No orders match criteria

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| statusFilter | cancelled | No cancelled orders |

**Steps:**
1. Apply filter with no matching orders

**Expected Result:**
- Empty state message displayed
- Clear indication no results

**Traceability:**
- AC: [AC-ADMIN-004](../l1/acceptance-criteria.md#ac-admin-004)
- BR: [BR-ADMIN](../l1/business-rules.md#br-admin), [BR-ORDER](../l1/business-rules.md#br-order)

---

### TC-AC-ADMIN-004-N02 – Show empty state for non-existent order search {#tc-ac-admin-004-n02}

**Preconditions:**
- Order ORD-99999 does not exist

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| searchTerm | ORD-99999 | Non-existent order |

**Steps:**
1. Search for non-existent order number

**Expected Result:**
- Empty state shown
- No false matches

**Traceability:**
- AC: [AC-ADMIN-004](../l1/acceptance-criteria.md#ac-admin-004)
- BR: [BR-ADMIN](../l1/business-rules.md#br-admin), [BR-ORDER](../l1/business-rules.md#br-order)

---

### TC-AC-ADMIN-005-N01 – Reject invalid transition from pending to shipped {#tc-ac-admin-005-n01}

**Preconditions:**
- Admin is authenticated
- Order exists with pending status

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| currentStatus | pending | Starting status |
| newStatus | shipped | Invalid skip transition |

**Steps:**
1. Navigate to order
2. Attempt to select shipped status
3. Submit change

**Expected Result:**
- Error displayed
- Message shows allowed transitions: pending→confirmed

**Traceability:**
- AC: [AC-ADMIN-005](../l1/acceptance-criteria.md#ac-admin-005)
- BR: [AMB-OP-047](../l1/business-rules.md#amb-op-047)

---

### TC-AC-ADMIN-005-N02 – Reject backward transition from delivered to shipped {#tc-ac-admin-005-n02}

**Preconditions:**
- Admin is authenticated
- Order is in delivered status

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| currentStatus | delivered | Final status |
| newStatus | shipped | Backward transition |

**Steps:**
1. Navigate to delivered order
2. Attempt to change to shipped
3. Submit

**Expected Result:**
- Error message displayed
- Status remains delivered

**Traceability:**
- AC: [AC-ADMIN-005](../l1/acceptance-criteria.md#ac-admin-005)
- BR: [AMB-OP-047](../l1/business-rules.md#amb-op-047)

---

### TC-AC-ADMIN-005-N03 – Reject transition from pending directly to delivered {#tc-ac-admin-005-n03}

**Preconditions:**
- Admin is authenticated
- Order is pending

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| currentStatus | pending | Initial status |
| newStatus | delivered | Skips two steps |

**Steps:**
1. Navigate to order
2. Select delivered status
3. Submit

**Expected Result:**
- Error with allowed transitions shown
- Order remains pending

**Traceability:**
- AC: [AC-ADMIN-005](../l1/acceptance-criteria.md#ac-admin-005)
- BR: [AMB-OP-047](../l1/business-rules.md#amb-op-047)

---

### TC-AC-ADMIN-006-N01 – Block adjustment resulting in negative inventory {#tc-ac-admin-006-n01}

**Preconditions:**
- Variant has stock of 10

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| currentStock | 10 | Current level |
| delta | -15 | Would result in -5 |

**Steps:**
1. Select variant
2. Apply delta -15
3. Submit

**Expected Result:**
- Error: Inventory cannot go below zero
- Stock remains at 10

**Traceability:**
- AC: [AC-ADMIN-006](../l1/acceptance-criteria.md#ac-admin-006)
- BR: [AMB-OP-054](../l1/business-rules.md#amb-op-054)

---

### TC-AC-ADMIN-006-N02 – Show row-by-row errors for invalid CSV import {#tc-ac-admin-006-n02}

**Preconditions:**
- Admin on inventory import page

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| csvRow1 | VAR-001,abc | Invalid non-numeric |
| csvRow2 | VAR-002,50 | Valid row |
| csvRow3 | INVALID,-10 | Invalid variant and negative |

**Steps:**
1. Upload CSV with mixed valid/invalid rows
2. Submit import

**Expected Result:**
- Row 1 error: invalid quantity
- Row 2: success
- Row 3 error: invalid variant and negative value

**Traceability:**
- AC: [AC-ADMIN-006](../l1/acceptance-criteria.md#ac-admin-006)
- BR: [AMB-OP-055](../l1/business-rules.md#amb-op-055)

---

### TC-AC-ADMIN-006-N03 – Reject setting absolute value to negative number {#tc-ac-admin-006-n03}

**Preconditions:**
- Admin managing inventory

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| adjustmentType | absolute | Direct set |
| newValue | -5 | Negative absolute value |

**Steps:**
1. Select variant
2. Choose absolute adjustment
3. Enter -5
4. Submit

**Expected Result:**
- Error: Inventory cannot go below zero
- Change rejected

**Traceability:**
- AC: [AC-ADMIN-006](../l1/acceptance-criteria.md#ac-admin-006)
- BR: [AMB-OP-054](../l1/business-rules.md#amb-op-054)

---

### TC-AC-VAR-001-N01 – Show out of stock for unavailable variant selection {#tc-ac-var-001-n01}

**Preconditions:**
- Variant M-Red has 0 stock

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| selectedVariant | M-Red | Zero stock variant |

**Steps:**
1. View product
2. Select M size and Red color

**Expected Result:**
- Out of stock message displayed for M-Red variant

**Traceability:**
- AC: [AC-VAR-001](../l1/acceptance-criteria.md#ac-var-001)
- BR: [AMB-ENT-009](../l1/business-rules.md#amb-ent-009)

---

### TC-AC-VAR-001-N02 – Allow selecting other variants when one is out of stock {#tc-ac-var-001-n02}

**Preconditions:**
- M-Red out of stock
- L-Red in stock

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| outOfStock | M-Red | Unavailable |
| inStock | L-Red | Available alternative |

**Steps:**
1. View product
2. See M-Red is out of stock
3. Select L-Red instead

**Expected Result:**
- L-Red variant selectable
- Add to cart enabled for L-Red

**Traceability:**
- AC: [AC-VAR-001](../l1/acceptance-criteria.md#ac-var-001)
- BR: [AMB-ENT-011](../l1/business-rules.md#amb-ent-011)

---

### TC-AC-CAT-001-N01 – Hide inactive category from customer view {#tc-ac-cat-001-n01}

**Preconditions:**
- Discontinued category marked inactive with 3 products

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| category | Discontinued | Inactive category |
| status | inactive | Disabled flag |

**Steps:**
1. Browse categories as customer
2. Search for Discontinued category

**Expected Result:**
- Discontinued category not visible
- Its products not shown

**Traceability:**
- AC: [AC-CAT-001](../l1/acceptance-criteria.md#ac-cat-001)
- BR: [AMB-ENT-037](../l1/business-rules.md#amb-ent-037)

---

### TC-AC-CAT-001-N02 – Hide products of inactive category from listings {#tc-ac-cat-001-n02}

**Preconditions:**
- Product only in inactive category

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| productId | PROD-OLD | Only in inactive category |
| category | Archived | Inactive |

**Steps:**
1. Search for product
2. Browse all categories

**Expected Result:**
- Product not found in any listing
- Not in search results

**Traceability:**
- AC: [AC-CAT-001](../l1/acceptance-criteria.md#ac-cat-001)
- BR: [AMB-ENT-037](../l1/business-rules.md#amb-ent-037)

---

## Boundary Tests

### TC-AC-PROD-001-B01 – Filter with price range at zero minimum {#tc-ac-prod-001-b01}

**Preconditions:**
- Products exist including free items (price 0)

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| priceMin | 0 | Boundary: zero price |
| priceMax | 10 | Small range |

**Steps:**
1. Set price range 0-10
2. Apply filters

**Expected Result:**
- Products with price 0-10 displayed
- Free items included

**Traceability:**
- AC: [AC-PROD-001](../l1/acceptance-criteria.md#ac-prod-001)
- BR: [BR-PROD](../l1/business-rules.md#br-prod)

---

### TC-AC-PROD-001-B02 – Display exactly 20th and 21st product across pages {#tc-ac-prod-001-b02}

**Preconditions:**
- At least 21 products exist

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| totalProducts | 21 | Just over page boundary |

**Steps:**
1. View first page
2. Note 20th product
3. Go to page 2
4. Verify 21st product

**Expected Result:**
- Page 1 has exactly 20 products
- Page 2 has 1 product

**Traceability:**
- AC: [AC-PROD-001](../l1/acceptance-criteria.md#ac-prod-001)
- BR: [BR-PROD](../l1/business-rules.md#br-prod)

---

### TC-AC-PROD-002-B01 – Product with exactly zero inventory shows out of stock {#tc-ac-prod-002-b01}

**Preconditions:**
- Product inventory is exactly 0

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| inventory | 0 | Boundary: exactly zero |

**Steps:**
1. Set product inventory to exactly 0
2. View product

**Expected Result:**
- Out of Stock indicator shown
- Add to cart disabled

**Traceability:**
- AC: [AC-PROD-002](../l1/acceptance-criteria.md#ac-prod-002)
- BR: [BR-PROD](../l1/business-rules.md#br-prod)

---

### TC-AC-PROD-002-B02 – Product with inventory 1 does NOT show out of stock {#tc-ac-prod-002-b02}

**Preconditions:**
- Product inventory is exactly 1

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| inventory | 1 | Boundary: just above zero |

**Steps:**
1. Set product inventory to 1
2. View product

**Expected Result:**
- No Out of Stock indicator
- Add to cart enabled

**Traceability:**
- AC: [AC-PROD-002](../l1/acceptance-criteria.md#ac-prod-002)
- BR: [BR-PROD](../l1/business-rules.md#br-prod)

---

### TC-AC-CART-001-B01 – Add quantity equal to exact available stock {#tc-ac-cart-001-b01}

**Preconditions:**
- Product has known stock level

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| inventory | 5 | Exact stock |
| quantity | 5 | Equals stock |

**Steps:**
1. Set quantity to exact available stock (5)
2. Add to cart

**Expected Result:**
- All 5 items added successfully
- No error shown

**Traceability:**
- AC: [AC-CART-001](../l1/acceptance-criteria.md#ac-cart-001)
- BR: [BR-CART](../l1/business-rules.md#br-cart)

---

### TC-AC-CART-001-B02 – Add minimum quantity of 1 {#tc-ac-cart-001-b02}

**Preconditions:**
- Product in stock

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| quantity | 1 | Minimum allowed |

**Steps:**
1. Set quantity to 1
2. Add to cart

**Expected Result:**
- 1 item added to cart

**Traceability:**
- AC: [AC-CART-001](../l1/acceptance-criteria.md#ac-cart-001)
- BR: [BR-CART](../l1/business-rules.md#br-cart)

---

### TC-AC-CART-002-B01 – Increment to exactly max stock {#tc-ac-cart-002-b01}

**Preconditions:**
- Product A qty 9 in cart
- Stock is 10

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| existingQty | 9 | One below max |
| addQty | 1 | Reaches exactly max |
| stock | 10 | Maximum |

**Steps:**
1. Add 1 more to reach stock limit

**Expected Result:**
- Quantity becomes 10
- No error shown

**Traceability:**
- AC: [AC-CART-002](../l1/acceptance-criteria.md#ac-cart-002)
- BR: [BR-CART](../l1/business-rules.md#br-cart)

---

### TC-AC-CART-002-B02 – Adding 0 quantity should not change cart {#tc-ac-cart-002-b02}

**Preconditions:**
- Product A in cart with quantity 2

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| existingQty | 2 | Current |
| addQty | 0 | Zero addition |

**Steps:**
1. Attempt to add 0 quantity

**Expected Result:**
- Quantity remains 2
- No change or error handled gracefully

**Traceability:**
- AC: [AC-CART-002](../l1/acceptance-criteria.md#ac-cart-002)
- BR: [BR-CART](../l1/business-rules.md#br-cart)

---

### TC-AC-CART-003-B01 – Cart persists at exactly 7 days minus 1 hour {#tc-ac-cart-003-b01}

**Preconditions:**
- Guest cart created

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| timeSinceActivity | 167 hours | Just under 7 days |

**Steps:**
1. Create guest cart
2. Return at 6 days 23 hours

**Expected Result:**
- Cart still exists
- Items preserved

**Traceability:**
- AC: [AC-CART-003](../l1/acceptance-criteria.md#ac-cart-003)
- BR: [BR-CART](../l1/business-rules.md#br-cart)

---

### TC-AC-CART-003-B02 – Cart cleared at exactly 7 days of inactivity {#tc-ac-cart-003-b02}

**Preconditions:**
- Guest cart created exactly 7 days ago

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| timeSinceActivity | 168 hours | Exactly 7 days |

**Steps:**
1. Create guest cart
2. Return at exactly 7 days

**Expected Result:**
- Cart is cleared per policy

**Traceability:**
- AC: [AC-CART-003](../l1/acceptance-criteria.md#ac-cart-003)
- BR: [BR-CART](../l1/business-rules.md#br-cart)

---

### TC-AC-CART-004-B01 – Merge when combined quantity exactly equals stock {#tc-ac-cart-004-b01}

**Preconditions:**
- Guest has 5, Account has 5, Stock is 10

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| guestQty | 5 | Half of stock |
| accountQty | 5 | Half of stock |
| stock | 10 | Exact boundary |

**Steps:**
1. Setup carts totaling exactly stock
2. Login
3. Verify merge

**Expected Result:**
- Cart contains quantity 10
- No stock warning shown

**Traceability:**
- AC: [AC-CART-004](../l1/acceptance-criteria.md#ac-cart-004)
- BR: [AMB-ENT-013](../l1/business-rules.md#amb-ent-013)

---

### TC-AC-CART-004-B02 – Merge when combined quantity is stock plus one {#tc-ac-cart-004-b02}

**Preconditions:**
- Guest has 6, Account has 5, Stock is 10

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| guestQty | 6 | Over half |
| accountQty | 5 | Half |
| stock | 10 | Max boundary |

**Steps:**
1. Setup carts totaling stock+1
2. Login
3. Verify limit applied

**Expected Result:**
- Cart limited to 10
- Notification shown

**Traceability:**
- AC: [AC-CART-004](../l1/acceptance-criteria.md#ac-cart-004)
- BR: [AMB-ENT-013](../l1/business-rules.md#amb-ent-013)

---

### TC-AC-CART-005-B01 – Set quantity to 1 - minimum valid {#tc-ac-cart-005-b01}

**Preconditions:**
- Cart has SKU-A qty 5

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| newQuantity | 1 | Minimum valid |

**Steps:**
1. Change quantity to 1
2. Verify update

**Expected Result:**
- Quantity set to 1
- Total recalculated

**Traceability:**
- AC: [AC-CART-005](../l1/acceptance-criteria.md#ac-cart-005)
- BR: [AMB-OP-011](../l1/business-rules.md#amb-op-011)

---

### TC-AC-CART-005-B02 – Set quantity to exactly available stock {#tc-ac-cart-005-b02}

**Preconditions:**
- Stock is 10

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| newQuantity | 10 | Exact stock |
| stock | 10 | Available |

**Steps:**
1. Set quantity to exactly stock level
2. Verify success

**Expected Result:**
- Quantity set to 10
- No warning shown

**Traceability:**
- AC: [AC-CART-005](../l1/acceptance-criteria.md#ac-cart-005)
- BR: [AMB-OP-012](../l1/business-rules.md#amb-op-012)

---

### TC-AC-CART-006-B01 – Remove one of two items - cart not empty {#tc-ac-cart-006-b01}

**Preconditions:**
- Cart has exactly 2 items

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| cartItems | [SKU-A, SKU-B] | Two items |

**Steps:**
1. Remove one item
2. Verify cart state

**Expected Result:**
- One item remains
- Empty cart state NOT shown

**Traceability:**
- AC: [AC-CART-006](../l1/acceptance-criteria.md#ac-cart-006)
- BR: [AMB-OP-014](../l1/business-rules.md#amb-op-014)

---

### TC-AC-CART-006-B02 – Undo at last moment before option disappears {#tc-ac-cart-006-b02}

**Preconditions:**
- Item removed
- Undo option about to expire

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| undoTimeout | near-expiry | Edge of window |

**Steps:**
1. Remove item
2. Wait until undo almost expires
3. Click undo

**Expected Result:**
- Item successfully restored

**Traceability:**
- AC: [AC-CART-006](../l1/acceptance-criteria.md#ac-cart-006)
- BR: [AMB-OP-015](../l1/business-rules.md#amb-op-015)

---

### TC-AC-CART-007-B01 – Stock drops to exactly zero - warning appears {#tc-ac-cart-007-b01}

**Preconditions:**
- Stock at 1
- Another user purchases last item

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| stockBefore | 1 | Last item |
| stockAfter | 0 | Now zero |

**Steps:**
1. Have item in cart
2. Stock drops to 0
3. Verify warning

**Expected Result:**
- Warning appears when stock hits zero

**Traceability:**
- AC: [AC-CART-007](../l1/acceptance-criteria.md#ac-cart-007)
- BR: [AMB-ENT-014](../l1/business-rules.md#amb-ent-014)

---

### TC-AC-CART-007-B02 – Stock at 1 - no warning shown yet {#tc-ac-cart-007-b02}

**Preconditions:**
- Stock at 1
- Cart qty is 1

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| stock | 1 | Not yet zero |

**Steps:**
1. View cart when stock is 1

**Expected Result:**
- No out of stock warning
- Checkout not blocked

**Traceability:**
- AC: [AC-CART-007](../l1/acceptance-criteria.md#ac-cart-007)
- BR: [AMB-ENT-014](../l1/business-rules.md#amb-ent-014)

---

### TC-AC-CART-008-B01 – Minimal price change triggers notification {#tc-ac-cart-008-b01}

**Preconditions:**
- Price changes by smallest unit

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| oldPrice | 10.00 | Original |
| newPrice | 10.01 | Minimal change |

**Steps:**
1. Price changes by $0.01
2. View cart

**Expected Result:**
- Notification shown for minimal change

**Traceability:**
- AC: [AC-CART-008](../l1/acceptance-criteria.md#ac-cart-008)
- BR: [AMB-ENT-016](../l1/business-rules.md#amb-ent-016)

---

### TC-AC-CART-008-B02 – Price changed to zero - free item {#tc-ac-cart-008-b02}

**Preconditions:**
- Price set to $0

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| oldPrice | 10.00 | Had cost |
| newPrice | 0.00 | Now free |

**Steps:**
1. Price drops to zero
2. View cart

**Expected Result:**
- Cart shows $0 price
- Notification shown

**Traceability:**
- AC: [AC-CART-008](../l1/acceptance-criteria.md#ac-cart-008)
- BR: [AMB-ENT-016](../l1/business-rules.md#amb-ent-016)

---

### TC-AC-ORDER-001-B01 – Place order with minimum cart (1 item) {#tc-ac-order-001-b01}

**Preconditions:**
- Verified customer
- Valid address and payment

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| cartItemCount | 1 | Minimum possible cart size |

**Steps:**
1. Add exactly 1 item to cart
2. Complete checkout

**Expected Result:**
- Order created successfully
- Single item order confirmed

**Traceability:**
- AC: [AC-ORDER-001](../l1/acceptance-criteria.md#ac-order-001)
- BR: [BR-ORDER](../l1/business-rules.md#br-order)

---

### TC-AC-ORDER-001-B02 – Place order with empty cart rejected {#tc-ac-order-001-b02}

**Preconditions:**
- Verified customer
- Empty cart

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| cartItemCount | 0 | Empty cart boundary |

**Steps:**
1. Navigate to checkout with empty cart

**Expected Result:**
- Checkout blocked
- Message to add items first

**Traceability:**
- AC: [AC-ORDER-001](../l1/acceptance-criteria.md#ac-order-001)
- BR: [BR-ORDER](../l1/business-rules.md#br-order)

---

### TC-AC-ORDER-002-B01 – Free shipping at exact $50.00 boundary {#tc-ac-order-002-b01}

**Preconditions:**
- Customer at checkout

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| subtotal | $50.00 | Exact boundary value |

**Steps:**
1. Add items totaling exactly $50.00
2. Check shipping at checkout

**Expected Result:**
- Free shipping is applied

**Traceability:**
- AC: [AC-ORDER-002](../l1/acceptance-criteria.md#ac-order-002)
- BR: [BR-SHIPPING](../l1/business-rules.md#br-shipping)

---

### TC-AC-ORDER-002-B02 – Standard shipping at $49.99 (just below threshold) {#tc-ac-order-002-b02}

**Preconditions:**
- Customer at checkout

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| subtotal | $49.99 | One cent below threshold |

**Steps:**
1. Add items totaling $49.99
2. Check shipping at checkout

**Expected Result:**
- Standard shipping fee of $5.99 applied

**Traceability:**
- AC: [AC-ORDER-002](../l1/acceptance-criteria.md#ac-order-002)
- BR: [BR-SHIPPING](../l1/business-rules.md#br-shipping)

---

### TC-AC-ORDER-002-B03 – Free shipping at $50.01 (just above threshold) {#tc-ac-order-002-b03}

**Preconditions:**
- Customer at checkout

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| subtotal | $50.01 | One cent above threshold |

**Steps:**
1. Add items totaling $50.01
2. Check shipping

**Expected Result:**
- Free shipping is applied

**Traceability:**
- AC: [AC-ORDER-002](../l1/acceptance-criteria.md#ac-order-002)
- BR: [BR-SHIPPING](../l1/business-rules.md#br-shipping)

---

### TC-AC-ORDER-003-B01 – Tax on minimum order amount {#tc-ac-order-003-b01}

**Preconditions:**
- Minimum value order

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| subtotal | $0.01 | Smallest possible amount |
| state | CA | State with tax |

**Steps:**
1. Create order with $0.01 subtotal
2. Check tax calculation

**Expected Result:**
- Tax calculated correctly on small amount
- Rounds appropriately

**Traceability:**
- AC: [AC-ORDER-003](../l1/acceptance-criteria.md#ac-order-003)
- BR: [BR-TAX](../l1/business-rules.md#br-tax)

---

### TC-AC-ORDER-003-B02 – Tax in zero-tax state {#tc-ac-order-003-b02}

**Preconditions:**
- Shipping to state with no sales tax

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| state | OR | Oregon has no sales tax |
| subtotal | $100.00 | Test amount |

**Steps:**
1. Enter Oregon shipping address
2. View tax calculation

**Expected Result:**
- Tax shows as $0.00
- Order total equals subtotal plus shipping

**Traceability:**
- AC: [AC-ORDER-003](../l1/acceptance-criteria.md#ac-order-003)
- BR: [BR-TAX](../l1/business-rules.md#br-tax)

---

### TC-AC-ORDER-004-B01 – Email content with single item order {#tc-ac-order-004-b01}

**Preconditions:**
- Order with single item

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| itemCount | 1 | Minimum items |

**Steps:**
1. Place single item order
2. Check email content

**Expected Result:**
- Email correctly shows single item
- No plural grammar issues

**Traceability:**
- AC: [AC-ORDER-004](../l1/acceptance-criteria.md#ac-order-004)
- BR: [BR-EMAIL](../l1/business-rules.md#br-email)

---

### TC-AC-ORDER-004-B02 – Email content with many items {#tc-ac-order-004-b02}

**Preconditions:**
- Order with many different items

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| itemCount | 50 | Large order |

**Steps:**
1. Place order with 50 different items
2. Verify email renders correctly

**Expected Result:**
- All 50 items listed in email
- Email renders properly

**Traceability:**
- AC: [AC-ORDER-004](../l1/acceptance-criteria.md#ac-order-004)
- BR: [BR-EMAIL](../l1/business-rules.md#br-email)

---

### TC-AC-ORDER-005-B01 – Order history with single order {#tc-ac-order-005-b01}

**Preconditions:**
- Customer has exactly one order

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| orderCount | 1 | Minimum orders |

**Steps:**
1. Login as customer with 1 order
2. View order history

**Expected Result:**
- Single order displayed correctly
- No pagination needed
- No empty state shown

**Traceability:**
- AC: [AC-ORDER-005](../l1/acceptance-criteria.md#ac-order-005)
- BR: [BR-TRACKING](../l1/business-rules.md#br-tracking)

---

### TC-AC-ORDER-005-B02 – Pagination boundary at page limit {#tc-ac-order-005-b02}

**Preconditions:**
- Customer has exactly page limit number of orders

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| orderCount | 20 | Assuming 20 per page limit |

**Steps:**
1. View order history with exactly page limit orders

**Expected Result:**
- All orders on single page
- No next page link or shows disabled

**Traceability:**
- AC: [AC-ORDER-005](../l1/acceptance-criteria.md#ac-order-005)
- BR: [BR-TRACKING](../l1/business-rules.md#br-tracking)

---

### TC-AC-ORDER-006-B01 – Cancel at exact moment status changes to shipped {#tc-ac-order-006-b01}

**Preconditions:**
- Order is confirmed
- Shipping process about to start

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| orderStatus | confirmed | About to ship |
| timing | concurrent | Race condition test |

**Steps:**
1. Initiate cancel and ship concurrently
2. Check final state

**Expected Result:**
- Only one operation succeeds
- Consistent final state

**Traceability:**
- AC: [AC-ORDER-006](../l1/acceptance-criteria.md#ac-order-006)
- BR: [AMB-ENT-044](../l1/business-rules.md#amb-ent-044)

---

### TC-AC-ORDER-006-B02 – Cancel order with zero total amount {#tc-ac-order-006-b02}

**Preconditions:**
- Order total is 0 (100% discount)
- Order is pending

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| orderTotal | 0.00 | Zero total order |

**Steps:**
1. Cancel zero-total order

**Expected Result:**
- Order cancelled
- No refund initiated (nothing to refund)

**Traceability:**
- AC: [AC-ORDER-006](../l1/acceptance-criteria.md#ac-order-006)
- BR: [AMB-OP-030](../l1/business-rules.md#amb-op-030)

---

### TC-AC-ORDER-007-B01 – Price snapshot with maximum decimal precision {#tc-ac-order-007-b01}

**Preconditions:**
- Product price has extended decimals

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| productPrice | 29.999 | 3 decimal price |

**Steps:**
1. Place order with extended decimal price
2. View order

**Expected Result:**
- Price correctly rounded/captured

**Traceability:**
- AC: [AC-ORDER-007](../l1/acceptance-criteria.md#ac-order-007)
- BR: [AMB-ENT-019](../l1/business-rules.md#amb-ent-019)

---

### TC-AC-ORDER-007-B02 – Price snapshot at zero price (free item) {#tc-ac-order-007-b02}

**Preconditions:**
- Free promotional item in cart

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| productPrice | 0.00 | Free item |

**Steps:**
1. Place order with free item
2. View order details

**Expected Result:**
- Order shows $0.00 for free item

**Traceability:**
- AC: [AC-ORDER-007](../l1/acceptance-criteria.md#ac-order-007)
- BR: [AMB-ENT-019](../l1/business-rules.md#amb-ent-019)

---

### TC-AC-CUST-001-B01 – Register with exactly 8 character password with number {#tc-ac-cust-001-b01}

**Preconditions:**
- Registration form open

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| password | Abcdefg1 | Exactly 8 chars |

**Steps:**
1. Enter 8-char password with number
2. Submit form

**Expected Result:**
- Password accepted
- Registration succeeds

**Traceability:**
- AC: [AC-CUST-001](../l1/acceptance-criteria.md#ac-cust-001)
- BR: [AMB-ENT-025](../l1/business-rules.md#amb-ent-025)

---

### TC-AC-CUST-001-B02 – Reject 7 character password even with number {#tc-ac-cust-001-b02}

**Preconditions:**
- Registration form open

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| password | Abcdef1 | Only 7 characters |

**Steps:**
1. Enter 7-char password
2. Submit form

**Expected Result:**
- Error: password must be at least 8 characters

**Traceability:**
- AC: [AC-CUST-001](../l1/acceptance-criteria.md#ac-cust-001)
- BR: [AMB-ENT-025](../l1/business-rules.md#amb-ent-025)

---

### TC-AC-CUST-002-B01 – Add address with minimum length values {#tc-ac-cust-002-b01}

**Preconditions:**
- Add address form open

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| street | 1 A | Minimum street |
| city | LA | 2 char city |

**Steps:**
1. Enter minimum length values
2. Submit form

**Expected Result:**
- Address saved successfully

**Traceability:**
- AC: [AC-CUST-002](../l1/acceptance-criteria.md#ac-cust-002)
- BR: [AMB-ENT-027](../l1/business-rules.md#amb-ent-027)

---

### TC-AC-CUST-002-B02 – Validate postal code for different countries {#tc-ac-cust-002-b02}

**Preconditions:**
- Add address form open

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| country | Canada | Different format |
| postalCode | K1A 0B1 | Canadian postal |

**Steps:**
1. Select Canada
2. Enter Canadian postal code
3. Submit

**Expected Result:**
- Canadian postal format accepted

**Traceability:**
- AC: [AC-CUST-002](../l1/acceptance-criteria.md#ac-cust-002)
- BR: [AMB-ENT-033](../l1/business-rules.md#amb-ent-033)

---

### TC-AC-CUST-003-B01 – Add card expiring this month {#tc-ac-cust-003-b01}

**Preconditions:**
- Add payment form open

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| expiryDate | current month/year | Edge of expiry |

**Steps:**
1. Enter card expiring this month
2. Save

**Expected Result:**
- Card accepted (still valid)

**Traceability:**
- AC: [AC-CUST-003](../l1/acceptance-criteria.md#ac-cust-003)
- BR: [AMB-ENT-035](../l1/business-rules.md#amb-ent-035)

---

### TC-AC-CUST-003-B02 – Reject card expiring last month {#tc-ac-cust-003-b02}

**Preconditions:**
- Add payment form open

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| expiryDate | last month | Just expired |

**Steps:**
1. Enter card expired last month
2. Save

**Expected Result:**
- Error: card expired

**Traceability:**
- AC: [AC-CUST-003](../l1/acceptance-criteria.md#ac-cust-003)
- BR: [AMB-ENT-035](../l1/business-rules.md#amb-ent-035)

---

### TC-AC-CUST-004-B01 – Delete account immediately after last order completes {#tc-ac-cust-004-b01}

**Preconditions:**
- Customer has one order just completed

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| customerId | cust-444 | Customer with just-completed order |
| orderStatus | completed | Just transitioned |

**Steps:**
1. Complete last pending order
2. Immediately request deletion

**Expected Result:**
- Deletion allowed
- No blocking error shown

**Traceability:**
- AC: [AC-CUST-004](../l1/acceptance-criteria.md#ac-cust-004)
- BR: [BR-CUST](../l1/business-rules.md#br-cust)

---

### TC-AC-ADMIN-001-B01 – Create product with minimum name length (2 chars) {#tc-ac-admin-001-b01}

**Preconditions:**
- Admin logged in

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| name | AB | Exactly 2 chars - minimum |

**Steps:**
1. Enter 2 character name
2. Fill other required fields
3. Submit

**Expected Result:**
- Product created successfully

**Traceability:**
- AC: [AC-ADMIN-001](../l1/acceptance-criteria.md#ac-admin-001)
- BR: [BR-ADMIN](../l1/business-rules.md#br-admin), [BR-PROD](../l1/business-rules.md#br-prod)

---

### TC-AC-ADMIN-001-B02 – Create product with maximum name length (200 chars) {#tc-ac-admin-001-b02}

**Preconditions:**
- Admin logged in

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| name | 200 character string | Exactly 200 chars |

**Steps:**
1. Enter exactly 200 character name
2. Submit

**Expected Result:**
- Product created successfully

**Traceability:**
- AC: [AC-ADMIN-001](../l1/acceptance-criteria.md#ac-admin-001)
- BR: [BR-ADMIN](../l1/business-rules.md#br-admin), [BR-PROD](../l1/business-rules.md#br-prod)

---

### TC-AC-ADMIN-001-B03 – Create product with minimum price ($0.01) {#tc-ac-admin-001-b03}

**Preconditions:**
- Admin logged in

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| price | 0.01 | Exactly minimum price |

**Steps:**
1. Enter valid name
2. Set price to $0.01
3. Submit

**Expected Result:**
- Product created with $0.01 price

**Traceability:**
- AC: [AC-ADMIN-001](../l1/acceptance-criteria.md#ac-admin-001)
- BR: [BR-ADMIN](../l1/business-rules.md#br-admin), [BR-PROD](../l1/business-rules.md#br-prod)

---

### TC-AC-ADMIN-001-B04 – Upload maximum 10 images to product {#tc-ac-admin-001-b04}

**Preconditions:**
- Admin creating product

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| images | 10 images | Maximum allowed |

**Steps:**
1. Create product
2. Upload exactly 10 images
3. Submit

**Expected Result:**
- All 10 images attached successfully

**Traceability:**
- AC: [AC-ADMIN-001](../l1/acceptance-criteria.md#ac-admin-001)
- BR: [BR-ADMIN](../l1/business-rules.md#br-admin), [BR-PROD](../l1/business-rules.md#br-prod)

---

### TC-AC-ADMIN-001-B05 – Reject uploading 11th image {#tc-ac-admin-001-b05}

**Preconditions:**
- Product has 10 images

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| images | 11 images | Exceeds maximum |

**Steps:**
1. Upload 10 images
2. Attempt to upload 11th

**Expected Result:**
- 11th image rejected
- Error message shown

**Traceability:**
- AC: [AC-ADMIN-001](../l1/acceptance-criteria.md#ac-admin-001)
- BR: [BR-ADMIN](../l1/business-rules.md#br-admin), [BR-PROD](../l1/business-rules.md#br-prod)

---

### TC-AC-ADMIN-001-B06 – Maximum description length (5000 chars) {#tc-ac-admin-001-b06}

**Preconditions:**
- Admin on add product form

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| description | 5000 char string | Maximum allowed |

**Steps:**
1. Enter 5000 character description
2. Submit

**Expected Result:**
- Product created with full description

**Traceability:**
- AC: [AC-ADMIN-001](../l1/acceptance-criteria.md#ac-admin-001)
- BR: [BR-ADMIN](../l1/business-rules.md#br-admin), [BR-PROD](../l1/business-rules.md#br-prod)

---

### TC-AC-ADMIN-002-B01 – Edit product name to exactly 2 characters {#tc-ac-admin-002-b01}

**Preconditions:**
- Product exists

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| name | AB | Minimum 2 chars |

**Steps:**
1. Edit product
2. Change name to 2 characters
3. Save

**Expected Result:**
- Product updated successfully

**Traceability:**
- AC: [AC-ADMIN-002](../l1/acceptance-criteria.md#ac-admin-002)
- BR: [BR-ADMIN](../l1/business-rules.md#br-admin), [BR-PROD](../l1/business-rules.md#br-prod)

---

### TC-AC-ADMIN-002-B02 – Edit product price to minimum $0.01 {#tc-ac-admin-002-b02}

**Preconditions:**
- Product exists with higher price

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| price | 0.01 | Minimum valid price |

**Steps:**
1. Edit product
2. Set price to $0.01
3. Save

**Expected Result:**
- Product updated with $0.01 price
- Audit trail logged

**Traceability:**
- AC: [AC-ADMIN-002](../l1/acceptance-criteria.md#ac-admin-002)
- BR: [BR-ADMIN](../l1/business-rules.md#br-admin), [BR-PROD](../l1/business-rules.md#br-prod)

---

### TC-AC-ADMIN-003-B01 – Delete product with zero order history {#tc-ac-admin-003-b01}

**Preconditions:**
- Product with no orders

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| productId | prod-404 | Product with zero orders |

**Steps:**
1. Delete product with no order history
2. Confirm

**Expected Result:**
- Product soft-deleted successfully

**Traceability:**
- AC: [AC-ADMIN-003](../l1/acceptance-criteria.md#ac-admin-003)
- BR: [BR-ADMIN](../l1/business-rules.md#br-admin), [BR-PROD](../l1/business-rules.md#br-prod)

---

### TC-AC-ADMIN-003-B02 – Delete product present in single cart {#tc-ac-admin-003-b02}

**Preconditions:**
- Product in exactly one cart

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| productId | prod-505 | In one cart |

**Steps:**
1. Delete product in one cart
2. Customer views their cart

**Expected Result:**
- Product removed from cart
- Notification shown

**Traceability:**
- AC: [AC-ADMIN-003](../l1/acceptance-criteria.md#ac-admin-003)
- BR: [BR-ADMIN](../l1/business-rules.md#br-admin), [BR-PROD](../l1/business-rules.md#br-prod)

---

### TC-AC-ADMIN-004-B01 – View orders at date range boundary {#tc-ac-admin-004-b01}

**Preconditions:**
- Order exists at exact boundary date

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| orderDate | 2024-01-31 23:59:59 | End of range |
| endDate | 2024-01-31 | Filter boundary |

**Steps:**
1. Set date range ending Jan 31
2. View orders

**Expected Result:**
- Order at 23:59:59 on Jan 31 included

**Traceability:**
- AC: [AC-ADMIN-004](../l1/acceptance-criteria.md#ac-admin-004)
- BR: [BR-ADMIN](../l1/business-rules.md#br-admin), [BR-ORDER](../l1/business-rules.md#br-order)

---

### TC-AC-ADMIN-004-B02 – Display single order in list {#tc-ac-admin-004-b02}

**Preconditions:**
- Only one order exists

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| orderCount | 1 | Single order |

**Steps:**
1. Navigate to order list with single order

**Expected Result:**
- Single order displayed correctly
- All fields visible

**Traceability:**
- AC: [AC-ADMIN-004](../l1/acceptance-criteria.md#ac-admin-004)
- BR: [BR-ADMIN](../l1/business-rules.md#br-admin), [BR-ORDER](../l1/business-rules.md#br-order)

---

### TC-AC-ADMIN-005-B01 – First valid transition in chain pending to confirmed {#tc-ac-admin-005-b01}

**Preconditions:**
- Order in initial pending state

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| currentStatus | pending | First state in chain |
| newStatus | confirmed | First valid transition |

**Steps:**
1. View pending order
2. Change to confirmed
3. Verify audit log

**Expected Result:**
- Transition succeeds
- Timestamp recorded accurately

**Traceability:**
- AC: [AC-ADMIN-005](../l1/acceptance-criteria.md#ac-admin-005)
- BR: [AMB-OP-047](../l1/business-rules.md#amb-op-047)

---

### TC-AC-ADMIN-005-B02 – Final valid transition shipped to delivered {#tc-ac-admin-005-b02}

**Preconditions:**
- Order in shipped state

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| currentStatus | shipped | Penultimate state |
| newStatus | delivered | Final state in chain |

**Steps:**
1. View shipped order
2. Change to delivered
3. Verify no further transitions

**Expected Result:**
- Transition succeeds
- No further status options available

**Traceability:**
- AC: [AC-ADMIN-005](../l1/acceptance-criteria.md#ac-admin-005)
- BR: [AMB-OP-048](../l1/business-rules.md#amb-op-048)

---

### TC-AC-ADMIN-006-B01 – Set stock to exactly zero via delta {#tc-ac-admin-006-b01}

**Preconditions:**
- Variant has stock of 10

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| currentStock | 10 | Starting value |
| delta | -10 | Exactly depletes stock |

**Steps:**
1. Apply delta -10 to reduce stock to exactly 0

**Expected Result:**
- Stock updated to 0
- No error - zero is valid minimum

**Traceability:**
- AC: [AC-ADMIN-006](../l1/acceptance-criteria.md#ac-admin-006)
- BR: [AMB-OP-051](../l1/business-rules.md#amb-op-051)

---

### TC-AC-ADMIN-006-B02 – Delta that would result in exactly -1 blocked {#tc-ac-admin-006-b02}

**Preconditions:**
- Variant has stock of 10

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| currentStock | 10 | Starting value |
| delta | -11 | Would result in -1 |

**Steps:**
1. Apply delta -11

**Expected Result:**
- Error: Inventory cannot go below zero
- Stock remains 10

**Traceability:**
- AC: [AC-ADMIN-006](../l1/acceptance-criteria.md#ac-admin-006)
- BR: [AMB-OP-054](../l1/business-rules.md#amb-op-054)

---

### TC-AC-VAR-001-B01 – Product with single variant type only size {#tc-ac-var-001-b01}

**Preconditions:**
- Product has only size variants, no color

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| variantTypes | size | Single type |
| options | S,M,L,XL | Size options only |

**Steps:**
1. View product with size-only variants

**Expected Result:**
- Only size selector displayed
- No color selector shown

**Traceability:**
- AC: [AC-VAR-001](../l1/acceptance-criteria.md#ac-var-001)
- BR: [AMB-ENT-008](../l1/business-rules.md#amb-ent-008)

---

### TC-AC-VAR-001-B02 – Variant with stock of exactly 1 unit {#tc-ac-var-001-b02}

**Preconditions:**
- Variant has exactly 1 unit in stock

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| variant | M-Blue | Target variant |
| stock | 1 | Minimum in-stock quantity |

**Steps:**
1. View product
2. Check M-Blue stock status

**Expected Result:**
- Shows as in stock or low stock
- Selectable for purchase

**Traceability:**
- AC: [AC-VAR-001](../l1/acceptance-criteria.md#ac-var-001)
- BR: [AMB-ENT-009](../l1/business-rules.md#amb-ent-009)

---

### TC-AC-CAT-001-B01 – Category at maximum depth level 3 {#tc-ac-cat-001-b01}

**Preconditions:**
- Category hierarchy exactly 3 levels deep

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| depth | 3 | Maximum allowed depth |
| path | L1 > L2 > L3 | Full hierarchy |

**Steps:**
1. Navigate to level 3 category
2. Verify products are accessible

**Expected Result:**
- Level 3 category displays correctly
- Products browsable

**Traceability:**
- AC: [AC-CAT-001](../l1/acceptance-criteria.md#ac-cat-001)
- BR: [AMB-ENT-006](../l1/business-rules.md#amb-ent-006)

---

### TC-AC-CAT-001-B02 – Single top-level category with no children {#tc-ac-cat-001-b02}

**Preconditions:**
- Category exists with no subcategories

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| category | Misc | Flat category |
| depth | 1 | No children |

**Steps:**
1. View category listing
2. Select Misc category

**Expected Result:**
- Category displays without expand option
- Products shown directly

**Traceability:**
- AC: [AC-CAT-001](../l1/acceptance-criteria.md#ac-cat-001)
- BR: [AMB-ENT-006](../l1/business-rules.md#amb-ent-006)

---

## Hallucination Prevention Tests

### TC-AC-PROD-001-H01 – Search text filter should NOT be available {#tc-ac-prod-001-h01}

**⚠️ Should NOT:** Text search or keyword filter should NOT be present

**Preconditions:**
- User is on product listing page

**Steps:**
1. Examine available filter options on product listing

**Expected Result:**
- Only category, price range, availability filters present

**Traceability:**
- AC: [AC-PROD-001](../l1/acceptance-criteria.md#ac-prod-001)
- BR: [BR-PROD](../l1/business-rules.md#br-prod)

---

### TC-AC-PROD-001-H02 – Rating filter should NOT be available {#tc-ac-prod-001-h02}

**⚠️ Should NOT:** Rating or review score filter should NOT exist

**Preconditions:**
- User is on product listing page

**Steps:**
1. Check all available sorting and filtering options

**Expected Result:**
- No rating-based filter exists

**Traceability:**
- AC: [AC-PROD-001](../l1/acceptance-criteria.md#ac-prod-001)
- BR: [BR-PROD](../l1/business-rules.md#br-prod)

---

### TC-AC-PROD-002-H01 – Waitlist or notify option should NOT be shown {#tc-ac-prod-002-h01}

**⚠️ Should NOT:** Waitlist signup or 'Notify when available' should NOT exist

**Preconditions:**
- Product is out of stock

**Steps:**
1. View out-of-stock product detail page
2. Check for notify/waitlist

**Expected Result:**
- Only Out of Stock indicator present

**Traceability:**
- AC: [AC-PROD-002](../l1/acceptance-criteria.md#ac-prod-002)
- BR: [BR-PROD](../l1/business-rules.md#br-prod)

---

### TC-AC-PROD-002-H02 – Pre-order option should NOT be available for out of stock {#tc-ac-prod-002-h02}

**⚠️ Should NOT:** Pre-order or backorder functionality should NOT be present

**Preconditions:**
- Product has zero inventory

**Steps:**
1. View out-of-stock product
2. Check for alternative purchase options

**Expected Result:**
- No pre-order or backorder option visible

**Traceability:**
- AC: [AC-PROD-002](../l1/acceptance-criteria.md#ac-prod-002)
- BR: [BR-PROD](../l1/business-rules.md#br-prod)

---

### TC-AC-CART-001-H01 – Wishlist option should NOT appear in add to cart flow {#tc-ac-cart-001-h01}

**⚠️ Should NOT:** Wishlist or Save for Later should NOT be in this flow

**Preconditions:**
- User viewing product

**Steps:**
1. View product page
2. Check available actions

**Expected Result:**
- Only Add to Cart functionality present

**Traceability:**
- AC: [AC-CART-001](../l1/acceptance-criteria.md#ac-cart-001)
- BR: [BR-CART](../l1/business-rules.md#br-cart)

---

### TC-AC-CART-001-H02 – Buy Now button should NOT be available {#tc-ac-cart-001-h02}

**⚠️ Should NOT:** Buy Now or instant checkout should NOT exist

**Preconditions:**
- User on product page

**Steps:**
1. View product detail page
2. Check for instant purchase options

**Expected Result:**
- Only Add to Cart button present

**Traceability:**
- AC: [AC-CART-001](../l1/acceptance-criteria.md#ac-cart-001)
- BR: [BR-CART](../l1/business-rules.md#br-cart)

---

### TC-AC-CART-002-H01 – Different variants should NOT merge into same line item {#tc-ac-cart-002-h01}

**⚠️ Should NOT:** Different variants should NOT merge into single line item

**Preconditions:**
- Product A (Size M) in cart
- Adding Product A (Size L)

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| variant1 | Size: M | In cart |
| variant2 | Size: L | Being added |

**Steps:**
1. Add Product A Size L when Size M already in cart

**Expected Result:**
- Two separate line items created

**Traceability:**
- AC: [AC-CART-002](../l1/acceptance-criteria.md#ac-cart-002)
- BR: [BR-CART](../l1/business-rules.md#br-cart)

---

### TC-AC-CART-002-H02 – Cart should NOT automatically remove lower-priced duplicate {#tc-ac-cart-002-h02}

**⚠️ Should NOT:** System should NOT auto-remove items or change prices unexpectedly

**Preconditions:**
- Same product at different prices due to sale

**Steps:**
1. Check cart behavior when adding same product

**Expected Result:**
- Quantities merge, price handled per business rules

**Traceability:**
- AC: [AC-CART-002](../l1/acceptance-criteria.md#ac-cart-002)
- BR: [BR-CART](../l1/business-rules.md#br-cart)

---

### TC-AC-CART-003-H01 – Guest cart should NOT require email to create {#tc-ac-cart-003-h01}

**⚠️ Should NOT:** Email or any personal info should NOT be required for guest cart

**Preconditions:**
- User not logged in

**Steps:**
1. Attempt to add to cart as guest
2. Check if email is requested

**Expected Result:**
- No email required to create cart

**Traceability:**
- AC: [AC-CART-003](../l1/acceptance-criteria.md#ac-cart-003)
- BR: [BR-CART](../l1/business-rules.md#br-cart)

---

### TC-AC-CART-003-H02 – Guest cart should NOT auto-merge with logged-in user cart {#tc-ac-cart-003-h02}

**⚠️ Should NOT:** Cart merge behavior is NOT specified - should NOT assume auto-merge

**Preconditions:**
- Guest has items in cart
- User logs in

**Steps:**
1. Add items as guest
2. Log in to existing account
3. Check cart

**Expected Result:**
- Behavior is undefined in spec, test actual implementation

**Traceability:**
- AC: [AC-CART-003](../l1/acceptance-criteria.md#ac-cart-003)
- BR: [BR-CART](../l1/business-rules.md#br-cart)

---

### TC-AC-CART-004-H01 – Guest cart items should NOT be lost after merge {#tc-ac-cart-004-h01}

**⚠️ Should NOT:** Guest cart items should NOT be discarded during merge

**Preconditions:**
- Guest has unique items not in account cart

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| guestItems | [SKU-X:1] | Unique to guest |

**Steps:**
1. Add unique item to guest cart
2. Login
3. Verify item preserved

**Expected Result:**
- Guest-only items appear in merged cart

**Traceability:**
- AC: [AC-CART-004](../l1/acceptance-criteria.md#ac-cart-004)
- BR: [AMB-ENT-013](../l1/business-rules.md#amb-ent-013)

---

### TC-AC-CART-004-H02 – Merge should NOT create duplicate line items {#tc-ac-cart-004-h02}

**⚠️ Should NOT:** Should NOT create two separate line items for same SKU

**Preconditions:**
- Same SKU exists in both carts

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| sku | SKU-A | Duplicate product |

**Steps:**
1. Have same SKU in both carts
2. Login
3. Check line items

**Expected Result:**
- Single line item with combined quantity

**Traceability:**
- AC: [AC-CART-004](../l1/acceptance-criteria.md#ac-cart-004)
- BR: [AMB-ENT-013](../l1/business-rules.md#amb-ent-013)

---

### TC-AC-CART-005-H01 – Update should NOT affect other items in cart {#tc-ac-cart-005-h01}

**⚠️ Should NOT:** Should NOT modify quantities of other cart items

**Preconditions:**
- Cart has SKU-A:2 and SKU-B:3

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| targetItem | SKU-A | Item to update |
| otherItem | SKU-B | Should not change |

**Steps:**
1. Update SKU-A quantity
2. Check SKU-B quantity

**Expected Result:**
- SKU-B quantity unchanged at 3

**Traceability:**
- AC: [AC-CART-005](../l1/acceptance-criteria.md#ac-cart-005)
- BR: [AMB-OP-011](../l1/business-rules.md#amb-op-011)

---

### TC-AC-CART-005-H02 – Negative quantity should NOT be accepted {#tc-ac-cart-005-h02}

**⚠️ Should NOT:** Should NOT accept negative quantities

**Preconditions:**
- Cart has item

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| newQuantity | -1 | Invalid negative |

**Steps:**
1. Attempt to set quantity to -1
2. Verify rejection

**Expected Result:**
- Negative quantity rejected or treated as 0

**Traceability:**
- AC: [AC-CART-005](../l1/acceptance-criteria.md#ac-cart-005)
- BR: [AMB-OP-011](../l1/business-rules.md#amb-op-011)

---

### TC-AC-CART-006-H01 – Removed item should NOT appear in cart after refresh {#tc-ac-cart-006-h01}

**⚠️ Should NOT:** Removed item should NOT reappear after page refresh

**Preconditions:**
- Item removed
- Undo not clicked

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| removedItem | SKU-A | Removed item |

**Steps:**
1. Remove item
2. Wait for undo to expire
3. Refresh page

**Expected Result:**
- Removed item not in cart

**Traceability:**
- AC: [AC-CART-006](../l1/acceptance-criteria.md#ac-cart-006)
- BR: [AMB-OP-014](../l1/business-rules.md#amb-op-014)

---

### TC-AC-CART-006-H02 – Remove should NOT require confirmation dialog {#tc-ac-cart-006-h02}

**⚠️ Should NOT:** Should NOT show confirmation dialog before removal

**Preconditions:**
- Cart has items

**Steps:**
1. Click remove on item
2. Observe behavior

**Expected Result:**
- Item removed immediately

**Traceability:**
- AC: [AC-CART-006](../l1/acceptance-criteria.md#ac-cart-006)
- BR: [AMB-OP-014](../l1/business-rules.md#amb-op-014)

---

### TC-AC-CART-007-H01 – Out of stock item should NOT be auto-removed {#tc-ac-cart-007-h01}

**⚠️ Should NOT:** Should NOT automatically remove out of stock items

**Preconditions:**
- Item goes out of stock

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| item | SKU-A | Out of stock |

**Steps:**
1. Item goes out of stock
2. Refresh cart
3. Verify presence

**Expected Result:**
- Item still present in cart

**Traceability:**
- AC: [AC-CART-007](../l1/acceptance-criteria.md#ac-cart-007)
- BR: [AMB-ENT-014](../l1/business-rules.md#amb-ent-014)

---

### TC-AC-CART-007-H02 – Other in-stock items should NOT be blocked {#tc-ac-cart-007-h02}

**⚠️ Should NOT:** Should NOT block checkout for items that ARE in stock

**Preconditions:**
- Cart has SKU-A (out of stock) and SKU-B (in stock)

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| outOfStock | SKU-A | Blocked item |
| inStock | SKU-B | Available item |

**Steps:**
1. Remove out of stock item
2. Proceed to checkout

**Expected Result:**
- Checkout succeeds for in-stock items

**Traceability:**
- AC: [AC-CART-007](../l1/acceptance-criteria.md#ac-cart-007)
- BR: [AMB-ENT-014](../l1/business-rules.md#amb-ent-014)

---

### TC-AC-CART-008-H01 – Should NOT lock in price from time of adding to cart {#tc-ac-cart-008-h01}

**⚠️ Should NOT:** Should NOT honor original price from when item was added

**Preconditions:**
- Item added at old price
- Price changed

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| addedAt | 10.00 | Price when added |
| current | 12.00 | Current price |

**Steps:**
1. Check price displayed in cart

**Expected Result:**
- Current price shown, not original

**Traceability:**
- AC: [AC-CART-008](../l1/acceptance-criteria.md#ac-cart-008)
- BR: [AMB-ENT-016](../l1/business-rules.md#amb-ent-016)

---

### TC-AC-CART-008-H02 – Price change should NOT auto-remove item from cart {#tc-ac-cart-008-h02}

**⚠️ Should NOT:** Should NOT auto-remove items due to price changes

**Preconditions:**
- Item price significantly increased

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| oldPrice | 10.00 | Original |
| newPrice | 50.00 | 5x increase |

**Steps:**
1. Significant price change occurs
2. View cart

**Expected Result:**
- Item remains in cart with new price

**Traceability:**
- AC: [AC-CART-008](../l1/acceptance-criteria.md#ac-cart-008)
- BR: [AMB-ENT-016](../l1/business-rules.md#amb-ent-016)

---

### TC-AC-ORDER-001-H01 – Order should NOT be created if payment not authorized {#tc-ac-order-001-h01}

**⚠️ Should NOT:** Order should NOT be created before payment authorization completes

**Preconditions:**
- Payment authorization pending or failed

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| paymentStatus | pending | Auth not completed |

**Steps:**
1. Interrupt payment authorization
2. Check order status

**Expected Result:**
- No order record exists in database

**Traceability:**
- AC: [AC-ORDER-001](../l1/acceptance-criteria.md#ac-order-001)
- BR: [BR-ORDER](../l1/business-rules.md#br-order)

---

### TC-AC-ORDER-001-H02 – Inventory should NOT decrement on failed payment {#tc-ac-order-001-h02}

**⚠️ Should NOT:** Inventory should NOT be decremented when payment fails

**Preconditions:**
- Payment will fail

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| initialStock | 10 | Track stock before attempt |

**Steps:**
1. Record initial inventory
2. Attempt order with failing payment
3. Check inventory

**Expected Result:**
- Inventory remains at initial level

**Traceability:**
- AC: [AC-ORDER-001](../l1/acceptance-criteria.md#ac-order-001)
- BR: [BR-ORDER](../l1/business-rules.md#br-order)

---

### TC-AC-ORDER-002-H01 – Tax should NOT be included in free shipping calculation {#tc-ac-order-002-h01}

**⚠️ Should NOT:** Tax should NOT be included when calculating free shipping eligibility

**Preconditions:**
- Subtotal $45, tax $8

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| subtotal | $45.00 | Below threshold |
| tax | $8.00 | Would push total over $50 |

**Steps:**
1. Add items with $45 subtotal
2. Tax calculated as $8
3. Check shipping

**Expected Result:**
- $5.99 shipping fee applied despite $53 total

**Traceability:**
- AC: [AC-ORDER-002](../l1/acceptance-criteria.md#ac-order-002)
- BR: [BR-SHIPPING](../l1/business-rules.md#br-shipping)

---

### TC-AC-ORDER-002-H02 – Previous order total should NOT affect current shipping {#tc-ac-order-002-h02}

**⚠️ Should NOT:** Previous order amounts should NOT affect current shipping calculation

**Preconditions:**
- Customer had $100 order last week
- Current cart is $30

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| currentSubtotal | $30.00 | Current order below threshold |

**Steps:**
1. Login as customer with order history
2. Add $30 to cart
3. Checkout

**Expected Result:**
- $5.99 shipping applied to current order

**Traceability:**
- AC: [AC-ORDER-002](../l1/acceptance-criteria.md#ac-order-002)
- BR: [BR-SHIPPING](../l1/business-rules.md#br-shipping)

---

### TC-AC-ORDER-003-H01 – Billing address should NOT affect tax calculation {#tc-ac-order-003-h01}

**⚠️ Should NOT:** Billing address state should NOT be used for tax calculation

**Preconditions:**
- Different billing and shipping states

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| billingState | CA | High tax state |
| shippingState | OR | No tax state |

**Steps:**
1. Enter CA billing address
2. Enter OR shipping address
3. Check tax

**Expected Result:**
- Tax is $0 based on OR shipping address

**Traceability:**
- AC: [AC-ORDER-003](../l1/acceptance-criteria.md#ac-order-003)
- BR: [BR-TAX](../l1/business-rules.md#br-tax)

---

### TC-AC-ORDER-003-H02 – Customer home state should NOT override shipping tax {#tc-ac-order-003-h02}

**⚠️ Should NOT:** Customer profile state should NOT be used for tax calculation

**Preconditions:**
- Customer profile has home state different from shipping

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| profileState | NY | Customer home state |
| shippingState | NH | No tax state for shipping |

**Steps:**
1. Ship to NH address
2. Verify tax calculation

**Expected Result:**
- Tax based on NH shipping, not NY profile

**Traceability:**
- AC: [AC-ORDER-003](../l1/acceptance-criteria.md#ac-order-003)
- BR: [BR-TAX](../l1/business-rules.md#br-tax)

---

### TC-AC-ORDER-004-H01 – Email should NOT contain payment card details {#tc-ac-order-004-h01}

**⚠️ Should NOT:** Email should NOT contain full credit card number or CVV

**Preconditions:**
- Order placed with credit card

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| cardNumber | 4111111111111111 | Test card used |

**Steps:**
1. Place order with credit card
2. Check email content

**Expected Result:**
- Email contains only specified fields

**Traceability:**
- AC: [AC-ORDER-004](../l1/acceptance-criteria.md#ac-order-004)
- BR: [BR-EMAIL](../l1/business-rules.md#br-email)

---

### TC-AC-ORDER-004-H02 – Email should NOT be sent synchronously {#tc-ac-order-004-h02}

**⚠️ Should NOT:** Order completion should NOT wait for email delivery confirmation

**Preconditions:**
- Order being placed

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| emailSendTime | measured | Track timing |

**Steps:**
1. Place order
2. Measure if order waits for email send

**Expected Result:**
- Order returns before email actually sent

**Traceability:**
- AC: [AC-ORDER-004](../l1/acceptance-criteria.md#ac-order-004)
- BR: [BR-EMAIL](../l1/business-rules.md#br-email)

---

### TC-AC-ORDER-005-H01 – Should NOT show other customers orders {#tc-ac-order-005-h01}

**⚠️ Should NOT:** Other customers orders should NOT appear in order history

**Preconditions:**
- Multiple customers with orders

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| customerAOrders | 5 | Current user orders |
| customerBOrders | 10 | Other user orders |

**Steps:**
1. Login as Customer A
2. View order history

**Expected Result:**
- Only 5 orders visible (Customer A orders)

**Traceability:**
- AC: [AC-ORDER-005](../l1/acceptance-criteria.md#ac-order-005)
- BR: [BR-TRACKING](../l1/business-rules.md#br-tracking)

---

### TC-AC-ORDER-005-H02 – Cancelled orders should NOT be hidden from history {#tc-ac-order-005-h02}

**⚠️ Should NOT:** Cancelled orders should NOT be filtered out from order history

**Preconditions:**
- Customer has cancelled orders

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| cancelledOrders | 2 | Orders with cancelled status |

**Steps:**
1. Login as customer with cancelled orders
2. View order history

**Expected Result:**
- Cancelled orders visible in history
- Status shows as cancelled

**Traceability:**
- AC: [AC-ORDER-005](../l1/acceptance-criteria.md#ac-order-005)
- BR: [BR-TRACKING](../l1/business-rules.md#br-tracking)

---

### TC-AC-ORDER-006-H01 – Cancelled order should NOT appear in active orders {#tc-ac-order-006-h01}

**⚠️ Should NOT:** Cancelled order should NOT appear in active/pending orders list

**Preconditions:**
- Order was cancelled

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| orderId | ORD-12345 | Cancelled order |

**Steps:**
1. View active orders list
2. Search for cancelled order

**Expected Result:**
- Order not in active orders list

**Traceability:**
- AC: [AC-ORDER-006](../l1/acceptance-criteria.md#ac-order-006)
- BR: [AMB-OP-028](../l1/business-rules.md#amb-op-028)

---

### TC-AC-ORDER-006-H02 – Partial cancellation should NOT be allowed {#tc-ac-order-006-h02}

**⚠️ Should NOT:** System should NOT allow cancelling individual items (whole order only)

**Preconditions:**
- Order has multiple items

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| orderItems | 3 | Multi-item order |

**Steps:**
1. Attempt to cancel single item from order

**Expected Result:**
- Partial cancellation not available

**Traceability:**
- AC: [AC-ORDER-006](../l1/acceptance-criteria.md#ac-order-006)
- BR: [AMB-OP-028](../l1/business-rules.md#amb-op-028)

---

### TC-AC-ORDER-007-H01 – Order should NOT auto-update to new prices {#tc-ac-order-007-h01}

**⚠️ Should NOT:** System should NOT have auto-sync or price update feature for orders

**Preconditions:**
- Order placed
- Product price changed

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| orderDate | 2024-01-01 | Order date |
| priceChangeDate | 2024-01-02 | Price update |

**Steps:**
1. Check order after price change
2. Run any batch jobs

**Expected Result:**
- Order prices remain at original captured values

**Traceability:**
- AC: [AC-ORDER-007](../l1/acceptance-criteria.md#ac-order-007)
- BR: [AMB-ENT-019](../l1/business-rules.md#amb-ent-019)

---

### TC-AC-ORDER-007-H02 – Customer should NOT see current price on order history {#tc-ac-order-007-h02}

**⚠️ Should NOT:** Order history should NOT display current catalog prices

**Preconditions:**
- Order placed at $29.99
- Current price is $39.99

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| orderPrice | 29.99 | Captured |
| currentPrice | 39.99 | Current catalog |

**Steps:**
1. View order history

**Expected Result:**
- Shows only captured order price

**Traceability:**
- AC: [AC-ORDER-007](../l1/acceptance-criteria.md#ac-order-007)
- BR: [AMB-ENT-019](../l1/business-rules.md#amb-ent-019)

---

### TC-AC-CUST-001-H01 – Phone number should NOT be required for registration {#tc-ac-cust-001-h01}

**⚠️ Should NOT:** System should NOT require phone number (not in spec)

**Preconditions:**
- Registration form open

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| phone |  | Not provided |

**Steps:**
1. Register without phone number

**Expected Result:**
- Registration succeeds without phone

**Traceability:**
- AC: [AC-CUST-001](../l1/acceptance-criteria.md#ac-cust-001)
- BR: [AMB-ENT-024](../l1/business-rules.md#amb-ent-024)

---

### TC-AC-CUST-001-H02 – Unverified user should NOT be able to place orders {#tc-ac-cust-001-h02}

**⚠️ Should NOT:** Unverified users should NOT be able to complete orders

**Preconditions:**
- User registered but not verified

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| verificationStatus | pending | Not verified |

**Steps:**
1. Login as unverified user
2. Attempt to place order

**Expected Result:**
- Order placement blocked
- Prompt to verify email

**Traceability:**
- AC: [AC-CUST-001](../l1/acceptance-criteria.md#ac-cust-001)
- BR: [AMB-ENT-030](../l1/business-rules.md#amb-ent-030)

---

### TC-AC-CUST-002-H01 – Billing address should NOT be required for shipping {#tc-ac-cust-002-h01}

**⚠️ Should NOT:** System should NOT require billing address when adding shipping

**Preconditions:**
- Add shipping address form open

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| billingAddress |  | Not provided |

**Steps:**
1. Add shipping address only

**Expected Result:**
- Shipping address saved without billing

**Traceability:**
- AC: [AC-CUST-002](../l1/acceptance-criteria.md#ac-cust-002)
- BR: [AMB-ENT-027](../l1/business-rules.md#amb-ent-027)

---

### TC-AC-CUST-002-H02 – Email should NOT be required per address {#tc-ac-cust-002-h02}

**⚠️ Should NOT:** System should NOT require email field for each address

**Preconditions:**
- Add address form open

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| addressEmail |  | Not in spec |

**Steps:**
1. Add address without separate email field

**Expected Result:**
- Address saved without requiring email

**Traceability:**
- AC: [AC-CUST-002](../l1/acceptance-criteria.md#ac-cust-002)
- BR: [AMB-ENT-027](../l1/business-rules.md#amb-ent-027)

---

### TC-AC-CUST-003-H01 – Raw card numbers should NOT be stored in database {#tc-ac-cust-003-h01}

**⚠️ Should NOT:** System should NOT store raw card numbers (only tokens)

**Preconditions:**
- Card has been added

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| cardNumber | 4242424242424242 | Original card |

**Steps:**
1. Add card
2. Query database directly

**Expected Result:**
- Only token stored, no raw card number

**Traceability:**
- AC: [AC-CUST-003](../l1/acceptance-criteria.md#ac-cust-003)
- BR: [AMB-ENT-028](../l1/business-rules.md#amb-ent-028)

---

### TC-AC-CUST-003-H02 – Bank account/ACH should NOT be supported {#tc-ac-cust-003-h02}

**⚠️ Should NOT:** System should NOT offer ACH/bank account (not in spec)

**Preconditions:**
- Payment methods page

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| paymentType | ACH | Bank transfer |

**Steps:**
1. Check available payment method options

**Expected Result:**
- ACH/bank account option not available

**Traceability:**
- AC: [AC-CUST-003](../l1/acceptance-criteria.md#ac-cust-003)
- BR: [AMB-ENT-028](../l1/business-rules.md#amb-ent-028)

---

### TC-AC-CUST-004-H01 – Soft-deleted customer should NOT appear in active lists {#tc-ac-cust-004-h01}

**⚠️ Should NOT:** Soft-deleted customer should NOT appear in active customer searches

**Preconditions:**
- Customer account soft-deleted

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| customerId | cust-555 | Soft-deleted customer |

**Steps:**
1. Soft delete customer
2. Search active customers
3. View customer lists

**Expected Result:**
- Customer not in active lists

**Traceability:**
- AC: [AC-CUST-004](../l1/acceptance-criteria.md#ac-cust-004)
- BR: [BR-CUST](../l1/business-rules.md#br-cust)

---

### TC-AC-CUST-004-H02 – GDPR erased data should NOT be recoverable {#tc-ac-cust-004-h02}

**⚠️ Should NOT:** Erased personal data should NOT be recoverable by any means

**Preconditions:**
- Customer selected full erasure

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| customerId | cust-666 | GDPR erased customer |

**Steps:**
1. Perform full GDPR erasure
2. Attempt to access old customer data

**Expected Result:**
- No personal data retrievable

**Traceability:**
- AC: [AC-CUST-004](../l1/acceptance-criteria.md#ac-cust-004)
- BR: [BR-CUST](../l1/business-rules.md#br-cust)

---

### TC-AC-ADMIN-001-H01 – Draft product should NOT appear in customer catalog {#tc-ac-admin-001-h01}

**⚠️ Should NOT:** Draft products should NOT be visible in customer-facing catalog

**Preconditions:**
- New product in draft status

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| status | draft | Initial status |

**Steps:**
1. Create product in draft
2. Browse catalog as customer

**Expected Result:**
- Product not visible to customers

**Traceability:**
- AC: [AC-ADMIN-001](../l1/acceptance-criteria.md#ac-admin-001)
- BR: [BR-ADMIN](../l1/business-rules.md#br-admin), [BR-PROD](../l1/business-rules.md#br-prod)

---

### TC-AC-ADMIN-001-H02 – Product should NOT require SKU field {#tc-ac-admin-001-h02}

**⚠️ Should NOT:** System should NOT require SKU as it is not in requirements

**Preconditions:**
- Admin on add product form

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| sku | null | Not in requirements |

**Steps:**
1. Create product without SKU
2. Submit with required fields only

**Expected Result:**
- Product created without SKU

**Traceability:**
- AC: [AC-ADMIN-001](../l1/acceptance-criteria.md#ac-admin-001)
- BR: [BR-ADMIN](../l1/business-rules.md#br-admin), [BR-PROD](../l1/business-rules.md#br-prod)

---

### TC-AC-ADMIN-002-H01 – Non-price changes should NOT appear in price audit {#tc-ac-admin-002-h01}

**⚠️ Should NOT:** Description-only changes should NOT create price audit entries

**Preconditions:**
- Product exists

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| description | New description | Non-price change |

**Steps:**
1. Edit product description only
2. Save
3. Check audit trail

**Expected Result:**
- No price audit entry created

**Traceability:**
- AC: [AC-ADMIN-002](../l1/acceptance-criteria.md#ac-admin-002)
- BR: [BR-ADMIN](../l1/business-rules.md#br-admin), [BR-PROD](../l1/business-rules.md#br-prod)

---

### TC-AC-ADMIN-002-H02 – Edit should NOT change product UUID {#tc-ac-admin-002-h02}

**⚠️ Should NOT:** Product UUID should NOT change when editing

**Preconditions:**
- Product exists with UUID

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| originalUUID | uuid-abc-123 | Original ID |

**Steps:**
1. Note product UUID
2. Edit and save product
3. Check UUID

**Expected Result:**
- UUID unchanged after edit

**Traceability:**
- AC: [AC-ADMIN-002](../l1/acceptance-criteria.md#ac-admin-002)
- BR: [BR-ADMIN](../l1/business-rules.md#br-admin), [BR-PROD](../l1/business-rules.md#br-prod)

---

### TC-AC-ADMIN-003-H01 – Soft-deleted product should NOT be purchasable {#tc-ac-admin-003-h01}

**⚠️ Should NOT:** Soft-deleted products should NOT be addable to cart

**Preconditions:**
- Product soft-deleted

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| productId | prod-606 | Deleted product |

**Steps:**
1. Soft delete product
2. Attempt direct add to cart via API

**Expected Result:**
- Add to cart fails

**Traceability:**
- AC: [AC-ADMIN-003](../l1/acceptance-criteria.md#ac-admin-003)
- BR: [BR-ADMIN](../l1/business-rules.md#br-admin), [BR-PROD](../l1/business-rules.md#br-prod)

---

### TC-AC-ADMIN-003-H02 – Deleted product should NOT be restored without admin action {#tc-ac-admin-003-h02}

**⚠️ Should NOT:** Soft-deleted products should NOT auto-restore

**Preconditions:**
- Product soft-deleted

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| productId | prod-707 | Deleted product |

**Steps:**
1. Soft delete product
2. Wait period
3. Check product status

**Expected Result:**
- Product remains deleted

**Traceability:**
- AC: [AC-ADMIN-003](../l1/acceptance-criteria.md#ac-admin-003)
- BR: [BR-ADMIN](../l1/business-rules.md#br-admin), [BR-PROD](../l1/business-rules.md#br-prod)

---

### TC-AC-ADMIN-004-H01 – Order list should NOT show deleted orders {#tc-ac-admin-004-h01}

**⚠️ Should NOT:** Deleted/cancelled orders should NOT appear in default view

**Preconditions:**
- Some orders are soft-deleted

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| deletedOrderId | ORD-DEL-001 | Deleted order |

**Steps:**
1. View order list
2. Search for deleted order

**Expected Result:**
- Deleted orders not in list

**Traceability:**
- AC: [AC-ADMIN-004](../l1/acceptance-criteria.md#ac-admin-004)
- BR: [BR-ADMIN](../l1/business-rules.md#br-admin), [BR-ORDER](../l1/business-rules.md#br-order)

---

### TC-AC-ADMIN-004-H02 – Filter should NOT allow SQL injection {#tc-ac-admin-004-h02}

**⚠️ Should NOT:** Filter/search should NOT execute injected SQL

**Preconditions:**
- Admin on order management

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| searchTerm | '; DROP TABLE orders;-- | Injection attempt |

**Steps:**
1. Enter SQL injection in search field
2. Execute search

**Expected Result:**
- Search fails safely or returns no results

**Traceability:**
- AC: [AC-ADMIN-004](../l1/acceptance-criteria.md#ac-admin-004)
- BR: [BR-ADMIN](../l1/business-rules.md#br-admin), [BR-ORDER](../l1/business-rules.md#br-order)

---

### TC-AC-ADMIN-005-H01 – Should NOT allow cancelled status not in spec {#tc-ac-admin-005-h01}

**⚠️ Should NOT:** Should NOT allow cancelled status as it is not in the defined transitions

**Preconditions:**
- Order exists in any status

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| newStatus | cancelled | Not in defined transitions |

**Steps:**
1. View order
2. Attempt to set cancelled status

**Expected Result:**
- Cancelled status not available or rejected

**Traceability:**
- AC: [AC-ADMIN-005](../l1/acceptance-criteria.md#ac-admin-005)

---

### TC-AC-ADMIN-005-H02 – Should NOT send email for pending status {#tc-ac-admin-005-h02}

**⚠️ Should NOT:** Should NOT send customer notification email when order is pending

**Preconditions:**
- New order created

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| status | pending | Initial status |

**Steps:**
1. Create new order
2. Order enters pending status
3. Check email logs

**Expected Result:**
- No email sent for pending status

**Traceability:**
- AC: [AC-ADMIN-005](../l1/acceptance-criteria.md#ac-admin-005)
- BR: [AMB-OP-049](../l1/business-rules.md#amb-op-049)

---

### TC-AC-ADMIN-006-H01 – Should NOT allow inventory for product without variants {#tc-ac-admin-006-h01}

**⚠️ Should NOT:** Should NOT allow inventory tracking at product level, only per variant

**Preconditions:**
- Product exists without any variants defined

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| productId | PROD-001 | Product without variants |

**Steps:**
1. Navigate to product
2. Attempt to adjust inventory at product level

**Expected Result:**
- No product-level inventory adjustment available

**Traceability:**
- AC: [AC-ADMIN-006](../l1/acceptance-criteria.md#ac-admin-006)
- BR: [AMB-ENT-009](../l1/business-rules.md#amb-ent-009)

---

### TC-AC-ADMIN-006-H02 – Should NOT skip audit trail for any adjustment {#tc-ac-admin-006-h02}

**⚠️ Should NOT:** Should NOT allow any inventory change without audit trail entry

**Preconditions:**
- Any inventory adjustment

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| delta | +1 | Minimal change |

**Steps:**
1. Make small adjustment
2. Check audit log

**Expected Result:**
- Audit entry exists with before/after values

**Traceability:**
- AC: [AC-ADMIN-006](../l1/acceptance-criteria.md#ac-admin-006)
- BR: [AMB-ENT-041](../l1/business-rules.md#amb-ent-041), [AMB-ENT-042](../l1/business-rules.md#amb-ent-042)

---

### TC-AC-VAR-001-H01 – Should NOT show aggregated product-level stock {#tc-ac-var-001-h01}

**⚠️ Should NOT:** Should NOT display combined/aggregated stock across all variants

**Preconditions:**
- Product has multiple variants with different stock levels

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| variant1Stock | 10 | First variant |
| variant2Stock | 0 | Second variant |

**Steps:**
1. View product page
2. Check for aggregated stock display

**Expected Result:**
- No total product stock shown

**Traceability:**
- AC: [AC-VAR-001](../l1/acceptance-criteria.md#ac-var-001)
- BR: [AMB-ENT-009](../l1/business-rules.md#amb-ent-009)

---

### TC-AC-VAR-001-H02 – Should NOT auto-substitute out of stock variant {#tc-ac-var-001-h02}

**⚠️ Should NOT:** Should NOT automatically substitute with different variant

**Preconditions:**
- Selected variant out of stock
- Similar variant in stock

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| selected | M-Red | Out of stock |
| similar | M-Blue | In stock alternative |

**Steps:**
1. Select M-Red variant
2. Observe system behavior

**Expected Result:**
- System does not auto-switch to M-Blue

**Traceability:**
- AC: [AC-VAR-001](../l1/acceptance-criteria.md#ac-var-001)

---

### TC-AC-CAT-001-H01 – Should NOT display fourth level categories {#tc-ac-cat-001-h01}

**⚠️ Should NOT:** Should NOT allow or display categories deeper than 3 levels

**Preconditions:**
- Attempt to create L1 > L2 > L3 > L4 hierarchy

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| attemptedDepth | 4 | Beyond maximum |

**Steps:**
1. Attempt to view 4-level deep category structure

**Expected Result:**
- Fourth level not displayed or not creatable

**Traceability:**
- AC: [AC-CAT-001](../l1/acceptance-criteria.md#ac-cat-001)
- BR: [AMB-ENT-006](../l1/business-rules.md#amb-ent-006)

---

### TC-AC-CAT-001-H02 – Should NOT require secondary category assignment {#tc-ac-cat-001-h02}

**⚠️ Should NOT:** Should NOT require secondary category to be assigned

**Preconditions:**
- Product with only primary category assigned

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| productId | PROD-SIMPLE | Simple product |
| primaryOnly | true | No secondary |

**Steps:**
1. Create product with primary category only
2. Save product

**Expected Result:**
- Product saves successfully
- No validation error

**Traceability:**
- AC: [AC-CAT-001](../l1/acceptance-criteria.md#ac-cat-001)
- BR: [AMB-ENT-039](../l1/business-rules.md#amb-ent-039)

---

