CREATE EXTENSION IF NOT EXISTS pgcrypto;

DROP TABLE IF EXISTS seckill_outbox;
DROP TABLE IF EXISTS seckill_events;

DROP TABLE IF EXISTS order_items;
DROP TABLE IF EXISTS orders;
DROP TYPE  IF EXISTS order_status;

DROP TABLE IF EXISTS cart_items;
DROP TABLE IF EXISTS carts;

DROP TABLE IF EXISTS stock_adjustments;
DROP TABLE IF EXISTS stock_levels;

DROP TABLE IF EXISTS product_categories;
DROP TABLE IF EXISTS categories;
DROP TABLE IF EXISTS products;

DROP TABLE IF EXISTS oauth_accounts;
DROP TABLE IF EXISTS contact_methods;
DROP TABLE IF EXISTS user_credentials;
DROP TABLE IF EXISTS users;

CREATE TABLE users (
	id            SERIAL PRIMARY KEY,
	created_at    TIMESTAMPTZ NOT NULL DEFAULT now(),
	last_login_at TIMESTAMPTZ
);

CREATE TABLE user_credentials (
	user_id       INTEGER PRIMARY KEY REFERENCES users(id) ON DELETE CASCADE,
	password_hash TEXT NOT NULL
);

CREATE TABLE contact_methods (
	id                 SERIAL      PRIMARY KEY,
	user_id            INTEGER     NOT NULL REFERENCES users(id) ON DELETE CASCADE,
	type               VARCHAR(31) NOT NULL CHECK (type IN ('email','phone')),
	value              TEXT        NOT NULL UNIQUE,
	is_verified        BOOLEAN     NOT NULL DEFAULT FALSE,
	created_at         TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE oauth_accounts (
	id                   SERIAL      PRIMARY KEY,
	user_id              INTEGER     NOT NULL REFERENCES users(id) ON DELETE CASCADE,
	provider             VARCHAR(31) NOT NULL CHECK (provider IN ('github' , 'google', 'facebook')),
	provider_user_id     TEXT        NOT NULL,
	created_at           TIMESTAMPTZ NOT NULL DEFAULT now(),
	UNIQUE (provider, provider_user_id)
);
CREATE INDEX ON oauth_accounts(user_id);

CREATE TABLE products (
  id             SERIAL PRIMARY KEY,
  name           TEXT NOT NULL,
  slug           TEXT NOT NULL UNIQUE,
  description    TEXT,
  price_cents    INTEGER NOT NULL,
  price_version  BIGINT NOT NULL DEFAULT 1,
  created_at     TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at     TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE categories (
  id         SERIAL  PRIMARY KEY,
  name       TEXT    NOT NULL UNIQUE,
  slug       TEXT    NOT NULL UNIQUE,
  parent_id  INTEGER REFERENCES categories(id) ON DELETE SET NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE product_categories (
  product_id  INTEGER  REFERENCES products(id) ON DELETE CASCADE,
  category_id INTEGER  REFERENCES categories(id) ON DELETE CASCADE,
  PRIMARY KEY (product_id, category_id)
);

CREATE TABLE stock_levels (
  product_id   INTEGER PRIMARY KEY REFERENCES products(id),
  on_hand      INTEGER NOT NULL DEFAULT 0,     -- physical units
  reserved     INTEGER NOT NULL DEFAULT 0,     -- soft holds
  updated_at   TIMESTAMPTZ NOT NULL DEFAULT now(),
  CHECK (on_hand >= 0),
  CHECK (reserved >= 0)
);

-- Immutable journal for adjustments & audit trail.
CREATE TABLE stock_adjustments (
  id            BIGSERIAL PRIMARY KEY,
	product_id    INTEGER NOT NULL REFERENCES products(id),
  delta         INTEGER NOT NULL, -- +receiving, -damage, etc.
  reason        TEXT NOT NULL,    -- 'receiving','return','correction',...
  created_by    TEXT NOT NULL,
  created_at    TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE carts (
  id           SERIAL PRIMARY KEY,
  user_id      INTEGER REFERENCES users(id) ON DELETE SET NULL,
  version      BIGINT NOT NULL DEFAULT 1,
  created_at   TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at   TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

create table cart_items (
  id                   SERIAL  PRIMARY KEY,
  cart_id              INTEGER NOT NULL REFERENCES carts(id) ON DELETE CASCADE,
  product_id           INTEGER NOT NULL REFERENCES products(id),
  qty                  INTEGER NOT NULL CHECK (qty > 0),
  price_cents_snapshot INTEGER NOT NULL,
  created_at           TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at           TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  unique (cart_id, product_id)
);

CREATE TYPE order_status AS ENUM (
  'waiting_payment',
  'paid',
  'canceled',
  'payment_failed',
  'refunded'
);

CREATE TABLE orders (
  id                  SERIAL PRIMARY KEY,
  order_number        TEXT UNIQUE NOT NULL, -- human-friendly; ULID/base36 OK

  user_id             INTEGER NOT NULL REFERENCES users(id),

  subtotal_cents      BIGINT NOT NULL,
  discount_cents      BIGINT NOT NULL DEFAULT 0,
  total_cents         BIGINT NOT NULL,
  status              order_status NOT NULL DEFAULT 'waiting_payment',

  payment_intent_id   TEXT,                 -- from Payment Service/PSP

  notes               TEXT,

  created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  version             BIGINT NOT NULL DEFAULT 1
);

CREATE TABLE order_items (
  id                  SERIAL PRIMARY KEY,
  order_id            INTEGER NOT NULL REFERENCES orders(id) ON DELETE CASCADE,
  product_id          INTEGER NOT NULL REFERENCES products(id),    -- still store for BI; don't join for reads
  product_name        TEXT NOT NULL,         -- snapshot for receipt
  qty                 INTEGER  NOT NULL CHECK (qty > 0),
  unit_price_cents    BIGINT NOT NULL,       -- snapshot at purchase time
  discount_cents      BIGINT NOT NULL DEFAULT 0,
  price_version       BIGINT,                -- from Pricing/Catalog for later forensics
  metadata            JSONB                  -- arbitrary snapshot (attributes, options)
);
CREATE INDEX ON order_items(order_id);

CREATE TABLE seckill_events (
    id              SERIAL PRIMARY KEY,
    product_id      INTEGER NOT NULL REFERENCES products(id),
    start_time      TIMESTAMPTZ NOT NULL,
    end_time        TIMESTAMPTZ NOT NULL,
    seckill_price_cents INTEGER NOT NULL CHECK (seckill_price_cents >= 0),
    seckill_stock   INTEGER NOT NULL CHECK (seckill_stock >= 0),
    created_at      TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT now(),
    CHECK (end_time > start_time),
    UNIQUE (product_id, start_time)
);

CREATE TABLE seckill_outbox (
    id              BIGSERIAL PRIMARY KEY,
    event_id        INTEGER NOT NULL REFERENCES products(id),
    created_at      TIMESTAMPTZ NOT NULL DEFAULT now(),
    preheated_at    TIMESTAMPTZ,
    scheduled_at    TIMESTAMPTZ NOT NULL
);
CREATE INDEX ON seckill_outbox(preheated_at);
