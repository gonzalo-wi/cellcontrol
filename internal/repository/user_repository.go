package repository

import (
	"github.com/gonzalo-wi/cellcontrol/internal/domain"
	"gorm.io/gorm"
)

type UserRepository interface {
	CreateUser(user *domain.User) error
	GetAllUsers() ([]domain.User, error)
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) CreateUser(user *domain.User) error {
	return r.db.Create(user).Error
}

func (r *userRepository) GetAllUsers() ([]domain.User, error) {
	var users []domain.User
	err := r.db.Find(&users).Error
	return users, err
}
