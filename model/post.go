package model

type Post struct {
	ID     int    `json:"-"`
	Title  string `json:"title"`
	Body   string `json:"body"`
	UserID int    `json:"user_id"`
}
