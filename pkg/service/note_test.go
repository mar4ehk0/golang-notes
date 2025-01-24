package service

import (
	"errors"
	"testing"

	"github.com/mar4ehk0/notes/mocks"
	"github.com/mar4ehk0/notes/model"
	"github.com/mar4ehk0/notes/pkg/dto"
	"github.com/mar4ehk0/notes/pkg/repository"
	"github.com/stretchr/testify/assert"
	gomock "go.uber.org/mock/gomock"
)

func TestCanCreateNote(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockNote(ctrl)
	noteService := NewNoteService(mockRepo)

	tests := []struct {
		name           string
		userID         int
		dto            dto.NoteDto
		expectedNoteID int
	}{
		{
			"Create note with Tags",
			1,
			dto.NoteDto{
				Title:  "Test Note",
				Body:   "lorem ipsum",
				TagsID: []int{1, 2},
			},
			123,
		},
		{
			"Create note without Tags",
			10,
			dto.NoteDto{
				Title: "Test Note",
				Body:  "lorem ipsum",
			},
			100,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if len(tc.dto.TagsID) > 0 {
				mockRepo.EXPECT().
					AddNoteWithTag(tc.userID, tc.dto).
					Return(tc.expectedNoteID, nil).
					Times(1)
			} else {
				mockRepo.EXPECT().
					AddNote(tc.userID, tc.dto).
					Return(tc.expectedNoteID, nil).
					Times(1)
			}

			noteID, err := noteService.CreateNote(tc.userID, tc.dto)

			assert.NoError(t, err)
			assert.Equal(t, tc.expectedNoteID, noteID)
		})
	}
}

func TestCanGetNote(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockNote(ctrl)
	noteService := NewNoteService(mockRepo)

	userID := 10
	noteID := 100

	expected := model.Note{ID: noteID, Title: "lorem", Body: "ipsum", UserID: userID}

	mockRepo.EXPECT().
		GetNoteByID(noteID).
		Return(expected, nil).
		Times(1)

	actual, err := noteService.GetNote(userID, noteID)
	assert.NoError(t, err)
	assert.Equal(t, expected, actual)
}

func TestFailGetNoteNotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockNote(ctrl)
	noteService := NewNoteService(mockRepo)

	userID := 10
	noteID := 100

	mockRepo.EXPECT().
		GetNoteByID(noteID).
		Return(model.Note{}, repository.NewNotFoundError("note", noteID)).
		Times(1)

	_, err := noteService.GetNote(userID, noteID)
	assert.Error(t, err)
	assert.ErrorIs(t, err, repository.NewNotFoundError("note", noteID))
}

func TestFailGetNoteForbiddenError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockNote(ctrl)
	noteService := NewNoteService(mockRepo)

	userID := 10
	noteID := 100
	otherUserID := 2

	expected := model.Note{ID: noteID, Title: "lorem", Body: "ipsum", UserID: otherUserID}

	mockRepo.EXPECT().
		GetNoteByID(noteID).
		Return(expected, nil).
		Times(1)

	_, err := noteService.GetNote(userID, noteID)
	assert.ErrorIs(t, err, NewForbiddenError(userID, noteID))
}

func TestCanGetNotes(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockNote(ctrl)
	noteService := NewNoteService(mockRepo)

	userID := 1

	expected := []model.Note{
		{ID: 1, Title: "lorem1", Body: "ipsum1", UserID: userID},
		{ID: 2, Title: "lorem2", Body: "ipsum2", UserID: userID},
	}

	mockRepo.EXPECT().
		GetNotesByUserID(userID).
		Return(expected, nil).
		Times(1)

	actual, err := noteService.GetNotes(userID)
	assert.NoError(t, err)
	assert.Equal(t, expected, actual)
}

func TestFailGetNotes(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockNote(ctrl)
	noteService := NewNoteService(mockRepo)

	userID := 1

	mockRepo.EXPECT().
		GetNotesByUserID(userID).
		Return([]model.Note{}, errors.New("")).
		Times(1)

	_, err := noteService.GetNotes(userID)
	assert.Error(t, err)
}

