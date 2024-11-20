package api

import (
	"os"

	"github.com/gin-gonic/gin"
	"github.com/market-league/internal/auth"
	"github.com/market-league/internal/db"
	"github.com/market-league/internal/league"
	"github.com/market-league/internal/portfolio"
	"github.com/market-league/internal/stock"
	"github.com/market-league/internal/trade"
	"github.com/market-league/internal/user"

	// "github.com/market-league/internal/scheduler"
	"github.com/market-league/internal/services"
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
	authService := auth.NewAuthService(authRepo, os.Getenv("JWT_KEY"))
	authHandler := auth.NewAuthHandler(authService)

	// Auth routes
	authRoutes := router.Group("/api/auth")
	{
		authRoutes.POST("/signup", authHandler.Signup)
		authRoutes.POST("/login", authHandler.Login)
		authRoutes.GET("/user-from-token", authHandler.GetUserFromToken) // New endpoint to get user from JWT
	}

	// Portfolio routes
	portfolioRepo := portfolio.NewPortfolioRepository(database)
	portfolioService := portfolio.NewPortfolioService(portfolioRepo)
	portfolioHandler := portfolio.NewPortfolioHandler(portfolioService)

	portfolioRoutes := router.Group("/api/portfolio")
	{
		portfolioRoutes.POST("/create-portfolio", portfolioHandler.CreatePortfolio)      // Create a portfolio
		portfolioRoutes.POST("/portfolio-with-id", portfolioHandler.GetPortfolioWithID)  // Fetch a portfolio by ID
		portfolioRoutes.POST("/league-portfolio", portfolioHandler.GetLeaguePortfolio)   // Fetch user's portfolio in a league
		portfolioRoutes.POST("/add-stock", portfolioHandler.AddStockToPortfolio)         // Add a stock to a portfolio
		portfolioRoutes.POST("/remove-stock", portfolioHandler.RemoveStockFromPortfolio) // Remove a stock from a portfolio
	}

	// Stocks routes
	stockRepo := stock.NewStockRepository(database)
	stockService := stock.NewStockService(stockRepo)
	stockHandler := stock.NewStockHandler(stockService)

	stockRoutes := router.Group("/api/stocks")
	{
		stockRoutes.POST("/create-stock", stockHandler.CreateStock)                // Create a new stock
		stockRoutes.POST("/create-stocks", stockHandler.CreateMultipleStocks)      // Create multiple stocks
		stockRoutes.POST("/stock-price", stockHandler.GetPrice)                    // Fetch stock price by ID
		stockRoutes.POST("/update-stock-price", stockHandler.UpdateStockPrice)     // Update stock price by ID
		stockRoutes.POST("/price-history", stockHandler.GetPriceHistory)           // Fetch price history by ID
		stockRoutes.POST("/update-price-history", stockHandler.UpdatePriceHistory) // Update price history by ID
	}

	userRepo := user.NewUserRepository(database)
	userService := user.NewUserService(userRepo)
	userHandler := user.NewUserHandler(userService)

	userRoutes := router.Group("/api/users")
	{
		userRoutes.POST("/user-info", userHandler.GetUserByID)
		userRoutes.POST("/user-leagues", userHandler.GetUserLeagues)
		userRoutes.POST("/user-trades", userHandler.GetUserTrades)
		userRoutes.POST("/user-portfolios", userHandler.GetUserPortfolios)
		// userRoutes.POST("/update-user", userHandler.GetUserByID)
	}

	// Trades routes
	tradeRepo := trade.NewTradeRepository(database)
	tradeService := trade.NewTradeService(tradeRepo, stockRepo, portfolioRepo, userRepo)
	tradeHandler := trade.NewTradeHandler(tradeService)

	tradeRoutes := router.Group("/api/trades")
	{
		tradeRoutes.POST("/create-trade", tradeHandler.CreateTrade) // Create a new trade
		tradeRoutes.POST("/confirm-trade", tradeHandler.ConfirmTrade)
		tradeRoutes.POST("/get-trades", tradeHandler.GetTrades)
	}

	// League routes
	leagueRepo := league.NewLeagueRepository(database)
	leagueService := league.NewLeagueService(leagueRepo, userRepo, portfolioRepo)
	leagueHandler := league.NewLeagueHandler(leagueService, portfolioService)

	leagueRoutes := router.Group("/api/leagues")
	{
		leagueRoutes.POST("/create-league", leagueHandler.CreateLeague)         // Create League
		leagueRoutes.POST("/remove-league", leagueHandler.RemoveLeague)         // Remove League
		leagueRoutes.POST("/add-user-to-league", leagueHandler.AddUserToLeague) // Add Users to League
		leagueRoutes.POST("/details", leagueHandler.GetLeagueDetails)           // Get League Details
		leagueRoutes.POST("/leaderboard", leagueHandler.GetLeaderboard)         // Get League Leaderboard
	}

	// go scheduler.StartDailyTask()
	router.GET("/api/services/stock-api", func(c *gin.Context) {
		quote, err := services.GetTestStock()
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		c.JSON(200, quote)
	})

}
