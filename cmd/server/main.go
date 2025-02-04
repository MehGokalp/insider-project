package server

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	pkgMessengerList "github.com/mehgokalp/insider-project/cmd/server/modules/messenger/list"
	pkgMessengerStartStop "github.com/mehgokalp/insider-project/cmd/server/modules/messenger/start_stop"
	"github.com/mehgokalp/insider-project/pkg/config"
	pkgDatabaseRepository "github.com/mehgokalp/insider-project/pkg/database/repository"
	"github.com/mehgokalp/insider-project/pkg/log"
	"github.com/mehgokalp/insider-project/pkg/meta"
	"github.com/spf13/cobra"
	"time"
)

func Server(
	cfg *config.Config,
	logger log.Logger,
	pkgDatabaseRepository *pkgDatabaseRepository.MessageRepository,
) *cobra.Command {
	cmdName := "server"

	return &cobra.Command{
		Use:   cmdName,
		Short: "Run backend server",
		RunE: func(cmd *cobra.Command, _ []string) error {
			r := getRouter(
				logger,
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
) *gin.Engine {
	r := gin.New()
	r.Use(gin.ErrorLogger())
	r.Use(jsonLoggerMiddleware())
	r.Use(gin.Recovery())

	v1 := r.Group("/v1")

	v1.GET("/messenger/", pkgMessengerList.NewHandler(logger))
	v1.PUT("/messenger/start-stop", pkgMessengerStartStop.NewHandler(logger))

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
