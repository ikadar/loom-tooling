# TDAI Test Cases

Generated: 2026-01-02T21:37:12+01:00

**Methodology:** Test-Driven AI Development (TDAI)

## Summary

| Category | Count | Ratio |
|----------|-------|-------|
| Positive | 90 | 31.8% |
| Negative | 72 | 25.4% |
| Boundary | 63 | - |
| Hallucination Prevention | 58 | - |
| **Total** | **283** | - |

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
3. Apply filter

**Expected Result:**
- Only Electronics products are displayed

**Traceability:**
- AC: [AC-PROD-001](../l1/acceptance-criteria.md#ac-prod-001)
- BR: [BR-PROD](../l1/business-rules.md#br-prod)

---

### TC-AC-PROD-001-P02 – Filter products by price range and category combined {#tc-ac-prod-001-p02}

**Preconditions:**
- Products exist with various prices
- User is on product listing page

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| category | Electronics | Valid category |
| priceMin | 50 | Minimum price |
| priceMax | 200 | Maximum price |

**Steps:**
1. Select category 'Electronics'
2. Set price range $50-$200
3. Apply filters

**Expected Result:**
- Only Electronics products priced $50-$200 displayed

**Traceability:**
- AC: [AC-PROD-001](../l1/acceptance-criteria.md#ac-prod-001)
- BR: [BR-PROD](../l1/business-rules.md#br-prod)

---

### TC-AC-PROD-001-P03 – Sort products by price ascending {#tc-ac-prod-001-p03}

**Preconditions:**
- Multiple products exist
- User is on product listing page

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| sortBy | price_asc | Sort by price low to high |

**Steps:**
1. Navigate to product listing
2. Select sort option 'Price: Low to High'

**Expected Result:**
- Products displayed in ascending price order

**Traceability:**
- AC: [AC-PROD-001](../l1/acceptance-criteria.md#ac-prod-001)
- BR: [BR-PROD](../l1/business-rules.md#br-prod)

---

### TC-AC-PROD-001-P04 – Display 20 products per page {#tc-ac-prod-001-p04}

**Preconditions:**
- More than 20 products exist
- User is on product listing page

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| totalProducts | 50 | Enough for pagination |

**Steps:**
1. Navigate to product listing
2. Count displayed products

**Expected Result:**
- Exactly 20 products displayed on first page

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
| productId | PROD-001 | Out of stock product |
| inventory | 0 | Zero available |

**Steps:**
1. Navigate to product listing
2. Locate out-of-stock product

**Expected Result:**
- 'Out of Stock' indicator visible on product card

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
- 'Out of Stock' indicator prominently displayed

**Traceability:**
- AC: [AC-PROD-002](../l1/acceptance-criteria.md#ac-prod-002)
- BR: [BR-PROD](../l1/business-rules.md#br-prod)

---

### TC-AC-PROD-002-P03 – Disable add to cart button for out of stock product {#tc-ac-prod-002-p03}

**Preconditions:**
- Product has zero inventory
- User on product detail page

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| inventory | 0 | No stock available |

**Steps:**
1. View out-of-stock product detail page
2. Check add-to-cart button state

**Expected Result:**
- Add to cart button is disabled/greyed out
- Button not clickable

**Traceability:**
- AC: [AC-PROD-002](../l1/acceptance-criteria.md#ac-prod-002)
- BR: [BR-PROD](../l1/business-rules.md#br-prod)

---

### TC-AC-CART-001-P01 – Add single quantity of in-stock product to cart {#tc-ac-cart-001-p01}

**Preconditions:**
- Product in stock
- User viewing product detail

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| productId | PROD-001 | In stock product |
| quantity | 1 | Single unit |
| availableStock | 10 | Sufficient stock |

**Steps:**
1. View in-stock product
2. Set quantity to 1
3. Click Add to Cart

**Expected Result:**
- Item added to cart
- Toast notification appears with View Cart option

**Traceability:**
- AC: [AC-CART-001](../l1/acceptance-criteria.md#ac-cart-001)
- BR: [BR-CART](../l1/business-rules.md#br-cart)

---

### TC-AC-CART-001-P02 – Add product with selected variant to cart {#tc-ac-cart-001-p02}

**Preconditions:**
- Product has variants
- Product in stock

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| productId | PROD-002 | Product with variants |
| variant | Size: Large, Color: Blue | Selected variant |
| quantity | 2 | Multiple units |

**Steps:**
1. View product with variants
2. Select variant
3. Add to cart

**Expected Result:**
- Item with variant added
- Cart icon shows updated count

**Traceability:**
- AC: [AC-CART-001](../l1/acceptance-criteria.md#ac-cart-001)
- BR: [BR-CART](../l1/business-rules.md#br-cart)

---

### TC-AC-CART-001-P03 – Cart icon updates with new item count {#tc-ac-cart-001-p03}

**Preconditions:**
- Cart is empty
- Product in stock

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| initialCartCount | 0 | Empty cart |
| quantity | 3 | Adding 3 items |

**Steps:**
1. Verify cart icon shows 0
2. Add 3 items to cart

**Expected Result:**
- Cart icon updates to show 3

**Traceability:**
- AC: [AC-CART-001](../l1/acceptance-criteria.md#ac-cart-001)
- BR: [BR-CART](../l1/business-rules.md#br-cart)

---

### TC-AC-CART-002-P01 – Increment quantity when adding same product again {#tc-ac-cart-002-p01}

**Preconditions:**
- Product A (qty 2) already in cart
- Product A in stock

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| productId | PROD-001 | Same product |
| existingQty | 2 | Already in cart |
| addQty | 1 | Adding 1 more |

**Steps:**
1. View cart with Product A qty 2
2. Add Product A with qty 1 again

**Expected Result:**
- Cart shows Product A with quantity 3

**Traceability:**
- AC: [AC-CART-002](../l1/acceptance-criteria.md#ac-cart-002)
- BR: [BR-CART](../l1/business-rules.md#br-cart)

---

### TC-AC-CART-002-P02 – Single line item maintained not duplicate created {#tc-ac-cart-002-p02}

**Preconditions:**
- Product in cart
- Same product added again

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| productId | PROD-001 | Same product |
| existingQty | 2 | In cart |
| addQty | 3 | Adding more |

**Steps:**
1. Add same product multiple times
2. View cart

**Expected Result:**
- Only one line item for product
- Combined quantity is 5

**Traceability:**
- AC: [AC-CART-002](../l1/acceptance-criteria.md#ac-cart-002)
- BR: [BR-CART](../l1/business-rules.md#br-cart)

---

### TC-AC-CART-002-P03 – Increment works for products with same variant {#tc-ac-cart-002-p03}

**Preconditions:**
- Product with variant in cart

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| productId | PROD-002 | Product with variant |
| variant | Size: Large | Same variant |
| existingQty | 1 | In cart |

**Steps:**
1. Add same product same variant again

**Expected Result:**
- Quantity incremented
- Single line item maintained

**Traceability:**
- AC: [AC-CART-002](../l1/acceptance-criteria.md#ac-cart-002)
- BR: [BR-CART](../l1/business-rules.md#br-cart)

---

### TC-AC-CART-003-P01 – Session-based cart created for guest user {#tc-ac-cart-003-p01}

**Preconditions:**
- User not logged in
- No existing session

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| userType | guest | Not authenticated |
| productId | PROD-001 | Product to add |

**Steps:**
1. Browse as guest user
2. Add item to cart

**Expected Result:**
- Session-based cart created
- Item visible in cart

**Traceability:**
- AC: [AC-CART-003](../l1/acceptance-criteria.md#ac-cart-003)
- BR: [BR-CART](../l1/business-rules.md#br-cart)

---

### TC-AC-CART-003-P02 – Guest cart persists across browser sessions within 7 days {#tc-ac-cart-003-p02}

**Preconditions:**
- Guest has items in cart
- Less than 7 days elapsed

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
- Session recognized

**Traceability:**
- AC: [AC-CART-003](../l1/acceptance-criteria.md#ac-cart-003)
- BR: [BR-CART](../l1/business-rules.md#br-cart)

---

### TC-AC-CART-003-P03 – Guest can add multiple items to session cart {#tc-ac-cart-003-p03}

**Preconditions:**
- User browsing as guest

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| products | PROD-001, PROD-002, PROD-003 | Multiple items |

**Steps:**
1. Add Product 1
2. Add Product 2
3. Add Product 3
4. View cart

**Expected Result:**
- All three products in cart
- Cart functions normally

**Traceability:**
- AC: [AC-CART-003](../l1/acceptance-criteria.md#ac-cart-003)
- BR: [BR-CART](../l1/business-rules.md#br-cart)

---

### TC-AC-CART-004-P01 – Merge guest cart with account cart on login {#tc-ac-cart-004-p01}

**Preconditions:**
- Guest user has 2 items in session cart
- Same user has registered account with 1 different item

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| guestCartItems | [{sku: A, qty: 2}, {sku: B, qty: 1}] | Guest session items |
| accountCartItems | [{sku: C, qty: 3}] | Existing account items |

**Steps:**
1. Add items A and B to guest cart
2. Log in with registered account containing item C
3. View merged cart

**Expected Result:**
- Cart contains all 3 items: A(2), B(1), C(3)
- Cart total reflects all items

**Traceability:**
- AC: [AC-CART-004](../l1/acceptance-criteria.md#ac-cart-004)
- BR: [AMB-ENT-013](../l1/business-rules.md#amb-ent-013)

---

### TC-AC-CART-004-P02 – Combine quantities for duplicate products on merge {#tc-ac-cart-004-p02}

**Preconditions:**
- Guest cart has product X with quantity 2
- Account cart has product X with quantity 3

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| guestProduct | {sku: X, qty: 2} | Same product in guest cart |
| accountProduct | {sku: X, qty: 3} | Same product in account cart |

**Steps:**
1. Add product X (qty 2) to guest cart
2. Log in with account containing product X (qty 3)
3. View merged cart

**Expected Result:**
- Product X shows combined quantity of 5
- Only one line item for product X

**Traceability:**
- AC: [AC-CART-004](../l1/acceptance-criteria.md#ac-cart-004)
- BR: [AMB-ENT-013](../l1/business-rules.md#amb-ent-013)

---

### TC-AC-CART-004-P03 – Merge guest cart into empty account cart {#tc-ac-cart-004-p03}

**Preconditions:**
- Guest cart has 3 items
- Account cart is empty

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| guestCartItems | [{sku: A}, {sku: B}, {sku: C}] | Guest items |
| accountCartItems | [] | Empty account cart |

**Steps:**
1. Add 3 items to guest cart
2. Log in with account having empty cart
3. View merged cart

**Expected Result:**
- All 3 guest items appear in account cart
- Guest session cart is cleared

**Traceability:**
- AC: [AC-CART-004](../l1/acceptance-criteria.md#ac-cart-004)
- BR: [AMB-ENT-013](../l1/business-rules.md#amb-ent-013)

---

### TC-AC-CART-005-P01 – Update cart item quantity successfully {#tc-ac-cart-005-p01}

**Preconditions:**
- Customer has item X with quantity 2 in cart
- Product X has sufficient stock

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| currentQty | 2 | Starting quantity |
| newQty | 5 | Target quantity |
| unitPrice | 10.00 | Price per unit |

**Steps:**
1. View cart with item X (qty 2)
2. Change quantity to 5
3. Submit update

**Expected Result:**
- Quantity updated to 5
- Cart total recalculated to 50.00

**Traceability:**
- AC: [AC-CART-005](../l1/acceptance-criteria.md#ac-cart-005)
- BR: [AMB-OP-011](../l1/business-rules.md#amb-op-011), [AMB-OP-012](../l1/business-rules.md#amb-op-012), [AMB-OP-013](../l1/business-rules.md#amb-op-013)

---

### TC-AC-CART-005-P02 – Cart total recalculates on quantity decrease {#tc-ac-cart-005-p02}

**Preconditions:**
- Customer has item with quantity 10

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| currentQty | 10 | Starting quantity |
| newQty | 3 | Decreased quantity |

**Steps:**
1. View cart with item (qty 10)
2. Decrease quantity to 3
3. Submit update

**Expected Result:**
- Quantity updated to 3
- Cart total reflects reduced amount

**Traceability:**
- AC: [AC-CART-005](../l1/acceptance-criteria.md#ac-cart-005)
- BR: [AMB-OP-011](../l1/business-rules.md#amb-op-011)

---

### TC-AC-CART-006-P01 – Remove item immediately without confirmation {#tc-ac-cart-006-p01}

**Preconditions:**
- Customer has item in cart

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| itemToRemove | SKU-123 | Target item |

**Steps:**
1. View cart with item
2. Click remove on item
3. Observe result

**Expected Result:**
- Item removed immediately
- No confirmation dialog shown

**Traceability:**
- AC: [AC-CART-006](../l1/acceptance-criteria.md#ac-cart-006)
- BR: [AMB-OP-014](../l1/business-rules.md#amb-op-014), [AMB-OP-015](../l1/business-rules.md#amb-op-015)

---

### TC-AC-CART-006-P02 – Show Undo option after item removal {#tc-ac-cart-006-p02}

**Preconditions:**
- Customer has item in cart

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| itemToRemove | SKU-456 | Item being removed |

**Steps:**
1. View cart with item
2. Click remove on item
3. Observe Undo option

**Expected Result:**
- Item removed from cart
- Undo option displayed briefly

**Traceability:**
- AC: [AC-CART-006](../l1/acceptance-criteria.md#ac-cart-006)
- BR: [AMB-OP-014](../l1/business-rules.md#amb-op-014)

---

### TC-AC-CART-006-P03 – Undo restores removed item {#tc-ac-cart-006-p03}

**Preconditions:**
- Customer just removed an item
- Undo option is visible

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| removedItem | {sku: X, qty: 3} | Item to restore |

**Steps:**
1. Remove item from cart
2. Click Undo while option visible
3. Verify cart state

**Expected Result:**
- Item restored to cart with original quantity
- Cart total recalculated

**Traceability:**
- AC: [AC-CART-006](../l1/acceptance-criteria.md#ac-cart-006)
- BR: [AMB-OP-014](../l1/business-rules.md#amb-op-014)

---

### TC-AC-CART-007-P01 – Out of stock item remains in cart with warning {#tc-ac-cart-007-p01}

**Preconditions:**
- Customer has item in cart
- Product inventory drops to zero

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| initialStock | 5 | Stock before |
| finalStock | 0 | Stock after others purchase |

**Steps:**
1. Add item to cart
2. Another user purchases all remaining stock
3. Refresh cart page

**Expected Result:**
- Item still visible in cart
- Warning indicator displayed on item

**Traceability:**
- AC: [AC-CART-007](../l1/acceptance-criteria.md#ac-cart-007)
- BR: [AMB-ENT-014](../l1/business-rules.md#amb-ent-014)

---

### TC-AC-CART-007-P02 – Checkout blocked for out of stock item {#tc-ac-cart-007-p02}

**Preconditions:**
- Cart contains out of stock item

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| itemStock | 0 | Out of stock |

**Steps:**
1. View cart with out of stock item
2. Attempt to proceed to checkout

**Expected Result:**
- Checkout blocked or item excluded
- Clear message about stock issue

**Traceability:**
- AC: [AC-CART-007](../l1/acceptance-criteria.md#ac-cart-007)
- BR: [AMB-ENT-014](../l1/business-rules.md#amb-ent-014)

---

### TC-AC-CART-007-P03 – Warning clears when stock replenished {#tc-ac-cart-007-p03}

**Preconditions:**
- Cart item was out of stock
- Stock is replenished

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| newStock | 10 | Stock replenished |

**Steps:**
1. Have out of stock item in cart
2. Stock gets replenished to 10
3. Refresh cart

**Expected Result:**
- Warning indicator removed
- Checkout now allowed for item

**Traceability:**
- AC: [AC-CART-007](../l1/acceptance-criteria.md#ac-cart-007)
- BR: [AMB-ENT-014](../l1/business-rules.md#amb-ent-014)

---

### TC-AC-CART-008-P01 – Cart displays current price after price increase {#tc-ac-cart-008-p01}

**Preconditions:**
- Customer has item in cart at $10
- Price increases to $12

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| originalPrice | 10.00 | Price when added |
| newPrice | 12.00 | Current price |

**Steps:**
1. Add item to cart at $10
2. Admin increases price to $12
3. View cart

**Expected Result:**
- Cart shows $12 per item
- Cart total reflects new price

**Traceability:**
- AC: [AC-CART-008](../l1/acceptance-criteria.md#ac-cart-008)
- BR: [AMB-ENT-016](../l1/business-rules.md#amb-ent-016)

---

### TC-AC-CART-008-P02 – Customer notified of price change {#tc-ac-cart-008-p02}

**Preconditions:**
- Item price changed since added to cart

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| originalPrice | 25.00 | Original |
| newPrice | 30.00 | Changed |

**Steps:**
1. Add item at $25
2. Price changes to $30
3. View cart

**Expected Result:**
- Notification shown about price change
- Original and new price visible

**Traceability:**
- AC: [AC-CART-008](../l1/acceptance-criteria.md#ac-cart-008)
- BR: [AMB-ENT-016](../l1/business-rules.md#amb-ent-016)

---

### TC-AC-CART-008-P03 – Cart displays current price after price decrease {#tc-ac-cart-008-p03}

**Preconditions:**
- Customer has item in cart at $50
- Price decreases to $40

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| originalPrice | 50.00 | Higher price |
| newPrice | 40.00 | Sale price |

**Steps:**
1. Add item to cart at $50
2. Price drops to $40
3. View cart

**Expected Result:**
- Cart shows $40 per item
- Customer benefits from lower price

**Traceability:**
- AC: [AC-CART-008](../l1/acceptance-criteria.md#ac-cart-008)
- BR: [AMB-ENT-016](../l1/business-rules.md#amb-ent-016)

---

### TC-AC-ORDER-001-P01 – Successfully place order with all valid prerequisites {#tc-ac-order-001-p01}

**Preconditions:**
- Customer is registered
- Email is verified
- Cart has items
- Valid shipping address
- Valid payment method

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| cartItems | 2 items | Multiple items in cart |
| paymentMethod | valid_card | Pre-authorized card |

**Steps:**
1. Login as verified customer
2. Add items to cart
3. Enter shipping address
4. Enter payment method
5. Submit order

**Expected Result:**
- Payment authorized
- Inventory decremented
- Order status is pending
- Cart cleared
- Confirmation page shown
- Email queued

**Traceability:**
- AC: [AC-ORDER-001](../l1/acceptance-criteria.md#ac-order-001)
- BR: [BR-ORDER](../l1/business-rules.md#br-order)

---

### TC-AC-ORDER-001-P02 – Order confirmation page displays correct order details {#tc-ac-order-001-p02}

**Preconditions:**
- Valid order submission prerequisites met

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| cartItems | 3 different products | Verify all items shown |
| totalAmount | 75.50 | Verify total displayed |

**Steps:**
1. Complete order submission
2. View confirmation page

**Expected Result:**
- Order number displayed
- All items listed
- Total matches cart
- Shipping address shown

**Traceability:**
- AC: [AC-ORDER-001](../l1/acceptance-criteria.md#ac-order-001)
- BR: [BR-ORDER](../l1/business-rules.md#br-order)

---

### TC-AC-ORDER-001-P03 – Cart is cleared after successful order placement {#tc-ac-order-001-p03}

**Preconditions:**
- Order successfully placed

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| cartItems | 5 items | Multiple items before order |

**Steps:**
1. Place order successfully
2. Navigate to cart page

**Expected Result:**
- Cart shows 0 items
- Cart total is $0.00

**Traceability:**
- AC: [AC-ORDER-001](../l1/acceptance-criteria.md#ac-order-001)
- BR: [BR-ORDER](../l1/business-rules.md#br-order)

---

### TC-AC-ORDER-002-P01 – Free shipping applied at exactly $50 subtotal {#tc-ac-order-002-p01}

**Preconditions:**
- Customer at checkout

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| subtotal | 50.00 | Exact threshold |

**Steps:**
1. Add items totaling $50
2. View shipping options

**Expected Result:**
- Free shipping automatically applied
- Shipping cost shows $0.00

**Traceability:**
- AC: [AC-ORDER-002](../l1/acceptance-criteria.md#ac-order-002)
- BR: [BR-SHIPPING](../l1/business-rules.md#br-shipping)

---

### TC-AC-ORDER-002-P02 – Free shipping applied above $50 subtotal {#tc-ac-order-002-p02}

**Preconditions:**
- Customer at checkout

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| subtotal | 75.00 | Above threshold |

**Steps:**
1. Add items totaling $75
2. View shipping options

**Expected Result:**
- Free shipping automatically applied

**Traceability:**
- AC: [AC-ORDER-002](../l1/acceptance-criteria.md#ac-order-002)
- BR: [BR-SHIPPING](../l1/business-rules.md#br-shipping)

---

### TC-AC-ORDER-002-P03 – Free shipping based on subtotal after discounts {#tc-ac-order-002-p03}

**Preconditions:**
- Discount code applied

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| originalSubtotal | 60.00 | Before discount |
| discount | 5.00 | Applied discount |
| finalSubtotal | 55.00 | After discount, still >= 50 |

**Steps:**
1. Add $60 of items
2. Apply $5 discount
3. View shipping

**Expected Result:**
- Free shipping applied based on $55 subtotal

**Traceability:**
- AC: [AC-ORDER-002](../l1/acceptance-criteria.md#ac-order-002)
- BR: [BR-SHIPPING](../l1/business-rules.md#br-shipping)

---

### TC-AC-ORDER-003-P01 – Calculate tax based on shipping state {#tc-ac-order-003-p01}

**Preconditions:**
- Customer has shipping address in taxable state

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| shippingState | CA | California has sales tax |
| subtotal | 100.00 | Easy calculation base |

**Steps:**
1. Enter California shipping address
2. View order total

**Expected Result:**
- California state tax applied to order

**Traceability:**
- AC: [AC-ORDER-003](../l1/acceptance-criteria.md#ac-order-003)
- BR: [BR-TAX](../l1/business-rules.md#br-tax)

---

### TC-AC-ORDER-003-P02 – Tax updates when shipping address changes {#tc-ac-order-003-p02}

**Preconditions:**
- Customer changing shipping address

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| initialState | CA | Higher tax state |
| newState | OR | No sales tax state |

**Steps:**
1. Set CA address
2. Note tax amount
3. Change to OR address

**Expected Result:**
- Tax recalculated for new state
- Oregon shows $0 tax

**Traceability:**
- AC: [AC-ORDER-003](../l1/acceptance-criteria.md#ac-order-003)
- BR: [BR-TAX](../l1/business-rules.md#br-tax)

---

### TC-AC-ORDER-003-P03 – Tax calculated for all US states with sales tax {#tc-ac-order-003-p03}

**Preconditions:**
- Customer shipping to state with sales tax

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| shippingState | NY | New York has sales tax |
| subtotal | 50.00 | Test amount |

**Steps:**
1. Enter NY shipping address
2. View tax calculation

**Expected Result:**
- NY state tax rate applied correctly

**Traceability:**
- AC: [AC-ORDER-003](../l1/acceptance-criteria.md#ac-order-003)
- BR: [BR-TAX](../l1/business-rules.md#br-tax)

---

### TC-AC-ORDER-004-P01 – Confirmation email queued after order placed {#tc-ac-order-004-p01}

**Preconditions:**
- Order successfully placed

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| orderNumber | ORD-12345 | New order |

**Steps:**
1. Place order successfully
2. Check email queue

**Expected Result:**
- Confirmation email added to queue
- Email not blocking order

**Traceability:**
- AC: [AC-ORDER-004](../l1/acceptance-criteria.md#ac-order-004)
- BR: [BR-EMAIL](../l1/business-rules.md#br-email)

---

### TC-AC-ORDER-004-P02 – Email contains all required order details {#tc-ac-order-004-p02}

**Preconditions:**
- Order placed
- Email sent

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| orderItems | 3 items | Multiple items |
| shippingAddress | 123 Main St | Full address |

**Steps:**
1. Place order
2. Receive confirmation email
3. Verify contents

**Expected Result:**
- Order number present
- All items listed
- Quantities shown
- Prices shown
- Total correct
- Shipping address included
- Estimated delivery shown

**Traceability:**
- AC: [AC-ORDER-004](../l1/acceptance-criteria.md#ac-order-004)
- BR: [BR-EMAIL](../l1/business-rules.md#br-email)

---

### TC-AC-ORDER-004-P03 – Email includes estimated delivery date {#tc-ac-order-004-p03}

**Preconditions:**
- Order placed with standard shipping

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| shippingMethod | standard | Has delivery estimate |

**Steps:**
1. Place order
2. Receive email
3. Check delivery estimate

**Expected Result:**
- Estimated delivery date included in email

**Traceability:**
- AC: [AC-ORDER-004](../l1/acceptance-criteria.md#ac-order-004)
- BR: [BR-EMAIL](../l1/business-rules.md#br-email)

---

### TC-AC-ORDER-005-P01 – Display order history with pagination {#tc-ac-order-005-p01}

**Preconditions:**
- Customer has multiple orders

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| totalOrders | 25 | Multiple pages |

**Steps:**
1. Login as customer
2. Navigate to order history

**Expected Result:**
- Orders displayed with pagination
- Can navigate between pages

**Traceability:**
- AC: [AC-ORDER-005](../l1/acceptance-criteria.md#ac-order-005)
- BR: [BR-ORDER-HISTORY](../l1/business-rules.md#br-order-history)

---

### TC-AC-ORDER-005-P02 – Order details show all required fields {#tc-ac-order-005-p02}

**Preconditions:**
- Customer has placed orders

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| orderStatus | shipped | Has tracking |

**Steps:**
1. View order history
2. Check order details

**Expected Result:**
- Order number shown
- Date shown
- Items listed
- Quantities shown
- Prices shown
- Status shown
- Shipping address shown

**Traceability:**
- AC: [AC-ORDER-005](../l1/acceptance-criteria.md#ac-order-005)
- BR: [BR-ORDER-HISTORY](../l1/business-rules.md#br-order-history)

---

### TC-AC-ORDER-005-P03 – Tracking number shown when available {#tc-ac-order-005-p03}

**Preconditions:**
- Order has been shipped

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| trackingNumber | 1Z999AA10123456784 | UPS tracking |

**Steps:**
1. View shipped order details

**Expected Result:**
- Tracking number displayed
- Tracking link functional

**Traceability:**
- AC: [AC-ORDER-005](../l1/acceptance-criteria.md#ac-order-005)
- BR: [BR-ORDER-HISTORY](../l1/business-rules.md#br-order-history)

---

### TC-AC-ORDER-006-P01 – Cancel pending order successfully {#tc-ac-order-006-p01}

**Preconditions:**
- Customer has existing order
- Order status is 'pending'

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| orderId | ORD-12345 | Valid pending order |
| orderStatus | pending | Cancellable status |

**Steps:**
1. Navigate to order details
2. Click cancel order button
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
| orderStatus | confirmed | Cancellable status |

**Steps:**
1. Navigate to order details
2. Click cancel order button
3. Confirm cancellation

**Expected Result:**
- Order status changes to 'cancelled'
- Success message displayed

**Traceability:**
- AC: [AC-ORDER-006](../l1/acceptance-criteria.md#ac-order-006)
- BR: [AMB-OP-028](../l1/business-rules.md#amb-op-028)

---

### TC-AC-ORDER-006-P03 – Inventory restored after cancellation {#tc-ac-order-006-p03}

**Preconditions:**
- Order with 5 units of product X
- Product X has 10 units in stock

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| productId | PROD-001 | Product in order |
| quantity | 5 | Units to restore |

**Steps:**
1. Cancel the order
2. Check inventory for product X

**Expected Result:**
- Product X inventory is now 15 units

**Traceability:**
- AC: [AC-ORDER-006](../l1/acceptance-criteria.md#ac-order-006)
- BR: [AMB-OP-029](../l1/business-rules.md#amb-op-029)

---

### TC-AC-ORDER-006-P04 – Refund initiated to original payment method {#tc-ac-order-006-p04}

**Preconditions:**
- Order paid with Visa ending 4242

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| paymentMethod | Visa *4242 | Original payment |
| orderTotal | 99.99 | Refund amount |

**Steps:**
1. Cancel the order
2. Check refund status

**Expected Result:**
- Refund initiated to Visa ending 4242
- Refund amount equals order total

**Traceability:**
- AC: [AC-ORDER-006](../l1/acceptance-criteria.md#ac-order-006)
- BR: [AMB-OP-030](../l1/business-rules.md#amb-op-030)

---

### TC-AC-ORDER-006-P05 – Cancellation confirmation email sent {#tc-ac-order-006-p05}

**Preconditions:**
- Customer has valid email on file

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| email | customer@test.com | Customer email |

**Steps:**
1. Cancel the order
2. Check email inbox

**Expected Result:**
- Cancellation confirmation email received

**Traceability:**
- AC: [AC-ORDER-006](../l1/acceptance-criteria.md#ac-order-006)
- BR: [AMB-OP-031](../l1/business-rules.md#amb-op-031)

---

### TC-AC-ORDER-007-P01 – Item prices captured at order creation {#tc-ac-order-007-p01}

**Preconditions:**
- Product A priced at $50.00
- Customer creates order

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| productPrice | 50.00 | Price at order time |

**Steps:**
1. Add Product A to cart
2. Place order
3. View order details

**Expected Result:**
- Order shows item price as $50.00

**Traceability:**
- AC: [AC-ORDER-007](../l1/acceptance-criteria.md#ac-order-007)
- BR: [AMB-ENT-019](../l1/business-rules.md#amb-ent-019)

---

### TC-AC-ORDER-007-P02 – Order prices unchanged after product price increase {#tc-ac-order-007-p02}

**Preconditions:**
- Order placed with Product A at $50.00

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| originalPrice | 50.00 | Price when ordered |
| newPrice | 75.00 | Updated catalog price |

**Steps:**
1. Update Product A price to $75.00
2. View existing order

**Expected Result:**
- Order still shows $50.00 for Product A

**Traceability:**
- AC: [AC-ORDER-007](../l1/acceptance-criteria.md#ac-order-007)
- BR: [AMB-ENT-019](../l1/business-rules.md#amb-ent-019)

---

### TC-AC-ORDER-007-P03 – Order prices unchanged after product price decrease {#tc-ac-order-007-p03}

**Preconditions:**
- Order placed with Product A at $50.00

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| originalPrice | 50.00 | Price when ordered |
| newPrice | 25.00 | Reduced catalog price |

**Steps:**
1. Update Product A price to $25.00
2. View existing order

**Expected Result:**
- Order still shows $50.00 for Product A

**Traceability:**
- AC: [AC-ORDER-007](../l1/acceptance-criteria.md#ac-order-007)
- BR: [AMB-ENT-019](../l1/business-rules.md#amb-ent-019)

---

### TC-AC-CUST-001-P01 – Successful registration with valid credentials {#tc-ac-cust-001-p01}

**Preconditions:**
- User is not registered

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| email | newuser@test.com | Valid email |
| password | SecurePass1 | Meets requirements |
| firstName | John | Required field |
| lastName | Doe | Required field |

**Steps:**
1. Navigate to registration
2. Enter all fields
3. Submit form

**Expected Result:**
- Account created successfully

**Traceability:**
- AC: [AC-CUST-001](../l1/acceptance-criteria.md#ac-cust-001)
- BR: [AMB-ENT-024](../l1/business-rules.md#amb-ent-024)

---

### TC-AC-CUST-001-P02 – Verification email sent after registration {#tc-ac-cust-001-p02}

**Preconditions:**
- User completes registration

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| email | newuser@test.com | Registration email |

**Steps:**
1. Complete registration
2. Check email inbox

**Expected Result:**
- Verification email received

**Traceability:**
- AC: [AC-CUST-001](../l1/acceptance-criteria.md#ac-cust-001)
- BR: [AMB-ENT-025](../l1/business-rules.md#amb-ent-025)

---

### TC-AC-CUST-001-P03 – User informed to verify email before ordering {#tc-ac-cust-001-p03}

**Preconditions:**
- Registration completed

**Steps:**
1. Complete registration
2. View confirmation message

**Expected Result:**
- Message states email verification required for orders

**Traceability:**
- AC: [AC-CUST-001](../l1/acceptance-criteria.md#ac-cust-001)
- BR: [AMB-ENT-026](../l1/business-rules.md#amb-ent-026)

---

### TC-AC-CUST-002-P01 – Add shipping address with all required fields {#tc-ac-cust-002-p01}

**Preconditions:**
- Customer is logged in

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| street | 123 Main St | Valid street |
| city | New York | Valid city |
| state | NY | Valid state |
| postalCode | 10001 | Valid postal code |
| country | USA | Valid country |
| recipientName | John Doe | Recipient |

**Steps:**
1. Navigate to address management
2. Enter address details
3. Save address

**Expected Result:**
- Address saved successfully

**Traceability:**
- AC: [AC-CUST-002](../l1/acceptance-criteria.md#ac-cust-002)
- BR: [AMB-ENT-027](../l1/business-rules.md#amb-ent-027)

---

### TC-AC-CUST-002-P02 – Set address as default {#tc-ac-cust-002-p02}

**Preconditions:**
- Customer has saved address

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| isDefault | true | Set as default |

**Steps:**
1. Add new address
2. Check 'Set as default' option
3. Save

**Expected Result:**
- Address saved as default
- Previous default unmarked

**Traceability:**
- AC: [AC-CUST-002](../l1/acceptance-criteria.md#ac-cust-002)
- BR: [AMB-ENT-031](../l1/business-rules.md#amb-ent-031)

---

### TC-AC-CUST-002-P03 – Add multiple addresses to account {#tc-ac-cust-002-p03}

**Preconditions:**
- Customer already has one address

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| addressCount | 2 | Second address |

**Steps:**
1. Add second address
2. Save

**Expected Result:**
- Both addresses appear in address list

**Traceability:**
- AC: [AC-CUST-002](../l1/acceptance-criteria.md#ac-cust-002)
- BR: [AMB-ENT-027](../l1/business-rules.md#amb-ent-027)

---

### TC-AC-CUST-003-P01 – Add Visa card via Stripe {#tc-ac-cust-003-p01}

**Preconditions:**
- Customer is logged in

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| cardType | Visa | Supported card |
| cardNumber | 4242424242424242 | Test Visa |

**Steps:**
1. Navigate to payment methods
2. Add Visa card
3. Save

**Expected Result:**
- Payment method tokenized and saved

**Traceability:**
- AC: [AC-CUST-003](../l1/acceptance-criteria.md#ac-cust-003)
- BR: [AMB-ENT-028](../l1/business-rules.md#amb-ent-028)

---

### TC-AC-CUST-003-P02 – Add PayPal payment method {#tc-ac-cust-003-p02}

**Preconditions:**
- Customer is logged in

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| paymentType | PayPal | Supported method |

**Steps:**
1. Select PayPal option
2. Authenticate with PayPal
3. Link account

**Expected Result:**
- PayPal account linked successfully

**Traceability:**
- AC: [AC-CUST-003](../l1/acceptance-criteria.md#ac-cust-003)
- BR: [AMB-ENT-028](../l1/business-rules.md#amb-ent-028)

---

### TC-AC-CUST-003-P03 – Set payment method as default {#tc-ac-cust-003-p03}

**Preconditions:**
- Customer has saved payment method

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| isDefault | true | Set as default |

**Steps:**
1. Add payment method
2. Set as default
3. Save

**Expected Result:**
- Payment method saved as default

**Traceability:**
- AC: [AC-CUST-003](../l1/acceptance-criteria.md#ac-cust-003)
- BR: [AMB-ENT-034](../l1/business-rules.md#amb-ent-034)

---

### TC-AC-CUST-004-P01 – Soft delete account when customer chooses deactivation {#tc-ac-cust-004-p01}

**Preconditions:**
- Customer is registered
- Customer has no pending orders

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| customerId | CUST-12345 | Valid customer |
| deletionType | soft_delete | Deactivation choice |

**Steps:**
1. Customer requests account deletion
2. Customer selects soft delete option
3. Customer confirms deletion request

**Expected Result:**
- Account status changed to deactivated
- Customer data preserved but inaccessible

**Traceability:**
- AC: [AC-CUST-004](../l1/acceptance-criteria.md#ac-cust-004)
- BR: [BR-CUST](../l1/business-rules.md#br-cust)

---

### TC-AC-CUST-004-P02 – Full data erasure when customer chooses GDPR deletion {#tc-ac-cust-004-p02}

**Preconditions:**
- Customer is registered
- Customer has no pending orders

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| customerId | CUST-67890 | Valid customer |
| deletionType | full_erasure | GDPR erasure choice |

**Steps:**
1. Customer requests account deletion
2. Customer selects full erasure option
3. Customer confirms deletion request

**Expected Result:**
- All personal data permanently deleted
- Anonymized order history retained

**Traceability:**
- AC: [AC-CUST-004](../l1/acceptance-criteria.md#ac-cust-004)
- BR: [BR-CUST](../l1/business-rules.md#br-cust), [BR-GDPR](../l1/business-rules.md#br-gdpr)

---

### TC-AC-CUST-004-P03 – Deletion proceeds after all pending orders completed {#tc-ac-cust-004-p03}

**Preconditions:**
- Customer had pending orders
- All orders now completed

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| customerId | CUST-11111 | Customer with completed orders |
| orderStatus | delivered | All orders completed |

**Steps:**
1. Verify all orders are completed
2. Customer requests deletion
3. Customer confirms deletion

**Expected Result:**
- Deletion request processed successfully

**Traceability:**
- AC: [AC-CUST-004](../l1/acceptance-criteria.md#ac-cust-004)
- BR: [BR-CUST](../l1/business-rules.md#br-cust)

---

### TC-AC-ADMIN-001-P01 – Create product with all valid required fields {#tc-ac-admin-001-p01}

**Preconditions:**
- Admin is logged in
- Admin has product management access

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| name | Test Product | Valid 12-char name |
| description | Product description text | Valid description |
| price | 29.99 | Valid price above minimum |
| category | Electronics | Valid category |

**Steps:**
1. Navigate to product management
2. Fill in all required fields
3. Submit new product

**Expected Result:**
- Product created with UUID
- Status is 'draft'
- Creator and timestamp logged

**Traceability:**
- AC: [AC-ADMIN-001](../l1/acceptance-criteria.md#ac-admin-001)
- BR: [BR-PROD](../l1/business-rules.md#br-prod)

---

### TC-AC-ADMIN-001-P02 – Create product with multiple images and primary marked {#tc-ac-admin-001-p02}

**Preconditions:**
- Admin is logged in
- Valid product data prepared

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| images | [img1.jpg, img2.jpg, img3.jpg] | 3 valid images |
| primaryImage | img1.jpg | First image as primary |

**Steps:**
1. Create product with valid data
2. Upload 3 images
3. Mark first as primary
4. Submit

**Expected Result:**
- Product created with 3 images
- Primary image correctly marked

**Traceability:**
- AC: [AC-ADMIN-001](../l1/acceptance-criteria.md#ac-admin-001)
- BR: [BR-PROD](../l1/business-rules.md#br-prod)

---

### TC-AC-ADMIN-001-P03 – Create product with duplicate name shows warning {#tc-ac-admin-001-p03}

**Preconditions:**
- Product with same name already exists

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| name | Existing Product | Name already in use |

**Steps:**
1. Enter duplicate product name
2. Submit product

**Expected Result:**
- Warning displayed about duplicate name
- Product creation allowed

**Traceability:**
- AC: [AC-ADMIN-001](../l1/acceptance-criteria.md#ac-admin-001)
- BR: [BR-PROD](../l1/business-rules.md#br-prod)

---

### TC-AC-ADMIN-002-P01 – Update product attributes successfully {#tc-ac-admin-002-p01}

**Preconditions:**
- Product exists
- Admin is logged in

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| productId | PROD-001 | Existing product |
| newName | Updated Product Name | Modified name |

**Steps:**
1. Open product for editing
2. Modify product name
3. Save changes

**Expected Result:**
- Product updated
- Last modified by and date recorded

**Traceability:**
- AC: [AC-ADMIN-002](../l1/acceptance-criteria.md#ac-admin-002)
- BR: [BR-PROD](../l1/business-rules.md#br-prod)

---

### TC-AC-ADMIN-002-P02 – Price change logged in audit trail {#tc-ac-admin-002-p02}

**Preconditions:**
- Product exists with price $29.99

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| oldPrice | 29.99 | Original price |
| newPrice | 34.99 | Updated price |

**Steps:**
1. Edit product price from $29.99 to $34.99
2. Save changes
3. Check audit log

**Expected Result:**
- Price updated
- Audit trail shows old and new price with timestamp

**Traceability:**
- AC: [AC-ADMIN-002](../l1/acceptance-criteria.md#ac-admin-002)
- BR: [BR-PROD](../l1/business-rules.md#br-prod), [BR-AUDIT](../l1/business-rules.md#br-audit)

---

### TC-AC-ADMIN-002-P03 – Cart prices reflect product price change {#tc-ac-admin-002-p03}

**Preconditions:**
- Product in customer cart
- Product price changed

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| productId | PROD-002 | Product in cart |
| newPrice | 24.99 | Updated price |

**Steps:**
1. Change product price
2. View cart containing this product

**Expected Result:**
- Cart shows updated price
- Cart total recalculated

**Traceability:**
- AC: [AC-ADMIN-002](../l1/acceptance-criteria.md#ac-admin-002)
- BR: [BR-PROD](../l1/business-rules.md#br-prod), [BR-CART](../l1/business-rules.md#br-cart)

---

### TC-AC-ADMIN-003-P01 – Soft delete product preserves order history {#tc-ac-admin-003-p01}

**Preconditions:**
- Product exists
- Product was in past orders

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| productId | PROD-005 | Product with order history |

**Steps:**
1. Click delete on product
2. Confirm deletion
3. View historical order

**Expected Result:**
- Product soft-deleted
- Order history still shows product details

**Traceability:**
- AC: [AC-ADMIN-003](../l1/acceptance-criteria.md#ac-admin-003)
- BR: [BR-PROD](../l1/business-rules.md#br-prod)

---

### TC-AC-ADMIN-003-P02 – Deleted product removed from active carts with notification {#tc-ac-admin-003-p02}

**Preconditions:**
- Product in customer cart
- Product deleted by admin

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| productId | PROD-006 | Product in cart |
| cartId | CART-001 | Customer cart with product |

**Steps:**
1. Delete product
2. Customer views their cart

**Expected Result:**
- Product removed from cart
- Notification shown to customer

**Traceability:**
- AC: [AC-ADMIN-003](../l1/acceptance-criteria.md#ac-admin-003)
- BR: [BR-PROD](../l1/business-rules.md#br-prod), [BR-CART](../l1/business-rules.md#br-cart)

---

### TC-AC-ADMIN-003-P03 – Deleted product no longer in listings {#tc-ac-admin-003-p03}

**Preconditions:**
- Product was visible in catalog

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| productId | PROD-007 | Previously visible product |

**Steps:**
1. Delete product
2. Search for product in catalog
3. Browse category

**Expected Result:**
- Product not found in search
- Product not in category listing

**Traceability:**
- AC: [AC-ADMIN-003](../l1/acceptance-criteria.md#ac-admin-003)
- BR: [BR-PROD](../l1/business-rules.md#br-prod)

---

### TC-AC-ADMIN-004-P01 – Display orders with all required columns {#tc-ac-admin-004-p01}

**Preconditions:**
- Orders exist in system
- Admin is logged in

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| orderCount | 50 | Multiple orders to display |

**Steps:**
1. Navigate to order management
2. View order list

**Expected Result:**
- Orders displayed with order number
- Date, customer, total, status shown

**Traceability:**
- AC: [AC-ADMIN-004](../l1/acceptance-criteria.md#ac-admin-004)
- BR: [BR-ORDER](../l1/business-rules.md#br-order)

---

### TC-AC-ADMIN-004-P02 – Filter orders by status {#tc-ac-admin-004-p02}

**Preconditions:**
- Orders exist with various statuses

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| statusFilter | shipped | Filter for shipped orders |

**Steps:**
1. View order list
2. Apply status filter for 'shipped'

**Expected Result:**
- Only shipped orders displayed

**Traceability:**
- AC: [AC-ADMIN-004](../l1/acceptance-criteria.md#ac-admin-004)
- BR: [BR-ORDER](../l1/business-rules.md#br-order)

---

### TC-AC-ADMIN-004-P03 – Filter orders by date range {#tc-ac-admin-004-p03}

**Preconditions:**
- Orders exist across multiple dates

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| startDate | 2024-01-01 | Range start |
| endDate | 2024-01-31 | Range end |

**Steps:**
1. View order list
2. Set date range filter
3. Apply filter

**Expected Result:**
- Only orders within date range displayed

**Traceability:**
- AC: [AC-ADMIN-004](../l1/acceptance-criteria.md#ac-admin-004)
- BR: [BR-ORDER](../l1/business-rules.md#br-order)

---

### TC-AC-ADMIN-004-P04 – Search orders by order number {#tc-ac-admin-004-p04}

**Preconditions:**
- Order with known number exists

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| searchTerm | ORD-2024-00123 | Specific order number |

**Steps:**
1. Enter order number in search
2. Execute search

**Expected Result:**
- Matching order displayed

**Traceability:**
- AC: [AC-ADMIN-004](../l1/acceptance-criteria.md#ac-admin-004)
- BR: [BR-ORDER](../l1/business-rules.md#br-order)

---

### TC-AC-ADMIN-005-P01 – Update order status from pending to confirmed {#tc-ac-admin-005-p01}

**Preconditions:**
- Admin is logged in
- Order exists with pending status

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| orderId | ORD-001 | Valid order ID |
| newStatus | confirmed | Valid next status |

**Steps:**
1. Navigate to order details
2. Change status to confirmed
3. Submit change

**Expected Result:**
- Status updated to confirmed
- Change logged with admin and timestamp
- Customer receives email notification

**Traceability:**
- AC: [AC-ADMIN-005](../l1/acceptance-criteria.md#ac-admin-005)
- BR: [AMB-ENT-018](../l1/business-rules.md#amb-ent-018), [AMB-OP-047](../l1/business-rules.md#amb-op-047)

---

### TC-AC-ADMIN-005-P02 – Update order status from confirmed to shipped {#tc-ac-admin-005-p02}

**Preconditions:**
- Admin is logged in
- Order exists with confirmed status

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| orderId | ORD-002 | Valid order ID |
| newStatus | shipped | Valid next status |

**Steps:**
1. Navigate to order details
2. Change status to shipped
3. Submit change

**Expected Result:**
- Status updated to shipped
- Change logged with admin and timestamp
- Customer receives email notification

**Traceability:**
- AC: [AC-ADMIN-005](../l1/acceptance-criteria.md#ac-admin-005)
- BR: [AMB-OP-048](../l1/business-rules.md#amb-op-048)

---

### TC-AC-ADMIN-005-P03 – Update order to shipped without tracking number {#tc-ac-admin-005-p03}

**Preconditions:**
- Admin is logged in
- Order exists with confirmed status

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| orderId | ORD-003 | Valid order ID |
| newStatus | shipped | Valid status |
| trackingNumber |  | Empty tracking |

**Steps:**
1. Navigate to order details
2. Change status to shipped
3. Leave tracking number empty
4. Submit change

**Expected Result:**
- Status updated to shipped
- No validation error for missing tracking

**Traceability:**
- AC: [AC-ADMIN-005](../l1/acceptance-criteria.md#ac-admin-005)
- BR: [AMB-OP-049](../l1/business-rules.md#amb-op-049)

---

### TC-AC-ADMIN-006-P01 – Set absolute inventory value successfully {#tc-ac-admin-006-p01}

**Preconditions:**
- Admin is logged in
- Product variant exists with stock 50

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| variantId | VAR-001 | Valid variant |
| newQuantity | 75 | Absolute value |
| reason | Recount | Optional reason |

**Steps:**
1. Navigate to inventory
2. Select variant
3. Set stock to 75
4. Submit

**Expected Result:**
- Stock updated to 75
- Audit trail shows before:50 after:75
- Reason recorded

**Traceability:**
- AC: [AC-ADMIN-006](../l1/acceptance-criteria.md#ac-admin-006)
- BR: [AMB-ENT-009](../l1/business-rules.md#amb-ent-009), [AMB-OP-051](../l1/business-rules.md#amb-op-051)

---

### TC-AC-ADMIN-006-P02 – Adjust inventory with delta value successfully {#tc-ac-admin-006-p02}

**Preconditions:**
- Admin is logged in
- Product variant exists with stock 50

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| variantId | VAR-002 | Valid variant |
| delta | +25 | Delta adjustment |

**Steps:**
1. Navigate to inventory
2. Select variant
3. Add 25 units
4. Submit

**Expected Result:**
- Stock updated to 75
- Audit trail shows before:50 after:75

**Traceability:**
- AC: [AC-ADMIN-006](../l1/acceptance-criteria.md#ac-admin-006)
- BR: [AMB-OP-052](../l1/business-rules.md#amb-op-052)

---

### TC-AC-ADMIN-006-P03 – Trigger low stock alert when threshold crossed {#tc-ac-admin-006-p03}

**Preconditions:**
- Product has low stock threshold of 10
- Current stock is 15

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| variantId | VAR-003 | Valid variant |
| delta | -10 | Reduces below threshold |

**Steps:**
1. Navigate to inventory
2. Reduce stock by 10
3. Submit

**Expected Result:**
- Stock updated to 5
- Low stock alert triggered

**Traceability:**
- AC: [AC-ADMIN-006](../l1/acceptance-criteria.md#ac-admin-006)
- BR: [AMB-ENT-042](../l1/business-rules.md#amb-ent-042), [AMB-OP-055](../l1/business-rules.md#amb-op-055)

---

### TC-AC-VAR-001-P01 – Display all variant options for product with size variants {#tc-ac-var-001-p01}

**Preconditions:**
- Product exists with size variants S, M, L, XL

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| productId | PROD-001 | Product with sizes |
| variants | S,M,L,XL | All size options |

**Steps:**
1. Navigate to product page
2. View variant selector

**Expected Result:**
- All four size options displayed
- Each variant selectable

**Traceability:**
- AC: [AC-VAR-001](../l1/acceptance-criteria.md#ac-var-001)
- BR: [AMB-ENT-008](../l1/business-rules.md#amb-ent-008)

---

### TC-AC-VAR-001-P02 – Display individual stock status per variant {#tc-ac-var-001-p02}

**Preconditions:**
- Product with variants having different stock levels

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| variantS | 10 in stock | In stock |
| variantM | 0 in stock | Out of stock |

**Steps:**
1. Navigate to product page
2. View each variant stock status

**Expected Result:**
- Size S shows in stock
- Size M shows out of stock

**Traceability:**
- AC: [AC-VAR-001](../l1/acceptance-criteria.md#ac-var-001)
- BR: [AMB-ENT-009](../l1/business-rules.md#amb-ent-009)

---

### TC-AC-VAR-001-P03 – Display price override for variant with different price {#tc-ac-var-001-p03}

**Preconditions:**
- Product with base price $50
- XL variant has $55 price override

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| basePrice | $50 | Base product price |
| xlPrice | $55 | XL override price |

**Steps:**
1. Navigate to product
2. Select XL variant

**Expected Result:**
- Price displayed as $55 for XL variant

**Traceability:**
- AC: [AC-VAR-001](../l1/acceptance-criteria.md#ac-var-001)
- BR: [AMB-ENT-010](../l1/business-rules.md#amb-ent-010)

---

### TC-AC-CAT-001-P01 – Display 3-level category hierarchy correctly {#tc-ac-cat-001-p01}

**Preconditions:**
- Categories: Electronics > Phones > Smartphones exist

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
- Hierarchy clearly displayed

**Traceability:**
- AC: [AC-CAT-001](../l1/acceptance-criteria.md#ac-cat-001)
- BR: [AMB-ENT-006](../l1/business-rules.md#amb-ent-006)

---

### TC-AC-CAT-001-P02 – Display product with primary and secondary categories {#tc-ac-cat-001-p02}

**Preconditions:**
- Product in primary: Phones, secondary: Accessories

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| primary | Phones | Primary category |
| secondary | Accessories | Secondary category |

**Steps:**
1. View product details
2. Check category assignment

**Expected Result:**
- Primary category shown
- Secondary category shown

**Traceability:**
- AC: [AC-CAT-001](../l1/acceptance-criteria.md#ac-cat-001)
- BR: [AMB-ENT-038](../l1/business-rules.md#amb-ent-038)

---

### TC-AC-CAT-001-P03 – Browse products in second-level category {#tc-ac-cat-001-p03}

**Preconditions:**
- Category Phones has products assigned

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| category | Phones | Second level category |

**Steps:**
1. Navigate to Phones category
2. View product listing

**Expected Result:**
- Products in Phones category displayed

**Traceability:**
- AC: [AC-CAT-001](../l1/acceptance-criteria.md#ac-cat-001)
- BR: [AMB-ENT-037](../l1/business-rules.md#amb-ent-037)

---

## Negative Tests (Error Cases)

### TC-AC-PROD-001-N01 – Show empty state when no products match filters {#tc-ac-prod-001-n01}

**Preconditions:**
- User is on product listing page

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| category | NonExistent | Category with no products |
| priceMin | 99999 | Price no product matches |

**Steps:**
1. Apply filters that match no products

**Expected Result:**
- 'No products found' message displayed
- Clear filters option shown

**Traceability:**
- AC: [AC-PROD-001](../l1/acceptance-criteria.md#ac-prod-001)
- BR: [BR-PROD](../l1/business-rules.md#br-prod)

---

### TC-AC-PROD-001-N02 – Reject invalid price range when min exceeds max {#tc-ac-prod-001-n02}

**Preconditions:**
- User is on product listing page

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| priceMin | 500 | Higher than max |
| priceMax | 100 | Lower than min |

**Steps:**
1. Set minimum price to $500
2. Set maximum price to $100
3. Try to apply filter

**Expected Result:**
- Validation error displayed
- Filter not applied

**Traceability:**
- AC: [AC-PROD-001](../l1/acceptance-criteria.md#ac-prod-001)
- BR: [BR-PROD](../l1/business-rules.md#br-prod)

---

### TC-AC-PROD-002-N01 – Cannot add out of stock product via direct API call {#tc-ac-prod-002-n01}

**Preconditions:**
- Product has zero inventory

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| productId | PROD-001 | Out of stock |
| quantity | 1 | Attempted quantity |

**Steps:**
1. Attempt to add out-of-stock product via API

**Expected Result:**
- API returns error
- Product not added to cart

**Traceability:**
- AC: [AC-PROD-002](../l1/acceptance-criteria.md#ac-prod-002)
- BR: [BR-PROD](../l1/business-rules.md#br-prod)

---

### TC-AC-PROD-002-N02 – Product becomes out of stock during session {#tc-ac-prod-002-n02}

**Preconditions:**
- Product initially in stock
- Another user purchases last item

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| productId | PROD-001 | Last item purchased by other |
| inventory | 0 | Now zero |

**Steps:**
1. View product page
2. Inventory drops to zero
3. Refresh page

**Expected Result:**
- Out of Stock indicator now shown
- Add to cart disabled

**Traceability:**
- AC: [AC-PROD-002](../l1/acceptance-criteria.md#ac-prod-002)
- BR: [BR-PROD](../l1/business-rules.md#br-prod)

---

### TC-AC-CART-001-N01 – Cannot add out of stock product to cart {#tc-ac-cart-001-n01}

**Preconditions:**
- Product has zero inventory

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| productId | PROD-003 | Out of stock |
| inventory | 0 | No stock |

**Steps:**
1. View out-of-stock product
2. Attempt to click Add to Cart

**Expected Result:**
- Button disabled
- 'Out of Stock' message shown

**Traceability:**
- AC: [AC-CART-001](../l1/acceptance-criteria.md#ac-cart-001)
- BR: [BR-CART](../l1/business-rules.md#br-cart)

---

### TC-AC-CART-001-N02 – Quantity exceeds stock limits to available {#tc-ac-cart-001-n02}

**Preconditions:**
- Product has limited stock

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| availableStock | 5 | Only 5 available |
| requestedQuantity | 10 | Exceeds stock |

**Steps:**
1. View product with 5 in stock
2. Try to add 10 to cart

**Expected Result:**
- Quantity limited to 5
- Notification shown about limit

**Traceability:**
- AC: [AC-CART-001](../l1/acceptance-criteria.md#ac-cart-001)
- BR: [BR-CART](../l1/business-rules.md#br-cart)

---

### TC-AC-CART-001-N03 – Cannot add product without selecting required variant {#tc-ac-cart-001-n03}

**Preconditions:**
- Product has required variants
- No variant selected

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| productId | PROD-002 | Has variants |
| variantSelected | false | Not selected |

**Steps:**
1. View product with variants
2. Attempt to add without variant selection

**Expected Result:**
- Add to cart button disabled until variant selected

**Traceability:**
- AC: [AC-CART-001](../l1/acceptance-criteria.md#ac-cart-001)
- BR: [BR-CART](../l1/business-rules.md#br-cart)

---

### TC-AC-CART-002-N01 – Combined quantity limited when exceeds stock {#tc-ac-cart-002-n01}

**Preconditions:**
- Product A (qty 3) in cart
- Only 5 total in stock

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| existingQty | 3 | Already in cart |
| addQty | 4 | Trying to add |
| availableStock | 5 | Total stock |

**Steps:**
1. Try to add 4 more when 3 already in cart and only 5 in stock

**Expected Result:**
- Quantity limited to 5
- Notification about limit shown

**Traceability:**
- AC: [AC-CART-002](../l1/acceptance-criteria.md#ac-cart-002)
- BR: [BR-CART](../l1/business-rules.md#br-cart)

---

### TC-AC-CART-002-N02 – Different variants create separate line items {#tc-ac-cart-002-n02}

**Preconditions:**
- Product with Size: Large in cart

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| productId | PROD-002 | Same product |
| existingVariant | Size: Large | In cart |
| newVariant | Size: Small | Different variant |

**Steps:**
1. Add same product with different variant

**Expected Result:**
- Two separate line items created
- Each variant tracked separately

**Traceability:**
- AC: [AC-CART-002](../l1/acceptance-criteria.md#ac-cart-002)
- BR: [BR-CART](../l1/business-rules.md#br-cart)

---

### TC-AC-CART-003-N01 – Cart cleared after 7 days of inactivity {#tc-ac-cart-003-n01}

**Preconditions:**
- Guest has items in cart
- 7+ days of inactivity

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| daysSinceActivity | 8 | Exceeds 7 day limit |

**Steps:**
1. Add items as guest
2. Wait 8 days without activity
3. Return to site

**Expected Result:**
- Cart is empty
- User sees empty cart state

**Traceability:**
- AC: [AC-CART-003](../l1/acceptance-criteria.md#ac-cart-003)
- BR: [BR-CART](../l1/business-rules.md#br-cart)

---

### TC-AC-CART-003-N02 – Expired session shows empty cart not error {#tc-ac-cart-003-n02}

**Preconditions:**
- Session has expired

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| sessionExpired | true | Past 7 days |

**Steps:**
1. Return after session expiry
2. Click cart icon

**Expected Result:**
- Empty cart displayed
- No error messages
- Can add new items

**Traceability:**
- AC: [AC-CART-003](../l1/acceptance-criteria.md#ac-cart-003)
- BR: [BR-CART](../l1/business-rules.md#br-cart)

---

### TC-AC-CART-004-N01 – Limit merged quantity to available stock {#tc-ac-cart-004-n01}

**Preconditions:**
- Product X has 5 units in stock
- Guest cart has X with qty 3, account cart has X with qty 4

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| stockLevel | 5 | Available inventory |
| guestQty | 3 | Guest cart quantity |
| accountQty | 4 | Account cart quantity |

**Steps:**
1. Set product X stock to 5
2. Add X (qty 3) to guest cart
3. Log in with account containing X (qty 4)
4. View merged cart

**Expected Result:**
- Product X quantity limited to 5 (available stock)
- Notification shown about quantity adjustment

**Traceability:**
- AC: [AC-CART-004](../l1/acceptance-criteria.md#ac-cart-004)
- BR: [AMB-ENT-013](../l1/business-rules.md#amb-ent-013)

---

### TC-AC-CART-004-N02 – Reject merge for out-of-stock product {#tc-ac-cart-004-n02}

**Preconditions:**
- Guest cart has product Y
- Product Y goes out of stock before login

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| productStock | 0 | Out of stock |
| guestQty | 2 | Guest had 2 items |

**Steps:**
1. Add product Y to guest cart
2. Product Y stock drops to 0
3. Log in with registered account
4. View merged cart

**Expected Result:**
- Product Y quantity set to 0 or marked unavailable
- Notification shown about stock issue

**Traceability:**
- AC: [AC-CART-004](../l1/acceptance-criteria.md#ac-cart-004)
- BR: [AMB-ENT-013](../l1/business-rules.md#amb-ent-013)

---

### TC-AC-CART-005-N01 – Remove item when quantity set to zero {#tc-ac-cart-005-n01}

**Preconditions:**
- Customer has item in cart

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| newQty | 0 | Zero quantity triggers removal |

**Steps:**
1. View cart with item
2. Set quantity to 0
3. Submit update

**Expected Result:**
- Item removed from cart
- Cart total recalculated without item

**Traceability:**
- AC: [AC-CART-005](../l1/acceptance-criteria.md#ac-cart-005)
- BR: [AMB-OP-011](../l1/business-rules.md#amb-op-011)

---

### TC-AC-CART-005-N02 – Limit quantity to available stock with notification {#tc-ac-cart-005-n02}

**Preconditions:**
- Product has 5 units in stock
- Customer tries to set quantity to 10

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| stockLevel | 5 | Available stock |
| requestedQty | 10 | Exceeds stock |

**Steps:**
1. View cart with item
2. Attempt to set quantity to 10
3. Submit update

**Expected Result:**
- Quantity limited to 5 (available stock)
- Notification shown about stock limitation

**Traceability:**
- AC: [AC-CART-005](../l1/acceptance-criteria.md#ac-cart-005)
- BR: [AMB-OP-012](../l1/business-rules.md#amb-op-012)

---

### TC-AC-CART-005-N03 – Idempotent update when quantity unchanged {#tc-ac-cart-005-n03}

**Preconditions:**
- Customer has item with quantity 3

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| currentQty | 3 | Current quantity |
| newQty | 3 | Same value |

**Steps:**
1. View cart with item (qty 3)
2. Set quantity to 3 again
3. Submit update

**Expected Result:**
- No change to cart
- No unnecessary notifications

**Traceability:**
- AC: [AC-CART-005](../l1/acceptance-criteria.md#ac-cart-005)
- BR: [AMB-OP-013](../l1/business-rules.md#amb-op-013)

---

### TC-AC-CART-006-N01 – Show empty cart state when last item removed {#tc-ac-cart-006-n01}

**Preconditions:**
- Cart has exactly one item

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| cartItems | 1 | Single item in cart |

**Steps:**
1. View cart with single item
2. Remove the item
3. Observe cart state

**Expected Result:**
- Empty cart state displayed
- Continue Shopping link visible

**Traceability:**
- AC: [AC-CART-006](../l1/acceptance-criteria.md#ac-cart-006)
- BR: [AMB-OP-015](../l1/business-rules.md#amb-op-015)

---

### TC-AC-CART-006-N02 – Cannot remove non-existent item {#tc-ac-cart-006-n02}

**Preconditions:**
- Item was already removed or never existed

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| itemId | INVALID-999 | Non-existent item |

**Steps:**
1. Attempt to remove non-existent item via API
2. Check response

**Expected Result:**
- Error returned or no-op
- No system crash

**Traceability:**
- AC: [AC-CART-006](../l1/acceptance-criteria.md#ac-cart-006)
- BR: [AMB-OP-014](../l1/business-rules.md#amb-op-014)

---

### TC-AC-CART-007-N01 – Cannot checkout with only out of stock items {#tc-ac-cart-007-n01}

**Preconditions:**
- Cart contains only out of stock items

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| allItemsStock | 0 | All items OOS |

**Steps:**
1. Have cart with only OOS items
2. Attempt checkout

**Expected Result:**
- Checkout completely blocked
- Message indicates no purchasable items

**Traceability:**
- AC: [AC-CART-007](../l1/acceptance-criteria.md#ac-cart-007)
- BR: [AMB-ENT-014](../l1/business-rules.md#amb-ent-014)

---

### TC-AC-CART-007-N02 – Cannot add more quantity to out of stock item {#tc-ac-cart-007-n02}

**Preconditions:**
- Cart item is out of stock

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| currentQty | 1 | Qty in cart |
| stockLevel | 0 | No stock |

**Steps:**
1. View cart with OOS item (qty 1)
2. Attempt to increase quantity to 2

**Expected Result:**
- Quantity increase rejected
- Stock unavailability message shown

**Traceability:**
- AC: [AC-CART-007](../l1/acceptance-criteria.md#ac-cart-007)
- BR: [AMB-ENT-014](../l1/business-rules.md#amb-ent-014)

---

### TC-AC-CART-008-N01 – Old price not charged at checkout {#tc-ac-cart-008-n01}

**Preconditions:**
- Price changed since item added to cart

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| addedAtPrice | 20.00 | Old price |
| currentPrice | 25.00 | New price |

**Steps:**
1. Add item at $20
2. Price increases to $25
3. Proceed to checkout and pay

**Expected Result:**
- Charged $25 (current price)
- Order reflects current pricing

**Traceability:**
- AC: [AC-CART-008](../l1/acceptance-criteria.md#ac-cart-008)
- BR: [AMB-ENT-016](../l1/business-rules.md#amb-ent-016)

---

### TC-AC-CART-008-N02 – Multiple price changes reflected correctly {#tc-ac-cart-008-n02}

**Preconditions:**
- Price changed multiple times

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| priceHistory | [10, 15, 12, 18] | Price changes |
| finalPrice | 18.00 | Latest price |

**Steps:**
1. Add item at $10
2. Price changes to $15, then $12, then $18
3. View cart

**Expected Result:**
- Cart shows $18 (latest price)
- Intermediate prices not relevant

**Traceability:**
- AC: [AC-CART-008](../l1/acceptance-criteria.md#ac-cart-008)
- BR: [AMB-ENT-016](../l1/business-rules.md#amb-ent-016)

---

### TC-AC-ORDER-001-N01 – Reject order from unregistered customer {#tc-ac-order-001-n01}

**Preconditions:**
- User is not logged in
- Cart has items

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| userStatus | guest | Unregistered user |

**Steps:**
1. Add items to cart as guest
2. Attempt to checkout

**Expected Result:**
- Redirected to registration/login page
- Cart preserved

**Traceability:**
- AC: [AC-ORDER-001](../l1/acceptance-criteria.md#ac-order-001)
- BR: [BR-ORDER](../l1/business-rules.md#br-order)

---

### TC-AC-ORDER-001-N02 – Reject order when email not verified {#tc-ac-order-001-n02}

**Preconditions:**
- Customer registered but email not verified
- Cart has items

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| emailVerified | false | Unverified email |

**Steps:**
1. Login as unverified customer
2. Add items to cart
3. Attempt checkout

**Expected Result:**
- Verification required message shown
- Order not created

**Traceability:**
- AC: [AC-ORDER-001](../l1/acceptance-criteria.md#ac-order-001)
- BR: [BR-ORDER](../l1/business-rules.md#br-order)

---

### TC-AC-ORDER-001-N03 – Handle payment authorization failure {#tc-ac-order-001-n03}

**Preconditions:**
- All prerequisites met except payment

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| paymentMethod | declined_card | Card that will decline |

**Steps:**
1. Complete checkout with failing payment
2. Submit order

**Expected Result:**
- Payment error displayed
- Cart remains intact
- Retry option available

**Traceability:**
- AC: [AC-ORDER-001](../l1/acceptance-criteria.md#ac-order-001)
- BR: [BR-ORDER](../l1/business-rules.md#br-order)

---

### TC-AC-ORDER-001-N04 – Handle out of stock during checkout {#tc-ac-order-001-n04}

**Preconditions:**
- Item becomes out of stock after adding to cart

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| stockLevel | 0 | Item sold out during checkout |

**Steps:**
1. Add item to cart
2. Another user buys last item
3. Submit order

**Expected Result:**
- Error shown for out of stock item
- Cart modification allowed

**Traceability:**
- AC: [AC-ORDER-001](../l1/acceptance-criteria.md#ac-order-001)
- BR: [BR-ORDER](../l1/business-rules.md#br-order)

---

### TC-AC-ORDER-001-N05 – Reject order with invalid shipping address {#tc-ac-order-001-n05}

**Preconditions:**
- All prerequisites met except shipping address

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| zipCode | invalid | Invalid zip code format |

**Steps:**
1. Enter invalid shipping address
2. Attempt to submit order

**Expected Result:**
- Validation errors shown
- Order not submitted

**Traceability:**
- AC: [AC-ORDER-001](../l1/acceptance-criteria.md#ac-order-001)
- BR: [BR-ORDER](../l1/business-rules.md#br-order)

---

### TC-AC-ORDER-002-N01 – Standard shipping fee applied below $50 {#tc-ac-order-002-n01}

**Preconditions:**
- Customer at checkout

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| subtotal | 49.99 | Just below threshold |

**Steps:**
1. Add items totaling $49.99
2. View shipping options

**Expected Result:**
- Standard shipping fee of $5.99 applied

**Traceability:**
- AC: [AC-ORDER-002](../l1/acceptance-criteria.md#ac-order-002)
- BR: [BR-SHIPPING](../l1/business-rules.md#br-shipping)

---

### TC-AC-ORDER-002-N02 – No free shipping when discount brings subtotal below $50 {#tc-ac-order-002-n02}

**Preconditions:**
- Discount will reduce below threshold

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| originalSubtotal | 55.00 | Above threshold |
| discount | 10.00 | Brings below $50 |
| finalSubtotal | 45.00 | After discount |

**Steps:**
1. Add $55 of items
2. Apply $10 discount
3. View shipping

**Expected Result:**
- $5.99 shipping fee applied
- Free shipping not available

**Traceability:**
- AC: [AC-ORDER-002](../l1/acceptance-criteria.md#ac-order-002)
- BR: [BR-SHIPPING](../l1/business-rules.md#br-shipping)

---

### TC-AC-ORDER-003-N01 – Zero tax for states without sales tax {#tc-ac-order-003-n01}

**Preconditions:**
- Shipping to no-sales-tax state

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| shippingState | OR | Oregon has no sales tax |
| subtotal | 100.00 | Test amount |

**Steps:**
1. Enter Oregon shipping address
2. View tax calculation

**Expected Result:**
- Tax shows $0.00
- Total equals subtotal plus shipping

**Traceability:**
- AC: [AC-ORDER-003](../l1/acceptance-criteria.md#ac-order-003)
- BR: [BR-TAX](../l1/business-rules.md#br-tax)

---

### TC-AC-ORDER-003-N02 – Reject checkout without shipping address for tax {#tc-ac-order-003-n02}

**Preconditions:**
- No shipping address provided

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| shippingAddress | null | Missing address |

**Steps:**
1. Attempt to view order total without shipping address

**Expected Result:**
- Tax cannot be calculated
- Shipping address required message

**Traceability:**
- AC: [AC-ORDER-003](../l1/acceptance-criteria.md#ac-order-003)
- BR: [BR-TAX](../l1/business-rules.md#br-tax)

---

### TC-AC-ORDER-004-N01 – Order creation not blocked by email failure {#tc-ac-order-004-n01}

**Preconditions:**
- Email service unavailable

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| emailServiceStatus | down | Service failure |

**Steps:**
1. Disable email service
2. Place order

**Expected Result:**
- Order created successfully
- Order confirmation page shown

**Traceability:**
- AC: [AC-ORDER-004](../l1/acceptance-criteria.md#ac-order-004)
- BR: [BR-EMAIL](../l1/business-rules.md#br-email)

---

### TC-AC-ORDER-004-N02 – Failed email queued for retry {#tc-ac-order-004-n02}

**Preconditions:**
- Email service temporarily fails

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| emailServiceStatus | error | Temporary failure |

**Steps:**
1. Place order during email outage
2. Check retry queue
3. Check logs

**Expected Result:**
- Email added to retry queue
- Failure logged

**Traceability:**
- AC: [AC-ORDER-004](../l1/acceptance-criteria.md#ac-order-004)
- BR: [BR-EMAIL](../l1/business-rules.md#br-email)

---

### TC-AC-ORDER-005-N01 – Show empty state for no orders {#tc-ac-order-005-n01}

**Preconditions:**
- Customer has no orders

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| orderCount | 0 | New customer |

**Steps:**
1. Login as new customer
2. View order history

**Expected Result:**
- Empty state message shown
- Start Shopping link present

**Traceability:**
- AC: [AC-ORDER-005](../l1/acceptance-criteria.md#ac-order-005)
- BR: [BR-ORDER-HISTORY](../l1/business-rules.md#br-order-history)

---

### TC-AC-ORDER-005-N02 – Unregistered user cannot view order history {#tc-ac-order-005-n02}

**Preconditions:**
- User not logged in

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| userStatus | guest | Not registered |

**Steps:**
1. Navigate to order history without login

**Expected Result:**
- Redirect to login page
- Order history not accessible

**Traceability:**
- AC: [AC-ORDER-005](../l1/acceptance-criteria.md#ac-order-005)
- BR: [BR-ORDER-HISTORY](../l1/business-rules.md#br-order-history)

---

### TC-AC-ORDER-006-N01 – Reject cancellation for shipped order {#tc-ac-order-006-n01}

**Preconditions:**
- Order status is 'shipped'

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| orderStatus | shipped | Non-cancellable status |

**Steps:**
1. Navigate to shipped order
2. Attempt to cancel

**Expected Result:**
- Cancellation is blocked
- Error message displayed

**Traceability:**
- AC: [AC-ORDER-006](../l1/acceptance-criteria.md#ac-order-006)
- BR: [AMB-OP-032](../l1/business-rules.md#amb-op-032)

---

### TC-AC-ORDER-006-N02 – Reject cancellation for delivered order {#tc-ac-order-006-n02}

**Preconditions:**
- Order status is 'delivered'

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| orderStatus | delivered | Non-cancellable status |

**Steps:**
1. Navigate to delivered order
2. Attempt to cancel

**Expected Result:**
- Cancellation is blocked
- Error message displayed

**Traceability:**
- AC: [AC-ORDER-006](../l1/acceptance-criteria.md#ac-order-006)
- BR: [AMB-OP-032](../l1/business-rules.md#amb-op-032)

---

### TC-AC-ORDER-006-N03 – Show message for already cancelled order {#tc-ac-order-006-n03}

**Preconditions:**
- Order status is 'cancelled'

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| orderStatus | cancelled | Already cancelled |

**Steps:**
1. Navigate to cancelled order
2. Attempt to cancel again

**Expected Result:**
- Already cancelled message displayed

**Traceability:**
- AC: [AC-ORDER-006](../l1/acceptance-criteria.md#ac-order-006)
- BR: [AMB-ENT-044](../l1/business-rules.md#amb-ent-044)

---

### TC-AC-ORDER-007-N01 – Order total not recalculated on catalog price change {#tc-ac-order-007-n01}

**Preconditions:**
- Order with total $100.00
- Catalog prices changed

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| orderTotal | 100.00 | Original total |

**Steps:**
1. Change all product prices
2. Retrieve order total

**Expected Result:**
- Order total remains $100.00

**Traceability:**
- AC: [AC-ORDER-007](../l1/acceptance-criteria.md#ac-order-007)
- BR: [AMB-ENT-019](../l1/business-rules.md#amb-ent-019)

---

### TC-AC-ORDER-007-N02 – Deleted product does not affect historical order {#tc-ac-order-007-n02}

**Preconditions:**
- Order contains Product A
- Product A later deleted

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| productId | PROD-A | Deleted product |

**Steps:**
1. Delete Product A from catalog
2. View historical order

**Expected Result:**
- Order shows Product A with original price

**Traceability:**
- AC: [AC-ORDER-007](../l1/acceptance-criteria.md#ac-order-007)
- BR: [AMB-ENT-019](../l1/business-rules.md#amb-ent-019)

---

### TC-AC-CUST-001-N01 – Reject registration with duplicate email {#tc-ac-cust-001-n01}

**Preconditions:**
- Email already registered in system

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| email | existing@test.com | Already registered |

**Steps:**
1. Submit registration with existing email

**Expected Result:**
- 'Email already in use' error displayed

**Traceability:**
- AC: [AC-CUST-001](../l1/acceptance-criteria.md#ac-cust-001)
- BR: [AMB-ENT-030](../l1/business-rules.md#amb-ent-030)

---

### TC-AC-CUST-001-N02 – Reject registration with weak password - no number {#tc-ac-cust-001-n02}

**Preconditions:**

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| password | NoNumbers | Missing number |

**Steps:**
1. Submit registration with password lacking number

**Expected Result:**
- Password requirements error displayed

**Traceability:**
- AC: [AC-CUST-001](../l1/acceptance-criteria.md#ac-cust-001)
- BR: [AMB-ENT-024](../l1/business-rules.md#amb-ent-024)

---

### TC-AC-CUST-001-N03 – Reject registration with missing required fields {#tc-ac-cust-001-n03}

**Preconditions:**

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| firstName |  | Missing required field |

**Steps:**
1. Submit registration without first name

**Expected Result:**
- Validation error for first name displayed

**Traceability:**
- AC: [AC-CUST-001](../l1/acceptance-criteria.md#ac-cust-001)
- BR: [AMB-ENT-024](../l1/business-rules.md#amb-ent-024)

---

### TC-AC-CUST-002-N01 – Reject address with missing required fields {#tc-ac-cust-002-n01}

**Preconditions:**
- Customer is logged in

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| city |  | Missing required city |

**Steps:**
1. Submit address without city

**Expected Result:**
- Validation error for city field displayed

**Traceability:**
- AC: [AC-CUST-002](../l1/acceptance-criteria.md#ac-cust-002)
- BR: [AMB-ENT-032](../l1/business-rules.md#amb-ent-032)

---

### TC-AC-CUST-002-N02 – Reject address with invalid postal code format {#tc-ac-cust-002-n02}

**Preconditions:**
- Customer is logged in
- Country selected is USA

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| postalCode | INVALID | Invalid format |

**Steps:**
1. Submit address with invalid postal code

**Expected Result:**
- Postal code format error displayed

**Traceability:**
- AC: [AC-CUST-002](../l1/acceptance-criteria.md#ac-cust-002)
- BR: [AMB-ENT-033](../l1/business-rules.md#amb-ent-033)

---

### TC-AC-CUST-003-N01 – Reject invalid card number {#tc-ac-cust-003-n01}

**Preconditions:**
- Customer is logged in

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| cardNumber | 1234567890123456 | Invalid card |

**Steps:**
1. Enter invalid card number
2. Submit

**Expected Result:**
- Card validation error from gateway displayed

**Traceability:**
- AC: [AC-CUST-003](../l1/acceptance-criteria.md#ac-cust-003)
- BR: [AMB-ENT-035](../l1/business-rules.md#amb-ent-035)

---

### TC-AC-CUST-003-N02 – Reject unsupported card type {#tc-ac-cust-003-n02}

**Preconditions:**
- Customer is logged in

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| cardType | Discover | Not supported |

**Steps:**
1. Attempt to add Discover card

**Expected Result:**
- Error showing supported types: Visa, MC, Amex

**Traceability:**
- AC: [AC-CUST-003](../l1/acceptance-criteria.md#ac-cust-003)
- BR: [AMB-ENT-036](../l1/business-rules.md#amb-ent-036)

---

### TC-AC-CUST-003-N03 – Reject expired card {#tc-ac-cust-003-n03}

**Preconditions:**
- Customer is logged in

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| expiry | 01/20 | Expired card |

**Steps:**
1. Enter expired card
2. Submit

**Expected Result:**
- Card expired error displayed

**Traceability:**
- AC: [AC-CUST-003](../l1/acceptance-criteria.md#ac-cust-003)
- BR: [AMB-ENT-035](../l1/business-rules.md#amb-ent-035)

---

### TC-AC-CUST-004-N01 – Block deletion when customer has pending orders {#tc-ac-cust-004-n01}

**Preconditions:**
- Customer is registered
- Customer has one or more pending orders

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| customerId | CUST-22222 | Customer with pending order |
| pendingOrderId | ORD-99999 | Order in processing state |

**Steps:**
1. Customer requests account deletion
2. System checks for pending orders

**Expected Result:**
- Deletion blocked with error message
- Customer informed of pending orders

**Traceability:**
- AC: [AC-CUST-004](../l1/acceptance-criteria.md#ac-cust-004)
- BR: [BR-CUST](../l1/business-rules.md#br-cust)

---

### TC-AC-CUST-004-N02 – Reject deletion without confirmation {#tc-ac-cust-004-n02}

**Preconditions:**
- Customer is registered
- Customer has no pending orders

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| customerId | CUST-33333 | Valid customer |
| confirmed | false | No confirmation given |

**Steps:**
1. Customer requests account deletion
2. Customer cancels confirmation dialog

**Expected Result:**
- Account remains active
- No data deleted or modified

**Traceability:**
- AC: [AC-CUST-004](../l1/acceptance-criteria.md#ac-cust-004)
- BR: [BR-CUST](../l1/business-rules.md#br-cust)

---

### TC-AC-CUST-004-N03 – Reject deletion for non-existent customer {#tc-ac-cust-004-n03}

**Preconditions:**
- Customer ID does not exist in system

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| customerId | CUST-INVALID | Non-existent customer |

**Steps:**
1. Attempt deletion with invalid customer ID

**Expected Result:**
- Error: Customer not found

**Traceability:**
- AC: [AC-CUST-004](../l1/acceptance-criteria.md#ac-cust-004)
- BR: [BR-CUST](../l1/business-rules.md#br-cust)

---

### TC-AC-ADMIN-001-N01 – Reject product with name too short {#tc-ac-admin-001-n01}

**Preconditions:**
- Admin is logged in

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| name | A | Only 1 character, minimum is 2 |

**Steps:**
1. Enter 1-character product name
2. Attempt to submit

**Expected Result:**
- Validation error: name too short
- Product not created

**Traceability:**
- AC: [AC-ADMIN-001](../l1/acceptance-criteria.md#ac-admin-001)
- BR: [BR-PROD](../l1/business-rules.md#br-prod)

---

### TC-AC-ADMIN-001-N02 – Reject product with name too long {#tc-ac-admin-001-n02}

**Preconditions:**
- Admin is logged in

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| name | [201 character string] | Exceeds 200 char limit |

**Steps:**
1. Enter 201-character product name
2. Attempt to submit

**Expected Result:**
- Validation error: name too long
- Product not created

**Traceability:**
- AC: [AC-ADMIN-001](../l1/acceptance-criteria.md#ac-admin-001)
- BR: [BR-PROD](../l1/business-rules.md#br-prod)

---

### TC-AC-ADMIN-001-N03 – Reject product with price below minimum {#tc-ac-admin-001-n03}

**Preconditions:**
- Admin is logged in

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| price | 0.00 | Below $0.01 minimum |

**Steps:**
1. Enter $0.00 price
2. Attempt to submit

**Expected Result:**
- Error: 'Price must be at least $0.01'
- Product not created

**Traceability:**
- AC: [AC-ADMIN-001](../l1/acceptance-criteria.md#ac-admin-001)
- BR: [BR-PROD](../l1/business-rules.md#br-prod)

---

### TC-AC-ADMIN-001-N04 – Reject product without category {#tc-ac-admin-001-n04}

**Preconditions:**
- Admin is logged in

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| category | null | No category selected |

**Steps:**
1. Fill product data without category
2. Attempt to submit

**Expected Result:**
- Error: 'Category is required'
- Product not created

**Traceability:**
- AC: [AC-ADMIN-001](../l1/acceptance-criteria.md#ac-admin-001)
- BR: [BR-PROD](../l1/business-rules.md#br-prod)

---

### TC-AC-ADMIN-002-N01 – Concurrent edit shows conflict warning {#tc-ac-admin-002-n01}

**Preconditions:**
- Two admins editing same product

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| productId | PROD-003 | Concurrently edited |
| admin1Edit | Name A | First admin change |
| admin2Edit | Name B | Second admin change |

**Steps:**
1. Admin 1 opens product
2. Admin 2 opens same product
3. Admin 1 saves
4. Admin 2 saves

**Expected Result:**
- Admin 2 sees conflict warning
- Last write (Admin 2) persisted

**Traceability:**
- AC: [AC-ADMIN-002](../l1/acceptance-criteria.md#ac-admin-002)
- BR: [BR-PROD](../l1/business-rules.md#br-prod)

---

### TC-AC-ADMIN-002-N02 – Reject edit with name too short {#tc-ac-admin-002-n02}

**Preconditions:**
- Product exists

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| name | X | Only 1 char, minimum is 2 |

**Steps:**
1. Edit product
2. Change name to 1 character
3. Attempt to save

**Expected Result:**
- Validation error shown
- Changes not saved

**Traceability:**
- AC: [AC-ADMIN-002](../l1/acceptance-criteria.md#ac-admin-002)
- BR: [BR-PROD](../l1/business-rules.md#br-prod)

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
2. Change price to $0.00
3. Attempt to save

**Expected Result:**
- Error: Price must be at least $0.01

**Traceability:**
- AC: [AC-ADMIN-002](../l1/acceptance-criteria.md#ac-admin-002)
- BR: [BR-PROD](../l1/business-rules.md#br-prod)

---

### TC-AC-ADMIN-003-N01 – Cancel confirmation takes no action {#tc-ac-admin-003-n01}

**Preconditions:**
- Product exists

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| productId | PROD-008 | Product to not delete |

**Steps:**
1. Click delete on product
2. Cancel in confirmation dialog

**Expected Result:**
- Product remains active
- Product still visible in listings

**Traceability:**
- AC: [AC-ADMIN-003](../l1/acceptance-criteria.md#ac-admin-003)
- BR: [BR-PROD](../l1/business-rules.md#br-prod)

---

### TC-AC-ADMIN-003-N02 – Reject delete for non-existent product {#tc-ac-admin-003-n02}

**Preconditions:**
- Product ID does not exist

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| productId | PROD-INVALID | Non-existent product |

**Steps:**
1. Attempt to delete non-existent product ID

**Expected Result:**
- Error: Product not found

**Traceability:**
- AC: [AC-ADMIN-003](../l1/acceptance-criteria.md#ac-admin-003)
- BR: [BR-PROD](../l1/business-rules.md#br-prod)

---

### TC-AC-ADMIN-004-N01 – Show empty state when no orders match filters {#tc-ac-admin-004-n01}

**Preconditions:**
- Orders exist but none match filter criteria

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| statusFilter | cancelled | No cancelled orders |

**Steps:**
1. Apply filter with no matching orders

**Expected Result:**
- Empty state message displayed
- Clear indication no results found

**Traceability:**
- AC: [AC-ADMIN-004](../l1/acceptance-criteria.md#ac-admin-004)
- BR: [BR-ORDER](../l1/business-rules.md#br-order)

---

### TC-AC-ADMIN-004-N02 – No results for non-existent order number search {#tc-ac-admin-004-n02}

**Preconditions:**
- Search term does not match any order

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| searchTerm | ORD-NONEXIST | Invalid order number |

**Steps:**
1. Search for non-existent order number

**Expected Result:**
- Empty state shown
- No orders in results

**Traceability:**
- AC: [AC-ADMIN-004](../l1/acceptance-criteria.md#ac-admin-004)
- BR: [BR-ORDER](../l1/business-rules.md#br-order)

---

### TC-AC-ADMIN-005-N01 – Reject invalid status transition pending to shipped {#tc-ac-admin-005-n01}

**Preconditions:**
- Admin is logged in
- Order exists with pending status

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| orderId | ORD-004 | Valid order ID |
| newStatus | shipped | Invalid skip transition |

**Steps:**
1. Navigate to order details
2. Attempt to change status to shipped

**Expected Result:**
- Error shown with allowed transitions
- Status remains pending

**Traceability:**
- AC: [AC-ADMIN-005](../l1/acceptance-criteria.md#ac-admin-005)
- BR: [AMB-OP-047](../l1/business-rules.md#amb-op-047)

---

### TC-AC-ADMIN-005-N02 – Reject invalid status transition shipped to confirmed {#tc-ac-admin-005-n02}

**Preconditions:**
- Admin is logged in
- Order exists with shipped status

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| orderId | ORD-005 | Valid order ID |
| newStatus | confirmed | Invalid backward transition |

**Steps:**
1. Navigate to order details
2. Attempt to change status to confirmed

**Expected Result:**
- Error shown with allowed transitions
- Status remains shipped

**Traceability:**
- AC: [AC-ADMIN-005](../l1/acceptance-criteria.md#ac-admin-005)
- BR: [AMB-OP-047](../l1/business-rules.md#amb-op-047)

---

### TC-AC-ADMIN-005-N03 – Reject status update by non-admin user {#tc-ac-admin-005-n03}

**Preconditions:**
- Customer user is logged in
- Order exists

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| userRole | customer | Non-admin role |
| newStatus | confirmed | Any status change |

**Steps:**
1. Attempt to access order status update as customer

**Expected Result:**
- Access denied
- Status unchanged

**Traceability:**
- AC: [AC-ADMIN-005](../l1/acceptance-criteria.md#ac-admin-005)
- BR: [AMB-ENT-018](../l1/business-rules.md#amb-ent-018)

---

### TC-AC-ADMIN-006-N01 – Block inventory adjustment resulting in negative quantity {#tc-ac-admin-006-n01}

**Preconditions:**
- Product variant exists with stock 10

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| variantId | VAR-004 | Valid variant |
| delta | -15 | Would result in -5 |

**Steps:**
1. Navigate to inventory
2. Attempt to reduce by 15
3. Submit

**Expected Result:**
- Error: Inventory cannot go below zero
- Stock remains at 10

**Traceability:**
- AC: [AC-ADMIN-006](../l1/acceptance-criteria.md#ac-admin-006)
- BR: [AMB-OP-053](../l1/business-rules.md#amb-op-053)

---

### TC-AC-ADMIN-006-N02 – Show row-by-row errors for invalid CSV import {#tc-ac-admin-006-n02}

**Preconditions:**
- Admin is logged in
- CSV file with mixed valid and invalid rows

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| csvRow1 | VAR-001,50 | Valid row |
| csvRow2 | INVALID,-10 | Invalid variant and negative |

**Steps:**
1. Navigate to bulk import
2. Upload CSV
3. Submit

**Expected Result:**
- Row 1 imported successfully
- Row 2 shows validation error with details

**Traceability:**
- AC: [AC-ADMIN-006](../l1/acceptance-criteria.md#ac-admin-006)
- BR: [AMB-OP-054](../l1/business-rules.md#amb-op-054)

---

### TC-AC-ADMIN-006-N03 – Block setting absolute value to negative {#tc-ac-admin-006-n03}

**Preconditions:**
- Admin is logged in
- Product variant exists

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| variantId | VAR-005 | Valid variant |
| newQuantity | -5 | Negative absolute value |

**Steps:**
1. Navigate to inventory
2. Set stock to -5
3. Submit

**Expected Result:**
- Error: Inventory cannot go below zero

**Traceability:**
- AC: [AC-ADMIN-006](../l1/acceptance-criteria.md#ac-admin-006)
- BR: [AMB-OP-053](../l1/business-rules.md#amb-op-053)

---

### TC-AC-VAR-001-N01 – Show out of stock status for unavailable variant {#tc-ac-var-001-n01}

**Preconditions:**
- Product exists
- Size M has 0 stock

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| variantId | VAR-M | Out of stock variant |
| stock | 0 | Zero stock |

**Steps:**
1. Navigate to product page
2. Select size M variant

**Expected Result:**
- Out of stock status shown for M
- Add to cart disabled

**Traceability:**
- AC: [AC-VAR-001](../l1/acceptance-criteria.md#ac-var-001)
- BR: [AMB-ENT-009](../l1/business-rules.md#amb-ent-009)

---

### TC-AC-VAR-001-N02 – Allow selecting other variants when one is out of stock {#tc-ac-var-001-n02}

**Preconditions:**
- Product exists
- Size M out of stock
- Size L in stock

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| variantM | 0 stock | Out of stock |
| variantL | 10 stock | In stock |

**Steps:**
1. View product with M out of stock
2. Select size L variant

**Expected Result:**
- Size L selectable
- Add to cart enabled for L

**Traceability:**
- AC: [AC-VAR-001](../l1/acceptance-criteria.md#ac-var-001)
- BR: [AMB-ENT-011](../l1/business-rules.md#amb-ent-011)

---

### TC-AC-CAT-001-N01 – Hide inactive category from customer view {#tc-ac-cat-001-n01}

**Preconditions:**
- Category exists but marked inactive

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| categoryId | CAT-INACTIVE | Inactive category |
| status | inactive | Category status |

**Steps:**
1. Browse as customer
2. Search for inactive category

**Expected Result:**
- Inactive category not visible
- Its products not shown

**Traceability:**
- AC: [AC-CAT-001](../l1/acceptance-criteria.md#ac-cat-001)
- BR: [AMB-ENT-039](../l1/business-rules.md#amb-ent-039)

---

### TC-AC-CAT-001-N02 – Hide products of inactive category from customer view {#tc-ac-cat-001-n02}

**Preconditions:**
- Product assigned to inactive category only

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| productId | PROD-HIDDEN | Product in inactive category |
| categoryStatus | inactive | Parent inactive |

**Steps:**
1. Browse catalog as customer
2. Search for product

**Expected Result:**
- Product not visible in browse
- Product not in search results

**Traceability:**
- AC: [AC-CAT-001](../l1/acceptance-criteria.md#ac-cat-001)
- BR: [AMB-ENT-039](../l1/business-rules.md#amb-ent-039)

---

## Boundary Tests

### TC-AC-PROD-001-B01 – Filter with price range minimum at zero {#tc-ac-prod-001-b01}

**Preconditions:**
- Products exist at various prices including $0
- User on listing page

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| priceMin | 0 | Minimum possible price |
| priceMax | 50 | Upper bound |

**Steps:**
1. Set price range $0-$50
2. Apply filter

**Expected Result:**
- Products priced $0-$50 displayed including free items

**Traceability:**
- AC: [AC-PROD-001](../l1/acceptance-criteria.md#ac-prod-001)
- BR: [BR-PROD](../l1/business-rules.md#br-prod)

---

### TC-AC-PROD-001-B02 – Page boundary at exactly 20 products {#tc-ac-prod-001-b02}

**Preconditions:**
- Exactly 20 products match filter criteria

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| matchingProducts | 20 | Exact page size |

**Steps:**
1. Apply filter matching exactly 20 products

**Expected Result:**
- All 20 products displayed
- No pagination controls shown

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
| inventory | 0 | Exactly zero boundary |

**Steps:**
1. View product with inventory = 0

**Expected Result:**
- Out of Stock indicator shown
- Add to cart disabled

**Traceability:**
- AC: [AC-PROD-002](../l1/acceptance-criteria.md#ac-prod-002)
- BR: [BR-PROD](../l1/business-rules.md#br-prod)

---

### TC-AC-PROD-002-B02 – Product with inventory of 1 does not show out of stock {#tc-ac-prod-002-b02}

**Preconditions:**
- Product inventory is exactly 1

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| inventory | 1 | Just above zero boundary |

**Steps:**
1. View product with inventory = 1

**Expected Result:**
- No Out of Stock indicator
- Add to cart enabled

**Traceability:**
- AC: [AC-PROD-002](../l1/acceptance-criteria.md#ac-prod-002)
- BR: [BR-PROD](../l1/business-rules.md#br-prod)

---

### TC-AC-CART-001-B01 – Add exactly 1 item minimum quantity {#tc-ac-cart-001-b01}

**Preconditions:**
- Product in stock
- User on product page

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| quantity | 1 | Minimum valid quantity |

**Steps:**
1. Set quantity to 1
2. Click Add to Cart

**Expected Result:**
- 1 item added successfully

**Traceability:**
- AC: [AC-CART-001](../l1/acceptance-criteria.md#ac-cart-001)
- BR: [BR-CART](../l1/business-rules.md#br-cart)

---

### TC-AC-CART-001-B02 – Add quantity equal to entire available stock {#tc-ac-cart-001-b02}

**Preconditions:**
- Product has exactly 5 in stock

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| availableStock | 5 | Total stock |
| quantity | 5 | Equals available |

**Steps:**
1. Set quantity to 5 (all available)
2. Click Add to Cart

**Expected Result:**
- All 5 items added successfully

**Traceability:**
- AC: [AC-CART-001](../l1/acceptance-criteria.md#ac-cart-001)
- BR: [BR-CART](../l1/business-rules.md#br-cart)

---

### TC-AC-CART-002-B01 – Add to reach exactly maximum available stock {#tc-ac-cart-002-b01}

**Preconditions:**
- Product (qty 3) in cart
- Exactly 5 in stock

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| existingQty | 3 | In cart |
| addQty | 2 | Adding to reach limit |
| availableStock | 5 | Max available |

**Steps:**
1. Add 2 more to existing 3
2. Combined equals stock of 5

**Expected Result:**
- Quantity updated to 5
- No error or limitation message

**Traceability:**
- AC: [AC-CART-002](../l1/acceptance-criteria.md#ac-cart-002)
- BR: [BR-CART](../l1/business-rules.md#br-cart)

---

### TC-AC-CART-002-B02 – Add 1 more when already at stock limit minus 1 {#tc-ac-cart-002-b02}

**Preconditions:**
- Product (qty 4) in cart
- 5 in stock

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| existingQty | 4 | One below limit |
| addQty | 1 | Adding 1 to reach limit |

**Steps:**
1. Add 1 more when at qty 4 with stock of 5

**Expected Result:**
- Quantity updated to 5 successfully

**Traceability:**
- AC: [AC-CART-002](../l1/acceptance-criteria.md#ac-cart-002)
- BR: [BR-CART](../l1/business-rules.md#br-cart)

---

### TC-AC-CART-003-B01 – Cart persists at exactly 7 days of inactivity {#tc-ac-cart-003-b01}

**Preconditions:**
- Guest cart exists
- Exactly 7 days since last activity

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| daysSinceActivity | 7 | Exactly at boundary |

**Steps:**
1. Add items
2. Return exactly 7 days later

**Expected Result:**
- Cart still contains items

**Traceability:**
- AC: [AC-CART-003](../l1/acceptance-criteria.md#ac-cart-003)
- BR: [BR-CART](../l1/business-rules.md#br-cart)

---

### TC-AC-CART-003-B02 – Cart cleared at 7 days plus 1 second of inactivity {#tc-ac-cart-003-b02}

**Preconditions:**
- Guest cart exists

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| timeSinceActivity | 7 days + 1 second | Just past boundary |

**Steps:**
1. Return just after 7 day expiry threshold

**Expected Result:**
- Cart is cleared
- Empty cart shown

**Traceability:**
- AC: [AC-CART-003](../l1/acceptance-criteria.md#ac-cart-003)
- BR: [BR-CART](../l1/business-rules.md#br-cart)

---

### TC-AC-CART-003-B03 – Activity resets the 7 day timer {#tc-ac-cart-003-b03}

**Preconditions:**
- Cart has items
- 6 days since last activity

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| initialInactivity | 6 | Days without activity |
| action | view cart | Activity performed |

**Steps:**
1. Return after 6 days
2. View cart
3. Wait 6 more days

**Expected Result:**
- Cart still exists
- Timer reset by activity

**Traceability:**
- AC: [AC-CART-003](../l1/acceptance-criteria.md#ac-cart-003)
- BR: [BR-CART](../l1/business-rules.md#br-cart)

---

### TC-AC-CART-004-B01 – Merge with combined quantity at exact stock limit {#tc-ac-cart-004-b01}

**Preconditions:**
- Product has exactly 10 units in stock
- Guest cart has 6, account cart has 4 of same product

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| stockLevel | 10 | Exact boundary |
| guestQty | 6 | Guest quantity |
| accountQty | 4 | Account quantity |

**Steps:**
1. Set stock to exactly 10
2. Add 6 to guest cart, 4 in account cart
3. Log in and merge

**Expected Result:**
- Merged quantity is exactly 10
- No stock warning shown

**Traceability:**
- AC: [AC-CART-004](../l1/acceptance-criteria.md#ac-cart-004)
- BR: [AMB-ENT-013](../l1/business-rules.md#amb-ent-013)

---

### TC-AC-CART-004-B02 – Merge with combined quantity one over stock limit {#tc-ac-cart-004-b02}

**Preconditions:**
- Product has exactly 10 units in stock
- Guest cart has 6, account cart has 5 of same product

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| stockLevel | 10 | Stock limit |
| guestQty | 6 | Guest quantity |
| accountQty | 5 | Account quantity (total 11) |

**Steps:**
1. Set stock to exactly 10
2. Add 6 to guest cart, 5 in account cart
3. Log in and merge

**Expected Result:**
- Merged quantity limited to 10
- Notification shown about quantity adjustment

**Traceability:**
- AC: [AC-CART-004](../l1/acceptance-criteria.md#ac-cart-004)
- BR: [AMB-ENT-013](../l1/business-rules.md#amb-ent-013)

---

### TC-AC-CART-005-B01 – Update quantity to minimum valid value (1) {#tc-ac-cart-005-b01}

**Preconditions:**
- Customer has item with quantity 5

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| newQty | 1 | Minimum valid quantity |

**Steps:**
1. View cart with item
2. Set quantity to 1
3. Submit update

**Expected Result:**
- Quantity updated to 1
- Item remains in cart

**Traceability:**
- AC: [AC-CART-005](../l1/acceptance-criteria.md#ac-cart-005)
- BR: [AMB-OP-011](../l1/business-rules.md#amb-op-011)

---

### TC-AC-CART-005-B02 – Update quantity to exactly available stock {#tc-ac-cart-005-b02}

**Preconditions:**
- Product has exactly 10 units in stock

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| stockLevel | 10 | Exact stock |
| newQty | 10 | At stock limit |

**Steps:**
1. View cart with item
2. Set quantity to exactly 10
3. Submit update

**Expected Result:**
- Quantity updated to 10
- No stock warning shown

**Traceability:**
- AC: [AC-CART-005](../l1/acceptance-criteria.md#ac-cart-005)
- BR: [AMB-OP-012](../l1/business-rules.md#amb-op-012)

---

### TC-AC-CART-006-B01 – Undo option expires after brief period {#tc-ac-cart-006-b01}

**Preconditions:**
- Customer just removed an item

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| waitTime | timeout+1 | Wait past Undo window |

**Steps:**
1. Remove item from cart
2. Wait until Undo option disappears
3. Attempt to click Undo

**Expected Result:**
- Undo option no longer available
- Item removal is permanent

**Traceability:**
- AC: [AC-CART-006](../l1/acceptance-criteria.md#ac-cart-006)
- BR: [AMB-OP-014](../l1/business-rules.md#amb-op-014)

---

### TC-AC-CART-006-B02 – Remove second-to-last item leaves one item {#tc-ac-cart-006-b02}

**Preconditions:**
- Cart has exactly 2 items

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| cartItems | 2 | Two items |

**Steps:**
1. View cart with 2 items
2. Remove one item
3. Check cart state

**Expected Result:**
- One item remains
- Normal cart view shown (not empty state)

**Traceability:**
- AC: [AC-CART-006](../l1/acceptance-criteria.md#ac-cart-006)
- BR: [AMB-OP-015](../l1/business-rules.md#amb-op-015)

---

### TC-AC-CART-007-B01 – Stock drops to exactly zero triggers warning {#tc-ac-cart-007-b01}

**Preconditions:**
- Product has 1 unit in stock
- Customer has 1 in cart

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| initialStock | 1 | Last unit |
| finalStock | 0 | Drops to zero |

**Steps:**
1. Have item in cart (stock = 1)
2. Another user purchases last unit
3. Refresh cart

**Expected Result:**
- Warning indicator appears
- Checkout blocked for item

**Traceability:**
- AC: [AC-CART-007](../l1/acceptance-criteria.md#ac-cart-007)
- BR: [AMB-ENT-014](../l1/business-rules.md#amb-ent-014)

---

### TC-AC-CART-007-B02 – Partial stock allows partial checkout {#tc-ac-cart-007-b02}

**Preconditions:**
- Cart has OOS item and in-stock item

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| itemA_stock | 0 | Out of stock |
| itemB_stock | 5 | In stock |

**Steps:**
1. Have cart with item A (OOS) and item B (in stock)
2. Proceed to checkout

**Expected Result:**
- Item A excluded from checkout
- Item B can be purchased

**Traceability:**
- AC: [AC-CART-007](../l1/acceptance-criteria.md#ac-cart-007)
- BR: [AMB-ENT-014](../l1/business-rules.md#amb-ent-014)

---

### TC-AC-CART-008-B01 – Price change of one cent triggers notification {#tc-ac-cart-008-b01}

**Preconditions:**
- Price changes by minimal amount

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| originalPrice | 9.99 | Original |
| newPrice | 10.00 | One cent more |

**Steps:**
1. Add item at $9.99
2. Price changes to $10.00
3. View cart

**Expected Result:**
- New price $10.00 displayed
- Price change notification shown

**Traceability:**
- AC: [AC-CART-008](../l1/acceptance-criteria.md#ac-cart-008)
- BR: [AMB-ENT-016](../l1/business-rules.md#amb-ent-016)

---

### TC-AC-CART-008-B02 – Price unchanged shows no notification {#tc-ac-cart-008-b02}

**Preconditions:**
- Price has not changed since added

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| originalPrice | 15.00 | Same price |
| currentPrice | 15.00 | Unchanged |

**Steps:**
1. Add item at $15
2. View cart later (price same)
3. Check for notifications

**Expected Result:**
- No price change notification
- Normal cart display

**Traceability:**
- AC: [AC-CART-008](../l1/acceptance-criteria.md#ac-cart-008)
- BR: [AMB-ENT-016](../l1/business-rules.md#amb-ent-016)

---

### TC-AC-ORDER-001-B01 – Place order with exactly 1 item in cart {#tc-ac-order-001-b01}

**Preconditions:**
- Customer verified
- Minimum cart contents

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| cartItems | 1 | Minimum valid cart size |

**Steps:**
1. Add single item to cart
2. Complete checkout

**Expected Result:**
- Order created successfully
- Single item in order

**Traceability:**
- AC: [AC-ORDER-001](../l1/acceptance-criteria.md#ac-order-001)
- BR: [BR-ORDER](../l1/business-rules.md#br-order)

---

### TC-AC-ORDER-001-B02 – Place order with empty cart rejected {#tc-ac-order-001-b02}

**Preconditions:**
- Customer verified
- Cart is empty

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| cartItems | 0 | Empty cart - below minimum |

**Steps:**
1. Navigate to checkout with empty cart

**Expected Result:**
- Checkout blocked
- Message to add items

**Traceability:**
- AC: [AC-ORDER-001](../l1/acceptance-criteria.md#ac-order-001)
- BR: [BR-ORDER](../l1/business-rules.md#br-order)

---

### TC-AC-ORDER-002-B01 – Boundary: $49.99 subtotal gets standard shipping {#tc-ac-order-002-b01}

**Preconditions:**
- Customer at checkout

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| subtotal | 49.99 | One cent below threshold |

**Steps:**
1. Add items totaling $49.99
2. View shipping

**Expected Result:**
- $5.99 shipping applied
- Not free shipping

**Traceability:**
- AC: [AC-ORDER-002](../l1/acceptance-criteria.md#ac-order-002)
- BR: [BR-SHIPPING](../l1/business-rules.md#br-shipping)

---

### TC-AC-ORDER-002-B02 – Boundary: $50.00 subtotal gets free shipping {#tc-ac-order-002-b02}

**Preconditions:**
- Customer at checkout

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| subtotal | 50.00 | Exact threshold |

**Steps:**
1. Add items totaling exactly $50.00
2. View shipping

**Expected Result:**
- Free shipping applied

**Traceability:**
- AC: [AC-ORDER-002](../l1/acceptance-criteria.md#ac-order-002)
- BR: [BR-SHIPPING](../l1/business-rules.md#br-shipping)

---

### TC-AC-ORDER-003-B01 – Tax calculation with minimum order amount {#tc-ac-order-003-b01}

**Preconditions:**
- Minimum possible order value

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| subtotal | 0.01 | Minimum order amount |
| shippingState | CA | Taxable state |

**Steps:**
1. Create order for $0.01
2. View tax calculation

**Expected Result:**
- Tax calculated correctly on small amount

**Traceability:**
- AC: [AC-ORDER-003](../l1/acceptance-criteria.md#ac-order-003)
- BR: [BR-TAX](../l1/business-rules.md#br-tax)

---

### TC-AC-ORDER-003-B02 – Tax rounds correctly to nearest cent {#tc-ac-order-003-b02}

**Preconditions:**
- Amount that produces fractional cents

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| subtotal | 10.33 | Produces fractional tax |
| taxRate | 0.0725 | 7.25% rate |

**Steps:**
1. Order $10.33 with 7.25% tax
2. View calculated tax

**Expected Result:**
- Tax rounded to nearest cent
- No fractional cents displayed

**Traceability:**
- AC: [AC-ORDER-003](../l1/acceptance-criteria.md#ac-order-003)
- BR: [BR-TAX](../l1/business-rules.md#br-tax)

---

### TC-AC-ORDER-004-B01 – Email handles order with single item {#tc-ac-order-004-b01}

**Preconditions:**
- Order with minimum items

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| orderItems | 1 | Single item order |

**Steps:**
1. Place single item order
2. Receive email

**Expected Result:**
- Email correctly formats single item
- No plural issues

**Traceability:**
- AC: [AC-ORDER-004](../l1/acceptance-criteria.md#ac-order-004)
- BR: [BR-EMAIL](../l1/business-rules.md#br-email)

---

### TC-AC-ORDER-004-B02 – Email handles order with many items {#tc-ac-order-004-b02}

**Preconditions:**
- Order with many items

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| orderItems | 50 | Large order |

**Steps:**
1. Place 50-item order
2. Receive email

**Expected Result:**
- All 50 items listed
- Email renders correctly

**Traceability:**
- AC: [AC-ORDER-004](../l1/acceptance-criteria.md#ac-order-004)
- BR: [BR-EMAIL](../l1/business-rules.md#br-email)

---

### TC-AC-ORDER-005-B01 – Pagination with exactly one page of orders {#tc-ac-order-005-b01}

**Preconditions:**
- Orders fit on single page

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| orderCount | 10 | Assuming 10 per page |

**Steps:**
1. View order history with 10 orders

**Expected Result:**
- All orders shown
- No pagination controls or disabled

**Traceability:**
- AC: [AC-ORDER-005](../l1/acceptance-criteria.md#ac-order-005)
- BR: [BR-ORDER-HISTORY](../l1/business-rules.md#br-order-history)

---

### TC-AC-ORDER-005-B02 – Pagination with page boundary orders {#tc-ac-order-005-b02}

**Preconditions:**
- Orders span page boundary

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| orderCount | 11 | One more than page size |

**Steps:**
1. View order history
2. Navigate to page 2

**Expected Result:**
- Page 1 shows 10 orders
- Page 2 shows 1 order

**Traceability:**
- AC: [AC-ORDER-005](../l1/acceptance-criteria.md#ac-order-005)
- BR: [BR-ORDER-HISTORY](../l1/business-rules.md#br-order-history)

---

### TC-AC-ORDER-006-B01 – Cancel order at exact moment before shipping {#tc-ac-order-006-b01}

**Preconditions:**
- Order about to transition to shipped

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| orderStatus | confirmed | Last cancellable state |

**Steps:**
1. Submit cancel request concurrent with ship process

**Expected Result:**
- One operation succeeds, other fails gracefully

**Traceability:**
- AC: [AC-ORDER-006](../l1/acceptance-criteria.md#ac-order-006)
- BR: [AMB-OP-028](../l1/business-rules.md#amb-op-028)

---

### TC-AC-ORDER-006-B02 – Cancel order with zero-value items {#tc-ac-order-006-b02}

**Preconditions:**
- Order total is $0.00 due to promotions

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| orderTotal | 0.00 | Zero value order |

**Steps:**
1. Cancel zero-value order

**Expected Result:**
- Order cancelled
- No refund initiated for zero amount

**Traceability:**
- AC: [AC-ORDER-006](../l1/acceptance-criteria.md#ac-order-006)
- BR: [AMB-OP-030](../l1/business-rules.md#amb-op-030)

---

### TC-AC-ORDER-007-B01 – Price captured for item at minimum price $0.01 {#tc-ac-order-007-b01}

**Preconditions:**
- Product priced at $0.01

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| productPrice | 0.01 | Minimum non-zero price |

**Steps:**
1. Place order with $0.01 product
2. View order

**Expected Result:**
- Order shows $0.01 price

**Traceability:**
- AC: [AC-ORDER-007](../l1/acceptance-criteria.md#ac-order-007)
- BR: [AMB-ENT-019](../l1/business-rules.md#amb-ent-019)

---

### TC-AC-ORDER-007-B02 – Price captured for high-value item {#tc-ac-order-007-b02}

**Preconditions:**
- Product priced at $999,999.99

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| productPrice | 999999.99 | High value item |

**Steps:**
1. Place order with expensive product
2. View order

**Expected Result:**
- Order shows $999,999.99 price correctly

**Traceability:**
- AC: [AC-ORDER-007](../l1/acceptance-criteria.md#ac-order-007)
- BR: [AMB-ENT-019](../l1/business-rules.md#amb-ent-019)

---

### TC-AC-CUST-001-B01 – Password at minimum length 8 characters with number {#tc-ac-cust-001-b01}

**Preconditions:**

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| password | Abcdefg1 | Exactly 8 chars with number |

**Steps:**
1. Submit registration with 8-char password

**Expected Result:**
- Registration succeeds

**Traceability:**
- AC: [AC-CUST-001](../l1/acceptance-criteria.md#ac-cust-001)
- BR: [AMB-ENT-024](../l1/business-rules.md#amb-ent-024)

---

### TC-AC-CUST-001-B02 – Password below minimum length 7 characters {#tc-ac-cust-001-b02}

**Preconditions:**

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| password | Abcde1x | Only 7 chars |

**Steps:**
1. Submit registration with 7-char password

**Expected Result:**
- Password too short error displayed

**Traceability:**
- AC: [AC-CUST-001](../l1/acceptance-criteria.md#ac-cust-001)
- BR: [AMB-ENT-024](../l1/business-rules.md#amb-ent-024)

---

### TC-AC-CUST-002-B01 – Add address with minimum field lengths {#tc-ac-cust-002-b01}

**Preconditions:**
- Customer is logged in

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| street | 1 A | Minimal street |
| city | LA | 2-char city |

**Steps:**
1. Submit address with minimum length values

**Expected Result:**
- Address saved successfully

**Traceability:**
- AC: [AC-CUST-002](../l1/acceptance-criteria.md#ac-cust-002)
- BR: [AMB-ENT-027](../l1/business-rules.md#amb-ent-027)

---

### TC-AC-CUST-002-B02 – Add address with maximum field lengths {#tc-ac-cust-002-b02}

**Preconditions:**
- Customer is logged in

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| street | A repeated 200 times | Max length street |

**Steps:**
1. Submit address with maximum length values

**Expected Result:**
- Address saved or truncation handled gracefully

**Traceability:**
- AC: [AC-CUST-002](../l1/acceptance-criteria.md#ac-cust-002)
- BR: [AMB-ENT-027](../l1/business-rules.md#amb-ent-027)

---

### TC-AC-CUST-003-B01 – Add card expiring this month {#tc-ac-cust-003-b01}

**Preconditions:**
- Customer is logged in

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| expiry | 01/26 | Expires this month |

**Steps:**
1. Add card expiring current month
2. Submit

**Expected Result:**
- Card accepted (still valid this month)

**Traceability:**
- AC: [AC-CUST-003](../l1/acceptance-criteria.md#ac-cust-003)
- BR: [AMB-ENT-028](../l1/business-rules.md#amb-ent-028)

---

### TC-AC-CUST-003-B02 – Add card with minimum CVV length {#tc-ac-cust-003-b02}

**Preconditions:**
- Customer is logged in

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| cvv | 123 | 3-digit CVV |

**Steps:**
1. Add Visa with 3-digit CVV
2. Submit

**Expected Result:**
- Card saved successfully

**Traceability:**
- AC: [AC-CUST-003](../l1/acceptance-criteria.md#ac-cust-003)
- BR: [AMB-ENT-028](../l1/business-rules.md#amb-ent-028)

---

### TC-AC-CUST-004-B01 – Deletion blocked with exactly one pending order {#tc-ac-cust-004-b01}

**Preconditions:**
- Customer has exactly 1 pending order

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| pendingOrderCount | 1 | Minimum blocking condition |

**Steps:**
1. Customer with 1 pending order requests deletion

**Expected Result:**
- Deletion blocked
- Single pending order ID shown in message

**Traceability:**
- AC: [AC-CUST-004](../l1/acceptance-criteria.md#ac-cust-004)
- BR: [BR-CUST](../l1/business-rules.md#br-cust)

---

### TC-AC-CUST-004-B02 – Deletion allowed with zero pending orders {#tc-ac-cust-004-b02}

**Preconditions:**
- Customer has exactly 0 pending orders

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| pendingOrderCount | 0 | Boundary: no blockers |

**Steps:**
1. Customer with 0 pending orders requests deletion
2. Customer confirms

**Expected Result:**
- Deletion proceeds successfully

**Traceability:**
- AC: [AC-CUST-004](../l1/acceptance-criteria.md#ac-cust-004)
- BR: [BR-CUST](../l1/business-rules.md#br-cust)

---

### TC-AC-ADMIN-001-B01 – Create product with minimum name length (2 chars) {#tc-ac-admin-001-b01}

**Preconditions:**
- Admin is logged in

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| name | AB | Exactly 2 characters - minimum |

**Steps:**
1. Enter 2-character product name
2. Fill other valid fields
3. Submit

**Expected Result:**
- Product created successfully

**Traceability:**
- AC: [AC-ADMIN-001](../l1/acceptance-criteria.md#ac-admin-001)
- BR: [BR-PROD](../l1/business-rules.md#br-prod)

---

### TC-AC-ADMIN-001-B02 – Create product with maximum name length (200 chars) {#tc-ac-admin-001-b02}

**Preconditions:**
- Admin is logged in

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| name | [200 character string] | Exactly 200 chars - max |

**Steps:**
1. Enter 200-character product name
2. Fill other valid fields
3. Submit

**Expected Result:**
- Product created successfully

**Traceability:**
- AC: [AC-ADMIN-001](../l1/acceptance-criteria.md#ac-admin-001)
- BR: [BR-PROD](../l1/business-rules.md#br-prod)

---

### TC-AC-ADMIN-001-B03 – Create product with minimum price ($0.01) {#tc-ac-admin-001-b03}

**Preconditions:**
- Admin is logged in

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| price | 0.01 | Minimum valid price |

**Steps:**
1. Enter $0.01 price
2. Fill other valid fields
3. Submit

**Expected Result:**
- Product created with $0.01 price

**Traceability:**
- AC: [AC-ADMIN-001](../l1/acceptance-criteria.md#ac-admin-001)
- BR: [BR-PROD](../l1/business-rules.md#br-prod)

---

### TC-AC-ADMIN-001-B04 – Create product with maximum 10 images {#tc-ac-admin-001-b04}

**Preconditions:**
- Admin is logged in

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| imageCount | 10 | Maximum allowed images |

**Steps:**
1. Create product
2. Upload exactly 10 images
3. Submit

**Expected Result:**
- Product created with 10 images

**Traceability:**
- AC: [AC-ADMIN-001](../l1/acceptance-criteria.md#ac-admin-001)
- BR: [BR-PROD](../l1/business-rules.md#br-prod)

---

### TC-AC-ADMIN-001-B05 – Reject product with 11 images {#tc-ac-admin-001-b05}

**Preconditions:**
- Admin is logged in

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| imageCount | 11 | Exceeds maximum |

**Steps:**
1. Create product
2. Attempt to upload 11 images

**Expected Result:**
- 11th image rejected
- Error shown for image limit

**Traceability:**
- AC: [AC-ADMIN-001](../l1/acceptance-criteria.md#ac-admin-001)
- BR: [BR-PROD](../l1/business-rules.md#br-prod)

---

### TC-AC-ADMIN-001-B06 – Create product with maximum description (5000 chars) {#tc-ac-admin-001-b06}

**Preconditions:**
- Admin is logged in

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| description | [5000 character string] | Maximum length |

**Steps:**
1. Enter 5000-character description
2. Submit

**Expected Result:**
- Product created with full description

**Traceability:**
- AC: [AC-ADMIN-001](../l1/acceptance-criteria.md#ac-admin-001)
- BR: [BR-PROD](../l1/business-rules.md#br-prod)

---

### TC-AC-ADMIN-002-B01 – Edit price to exactly minimum ($0.01) {#tc-ac-admin-002-b01}

**Preconditions:**
- Product exists with higher price

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| newPrice | 0.01 | Minimum valid price |

**Steps:**
1. Edit product price to $0.01
2. Save

**Expected Result:**
- Price updated to $0.01
- Change logged in audit trail

**Traceability:**
- AC: [AC-ADMIN-002](../l1/acceptance-criteria.md#ac-admin-002)
- BR: [BR-PROD](../l1/business-rules.md#br-prod)

---

### TC-AC-ADMIN-002-B02 – Edit name to exactly 200 characters {#tc-ac-admin-002-b02}

**Preconditions:**
- Product exists

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| name | [200 character string] | Maximum valid length |

**Steps:**
1. Edit product name to 200 chars
2. Save

**Expected Result:**
- Name updated successfully

**Traceability:**
- AC: [AC-ADMIN-002](../l1/acceptance-criteria.md#ac-admin-002)
- BR: [BR-PROD](../l1/business-rules.md#br-prod)

---

### TC-AC-ADMIN-003-B01 – Delete product that is in exactly one cart {#tc-ac-admin-003-b01}

**Preconditions:**
- Product exists in exactly 1 cart

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| cartCount | 1 | Single cart with product |

**Steps:**
1. Delete product
2. Check the single affected cart

**Expected Result:**
- Product removed from cart
- Notification queued for customer

**Traceability:**
- AC: [AC-ADMIN-003](../l1/acceptance-criteria.md#ac-admin-003)
- BR: [BR-PROD](../l1/business-rules.md#br-prod), [BR-CART](../l1/business-rules.md#br-cart)

---

### TC-AC-ADMIN-003-B02 – Delete product in multiple carts {#tc-ac-admin-003-b02}

**Preconditions:**
- Product in 100 customer carts

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| cartCount | 100 | Many affected carts |

**Steps:**
1. Delete product
2. Check all affected carts

**Expected Result:**
- Product removed from all 100 carts
- All customers get notification

**Traceability:**
- AC: [AC-ADMIN-003](../l1/acceptance-criteria.md#ac-admin-003)
- BR: [BR-PROD](../l1/business-rules.md#br-prod), [BR-CART](../l1/business-rules.md#br-cart)

---

### TC-AC-ADMIN-004-B01 – Filter with date range of exactly one day {#tc-ac-admin-004-b01}

**Preconditions:**
- Orders exist on specific date

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| startDate | 2024-06-15 | Same day |
| endDate | 2024-06-15 | Same day |

**Steps:**
1. Set start and end date to same day
2. Apply filter

**Expected Result:**
- Only orders from that specific day shown

**Traceability:**
- AC: [AC-ADMIN-004](../l1/acceptance-criteria.md#ac-admin-004)
- BR: [BR-ORDER](../l1/business-rules.md#br-order)

---

### TC-AC-ADMIN-004-B02 – Display order list with exactly one order {#tc-ac-admin-004-b02}

**Preconditions:**
- Only one order exists in system

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| orderCount | 1 | Single order |

**Steps:**
1. View order management with only one order

**Expected Result:**
- Single order displayed correctly
- All columns visible

**Traceability:**
- AC: [AC-ADMIN-004](../l1/acceptance-criteria.md#ac-admin-004)
- BR: [BR-ORDER](../l1/business-rules.md#br-order)

---

### TC-AC-ADMIN-005-B01 – Update order through complete status lifecycle {#tc-ac-admin-005-b01}

**Preconditions:**
- Admin is logged in
- Order exists with pending status

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| orderId | ORD-006 | Valid order ID |
| transitions | all four states | Complete lifecycle |

**Steps:**
1. Update pending to confirmed
2. Update confirmed to shipped
3. Update shipped to delivered

**Expected Result:**
- All transitions succeed
- Four audit log entries created
- Three emails sent

**Traceability:**
- AC: [AC-ADMIN-005](../l1/acceptance-criteria.md#ac-admin-005)
- BR: [AMB-OP-050](../l1/business-rules.md#amb-op-050)

---

### TC-AC-ADMIN-005-B02 – Reject status change from delivered (final state) {#tc-ac-admin-005-b02}

**Preconditions:**
- Admin is logged in
- Order exists with delivered status

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| orderId | ORD-007 | Valid order ID |
| currentStatus | delivered | Final state |

**Steps:**
1. Navigate to delivered order
2. Attempt any status change

**Expected Result:**
- No valid transitions available
- Status remains delivered

**Traceability:**
- AC: [AC-ADMIN-005](../l1/acceptance-criteria.md#ac-admin-005)
- BR: [AMB-OP-047](../l1/business-rules.md#amb-op-047)

---

### TC-AC-ADMIN-006-B01 – Set inventory to exactly zero {#tc-ac-admin-006-b01}

**Preconditions:**
- Product variant exists with stock 10

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| variantId | VAR-006 | Valid variant |
| newQuantity | 0 | Boundary: minimum valid |

**Steps:**
1. Navigate to inventory
2. Set stock to 0
3. Submit

**Expected Result:**
- Stock updated to 0
- Audit trail recorded

**Traceability:**
- AC: [AC-ADMIN-006](../l1/acceptance-criteria.md#ac-admin-006)
- BR: [AMB-OP-051](../l1/business-rules.md#amb-op-051)

---

### TC-AC-ADMIN-006-B02 – Reject delta that results in exactly -1 {#tc-ac-admin-006-b02}

**Preconditions:**
- Product variant exists with stock 5

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| variantId | VAR-007 | Valid variant |
| delta | -6 | Results in -1 |

**Steps:**
1. Navigate to inventory
2. Attempt to reduce by 6
3. Submit

**Expected Result:**
- Error: Inventory cannot go below zero
- Stock remains at 5

**Traceability:**
- AC: [AC-ADMIN-006](../l1/acceptance-criteria.md#ac-admin-006)
- BR: [AMB-OP-053](../l1/business-rules.md#amb-op-053)

---

### TC-AC-VAR-001-B01 – Display product with single variant only {#tc-ac-var-001-b01}

**Preconditions:**
- Product exists with only one variant

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| productId | PROD-002 | Single variant product |
| variantCount | 1 | Minimum variants |

**Steps:**
1. Navigate to single-variant product page

**Expected Result:**
- Variant auto-selected or hidden
- Stock status displayed

**Traceability:**
- AC: [AC-VAR-001](../l1/acceptance-criteria.md#ac-var-001)
- BR: [AMB-ENT-008](../l1/business-rules.md#amb-ent-008)

---

### TC-AC-VAR-001-B02 – Display product with many variant combinations {#tc-ac-var-001-b02}

**Preconditions:**
- Product with 4 sizes x 5 colors = 20 variants

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| sizes | S,M,L,XL | 4 sizes |
| colors | 5 colors | Multiple colors |

**Steps:**
1. Navigate to product page
2. View all variant combinations

**Expected Result:**
- All combinations accessible
- Stock shown per combination

**Traceability:**
- AC: [AC-VAR-001](../l1/acceptance-criteria.md#ac-var-001)
- BR: [AMB-ENT-008](../l1/business-rules.md#amb-ent-008)

---

### TC-AC-CAT-001-B01 – Display exactly 3 levels of hierarchy (maximum) {#tc-ac-cat-001-b01}

**Preconditions:**
- Full 3-level hierarchy exists

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| levels | 3 | Maximum depth |
| hierarchy | L1>L2>L3 | Full depth |

**Steps:**
1. Navigate through all 3 levels

**Expected Result:**
- All 3 levels navigable
- Level 3 categories show products

**Traceability:**
- AC: [AC-CAT-001](../l1/acceptance-criteria.md#ac-cat-001)
- BR: [AMB-ENT-006](../l1/business-rules.md#amb-ent-006)

---

### TC-AC-CAT-001-B02 – Display single-level category (minimum) {#tc-ac-cat-001-b02}

**Preconditions:**
- Top-level category with no children exists

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| levels | 1 | Minimum depth |
| category | Clearance | Top-level only |

**Steps:**
1. Navigate to single-level category

**Expected Result:**
- Category displays products directly
- No subcategory UI shown

**Traceability:**
- AC: [AC-CAT-001](../l1/acceptance-criteria.md#ac-cat-001)
- BR: [AMB-ENT-006](../l1/business-rules.md#amb-ent-006)

---

## Hallucination Prevention Tests

### TC-AC-PROD-001-H01 – Should not require login to browse products {#tc-ac-prod-001-h01}

**⚠️ Should NOT:** Should NOT require authentication to view product listing

**Preconditions:**
- User is not logged in

**Steps:**
1. Navigate to product listing without authentication

**Expected Result:**
- Products displayed without login requirement

**Traceability:**
- AC: [AC-PROD-001](../l1/acceptance-criteria.md#ac-prod-001)
- BR: [BR-PROD](../l1/business-rules.md#br-prod)

---

### TC-AC-PROD-001-H02 – Should not auto-apply previously used filters {#tc-ac-prod-001-h02}

**⚠️ Should NOT:** Should NOT persist filter selections across sessions

**Preconditions:**
- User previously applied filters
- User returns to listing page

**Steps:**
1. Apply filters
2. Navigate away
3. Return to product listing

**Expected Result:**
- No filters pre-applied
- All products shown by default

**Traceability:**
- AC: [AC-PROD-001](../l1/acceptance-criteria.md#ac-prod-001)
- BR: [BR-PROD](../l1/business-rules.md#br-prod)

---

### TC-AC-PROD-002-H01 – Should not hide out of stock products from listing {#tc-ac-prod-002-h01}

**⚠️ Should NOT:** Should NOT hide or remove out of stock products from listing

**Preconditions:**
- Out of stock products exist

**Steps:**
1. Navigate to product listing
2. Search for out-of-stock products

**Expected Result:**
- Out of stock products visible in listing

**Traceability:**
- AC: [AC-PROD-002](../l1/acceptance-criteria.md#ac-prod-002)
- BR: [BR-PROD](../l1/business-rules.md#br-prod)

---

### TC-AC-PROD-002-H02 – Should not show notify me or waitlist option {#tc-ac-prod-002-h02}

**⚠️ Should NOT:** Should NOT show 'Notify Me' or waitlist functionality

**Preconditions:**
- Product is out of stock

**Steps:**
1. View out-of-stock product detail page

**Expected Result:**
- Only Out of Stock indicator and disabled button shown

**Traceability:**
- AC: [AC-PROD-002](../l1/acceptance-criteria.md#ac-prod-002)
- BR: [BR-PROD](../l1/business-rules.md#br-prod)

---

### TC-AC-CART-001-H01 – Should not require account to add to cart {#tc-ac-cart-001-h01}

**⚠️ Should NOT:** Should NOT require login or account creation to add to cart

**Preconditions:**
- User is not logged in
- Product in stock

**Steps:**
1. Browse as guest
2. Add item to cart

**Expected Result:**
- Item added without login prompt

**Traceability:**
- AC: [AC-CART-001](../l1/acceptance-criteria.md#ac-cart-001)
- BR: [BR-CART](../l1/business-rules.md#br-cart)

---

### TC-AC-CART-001-H02 – Should not reserve inventory when adding to cart {#tc-ac-cart-001-h02}

**⚠️ Should NOT:** Should NOT reserve or hold inventory just from cart addition

**Preconditions:**
- Product has 5 in stock

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| quantity | 3 | Added to cart |

**Steps:**
1. Add 3 items to cart
2. Check available stock for other users

**Expected Result:**
- Other users still see 5 available

**Traceability:**
- AC: [AC-CART-001](../l1/acceptance-criteria.md#ac-cart-001)
- BR: [BR-CART](../l1/business-rules.md#br-cart)

---

### TC-AC-CART-002-H01 – Should not split into separate line items based on add time {#tc-ac-cart-002-h01}

**⚠️ Should NOT:** Should NOT create separate line items based on when items were added

**Preconditions:**
- Product in cart from yesterday

**Steps:**
1. Add same product today that was added yesterday

**Expected Result:**
- Quantities combined into single line item

**Traceability:**
- AC: [AC-CART-002](../l1/acceptance-criteria.md#ac-cart-002)
- BR: [BR-CART](../l1/business-rules.md#br-cart)

---

### TC-AC-CART-002-H02 – Should not auto-remove items when combining exceeds stock {#tc-ac-cart-002-h02}

**⚠️ Should NOT:** Should NOT remove existing cart items to accommodate new addition

**Preconditions:**
- Product qty 3 in cart
- Only 4 in stock

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| existingQty | 3 | In cart |
| addQty | 5 | Exceeds when combined |

**Steps:**
1. Try to add 5 more when 3 in cart with only 4 stock

**Expected Result:**
- Quantity limited to 4
- Original items preserved

**Traceability:**
- AC: [AC-CART-002](../l1/acceptance-criteria.md#ac-cart-002)
- BR: [BR-CART](../l1/business-rules.md#br-cart)

---

### TC-AC-CART-003-H01 – Should not merge guest cart with account on login {#tc-ac-cart-003-h01}

**⚠️ Should NOT:** Should NOT auto-merge carts without explicit AC requirement

**Preconditions:**
- Guest has items in cart
- Guest logs in to account

**Steps:**
1. Add items as guest
2. Log in to existing account

**Expected Result:**
- Behavior not specified in AC

**Traceability:**
- AC: [AC-CART-003](../l1/acceptance-criteria.md#ac-cart-003)
- BR: [BR-CART](../l1/business-rules.md#br-cart)

---

### TC-AC-CART-003-H02 – Should not send cart recovery emails to guests {#tc-ac-cart-003-h02}

**⚠️ Should NOT:** Should NOT send cart recovery or reminder emails to guest users

**Preconditions:**
- Guest has abandoned cart
- Session approaching expiry

**Steps:**
1. Add items as guest
2. Do not return for 6 days

**Expected Result:**
- No email sent about cart

**Traceability:**
- AC: [AC-CART-003](../l1/acceptance-criteria.md#ac-cart-003)
- BR: [BR-CART](../l1/business-rules.md#br-cart)

---

### TC-AC-CART-004-H01 – Guest cart items should NOT be lost on merge {#tc-ac-cart-004-h01}

**⚠️ Should NOT:** Guest cart items should NOT be discarded or replaced by account cart

**Preconditions:**
- Guest cart has unique items not in account cart

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| guestItems | [{sku: UNIQUE1}, {sku: UNIQUE2}] | Unique guest items |

**Steps:**
1. Add unique items to guest cart
2. Log in with account
3. Verify all items present

**Expected Result:**
- All guest items preserved in merged cart

**Traceability:**
- AC: [AC-CART-004](../l1/acceptance-criteria.md#ac-cart-004)
- BR: [AMB-ENT-013](../l1/business-rules.md#amb-ent-013)

---

### TC-AC-CART-004-H02 – Account cart items should NOT be replaced by guest cart {#tc-ac-cart-004-h02}

**⚠️ Should NOT:** Account cart items should NOT be deleted or overwritten by guest cart

**Preconditions:**
- Account cart has items before login

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| accountItems | [{sku: ACC1, qty: 5}] | Pre-existing account items |

**Steps:**
1. Have account cart with item ACC1 (qty 5)
2. Add different items to guest cart
3. Log in

**Expected Result:**
- Account cart item ACC1 still present with qty 5

**Traceability:**
- AC: [AC-CART-004](../l1/acceptance-criteria.md#ac-cart-004)
- BR: [AMB-ENT-013](../l1/business-rules.md#amb-ent-013)

---

### TC-AC-CART-005-H01 – Negative quantity should NOT be accepted {#tc-ac-cart-005-h01}

**⚠️ Should NOT:** System should NOT accept negative quantities or credit customer

**Preconditions:**
- Customer has item in cart

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| newQty | -1 | Invalid negative |

**Steps:**
1. View cart with item
2. Attempt to set quantity to -1
3. Submit update

**Expected Result:**
- Negative quantity rejected

**Traceability:**
- AC: [AC-CART-005](../l1/acceptance-criteria.md#ac-cart-005)
- BR: [AMB-OP-011](../l1/business-rules.md#amb-op-011)

---

### TC-AC-CART-005-H02 – Update should NOT affect other cart items {#tc-ac-cart-005-h02}

**⚠️ Should NOT:** Updating one item should NOT modify quantities of other items

**Preconditions:**
- Cart has items A, B, and C

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| targetItem | A | Item being updated |
| otherItems | [B, C] | Should be unaffected |

**Steps:**
1. View cart with multiple items
2. Update quantity of item A only
3. Check other items

**Expected Result:**
- Items B and C unchanged

**Traceability:**
- AC: [AC-CART-005](../l1/acceptance-criteria.md#ac-cart-005)
- BR: [AMB-OP-011](../l1/business-rules.md#amb-op-011)

---

### TC-AC-CART-006-H01 – Confirmation dialog should NOT be shown {#tc-ac-cart-006-h01}

**⚠️ Should NOT:** System should NOT show 'Are you sure?' confirmation dialog

**Preconditions:**
- Customer has item in cart

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| itemToRemove | any | Any cart item |

**Steps:**
1. Click remove on item
2. Check for confirmation dialog

**Expected Result:**
- Item removed without dialog

**Traceability:**
- AC: [AC-CART-006](../l1/acceptance-criteria.md#ac-cart-006)
- BR: [AMB-OP-014](../l1/business-rules.md#amb-op-014)

---

### TC-AC-CART-006-H02 – Removed item should NOT appear in order {#tc-ac-cart-006-h02}

**⚠️ Should NOT:** Removed items should NOT appear in final order or be charged

**Preconditions:**
- Customer removed item, then checks out

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| removedItem | SKU-789 | Previously removed |

**Steps:**
1. Remove item from cart (no Undo)
2. Proceed to checkout
3. Complete order

**Expected Result:**
- Removed item not in order

**Traceability:**
- AC: [AC-CART-006](../l1/acceptance-criteria.md#ac-cart-006)
- BR: [AMB-OP-014](../l1/business-rules.md#amb-op-014)

---

### TC-AC-CART-007-H01 – Out of stock item should NOT be auto-removed {#tc-ac-cart-007-h01}

**⚠️ Should NOT:** System should NOT automatically remove out of stock items from cart

**Preconditions:**
- Item goes out of stock

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| stock | 0 | Out of stock |

**Steps:**
1. Have item in cart
2. Stock drops to zero
3. Check cart

**Expected Result:**
- Item still in cart (with warning)

**Traceability:**
- AC: [AC-CART-007](../l1/acceptance-criteria.md#ac-cart-007)
- BR: [AMB-ENT-014](../l1/business-rules.md#amb-ent-014)

---

### TC-AC-CART-007-H02 – Out of stock item should NOT be purchasable {#tc-ac-cart-007-h02}

**⚠️ Should NOT:** System should NOT allow purchase of zero-stock items via any method

**Preconditions:**
- Cart contains OOS item

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| stock | 0 | No inventory |

**Steps:**
1. Have OOS item in cart
2. Attempt various checkout methods

**Expected Result:**
- OOS item cannot be ordered

**Traceability:**
- AC: [AC-CART-007](../l1/acceptance-criteria.md#ac-cart-007)
- BR: [AMB-ENT-014](../l1/business-rules.md#amb-ent-014)

---

### TC-AC-CART-008-H01 – Old price should NOT be honored at checkout {#tc-ac-cart-008-h01}

**⚠️ Should NOT:** System should NOT allow customers to pay the old lower price

**Preconditions:**
- Price increased after adding to cart

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| oldPrice | 10.00 | When added |
| newPrice | 15.00 | Current |

**Steps:**
1. Add at $10, price increases to $15
2. Request old price at checkout

**Expected Result:**
- Only current price honored

**Traceability:**
- AC: [AC-CART-008](../l1/acceptance-criteria.md#ac-cart-008)
- BR: [AMB-ENT-016](../l1/business-rules.md#amb-ent-016)

---

### TC-AC-CART-008-H02 – Cached prices should NOT override current prices {#tc-ac-cart-008-h02}

**⚠️ Should NOT:** Stale cached prices should NOT be displayed or charged

**Preconditions:**
- Browser has cached cart data

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| cachedPrice | 20.00 | Old cached |
| currentPrice | 25.00 | Server price |

**Steps:**
1. Load cart from cache
2. Verify price displayed

**Expected Result:**
- Current server price shown

**Traceability:**
- AC: [AC-CART-008](../l1/acceptance-criteria.md#ac-cart-008)
- BR: [AMB-ENT-016](../l1/business-rules.md#amb-ent-016)

---

### TC-AC-ORDER-001-H01 – Order should NOT be created if payment not authorized {#tc-ac-order-001-h01}

**⚠️ Should NOT:** Create order before payment authorization completes

**Preconditions:**
- Payment authorization pending or failed

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| paymentStatus | pending | Not yet authorized |

**Steps:**
1. Submit order
2. Payment auth in progress
3. Check database

**Expected Result:**
- No order record created until payment authorized

**Traceability:**
- AC: [AC-ORDER-001](../l1/acceptance-criteria.md#ac-order-001)
- BR: [BR-ORDER](../l1/business-rules.md#br-order)

---

### TC-AC-ORDER-001-H02 – Inventory should NOT decrement on failed orders {#tc-ac-order-001-h02}

**⚠️ Should NOT:** Decrement inventory when payment fails

**Preconditions:**
- Payment will fail
- Item has stock of 10

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| initialStock | 10 | Track inventory |

**Steps:**
1. Attempt order with failing payment
2. Check inventory levels

**Expected Result:**
- Inventory remains at 10

**Traceability:**
- AC: [AC-ORDER-001](../l1/acceptance-criteria.md#ac-order-001)
- BR: [BR-ORDER](../l1/business-rules.md#br-order)

---

### TC-AC-ORDER-002-H01 – Tax should NOT be included in free shipping calculation {#tc-ac-order-002-h01}

**⚠️ Should NOT:** Include tax in free shipping threshold calculation

**Preconditions:**
- Subtotal $45, tax would bring total over $50

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| subtotal | 45.00 | Below threshold |
| tax | 8.00 | Would make total $53 |

**Steps:**
1. Add $45 items
2. Tax calculated as $8
3. View shipping

**Expected Result:**
- $5.99 shipping applied despite total being $53

**Traceability:**
- AC: [AC-ORDER-002](../l1/acceptance-criteria.md#ac-order-002)
- BR: [BR-SHIPPING](../l1/business-rules.md#br-shipping)

---

### TC-AC-ORDER-002-H02 – Free shipping should NOT require manual selection {#tc-ac-order-002-h02}

**⚠️ Should NOT:** Require customer to manually select free shipping option

**Preconditions:**
- Subtotal is $60

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| subtotal | 60.00 | Above threshold |

**Steps:**
1. Add $60 of items
2. Proceed to checkout

**Expected Result:**
- Free shipping pre-selected without user action

**Traceability:**
- AC: [AC-ORDER-002](../l1/acceptance-criteria.md#ac-order-002)
- BR: [BR-SHIPPING](../l1/business-rules.md#br-shipping)

---

### TC-AC-ORDER-003-H01 – Tax should NOT be based on billing address {#tc-ac-order-003-h01}

**⚠️ Should NOT:** Use billing address state for tax calculation

**Preconditions:**
- Different billing and shipping addresses

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| billingState | CA | High tax state |
| shippingState | OR | No tax state |

**Steps:**
1. Enter CA billing, OR shipping
2. View tax calculation

**Expected Result:**
- $0 tax (Oregon rate)
- CA billing address ignored for tax

**Traceability:**
- AC: [AC-ORDER-003](../l1/acceptance-criteria.md#ac-order-003)
- BR: [BR-TAX](../l1/business-rules.md#br-tax)

---

### TC-AC-ORDER-003-H02 – International addresses should NOT calculate US tax {#tc-ac-order-003-h02}

**⚠️ Should NOT:** Apply US state tax to international shipments

**Preconditions:**
- Shipping to non-US address

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| shippingCountry | Canada | Non-US country |

**Steps:**
1. Enter Canadian shipping address
2. View tax calculation

**Expected Result:**
- No US state tax applied

**Traceability:**
- AC: [AC-ORDER-003](../l1/acceptance-criteria.md#ac-order-003)
- BR: [BR-TAX](../l1/business-rules.md#br-tax)

---

### TC-AC-ORDER-004-H01 – Email should NOT be sent synchronously {#tc-ac-order-004-h01}

**⚠️ Should NOT:** Block order completion while sending email

**Preconditions:**
- Order being placed

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| emailDelay | 5000ms | Slow email service |

**Steps:**
1. Place order
2. Measure order completion time

**Expected Result:**
- Order completes without waiting for email

**Traceability:**
- AC: [AC-ORDER-004](../l1/acceptance-criteria.md#ac-order-004)
- BR: [BR-EMAIL](../l1/business-rules.md#br-email)

---

### TC-AC-ORDER-004-H02 – Email should NOT include payment card details {#tc-ac-order-004-h02}

**⚠️ Should NOT:** Include full payment card details in confirmation email

**Preconditions:**
- Order placed with credit card

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| paymentMethod | card_ending_4242 | Credit card |

**Steps:**
1. Place order with card
2. Receive email
3. Check contents

**Expected Result:**
- No card number in email
- No CVV in email

**Traceability:**
- AC: [AC-ORDER-004](../l1/acceptance-criteria.md#ac-order-004)
- BR: [BR-EMAIL](../l1/business-rules.md#br-email)

---

### TC-AC-ORDER-005-H01 – Should NOT show other customers orders {#tc-ac-order-005-h01}

**⚠️ Should NOT:** Display orders belonging to other customers

**Preconditions:**
- Multiple customers with orders exist

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| customerAOrders | 5 | Current user |
| customerBOrders | 10 | Other user |

**Steps:**
1. Login as Customer A
2. View order history

**Expected Result:**
- Only 5 orders shown
- Customer B orders not visible

**Traceability:**
- AC: [AC-ORDER-005](../l1/acceptance-criteria.md#ac-order-005)
- BR: [BR-ORDER-HISTORY](../l1/business-rules.md#br-order-history)

---

### TC-AC-ORDER-005-H02 – Tracking number should NOT show for pending orders {#tc-ac-order-005-h02}

**⚠️ Should NOT:** Show tracking number placeholder for unshipped orders

**Preconditions:**
- Order in pending status

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| orderStatus | pending | Not yet shipped |

**Steps:**
1. View pending order details

**Expected Result:**
- Tracking number field empty or hidden

**Traceability:**
- AC: [AC-ORDER-005](../l1/acceptance-criteria.md#ac-order-005)
- BR: [BR-ORDER-HISTORY](../l1/business-rules.md#br-order-history)

---

### TC-AC-ORDER-006-H01 – Cancelled order should NOT appear in active orders {#tc-ac-order-006-h01}

**⚠️ Should NOT:** Cancelled order should NOT appear in active/open orders list

**Preconditions:**
- Order has been cancelled

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| orderId | ORD-12345 | Cancelled order |

**Steps:**
1. Navigate to active orders list

**Expected Result:**
- Cancelled order not in active list

**Traceability:**
- AC: [AC-ORDER-006](../l1/acceptance-criteria.md#ac-order-006)
- BR: [AMB-ENT-044](../l1/business-rules.md#amb-ent-044)

---

### TC-AC-ORDER-006-H02 – Customer should NOT cancel another customer order {#tc-ac-order-006-h02}

**⚠️ Should NOT:** System should NOT allow cancelling orders owned by others

**Preconditions:**
- Customer A logged in
- Order belongs to Customer B

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| customerId | CUST-A | Logged in customer |
| orderOwner | CUST-B | Different customer |

**Steps:**
1. Attempt to cancel another customers order

**Expected Result:**
- Access denied or order not found

**Traceability:**
- AC: [AC-ORDER-006](../l1/acceptance-criteria.md#ac-order-006)
- BR: [AMB-OP-028](../l1/business-rules.md#amb-op-028)

---

### TC-AC-ORDER-007-H01 – Order should NOT auto-apply later discounts {#tc-ac-order-007-h01}

**⚠️ Should NOT:** System should NOT retroactively apply discounts to orders

**Preconditions:**
- Order placed at full price
- Discount added later

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| discount | 20% | Later promotion |

**Steps:**
1. Add store-wide 20% discount
2. View existing order

**Expected Result:**
- Order prices unchanged

**Traceability:**
- AC: [AC-ORDER-007](../l1/acceptance-criteria.md#ac-order-007)
- BR: [AMB-ENT-019](../l1/business-rules.md#amb-ent-019)

---

### TC-AC-ORDER-007-H02 – Order should NOT reflect exchange rate changes {#tc-ac-order-007-h02}

**⚠️ Should NOT:** System should NOT recalculate order based on exchange rates

**Preconditions:**
- Order placed in USD
- Exchange rates change

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| currency | USD | Order currency |

**Steps:**
1. Change exchange rates
2. View order in different currency

**Expected Result:**
- Original USD amounts unchanged

**Traceability:**
- AC: [AC-ORDER-007](../l1/acceptance-criteria.md#ac-order-007)
- BR: [AMB-ENT-019](../l1/business-rules.md#amb-ent-019)

---

### TC-AC-CUST-001-H01 – Phone number should NOT be required for registration {#tc-ac-cust-001-h01}

**⚠️ Should NOT:** System should NOT require phone number for registration

**Preconditions:**

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

### TC-AC-CUST-001-H02 – Unverified user should NOT place orders {#tc-ac-cust-001-h02}

**⚠️ Should NOT:** Unverified users should NOT be able to place orders

**Preconditions:**
- User registered but not verified

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| emailVerified | false | Unverified user |

**Steps:**
1. Login as unverified user
2. Attempt to place order

**Expected Result:**
- Order blocked with verification prompt

**Traceability:**
- AC: [AC-CUST-001](../l1/acceptance-criteria.md#ac-cust-001)
- BR: [AMB-ENT-026](../l1/business-rules.md#amb-ent-026)

---

### TC-AC-CUST-002-H01 – Address should NOT require apartment/unit number {#tc-ac-cust-002-h01}

**⚠️ Should NOT:** System should NOT require apartment/unit number

**Preconditions:**
- Customer is logged in

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| apartment |  | Not provided |

**Steps:**
1. Add address without apartment number

**Expected Result:**
- Address saves without apartment field

**Traceability:**
- AC: [AC-CUST-002](../l1/acceptance-criteria.md#ac-cust-002)
- BR: [AMB-ENT-027](../l1/business-rules.md#amb-ent-027)

---

### TC-AC-CUST-002-H02 – Deleted address should NOT appear in checkout {#tc-ac-cust-002-h02}

**⚠️ Should NOT:** Deleted addresses should NOT appear in checkout options

**Preconditions:**
- Customer deleted an address

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| addressId | ADDR-001 | Deleted address |

**Steps:**
1. Delete address
2. Go to checkout

**Expected Result:**
- Deleted address not in shipping options

**Traceability:**
- AC: [AC-CUST-002](../l1/acceptance-criteria.md#ac-cust-002)
- BR: [AMB-ENT-027](../l1/business-rules.md#amb-ent-027)

---

### TC-AC-CUST-003-H01 – Raw card number should NOT be stored {#tc-ac-cust-003-h01}

**⚠️ Should NOT:** System should NOT store raw card numbers, only tokens

**Preconditions:**
- Customer adds payment method

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| cardNumber | 4242424242424242 | Full number |

**Steps:**
1. Add card
2. Query database for card data

**Expected Result:**
- Only token stored, not full card number

**Traceability:**
- AC: [AC-CUST-003](../l1/acceptance-criteria.md#ac-cust-003)
- BR: [AMB-ENT-028](../l1/business-rules.md#amb-ent-028)

---

### TC-AC-CUST-003-H02 – Debit cards should NOT be blocked if gateway accepts {#tc-ac-cust-003-h02}

**⚠️ Should NOT:** System should NOT explicitly block debit cards

**Preconditions:**
- Customer has Visa debit card

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| cardType | Visa Debit | Debit not credit |

**Steps:**
1. Add Visa debit card

**Expected Result:**
- Card accepted if gateway validates

**Traceability:**
- AC: [AC-CUST-003](../l1/acceptance-criteria.md#ac-cust-003)
- BR: [AMB-ENT-028](../l1/business-rules.md#amb-ent-028)

---

### TC-AC-CUST-004-H01 – Deleted customer should NOT appear in customer listings {#tc-ac-cust-004-h01}

**⚠️ Should NOT:** Deleted customer should NOT appear in customer search or listings

**Preconditions:**
- Customer account was deleted via full erasure

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| deletedCustomerId | CUST-ERASED | Previously deleted |

**Steps:**
1. Search for deleted customer in admin panel
2. Query customer database

**Expected Result:**
- Customer not found in any listing

**Traceability:**
- AC: [AC-CUST-004](../l1/acceptance-criteria.md#ac-cust-004)
- BR: [BR-CUST](../l1/business-rules.md#br-cust)

---

### TC-AC-CUST-004-H02 – Deleted customer should NOT be able to login {#tc-ac-cust-004-h02}

**⚠️ Should NOT:** Soft-deleted customer should NOT be able to authenticate

**Preconditions:**
- Customer account was soft deleted

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| email | deleted@example.com | Deactivated account |

**Steps:**
1. Attempt login with deleted account credentials

**Expected Result:**
- Login rejected with appropriate message

**Traceability:**
- AC: [AC-CUST-004](../l1/acceptance-criteria.md#ac-cust-004)
- BR: [BR-CUST](../l1/business-rules.md#br-cust)

---

### TC-AC-ADMIN-001-H01 – New product should NOT be visible to customers {#tc-ac-admin-001-h01}

**⚠️ Should NOT:** Draft products should NOT appear in customer product listings

**Preconditions:**
- New product created in draft status

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| productId | PROD-NEW-001 | Draft product |

**Steps:**
1. Create product in draft
2. Search for product in customer-facing catalog

**Expected Result:**
- Product not visible to customers

**Traceability:**
- AC: [AC-ADMIN-001](../l1/acceptance-criteria.md#ac-admin-001)
- BR: [BR-PROD](../l1/business-rules.md#br-prod)

---

### TC-AC-ADMIN-001-H02 – Product should NOT require SKU for creation {#tc-ac-admin-001-h02}

**⚠️ Should NOT:** SKU should NOT be a required field (not in specification)

**Preconditions:**
- Admin is logged in

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| sku | null | SKU not mentioned in requirements |

**Steps:**
1. Create product without SKU field
2. Submit with required fields only

**Expected Result:**
- Product created successfully without SKU

**Traceability:**
- AC: [AC-ADMIN-001](../l1/acceptance-criteria.md#ac-admin-001)
- BR: [BR-PROD](../l1/business-rules.md#br-prod)

---

### TC-AC-ADMIN-002-H01 – Edit should NOT change product UUID {#tc-ac-admin-002-h01}

**⚠️ Should NOT:** Product UUID should NOT change on edit

**Preconditions:**
- Product exists with UUID

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| originalUUID | uuid-12345 | Original product ID |

**Steps:**
1. Record product UUID
2. Edit all product fields
3. Save
4. Check UUID

**Expected Result:**
- UUID remains unchanged

**Traceability:**
- AC: [AC-ADMIN-002](../l1/acceptance-criteria.md#ac-admin-002)
- BR: [BR-PROD](../l1/business-rules.md#br-prod)

---

### TC-AC-ADMIN-002-H02 – Edit should NOT affect completed order history {#tc-ac-admin-002-h02}

**⚠️ Should NOT:** Completed orders should NOT show updated product price

**Preconditions:**
- Product was in completed orders
- Product price changed

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| orderId | ORD-COMPLETED | Historical order |

**Steps:**
1. Change product price
2. View completed order with this product

**Expected Result:**
- Order shows original purchase price

**Traceability:**
- AC: [AC-ADMIN-002](../l1/acceptance-criteria.md#ac-admin-002)
- BR: [BR-PROD](../l1/business-rules.md#br-prod), [BR-ORDER](../l1/business-rules.md#br-order)

---

### TC-AC-ADMIN-003-H01 – Deleted product should NOT be purchasable {#tc-ac-admin-003-h01}

**⚠️ Should NOT:** Deleted product should NOT be addable to cart or purchasable

**Preconditions:**
- Product was deleted

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| productId | PROD-DELETED | Soft-deleted product |

**Steps:**
1. Attempt direct add to cart via API
2. Attempt checkout with cached product

**Expected Result:**
- Add to cart fails
- Checkout blocked

**Traceability:**
- AC: [AC-ADMIN-003](../l1/acceptance-criteria.md#ac-admin-003)
- BR: [BR-PROD](../l1/business-rules.md#br-prod)

---

### TC-AC-ADMIN-003-H02 – Deleted product should NOT be restored automatically {#tc-ac-admin-003-h02}

**⚠️ Should NOT:** Soft-deleted products should NOT auto-restore

**Preconditions:**
- Product was soft-deleted

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| productId | PROD-SOFT-DEL | Soft-deleted product |

**Steps:**
1. Wait for system processes
2. Check product status after time passes

**Expected Result:**
- Product remains soft-deleted

**Traceability:**
- AC: [AC-ADMIN-003](../l1/acceptance-criteria.md#ac-admin-003)
- BR: [BR-PROD](../l1/business-rules.md#br-prod)

---

### TC-AC-ADMIN-004-H01 – Order list should NOT show deleted customer personal data {#tc-ac-admin-004-h01}

**⚠️ Should NOT:** Deleted customer personal data should NOT be visible in order list

**Preconditions:**
- Order placed by now-deleted customer

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| orderId | ORD-ANON | Order from deleted customer |

**Steps:**
1. View order from deleted customer account

**Expected Result:**
- Order shown with anonymized customer info

**Traceability:**
- AC: [AC-ADMIN-004](../l1/acceptance-criteria.md#ac-admin-004)
- BR: [BR-ORDER](../l1/business-rules.md#br-order), [BR-GDPR](../l1/business-rules.md#br-gdpr)

---

### TC-AC-ADMIN-004-H02 – Order list should NOT include customer payment details {#tc-ac-admin-004-h02}

**⚠️ Should NOT:** Payment card numbers should NOT appear in order list view

**Preconditions:**
- Orders exist with payment info

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| orderId | ORD-001 | Order with payment |

**Steps:**
1. View order list
2. Check visible columns

**Expected Result:**
- Only specified columns shown

**Traceability:**
- AC: [AC-ADMIN-004](../l1/acceptance-criteria.md#ac-admin-004)
- BR: [BR-ORDER](../l1/business-rules.md#br-order)

---

### TC-AC-ADMIN-005-H01 – Pending status should NOT trigger email notification {#tc-ac-admin-005-h01}

**⚠️ Should NOT:** Send email notification when order is created or stays pending

**Preconditions:**
- New order created with pending status

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| status | pending | Initial status |

**Steps:**
1. Create new order
2. Verify email queue

**Expected Result:**
- No email notification sent for pending status

**Traceability:**
- AC: [AC-ADMIN-005](../l1/acceptance-criteria.md#ac-admin-005)
- BR: [AMB-OP-048](../l1/business-rules.md#amb-op-048)

---

### TC-AC-ADMIN-005-H02 – Status change should NOT allow custom status values {#tc-ac-admin-005-h02}

**⚠️ Should NOT:** Accept any status value outside the defined four states

**Preconditions:**
- Admin is logged in
- Order exists

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| newStatus | processing | Custom status not in spec |

**Steps:**
1. Attempt to set custom status value

**Expected Result:**
- Custom status rejected

**Traceability:**
- AC: [AC-ADMIN-005](../l1/acceptance-criteria.md#ac-admin-005)
- BR: [AMB-OP-047](../l1/business-rules.md#amb-op-047)

---

### TC-AC-ADMIN-006-H01 – Reason field should NOT be required for adjustments {#tc-ac-admin-006-h01}

**⚠️ Should NOT:** Require reason field for inventory adjustments

**Preconditions:**
- Admin is logged in
- Product variant exists

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| variantId | VAR-008 | Valid variant |
| reason |  | Empty reason field |

**Steps:**
1. Navigate to inventory
2. Adjust stock
3. Leave reason empty
4. Submit

**Expected Result:**
- Adjustment succeeds without reason

**Traceability:**
- AC: [AC-ADMIN-006](../l1/acceptance-criteria.md#ac-admin-006)
- BR: [AMB-ENT-043](../l1/business-rules.md#amb-ent-043)

---

### TC-AC-ADMIN-006-H02 – Inventory should NOT be tracked at product level only {#tc-ac-admin-006-h02}

**⚠️ Should NOT:** Apply inventory changes to product level affecting all variants

**Preconditions:**
- Product has multiple variants

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| variantA | VAR-A | First variant |
| variantB | VAR-B | Second variant |

**Steps:**
1. Adjust stock for variant A
2. Check stock for variant B

**Expected Result:**
- Variant A stock changed
- Variant B stock unchanged

**Traceability:**
- AC: [AC-ADMIN-006](../l1/acceptance-criteria.md#ac-admin-006)
- BR: [AMB-ENT-040](../l1/business-rules.md#amb-ent-040)

---

### TC-AC-VAR-001-H01 – Out of stock variant should NOT be hidden from selection {#tc-ac-var-001-h01}

**⚠️ Should NOT:** Hide or remove out of stock variants from the selection options

**Preconditions:**
- Product exists
- One variant out of stock

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| outOfStockVariant | VAR-M | Out of stock |

**Steps:**
1. Navigate to product page
2. Check variant selector options

**Expected Result:**
- Out of stock variant visible in selector

**Traceability:**
- AC: [AC-VAR-001](../l1/acceptance-criteria.md#ac-var-001)
- BR: [AMB-ENT-009](../l1/business-rules.md#amb-ent-009)

---

### TC-AC-VAR-001-H02 – Variant without price override should NOT show zero price {#tc-ac-var-001-h02}

**⚠️ Should NOT:** Show zero, null, or empty price for variants without override

**Preconditions:**
- Product with base price
- Variant without override

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| basePrice | $50 | Product base price |
| variantOverride | null | No override |

**Steps:**
1. Navigate to product
2. Select variant without override

**Expected Result:**
- Base price $50 displayed

**Traceability:**
- AC: [AC-VAR-001](../l1/acceptance-criteria.md#ac-var-001)
- BR: [AMB-ENT-010](../l1/business-rules.md#amb-ent-010)

---

### TC-AC-CAT-001-H01 – Should NOT display more than 3 levels of hierarchy {#tc-ac-cat-001-h01}

**⚠️ Should NOT:** Display or allow navigation beyond 3 levels of category depth

**Preconditions:**
- Attempt to create 4-level hierarchy in admin

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| level4 | SubSubSubCategory | Fourth level attempt |

**Steps:**
1. Attempt to create category at level 4
2. Check customer view

**Expected Result:**
- Fourth level creation blocked or flattened

**Traceability:**
- AC: [AC-CAT-001](../l1/acceptance-criteria.md#ac-cat-001)
- BR: [AMB-ENT-006](../l1/business-rules.md#amb-ent-006)

---

### TC-AC-CAT-001-H02 – Secondary category should NOT be required {#tc-ac-cat-001-h02}

**⚠️ Should NOT:** Require secondary category assignment for products

**Preconditions:**
- Product exists with primary category only

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| primary | Electronics | Primary set |
| secondary | null | No secondary |

**Steps:**
1. Create product with primary category only
2. Save product

**Expected Result:**
- Product saved successfully
- Displays in primary category

**Traceability:**
- AC: [AC-CAT-001](../l1/acceptance-criteria.md#ac-cat-001)
- BR: [AMB-ENT-038](../l1/business-rules.md#amb-ent-038)

---

