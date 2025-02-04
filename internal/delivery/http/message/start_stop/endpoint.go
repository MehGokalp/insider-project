package start_stop

import (
	"github.com/gin-gonic/gin"
	"github.com/mehgokalp/insider-project/internal/domain/redis"
	pkgRedisRepository "github.com/mehgokalp/insider-project/internal/repository/redis"
	"github.com/mehgokalp/insider-project/pkg/log"
	"github.com/rotisserie/eris"
	"net/http"
)

type Handler struct {
	ctx        *gin.Context
	logger     log.Logger
	repository pkgRedisRepository.RedisMessageEngineRepository
}

// NewHandler godoc
// @Summary Update message status
// @Description Update the status of the message engine
// @Tags messages
// @Accept  json
// @Produce  json
// @Param action body startStopForm true "Action to start or stop the message engine"
// @Success 202 {object} nil
// @Failure 406 {object} error
// @Failure 500 {object} error
// @Router /messages/ [patch]
func NewHandler(logger log.Logger, repository pkgRedisRepository.RedisMessageEngineRepository) func(*gin.Context) {
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
	Action string `json:"action" binding:"required,oneof=start stop"`
}

func (h *Handler) Handle() {
	var form startStopForm
	err := h.ctx.BindJSON(&form)

	if err != nil {
		h.logger.Errorf(eris.Wrap(err, "failed to bind json").Error())

		h.ctx.AbortWithStatusJSON(http.StatusNotAcceptable, err)
		return
	}

	err = h.repository.UpdateStatus(h.ctx, redis.MessageEngineRunningStatus{
		Consume: form.Action == "start",
	})
	if err != nil {
		h.logger.Errorf(eris.Wrap(err, "failed to save status").Error())

		h.ctx.AbortWithStatusJSON(http.StatusInternalServerError, err)
		return
	}

	h.ctx.JSON(http.StatusAccepted, nil)
}
