package auth

// Import necessary packages
import (
	"errors"

	"github.com/market-league/internal/models"
	"golang.org/x/crypto/bcrypt"
)

// AuthService interface defines methods related to authentication services.
type AuthService interface {
	Signup(user *models.User) error
	Login(username string, password string) (*models.User, error)
}

// authService struct implements the AuthService interface.
type authService struct {
	repo AuthRepository
}

// NewAuthService returns an instance of AuthService with the provided AuthRepository.
func NewAuthService(repo AuthRepository) AuthService {
	return &authService{
		repo: repo,
	}
}

// Signup handles user registration, encrypting the password, and saving the user in the repository.
func (s *authService) Signup(newUser *models.User) error {
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
func (s *authService) Login(username string, password string) (*models.User, error) {
	// Find user by email
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
