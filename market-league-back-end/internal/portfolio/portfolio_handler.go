package portfolio

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

// PortfolioHandler defines the HTTP handler for portfolio-related operations.
type PortfolioHandler struct {
	service *PortfolioService
}

// NewPortfolioHandler creates a new instance of PortfolioHandler.
func NewPortfolioHandler(service *PortfolioService) *PortfolioHandler {
	return &PortfolioHandler{service: service}
}

// GetPortfolio handles fetching a portfolio by its ID.
func (h *PortfolioHandler) GetPortfolioWithID(c *gin.Context) {
	var request struct {
		PortfolioID uint `json:"portfolio_id" binding:"required"`
	}

	// Bind the JSON input to the request struct
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	// Fetch the portfolio with the given ID from the service
	portfolio, err := h.service.GetPortfolioWithID(request.PortfolioID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	// Return the portfolio details
	c.JSON(http.StatusOK, portfolio)
}

// GetUserPortfolio handles fetching a user's portfolio in a specific league.
func (h *PortfolioHandler) GetLeaguePortfolio(c *gin.Context) {
	var request struct {
		UserID   uint `json:"user_id" binding:"required"`
		LeagueID uint `json:"league_id" binding:"required"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	portfolio, err := h.service.GetLeaguePortfolio(request.UserID, request.LeagueID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, portfolio)
}

// CreatePortfolio handles the creation of a new portfolio for a user in a league.
func (h *PortfolioHandler) CreatePortfolio(c *gin.Context) {
	var request struct {
		UserID   uint `json:"user_id" binding:"required"`
		LeagueID uint `json:"league_id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	portfolio, err := h.service.CreatePortfolio(request.UserID, request.LeagueID)
	if err != nil {
		c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, portfolio)
}

// AddStockToPortfolio handles adding a stock to a user's portfolio.
func (h *PortfolioHandler) AddStockToPortfolio(c *gin.Context) {
	var request struct {
		PortfolioID uint `json:"portfolio_id" binding:"required"`
		StockID     uint `json:"stock_id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.service.AddStockToPortfolio(request.PortfolioID, request.StockID)
	if err != nil {
		// Check for specific error about stock already being in portfolio
		if err.Error() == fmt.Sprintf("stock with ID %d is already in the portfolio", request.StockID) {
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Stock added to portfolio successfully"})
}

// RemoveStockFromPortfolio handles removing a stock from a user's portfolio.
func (h *PortfolioHandler) RemoveStockFromPortfolio(c *gin.Context) {
	var request struct {
		PortfolioID uint `json:"portfolio_id" binding:"required"`
		StockID     uint `json:"stock_id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.service.RemoveStockFromPortfolio(request.PortfolioID, request.StockID)
	if err != nil {
		// Check for specific error about stock not found
		if err.Error() == fmt.Sprintf("stock with ID %d is not in the portfolio", request.StockID) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Stock removed from portfolio successfully"})
}
