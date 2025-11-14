package dummy

import (
	"fmt"

	"challenge/internal/db"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
	Login     string `json:"login"`
	Password  string `json:"-"`
}

// To manage more than one user object
var users []User

// Add one user
func UserAdd(u *User) (err error) {

	// Hash password before saving
	if u.Password != "" {
		hashed, hashErr := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
		if hashErr != nil {
			return hashErr
		}
		u.Password = string(hashed)
	}

	if err = db.Conn.Create(u).Error; err != nil {
		return err
	}
	return nil

}

// Delete one user
func UserDel(id string) (err error) {
	result := db.Conn.Delete(&User{}, id)

	if result.Error != nil {
		return result.Error
	}
	return result.Error
}

// Update one user
func UserUpdate() {
	fmt.Println("Not implemented yet.")
}

// Get one user
func UserGet(u *User, id string) (err error) {
	err = db.Conn.Where("ID = ?", id).First(u).Error
	return err
}

// List all users
func UserList() ([]User, int64, error) {

	results := db.Conn.Find(&users)
	return users, results.RowsAffected, results.Error

}

// Find one user
func UserFind(u User) {
	fmt.Printf("You are looking for user: %v, but that function is not implemented yet.\n", u)
}
