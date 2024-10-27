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
		UserID     uint `json:"user_id" binding:"required"` // User ID to fetch information for
		Portfolios bool `json:"portfolios"`                 // Whether to include user's portfolios
		Leagues    bool `json:"leagues"`                    // Whether to include user's leagues
		Trades     bool `json:"trades"`                     // Whether to include user's trades
	}

	// Bind the JSON request data to the struct
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Call the service to get the user info based on the filter criteria
	userInfo, err := h.service.GetUserByID(request.UserID, request.Portfolios, request.Leagues, request.Trades)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch user information"})
		return
	}

	c.JSON(http.StatusOK, userInfo)
}
