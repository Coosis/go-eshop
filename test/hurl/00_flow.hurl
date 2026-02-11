POST {{base_url}}/v1/admin/catalog/categories
Content-Type: application/json
{
  "name": "{{category_name}}",
  "slug": "{{category_slug}}"
}
HTTP 200
[Captures]
category_id: jsonpath "$.id"

PATCH {{base_url}}/v1/admin/catalog/categories/{{category_id}}
Content-Type: application/json
{
  "id": {{category_id}},
  "name": "{{category_name}} Updated",
  "slug": "{{category_slug}}",
  "parent_id": null
}
HTTP 200

POST {{base_url}}/v1/admin/catalog/products
Content-Type: application/json
{
  "name": "{{product_name}}",
  "slug": "{{product_slug}}",
  "description": "Product created by hurl",
  "price_cents": 1999,
  "category_ids": [{{category_id}}]
}
HTTP 200
[Captures]
product_id: jsonpath "$.id"
product_slug: jsonpath "$.slug"

PATCH {{base_url}}/v1/admin/catalog/products/{{product_id}}
Content-Type: application/json
{
  "id": {{product_id}},
  "name": "{{product_name}} Updated",
  "slug": "{{product_slug}}",
  "description": "Product updated by hurl",
  "price_cents": 2099,
  "category_ids": [{{category_id}}]
}
HTTP 200

POST {{base_url}}/v1/admin/stock/adjustments
Content-Type: application/json
{
  "product_id": {{product_id}},
  "delta": 100,
  "reason": "initial stock",
  "created_by": "hurl"
}
HTTP 200

GET {{base_url}}/v1/admin/stock/adjustments?product_id={{product_id}}&delta_min=-100000&delta_max=100000&page=1&per_page=10
HTTP 200

GET {{base_url}}/v1/catalog/products?page=1&per_page=20&min_price=0&max_price=999999
HTTP 200

GET {{base_url}}/v1/catalog/products/{{product_id}}
HTTP 200

GET {{base_url}}/v1/catalog/products/slug/{{product_slug}}
HTTP 200

GET {{base_url}}/v1/catalog/categories?page=1&per_page=20
HTTP 200

GET {{base_url}}/v1/catalog/categories/{{category_id}}
HTTP 200

GET {{base_url}}/v1/catalog/categories/slug/{{category_slug}}
HTTP 200

GET {{base_url}}/v1/catalog/categories/{{category_id}}/products?page=1&per_page=20
HTTP 200

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
