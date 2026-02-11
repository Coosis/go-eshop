package handlers

import (
	log "github.com/sirupsen/logrus"

	"github.com/Coosis/go-eshop/internal/cart"
	"github.com/Coosis/go-eshop/internal/auth"
	"github.com/labstack/echo/v4"
)

type CartHandler struct {
	Svc cart.CartService
}

func RegisterCartRoutes(e *echo.Echo, svc cart.CartService) {
	handler := &CartHandler{Svc: svc}
	g := e.Group("/v1/cart")
	g.Use(auth.RequireUserID)
	g.GET("", handler.GetCart)
	g.PUT("/items/:product_id", handler.UpdateItem)
	g.POST("/items", handler.AddItem)
	g.PATCH("/items/:product_id", handler.ChangeItemQuantity)
	g.DELETE("/items/:product_id", handler.RemoveItem)
	g.DELETE("", handler.Clear)
	g.POST("/refresh-cart", handler.Refresh)
}

func (h *CartHandler) GetCart(c echo.Context) error {
	uid, _ := c.Get(auth.UserIDKey).(int32)
	type queryParams struct {
		Page int32 `query:"page" validate:"min=1,required"`
		PerPage int32 `query:"per_page" validate:"min=1,max=250,required"`
	}
	var params queryParams
	if err := c.Bind(&params); err != nil {
		log.Printf("error binding query params: %v", err)
		return c.JSON(400, map[string]string{"error": "malformed query parameters"})
	}
	log.Infof("getting cart for user ID: %d, page: %d, per_page: %d", uid, params.Page, params.PerPage)
	cart, err := h.Svc.GetCurrentCart(c.Request().Context(), uid, cart.CartPaging{
		Page:   params.Page,
		PerPage: params.PerPage,
	})
	if err != nil {
		log.Printf("error getting cart: %v", err)
		return c.JSON(500, map[string]string{"error": "internal server error"})
	}
	return c.JSON(200, cart)
}

func (h *CartHandler) UpdateItem(c echo.Context) error {
	uid, _ := c.Get(auth.UserIDKey).(int32)
	type queryParams struct {
		ProductID int32 `param:"product_id" validate:"min=1,required"`
		Quantity  int32 `json:"quantity" validate:"min=1,required"`
		Page int32 `json:"page" validate:"min=1,required"`
		PerPage int32 `json:"per_page" validate:"min=1,max=250,required"`
	}
	var params queryParams
	if err := c.Bind(&params); err != nil {
		log.Printf("error binding query params: %v", err)
		return c.JSON(400, map[string]string{"error": "malformed query parameters"})
	}
	log.Infof("updating cart item for user ID: %d, page: %d, per_page: %d", uid, params.Page, params.PerPage)
	mut, err := h.Svc.UpdateCartItem(
		c.Request().Context(),
		uid,
		cart.UpdateCartItemRequest{
			ProductID: params.ProductID,
			Quantity:  params.Quantity,
			CartPaging: cart.CartPaging{
				Page:    params.Page,
				PerPage: params.PerPage,
			},
		},
	)
	if err != nil {
		log.Printf("error updating cart item: %v", err)
		return c.JSON(500, map[string]string{"error": "internal server error"})
	}
	return c.JSON(200, mut)
}

func (h *CartHandler) AddItem(c echo.Context) error {
	uid, _ := c.Get(auth.UserIDKey).(int32)
	type queryParams struct {
		ProductID int32 `json:"product_id" validate:"min=1,required"`
		Quantity  int32 `json:"quantity" validate:"min=1,required"`
		Page int32 `json:"page" validate:"min=1,required"`
		PerPage int32 `json:"per_page" validate:"min=1,max=250,required"`
	}
	var params queryParams
	if err := c.Bind(&params); err != nil {
		log.Printf("error binding query params: %v", err)
		return c.JSON(400, map[string]string{"error": "malformed query parameters"})
	}
	log.Infof("adding an item for user ID: %d, page: %d, per_page: %d", uid, params.Page, params.PerPage)
	mut, err := h.Svc.AddCartItem(
		c.Request().Context(),
		uid,
		cart.AddCartItemRequest{
			ProductID: params.ProductID,
			Quantity:  params.Quantity,
			CartPaging: cart.CartPaging{
				Page:    params.Page,
				PerPage: params.PerPage,
			},
		},
	); 
	if err != nil {
		log.Printf("error adding item to cart: %v", err)
		return c.JSON(500, map[string]string{"error": "internal server error"})
	}
	return c.JSON(200, mut)
}

