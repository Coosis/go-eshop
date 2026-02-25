package comm

import (
	"context"
	log "github.com/sirupsen/logrus"
	"strings"

	"github.com/redis/go-redis/v9"
)

const (
	BF_ErrorRate   = 0.01
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
	res, err := client.Do(ctx, "BF.EXISTS", key, item).Result()
	if err != nil {
		log.Errorf("Failed to check Bloom filter existence: %v", err)
		return false
	}
	switch v := res.(type) {
	case bool:
		return v
	case int64:
		return v == 1
	case string:
		s := strings.TrimSpace(strings.ToLower(v))
		return s == "1" || s == "true"
	default:
		log.Errorf("Unexpected Bloom EXISTS response type: %T (%v)", res, res)
		return false
	}
}

func BFAdd(
	ctx context.Context,
	client *redis.Client,
	key string,
	item any,
) *redis.Cmd {
	return client.Do(ctx, "BF.ADD", key, item)
}
