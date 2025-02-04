package server

import (
	"bytes"
	"encoding/json"
	"github.com/mehgokalp/insider-project/internal/domain/mysql"
	mocks2 "github.com/mehgokalp/insider-project/internal/mocks"
	"github.com/mehgokalp/insider-project/internal/server"
	"github.com/stretchr/testify/mock"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestGetMessages(t *testing.T) {
	gin.SetMode(gin.TestMode)

	logger := mocks2.NewLoggerMock("Infof", "Errorf")
	messageRepository := mocks2.NewMessageRepository(t)
	redisMessageEngineRepository := mocks2.NewRedisMessageEngineRepository(t)

	// Mock the List method
	messages := []mysql.Message{
		{ID: 1, To: "1234567890", Content: "Test message 1"},
		{ID: 2, To: "0987654321", Content: "Test message 2"},
	}
	messageRepository.On("List").Return(messages, nil)

	router := server.GetRouter(logger, messageRepository, redisMessageEngineRepository)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/v1/messages/", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	// Add more assertions based on the expected response
}

func TestPatchMessages(t *testing.T) {
	gin.SetMode(gin.TestMode)

	logger := mocks2.NewLoggerMock("Infof", "Errorf")
	messageRepository := mocks2.NewMessageRepository(t)
	redisMessageEngineRepository := mocks2.NewRedisMessageEngineRepository(t)

	// Mock the UpdateStatus method
	redisMessageEngineRepository.On("UpdateStatus", mock.Anything, mock.Anything).Return(nil)

	router := server.GetRouter(logger, messageRepository, redisMessageEngineRepository)

	body := map[string]interface{}{
		"action": "start",
	}
	jsonBody, _ := json.Marshal(body)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PATCH", "/v1/messages/", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusAccepted, w.Code)
	// Add more assertions based on the expected response
}
