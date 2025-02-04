package message

import (
	"context"
	"github.com/mehgokalp/insider-project/internal/message/service"
	"github.com/mehgokalp/insider-project/internal/repository/redis"
	pkgLog "github.com/mehgokalp/insider-project/pkg/log"
	"github.com/spf13/cobra"
	"time"
)

const batchMessageLimit = 2

func MessageCmd(
	ctx context.Context,
	logger pkgLog.Logger,
	redisMessageEngineRepository redis.RedisMessageEngineRepository,
	service service.SendUnsentSmsService,
) *cobra.Command {
	cmdName := "engine:message"

	return &cobra.Command{
		Use:   cmdName,
		Short: "Run message engine",
		RunE: func(cmd *cobra.Command, _ []string) error {

			logger.Infof("Message engine is running")
			handler := handler{
				consume:                      true,
				logger:                       logger,
				batchMessageLimit:            batchMessageLimit,
				tickerInterval:               2 * time.Minute,
				redisMessageEngineRepository: redisMessageEngineRepository,
				service:                      service,
			}

			handler.listen(ctx)

			return nil
		},
	}
}
