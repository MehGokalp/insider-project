package start_stop

import (
	"github.com/gin-gonic/gin"
	"github.com/mehgokalp/insider-project/pkg/log"
	"github.com/mehgokalp/insider-project/pkg/redis"
	pkgRedisRepository "github.com/mehgokalp/insider-project/pkg/redis/repository"
	"github.com/rotisserie/eris"
	"net/http"
)

type Handler struct {
	ctx        *gin.Context
	logger     log.Logger
	repository *pkgRedisRepository.MessageEngineRepository
}

func NewHandler(logger log.Logger, repository *pkgRedisRepository.MessageEngineRepository) func(*gin.Context) {
	return func(c *gin.Context) {
		h := Handler{
			ctx:        c,
			logger:     logger,
			repository: repository,
		}

		h.Handle()
	}
}

type startStopForm struct {
	Status string `json:"status" binding:"required,oneof=start stop"`
}

func (h *Handler) Handle() {
	var form startStopForm
	err := h.ctx.BindJSON(&form)

	if err != nil {
		h.logger.Errorf(eris.Wrap(err, "failed to bind json").Error())

		h.ctx.AbortWithStatusJSON(http.StatusNotAcceptable, err)
		return
	}

	err = h.repository.Save(h.ctx, redis.MessageEngineRunningStatus{
		Consume: form.Status == "start",
	})
	if err != nil {
		h.logger.Errorf(eris.Wrap(err, "failed to save status").Error())

		h.ctx.AbortWithStatusJSON(http.StatusInternalServerError, err)
		return
	}

	h.ctx.JSON(http.StatusAccepted, nil)
}
