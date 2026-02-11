# Generic Error Response
```json
{ "error": "Error message here" }
```

# Auth & Users
1. POST /v1/auth/register (create users, user_credentials, optional contact_methods)
2. POST /v1/auth/login
3. POST /v1/auth/logout (optional; usually token revocation / client-side)
4. GET  /v1/me
5. PATCH /v1/me (optional: notes, profile fields if you add later)

Contact methods
1. POST /v1/me/contacts (add email/phone)
2. POST /v1/me/contacts/{contact_id}/verify (mock verify flow)
3. DELETE /v1/me/contacts/{contact_id}

OAuth accounts (optional, can be stubbed)
1. POST /v1/auth/oauth/{provider}/callback
2. GET  /v1/me/oauth-accounts
3. DELETE /v1/me/oauth-accounts/{id}


# Catalog
Products
1. GET  /v1/catalog/products?per_page=20&page=1&q=shirt&min_price=10&max_price=20&category=somecategory
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
            "category_ids": [1, 2],
        }
    ],
    "page": 1,
    "per_page": 20,
    "total": 5
}
```
2. GET  /v1/catalog/products/{id}
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
3. GET  /v1/catalog/products/slug/{slug}
returns:
```json
{
    "id": 123,
    "name": "Cool Shirt",
    "slug": "cool-shirt",
    "description": "A very cool shirt.",
    "price_cents": 1999,
    "price_version": 1,
    "category_ids": [1, 2],
}
```

Categories
1. GET  /v1/catalog/categories
returns:
```json
{
    "items": [
        {
            "id": 1,
            "name": "Shirts",
            "slug": "shirts",
            "parent_id": 1234,
        }
    ],
    "page": 1,
    "per_page": 20,
    "total": 5,
}
```
2. GET  /v1/catalog/categories/{id}
returns:
```json
{
    "id": 1,
    "name": "Shirts",
    "slug": "shirts",
    "parent_id": 1234,
}
```
3. GET  /v1/catalog/categories/slug/{slug}
```json
{
    "id": 1,
    "name": "Shirts",
    "slug": "shirts",
    "parent_id": 1234,
}
```
4. GET  /v1/catalog/categories/{id}/products
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
            "category_ids": [1, 2],
        }
    ],
    "page": 1,
    "per_page": 20,
    "total": 5,
}
```

(Admin-only)
TODO! add delete
1. POST /v1/admin/products
2. PATCH /v1/admin/products/{id}
3. POST /v1/admin/categories
4. PATCH /v1/admin/categories/{id}


