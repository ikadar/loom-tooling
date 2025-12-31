# TDAI Test Cases

Generated: 2025-12-31T12:00:34+01:00

**Methodology:** Test-Driven AI Development (TDAI)

## Summary

| Category | Count | Ratio |
|----------|-------|-------|
| Positive | 72 | 41.9% |
| Negative | 54 | 31.4% |
| Boundary | 22 | - |
| Hallucination Prevention | 24 | - |
| **Total** | **172** | - |

**Coverage:** 40 ACs covered
**Hallucination Prevention:** ✓ Enabled

---

## Positive Tests (Happy Path)

### TC-AC-CUST-001-P01 – Valid registration with all required fields

**Preconditions:**
- Email not registered in system

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| email | john@example.com | Valid email |
| password | SecurePass1 | Valid password |
| firstName | John | Valid name |
| lastName | Doe | Valid name |

**Steps:**
1. Submit registration form with all fields

**Expected Result:**
- Account created with status 'unverified'
- Verification email sent

**Traceability:**
- AC: AC-CUST-001
- BR: [AMB-ENT-024 AMB-ENT-025]

---

### TC-AC-CUST-001-P02 – Valid registration with complex password

**Preconditions:**
- Email not registered in system

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| email | jane@example.com | Valid email |
| password | C0mpl3x!Pass99 | Strong password |
| firstName | Jane | Valid name |
| lastName | Smith | Valid name |

**Steps:**
1. Submit registration form with complex password

**Expected Result:**
- Account created with status 'unverified'
- Verification email sent

**Traceability:**
- AC: AC-CUST-001
- BR: [AMB-ENT-024 AMB-ENT-026]

---

### TC-AC-CUST-002-P01 – Verified customer can place order

**Preconditions:**
- Customer email is verified
- Cart has items

**Steps:**
1. Attempt to place order

**Expected Result:**
- Order placement proceeds

**Traceability:**
- AC: AC-CUST-002
- BR: [AMB-ENT-030]

---

### TC-AC-CUST-002-P02 – Prompt to verify shown for unverified customer

**Preconditions:**
- Customer email is unverified

**Steps:**
1. Attempt to place order

**Expected Result:**
- EMAIL_NOT_VERIFIED error
- Prompt to verify email shown

**Traceability:**
- AC: AC-CUST-002
- BR: [AMB-ENT-030]

---

### TC-AC-CUST-003-P01 – Add valid shipping address

**Preconditions:**
- Customer has one existing address

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| street | 456 Oak Ave | Valid street |
| city | Boston | Valid city |
| state | MA | Valid state |
| postalCode | 02101 | Valid ZIP |
| country | US | Valid country |
| recipient | John Doe | Valid name |

**Steps:**
1. Add new shipping address

**Expected Result:**
- Address saved
- Address selectable for orders

**Traceability:**
- AC: AC-CUST-003
- BR: [AMB-ENT-027 AMB-ENT-031]

---

### TC-AC-CUST-003-P02 – Add second shipping address

**Preconditions:**
- Customer has one existing address

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| street | 789 Elm St | Different address |
| city | Seattle | Different city |
| state | WA | Valid state |
| postalCode | 98101 | Valid ZIP |
| country | US | Valid country |

**Steps:**
1. Add second shipping address

**Expected Result:**
- Both addresses available
- Can select either for orders

**Traceability:**
- AC: AC-CUST-003
- BR: [AMB-ENT-027]

---

### TC-AC-CUST-004-P01 – Successful data erasure request

**Preconditions:**
- Customer registered
- No pending orders

**Steps:**
1. Request account deletion
2. Confirm erasure

**Expected Result:**
- Personal data deleted
- Order history anonymized
- Confirmation sent

**Traceability:**
- AC: AC-CUST-004
- BR: [AMB-ENT-029]

---

### TC-AC-CUST-004-P02 – Order history anonymized after erasure

**Preconditions:**
- Customer with order history
- No pending orders

**Steps:**
1. Complete data erasure
2. Check order records

**Expected Result:**
- Orders exist but anonymized
- No personal data in orders

**Traceability:**
- AC: AC-CUST-004
- BR: [AMB-ENT-029]

---

### TC-AC-PROD-001-P01 – Filter by category, price, and availability

**Preconditions:**
- Products exist in multiple categories

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| category | Electronics | Valid category |
| priceMin | 50 | Min price |
| priceMax | 200 | Max price |
| availability | in stock | In stock only |

**Steps:**
1. Apply all filters

**Expected Result:**
- Only matching products shown
- Sorted by newest
- 20 per page

**Traceability:**
- AC: AC-PROD-001
- BR: [AMB-OP-001 AMB-OP-002]

---

### TC-AC-PROD-001-P02 – Pagination works with 20 products per page

**Preconditions:**
- More than 20 products match filters

**Steps:**
1. Apply filters
2. Navigate to page 2

**Expected Result:**
- Page 2 shows next 20 products
- Proper pagination controls

**Traceability:**
- AC: AC-PROD-001
- BR: [AMB-OP-003]

---

### TC-AC-PROD-002-P01 – Out of stock indicator displayed

**Preconditions:**
- Product has zero stock

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| productName | Widget Pro | Out of stock |

**Steps:**
1. View product listing

**Expected Result:**
- 'Out of Stock' indicator shown
- Add to Cart disabled

**Traceability:**
- AC: AC-PROD-002
- BR: [AMB-OP-004]

---

### TC-AC-PROD-002-P02 – Product details still viewable when out of stock

**Preconditions:**
- Product has zero stock

**Steps:**
1. Click on out of stock product

**Expected Result:**
- Product details page loads
- Description visible

**Traceability:**
- AC: AC-PROD-002
- BR: [AMB-OP-004]

---

### TC-AC-PROD-003-P01 – Variant options displayed correctly

**Preconditions:**
- Product has size and color variants

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| product | T-Shirt | Has variants |
| sizes | S, M, L | Size options |
| colors | Red, Blue | Color options |

**Steps:**
1. View product detail page

**Expected Result:**
- Size options shown
- Color options shown
- Add to Cart disabled

**Traceability:**
- AC: AC-PROD-003
- BR: [AMB-ENT-008 AMB-ENT-009]

---

### TC-AC-PROD-003-P02 – Price and stock update on variant selection

**Preconditions:**
- Product variants have different prices/stock

**Steps:**
1. Select size M and color Red

**Expected Result:**
- Price updates for variant
- Stock updates for variant
- Add to Cart enabled

**Traceability:**
- AC: AC-PROD-003
- BR: [AMB-ENT-010 AMB-ENT-011]

---

### TC-AC-CART-001-P01 – Add product to cart successfully

**Preconditions:**
- Product has 10 units in stock

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| product | Widget | $25.00 price |
| quantity | 2 | Quantity to add |

**Steps:**
1. Add 2 units to cart

**Expected Result:**
- Cart shows Widget qty 2
- Unit price $25.00
- Subtotal $50.00
- Toast with View Cart

**Traceability:**
- AC: AC-CART-001
- BR: [AMB-OP-005 AMB-OP-006]

---

### TC-AC-CART-001-P02 – Toast notification appears with View Cart option

**Preconditions:**
- Product available

**Steps:**
1. Add any product to cart

**Expected Result:**
- Toast notification shown
- 'View Cart' option in toast

**Traceability:**
- AC: AC-CART-001
- BR: [AMB-OP-009]

---

### TC-AC-CART-002-P01 – Quantity increases for existing item

**Preconditions:**
- Cart has Widget qty 2
- Stock is 10

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| additionalQty | 3 | Add 3 more |

**Steps:**
1. Add 3 more Widget to cart

**Expected Result:**
- Cart shows Widget qty 5
- Single line item

**Traceability:**
- AC: AC-CART-002
- BR: [AMB-OP-005 AMB-OP-006]

---

### TC-AC-CART-002-P02 – No duplicate line items created

**Preconditions:**
- Cart has Widget qty 2

**Steps:**
1. Add more of same product

**Expected Result:**
- Still single line item
- Quantity updated

**Traceability:**
- AC: AC-CART-002
- BR: [AMB-OP-005]

---

### TC-AC-CART-003-P01 – Guest cart merged on login

**Preconditions:**
- Guest has items in cart
- Account has existing cart

**Steps:**
1. Add products as guest
2. Log in to account

**Expected Result:**
- Guest items merged with account cart

**Traceability:**
- AC: AC-CART-003
- BR: [AMB-ENT-012 AMB-OP-007]

---

### TC-AC-CART-003-P02 – Guest cart items preserved after login

**Preconditions:**
- Guest has Widget qty 2
- Account cart empty

**Steps:**
1. Log in to account

**Expected Result:**
- Widget qty 2 now in account cart

**Traceability:**
- AC: AC-CART-003
- BR: [AMB-ENT-013]

---

### TC-AC-CART-004-P01 – Update quantity successfully

**Preconditions:**
- Cart has Widget qty 3
- Stock is 10

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| newQuantity | 5 | Valid update |

**Steps:**
1. Change quantity to 5

**Expected Result:**
- Quantity updates to 5
- Subtotal recalculates
- Cart total updates

**Traceability:**
- AC: AC-CART-004
- BR: [AMB-OP-011]

---

### TC-AC-CART-004-P02 – Subtotal recalculates on quantity change

**Preconditions:**
- Widget at $25 qty 3 = $75

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| newQuantity | 5 | Change to 5 |

**Steps:**
1. Update quantity

**Expected Result:**
- Subtotal becomes $125.00

**Traceability:**
- AC: AC-CART-004
- BR: [AMB-OP-011]

---

### TC-AC-CART-005-P01 – Remove item with undo option

**Preconditions:**
- Cart has Widget and Gadget

**Steps:**
1. Remove Widget from cart

**Expected Result:**
- Widget removed immediately
- Undo option for 5 seconds
- Cart total recalculates

**Traceability:**
- AC: AC-CART-005
- BR: [AMB-OP-014]

