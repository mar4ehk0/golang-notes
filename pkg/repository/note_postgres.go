package repository

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/mar4ehk0/notes/model"
	"github.com/mar4ehk0/notes/pkg/dto"
)

type NotePostgres struct {
	db *sqlx.DB
}

func NewNotePostgres(db *sqlx.DB) *NotePostgres {
	return &NotePostgres{db: db}
}

func (r *NotePostgres) AddNote(userId int, input dto.NoteCreateDto) (model.Note, error) {
	var id int
	var note model.Note

	query := fmt.Sprintf("INSERT INTO %s (title, body, user_id) VALUES ($1, $2, $3) RETURNING id", notesTable)
	row := r.db.QueryRow(query, input.Title, input.Body, userId)

	err := row.Scan(&id)
	if err != nil {
		return note, fmt.Errorf("scan id {%d %v}: %w", userId, input, err)
	}

	return model.Note{ID: id, Title: input.Title, Body: input.Body}, nil
}

func (r *NotePostgres) GetNoteByID(noteId int) (model.Note, error) {
	var note model.Note

	query := fmt.Sprintf("SELECT id, title, body, user_id FROM %s WHERE id=$1", notesTable)
	err := r.db.QueryRowx(query, noteId).StructScan(&note)
	if err != nil {
		return note, fmt.Errorf("struct scan {%v}: %w", noteId, err)
	}

	return note, nil
}

func (r *NotePostgres) GetNotesByUserId(userID int) ([]model.Note, error) {
	notes := []model.Note{}

	query := fmt.Sprintf("SELECT id, title, body, user_id FROM %s WHERE user_id=$1", notesTable)
	err := r.db.Select(&notes, query, userID)
	if err != nil {
		return notes, fmt.Errorf("select notes for userID {%v}: %w", userID, err)
	}

	return notes, nil
}
