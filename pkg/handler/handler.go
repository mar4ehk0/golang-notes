package handler

import "github.com/gin-gonic/gin"

type Handler struct {
	router *gin.Engine
}

func New(router *gin.Engine) *Handler {
	router.LoadHTMLGlob("templates/**/*")

	router.Static("/img", "./static/img")
	router.Static("/css", "./static/css")
	router.StaticFile("/favicon.ico", "./static/img/favicon.ico")

	return &Handler{router: router}
}

func (h *Handler) InitRoutes() *gin.Engine {
	auth := h.router.Group("/auth")
	{
		auth.GET("/sign-in", h.renderFormSignIn)
		auth.POST("/sign-in", h.processFormSignIn)
	}

	return h.router
}
