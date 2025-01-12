package service

import "github.com/mar4ehk0/notes/pkg/repository"

type Authorization interface {
}

type Note interface {
}

type Tag interface {
}

type Service struct {
	Authorization
	Note
	Tag
}

func NewService(repository *repository.Repository) *Service {
	return &Service{}
}
