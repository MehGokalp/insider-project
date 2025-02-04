package redis

import (
	"context"
	"github.com/go-redis/redis/v8"
	"time"
)

type (
	Client interface {
		Get(ctx context.Context, key string) *redis.StringCmd
		Set(ctx context.Context, key string, value interface{}, expiration time.Duration) *redis.StatusCmd
		Publish(ctx context.Context, channel string, message interface{}) *redis.IntCmd
		Subscribe(ctx context.Context, channels ...string) *redis.PubSub
	}
)
