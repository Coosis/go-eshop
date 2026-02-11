# Stock (admin-ish, good for showcasing transactions)

1. GET `/v1/products/{id}/stock` (reads stock_levels)
returns:
```json
{
  "product_id": 123,
  "stock_level": 100
}
```

2. POST `/v1/admin/stock/adjustments` (writes stock_adjustments, updates stock_levels)
json:
```json
{
  "product_id": 123,
  "delta": -5,
  "reason": "damaged in shipping",
  "created_by": "some admin"
}
```
returns:
```json
{
  "product_id": 123,
  "stock_level": 95
}
```

3. GET `/v1/admin/stock/adjustments?product_id=123` (product_id required; optional filters below)
Optional filters: `created_after`, `created_before`, `created_by`, `delta_min`, `delta_max`, `page`, `per_page`
returns:
```json
{
  "items": [
    {
      "product_id": 123,
      "stock_level": 100
    }
  ],
  "page": 1,
  "per_page": 20,
  "total": 5
}
```

4. GET `/v1/admin/stock/adjustments/{id}`
returns:
```json
{
  "product_id": 123,
  "stock_level": 100
}
```
