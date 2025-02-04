package list

import (
	"github.com/gin-gonic/gin"
	pkgDatabaseRepository "github.com/mehgokalp/insider-project/internal/repository/mysql"
	"github.com/mehgokalp/insider-project/pkg/log"
	"net/http"
)

type Handler struct {
	ctx        *gin.Context
	logger     log.Logger
	repository pkgDatabaseRepository.MessageRepository
}

// NewHandler godoc
// @Summary List messages
// @Description Get all messages
// @Tags messages
// @Accept  json
// @Produce  json
// @Success 200 {array} Message
// @Router /messages/ [get]
func NewHandler(logger log.Logger, repository pkgDatabaseRepository.MessageRepository) func(*gin.Context) {
	return func(c *gin.Context) {
		h := Handler{
			ctx:        c,
			logger:     logger,
			repository: repository,
		}

		h.Handle()
	}
}

func (h *Handler) Handle() {
	list, err := h.repository.List()
	if err != nil {
		h.ctx.JSON(http.StatusInternalServerError, err)
		return
	}

	h.ctx.JSON(http.StatusOK, mapMessageList(list))
}
