package leagueportfolio

import (
	"fmt"

	"github.com/market-league/internal/models"
	"gorm.io/gorm"
)

type LeaguePortfolioRepository struct {
	db *gorm.DB
}

func NewLeaguePortfolioRepository(db *gorm.DB) *LeaguePortfolioRepository {
	return &LeaguePortfolioRepository{db: db}
}

func (r *LeaguePortfolioRepository) CreateLeaguePortfolio(portfolio *models.LeaguePortfolio) (*models.LeaguePortfolio, error) {
	if err := r.db.Create(portfolio).Error; err != nil {
		return nil, err
	}
	return portfolio, nil
}

// GetLeagueDetails retrieves details for a specific league by ID.
func (r *LeaguePortfolioRepository) GetLeagueDetails(leagueID uint) (*models.League, error) {
	var league models.League
	err := r.db.Where("id = ?", leagueID).First(&league).Error
	return &league, err
}

// AddStocksToLeaguePortfolio associates stocks with a league portfolio.
func (r *LeaguePortfolioRepository) AddStocksToLeaguePortfolio(portfolioID uint, stocks []models.Stock) error {
	// Check if there are stocks to add
	if len(stocks) == 0 {
		return nil
	}

	// Check if the portfolio exists
	var portfolio models.LeaguePortfolio
	if err := r.db.Preload("Stocks").First(&portfolio, portfolioID).Error; err != nil {
		return err
	}

	// Add stocks one by one to avoid bulk association issues
	for _, stock := range stocks {
		if err := r.db.Model(&portfolio).Association("Stocks").Append(&stock); err != nil {
			return err
		}
	}

	return nil
}

// GetLeaguePortfolioWithID retrieves a league portfolio by its ID.
func (r *LeaguePortfolioRepository) GetLeaguePortfolioWithID(leaguePortfolioID uint) (*models.LeaguePortfolio, error) {
	var leaguePortfolio models.LeaguePortfolio

	// Preload the Stocks association for the portfolio
	err := r.db.Preload("Stocks").First(&leaguePortfolio, leaguePortfolioID).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("league portfolio with ID %d not found", leaguePortfolioID)
		}
		return nil, fmt.Errorf("failed to fetch league portfolio: %w", err)
	}

	// Fetch the League details using the LeagueRepository
	league, err := r.GetLeagueDetails(leaguePortfolio.LeagueID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch associated league: %w", err)
	}

	// Attach the fetched League to the LeaguePortfolio
	leaguePortfolio.League = *league

	return &leaguePortfolio, nil
}

func (r *LeaguePortfolioRepository) GetLeaguePortfolioIDByLeagueID(leagueID uint) (uint, error) {
	var leaguePortfolioID uint

	// Assuming the LeaguePortfolio model has a column named `league_id` for the relationship
	err := r.db.Table("league_portfolios").
		Select("id").
		Where("league_id = ?", leagueID).
		Scan(&leaguePortfolioID).Error

	if err != nil {
		return 0, fmt.Errorf("failed to fetch LeaguePortfolioID for LeagueID %d: %w", leagueID, err)
	}

	return leaguePortfolioID, nil
}

// UpdateLeaguePortfolio updates an existing league portfolio in the database.
func (r *LeaguePortfolioRepository) UpdateLeaguePortfolio(portfolio *models.LeaguePortfolio) error {
	// Start a transaction to ensure atomicity
	tx := r.db.Begin()

	// Update the LeaguePortfolio itself (e.g., name, etc.)
	if err := tx.Save(portfolio).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to update league portfolio: %w", err)
	}

	// Explicitly update the Stocks association
	if err := tx.Model(portfolio).Association("Stocks").Replace(portfolio.Stocks); err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to update stocks for league portfolio: %w", err)
	}

	// Commit the transaction
	if err := tx.Commit().Error; err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}
