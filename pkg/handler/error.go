package handler

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mar4ehk0/notes/pkg/repository"
	"github.com/mar4ehk0/notes/pkg/service"
)

func (h *Handler) checkError(err error, c *gin.Context) {
	var notFoundErr *repository.NotFoundError
	var forbiddenErr *service.ForbiddenError

	if errors.As(err, &forbiddenErr) {
		c.Redirect(http.StatusFound, "/403")
		return
	}

	if errors.As(err, &notFoundErr) {
		h.saveItemToSession(c, flashError, notFoundErr.Error())
		c.Redirect(http.StatusFound, "/404")
		return
	}
}