---

### TC-AC-CART-005-P02 – Undo restores removed item

**Preconditions:**
- Item just removed
- Undo visible

**Steps:**
1. Click Undo within 5 seconds

**Expected Result:**
- Item restored to cart
- Same quantity as before

**Traceability:**
- AC: AC-CART-005
- BR: [AMB-OP-014]

---

### TC-AC-CART-006-P01 – Empty cart state shown

**Preconditions:**
- Cart has only Widget

**Steps:**
1. Remove Widget

**Expected Result:**
- Empty cart state displayed
- 'Continue Shopping' link shown

**Traceability:**
- AC: AC-CART-006
- BR: [AMB-OP-015]

---

### TC-AC-CART-006-P02 – Continue Shopping link works

**Preconditions:**
- Cart is empty

**Steps:**
1. Click Continue Shopping

**Expected Result:**
- Navigates to product catalog

**Traceability:**
- AC: AC-CART-006
- BR: [AMB-OP-015]

---

### TC-AC-CART-007-P01 – Out of stock warning on cart item

**Preconditions:**
- Cart has Widget qty 2
- Widget goes out of stock

**Steps:**
1. View cart after stock depleted

**Expected Result:**
- Item remains in cart
- 'Out of Stock' warning shown

**Traceability:**
- AC: AC-CART-007
- BR: [AMB-ENT-014]

---

### TC-AC-CART-007-P02 – Checkout blocked with out of stock item

**Preconditions:**
- Cart has out of stock item

**Steps:**
1. Attempt checkout

**Expected Result:**
- Checkout blocked
- Prompt to remove or wait

**Traceability:**
- AC: AC-CART-007
- BR: [AMB-ENT-014]

---

### TC-AC-CART-008-P01 – Price change notification shown

**Preconditions:**
- Cart has Widget at $25.00
- Admin changes to $30.00

**Steps:**
1. View cart after price change

**Expected Result:**
- New price $30.00 shown
- 'Price updated from $25 to $30' notification

**Traceability:**
- AC: AC-CART-008
- BR: [AMB-ENT-016]

---

### TC-AC-CART-008-P02 – Cart total reflects new price

**Preconditions:**
- Price changed from $25 to $30
- Qty is 2

**Steps:**
1. View cart totals

**Expected Result:**
- Subtotal shows $60.00 (not $50.00)

**Traceability:**
- AC: AC-CART-008
- BR: [AMB-ENT-016]

---

### TC-AC-CART-009-P01 – Cart cleared after 30 days inactivity

**Preconditions:**
- Logged-in customer cart inactive 30 days

**Steps:**
1. Wait 30 days
2. Log in

**Expected Result:**
- Cart items cleared
- Notification on login

**Traceability:**
- AC: AC-CART-009
- BR: [AMB-ENT-012]

---

### TC-AC-CART-009-P02 – Customer notified of cart expiration

**Preconditions:**
- Cart expired

**Steps:**
1. Log in after expiration

**Expected Result:**
- Clear notification about expired cart

**Traceability:**
- AC: AC-CART-009
- BR: [AMB-ENT-012]

---

### TC-AC-ORDER-001-P01 – Order placed with free shipping over $75

**Preconditions:**
- Verified customer
- Cart subtotal $75+

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| subtotal | $75.00 | Free shipping threshold |

**Steps:**
1. Place order with valid address and credit card

**Expected Result:**
- Order status 'pending'
- Shipping $0.00
- Payment authorized
- Inventory decremented

**Traceability:**
- AC: AC-ORDER-001
- BR: [AMB-OP-016 AMB-OP-020]

---

### TC-AC-ORDER-001-P02 – Cart cleared after successful order

**Preconditions:**
- Order placed successfully

**Steps:**
1. View cart after order

**Expected Result:**
- Cart is empty
- Confirmation email queued

**Traceability:**
- AC: AC-ORDER-001
- BR: [AMB-OP-017]

---

### TC-AC-ORDER-002-P01 – Order includes $5.99 shipping under $75

**Preconditions:**
- Verified customer
- Cart subtotal $35.00

**Steps:**
1. Place order

**Expected Result:**
- Shipping $5.99
- Total $40.99 plus tax

**Traceability:**
- AC: AC-ORDER-002
- BR: [AMB-OP-020 AMB-OP-021]

---

### TC-AC-ORDER-002-P02 – Tax calculated on order total

**Preconditions:**
- Subtotal + shipping calculated

**Steps:**
1. View order total

**Expected Result:**
- Tax added to subtotal
- Final total displayed

**Traceability:**
- AC: AC-ORDER-002
- BR: [AMB-OP-021]

---

### TC-AC-ORDER-003-P01 – Sequential order number assigned

**Preconditions:**
- Last order was ORD-2024-00042

**Steps:**
1. Place new order

**Expected Result:**
- Order number ORD-2024-00043 assigned

**Traceability:**
- AC: AC-ORDER-003
- BR: [AMB-ENT-017]

---

### TC-AC-ORDER-003-P02 – Order number follows year format

**Preconditions:**
- Current year is 2024

**Steps:**
1. Place order

**Expected Result:**
- Format is ORD-2024-XXXXX

**Traceability:**
- AC: AC-ORDER-003
- BR: [AMB-ENT-017]

---

### TC-AC-ORDER-004-P01 – Order captures price at order time

**Preconditions:**
- Widget price is $25.00

**Steps:**
1. Place order
2. Admin changes price to $30

**Expected Result:**
- Order line item shows $25.00

**Traceability:**
- AC: AC-ORDER-004
- BR: [AMB-ENT-019]

---

### TC-AC-ORDER-004-P02 – Order total unaffected by later price changes

**Preconditions:**
- Order placed at $25 per unit

**Steps:**
1. View order after product price change

**Expected Result:**
- Original price preserved in order

**Traceability:**
- AC: AC-ORDER-004
- BR: [AMB-ENT-019]

---

### TC-AC-ORDER-005-P01 – Tax calculated correctly for California

**Preconditions:**
- Subtotal $100.00
- Shipping to CA (7.25%)

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| subtotal | $100.00 | Before tax |
| state | CA | 7.25% rate |

**Steps:**
1. Calculate order totals

**Expected Result:**
- Tax is $7.25

**Traceability:**
- AC: AC-ORDER-005
- BR: [AMB-OP-022]

---

### TC-AC-ORDER-005-P02 – Tax calculated on subtotal before shipping

**Preconditions:**
- Subtotal $50
- Shipping $5.99

**Steps:**
1. View tax amount

**Expected Result:**
- Tax based on $50 only

**Traceability:**
- AC: AC-ORDER-005
- BR: [AMB-OP-022]

---

### TC-AC-ORDER-006-P01 – Order history shows all orders

**Preconditions:**
- Customer has multiple orders

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| order1 | ORD-2024-00001 | Delivered |
| order2 | ORD-2024-00042 | Shipped |

**Steps:**
1. View order history

**Expected Result:**
- All orders displayed
- Shows number, date, items, total, status

**Traceability:**
- AC: AC-ORDER-006
- BR: [AMB-OP-025 AMB-OP-026]

---

### TC-AC-ORDER-006-P02 – Order history has pagination

**Preconditions:**
- Customer has 50+ orders

**Steps:**
1. View order history

**Expected Result:**
- Pagination controls shown
- Can navigate pages

**Traceability:**
- AC: AC-ORDER-006
- BR: [AMB-OP-026]

---

### TC-AC-ORDER-007-P01 – Tracking number displayed for shipped order

**Preconditions:**
- Order status 'shipped'
- Has tracking number

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| orderNumber | ORD-2024-00042 | Shipped |
| tracking | 1Z999AA10123456784 | UPS |

**Steps:**
1. View order details

**Expected Result:**
- Tracking number displayed
- Link to carrier tracking page

**Traceability:**
- AC: AC-ORDER-007
- BR: [AMB-OP-027]

---

### TC-AC-ORDER-007-P02 – Carrier tracking link is clickable

**Preconditions:**
- Shipped order with tracking

**Steps:**
1. Click tracking link

**Expected Result:**
- Opens carrier tracking page

**Traceability:**
- AC: AC-ORDER-007
- BR: [AMB-OP-027]

---

### TC-AC-ORDER-008-P01 – Successfully cancel pending order

**Preconditions:**
- Order status 'pending'

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| reason | Changed mind | Cancellation reason |

**Steps:**
1. Request cancellation with reason

**Expected Result:**
- Status changes to 'cancelled'
- Inventory restored
- Refund initiated
- Email sent

**Traceability:**
- AC: AC-ORDER-008
- BR: [AMB-OP-028 AMB-OP-029]

---

### TC-AC-ORDER-008-P02 – Refund to original payment method

**Preconditions:**
- Order cancelled

**Steps:**
1. Check refund details

**Expected Result:**
- Refund to original card
- Full amount refunded

**Traceability:**
- AC: AC-ORDER-008
- BR: [AMB-OP-030 AMB-OP-031]

---

### TC-AC-ORDER-009-P01 – Cancel confirmed order successfully

**Preconditions:**
- Order status 'confirmed'

**Steps:**
1. Request cancellation

**Expected Result:**
- Status changes to 'cancelled'
- Inventory restored
- Refund initiated

**Traceability:**
- AC: AC-ORDER-009
- BR: [AMB-OP-028]

---

### TC-AC-ORDER-009-P02 – Inventory restored on confirmed order cancel

**Preconditions:**
- Confirmed order for 5 units
- Stock was 10 (now 5)

**Steps:**
1. Cancel order
2. Check inventory

**Expected Result:**
- Stock returns to 10

**Traceability:**
- AC: AC-ORDER-009
- BR: [AMB-OP-029]

---

### TC-AC-ORDER-010-P01 – Shipped order shows cancel disabled

**Preconditions:**
- Order status 'shipped'

**Steps:**
1. View order details

**Expected Result:**
- Cancel button disabled or hidden

