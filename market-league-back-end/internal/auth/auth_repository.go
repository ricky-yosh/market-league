package auth

import (
	"github.com/market-league/internal/user"
	"gorm.io/gorm"
)

// AuthRepository is an interface that defines the methods for interacting with the user data in the database.
type AuthRepository interface {
	CreateUser(newUser *user.User) error
	FindUserByEmail(email string) (*user.User, error)
}

// authRepository is a struct that implements the AuthRepository interface.
type authRepository struct {
	db *gorm.DB
}

// NewAuthRepository creates a new instance of AuthRepository with a given database connection.
func NewAuthRepository(db *gorm.DB) AuthRepository {
	return &authRepository{db: db}
}

// CreateUser creates a new user in the database.
func (r *authRepository) CreateUser(newUser *user.User) error {
	return r.db.Create(newUser).Error
}

// GetUserByEmail retrieves a user by their email.
func (r *authRepository) FindUserByEmail(email string) (*user.User, error) {
	var foundUser user.User
	err := r.db.Where("email = ?", email).First(&foundUser).Error
	if err != nil {
		return nil, err
	}
	return &foundUser, nil
}
