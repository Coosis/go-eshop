-- name: GetCurrentCart :many
WITH old_cart AS (
  SELECT id FROM carts cts WHERE cts.user_id = $1
), refreshed_items AS (
  UPDATE cart_items
  SET price_cents_snapshot = p.price_cents,
      updated_at = NOW()
  FROM products p
  WHERE cart_id = (SELECT id FROM old_cart)
    AND product_id = p.id
  RETURNING *
), new_cart AS (
  UPDATE carts SET 
    updated_at = NOW(),
    version = version + 1
  WHERE id = (SELECT id FROM old_cart)
  RETURNING id, version
), cart_joined AS (
  SELECT nc.version, ci.product_id, ci.qty,
    ci.price_cents_snapshot, ci.updated_at
  FROM new_cart nc
  LEFT JOIN cart_items ci ON ci.cart_id = nc.id
), total_count AS (
  SELECT COUNT(*) AS total_count
  FROM cart_items ci
  WHERE ci.cart_id = (SELECT id FROM old_cart)
) SELECT cj.version, cj.product_id, cj.qty,
  cj.price_cents_snapshot, cj.updated_at, tc.total_count
FROM cart_joined cj, total_count tc
LIMIT @page_size
OFFSET ((@page_number::int)-1)*(@page_size);

-- name: UpdateCartItem :many
WITH old_cart AS (
  SELECT id
  FROM carts cts
  WHERE cts.user_id = $1
), updated_qty AS (
  UPDATE cart_items cit
  SET qty = $3,
      updated_at = NOW()
  WHERE cit.cart_id = (SELECT id FROM old_cart)
    AND cit.product_id = $2
  RETURNING cit.cart_id, cit.product_id, cit.qty, cit.price_cents_snapshot, cit.updated_at
), refreshed_items AS (
  UPDATE cart_items ci
  SET price_cents_snapshot = p.price_cents,
      updated_at = NOW()
  FROM products p
  WHERE ci.cart_id = (SELECT id FROM old_cart)
    AND ci.product_id = p.id
    AND ci.product_id <> $2
  RETURNING ci.cart_id, ci.product_id, ci.qty, ci.price_cents_snapshot, ci.updated_at
), all_items AS (
  SELECT * FROM updated_qty
  UNION ALL
  SELECT * FROM refreshed_items
), new_cart AS (
  UPDATE carts
  SET updated_at = NOW(),
      version = version + 1
  WHERE id = (SELECT id FROM old_cart)
  RETURNING id, version
), cart_joined AS (
  SELECT nc.version, ci.product_id, ci.qty,
         ci.price_cents_snapshot, ci.updated_at
  FROM new_cart nc
  LEFT JOIN all_items ci ON ci.cart_id = nc.id
), total_count AS (
  SELECT COUNT(*) AS total_count
  FROM all_items ci
) SELECT cj.version, cj.product_id, cj.qty,
    cj.price_cents_snapshot, cj.updated_at, tc.total_count
FROM cart_joined cj, total_count tc
LIMIT @page_size
OFFSET ((@page_number::int)-1)*(@page_size);

-- name: AddCartItem :many
WITH old_cart AS (
  SELECT id FROM carts cts WHERE cts.user_id = $1
), refreshed_items AS (
  UPDATE cart_items ci
  SET price_cents_snapshot = p.price_cents,
      updated_at = NOW()
  FROM products p
  WHERE cart_id = (SELECT id FROM old_cart)
    AND product_id = p.id
    AND p.id <> $2
  RETURNING ci.cart_id, ci.product_id, ci.qty, ci.price_cents_snapshot, ci.updated_at
), inserted_item AS (
  INSERT INTO cart_items (cart_id, product_id, qty, price_cents_snapshot)
  SELECT (SELECT id FROM old_cart), $2, $3, p.price_cents
  FROM products p
  WHERE p.id = $2
  ON CONFLICT (cart_id, product_id) DO UPDATE 
    SET qty = cart_items.qty + EXCLUDED.qty,
        updated_at = NOW()
  RETURNING cart_id, product_id, qty, price_cents_snapshot, updated_at
), new_cart AS (
  UPDATE carts SET 
    updated_at = NOW(),
    version = version + 1
  WHERE id = (SELECT id FROM old_cart)
  RETURNING id, version
), all_items AS (
  SELECT * FROM refreshed_items
  UNION ALL
  SELECT * FROM inserted_item
), cart_joined AS (
  SELECT nc.version, ci.product_id, ci.qty,
    ci.price_cents_snapshot, ci.updated_at
  FROM new_cart nc
  LEFT JOIN all_items ci ON ci.cart_id = nc.id
), total_count AS (
  SELECT COUNT(*) AS total_count
  FROM all_items ci
) SELECT cj.version, cj.product_id, cj.qty,
  cj.price_cents_snapshot, cj.updated_at, tc.total_count
