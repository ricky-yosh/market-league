package user

import (
	"fmt"

	"github.com/market-league/internal/models"
)

// UserService handles business logic related to users.
type UserService struct {
	repo *UserRepository // Use the repository to access the database
}

// NewUserService creates a new instance of UserService.
func NewUserService(repo *UserRepository) *UserService {
	return &UserService{repo: repo}
}

// GetUserInfo represents the aggregated information for a user.
type GetUserInfo struct {
	UserID     uint               `json:"user_id"`
	Username   string             `json:"username"`
	Email      string             `json:"email"`
	Portfolios []models.Portfolio `json:"portfolios,omitempty"`
	Leagues    []models.League    `json:"leagues,omitempty"`
	Trades     []models.Trade     `json:"trades,omitempty"`
}

// GetUserByID fetches user information based on filter criteria.
func (s *UserService) GetUserByID(userID uint) (*GetUserInfo, error) {
	// Fetch the user details (username, email, etc.)
	user, err := s.repo.GetUserByID(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch user details: %v", err)
	}

	// Prepare the GetUserInfo response
	userInfo := &GetUserInfo{
		UserID:   user.ID,
		Username: user.Username,
		Email:    user.Email,
	}

	return userInfo, nil
}

// GetUserLeagues retrieves all leagues for a given user.
func (s *UserService) GetUserLeagues(userID uint) ([]models.League, error) {
	leagues, err := s.repo.GetUserLeagues(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get leagues for user: %v", err)
	}

	return leagues, nil
}

// GetUserTrades retrieves all trades involving a given user within a specific league.
func (s *UserService) GetUserTrades(userID uint, leagueID uint) ([]models.Trade, error) {
	trades, err := s.repo.GetUserTrades(userID, leagueID)
	if err != nil {
		return nil, fmt.Errorf("failed to get trades for user in league %d: %v", leagueID, err)
	}

	return trades, nil
}

// GetUserPortfolios retrieves all portfolios for a given user.
func (s *UserService) GetUserPortfolios(userID uint) ([]models.Portfolio, error) {
	portfolios, err := s.repo.GetUserPortfolios(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get portfolios for user: %v", err)
	}

	return portfolios, nil
}
