package orders

import (
	"context"
	"fmt"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/Coosis/go-eshop/internal/comm"
	db "github.com/Coosis/go-eshop/sqlc"
	sqlc "github.com/Coosis/go-eshop/sqlc"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/bwmarrin/snowflake"
)

type OrderActor struct {
	Pool *pgxpool.Pool
	Node *snowflake.Node
}

func (o *OrderActor) GenerateOrderNumber() string {
	return o.Node.Generate().String()
}

func (o *OrderActor) PlaceOrder(
	ctx context.Context,
	req *PlaceOrderRequest,
) (*OrderInfo, error) {
	log.Infof("PlaceOrder called with request: %+v", req)
	tx, err := o.Pool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		log.Errorf("error starting transaction: %v", err)
		return nil, comm.InternalError
	}
	defer tx.Rollback(ctx)

	queries := sqlc.New(o.Pool).WithTx(tx)

	log.Infof("Trying to soft hold stock...")
	row, err := queries.SoftHoldStockForCart(ctx, pgtype.Int4{Int32: req.UserID, Valid: true})
	if err != nil {
		log.Errorf("error soft holding stock for userID=%d: %v", req.UserID, err)
		return nil, comm.InternalError
	}

	if row.SuccessfullyHeldItems != row.TotalItems {
		log.Errorf(
			"could not hold all items for userID=%d: held %d out of %d, error: %v",
			req.UserID,
			row.SuccessfullyHeldItems,
			row.TotalItems,
			err,
		)
		return nil, fmt.Errorf("Some items are out of stock")
	}

	paymentIntentID := pgtype.Text{Valid: false}
	if req.PaymentIntentID != nil {
		paymentIntentID = pgtype.Text{String: *req.PaymentIntentID, Valid: true}
	}

	notes := pgtype.Text{Valid: false}
	if req.Notes != nil {
		notes = pgtype.Text{String: *req.Notes, Valid: true}
	}

	r, err := queries.CreateOrder(ctx, sqlc.CreateOrderParams{
		UserID:          pgtype.Int4{Int32: req.UserID, Valid: true},
		OrderNumber:     o.GenerateOrderNumber(),
		DiscountCents:   0,
		PaymentIntentID: paymentIntentID,
		Notes:           notes,
		IdempotencyKey:  req.IdempotencyKey,
		Version:         req.CartVersion,
	})
	if err != nil {
		log.Errorf("error creating order: %v", err)
		return nil, comm.InternalError
	}

	status := string(r.Status)
	if err := tx.Commit(ctx); err != nil {
		log.Errorf("error committing transaction: %v", err)
		return nil, comm.InternalError
	}
	return &OrderInfo{
		OrderID:         r.ID,
		OrderNumber:     r.OrderNumber,
		SubtotalCents:   r.SubtotalCents,
		DiscountCents:   r.DiscountCents,
		TotalCents:      r.TotalCents,
		Status:          status,
		PaymentIntentID: &r.PaymentIntentID.String,
		Notes:           &r.Notes.String,
		CreatedAt:       r.CreatedAt.Time.UnixMilli(),
		Version:         r.Version,
	}, nil
}

