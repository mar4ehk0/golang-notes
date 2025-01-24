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

func (err NotFoundError) Is(target error) bool {
	t, ok := target.(*NotFoundError)
	if !ok {
		return false
	}
	return err.Entity == t.Entity && err.ID == t.ID
}
