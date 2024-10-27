package auth

// Import necessary packages
import (
	"errors"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/market-league/internal/models"
	"golang.org/x/crypto/bcrypt"
)

// AuthService handles authentication logic.
type AuthService struct {
	repo         *AuthRepository
	jwtSecretKey []byte
}

// NewAuthService creates a new instance of AuthService.
func NewAuthService(repo *AuthRepository, secretKey string) *AuthService {
	return &AuthService{
		repo:         repo,
		jwtSecretKey: []byte(secretKey),
	}
}

// Signup handles user registration, encrypting the password, and saving the user in the repository.
func (s *AuthService) Signup(newUser *models.User) error {
	// Check if the user already exists
	existingUser, err := s.repo.FindUserByUsername(newUser.Username)
	if err == nil && existingUser != nil {
		return errors.New("user with this email already exists")
	}

	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newUser.Password), bcrypt.DefaultCost)
	if err != nil {
		return errors.New("failed to hash password")
	}
	newUser.Password = string(hashedPassword)

	// Save the new user
	return s.repo.CreateUser(newUser)
}

// Login authenticates the user by verifying the email and password.
func (s *AuthService) Login(username string, password string) (*models.User, error) {
	// Find user by username
	existingUser, err := s.repo.FindUserByUsername(username)
	if err != nil {
		return nil, errors.New("user not found")
	}

	// Verify password
	err = bcrypt.CompareHashAndPassword([]byte(existingUser.Password), []byte(password))
	if err != nil {
		return nil, errors.New("incorrect password")
	}

	return existingUser, nil
}

// GenerateJWT generates a JWT for the authenticated user.
func (s *AuthService) GenerateJWT(userID uint) (string, error) {
	userIDStr := strconv.FormatUint(uint64(userID), 10)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		Subject:   userIDStr,
		ExpiresAt: time.Now().Add(time.Hour * 24).Unix(), // Set the expiration time (e.g., 24 hours)
	})

	// Sign the token with the secret key
	tokenString, err := token.SignedString(s.jwtSecretKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// ParseJWT parses the provided token and extracts the user ID.
func (s *AuthService) ParseJWT(tokenString string) (uint, error) {
	// Parse the token
	token, err := jwt.ParseWithClaims(tokenString, &jwt.StandardClaims{}, func(token *jwt.Token) (interface{}, error) {
		return s.jwtSecretKey, nil
	})

	if err != nil || !token.Valid {
		return 0, errors.New("invalid token")
	}

	// Extract user ID from the token claims
	claims, ok := token.Claims.(*jwt.StandardClaims)
	if !ok {
		return 0, errors.New("invalid token claims")
	}

	// Convert user ID from string to uint
	userID, err := strconv.ParseUint(claims.Subject, 10, 32)
	if err != nil {
		return 0, errors.New("invalid user ID in token")
	}

	return uint(userID), nil
}

// GetUserByID retrieves a user by their ID.
func (s *AuthService) GetUserByID(userID uint) (*models.User, error) {
	return s.repo.GetUserByID(userID)
}
