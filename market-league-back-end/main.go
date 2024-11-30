package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/market-league/api"
	"github.com/market-league/internal/db"
)

func main() {
	// Initialize the database
	db.InitDB()
	// Initializes Gin router instance with default middleware attached
	router := gin.Default()
	// CORS configuration
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:4200"}, // Allow your frontend origin
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))
	// Create all api endpoints
	api.RegisterRoutes(router)
	// Call stock updater loop
	api.StartDailyTask()
	// Enable CORS for all origins and methods
	router.Use(cors.Default())
	// Port that go backend will run on
	router.Run(":9000")
}
