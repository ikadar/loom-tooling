# Test Cases

Generated: 2025-12-29T15:49:16+01:00

---

## TC-CUST-001-01 – Register new customer successfully

**Type:** happy_path

**Preconditions:**
- Email not registered

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| email | john@example.com | Valid |
| password | SecurePass1 | Valid |
| first_name | John |  |
| last_name | Doe |  |

**Steps:**
1. Submit registration form

**Expected Result:**
- Account created with unverified status
- Verification email sent

**Traceability:**
- AC: AC-CUST-001

---

## TC-CUST-001-02 – Register with duplicate email fails

**Type:** error_case

**Preconditions:**
- Email already registered

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| email | existing@example.com | Existing |

**Steps:**
1. Submit registration with existing email

**Expected Result:**
- DUPLICATE_EMAIL error
- Login suggestion shown

**Traceability:**
- AC: AC-CUST-001

---

## TC-CUST-002-01 – Unverified user cannot place order

**Type:** happy_path

**Preconditions:**
- Customer registered
- Email unverified
- Cart has items

**Steps:**
1. Attempt to place order

**Expected Result:**
- EMAIL_NOT_VERIFIED error
- Verify email prompt shown

**Traceability:**
- AC: AC-CUST-002

---

## TC-CUST-003-01 – Add new shipping address

**Type:** happy_path

**Preconditions:**
- Customer logged in
- One address exists

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| street | 456 Oak Ave |  |
| city | Boston |  |
| state | MA |  |
| postal_code | 02101 |  |
| country | US |  |

**Steps:**
1. Add new address

**Expected Result:**
- Address saved
- Available for selection

**Traceability:**
- AC: AC-CUST-003

---

## TC-CUST-003-02 – Add address with invalid postal code

**Type:** error_case

**Preconditions:**
- Customer logged in

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| postal_code | INVALID | Invalid format |

**Steps:**
1. Submit address with bad postal code

**Expected Result:**
- INVALID_POSTAL_CODE error

**Traceability:**
- AC: AC-CUST-003

---

## TC-CUST-004-01 – GDPR data erasure request

**Type:** happy_path

**Preconditions:**
- Customer registered
- No pending orders

**Steps:**
1. Request account deletion
2. Confirm erasure

**Expected Result:**
- Personal data deleted
- Orders anonymized
- Confirmation sent

**Traceability:**
- AC: AC-CUST-004

---

## TC-CUST-004-02 – GDPR erasure blocked by pending order

**Type:** error_case

**Preconditions:**
- Customer has pending order

**Steps:**
1. Request account deletion

**Expected Result:**
- PENDING_ORDERS_EXIST error

**Traceability:**
- AC: AC-CUST-004

---

## TC-PROD-001-01 – Filter products by category and price

**Type:** happy_path

**Preconditions:**
- Products exist in multiple categories

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| category | Electronics |  |
| price_min | 50 |  |
| price_max | 200 |  |
| availability | in stock |  |

**Steps:**
1. Apply filters

**Expected Result:**
- Matching products shown
- Sorted by newest
- 20 per page

**Traceability:**
- AC: AC-PROD-001

---

## TC-PROD-001-02 – No products match filters

**Type:** error_case

**Preconditions:**
- No products match criteria

**Steps:**
1. Apply restrictive filters

**Expected Result:**
- No products found message
- Clear filters option

**Traceability:**
- AC: AC-PROD-001

---

## TC-PROD-002-01 – View out-of-stock product

**Type:** happy_path

**Preconditions:**
- Product has zero stock

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| product | Widget Pro | Out of stock |

**Steps:**
1. View product listing

**Expected Result:**
- Out of Stock indicator shown
- Add to Cart disabled

**Traceability:**
- AC: AC-PROD-002

---

## TC-PROD-003-01 – View product with variants

**Type:** happy_path

**Preconditions:**
- Product has size and color variants

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| product | T-Shirt | Has variants |