**Traceability:**
- AC: AC-ORDER-010
- BR: [AMB-OP-028 AMB-ENT-018]

---

### TC-AC-ORDER-010-P02 – Return process suggested for shipped orders

**Preconditions:**
- Order status 'shipped'

**Steps:**
1. Attempt to cancel

**Expected Result:**
- Return/refund instructions shown instead

**Traceability:**
- AC: AC-ORDER-010

---

### TC-AC-ORDER-011-P01 – Confirmation email with order details

**Preconditions:**
- Order just placed successfully

**Steps:**
1. Check email queue

**Expected Result:**
- Email queued with order number, items, quantities, prices, total

**Traceability:**
- AC: AC-ORDER-011
- BR: [AMB-OP-056 AMB-OP-057]

---

### TC-AC-ORDER-011-P02 – Email includes shipping and delivery info

**Preconditions:**
- Order placed

**Steps:**
1. View email content

**Expected Result:**
- Shipping address shown
- Estimated delivery date included

**Traceability:**
- AC: AC-ORDER-011
- BR: [AMB-OP-058 AMB-OP-059]

---

### TC-AC-ADMIN-001-P01 – Create product with required fields

**Preconditions:**
- Admin on product management page

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| name | Super Widget | Valid name |
| description | A great widget | Valid desc |
| price | 29.99 | Valid price |
| category | Electronics | Valid category |

**Steps:**
1. Create product with all required fields

**Expected Result:**
- Product created with status 'draft'
- UUID assigned
- Audit logged

**Traceability:**
- AC: AC-ADMIN-001
- BR: [AMB-ENT-001 AMB-OP-033]

---

### TC-AC-ADMIN-001-P02 – Audit log includes creator and timestamp

**Preconditions:**
- Product creation

**Steps:**
1. Check audit log

**Expected Result:**
- Creator username logged
- Timestamp logged

**Traceability:**
- AC: AC-ADMIN-001
- BR: [AMB-OP-037]

---

### TC-AC-ADMIN-002-P01 – Upload images and set primary

**Preconditions:**
- Product in draft status

**Steps:**
1. Upload 3 images
2. Set one as primary

**Expected Result:**
- Images uploaded to cloud
- URLs stored
- Primary image marked

**Traceability:**
- AC: AC-ADMIN-002
- BR: [AMB-ENT-004 AMB-OP-036]

---

### TC-AC-ADMIN-002-P02 – Multiple images displayed in order

**Preconditions:**
- Product has 3 images

**Steps:**
1. View product images

**Expected Result:**
- Primary shown first
- All images accessible

**Traceability:**
- AC: AC-ADMIN-002
- BR: [AMB-ENT-004]

---

### TC-AC-ADMIN-003-P01 – Update price with audit log

**Preconditions:**
- Product price is $25.00

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| newPrice | 30.00 | New price |

**Steps:**
1. Change price to $30.00

**Expected Result:**
- Price updated
- Change logged with admin and timestamp

**Traceability:**
- AC: AC-ADMIN-003
- BR: [AMB-ENT-007 AMB-OP-038]

---

### TC-AC-ADMIN-003-P02 – Cart prices update after product price change

**Preconditions:**
- Product in customer carts at $25

**Steps:**
1. Admin changes price to $30
2. Customer views cart

**Expected Result:**
- Cart shows $30 on next view

**Traceability:**
- AC: AC-ADMIN-003
- BR: [AMB-OP-040]

---

### TC-AC-ADMIN-004-P01 – Last save wins with conflict warning

**Preconditions:**
- Admin A and B editing same product

**Steps:**
1. Admin A saves
2. Admin B saves different changes

**Expected Result:**
- Admin B's changes saved
- Conflict warning shown to B

**Traceability:**
- AC: AC-ADMIN-004
- BR: [AMB-OP-039]

---

### TC-AC-ADMIN-004-P02 – Warning shows overwritten fields

**Preconditions:**
- Concurrent edit conflict

**Steps:**
1. Second admin saves

**Expected Result:**
- Warning indicates which fields were overwritten

**Traceability:**
- AC: AC-ADMIN-004
- BR: [AMB-OP-039]

---

### TC-AC-ADMIN-005-P01 – Soft delete hides from customers

**Preconditions:**
- Product has order history

**Steps:**
1. Admin confirms removal

**Expected Result:**
- Product soft-deleted
- Not visible to customers

**Traceability:**
- AC: AC-ADMIN-005
- BR: [AMB-ENT-005 AMB-OP-041]

---

### TC-AC-ADMIN-005-P02 – Order history preserved after removal

**Preconditions:**
- Product in past orders

**Steps:**
1. Remove product
2. View order history

**Expected Result:**
- Order history still shows product details

**Traceability:**
- AC: AC-ADMIN-005
- BR: [AMB-OP-042]

---

### TC-AC-ADMIN-005-P03 – Cart items removed with notification

**Preconditions:**
- Product in customer carts

**Steps:**
1. Admin removes product

**Expected Result:**
- Cart items removed
- Customer notification on next view

**Traceability:**
- AC: AC-ADMIN-005
- BR: [AMB-OP-043]

---

### TC-AC-ADMIN-006-P01 – Filter orders by status and date

**Preconditions:**
- 100 orders in system

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| status | shipped | Filter value |
| dateRange | last 7 days | Date filter |

**Steps:**
1. Apply status and date filters

**Expected Result:**
- Only matching orders shown
- Displays order number, date, customer, total, status

**Traceability:**
- AC: AC-ADMIN-006
- BR: [AMB-OP-044 AMB-OP-045]

---

### TC-AC-ADMIN-006-P02 – Search by order number

**Preconditions:**
- Orders exist

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| search | ORD-2024 | Partial number |

**Steps:**
1. Search for ORD-2024

**Expected Result:**
- All 2024 orders shown

**Traceability:**
- AC: AC-ADMIN-006
- BR: [AMB-OP-044]

---

### TC-AC-ADMIN-007-P01 – Export filtered orders to CSV

**Preconditions:**
- Filtered list shows 25 orders

**Steps:**
1. Click Export CSV

**Expected Result:**
- CSV file downloads
- Contains all 25 filtered orders

**Traceability:**
- AC: AC-ADMIN-007
- BR: [AMB-OP-046]

---

### TC-AC-ADMIN-007-P02 – CSV includes order details

**Preconditions:**
- Export completed

**Steps:**
1. Open CSV file

**Expected Result:**
- Order number, date, customer, items, totals included

**Traceability:**
- AC: AC-ADMIN-007
- BR: [AMB-OP-046]

---

### TC-AC-ADMIN-008-P01 – Update to shipped with tracking

**Preconditions:**
- Order status 'confirmed'

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| status | shipped | New status |
| tracking | 1Z999AA10123456784 | Tracking number |

**Steps:**
1. Update status with tracking number

**Expected Result:**
- Status changes to 'shipped'
- Tracking stored
- Change logged

**Traceability:**
- AC: AC-ADMIN-008
- BR: [AMB-OP-047 AMB-OP-048]

---

### TC-AC-ADMIN-008-P02 – Customer email sent on shipment

**Preconditions:**
- Order marked shipped

**Steps:**
1. Check email queue

**Expected Result:**
- Shipment notification email queued

**Traceability:**
- AC: AC-ADMIN-008
- BR: [AMB-OP-050]

---

### TC-AC-ADMIN-009-P01 – Set stock with audit log

**Preconditions:**
- Product stock is 50

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| newStock | 100 | New quantity |
| reason | Shipment received | Reason |

**Steps:**
1. Set stock to 100 with reason

**Expected Result:**
- Stock updated to 100
- Audit log: before 50, after 100, reason, admin, timestamp

**Traceability:**
- AC: AC-ADMIN-009
- BR: [AMB-OP-051 AMB-OP-052]

---

### TC-AC-ADMIN-009-P02 – Audit log includes all details

**Preconditions:**
- Stock updated

**Steps:**
1. View audit log

**Expected Result:**
- Before value, after value, reason, admin, timestamp all logged

**Traceability:**
- AC: AC-ADMIN-009
- BR: [AMB-OP-055]

---

### TC-AC-ADMIN-010-P01 – Adjust stock by negative delta

**Preconditions:**
- Stock is 50

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| delta | -5 | Reduce by 5 |
| reason | Damaged goods | Reason |

**Steps:**
1. Adjust stock by -5

**Expected Result:**
- Stock updated to 45
- Audit log created

**Traceability:**
- AC: AC-ADMIN-010
- BR: [AMB-OP-051 AMB-OP-053]

---

### TC-AC-ADMIN-010-P02 – Adjust stock by positive delta

**Preconditions:**
- Stock is 50

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| delta | +10 | Add 10 |

**Steps:**
1. Adjust stock by +10

**Expected Result:**
- Stock updated to 60

**Traceability:**
- AC: AC-ADMIN-010
- BR: [AMB-OP-051]

---

### TC-AC-ADMIN-011-P01 – Bulk update from valid CSV

**Preconditions:**
- CSV with 50 valid SKUs and quantities

**Steps:**
1. Import CSV file

**Expected Result:**
- All 50 rows updated
- Summary shown

**Traceability:**
- AC: AC-ADMIN-011
- BR: [AMB-OP-054]

---

### TC-AC-ADMIN-011-P02 – Summary shows success count

**Preconditions:**
- CSV imported

**Steps:**
1. View import results

**Expected Result:**
- Success count displayed
- Error count if any

**Traceability:**
- AC: AC-ADMIN-011
- BR: [AMB-OP-054]

---

### TC-AC-ADMIN-012-P01 – Alert triggered when stock drops below threshold

**Preconditions:**
- Threshold is 10
- Stock is 15

**Steps:**
1. Stock drops to 8 (via order or adjustment)

**Expected Result:**
- Low stock alert triggered

**Traceability:**
- AC: AC-ADMIN-012
- BR: [AMB-ENT-041]

---

