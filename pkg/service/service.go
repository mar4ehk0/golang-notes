package service

import (
	"github.com/mar4ehk0/notes/model"
	"github.com/mar4ehk0/notes/pkg/dto"
	"github.com/mar4ehk0/notes/pkg/repository"
)

type Authorization interface {
	CreateUser(d dto.UserSingUpDto) (model.User, error)
	Authorize(d dto.UserSingInDto) (model.User, bool, error)
}

type Note interface {
	CreateNote(userID int, input dto.NoteDto) (int, error)
	GetNote(userID int, noteId int) (model.Note, error)
	GetNotes(userID int) ([]model.Note, error)
	UpdateNote(userID int, noteID int, input dto.NoteDto) error
	DeleteNote(userID int, noteID int) error
}

type Tag interface {
	GetTags() ([]model.Tag, error)
	GetTagByID(tagID int) (model.Tag, error)
	GetTagsByNoteId(noteID int) ([]model.Tag, error)
	GetTagsWithTaggedByNoteID(noteID int) ([]dto.TagsWithTagged, error)
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
		Tag:           NewTagService(repository),
	}
}
