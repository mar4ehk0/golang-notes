package dto

import "github.com/mar4ehk0/notes/model"

type TagsWithTagged struct {
	model.Tag
	Tagged bool
}

func NewTagsWithTagged(tag model.Tag, tagged bool) *TagsWithTagged {
	return &TagsWithTagged{Tag: tag, Tagged: tagged}
}