### TC-AC-ADMIN-012-P02 – Admin notified of low stock

**Preconditions:**
- Low stock alert triggered

**Steps:**
1. Check admin notifications

**Expected Result:**
- Low stock notification visible

**Traceability:**
- AC: AC-ADMIN-012
- BR: [AMB-ENT-041]

---

### TC-AC-ADMIN-013-P01 – Create subcategory under parent

**Preconditions:**
- Category 'Electronics' exists

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| name | Smartphones | Subcategory |
| parent | Electronics | Parent category |

**Steps:**
1. Create subcategory Smartphones under Electronics

**Expected Result:**
- Subcategory created
- Parent reference set

**Traceability:**
- AC: AC-ADMIN-013
- BR: [AMB-ENT-037]

---

### TC-AC-ADMIN-013-P02 – Display order configurable

**Preconditions:**
- Subcategory created

**Steps:**
1. Set display order

**Expected Result:**
- Order saved
- Categories display in configured order

**Traceability:**
- AC: AC-ADMIN-013
- BR: [AMB-ENT-038]

---

### TC-AC-PAY-001-P01 – Successful Visa payment via Stripe

**Preconditions:**
- Customer at checkout
- Valid Visa card

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| cardLast4 | 4242 | Visa test card |

**Steps:**
1. Submit payment

**Expected Result:**
- Payment authorized via Stripe
- Tokenized reference stored
- Order created

**Traceability:**
- AC: AC-PAY-001
- BR: [AMB-ENT-034 AMB-ENT-036]

---

### TC-AC-PAY-001-P02 – Mastercard payment successful

**Preconditions:**
- Valid Mastercard

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| cardLast4 | 5555 | Mastercard |

**Steps:**
1. Submit payment

**Expected Result:**
- Payment authorized
- Order created

**Traceability:**
- AC: AC-PAY-001
- BR: [AMB-ENT-034]

---

### TC-AC-PAY-002-P01 – Successful PayPal authorization

**Preconditions:**
- Customer selects PayPal at checkout

**Steps:**
1. Complete PayPal authorization

**Expected Result:**
- Payment authorized
- PayPal reference stored
- Order created

**Traceability:**
- AC: AC-PAY-002
- BR: [AMB-ENT-035]

---

### TC-AC-PAY-002-P02 – PayPal reference stored for order

**Preconditions:**
- PayPal payment completed

**Steps:**
1. View order details

**Expected Result:**
- PayPal transaction reference visible

**Traceability:**
- AC: AC-PAY-002
- BR: [AMB-ENT-035]

---

### TC-AC-PAY-003-P01 – Save card for future purchases

**Preconditions:**
- Customer completing order

**Steps:**
1. Select 'Save for future purchases'
2. Complete order

**Expected Result:**
- Tokenized payment method stored
- Available for future checkouts

**Traceability:**
- AC: AC-PAY-003
- BR: [AMB-ENT-028 AMB-ENT-036]

---

### TC-AC-PAY-003-P02 – Saved card available at next checkout

**Preconditions:**
- Customer has saved card

**Steps:**
1. Start new checkout

**Expected Result:**
- Saved card shown as option
- Can select for payment

**Traceability:**
- AC: AC-PAY-003
- BR: [AMB-ENT-028]

---

## Negative Tests (Error Cases)

### TC-AC-CUST-001-N01 – Reject duplicate email registration

**Preconditions:**
- Email john@example.com already registered

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| email | john@example.com | Existing email |

**Steps:**
1. Submit registration with existing email

**Expected Result:**
- DUPLICATE_EMAIL error
- Suggest login or password reset

**Traceability:**
- AC: AC-CUST-001
- BR: [AMB-ENT-024]

---

### TC-AC-CUST-001-N02 – Reject password under 8 characters

**Preconditions:**
- Email not registered

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| password | Pass1 | Only 5 characters |

**Steps:**
1. Submit registration with short password

**Expected Result:**
- INVALID_PASSWORD error
- Requirements message shown

**Traceability:**
- AC: AC-CUST-001
- BR: [AMB-ENT-025]

---

### TC-AC-CUST-001-N03 – Reject password without number

**Preconditions:**
- Email not registered

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| password | SecurePass | No numbers |

**Steps:**
1. Submit registration with password lacking number

**Expected Result:**
- INVALID_PASSWORD error
- Requirements message shown

**Traceability:**
- AC: AC-CUST-001
- BR: [AMB-ENT-025]

---

### TC-AC-CUST-001-N04 – Reject invalid email format

**Preconditions:**
- User on registration page

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| email | invalid-email | Missing @ and domain |

**Steps:**
1. Submit registration with invalid email

**Expected Result:**
- INVALID_EMAIL error

**Traceability:**
- AC: AC-CUST-001
- BR: [AMB-ENT-024]

---

### TC-AC-CUST-002-N01 – Unverified customer order rejected

**Preconditions:**
- Customer email is unverified
- Cart has items

**Steps:**
1. Attempt to place order

**Expected Result:**
- Order rejected
- EMAIL_NOT_VERIFIED error

**Traceability:**
- AC: AC-CUST-002
- BR: [AMB-ENT-030]

---

### TC-AC-CUST-002-N02 – Expired verification link rejected

**Preconditions:**
- Verification link is expired

**Steps:**
1. Click expired verification link

**Expected Result:**
- VERIFICATION_EXPIRED error
- Resend option shown

**Traceability:**
- AC: AC-CUST-002

---

### TC-AC-CUST-003-N01 – Reject address with missing required field

**Preconditions:**
- Customer on address form

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| street |  | Missing required field |

**Steps:**
1. Submit address without street

**Expected Result:**
- INVALID_ADDRESS error
- Missing field specified

**Traceability:**
- AC: AC-CUST-003
- BR: [AMB-ENT-031]

---

### TC-AC-CUST-003-N02 – Reject invalid postal code format

**Preconditions:**
- Customer on address form

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| postalCode | ABCDE | Invalid format |

**Steps:**
1. Submit address with invalid postal code

**Expected Result:**
- INVALID_POSTAL_CODE error

**Traceability:**
- AC: AC-CUST-003
- BR: [AMB-ENT-032]

---

### TC-AC-CUST-004-N01 – Reject erasure with pending orders

**Preconditions:**
- Customer has pending order

**Steps:**
1. Request data erasure

**Expected Result:**
- PENDING_ORDERS_EXIST error
- Must wait for completion

**Traceability:**
- AC: AC-CUST-004

---

### TC-AC-CUST-004-N02 – Reject erasure with shipped but not delivered order

**Preconditions:**
- Customer has shipped order (not delivered)

**Steps:**
1. Request data erasure

**Expected Result:**
- PENDING_ORDERS_EXIST error

**Traceability:**
- AC: AC-CUST-004

---

### TC-AC-PROD-001-N01 – Empty result with no matching products

**Preconditions:**
- No products match filter criteria

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| category | Electronics | Valid |
| priceMax | 1 | Unrealistic price |

**Steps:**
1. Apply filters with no matches

**Expected Result:**
- 'No products found' message
- Clear filters option shown

**Traceability:**
- AC: AC-PROD-001

---

### TC-AC-PROD-001-N02 – Invalid price range handled

**Preconditions:**
- User on browse page

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| priceMin | 200 | Higher than max |
| priceMax | 50 | Lower than min |

**Steps:**
1. Apply inverted price range

**Expected Result:**
- Error or auto-correction
- No results returned

**Traceability:**
- AC: AC-PROD-001

---

### TC-AC-PROD-002-N01 – Cannot add out of stock product to cart

**Preconditions:**
- Product has zero stock

**Steps:**
1. Attempt to click disabled Add to Cart

**Expected Result:**
- Button is disabled
- No cart action occurs

**Traceability:**
- AC: AC-PROD-002
- BR: [AMB-OP-004]

---

### TC-AC-PROD-002-N02 – API rejects adding out of stock via direct call

**Preconditions:**
- Product has zero stock

**Steps:**
1. Call add-to-cart API directly

**Expected Result:**
- OUT_OF_STOCK error returned

**Traceability:**
- AC: AC-PROD-002
- BR: [AMB-OP-004]

---

### TC-AC-PROD-003-N01 – Selected variant out of stock

**Preconditions:**
- Specific variant (M, Red) is out of stock

**Steps:**
1. Select out of stock variant

**Expected Result:**
- 'Out of Stock' for variant
- Can select different variant

**Traceability:**
- AC: AC-PROD-003
- BR: [AMB-OP-010]

---

### TC-AC-PROD-003-N02 – Cannot add to cart without variant selection

**Preconditions:**
- No variant selected

**Steps:**
1. Attempt to add to cart without selection

**Expected Result:**
- Add to Cart is disabled
- Prompt to select variant

**Traceability:**
- AC: AC-PROD-003
- BR: [AMB-OP-010]

---

### TC-AC-CART-001-N01 – Cannot add out of stock product

**Preconditions:**
- Product stock is 0

**Steps:**
1. Attempt to add to cart

**Expected Result:**
- OUT_OF_STOCK error
- Product not added

**Traceability:**
- AC: AC-CART-001
- BR: [AMB-OP-006]

---

### TC-AC-CART-001-N02 – Cannot add more than available stock

**Preconditions:**
- Product has 5 units in stock

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| quantity | 10 | Exceeds stock |

**Steps:**
1. Attempt to add 10 units

**Expected Result:**
- INSUFFICIENT_STOCK error
- Available quantity shown

**Traceability:**
- AC: AC-CART-001
- BR: [AMB-OP-006]

---

### TC-AC-CART-002-N01 – Total quantity limited to available stock

**Preconditions:**
- Cart has Widget qty 8
- Stock is 10

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| additionalQty | 5 | Would exceed stock |

**Steps:**
1. Attempt to add 5 more

**Expected Result:**
- Limited to 10 total
- Notification shown

**Traceability:**
- AC: AC-CART-002
- BR: [AMB-OP-006]

---

