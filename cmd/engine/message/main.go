package message

import (
	"context"
	"encoding/json"
	pkgDatabaseRepository "github.com/mehgokalp/insider-project/pkg/database/repository"
	pkgLog "github.com/mehgokalp/insider-project/pkg/log"
	"github.com/mehgokalp/insider-project/pkg/provider/webhook"
	"github.com/mehgokalp/insider-project/pkg/redis"
	pkgRedisRepository "github.com/mehgokalp/insider-project/pkg/redis/repository"
	"github.com/rotisserie/eris"
	"github.com/spf13/cobra"
	"time"
)

const batchMessageLimit = 2

func MessageCmd(
	ctx context.Context,
	logger pkgLog.Logger,
	requester webhook.Requester,
	messageRepository pkgDatabaseRepository.MessageRepository,
	redisMessageRepository pkgRedisRepository.RedisMessageRepository,
	redisMessageEngineRepository pkgRedisRepository.RedisMessageEngineRepository,
) *cobra.Command {
	cmdName := "engine:message"

	return &cobra.Command{
		Use:   cmdName,
		Short: "Run message engine",
		RunE: func(cmd *cobra.Command, _ []string) error {

			logger.Infof("Message engine is running")
			handler := handler{
				consume:                true,
				requester:              requester,
				logger:                 logger,
				messageRepository:      messageRepository,
				redisMessageRepository: redisMessageRepository,
				batchMessageLimit:      batchMessageLimit,
				tickerInterval:         2 * time.Minute,
			}

			go func() {
				pubsub := redisMessageEngineRepository.ListenStatusUpdates(ctx)
				defer pubsub.Close()

				logger.Infof("listening status updates")
				for msg := range pubsub.Channel() {
					var status redis.MessageEngineRunningStatus
					err := json.Unmarshal([]byte(msg.Payload), &status)
					if err != nil {
						panic(eris.Wrap(err, "can not read status update"))
					}

					handler.consume = status.Consume
					logger.Infof("consume status updated")
				}
			}()

			handler.listen(ctx)

			return nil
		},
	}
}
