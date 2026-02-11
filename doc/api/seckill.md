# Seckill

## Events
1. GET `/v1/seckill/events` (active/upcoming)
returns:
```json
{
  "items": [
    {
      "id": 123,
      "product_id": 456,
      "start_time": 1693555200000,
      "end_time": 1693562400000,
      "seckill_price_cents": 999,
      "seckill_stock": 100
    }
  ],
  "page": 1,
  "per_page": 20,
  "total": 1
}
```

2. GET `/v1/seckill/events/{event_id}`
returns:
```json
{
  "id": 123,
  "product_id": 456,
  "start_time": 1693555200000,
  "end_time": 1693562400000,
  "seckill_price_cents": 999,
  "seckill_stock": 100
}
```

3. POST `/v1/admin/seckill/events`
json:
```json
{
  "product_id": 456,
  "start_time": 1693555200000,
  "end_time": 1693562400000,
  "seckill_price_cents": 999,
  "seckill_stock": 100
}
```
returns:
```json
{
  "id": 123,
  "product_id": 456,
  "start_time": 1693555200000,
  "end_time": 1693562400000,
  "seckill_price_cents": 999,
  "seckill_stock": 100
}
```

4. PATCH `/v1/admin/seckill/events/{event_id}`
returns:
```json
{
  "id": 123,
  "product_id": 456,
  "start_time": 1693555200000,
  "end_time": 1693562400000,
  "seckill_price_cents": 999,
  "seckill_stock": 100
}
```

## Purchase flow
1. POST `/v1/seckill/events/{event_id}/attempt`
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
  "state": "OK:95"
}
```

2. GET `/v1/seckill/attempts/{idempotency_key}/status` (polling status if async)
returns:
```json
{
  "state": "queued"
}
```
