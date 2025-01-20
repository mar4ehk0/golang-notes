package service

import (
	"fmt"

	"github.com/mar4ehk0/notes/model"
	"github.com/mar4ehk0/notes/pkg/repository"
)

type TagService struct {
	repo repository.Tag
}

func NewTagService(repo repository.Tag) *TagService {
	return &TagService{repo: repo}
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
