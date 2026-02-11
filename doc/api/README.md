# API Docs (Core)

These docs describe the core API surface as implemented.
Auth/user flows are intentionally omitted for now.

# Conventions
1. Success responses use HTTP 200 with JSON bodies.
2. Error responses use HTTP 4xx/5xx with a JSON body:
```json
{ "error": "Error message here" }
```
3. Pagination uses:
```json
{
  "items": [],
  "page": 1,
  "per_page": 20,
  "total": 0
}
```
`total` is the number of items in the current page (not the global total).
4. Timestamps are milliseconds since epoch unless noted.

# Sections
1. Catalog: `doc/api/catalog.md`
2. Cart: `doc/api/cart.md`
3. Orders: `doc/api/orders.md`
4. Stock: `doc/api/stock.md`
5. Seckill: `doc/api/seckill.md`
