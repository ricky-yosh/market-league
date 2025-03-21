package api

import (
	"os"

	"github.com/gin-gonic/gin"
	"github.com/market-league/internal/auth"
	"github.com/market-league/internal/db"
	"github.com/market-league/internal/league"
	league_portfolio "github.com/market-league/internal/leagueportfolio"
	"github.com/market-league/internal/portfolio"
	"github.com/market-league/internal/stock"
	"github.com/market-league/internal/trade"
	"github.com/market-league/internal/user"
)

func RegisterRoutes(router *gin.Engine) {
	// Shared database instance
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

	// * DEPENDENCIES *

	// Initialize Portfolio Dependencies
	portfolioRepo := portfolio.NewPortfolioRepository(database)
	portfolioService := portfolio.NewPortfolioService(portfolioRepo)
	portfolioHandler := portfolio.NewPortfolioHandler(portfolioService)

	// Initialize Stock Dependencies
	stockRepo := stock.NewStockRepository(database)
	stockService := stock.NewStockService(stockRepo)
	stockHandler := stock.NewStockHandler(stockService)

	// Initialize User Dependencies
	userRepo := user.NewUserRepository(database)
	userService := user.NewUserService(userRepo)
	userHandler := user.NewUserHandler(userService)

	// Initialize Trade Dependencies
	tradeRepo := trade.NewTradeRepository(database)
	tradeService := trade.NewTradeService(tradeRepo, stockRepo, portfolioRepo, userRepo)
	tradeHandler := trade.NewTradeHandler(tradeService)

	// Initialize LeaguePortfolio Dependencies
	leaguePortfolioRepository := league_portfolio.NewLeaguePortfolioRepository(database)
	leaguePortfolioService := league_portfolio.NewLeaguePortfolioService(leaguePortfolioRepository, stockRepo, portfolioRepo)
	leaguePortfolioHandler := league_portfolio.NewLeaguePortfolioHandler(leaguePortfolioService)

	// Initialize League Dependencies
	leagueRepo := league.NewLeagueRepository(database)
	leagueService := league.NewLeagueService(leagueRepo, userRepo, portfolioRepo)
	leagueHandler := league.NewLeagueHandler(leagueService, portfolioService, leaguePortfolioService)

	webSocketHandler := NewWebSocketHandler(
		portfolioHandler,
		stockHandler,
		userHandler,
		tradeHandler,
		leaguePortfolioHandler,
		leagueHandler,
	)

	// WebSocket endpoint
	router.GET("/ws", webSocketHandler.HandleWebSocket) // Route for WebSocket connection

	// Initialize the scheduler and start it
	scheduler := &Scheduler{
		db:           database,
		StockService: stockService,
		stockRepo:    stockRepo,
	}
	scheduler.StartDailyTask()

}
