package portfolio

import (
	"encoding/json"
	"fmt"

	"github.com/gorilla/websocket"
	ws "github.com/market-league/api/websocket"
)

// PortfolioHandler Interface
type PortfolioHandlerInterface interface {
	GetPortfolioWithID(conn *websocket.Conn, rawData json.RawMessage) error
	GetLeaguePortfolio(conn *websocket.Conn, rawData json.RawMessage) error
	CreatePortfolio(conn *websocket.Conn, rawData json.RawMessage) error
	AddStockToPortfolio(conn *websocket.Conn, rawData json.RawMessage) error
	RemoveStockFromPortfolio(conn *websocket.Conn, rawData json.RawMessage) error
}

// Compile-time check
var _ PortfolioHandlerInterface = (*PortfolioHandler)(nil)

// PortfolioHandler defines the HTTP handler for portfolio-related operations.
type PortfolioHandler struct {
	service *PortfolioService
}

// NewPortfolioHandler creates a new instance of PortfolioHandler.
func NewPortfolioHandler(service *PortfolioService) *PortfolioHandler {
	return &PortfolioHandler{service: service}
}

// * Implementation of Interface

// GetPortfolio handles fetching a portfolio by its ID.
func (h *PortfolioHandler) GetPortfolioWithID(conn *websocket.Conn, rawData json.RawMessage) error {
	// Step 1: Parse the WebSocket message
	var request struct {
		PortfolioID uint `json:"portfolio_id" binding:"required"`
	}

	// Step 2: Parse data from WebSocket JSON payload
	if err := json.Unmarshal(rawData, &request); err != nil {
		ws.SendError(conn, "Invalid input: "+err.Error())
		return fmt.Errorf("invalid input: %v", err)
	}

	// Step 3: Process business logic (reuse the service layer)
	portfolio, err := h.service.GetPortfolioWithID(request.PortfolioID)
	if err != nil {
		ws.SendError(conn, err.Error())
		return fmt.Errorf("failed to retrieve portfolio with ID: %v", err)
	}

	// Step 4: Marshal the portfolio into JSON
	portfolioJSON, err := json.Marshal(portfolio)
	if err != nil {
		ws.SendError(conn, "Failed to serialize portfolio")
		return fmt.Errorf("serialization error: %v", err)
	}

	// Step 5: Send success response back via WebSocket
	response := ws.WebsocketMessage{
		Type: ws.MessageType_LeaguePortfolio_DraftStock,
		Data: json.RawMessage(portfolioJSON), // Use marshaled JSON bytes
	}
	if err := conn.WriteJSON(response); err != nil {
		return fmt.Errorf("failed to send response: %v", err)
	}

	return nil
}

// GetUserPortfolio handles fetching a user's portfolio in a specific league.
func (h *PortfolioHandler) GetLeaguePortfolio(conn *websocket.Conn, rawData json.RawMessage) error {
	// Step 1: Parse the WebSocket message
	var request struct {
		UserID   uint `json:"user_id" binding:"required"`
		LeagueID uint `json:"league_id" binding:"required"`
	}

	// Step 2: Parse data from WebSocket JSON payload
	if err := json.Unmarshal(rawData, &request); err != nil {
		ws.SendError(conn, "Invalid input: "+err.Error())
		return fmt.Errorf("invalid input: %v", err)
	}

	// Step 3: Process business logic (reuse the service layer)
	portfolio, err := h.service.GetLeaguePortfolio(request.UserID, request.LeagueID)
	if err != nil {
		ws.SendError(conn, err.Error())
		return fmt.Errorf("failed to retrieve User's Portfolio from a specific league: %v", err)
	}

	// Step 4: Marshal the portfolio into JSON
	portfolioJSON, err := json.Marshal(portfolio)
	if err != nil {
		ws.SendError(conn, "Failed to serialize portfolio")
		return fmt.Errorf("serialization error: %v", err)
	}

	// Step 5: Send success response back via WebSocket
	response := ws.WebsocketMessage{
		Type: ws.MessageType_LeaguePortfolio_DraftStock,
		Data: json.RawMessage(portfolioJSON), // Use marshaled JSON bytes
	}
	if err := conn.WriteJSON(response); err != nil {
		return fmt.Errorf("failed to send response: %v", err)
	}

	return nil
}

