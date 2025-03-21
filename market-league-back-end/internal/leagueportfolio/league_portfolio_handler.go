package leagueportfolio

import (
	"encoding/json"
	"fmt"

	"github.com/gorilla/websocket"
	ws "github.com/market-league/api/websocket"
)

// LeaguePortfolioHandler Interface
type LeaguePortfolioHandlerInterface interface {
	DraftStock(conn *websocket.Conn, rawData json.RawMessage) error
	GetLeaguePortfolioInfo(conn *websocket.Conn, rawData json.RawMessage) error
}

// Compile-time check
var _ LeaguePortfolioHandlerInterface = (*LeaguePortfolioHandler)(nil)

type LeaguePortfolioHandler struct {
	leaguePortfolioService *LeaguePortfolioService
}

func NewLeaguePortfolioHandler(leaguePortfolioService *LeaguePortfolioService) *LeaguePortfolioHandler {
	return &LeaguePortfolioHandler{
		leaguePortfolioService: leaguePortfolioService,
	}
}

// * Implementation of Interface

// DraftStock handles the drafting of a stock by transferring it from the LeaguePortfolio to a user's portfolio
func (h *LeaguePortfolioHandler) DraftStock(conn *websocket.Conn, rawData json.RawMessage) error {
	// Step 1: Parse the WebSocket message
	var request struct {
		LeagueID uint `json:"league_id" binding:"required"`
		UserID   uint `json:"user_id" binding:"required"`
		StockID  uint `json:"stock_id" binding:"required"`
	}

	// Step 2: Parse data from WebSocket JSON payload
	if err := json.Unmarshal(rawData, &request); err != nil {
		ws.SendError(conn, ws.MessageType_LeaguePortfolio_DraftStock, "Invalid input: "+err.Error())
		return fmt.Errorf("invalid input: %v", err)
	}

	// Step 3: Process business logic (reuse the service layer)
	if err := h.leaguePortfolioService.DraftStock(request.LeagueID, request.UserID, request.StockID); err != nil {
		ws.SendError(conn, ws.MessageType_LeaguePortfolio_DraftStock, err.Error())
		return fmt.Errorf("failed to retrieve portfolio with ID: %v", err)
	}

	// Step 4: Send success response (no data, just confirmation)
	response := ws.WebsocketMessage{
		Type: ws.MessageType_LeaguePortfolio_DraftStock,
		Data: json.RawMessage(`{"message": "Stock price updated successfully"}`), // Simple JSON message
	}
	if err := conn.WriteJSON(response); err != nil {
		return fmt.Errorf("failed to send response: %v", err)
	}

	return nil
}

// GetLeaguePortfolioInfo handles retrieving the League's LeaguePortfolio
func (h *LeaguePortfolioHandler) GetLeaguePortfolioInfo(conn *websocket.Conn, rawData json.RawMessage) error {
	// Step 1: Parse the WebSocket message
	var request struct {
		LeaguePortfolioID uint `json:"league_id" binding:"required"`
	}

	// Step 2: Parse data from WebSocket JSON payload
	if err := json.Unmarshal(rawData, &request); err != nil {
		ws.SendError(conn, ws.MessageType_LeaguePortfolio_GetLeaguePortfolioInfo, "Invalid input: "+err.Error())
		return fmt.Errorf("invalid input: %v", err)
	}

	// Step 3: Process business logic (reuse the service layer)
	leaguePortfolio, err := h.leaguePortfolioService.GetLeaguePortfolioInfo(request.LeaguePortfolioID)
	if err != nil {
		ws.SendError(conn, ws.MessageType_LeaguePortfolio_GetLeaguePortfolioInfo, err.Error())
		return fmt.Errorf("failed to retrieve portfolio with ID: %v", err)
	}

	// Step 4: Marshal the portfolio into JSON
	leaguePortfolioJSON, err := json.Marshal(leaguePortfolio)
	if err != nil {
		ws.SendError(conn, ws.MessageType_LeaguePortfolio_GetLeaguePortfolioInfo, "Failed to serialize portfolio")
		return fmt.Errorf("serialization error: %v", err)
	}

	// Step 5: Send success response back via WebSocket
	response := ws.WebsocketMessage{
		Type: ws.MessageType_LeaguePortfolio_GetLeaguePortfolioInfo,
		Data: json.RawMessage(leaguePortfolioJSON), // Use marshaled JSON bytes
	}
	if err := conn.WriteJSON(response); err != nil {
		return fmt.Errorf("failed to send response: %v", err)
	}

	return nil
}
