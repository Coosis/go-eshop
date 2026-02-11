package handlers

import (
	log "github.com/sirupsen/logrus"

	"github.com/Coosis/go-eshop/internal/catalog"
	"github.com/labstack/echo/v4"
)

func RegisterCatalogAdminRoutes(e *echo.Echo, svc catalog.CatalogService) {
	handler := &CatalogHandler{Svc: svc}
	g := e.Group("/v1/admin/catalog")
	g.POST("/products", handler.CreateProduct)
	g.PATCH("/products/:id", handler.UpdateProduct)
	// g.DELETE("/products/:id", handler.DeleteProduct)
	g.POST("/categories", handler.CreateCategory)
	g.PATCH("/categories/:id", handler.UpdateCategory)
	// g.DELETE("/categories/:id", handler.DeleteCategory)
}

func (h *CatalogHandler) CreateProduct(c echo.Context) error {
	var req catalog.CreateProductRequest
	if err := c.Bind(&req); err != nil {
		log.Printf("error binding request body: %v", err)
		return c.JSON(400, map[string]string{"error": "malformed request body"})
	}
	p, err := h.Svc.CreateProduct(c.Request().Context(), req)
	if err != nil {
		log.Printf("error creating product: %v", err)
		return c.JSON(500, map[string]string{"error": "internal server error"})
	}
	return c.JSON(200, p)
}
func (h *CatalogHandler) UpdateProduct(c echo.Context) error {
	var req catalog.UpdateProductRequest
	if err := c.Bind(&req); err != nil {
		log.Printf("error binding request body: %v", err)
		return c.JSON(400, map[string]string{"error": "malformed request body"})
	}
	p, err := h.Svc.UpdateProductByID(c.Request().Context(), req)
	if err != nil {
		log.Printf("error updating product: %v", err)
		return c.JSON(500, map[string]string{"error": "internal server error"})
	}
	return c.JSON(200, p)
}
func (h *CatalogHandler) CreateCategory(c echo.Context) error {
	var req catalog.CreateCategoryRequest
	if err := c.Bind(&req); err != nil {
		log.Printf("error binding request body: %v", err)
		return c.JSON(400, map[string]string{"error": "malformed request body"})
	}
	category, err := h.Svc.CreateCategory(c.Request().Context(), req)
	if err != nil {
		log.Printf("error creating category: %v", err)
		return c.JSON(500, map[string]string{"error": "internal server error"})
	}
	return c.JSON(200, category)
}
func (h *CatalogHandler) UpdateCategory(c echo.Context) error {
	var req catalog.UpdateCategoryRequest
	if err := c.Bind(&req); err != nil {
		log.Printf("error binding request body: %v", err)
		return c.JSON(400, map[string]string{"error": "malformed request body"})
	}
	category, err := h.Svc.UpdateCategoryByID(c.Request().Context(), req)
	if err != nil {
		log.Printf("error creating category: %v", err)
		return c.JSON(500, map[string]string{"error": "internal server error"})
	}
	return c.JSON(200, category)
}
