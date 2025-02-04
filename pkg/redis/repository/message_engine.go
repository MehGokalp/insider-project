package repository

import (
	"context"
	"encoding/json"
	"github.com/mehgokalp/insider-project/pkg/redis"
	"github.com/rotisserie/eris"
)

const messageEngineStatusKey = "message_engine_status"

type MessageEngineRepository struct {
	client redis.Client
}

func NewMessageEngineRepository(client redis.Client) *MessageEngineRepository {
	return &MessageEngineRepository{client: client}
}

func (r *MessageEngineRepository) Save(ctx context.Context, status redis.MessageEngineRunningStatus) error {
	parsed, err := json.Marshal(status)
	if err != nil {
		return eris.Wrap(err, "failed to marshal message")
	}

	return r.client.Set(ctx, messageEngineStatusKey, parsed, status.Duration()).Err()
}
