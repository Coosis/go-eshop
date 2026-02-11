# Catalog

## Products
1. GET `/v1/catalog/products?per_page=20&page=1&min_price=10&max_price=20`
returns:
```json
{
  "items": [
    {
      "id": 123,
      "name": "Cool Shirt",
      "slug": "cool-shirt",
      "description": "A very cool shirt.",
      "price_cents": 1999,
      "price_version": 1,
      "category_ids": [1, 2]
    }
  ],
  "page": 1,
  "per_page": 20,
  "total": 1
}
```

2. GET `/v1/catalog/products/{id}`
returns:
```json
{
  "id": 123,
  "name": "Cool Shirt",
  "slug": "cool-shirt",
  "description": "A very cool shirt.",
  "price_cents": 1999,
  "price_version": 1,
  "category_ids": [1, 2]
}
```

3. GET `/v1/catalog/products/slug/{slug}`
returns:
```json
{
  "id": 123,
  "name": "Cool Shirt",
  "slug": "cool-shirt",
  "description": "A very cool shirt.",
  "price_cents": 1999,
  "price_version": 1,
  "category_ids": [1, 2]
}
```

## Categories
1. GET `/v1/catalog/categories`
returns:
```json
{
  "items": [
    {
      "id": 1,
      "name": "Shirts",
      "slug": "shirts",
      "parent_id": 1234
    }
  ],
  "page": 1,
  "per_page": 20,
  "total": 1
}
```

2. GET `/v1/catalog/categories/{id}`
returns:
```json
{
  "id": 1,
  "name": "Shirts",
  "slug": "shirts",
  "parent_id": 1234
}
```

3. GET `/v1/catalog/categories/slug/{slug}`
returns:
```json
{
  "id": 1,
  "name": "Shirts",
  "slug": "shirts",
  "parent_id": 1234
}
```

4. GET `/v1/catalog/categories/{id}/products`
returns:
```json
{
  "items": [
    {
      "id": 123,
      "name": "Cool Shirt",
      "slug": "cool-shirt",
      "description": "A very cool shirt.",
      "price_cents": 1999,
      "price_version": 1,
      "category_ids": [1, 2]
    }
  ],
  "page": 1,
  "per_page": 20,
  "total": 1
}
```

## Admin (Catalog)
1. POST `/v1/admin/catalog/products`
2. PATCH `/v1/admin/catalog/products/{id}`
3. POST `/v1/admin/catalog/categories`
4. PATCH `/v1/admin/catalog/categories/{id}`
