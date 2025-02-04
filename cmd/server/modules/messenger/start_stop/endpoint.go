package start_stop

import (
	"github.com/gin-gonic/gin"
	"github.com/mehgokalp/insider-project/pkg/log"
	"net/http"
)

type Handler struct {
	ctx    *gin.Context
	logger log.Logger
}

func NewHandler(logger log.Logger) func(*gin.Context) {
	return func(c *gin.Context) {
		h := Handler{
			ctx:    c,
			logger: logger,
		}

		h.Handle()
	}
}

func (h *Handler) Handle() {
	h.ctx.JSON(http.StatusOK, "endpoint")
}
