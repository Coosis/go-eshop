package handlers

import (
	log "github.com/sirupsen/logrus"

	"github.com/Coosis/go-eshop/internal/stock"
	"github.com/labstack/echo/v4"
)

type StockHandler struct { 
	Svc stock.StockService
}

func RegisterStockRoutes(e *echo.Echo, svc stock.StockService) {
	handler := &StockHandler{Svc: svc}
	e.GET("/v1/products/:id/stock", handler.GetStockLevel)
	// TODO! add admin auth middleware
	g := e.Group("/v1/admin/stock")
	g.POST("/adjustments", handler.AdjustStockLevel)
	g.GET("/adjustments", handler.GetStockAdjustments)
	g.GET("/adjustments/:id", handler.GetStockAdjustmentByID)
}

func(h *StockHandler) GetStockLevel(c echo.Context) error {
	type queryParams struct {
		ID int32 `param:"id" validate:"min=1,required"`
	}
	var params queryParams
	if err := c.Bind(&params); err != nil {
		log.Errorf("error binding path params: %v", err)
		return c.JSON(400, map[string]string{"error": "invalid product ID"})
	}
	level, err := h.Svc.GetStockLevel(c.Request().Context(), params.ID)
	if err != nil {
		log.Errorf("error getting stock level: %v", err)
		return c.JSON(500, map[string]string{"error": "internal server error"})
	}
	return c.JSON(200, level)
}

func(h *StockHandler) AdjustStockLevel(c echo.Context) error {
	var req stock.AdjustStockRequest
	if err := c.Bind(&req); err != nil {
		log.Errorf("error binding request body: %v", err)
		return c.JSON(400, map[string]string{"error": "malformed request body"})
	}
	level, err := h.Svc.AdjustStockLevel(c.Request().Context(), req)
	if err != nil {
		log.Errorf("error adjusting stock level: %v", err)
		return c.JSON(500, map[string]string{"error": "internal server error"})
	}
	return c.JSON(200, level)
}

func(h *StockHandler) GetStockAdjustments(c echo.Context) error {
	var f stock.StockAdjustmentFilter
	if err := c.Bind(&f); err != nil {
		log.Errorf("error binding filter params: %v", err)
		return c.JSON(400, map[string]string{"error": "malformed filter parameters"})
	}
	levels, err := h.Svc.GetStockAdjustments(c.Request().Context(), f)
	if err != nil {
		log.Errorf("error getting stock adjustments: %v", err)
		return c.JSON(500, map[string]string{"error": "internal server error"})
	}
	return c.JSON(200, levels)
}

func(h *StockHandler) GetStockAdjustmentByID(c echo.Context) error {
	type Param struct {
		ID int64 `param:"id" validate:"min=1,required"`
	}
	var p Param
	if err := c.Bind(&p); err != nil {
		log.Errorf("error binding request params: %v", err)
		return c.JSON(400, map[string]string{"error": "malformed adjustment ID"})
	}
	level, err := h.Svc.GetStockAdjustmentByID(c.Request().Context(), p.ID)
	if err != nil {
		log.Errorf("error getting stock adjustment by ID: %v", err)
		return c.JSON(500, map[string]string{"error": "internal server error"})
	}
	return c.JSON(200, level)
}
