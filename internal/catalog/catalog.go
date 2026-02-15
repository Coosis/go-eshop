package catalog

import (
	"context"
	"fmt"
	"strings"

	"github.com/redis/go-redis/v9"
	log "github.com/sirupsen/logrus"

	"github.com/Coosis/go-eshop/internal/comm"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"

	sqlc "github.com/Coosis/go-eshop/sqlc"
)

type CatalogActor struct {
	Pool *pgxpool.Pool
	Client *redis.Client
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
	if !comm.BFExists(
		ctx,
		c.Client,
		comm.BF_products_id,
		id,
	) {
		log.Warnf("product ID %d not found in bloom filter", id)
		return Product{}, fmt.Errorf("product ID %d not found", id)
	}

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
	if !comm.BFExists(
		ctx,
		c.Client,
		comm.BF_products_slug,
		slug,
	) {
		log.Warnf("product slug %s not found in bloom filter", slug)
		return Product{}, fmt.Errorf("product slug %s not found", slug)
	}
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
	if !comm.BFExists(
		ctx,
		c.Client,
		comm.BF_categories_id,
		id,
	) {
		log.Warnf("category ID %d not found in bloom filter", id)
		return Category{}, fmt.Errorf("category not found")
	}

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
	if !comm.BFExists(
		ctx,
		c.Client,
		comm.BF_categories_slug,
		slug,
	) {
		log.Warnf("category slug %s not found in bloom filter", slug)
		return Category{}, fmt.Errorf("category not found")
	}

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
	} else {
		if !comm.BFExists(
			ctx,
			c.Client,
			comm.BF_categories_id,
			*filter.CategoryID,
		) {
			log.Warnf("category ID %d not found in bloom filter", *filter.CategoryID)
			return comm.Page[Product]{}, fmt.Errorf("category not found")
		}
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

	// best effort
	if err := comm.BFReserve(ctx, c.Client, comm.BF_products_id).Err(); err != nil && !strings.Contains(err.Error(), "exist") {
		log.Errorf("error reserving bloom filter for products id: %v", err)
		// keep going, doesn't matter
	}
	if err := comm.BFReserve(ctx, c.Client, comm.BF_products_slug).Err(); err != nil && !strings.Contains(err.Error(), "exist") {
		log.Errorf("error reserving bloom filter for products slug: %v", err)
	}
	if err := comm.BFAdd(ctx, c.Client, comm.BF_products_slug, req.Slug).Err(); err != nil {
		log.Errorf("error adding product slug to bloom filter: %v", err)
	}
	if err := comm.BFAdd(ctx, c.Client, comm.BF_products_id, row.ID).Err(); err != nil {
		log.Errorf("error adding product id to bloom filter: %v", err)
	}

	return product, nil
}

func (c *CatalogActor) UpdateProductByID(
	ctx context.Context,
	req UpdateProductRequest,
) (Product, error) {
	if !comm.BFExists(
		ctx,
		c.Client,
		comm.BF_products_id,
		req.ID,
	) {
		log.Warnf("product ID %d not found in bloom filter", req.ID)
		return Product{}, fmt.Errorf("product not found")
	}

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
	// slug for bf
	if err := comm.BFAdd(ctx, c.Client, comm.BF_products_slug, req.Slug).Err(); err != nil {
		log.Errorf("error adding product slug to bloom filter: %v", err)
	}
	return product, nil
}

func (c *CatalogActor) CreateCategory(
	ctx context.Context,
	req CreateCategoryRequest,
) (Category, error) {
	if req.ParentID != nil && !comm.BFExists(
		ctx,
		c.Client,
		comm.BF_categories_id,
		*req.ParentID,
	) {
		log.Warnf("parent category ID %d not found in bloom filter", *req.ParentID)
		return Category{}, fmt.Errorf("parent category not found")
	}

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
	if err := comm.BFReserve(ctx, c.Client, comm.BF_categories_id).Err(); err != nil && !strings.Contains(err.Error(), "exist") {
		log.Errorf("error reserving bloom filter for categories id: %v", err)
	}
	if err := comm.BFReserve(ctx, c.Client, comm.BF_categories_slug).Err(); err != nil && !strings.Contains(err.Error(), "exist") {
		log.Errorf("error reserving bloom filter for categories slug: %v", err)
	}
	if err := comm.BFAdd(ctx, c.Client, comm.BF_categories_id, row.ID).Err(); err != nil {
		log.Errorf("error adding category id to bloom filter: %v", err)
	}
	if err := comm.BFAdd(ctx, c.Client, comm.BF_categories_slug, req.Slug).Err(); err != nil {
		log.Errorf("error adding category slug to bloom filter: %v", err)
	}
	return category, nil
}

func (c *CatalogActor) UpdateCategoryByID(
	ctx context.Context,
	req UpdateCategoryRequest,
) (Category, error) {
	if !comm.BFExists(
		ctx,
		c.Client,
		comm.BF_categories_id,
		req.ID,
	) {
		log.Warnf("category ID %d not found in bloom filter", req.ID)
		return Category{}, fmt.Errorf("category not found")
	}

	if req.ParentID != nil && !comm.BFExists(
		ctx,
		c.Client,
		comm.BF_categories_id,
		*req.ParentID,
	) {
		log.Warnf("parent category ID %d not found in bloom filter", *req.ParentID)
		return Category{}, fmt.Errorf("parent category not found")
	}

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
	// slug for bf
	if err := comm.BFAdd(ctx, c.Client, comm.BF_categories_slug, req.Slug).Err(); err != nil {
		log.Errorf("error adding category slug to bloom filter: %v", err)
	}
	return category, nil
}
