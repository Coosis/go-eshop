package cart

import (
	"context"

	"github.com/Coosis/go-eshop/internal/comm"
)

type CartItem struct {
	ProductID  int32 `json:"product_id"`
	Quantity   int32 `json:"quantity"`
	PriceCents int32 `json:"price_cents"`
}

type Cart struct {
	Version int64               `json:"version"`
	Items   comm.Page[CartItem] `json:"paged_items"`
}

type CartPaging struct {
	Page    int32 `json:"page" query:"page" validate:"omitempty,min=1"`
	PerPage int32 `json:"per_page" query:"per_page" validate:"omitempty,min=1,max=250"`
}

type UpdateCartItemRequest struct {
	ProductID int32 `json:"product_id" validate:"min=1,required"`
	Quantity  int32 `json:"quantity" validate:"min=0,required"`
	CartPaging
}

type AddCartItemRequest struct {
	ProductID int32 `json:"product_id" validate:"min=1,required"`
	Quantity  int32 `json:"quantity" validate:"min=1,required"`
	CartPaging
}

type ChangeCartItemQuantityRequest struct {
	ProductID int32 `json:"product_id" validate:"min=1,required"`
	Delta     int32 `json:"delta" validate:"ne=0,required"`
	CartPaging
}

type RemoveCartItemRequest struct {
	ProductID int32 `json:"product_id" validate:"min=1,required"`
	CartPaging
}

type RefreshCartRequest struct {
	CartPaging
}

type CartService interface {
	GetCurrentCart(ctx context.Context, userID int32, p CartPaging) (Cart, error)
	UpdateCartItem(ctx context.Context, userID int32, req UpdateCartItemRequest) (Cart, error)
	AddCartItem(ctx context.Context, userID int32, req AddCartItemRequest) (Cart, error)
	ChangeCartItemQuantity(ctx context.Context, userID int32, req ChangeCartItemQuantityRequest) (Cart, error)
	RemoveCartItem(ctx context.Context, userID int32, req RemoveCartItemRequest) (Cart, error)
	ClearCart(ctx context.Context, userID int32) (Cart, error)
}