FROM cart_joined cj, total_count tc
LIMIT @page_size
OFFSET ((@page_number::int)-1)*(@page_size);

-- name: ChangeCartItemQty :many
WITH old_cart AS (
  SELECT id FROM carts cts WHERE cts.user_id = $1
), refreshed_items AS (
  UPDATE cart_items ci
  SET price_cents_snapshot = p.price_cents,
      updated_at = NOW(),
      qty = CASE 
      WHEN ci.product_id = $2
        THEN ci.qty + @delta
        ELSE ci.qty END
  FROM products p
  WHERE cart_id = (SELECT id FROM old_cart)
    AND product_id = p.id
  RETURNING ci.cart_id, ci.product_id, ci.qty, ci.price_cents_snapshot, ci.updated_at
), new_cart AS (
  UPDATE carts SET 
    updated_at = NOW(),
    version = version + 1
  WHERE id = (SELECT id FROM old_cart)
  RETURNING id, version
), cart_joined AS (
  SELECT nc.version, ci.product_id, ci.qty,
    ci.price_cents_snapshot, ci.updated_at
  FROM new_cart nc
  LEFT JOIN refreshed_items ci ON ci.cart_id = nc.id
), total_count AS (
  SELECT COUNT(*) AS total_count
  FROM refreshed_items ci
) SELECT cj.version, cj.product_id, cj.qty,
  cj.price_cents_snapshot, cj.updated_at, tc.total_count
FROM cart_joined cj, total_count tc
LIMIT @page_size
OFFSET ((@page_number::int)-1)*(@page_size);

-- name: RemoveCartItem :many
WITH old_cart AS (
  SELECT id FROM carts cts WHERE cts.user_id = $1
), refreshed_items AS (
  UPDATE cart_items ci
  SET price_cents_snapshot = p.price_cents,
      updated_at = NOW()
  FROM products p
  WHERE cart_id = (SELECT id FROM old_cart)
    AND product_id = p.id
    AND p.id <> $2
  RETURNING ci.cart_id, ci.product_id, ci.qty, ci.price_cents_snapshot, ci.updated_at
), deleted_item AS (
  DELETE FROM cart_items ci
  WHERE ci.cart_id = (SELECT id FROM old_cart)
    AND ci.product_id = $2
  RETURNING *
), new_cart AS (
  UPDATE carts SET 
    updated_at = NOW(),
    version = version + 1
  WHERE id = (SELECT id FROM old_cart)
  RETURNING id, version
), cart_joined AS (
  SELECT nc.version, ci.product_id, ci.qty,
    ci.price_cents_snapshot, ci.updated_at
  FROM new_cart nc
  LEFT JOIN refreshed_items ci ON ci.cart_id = nc.id
), total_count AS (
  SELECT COUNT(*) AS total_count
  FROM refreshed_items ci
) SELECT cj.version, cj.product_id, cj.qty,
  cj.price_cents_snapshot, cj.updated_at, tc.total_count
FROM cart_joined cj, total_count tc
LIMIT @page_size
OFFSET ((@page_number::int)-1)*(@page_size);

-- name: ClearCart :many
WITH old_cart AS (
  SELECT id FROM carts cts WHERE cts.user_id = $1
), deleted_items AS (
  DELETE FROM cart_items
  WHERE cart_id = (SELECT id FROM old_cart)
  RETURNING *
), new_cart AS (
  UPDATE carts SET
    updated_at = NOW(),
    version = version + 1
  WHERE id = (SELECT id FROM old_cart)
  RETURNING id, version
) SELECT nc.version
FROM new_cart nc;

