package service

import (
	"challenge/internal/domain"
	"challenge/internal/repository"
)

// UserService handles user business logic
type UserService interface {
	AddUser(u *domain.User) error
	GetUser(id string) (*domain.User, error)
	ListUsers() ([]domain.User, error)
	DeleteUser(id string) error
}

// userService implements UserService
type userService struct {
	repo repository.UserRepository
}

// NewUserService creates a new user service
func NewUserService(userRepo repository.UserRepository) UserService {
	return &userService{repo: userRepo}
}

// AddUser adds a new user
func (s *userService) AddUser(u *domain.User) error {
	return s.repo.Add(u)
}

// GetUser retrieves a user by ID
func (s *userService) GetUser(id string) (*domain.User, error) {
	return s.repo.Get(id)
}

// ListUsers retrieves all users
func (s *userService) ListUsers() ([]domain.User, error) {
	return s.repo.List()
}

// DeleteUser deletes a user
func (s *userService) DeleteUser(id string) error {
	return s.repo.Delete(id)
}
