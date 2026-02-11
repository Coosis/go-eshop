package cart

import (
	"context"
	log "github.com/sirupsen/logrus"

	"github.com/Coosis/go-eshop/internal/comm"
	sqlc "github.com/Coosis/go-eshop/sqlc"
	"github.com/jackc/pgx/v5/pgtype"
)

type CartActor struct { 
	Pool sqlc.DBTX
}

func (ca *CartActor) GetCurrentCart(
	ctx context.Context,
	userID int32,
	p CartPaging,
) (Cart, error) {
	pg_userID := pgtype.Int4{
		Int32: userID,
		Valid: true,
	}
	queries := sqlc.New(ca.Pool)
	log.Infof("GetCurrentCart userID=%d page=%d perPage=%d", userID, p.Page, p.PerPage)
	rows, err := queries.GetCurrentCart(ctx, sqlc.GetCurrentCartParams{
		UserID: pg_userID,
		PageNumber: p.Page,
		PageSize: p.PerPage,
	})
	if err != nil {
		log.Errorf("error getting current cart: %v", err)
		return Cart{}, err
	}
	items := []CartItem{}
	if rows[0].TotalCount != 0 {
		log.Infof("GetCurrentCart found %d items", rows[0].TotalCount)
		for _, row := range rows {
			item := CartItem{
				ProductID: row.ProductID.Int32,
				Quantity: row.Qty.Int32,
				PriceCents: row.PriceCentsSnapshot.Int32,
			}
			items = append(items, item)
		}
	}
	paged := comm.Page[CartItem]{
		Items: items,
		Page: p.Page,
		PerPage: p.PerPage,
		Total: rows[0].TotalCount,
	}
	c := Cart{
		Version: rows[0].Version,
		Items:  paged,
	}
	return c, nil
}

func (ca *CartActor) UpdateCartItem(
	ctx context.Context,
	userID int32,
	req UpdateCartItemRequest,
) (Cart, error) {
	queries := sqlc.New(ca.Pool)
	rows, err := queries.UpdateCartItem(ctx, sqlc.UpdateCartItemParams{
		UserID: pgtype.Int4{
			Int32: userID,
			Valid: true,
		},
		ProductID: req.ProductID,
		Qty: req.Quantity,
		PageNumber: req.Page,
		PageSize: req.PerPage,
	})
	if err != nil {
		log.Errorf("error updating cart item: %v", err)
		return Cart{}, err
	}
	items := []CartItem{}
	if rows[0].TotalCount != 0 {
		log.Infof("GetCurrentCart found %d items", rows[0].TotalCount)
		for _, row := range rows {
			item := CartItem{
				ProductID: row.ProductID.Int32,
				Quantity: row.Qty.Int32,
				PriceCents: row.PriceCentsSnapshot.Int32,
			}
			items = append(items, item)
		}
	}
	paged := comm.Page[CartItem]{
		Items: items,
		Page: req.Page,
		PerPage: req.PerPage,
		Total: rows[0].TotalCount,
	}
	c := Cart{
		Version: rows[0].Version,
		Items: paged,
	}
	return c, nil
}

func (ca *CartActor) AddCartItem(
	ctx context.Context,
	userID int32,
	req AddCartItemRequest,
) (Cart, error) {
	queries := sqlc.New(ca.Pool)
	rows, err := queries.AddCartItem(ctx, sqlc.AddCartItemParams{
		UserID: pgtype.Int4{
			Int32: userID,
			Valid: true,
		},
		ID: req.ProductID,
		Qty: req.Quantity,
		PageNumber: req.Page,
		PageSize: req.PerPage,
	})
	if err != nil {
		log.Errorf("error adding cart item: %v", err)
		return Cart{}, err
	}
	items := []CartItem{}
	if rows[0].TotalCount != 0 {
		log.Infof("GetCurrentCart found %d items", rows[0].TotalCount)
		for _, row := range rows {
			item := CartItem{
				ProductID: row.ProductID.Int32,
				Quantity: row.Qty.Int32,
				PriceCents: row.PriceCentsSnapshot.Int32,
			}
			items = append(items, item)
		}
	}
	paged := comm.Page[CartItem]{
		Items: items,
		Page: req.Page,
		PerPage: req.PerPage,
		Total: rows[0].TotalCount,
	}
	c := Cart{
		Version: rows[0].Version,
		Items: paged,
	}
	return c, nil
}
func (ca *CartActor) ChangeCartItemQuantity(
	ctx context.Context,
	userID int32,
	req ChangeCartItemQuantityRequest,
) (Cart, error) {
	queries := sqlc.New(ca.Pool)
	rows, err := queries.ChangeCartItemQty(ctx, sqlc.ChangeCartItemQtyParams{
		UserID: pgtype.Int4{
			Int32: userID,
			Valid: true,
		},
		ProductID: req.ProductID,
		Delta: req.Delta,
		PageNumber: req.Page,
		PageSize: req.PerPage,
	})
	if err != nil {
		log.Errorf("error changing cart item quantity: %v", err)
		return Cart{}, err
	}
	items := []CartItem{}
	if rows[0].TotalCount != 0 {
		log.Infof("GetCurrentCart found %d items", rows[0].TotalCount)
		for _, row := range rows {
			item := CartItem{
				ProductID: row.ProductID.Int32,
				Quantity: row.Qty.Int32,
				PriceCents: row.PriceCentsSnapshot.Int32,
			}
			items = append(items, item)
		}
	}
	paged := comm.Page[CartItem]{
		Items: items,
		Page: req.Page,
		PerPage: req.PerPage,
		Total: rows[0].TotalCount,
	}
	c := Cart{
		Version: rows[0].Version,
		Items: paged,
	}
	return c, nil
}

func (ca *CartActor) RemoveCartItem(
	ctx context.Context,
	userID int32,
	req RemoveCartItemRequest,
) (Cart, error) {
	queries := sqlc.New(ca.Pool)
	rows, err := queries.RemoveCartItem(ctx, sqlc.RemoveCartItemParams{
		UserID: pgtype.Int4{
			Int32: userID,
			Valid: true,
		},
		ID: req.ProductID,
		PageNumber: req.Page,
		PageSize: req.PerPage,
	})
	if err != nil {
		log.Errorf("error removing cart item: %v", err)
		return Cart{}, err
	}
	items := []CartItem{}
	if rows[0].TotalCount != 0 {
		log.Infof("GetCurrentCart found %d items", rows[0].TotalCount)
		for _, row := range rows {
			item := CartItem{
				ProductID: row.ProductID.Int32,
				Quantity: row.Qty.Int32,
				PriceCents: row.PriceCentsSnapshot.Int32,
			}
			items = append(items, item)
		}
	}
	paged := comm.Page[CartItem]{
		Items: items,
		Page: req.Page,
		PerPage: req.PerPage,
		Total: rows[0].TotalCount,
	}
	c := Cart{
		Version: rows[0].Version,
		Items: paged,
	}
	return c, nil
}

func (ca *CartActor) ClearCart(
	ctx context.Context,
	userID int32,
) (Cart, error) {
	queries := sqlc.New(ca.Pool)
	rows, err := queries.ClearCart(
		ctx, 
		pgtype.Int4{
			Int32: userID,
			Valid: true,
		},
	)
	if err != nil {
		log.Printf("error clearing cart: %v", err)
		return Cart{}, err
	}
	paged := comm.Page[CartItem]{
		Items: []CartItem{},
		Page: 0,
		PerPage: 0,
		Total: 0,
	}
	c := Cart{
		Version: rows[0],
		Items: paged,
	}
	return c, nil
}
