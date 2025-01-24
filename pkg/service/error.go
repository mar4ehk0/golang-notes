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

func (err ForbiddenError) Is(target error) bool {
	t, ok := target.(*ForbiddenError)
	if !ok {
		return false
	}
	return err.UserID == t.UserID && err.EntityID == t.EntityID
}
