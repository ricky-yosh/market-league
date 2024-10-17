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
	// Create all api endpoints
	api.RegisterRoutes(router)
	// Enable CORS for all origins and methods
	router.Use(cors.Default())
	// Port that go backend will run on
	router.Run(":9000")
}
