# Test Cases

Generated: 2025-12-29T12:39:52+01:00

---

## TC-AC-PROD-001-01 – Filter products by single category

**Type:** happy_path

**Preconditions:**
- Product listing page is accessible
- Multiple products exist in different categories

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| category | Electronics | Valid category with products |
| sort_by | name | Default sort |
| page_size | 20 | Default pagination |

**Steps:**
1. Navigate to product listing page
2. Select 'Electronics' from category filter
3. Observe the filtered results

**Expected Result:**
- Only products in Electronics category are displayed
- Products are sorted by name
- Maximum 20 products per page

**Traceability:**
- AC: AC-PROD-001

---

## TC-AC-PROD-001-02 – Filter products by price range and availability

**Type:** happy_path

**Preconditions:**
- Product listing page is accessible
- Products exist in various price ranges

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| price_min | 10 | Minimum price |
| price_max | 50 | Maximum price |
| availability | in_stock | Only available products |
| sort_by | price_asc | Price ascending |

**Steps:**
1. Navigate to product listing page
2. Set minimum price to $10
3. Set maximum price to $50
4. Select 'In Stock' availability filter
5. Select 'Price: Low to High' sort option

**Expected Result:**
- Only products priced between $10-$50 are displayed
- All displayed products are in stock
- Products are sorted by price ascending
- Maximum 20 products per page

**Traceability:**
- AC: AC-PROD-001

---

## TC-AC-PROD-001-03 – Sort products by newest

**Type:** happy_path

**Preconditions:**
- Product listing page is accessible
- Products with different creation dates exist

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| sort_by | newest | Sort by creation date descending |

**Steps:**
1. Navigate to product listing page
2. Select 'Newest' sort option

**Expected Result:**
- Products are sorted by creation date descending
- Most recently added products appear first

**Traceability:**
- AC: AC-PROD-001

---

## TC-AC-PROD-001-04 – No products match applied filters

**Type:** error_case

**Preconditions:**
- Product listing page is accessible

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| category | Electronics | Valid category |
| price_min | 10000 | Very high minimum price |
| price_max | 20000 | Very high maximum price |

**Steps:**
1. Navigate to product listing page
2. Select 'Electronics' category
3. Set price range to $10,000 - $20,000 where no products exist

**Expected Result:**
- Empty state is displayed
- 'No products found' message is shown
- Option to clear filters is available

**Traceability:**
- AC: AC-PROD-001

---

## TC-AC-PROD-001-05 – Invalid price range - minimum greater than maximum

**Type:** error_case

**Preconditions:**
- Product listing page is accessible

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| price_min | 100 | Higher than max |
| price_max | 50 | Lower than min |

**Steps:**
1. Navigate to product listing page
2. Set minimum price to $100
3. Set maximum price to $50
4. Attempt to apply filter

**Expected Result:**
- Validation error is displayed
- Filter is not applied until corrected

**Traceability:**
- AC: AC-PROD-001

---

## TC-AC-PROD-002-01 – View out-of-stock product in listing

**Type:** happy_path

**Preconditions:**
- Product exists with zero inventory

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| product_id | PROD-OOS-001 | Out of stock product |
| inventory | 0 | Zero available |

**Steps:**
1. Navigate to product listing page
2. Locate out-of-stock product in listing

**Expected Result:**
- Product is displayed in listing
- 'Out of stock' indicator is visible

**Traceability:**
- AC: AC-PROD-002

---

## TC-AC-PROD-002-02 – View out-of-stock product detail page

**Type:** happy_path

**Preconditions:**
- Product exists with zero inventory

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| product_id | PROD-OOS-001 | Out of stock product |

**Steps:**
1. Navigate to product detail page for out-of-stock product

**Expected Result:**
- Product details are displayed
- 'Out of stock' indicator is visible
- Add to cart button is disabled

**Traceability:**
- AC: AC-PROD-002

---

## TC-AC-CART-001-01 – Add single in-stock product to cart

**Type:** happy_path

**Preconditions:**
- User is viewing an in-stock product
- Product has no variants or variant is selected

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| product_id | PROD-001 | Valid in-stock product |
| quantity | 1 | Single item |
| available_stock | 10 | Sufficient stock |

**Steps:**
1. Navigate to product detail page
2. Set quantity to 1
3. Click 'Add to Cart' button

**Expected Result:**
- Item is added to cart
- Toast notification appears
- 'View Cart' option is shown in toast
- Cart item count updates to reflect new item

**Traceability:**
- AC: AC-CART-001

---

## TC-AC-CART-001-02 – Add product with maximum available quantity

**Type:** happy_path

**Preconditions:**
- User is viewing an in-stock product

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| product_id | PROD-002 | Product with limited stock |
| quantity | 5 | Equal to available stock |
| available_stock | 5 | Limited inventory |

**Steps:**
1. Navigate to product detail page
2. Set quantity to 5 (max available)
3. Click 'Add to Cart' button

**Expected Result:**
- Item is added to cart with quantity 5
- Toast notification appears
- Cart item count updates

**Traceability:**
- AC: AC-CART-001

---

## TC-AC-CART-001-03 – Add to cart disabled for out-of-stock product

**Type:** error_case

**Preconditions:**
- Product has zero inventory

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| product_id | PROD-OOS-001 | Out of stock product |
| available_stock | 0 | No inventory |