func (h *CartHandler) ChangeItemQuantity(c echo.Context) error {
	uid, _ := c.Get(auth.UserIDKey).(int32)
	type queryParams struct {
		ProductID int32 `param:"product_id" validate:"min=1,required"`
		Delta     int32 `query:"delta" validate:"min=-99,max=99,required"`
		Page int32 `query:"page" validate:"min=1,required"`
		PerPage int32 `query:"per_page" validate:"min=1,max=250,required"`
	}
	var params queryParams
	if err := c.Bind(&params); err != nil {
		log.Printf("error binding query params: %v", err)
		return c.JSON(400, map[string]string{"error": "malformed query parameters"})
	}
	mut, err := h.Svc.ChangeCartItemQuantity(
		c.Request().Context(),
		uid,
		cart.ChangeCartItemQuantityRequest{
			ProductID: params.ProductID,
			Delta:     params.Delta,
			CartPaging: cart.CartPaging{
				Page:    params.Page,
				PerPage: params.PerPage,
			},
		},
	);
	if err != nil {
		log.Printf("error changing cart item quantity: %v", err)
		return c.JSON(500, map[string]string{"error": "internal server error"})
	}
	return c.JSON(200, mut)
}

func (h *CartHandler) RemoveItem(c echo.Context) error {
	uid, _ := c.Get(auth.UserIDKey).(int32)
	type queryParams struct {
		ProductID int32 `param:"product_id" validate:"min=1,required"`
		Page int32 `query:"page" validate:"min=1,required"`
		PerPage int32 `query:"per_page" validate:"min=1,max=250,required"`
	}
	var params queryParams
	if err := c.Bind(&params); err != nil {
		log.Printf("error binding query params: %v", err)
		return c.JSON(400, map[string]string{"error": "malformed query parameters"})
	}
	mut, err := h.Svc.RemoveCartItem(
		c.Request().Context(),
		uid,
		cart.RemoveCartItemRequest{
			ProductID: params.ProductID,
			CartPaging: cart.CartPaging{
				Page:    params.Page,
				PerPage: params.PerPage,
			},
		},
	)
	if err != nil {
		log.Printf("error removing cart item: %v", err)
		return c.JSON(500, map[string]string{"error": "internal server error"})
	}
	return c.JSON(200, mut)
}

func (h *CartHandler) Clear(c echo.Context) error {
	uid, _ := c.Get(auth.UserIDKey).(int32)
	mut, err := h.Svc.ClearCart(c.Request().Context(), uid)
	if err != nil {
		log.Printf("error clearing cart: %v", err)
		return c.JSON(500, map[string]string{"error": "internal server error"})
	}
	return c.JSON(200, mut)
}

func (h *CartHandler) Refresh(c echo.Context) error {
	uid, _ := c.Get(auth.UserIDKey).(int32)
	type queryParams struct {
		ProductID int32 `param:"product_id" validate:"min=1,required"`
		Page int32 `query:"page" validate:"min=1,required"`
		PerPage int32 `query:"per_page" validate:"min=1,max=250,required"`
	}
	var params queryParams
	if err := c.Bind(&params); err != nil {
		log.Printf("error binding query params: %v", err)
		return c.JSON(400, map[string]string{"error": "malformed query parameters"})
	}
	mut, err := h.Svc.RefreshCart(c.Request().Context(), uid, cart.RefreshCartRequest{
		CartPaging: cart.CartPaging{
			Page:    params.Page,
			PerPage: params.PerPage,
		},
	})
	if err != nil {
		log.Printf("error refreshing cart: %v", err)
		return c.JSON(500, map[string]string{"error": "internal server error"})
	}
	return c.JSON(200, mut)
}