-- name: RefreshCart :many
WITH old_cart AS (
  SELECT id FROM carts cts WHERE cts.user_id = $1
), refreshed_items AS (
  UPDATE cart_items ci
  SET price_cents_snapshot = p.price_cents,
      updated_at = NOW()
  FROM products p
  WHERE cart_id = (SELECT id FROM old_cart)
    AND product_id = p.id
  RETURNING ci.cart_id, ci.product_id, ci.qty, ci.price_cents_snapshot, ci.updated_at
), new_cart AS (
  UPDATE carts SET 
    updated_at = NOW(),
    version = version + 1
  WHERE id = (SELECT id FROM old_cart)
  RETURNING id, version
), cart_joined AS (
  SELECT nc.version, ci.product_id, ci.qty,
    ci.price_cents_snapshot, ci.updated_at
  FROM new_cart nc
  LEFT JOIN refreshed_items ci ON ci.cart_id = nc.id
), total_count AS (
  SELECT COUNT(*) AS total_count
  FROM refreshed_items ci
) SELECT cj.version, cj.product_id, cj.qty,
  cj.price_cents_snapshot, cj.updated_at, tc.total_count
FROM cart_joined cj, total_count tc
LIMIT @page_size
OFFSET ((@page_number::int)-1)*(@page_size);


-- name: GetStockLevel :one
SELECT product_id, on_hand FROM stock_levels WHERE product_id = $1;

-- name: AdjustStockLevel :one
WITH updated AS (
  UPDATE stock_levels stl
  SET on_hand = stl.on_hand + @delta,
      updated_at = now()
  WHERE stl.product_id = $1
    AND stl.on_hand + @delta >= 0
  RETURNING stl.product_id, stl.on_hand
),
adjustment AS (
  INSERT INTO stock_adjustments (product_id, delta, reason, created_by)
  SELECT $1, @delta, $2, $3
  FROM updated
  RETURNING id
)
SELECT u.product_id, u.on_hand
FROM updated u;

-- name: GetStockAdjustments :many
SELECT *
FROM stock_adjustments
WHERE product_id = $1
  AND (@created_after::timestamptz IS NULL OR created_at >= @created_after::timestamptz)
  AND (@created_before::timestamptz IS NULL OR created_at <= @created_before::timestamptz)
  AND (@created_by::TEXT IS NULL OR @created_by = '' OR created_by = @created_by)
  AND (@min_delta::int IS NULL OR delta >= @min_delta::int)
  AND (@max_delta::int IS NULL OR delta <= @max_delta::int)
ORDER BY created_at DESC
LIMIT @page_size
OFFSET ((@page_number::int)-1) * @page_size;

-- name: GetStockAdjustmentByID :one
SELECT * FROM stock_adjustments
WHERE id = $1 LIMIT 1;

-- name: GetSeckillEvents :many
SELECT * FROM seckill_events
WHERE end_time >= NOW()
ORDER BY start_time ASC
LIMIT @page_size
OFFSET ((@page_number::int)-1)*(@page_size);

-- name: GetSeckillEventByID :one
SELECT * FROM seckill_events
WHERE id = $1 LIMIT 1;

-- name: AddSeckillEvent :one
INSERT INTO seckill_events (
  product_id, start_time, end_time, 
  seckill_price_cents, seckill_stock, 
  preheated_at
)
VALUES (
  $1, $2, $3, 
  $4, $5, 
  NULL
)
RETURNING *;

-- name: UpdateSeckillEventByID :one
UPDATE seckill_events SET
  product_id = $2,
  start_time = $3,
  end_time = $4,
  seckill_price_cents = $5,
  seckill_stock = $6,
  updated_at = NOW(),
  preheated_at = NULL
WHERE id = $1
RETURNING *;

-- name: MarkPreheated :one
UPDATE seckill_events SET preheated_at = NOW() WHERE id = $1
RETURNING *;

-- name: GetDueNotPreheated :many
SELECT id, product_id, start_time, end_time, seckill_price_cents,
seckill_stock
FROM seckill_events
WHERE preheated_at IS NULL
  AND start_time <= NOW() + interval '5 minutes'
ORDER BY start_time ASC
LIMIT $1
FOR UPDATE SKIP LOCKED;

-- catalog ---------------------------------------------------------------

