package main

import (
	"context"
	"crypto/tls"
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/go-redis/redis/v8"
	"github.com/mehgokalp/insider-project/cmd"
	"github.com/mehgokalp/insider-project/cmd/engine/message"
	"github.com/mehgokalp/insider-project/cmd/server"
	"github.com/mehgokalp/insider-project/config"
	"github.com/mehgokalp/insider-project/internal/message/service"
	"github.com/mehgokalp/insider-project/internal/message/usecase"
	pkgWebhookService "github.com/mehgokalp/insider-project/internal/provider/webhook/service"
	pkgDatabaseRepository "github.com/mehgokalp/insider-project/internal/repository/mysql"
	pkgRedisRepository "github.com/mehgokalp/insider-project/internal/repository/redis"
	"github.com/mehgokalp/insider-project/migrations"
	"github.com/mehgokalp/insider-project/pkg/log"
	"github.com/rotisserie/eris"
	"github.com/spf13/cobra"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"net/http"
)

func main() {
	rootCmd := &cobra.Command{
		Use:   "insider-project",
		Short: "Main entry-point command for the application",
	}

	ctx := context.Background()
	cfg := config.New()
	logger := log.New()

	db, err := gorm.Open(mysql.Open(cfg.Mysql.DSN), &gorm.Config{})
	if err != nil {
		panic(eris.Wrap(err, "failed to connect to database"))
	}

	err = migrations.AutoMigrate(db)
	if err != nil {
		panic(eris.Wrap(err, "migration failed"))
	}

	messageRepository := pkgDatabaseRepository.NewMessageRepository(db)

	requester := pkgWebhookService.NewRequester(&http.Client{}, cfg.MessageProvider.BaseUrl, logger, validator.New())

	redisOpt, err := redis.ParseURL(cfg.Redis.DSN)
	if err != nil {
		panic(err)
	}
	if cfg.Env != "dev" {
		redisOpt.TLSConfig = &tls.Config{}
	}
	redisClient := redis.NewClient(redisOpt)

	redisMessageRepository := pkgRedisRepository.NewMessageRepository(redisClient, pkgRedisRepository.MessageRepositoryPrefix)
	redisMessageEngineRepository := pkgRedisRepository.NewMessageEngineRepository(redisClient)

	rootCmd.AddCommand(
		server.Server(
			cfg,
			logger,
			messageRepository,
			redisMessageEngineRepository,
		),
	)

	sendSmsUseCase := usecase.NewSendSMSUsecase(
		requester,
		messageRepository,
		redisMessageRepository,
	)
	sendUnsentSmsUseCase := usecase.NewSendUnsentSMSUsecase(
		sendSmsUseCase,
		messageRepository,
		logger,
	)

	sendUnsentSmsService := service.NewSendUnsentSmsService(sendUnsentSmsUseCase)

	rootCmd.AddCommand(message.MessageCmd(ctx, logger, redisMessageEngineRepository, sendUnsentSmsService))
	rootCmd.AddCommand(cmd.PopulateCmd(db))

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(eris.ToString(err, true))
	}
}
