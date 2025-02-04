package message

import (
	"context"
	pkgLog "github.com/mehgokalp/insider-project/pkg/log"
	"github.com/mehgokalp/insider-project/pkg/provider/webhook"
	pkgRedisRepository "github.com/mehgokalp/insider-project/pkg/redis/repository"
	"github.com/spf13/cobra"
)

func MessageCmd(ctx context.Context, logger pkgLog.Logger, requester *webhook.Requester, redisMessageRepository *pkgRedisRepository.MessageRepository) *cobra.Command {
	cmdName := "engine:message"

	return &cobra.Command{
		Use:   cmdName,
		Short: "Run message engine",
		RunE: func(cmd *cobra.Command, _ []string) error {

			logger.Infof("Message engine is running")

			return nil
		},
	}
}
