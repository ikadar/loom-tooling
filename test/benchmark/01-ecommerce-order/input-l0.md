# E-Commerce Order System

## Overview

A simple online store where customers can browse products, add them to cart, and place orders.

## User Stories

### US-001: Browse Products
As a customer, I want to browse available products so that I can decide what to buy.

Products have a name, description, price, and category. Some products may have variants (like size or color).

### US-002: Add to Cart
As a customer, I want to add products to my shopping cart so that I can purchase multiple items at once.

The cart should show the items, quantities, and total price. Customers can update quantities or remove items.

### US-003: Place Order
As a customer, I want to place an order for the items in my cart so that I can receive the products.

When placing an order, the customer provides:
- Shipping address
- Payment method (credit card or PayPal)

After the order is placed, the customer receives a confirmation email.

### US-004: Track Order Status
As a customer, I want to see the status of my orders so that I know when they will arrive.

Order statuses include: pending, confirmed, shipped, delivered.

### US-005: Cancel Order
As a customer, I want to cancel my order if I change my mind.

## Business Rules

- Customers must be registered to place orders
- Products must be in stock to be added to cart
- Orders over $50 qualify for free shipping
- Orders can only be cancelled before shipping

## Admin Features

- Add/edit/remove products
- View all orders
- Update order status
- Manage inventory
