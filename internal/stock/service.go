package stock

import (
	"context"
	"github.com/Coosis/go-eshop/internal/comm"
)

type StockAdjustment struct {
	ID        int64  `json:"id"`
	ProductID int32  `json:"product_id"`
	Delta     int32  `json:"delta"`
	Reason    string `json:"reason"`
	CreatedBy string `json:"created_by"`
	CreatedAt int64  `json:"created_at"`
}

type StockAdjustmentFilter struct {
	PerPage       int32  `query:"per_page" validate:"omitempty,min=1,max=250"`
	Page          int32  `query:"page" validate:"omitempty,min=1"`
	ProductID     int32  `query:"product_id" validate:"omitempty,min=1"`
	CreatedBy     string `query:"created_by"`
	CreatedAfter  *int64  `query:"created_after" validate:"omitempty,min=0"`
	CreatedBefore *int64  `query:"created_before" validate:"omitempty,min=0"`
	DeltaMin      int32  `query:"delta_min"`
	DeltaMax      int32  `query:"delta_max"`
}

type AdjustStockRequest struct {
	ProductID int32   `json:"product_id" validate:"min=1,required"`
	Delta     int32   `json:"delta" validate:"required,ne=0"`
	Reason    *string `json:"reason,omitempty"`
	CreatedBy string  `json:"created_by" validate:"required"`
}

type StockLevel struct {
	ProductID  int32 `json:"product_id"`
	StockLevel int32 `json:"stock_level"`
}

type StockService interface {
	GetStockLevel(ctx context.Context, productID int32) (StockLevel, error)
	AdjustStockLevel(ctx context.Context, req AdjustStockRequest) (StockLevel, error)
	GetStockAdjustments(ctx context.Context, filter StockAdjustmentFilter) (comm.Page[StockAdjustment], error)
	GetStockAdjustmentByID(ctx context.Context, id int64) (StockAdjustment, error)
	SoftHold(ctx context.Context, productID int32, quantity int32) error
	ReleaseHold(ctx context.Context, productID int32, quantity int32) error
	CommitSoftHold(ctx context.Context, productID int32, created_by string, quantity int32) error
}
