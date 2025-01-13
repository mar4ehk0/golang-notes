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
}

type Tag interface {
}

type Repository struct {
	Authorization
	Note
	Tag
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		Authorization: NewAuthPostgres(db),
	}
}
