# L2 to L3 Derivation Prompt

You are an expert software architect. Generate L3 implementation artifacts from L2 inputs.

OUTPUT REQUIREMENT: Wrap your JSON response in ```json code blocks. No explanations.

## Your Task

From Test Cases (TC) and Technical Specifications (TS), generate:
1. **API Specification** - OpenAPI-style endpoint definitions
2. **Implementation Skeletons** - Function signatures and key logic outlines

## API Endpoint Structure

For each unique endpoint discovered in tech specs:
- method: HTTP method
- path: URL path with parameters
- summary: Brief description
- request_body: Input schema
- responses: Success and error responses
- related_specs: Tech spec IDs

## Implementation Skeleton Structure

For each major component/service:
- name: Component/service name
- functions: Key function signatures
- dependencies: Required services/modules
- related_specs: Tech spec IDs

## Output Format

```json
{
  "api_spec": {
    "openapi": "3.0.0",
    "info": {"title": "Generated API", "version": "1.0.0"},
    "paths": {
      "/api/cart/items": {
        "post": {
          "summary": "Add item to cart",
          "operationId": "addCartItem",
          "requestBody": {
            "content": {
              "application/json": {
                "schema": {
                  "type": "object",
                  "properties": {
                    "product_id": {"type": "string"},
                    "quantity": {"type": "integer"}
                  },
                  "required": ["product_id", "quantity"]
                }
              }
            }
          },
          "responses": {
            "201": {"description": "Item added"},
            "400": {"description": "Validation error"}
          },
          "x-related-specs": ["TS-BR-STOCK-001"]
        }
      }
    }
  },
  "implementation_skeletons": [
    {
      "name": "CartService",
      "type": "service",
      "functions": [
        {
          "name": "addItem",
          "signature": "addItem(customerId: string, productId: string, quantity: number): CartItem",
          "description": "Add item to cart with stock validation",
          "steps": [
            "Validate stock availability",
            "Check for existing cart item",
            "Create or update cart item",
            "Return updated item"
          ],
          "error_cases": ["OUT_OF_STOCK", "INVALID_QUANTITY"]
        }
      ],
      "dependencies": ["InventoryService", "ProductService"],
      "related_specs": ["TS-BR-STOCK-001", "TS-BR-CART-001"]
    }
  ],
  "summary": {
    "endpoints_count": 10,
    "services_count": 5,
    "functions_count": 20
  }
}
```

---

REMINDER: Output ONLY a ```json code block. No explanations.

L2 INPUT (Test Cases + Tech Specs):
