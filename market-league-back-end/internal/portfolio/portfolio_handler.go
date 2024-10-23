package portfolio

import (
	"net/http"
	"strconv"

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
func (h *PortfolioHandler) GetPortfolio(c *gin.Context) {
	portfolioID, err := strconv.ParseUint(c.Param("portfolioID"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid portfolio ID"})
		return
	}

	portfolio, err := h.service.GetPortfolio(uint(portfolioID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Portfolio not found"})
		return
	}

	c.JSON(http.StatusOK, portfolio)
}

// GetUserPortfolio handles fetching a user's portfolio in a specific league.
func (h *PortfolioHandler) GetUserPortfolio(c *gin.Context) {
	userID, err := strconv.ParseUint(c.Query("userID"), 10, 64)
	leagueID, err := strconv.ParseUint(c.Query("leagueID"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID or league ID"})
		return
	}

	portfolio, err := h.service.GetUserPortfolio(uint(userID), uint(leagueID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Portfolio not found"})
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
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
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
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Stock removed from portfolio successfully"})
}