**Steps:**
1. Navigate to out-of-stock product detail page
2. Observe add to cart button state

**Expected Result:**
- Add to cart button is disabled
- 'Out of Stock' message is displayed

**Traceability:**
- AC: AC-CART-001

---

## TC-AC-CART-001-04 – Quantity exceeds available stock

**Type:** error_case

**Preconditions:**
- User is viewing an in-stock product

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| product_id | PROD-003 | Limited stock product |
| quantity | 15 | Exceeds available |
| available_stock | 10 | Available inventory |

**Steps:**
1. Navigate to product detail page
2. Attempt to set quantity to 15
3. Click 'Add to Cart' button

**Expected Result:**
- Quantity is limited to available stock (10)
- Notification is shown about stock limitation

**Traceability:**
- AC: AC-CART-001

---

## TC-AC-CART-001-05 – Add to cart disabled when variant not selected

**Type:** error_case

**Preconditions:**
- Product has multiple variants
- No variant selected

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| product_id | PROD-VAR-001 | Product with size variants |
| variant_selected | false | No variant chosen |

**Steps:**
1. Navigate to product with variants
2. Do not select any variant
3. Observe add to cart button state

**Expected Result:**
- Add to cart button is disabled
- Prompt to select variant is displayed

**Traceability:**
- AC: AC-CART-001

---

## TC-AC-CART-002-01 – Add existing product increments quantity

**Type:** happy_path

**Preconditions:**
- User has Product A with quantity 2 in cart
- Product A is in stock

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| product_id | PROD-A | Already in cart |
| cart_quantity | 2 | Current cart quantity |
| add_quantity | 1 | Additional quantity |
| available_stock | 10 | Sufficient stock |

**Steps:**
1. Navigate to Product A detail page
2. Set quantity to 1
3. Click 'Add to Cart' button

**Expected Result:**
- Cart item quantity increases to 3
- No new line item created
- Cart shows single entry for Product A with quantity 3

**Traceability:**
- AC: AC-CART-002

---

## TC-AC-CART-002-02 – Add existing product - merged total exceeds stock

**Type:** error_case

**Preconditions:**
- User has Product B with quantity 8 in cart
- Product B has 10 available stock

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| product_id | PROD-B | Already in cart |
| cart_quantity | 8 | Current cart quantity |
| add_quantity | 5 | Would exceed stock |
| available_stock | 10 | Available inventory |

**Steps:**
1. Navigate to Product B detail page
2. Set quantity to 5
3. Click 'Add to Cart' button

**Expected Result:**
- Cart quantity set to maximum available (10)
- Notification shown about stock limitation

**Traceability:**
- AC: AC-CART-002

---

## TC-AC-CART-003-01 – Guest user adds item to session cart

**Type:** happy_path

**Preconditions:**
- User is not logged in
- No existing session cart

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| product_id | PROD-001 | Valid product |
| quantity | 1 | Single item |

**Steps:**
1. As guest user, navigate to product page
2. Add product to cart

**Expected Result:**
- Session-based cart is created
- Item is added to cart
- Cart persists across page navigation

**Traceability:**
- AC: AC-CART-003

---

## TC-AC-CART-003-02 – Guest cart persists for 7 days of inactivity

**Type:** edge_case

**Preconditions:**
- Guest user has items in session cart

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| cart_items | 2 | Items in cart |
| inactivity_days | 6 | Less than expiry |

**Steps:**
1. As guest user, add items to cart
2. Simulate 6 days of inactivity
3. Return to site

**Expected Result:**
- Cart items are preserved
- Cart is accessible

**Traceability:**
- AC: AC-CART-003

---

## TC-AC-CART-004-01 – Cart merge on login - no duplicates

**Type:** happy_path

**Preconditions:**
- Guest user has items in session cart
- User account exists with different items in cart

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| guest_cart_items | [map[product:PROD-A qty:2]] | Guest cart |
| account_cart_items | [map[product:PROD-B qty:1]] | Account cart |

**Steps:**
1. As guest, add Product A (qty 2) to cart
2. Log in to account that has Product B (qty 1) in cart

**Expected Result:**
- Cart contains both Product A (qty 2) and Product B (qty 1)
- Total of 2 line items in cart

**Traceability:**
- AC: AC-CART-004

---

## TC-AC-CART-004-02 – Cart merge on login - duplicate products combined

**Type:** happy_path

**Preconditions:**
- Guest user has Product A in session cart
- User account has same Product A in cart

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| guest_cart_items | [map[product:PROD-A qty:2]] | Guest cart |
| account_cart_items | [map[product:PROD-A qty:3]] | Account cart |

**Steps:**
1. As guest, add Product A (qty 2) to cart
2. Log in to account that has Product A (qty 3) in cart

**Expected Result:**
- Cart contains Product A with combined quantity 5
- Single line item for Product A

**Traceability:**
- AC: AC-CART-004

---

## TC-AC-CART-004-03 – Cart merge - merged quantity exceeds stock

**Type:** error_case

**Preconditions:**
- Guest and account carts have same product
- Combined quantity exceeds stock

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| guest_cart_items | [map[product:PROD-A qty:6]] | Guest cart |
| account_cart_items | [map[product:PROD-A qty:7]] | Account cart |
| available_stock | 10 | Max available |

**Steps:**
1. As guest, add Product A (qty 6) to cart
2. Log in to account that has Product A (qty 7) in cart

