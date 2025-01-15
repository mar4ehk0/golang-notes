package handler

import (
	"errors"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/mar4ehk0/notes/pkg/repository"
	"github.com/mar4ehk0/notes/pkg/service"
)

func checkError(err error, c *gin.Context) {
	session := sessions.Default(c)
	var notFoundErr *repository.NotFoundError
	var forbiddenErr *service.ForbiddenError

	if errors.As(err, &forbiddenErr) {
		c.Redirect(http.StatusFound, "/403")
		return
	}

	if errors.As(err, &notFoundErr) {
		saveItemToSession(&session, flashError, notFoundErr.Error())
		c.Redirect(http.StatusFound, "/404")
		return
	}
}
