package db

import (
	"challenge/internal/domain"
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// Conn is the global database connection
var Conn *gorm.DB

// sqliteConnector creates a SQLite database connection
func sqliteConnector() error {
	var err error
	Conn, err = gorm.Open(sqlite.Open("./test.db"), &gorm.Config{})
	return err
}

// postgresConnection creates a PostgreSQL database connection
func postgresConnection(host, user, password, dbname string) error {
	if password == "" {
		return fmt.Errorf("POSTGRES_PASSWORD not set")
	}
	dsn := fmt.Sprintf("host=%v user=%v password=%v dbname=%v port=5432 sslmode=disable", host, user, password, dbname)
	var err error
	Conn, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	return err
}

// Initialize sets up the database connection and runs migrations
func Initialize(host, user, password, dbname string) error {
	fmt.Println("DevSecOps challenge")

	// Connect to database
	var dbErr error
	if host == "" {
		dbErr = sqliteConnector()
	} else {
		fmt.Println("Using postgresql DB driver")
		dbErr = postgresConnection(host, user, password, dbname)
	}

	if dbErr != nil {
		return dbErr
	}

	// Run migrations
	return Conn.AutoMigrate(&domain.User{})
}

// SeedAdminUser creates an admin user if none exist
func SeedAdminUser(adminPassword string) error {
	var users []domain.User
	if err := Conn.Find(&users).Error; err != nil {
		return err
	}

	if len(users) == 0 {
		fmt.Println("Could not find any users, bootstrapping an admin account")
		if adminPassword == "" {
			adminPassword = "changeme"
			fmt.Println("Warning: ADMIN_PASSWORD not set, using default")
		}
		admin := domain.User{
			Firstname: "Admin",
			Lastname:  "Istrator",
			Login:     "admin",
			Password:  adminPassword,
		}
		if err := Conn.Create(&admin).Error; err != nil {
			return fmt.Errorf("could not create admin user, reason %v", err)
		}
	} else {
		fmt.Println("Found users, skipping admin account bootstrapping")
	}

	return nil
}
