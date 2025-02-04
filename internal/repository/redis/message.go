package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/mehgokalp/insider-project/internal/db/redis"
	redis2 "github.com/mehgokalp/insider-project/internal/domain/redis"
	"github.com/rotisserie/eris"
)

const MessageRepositoryPrefix = "msg"

type Repository interface {
	Save(ctx context.Context, message redis2.Message) error
}

type MessageRepository struct {
	client redis.Client
	prefix string
}

func NewMessageRepository(client redis.Client, prefix string) *MessageRepository {
	return &MessageRepository{client: client, prefix: prefix}
}

func (r *MessageRepository) Save(ctx context.Context, message redis2.Message) error {
	parsed, err := json.Marshal(message)
	if err != nil {
		return eris.Wrap(err, "failed to marshal message")
	}

	return r.client.Set(ctx, fmt.Sprintf("%v_%v", r.prefix, message.ID), string(parsed), message.Duration()).Err()
}