**Steps:**
1. View product detail
2. Select variant

**Expected Result:**
- Variant options displayed
- Price/stock updates

**Traceability:**
- AC: AC-PROD-003

---

## TC-PROD-003-02 – Selected variant out of stock

**Type:** error_case

**Preconditions:**
- Specific variant has no stock

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| variant | Size L, Red | Out of stock |

**Steps:**
1. Select out-of-stock variant

**Expected Result:**
- Out of Stock shown for variant

**Traceability:**
- AC: AC-PROD-003

---

## TC-CART-001-01 – Add product to cart

**Type:** happy_path

**Preconditions:**
- Product in stock

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| product | Widget | 10 in stock |
| quantity | 2 |  |
| price | 25.00 |  |

**Steps:**
1. Add 2 units to cart

**Expected Result:**
- Cart shows qty 2
- Subtotal $50.00
- Toast notification

**Traceability:**
- AC: AC-CART-001

---

## TC-CART-001-02 – Add out-of-stock product fails

**Type:** error_case

**Preconditions:**
- Product out of stock

**Steps:**
1. Attempt to add to cart

**Expected Result:**
- OUT_OF_STOCK error

**Traceability:**
- AC: AC-CART-001

---

## TC-CART-002-01 – Add existing item increases quantity

**Type:** happy_path

**Preconditions:**
- Cart has Widget qty 2

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| add_quantity | 3 |  |

**Steps:**
1. Add 3 more Widget

**Expected Result:**
- Cart shows Widget qty 5
- Single line item

**Traceability:**
- AC: AC-CART-002

---

## TC-CART-003-01 – Guest cart merges on login

**Type:** happy_path

**Preconditions:**
- Guest has cart items
- Account has cart items

**Steps:**
1. Add items as guest
2. Log in

**Expected Result:**
- Carts merged

**Traceability:**
- AC: AC-CART-003

---

## TC-CART-004-01 – Update cart item quantity

**Type:** happy_path

**Preconditions:**
- Cart has Widget qty 3
- 10 in stock

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| new_quantity | 5 |  |

**Steps:**
1. Change quantity to 5

**Expected Result:**
- Quantity updates to 5
- Totals recalculate

**Traceability:**
- AC: AC-CART-004

---

## TC-CART-004-02 – Quantity exceeds stock limited

**Type:** error_case

**Preconditions:**
- 10 units in stock

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| new_quantity | 15 | Exceeds stock |

**Steps:**
1. Set quantity to 15

**Expected Result:**
- Limited to 10
- Only 10 available notification

**Traceability:**
- AC: AC-CART-004

---

## TC-CART-005-01 – Remove item from cart

**Type:** happy_path

**Preconditions:**
- Cart has Widget and Gadget

**Steps:**
1. Remove Widget

**Expected Result:**
- Widget removed
- Undo option for 5 seconds

**Traceability:**
- AC: AC-CART-005

---

## TC-CART-006-01 – Remove last item shows empty cart

**Type:** happy_path

**Preconditions:**
- Cart has only Widget

**Steps:**
1. Remove Widget

**Expected Result:**
- Empty cart displayed
- Continue Shopping link

**Traceability:**
- AC: AC-CART-006

---

## TC-CART-007-01 – Cart item goes out of stock

**Type:** happy_path

**Preconditions:**
- Cart has Widget qty 2

**Steps:**
1. Inventory drops to 0
2. View cart

**Expected Result:**
- Out of Stock warning
- Checkout blocked

**Traceability:**
- AC: AC-CART-007

---

## TC-CART-008-01 – Cart price update notification

**Type:** happy_path

**Preconditions:**
- Cart has Widget at $25.00

**Steps:**
1. Admin changes price to $30.00
2. View cart

**Expected Result:**
- Price shows $30.00
- Price change notification

**Traceability:**
- AC: AC-CART-008

---

