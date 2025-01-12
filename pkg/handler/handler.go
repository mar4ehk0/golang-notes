package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/mar4ehk0/notes/pkg/service"
)

type Handler struct {
	router  *gin.Engine
	service *service.Service
}

func New(router *gin.Engine, service *service.Service) *Handler {
	router.LoadHTMLGlob("templates/**/*")

	router.Static("/img", "./static/img")
	router.Static("/css", "./static/css")
	router.StaticFile("/favicon.ico", "./static/img/favicon.ico")

	return &Handler{router: router, service: service}
}

func (h *Handler) InitRoutes() *gin.Engine {
	auth := h.router.Group("/auth")
	{
		auth.GET("/sign-in", h.renderFormSignIn)
		auth.POST("/sign-in", h.processFormSignIn)
	}

	return h.router
}
