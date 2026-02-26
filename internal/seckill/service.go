package seckill

import (
	"context"

	"github.com/Coosis/go-eshop/internal/comm"
)

type SeckillEvent struct {
	ID                int32 `json:"id"`
	ProductID         int32 `json:"product_id"`
	StartTime         int64 `json:"start_time"`
	EndTime           int64 `json:"end_time"`
	SeckillPriceCents int32 `json:"seckill_price_cents"`
	SeckillStock      int32 `json:"seckill_stock"`
}

// basically identical to SeckillEvent, but without ID
type SeckillEventInfo struct {
	ProductID         int32 `json:"product_id"`
	StartTime         int64 `json:"start_time"`
	EndTime           int64 `json:"end_time"`
	SeckillPriceCents int32 `json:"seckill_price_cents"`
	SeckillStock      int32 `json:"seckill_stock"`
}

type SeckillAttempt struct {
	EventID        int32 `json:"event_id"`
	Quantity       int64 `json:"quantity"`
	IdempotencyKey string `json:"idempotency_key"`
}

type SeckillAttemptStatus struct {
	State   string `json:"state"`
}

type SeckillService interface {
	// events
	GetSeckillEvents(
		ctx context.Context,
		page int32,
		pageSize int32,
	) (comm.Page[SeckillEvent], error)
	GetSeckillEventByID(ctx context.Context, id int32) (SeckillEvent, error)
	// admin
	AddSeckillEvent(ctx context.Context, new_event SeckillEventInfo) (SeckillEvent, error)
	UpdateSeckillEventByID(ctx context.Context, id int32, new_event SeckillEventInfo) (SeckillEvent, error)

	// purchase
	PurchaseSeckillProduct(
		ctx context.Context,
		userID int32,
		ip string,
		attempt SeckillAttempt,
	) (SeckillAttemptStatus, error)
	GetSeckillPurchase(
		ctx context.Context,
		userID int32,
		attempt_id string,
	) (SeckillAttemptStatus, error)
}
