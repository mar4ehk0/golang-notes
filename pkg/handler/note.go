package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *Handler) renderListNote(c *gin.Context) {
	c.HTML(http.StatusOK, "note/list.tmpl", gin.H{})
}

func (h *Handler) renderFormNoteCreate(c *gin.Context) {
	c.HTML(http.StatusOK, "note/create.tmpl", gin.H{})
}
