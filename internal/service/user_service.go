package service

import (
	"strings"

	"github.com/gonzalo-wi/cellcontrol/internal/domain"
	"github.com/gonzalo-wi/cellcontrol/internal/repository"
)

type UserService interface {
	CreateUser(nombre, apellido, email, reparto string) error
	GetAllUsers() ([]domain.User, error)
}

type userService struct {
	repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) UserService {
	return &userService{repo: repo}
}

func (s *userService) CreateUser(nombre, apellido, email, reparto string) error {
	u := &domain.User{
		Nombre:   strings.TrimSpace(nombre),
		Apellido: strings.TrimSpace(apellido),
		Email:    strings.ToLower(strings.TrimSpace(email)),
		Reparto:  strings.TrimSpace(reparto),
	}
	if err := s.repo.CreateUser(u); err != nil {
		return err
	}
	return nil
}

func (s *userService) GetAllUsers() ([]domain.User, error) {
	return s.repo.GetAllUsers()
}
