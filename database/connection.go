package database

import (
	"errors"
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func DBConnection() (*gorm.DB, error) {

	if err := godotenv.Load(); err != nil {
		panic(err)
	}

	PG_HOST := os.Getenv("APP_PG_HOST")
	PG_USER := os.Getenv("APP_PG_USER")
	PG_PASSWORD := os.Getenv("APP_PG_PASSWORD")
	PG_DBNAME := os.Getenv("APP_PG_DBNAME")
	PG_PORT := os.Getenv("APP_PG_PORT")
	PG_SSLMODE := os.Getenv("APP_PG_SSLMODE")
	PG_TIMEZONE := os.Getenv("APP_PG_TIMEZONE")

	PGConnection := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=%s",
		PG_HOST,
		PG_USER,
		PG_PASSWORD,
		PG_DBNAME,
		PG_PORT,
		PG_SSLMODE,
		PG_TIMEZONE,
	)

	return gorm.Open(postgres.Open(PGConnection), &gorm.Config{})
}

func IsNotFound(row *gorm.DB) bool {

	return errors.Is(row.Error, gorm.ErrRecordNotFound)
}
