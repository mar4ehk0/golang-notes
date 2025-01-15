package handler

import (
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
	sessionName   = "mysession"
	userIdCtx     = "userId"
)

func New(router *gin.Engine, service *service.Service) *Handler {
	router.LoadHTMLGlob("templates/**/*")

	router.Static("/img", "./static/img")
	router.Static("/css", "./static/css")
	router.StaticFile("/favicon.ico", "./static/img/favicon.ico")

	store := cookie.NewStore([]byte("secret"))
	router.Use(sessions.Sessions(sessionName, store))

	return &Handler{router: router, services: service}
}

func (h *Handler) InitRoutes() *gin.Engine {
	h.router.GET("/", h.renderHomePage)
	h.router.GET("/403", h.render403)
	h.router.GET("/404", h.render404)

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
		workspace.GET("/notes", h.renderNoteList)
		workspace.GET("/notes/create", h.renderFormNoteCreate)
		workspace.POST("/notes", h.processFormNoteCreate)
		workspace.GET("/notes/:id", h.renderNote)
		workspace.GET("/notes/:id/update", h.renderNoteUpdate)
		workspace.POST("/notes/:id", h.processFormNoteUpdate)
	}
	return h.router
}

func AuthRequired(c *gin.Context) {
	session := sessions.Default(c)
	userId := session.Get(authenticated)
	if userId == nil {
		logrus.Errorf("userID is empty")
		c.Abort()
		return
	}
	c.Set(userIdCtx, userId)
	c.Next()
}

func getItemFromSession(s *sessions.Session, key string) interface{} {
	session := *s

	value := session.Get(key)
	session.Delete(key)
	if err := session.Save(); err != nil {
		logrus.Errorf("failed get session: %s", err.Error())
	}

	return value
}

func saveItemToSession(s *sessions.Session, key string, value interface{}) {
	session := *s

	session.Set(key, value)
	if err := session.Save(); err != nil {
		logrus.Errorf("failed save session: %s", err.Error())
	}
}
