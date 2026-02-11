# Cart + Orders flow

GET {{base_url}}/v1/cart?page=1&per_page=20
HTTP 200
[Captures]
cart_version: jsonpath "$.version"

POST {{base_url}}/v1/cart/items
Content-Type: application/json
{
  "product_id": {{product_id}},
  "quantity": 1,
  "page": 1,
  "per_page": 20
}
HTTP 200

GET {{base_url}}/v1/cart?page=1&per_page=20
HTTP 200
[Captures]
cart_version: jsonpath "$.version"

PUT {{base_url}}/v1/cart/items/{{product_id}}
Content-Type: application/json
{
  "quantity": 2,
  "page": 1,
  "per_page": 20
}
HTTP 200

PATCH {{base_url}}/v1/cart/items/{{product_id}}?delta=1&page=1&per_page=20
HTTP 200

DELETE {{base_url}}/v1/cart/items/{{product_id}}?page=1&per_page=20
HTTP 200

POST {{base_url}}/v1/cart/items
Content-Type: application/json
{
  "product_id": {{product_id}},
  "quantity": 1,
  "page": 1,
  "per_page": 20
}
HTTP 200

GET {{base_url}}/v1/cart?page=1&per_page=20
HTTP 200
[Captures]
cart_version: jsonpath "$.version"

POST {{base_url}}/v1/orders
Content-Type: application/json
{
  "idempotency_key": "{{idempotency_key}}",
  "cart_version": {{cart_version}}
}
HTTP 200
[Captures]
order_id: jsonpath "$.order_id"

GET {{base_url}}/v1/orders?per_page=20&page=1
HTTP 200

GET {{base_url}}/v1/orders/{{order_id}}
HTTP 200

POST {{base_url}}/v1/orders/{{order_id}}/pay
Content-Type: application/json
{
  "payment_intent_id": "pi_hurl_test"
}
HTTP 200
