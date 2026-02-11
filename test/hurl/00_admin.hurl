# Admin setup: create category, product, stock

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
  "name": "Hurl Category Updated",
  "slug": "hurl-category",
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

GET {{base_url}}/v1/admin/stock/adjustments?product_id={{product_id}}&page=1&per_page=10
HTTP 200
