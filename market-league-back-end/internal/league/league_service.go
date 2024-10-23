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

// CreateLeague creates a new league with the given details.
func (s *LeagueService) CreateLeague(leagueName, ownerUser, startDate, endDate string) (*models.League, error) {
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
	owner, err := s.userRepo.GetUserByUsername(ownerUser)
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

	return league, nil
}

// AddUserToLeague associates a user with a specific league.
func (s *LeagueService) AddUserToLeague(userID, leagueID uint) error {
	// Delegate the logic to the repository
	err := s.repo.AddUserToLeague(userID, leagueID)
	if err != nil {
		return fmt.Errorf("failed to add user to league: %v", err)
	}

	return nil
}

// GetLeagueDetails retrieves details for a specific league by ID.
func (s *LeagueService) GetLeagueDetails(leagueID uint) (*models.League, error) {
	league, err := s.repo.GetLeagueDetails(leagueID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch league details: %v", err)
	}
	return league, nil
}

// GetLeaderboard retrieves the leaderboard for a specific league.
func (s *LeagueService) GetLeaderboard(leagueID uint, portfolioService *portfolio.PortfolioService) ([]models.LeaderboardEntry, error) {
	// Delegate the leaderboard retrieval to the repository and pass the portfolio service for calculations
	return s.repo.GetLeaderboard(leagueID, portfolioService)
}
