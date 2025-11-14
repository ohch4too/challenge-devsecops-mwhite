package repository

import (
	"challenge/internal/domain"
	"testing"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func setupTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		t.Fatalf("failed to open test database: %v", err)
	}
	if err := db.AutoMigrate(&domain.User{}); err != nil {
		t.Fatalf("failed to migrate: %v", err)
	}
	return db
}

func TestUserRepository_Add(t *testing.T) {
	db := setupTestDB(t)
	repo := NewUserRepository(db)

	user := &domain.User{
		Firstname: "John",
		Lastname:  "Doe",
		Login:     "johndoe",
		Password:  "password123",
	}

	if err := repo.Add(user); err != nil {
		t.Errorf("Add() error = %v", err)
	}

	if user.ID == 0 {
		t.Error("Add() did not set user ID")
	}

	if user.Password == "password123" {
		t.Error("Add() did not hash password")
	}
}

func TestUserRepository_Get(t *testing.T) {
	db := setupTestDB(t)
	repo := NewUserRepository(db)

	user := &domain.User{
		Firstname: "Jane",
		Lastname:  "Doe",
		Login:     "janedoe",
		Password:  "password123",
	}
	repo.Add(user)

	retrieved, err := repo.Get("1")
	if err != nil {
		t.Errorf("Get() error = %v", err)
	}

	if retrieved.Login != "janedoe" {
		t.Errorf("Get() login = %v, want janedoe", retrieved.Login)
	}

	if retrieved.Password != "" {
		t.Error("Get() did not clear password")
	}
}

func TestUserRepository_List(t *testing.T) {
	db := setupTestDB(t)
	repo := NewUserRepository(db)

	users := []*domain.User{
		{Firstname: "User1", Lastname: "Test", Login: "user1", Password: "password123"},
		{Firstname: "User2", Lastname: "Test", Login: "user2", Password: "password123"},
	}

	for _, u := range users {
		repo.Add(u)
	}

	result, err := repo.List()
	if err != nil {
		t.Errorf("List() error = %v", err)
	}

	if len(result) != 2 {
		t.Errorf("List() returned %d users, want 2", len(result))
	}

	for _, u := range result {
		if u.Password != "" {
			t.Error("List() did not clear passwords")
		}
	}
}

func TestUserRepository_Delete(t *testing.T) {
	db := setupTestDB(t)
	repo := NewUserRepository(db)

	user := &domain.User{
		Firstname: "Delete",
		Lastname:  "Me",
		Login:     "deleteme",
		Password:  "password123",
	}
	repo.Add(user)

	if err := repo.Delete("1"); err != nil {
		t.Errorf("Delete() error = %v", err)
	}

	if _, err := repo.Get("1"); err == nil {
		t.Error("Delete() did not delete user")
	}
}
