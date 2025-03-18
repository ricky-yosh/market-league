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

// GetUserLeagues retrieves all leagues that the user is a member of.
func (r *UserRepository) GetUserLeagues(userID uint) ([]models.League, error) {
	var user models.User
	err := r.db.
		Preload("Leagues").
		Preload("Leagues.LeaguePlayers").
		First(&user, userID).Error
	if err != nil {
		return nil, fmt.Errorf("failed to find user: %w", err)
	}

	return user.Leagues, nil
}

// GetUserTrades retrieves all trades involving a specific user in a specific league.
func (r *UserRepository) GetUserTrades(userID uint, leagueID uint) ([]models.Trade, error) {
	var trades []models.Trade
	// Fetch trades where the user is involved in either side and the league matches
	err := r.db.Where("(player1_id = ? OR player2_id = ?) AND league_id = ?", userID, userID, leagueID).Find(&trades).Error
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve trades for user in league %d: %w", leagueID, err)
	}
	return trades, nil
}

// GetUserPortfolios retrieves all portfolios for a specific user.
func (r *UserRepository) GetUserPortfolios(userID uint) ([]models.Portfolio, error) {
	var portfolios []models.Portfolio
	err := r.db.Where("user_id = ?", userID).Find(&portfolios).Error
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve portfolios: %w", err)
	}
	return portfolios, nil
}
