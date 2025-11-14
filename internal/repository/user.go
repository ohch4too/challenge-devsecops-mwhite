package repository

import (
	"challenge/internal/domain"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// UserRepository defines data access operations for users
type UserRepository interface {
	Add(u *domain.User) error
	Get(id string) (*domain.User, error)
	List() ([]domain.User, error)
	Delete(id string) error
}

// userRepository implements UserRepository
type userRepository struct {
	db *gorm.DB
}

// NewUserRepository creates a new user repository
func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

// Add creates a new user
func (r *userRepository) Add(u *domain.User) error {
	// Hash password before saving
	if u.Password != "" {
		hashed, hashErr := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
		if hashErr != nil {
			return hashErr
		}
		u.Password = string(hashed)
	}

	return r.db.Create(u).Error
}

// Get retrieves a user by ID
func (r *userRepository) Get(id string) (*domain.User, error) {
	var user domain.User
	err := r.db.Where("ID = ?", id).First(&user).Error
	if err != nil {
		return nil, err
	}
	// Clear password before returning for security
	user.Password = ""
	return &user, nil
}

// List retrieves all users
func (r *userRepository) List() ([]domain.User, error) {
	var users []domain.User
	result := r.db.Find(&users)
	// Clear passwords before returning for security
	for i := range users {
		users[i].Password = ""
	}
	return users, result.Error
}

// Delete removes a user by ID
func (r *userRepository) Delete(id string) error {
	return r.db.Delete(&domain.User{}, id).Error
}
