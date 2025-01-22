package repository

import (
	"errors"
	"fmt"

	"github.com/jackc/pgx"
	"github.com/jmoiron/sqlx"
)

const (
	usersTable     = "users"
	notesTable     = "notes"
	tagsTable      = "tags"
	tagsNotesTable = "tags_notes"
)

var (
	ErrDBDuplicateKey = errors.New("duplicate value for unique index")
)

type ConfigPostgres struct {
	Host     string
	Port     string
	Username string
	Password string
	DBName   string
	SSLMode  string
}

func NewPostgresDB(c ConfigPostgres) (*sqlx.DB, error) {
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s dbname=%s sslmode=%s password=%s",
		c.Host, c.Port, c.Username, c.DBName, c.SSLMode, c.Password,
	)
	sqlxDB, err := sqlx.Connect("pgx", dsn)

	if err != nil {
		return nil, err
	}

	return sqlxDB, err
}

func isErrDuplicate(err error) bool {
	var pgErr pgx.PgError
	if errors.As(err, &pgErr) {
		return pgErr.Code == "23505"
	}

	return false
}