### TC-AC-CART-003-N01 – Merged quantity limited to stock

**Preconditions:**
- Guest cart has 8 units
- Account cart has 5
- Stock is 10

**Steps:**
1. Log in to account

**Expected Result:**
- Total limited to 10
- Notification about limit

**Traceability:**
- AC: AC-CART-003
- BR: [AMB-OP-007]

---

### TC-AC-CART-004-N01 – Setting quantity to 0 removes item

**Preconditions:**
- Cart has Widget qty 3

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| newQuantity | 0 | Zero quantity |

**Steps:**
1. Set quantity to 0

**Expected Result:**
- Widget removed from cart

**Traceability:**
- AC: AC-CART-004
- BR: [AMB-OP-012]

---

### TC-AC-CART-004-N02 – Quantity exceeding stock limited

**Preconditions:**
- Cart has Widget qty 3
- Stock is 10

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| newQuantity | 15 | Exceeds stock |

**Steps:**
1. Attempt to set quantity to 15

**Expected Result:**
- Limited to 10
- 'Only 10 available' notification

**Traceability:**
- AC: AC-CART-004
- BR: [AMB-OP-013]

---

### TC-AC-CART-005-N01 – Undo not available after 5 seconds

**Preconditions:**
- Item removed

**Steps:**
1. Wait more than 5 seconds

**Expected Result:**
- Undo option disappears
- Removal is permanent

**Traceability:**
- AC: AC-CART-005
- BR: [AMB-OP-014]

---

### TC-AC-CART-006-N01 – Cannot proceed to checkout with empty cart

**Preconditions:**
- Cart is empty

**Steps:**
1. Attempt to proceed to checkout

**Expected Result:**
- Checkout blocked
- Message to add items

**Traceability:**
- AC: AC-CART-006

---

### TC-AC-CART-007-N01 – Cannot increase quantity of out of stock item

**Preconditions:**
- Cart item is out of stock

**Steps:**
1. Attempt to increase quantity

**Expected Result:**
- Increase blocked
- Still at 0 available

**Traceability:**
- AC: AC-CART-007
- BR: [AMB-ENT-014]

---

### TC-AC-CART-008-N01 – Price decrease also notified

**Preconditions:**
- Cart has Widget at $25.00
- Price drops to $20.00

**Steps:**
1. View cart

**Expected Result:**
- New price $20.00
- Notification of price change

**Traceability:**
- AC: AC-CART-008
- BR: [AMB-ENT-016]

---

### TC-AC-CART-009-N01 – Cart activity resets expiration timer

**Preconditions:**
- Cart has items
- 29 days of inactivity

**Steps:**
1. Update cart quantity
2. Wait another 29 days

**Expected Result:**
- Cart still active
- Timer reset by activity

**Traceability:**
- AC: AC-CART-009
- BR: [AMB-ENT-012]

---

### TC-AC-ORDER-001-N01 – Order rejected on payment failure

**Preconditions:**
- Invalid credit card

**Steps:**
1. Attempt to place order

**Expected Result:**
- PAYMENT_FAILED error
- Cart preserved
- Retry allowed

**Traceability:**
- AC: AC-ORDER-001
- BR: [AMB-OP-018]

---

### TC-AC-ORDER-001-N02 – Order rejected when item out of stock

**Preconditions:**
- Item goes out of stock at checkout

**Steps:**
1. Attempt to place order

**Expected Result:**
- STOCK_UNAVAILABLE error
- Modification allowed

**Traceability:**
- AC: AC-ORDER-001
- BR: [AMB-OP-019]

---

### TC-AC-ORDER-001-N03 – Order rejected for unverified customer

**Preconditions:**
- Customer email not verified

**Steps:**
1. Attempt to place order

**Expected Result:**
- EMAIL_NOT_VERIFIED error

**Traceability:**
- AC: AC-ORDER-001
- BR: [AMB-OP-016]

---

### TC-AC-ORDER-002-N01 – Shipping cannot be waived below threshold

**Preconditions:**
- Subtotal $74.99

**Steps:**
1. View shipping cost

**Expected Result:**
- Shipping is $5.99
- Not free

**Traceability:**
- AC: AC-ORDER-002
- BR: [AMB-OP-020]

---

### TC-AC-ORDER-003-N01 – Concurrent orders get unique numbers

**Preconditions:**
- Multiple orders placed simultaneously

**Steps:**
1. Place 3 orders at same time

**Expected Result:**
- Each gets unique sequential number

**Traceability:**
- AC: AC-ORDER-003
- BR: [AMB-ENT-017]

---

### TC-AC-ORDER-004-N01 – Price change does not affect existing orders

**Preconditions:**
- Order exists with price $25

**Steps:**
1. Admin changes price
2. View existing order

**Expected Result:**
- Order still shows $25

**Traceability:**
- AC: AC-ORDER-004
- BR: [AMB-ENT-019]

---

### TC-AC-ORDER-005-N01 – Zero tax for tax-exempt states

**Preconditions:**
- Shipping to tax-exempt location

**Steps:**
1. Calculate order totals

**Expected Result:**
- Tax is $0.00

**Traceability:**
- AC: AC-ORDER-005
- BR: [AMB-OP-022]

---

### TC-AC-ORDER-006-N01 – Empty order history for new customer

**Preconditions:**
- New customer with no orders

**Steps:**
1. View order history

**Expected Result:**
- Empty state message
- Link to start shopping

**Traceability:**
- AC: AC-ORDER-006
- BR: [AMB-OP-025]

---

### TC-AC-ORDER-007-N01 – No tracking shown for pending order

**Preconditions:**
- Order status 'pending'

**Steps:**
1. View order details

**Expected Result:**
- No tracking number field
- Status shows pending

**Traceability:**
- AC: AC-ORDER-007
- BR: [AMB-OP-027]

---

### TC-AC-ORDER-008-N01 – Cancellation reason required

**Preconditions:**
- Order status 'pending'

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| reason |  | Empty reason |

**Steps:**
1. Submit cancellation without reason

**Expected Result:**
- Reason required error or prompt

**Traceability:**
- AC: AC-ORDER-008
- BR: [AMB-OP-032]

---

### TC-AC-ORDER-009-N01 – Cannot cancel after status advances beyond confirmed

**Preconditions:**
- Order already shipped

**Steps:**
1. Attempt cancellation

**Expected Result:**
- Cancellation rejected

**Traceability:**
- AC: AC-ORDER-009
- BR: [AMB-OP-028]

---

### TC-AC-ORDER-010-N01 – Cancellation rejected for shipped order

**Preconditions:**
- Order status 'shipped'

**Steps:**
1. Request cancellation via API

**Expected Result:**
- ORDER_ALREADY_SHIPPED error

**Traceability:**
- AC: AC-ORDER-010
- BR: [AMB-OP-028]

---

### TC-AC-ORDER-010-N02 – Cancellation rejected for delivered order

**Preconditions:**
- Order status 'delivered'

**Steps:**
1. Request cancellation

**Expected Result:**
- Cancellation not allowed

**Traceability:**
- AC: AC-ORDER-010
- BR: [AMB-ENT-018]

---

### TC-AC-ORDER-011-N01 – Email failure does not block order

**Preconditions:**
- Email service temporarily down

**Steps:**
1. Place order

**Expected Result:**
- Order succeeds
- Email queued for retry

**Traceability:**
- AC: AC-ORDER-011
- BR: [AMB-OP-060]

---

### TC-AC-ADMIN-001-N01 – Reject product with empty name

**Preconditions:**
- Admin creating product

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| name |  | Empty |

**Steps:**
1. Submit with empty name

**Expected Result:**
- NAME_REQUIRED error

**Traceability:**
- AC: AC-ADMIN-001
- BR: [AMB-OP-034]

---

### TC-AC-ADMIN-001-N02 – Reject name over 200 characters

**Preconditions:**
- Admin creating product

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| name | [201 chars] | Too long |

**Steps:**
1. Submit with long name

**Expected Result:**
- NAME_TOO_LONG error

**Traceability:**
- AC: AC-ADMIN-001
- BR: [AMB-OP-034]

---

### TC-AC-ADMIN-001-N03 – Reject zero or negative price

**Preconditions:**
- Admin creating product

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| price | 0 | Zero price |

**Steps:**
1. Submit with zero price

**Expected Result:**
- INVALID_PRICE error

**Traceability:**
- AC: AC-ADMIN-001
- BR: [AMB-OP-035]

---

### TC-AC-ADMIN-001-N04 – Reject without category

**Preconditions:**
- Admin creating product

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| category |  | No category |

**Steps:**
1. Submit without category

**Expected Result:**
- CATEGORY_REQUIRED error

**Traceability:**
- AC: AC-ADMIN-001
- BR: [AMB-OP-035]

---

### TC-AC-ADMIN-001-N05 – Duplicate name shows warning (not blocking)

**Preconditions:**
- Product named 'Super Widget' exists

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| name | Super Widget | Duplicate |

**Steps:**
1. Create product with same name

**Expected Result:**
- Warning shown
- Creation allowed

**Traceability:**
- AC: AC-ADMIN-001
- BR: [AMB-ENT-002]

---

### TC-AC-ADMIN-002-N01 – Reject more than 10 images

**Preconditions:**
- Product already has 10 images

**Steps:**
1. Attempt to upload 11th image

**Expected Result:**
- MAX_IMAGES_EXCEEDED error

**Traceability:**
- AC: AC-ADMIN-002
- BR: [AMB-OP-036]

---

### TC-AC-ADMIN-002-N02 – Reject invalid image format

**Preconditions:**
- Admin uploading image

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| file | document.pdf | Not an image |

**Steps:**
1. Upload PDF file

**Expected Result:**
- INVALID_IMAGE_FORMAT error

**Traceability:**
- AC: AC-ADMIN-002
- BR: [AMB-OP-036]

---

### TC-AC-ADMIN-003-N01 – Reject negative price update

