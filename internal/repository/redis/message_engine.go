package redis

import (
	"context"
	"encoding/json"
	"github.com/go-redis/redis/v8"
	pkgRedis "github.com/mehgokalp/insider-project/internal/db/redis"
	redis2 "github.com/mehgokalp/insider-project/internal/domain/redis"
	"github.com/rotisserie/eris"
)

const messageEngineStatusChannel = "message_engine_status"

type RedisMessageEngineRepository interface {
	UpdateStatus(ctx context.Context, status redis2.MessageEngineRunningStatus) error
	ListenStatusUpdates(ctx context.Context) *redis.PubSub
}

type MessageEngineRepository struct {
	client pkgRedis.Client
}

func NewMessageEngineRepository(client pkgRedis.Client) *MessageEngineRepository {
	return &MessageEngineRepository{client: client}
}

func (r *MessageEngineRepository) UpdateStatus(ctx context.Context, status redis2.MessageEngineRunningStatus) error {
	parsed, err := json.Marshal(status)
	if err != nil {
		return eris.Wrap(err, "failed to marshal message")
	}

	return r.client.Publish(ctx, messageEngineStatusChannel, parsed).Err()
}

func (r *MessageEngineRepository) ListenStatusUpdates(ctx context.Context) *redis.PubSub {
	return r.client.Subscribe(ctx, messageEngineStatusChannel)
}
