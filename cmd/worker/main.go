package main

import (
	"context"
	"os"
	"os/signal"
	"strconv"
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
	RunE:  runWorker,
}

func runWorker(cmd *cobra.Command, args []string) error {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	worker, err := NewWorker(client, pool, workerID)
	if err != nil {
		log.Errorf("Failed to create worker: %v", err)
		return err
	}

	done := make(chan struct{})
	if orderWorkerEnabled {
		log.Info("Starting order worker...")
		worker.createOrder(ctx)
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

var (
	ENV_redisURL = "REDIS_URL"
	redisURL     = "redis://localhost:6380/0"

	ENV_dbURL = "DB_URL"
	dbURL     = "postgres://postgres:passwd@localhost:5433/postgres"

	ENV_workerID       = "WORKER_ID"
	workerID     int64 = 2
)

func main() {
	if envWorkerID := os.Getenv(ENV_workerID); envWorkerID != "" {
		i, err := strconv.ParseInt(envWorkerID, 10, 64)
		if err != nil {
			log.Warnf("Invalid worker ID in environment variable, should be 64 bit integer. Using default: %d", workerID)
		} else {
			workerID = i
			log.Infof("Using worker ID from environment variable: %d", workerID)
		}
	} else {
		log.Warnf("Environment variable for worker ID not set, using default: %d", workerID)
	}

	if envRedisUrl := os.Getenv(ENV_redisURL); envRedisUrl != "" {
		redisURL = envRedisUrl
		log.Infof("Using Redis URL from environment variable: %s", redisURL)
	} else {
		log.Warnf("Environment variable for Redis URL not set, using default: %s", redisURL)
	}

	if envDbUrl := os.Getenv(ENV_dbURL); envDbUrl != "" {
		dbURL = envDbUrl
		log.Infof("Using DB URL from environment variable: %s", dbURL)
	} else {
		log.Warnf("Environment variable for DB URL not set, using default: %s", dbURL)
	}
	opt, err := redis.ParseURL(redisURL)
	if err != nil {
		log.Fatalf("parsing redis url failed: %v", err)
	}
	client = redis.NewClient(opt)

	cfg, err := pgxpool.ParseConfig(dbURL)
	if err != nil {
		log.Fatalf("Failed to parse DB config: %v", err)
		return
	}
	attemptCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	pool, err = pgxpool.NewWithConfig(attemptCtx, cfg)
	cancel()
	if err != nil {
		log.Fatalf("Failed to create DB pool: %v", err)
		return
	}

	if pingErr := pool.Ping(context.Background()); pingErr != nil {
		pool.Close()
		log.Fatalf("Failed to ping DB: %v", pingErr)
		return
	}

	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
