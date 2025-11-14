package db

import (
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"os"
)

var Conn *gorm.DB
var err error

var DBHost string
var DBUser string

func init() {

	if len(os.Getenv("POSTGRES_USER")) != 0 {
		DBUser = os.Getenv("POSTGRES_USER")
	} else {
		DBUser = "challenge"
	}

	if len(os.Getenv("POSTGRES_HOST")) != 0 {
		DBHost = os.Getenv("POSTGRES_HOST")
	} else {
		DBHost = "localhost"
	}

}

func SqliteConnector() error {

	Conn, err = gorm.Open(sqlite.Open("./test.db"), &gorm.Config{})
	if err != nil {
		return err
	}

	return nil
}

func PostgresConnection() error {

	dbPassword := os.Getenv("POSTGRES_PASSWORD")
	if dbPassword == "" {
		fmt.Println("POSTGRES_PASSWORD not set")
		return fmt.Errorf("POSTGRES_PASSWORD not set")
	}
	dsn := fmt.Sprintf("host=%v user=%v password=%v dbname=challenge port=5432 sslmode=disable", DBHost, DBUser, dbPassword)
	Conn, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		return err
	}

	return nil

}
