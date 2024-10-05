package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	// GORM
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// Counter model
type Counter struct {
	ID    uint `gorm:"primaryKey"`
	Value int  `gorm:"default:0"`
}

var db *gorm.DB

// Initialize GORM with PostgreSQL
func initDB() {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_SSLMODE"),
	)

	var err error
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Automatically run migrations
	err = db.AutoMigrate(&Counter{})
	if err != nil {
		log.Fatalf("Failed to migrate database schema: %v", err)
	}
}

func main() {
	// Initialize the database
	initDB()
	// Initializes Gin router instance with default middleware attached
	router := gin.Default()
	// Enable CORS for all origins and methods
	router.Use(cors.Default())

	// Route to increment the counter
	router.GET("/api/increment", incrementCounter)

	// Route to get the counter value
	router.GET("/api/counter", getCounterValue)

	router.Run(":9000")
}

// Handler to increment the counter
func incrementCounter(c *gin.Context) {
	fmt.Println("Increment Counter")
	// Retrieve the counter (assuming single row)
	var counter Counter
	db.First(&counter)

	// Increment the counter value
	counter.Value++

	// Save updated counter back to the database
	db.Save(&counter)

	// Return the updated counter value as JSON
	c.JSON(http.StatusOK, gin.H{
		"value": counter.Value,
	})
}

// Handler to get the counter value
func getCounterValue(c *gin.Context) {
	// Retrieve the counter
	var counter Counter
	db.First(&counter)

	// Return the counter value as JSON
	c.JSON(http.StatusOK, gin.H{
		"value": counter.Value,
	})
}
