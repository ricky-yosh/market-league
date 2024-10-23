package league

import (
	"github.com/market-league/internal/models"
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
func (r *LeagueRepository) AddUserToLeague(userLeague *models.UserLeague) error {
	// TODO:
}

// GetLeaderboard retrieves the leaderboard for a given league ID.
func (r *LeagueRepository) GetLeaderboard(leagueID uint) (*models.Leaderboard, error) {
	// TODO:
}

// GetLeagueDetails retrieves details for a specific league by ID.
func (r *LeagueRepository) GetLeagueDetails(leagueID uint) (*models.League, error) {
	var league models.League
	err := r.db.Preload("Users").Where("id = ?", leagueID).First(&league).Error
	return &league, err
}
