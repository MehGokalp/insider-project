package usecase

import (
	"context"
	"github.com/mehgokalp/insider-project/internal/domain/mysql"
	"github.com/mehgokalp/insider-project/internal/domain/redis"
	"github.com/mehgokalp/insider-project/internal/provider/webhook/dto"
	pkgWebhookService "github.com/mehgokalp/insider-project/internal/provider/webhook/service"
	pkgDatabaseRepository "github.com/mehgokalp/insider-project/internal/repository/mysql"
	pkgRedisRepository "github.com/mehgokalp/insider-project/internal/repository/redis"
	"github.com/rotisserie/eris"
	"time"
)

type SendSMSUsecase struct {
	requester              pkgWebhookService.Requester
	messageRepository      pkgDatabaseRepository.MessageRepository
	redisMessageRepository pkgRedisRepository.Repository
}

func NewSendSMSUsecase(
	requester pkgWebhookService.Requester,
	messageRepository pkgDatabaseRepository.MessageRepository,
	redisMessageRepository pkgRedisRepository.Repository,
) SendSMSUsecase {
	return SendSMSUsecase{
		requester:              requester,
		messageRepository:      messageRepository,
		redisMessageRepository: redisMessageRepository,
	}
}

func (u *SendSMSUsecase) SendSMS(ctx context.Context, message mysql.Message) error {
	resp, err := u.requester.SendSMS(dto.SendMessageRequest{
		To:      message.To,
		Content: message.Content,
	})
	if err != nil {
		return eris.Wrap(err, "error while sending message")
	}
	message.MessageId = resp.ID

	err = u.messageRepository.UpdateSentStatus(message)
	if err != nil {
		return eris.Wrap(err, "error while updating message")
	}

	err = u.redisMessageRepository.Save(ctx, redis.Message{
		ID:   resp.ID,
		Time: time.Now().Format(time.RFC3339),
	})

	if err != nil {
		return eris.Wrap(err, "error while saving message to redis")
	}

	return nil
}
