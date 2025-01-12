package model

type TagNote struct {
	ID     int `json:"-"`
	TagID  int `json:"tag_id"`
	NoteID int `json:"note_id"`
}
