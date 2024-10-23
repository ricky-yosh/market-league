package league

import (
	"fmt"
	"sort"
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
func (s *LeagueService) CreateLeague(leagueName, startDate, endDate string) (*models.League, error) {
	// Parse start and end dates into time.Time
	start, err := time.Parse("2006-01-02", startDate)
	if err != nil {
		return nil, fmt.Errorf("invalid start date format: %v", err)
	}

	end, err := time.Parse("2006-01-02", endDate)
	if err != nil {
		return nil, fmt.Errorf("invalid end date format: %v", err)
	}

	// Create a new league instance
	league := &models.League{
		LeagueName: leagueName,
		StartDate:  start,
		EndDate:    end,
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
	// Step 1: Retrieve the league by ID
	league, err := s.repo.GetLeagueDetails(leagueID)
	if err != nil {
		return fmt.Errorf("league not found: %v", err)
	}

	// Step 2: Retrieve the user by ID using the UserRepository
	user, err := s.userRepo.GetUserByID(userID)
	if err != nil {
		return fmt.Errorf("user not found: %v", err)
	}

	// Step 3: Add the user to the league using GORM's association handling
	err = s.repo.db.Model(&league).Association("Users").Append(&user)
	if err != nil {
		return fmt.Errorf("failed to add user to league: %v", err)
	}

	return nil
}

// GetLeaderboard retrieves a sorted list of users based on their portfolio value in a specific league.
func (s *LeagueService) GetLeaderboard(leagueID uint) ([]models.UserPortfolioValue, error) {
	// Step 1: Get all users in the league
	league, err := s.repo.GetLeagueDetails(leagueID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch league details: %v", err)
	}

	// Step 2: Get each user's portfolio value in the league dynamically
	var userPortfolioValues []models.UserPortfolioValue
	for _, user := range league.Users {
		// Get the user's portfolio for this league
		portfolio, err := s.portfolioRepo.GetUserPortfolioInLeague(user.ID, leagueID)
		if err != nil {
			return nil, fmt.Errorf("failed to fetch portfolio for user %d: %v", user.ID, err)
		}

		// Calculate the total portfolio value dynamically
		totalValue := portfolio.GetPortfolioValue()

		// Add the user's portfolio value to the list
		userPortfolioValues = append(userPortfolioValues, models.UserPortfolioValue{
			UserID:     user.ID,
			Username:   user.Username,
			TotalValue: totalValue,
		})
	}

	// Step 3: Sort users based on total portfolio value in descending order
	sort.Slice(userPortfolioValues, func(i, j int) bool {
		return userPortfolioValues[i].TotalValue > userPortfolioValues[j].TotalValue
	})

	return userPortfolioValues, nil
}

// GetLeagueDetails retrieves details for a specific league by ID.
func (s *LeagueService) GetLeagueDetails(leagueID uint) (*models.League, error) {
	league, err := s.repo.GetLeagueDetails(leagueID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch league details: %v", err)
	}
	return league, nil
}

// Helper method to calculate the total value of a portfolio
func (s *LeagueService) calculatePortfolioValue(portfolio *models.Portfolio) (float64, error) {
	var totalValue float64
	for _, stock := range portfolio.Stocks {
		totalValue += stock.CurrentPrice // Assuming quantity is 1 per stock
	}
	return totalValue, nil
}