**Expected Result:**
- Cart contains Product A with quantity limited to 10
- Notification shown about stock limitation

**Traceability:**
- AC: AC-CART-004

---

## TC-AC-CART-005-01 – Update cart item quantity within valid range

**Type:** happy_path

**Preconditions:**
- Customer has product in cart
- Product has sufficient stock

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| product_id | PROD-001 | Product in cart |
| current_quantity | 2 | Current quantity |
| new_quantity | 5 | Updated quantity |
| available_stock | 10 | Sufficient stock |

**Steps:**
1. Navigate to cart
2. Update quantity from 2 to 5
3. Observe cart update

**Expected Result:**
- Cart item quantity updates to 5
- Cart total recalculates
- Current product price is used in calculation

**Traceability:**
- AC: AC-CART-005

---

## TC-AC-CART-005-02 – Update cart quantity to zero removes item

**Type:** error_case

**Preconditions:**
- Customer has product in cart

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| product_id | PROD-001 | Product in cart |
| current_quantity | 2 | Current quantity |
| new_quantity | 0 | Zero quantity |

**Steps:**
1. Navigate to cart
2. Set quantity to 0

**Expected Result:**
- Item is removed from cart
- Cart total recalculates

**Traceability:**
- AC: AC-CART-005

---

## TC-AC-CART-005-03 – Update cart quantity exceeds available stock

**Type:** error_case

**Preconditions:**
- Customer has product in cart

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| product_id | PROD-001 | Product in cart |
| new_quantity | 20 | Exceeds stock |
| available_stock | 10 | Max available |

**Steps:**
1. Navigate to cart
2. Attempt to set quantity to 20

**Expected Result:**
- Quantity limited to available stock (10)
- Notification shown about stock limitation

**Traceability:**
- AC: AC-CART-005

---

## TC-AC-CART-005-04 – Update cart with negative quantity

**Type:** error_case

**Preconditions:**
- Customer has product in cart

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| product_id | PROD-001 | Product in cart |
| new_quantity | -1 | Invalid negative quantity |

**Steps:**
1. Navigate to cart
2. Attempt to set quantity to -1

**Expected Result:**
- Validation error is displayed
- Quantity remains unchanged

**Traceability:**
- AC: AC-CART-005

---

## TC-AC-CART-006-01 – Remove item from cart with undo option

**Type:** happy_path

**Preconditions:**
- Customer has multiple items in cart

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| cart_items | 3 | Multiple items |
| remove_product | PROD-001 | Item to remove |

**Steps:**
1. Navigate to cart
2. Click remove on PROD-001

**Expected Result:**
- Item is immediately removed
- No confirmation dialog shown
- 'Undo' option appears temporarily
- Cart total recalculates

**Traceability:**
- AC: AC-CART-006

---

## TC-AC-CART-007-01 – Remove last item shows empty cart state

**Type:** happy_path

**Preconditions:**
- Customer has exactly one item in cart

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| cart_items | 1 | Single item |
| product_id | PROD-001 | Only item |

**Steps:**
1. Navigate to cart with single item
2. Click remove on the item

**Expected Result:**
- Item is removed
- Empty cart state is displayed
- 'Continue Shopping' link is shown

**Traceability:**
- AC: AC-CART-007

---

## TC-AC-CART-008-01 – Cart item goes out of stock

**Type:** happy_path

**Preconditions:**
- Customer has product in cart
- Product subsequently goes out of stock

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| product_id | PROD-001 | Initially in stock |
| new_stock | 0 | Becomes out of stock |

**Steps:**
1. Add in-stock product to cart
2. Product goes out of stock (inventory updated to 0)
3. View cart

**Expected Result:**
- Item remains in cart
- Warning indicator is displayed on item
- Checkout is blocked until resolved

**Traceability:**
- AC: AC-CART-008

---

## TC-AC-CART-009-01 – Price change notification in cart

**Type:** happy_path

**Preconditions:**
- Customer has product in cart
- Product price changes

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| product_id | PROD-001 | Product in cart |
| original_price | 29.99 | Price when added |
| new_price | 34.99 | Updated price |

**Steps:**
1. Add product at $29.99 to cart
2. Product price changes to $34.99
3. View cart

**Expected Result:**
- Cart shows current price ($34.99)
- Notification indicates price has changed

**Traceability:**
- AC: AC-CART-009

---

## TC-AC-ORDER-001-01 – Place order as registered customer - success

**Type:** happy_path

**Preconditions:**
- Registered customer with verified email
- Items in cart
- Valid shipping address
- Valid payment method

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| customer_email | verified@test.com | Verified email |
| cart_subtotal | 60 | Before tax |
| shipping_address | 123 Main St, NY 10001 | Valid US address |
| payment_method | Visa ending 4242 | Valid card |

**Steps:**
1. Log in as registered customer with verified email
2. Navigate to checkout with items in cart
3. Confirm shipping address
4. Confirm payment method
5. Submit order

**Expected Result:**
- Payment is authorized
- Inventory is decremented
- Order created with status 'pending'
- Order number format is ORD-YYYY-NNNNN
- Cart is cleared
- Confirmation email is queued

**Traceability:**
- AC: AC-ORDER-001

---

## TC-AC-ORDER-001-02 – Place order - customer not registered

**Type:** error_case