**Preconditions:**
- Product exists

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| newPrice | -5.00 | Negative |

**Steps:**
1. Attempt negative price

**Expected Result:**
- INVALID_PRICE error

**Traceability:**
- AC: AC-ADMIN-003
- BR: [AMB-OP-038]

---

### TC-AC-ADMIN-004-N01 – First admin's changes lost

**Preconditions:**
- Both admins edit

**Steps:**
1. A saves
2. B saves

**Expected Result:**
- A's specific changes overwritten

**Traceability:**
- AC: AC-ADMIN-004
- BR: [AMB-OP-039]

---

### TC-AC-ADMIN-005-N01 – Cannot permanently delete product with orders

**Preconditions:**
- Product has order history

**Steps:**
1. Attempt hard delete

**Expected Result:**
- Only soft delete allowed

**Traceability:**
- AC: AC-ADMIN-005
- BR: [AMB-OP-042]

---

### TC-AC-ADMIN-006-N01 – No results for invalid search

**Preconditions:**
- Orders exist

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| search | INVALID-99999 | No match |

**Steps:**
1. Search for non-existent order

**Expected Result:**
- No results message

**Traceability:**
- AC: AC-ADMIN-006
- BR: [AMB-OP-044]

---

### TC-AC-ADMIN-007-N01 – Empty export for no results

**Preconditions:**
- Filter returns 0 orders

**Steps:**
1. Attempt export

**Expected Result:**
- Empty CSV or appropriate message

**Traceability:**
- AC: AC-ADMIN-007
- BR: [AMB-OP-046]

---

### TC-AC-ADMIN-008-N01 – Invalid status transition rejected

**Preconditions:**
- Order status 'pending'

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| status | shipped | Skip confirmed |

**Steps:**
1. Attempt pending→shipped transition

**Expected Result:**
- INVALID_STATUS_TRANSITION error

**Traceability:**
- AC: AC-ADMIN-008
- BR: [AMB-OP-049]

---

### TC-AC-ADMIN-008-N02 – Shipped without tracking number

**Preconditions:**
- Order confirmed

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| tracking |  | Empty |

**Steps:**
1. Mark shipped without tracking

**Expected Result:**
- Warning or requirement for tracking

**Traceability:**
- AC: AC-ADMIN-008
- BR: [AMB-OP-048]

---

### TC-AC-ADMIN-009-N01 – Reject negative stock quantity

**Preconditions:**
- Admin setting stock

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| newStock | -10 | Negative |

**Steps:**
1. Attempt to set negative stock

**Expected Result:**
- INVALID_QUANTITY error

**Traceability:**
- AC: AC-ADMIN-009
- BR: [AMB-OP-053]

---

### TC-AC-ADMIN-010-N01 – Reject adjustment resulting in negative stock

**Preconditions:**
- Stock is 10

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| delta | -15 | Would be -5 |

**Steps:**
1. Attempt to adjust by -15

**Expected Result:**
- INSUFFICIENT_STOCK error

**Traceability:**
- AC: AC-ADMIN-010
- BR: [AMB-OP-053]

---

### TC-AC-ADMIN-011-N01 – Invalid SKU rows skipped with error

**Preconditions:**
- CSV has non-existent SKU

**Steps:**
1. Import CSV

**Expected Result:**
- Row skipped
- SKU not found error reported

**Traceability:**
- AC: AC-ADMIN-011
- BR: [AMB-OP-054]

---

### TC-AC-ADMIN-011-N02 – Invalid quantity rows skipped

**Preconditions:**
- CSV has negative quantity

**Steps:**
1. Import CSV

**Expected Result:**
- Row skipped
- Invalid quantity error reported

**Traceability:**
- AC: AC-ADMIN-011
- BR: [AMB-OP-054]

---

### TC-AC-ADMIN-012-N01 – No alert when stock above threshold

**Preconditions:**
- Threshold is 10
- Stock is 15

**Steps:**
1. Stock drops to 11

**Expected Result:**
- No alert triggered

**Traceability:**
- AC: AC-ADMIN-012
- BR: [AMB-ENT-041]

---

### TC-AC-ADMIN-013-N01 – Reject more than 3 levels of nesting

**Preconditions:**
- 3 levels already exist

**Steps:**
1. Attempt to create 4th level

**Expected Result:**
- MAX_NESTING_EXCEEDED error

**Traceability:**
- AC: AC-ADMIN-013
- BR: [AMB-ENT-037]

---

### TC-AC-PAY-001-N01 – Declined card rejected

**Preconditions:**
- Card will be declined

**Steps:**
1. Submit payment

**Expected Result:**
- PAYMENT_DECLINED error
- Reason shown

**Traceability:**
- AC: AC-PAY-001
- BR: [AMB-ENT-035]

---

### TC-AC-PAY-001-N02 – Invalid card number rejected

**Preconditions:**
- Invalid card number format

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| cardNumber | 1234 | Invalid |

**Steps:**
1. Submit payment

**Expected Result:**
- INVALID_CARD error

**Traceability:**
- AC: AC-PAY-001
- BR: [AMB-ENT-035]

---

### TC-AC-PAY-001-N03 – Unsupported card type rejected

**Preconditions:**
- Discover card (not supported)

**Steps:**
1. Attempt payment with Discover

**Expected Result:**
- UNSUPPORTED_CARD_TYPE error

**Traceability:**
- AC: AC-PAY-001
- BR: [AMB-ENT-034]

---

### TC-AC-PAY-002-N01 – Cancelled PayPal authorization

**Preconditions:**
- Customer cancels PayPal

**Steps:**
1. Cancel PayPal authorization

**Expected Result:**
- PAYMENT_CANCELLED
- Return to checkout

**Traceability:**
- AC: AC-PAY-002

---

### TC-AC-PAY-002-N02 – PayPal timeout handled

**Preconditions:**
- PayPal session expires

**Steps:**
1. Wait for timeout

**Expected Result:**
- Timeout error
- Return to checkout

**Traceability:**
- AC: AC-PAY-002

---

### TC-AC-PAY-003-N01 – Opt-out does not save card

**Preconditions:**
- Customer completing order

**Steps:**
1. Do NOT select save option
2. Complete order

**Expected Result:**
- Card not saved
- Not available at future checkout

**Traceability:**
- AC: AC-PAY-003
- BR: [AMB-ENT-028]

---

## Boundary Tests

### TC-AC-CUST-001-B01 – Password at minimum length (8 chars)

**Preconditions:**
- Email not registered

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| password | Secure1! | Exactly 8 characters |

**Steps:**
1. Submit registration with 8-character password

**Expected Result:**
- Registration succeeds

**Traceability:**
- AC: AC-CUST-001
- BR: [AMB-ENT-025]

---

### TC-AC-CUST-001-B02 – Password at 7 chars rejected

**Preconditions:**
- Email not registered

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| password | Secur1! | Only 7 characters |

**Steps:**
1. Submit registration with 7-character password

**Expected Result:**
- INVALID_PASSWORD error

**Traceability:**
- AC: AC-CUST-001
- BR: [AMB-ENT-025]

---

### TC-AC-CUST-003-B01 – 5-digit ZIP code accepted (US minimum)

**Preconditions:**
- Customer on address form

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| postalCode | 02101 | Exactly 5 digits |

**Steps:**
1. Submit address with 5-digit ZIP

**Expected Result:**
- Address saved successfully

**Traceability:**
- AC: AC-CUST-003
- BR: [AMB-ENT-032]

---

### TC-AC-PROD-001-B01 – Products at exact price boundary included

**Preconditions:**
- Product exists at exactly $50.00

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| priceMin | 50 | Exact boundary |

**Steps:**
1. Filter with min price $50

**Expected Result:**
- $50 product is included

**Traceability:**
- AC: AC-PROD-001
- BR: [AMB-OP-002]

---

### TC-AC-PROD-003-B01 – Last variant in stock selectable

**Preconditions:**
- Only one variant has stock (1 unit)

**Steps:**
1. Select the only available variant

**Expected Result:**
- Variant selectable
- Shows stock count 1

**Traceability:**
- AC: AC-PROD-003
- BR: [AMB-ENT-011]

---

### TC-AC-CART-001-B01 – Add exactly available stock quantity

**Preconditions:**
- Product has exactly 10 units

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| quantity | 10 | Exact stock |

**Steps:**
1. Add 10 units to cart

**Expected Result:**
- All 10 added successfully

**Traceability:**
- AC: AC-CART-001
- BR: [AMB-OP-006]

---

### TC-AC-CART-002-B01 – Adding to reach exactly max stock

**Preconditions:**
- Cart has Widget qty 8
- Stock is 10

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| additionalQty | 2 | Reaches exactly 10 |

**Steps:**
1. Add 2 more units

**Expected Result:**
- Quantity becomes 10
- Success

**Traceability:**
- AC: AC-CART-002
- BR: [AMB-OP-006]

---

### TC-AC-CART-004-B01 – Set quantity to exactly max stock

**Preconditions:**
- Stock is exactly 10

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| newQuantity | 10 | Exact stock |

**Steps:**
1. Set quantity to 10

**Expected Result:**
- Quantity accepted
- No error

**Traceability:**
- AC: AC-CART-004
- BR: [AMB-OP-013]

---

### TC-AC-CART-005-B01 – Undo at exactly 5 seconds

**Preconditions:**
- Item just removed

**Steps:**
1. Click Undo at 5 second mark

**Expected Result:**
- Undo succeeds at boundary

**Traceability:**
- AC: AC-CART-005
- BR: [AMB-OP-014]

---

### TC-AC-CART-009-B01 – Cart at exactly 30 days expires

**Preconditions:**
- Cart inactive for exactly 30 days

**Steps:**
1. Check cart status

**Expected Result:**
- Cart expired at 30-day mark

**Traceability:**
- AC: AC-CART-009
- BR: [AMB-ENT-012]

---

### TC-AC-ORDER-001-B01 – Free shipping at exactly $75 subtotal

