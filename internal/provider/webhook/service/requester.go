package service

import (
	"bytes"
	"encoding/json"
	"github.com/go-playground/validator/v10"
	"github.com/mehgokalp/insider-project/internal/provider/webhook/dto"
	"github.com/mehgokalp/insider-project/internal/provider/webhook/interfaces"
	pkgLog "github.com/mehgokalp/insider-project/pkg/log"
	"github.com/rotisserie/eris"
	"io"
	"net/http"
)

type Requester interface {
	SendSMS(request dto.SendMessageRequest) (*dto.SendMessageResponse, error)
}

type HttpRequester struct {
	client    interfaces.HTTPClient
	baseUri   string
	logger    pkgLog.Logger
	validator *validator.Validate
}

func NewRequester(client interfaces.HTTPClient, baseUri string, logger pkgLog.Logger, validator *validator.Validate) *HttpRequester {
	return &HttpRequester{
		client:    client,
		baseUri:   baseUri,
		logger:    logger,
		validator: validator,
	}
}

func (r *HttpRequester) SendSMS(request dto.SendMessageRequest) (*dto.SendMessageResponse, error) {
	// validate request
	if err := r.validator.Struct(request); err != nil {
		return nil, eris.Wrap(err, "request validation error")
	}

	parsed, _ := json.Marshal(request)

	req, _ := http.NewRequest("GET", r.baseUri, bytes.NewBuffer(parsed))
	req.Header.Set("Content-Type", "application/json")

	resp, err := r.client.Do(req)
	// Handle Error
	if err != nil {
		return nil, eris.Wrap(err, "unknown error returned")
	}

	body, _ := io.ReadAll(resp.Body)
	defer resp.Body.Close()

	if resp.StatusCode != 202 {
		return nil, eris.New("status code is not expected")
	}

	var res dto.SendMessageResponse
	if err := json.Unmarshal(body, &res); err != nil {
		return nil, eris.Wrap(err, "json unmarshall err returned")
	}

	if err := r.validator.Struct(res); err != nil {
		return nil, eris.Wrap(err, "response validation error")
	}

	r.logger.Debugf("send sms response: %v", string(body))

	return &res, nil
}
