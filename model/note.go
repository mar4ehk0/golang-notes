package model

type Note struct {
	ID     int    `db:"id"`
	Title  string `db:"title"`
	Body   string `db:"body"`
	UserID int    `db:"user_id"`
}