**Preconditions:**
- Cart subtotal exactly $75.00

**Steps:**
1. Place order

**Expected Result:**
- Shipping cost $0.00

**Traceability:**
- AC: AC-ORDER-001
- BR: [AMB-OP-020]

---

### TC-AC-ORDER-002-B01 – Shipping charged at $74.99 subtotal

**Preconditions:**
- Cart subtotal exactly $74.99

**Steps:**
1. View order summary

**Expected Result:**
- Shipping $5.99 (not free)

**Traceability:**
- AC: AC-ORDER-002
- BR: [AMB-OP-020]

---

### TC-AC-ORDER-003-B01 – Year rollover resets sequence

**Preconditions:**
- Last 2024 order was ORD-2024-99999
- Now 2025

**Steps:**
1. Place first order in 2025

**Expected Result:**
- Order number ORD-2025-00001

**Traceability:**
- AC: AC-ORDER-003
- BR: [AMB-ENT-017]

---

### TC-AC-ORDER-005-B01 – Tax on minimum purchase

**Preconditions:**
- Subtotal $0.01
- CA 7.25%

**Steps:**
1. Calculate tax

**Expected Result:**
- Tax calculated even on small amount

**Traceability:**
- AC: AC-ORDER-005
- BR: [AMB-OP-022]

---

### TC-AC-ADMIN-001-B01 – Name at exactly 200 characters accepted

**Preconditions:**
- Admin creating product

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| name | [200 chars] | Max length |

**Steps:**
1. Submit with 200-char name

**Expected Result:**
- Product created successfully

**Traceability:**
- AC: AC-ADMIN-001
- BR: [AMB-OP-034]

---

### TC-AC-ADMIN-002-B01 – Upload exactly 10 images

**Preconditions:**
- Product has 0 images

**Steps:**
1. Upload 10 images

**Expected Result:**
- All 10 images accepted

**Traceability:**
- AC: AC-ADMIN-002
- BR: [AMB-OP-036]

---

### TC-AC-ADMIN-009-B01 – Set stock to zero

**Preconditions:**
- Product has stock

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| newStock | 0 | Zero |

**Steps:**
1. Set stock to 0

**Expected Result:**
- Stock set to 0
- Product shows out of stock

**Traceability:**
- AC: AC-ADMIN-009
- BR: [AMB-OP-051]

---

### TC-AC-ADMIN-010-B01 – Adjust to exactly zero

**Preconditions:**
- Stock is 10

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| delta | -10 | Exactly to 0 |

**Steps:**
1. Adjust by -10

**Expected Result:**
- Stock becomes 0
- Product out of stock

**Traceability:**
- AC: AC-ADMIN-010
- BR: [AMB-OP-053]

---

### TC-AC-ADMIN-011-B01 – Single row CSV imports successfully

**Preconditions:**
- CSV has 1 valid row

**Steps:**
1. Import single-row CSV

**Expected Result:**
- Row processed successfully

**Traceability:**
- AC: AC-ADMIN-011
- BR: [AMB-OP-054]

---

### TC-AC-ADMIN-012-B01 – Alert at exactly threshold

**Preconditions:**
- Threshold is 10
- Stock is 11

**Steps:**
1. Stock drops to exactly 10

**Expected Result:**
- Alert triggered at threshold

**Traceability:**
- AC: AC-ADMIN-012
- BR: [AMB-ENT-041]

---

### TC-AC-ADMIN-013-B01 – Create exactly 3 levels deep

**Preconditions:**
- 2 levels exist

**Steps:**
1. Create 3rd level category

**Expected Result:**
- 3rd level created successfully

**Traceability:**
- AC: AC-ADMIN-013
- BR: [AMB-ENT-037]

---

### TC-AC-PAY-001-B01 – Card expiring this month accepted

**Preconditions:**
- Card expires end of current month

**Steps:**
1. Submit payment

**Expected Result:**
- Payment accepted (still valid)

**Traceability:**
- AC: AC-PAY-001

---

## Hallucination Prevention Tests

### TC-AC-CUST-001-H01 – Phone should NOT be required

**⚠️ Should NOT:** Require phone number (not specified in AC)

**Preconditions:**
- Email not registered

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| phone |  | Empty - not in AC |

**Steps:**
1. Submit registration without phone number

**Expected Result:**
- Registration succeeds

**Traceability:**
- AC: AC-CUST-001

---

### TC-AC-CUST-001-H02 – Account should NOT be auto-verified

**⚠️ Should NOT:** Set account to 'verified' without email verification

**Preconditions:**
- Valid registration data

**Steps:**
1. Complete registration
2. Check account status

**Expected Result:**
- Account status is 'unverified'

**Traceability:**
- AC: AC-CUST-001
- BR: [AMB-ENT-030]

---

### TC-AC-CUST-002-H01 – Unverified order should NOT proceed to payment

**⚠️ Should NOT:** Allow payment processing for unverified customers

**Preconditions:**
- Customer email is unverified

**Steps:**
1. Attempt checkout flow

**Expected Result:**
- Blocked before payment

**Traceability:**
- AC: AC-CUST-002

---

### TC-AC-CUST-003-H01 – Billing address should NOT be required

**⚠️ Should NOT:** Require billing address when adding shipping address

**Preconditions:**
- Customer adding shipping address

**Steps:**
1. Add only shipping address

**Expected Result:**
- Address saved without billing info

**Traceability:**
- AC: AC-CUST-003

---

### TC-AC-CUST-004-H01 – Deleted data should NOT be recoverable

**⚠️ Should NOT:** Allow recovery of erased personal data

**Preconditions:**
- Customer completed data erasure

**Steps:**
1. Attempt to access deleted customer data

**Expected Result:**
- Data not found

**Traceability:**
- AC: AC-CUST-004
- BR: [AMB-ENT-029]

---

### TC-AC-PROD-001-H01 – Deleted products should NOT appear

**⚠️ Should NOT:** Show soft-deleted products in browse results

**Preconditions:**
- Product was soft-deleted

**Steps:**
1. Browse all products

**Expected Result:**
- Deleted product not visible

**Traceability:**
- AC: AC-PROD-001

---

### TC-AC-PROD-002-H01 – Out of stock product should NOT be hidden

**⚠️ Should NOT:** Hide out-of-stock products from listings

**Preconditions:**
- Product has zero stock

**Steps:**
1. Browse product listings

**Expected Result:**
- Product is visible with indicator

**Traceability:**
- AC: AC-PROD-002
- BR: [AMB-OP-004]

---

### TC-AC-PROD-003-H01 – Base product price should NOT be used for variants

**⚠️ Should NOT:** Use base product price instead of variant price

**Preconditions:**
- Variant has different price than base

**Steps:**
1. Select variant
2. Add to cart

**Expected Result:**
- Variant price used in cart

**Traceability:**
- AC: AC-PROD-003
- BR: [AMB-ENT-010]

---

### TC-AC-CART-001-H01 – Cart total should NOT include shipping yet

**⚠️ Should NOT:** Include shipping cost in cart subtotal

**Preconditions:**
- Product added to cart

**Steps:**
1. View cart totals

**Expected Result:**
- Only subtotal shown
- No shipping in cart view

**Traceability:**
- AC: AC-CART-001

---

### TC-AC-CART-002-H01 – Should NOT create separate line items

**⚠️ Should NOT:** Create multiple line items for same product

**Preconditions:**
- Same product added twice

**Steps:**
1. Add same product multiple times

**Expected Result:**
- Single line item exists

**Traceability:**
- AC: AC-CART-002
- BR: [AMB-OP-005]

---

### TC-AC-CART-003-H01 – Guest cart should NOT be lost on login

**⚠️ Should NOT:** Discard guest cart items on login

**Preconditions:**
- Guest has items in cart

**Steps:**
1. Log in to account

**Expected Result:**
- All guest items present in merged cart

**Traceability:**
- AC: AC-CART-003
- BR: [AMB-ENT-013]

---

### TC-AC-CART-004-H01 – Negative quantity should NOT be accepted

**⚠️ Should NOT:** Accept negative quantity values

**Preconditions:**
- Cart has item

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| newQuantity | -1 | Negative |

**Steps:**
1. Attempt negative quantity

**Expected Result:**
- Rejected or treated as 0

**Traceability:**
- AC: AC-CART-004

---

### TC-AC-CART-005-H01 – Removed item should NOT reappear

**⚠️ Should NOT:** Show removed items after undo expires

**Preconditions:**
- Item removed
- Undo expired

**Steps:**
1. Refresh cart page

**Expected Result:**
- Removed item stays removed

**Traceability:**
- AC: AC-CART-005

---

### TC-AC-CART-006-H01 – Empty cart should NOT show checkout button

**⚠️ Should NOT:** Show enabled checkout button for empty cart

**Preconditions:**
- Cart is empty

**Steps:**
1. View empty cart

**Expected Result:**
- Checkout button hidden or disabled

**Traceability:**
- AC: AC-CART-006

---

### TC-AC-CART-007-H01 – Out of stock item should NOT be auto-removed

**⚠️ Should NOT:** Automatically remove items when stock depletes

**Preconditions:**
- Cart item goes out of stock

**Steps:**
1. View cart

**Expected Result:**
- Item still in cart with warning

**Traceability:**
- AC: AC-CART-007
- BR: [AMB-ENT-014]

---

### TC-AC-CART-008-H01 – Old price should NOT be charged

**⚠️ Should NOT:** Use original cart price at checkout

**Preconditions:**
- Price changed since added to cart

**Steps:**
1. Proceed to checkout

**Expected Result:**
- Current price used for checkout

**Traceability:**
- AC: AC-CART-008
- BR: [AMB-ENT-016]

---

### TC-AC-CART-009-H01 – Cart at 29 days should NOT expire

**⚠️ Should NOT:** Expire cart before 30 days

**Preconditions:**
- Cart inactive for 29 days

**Steps:**
1. Log in

