package service

import "fmt"

type ForbiddenError struct {
	UserID   int
	EntityID int
}

func NewForbiddenError(userID int, entityID int) *ForbiddenError {
	return &ForbiddenError{UserID: userID, EntityID: entityID}
}

func (err ForbiddenError) Error() string {
	return fmt.Sprintf("user id %d access denied to note id %d", err.UserID, err.EntityID)
}
