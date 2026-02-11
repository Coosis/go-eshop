package main

import (
	"context"
	"os"
	"os/signal"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var (
	orderWorkerEnabled bool
)

var pool *pgxpool.Pool
var client *redis.Client

func init() {
	rootCmd.Flags().BoolVarP(
		&orderWorkerEnabled,
		"order",
		"o",
		true,
		"whether to run the order worker",
	)
}

var rootCmd = &cobra.Command{
	Use:   "eshop-worker",
	Short: "A simple e-commerce application",
	Long:  `This is a simple e-commerce application built with Go, Echo, PostgreSQL, and Redis.`,
	RunE: runWorker,
}

func runWorker(cmd *cobra.Command, args []string) error {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	done := make(chan struct{})
	if orderWorkerEnabled {
		log.Info("Starting order worker...")
		createOrder(
			ctx,
			client,
			pool,
			"worker-1",
		)
		close(done)
	}

	<-ctx.Done()
	log.Info("Shutting down worker...")

	select {
	case <-done:
		log.Info("Worker stopped gracefully.")
	case <-time.After(5 * time.Second):
		log.Warn("Worker shutdown timed out, forcing exit.")
	}
	return nil
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
		panic(err)
	}
}
