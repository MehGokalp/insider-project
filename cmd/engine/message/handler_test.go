package message

import (
	"context"
	"github.com/mehgokalp/insider-project/internal/domain/mysql"
	"github.com/mehgokalp/insider-project/internal/message/service"
	"github.com/mehgokalp/insider-project/internal/message/usecase"
	mocks2 "github.com/mehgokalp/insider-project/internal/mocks"
	"github.com/mehgokalp/insider-project/internal/provider/webhook/dto"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestHandler_Handle(t *testing.T) {
	ctx := context.Background()
	logger := mocks2.NewLoggerMock("Infof", "Errorf")

	mockMessageRepository := mocks2.NewMessageRepository(t)
	mockRequester := new(mocks2.Requester)
	mockRedisMessageRepository := mocks2.NewRedisMessageRepository(t)

	sendSmsUseCase := usecase.NewSendSMSUsecase(
		mockRequester,
		mockMessageRepository,
		mockRedisMessageRepository,
	)
	sendUnsentSmsUseCase := usecase.NewSendUnsentSMSUsecase(
		sendSmsUseCase,
		mockMessageRepository,
		logger,
	)

	sendUnsentSmsService := service.NewSendUnsentSmsService(sendUnsentSmsUseCase)

	h := &handler{
		consume:           true,
		logger:            logger,
		batchMessageLimit: 2,
		tickerInterval:    2 * time.Second,
		service:           sendUnsentSmsService,
	}

	messages := []mysql.Message{
		{ID: 1, To: "1234567890", Content: "Test message 1"},
		{ID: 2, To: "0987654321", Content: "Test message 2"},
	}

	mockMessageRepository.On("GetUnsentMessages", 2).Return(messages, nil)
	mockRequester.On("SendSMS", mock.Anything).Return(&dto.SendMessageResponse{ID: "msg-1"}, nil)
	mockMessageRepository.On("UpdateSentStatus", mock.Anything).Return(nil)
	mockRedisMessageRepository.On("Save", ctx, mock.Anything).Return(nil)

	err := h.handle(ctx)
	assert.NoError(t, err)

	mockMessageRepository.AssertExpectations(t)
	mockRequester.AssertExpectations(t)
	mockRedisMessageRepository.AssertExpectations(t)
}
