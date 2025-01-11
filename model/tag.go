package model

type Tag struct {
	ID int `json:"-"`
	Name string `json:"name"`
}
