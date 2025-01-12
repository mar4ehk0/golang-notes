package repository

import "github.com/jmoiron/sqlx"

type Authorization interface {
}

type Note interface {
}

type Tag interface {
}

type Repository struct {
	Authorization
	Note
	Tag
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{}
}
