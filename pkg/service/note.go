package service

import (
	"fmt"

	"github.com/mar4ehk0/notes/model"
	"github.com/mar4ehk0/notes/pkg/dto"
	"github.com/mar4ehk0/notes/pkg/repository"
)

type NoteService struct {
	repo repository.Note
}

func NewNoteService(repo repository.Note) *NoteService {
	return &NoteService{repo: repo}
}

func (s *NoteService) CreateNote(userId int, input dto.NoteCreateDto) (model.Note, error) {
	note, err := s.repo.AddNote(userId, input)
	if err != nil {
		return model.Note{}, fmt.Errorf("repo add note: %w", err)
	}

	return note, nil
}

func (s *NoteService) GetNote(userId int, noteId int) (model.Note, error) {
	note, err := s.repo.GetNoteByID(noteId)
	if err != nil {
		return model.Note{}, fmt.Errorf("repo get note by noteID{%v}: %w", noteId, err)
	}
	if note.UserID != userId {
		return model.Note{}, ErrForbidden
	}

	return note, nil
}

func (s *NoteService) GetNotes(userId int) ([]model.Note, error) {
	notes, err := s.repo.GetNotesByUserId(userId)
	if err != nil {
		return make([]model.Note, 0), fmt.Errorf("repo get notes by userID: %w", err)
	}

	return notes, nil
}