**Expected Result:**
- Cart items still present

**Traceability:**
- AC: AC-CART-009
- BR: [AMB-ENT-012]

---

### TC-AC-ORDER-001-H01 – Inventory should NOT decrement on failure

**⚠️ Should NOT:** Decrement inventory when order fails

**Preconditions:**
- Payment fails

**Steps:**
1. Check inventory after failed order

**Expected Result:**
- Inventory unchanged

**Traceability:**
- AC: AC-ORDER-001
- BR: [AMB-OP-017]

---

### TC-AC-ORDER-002-H01 – Shipping should NOT be included in tax calc

**⚠️ Should NOT:** Calculate tax on shipping cost

**Preconditions:**
- Order with shipping

**Steps:**
1. Check tax calculation

**Expected Result:**
- Tax on subtotal only

**Traceability:**
- AC: AC-ORDER-002
- BR: [AMB-OP-022]

---

### TC-AC-ORDER-003-H01 – Order numbers should NOT have gaps

**⚠️ Should NOT:** Skip order numbers for failed attempts

**Preconditions:**
- Failed order attempt

**Steps:**
1. Place successful order after failure

**Expected Result:**
- Next sequential number used

**Traceability:**
- AC: AC-ORDER-003
- BR: [AMB-ENT-017]

---

### TC-AC-ORDER-004-H01 – Order price should NOT update retroactively

**⚠️ Should NOT:** Update historical order prices when catalog changes

**Preconditions:**
- Order placed
- Price later changes

**Steps:**
1. View order details

**Expected Result:**
- Original price shown

**Traceability:**
- AC: AC-ORDER-004
- BR: [AMB-ENT-019]

---

### TC-AC-ORDER-005-H01 – Tax should NOT be calculated on shipping

**⚠️ Should NOT:** Include shipping in tax calculation base

**Preconditions:**
- Order with shipping cost

**Steps:**
1. Verify tax base

**Expected Result:**
- Tax excludes shipping amount

**Traceability:**
- AC: AC-ORDER-005
- BR: [AMB-OP-022]

---

### TC-AC-ORDER-006-H01 – Other customers' orders should NOT appear

**⚠️ Should NOT:** Show orders belonging to other customers

**Preconditions:**
- Customer viewing order history

**Steps:**
1. View order history

**Expected Result:**
- Only own orders shown

**Traceability:**
- AC: AC-ORDER-006

---

### TC-AC-ORDER-007-H01 – Tracking should NOT appear before shipment

**⚠️ Should NOT:** Display placeholder or fake tracking before shipment

**Preconditions:**
- Order not yet shipped

**Steps:**
1. View order details

**Expected Result:**
- No tracking information shown

**Traceability:**
- AC: AC-ORDER-007

---

### TC-AC-ORDER-008-H01 – Cancelled order should NOT be shipped

**⚠️ Should NOT:** Ship cancelled orders

**Preconditions:**
- Order cancelled

**Steps:**
1. Check fulfillment status

**Expected Result:**
- Order not in fulfillment queue

**Traceability:**
- AC: AC-ORDER-008
- BR: [AMB-OP-028]

---

### TC-AC-ORDER-009-H01 – Confirmed order refund should NOT be partial

**⚠️ Should NOT:** Issue partial refund for non-shipped cancellation

**Preconditions:**
- Confirmed order cancelled

**Steps:**
1. Check refund amount

**Expected Result:**
- Full refund issued

**Traceability:**
- AC: AC-ORDER-009
- BR: [AMB-OP-030]

---

### TC-AC-ORDER-010-H01 – Shipped order status should NOT revert

**⚠️ Should NOT:** Allow status regression from shipped

**Preconditions:**
- Order shipped

**Steps:**
1. Check order status transitions

**Expected Result:**
- Cannot go back to pending/confirmed

**Traceability:**
- AC: AC-ORDER-010
- BR: [AMB-ENT-018]

---

### TC-AC-ORDER-011-H01 – Email should NOT contain payment details

**⚠️ Should NOT:** Include full credit card number in email

**Preconditions:**
- Order confirmation email

**Steps:**
1. Check email content

**Expected Result:**
- No full card number
- Only last 4 digits if any

**Traceability:**
- AC: AC-ORDER-011

---

### TC-AC-ADMIN-001-H01 – New product should NOT be visible to customers

**⚠️ Should NOT:** Show draft products in customer catalog

**Preconditions:**
- Product just created

**Steps:**
1. Browse products as customer

**Expected Result:**
- Draft product not visible

**Traceability:**
- AC: AC-ADMIN-001
- BR: [AMB-ENT-003]

---

### TC-AC-ADMIN-002-H01 – Product should NOT require images

**⚠️ Should NOT:** Block product creation without images

**Preconditions:**
- Draft product

**Steps:**
1. Publish product without images

**Expected Result:**
- Should not require images

**Traceability:**
- AC: AC-ADMIN-002

---

### TC-AC-ADMIN-003-H01 – Existing orders should NOT reflect new price

**⚠️ Should NOT:** Update prices in existing orders

**Preconditions:**
- Order placed at $25
- Price changed to $30

**Steps:**
1. View existing order

**Expected Result:**
- Order still shows $25

**Traceability:**
- AC: AC-ADMIN-003
- BR: [AMB-ENT-019]

---

### TC-AC-ADMIN-004-H01 – Concurrent edit should NOT cause data corruption

**⚠️ Should NOT:** Corrupt product data on concurrent saves

**Preconditions:**
- Concurrent edits

**Steps:**
1. Both admins save simultaneously

**Expected Result:**
- Product data remains valid

**Traceability:**
- AC: AC-ADMIN-004

---

### TC-AC-ADMIN-005-H01 – Deleted product should NOT appear in search

**⚠️ Should NOT:** Show soft-deleted products in search

**Preconditions:**
- Product soft-deleted

**Steps:**
1. Customer searches for product

**Expected Result:**
- Product not in search results

**Traceability:**
- AC: AC-ADMIN-005
- BR: [AMB-OP-041]

---

### TC-AC-ADMIN-006-H01 – Admin should NOT see customer payment details

**⚠️ Should NOT:** Display full payment credentials in order list

**Preconditions:**
- Admin viewing orders

**Steps:**
1. View order list

**Expected Result:**
- No full card numbers shown

**Traceability:**
- AC: AC-ADMIN-006

---

### TC-AC-ADMIN-007-H01 – CSV should NOT contain sensitive data

**⚠️ Should NOT:** Export sensitive payment or password data

**Preconditions:**
- Export generated

**Steps:**
1. Review CSV contents

**Expected Result:**
- No full card numbers or passwords

**Traceability:**
- AC: AC-ADMIN-007

---

### TC-AC-ADMIN-008-H01 – Cancelled order should NOT become shipped

**⚠️ Should NOT:** Allow cancelled orders to be shipped

**Preconditions:**
- Order is cancelled

**Steps:**
1. Attempt to mark shipped

**Expected Result:**
- Transition blocked

**Traceability:**
- AC: AC-ADMIN-008
- BR: [AMB-OP-049]

---

### TC-AC-ADMIN-009-H01 – Stock change should NOT affect in-flight orders

**⚠️ Should NOT:** Cancel or modify pending orders on stock reduction

**Preconditions:**
- Pending order exists for product

**Steps:**
1. Reduce stock below order quantity

**Expected Result:**
- Pending order unaffected

**Traceability:**
- AC: AC-ADMIN-009

---

### TC-AC-ADMIN-010-H01 – Delta adjustment should NOT skip audit

**⚠️ Should NOT:** Skip audit logging for delta adjustments

**Preconditions:**
- Stock adjusted

**Steps:**
1. Check audit log

**Expected Result:**
- Adjustment logged

**Traceability:**
- AC: AC-ADMIN-010
- BR: [AMB-OP-052]

---

### TC-AC-ADMIN-011-H01 – Partial failures should NOT rollback successes

**⚠️ Should NOT:** Rollback successful updates when some rows fail

**Preconditions:**
- CSV has mix of valid and invalid rows

**Steps:**
1. Import CSV
2. Check updates

**Expected Result:**
- Valid rows are updated

**Traceability:**
- AC: AC-ADMIN-011

---

### TC-AC-ADMIN-012-H01 – Alert should NOT trigger repeatedly

**⚠️ Should NOT:** Send multiple alerts for same low-stock condition

**Preconditions:**
- Already below threshold

**Steps:**
1. Stock drops further

**Expected Result:**
- No duplicate alert

**Traceability:**
- AC: AC-ADMIN-012

---

### TC-AC-ADMIN-013-H01 – Circular parent reference should NOT be allowed

**⚠️ Should NOT:** Allow circular parent-child relationships

**Preconditions:**
- Category A exists under Category B

**Steps:**
1. Try to set B as child of A

**Expected Result:**
- Circular reference blocked

**Traceability:**
- AC: AC-ADMIN-013

---

### TC-AC-PAY-001-H01 – Full card number should NOT be stored

**⚠️ Should NOT:** Store full card number in database

**Preconditions:**
- Payment completed

**Steps:**
1. Check stored payment data

**Expected Result:**
- Only token and last 4 digits stored

**Traceability:**
- AC: AC-PAY-001
- BR: [AMB-ENT-036]

---

### TC-AC-PAY-002-H01 – PayPal credentials should NOT be stored locally

**⚠️ Should NOT:** Store PayPal email or password

**Preconditions:**
- PayPal payment completed

**Steps:**
1. Check stored data

**Expected Result:**
- Only transaction reference stored

**Traceability:**
- AC: AC-PAY-002

---

### TC-AC-PAY-003-H01 – Saved card should NOT store full number

**⚠️ Should NOT:** Display or store full card number

**Preconditions:**
- Card saved for future

**Steps:**
1. View saved payment methods

**Expected Result:**
- Only last 4 digits shown

**Traceability:**
- AC: AC-PAY-003
- BR: [AMB-ENT-036]

---