-- name: GetProducts :many
WITH products AS (
  SELECT prod.* FROM products prod
  WHERE (sqlc.narg(min_price_cents)::int IS NULL OR prod.price_cents >= @min_price_cents)
    AND (sqlc.narg(max_price_cents)::int IS NULL OR prod.price_cents <= @max_price_cents)
    AND (sqlc.narg(category_id)::int IS NULL OR EXISTS (
      SELECT 1 FROM product_categories pc
      WHERE pc.product_id = prod.id AND pc.category_id = @category_id::int
    ))
  ORDER BY prod.id
  LIMIT @page_size
  OFFSET ((@page_number::int)-1)*(@page_size)
) SELECT 
  p.*,
  COALESCE(
    array_agg(category_id ORDER BY pc.category_id)
      FILTER (WHERE pc.category_id IS NOT NULL),
      '{}'
  )::int[] AS category_ids
FROM products p
LEFT JOIN product_categories pc ON pc.product_id = p.id
GROUP BY p.id, p.name, p.slug, p.description, p.price_cents, p.price_version, p.created_at, p.updated_at;

-- name: GetProductByID :one
WITH product AS (
  SELECT prod.*
  FROM products prod
  WHERE id = $1 LIMIT 1
) SELECT 
  p.*,
  COALESCE(
    array_agg(category_id ORDER BY pc.category_id)
    FILTER (WHERE pc.category_id IS NOT NULL),
    '{}'
  )::int[] AS category_ids
FROM product p
LEFT JOIN product_categories pc ON pc.product_id = p.id
GROUP BY p.id, p.name, p.slug, p.description, p.price_cents, p.price_version, p.created_at, p.updated_at;

-- name: GetProductBySlug :one
WITH product AS (
  SELECT prod.*
  FROM products prod
  WHERE slug = $1 LIMIT 1
) SELECT 
  p.*,
  COALESCE(
    array_agg(category_id ORDER BY pc.category_id)
      FILTER (WHERE pc.category_id IS NOT NULL),
      '{}'
  )::int[] AS category_ids
FROM product p
LEFT JOIN product_categories pc ON pc.product_id = p.id
GROUP BY p.id, p.name, p.slug, p.description, p.price_cents, p.price_version, p.created_at, p.updated_at;

-- name: GetCategories :many
SELECT c.id, c.name, c.slug, c.parent_id FROM categories c
ORDER BY c.id
LIMIT @page_size
OFFSET ((@page_number::int)-1)*(@page_size);

-- name: GetCategoryByID :one
SELECT c.id, c.name, c.slug, c.parent_id FROM categories c
WHERE c.id = $1 LIMIT 1;

-- name: GetCategoryBySlug :one
SELECT c.id, c.name, c.slug, c.parent_id FROM categories c
WHERE c.slug = $1 LIMIT 1;

-- name: GetProductByCategoryID :many
WITH products AS (
  SELECT p.* FROM products p
  WHERE EXISTS (
    SELECT 1 FROM product_categories pc
    WHERE pc.product_id = p.id AND pc.category_id = $1
  ) ORDER BY p.id
  LIMIT @page_size
  OFFSET ((@page_number::int)-1)*(@page_size)
) SELECT 
  p.*,
  COALESCE(
    array_agg(category_id ORDER BY pc.category_id)
      FILTER (WHERE pc.category_id IS NOT NULL),
      '{}'
  )::int[] AS category_ids
FROM products p
LEFT JOIN product_categories pc ON pc.product_id = p.id
GROUP BY p.id, p.name, p.slug, p.description, p.price_cents, p.price_version, p.created_at, p.updated_at;

-- name: CreateProduct :one
WITH products AS (
  INSERT INTO products (name, slug, description, price_cents, created_at, updated_at)
  VALUES ($1, $2, $3, $4, NOW(), NOW())
  RETURNING *
), inserted_stock AS (
  INSERT INTO stock_levels (product_id, on_hand, reserved, updated_at)
  SELECT id, 0, 0, NOW() FROM products
) SELECT 
  p.*,
  COALESCE(
    array_agg(category_id ORDER BY pc.category_id)
      FILTER (WHERE pc.category_id IS NOT NULL),
      '{}'
  )::int[] AS category_ids
