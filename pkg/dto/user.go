package dto

import (
	"errors"
	"fmt"
)

const minLenPassword = 6

type UserSingUpDto struct {
	Email           string `form:"email" binding:"required"`
	Password        string `form:"password" binding:"required"`
	ConfirmPassword string `form:"confirm_password" binding:"required"`
}

func (d *UserSingUpDto) Validate() error {
	if d.Password != d.ConfirmPassword {
		return errors.New("password not equal confirm password")
	}

	if len([]rune(d.Password)) < minLenPassword {
		return fmt.Errorf("password must be more %d character", minLenPassword)
	}

	return nil
}

type UserSingInDto struct {
	Email    string `form:"email" binding:"required"`
	Password string `form:"password" binding:"required"`
}
