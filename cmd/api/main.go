package main

import (
	"context"
	// "fmt"
	"strings"
	"sync/atomic"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/Coosis/go-eshop/internal/auth"
	"github.com/Coosis/go-eshop/internal/cart"
	"github.com/Coosis/go-eshop/internal/catalog"
	"github.com/Coosis/go-eshop/internal/comm"
	"github.com/Coosis/go-eshop/internal/handlers"
	"github.com/Coosis/go-eshop/internal/orders"
	"github.com/Coosis/go-eshop/internal/seckill"
	"github.com/Coosis/go-eshop/internal/stock"
	validator "github.com/go-playground/validator/v10"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/redis/go-redis/v9"
)

var (
	ServiceWarm atomic.Bool
)

func init() {
	log.SetReportCaller(true)
}

func main() {
	ServiceWarm.Store(false)

	opt, err := redis.ParseURL("redis://127.0.0.1:6380/0")
	if err != nil {
		panic(err)
	}
	client := redis.NewClient(opt)

	if err := client.Ping(context.Background()).Err(); err != nil {
		log.Errorf("Failed to connect to Redis: %v", err)
		return
	}

	if err := comm.BFReserve(
		context.Background(), 
		client,
		comm.BF_seckill_events,
	).Err(); err != nil {
		if !strings.Contains(err.Error(), "exist") {
			log.Errorf("Failed to create Bloom filter: %v", err)
			return
		}
	}

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

	go func() {
		if err := warmup(context.Background(), pool, client); err != nil {
			return
		}
		ServiceWarm.Store(true)
		log.Info("Service warmup completed, now accepting requests.")
	}()

	catalogActor := catalog.CatalogActor{Pool: pool, Client: client}
	cartActor := cart.CartActor{Pool: pool}
	orderActor := orders.OrderActor{Pool: pool}
	stockActor := stock.StockActor{Pool: pool}
	seckillActor := seckill.SeckillActor{Pool: pool, Client: client}

	e := echo.New()
	e.GET("/health", func(c echo.Context) error {
		healthy := ServiceWarm.Load()
		if healthy {
			return c.String(200, "OK")
		}
		return c.String(503, "Service is warming up, please try again later.")
	})
	e.Validator = &CustomValidator{validator: validator.New()}
	e.Use(middleware.RecoverWithConfig(
		middleware.RecoverConfig{
			LogErrorFunc: func(c echo.Context, err error, stack []byte) error {
				log.Errorf("Panic recovered: %v\nStack trace:\n%s", err, string(stack))
				return err
			},
		},
	))
	e.Use(auth.DevWithUserID(1))
	e.GET("/panic", func(c echo.Context) error {
		panic("Intentional panic for testing recovery middleware")
	})
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
