package list

import (
	"github.com/gin-gonic/gin"
	pkgDatabaseRepository "github.com/mehgokalp/insider-project/pkg/database/repository"
	"github.com/mehgokalp/insider-project/pkg/log"
	"net/http"
)

type Handler struct {
	ctx        *gin.Context
	logger     log.Logger
	repository *pkgDatabaseRepository.MessageRepository
}

func NewHandler(logger log.Logger, repository *pkgDatabaseRepository.MessageRepository) func(*gin.Context) {
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
