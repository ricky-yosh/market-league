package leagueportfolio

import (
	"encoding/json"
	"fmt"
	"log"

	ws "github.com/market-league/api/websocket"
)

// LeaguePortfolioHandler Interface
type LeaguePortfolioHandlerInterface interface {
	DraftStock(conn *ws.Connection, rawData json.RawMessage) error
	GetLeaguePortfolioInfo(conn *ws.Connection, rawData json.RawMessage) error
}

// Compile-time check
var _ LeaguePortfolioHandlerInterface = (*LeaguePortfolioHandler)(nil)

type LeaguePortfolioHandler struct {
	leaguePortfolioService *LeaguePortfolioService
}

func NewLeaguePortfolioHandler(
	leaguePortfolioService *LeaguePortfolioService,
) *LeaguePortfolioHandler {
	return &LeaguePortfolioHandler{
		leaguePortfolioService: leaguePortfolioService,
	}
}

// * Implementation of Interface

// DraftStock handles the drafting of a stock by transferring it from the LeaguePortfolio to a user's portfolio
func (h *LeaguePortfolioHandler) DraftStock(conn *ws.Connection, rawData json.RawMessage) error {
	var request struct {
		LeagueID uint `json:"league_id" binding:"required"`
		UserID   uint `json:"user_id" binding:"required"`
		StockID  uint `json:"stock_id" binding:"required"`
	}
	if err := json.Unmarshal(rawData, &request); err != nil {
		ws.SendError(conn, ws.MessageType_LeaguePortfolio_DraftStock, "Invalid input: "+err.Error())
		return fmt.Errorf("invalid input: %v", err)
	}

	// Add logging to see the request
	log.Printf("DraftStock request: LeagueID=%d, UserID=%d, StockID=%d",
		request.LeagueID, request.UserID, request.StockID)

	// Check if there's an active draft channel for this league.
	draftChan := h.leaguePortfolioService.GetDraftSelectionChannel(request.LeagueID)

	// Add logging to see if the channel is nil
	if draftChan == nil {
		log.Printf("No active draft channel found for league %d", request.LeagueID)
		return fmt.Errorf("no active draft for league %d", request.LeagueID)
	}

	log.Printf("Found active draft channel for league %d, sending stock ID %d",
		request.LeagueID, request.StockID)

	// Send the selection into the channel.
	draftChan <- request.StockID
	log.Printf("Successfully sent stock ID %d to draft channel for league %d",
		request.StockID, request.LeagueID)

	return nil
}

// GetLeaguePortfolioInfo handles retrieving the League's LeaguePortfolio
func (h *LeaguePortfolioHandler) GetLeaguePortfolioInfo(conn *ws.Connection, rawData json.RawMessage) error {
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
	if err := conn.Ws.WriteJSON(response); err != nil {
		return fmt.Errorf("failed to send response: %v", err)
	}

	return nil
}