FROM products p
LEFT JOIN product_categories pc ON pc.product_id = p.id
GROUP BY p.id, p.name, p.slug, p.description, p.price_cents, p.price_version, p.created_at, p.updated_at;

-- name: UpdateProductByID :one
WITH products AS (
  UPDATE products prod
  SET name = $2,
      slug = $3,
      description = $4,
      price_cents = $5,
      price_version = price_version + 1,
      updated_at = NOW()
  WHERE prod.id = $1
  RETURNING *
) SELECT 
  p.*,
  COALESCE(
    array_agg(category_id ORDER BY pc.category_id)
      FILTER (WHERE pc.category_id IS NOT NULL),
      '{}'
  )::int[] AS category_ids
FROM products p
LEFT JOIN product_categories pc ON pc.product_id = p.id
GROUP BY p.id, p.name, p.slug, p.description, p.price_cents, p.price_version, p.created_at, p.updated_at;

-- name: CreateCategory :one
INSERT INTO categories (name, slug, parent_id, created_at)
VALUES ($1, $2, $3, NOW())
RETURNING *;

-- name: UpdateCategoryByID :one
UPDATE categories
SET name = $2,
    slug = $3,
    parent_id = $4
WHERE id = $1
RETURNING *;

-- order ---------------------------------------------------------------
-- name: CreateOrder :one
WITH new_cart AS (
  UPDATE carts SET 
    updated_at = NOW(),
    version = version + 1
  WHERE carts.user_id = $1
    AND carts.version = $7
  RETURNING id
), cart_items_with_price AS (
  SELECT *
  FROM new_cart c
  JOIN cart_items ci ON ci.cart_id = c.id
  JOIN products p ON p.id = ci.product_id
  WHERE c.id = (SELECT id FROM new_cart)
), item_summary AS (
  SELECT SUM(cip.qty * cip.price_cents) AS subtotal_cents
  FROM cart_items_with_price cip
), new_order AS (
  INSERT INTO orders (
    order_number, user_id, 
    subtotal_cents, discount_cents, total_cents,
    status, idempotency_key, payment_intent_id, 
    notes, version
  )
  SELECT 
    $2, $1,
    its.subtotal_cents, $3, (its.subtotal_cents - $3),
    'waiting_payment', $6, $4, 
    $5, 1
  FROM item_summary its
  JOIN new_cart nc ON TRUE
  WHERE its.subtotal_cents IS NOT NULL
  RETURNING id, order_number, subtotal_cents, discount_cents, total_cents,
    status, payment_intent_id, notes, created_at, version
), order_items AS (
  INSERT INTO order_items (order_id, product_id, product_name, qty, unit_price_cents,
    price_version, metadata)
  SELECT no.id, cip.product_id, cip.name, cip.qty, cip.price_cents,
    cip.price_version, '{}'::jsonb
  FROM cart_items_with_price cip
  JOIN new_order no ON TRUE
  RETURNING *
), cleared_cart AS (
  DELETE FROM cart_items
  WHERE cart_id = (SELECT id FROM new_cart)
) SELECT * FROM new_order;

-- name: CreateSeckillOrder :one
INSERT INTO orders (
  order_number, user_id, 
  subtotal_cents, discount_cents, total_cents,
  status, idempotency_key, payment_intent_id, 
  notes, version
) VALUES (
  $1, $2,
  $3, 0, $3,
  'waiting_payment', $4, NULL,
  $5, 1
) RETURNING id, order_number, subtotal_cents, discount_cents, total_cents,
  status, payment_intent_id, notes, created_at, version;

---- name: AddOrderItem :one
-- INSERT INTO order_items (order_id, product_id, product_name, qty, unit_price_cents,
--   price_version, metadata)
-- VALUES ($1, $2, $3, $4, $5, $6, $7)
-- RETURNING *;
  
-- name: GetOrders :many
SELECT id, order_number, subtotal_cents, discount_cents, total_cents,
  status, payment_intent_id, notes, created_at, version
FROM orders
WHERE user_id = $1
  AND (sqlc.narg(before)::TIMESTAMPTZ IS NULL OR created_at <= sqlc.narg(before))
  AND (sqlc.narg(after)::TIMESTAMPTZ IS NULL OR created_at >= sqlc.narg(after))
  AND (sqlc.narg(status)::order_status IS NULL OR status = sqlc.narg(status)::order_status)