func TestCanUpdateNote(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockNote(ctrl)
	noteService := NewNoteService(mockRepo)

	tests := []struct {
		name   string
		userID int
		noteID int
		dto    dto.NoteDto
	}{
		{
			"Update note with Tags",
			1,
			2,
			dto.NoteDto{
				Title:  "Test Note",
				Body:   "lorem ipsum",
				TagsID: []int{1, 2},
			},
		},
		{
			"Update note without Tags",
			1,
			2,
			dto.NoteDto{
				Title: "Test Note",
				Body:  "lorem ipsum",
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			value := model.Note{
				ID:     2,
				Title:  "lorem",
				Body:   "ipsum",
				UserID: tc.userID,
			}

			mockRepo.EXPECT().
				GetNoteByID(tc.noteID).
				Return(value, nil).
				Times(1)

			mockRepo.EXPECT().
				UpdateNote(tc.noteID, tc.dto).
				Return(nil).
				Times(1)

			err := noteService.UpdateNote(tc.userID, tc.noteID, tc.dto)
			assert.NoError(t, err)
		})
	}
}

func TestFailUpdateNoteForbiddenError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockNote(ctrl)
	noteService := NewNoteService(mockRepo)

	userID := 1
	otherUserID := 55
	noteID := 1
	dto := dto.NoteDto{
		Title: "New Title",
		Body:  "New Body",
	}

	value := model.Note{
		ID:     1,
		Title:  "lorem",
		Body:   "ipsum",
		UserID: otherUserID,
	}

	mockRepo.EXPECT().
		GetNoteByID(noteID).
		Return(value, nil).
		Times(1)

	err := noteService.UpdateNote(userID, noteID, dto)
	assert.ErrorIs(t, err, NewForbiddenError(userID, noteID))
}

func TestFailUpdateNoteNotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockNote(ctrl)
	noteService := NewNoteService(mockRepo)

	userID := 1
	noteID := 1
	dto := dto.NoteDto{
		Title: "New Title",
		Body:  "New Body",
	}

	mockRepo.EXPECT().
		GetNoteByID(noteID).
		Return(model.Note{}, repository.NewNotFoundError("note", noteID)).
		Times(1)

	err := noteService.UpdateNote(userID, noteID, dto)
	assert.ErrorIs(t, err, repository.NewNotFoundError("note", noteID))
}

func TestCanDeleteNote(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockNote(ctrl)
	noteService := NewNoteService(mockRepo)

	userID := 1
	noteID := 1
	value := model.Note{
		ID:     1,
		Title:  "lorem",
		Body:   "ipsum",
		UserID: userID,
	}

	mockRepo.EXPECT().
		GetNoteByID(noteID).
		Return(value, nil).
		Times(1)

	mockRepo.EXPECT().
		DeleteNote(noteID).
		Return(true, nil).
		Times(1)

	err := noteService.DeleteNote(userID, noteID)
	assert.NoError(t, err)
}

func TestFailDeleteNote(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockNote(ctrl)
	noteService := NewNoteService(mockRepo)

	userID := 1
	noteID := 1
	value := model.Note{
		ID:     1,
		Title:  "lorem",
		Body:   "ipsum",
		UserID: userID,
	}

	mockRepo.EXPECT().
		GetNoteByID(noteID).
		Return(value, nil).
		Times(1)

	mockRepo.EXPECT().
		DeleteNote(noteID).
		Return(false, nil).
		Times(1)

	err := noteService.DeleteNote(userID, noteID)
	assert.Error(t, err)
}

func TestFailDeleteNoteForbiddenError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockNote(ctrl)
	noteService := NewNoteService(mockRepo)

	userID := 1
	otherUserID := 2
	noteID := 1
	value := model.Note{
		ID:     noteID,
		Title:  "lorem",
		Body:   "ipsum",
		UserID: otherUserID,
	}

	mockRepo.EXPECT().
		GetNoteByID(noteID).
		Return(value, nil).
		Times(1)

	err := noteService.DeleteNote(userID, noteID)
	assert.ErrorIs(t, err, NewForbiddenError(userID, noteID))
}
