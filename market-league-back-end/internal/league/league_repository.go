package league

import (
	"fmt"

	"github.com/market-league/internal/models"
	"github.com/market-league/internal/portfolio"
	"gorm.io/gorm"
)

// LeagueRepository defines the interface for league-related database operations.
type LeagueRepository struct {
	db *gorm.DB
}

// NewLeagueRepository creates a new instance of LeagueRepository.
func NewLeagueRepository(db *gorm.DB) *LeagueRepository {
	return &LeagueRepository{db: db}
}

// CreateLeague creates a new league in the database.
func (r *LeagueRepository) CreateLeague(league *models.League) error {
	return r.db.Create(league).Error
}

// AddUserToLeague adds a user to a specific league by creating a record in the User_Leagues table.
func (r *LeagueRepository) AddUserToLeague(userID, leagueID uint) error {
	// Fetch the league
	var league models.League
	if err := r.db.First(&league, leagueID).Error; err != nil {
		return fmt.Errorf("failed to find league: %w", err)
	}

	// Fetch the user
	var user models.User
	if err := r.db.First(&user, userID).Error; err != nil {
		return fmt.Errorf("failed to find user: %w", err)
	}

	// Append the user to the league's Users association
	if err := r.db.Model(&league).Association("Users").Append(&user); err != nil {
		return fmt.Errorf("failed to add user to league: %w", err)
	}

	return nil
}

// GetLeaderboard retrieves the leaderboard for a given league ID.
func (r *LeagueRepository) GetLeaderboard(leagueID uint, portfolioService *portfolio.PortfolioService) ([]models.LeaderboardEntry, error) {
	// Leaderboard will consist of portfolios ordered by total value within a given league.
	var portfolios []models.Portfolio
	var leaderboard []models.LeaderboardEntry

	// Fetch portfolios related to the given league ID and preload associated users and stocks.
	err := r.db.Preload("User").Preload("Stocks").
		Where("league_id = ?", leagueID).
		Find(&portfolios).Error

	if err != nil {
		return nil, err
	}

	// Map the result into the leaderboard slice
	for _, portfolio := range portfolios {
		leaderboard = append(leaderboard, models.LeaderboardEntry{
			Username:   portfolio.User.Username,
			TotalValue: portfolioService.CalculateTotalValue(&portfolio),
		})
	}

	return leaderboard, nil
}

// GetLeagueDetails retrieves details for a specific league by ID.
func (r *LeagueRepository) GetLeagueDetails(leagueID uint) (*models.League, error) {
	var league models.League
	err := r.db.Preload("Users").Where("id = ?", leagueID).First(&league).Error
	return &league, err
}
