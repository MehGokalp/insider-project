package server

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	_ "github.com/mehgokalp/insider-project/cmd/server/docs"
	pkgMessageList "github.com/mehgokalp/insider-project/cmd/server/modules/message/list"
	pkgMessageStartStop "github.com/mehgokalp/insider-project/cmd/server/modules/message/start_stop"
	"github.com/mehgokalp/insider-project/pkg/config"
	pkgDatabaseRepository "github.com/mehgokalp/insider-project/pkg/database/repository"
	"github.com/mehgokalp/insider-project/pkg/log"
	"github.com/mehgokalp/insider-project/pkg/meta"
	pkgRedisRepository "github.com/mehgokalp/insider-project/pkg/redis/repository"
	"github.com/spf13/cobra"
	"github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"time"
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
			r := getRouter(
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

func getRouter(
	logger log.Logger,
	messageRepository pkgDatabaseRepository.MessageRepository,
	redisMessageEngineRepository pkgRedisRepository.RedisMessageEngineRepository,
) *gin.Engine {
	r := gin.New()
	r.Use(gin.ErrorLogger())
	r.Use(jsonLoggerMiddleware())
	r.Use(gin.Recovery())

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	v1 := r.Group("/v1")

	v1.GET("/messages/", pkgMessageList.NewHandler(logger, messageRepository))
	v1.PATCH("/messages/", pkgMessageStartStop.NewHandler(logger, redisMessageEngineRepository))

	return r
}

func jsonLoggerMiddleware() gin.HandlerFunc {
	return gin.LoggerWithFormatter(
		func(params gin.LogFormatterParams) string {
			line := make(map[string]interface{})

			line["app_name"] = meta.AppName
			line["level"] = "debug"
			line["status_code"] = params.StatusCode
			line["path"] = params.Path
			line["method"] = params.Method
			line["remote_addr"] = params.ClientIP
			line["response_time"] = params.Latency.String()
			line["time"] = params.TimeStamp.Format(time.RFC3339)

			s, _ := json.Marshal(line)
			return string(s) + "\n"
		},
	)
}
