package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/mehgokalp/insider-project/pkg/redis"
	"github.com/rotisserie/eris"
	"time"
)

const MessageRepositoryPrefix = "msg"

type MessageRepository struct {
	client redis.Client
	prefix string
}

func NewMessageRepository(client redis.Client, prefix string) *MessageRepository {
	return &MessageRepository{client: client, prefix: prefix}
}

func (r *MessageRepository) Save(ctx context.Context, message redis.Message, duration time.Duration) error {
	parsed, err := json.Marshal(message)
	if err != nil {
		return eris.Wrap(err, "failed to marshal message")
	}

	return r.client.Set(ctx, fmt.Sprintf("%v_%v", r.prefix, message.ID), string(parsed), duration).Err()
}
