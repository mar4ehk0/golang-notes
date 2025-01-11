package model

type TagPost struct {
	ID     int `json:"-"`
	TagID  int `json:"tag_id"`
	PostID int `json:"post_id"`
}
