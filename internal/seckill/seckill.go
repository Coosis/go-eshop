package seckill

import (
	"context"
	"fmt"
	"time"

	sqlc "github.com/Coosis/go-eshop/sqlc"
	log "github.com/sirupsen/logrus"

	"github.com/Coosis/go-eshop/internal/comm"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/redis/go-redis/v9"
)

type SeckillActor struct {
	Pool   sqlc.DBTX
	Client *redis.Client
}

func (s *SeckillActor) GetSeckillEvents(
	ctx context.Context,
	page int32,
	pageSize int32,
) (comm.Page[SeckillEvent], error) {
	log.Infof("Getting seckill events: page=%d, pageSize=%d", page, pageSize)
	queries := sqlc.New(s.Pool)
	rows, err := queries.GetSeckillEvents(ctx, sqlc.GetSeckillEventsParams{
		PageNumber: page,
		PageSize:   pageSize,
	})
	if err != nil {
		if err == pgx.ErrNoRows {
			return comm.Page[SeckillEvent]{
				Items:   []SeckillEvent{},
				Page:    page,
				PerPage: pageSize,
				Total:   0,
			}, nil
		}
		return comm.Page[SeckillEvent]{}, err
	}
	items := make([]SeckillEvent, len(rows))
	for i, row := range rows {
		items[i] = SeckillEvent{
			ID:                row.ID,
			ProductID:         row.ProductID,
			StartTime:         row.StartTime.Time.UnixMilli(),
			EndTime:           row.EndTime.Time.UnixMilli(),
			SeckillPriceCents: row.SeckillPriceCents,
			SeckillStock:      row.SeckillStock,
		}
	}
	return comm.Page[SeckillEvent]{
		Items:   items,
		Page:    page,
		PerPage: pageSize,
		Total:   int64(len(items)),
	}, nil
}

func (s *SeckillActor) GetSeckillEventByID(ctx context.Context, id int32) (SeckillEvent, error) {
	if !comm.BFExists(
		ctx,
		s.Client,
		comm.BF_products_id,
		id,
	) {
		log.Warnf("product ID %d not found in bloom filter", id)
		return SeckillEvent{}, fmt.Errorf("product ID %d not found", id)
	}

	log.Infof("Getting seckill event by ID: %d", id)
	queries := sqlc.New(s.Pool)
	row, err := queries.GetSeckillEventByID(ctx, id)
	if err != nil {
		return SeckillEvent{}, fmt.Errorf("failed to get seckill event by ID: %w", err)
	}
	return SeckillEvent{
		ID:                row.ID,
		ProductID:         row.ProductID,
		StartTime:         row.StartTime.Time.UnixMilli(),
		EndTime:           row.EndTime.Time.UnixMilli(),
		SeckillPriceCents: row.SeckillPriceCents,
		SeckillStock:      row.SeckillStock,
	}, nil

}

func (s *SeckillActor) AddSeckillEvent(ctx context.Context, new_event SeckillEventInfo) (SeckillEvent, error) {
	log.Infof("Adding seckill event: ProductID=%d, StartTime=%d, EndTime=%d, PriceCents=%d, Stock=%d",
		new_event.ProductID,
		new_event.StartTime,
		new_event.EndTime,
		new_event.SeckillPriceCents,
		new_event.SeckillStock,
	)
	queries := sqlc.New(s.Pool)
	pg_starttime := pgtype.Timestamptz{Valid: true, Time: time.UnixMilli(new_event.StartTime)}
	pg_endtime := pgtype.Timestamptz{Valid: true, Time: time.UnixMilli(new_event.EndTime)}
	row, err := queries.AddSeckillEvent(ctx, sqlc.AddSeckillEventParams{
		ProductID:         new_event.ProductID,
		StartTime:         pg_starttime,
		EndTime:           pg_endtime,
		SeckillPriceCents: new_event.SeckillPriceCents,
		SeckillStock:      new_event.SeckillStock,
	})
	if err != nil {
		return SeckillEvent{}, fmt.Errorf("failed to add seckill event: %w", err)
	}
	_, err = comm.BFAdd(
		ctx,
		s.Client,
		comm.BF_seckill_events,
		fmt.Sprintf("%d", row.ID),
	).Result()
	if err != nil {
		log.Errorf("Failed to add seckill event ID=%d to Bloom filter: %v", row.ID, err)
	}
	return SeckillEvent{
		ID:                row.ID,
		ProductID:         row.ProductID,
		StartTime:         row.StartTime.Time.UnixMilli(),
		EndTime:           row.EndTime.Time.UnixMilli(),
		SeckillPriceCents: row.SeckillPriceCents,
		SeckillStock:      row.SeckillStock,
	}, nil
}

