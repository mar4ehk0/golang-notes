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
	errMsg := h.getItemFromSession(c, flashError)
	infoMsg := h.getItemFromSession(c, flashInfo)

	c.HTML(http.StatusOK, "auth/sign_up.tmpl", gin.H{
		"Error": errMsg,
		"Info":  infoMsg,
	})
}

func (h *Handler) processFormSignUp(c *gin.Context) {
	urlRedirect := "/auth/sign-up"
	var input dto.UserSingUpDto

	if err := c.ShouldBind(&input); err != nil {
		h.saveItemToSession(c, flashError, "All fields are required.")
		c.Redirect(http.StatusFound, urlRedirect)
		return
	}

	err := input.Validate()
	if err != nil {
		h.saveItemToSession(c, flashError, err.Error())
		c.Redirect(http.StatusFound, urlRedirect)
		return
	}

	user, err := h.services.Authorization.CreateUser(input)
	if err != nil {
		logrus.Errorf("process form sign-up: create user: %s", err.Error())

		msg := "Something went wrong"
		if errors.Is(err, repository.ErrDBDuplicateKey) {
			msg = fmt.Sprintf("User already exist with same email: %s", input.Email)
		}

		h.saveItemToSession(c, flashError, msg)
		c.Redirect(http.StatusFound, urlRedirect)
		return
	}

	h.saveItemToSession(c, flashInfo, fmt.Sprintf("User created: %s", user.Email))
	c.Redirect(http.StatusFound, "/auth/sign-in")
}

func (h *Handler) renderFormSignIn(c *gin.Context) {
	// session := sessions.Default(c)

	errMsg := h.getItemFromSession(c, flashError)
	infoMsg := h.getItemFromSession(c, flashInfo)
	email := h.getItemFromSession(c, emailFormSignIn)

	c.HTML(http.StatusOK, "auth/sign_in.tmpl", gin.H{
		"Error": errMsg,
		"Info":  infoMsg,
		"Email": email,
	})
}

func (h *Handler) processFormSignIn(c *gin.Context) {
	session := sessions.Default(c)
	urlRedirect := "/auth/sign-up"
	var input dto.UserSingInDto

	if err := c.ShouldBind(&input); err != nil {
		h.saveItemToSession(c, flashError, "Email and Password are required")
		c.Redirect(http.StatusFound, urlRedirect)
		return
	}

	user, canAuthorize, err := h.services.Authorization.Authorize(input)
	if err != nil {
		logrus.Errorf("process form sign-in: can authorize: %s", err.Error())

		h.saveItemToSession(c, flashError, "Something went wrong")
		c.Redirect(http.StatusFound, urlRedirect)
		return
	}

	if !canAuthorize {
		h.saveItemToSession(c, flashError, "Email or password wrong")
		c.Redirect(http.StatusFound, urlRedirect)
		return
	}

	session.Set(authenticated, user.ID)
	if err := session.Save(); err != nil {
		logrus.Errorf("process form sign-in: save session: %s", err.Error())

		h.saveItemToSession(c, flashError, "Something went wrong")
		c.Redirect(http.StatusFound, urlRedirect)
		return
	}

	c.Redirect(http.StatusFound, "/workspace/notes/create")
}
