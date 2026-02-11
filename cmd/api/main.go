package main

import (
	"context"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/Coosis/go-eshop/internal/auth"
	"github.com/Coosis/go-eshop/internal/cart"
	"github.com/Coosis/go-eshop/internal/catalog"
	"github.com/Coosis/go-eshop/internal/handlers"
	"github.com/Coosis/go-eshop/internal/orders"
	"github.com/Coosis/go-eshop/internal/seckill"
	"github.com/Coosis/go-eshop/internal/stock"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/redis/go-redis/v9"
	validator "github.com/go-playground/validator/v10"
)

func init() {
	log.SetReportCaller(true)
}

func main() {
	opt, err := redis.ParseURL("redis://localhost:6380/0")
	if err != nil {
		panic(err)
	}
	client := redis.NewClient(opt)

	url := "postgres://postgres:passwd@localhost:5433/postgres"
	cfg, err := pgxpool.ParseConfig(url)
	if err != nil {
		log.Errorf("Failed to parse DB config: %v", err)
		return
	}

	attemptCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	pool, err := pgxpool.NewWithConfig(attemptCtx, cfg)
	cancel()
	if err != nil {
		log.Errorf("Failed to create DB pool: %v", err)
		return
	}

	if pingErr := pool.Ping(context.Background()); pingErr != nil {
		pool.Close()
		log.Errorf("Failed to ping DB: %v", pingErr)
		return
	}

	catalogActor := catalog.CatalogActor{Pool: pool}
	cartActor := cart.CartActor{Pool: pool}
	orderActor := orders.OrderActor{Pool: pool}
	stockActor := stock.StockActor{Pool: pool}
	seckillActor := seckill.SeckillActor{Pool: pool, Client: client}

	e := echo.New()
	e.Validator = &CustomValidator{validator: validator.New()}
	e.Use(middleware.Recover())
	e.Use(auth.DevWithUserID(1))
	handlers.RegisterCatalogProductRoutes(e, &catalogActor)
	handlers.RegisterCatalogAdminRoutes(e, &catalogActor)
	handlers.RegisterCatalogCategoryRoutes(e, &catalogActor)
	handlers.RegisterCartRoutes(e, &cartActor)
	handlers.RegisterOrderRoutes(e, &orderActor)
	handlers.RegisterStockRoutes(e, &stockActor)
	handlers.RegisterSeckillRoutes(e, &seckillActor)
	if err := e.Start(":8144"); err != nil {
		panic(err)
	}
}