func (o *OrderActor) GetOrders(
	ctx context.Context,
	f GetOrderRequest,
) (comm.Page[OrderInfo], error) {
	log.Infof("GetOrders called with filter: %+v", f)
	empty_page := comm.Page[OrderInfo]{}
	queries := sqlc.New(o.Pool)
	pg_before := pgtype.Timestamptz{Valid: false}
	if f.Before != nil {
		beft := time.UnixMilli(*f.Before)
		log.Infof("Order Before: unix millis: %v, timestamp: %v", *f.Before, beft)
		pg_before = pgtype.Timestamptz{Time: beft, Valid: true}
	}
	pg_after := pgtype.Timestamptz{Valid: false}
	if f.After != nil {
		aftt := time.UnixMilli(*f.After)
		log.Infof("Order After: unix millis: %v, timestamp: %v", *f.After, aftt)
		pg_after = pgtype.Timestamptz{Time: aftt, Valid: true}
	}
	pg_status := db.NullOrderStatus{Valid: false}
	if f.Status != nil {
		var odstat sqlc.OrderStatus
		if err := odstat.Scan(*f.Status); err != nil {
			log.Errorf("error scanning order status: %v", err)
			return empty_page, comm.InternalError
		}
		pg_status = db.NullOrderStatus{OrderStatus: odstat, Valid: true}
	}
	rows, err := queries.GetOrders(ctx, sqlc.GetOrdersParams{
		UserID:     f.UserID,
		Before:     pg_before,
		After:      pg_after,
		Status:     pg_status,
		PageNumber: f.Page,
		PageSize:   f.PerPage,
	})
	if err != nil {
		log.Errorf("error getting orders: %v", err)
		return empty_page, comm.InternalError
	}
	order_infos := []OrderInfo{}
	for _, r := range rows {
		status := string(r.Status)
		order_info := OrderInfo{
			OrderID:         r.ID,
			OrderNumber:     r.OrderNumber,
			SubtotalCents:   r.SubtotalCents,
			DiscountCents:   r.DiscountCents,
			TotalCents:      r.TotalCents,
			Status:          status,
			PaymentIntentID: &r.PaymentIntentID.String,
			Notes:           &r.Notes.String,
			CreatedAt:       r.CreatedAt.Time.UnixMilli(),
			Version:         r.Version,
		}
		order_infos = append(order_infos, order_info)
	}
	return comm.Page[OrderInfo]{
		Items:   order_infos,
		Page:    f.Page,
		PerPage: f.PerPage,
		Total: int64(len(rows)),
	}, nil
}

func (o *OrderActor) GetOrderByID(
	ctx context.Context,
	orderID int32,
	userID int32,
) (*OrderInfo, error) {
	log.Infof("GetOrderByID called with orderID: %d, userID: %d", orderID, userID)
	queries := sqlc.New(o.Pool)
	rows, err := queries.GetOrderByID(ctx, sqlc.GetOrderByIDParams{
		ID:     orderID,
		UserID: userID,
	})
	if err != nil {
		log.Errorf("error getting order by ID: %v", err)
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("order not found")
		}
		return nil, comm.InternalError
	}
	status := string(rows.Status)
	return &OrderInfo{
		OrderID:         rows.ID,
		OrderNumber:     rows.OrderNumber,
		SubtotalCents:   rows.SubtotalCents,
		DiscountCents:   rows.DiscountCents,
		TotalCents:      rows.TotalCents,
		Status:          status,
		PaymentIntentID: &rows.PaymentIntentID.String,
		Notes:           &rows.Notes.String,
		CreatedAt:       rows.CreatedAt.Time.UnixMilli(),
		Version:         rows.Version,
	}, nil
}
func (o *OrderActor) CancelOrder(
	ctx context.Context,
	orderID int32,
	userID int32,
) (*OrderInfo, error) {
	log.Infof("CancelOrder called with orderID: %d, userID: %d", orderID, userID)
	queries := sqlc.New(o.Pool)
	rows, err := queries.CancelOrder(ctx, sqlc.CancelOrderParams{
		ID:     orderID,
		UserID: userID,
	})
	if err != nil {
		log.Errorf("error cancelling order by ID: %v", err)
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("order not found")
		}
		return nil, comm.InternalError
	}
	status := string(rows.Status)
	return &OrderInfo{
		OrderID:         rows.ID,
		OrderNumber:     rows.OrderNumber,
		SubtotalCents:   rows.SubtotalCents,
		DiscountCents:   rows.DiscountCents,
		TotalCents:      rows.TotalCents,
		Status:          status,
		PaymentIntentID: &rows.PaymentIntentID.String,
		Notes:           &rows.Notes.String,
		CreatedAt:       rows.CreatedAt.Time.UnixMilli(),
		Version:         rows.Version,
	}, nil
}