# Stock (admin-ish, good for showcasing transactions)
1. GET  /v1/products/{id}/stock (reads stock_levels)
returns:
```json
{
    "product_id": 123,
    "stock_level": 100
}
```
2. POST /v1/admin/stock/adjustments (writes stock_adjustments, updates stock_levels)
json: 
```json
{
    "product_id": 123,
    "delta": -5,
    "reason": "damaged in shipping",
    "created_by": "some admin",
}
```
returns:
```json
{
    "product_id": 123,
    "stock_level": 95
}
```
3. GET  /v1/admin/stock/adjustments (filter by product/date)
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
    "total": 5,
}
```
4. GET  /v1/admin/stock/adjustments/{id}
returns:
```json
{
    "product_id": 123,
    "stock_level": 100
}
```


# Cart
1. GET  /v1/cart (get current cart; create if missing)
2. PUT  /v1/cart/items/{product_id} (set qty; upsert cart_items)
3. POST /v1/cart/items (add item; body has product_id, qty)
4. PATCH /v1/cart/items/{product_id} (change qty)
5. DELETE /v1/cart/items/{product_id}
6. DELETE /v1/cart (clear)
7. POST /v1/cart/refresh-cart (re-snapshot price_cents_snapshot)


# Orders
1. POST /v1/orders
[auth required]
json:
```json
{
    "idempotency_key": "somekey",
    "cart_version": 234,
    "notes": "notesifneeded",
    "payment_intent_id": "somepaymentintent"
}
```
returns:
```json
{
    "order_id": 1234,
    "order_number": "ORD-20240901-0001",
    "subtotal_cents": 1234,
    "discount_cents": 1234,
    "total_cents": 1234,
    "status": "some status",
    "payment_intent_id": "some id",
    "notes": "some notes",
    "created_at": 1234,
    "version": 1234
}
```

2. GET  /v1/orders?before=1234&after=1234&status=paid&per_page=20&page=1
[auth required]
returns:
```json
[
    {
	    "order_id": 1234,
	    "order_number": "ORD-20240901-0001",
	    "subtotal_cents": 1234,
	    "discount_cents": 1234,
	    "total_cents": 1234,
	    "status": "some status",
	    "payment_intent_id": "some id",
	    "notes": "some notes",
	    "created_at": 1234,
	    "version": 1234
    }
]
```
3. GET  /v1/orders/{order_id}
[auth required]
returns:
```json
{
    "order_id": 1234,
    "order_number": "ORD-20240901-0001",
    "subtotal_cents": 1234,
    "discount_cents": 1234,
    "total_cents": 1234,
    "status": "some status",
    "payment_intent_id": "some id",
    "notes": "some notes",
    "created_at": 1234,
    "version": 1234
}
```

4. POST /v1/orders/{order_id}/cancel
[auth required]
json:
```json
{ "order_id": 1234 }
```
returns: 
```json
{
    "order_id": 1234,
    "order_number": "ORD-20240901-0001",
    "subtotal_cents": 1234,
    "discount_cents": 1234,
    "total_cents": 1234,
    "status": "some status",
    "payment_intent_id": "some id",
    "notes": "some notes",
    "created_at": 1234,
    "version": 1234
}
```

5. POST /v1/orders/{order_id}/pay (mock payment intent)
[auth required]
```json
{ "order_id": 1234 }
```
returns:
```json
{
    "order_id": 1234,
    "order_number": "ORD-20240901-0001",
    "subtotal_cents": 1234,
    "discount_cents": 1234,
    "total_cents": 1234,
    "status": "some status",
    "payment_intent_id": "some id",
    "notes": "some notes",
    "created_at": 1234,
    "version": 1234
}
```

6. POST /v1/orders/{order_id}/refund (optional)
[auth required]
```json
{ "order_id": 1234 }
```
returns:
```json
{
    "order_id": 1234,
    "order_number": "ORD-20240901-0001",
    "subtotal_cents": 1234,
    "discount_cents": 1234,
    "total_cents": 1234,
    "status": "some status",
    "payment_intent_id": "some id",
    "notes": "some notes",
    "created_at": 1234,
    "version": 1234
}
```

7. POST /v1/orders/{order_id}/payment-webhook (internal/mock) to transition waiting_payment -> paid / payment_failed


# Seckill
Events
1. GET  /v1/seckill/events (active/upcoming)
returns:
```json
{
    "items": [
        {
            "id": 123,
            "product_id": 456,
            "start_time": 1693555200,
            "end_time": 1693562400,
            "seckill_price_cents": 999,
            "seckill_stock": 100,
        }
    ],
    "page": 1,
    "per_page": 20,
    "total": 5,
}
```
2. GET  /v1/seckill/events/{event_id}
returns:
```json
{
    "id": 123,
    "product_id": 456,
    "start_time": 1693555200,
    "end_time": 1693562400,
    "seckill_price_cents": 999,
    "seckill_stock": 100,
}
```
3. POST /v1/admin/seckill/events
json:
```json
{
    "product_id": 456,
    "start_time": 1693555200,
    "end_time": 1693562400,
    "seckill_price_cents": 999,
    "seckill_stock": 100,
}
```
returns:
```json
{
    "id": 123,
    "product_id": 456,
    "start_time": 1693555200,
    "end_time": 1693562400,
    "seckill_price_cents": 999,
    "seckill_stock": 100,
}
```
4. PATCH /v1/admin/seckill/events/{event_id}
returns:
```json
{
    "id": 123,
    "product_id": 456,
    "start_time": 1693555200,
    "end_time": 1693562400,
    "seckill_price_cents": 999,
    "seckill_stock": 100,
}
```

Purchase flow
1. POST /v1/seckill/events/{event_id}/attempt (Idempotency-Key required)
[auth required]
json:
```json
{
    "event_id": 123,
    "quantity": 1,
    "idempotency_key": "somekey"
}
```
returns:
```json
{
    "state": "queued",
    "order_id": null | 1234
}
```
internally: redis lock + lua stock decrement + enqueue outbox/MQ
2. GET  /v1/seckill/attempts/{request_id} (polling status if you do async)
[auth required]
returns:
```json
{
    "state": "queued",
    "order_id": null | 1234
}
```

Outbox / ops (internal/admin)
1. POST /v1/admin/seckill/outbox/preheat (mark preheated)
json:
```json
{
    "event_id": 123
}
```
HTTP 200 on success
2. GET  /v1/admin/seckill/outbox (inspect scheduled/preheated rows)
returns:
```json

```


# Infra / Observability
1. GET /healthz
2. GET /readyz
3. GET /metrics (Prometheus if you do it)
