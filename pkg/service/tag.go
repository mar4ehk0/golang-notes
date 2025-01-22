package service

import (
	"fmt"

	"github.com/mar4ehk0/notes/model"
	"github.com/mar4ehk0/notes/pkg/dto"
	"github.com/mar4ehk0/notes/pkg/repository"
)

type TagService struct {
	repo repository.Tag
}

func NewTagService(repo repository.Tag) *TagService {
	return &TagService{repo: repo}
}

func (s *TagService) GetTagsWithTaggedByNoteID(noteID int) ([]dto.TagsWithTagged, error) {
	tags, err := s.GetTags()
	if err != nil {
		return []dto.TagsWithTagged{}, fmt.Errorf("repo get tags: %w", err)
	}

	tagsForNote, err := s.GetTagsByNoteId(noteID)
	if err != nil {
		return []dto.TagsWithTagged{}, fmt.Errorf("repo get tags by nodeID{%d}: %w", noteID, err)
	}

	tagsTagged := make([]dto.TagsWithTagged, 0)
	var tagged bool

	for _, tag := range tags {
		for _, tagForNote := range tagsForNote {
			tagged = false
			if tag.ID == tagForNote.ID {
				tagged = true
			}

			tagsTagged = append(tagsTagged, *dto.NewTagsWithTagged(tag, tagged))
		}
	}

	return tagsTagged, nil
}

func (s *TagService) GetTags() ([]model.Tag, error) {
	tags, err := s.repo.GetTags()
	if err != nil {
		return make([]model.Tag, 0), fmt.Errorf("repo get tags: %w", err)
	}

	return tags, nil
}

func (s *TagService) GetTagByID(tagID int) (model.Tag, error) {
	tag, err := s.repo.GetTagByID(tagID)
	if err != nil {
		return model.Tag{}, fmt.Errorf("repo get tag: %w", err)
	}

	return tag, nil
}

func (s *TagService) GetTagsByNoteId(noteID int) ([]model.Tag, error) {
	tags, err := s.repo.GetTagsByNoteId(noteID)
	if err != nil {
		return make([]model.Tag, 0), fmt.Errorf("repo get tags by nodeID{%d}: %w", noteID, err)
	}

	return tags, nil
}
