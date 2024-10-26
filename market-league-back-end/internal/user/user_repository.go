package user

import (
	"fmt"

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

// GetUserByID fetches basic user details by ID.
func (r *UserRepository) GetUserByID(userID uint) (*models.User, error) {
	var user models.User
	if err := r.db.First(&user, userID).Error; err != nil {
		return nil, fmt.Errorf("failed to find user with ID %d: %w", userID, err)
	}
	return &user, nil
}

// GetUserByUsername finds a user by their username in the database.
func (r *UserRepository) GetUserByUsername(username string) (*models.User, error) {
	var user models.User
	if err := r.db.Where("username = ?", username).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

// GetUserPortfolios fetches all portfolios for a given user.
func (r *UserRepository) GetUserPortfolios(userID uint) ([]models.Portfolio, error) {
	var portfolios []models.Portfolio
	if err := r.db.Where("user_id = ?", userID).Find(&portfolios).Error; err != nil {
		return nil, fmt.Errorf("failed to find portfolios for user with ID %d: %w", userID, err)
	}
	return portfolios, nil
}

// GetUserLeagues fetches all leagues that the user is part of.
func (r *UserRepository) GetUserLeagues(userID uint) ([]models.League, error) {
	var leagues []models.League
	if err := r.db.Model(&models.User{}).Where("id = ?", userID).Association("Leagues").Find(&leagues); err != nil {
		return nil, fmt.Errorf("failed to find leagues for user with ID %d: %w", userID, err)
	}
	return leagues, nil
}

// GetUserTrades fetches all trades involving a given user.
func (r *UserRepository) GetUserTrades(userID uint) ([]models.Trade, error) {
	var trades []models.Trade
	if err := r.db.Where("player1_id = ? OR player2_id = ?", userID, userID).Find(&trades).Error; err != nil {
		return nil, fmt.Errorf("failed to find trades for user with ID %d: %w", userID, err)
	}
	return trades, nil
}