**Preconditions:**
- Guest user with items in cart

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| user_type | guest | Not registered |

**Steps:**
1. As guest user, add items to cart
2. Attempt to proceed to checkout

**Expected Result:**
- REGISTRATION_REQUIRED error is returned
- User is redirected to registration

**Traceability:**
- AC: AC-ORDER-001

---

## TC-AC-ORDER-001-03 – Place order - email not verified

**Type:** error_case

**Preconditions:**
- Registered customer with unverified email
- Items in cart

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| customer_email | unverified@test.com | Not verified |
| email_verified | false | Pending verification |

**Steps:**
1. Log in as customer with unverified email
2. Attempt to place order

**Expected Result:**
- EMAIL_NOT_VERIFIED error is returned
- Prompt to verify email is displayed

**Traceability:**
- AC: AC-ORDER-001

---

## TC-AC-ORDER-001-04 – Place order - payment authorization fails

**Type:** error_case

**Preconditions:**
- Registered customer with verified email
- Items in cart
- Payment method will fail

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| payment_method | Declined card | Card will be declined |

**Steps:**
1. Proceed to checkout
2. Submit order with card that will be declined

**Expected Result:**
- PAYMENT_FAILED error is returned
- Cart is preserved
- Option to retry with different payment is available

**Traceability:**
- AC: AC-ORDER-001

---

## TC-AC-ORDER-001-05 – Place order - item out of stock at checkout

**Type:** error_case

**Preconditions:**
- Customer in checkout
- Item goes out of stock during checkout

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| product_id | PROD-001 | Becomes out of stock |
| available_stock | 0 | Sold out |

**Steps:**
1. Proceed to checkout with items
2. Item goes out of stock
3. Submit order

**Expected Result:**
- INSUFFICIENT_STOCK error for affected items
- User can modify cart

**Traceability:**
- AC: AC-ORDER-001

---

## TC-AC-ORDER-001-06 – Place order - invalid shipping address

**Type:** error_case

**Preconditions:**
- Customer in checkout

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| shipping_address | Invalid address | Missing required fields |

**Steps:**
1. Proceed to checkout
2. Enter invalid shipping address
3. Attempt to submit

**Expected Result:**
- ADDRESS_VALIDATION_ERROR is returned
- Invalid fields are highlighted

**Traceability:**
- AC: AC-ORDER-001

---

## TC-AC-ORDER-002-01 – Free shipping for orders $50 or more

**Type:** happy_path

**Preconditions:**
- Customer has items in cart totaling $50 or more

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| cart_subtotal | 50 | Exactly at threshold |
| discount | 0 | No discounts applied |

**Steps:**
1. Add items totaling $50 to cart
2. Proceed to checkout

**Expected Result:**
- Free shipping is automatically applied
- Shipping cost shows $0.00

**Traceability:**
- AC: AC-ORDER-002

---

## TC-AC-ORDER-002-02 – Free shipping after discount still qualifies

**Type:** edge_case

**Preconditions:**
- Customer has items totaling $60 before discount

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| cart_subtotal_before_discount | 60 | Before discount |
| discount | 5 | Applied discount |
| cart_subtotal_after_discount | 55 | After discount, still >= $50 |

**Steps:**
1. Add items totaling $60 to cart
2. Apply $5 discount
3. Proceed to checkout

**Expected Result:**
- Free shipping is applied
- Subtotal after discount ($55) qualifies for free shipping

**Traceability:**
- AC: AC-ORDER-002

---

## TC-AC-ORDER-003-01 – Standard shipping charge for orders under $50

**Type:** happy_path

**Preconditions:**
- Customer has items in cart totaling less than $50

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| cart_subtotal | 45 | Below threshold |

**Steps:**
1. Add items totaling $45 to cart
2. Proceed to checkout

**Expected Result:**
- Standard shipping fee of $5.99 is applied

**Traceability:**
- AC: AC-ORDER-003

---

## TC-AC-ORDER-004-01 – Tax calculation for US shipping address

**Type:** happy_path

**Preconditions:**
- Customer has items in cart
- US shipping address entered

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| shipping_state | CA | California |
| cart_subtotal | 100 | Base amount |

**Steps:**
1. Add items to cart totaling $100
2. Enter California shipping address
3. View order total

**Expected Result:**
- State tax is calculated based on California tax rate
- Tax amount is displayed in order summary

**Traceability:**
- AC: AC-ORDER-004

---

## TC-AC-ORDER-005-01 – View order history with pagination

**Type:** happy_path

**Preconditions:**
- Logged-in customer has placed multiple orders

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| total_orders | 25 | Multiple orders |

**Steps:**
1. Log in as customer with order history
2. Navigate to order history

**Expected Result:**
- Orders are displayed with pagination
- Each order shows: order number, items, quantities, prices, status, shipping address
- Tracking number shown if available

**Traceability:**
- AC: AC-ORDER-005

---

## TC-AC-ORDER-006-01 – Cancel pending order successfully

**Type:** happy_path

**Preconditions:**
- Customer has order with status 'pending'

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| order_id | ORD-2025-00001 | Pending order |
| order_status | pending | Cancellable status |
| cancellation_reason | Changed my mind | Optional reason |

**Steps:**
1. Navigate to order details for pending order
2. Click cancel order
3. Enter optional cancellation reason
4. Confirm cancellation