// CreatePortfolio handles the creation of a new portfolio for a user in a league.
func (h *PortfolioHandler) CreatePortfolio(conn *websocket.Conn, rawData json.RawMessage) error {
	// Step 1: Parse the WebSocket message
	var request struct {
		UserID   uint `json:"user_id"`
		LeagueID uint `json:"league_id"`
	}

	// Step 2: Parse data from WebSocket JSON payload
	if err := json.Unmarshal(rawData, &request); err != nil {
		ws.SendError(conn, "Invalid input: "+err.Error())
		return fmt.Errorf("invalid input: %v", err)
	}

	// Step 3: Process business logic (reuse the service layer)
	portfolio, err := h.service.CreatePortfolio(request.UserID, request.LeagueID)
	if err != nil {
		ws.SendError(conn, err.Error())
		return fmt.Errorf("failed to create portfolio: %v", err)
	}

	// Step 4: Marshal the portfolio into JSON
	portfolioJSON, err := json.Marshal(portfolio)
	if err != nil {
		ws.SendError(conn, "Failed to serialize portfolio")
		return fmt.Errorf("serialization error: %v", err)
	}

	// Step 5: Send success response back via WebSocket
	response := ws.WebsocketMessage{
		Type: ws.MessageType_LeaguePortfolio_DraftStock,
		Data: json.RawMessage(portfolioJSON), // Use marshaled JSON bytes
	}
	if err := conn.WriteJSON(response); err != nil {
		return fmt.Errorf("failed to send response: %v", err)
	}

	return nil
}

// AddStockToPortfolio handles adding a stock to a user's portfolio.
func (h *PortfolioHandler) AddStockToPortfolio(conn *websocket.Conn, rawData json.RawMessage) error {
	// Step 1: Parse the WebSocket message
	var request struct {
		PortfolioID uint `json:"portfolio_id" binding:"required"`
		StockID     uint `json:"stock_id" binding:"required"`
	}

	// Step 2: Parse data from WebSocket JSON payload
	if err := json.Unmarshal(rawData, &request); err != nil {
		ws.SendError(conn, "Invalid input: "+err.Error())
		return fmt.Errorf("invalid input: %v", err)
	}

	// Step 3: Process business logic (reuse the service layer)
	if err := h.service.AddStockToPortfolio(request.PortfolioID, request.StockID); err != nil {
		ws.SendError(conn, err.Error())
		return fmt.Errorf("failed to add stock to portfolio: %v", err)
	}

	// Step 4: Send success response (no data, just confirmation)
	response := ws.WebsocketMessage{
		Type: ws.MessageType_LeaguePortfolio_DraftStock,
		Data: json.RawMessage(`{"message": "Stock added successfully"}`), // Simple JSON message
	}
	if err := conn.WriteJSON(response); err != nil {
		return fmt.Errorf("failed to send response: %v", err)
	}

	return nil
}

// RemoveStockFromPortfolio handles removing a stock from a user's portfolio.
func (h *PortfolioHandler) RemoveStockFromPortfolio(conn *websocket.Conn, rawData json.RawMessage) error {
	// Step 1: Parse the WebSocket message
	var request struct {
		PortfolioID uint `json:"portfolio_id" binding:"required"`
		StockID     uint `json:"stock_id" binding:"required"`
	}

	// Step 2: Parse data from WebSocket JSON payload
	if err := json.Unmarshal(rawData, &request); err != nil {
		ws.SendError(conn, "Invalid input: "+err.Error())
		return fmt.Errorf("invalid input: %v", err)
	}

	// Step 3: Process business logic (reuse the service layer)
	if err := h.service.RemoveStockFromPortfolio(request.PortfolioID, request.StockID); err != nil {
		ws.SendError(conn, err.Error())
		return fmt.Errorf("failed to remove a stock from a portfolio: %v", err)
	}

	// Step 4: Send success response (no data, just confirmation)
	response := ws.WebsocketMessage{
		Type: ws.MessageType_LeaguePortfolio_DraftStock,
		Data: json.RawMessage(`{"message": "Stock removed successfully"}`), // Simple JSON message
	}
	if err := conn.WriteJSON(response); err != nil {
		return fmt.Errorf("failed to send response: %v", err)
	}

	return nil
}
