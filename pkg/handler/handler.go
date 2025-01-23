package handler

import (
	"net/http"
	"strconv"

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
	userIDCtx     = "userId"
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
		workspace.GET("/notes/:id/update", h.renderFormNoteUpdate)
		workspace.POST("/notes/:id", h.processFormNoteUpdate)
		workspace.GET("/notes/:id/delete", h.renderNoteDelete)
		workspace.POST("/notes/:id/delete", h.processNoteDelete)
	}
	return h.router
}

func AuthRequired(c *gin.Context) {
	session := sessions.Default(c)
	userID := session.Get(authenticated)
	if userID == nil {
		logrus.Errorf("userID is empty")
		c.Abort()
		return
	}
	c.Set(userIDCtx, userID)
	c.Next()
}

func (h *Handler) getParamIDInt(c *gin.Context) int {
	param, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		logrus.Errorf("atoi: %s", err.Error())

		h.saveItemToSession(c, flashError, "Something went wrong")
		c.Redirect(http.StatusFound, "/workspace/notes")
		return -1
	}

	return param
}

func (h *Handler) getItemFromSession(c *gin.Context, key string) interface{} {
	session := sessions.Default(c)

	value := session.Get(key)
	session.Delete(key)
	if err := session.Save(); err != nil {
		logrus.Errorf("failed get session: %s", err.Error())
	}

	return value
}

func (h *Handler) saveItemToSession(c *gin.Context, key string, value interface{}) {
	session := sessions.Default(c)

	session.Set(key, value)
	if err := session.Save(); err != nil {
		logrus.Errorf("failed save session: %s", err.Error())
	}
}