**Expected Result:**
- Order status changes to 'cancelled'
- Payment is refunded to original method
- Inventory is restored
- Cancellation email is sent

**Traceability:**
- AC: AC-ORDER-006

---

## TC-AC-ORDER-006-02 – Cancel confirmed order successfully

**Type:** happy_path

**Preconditions:**
- Customer has order with status 'confirmed'

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| order_id | ORD-2025-00002 | Confirmed order |
| order_status | confirmed | Cancellable status |

**Steps:**
1. Navigate to order details for confirmed order
2. Click cancel order
3. Confirm cancellation

**Expected Result:**
- Order status changes to 'cancelled'
- Payment is refunded
- Inventory is restored
- Cancellation email is sent

**Traceability:**
- AC: AC-ORDER-006

---

## TC-AC-ORDER-006-03 – Cannot cancel shipped order

**Type:** error_case

**Preconditions:**
- Customer has order with status 'shipped'

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| order_id | ORD-2025-00003 | Shipped order |
| order_status | shipped | Non-cancellable status |

**Steps:**
1. Navigate to order details for shipped order
2. Attempt to cancel order

**Expected Result:**
- CANCELLATION_NOT_ALLOWED error is returned
- Suggestion to contact support is displayed

**Traceability:**
- AC: AC-ORDER-006

---

## TC-AC-ORDER-007-01 – Order confirmation email sent successfully

**Type:** happy_path

**Preconditions:**
- Order has been successfully placed

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| order_id | ORD-2025-00001 | New order |
| customer_email | customer@test.com | Verified email |

**Steps:**
1. Place an order successfully

**Expected Result:**
- Confirmation email is queued asynchronously
- Email contains: order number, items, total, shipping address, estimated delivery

**Traceability:**
- AC: AC-ORDER-007

---

## TC-AC-ORDER-007-02 – Email service unavailable - order not blocked

**Type:** error_case

**Preconditions:**
- Email service is unavailable

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| email_service_status | unavailable | Service down |

**Steps:**
1. Place an order when email service is down

**Expected Result:**
- Order creation is not blocked
- Email is queued for retry
- Failure is logged

**Traceability:**
- AC: AC-ORDER-007

---

## TC-AC-CUST-001-01 – Customer registration with valid data

**Type:** happy_path

**Preconditions:**
- Guest user on registration page

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| email | newuser@test.com | Valid unique email |
| password | SecurePass1 | 8+ chars with number |
| first_name | John | Required |
| last_name | Doe | Required |
| phone | +1234567890 | Optional |

**Steps:**
1. Navigate to registration page
2. Enter valid registration data
3. Submit registration form

**Expected Result:**
- Account is created
- Verification email is sent
- User cannot place orders until email is verified

**Traceability:**
- AC: AC-CUST-001

---

## TC-AC-CUST-001-02 – Registration - duplicate email

**Type:** error_case

**Preconditions:**
- Email already registered in system

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| email | existing@test.com | Already registered |

**Steps:**
1. Navigate to registration page
2. Enter email that already exists
3. Submit registration form

**Expected Result:**
- DUPLICATE_EMAIL error is returned

**Traceability:**
- AC: AC-CUST-001

---

## TC-AC-CUST-001-03 – Registration - password too short

**Type:** error_case

**Preconditions:**
- Guest user on registration page

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| password | Short1 | Less than 8 characters |

**Steps:**
1. Navigate to registration page
2. Enter password with less than 8 characters
3. Submit registration form

**Expected Result:**
- PASSWORD_TOO_SHORT error is returned

**Traceability:**
- AC: AC-CUST-001

---

## TC-AC-CUST-001-04 – Registration - password missing number

**Type:** error_case

**Preconditions:**
- Guest user on registration page

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| password | SecurePassword | No number included |

**Steps:**
1. Navigate to registration page
2. Enter password without a number
3. Submit registration form

**Expected Result:**
- PASSWORD_MISSING_NUMBER error is returned

**Traceability:**
- AC: AC-CUST-001

---

## TC-AC-CUST-001-05 – Registration - invalid email format

**Type:** error_case

**Preconditions:**
- Guest user on registration page

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| email | invalid-email | Not valid email format |

**Steps:**
1. Navigate to registration page
2. Enter invalid email format
3. Submit registration form

**Expected Result:**
- INVALID_EMAIL error is returned

**Traceability:**
- AC: AC-CUST-001

---

## TC-AC-CUST-002-01 – Add new shipping address

**Type:** happy_path

**Preconditions:**
- Registered customer is logged in

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| street | 123 Main Street | Required |
| city | New York | Required |
| state | NY | Required |
| zip | 10001 | Required |
| country | US | Domestic |

**Steps:**
1. Navigate to address management
2. Click add new address
3. Enter all required fields
4. Save address

**Expected Result:**
- Address is saved
- Option to set as default is available

**Traceability:**
- AC: AC-CUST-002

---

## TC-AC-CUST-002-02 – Add shipping address - missing required field

**Type:** error_case

**Preconditions:**
- Registered customer is logged in

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| street |  | Missing required field |
| city | New York | Provided |

**Steps:**
1. Navigate to address management
2. Attempt to save address with missing street

**Expected Result:**
- INVALID_ADDRESS error with field specification

**Traceability:**
- AC: AC-CUST-002

---

## TC-AC-CUST-002-03 – Add shipping address - unsupported country

**Type:** error_case

