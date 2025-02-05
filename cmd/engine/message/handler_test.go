package message

import (
	"context"
	"testing"
	"time"

	"github.com/mehgokalp/insider-project/pkg/database"
	"github.com/mehgokalp/insider-project/pkg/mocks"
	"github.com/mehgokalp/insider-project/pkg/provider/webhook"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestHandler_Handle(t *testing.T) {
	ctx := context.Background()
	logger := mocks.NewLoggerMock("Infof", "Errorf")

	mockMessageRepository := mocks.NewMessageRepository(t)
	mockRequester := new(mocks.Requester)
	mockRedisMessageRepository := mocks.NewRedisMessageRepository(t)

	h := &handler{
		consume:                true,
		logger:                 logger,
		requester:              mockRequester,
		messageRepository:      mockMessageRepository,
		redisMessageRepository: mockRedisMessageRepository,
		batchMessageLimit:      2,
		tickerInterval:         2 * time.Second,
	}

	messages := []database.Message{
		{ID: 1, To: "1234567890", Content: "Test message 1"},
		{ID: 2, To: "0987654321", Content: "Test message 2"},
	}

	mockMessageRepository.On("GetUnsentMessages", 2).Return(messages, nil)
	mockRequester.On("SendSMS", mock.Anything).Return(&webhook.SendMessageResponse{ID: "msg-1"}, nil)
	mockMessageRepository.On("UpdateSentStatus", mock.Anything).Return(nil)
	mockRedisMessageRepository.On("Save", ctx, mock.Anything).Return(nil)

	err := h.handle(ctx)
	assert.NoError(t, err)

	mockMessageRepository.AssertExpectations(t)
	mockRequester.AssertExpectations(t)
	mockRedisMessageRepository.AssertExpectations(t)
}
