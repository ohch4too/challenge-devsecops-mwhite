package domain

import (
	"errors"
	"strings"

	"gorm.io/gorm"
)

// User represents a user in the system
type User struct {
	gorm.Model
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
	Login     string `json:"login"`
	Password  string `json:"password,omitempty"` // omitempty prevents empty password in responses
}

// Validate checks if the user has required fields
func (u *User) Validate() error {
	if strings.TrimSpace(u.Firstname) == "" {
		return errors.New("firstname is required")
	}
	if strings.TrimSpace(u.Lastname) == "" {
		return errors.New("lastname is required")
	}
	if strings.TrimSpace(u.Login) == "" {
		return errors.New("login is required")
	}
	if len(u.Login) < 3 {
		return errors.New("login must be at least 3 characters")
	}
	if strings.TrimSpace(u.Password) == "" {
		return errors.New("password is required")
	}
	if len(u.Password) < 8 {
		return errors.New("password must be at least 8 characters")
	}
	return nil
}