ORDER BY updated_at DESC
LIMIT @page_size
OFFSET ((@page_number::int)-1)*(@page_size);

-- name: GetOrderByID :one
SELECT id, order_number, subtotal_cents, discount_cents, total_cents,
  status, payment_intent_id, notes, created_at, version
FROM orders
WHERE id = $1 
  AND user_id = $2
LIMIT 1;

-- name: CancelOrder :one
UPDATE orders
SET status = 'canceled',
    updated_at = NOW(),
    version = version + 1
WHERE id = $1
  AND user_id = $2
  AND status = 'waiting_payment'
RETURNING id, order_number, subtotal_cents, discount_cents, total_cents,
  status, payment_intent_id, notes, created_at, version;

-- name: PayOrder :one
UPDATE orders SET
  status = 'paid',
  payment_intent_id = $2,
  updated_at = NOW(),
  version = version + 1
WHERE id = $1
  AND user_id = $3
  AND status = 'waiting_payment'
RETURNING id, order_number, subtotal_cents, discount_cents, total_cents,
  status, payment_intent_id, notes, created_at, version;

-- name: RefundOrder :one
UPDATE orders SET
  status = 'refunded',
  updated_at = NOW(),
  version = version + 1
WHERE id = $1
  AND user_id = $2
  AND status = 'paid'
RETURNING id, order_number, subtotal_cents, discount_cents, total_cents,
  status, payment_intent_id, notes, created_at, version;

-- stock ---------------------------------------------------------------
-- name: SoftHoldStock :exec
UPDATE stock_levels 
SET on_hand = on_hand - @delta,
    reserved = reserved + @delta,
    updated_at = NOW()
WHERE product_id = $1
  AND on_hand >= @delta;

-- name: ReleaseStockHold :exec
UPDATE stock_levels
SET on_hand = on_hand + @delta,
    reserved = reserved - @delta,
    updated_at = NOW()
WHERE product_id = $1
  AND reserved >= @delta;

-- name: FinalizeStockDeduction :exec
WITH updated_stock AS (
  UPDATE stock_levels sl
  SET reserved = reserved - @delta,
      updated_at = NOW()
  WHERE sl.product_id = $1
    AND sl.reserved >= @delta
  RETURNING product_id, on_hand, reserved
) INSERT INTO stock_adjustments (product_id, delta, reason, created_at, created_by)
SELECT $1, -@delta, 'order_fulfillment', NOW(), $2
FROM updated_stock us
RETURNING id;

-- name: SoftHoldStockForCart :one
WITH need AS (
  SELECT ci.product_id, ci.qty
  FROM carts c
  JOIN cart_items ci ON ci.cart_id = c.id
  WHERE c.user_id = $1
), upd AS (
  UPDATE stock_levels sl
  SET on_hand = sl.on_hand - ned.qty,
      reserved = sl.reserved + ned.qty,
      updated_at = NOW()
  FROM need ned
  WHERE sl.product_id = ned.product_id
    AND sl.on_hand >= ned.qty
  RETURNING sl.product_id
) SELECT
  (SELECT COUNT(*) FROM need) AS total_items,
  (SELECT COUNT(*) FROM upd) AS successfully_held_items;

-- name: FinalizeStockHoldForOrder :one
WITH order_items AS (
  SELECT oi.product_id, oi.qty
  FROM orders o
  JOIN order_items oi ON oi.order_id = o.id
  WHERE o.id = $1
), upd AS (
  UPDATE stock_levels sl
  SET reserved = sl.reserved - oi.qty,
      updated_at = NOW()
  FROM order_items oi
  WHERE sl.product_id = oi.product_id
    AND sl.reserved >= oi.qty
  RETURNING sl.product_id
), adj AS (
  INSERT INTO stock_adjustments (product_id, delta, reason, created_at, created_by)
  SELECT oi.product_id, -oi.qty, 'order_fulfillment', NOW(), $2
  FROM order_items oi
  JOIN upd u ON u.product_id = oi.product_id
  RETURNING id
) SELECT
  (SELECT COUNT(*) FROM order_items) AS total_items,
  (SELECT COUNT(*) FROM upd) AS successfully_finalized_items;
