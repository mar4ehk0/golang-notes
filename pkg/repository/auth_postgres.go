package repository

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/mar4ehk0/notes/pkg/dto"
)

type AuthPostgres struct {
	db *sqlx.DB
}

func NewAuthPostgres(db *sqlx.DB) *AuthPostgres {
	return &AuthPostgres{db: db}
}

func (r *AuthPostgres) CreateUser(user dto.UserSingUpDto) (int, error) {
	var id int

	query := fmt.Sprintf("INSERT INTO %s (email, password) VALUES ($1, $2) RETURNING id", usersTable)
	row := r.db.QueryRow(query, user.Email, user.Password)

	err := row.Scan(&id)
	if err != nil {
		if isErrDuplicate(err) {
			return 0, fmt.Errorf("user with same email - %s exist: %w", user.Email, ErrDBDuplicateKey)
		}
		return 0, fmt.Errorf("scan id {%s}: %w", user.Email, err)
	}

	return id, nil
}
