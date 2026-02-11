package handlers

import (
	"log"

	"github.com/Coosis/go-eshop/internal/catalog"
	"github.com/labstack/echo/v4"
)

func RegisterCatalogCategoryRoutes(e *echo.Echo, svc catalog.CatalogService) {
	handler := &CatalogHandler{Svc: svc}
	category_router := e.Group("/v1/catalog/categories")
	category_router.GET("", handler.GetCategories)
	category_router.GET("/:id", handler.GetCategoryByID)
	category_router.GET("/slug/:slug", handler.GetCategoryBySlug)
	category_router.GET("/:id/products", handler.GetProductsByCategoryID)
}

func (h *CatalogHandler) GetCategories(c echo.Context) error {
	type Params struct {
		Page	int32 `query:"page" validate:"omitempty,min=1"`
		PerPage int32 `query:"per_page" validate:"omitempty,min=1,max=250"`
	}
	var p Params
	if err := c.Bind(&p); err != nil {
		log.Printf("error binding query params: %v", err)
		return c.JSON(400, map[string]string{"error": "malformed query parameters"})
	}
	categories, err := h.Svc.GetCategories(c.Request().Context(), catalog.GetCategoriesRequest{
		Page:    p.Page,
		PerPage: p.PerPage,
	})
	if err != nil {
		log.Printf("error getting categories: %v", err)
		return c.JSON(500, map[string]string{"error": "internal server error"})
	}
	return c.JSON(200, categories)
}

func (h *CatalogHandler) GetCategoryByID(c echo.Context) error {
	type Params struct {
		ID int32 `param:"id" validate:"required"`
	}
	var p Params
	if err := c.Bind(&p); err != nil {
		log.Printf("error binding path params: %v", err)
		return c.JSON(400, map[string]string{"error": "invalid category ID"})
	}

	category, err := h.Svc.GetCategoryByID(c.Request().Context(), p.ID)
	if err != nil {
		log.Printf("error getting category by ID: %v", err)
		return c.JSON(500, map[string]string{"error": "internal server error"})
	}
	return c.JSON(200, category)
}

func (h *CatalogHandler) GetCategoryBySlug(c echo.Context) error {
	type Params struct {
		Slug string `param:"slug" validate:"required"`
	}
	var p Params
	if err := c.Bind(&p); err != nil {
		log.Printf("error binding path params: %v", err)
		return c.JSON(400, map[string]string{"error": "invalid category slug"})
	}
	category, err := h.Svc.GetCategoryBySlug(c.Request().Context(), p.Slug)
	if err != nil {
		log.Printf("error getting category by slug: %v", err)
		return c.JSON(500, map[string]string{"error": "internal server error"})
	}
	return c.JSON(200, category)
}

func (h *CatalogHandler) GetProductsByCategoryID(c echo.Context) error {
	var f catalog.ProductFilter
	if err := c.Bind(&f); err != nil {
		log.Printf("error binding filter params: %v", err)
		return c.JSON(400, map[string]string{"error": "malformed filter parameters"})
	}

	products, err := h.Svc.GetProductsByCategoryID(c.Request().Context(), f)
	if err != nil {
		log.Printf("error getting products by category ID: %v", err)
		return c.JSON(500, map[string]string{"error": "internal server error"})
	}
	return c.JSON(200, products)
}
