package repository

import (
	"github.com/jmoiron/sqlx"
	"github.com/mar4ehk0/notes/model"
	"github.com/mar4ehk0/notes/pkg/dto"
)

type Authorization interface {
	CreateUser(user dto.UserSingUpDto) (int, error)
	GetUserByEmail(email string) (model.User, error)
}

type Note interface {
	AddNote(userID int, input dto.NoteDto) (int, error)
	AddNoteWithTag(userID int, input dto.NoteDto) (int, error)
	GetNoteByID(noteID int) (model.Note, error)
	GetNotesByUserID(userID int) ([]model.Note, error)
	UpdateNote(noteID int, input dto.NoteDto) error
	DeleteNote(noteID int) (bool, error)
}

type Tag interface {
	GetTags() ([]model.Tag, error)
	GetTagByID(tagID int) (model.Tag, error)
	GetTagsByNoteID(noteID int) ([]model.Tag, error)
}

type Repository struct {
	Authorization
	Note
	Tag
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		Authorization: NewAuthPostgres(db),
		Note:          NewNotePostgres(db),
		Tag:           NewTagPostgres(db),
	}
}
