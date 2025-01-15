package handler

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/mar4ehk0/notes/pkg/dto"
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

	errMsg := getItemFromSession(&session, flashError)
	infoMsg := getItemFromSession(&session, flashInfo)

	c.HTML(http.StatusOK, "note/list.tmpl", gin.H{
		"Notes": notes,
		"Error": errMsg,
		"Info":  infoMsg,
	})
}

func (h *Handler) renderNote(c *gin.Context) {
	session := sessions.Default(c)
	userID := c.GetInt(userIdCtx)
	noteID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		logrus.Errorf("render note item: atoi: %s", err.Error())

		saveItemToSession(&session, flashError, "Something went wrong")
		c.Redirect(http.StatusFound, "/workspace/notes")
		return
	}

	note, err := h.services.Note.GetNote(userID, noteID)
	if err != nil {
		logrus.Errorf("render note item: get note: %s", err.Error())

		checkError(err, c)

		saveItemToSession(&session, flashError, "Something went wrong")
		c.Redirect(http.StatusFound, "/workspace/notes")
		return
	}

	errMsg := getItemFromSession(&session, flashError)
	infoMsg := getItemFromSession(&session, flashInfo)

	c.HTML(http.StatusOK, "note/item.tmpl", gin.H{
		"ID":    note.ID,
		"Title": note.Title,
		"Body":  note.Body,
		"Error": errMsg,
		"Info":  infoMsg,
	})
}

func (h *Handler) renderFormNoteCreate(c *gin.Context) {
	session := sessions.Default(c)
	errMsg := getItemFromSession(&session, flashError)
	infoMsg := getItemFromSession(&session, flashInfo)

	c.HTML(http.StatusOK, "note/create.tmpl", gin.H{
		"Error": errMsg,
		"Info":  infoMsg,
	})
}

func (h *Handler) processFormNoteCreate(c *gin.Context) {
	session := sessions.Default(c)

	userId := c.GetInt(userIdCtx)

	var input dto.NoteDto

	if err := c.ShouldBind(&input); err != nil {
		saveItemToSession(&session, flashError, "Title and Body are required")
		c.Redirect(http.StatusFound, "/workspace/note/create")
		return
	}

	noteID, err := h.services.Note.CreateNote(userId, input)
	if err != nil {
		logrus.Errorf("process form note create: create note: %s", err.Error())

		saveItemToSession(&session, flashError, "Something went wrong")
		c.Redirect(http.StatusFound, "/workspace/note/create")
		return
	}

	c.Redirect(http.StatusFound, fmt.Sprintf("/workspace/notes/%d", noteID))
}

func (h *Handler) renderNoteUpdate(c *gin.Context) {
	session := sessions.Default(c)
	userID := c.GetInt(userIdCtx)
	noteID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		logrus.Errorf("render note item update: atoi: %s", err.Error())

		saveItemToSession(&session, flashError, "Something went wrong")
		c.Redirect(http.StatusFound, "/workspace/notes")
		return
	}

	note, err := h.services.Note.GetNote(userID, noteID)
	if err != nil {
		logrus.Errorf("render note update: get note: %s", err.Error())

		checkError(err, c)

		saveItemToSession(&session, flashError, "Something went wrong")
		c.Redirect(http.StatusFound, "/workspace/notes")
		return
	}

	errMsg := getItemFromSession(&session, flashError)
	infoMsg := getItemFromSession(&session, flashInfo)

	c.HTML(http.StatusOK, "note/update.tmpl", gin.H{
		"ID":    note.ID,
		"Title": note.Title,
		"Body":  note.Body,
		"Error": errMsg,
		"Info":  infoMsg,
	})
}

func (h *Handler) processFormNoteUpdate(c *gin.Context) {
	session := sessions.Default(c)

	userID := c.GetInt(userIdCtx)
	noteID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		logrus.Errorf("process note update: atoi: %s", err.Error())

		saveItemToSession(&session, flashError, "Something went wrong")
		c.Redirect(http.StatusFound, "/workspace/notes")
		return
	}

	var input dto.NoteDto

	if err := c.ShouldBind(&input); err != nil {
		saveItemToSession(&session, flashError, "Title and Body are required")
		c.Redirect(http.StatusFound, fmt.Sprintf("/workspace/notes/%d/update", noteID))
		return
	}

	err = h.services.Note.UpdateNote(userID, noteID, input)
	if err != nil {
		logrus.Errorf("process form note update: update note: %s", err.Error())

		checkError(err, c)

		saveItemToSession(&session, flashError, "Something went wrong")
		c.Redirect(http.StatusFound, fmt.Sprintf("/workspace/notes/%d/update", noteID))
		return
	}

	c.Redirect(http.StatusFound, fmt.Sprintf("/workspace/notes/%d", noteID))
}
