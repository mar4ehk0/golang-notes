package service

import (
	"fmt"

	"github.com/mar4ehk0/notes/model"
	"github.com/mar4ehk0/notes/pkg/dto"
	"github.com/mar4ehk0/notes/pkg/repository"
	"golang.org/x/crypto/bcrypt"
)

const salt = "f1a2b3c4d5e6f7a8b9c0d1e2f3a4b5c6"

type AuthService struct {
	repo repository.Authorization
}

func NewAuthService(repo repository.Authorization) *AuthService {
	return &AuthService{repo: repo}
}

func (s *AuthService) CreateUser(d dto.UserSingUpDto) (model.User, error) {
	var user model.User

	hashedPassword, err := s.generatePasswordHash(d.Password)
	if err != nil {
		return user, fmt.Errorf("generate password hash: %w", err)
	}

	d.Password = hashedPassword

	id, err := s.repo.CreateUser(d)
	if err != nil {
		return user, fmt.Errorf("repo create user: %w", err)
	}

	user.ID = id
	user.Email = d.Email
	user.Password = d.Password
	return user, nil
}

func (s *AuthService) Authorize(d dto.UserSingInDto) (model.User, bool, error) {
	var user model.User

	user, err := s.repo.GetUserByEmail(d.Email)
	if err != nil {
		return user, false, fmt.Errorf("repo get user by email: %w", err)
	}

	if s.comparePassword(user.Password, d.Password) {
		return user, true, nil
	}

	return user, false, nil
}

func (s *AuthService) generatePasswordHash(password string) (string, error) {
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	return string(hashedBytes), err
}

func (s *AuthService) comparePassword(hashedPassword string, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))

	return err == nil
}