func (s *SeckillActor) UpdateSeckillEventByID(
	ctx context.Context,
	id int32,
	new_event SeckillEventInfo,
) (SeckillEvent, error) {
	log.Infof("Updating seckill event by ID: %d, ProductID=%d, StartTime=%d, EndTime=%d, PriceCents=%d, Stock=%d",
		id,
		new_event.ProductID,
		new_event.StartTime,
		new_event.EndTime,
		new_event.SeckillPriceCents,
		new_event.SeckillStock,
	)
	queries := sqlc.New(s.Pool)
	pg_starttime := pgtype.Timestamptz{Valid: true, Time: time.UnixMilli(new_event.StartTime)}
	pg_endtime := pgtype.Timestamptz{Valid: true, Time: time.UnixMilli(new_event.EndTime)}
	row, err := queries.UpdateSeckillEventByID(ctx, sqlc.UpdateSeckillEventByIDParams{
		ID:                id,
		ProductID:         new_event.ProductID,
		StartTime:         pg_starttime,
		EndTime:           pg_endtime,
		SeckillPriceCents: new_event.SeckillPriceCents,
		SeckillStock:      new_event.SeckillStock,
	})
	if err != nil {
		return SeckillEvent{}, fmt.Errorf("failed to update seckill event by ID: %w", err)
	}
	_, err = comm.BFAdd(
		ctx,
		s.Client,
		comm.BF_seckill_events,
		fmt.Sprintf("%d", row.ID),
	).Result()
	if err != nil {
		log.Errorf("Failed to add seckill event ID=%d to Bloom filter after update: %v", row.ID, err)
	}
	return SeckillEvent{
		ID:                row.ID,
		ProductID:         row.ProductID,
		StartTime:         row.StartTime.Time.UnixMilli(),
		EndTime:           row.EndTime.Time.UnixMilli(),
		SeckillPriceCents: row.SeckillPriceCents,
		SeckillStock:      row.SeckillStock,
	}, nil
}

func (s *SeckillActor) PurchaseSeckillProduct(
	ctx context.Context,
	userID int32,
	attempt SeckillAttempt,
) (SeckillAttemptStatus, error) {
	log.Infof("User %d attempting to purchase seckill product: EventID=%d, Quantity=%d, IdempotencyKey=%s",
		userID,
		attempt.EventID,
		attempt.Quantity,
		attempt.IdempotencyKey,
	)
	lua := `
	local event_id = KEYS[1]

	local exists = redis.call("BF.EXISTS", "bf_seckill_events", event_id)
	if exists == 0 then
		return "ERR:event_not_found"
	end

	local idempotency_key = KEYS[2]
	local qty = tonumber(ARGV[1])
	
	local attempt_key = "seckill_attempt:" .. idempotency_key

	local not_started_key = "seckill_event:" .. event_id .. ":not_started"
	if redis.call("EXISTS", not_started_key) == 1 then
		return "ERR:event_not_started"
	end

	local prev = redis.call("GET", attempt_key)
	if prev ~= nil and prev ~= false then
		return prev
	end

	if qty == nil or qty <= 0 then
		return "ERR:invalid_qty"
	end

	local idem_ttl_ms = tonumber(redis.call("GET", "seckill_event:" .. event_id .. ":event_len_ms") or "0")
	if idem_ttl_ms <= 0 then
		idem_ttl_ms = 300000
	end
	local stock_key = "seckill_event:" .. event_id .. ":stock"
	local stock = tonumber(redis.call("GET", stock_key) or "0")
	if stock < qty then
		local res = "ERR:out_of_stock"
	    redis.call("PSETEX", attempt_key, idem_ttl_ms, res)
		return res
	end

	local price_cents = tonumber(redis.call("GET", "seckill_event:" .. event_id .. ":price") or "0")
	if price_cents <= 0 then
		return "ERR:event_not_active"
	end
	redis.call(
		"XADD", "seckill_order_stream", "*", 
		"event_id", event_id,
		"quantity", qty,
		"idempotency_key", idempotency_key,
		"price_cents", price_cents
	)

	local new_stock = redis.call("DECRBY", stock_key, qty)
	local res = "OK:" .. new_stock
	redis.call("PSETEX", attempt_key, idem_ttl_ms, res)
	return res
	`
	result, err := s.Client.Eval(
		ctx, 
		lua, 
		// keys
		[]string{
			fmt.Sprintf("%d", attempt.EventID),
			fmt.Sprintf("%d:%s", userID, attempt.IdempotencyKey),
		},
		// args
		fmt.Sprintf("%d", attempt.Quantity),
	).Result()
	if err != nil {
		log.Errorf("Failed to execute Lua script for seckill purchase attempt: %v", err)
		return SeckillAttemptStatus{}, err
	}
	log.Infof("Seckill purchase attempt result: %s", result)
	return SeckillAttemptStatus{
		State:   result.(string),
	}, nil
}

func (s *SeckillActor) GetSeckillPurchase(
	ctx context.Context,
	userID int32,
	attempt_id string,
) (SeckillAttemptStatus, error) {
	// only 3 states:
	// - OK:xxx
	// - ERR:out_of_stock
	// - nil
	// if nil, we treat it as not processed yet
	lua := `
	local attempt_key = "seckill_attempt:" .. KEYS[1]
	local res = redis.call("GET", attempt_key)
	if res ~= nil and res ~= false then
		return res
	end
	return nil
	`
	result, err := s.Client.Eval(
		ctx, 
		lua,
		[]string{fmt.Sprintf("%d:%s", userID, attempt_id)},
	).Result()
	if err != nil {
		if err == redis.Nil {
			return SeckillAttemptStatus{
				State: "not_received",
			}, nil
		}
		log.Errorf("Failed to execute Lua script for getting seckill purchase status: %v", err)
		return SeckillAttemptStatus{}, err
	}
	state := "queued"
	if result == nil {
		state = "not_received"
	} else if resStr, ok := result.(string); ok {
		if len(resStr) >= 3 && resStr[:3] == "ERR" {
			state = "out_of_stock"
		}
	}
	return SeckillAttemptStatus{
		State: state,
	}, nil
}
