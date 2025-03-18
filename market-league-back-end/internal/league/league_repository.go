package league

import (
	"fmt"
	"log"

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

// CreateLeaguePlayer inserts a new LeaguePlayer record into the database.
func (r *LeagueRepository) CreateLeaguePlayer(lp *models.LeaguePlayer) error {
	return r.db.Create(lp).Error
}

// AddUserToLeague adds a user to a specific league by creating a record in the User_Leagues table.
func (r *LeagueRepository) AddUserToLeague(userID, leagueID uint) error {
	// Fetch the league to ensure it exists
	var league models.League
	if err := r.db.First(&league, leagueID).Error; err != nil {
		return fmt.Errorf("failed to find league: %w", err)
	}

	// Fetch the user to ensure it exists
	var user models.User
	if err := r.db.First(&user, userID).Error; err != nil {
		return fmt.Errorf("failed to find user: %w", err)
	}

	// Check if the user is already in the league by querying the join table
	var count int64
	if err := r.db.Table("user_leagues").Where("user_id = ? AND league_id = ?", userID, leagueID).Count(&count).Error; err != nil {
		return fmt.Errorf("failed to check if user is already in league: %w", err)
	}

	if count > 0 {
		log.Println("User already in league")
		return fmt.Errorf("user already in league")
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
			TotalValue: portfolio.Points,
		})
	}

	return leaderboard, nil
}

// GetLeagueDetails retrieves details for a specific league by ID.
func (r *LeagueRepository) GetLeagueDetails(leagueID uint) (*models.League, error) {
	var league models.League
	err := r.db.
		Preload("LeaguePlayers").
		Preload("Users").
		Where("id = ?", leagueID).First(&league).Error
	return &league, err
}

// RemovePortfolioStocksByLeagueID removes stocks associated with portfolios in a league
func (r *LeagueRepository) RemovePortfolioStocksByLeagueID(tx *gorm.DB, leagueID uint) error {
	return tx.Exec(`
        DELETE FROM portfolio_stocks
        WHERE portfolio_id IN (SELECT id FROM portfolios WHERE league_id = ?)`, leagueID).Error
}

// RemovePortfoliosByLeagueID removes portfolios associated with a league
func (r *LeagueRepository) RemovePortfoliosByLeagueID(tx *gorm.DB, leagueID uint) error {
	return tx.Where("league_id = ?", leagueID).Delete(&models.Portfolio{}).Error
}

// RemoveTradesByLeagueID removes trades associated with a league
func (r *LeagueRepository) RemoveTradesByLeagueID(tx *gorm.DB, leagueID uint) error {
	return tx.Where("league_id = ?", leagueID).Delete(&models.Trade{}).Error
}

func (r *LeagueRepository) RemovePortfolioPointsHistoryByLeagueID(tx *gorm.DB, leagueID uint) error {
	return tx.Exec("DELETE FROM portfolio_points_histories WHERE portfolio_id IN (SELECT id FROM portfolios WHERE league_id = ?)", leagueID).Error
}

// RemoveLeaguePortfolioByLeagueID removes league portfolios for a specific league
func (r *LeagueRepository) RemoveLeaguePortfolioByLeagueID(tx *gorm.DB, leagueID uint) error {
	return tx.Exec("DELETE FROM league_portfolios WHERE league_id = ?", leagueID).Error
}

// RemoveOwnershipHistoriesByLeagueID removes ownership history for a specific league
func (r *LeagueRepository) RemoveOwnershipHistoriesByLeagueID(tx *gorm.DB, leagueID uint) error {
	return tx.Exec("DELETE FROM ownership_histories WHERE portfolio_id IN (SELECT id FROM portfolios WHERE league_id = ?)", leagueID).Error
}

// RemovePortfolioStocksByLeagueID removes all portfolio stocks associated with a league
func (r *LeagueRepository) RemoveLeaguePortfolioStocksByLeagueID(tx *gorm.DB, leagueID uint) error {
	query := `
		DELETE FROM league_portfolio_stocks
		WHERE league_portfolio_id IN (
			SELECT id FROM league_portfolios WHERE league_id = ?
		)
	`
	return tx.Exec(query, leagueID).Error
}

// RemoveUserLeaguesByLeagueID removes user-league associations for a league
func (r *LeagueRepository) RemoveUserLeaguesByLeagueID(tx *gorm.DB, leagueID uint) error {
	return tx.Exec("DELETE FROM user_leagues WHERE league_id = ?", leagueID).Error
}

// RemoveLeaguePlayerByLeagueID removes user-league associations for a league
func (r *LeagueRepository) RemoveLeaguePlayerByLeagueID(tx *gorm.DB, leagueID uint) error {
	return tx.Exec("DELETE FROM league_players WHERE league_id = ?", leagueID).Error
}

// RemoveLeague removes the league itself
func (r *LeagueRepository) RemoveLeague(tx *gorm.DB, leagueID uint) error {
	return tx.Where("id = ?", leagueID).Delete(&models.League{}).Error
}

// QueueUpPlayer updates the player's draft status to "ready"
func (r *LeagueRepository) QueueUpPlayer(leagueID uint, playerID uint) error {
	var leaguePlayer models.LeaguePlayer
	if err := r.db.Where("league_id = ? AND player_id = ?", leagueID, playerID).First(&leaguePlayer).Error; err != nil {
		return err
	}

	leaguePlayer.DraftStatus = models.DraftReady
	return r.db.Save(&leaguePlayer).Error
}

// AllPlayersReady checks if every player in the league is ready
func (r *LeagueRepository) AllPlayersReady(leagueID uint) (bool, error) {
	var count int64
	err := r.db.Model(&models.LeaguePlayer{}).
		Where("league_id = ? AND draft_status != ?", leagueID, models.DraftReady).
		Count(&count).Error
	return count == 0, err
}

// GetLeague retrieves a league along with its players
func (r *LeagueRepository) GetLeague(leagueID uint) (*models.League, error) {
	var league models.League
	err := r.db.Preload("LeaguePlayers").First(&league, leagueID).Error
	return &league, err
}

// UpdateLeague updates a league record
func (r *LeagueRepository) UpdateLeague(league *models.League) error {
	return r.db.Save(league).Error
}