func (o *OrderActor) PayOrder(
	ctx context.Context,
	req *PayOrderRequest,
) (*OrderInfo, error) {
	log.Infof("PayOrder called with request: %+v", req)
	tx, err := o.Pool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		log.Errorf("error starting transaction: %v", err)
		return nil, comm.InternalError
	}
	defer tx.Rollback(ctx)
	queries := sqlc.New(o.Pool).WithTx(tx)
	finrows, err := queries.FinalizeStockHoldForOrder(ctx, sqlc.FinalizeStockHoldForOrderParams{
		ID:    req.OrderID,
		CreatedBy: fmt.Sprintf("userID:%d", req.UserID),
	})
	if err != nil {
		log.Errorf("error finalizing stock hold for order: %v", err)
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("order not found")
		}
		return nil, comm.InternalError
	}
	if finrows.SuccessfullyFinalizedItems != finrows.TotalItems {
		log.Errorf(
			"could not finalize hold for all items for orderID=%d: finalized %d out of %d, error: %v",
			req.OrderID,
			finrows.SuccessfullyFinalizedItems,
			finrows.TotalItems,
			err,
		)
		return nil, fmt.Errorf("Some items are having stock issues...")
	}
	rows, err := queries.PayOrder(ctx, sqlc.PayOrderParams{
		ID:              req.OrderID,
		UserID:          req.UserID,
		PaymentIntentID: pgtype.Text{String: req.PaymentIntentID, Valid: true},
	})
	if err != nil {
		log.Errorf("error paying for order: %v", err)
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("order not found")
		}
		return nil, comm.InternalError
	}
	status := string(rows.Status)
	if err := tx.Commit(ctx); err != nil {
		log.Errorf("error committing transaction: %v", err)
		return nil, comm.InternalError
	}
	return &OrderInfo{
		OrderID:         rows.ID,
		OrderNumber:     rows.OrderNumber,
		SubtotalCents:   rows.SubtotalCents,
		DiscountCents:   rows.DiscountCents,
		TotalCents:      rows.TotalCents,
		Status:          status,
		PaymentIntentID: &rows.PaymentIntentID.String,
		Notes:           &rows.Notes.String,
		CreatedAt:       rows.CreatedAt.Time.UnixMilli(),
		Version:         rows.Version,
	}, nil
}

func (o *OrderActor) RefundOrder(
	ctx context.Context,
	orderID int32,
	userID int32,
) (*OrderInfo, error) {
	log.Infof("RefundOrder called with orderID: %d, userID: %d", orderID, userID)
	queries := sqlc.New(o.Pool)
	rows, err := queries.RefundOrder(ctx, sqlc.RefundOrderParams{
		ID:     orderID,
		UserID: userID,
	})
	if err != nil {
		log.Errorf("error refunding order by ID: %v", err)
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("order not found")
		}
		return nil, comm.InternalError
	}
	status := string(rows.Status)
	return &OrderInfo{
		OrderID:         rows.ID,
		OrderNumber:     rows.OrderNumber,
		SubtotalCents:   rows.SubtotalCents,
		DiscountCents:   rows.DiscountCents,
		TotalCents:      rows.TotalCents,
		Status:          status,
		PaymentIntentID: &rows.PaymentIntentID.String,
		Notes:           &rows.Notes.String,
		CreatedAt:       rows.CreatedAt.Time.UnixMilli(),
		Version:         rows.Version,
	}, nil
}

func (o *OrderActor) PayOrderWebhook(
	ctx context.Context,
	orderID int32,
	paymentIntentID string,
) error {
	log.Infof("PayOrderWebhook called with orderID: %d, paymentIntentID: %s", orderID, paymentIntentID)
	// TODO! implement webhook handling
	return nil
}
