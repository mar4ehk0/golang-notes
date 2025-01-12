package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *Handler) renderHomePage(c *gin.Context) {
	c.HTML(http.StatusOK, "homepage/index.tmpl", nil)
}
