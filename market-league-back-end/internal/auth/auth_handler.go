package auth

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/market-league/internal/user"
)

// AuthHandler defines a struct for handling authentication requests.
type AuthHandler struct {
	service AuthService
}

// NewAuthHandler initializes and returns an AuthHandler instance.
func NewAuthHandler(service AuthService) *AuthHandler {
	return &AuthHandler{
		service: service,
	}
}

// Signup handles user registration requests.
func (h *AuthHandler) Signup(c *gin.Context) {
	var newUser user.User

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
		Email    string `json:"email" binding:"required"`
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
	user, err := h.service.Login(loginDetails.Email, loginDetails.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Invalid email or password",
		})
		return
	}

	// extract username
	username := user.Name

	// Respond with the authenticated user (for simplicity, you could return a token in a real-world app)
	c.JSON(http.StatusOK, gin.H{
		"message":  "Login successful",
		"username": username,
	})
}
