package auth

import (
	"github.com/market-league/internal/models"
	"gorm.io/gorm"
)

// AuthRepository handles database operations related to authentication.
type AuthRepository struct {
	db *gorm.DB
}

// NewAuthRepository creates a new instance of AuthRepository.
func NewAuthRepository(db *gorm.DB) *AuthRepository {
	return &AuthRepository{db: db}
}

// CreateUser creates a new user in the database.
func (r *AuthRepository) CreateUser(newUser *models.User) error {
	return r.db.Create(newUser).Error
}

// FindUserByUsername retrieves a user by their username.
func (r *AuthRepository) FindUserByUsername(username string) (*models.User, error) {
	var foundUser models.User
	err := r.db.Where("username = ?", username).First(&foundUser).Error
	if err != nil {
		return nil, err
	}
	return &foundUser, nil
}

// GetUserByID retrieves a user by their ID.
func (r *AuthRepository) GetUserByID(userID uint) (*models.User, error) {
	var user models.User
	err := r.db.First(&user, userID).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}
