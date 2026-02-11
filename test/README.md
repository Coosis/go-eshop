# Hurl Tests

These tests assume the API server is running locally and seeded with an empty-ish DB.

## Run (single flow)
Captures don’t carry across separate Hurl files, so run the combined flow:

```sh
suffix=$(date +%s%N)

hurl --test \
  --variable base_url=http://localhost:8144 \
  --variable idempotency_key=order-$suffix \
  --variable category_name=Hurl-Category-$suffix \
  --variable category_slug=hurl-category-$suffix \
  --variable product_name=Hurl-Product-$suffix \
  --variable product_slug=hurl-product-$suffix \
  test/hurl/00_flow.hurl
```

## Negative tests
Requires captures from the flow (product_id), so run after `00_flow.hurl` or provide variables manually.

```sh
suffix=$(date +%s%N)

hurl --test \
  --variable base_url=http://localhost:8144 \
  --variable product_id=123 \
  test/hurl/90_cart_validation.hurl
```

## Seckill (optional)
Requires redis + worker/scheduler if you use them:

```sh
suffix=$(date +%s%N)

hurl --test \
  --variable base_url=http://localhost:8144 \
  --variable seckill_idempotency_key=seckill-$suffix \
  test/hurl/30_seckill.hurl
```
