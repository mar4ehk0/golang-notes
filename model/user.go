package model

type User struct {
	ID       int    `form:"-"`
	Email    string `form:"email" binding:"required"`
	Password string `form:"password" binding:"required"`
}
