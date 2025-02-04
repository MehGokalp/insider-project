package server

import (
	"github.com/gin-gonic/gin"
	pkgMessageList "github.com/mehgokalp/insider-project/internal/delivery/http/message/list"
	pkgMessageStartStop "github.com/mehgokalp/insider-project/internal/delivery/http/message/start_stop"
	pkgDatabaseRepository "github.com/mehgokalp/insider-project/internal/repository/mysql"
	pkgRedisRepository "github.com/mehgokalp/insider-project/internal/repository/redis"
	"github.com/mehgokalp/insider-project/pkg/log"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func GetRouter(
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
