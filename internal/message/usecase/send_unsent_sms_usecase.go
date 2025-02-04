package usecase

import (
	"context"
	pkgDatabaseRepository "github.com/mehgokalp/insider-project/internal/repository/mysql"
	pkgLog "github.com/mehgokalp/insider-project/pkg/log"
	"github.com/rotisserie/eris"
)

type SendUnsentSMSUsecase struct {
	sendSmsUseCase    SendSMSUsecase
	messageRepository pkgDatabaseRepository.MessageRepository
	logger            pkgLog.Logger
}

func NewSendUnsentSMSUsecase(
	sendSmsUseCase SendSMSUsecase,
	messageRepository pkgDatabaseRepository.MessageRepository,
	logger pkgLog.Logger,
) SendUnsentSMSUsecase {
	return SendUnsentSMSUsecase{
		sendSmsUseCase:    sendSmsUseCase,
		messageRepository: messageRepository,
		logger:            logger,
	}
}

func (u *SendUnsentSMSUsecase) SendUnsentSms(ctx context.Context, limit int) []error {
	messages, err := u.messageRepository.GetUnsentMessages(limit)
	if err != nil {
		return []error{eris.Wrap(err, "error while getting unsent messages")}
	}

	var errors []error
	for _, message := range messages {
		err = u.sendSmsUseCase.SendSMS(ctx, message)
		if err != nil {
			errors = append(errors, eris.Wrap(err, "error while sending message"))
		}

		u.logger.Infof("message sent to: %v", message.To)
	}

	return errors
}
