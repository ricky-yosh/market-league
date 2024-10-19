package db

import (
	// GORM
	"fmt"
	"log"
	"os"

	"github.com/market-league/internal/user"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

// Initialize GORM with PostgreSQL
func InitDB() {
	// CONNECTION
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_SSLMODE"),
	)

	//
	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// MIGRATIONS
	err = DB.AutoMigrate(
		// Add migrations go here
		&user.User{},
	)

	if err != nil {
		log.Fatalf("Failed to migrate database schema: %v", err)
	}
}

// GetDB returns the initialized *gorm.DB instance
func GetDB() *gorm.DB {
	if DB == nil {
		log.Fatal("Database is not initialized. Call InitDB() first.")
	}
	return DB
}