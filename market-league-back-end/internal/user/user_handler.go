package user

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/market-league/internal/models"
)

// UserHandler defines the HTTP handler for user-related operations.
type UserHandler struct {
	service *UserService
}

// NewUserHandler creates a new instance of UserHandler.
func NewUserHandler(service *UserService) *UserHandler {
	return &UserHandler{service: service}
}

// GetUserByID handles fetching a user by their ID.
func (h *UserHandler) GetUserByID(c *gin.Context) {
	// Parse the user ID from the URL parameter
	userID, err := strconv.ParseUint(c.Param("userID"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	// Get the user by ID using the service
	user, err := h.service.GetUserByID(uint(userID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, user)
}

// UpdateUser handles updating user details.
func (h *UserHandler) UpdateUser(c *gin.Context) {
	// Parse the user ID from the URL parameter
	userID, err := strconv.ParseUint(c.Param("userID"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	// Bind the JSON data to the user model
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Update the user using the service
	err = h.service.UpdateUser(uint(userID), &user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User updated successfully"})
}

// GetUserLeagues handles fetching all leagues a user is in.
func (h *UserHandler) GetUserLeagues(c *gin.Context) {
	// Parse the user ID from the URL parameter
	userID, err := strconv.ParseUint(c.Param("userID"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	// Get the user's leagues using the service
	leagues, err := h.service.GetUserLeagues(uint(userID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch leagues"})
		return
	}

	c.JSON(http.StatusOK, leagues)
}

// GetUserPortfolios handles fetching all portfolios a user has.
func (h *UserHandler) GetUserPortfolios(c *gin.Context) {
	// Parse the user ID from the URL parameter
	userID, err := strconv.ParseUint(c.Param("userID"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	// Get the user's portfolios using the service
	portfolios, err := h.service.GetUserPortfolios(uint(userID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch portfolios"})
		return
	}

	c.JSON(http.StatusOK, portfolios)
}
