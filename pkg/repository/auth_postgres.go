package repository

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/mar4ehk0/notes/model"
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
			return 0, fmt.Errorf("user with email - %s exist: %w", user.Email, ErrDBDuplicateKey)
		}
		return 0, fmt.Errorf("scan id {%s}: %w", user.Email, err)
	}

	return id, nil
}

func (r *AuthPostgres) GetUserByEmail(email string) (model.User, error) {
	var user model.User

	query := fmt.Sprintf("SELECT id, email, password FROM %s WHERE email=$1", usersTable)
	err := r.db.QueryRowx(query, email).StructScan(&user)
	if err != nil {
		return user, fmt.Errorf("struct scan {%v}: %w", email, err)
	}

	return user, nil
}