## TC-CART-009-01 – Cart expires after 30 days

**Type:** happy_path

**Preconditions:**
- Cart inactive 30 days

**Steps:**
1. Expiration passes
2. Customer logs in

**Expected Result:**
- Cart cleared
- Customer notified

**Traceability:**
- AC: AC-CART-009

---

## TC-ORDER-001-01 – Place order with free shipping

**Type:** happy_path

**Preconditions:**
- Verified customer
- Cart $75+
- Valid payment

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| subtotal | 75.00 | Free shipping |

**Steps:**
1. Place order

**Expected Result:**
- Order pending
- Shipping $0
- Cart cleared

**Traceability:**
- AC: AC-ORDER-001

---

## TC-ORDER-001-02 – Order fails with payment declined

**Type:** error_case

**Preconditions:**
- Payment will fail

**Steps:**
1. Place order

**Expected Result:**
- PAYMENT_FAILED error
- Cart preserved

**Traceability:**
- AC: AC-ORDER-001

---

## TC-ORDER-002-01 – Place order with paid shipping

**Type:** happy_path

**Preconditions:**
- Verified customer
- Cart $35

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| subtotal | 35.00 | Below free shipping |

**Steps:**
1. Place order

**Expected Result:**
- Shipping $5.99
- Total $40.99 + tax

**Traceability:**
- AC: AC-ORDER-002

---

## TC-ORDER-003-01 – Order number generated correctly

**Type:** happy_path

**Preconditions:**
- Last order ORD-2024-00042

**Steps:**
1. Place new order

**Expected Result:**
- Order number ORD-2024-00043

**Traceability:**
- AC: AC-ORDER-003

---

## TC-ORDER-004-01 – Order stores price snapshot

**Type:** happy_path

**Preconditions:**
- Widget price $25.00

**Steps:**
1. Place order
2. Admin changes price

**Expected Result:**
- Order line item shows $25.00

**Traceability:**
- AC: AC-ORDER-004

---

## TC-ORDER-005-01 – Tax calculated for California

**Type:** happy_path

**Preconditions:**
- Subtotal $100
- Shipping to CA

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| tax_rate | 7.25% | CA rate |

**Steps:**
1. Calculate order total

**Expected Result:**
- Tax $7.25

**Traceability:**
- AC: AC-ORDER-005

---

## TC-ORDER-006-01 – View order history

**Type:** happy_path

**Preconditions:**
- Customer has multiple orders

**Steps:**
1. View order history

**Expected Result:**
- Orders with pagination
- Shows status and details

**Traceability:**
- AC: AC-ORDER-006

---

## TC-ORDER-007-01 – View shipped order with tracking

**Type:** happy_path

**Preconditions:**
- Order shipped with tracking

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| tracking | 1Z999AA10123456784 |  |

**Steps:**
1. View order details

**Expected Result:**
- Tracking number shown
- Carrier link provided

**Traceability:**
- AC: AC-ORDER-007

---

## TC-ORDER-008-01 – Cancel pending order

**Type:** happy_path

**Preconditions:**
- Order status pending

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| reason | Changed mind |  |

**Steps:**
1. Request cancellation

**Expected Result:**
- Status cancelled
- Inventory restored
- Refund initiated

**Traceability:**
- AC: AC-ORDER-008

---

## TC-ORDER-009-01 – Cancel confirmed order

**Type:** happy_path

**Preconditions:**
- Order status confirmed

**Steps:**
1. Request cancellation

**Expected Result:**
- Status cancelled
- Refund initiated

**Traceability:**
- AC: AC-ORDER-009

---

## TC-ORDER-010-01 – Cannot cancel shipped order

**Type:** error_case

**Preconditions:**
- Order status shipped

**Steps:**
1. Attempt cancellation

**Expected Result:**
- ORDER_ALREADY_SHIPPED error

**Traceability:**
- AC: AC-ORDER-010

---

