package handler

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mar4ehk0/notes/model"
	"github.com/mar4ehk0/notes/pkg/dto"
	"github.com/sirupsen/logrus"
)

func (h *Handler) renderNoteList(c *gin.Context) {
	userID := c.GetInt(userIDCtx)

	notes, err := h.services.Note.GetNotes(userID)
	if err != nil {
		logrus.Errorf("render note list: get notes: %s", err.Error())

		h.saveItemToSession(c, flashError, "Something went wrong")
		c.Redirect(http.StatusFound, "/workspace/notes/create")
		return
	}

	errMsg := h.getItemFromSession(c, flashError)
	infoMsg := h.getItemFromSession(c, flashInfo)

	c.HTML(http.StatusOK, "note/list.tmpl", gin.H{
		"Notes": notes,
		"Error": errMsg,
		"Info":  infoMsg,
	})
}

func (h *Handler) renderNote(c *gin.Context) {
	msg := "render note item"
	note := h.getNote(c, msg)

	tags, err := h.services.Tag.GetTagsByNoteID(note.ID)
	if err != nil {
		logrus.Errorf("%s: get tags by node id: %s", msg, err.Error())

		h.saveItemToSession(c, flashError, "Something went wrong")
		c.Redirect(http.StatusFound, "/workspace/notes")
		return
	}

	errMsg := h.getItemFromSession(c, flashError)
	infoMsg := h.getItemFromSession(c, flashInfo)

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
	errMsg := h.getItemFromSession(c, flashError)
	infoMsg := h.getItemFromSession(c, flashInfo)

	tags, err := h.services.Tag.GetTags()
	if err != nil {
		logrus.Errorf("render form note create: get tags: %s", err.Error())

		h.saveItemToSession(c, flashError, "Something went wrong")
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
	userId := c.GetInt(userIDCtx)

	var input dto.NoteDto

	if err := c.ShouldBind(&input); err != nil {
		logrus.Printf("%v \n", err)

		h.saveItemToSession(c, flashError, "Title and Body are required")
		c.Redirect(http.StatusFound, "/workspace/notes/create")
		return
	}

	noteID, err := h.services.Note.CreateNote(userId, input)
	if err != nil {
		logrus.Errorf("process form note create: create note: %s", err.Error())

		h.saveItemToSession(c, flashError, "Something went wrong")
		c.Redirect(http.StatusFound, "/workspace/notes/create")
		return
	}

	c.Redirect(http.StatusFound, fmt.Sprintf("/workspace/notes/%d", noteID))
}

func (h *Handler) renderFormNoteUpdate(c *gin.Context) {
	msg := "render form note update"

	note := h.getNote(c, msg)

	tags, err := h.services.Tag.GetTagsWithTaggedByNoteID(note.ID)
	if err != nil {
		logrus.Errorf("%s: get tags with tagged by note id: %s", msg, err.Error())

		h.saveItemToSession(c, flashError, "Something went wrong")
		c.Redirect(http.StatusFound, "/workspace/notes")
		c.Abort()
	}

	errMsg := h.getItemFromSession(c, flashError)
	infoMsg := h.getItemFromSession(c, flashInfo)

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
	userID := c.GetInt(userIDCtx)
	noteID := h.getParamIDInt(c)

	var input dto.NoteDto

	if err := c.ShouldBind(&input); err != nil {
		h.saveItemToSession(c, flashError, "Title and Body are required")
		c.Redirect(http.StatusFound, fmt.Sprintf("/workspace/notes/%d/update", noteID))
		return
	}

	err := h.services.Note.UpdateNote(userID, noteID, input)
	if err != nil {
		logrus.Errorf("process form note update: update note: %s", err.Error())

		h.checkError(err, c)

		h.saveItemToSession(c, flashError, "Something went wrong")
		c.Redirect(http.StatusFound, fmt.Sprintf("/workspace/notes/%d/update", noteID))
		return
	}

	c.Redirect(http.StatusFound, fmt.Sprintf("/workspace/notes/%d", noteID))
}

func (h *Handler) renderNoteDelete(c *gin.Context) {
	userID := c.GetInt(userIDCtx)
	noteID := h.getParamIDInt(c)

	note, err := h.services.Note.GetNote(userID, noteID)
	if err != nil {
		logrus.Errorf("render note delete: get note: %s", err.Error())

		h.checkError(err, c)

		h.saveItemToSession(c, flashError, "Something went wrong")
		c.Redirect(http.StatusFound, "/workspace/notes")
		return
	}

	errMsg := h.getItemFromSession(c, flashError)
	infoMsg := h.getItemFromSession(c, flashInfo)

	c.HTML(http.StatusOK, "note/delete.tmpl", gin.H{
		"ID":    note.ID,
		"Title": note.Title,
		"Error": errMsg,
		"Info":  infoMsg,
	})
}

func (h *Handler) processNoteDelete(c *gin.Context) {
	userID := c.GetInt(userIDCtx)
	noteID := h.getParamIDInt(c)

	err := h.services.Note.DeleteNote(userID, noteID)
	if err != nil {
		logrus.Errorf("process form note delete: delete note: %s", err.Error())

		h.checkError(err, c)

		h.saveItemToSession(c, flashError, "Something went wrong")
		c.Redirect(http.StatusFound, fmt.Sprintf("/workspace/notes/%d/delete", noteID))
		return
	}

	h.saveItemToSession(c, flashInfo, fmt.Sprintf("note: %d was delete", noteID))
	c.Redirect(http.StatusFound, "/workspace/notes")
}

func (h *Handler) getNote(c *gin.Context, msg string) model.Note {
	userID := c.GetInt(userIDCtx)
	noteID := h.getParamIDInt(c)

	note, err := h.services.Note.GetNote(userID, noteID)
	if err != nil {
		logrus.Errorf("%s: get note: %s", msg, err.Error())

		h.checkError(err, c)

		h.saveItemToSession(c, flashError, "Something went wrong")
		c.Redirect(http.StatusFound, "/workspace/notes")
		return model.Note{}
	}

	return note
}
