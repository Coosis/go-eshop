# Cart

All cart endpoints require an authenticated user (currently via dev middleware).

1. GET `/v1/cart?page=1&per_page=20` (get current cart; create if missing)

2. PUT `/v1/cart/items/{product_id}` (update qty; fails if item missing)
json:
```json
{
  "quantity": 2,
  "page": 1,
  "per_page": 20
}
```

3. POST `/v1/cart/items` (add item)
json:
```json
{
  "product_id": 123,
  "quantity": 2,
  "page": 1,
  "per_page": 20
}
```

4. PATCH `/v1/cart/items/{product_id}?delta=1&page=1&per_page=20` (change qty)

5. DELETE `/v1/cart/items/{product_id}?page=1&per_page=20`

6. DELETE `/v1/cart` (clear)