## TC-ORDER-011-01 – Order confirmation email sent

**Type:** happy_path

**Preconditions:**
- Order just placed

**Steps:**
1. Complete order placement

**Expected Result:**
- Email queued with order details

**Traceability:**
- AC: AC-ORDER-011

---

## TC-ADMIN-001-01 – Admin adds new product

**Type:** happy_path

**Preconditions:**
- Admin logged in

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| name | Super Widget |  |
| price | 29.99 |  |
| category | Electronics |  |

**Steps:**
1. Create product

**Expected Result:**
- Product created as draft
- UUID assigned

**Traceability:**
- AC: AC-ADMIN-001

---

## TC-ADMIN-001-02 – Add product with invalid price

**Type:** error_case

**Preconditions:**
- Admin logged in

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| price | -5.00 | Invalid |

**Steps:**
1. Submit product with negative price

**Expected Result:**
- INVALID_PRICE error

**Traceability:**
- AC: AC-ADMIN-001

---

## TC-ADMIN-002-01 – Upload product images

**Type:** happy_path

**Preconditions:**
- Product in draft status

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| images | 3 images |  |

**Steps:**
1. Upload images
2. Set primary

**Expected Result:**
- Images uploaded
- Primary marked

**Traceability:**
- AC: AC-ADMIN-002

---

## TC-ADMIN-002-02 – Upload exceeds max images

**Type:** error_case

**Preconditions:**
- Product exists

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| images | 11 images | Exceeds limit |

**Steps:**
1. Upload 11 images

**Expected Result:**
- MAX_IMAGES_EXCEEDED error

**Traceability:**
- AC: AC-ADMIN-002

---

## TC-ADMIN-003-01 – Edit product price

**Type:** happy_path

**Preconditions:**
- Widget price $25.00

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| new_price | 30.00 |  |

**Steps:**
1. Change price

**Expected Result:**
- Price updated
- Change logged

**Traceability:**
- AC: AC-ADMIN-003

---

## TC-ADMIN-004-01 – Concurrent edit shows conflict

**Type:** happy_path

**Preconditions:**
- Two admins editing same product

**Steps:**
1. Admin A saves
2. Admin B saves

**Expected Result:**
- Admin B overwrites
- Conflict warning shown

**Traceability:**
- AC: AC-ADMIN-004

---

## TC-ADMIN-005-01 – Soft delete product with orders

**Type:** happy_path

**Preconditions:**
- Product has order history

**Steps:**
1. Confirm removal

**Expected Result:**
- Product soft-deleted
- History preserved

**Traceability:**
- AC: AC-ADMIN-005

---

## TC-ADMIN-006-01 – Filter orders by status and date

**Type:** happy_path

**Preconditions:**
- 100 orders exist

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| status | shipped |  |
| date_range | last 7 days |  |

**Steps:**
1. Apply filters

**Expected Result:**
- Matching orders displayed

**Traceability:**
- AC: AC-ADMIN-006

---

## TC-ADMIN-007-01 – Export orders to CSV

**Type:** happy_path

**Preconditions:**
- Filtered list with 25 orders

**Steps:**
1. Click Export CSV

**Expected Result:**
- CSV file downloads

**Traceability:**
- AC: AC-ADMIN-007

---

## TC-ADMIN-008-01 – Update order status to shipped

**Type:** happy_path

**Preconditions:**
- Order status confirmed

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| tracking | 1Z999AA10123456784 |  |

**Steps:**
1. Update to shipped with tracking

**Expected Result:**
- Status shipped
- Tracking stored
- Email sent

**Traceability:**
- AC: AC-ADMIN-008

---

## TC-ADMIN-008-02 – Invalid status transition rejected

**Type:** error_case

**Preconditions:**
- Order status pending

**Steps:**
1. Try to set status to shipped

**Expected Result:**
- INVALID_STATUS_TRANSITION error

**Traceability:**
- AC: AC-ADMIN-008

