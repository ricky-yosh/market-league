package api

import (

	// "github.com/market-league/internal/league"
	// "github.com/market-league/internal/user"
	// "github.com/market-league/internal/portfolio"
	// "github.com/market-league/internal/trade"

	"github.com/gin-gonic/gin"
	"github.com/market-league/internal/auth"
	"github.com/market-league/internal/db"
)

func RegisterRoutes(router *gin.Engine) {
	// Test route
	router.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Welcome to MarketLeague API!",
		})
	})
	database := db.GetDB()

	// Set up the authentication flow by initializing the repository, service, and handler layers
	authRepo := auth.NewAuthRepository(database)
	authService := auth.NewAuthService(authRepo)
	authHandler := auth.NewAuthHandler(authService)

	// Auth routes
	authRoutes := router.Group("/api/auth")
	{
		authRoutes.POST("/signup", authHandler.Signup)
		authRoutes.POST("/login", authHandler.Login)
	}
}
