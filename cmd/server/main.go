package server

import (
	"fmt"
	_ "github.com/mehgokalp/insider-project/cmd/server/docs"
	"github.com/mehgokalp/insider-project/config"
	pkgDatabaseRepository "github.com/mehgokalp/insider-project/internal/repository/mysql"
	pkgRedisRepository "github.com/mehgokalp/insider-project/internal/repository/redis"
	"github.com/mehgokalp/insider-project/internal/server"
	"github.com/mehgokalp/insider-project/pkg/log"
	"github.com/spf13/cobra"
)

// @title Messages API
// @version 1.0
// @description This is a sample server for managing messages.
// @host localhost:8081
// @BasePath /v1

func Server(
	cfg *config.Config,
	logger log.Logger,
	messageRepository pkgDatabaseRepository.MessageRepository,
	redisMessageEngineRepository pkgRedisRepository.RedisMessageEngineRepository,
) *cobra.Command {
	cmdName := "server"

	return &cobra.Command{
		Use:   cmdName,
		Short: "Run backend server",
		RunE: func(cmd *cobra.Command, _ []string) error {
			r := server.GetRouter(
				logger,
				messageRepository,
				redisMessageEngineRepository,
			)

			if err := r.Run(fmt.Sprintf(":%v", cfg.Port)); err != nil {
				return err
			}

			return nil
		},
	}
}
