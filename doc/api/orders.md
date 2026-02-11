# Orders

All order endpoints require an authenticated user (currently via dev middleware), except the webhook.

1. POST `/v1/orders`
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

2. GET `/v1/orders?before=1234&after=1234&status=paid&per_page=20&page=1`
returns:
```json
{
  "items": [
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
  ],
  "page": 1,
  "per_page": 20,
  "total": 1
}
```

3. GET `/v1/orders/{order_id}`
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

4. POST `/v1/orders/{order_id}/cancel`
returns:
```json
"ok"
```

5. POST `/v1/orders/{order_id}/pay`
json:
```json
{ "payment_intent_id": "some id" }
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

6. POST `/v1/orders/{order_id}/refund`
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

7. POST `/{order_id}/webhook` (internal/mock) to transition waiting_payment -> paid / payment_failed
