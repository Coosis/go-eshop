package main

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"sync"
	"time"

	sqlc "github.com/Coosis/go-eshop/sqlc"
	"github.com/jackc/pgx/v5/pgtype"
	log "github.com/sirupsen/logrus"

	"github.com/redis/go-redis/v9"
)

const (
	seckillOrderStream = "seckill_order_stream"
	groupName          = "workers"
	stealCnt           = 10
)

func (w *Worker) createOrder(
	ctx context.Context,
) error {
	res, err := client.XGroupCreateMkStream(
		ctx,
		seckillOrderStream,
		groupName,
		"0",
	).Result()

	if err != nil {
		if redis.HasErrorPrefix(err, "BUSYGROUP") {
			log.Warnf("Consumer group already exists: %s, but it's probably ok", groupName)
		} else {
			log.Errorf("Failed to create consumer group: %v", err)
			return err
		}
	}

	var worker_identifier string = fmt.Sprintf("worker-%d", w.workerID)

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()

		ticker := time.NewTicker(time.Second * 5)
		defer ticker.Stop()
		scan_head := "0-0"
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				msgs, nxt, err := client.XAutoClaim(ctx, &redis.XAutoClaimArgs{
					Stream:   seckillOrderStream,
					Group:    groupName,
					MinIdle:  time.Second * 30,
					Start:    scan_head,
					Count:    stealCnt,
					Consumer: worker_identifier,
				}).Result()
				if err != nil {
					log.Errorf("Failed to auto claim messages: %v", err)
					continue
				}
				scan_head = nxt

				for _, msg := range msgs {
					log.Infof("%s processing message: %v", worker_identifier, msg)

					err := w.handleCreateOrder(ctx, msg)
					if err != nil {
						log.Errorf("Failed to handle message: %v", err)
						continue
					}

					client.XAck(ctx, seckillOrderStream, groupName, msg.ID)
				}
			}
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			select {
			case <-ctx.Done():
				return
			default:
			}
			msgs, err := client.XReadGroup(ctx, &redis.XReadGroupArgs{
				Group:    groupName,
				Consumer: worker_identifier,
				Streams:  []string{seckillOrderStream, ">"},
				Count:    10,
				Block:    time.Second * 3,
				NoAck:    false,
			}).Result()
			if err != nil {
				if err == redis.Nil {
				} else {
					log.Errorf("Failed to read messages: %v", err)
				}
				continue
			}
			for _, m := range msgs {
				for _, msg := range m.Messages {
					log.Infof("%s processing message: %v", worker_identifier, msg)

					err := w.handleCreateOrder(ctx, msg)
					if err != nil {
						log.Errorf("Failed to handle message: %v", err)
						continue
					} else {
						log.Infof("Successfully handled message ID=%s", msg.ID)
					}

					acked, err := client.XAck(ctx, seckillOrderStream, groupName, msg.ID).Result()
					if err != nil {
						log.Errorf("Failed to acknowledge message: %v", err)
					} else {
						log.Infof("Acknowledged message ID=%s, acked=%d", msg.ID, acked)
					}
				}
			}
		}
	}()

	log.Infof("Consumer group created: %s", res)
	wg.Wait()
	return nil
}

func (w *Worker) handleCreateOrder(
	ctx context.Context,
	msg redis.XMessage,
) error {
	event_id := msg.Values["event_id"]
	quantity := msg.Values["quantity"].(string)
	idempotency_key := msg.Values["idempotency_key"].(string)
	price := msg.Values["price_cents"].(string)
	log.Infof(
		"worker-%d handling order: EventID=%s, Quantity=%s, IdempotencyKey=%s",
		w.workerID,
		event_id,
		quantity,
		idempotency_key,
	)

	userID := strings.SplitN(idempotency_key, ":", 2)[0]
	userIDInt, err := strconv.ParseInt(userID, 10, 32)
	if err != nil {
		log.Errorf("Invalid user ID in idempotency key: %s", userID)
		return err
	}

	priceInt, err := strconv.ParseInt(price, 10, 32)
	if err != nil {
		log.Errorf("Invalid price in message: %s", price)
		return err
	}
	qtyInt, err := strconv.ParseInt(quantity, 10, 64)
	if err != nil {
		log.Errorf("Invalid quantity in message: %s", quantity)
		return err
	}
	queries := sqlc.New(w.db)
	if _, err := queries.CreateSeckillOrder(ctx, sqlc.CreateSeckillOrderParams{
		UserID:         int32(userIDInt),
		OrderNumber:    w.GenerateOrderNumber(),
		SubtotalCents:  priceInt * qtyInt,
		Notes:          pgtype.Text{Valid: false},
		IdempotencyKey: idempotency_key,
	}); err != nil {
		log.Errorf("Failed to create order in database: %v", err)
		return err
	}
	return nil
}
