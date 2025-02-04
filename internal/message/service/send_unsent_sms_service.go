package service

import (
	"context"
	"github.com/mehgokalp/insider-project/internal/message/usecase"
)

type SendUnsentSmsService struct {
	usecase usecase.SendUnsentSMSUsecase
}

func NewSendUnsentSmsService(usecase usecase.SendUnsentSMSUsecase) SendUnsentSmsService {
	return SendUnsentSmsService{
		usecase: usecase,
	}
}

func (s SendUnsentSmsService) SendUnsentSms(ctx context.Context, limit int) []error {
	return s.usecase.SendUnsentSms(ctx, limit)
}
