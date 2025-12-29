# Sequence Design Derivation Prompt

You are an expert system architect. Generate Sequence Designs from L1 and L2 documents.

OUTPUT REQUIREMENT: Wrap your JSON response in ```json code blocks. No explanations.

## Your Task

From Domain Model, Interface Contracts, and Business Rules, derive Sequence Designs that show:
1. **End-to-end flows** - How services collaborate to fulfill business processes
2. **Participants** - All components involved in each sequence
3. **Steps** - Chronological interactions with events and commands
4. **Outcomes** - Final state after sequence completion

## Sequence Structure

For each sequence, include:
- id: Unique identifier (SEQ-{FLOW}-{NNN})
- name: Descriptive name
- trigger: What starts the sequence
- participants: Services and aggregates involved
- steps: Ordered list of interactions
- outcome: Final state description
- exceptions: Alternative paths for errors

## Output Format

```json
{
  "sequences": [
    {
      "id": "SEQ-ORDER-001",
      "name": "Place Order Flow",
      "description": "Complete flow from cart checkout to order confirmation",
      "trigger": {
        "type": "user_action",
        "description": "Customer clicks 'Place Order' button"
      },
      "participants": [
        {"name": "Customer", "type": "actor"},
        {"name": "Cart Service", "type": "service"},
        {"name": "Order Service", "type": "service"},
        {"name": "Inventory Service", "type": "service"},
        {"name": "Payment Service", "type": "service"}
      ],
      "steps": [
        {
          "step": 1,
          "actor": "Customer",
          "action": "Submit order request",
          "target": "Order Service",
          "data": ["cartId", "shippingAddress", "paymentMethod"]
        },
        {
          "step": 2,
          "actor": "Order Service",
          "action": "Fetch cart contents",
          "target": "Cart Service",
          "returns": "Cart with items"
        },
        {
          "step": 3,
          "actor": "Order Service",
          "action": "Reserve inventory for each item",
          "target": "Inventory Service",
          "event": "InventoryReserved"
        },
        {
          "step": 4,
          "actor": "Order Service",
          "action": "Create order with status 'pending'",
          "target": "Order Aggregate",
          "event": "OrderCreated"
        },
        {
          "step": 5,
          "actor": "Order Service",
          "action": "Clear customer cart",
          "target": "Cart Service",
          "event": "CartCleared"
        },
        {
          "step": 6,
          "actor": "Order Service",
          "action": "Return order confirmation",
          "target": "Customer",
          "returns": "Order details with orderId"
        }
      ],
      "outcome": {
        "success": "Order created with 'pending' status, inventory reserved, cart emptied",
        "state_changes": ["Order.status = pending", "Inventory.reserved += quantities", "Cart.items = empty"]
      },
      "exceptions": [
        {
          "condition": "Cart is empty",
          "step": 2,
          "handling": "Return error CART_EMPTY, abort flow"
        },
        {
          "condition": "Insufficient stock",
          "step": 3,
          "handling": "Return error INSUFFICIENT_STOCK with affected items, abort flow"
        }
      ],
      "relatedACs": ["AC-ORDER-001", "AC-ORDER-002"],
      "relatedBRs": ["BR-ORDER-001", "BR-INV-001"]
    }
  ],
  "summary": {
    "total_sequences": 8,
    "sequences_by_domain": {
      "order": 3,
      "cart": 2,
      "inventory": 2,
      "customer": 1
    }
  }
}
```

## Important Sequences to Identify

Look for these common flows:
- Create/Place Order
- Cancel Order
- Update Cart
- Process Payment
- Register Customer
- Update Inventory
- Ship Order

## Quality Checklist

Before output, verify:
- [ ] All participants are identified
- [ ] Steps are in logical order
- [ ] Events match domain model events
- [ ] Exception paths are documented
- [ ] State changes are explicit

---

REMINDER: Output ONLY a ```json code block. No explanations.

L1/L2 INPUT (Domain Model + Interface Contracts + Business Rules):

