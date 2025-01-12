package handler

import (
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/mar4ehk0/notes/model"
	"github.com/sirupsen/logrus"
)

const (
	errorFormSignIn = "error_form_sign__in"
	emailFormSignIn = "email_form_sign__in"
)

func (h *Handler) renderFormSignIn(c *gin.Context) {
	session := sessions.Default(c)

	errorMessage := session.Get(errorFormSignIn)
	session.Delete(errorFormSignIn)
	email := session.Get(emailFormSignIn)
	session.Delete(emailFormSignIn)
	err := session.Save()
	if err != nil {
		logrus.Fatalf("failed session save: %s", err.Error())
	}

	c.HTML(http.StatusOK, "auth/sign_in.tmpl", gin.H{
		"Error": errorMessage,
		"Email": email,
	})
}

func (h *Handler) processFormSignIn(c *gin.Context) {
	session := sessions.Default(c)

	var input model.User

	// logrus.Printf("input %v", c.Request.Body)

	// body, _ := io.ReadAll(c.Request.Body)
	// logrus.Printf("Raw body: %s", body)

	if err := c.ShouldBind(&input); err != nil {
		session.Set(errorFormSignIn, "Email and Password are required")
		session.Save()
		c.Redirect(http.StatusFound, "/auth/sign-in")
		return
	}
	/// процесс логина
	logrus.Println("ok")
}
