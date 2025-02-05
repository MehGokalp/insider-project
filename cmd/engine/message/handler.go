package message

import (
	"context"
	"github.com/mehgokalp/insider-project/pkg/database"
	pkgDatabaseRepository "github.com/mehgokalp/insider-project/pkg/database/repository"
	pkgLog "github.com/mehgokalp/insider-project/pkg/log"
	"github.com/mehgokalp/insider-project/pkg/provider/webhook"
	"github.com/mehgokalp/insider-project/pkg/redis"
	pkgRedisRepository "github.com/mehgokalp/insider-project/pkg/redis/repository"
	"github.com/rotisserie/eris"
	"time"
)

type handler struct {
	consume                bool
	logger                 pkgLog.Logger
	requester              webhook.Requester
	messageRepository      pkgDatabaseRepository.MessageRepository
	redisMessageRepository pkgRedisRepository.RedisMessageRepository
	batchMessageLimit      int
	tickerInterval         time.Duration
}

func (h *handler) listen(ctx context.Context) {
	ticker := time.NewTicker(h.tickerInterval)
	for range ticker.C {
		if !h.consume {
			continue
		}

		err := h.handle(ctx)
		if err != nil {
			h.logger.Errorf("error on message engine: %v", err)
		}
	}
}

func (h *handler) handle(ctx context.Context) error {
	if !h.consume {
		return nil
	}

	messages, err := h.messageRepository.GetUnsentMessages(h.batchMessageLimit)
	if err != nil {
		return eris.Wrap(err, "error while getting unsent messages")
	}

	for _, message := range messages {
		err = h.sendMessage(ctx, message)
		if err != nil {
			h.logger.Errorf("error while sending message: %v", err)
			continue
		}

		h.logger.Infof("message sent to: %v", message.To)
	}

	return nil
}

func (h *handler) sendMessage(ctx context.Context, message database.Message) error {
	resp, err := h.requester.SendSMS(webhook.SendMessageRequest{
		To:      message.To,
		Content: message.Content,
	})
	if err != nil {
		return eris.Wrap(err, "error while sending message")
	}
	message.MessageId = resp.ID

	err = h.messageRepository.UpdateSentStatus(message)
	if err != nil {
		return eris.Wrap(err, "error while updating message")
	}

	err = h.redisMessageRepository.Save(ctx, redis.Message{
		ID:   resp.ID,
		Time: time.Now().Format(time.RFC3339),
	})

	if err != nil {
		return eris.Wrap(err, "error while saving message to redis")
	}

	return nil
}
