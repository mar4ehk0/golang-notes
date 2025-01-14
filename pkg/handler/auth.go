package handler

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/mar4ehk0/notes/pkg/dto"
	"github.com/mar4ehk0/notes/pkg/repository"
	"github.com/sirupsen/logrus"
)

const (
	emailFormSignIn = "email_form_sign_in"
)

func (h *Handler) renderFormSignUp(c *gin.Context) {
	session := sessions.Default(c)

	errMsg := getItemFromSession(&session, flashError)
	infoMsg := getItemFromSession(&session, flashInfo)

	c.HTML(http.StatusOK, "auth/sign_up.tmpl", gin.H{
		"Error": errMsg,
		"Info":  infoMsg,
	})
}

func (h *Handler) processFormSignUp(c *gin.Context) {
	session := sessions.Default(c)

	var input dto.UserSingUpDto

	if err := c.ShouldBind(&input); err != nil {
		saveItemToSession(&session, flashError, "All fields are required.")
		c.Redirect(http.StatusFound, "/auth/sign-up")
		return
	}

	err := input.Validate()
	if err != nil {
		saveItemToSession(&session, flashError, err.Error())
		c.Redirect(http.StatusFound, "/auth/sign-up")
		return
	}

	user, err := h.services.Authorization.CreateUser(input)
	if err != nil {
		logrus.Errorf("process form sign-up: create user: %s", err.Error())

		msg := "Something went wrong"
		if errors.Is(err, repository.ErrDBDuplicateKey) {
			msg = fmt.Sprintf("User already exist with same email: %s", input.Email)
		}

		saveItemToSession(&session, flashError, msg)
		c.Redirect(http.StatusFound, "/auth/sign-up")
		return
	}

	saveItemToSession(&session, flashInfo, fmt.Sprintf("User created - %s", user.Email))
	c.Redirect(http.StatusFound, "/auth/sign-in")
}

func (h *Handler) renderFormSignIn(c *gin.Context) {
	session := sessions.Default(c)

	errMsg := getItemFromSession(&session, flashError)
	infoMsg := getItemFromSession(&session, flashInfo)
	email := getItemFromSession(&session, emailFormSignIn)

	c.HTML(http.StatusOK, "auth/sign_in.tmpl", gin.H{
		"Error": errMsg,
		"Info":  infoMsg,
		"Email": email,
	})
}

func (h *Handler) processFormSignIn(c *gin.Context) {
	session := sessions.Default(c)

	var input dto.UserSingInDto

	if err := c.ShouldBind(&input); err != nil {
		saveItemToSession(&session, flashError, "Email and Password are required")
		c.Redirect(http.StatusFound, "/auth/sign-in")
		return
	}

	user, canAuthorize, err := h.services.Authorization.Authorize(input)
	if err != nil {
		logrus.Errorf("process form sign-in: can authorize: %s", err.Error())

		saveItemToSession(&session, flashError, "Something went wrong")
		c.Redirect(http.StatusFound, "/auth/sign-in")
		return
	}

	if !canAuthorize {
		saveItemToSession(&session, flashError, "Email or password wrong")
		c.Redirect(http.StatusFound, "/auth/sign-in")
		return
	}

	session.Set(authenticated, user.ID)
	if err := session.Save(); err != nil {
		logrus.Errorf("process form sign-in: save session: %s", err.Error())

		saveItemToSession(&session, flashError, "Something went wrong")
		c.Redirect(http.StatusFound, "/auth/sign-in")
		return
	}

	c.Redirect(http.StatusFound, "/workspace/notes")
}
