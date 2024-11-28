package user

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// UserHandler defines the HTTP handler for user-related operations.
type UserHandler struct {
	service *UserService
}

// NewUserHandler creates a new instance of UserHandler.
func NewUserHandler(service *UserService) *UserHandler {
	return &UserHandler{service: service}
}

// GetUserByID fetches user information based on filter criteria.
func (h *UserHandler) GetUserByID(c *gin.Context) {
	var request struct {
		UserID uint `json:"user_id" binding:"required"` // User ID to fetch information for
	}

	// Bind the JSON request data to the struct
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Call the service to get the user info based on the filter criteria
	userInfo, err := h.service.GetUserByID(request.UserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch user information"})
		return
	}

	c.JSON(http.StatusOK, userInfo)
}

// GetUserLeagues handles requests to retrieve leagues that a user is a member of.
func (h *UserHandler) GetUserLeagues(c *gin.Context) {
	var request struct {
		UserID uint `json:"user_id" binding:"required"`
	}

	// Bind the incoming JSON to the request struct
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	// Call the service to get user leagues
	leagues, err := h.service.GetUserLeagues(request.UserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Respond with the list of leagues
	c.JSON(http.StatusOK, gin.H{"leagues": leagues})
}

// GetUserTrades handles requests to retrieve trades that a user is involved in within a specific league.
func (h *UserHandler) GetUserTrades(c *gin.Context) {
	var request struct {
		UserID   uint `json:"user_id" binding:"required"`
		LeagueID uint `json:"league_id" binding:"required"`
	}

	// Bind the incoming JSON to the request struct
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	// Call the service to get user trades for the specified league
	trades, err := h.service.GetUserTrades(request.UserID, request.LeagueID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Respond with the list of trades
	c.JSON(http.StatusOK, gin.H{"trades": trades})
}

// GetUserPortfolios handles requests to retrieve portfolios that a user is associated with.
func (h *UserHandler) GetUserPortfolios(c *gin.Context) {
	var request struct {
		UserID uint `json:"user_id" binding:"required"`
	}

	// Bind the incoming JSON to the request struct
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	// Call the service to get user portfolios
	portfolios, err := h.service.GetUserPortfolios(request.UserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Respond with the list of portfolios
	c.JSON(http.StatusOK, gin.H{"portfolios": portfolios})
}
