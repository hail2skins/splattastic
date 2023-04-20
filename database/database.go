package database

import (
	"fmt"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var Database *gorm.DB

func Connect() {
	var err error
	appEnv := os.Getenv("APP_ENV")
	host := os.Getenv("DB_HOST")
	username := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	databaseName := os.Getenv("DB_NAME")
	port := os.Getenv("DB_PORT")
	sslrootcert := os.Getenv("DB_SSLROOTCERT")

	if appEnv == "DEV" {
		dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=verify-full sslrootcert=%s TimeZone=America/Chicago", host, username, password, databaseName, port, sslrootcert)
		Database, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	} else if appEnv == "PROD" {
		dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s TimeZone=America/Chicago", host, username, password, databaseName, port)
		Database, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	} else {
		panic("APP_ENV not set or has an invalid value. Must be either 'DEV' or 'PROD'.")
	}

	if err != nil {
		panic(err)
	} else {
		fmt.Println("Successfully connected to the database")
	}
}
