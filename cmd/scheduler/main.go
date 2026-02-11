package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"time"

	sqlc "github.com/Coosis/go-eshop/sqlc"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

const (
	redisKeyEventStockFmt      = "seckill_event:%d:stock"
	redisKeyEventPriceFmt      = "seckill_event:%d:price"
	redisKeyEventNotStartedFmt = "seckill_event:%d:not_started"
	redisKeyEventLenMsFmt      = "seckill_event:%d:event_len_ms"
)

var pool *pgxpool.Pool
var client *redis.Client

var rootCmd = &cobra.Command{
	Use:   "seckill-scheduler",
	Short: "A scheduler for seckill events",
	Long:  `This service is responsible for scheduling seckill events and managing their lifecycle.`,
}

func init() {
	rootCmd.AddCommand(&cobra.Command{
		Use:   "start",
		Short: "Start the seckill scheduler",
		RunE:  runScheduler,
	})
}

func main() {
	opt, err := redis.ParseURL("redis://localhost:6380/0")
	if err != nil {
		panic(err)
	}
	client = redis.NewClient(opt)

	url := "postgres://postgres:passwd@localhost:5433/postgres"
	cfg, err := pgxpool.ParseConfig(url)
	if err != nil {
		log.Errorf("Failed to parse DB config: %v", err)
		return
	}
	attemptCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	pool, err = pgxpool.NewWithConfig(attemptCtx, cfg)
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

	if err := rootCmd.Execute(); err != nil {
		log.Fatalf("Error executing command: %v", err)
	}
}

func runScheduler(cmd *cobra.Command, args []string) error {
	log.Infof("Starting the scheduler for seckill events...")

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()
	var wg sync.WaitGroup
	for {
		select {
		case <-ctx.Done():
			log.Infof("Received interrupt signal, shutting down scheduler...")
			ticker.Stop()
			sdctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			done := make(chan struct{})
			go func() {
				wg.Wait()
				close(done)
			}()
			select {
			case <-done:
				log.Infof("All tasks completed, exiting.")
			case <-sdctx.Done():
				log.Warnf("Timeout reached, exiting with pending tasks.")
			}
			return nil
		case <-ticker.C:
			log.Infof("Checking for upcoming seckill events...")
			wg.Add(1)
			go func() {
				defer wg.Done()
				schedulerCheck()
			}()
		}
	}
}

func schedulerCheck() {
	log.Infof("Querying database for upcoming seckill events...")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	tx, err := pool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		log.Errorf("Failed to begin transaction: %v", err)
		return
	}
	defer tx.Rollback(ctx)
	queries := sqlc.New(pool).WithTx(tx)
	rows, err := queries.GetDueNotPreheated(ctx, 10) // only take 10
	if err != nil {
		log.Errorf("Failed to query upcoming seckill events: %v", err)
		return
	}
	for _, r := range rows {
		log.Infof(`Found upcoming seckill event: ID=%d,  StartTime=%s, 
			EndTime=%s, ProductID=%d, PriceCents=%d, Stock=%d`,
			r.ID,
			r.StartTime.Time.Format(time.RFC3339),
			r.EndTime.Time.Format(time.RFC3339),
			r.ProductID,

			r.SeckillPriceCents,
			r.SeckillStock,
		)

		log.Infof("Preheating event ID=%d...", r.ID)

		// redis preheat
		lua := `
			local stock_key = KEYS[1]
			local price_key = KEYS[2]
			local not_started_key = KEYS[3]
			local event_len_ms_key = KEYS[4]

			local stock = ARGV[1]
			local price = ARGV[2]

			local until_start_ms = tonumber(ARGV[3])
			local until_end_ms = tonumber(ARGV[4])

			local event_len_ms = until_end_ms - until_start_ms

			redis.call("SET", event_len_ms_key, event_len_ms, "NX", "PX", until_end_ms)
			redis.call("SET", stock_key, stock, "NX", "PX", until_end_ms)
			redis.call("SET", price_key, price, "NX", "PX", until_end_ms)
			redis.call("SET", not_started_key, 1, "NX", "PX", until_start_ms)
			return 1
		`

		if ret, err := client.Eval(
			ctx,
			lua,
			[]string{
				fmt.Sprintf(redisKeyEventStockFmt, r.ID),
				fmt.Sprintf(redisKeyEventPriceFmt, r.ID),
				fmt.Sprintf(redisKeyEventNotStartedFmt, r.ID),
				fmt.Sprintf(redisKeyEventLenMsFmt, r.ID),
			},
			r.SeckillStock,
			r.SeckillPriceCents,
			time.Until(r.StartTime.Time).Milliseconds(),
			time.Until(r.EndTime.Time).Milliseconds(),
		).Result(); err != nil || ret.(int64) != 1 {
			log.Errorf("Failed to preheat Redis for event ID=%d: %v", r.ID, err)
			continue
		}

		if _, err := queries.MarkPreheated(ctx, r.ID); err != nil {
			log.Errorf("Failed to mark event ID=%d as preheated: %v", r.ID, err)
			return
		}
	}
	if err := tx.Commit(ctx); err != nil {
		log.Errorf("Failed to commit transaction: %v", err)
		return
	}
}
