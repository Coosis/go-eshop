# Catalog read endpoints

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
