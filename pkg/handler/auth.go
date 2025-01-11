package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *Handler) renderFormSignIn(c *gin.Context) {
	c.HTML(http.StatusOK, "auth/sign_in.tmpl", gin.H{})
}

func (h *Handler) processFormSignIn(c *gin.Context) {
	
}
