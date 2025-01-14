package handler

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/mar4ehk0/notes/pkg/dto"
	"github.com/mar4ehk0/notes/pkg/service"
	"github.com/sirupsen/logrus"
)

func (h *Handler) renderNoteList(c *gin.Context) {
	session := sessions.Default(c)
	userID := c.GetInt(userIdCtx)

	notes, err := h.services.Note.GetNotes(userID)
	if err != nil {
		logrus.Errorf("render note list: get notes: %s", err.Error())

		saveItemToSession(&session, flashError, "Something went wrong")
		c.Redirect(http.StatusFound, "/") // 500
		return
	}

	c.HTML(http.StatusOK, "note/list.tmpl", gin.H{
		"Notes": notes,
	})
}

func (h *Handler) renderNoteItem(c *gin.Context) {
	session := sessions.Default(c)
	userID := c.GetInt(userIdCtx)
	noteID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		logrus.Errorf("render note item: atoi: %s", err.Error())

		saveItemToSession(&session, flashError, "Something went wrong")
		c.Redirect(http.StatusFound, "/workspace/notes/list")
		return
	}

	note, err := h.services.Note.GetNote(userID, noteID)
	if err != nil {
		logrus.Errorf("render note item: get note: %s", err.Error())

		if errors.Is(err, service.ErrForbidden) {
			c.Redirect(http.StatusFound, "/403")
			return
		}
		saveItemToSession(&session, flashError, "Something went wrong")
		c.Redirect(http.StatusFound, "/workspace/notes/list")
		return
	}

	c.HTML(http.StatusOK, "note/item.tmpl", gin.H{
		"Title": note.Title,
		"Body":  note.Body,
	})
}

func (h *Handler) renderFormNoteCreate(c *gin.Context) {
	c.HTML(http.StatusOK, "note/create.tmpl", gin.H{})
}

func (h *Handler) processFormNoteCreate(c *gin.Context) {
	session := sessions.Default(c)

	userId := c.GetInt(userIdCtx)

	var input dto.NoteCreateDto

	if err := c.ShouldBind(&input); err != nil {
		saveItemToSession(&session, flashError, "Title and Body are required")
		c.Redirect(http.StatusFound, "/workspace/note/create")
		return
	}

	note, err := h.services.Note.CreateNote(userId, input)
	if err != nil {
		logrus.Errorf("process form note create: %s", err.Error())

		saveItemToSession(&session, flashError, "Something went wrong")
		c.Redirect(http.StatusFound, "/workspace/note/create")
		return
	}

	c.Redirect(http.StatusFound, fmt.Sprintf("/workspace/notes/%d", note.ID))
}
