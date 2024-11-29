package leagueportfolio

import (
	"fmt"

	"github.com/market-league/internal/models"
	"gorm.io/gorm"
)

type LeaguePortfolioRepository interface {
	CreateLeaguePortfolio(portfolio *models.LeaguePortfolio) (*models.LeaguePortfolio, error)
	GetLeagueDetails(leagueID uint) (*models.League, error)
	AddStocksToLeaguePortfolio(portfolioID uint, stocks []models.Stock) error
}

type leaguePortfolioRepository struct {
	db *gorm.DB
}

func NewLeaguePortfolioRepository(db *gorm.DB) LeaguePortfolioRepository {
	return &leaguePortfolioRepository{db}
}

func (r *leaguePortfolioRepository) CreateLeaguePortfolio(portfolio *models.LeaguePortfolio) (*models.LeaguePortfolio, error) {
	if err := r.db.Create(portfolio).Error; err != nil {
		return nil, err
	}
	return portfolio, nil
}

// GetLeagueDetails retrieves details for a specific league by ID.
func (r *leaguePortfolioRepository) GetLeagueDetails(leagueID uint) (*models.League, error) {
	var league models.League
	err := r.db.Preload("Users").Where("id = ?", leagueID).First(&league).Error
	return &league, err
}

// AddStocksToLeaguePortfolio associates stocks with a league portfolio.
func (r *leaguePortfolioRepository) AddStocksToLeaguePortfolio(portfolioID uint, stocks []models.Stock) error {
	// Check if there are stocks to add
	if len(stocks) == 0 {
		return nil
	}

	// Check if the portfolio exists
	var portfolio models.LeaguePortfolio
	if err := r.db.Preload("Stocks").First(&portfolio, portfolioID).Error; err != nil {
		return err
	}

	fmt.Println("Stocks: ", stocks)
	fmt.Println("Portfolio: ", stocks)

	// Add stocks one by one to avoid bulk association issues
	for _, stock := range stocks {
		if err := r.db.Model(&portfolio).Association("Stocks").Append(&stock); err != nil {
			return err
		}
	}

	return nil
}
