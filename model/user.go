package model

type User struct {
	ID       int    `form:"-" db:"id"`
	Email    string `form:"email" binding:"required"`
	Password string `form:"password" binding:"required"`
}
