package leagueportfolio

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type LeaguePortfolioHandler struct {
	leaguePortfolioService LeaguePortfolioService
}

func NewLeaguePortfolioHandler(leaguePortfolioService LeaguePortfolioService) *LeaguePortfolioHandler {
	return &LeaguePortfolioHandler{
		leaguePortfolioService: leaguePortfolioService,
	}
}

func (h *LeaguePortfolioHandler) DraftStock(c *gin.Context) {
	var request struct {
		LeaguePortfolioID uint `json:"league_portfolio_id" binding:"required"`
		UserPortfolioID   uint `json:"user_portfolio_id" binding:"required"`
		StockID           uint `json:"stock_id" binding:"required"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	err := h.leaguePortfolioService.DraftStock(request.LeaguePortfolioID, request.UserPortfolioID, request.StockID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Stock drafted successfully"})
}
