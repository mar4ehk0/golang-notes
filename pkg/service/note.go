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

func (s *NoteService) CreateNote(userID int, input dto.NoteDto) (int, error) {
	var noteID int
	var err error

	if len(input.TagsID) > 0 {
		noteID, err = s.repo.AddNoteWithTag(userID, input)
		if err != nil {
			return 0, fmt.Errorf("repo add note with tag: %w", err)
		}
	} else {
		noteID, err = s.repo.AddNote(userID, input)
		if err != nil {
			return 0, fmt.Errorf("repo add note: %w", err)
		}
	}

	return noteID, nil
}

func (s *NoteService) GetNote(userID int, noteID int) (model.Note, error) {
	note, err := s.repo.GetNoteByID(noteID)
	if err != nil {
		return model.Note{}, fmt.Errorf("repo get note by noteID{%v}: %w", noteID, err)
	}
	if note.UserID != userID {
		return model.Note{}, NewForbiddenError(userID, noteID)
	}

	return note, nil
}

func (s *NoteService) GetNotes(userID int) ([]model.Note, error) {
	notes, err := s.repo.GetNotesByUserID(userID)
	if err != nil {
		return make([]model.Note, 0), fmt.Errorf("repo get notes by userID: %w", err)
	}

	return notes, nil
}

func (s *NoteService) UpdateNote(userID int, noteID int, input dto.NoteDto) error {
	note, err := s.repo.GetNoteByID(noteID)
	if err != nil {
		return fmt.Errorf("repo get note by noteID{%v}: %w", noteID, err)
	}
	if note.UserID != userID {
		return NewForbiddenError(userID, noteID)
	}

	err = s.repo.UpdateNote(noteID, input)
	if err != nil {
		return fmt.Errorf("repo update note noteID{%v}: %w", noteID, err)
	}

	return nil
}

func (s *NoteService) DeleteNote(userID int, noteID int) error {
	note, err := s.repo.GetNoteByID(noteID)
	if err != nil {
		return fmt.Errorf("repo get note by noteID{%v}: %w", noteID, err)
	}
	if note.UserID != userID {
		return NewForbiddenError(userID, noteID)
	}

	isDeleted, err := s.repo.DeleteNote(noteID)
	if err != nil {
		return fmt.Errorf("repo delete note by noteID{%v}: %w", noteID, err)
	}

	if !isDeleted {
		return fmt.Errorf("does not delete note by noteID{%v}: %w", noteID, err)
	}

	return nil
}
