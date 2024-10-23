package portfolio

import (
	"github.com/market-league/internal/models"
	"gorm.io/gorm"
)

// PortfolioRepository provides access to portfolio-related operations in the database.
type PortfolioRepository struct {
	db *gorm.DB
}

// NewPortfolioRepository creates a new instance of PortfolioRepository.
func NewPortfolioRepository(db *gorm.DB) *PortfolioRepository {
	return &PortfolioRepository{db: db}
}

// GetPortfolioByID fetches a portfolio by its ID.
func (r *PortfolioRepository) GetPortfolioByID(portfolioID uint) (*models.Portfolio, error) {
	var portfolio models.Portfolio
	err := r.db.Preload("Stocks").First(&portfolio, portfolioID).Error
	if err != nil {
		return nil, err
	}
	return &portfolio, nil
}

// GetUserPortfolioInLeague fetches a user's portfolio in a specific league.
func (r *PortfolioRepository) GetUserPortfolioInLeague(userID, leagueID uint) (*models.Portfolio, error) {
	var portfolio models.Portfolio
	err := r.db.Preload("Stocks").Where("user_id = ? AND league_id = ?", userID, leagueID).First(&portfolio).Error
	if err != nil {
		return nil, err
	}
	return &portfolio, nil
}

// CreatePortfolio creates a new portfolio for a user in a league.
func (r *PortfolioRepository) CreatePortfolio(portfolio *models.Portfolio) error {
	return r.db.Create(portfolio).Error
}

// UpdatePortfolio updates an existing portfolio in the database.
func (r *PortfolioRepository) UpdatePortfolio(portfolio *models.Portfolio) error {
	return r.db.Save(portfolio).Error
}

// DeletePortfolio deletes a portfolio by its ID.
func (r *PortfolioRepository) DeletePortfolio(portfolioID uint) error {
	return r.db.Delete(&models.Portfolio{}, portfolioID).Error
}
