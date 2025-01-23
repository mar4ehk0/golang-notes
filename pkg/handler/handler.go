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

	prefixUrlAuth = "/auth"
	signIn        = "/sign-in"
	urlSignIn     = prefixUrlAuth + signIn
	signUp        = "/sign-up"
	urlSignUp     = prefixUrlAuth + signUp

	prefixUrlWorkspace = "/workspace"
	notes              = "/notes"
	urlNotes           = prefixUrlWorkspace + notes
	notesCreate        = "/notes/create"
	urlNotesCreate     = prefixUrlWorkspace + notesCreate
	notesID            = "/notes/:id"
	urlNotesID         = prefixUrlWorkspace + notesID
	notesIDUpdate      = notesID + "/update"
	urlNotesIDUpdate   = prefixUrlWorkspace + notesIDUpdate
	notesIDDelete      = notesID + "/delete"
	urlNotesIDDelete   = prefixUrlWorkspace + notesIDDelete
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

	auth := h.router.Group(prefixUrlAuth)
	{
		auth.GET(signIn, h.renderFormSignIn)
		auth.POST(signIn, h.processFormSignIn)

		auth.GET(signUp, h.renderFormSignUp)
		auth.POST(signUp, h.processFormSignUp)
	}

	workspace := h.router.Group(prefixUrlWorkspace)
	workspace.Use(AuthRequired)
	{
		workspace.GET(notes, h.renderNoteList)
		workspace.GET(notesCreate, h.renderFormNoteCreate)
		workspace.POST(notes, h.processFormNoteCreate)
		workspace.GET(notesID, h.renderNote)
		workspace.GET(notesIDUpdate, h.renderFormNoteUpdate)
		workspace.POST(notesID, h.processFormNoteUpdate)
		workspace.GET(notesIDDelete, h.renderNoteDelete)
		workspace.POST(notesIDDelete, h.processNoteDelete)
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

func (h *Handler) getParamInt(key string, c *gin.Context) int {
	session := sessions.Default(c)
	param, err := strconv.Atoi(c.Param(key))
	if err != nil {
		logrus.Errorf("atoi: %s", err.Error())

		saveItemToSession(&session, flashError, "Something went wrong")
		c.Redirect(http.StatusFound, "/workspace/notes")
		c.Abort()
	}

	return param
}

func (h *Handler) RedirectAndAbort(c *gin.Context, url string) {
	logrus.Printf("%p", c)
	c.Redirect(http.StatusFound, url)
	c.Abort()
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
