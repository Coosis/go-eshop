# Seckill flow (requires redis + seckill worker/scheduler if you use them)

POST {{base_url}}/v1/admin/seckill/events
Content-Type: application/json
{
  "product_id": {{product_id}},
  "start_time": 1710000000000,
  "end_time": 1710003600000,
  "seckill_price_cents": 999,
  "seckill_stock": 10
}
HTTP 200
[Captures]
seckill_event_id: jsonpath "$.id"

GET {{base_url}}/v1/seckill/events?page=1&per_page=20
HTTP 200

POST {{base_url}}/v1/seckill/events/{{seckill_event_id}}/attempt
Content-Type: application/json
{
  "event_id": {{seckill_event_id}},
  "quantity": 1,
  "idempotency_key": "{{seckill_idempotency_key}}"
}
HTTP 200

GET {{base_url}}/v1/seckill/attempts/{{seckill_idempotency_key}}/status
HTTP 200
