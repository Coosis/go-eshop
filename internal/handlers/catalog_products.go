package handlers

import (
	"github.com/jackc/pgx/v5"
	log "github.com/sirupsen/logrus"

	"github.com/Coosis/go-eshop/internal/catalog"
	"github.com/labstack/echo/v4"
)

type CatalogHandler struct {
	Svc catalog.CatalogService
}

func RegisterCatalogProductRoutes(e *echo.Echo, svc catalog.CatalogService) {
	handler := &CatalogHandler{Svc: svc}
	product_router := e.Group("/v1/catalog/products")
	product_router.GET("", handler.GetProducts)
	product_router.GET("/:id", handler.GetProductByID)
	product_router.GET("/slug/:slug", handler.GetProductBySlug)
}

func (h *CatalogHandler) GetProducts(c echo.Context) error {
	var f catalog.ProductFilter
	if err := c.Bind(&f); err != nil {
		log.Errorf("error binding filter params: %v", err)
		return c.JSON(400, map[string]string{"error": "malformed filter parameters"})
	}
	if err := c.Validate(f); err != nil {
		log.Errorf("validation error: %v", err)
		return c.JSON(400, map[string]string{"error": "invalid filter parameters"})
	}
	page, err := h.Svc.GetProducts(c.Request().Context(), f)
	if err != nil {
		if err == pgx.ErrNoRows {
			log.Warnf("no products found: %v", err)
			return c.JSON(404, map[string]string{"error": "no products found"})
		}
		log.Errorf("error getting products: %v", err)
		return c.JSON(500, map[string]string{"error": "internal server error"})
	}
	return c.JSON(200, page)
}

func (h *CatalogHandler) GetProductByID(c echo.Context) error {
	type params struct {
		ID int32 `param:"id" validate:"required"`
	}
	var p params
	if err := c.Bind(&p); err != nil {
		log.Errorf("error binding path params: %v", err)
		return c.JSON(400, map[string]string{"error": "invalid product ID"})
	}
	if err := c.Validate(p); err != nil {
		log.Errorf("validation error: %v", err)
		return c.JSON(400, map[string]string{"error": "invalid product ID"})
	}
	product, err := h.Svc.GetProductByID(c.Request().Context(), p.ID)
	if err != nil {
		if err == pgx.ErrNoRows {
			log.Warnf("no products found: %v", err)
			return c.JSON(404, map[string]string{"error": "no products found"})
		}
		log.Errorf("error getting product by ID: %v", err)
		return c.JSON(500, map[string]string{"error": "internal server error"})
	}
	return c.JSON(200, product)
}

func (h *CatalogHandler) GetProductBySlug(c echo.Context) error {
	type params struct {
		Slug string `param:"slug" validate:"required"`
	}
	var p params
	if err := c.Bind(&p); err != nil {
		log.Errorf("error binding path params: %v", err)
		return c.JSON(400, map[string]string{"error": "invalid product slug"})
	}
	if err := c.Validate(p); err != nil {
		log.Errorf("validation error: %v", err)
		return c.JSON(400, map[string]string{"error": "invalid product slug"})
	}
	product, err := h.Svc.GetProductBySlug(c.Request().Context(), p.Slug)
	if err != nil {
		if err == pgx.ErrNoRows {
			log.Warnf("no products found: %v", err)
			return c.JSON(404, map[string]string{"error": "no products found"})
		}
		log.Errorf("error getting product by slug: %v", err)
		return c.JSON(500, map[string]string{"error": "internal server error"})
	}
	return c.JSON(200, product)
}
