package main

import (
	"fmt"
	"log"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres" // Underscore to prevent unused import error
	_ "github.com/golang-migrate/migrate/v4/source/file"       // Underscore to prevent unused import error
)

var databaseURL string
var migrationsPath string

func init() {
	databaseURL = fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_SSLMODE"),
	)

	migrationsPath = "file://./migrations"
}

func main() {
	// Run migrations when the app starts
	runMigrations()
	// Initializes Gin router instance with default middleware attached
	router := gin.Default()
	// Enable CORS for all origins and methods
	router.Use(cors.Default())

	router.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "test",
		})
	})

	router.Run(":9000")
}

func runMigrations() {
	// Create a new migrate instance
	m, err := migrate.New(migrationsPath, databaseURL)
	if err != nil {
		log.Fatalf("Could not create migrate instance: %v", err)
	}

	// Run all "up" migrations
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatalf("Migration failed: %v", err)
	} else if err == migrate.ErrNoChange {
		log.Println("No migrations to apply")
	} else {
		log.Println("Migrations applied successfully")
	}
}
