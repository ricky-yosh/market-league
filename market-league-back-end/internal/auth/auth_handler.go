package auth

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/market-league/internal/league"
	"github.com/market-league/internal/models"
)

// AuthHandler defines a struct for handling authentication requests.
type AuthHandler struct {
	service *AuthService
}

// NewAuthHandler creates a new instance of AuthHandler.
func NewAuthHandler(service *AuthService) *AuthHandler {
	return &AuthHandler{service: service}
}

// Signup handles user registration requests.
func (h *AuthHandler) Signup(c *gin.Context) {
	var newUser models.User

	// Bind incoming JSON to the newUser struct
	if err := c.ShouldBindJSON(&newUser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request body",
		})
		return
	}

	// Call the service to sign up the user
	err := h.service.Signup(&newUser)
	if err != nil {
		c.JSON(http.StatusConflict, gin.H{
			"error": err.Error(),
		})
		return
	}

	// Respond with success
	c.JSON(http.StatusCreated, gin.H{
		"message": "User registered successfully",
	})
}

// Login handles user authentication requests.
func (h *AuthHandler) Login(c *gin.Context) {
	var loginDetails struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	// Bind incoming JSON to the loginDetails struct
	if err := c.ShouldBindJSON(&loginDetails); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request body",
		})
		return
	}

	// Call the service to authenticate the user
	user, err := h.service.Login(loginDetails.Username, loginDetails.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Invalid username or password",
		})
		return
	}

	// Generate JWT using the service
	jwtToken, err := h.service.GenerateJWT(user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate JWT"})
		return
	}

	// Respond with the authenticated user and token
	c.JSON(http.StatusOK, gin.H{
		"token":    jwtToken,
		"username": user.Username,
		"message":  "Login successful",
	})
}

// GetUserFromToken extracts user information based on the provided JWT token.
func (h *AuthHandler) GetUserFromToken(c *gin.Context) {
	// Extract the token from the Authorization header
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Authorization header is missing"})
		return
	}

	// Parse the token
	tokenString := authHeader[len("Bearer "):] // Remove "Bearer " from the header
	userID, err := h.service.ParseJWT(tokenString)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
		return
	}

	// Fetch the user based on the userID
	user, err := h.service.GetUserByID(userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// Sanitize and respond with user data
	sanitizedUser := league.SanitizeUser(*user)
	c.JSON(http.StatusOK, sanitizedUser)
}
