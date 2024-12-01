package league

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	league_portfolio "github.com/market-league/internal/leagueportfolio"
	"github.com/market-league/internal/portfolio"
)

// LeagueHandler defines the HTTP handler for league-related operations.
type LeagueHandler struct {
	service                *LeagueService
	portfolioService       *portfolio.PortfolioService
	leaguePortfolioService *league_portfolio.LeaguePortfolioService
}

// NewLeagueHandler creates a new instance of LeagueHandler.
func NewLeagueHandler(service *LeagueService,
	portfolioService *portfolio.PortfolioService,
	leaguePortfolioService *league_portfolio.LeaguePortfolioService) *LeagueHandler {
	return &LeagueHandler{
		service:                service,
		portfolioService:       portfolioService,
		leaguePortfolioService: leaguePortfolioService,
	}
}

// CreateLeague handles the creation of a new league.
func (h *LeagueHandler) CreateLeague(c *gin.Context) {
	var leagueRequest struct {
		LeagueName string `json:"league_name" binding:"required"`
		OwnerUser  uint   `json:"owner_user" binding:"required"`
		EndDate    string `json:"end_date" binding:"required"`
	}

	if err := c.ShouldBindJSON(&leagueRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Set the start date to the current date and time
	startDate := time.Now().Format(time.RFC3339)

	// Pass the values to the service to create the league
	league, err := h.service.CreateLeague(leagueRequest.LeagueName, leagueRequest.OwnerUser, startDate, leagueRequest.EndDate)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create league"})
		return
	}

	// Create a portfolio for the user in the league
	portfolio, err := h.portfolioService.CreatePortfolio(leagueRequest.OwnerUser, league.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Create a league portfolio using the new LeaguePortfolioService
	leaguePortfolio, err := h.leaguePortfolioService.CreateLeaguePortfolio(league.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Construct response with sanitized user details
	response := gin.H{
		"message":         "League successfully created",
		"league":          league,
		"userPortfolio":   portfolio,
		"leaguePortfolio": leaguePortfolio,
	}

	c.JSON(http.StatusCreated, response)
}

// AddUserToLeague handles adding a user to a league.
func (h *LeagueHandler) AddUserToLeague(c *gin.Context) {
	var request struct {
		UserID   uint `json:"user_id" binding:"required"`
		LeagueID uint `json:"league_id" binding:"required"`
	}

	// Bind JSON input to the request struct
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	// Call the service to add the user to the league
	err := h.service.AddUserToLeague(request.UserID, request.LeagueID)
	if err != nil {
		if err.Error() == "user already in league" {
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Create a portfolio for the user in the league
	portfolio, err := h.portfolioService.CreatePortfolio(request.UserID, request.LeagueID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Return success response with portfolio details
	c.JSON(http.StatusOK, gin.H{
		"message":   "User successfully added to league",
		"portfolio": portfolio,
	})
}

// GetLeagueDetails handles fetching the details of a specific league.
func (h *LeagueHandler) GetLeagueDetails(c *gin.Context) {

	var request struct {
		LeagueID uint `json:"league_id" binding:"required"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	league, users, err := h.service.GetLeagueDetails(request.LeagueID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch league details"})
		return
	}

	// Construct response with sanitized user details
	response := gin.H{
		"id":          league.ID,
		"league_name": league.LeagueName,
		"start_date":  league.StartDate,
		"end_date":    league.EndDate,
		"users":       users,
	}

	c.JSON(http.StatusOK, response)
}

// GetLeaderboard handles fetching the leaderboard for a specific league.
func (h *LeagueHandler) GetLeaderboard(c *gin.Context) {
	var request struct {
		LeagueID uint `json:"league_id" binding:"required"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Pass the LeagueID and the PortfolioService to the service method
	leaderboard, err := h.service.GetLeaderboard(request.LeagueID, h.portfolioService)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch leaderboard"})
		return
	}

	c.JSON(http.StatusOK, leaderboard)
}

// RemoveLeague handles the removal of a league and all associated records
func (h *LeagueHandler) RemoveLeague(c *gin.Context) {
	var request struct {
		LeagueID uint `json:"league_id" binding:"required"`
	}

	// Bind JSON input to the request struct
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	// Call the service to remove the league
	if err := h.service.RemoveLeague(request.LeagueID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Return success response
	c.JSON(http.StatusOK, gin.H{"message": "League and associated data removed successfully"})
}
