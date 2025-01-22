package repository

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/mar4ehk0/notes/model"
)

type TagPostgres struct {
	db *sqlx.DB
}

func NewTagPostgres(db *sqlx.DB) *TagPostgres {
	return &TagPostgres{db: db}
}

func (r *TagPostgres) GetTags() ([]model.Tag, error) {
	tags := []model.Tag{}

	query := fmt.Sprintf("SELECT id, name FROM %s", tagsTable)
	err := r.db.Select(&tags, query)
	if err != nil {
		return tags, fmt.Errorf("select tags: %w", err)
	}

	return tags, nil
}

func (r *TagPostgres) GetTagByID(tagID int) (model.Tag, error) {
	var tag model.Tag

	query := fmt.Sprintf("SELECT id, name FROM %s WHERE id=$1", tagsTable)
	err := r.db.QueryRowx(query, tagID).StructScan(&tag)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return tag, NewNotFoundError("tag", tagID)
		}
		return tag, fmt.Errorf("struct scan {%v}: %w", tagID, err)
	}

	return tag, nil
}

func (r *TagPostgres) GetTagsByNoteID(noteID int) ([]model.Tag, error) {
	tags := []model.Tag{}

	query := fmt.Sprintf(
		"SELECT t.id, t.name FROM %s as tn INNER JOIN %s as t ON t.id = tn.tag_id WHERE tn.note_id=$1",
		tagsNotesTable,
		tagsTable,
	)
	err := r.db.Select(&tags, query, noteID)
	if err != nil {
		return tags, fmt.Errorf("select: %w", err)
	}

	return tags, nil
}
