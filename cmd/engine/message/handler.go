package message

import (
	"context"
	"encoding/json"
	pkgRedisDomain "github.com/mehgokalp/insider-project/internal/domain/redis"
	"github.com/mehgokalp/insider-project/internal/message/service"
	"github.com/mehgokalp/insider-project/internal/repository/redis"
	pkgLog "github.com/mehgokalp/insider-project/pkg/log"
	"github.com/rotisserie/eris"
	"time"
)

type handler struct {
	consume                      bool
	logger                       pkgLog.Logger
	batchMessageLimit            int
	tickerInterval               time.Duration
	redisMessageEngineRepository redis.RedisMessageEngineRepository
	service                      service.SendUnsentSmsService
}

func (h *handler) listen(ctx context.Context) {
	go func() {
		pubsub := h.redisMessageEngineRepository.ListenStatusUpdates(ctx)
		defer pubsub.Close()

		h.logger.Infof("listening status updates")
		for msg := range pubsub.Channel() {
			var status pkgRedisDomain.MessageEngineRunningStatus
			err := json.Unmarshal([]byte(msg.Payload), &status)
			if err != nil {
				panic(eris.Wrap(err, "can not read status update"))
			}

			h.consume = status.Consume
			h.logger.Infof("consume status updated")
		}
	}()

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

	errors := h.service.SendUnsentSms(ctx, h.batchMessageLimit)
	if len(errors) > 0 {
		return eris.New("error while sending unsent messages")
	}

	return nil
}
