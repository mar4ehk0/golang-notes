package model

const (
	DefaultTagId = -1
)

type Tag struct {
	ID   int    `db:"id"`
	Name string `db:"name"`
}
