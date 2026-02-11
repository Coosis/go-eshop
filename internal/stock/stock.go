package stock

import (
	"context"
	"time"

	"github.com/Coosis/go-eshop/internal/comm"
	sqlc "github.com/Coosis/go-eshop/sqlc"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	log "github.com/sirupsen/logrus"
)

type StockActor struct {
	Pool sqlc.DBTX
}

// admin

func (s *StockActor) GetStockLevel(
	ctx context.Context,
	productID int32,
) (StockLevel, error) {
	log.Infof("GetStockLevel called for productID=%d", productID)
	queries := sqlc.New(s.Pool)
	row, err := queries.GetStockLevel(ctx, productID)
	if err != nil {
		log.Errorf("error getting stock level for productID=%d: %v", productID, err)
		return StockLevel{}, err
	}
	return StockLevel{
		ProductID: row.ProductID,
		StockLevel: row.OnHand,
	}, nil
}

func (s *StockActor) AdjustStockLevel(
	ctx context.Context,
	req AdjustStockRequest,
) (StockLevel, error) {
	log.Infof("AdjustStockLevel called for productID=%d delta=%d", req.ProductID, req.Delta)
	queries := sqlc.New(s.Pool)
	reason := ""
	if req.Reason != nil {
		reason = *req.Reason
	}

	row, err := queries.AdjustStockLevel(ctx, sqlc.AdjustStockLevelParams{
		ProductID: req.ProductID,
		Delta:     req.Delta,
		Reason:    reason,
		CreatedBy: req.CreatedBy,
	})
	if err != nil {
		log.Errorf("error adjusting stock level for productID=%d: %v", req.ProductID, err)
		return StockLevel{}, err
	}
	return StockLevel{
		ProductID: row.ProductID,
		StockLevel: row.OnHand,
	}, nil
}

// audits

func (s *StockActor) GetStockAdjustments(
	ctx context.Context,
	filter StockAdjustmentFilter,
) (comm.Page[StockAdjustment], error) {
	log.Infof("GetStockAdjustments called with filter: %+v", filter)
	queries := sqlc.New(s.Pool)
	pg_createdAfter := pgtype.Timestamptz{Valid: false}
	if filter.CreatedAfter != nil {
		pg_createdAfter = pgtype.Timestamptz{
			Time:  time.Unix(*filter.CreatedAfter, 0),
			Valid: true,
		}
	}
	pg_createdBefore := pgtype.Timestamptz{Valid: false}
	if filter.CreatedBefore != nil {
		pg_createdBefore = pgtype.Timestamptz{
			Time:  time.Unix(*filter.CreatedBefore, 0),
			Valid: true,
		}
	}
	rows, err := queries.GetStockAdjustments(ctx, sqlc.GetStockAdjustmentsParams{
		ProductID:     filter.ProductID,
		CreatedBy:     filter.CreatedBy,
		CreatedAfter:  pg_createdAfter,
		CreatedBefore: pg_createdBefore,
		MinDelta:      filter.DeltaMin,
		MaxDelta:      filter.DeltaMax,
		PageNumber:    filter.Page,
		PageSize:      filter.PerPage,
	})
	if err != nil {
		if err == pgx.ErrNoRows {
			log.Infof("no stock adjustments found for given filter")
			return comm.Page[StockAdjustment]{
				Items:   []StockAdjustment{},
				Page:    filter.Page,
				PerPage: filter.PerPage,
				Total:   0,
			}, nil
		}
		log.Errorf("error getting stock adjustments: %v", err)
		return comm.Page[StockAdjustment]{}, err
	}
	adjustments := []StockAdjustment{}
	for _, row := range rows {
		adj := StockAdjustment{
			ID:        row.ID,
			ProductID: row.ProductID,
			Delta:     row.Delta,
			Reason:    row.Reason,
			CreatedBy: row.CreatedBy,
			CreatedAt: row.CreatedAt.Time.UnixMilli(),
		}
		adjustments = append(adjustments, adj)
	}
	return comm.Page[StockAdjustment]{
		Items:   adjustments,
		Page:    filter.Page,
		PerPage: filter.PerPage,
		Total:   int64(len(adjustments)),
	}, nil
}

func (s *StockActor) GetStockAdjustmentByID(
	ctx context.Context,
	id int64,
) (StockAdjustment, error) {
	log.Infof("GetStockAdjustmentByID called for id=%d", id)
	queries := sqlc.New(s.Pool)
	row, err := queries.GetStockAdjustmentByID(ctx, id)
	if err != nil {
		log.Errorf("error getting stock adjustment by id=%d: %v", id, err)
		return StockAdjustment{}, err
	}
	return StockAdjustment{
		ID:        row.ID,
		ProductID: row.ProductID,
		Delta:     row.Delta,
		Reason:    row.Reason,
		CreatedBy: row.CreatedBy,
		CreatedAt: row.CreatedAt.Time.UnixMilli(),
	}, nil
}

// softhold related methods

func (s *StockActor) SoftHold(
	ctx context.Context,
	productID int32,
	quantity int32,
) error {
	log.Infof("SoftHold called for productID=%d quantity=%d", productID, quantity)
	queries := sqlc.New(s.Pool)
	err := queries.SoftHoldStock(ctx, sqlc.SoftHoldStockParams{
		ProductID: productID,
		Delta: quantity,
	})
	if err != nil {
		log.Errorf("error placing soft hold for productID=%d: %v", productID, err)
		return err
	}
	return nil
}

func (s *StockActor) ReleaseHold(
	ctx context.Context,
	productID int32,
	quantity int32,
) error {
	log.Infof("ReleaseHold called for productID=%d quantity=%d", productID, quantity)
	queries := sqlc.New(s.Pool)
	err := queries.ReleaseStockHold(ctx, sqlc.ReleaseStockHoldParams{
		ProductID: productID,
		Delta:  quantity,
	})
	if err != nil {
		log.Errorf("error releasing hold for productID=%d: %v", productID, err)
		return err
	}
	return nil
}

func (s *StockActor) CommitSoftHold(
	ctx context.Context,
	productID int32,
	created_by string,
	quantity int32,
) error {
	log.Infof("CommitSoftHold called for productID=%d quantity=%d", productID, quantity)
	queries := sqlc.New(s.Pool)
	err := queries.FinalizeStockDeduction(ctx, sqlc.FinalizeStockDeductionParams{
		ProductID: productID,
		CreatedBy: created_by,
		Delta:  quantity,
	})
	if err != nil {
		log.Errorf("error committing soft hold for productID=%d: %v", productID, err)
		return err
	}
	return nil
}
