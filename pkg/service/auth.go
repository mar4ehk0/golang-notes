package service

import (
	"crypto/sha1"
	"fmt"

	"github.com/mar4ehk0/notes/model"
	"github.com/mar4ehk0/notes/pkg/dto"
	"github.com/mar4ehk0/notes/pkg/repository"
)

const salt = "f1a2b3c4d5e6f7a8b9c0d1e2f3a4b5c6"

type AuthService struct {
	repo repository.Authorization
}

func NewAuthService(repo repository.Authorization) *AuthService {
	return &AuthService{repo: repo}
}

func (s *AuthService) CreateUser(d dto.UserSingUpDto) (model.User, error) {
	d.Password = s.generatePasswordHash(d.Password)

	result := model.User{}
	id, err := s.repo.CreateUser(d)
	if err != nil {
		return result, fmt.Errorf("service create user: %w", err)
	}

	result.ID = id
	result.Email = d.Email
	result.Password = d.Password
	return result, nil
}

func (s *AuthService) generatePasswordHash(password string) string {
	hash := sha1.New()
	hash.Write([]byte(password))

	return fmt.Sprintf("%x", hash.Sum([]byte(salt)))
}
