package orders

import (
	"context"

	"github.com/Coosis/go-eshop/internal/comm"
)

type PlaceOrderItem struct {
	ProductID int32 `json:"product_id"`
	Quantity  int32 `json:"quantity"`
}

type PlaceOrderRequest struct {
	UserID          int32

	IdempotencyKey  string  `json:"idempotency_key"`
	CartVersion     int64   `json:"cart_version"`
	Notes           *string `json:"notes,omitempty"`
	PaymentIntentID *string `json:"payment_intent_id,omitempty"`
}

type GetOrderRequest struct {
	UserID int32

	Before *int64 `query:"before" validate:"omitempty,min=0"`
	After  *int64 `query:"after" validate:"omitempty,min=0"`
	// TODO! verify states
	Status  *string `query:"status" validate:"omitempty,oneof= pending paid shipped delivered cancelled refunded"`
	Page    int32  `query:"page" validate:"min=1"`
	PerPage int32  `query:"per_page" validate:"min=1,max=250"`
}

type PayOrderRequest struct {
	UserID          int32

	OrderID         int32  `json:"order_id" validate:"min=1,required"`
	PaymentIntentID string `json:"payment_intent_id"`
}

type OrderInfo struct {
	OrderID       int32  `json:"order_id"`
	OrderNumber   string `json:"order_number"`
	SubtotalCents int64  `json:"subtotal_cents"`
	DiscountCents int64  `json:"discount_cents"`
	TotalCents    int64  `json:"total_cents"`

	Status          string  `json:"status"`
	PaymentIntentID *string `json:"payment_intent_id,omitempty"`
	Notes           *string `json:"notes,omitempty"`
	CreatedAt       int64   `json:"created_at"`
	Version         int64   `json:"version"`
}

// TODO! viewing products for order details
type OrderService interface {
	PlaceOrder(ctx context.Context, req *PlaceOrderRequest) (*OrderInfo, error)
	GetOrders(ctx context.Context, f GetOrderRequest) (comm.Page[OrderInfo], error)
	GetOrderByID(ctx context.Context, orderID int32, userID int32) (*OrderInfo, error)
	CancelOrder(ctx context.Context, orderID int32, userID int32) (*OrderInfo, error)
	PayOrder(ctx context.Context, req *PayOrderRequest) (*OrderInfo, error)
	RefundOrder(ctx context.Context, orderID int32, userID int32) (*OrderInfo, error)

	PayOrderWebhook(ctx context.Context, orderID int32, paymentIntentID string) error
}
