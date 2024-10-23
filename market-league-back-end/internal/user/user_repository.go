package user

import (
	"github.com/market-league/internal/models"
	"gorm.io/gorm"
)

// UserRepository provides access to user-related operations in the database.
type UserRepository struct {
	db *gorm.DB
}

// NewUserRepository creates a new instance of UserRepository.
func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

// GetUserByID fetches a user by their ID.
func (r *UserRepository) GetUserByID(userID uint) (*models.User, error) {
	var user models.User
	err := r.db.First(&user, userID).Error
	return &user, err
}

// GetUserByUsername finds a user by their username in the database.
func (r *UserRepository) GetUserByUsername(username string) (*models.User, error) {
	var user models.User
	if err := r.db.Where("username = ?", username).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

// UpdateUser updates user details in the database.
func (r *UserRepository) UpdateUser(userID uint, user *models.User) error {
	return r.db.Model(&models.User{}).Where("id = ?", userID).Updates(user).Error
}

// GetUserLeagues fetches all leagues that a user is in.
func (r *UserRepository) GetUserLeagues(userID uint) ([]models.League, error) {
	var user models.User
	err := r.db.Preload("Leagues").First(&user, userID).Error
	return user.Leagues, err
}

// GetUserPortfolios fetches all portfolios that belong to a user.
func (r *UserRepository) GetUserPortfolios(userID uint) ([]models.Portfolio, error) {
	var portfolios []models.Portfolio
	err := r.db.Where("user_id = ?", userID).Find(&portfolios).Error
	return portfolios, err
}
