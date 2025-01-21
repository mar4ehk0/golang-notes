package dto

type NoteDto struct {
	Title string `form:"title" binding:"required"`
	Body  string `form:"body" binding:"required"`
	TagsID []int `form:"tags_id[]"`
}
