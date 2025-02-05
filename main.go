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
	"github.com/mehgokalp/insider-project/pkg/config"
	"github.com/mehgokalp/insider-project/pkg/database"
	pkgDatabaseRepository "github.com/mehgokalp/insider-project/pkg/database/repository"
	"github.com/mehgokalp/insider-project/pkg/log"
	"github.com/mehgokalp/insider-project/pkg/provider/webhook"
	pkgRedisRepository "github.com/mehgokalp/insider-project/pkg/redis/repository"
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

	err = database.AutoMigrate(db)
	if err != nil {
		panic(eris.Wrap(err, "migration failed"))
	}

	messageRepository := pkgDatabaseRepository.NewMessageRepository(db)

	requester := webhook.NewRequester(&http.Client{}, cfg.MessageProvider.BaseUrl, logger, validator.New())

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

	rootCmd.AddCommand(message.MessageCmd(ctx, logger, requester, messageRepository, redisMessageRepository, redisMessageEngineRepository))
	rootCmd.AddCommand(cmd.PopulateCmd(db))

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(eris.ToString(err, true))
	}
}
