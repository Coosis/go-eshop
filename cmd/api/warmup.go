package main

import (
	"context"
	"strings"
	"time"

	comm "github.com/Coosis/go-eshop/internal/comm"
	sqlc "github.com/Coosis/go-eshop/sqlc"
	"github.com/redis/go-redis/v9"
	"golang.org/x/sync/errgroup"
	log "github.com/sirupsen/logrus"
)

const (
	CHUNK_SIZE = 20
)

func warmup(
	ctx context.Context,
	db sqlc.DBTX,
	client *redis.Client,
) error {
	queries := sqlc.New(db)
	var wg errgroup.Group
	wg.Go(func() error {
		if err := comm.BFReserve(
			ctx,
			client,
			comm.BF_products_id,
		).Err(); err != nil {
			if !strings.Contains(err.Error(), "exist") {
				log.Errorf("Failed to create Bloom filter: %v", err)
				return err
			}
		}
		if err := comm.BFReserve(
			ctx,
			client,
			comm.BF_products_slug,
		).Err(); err != nil {
			if !strings.Contains(err.Error(), "exist") {
				log.Errorf("Failed to create Bloom filter: %v", err)
				return err
			}
		}
		var lastID int32
		lastID = 0
		for {
			rows, err := queries.GetProductsChunk(
				ctx,
				sqlc.GetProductsChunkParams{
					Limit: 100,
					ID: lastID,
				},
			)
			if err != nil {
				log.Errorf("Failed to get products chunk: %v", err)
				time.Sleep(time.Second * 5)
				continue
			}

			if len(rows) == 0 {
				break
			}

			ids := make([]any, 0, len(rows))
			slugs := make([]any, 0, len(rows))
			for _, row := range rows {
				ids = append(ids, row.ID)
				slugs = append(slugs, row.Slug)
				lastID = max(lastID, row.ID)
			}
			_, err = comm.BFMAdd(
				ctx,
				client,
				comm.BF_products_id,
				ids...,
			).Result()
			if err != nil {
				log.Errorf("Failed to add product IDs to Bloom filter: %v", err)
				log.Errorf("IDs: %v", ids)
				log.Errorf("Best effort anyway, continue...")
			}
			_, err = comm.BFMAdd(
				ctx,
				client,
				comm.BF_products_slug,
				slugs...,
			).Result()
			if err != nil {
				log.Errorf("Failed to add product slugs to Bloom filter: %v", err)
				log.Errorf("Slugs: %v", slugs)
				log.Errorf("Best effort anyway, continue...")
			}
			log.Infof("Cached products up to ID %d", lastID)
		}
		return nil
	})
	wg.Go(func() error {
		if err := comm.BFReserve(
			ctx,
			client,
			comm.BF_categories_id,
		).Err(); err != nil {
			if !strings.Contains(err.Error(), "exist") {
				log.Errorf("Failed to create Bloom filter: %v", err)
				return err
			}
		}
		if err := comm.BFReserve(
			ctx,
			client,
			comm.BF_categories_slug,
		).Err(); err != nil {
			if !strings.Contains(err.Error(), "exist") {
				log.Errorf("Failed to create Bloom filter: %v", err)
				return err
			}
		}
		var lastID int32
		lastID = 0
		for {
			rows, err := queries.GetCategoriesChunk(
				ctx,
				sqlc.GetCategoriesChunkParams{
					Limit: 100,
					ID: lastID,
				},
			)
			if err != nil {
				log.Errorf("Failed to get categories chunk: %v", err)
				time.Sleep(time.Second * 5)
				continue
			}

			if len(rows) == 0 {
				break
			}

			ids := make([]any, 0, len(rows))
			slugs := make([]any, 0, len(rows))
			for _, row := range rows {
				ids = append(ids, row.ID)
				slugs = append(slugs, row.Slug)
				lastID = max(lastID, row.ID)
			}
			_, err = comm.BFMAdd(
				ctx,
				client,
				comm.BF_categories_id,
				ids...,
			).Result()
			if err != nil {
				log.Errorf("Failed to add category IDs to Bloom filter: %v", err)
				log.Errorf("IDs: %v", ids)
				log.Errorf("Best effort anyway, continue...")
			}
			_, err = comm.BFMAdd(
				ctx,
				client,
				comm.BF_categories_slug,
				slugs...,
			).Result()
			if err != nil {
				log.Errorf("Failed to add category slugs to Bloom filter: %v", err)
				log.Errorf("Slugs: %v", slugs)
				log.Errorf("Best effort anyway, continue...")
			}
			log.Infof("Cached categories up to ID %d", lastID)
		}
		return nil
	})
	err := wg.Wait()
	if err != nil {
		log.Errorf("Warmup failed: %v", err)
		return err
	}
	log.Infof("Service warmup completed")
	return nil
}