**Preconditions:**
- Registered customer is logged in

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| country | FR | Non-domestic country |

**Steps:**
1. Navigate to address management
2. Attempt to add international address

**Expected Result:**
- UNSUPPORTED_COUNTRY error is returned

**Traceability:**
- AC: AC-CUST-002

---

## TC-AC-CUST-003-01 – Add payment method - Visa card

**Type:** happy_path

**Preconditions:**
- Registered customer is logged in

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| card_type | Visa | Supported card type |
| card_number | 4242424242424242 | Valid test card |

**Steps:**
1. Navigate to payment methods
2. Click add new payment method
3. Enter Visa card details
4. Save payment method

**Expected Result:**
- Payment method is tokenized via Stripe
- Card is stored
- Option to set as default is available

**Traceability:**
- AC: AC-CUST-003

---

## TC-AC-CUST-003-02 – Add payment method - PayPal

**Type:** happy_path

**Preconditions:**
- Registered customer is logged in

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| payment_type | PayPal | Supported payment method |

**Steps:**
1. Navigate to payment methods
2. Select PayPal option
3. Complete PayPal authorization

**Expected Result:**
- PayPal account is linked
- Payment method is stored

**Traceability:**
- AC: AC-CUST-003

---

## TC-AC-CUST-003-03 – Add payment method - invalid card number

**Type:** error_case

**Preconditions:**
- Registered customer is logged in

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| card_number | 1234567890123456 | Invalid card number |

**Steps:**
1. Navigate to payment methods
2. Enter invalid card number
3. Attempt to save

**Expected Result:**
- INVALID_CARD error is returned

**Traceability:**
- AC: AC-CUST-003

---

## TC-AC-CUST-003-04 – Add payment method - unsupported card type

**Type:** error_case

**Preconditions:**
- Registered customer is logged in

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| card_type | Discover | Not supported |

**Steps:**
1. Navigate to payment methods
2. Attempt to add Discover card

**Expected Result:**
- UNSUPPORTED_CARD_TYPE error is returned

**Traceability:**
- AC: AC-CUST-003

---

## TC-AC-ADMIN-001-01 – Create new product with valid data

**Type:** happy_path

**Preconditions:**
- Admin user is logged in
- On product management page

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| name | New Test Product | 2-200 chars |
| price | 29.99 | >$0.01 |
| category | Electronics | Required |
| description | Product description text | Up to 5000 chars |
| images | 3 | Up to 10 images |

**Steps:**
1. Navigate to product management
2. Click create new product
3. Enter valid product data
4. Save product

**Expected Result:**
- Product is created with 'draft' status
- UUID is assigned
- Creator and timestamp are logged

**Traceability:**
- AC: AC-ADMIN-001

---

## TC-AC-ADMIN-001-02 – Create product - name too short

**Type:** error_case

**Preconditions:**
- Admin user is logged in

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| name | A | Less than 2 characters |

**Steps:**
1. Attempt to create product with single character name

**Expected Result:**
- INVALID_NAME error is returned

**Traceability:**
- AC: AC-ADMIN-001

---

## TC-AC-ADMIN-001-03 – Create product - name too long

**Type:** error_case

**Preconditions:**
- Admin user is logged in

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| name | [201 character string] | Exceeds 200 characters |

**Steps:**
1. Attempt to create product with name over 200 characters

**Expected Result:**
- NAME_TOO_LONG error is returned

**Traceability:**
- AC: AC-ADMIN-001

---

## TC-AC-ADMIN-001-04 – Create product - price zero or negative

**Type:** error_case

**Preconditions:**
- Admin user is logged in

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| price | 0 | Invalid price |

**Steps:**
1. Attempt to create product with price of $0

**Expected Result:**
- INVALID_PRICE error is returned

**Traceability:**
- AC: AC-ADMIN-001

---

## TC-AC-ADMIN-001-05 – Create product - no category selected

**Type:** error_case

**Preconditions:**
- Admin user is logged in

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| category | <nil> | No category |

**Steps:**
1. Attempt to create product without selecting category

**Expected Result:**
- CATEGORY_REQUIRED error is returned

**Traceability:**
- AC: AC-ADMIN-001

---

## TC-AC-ADMIN-001-06 – Create product - duplicate name warning

**Type:** edge_case

**Preconditions:**
- Admin user is logged in
- Product with same name exists

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| name | Existing Product Name | Already exists |

**Steps:**
1. Attempt to create product with name that already exists

**Expected Result:**
- Warning is displayed
- Creation is not blocked

**Traceability:**
- AC: AC-ADMIN-001

---

## TC-AC-ADMIN-002-01 – Edit product attributes

**Type:** happy_path

**Preconditions:**
- Admin user is logged in
- Product exists

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| product_id | PROD-001 | Existing product |
| original_price | 29.99 | Before edit |
| new_price | 34.99 | After edit |

**Steps:**
1. Navigate to product edit page
2. Modify product price
3. Save changes

**Expected Result:**
- Changes are saved
- Last modified by/date is updated
- Price change is logged
- Existing orders retain original price
- Carts get updated price

**Traceability:**
- AC: AC-ADMIN-002

---

## TC-AC-ADMIN-002-02 – Edit product - concurrent edit conflict

**Type:** error_case

**Preconditions:**
- Two admins editing same product

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| product_id | PROD-001 | Same product |
| admin_1 | Admin A | First editor |
| admin_2 | Admin B | Second editor |

