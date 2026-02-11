package catalog

import (
	"context"
	"fmt"

	log "github.com/sirupsen/logrus"

	"github.com/Coosis/go-eshop/internal/comm"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"

	sqlc "github.com/Coosis/go-eshop/sqlc"
)

type CatalogActor struct {
	Pool *pgxpool.Pool
}

// products
func (c *CatalogActor) GetProducts(
	ctx context.Context,
	filter ProductFilter,
) (comm.Page[Product], error) {
	page := comm.Page[Product]{}
	page.Page = filter.Page
	page.PerPage = filter.PerPage
	queries := sqlc.New(c.Pool)
	pg_minPrice := pgtype.Int4{Valid: false}
	if filter.MinPrice != nil {
		pg_minPrice = pgtype.Int4{Int32: *filter.MinPrice, Valid: true}
	}
	pg_maxPrice := pgtype.Int4{Valid: false}
	if filter.MaxPrice != nil {
		pg_maxPrice = pgtype.Int4{Int32: *filter.MaxPrice, Valid: true}
	}
	pg_categoryID := pgtype.Int4{Valid: false}
	if filter.CategoryID != nil {
		pg_categoryID = pgtype.Int4{Int32: *filter.CategoryID, Valid: true}
	}
	rows, err := queries.GetProducts(ctx, sqlc.GetProductsParams{
		MinPriceCents: pg_minPrice,
		MaxPriceCents: pg_maxPrice,
		CategoryID:    pg_categoryID,
		PageNumber:    filter.Page,
		PageSize:      filter.PerPage,
	})
	if err != nil {
		log.Errorf("error getting products: %v", err)
		return page, err
	}

	page.Total = int64(len(rows))
	var items []Product
	for _, row := range rows {
		product := Product{
			ID: row.ID,
			ProductProperties: ProductProperties{
				Name:        row.Name,
				Slug:        row.Slug,
				Description: &row.Description.String,
				PriceCents:  row.PriceCents,
				CategoryIDs: row.CategoryIds,
			},
			PriceVersion: row.PriceVersion,
		}
		items = append(items, product)
	}

	page.Items = items
	log.Infof("GetProducts called with filter: %+v", filter)
	return page, nil
}

func (c *CatalogActor) GetProductByID(
	ctx context.Context,
	id int32,
) (Product, error) {
	queries := sqlc.New(c.Pool)
	row, err := queries.GetProductByID(ctx, id)
	log.Infof("GetProductByID called with id: %d", id)
	if err != nil {
		log.Errorf("error getting product by id: %v", err)
		return Product{}, err
	}

	product := Product{
		ID: row.ID,
		ProductProperties: ProductProperties{
			Name:        row.Name,
			Slug:        row.Slug,
			Description: &row.Description.String,
			PriceCents:  row.PriceCents,
			CategoryIDs: row.CategoryIds,
		},
		PriceVersion: row.PriceVersion,
	}

	return product, nil
}

func (c *CatalogActor) GetProductBySlug(
	ctx context.Context,
	slug string,
) (Product, error) {
	quieries := sqlc.New(c.Pool)
	row, err := quieries.GetProductBySlug(ctx, slug)
	log.Infof("GetProductBySlug called with slug: %s", slug)
	if err != nil {
		log.Errorf("error getting product by slug: %v", err)
		return Product{}, err
	}

	product := Product{
		ID: row.ID,
		ProductProperties: ProductProperties{
			Name:        row.Name,
			Slug:        row.Slug,
			Description: &row.Description.String,
			PriceCents:  row.PriceCents,
			CategoryIDs: row.CategoryIds,
		},
		PriceVersion: row.PriceVersion,
	}

	return product, nil
}

// categories
func (c *CatalogActor) GetCategories(
	ctx context.Context,
	req GetCategoriesRequest,
) (comm.Page[Category], error) {
	page := comm.Page[Category]{}
	page.Page = req.Page
	page.PerPage = req.PerPage
	queries := sqlc.New(c.Pool)
	rows, err := queries.GetCategories(ctx, sqlc.GetCategoriesParams{
		PageNumber: req.Page,
		PageSize:   req.PerPage,
	})
	log.Infof("GetCategories called with req: %+v", req)
	if err != nil {
		log.Errorf("error getting categories: %v", err)
		return page, err
	}
	page.Total = int64(len(rows))
	var items []Category
	for _, row := range rows {
		category := Category{
			ID: row.ID,
			CategoryProperties: CategoryProperties{
				Name:     row.Name,
				Slug:     row.Slug,
				ParentID: &row.ParentID.Int32,
			},
		}
		items = append(items, category)
	}
	page.Items = items
	return page, nil
}

func (c *CatalogActor) GetCategoryByID(
	ctx context.Context,
	id int32,
) (Category, error) {
	queries := sqlc.New(c.Pool)
	row, err := queries.GetCategoryByID(ctx, id)
	log.Infof("GetCategoryByID called with id: %d", id)
	if err != nil {
		log.Errorf("error getting category by id: %v", err)
		return Category{}, err
	}

	category := Category{
		ID: row.ID,
		CategoryProperties: CategoryProperties{
			Name:     row.Name,
			Slug:     row.Slug,
			ParentID: &row.ParentID.Int32,
		},
	}
	return category, nil
}

func (c *CatalogActor) GetCategoryBySlug(
	ctx context.Context,
	slug string,
) (Category, error) {
	queries := sqlc.New(c.Pool)
	row, err := queries.GetCategoryBySlug(ctx, slug)
	log.Infof("GetCategoryByID called with slug: %v", slug)
	if err != nil {
		log.Errorf("error getting category by slug: %v", err)
		return Category{}, err
	}

	category := Category{
		ID: row.ID,
		CategoryProperties: CategoryProperties{
			Name:     row.Name,
			Slug:     row.Slug,
			ParentID: &row.ParentID.Int32,
		},
	}
	return category, nil
}

