package comm

import (
	log "github.com/sirupsen/logrus"
	"context"

	"github.com/redis/go-redis/v9"
)

const (
	BF_ErrorRate = 0.01
	BF_InitialSize = 1000000
)

func BFReserve(
	ctx context.Context,
	client *redis.Client,
	key string,
) *redis.Cmd {
	return client.Do(
		ctx,
		"BF.RESERVE",
		key,
		BF_ErrorRate,
		BF_InitialSize,
	)
}

func BFMAdd(
	ctx context.Context,
	client *redis.Client,
	key string,
	items ...any,
) *redis.Cmd {
	args := make([]any, 2+len(items))
	args[0] = "BF.MADD"
	args[1] = key
	copy(args[2:], items)

	return client.Do(ctx, args...)
}

func BFExists(
	ctx context.Context,
	client *redis.Client,
	key string,
	item any,
) bool {
	res, err := client.Do(ctx, "BF.EXISTS", key, item).Int()
	if err != nil {
		log.Errorf("Failed to check Bloom filter existence: %v", err)
		return false
	}
	return res == 1
}

func BFAdd(
	ctx context.Context,
	client *redis.Client,
	key string,
	item any,
) *redis.Cmd {
	return client.Do(ctx, "BF.ADD", key, item)
}
