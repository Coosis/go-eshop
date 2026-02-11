package handlers

import (
	"log"

	"github.com/Coosis/go-eshop/internal/auth"
	"github.com/Coosis/go-eshop/internal/orders"
	"github.com/labstack/echo/v4"
)

type OrderHandler struct {
	Svc orders.OrderService
}

func RegisterOrderRoutes(e *echo.Echo, Svc orders.OrderService) {
	handler := &OrderHandler{Svc: Svc}
	e.POST("/:id/webhook", handler.OrderWebhook)

	g := e.Group("/v1/orders")
	g.Use(auth.RequireUserID)
	g.POST("", handler.PlaceOrder)
	g.GET("", handler.GetOrders)
	g.GET("/:order_id", handler.GetOrderByID)
	g.POST("/:order_id/cancel", handler.CancelOrder)
	g.POST("/:order_id/pay", handler.PayOrder)
	g.POST("/:order_id/refund", handler.RefundOrder)
}

func (h *OrderHandler) PlaceOrder(e echo.Context) error {
	userID, _ := e.Get(auth.UserIDKey).(int32)
	var p orders.PlaceOrderRequest
	if err := e.Bind(&p); err != nil {
		log.Printf("error binding request params: %v", err)
		return e.JSON(400, map[string]string{"error": "malformed request parameters"})
	}
	p.UserID = userID
	orderInfo, err := h.Svc.PlaceOrder(e.Request().Context(), &p)
	if err != nil {
		log.Printf("error placing order: %v", err)
		return e.JSON(500, map[string]string{"error": err.Error()})
	}
	return e.JSON(200, orderInfo)
}

func (h *OrderHandler) GetOrders(e echo.Context) error {
	userID, _ := e.Get(auth.UserIDKey).(int32)
	var f orders.GetOrderRequest
	if err := e.Bind(&f); err != nil {
		log.Printf("error binding filter params: %v", err)
		return e.JSON(400, map[string]string{"error": "malformed filter parameters"})
	}
	f.UserID = userID
	infolist, err := h.Svc.GetOrders(e.Request().Context(), f)
	if err != nil {
		log.Printf("error getting orders: %v", err)
		return e.JSON(500, map[string]string{"error": err.Error()})
	}
	return e.JSON(200, infolist)
}

func (h *OrderHandler) GetOrderByID(e echo.Context) error {
	userID, _ := e.Get(auth.UserIDKey).(int32)
	type Params struct {
		ID int32 `param:"order_id" validate:"min=1,required"`
	}
	var p Params
	if err := e.Bind(&p); err != nil {
		log.Printf("error binding path params: %v", err)
		return e.JSON(400, map[string]string{"error": "invalid order ID"})
	}
	infolist, err := h.Svc.GetOrderByID(e.Request().Context(), p.ID, userID)
	if err != nil {
		log.Printf("error getting order by ID: %v", err)
		return e.JSON(500, map[string]string{"error": err.Error()})
	}
	return e.JSON(200, infolist)
}

func (h *OrderHandler) CancelOrder(e echo.Context) error {
	userID, _ := e.Get(auth.UserIDKey).(int32)
	type Params struct {
		ID int32 `param:"order_id" validate:"min=1,required"`
	}
	var p Params
	if err := e.Bind(&p); err != nil {
		log.Printf("error binding path params: %v", err)
		return e.JSON(400, map[string]string{"error": "invalid order ID"})
	}
	_, err := h.Svc.CancelOrder(e.Request().Context(), p.ID, userID)
	if err != nil {
		log.Printf("error canceling order: %v", err)
		return e.JSON(500, map[string]string{"error": err.Error()})
	}
	return e.JSON(200, "ok")
}

func (h *OrderHandler) PayOrder(e echo.Context) error {
	userID, _ := e.Get(auth.UserIDKey).(int32)
	var p orders.PayOrderRequest
	if err := e.Bind(&p); err != nil {
		log.Printf("error binding path params: %v", err)
		return e.JSON(400, map[string]string{"error": "invalid order ID"})
	}
	p.UserID = userID
	info, err := h.Svc.PayOrder(e.Request().Context(), &p)
	if err != nil {
		log.Printf("error paying for order: %v", err)
		return e.JSON(500, map[string]string{"error": err.Error()})
	}
	return e.JSON(200, info)
}

func (h *OrderHandler) RefundOrder(e echo.Context) error {
	userID, _ := e.Get(auth.UserIDKey).(int32)
	type Params struct {
		ID int32 `param:"order_id" validate:"min=1,required"`
	}
	var p Params
	if err := e.Bind(&p); err != nil {
		log.Printf("error binding path params: %v", err)
		return e.JSON(400, map[string]string{"error": "invalid order ID"})
	}

	info, err := h.Svc.RefundOrder(e.Request().Context(), p.ID, userID)
	if err != nil {
		log.Printf("error refunding order: %v", err)
		return e.JSON(500, map[string]string{"error": err.Error()})
	}
	return e.JSON(200, info)
}

func (h *OrderHandler) OrderWebhook(e echo.Context) error {
	// TODO! implement webhook handling
	return e.JSON(200, "order webhook received!")
}
