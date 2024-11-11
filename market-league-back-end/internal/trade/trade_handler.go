package trade

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// TradeHandler handles HTTP requests for trades
type TradeHandler struct {
	TradeService *TradeService
}

// NewTradeHandler creates a new instance of TradeHandler
func NewTradeHandler(tradeService *TradeService) *TradeHandler {
	return &TradeHandler{
		TradeService: tradeService,
	}
}

// CreateTradeHandler handles the creation of a new trade
func (h *TradeHandler) CreateTrade(c *gin.Context) {
	var request struct {
		LeagueID   uint   `json:"league_id"`
		User1ID    uint   `json:"user1_id"`
		User2ID    uint   `json:"user2_id"`
		Stocks1IDs []uint `json:"stocks1_ids"`
		Stocks2IDs []uint `json:"stocks2_ids"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	trade, err := h.TradeService.CreateTrade(request.LeagueID, request.User1ID, request.User2ID, request.Stocks1IDs, request.Stocks2IDs)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, trade)
}
