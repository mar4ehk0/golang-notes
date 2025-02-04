package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/jmoiron/sqlx"
	"github.com/mar4ehk0/notes/model"
	"github.com/mar4ehk0/notes/pkg/dto"
	"github.com/sirupsen/logrus"
)

type NotePostgres struct {
	db *sqlx.DB
}

func NewNotePostgres(db *sqlx.DB) *NotePostgres {
	return &NotePostgres{db: db}
}

func (r *NotePostgres) AddNoteWithTag(userID int, input dto.NoteDto) (int, error) {
	var err error

	tx, err := r.db.Begin()
	if err != nil {
		return 0, fmt.Errorf("begin transaction {%v}: %w", input, err)
	}
	defer func() {
		if err != nil {
			if rbErr := tx.Rollback(); rbErr != nil {
				logrus.Errorf("transaction rollback error: %v", rbErr)
			}
		} else {
			if cmtErr := tx.Commit(); cmtErr != nil {
				logrus.Errorf("transaction commit error: %v", cmtErr)
				err = cmtErr
			}
		}
	}()

	var noteID int

	query := fmt.Sprintf("INSERT INTO %s (title, body, user_id) VALUES ($1, $2, $3) RETURNING id", notesTable)
	row := tx.QueryRow(query, input.Title, input.Body, userID)
	err = row.Scan(&noteID)
	if err != nil {
		err = fmt.Errorf("scan id {%d %v}: %w", userID, input, err)
		return 0, err
	}

	err = r.addTagNoteTx(tx, input.TagsID, noteID)
	if err != nil {
		return 0, fmt.Errorf("exec: %w", err)
	}

	return noteID, nil
}

func (r *NotePostgres) AddNote(userID int, input dto.NoteDto) (int, error) {
	var id int

	query := fmt.Sprintf("INSERT INTO %s (title, body, user_id) VALUES ($1, $2, $3) RETURNING id", notesTable)
	row := r.db.QueryRow(query, input.Title, input.Body, userID)

	err := row.Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("scan id {%d %v}: %w", userID, input, err)
	}

	return id, nil
}

func (r *NotePostgres) GetNoteByID(noteID int) (model.Note, error) {
	var note model.Note

	query := fmt.Sprintf("SELECT id, title, body, user_id FROM %s WHERE id=$1", notesTable)
	err := r.db.QueryRowx(query, noteID).StructScan(&note)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return note, NewNotFoundError("note", noteID)
		}
		return note, fmt.Errorf("struct scan {%v}: %w", noteID, err)
	}

	return note, nil
}

func (r *NotePostgres) GetNotesByUserID(userID int) ([]model.Note, error) {
	notes := []model.Note{}

	query := fmt.Sprintf("SELECT id, title, body, user_id FROM %s WHERE user_id=$1", notesTable)
	err := r.db.Select(&notes, query, userID)
	if err != nil {
		return notes, fmt.Errorf("select notes for userID {%v}: %w", userID, err)
	}

	return notes, nil
}

func (r *NotePostgres) UpdateNote(noteID int, d dto.NoteDto) error {
	tx, err := r.db.Begin()
	if err != nil {
		return fmt.Errorf("begin transaction {%v}: %w", d, err)
	}
	defer func() {
		if err != nil {
			if rbErr := tx.Rollback(); rbErr != nil {
				logrus.Errorf("transaction rollback error: %v", rbErr)
			}
		} else {
			if cmtErr := tx.Commit(); cmtErr != nil {
				logrus.Errorf("transaction commit error: %v", cmtErr)
				err = cmtErr
			}
		}
	}()

	query := fmt.Sprintf("UPDATE %s SET title=$1, body=$2 WHERE id=$3", notesTable)
	_, err = tx.Exec(query, d.Title, d.Body, noteID)
	if err != nil {
		return fmt.Errorf("update %s exec: %w", notesTable, err)
	}

	query = fmt.Sprintf("DELETE FROM %s WHERE note_id=$1", tagsNotesTable)
	_, err = tx.Exec(query, noteID)
	if err != nil {
		return fmt.Errorf("delete %s exec: %w", tagsNotesTable, err)
	}

	if len(d.TagsID) == 0 {
		return nil
	}

	err = r.addTagNoteTx(tx, d.TagsID, noteID)
	if err != nil {
		return fmt.Errorf("add %s exec: %w", tagsNotesTable, err)
	}

	return nil
}

func (r *NotePostgres) DeleteNote(noteID int) (bool, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return false, fmt.Errorf("begin transaction noteID{%v}: %w", noteID, err)
	}
	defer func() {
		if err != nil {
			if rbErr := tx.Rollback(); rbErr != nil {
				logrus.Errorf("transaction rollback error: %v", rbErr)
			}
		} else {
			if cmtErr := tx.Commit(); cmtErr != nil {
				logrus.Errorf("transaction commit error: %v", cmtErr)
				err = cmtErr
			}
		}
	}()

	var query string

	query = fmt.Sprintf("DELETE FROM %s WHERE note_id=$1", tagsNotesTable)
	_, err = tx.Exec(query, noteID)
	if err != nil {
		return false, fmt.Errorf("delete %s exec: %w", tagsNotesTable, err)
	}

	query = fmt.Sprintf("DELETE FROM %s WHERE id=$1", notesTable)
	result, err := tx.Exec(query, noteID)
	if err != nil {
		return false, fmt.Errorf("delete %s exec: %w", notesTable, err)
	}

	countDeleted, err := result.RowsAffected()
	if err != nil {
		return false, fmt.Errorf("get rows affected: %w", err)
	}

	return countDeleted > 0, nil
}

func (r *NotePostgres) addTagNoteTx(tx *sql.Tx, tagsID []int, noteID int) error {
	dto := dto.NewTagsNotesForNote(tagsID, noteID)

	query := fmt.Sprintf("INSERT INTO %s (tag_id, note_id) VALUES ", tagsNotesTable)
	placeholders := make([]string, 0)
	values := make([]interface{}, 0)
	for i, v := range dto {
		placeholders = append(placeholders, fmt.Sprintf("($%d, $%d)", i*2+1, i*2+2))
		values = append(values, v.TagID, v.NoteID)
	}
	query += strings.Join(placeholders, ",")

	_, err := tx.Exec(query, values...)
	if err != nil {
		return fmt.Errorf("exec: %w", err)
	}

	return nil
}
