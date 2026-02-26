package handlers

import (
	"net"

	log "github.com/sirupsen/logrus"

	"github.com/Coosis/go-eshop/internal/auth"
	"github.com/Coosis/go-eshop/internal/seckill"
	"github.com/labstack/echo/v4"
)

type SeckillHandler struct {
	Svc seckill.SeckillService
}

func RegisterSeckillRoutes(e *echo.Echo, Svc seckill.SeckillService) {
	handler := &SeckillHandler{Svc: Svc}

	g := e.Group("/v1/seckill")
	g.GET("/events", handler.GetSeckillEvents)
	g.GET("/events/:event_id", handler.GetSeckillEventByID)

	// admin routes
	// TODO! add admin auth middleware
	e.POST("/v1/admin/seckill/events", handler.AddSeckillEvent)
	e.PATCH("/v1/admin/seckill/events/:event_id", handler.UpdateSeckillEventByID)
	// e.POST("/v1/admin/seckill/outbox/preheat", handler.MarkPreheated)
	// e.GET("/v1/admin/seckill/outbox", handler.GetEventPreheat)

	g.POST("/events/:event_id/attempt", auth.RequireUserID(handler.AttemptSeckill))
	g.GET("/attempts/:idempotency_key/status", auth.RequireUserID(handler.GetSeckillPurchase))
}

func (h *SeckillHandler) GetSeckillEvents(e echo.Context) error {
	type Param struct {
		Page     int32 `query:"page"`
		PageSize int32 `query:"per_page"`
	}
	var p Param
	if err := e.Bind(&p); err != nil {
		log.Errorf("error binding request params: %v", err)
		return e.JSON(400, map[string]string{"error": "malformed query parameters"})
	}
	paged, err := h.Svc.GetSeckillEvents(e.Request().Context(), p.Page, p.PageSize)
	if err != nil {
		log.Errorf("error getting seckill events: %v", err)
		return e.JSON(500, map[string]string{"error": "internal server error"})
	}
	return e.JSON(200, paged)
}

func (h *SeckillHandler) GetSeckillEventByID(e echo.Context) error {
	type Param struct {
		ID int32 `param:"event_id"`
	}
	var p Param
	if err := e.Bind(&p); err != nil {
		log.Errorf("error binding request params: %v", err)
		return e.JSON(400, map[string]string{"error": "malformed event ID"})
	}
	event, err := h.Svc.GetSeckillEventByID(e.Request().Context(), p.ID)
	if err != nil {
		log.Errorf("error getting seckill event by ID: %v", err)
		return e.JSON(500, map[string]string{"error": "internal server error"})
	}
	return e.JSON(200, event)
}

func (h *SeckillHandler) AddSeckillEvent(e echo.Context) error {
	var req seckill.SeckillEventInfo
	if err := e.Bind(&req); err != nil {
		log.Errorf("error binding request body: %v", err)
		return e.JSON(400, map[string]string{"error": "malformed request body"})
	}
	log.Infof("received request to add seckill event: %+v", req)
	event, err := h.Svc.AddSeckillEvent(e.Request().Context(), req)
	if err != nil {
		log.Errorf("error adding seckill event: %v", err)
		return e.JSON(500, map[string]string{"error": "internal server error"})
	}
	return e.JSON(200, event)
}

func (h *SeckillHandler) UpdateSeckillEventByID(e echo.Context) error {
	type Param struct { ID int32 `param:"event_id"` }
	var p Param
	if err := e.Bind(&p); err != nil {
		log.Errorf("error binding request params: %v", err)
		return e.JSON(400, map[string]string{"error": "malformed event ID"})
	}

	var eventinfo seckill.SeckillEventInfo
	if err := e.Bind(&eventinfo); err != nil {
		log.Errorf("error binding request body: %v", err)
		return e.JSON(400, map[string]string{"error": "malformed request body"})
	}
	event, err := h.Svc.UpdateSeckillEventByID(e.Request().Context(), p.ID, eventinfo)
	if err != nil {
		log.Errorf("error updating seckill event: %v", err)
		return e.JSON(500, map[string]string{"error": "internal server error"})
	}
	return e.JSON(200, event)
}

func (h *SeckillHandler) AttemptSeckill(e echo.Context) error {
	userID, _ := e.Get(auth.UserIDKey).(int32)
	var attempt seckill.SeckillAttempt
	if err := e.Bind(&attempt); err != nil {
		log.Errorf("error binding request body: %v", err)
		return e.JSON(400, map[string]string{"error": "malformed request body"})
	}
	client_addr := e.Request().RemoteAddr
	client_ip, _, err := net.SplitHostPort(client_addr);
	if err != nil {
		log.Warnf("error parsing client IP from RemoteAddr: %v", err)
		return e.JSON(400, map[string]string{"error": "invalid client IP address"})
	}
	stats, err := h.Svc.PurchaseSeckillProduct(e.Request().Context(), userID, client_ip, attempt)
	if err != nil {
		log.Errorf("error attempting seckill: %v", err)
		return e.JSON(500, map[string]string{"error": "internal server error"})
	}
	return e.JSON(200, stats)
}

func (h *SeckillHandler) GetSeckillPurchase(e echo.Context) error {
	userID, _ := e.Get(auth.UserIDKey).(int32)
	type Param struct {
		Key string `param:"idempotency_key"`
	}
	var p Param
	if err := e.Bind(&p); err != nil {
		log.Errorf("error binding request params: %v", err)
		return e.JSON(400, map[string]string{"error": "malformed request ID"})
	}
	stats, err := h.Svc.GetSeckillPurchase(e.Request().Context(), userID, p.Key)
	if err != nil {
		log.Errorf("error getting seckill purchase status: %v", err)
		return e.JSON(500, map[string]string{"error": "internal server error"})
	}
	return e.JSON(200, stats)
}
