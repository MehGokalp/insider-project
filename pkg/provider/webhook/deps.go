package webhook

import (
	"bytes"
	"encoding/json"
	"github.com/go-playground/validator/v10"
	pkgLog "github.com/mehgokalp/insider-project/pkg/log"
	"github.com/rotisserie/eris"
	"io"
	"net/http"
)

type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type Requester struct {
	client    HTTPClient
	baseUri   string
	logger    pkgLog.Logger
	validator *validator.Validate
}

func NewRequester(client HTTPClient, baseUri string, logger pkgLog.Logger) *Requester {
	return &Requester{
		client:  client,
		baseUri: baseUri,
		logger:  logger,
	}
}

func (r *Requester) SendSMS(request SendMessageRequest) (*SendMessageResponse, error) {
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

	if resp.StatusCode != 200 {
		return nil, eris.New("status code is not expected")
	}

	var res SendMessageResponse
	if err := json.Unmarshal(body, &res); err != nil {
		return nil, eris.Wrap(err, "json unmarshall err returned")
	}

	if err := r.validator.Struct(res); err != nil {
		return nil, eris.Wrap(err, "response validation error")
	}

	r.logger.Debugf("send sms response: %v", string(body))

	return &res, nil
}
