package handler

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/mar4ehk0/notes/model"
	"github.com/mar4ehk0/notes/pkg/dto"
	"github.com/sirupsen/logrus"
)

func (h *Handler) renderNoteList(c *gin.Context) {
	session := sessions.Default(c)
	userID := c.GetInt(userIDCtx)

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
	msg := "render note item"
	note := h.getNote(c, msg)

	tags, err := h.services.Tag.GetTagsByNoteID(note.ID)
	if err != nil {
		logrus.Errorf("%s: get tags by node id: %s", msg, err.Error())

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
		"Tags":  tags,
	})
}

func (h *Handler) renderFormNoteCreate(c *gin.Context) {
	session := sessions.Default(c)
	errMsg := getItemFromSession(&session, flashError)
	infoMsg := getItemFromSession(&session, flashInfo)

	tags, err := h.services.Tag.GetTags()
	if err != nil {
		logrus.Errorf("render form note create: get tags: %s", err.Error())

		saveItemToSession(&session, flashError, "Something went wrong")
		c.Redirect(http.StatusFound, "/workspace/notes")
		return
	}

	c.HTML(http.StatusOK, "note/create.tmpl", gin.H{
		"Error": errMsg,
		"Info":  infoMsg,
		"Tags":  tags,
	})
}

func (h *Handler) processFormNoteCreate(c *gin.Context) {
	session := sessions.Default(c)

	userID := c.GetInt(userIDCtx)

	var input dto.NoteDto

	if err := c.ShouldBind(&input); err != nil {
		logrus.Printf("%v \n", err)

		saveItemToSession(&session, flashError, "Title and Body are required")
		c.Redirect(http.StatusFound, "/workspace/notes/create")
		return
	}

	noteID, err := h.services.Note.CreateNote(userID, input)
	if err != nil {
		logrus.Errorf("process form note create: create note: %s", err.Error())

		saveItemToSession(&session, flashError, "Something went wrong")
		c.Redirect(http.StatusFound, "/workspace/notes/create")
		return
	}

	c.Redirect(http.StatusFound, fmt.Sprintf("/workspace/notes/%d", noteID))
}

func (h *Handler) renderFormNoteUpdate(c *gin.Context) {
	session := sessions.Default(c)
	msg := "render form note update"

	note := h.getNote(c, msg)

	tags, err := h.services.Tag.GetTagsWithTaggedByNoteID(note.ID)
	if err != nil {
		logrus.Errorf("%s: get tags with tagged by note id: %s", msg, err.Error())

		saveItemToSession(&session, flashError, "Something went wrong")
		c.Redirect(http.StatusFound, "/workspace/notes")
		c.Abort()
	}

	errMsg := getItemFromSession(&session, flashError)
	infoMsg := getItemFromSession(&session, flashInfo)

	c.HTML(http.StatusOK, "note/update.tmpl", gin.H{
		"ID":    note.ID,
		"Title": note.Title,
		"Body":  note.Body,
		"Error": errMsg,
		"Info":  infoMsg,
		"Tags":  tags,
	})
}

func (h *Handler) processFormNoteUpdate(c *gin.Context) {
	session := sessions.Default(c)

	userID := c.GetInt(userIDCtx)
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

func (h *Handler) renderNoteDelete(c *gin.Context) {
	session := sessions.Default(c)
	userID := c.GetInt(userIDCtx)
	noteID := h.getParamInt("id", c)

	note, err := h.services.Note.GetNote(userID, noteID)
	if err != nil {
		logrus.Errorf("render note delete: get note: %s", err.Error())

		checkError(err, c)

		saveItemToSession(&session, flashError, "Something went wrong")
		// создать один метод c Abort()
		c.Redirect(http.StatusFound, "/workspace/notes")
		return
	}

	errMsg := getItemFromSession(&session, flashError)
	infoMsg := getItemFromSession(&session, flashInfo)

	c.HTML(http.StatusOK, "note/delete.tmpl", gin.H{
		"ID":    note.ID,
		"Title": note.Title,
		"Error": errMsg,
		"Info":  infoMsg,
	})
}

func (h *Handler) processNoteDelete(c *gin.Context) {
	session := sessions.Default(c)
	userID := c.GetInt(userIDCtx)
	noteID := h.getParamInt("id", c)

	err := h.services.Note.DeleteNote(userID, noteID)
	if err != nil {
		logrus.Errorf("process form note delete: delete note: %s", err.Error())

		checkError(err, c)

		saveItemToSession(&session, flashError, "Something went wrong")
		c.Redirect(http.StatusFound, fmt.Sprintf("/workspace/notes/%d/delete", noteID))
		return
	}

	saveItemToSession(&session, flashInfo, fmt.Sprintf("note: %d was delete", noteID))
	c.Redirect(http.StatusFound, "/workspace/notes")
}

func (h *Handler) getNote(c *gin.Context, msg string) model.Note {
	session := sessions.Default(c)
	userID := c.GetInt(userIDCtx)
	noteID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		logrus.Errorf("%s: atoi: %s", msg, err.Error())

		saveItemToSession(&session, flashError, "Something went wrong")
		c.Redirect(http.StatusFound, "/workspace/notes")
		c.Abort()
	}

	note, err := h.services.Note.GetNote(userID, noteID)
	if err != nil {
		logrus.Errorf("%s: get note: %s", msg, err.Error())

		checkError(err, c)

		saveItemToSession(&session, flashError, "Something went wrong")
		c.Redirect(http.StatusFound, "/workspace/notes")
		c.Abort()
	}

	return note
}
