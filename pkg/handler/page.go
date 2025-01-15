package handler

import (
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func (h *Handler) renderHomePage(c *gin.Context) {
	c.HTML(http.StatusOK, "page/homepage.tmpl", nil)
}

func (h *Handler) render403(c *gin.Context) {
	c.HTML(http.StatusOK, "page/403.tmpl", nil)
}

func (h *Handler) render404(c *gin.Context) {
	session := sessions.Default(c)
	info := getItemFromSession(&session, flashError)

	c.HTML(http.StatusOK, "page/404.tmpl", gin.H{
		"Info": info,
	})
}
