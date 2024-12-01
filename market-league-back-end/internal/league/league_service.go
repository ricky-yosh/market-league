package league

import (
	"fmt"
	"time"

	"github.com/market-league/internal/models"
	"github.com/market-league/internal/portfolio"
	"github.com/market-league/internal/user"
)

// LeagueService handles the business logic for managing leagues.
type LeagueService struct {
	repo          *LeagueRepository
	userRepo      *user.UserRepository           // Reference to UserRepository
	portfolioRepo *portfolio.PortfolioRepository // Reference to PortfolioRepository
}

// NewLeagueService creates a new instance of LeagueService.
func NewLeagueService(repo *LeagueRepository, userRepo *user.UserRepository, portfolioRepo *portfolio.PortfolioRepository) *LeagueService {
	return &LeagueService{
		repo:          repo,
		userRepo:      userRepo,
		portfolioRepo: portfolioRepo,
	}
}

// LeagueResponse represents the response with sanitized users.
type LeagueResponse struct {
	ID         uint                   `json:"id"`
	LeagueName string                 `json:"league_name"`
	StartDate  time.Time              `json:"start_date"`
	EndDate    time.Time              `json:"end_date"`
	Users      []models.SanitizedUser `json:"users"`
}

// CreateLeague creates a new league with the given details.
func (s *LeagueService) CreateLeague(leagueName string, ownerUser uint, startDate, endDate string) (*LeagueResponse, error) {
	// Parse start and end dates into time.Time
	start, err := time.Parse(time.RFC3339, startDate)
	if err != nil {
		return nil, fmt.Errorf("invalid start date format: %v", err)
	}

	end, err := time.Parse(time.RFC3339, endDate)
	if err != nil {
		return nil, fmt.Errorf("invalid end date format: %v", err)
	}

	// Fetch the user by username
	owner, err := s.userRepo.GetUserByID(ownerUser)
	if err != nil {
		return nil, fmt.Errorf("failed to find owner user: %v", err)
	}

	// Create a new league instance and add the owner to the Users slice
	league := &models.League{
		LeagueName: leagueName,
		StartDate:  start,
		EndDate:    end,
		Users:      []models.User{*owner}, // Add the owner to the Users slice
	}

	// Save the league to the repository
	err = s.repo.CreateLeague(league)
	if err != nil {
		return nil, fmt.Errorf("failed to create league: %v", err)
	}

	// Use the existing SanitizeUsers function to sanitize the user data
	sanitizedUsers := SanitizeUsers(league.Users)

	// Return the league response with sanitized users
	return &LeagueResponse{
		ID:         league.ID,
		LeagueName: league.LeagueName,
		StartDate:  league.StartDate,
		EndDate:    league.EndDate,
		Users:      sanitizedUsers,
	}, nil
}

// AddUserToLeague associates a user with a specific league.
func (s *LeagueService) AddUserToLeague(userID, leagueID uint) error {
	// Delegate the logic to the repository
	err := s.repo.AddUserToLeague(userID, leagueID)
	if err != nil {
		return fmt.Errorf("%v", err)
	}

	return nil
}

// GetLeagueDetails retrieves details for a specific league by ID.
func (s *LeagueService) GetLeagueDetails(leagueID uint) (*models.League, []models.SanitizedUser, error) {
	// Fetch the league details from the repository
	league, err := s.repo.GetLeagueDetails(leagueID)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to fetch league details: %v", err)
	}

	// Convert the league's users to sanitized DTOs
	sanitized_users := SanitizeUsers(league.Users)

	return league, sanitized_users, nil
}

// Helper Functions
func SanitizeUser(user models.User) models.SanitizedUser {
	return models.SanitizedUser{
		ID:        user.ID,
		Username:  user.Username,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
	}
}

// Sanitize User Object to make object more secure
func SanitizeUsers(users []models.User) []models.SanitizedUser {
	sanitized_users := make([]models.SanitizedUser, len(users))
	for i, user := range users {
		sanitized_users[i] = SanitizeUser(user)
	}
	return sanitized_users
}

// GetLeaderboard retrieves the leaderboard for a specific league.
func (s *LeagueService) GetLeaderboard(leagueID uint, portfolioService *portfolio.PortfolioService) ([]models.LeaderboardEntry, error) {
	// Delegate the leaderboard retrieval to the repository and pass the portfolio service for calculations
	return s.repo.GetLeaderboard(leagueID, portfolioService)
}

// RemoveLeague removes a league and all associated data in a transaction
func (s *LeagueService) RemoveLeague(leagueID uint) error {
	// Start a transaction
	tx := s.repo.db.Begin()
	if err := tx.Error; err != nil {
		return err
	}

	// Execute deletions in the correct order
	if err := s.repo.RemovePortfolioStocksByLeagueID(tx, leagueID); err != nil {
		tx.Rollback()
		return err
	}

	if err := s.repo.RemovePortfoliosByLeagueID(tx, leagueID); err != nil {
		tx.Rollback()
		return err
	}

	if err := s.repo.RemoveTradesByLeagueID(tx, leagueID); err != nil {
		tx.Rollback()
		return err
	}

	if err := s.repo.RemoveLeaguePortfolioByLeagueID(tx, leagueID); err != nil {
		tx.Rollback()
		return err
	}

	if err := s.repo.RemoveUserLeaguesByLeagueID(tx, leagueID); err != nil {
		tx.Rollback()
		return err
	}

	if err := s.repo.RemoveLeague(tx, leagueID); err != nil {
		tx.Rollback()
		return err
	}

	// Commit the transaction
	return tx.Commit().Error
}
