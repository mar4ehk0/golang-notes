package dto

type NoteDto struct {
	Title string `form:"title" binding:"required"`
	Body  string `form:"body" binding:"required"`
	TagID int `form:"tag" binding:"required"`
}