**Steps:**
1. Admin A opens product for editing
2. Admin B opens same product for editing
3. Admin A saves changes
4. Admin B saves changes

**Expected Result:**
- Last write wins
- Conflict warning is displayed

**Traceability:**
- AC: AC-ADMIN-002

---

## TC-AC-ADMIN-003-01 – Remove product with soft delete

**Type:** happy_path

**Preconditions:**
- Admin user is logged in
- Product exists

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| product_id | PROD-001 | Product to remove |

**Steps:**
1. Navigate to product management
2. Select product to remove
3. View impact summary
4. Confirm removal

**Expected Result:**
- Product is soft-deleted
- Product removed from all active carts
- Users notified on next cart view
- Product preserved for order history

**Traceability:**
- AC: AC-ADMIN-003

---

## TC-AC-ADMIN-004-01 – View all orders with filters

**Type:** happy_path

**Preconditions:**
- Admin user is logged in
- Orders exist in system

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| status_filter | pending | Filter by status |
| date_range | last_30_days | Date filter |

**Steps:**
1. Navigate to order management
2. Apply status filter for 'pending'
3. Apply date range filter

**Expected Result:**
- Matching orders are displayed
- Shows: order number, date, customer, total, status
- Results are paginated

**Traceability:**
- AC: AC-ADMIN-004

---

## TC-AC-ADMIN-004-02 – Search orders by order number

**Type:** happy_path

**Preconditions:**
- Admin user is logged in
- Order exists

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| search_term | ORD-2025-00001 | Specific order number |

**Steps:**
1. Navigate to order management
2. Enter order number in search
3. Submit search

**Expected Result:**
- Matching order is displayed

**Traceability:**
- AC: AC-ADMIN-004

---

## TC-AC-ADMIN-005-01 – Update order status from pending to confirmed

**Type:** happy_path

**Preconditions:**
- Admin user is logged in
- Order with status 'pending' exists

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| order_id | ORD-2025-00001 | Pending order |
| current_status | pending | Current status |
| new_status | confirmed | Target status |

**Steps:**
1. Navigate to order details
2. Update status to 'confirmed'

**Expected Result:**
- Status changes to confirmed
- Customer receives email notification
- Change is logged with admin user and timestamp

**Traceability:**
- AC: AC-ADMIN-005

---

## TC-AC-ADMIN-005-02 – Invalid status transition

**Type:** error_case

**Preconditions:**
- Admin user is logged in
- Order with status 'pending' exists

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| order_id | ORD-2025-00001 | Pending order |
| current_status | pending | Current status |
| new_status | delivered | Invalid target |

**Steps:**
1. Navigate to order details
2. Attempt to update status from 'pending' directly to 'delivered'

**Expected Result:**
- INVALID_STATUS_TRANSITION error is returned

**Traceability:**
- AC: AC-ADMIN-005

---

## TC-AC-ADMIN-006-01 – Mark order as shipped with tracking

**Type:** happy_path

**Preconditions:**
- Admin user is logged in
- Order with status 'confirmed' exists

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| order_id | ORD-2025-00001 | Confirmed order |
| tracking_number | 1Z999AA10123456784 | UPS tracking |

**Steps:**
1. Navigate to order details
2. Update status to 'shipped'
3. Enter tracking number

**Expected Result:**
- Status changes to shipped
- Tracking number is stored
- Customer receives email with tracking link to carrier website

**Traceability:**
- AC: AC-ADMIN-006

---

## TC-AC-ADMIN-007-01 – Update inventory with absolute adjustment

**Type:** happy_path

**Preconditions:**
- Admin user is logged in
- Product variant exists

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| variant_sku | SKU-001-M | Product variant |
| current_stock | 50 | Before update |
| new_stock | 75 | Absolute value |
| reason | Inventory recount | Optional reason |

**Steps:**
1. Navigate to inventory management
2. Select product variant
3. Set stock to 75 (absolute)
4. Enter reason
5. Save changes

**Expected Result:**
- Stock is updated to 75
- Audit trail records before/after values
- Timestamp and admin user logged

**Traceability:**
- AC: AC-ADMIN-007

---

## TC-AC-ADMIN-007-02 – Update inventory with delta adjustment

**Type:** happy_path

**Preconditions:**
- Admin user is logged in
- Product variant exists

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| variant_sku | SKU-001-M | Product variant |
| current_stock | 50 | Before update |
| delta | -10 | Decrement by 10 |

**Steps:**
1. Navigate to inventory management
2. Select product variant
3. Apply delta adjustment of -10
4. Save changes

**Expected Result:**
- Stock is updated to 40
- Audit trail records change

**Traceability:**
- AC: AC-ADMIN-007

---

## TC-AC-ADMIN-007-03 – Low stock alert triggered

**Type:** happy_path

**Preconditions:**
- Admin user is logged in
- Low stock threshold set to 10

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| variant_sku | SKU-001-M | Product variant |
| current_stock | 15 | Above threshold |
| new_stock | 8 | Below threshold |

**Steps:**
1. Update stock from 15 to 8
2. Save changes

**Expected Result:**
- Stock is updated
- Low stock alert is triggered

**Traceability:**
- AC: AC-ADMIN-007

---

## TC-AC-ADMIN-007-04 – Update inventory - result would be negative

**Type:** error_case

