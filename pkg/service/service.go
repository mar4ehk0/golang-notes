package service

import (
	"errors"

	"github.com/mar4ehk0/notes/model"
	"github.com/mar4ehk0/notes/pkg/dto"
	"github.com/mar4ehk0/notes/pkg/repository"
)

var ErrForbidden = errors.New("forbidden")

type Authorization interface {
	CreateUser(d dto.UserSingUpDto) (model.User, error)
	Authorize(d dto.UserSingInDto) (model.User, bool, error)
}

type Note interface {
	CreateNote(userID int, input dto.NoteCreateDto) (model.Note, error)
	GetNote(userID int, noteId int) (model.Note, error)
	GetNotes(userID int) ([]model.Note, error)
}

type Tag interface {
}

type Service struct {
	Authorization
	Note
	Tag
}

func NewService(repository *repository.Repository) *Service {
	return &Service{
		Authorization: NewAuthService(repository),
		Note:          NewNoteService(repository),
	}
}
