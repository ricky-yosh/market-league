package auth

import (
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/market-league/internal/models"
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

var jwtSecretKey = []byte(os.Getenv("JWT_KEY"))

// GenerateJWT generates a JWT for the authenticated user
func GenerateJWT(userID uint) (string, error) {
	userIDStr := strconv.FormatUint(uint64(userID), 10)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		Subject:   userIDStr,
		ExpiresAt: time.Now().Add(time.Hour * 24).Unix(), // Set the expiration time (e.g., 24 hours)
	})

	// Sign the token with your secret key
	tokenString, err := token.SignedString(jwtSecretKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
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
	// Generate JWT
	jwtToken, err := GenerateJWT(user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate JWT"})
		return
	}

	// extract username
	username := user.Username

	// Respond with the authenticated user (for simplicity, you could return a token in a real-world app)
	c.JSON(http.StatusOK, gin.H{
		"token":    jwtToken,
		"username": username,
		"message":  "Login successful",
	})
}