**Preconditions:**
- Admin user is logged in

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| current_stock | 5 | Current stock |
| delta | -10 | Would result in -5 |

**Steps:**
1. Attempt to apply delta of -10 to stock of 5

**Expected Result:**
- INVALID_STOCK_LEVEL error is returned
- Minimum stock is zero

**Traceability:**
- AC: AC-ADMIN-007

---

## TC-AC-ADMIN-008-01 – Bulk inventory update via CSV

**Type:** happy_path

**Preconditions:**
- Admin user is logged in
- Valid CSV file prepared

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| csv_file | inventory_update.csv | SKU and quantity columns |
| rows | 50 | Number of records |

**Steps:**
1. Navigate to inventory management
2. Select bulk update option
3. Upload CSV file
4. Confirm import

**Expected Result:**
- All valid rows are processed
- Audit trail created for each change
- Summary shows success count

**Traceability:**
- AC: AC-ADMIN-008

---

## TC-AC-ADMIN-008-02 – Bulk inventory update - invalid SKU rows

**Type:** error_case

**Preconditions:**
- Admin user is logged in
- CSV with invalid SKUs

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| invalid_sku | INVALID-SKU | Non-existent SKU |

**Steps:**
1. Upload CSV with some invalid SKUs

**Expected Result:**
- Invalid SKU rows are skipped
- Summary shows errors for invalid rows
- Valid rows are still processed

**Traceability:**
- AC: AC-ADMIN-008

---

## TC-AC-ADMIN-008-03 – Bulk inventory update - invalid quantity format

**Type:** error_case

**Preconditions:**
- Admin user is logged in

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| quantity | abc | Non-numeric value |

**Steps:**
1. Upload CSV with non-numeric quantity values

**Expected Result:**
- Rows with invalid quantity are skipped
- Error shown in summary

**Traceability:**
- AC: AC-ADMIN-008

---

## TC-AC-VAR-001-01 – Add product variant with unique SKU

**Type:** happy_path

**Preconditions:**
- Admin is editing a product

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| sku | PROD-001-RED-M | Unique SKU |
| size | M | Variant attribute |
| color | Red | Variant attribute |
| price_override | <nil> | Inherits parent price |

**Steps:**
1. Navigate to product edit page
2. Add new variant
3. Enter unique SKU and attributes
4. Save variant

**Expected Result:**
- Variant is created with own inventory tracking
- Inherits parent price (no override)

**Traceability:**
- AC: AC-VAR-001

---

## TC-AC-VAR-001-02 – Add variant with price override

**Type:** happy_path

**Preconditions:**
- Admin is editing a product with base price $29.99

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| sku | PROD-001-RED-XL | Unique SKU |
| price_override | 34.99 | Higher price for XL |

**Steps:**
1. Add new variant
2. Enter SKU and attributes
3. Set price override to $34.99
4. Save variant

**Expected Result:**
- Variant created with override price of $34.99

**Traceability:**
- AC: AC-VAR-001

---

## TC-AC-VAR-001-03 – Add variant - duplicate SKU

**Type:** error_case

**Preconditions:**
- SKU already exists in system

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| sku | EXISTING-SKU | Already exists |

**Steps:**
1. Attempt to create variant with existing SKU

**Expected Result:**
- DUPLICATE_SKU error is returned

**Traceability:**
- AC: AC-VAR-001

---

## TC-AC-CAT-001-01 – Create category with all fields

**Type:** happy_path

**Preconditions:**
- Admin is on category management

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| name | Electronics | Category name |
| description | Electronic devices | Description |
| display_order | 1 | Sort order |
| status | active | Active status |

**Steps:**
1. Navigate to category management
2. Create new category
3. Enter all fields
4. Save category

**Expected Result:**
- Category is created with all attributes
- Category is active

**Traceability:**
- AC: AC-CAT-001

---

## TC-AC-CAT-001-02 – Create nested category up to 3 levels

**Type:** happy_path

**Preconditions:**
- Level 1 and Level 2 categories exist

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| parent_level_1 | Electronics | Top level |
| parent_level_2 | Phones | Level 2 |
| new_category | Smartphones | Level 3 |

**Steps:**
1. Create category 'Smartphones' under 'Phones'

**Expected Result:**
- Level 3 category is created
- Hierarchy: Electronics > Phones > Smartphones

**Traceability:**
- AC: AC-CAT-001

---

## TC-AC-CAT-001-03 – Create category - nesting exceeds 3 levels

**Type:** error_case

**Preconditions:**
- 3 levels of categories already exist

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| parent | Level 3 Category | Already at max depth |

**Steps:**
1. Attempt to create category under level 3 category

**Expected Result:**
- CATEGORY_DEPTH_EXCEEDED error is returned

**Traceability:**
- AC: AC-CAT-001

---

## TC-AC-GDPR-001-01 – Customer data erasure request

**Type:** happy_path

**Preconditions:**
- Customer account exists with orders

**Test Data:**
| Field | Value | Notes |
|-------|-------|-------|
| customer_id | CUST-001 | Customer requesting deletion |
| orders | 5 | Historical orders |

**Steps:**
1. Customer submits GDPR data deletion request
2. Request is processed

**Expected Result:**
- All personal data is erased
- Order history is anonymized
- Business records preserved with anonymized data

**Traceability:**
- AC: AC-GDPR-001

---

