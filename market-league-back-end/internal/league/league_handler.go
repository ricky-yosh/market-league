package league

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// LeagueHandler defines the HTTP handler for league-related operations.
type LeagueHandler struct {
	service *LeagueService
}

// NewLeagueHandler creates a new instance of LeagueHandler.
func NewLeagueHandler(service *LeagueService) *LeagueHandler {
	return &LeagueHandler{service: service}
}

// CreateLeague handles the creation of a new league.
func (h *LeagueHandler) CreateLeague(c *gin.Context) {
	var leagueRequest struct {
		LeagueName string `json:"league_name" binding:"required"`
		StartDate  string `json:"start_date" binding:"required"`
		EndDate    string `json:"end_date" binding:"required"`
	}

	if err := c.ShouldBindJSON(&leagueRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	league, err := h.service.CreateLeague(leagueRequest.LeagueName, leagueRequest.StartDate, leagueRequest.EndDate)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create league"})
		return
	}

	c.JSON(http.StatusOK, league)
}

// AddUserToLeague handles adding a user to a league.
func (h *LeagueHandler) AddUserToLeague(c *gin.Context) {
	var request struct {
		UserID   uint `json:"user_id" binding:"required"`
		LeagueID uint `json:"league_id" binding:"required"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.service.AddUserToLeague(request.UserID, request.LeagueID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add user to league"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User added to league successfully"})
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

	league, err := h.service.GetLeagueDetails(request.LeagueID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch league details"})
		return
	}

	c.JSON(http.StatusOK, league)
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

	leaderboard, err := h.service.GetLeaderboard(request.LeagueID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch leaderboard"})
		return
	}

	c.JSON(http.StatusOK, leaderboard)
}
