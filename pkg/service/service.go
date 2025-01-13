package service

import (
	"github.com/mar4ehk0/notes/model"
	"github.com/mar4ehk0/notes/pkg/dto"
	"github.com/mar4ehk0/notes/pkg/repository"
)

type Authorization interface {
	CreateUser(d dto.UserSingUpDto) (model.User, error)
	CanAuthorize(d dto.UserSingInDto) (bool, error)
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
	return &Service{
		Authorization: NewAuthService(repository),
	}
}
