package leagueportfolio

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type LeaguePortfolioHandler struct {
	leaguePortfolioService *LeaguePortfolioService
}

func NewLeaguePortfolioHandler(leaguePortfolioService *LeaguePortfolioService) *LeaguePortfolioHandler {
	return &LeaguePortfolioHandler{
		leaguePortfolioService: leaguePortfolioService,
	}
}

func (h *LeaguePortfolioHandler) DraftStock(c *gin.Context) {
	var request struct {
		LeagueID uint `json:"league_id" binding:"required"`
		UserID   uint `json:"user_id" binding:"required"`
		StockID  uint `json:"stock_id" binding:"required"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	err := h.leaguePortfolioService.DraftStock(request.LeagueID, request.UserID, request.StockID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Stock drafted successfully"})
}

func (h *LeaguePortfolioHandler) GetLeaguePortfolioInfo(c *gin.Context) {
	var request struct {
		LeaguePortfolioID uint `json:"league_id" binding:"required"`
	}

	// Bind the incoming JSON request to the struct
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	// Call the service layer to get the stocks
	leaguePortfolio, err := h.leaguePortfolioService.GetLeaguePortfolioInfo(request.LeaguePortfolioID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Respond with the stocks
	c.JSON(http.StatusOK, leaguePortfolio)
}