func (c *CatalogActor) GetProductsByCategoryID(
	ctx context.Context,
	filter ProductFilter,
) (comm.Page[Product], error) {
	if filter.CategoryID == nil {
		return comm.Page[Product]{}, fmt.Errorf("CategoryID is required in filter")
	}
	page := comm.Page[Product]{}
	page.Page = filter.Page
	page.PerPage = filter.PerPage
	queries := sqlc.New(c.Pool)
	prods, err := queries.GetProductByCategoryID(ctx, sqlc.GetProductByCategoryIDParams{
		CategoryID: *filter.CategoryID,
		PageNumber: filter.Page,
		PageSize:   filter.PerPage,
	})
	log.Infof("GetProductsByCategoryID called with filter: %+v", filter)
	if err != nil {
		log.Errorf("error getting products by category id: %v", err)
		return page, err
	}

	for _, row := range prods {
		product := Product{
			ID: row.ID,
			ProductProperties: ProductProperties{
				Name:        row.Name,
				Slug:        row.Slug,
				Description: &row.Description.String,
				PriceCents:  row.PriceCents,
				CategoryIDs: row.CategoryIds,
			},
			PriceVersion: row.PriceVersion,
		}
		page.Items = append(page.Items, product)
	}

	page.Total = int64(len(prods))
	return page, nil
}

// admin-only
func (c *CatalogActor) CreateProduct(
	ctx context.Context,
	req CreateProductRequest,
) (Product, error) {
	var desc pgtype.Text
	if req.Description != nil {
		desc = pgtype.Text{String: *req.Description, Valid: true}
	} else {
		desc = pgtype.Text{Valid: false}
	}
	queries := sqlc.New(c.Pool)
	row, err := queries.CreateProduct(ctx, sqlc.CreateProductParams{
		Name:        req.Name,
		Slug:        req.Slug,
		Description: desc,
		PriceCents:  req.PriceCents,
	})
	log.Infof("CreateProduct called with req: %+v", req)
	if err != nil {
		log.Errorf("error creating product: %v", err)
		return Product{}, err
	}

	product := Product{
		ID: row.ID,
		ProductProperties: ProductProperties{
			Name:        row.Name,
			Slug:        row.Slug,
			Description: &row.Description.String,
			PriceCents:  row.PriceCents,
			CategoryIDs: row.CategoryIds,
		},
		PriceVersion: row.PriceVersion,
	}
	return product, nil
}

func (c *CatalogActor) UpdateProductByID(
	ctx context.Context,
	req UpdateProductRequest,
) (Product, error) {
	var desc pgtype.Text
	if req.Description != nil {
		desc = pgtype.Text{String: *req.Description, Valid: true}
	} else {
		desc = pgtype.Text{Valid: false}
	}

	queries := sqlc.New(c.Pool)
	row, err := queries.UpdateProductByID(ctx, sqlc.UpdateProductByIDParams{
		ID:          req.ID,
		Name:        req.Name,
		Slug:        req.Slug,
		Description: desc,
		PriceCents:  req.PriceCents,
	})
	log.Infof("UpdateProductByID called with req: %+v", req)
	if err != nil {
		log.Errorf("error updating product by id: %v", err)
		return Product{}, err
	}

	product := Product{
		ID: row.ID,
		ProductProperties: ProductProperties{
			Name:        row.Name,
			Slug:        row.Slug,
			Description: &row.Description.String,
			PriceCents:  row.PriceCents,
			CategoryIDs: row.CategoryIds,
		},
		PriceVersion: row.PriceVersion,
	}
	return product, nil
}

func (c *CatalogActor) CreateCategory(
	ctx context.Context,
	req CreateCategoryRequest,
) (Category, error) {
	var parentID pgtype.Int4
	if req.ParentID != nil {
		parentID = pgtype.Int4{Int32: *req.ParentID, Valid: true}
	} else {
		parentID = pgtype.Int4{Valid: false}
	}

	queries := sqlc.New(c.Pool)
	row, err := queries.CreateCategory(ctx, sqlc.CreateCategoryParams{
		Name:     req.Name,
		Slug:     req.Slug,
		ParentID: parentID,
	})
	log.Infof("CreateCategory called with req: %+v", req)
	if err != nil {
		log.Errorf("error creating category: %v", err)
		return Category{}, err
	}

	category := Category{
		ID: row.ID,
		CategoryProperties: CategoryProperties{
			Name:     row.Name,
			Slug:     row.Slug,
			ParentID: &row.ParentID.Int32,
		},
	}
	return category, nil
}

func (c *CatalogActor) UpdateCategoryByID(
	ctx context.Context,
	req UpdateCategoryRequest,
) (Category, error) {
	var parentID pgtype.Int4
	if req.ParentID != nil {
		parentID = pgtype.Int4{Int32: *req.ParentID, Valid: true}
	} else {
		parentID = pgtype.Int4{Valid: false}
	}

	queries := sqlc.New(c.Pool)
	row, err := queries.UpdateCategoryByID(ctx, sqlc.UpdateCategoryByIDParams{
		ID:       req.ID,
		Name:     req.Name,
		Slug:     req.Slug,
		ParentID: parentID,
	})
	log.Infof("UpdateCategoryByID called with req: %+v", req)
	if err != nil {
		log.Errorf("error updating category by id: %v", err)
		return Category{}, err
	}

	category := Category{
		ID: row.ID,
		CategoryProperties: CategoryProperties{
			Name:     row.Name,
			Slug:     row.Slug,
			ParentID: &row.ParentID.Int32,
		},
	}
	return category, nil
}
