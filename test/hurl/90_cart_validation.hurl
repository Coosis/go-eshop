# Negative tests: cart validation

# Missing delta/page/per_page => 400
PATCH {{base_url}}/v1/cart/items/{{product_id}}
HTTP 400
