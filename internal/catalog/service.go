package catalog

import (
	"context"
	"errors"

	"github.com/Coosis/go-eshop/internal/comm"
)

var ErrNotFound = errors.New("not found")

// product related structs:
// so other structs can share
type ProductProperties struct {
	Name        string  `json:"name"`
	Slug        string  `json:"slug"`
	Description *string `json:"description,omitempty"`
	PriceCents  int32   `json:"price_cents"`
	CategoryIDs []int32 `json:"category_ids,omitempty"`
}

type Product struct {
	ID int32 `json:"id"`
	ProductProperties
	PriceVersion int64 `json:"price_version"`
}
type CreateProductRequest struct {
	ProductProperties
}

type UpdateProductRequest struct {
	ID int32 `json:"id"`
	ProductProperties
}

// category related structs:
type CategoryProperties struct {
	Name     string `json:"name"`
	Slug     string `json:"slug"`
	ParentID *int32 `json:"parent_id,omitempty"`
}

type Category struct {
	ID int32 `json:"id"`
	CategoryProperties
}

type GetCategoriesRequest struct {
	PerPage int32 `query:"per_page" validate:"omitempty,min=1,max=250"`
	Page    int32 `query:"page" validate:"omitempty,min=1"`
}

type CreateCategoryRequest struct {
	CategoryProperties
}

type UpdateCategoryRequest struct {
	ID int32 `json:"id"`
	CategoryProperties
}

type CatalogService interface {
	// products
	GetProducts(ctx context.Context, filter ProductFilter) (comm.Page[Product], error)
	GetProductByID(ctx context.Context, id int32) (Product, error)
	GetProductBySlug(ctx context.Context, slug string) (Product, error)

	// categories
	GetCategories(ctx context.Context, req GetCategoriesRequest) (comm.Page[Category], error)
	GetCategoryByID(ctx context.Context, id int32) (Category, error)
	GetCategoryBySlug(ctx context.Context, slug string) (Category, error)
	GetProductsByCategoryID(
		ctx context.Context,
		filter ProductFilter,
	) (comm.Page[Product], error)

	// admin-only
	CreateProduct(ctx context.Context, req CreateProductRequest) (Product, error)
	UpdateProductByID(ctx context.Context, req UpdateProductRequest) (Product, error)
	CreateCategory(ctx context.Context, req CreateCategoryRequest) (Category, error)
	UpdateCategoryByID(ctx context.Context, req UpdateCategoryRequest) (Category, error)
}

type ProductFilter struct {
	CategoryID *int32 `param:"id" validate:"min=1,omitempty"`
	PerPage    int32  `query:"per_page" validate:"required,min=1,max=250"`
	Page       int32  `query:"page" validate:"required,min=1"`
	Q          string `query:"q" validate:"omitempty"`
	MinPrice   *int32 `query:"min_price" validate:"omitempty,min=0"`
	MaxPrice   *int32 `query:"max_price" validate:"omitempty,min=0"`
}
