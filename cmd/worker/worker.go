package main

import (
	"fmt"
	sqlc "github.com/Coosis/go-eshop/sqlc"
	log "github.com/sirupsen/logrus"
	"github.com/redis/go-redis/v9"
	"github.com/bwmarrin/snowflake"
)

type Worker struct {
	client *redis.Client
	db sqlc.DBTX
	workerID int64
	node *snowflake.Node
}

func NewWorker(client *redis.Client, db sqlc.DBTX, workerID int64) (*Worker, error) {
	if workerID < 0 || workerID > 1023 {
		return nil, fmt.Errorf("worker ID must be between 0 and 1023")
	}
	node, err := snowflake.NewNode(workerID)
	if err != nil {
		log.Errorf("Failed to create Snowflake node: %v", err)
		return nil, fmt.Errorf("failed to create Snowflake node: %w", err)
	}
	return &Worker{
		client: client,
		db: db,
		workerID: workerID,
		node: node,
	}, nil
}

func (w *Worker) GenerateOrderNumber() string {
	return w.node.Generate().String()
}
