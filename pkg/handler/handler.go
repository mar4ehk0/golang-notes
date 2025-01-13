package handler

import (
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/mar4ehk0/notes/pkg/service"
	"github.com/sirupsen/logrus"
)

type Handler struct {
	router   *gin.Engine
	services *service.Service
}

const (
	flashError    = "Error"
	flashInfo     = "Info"
	authenticated = "authenticated"
)

func New(router *gin.Engine, service *service.Service) *Handler {
	router.LoadHTMLGlob("templates/**/*")

	router.Static("/img", "./static/img")
	router.Static("/css", "./static/css")
	router.StaticFile("/favicon.ico", "./static/img/favicon.ico")

	store := cookie.NewStore([]byte("secret"))
	router.Use(sessions.Sessions("mysession", store))

	return &Handler{router: router, services: service}
}

func (h *Handler) InitRoutes() *gin.Engine {

	h.router.GET("/", h.renderHomePage)
	auth := h.router.Group("/auth")
	{
		auth.GET("/sign-in", h.renderFormSignIn)
		auth.POST("/sign-in", h.processFormSignIn)

		auth.GET("/sign-up", h.renderFormSignUp)
		auth.POST("/sign-up", h.processFormSignUp)
	}

	workspace := h.router.Group("/workspace")
	workspace.Use(AuthRequired)
	{
		// workspace.GET("/notes/list", h.)
		workspace.GET("/notes/create", h.renderFormNoteCreate)
		// workspace.POST("/notes", h.)
		// workspace.GET("/notes/:id", h.)
	}
	return h.router
}

func AuthRequired(c *gin.Context) {
	session := sessions.Default(c)
	user := session.Get(authenticated)
	if user == nil {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	c.Next()
}

func getItemFromSession(s *sessions.Session, key string) interface{} {
	session := *s

	value := session.Get(key)
	session.Delete(key)
	err := session.Save()
	if err != nil {
		logrus.Error("failed get session: %s", err.Error())
	}

	return value
}

func saveItemToSession(s *sessions.Session, key string, value interface{}) {
	session := *s

	session.Set(key, value)
	err := session.Save()
	if err != nil {
		logrus.Error("failed get session: %s", err.Error())
	}
}
