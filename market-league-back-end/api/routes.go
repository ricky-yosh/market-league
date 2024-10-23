package api

import (
	"github.com/gin-gonic/gin"
	"github.com/market-league/internal/auth"
	"github.com/market-league/internal/db"
	"github.com/market-league/internal/league"
	"github.com/market-league/internal/portfolio"
	"github.com/market-league/internal/stock"
	"github.com/market-league/internal/trade"
	"github.com/market-league/internal/user"
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

	// Portfolio routes
	portfolioRepo := portfolio.NewPortfolioRepository(database)
	portfolioService := portfolio.NewPortfolioService(portfolioRepo)
	portfolioHandler := portfolio.NewPortfolioHandler(portfolioService)

	portfolioRoutes := router.Group("/api/portfolio")
	{
		portfolioRoutes.POST("/:portfolioID", portfolioHandler.GetPortfolio)             // Fetch a portfolio by ID
		portfolioRoutes.POST("/user-portfolio", portfolioHandler.GetUserPortfolio)       // Fetch user's portfolio in a league
		portfolioRoutes.POST("/create", portfolioHandler.CreatePortfolio)                // Create a portfolio
		portfolioRoutes.POST("/add-stock", portfolioHandler.AddStockToPortfolio)         // Add a stock to a portfolio
		portfolioRoutes.POST("/remove-stock", portfolioHandler.RemoveStockFromPortfolio) // Remove a stock from a portfolio
	}

	// Stocks routes
	stockRepo := stock.NewStockRepository(database)
	stockService := stock.NewStockService(stockRepo)
	stockHandler := stock.NewStockHandler(stockService)

	stockRoutes := router.Group("/api/stocks")
	{
		stockRoutes.GET("/:stockID/price", stockHandler.GetPrice)                // Fetch stock price by ID
		stockRoutes.GET("/:stockID/price-history", stockHandler.GetPriceHistory) // Fetch price history by ID
		stockRoutes.POST("/", stockHandler.CreateStock)                          // Create a new stock
		stockRoutes.PUT("/:stockID/update-price", stockHandler.UpdateStockPrice) // Update stock price by ID
	}

	// Trades routes
	tradeRepo := trade.NewTradeRepository(database)
	tradeService := trade.NewTradeService(tradeRepo)
	tradeHandler := trade.NewTradeHandler(tradeService)

	tradeRoutes := router.Group("/api/trades")
	{
		tradeRoutes.POST("/create", tradeHandler.CreateTrade)                         // Create a new trade
		tradeRoutes.GET("/user/:userID", tradeHandler.GetTradesByUser)                // Get all trades by a specific user
		tradeRoutes.GET("/portfolio/:portfolioID", tradeHandler.GetTradesByPortfolio) // Get all trades by a specific portfolio
	}

	userRepo := user.NewUserRepository(database)
	userService := user.NewUserService(userRepo)
	userHandler := user.NewUserHandler(userService)

	userRoutes := router.Group("/api/users")
	{
		userRoutes.GET("/:userID", userHandler.GetUserByID)                  // Fetch a user by their ID
		userRoutes.PUT("/:userID/update", userHandler.UpdateUser)            // Update user details
		userRoutes.GET("/:userID/leagues", userHandler.GetUserLeagues)       // Fetch all leagues a user is in
		userRoutes.GET("/:userID/portfolios", userHandler.GetUserPortfolios) // Fetch all portfolios a user has
	}

	// League routes
	leagueRepo := league.NewLeagueRepository(database)
	leagueService := league.NewLeagueService(leagueRepo, userRepo, portfolioRepo)
	leagueHandler := league.NewLeagueHandler(leagueService, portfolioService)

	leagueRoutes := router.Group("/api/leagues")
	{
		leagueRoutes.POST("/create-league", leagueHandler.CreateLeague)         // Create League
		leagueRoutes.POST("/add-user-to-league", leagueHandler.AddUserToLeague) // Add Users to League
		leagueRoutes.POST("/details", leagueHandler.GetLeagueDetails)           // Get League Details
		leagueRoutes.POST("/leaderboard", leagueHandler.GetLeaderboard)         // Get League Leaderboard
	}

}
