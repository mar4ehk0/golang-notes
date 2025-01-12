package repository

import (
	"fmt"

	"github.com/jmoiron/sqlx"
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