---

## TC-ADMIN-009-01 – Set inventory quantity

**Type:** happy_path

**Preconditions:**
- Widget stock 50

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| new_stock | 100 |  |
| reason | Shipment received |  |

**Steps:**
1. Set stock to 100

**Expected Result:**
- Stock 100
- Audit log created

**Traceability:**
- AC: AC-ADMIN-009

---

## TC-ADMIN-010-01 – Adjust inventory delta

**Type:** happy_path

**Preconditions:**
- Widget stock 50

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| adjustment | -5 |  |
| reason | Damaged goods |  |

**Steps:**
1. Adjust by -5

**Expected Result:**
- Stock 45
- Audit log created

**Traceability:**
- AC: AC-ADMIN-010

---

## TC-ADMIN-010-02 – Adjustment causing negative stock fails

**Type:** error_case

**Preconditions:**
- Widget stock 5

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| adjustment | -10 | Would go negative |

**Steps:**
1. Adjust by -10

**Expected Result:**
- INSUFFICIENT_STOCK error

**Traceability:**
- AC: AC-ADMIN-010

---

## TC-ADMIN-011-01 – Bulk inventory update via CSV

**Type:** happy_path

**Preconditions:**
- CSV with 50 SKUs

**Steps:**
1. Import CSV

**Expected Result:**
- Valid rows updated
- Errors reported

**Traceability:**
- AC: AC-ADMIN-011

---

## TC-ADMIN-012-01 – Low stock alert triggered

**Type:** happy_path

**Preconditions:**
- Threshold 10
- Current stock 15

**Steps:**
1. Stock drops to 8

**Expected Result:**
- Low stock alert triggered

**Traceability:**
- AC: AC-ADMIN-012

---

## TC-ADMIN-013-01 – Create category hierarchy

**Type:** happy_path

**Preconditions:**
- Electronics category exists

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| name | Smartphones |  |
| parent | Electronics |  |

**Steps:**
1. Create subcategory

**Expected Result:**
- Subcategory created with parent ref

**Traceability:**
- AC: AC-ADMIN-013

---

## TC-ADMIN-013-02 – Category exceeds max nesting

**Type:** error_case

**Preconditions:**
- 3 levels already exist

**Steps:**
1. Create 4th level category

**Expected Result:**
- MAX_NESTING_EXCEEDED error

**Traceability:**
- AC: AC-ADMIN-013

---

## TC-PAY-001-01 – Pay with credit card

**Type:** happy_path

**Preconditions:**
- Customer at checkout

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| card | Visa 4242 | Valid |

**Steps:**
1. Submit payment

**Expected Result:**
- Payment authorized
- Order created

**Traceability:**
- AC: AC-PAY-001

---

## TC-PAY-001-02 – Card declined

**Type:** error_case

**Preconditions:**
- Invalid card

**Steps:**
1. Submit payment

**Expected Result:**
- PAYMENT_DECLINED error with reason

**Traceability:**
- AC: AC-PAY-001

---

## TC-PAY-002-01 – Pay with PayPal

**Type:** happy_path

**Preconditions:**
- Customer selects PayPal

**Steps:**
1. Complete PayPal authorization

**Expected Result:**
- Payment authorized
- Order created

**Traceability:**
- AC: AC-PAY-002

---

## TC-PAY-002-02 – PayPal cancelled

**Type:** error_case

**Preconditions:**
- Customer at PayPal

**Steps:**
1. Cancel PayPal authorization

**Expected Result:**
- PAYMENT_CANCELLED
- Return to checkout

**Traceability:**
- AC: AC-PAY-002

---

## TC-PAY-003-01 – Save payment method for future

**Type:** happy_path

**Preconditions:**
- Customer completing order

**Steps:**
1. Select save for future
2. Complete order

**Expected Result:**
- Token stored
- Available for future use

**Traceability:**
- AC: AC-PAY-003

---

