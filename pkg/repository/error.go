package repository

import "fmt"

type NotFoundError struct {
	Entity string
	ID     int
}

func NewNotFoundError(e string, id int) *NotFoundError {
	return &NotFoundError{Entity: e, ID: id}
}

func (err NotFoundError) Error() string {
	return fmt.Sprintf("%s id %d not found", err.Entity, err.ID)
}
