package league

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/gin-gonic/gin"
	ws "github.com/market-league/api/websocket"
	league_portfolio "github.com/market-league/internal/league_portfolio"
	"github.com/market-league/internal/models"
	"github.com/market-league/internal/portfolio"
)

// LeagueHandler Interface
type LeagueHandlerInterface interface {
	CreateLeague(conn *ws.Connection, rawData json.RawMessage) error
	AddUserToLeague(conn *ws.Connection, rawData json.RawMessage) error
	GetLeagueDetails(conn *ws.Connection, rawData json.RawMessage) error
	GetLeaderboard(conn *ws.Connection, rawData json.RawMessage) error
	RemoveLeague(conn *ws.Connection, rawData json.RawMessage) error
	QueueUp(conn *ws.Connection, rawData json.RawMessage) error
	GetPlayerPortfoliosInLeague(conn *ws.Connection, rawData json.RawMessage) error
	GetAllLeagues(conn *ws.Connection, rawData json.RawMessage) error
	SubscribeToLeague(conn *ws.Connection, rawData json.RawMessage) error
	UnsubscribeToLeague(conn *ws.Connection, rawData json.RawMessage) error
	HandleDisconnect(leagueID uint, conn *ws.Connection) error
}

// Compile-time check
var _ LeagueHandlerInterface = (*LeagueHandler)(nil)

// LeagueHandler defines the HTTP handler for league-related operations.
type LeagueHandler struct {
	service                *LeagueService
	portfolioService       *portfolio.PortfolioService
	leaguePortfolioService *league_portfolio.LeaguePortfolioService
}

// NewLeagueHandler creates a new instance of LeagueHandler.
func NewLeagueHandler(
	service *LeagueService,
	portfolioService *portfolio.PortfolioService,
	leaguePortfolioService *league_portfolio.LeaguePortfolioService,
) *LeagueHandler {
	return &LeagueHandler{
		service:                service,
		portfolioService:       portfolioService,
		leaguePortfolioService: leaguePortfolioService,
	}
}

// * Implementation of Interface

// CreateLeague handles the creation of a new league.
func (h *LeagueHandler) CreateLeague(conn *ws.Connection, rawData json.RawMessage) error {
	// Step 1: Parse the WebSocket message
	var request struct {
		LeagueName string `json:"league_name" binding:"required"`
		OwnerUser  uint   `json:"owner_user" binding:"required"`
		EndDate    string `json:"end_date" binding:"required"`
	}

	// Step 2: Parse data from WebSocket JSON payload
	if err := json.Unmarshal(rawData, &request); err != nil {
		ws.SendError(conn, ws.MessageType_League_CreateLeague, "Invalid input: "+err.Error())
		return fmt.Errorf("invalid input: %v", err)
	}

	// Step 3: Process business logic (reuse the service layer)

	// Step 3a: Pass the values to the service to create the league
	startDate := time.Now().Format(time.RFC3339) // Set the start date to the current date and time
	league, err := h.service.CreateLeague(request.LeagueName, request.OwnerUser, startDate, request.EndDate)
	if err != nil {
		ws.SendError(conn, ws.MessageType_League_CreateLeague, err.Error())
		return fmt.Errorf("failed to retrieve portfolio with ID: %v", err)
	}

	// Step 3b: Create a portfolio for the user in the league
	portfolio, err := h.portfolioService.CreatePortfolio(request.OwnerUser, league.ID)
	if err != nil {
		ws.SendError(conn, ws.MessageType_League_CreateLeague, err.Error())
		return fmt.Errorf("failed to retrieve portfolio with ID: %v", err)
	}

	// Step 3c: Create a league portfolio using the new LeaguePortfolioService
	leaguePortfolio, err := h.leaguePortfolioService.CreateLeaguePortfolio(league.ID)
	if err != nil {
		ws.SendError(conn, ws.MessageType_League_CreateLeague, err.Error())
		return fmt.Errorf("failed to retrieve portfolio with ID: %v", err)
	}

	// Step 5: Send success response back via WebSocket

	// Step 5a: Construct response with sanitized user details
	data := gin.H{
		"league":          league,
		"userPortfolio":   portfolio,
		"leaguePortfolio": leaguePortfolio,
	}
	dataJSON, err := json.Marshal(data) // Marshal the payload into JSON bytes
	if err != nil {
		ws.SendError(conn, ws.MessageType_League_CreateLeague, "Failed to serialize response")
		return fmt.Errorf("serialization error: %v", err)
	}

	// Step 5b: Construct response with sanitized user details
	response := ws.WebsocketMessage{
		Type: ws.MessageType_League_CreateLeague,
		Data: json.RawMessage(dataJSON), // Use marshaled JSON bytes
	}
	if err := conn.Ws.WriteJSON(response); err != nil {
		return fmt.Errorf("failed to send response: %v", err)
	}

	return nil

}

// AddUserToLeague handles adding a user to a league.
func (h *LeagueHandler) AddUserToLeague(conn *ws.Connection, rawData json.RawMessage) error {
	// Step 1: Parse the WebSocket message
	var request struct {
		UserID   uint `json:"user_id" binding:"required"`
		LeagueID uint `json:"league_id" binding:"required"`
	}

	// Step 2: Parse data from WebSocket JSON payload
	if err := json.Unmarshal(rawData, &request); err != nil {
		ws.SendError(conn, ws.MessageType_League_AddUserToLeague, "Invalid input: "+err.Error())
		return fmt.Errorf("invalid input: %v", err)
	}

	// Step 3: Process business logic (reuse the service layer)

	// Step 3a: Process business logic (reuse the service layer)
	if err := h.service.AddUserToLeague(request.UserID, request.LeagueID); err != nil {
		ws.SendError(conn, ws.MessageType_League_AddUserToLeague, err.Error())
		return fmt.Errorf("failed to retrieve portfolio with ID: %v", err)
	}

	// Step 3b: Process business logic (reuse the service layer)
	portfolio, err := h.portfolioService.CreatePortfolio(request.UserID, request.LeagueID)
	if err != nil {
		ws.SendError(conn, ws.MessageType_League_AddUserToLeague, err.Error())
		return fmt.Errorf("failed to retrieve portfolio with ID: %v", err)
	}

	// Step 4: Marshal the portfolio into JSON
	portfolioJSON, err := json.Marshal(portfolio)
	if err != nil {
		ws.SendError(conn, ws.MessageType_League_AddUserToLeague, "Failed to serialize portfolio")
		return fmt.Errorf("serialization error: %v", err)
	}

	// Step 5: Send success response back via WebSocket
	response := ws.WebsocketMessage{
		Type: ws.MessageType_League_AddUserToLeague,
		Data: json.RawMessage(portfolioJSON), // Use marshaled JSON bytes
	}
	if err := conn.Ws.WriteJSON(response); err != nil {
		return fmt.Errorf("failed to send response: %v", err)
	}

	return nil

}

// GetLeagueDetails handles fetching the details of a specific league.
func (h *LeagueHandler) GetLeagueDetails(conn *ws.Connection, rawData json.RawMessage) error {
	// Step 1: Parse the WebSocket message
	var request struct {
		LeagueID uint `json:"league_id" binding:"required"`
	}

	// Step 2: Parse data from WebSocket JSON payload
	if err := json.Unmarshal(rawData, &request); err != nil {
		ws.SendError(conn, ws.MessageType_League_GetDetails, "Invalid input: "+err.Error())
		return fmt.Errorf("invalid input: %v", err)
	}

	// Step 3: Process business logic (reuse the service layer)
	league, users, err := h.service.GetLeagueDetails(request.LeagueID)
	if err != nil {
		ws.SendError(conn, ws.MessageType_League_GetDetails, err.Error())
		return fmt.Errorf("failed to retrieve portfolio with ID: %v", err)
	}

	// Step 4: Marshal the portfolio into JSON
	data := gin.H{
		"id":             league.ID,
		"league_name":    league.LeagueName,
		"start_date":     league.StartDate,
		"end_date":       league.EndDate,
		"league_state":   league.LeagueState,
		"users":          users,
		"max_players":    league.MaxPlayers,
		"league_players": league.LeaguePlayers,
	}
	// Construct response with sanitized user details
	dataJSON, err := json.Marshal(data)
	if err != nil {
		ws.SendError(conn, ws.MessageType_League_GetDetails, "Failed to serialize portfolio")
		return fmt.Errorf("serialization error: %v", err)
	}

	// Step 5: Send success response back via WebSocket
	response := ws.WebsocketMessage{
		Type: ws.MessageType_League_GetDetails,
		Data: json.RawMessage(dataJSON), // Use marshaled JSON bytes
	}
	if err := conn.Ws.WriteJSON(response); err != nil {
		return fmt.Errorf("failed to send response: %v", err)
	}

	return nil
}

// GetLeaderboard handles fetching the leaderboard for a specific league.
func (h *LeagueHandler) GetLeaderboard(conn *ws.Connection, rawData json.RawMessage) error {
	// Step 1: Parse the WebSocket message
	var request struct {
		LeagueID uint `json:"league_id" binding:"required"`
	}

	// Step 2: Parse data from WebSocket JSON payload
	if err := json.Unmarshal(rawData, &request); err != nil {
		ws.SendError(conn, ws.MessageType_League_GetLeaderboard, "Invalid input: "+err.Error())
		return fmt.Errorf("invalid input: %v", err)
	}

	// Step 3: Process business logic (reuse the service layer)
	leaderboard, err := h.service.GetLeaderboard(request.LeagueID, h.portfolioService)
	if err != nil {
		ws.SendError(conn, ws.MessageType_League_GetLeaderboard, err.Error())
		return fmt.Errorf("failed to retrieve portfolio with ID: %v", err)
	}

	// Step 4: Marshal the portfolio into JSON
	leaderboardJSON, err := json.Marshal(leaderboard)
	if err != nil {
		ws.SendError(conn, ws.MessageType_League_GetLeaderboard, "Failed to serialize portfolio")
		return fmt.Errorf("serialization error: %v", err)
	}

	// Step 5: Send success response back via WebSocket
	response := ws.WebsocketMessage{
		Type: ws.MessageType_League_GetLeaderboard,
		Data: json.RawMessage(leaderboardJSON), // Use marshaled JSON bytes
	}
	if err := conn.Ws.WriteJSON(response); err != nil {
		return fmt.Errorf("failed to send response: %v", err)
	}

	return nil

}

// RemoveLeague handles the removal of a league and all associated records
func (h *LeagueHandler) RemoveLeague(conn *ws.Connection, rawData json.RawMessage) error {
	// Step 1: Parse the WebSocket message
	var request struct {
		LeagueID uint `json:"league_id" binding:"required"`
	}

	// Step 2: Parse data from WebSocket JSON payload
	if err := json.Unmarshal(rawData, &request); err != nil {
		ws.SendError(conn, ws.MessageType_League_RemoveLeague, "Invalid input: "+err.Error())
		return fmt.Errorf("invalid input: %v", err)
	}

	// Step 3: Process business logic (reuse the service layer)
	if err := h.service.RemoveLeague(request.LeagueID); err != nil {
		ws.SendError(conn, ws.MessageType_League_RemoveLeague, err.Error())
		return fmt.Errorf("failed to retrieve portfolio with ID: %v", err)
	}

	// Step 4: Send success response (no data, just confirmation)
	response := ws.WebsocketMessage{
		Type: ws.MessageType_League_RemoveLeague,
		Data: json.RawMessage(`{"message": "League removed successfully"}`), // Simple JSON message
	}
	if err := conn.Ws.WriteJSON(response); err != nil {
		return fmt.Errorf("failed to send response: %v", err)
	}

	return nil
}

func (h *LeagueHandler) GetPlayerPortfoliosInLeague(conn *ws.Connection, rawData json.RawMessage) error {
	var request struct {
		LeagueID uint `json:"league_id" binding:"required"`
	}

	if err := json.Unmarshal(rawData, &request); err != nil {
		ws.SendError(conn, ws.MessageType_League_Portfolios, "Invalid input: "+err.Error())
		return fmt.Errorf("invalid input: %v", err)
	}

	// Fetch all portfolios for the league from the service layer
	portfolios, err := h.service.GetPlayerPortfoliosInLeague(request.LeagueID)
	if err != nil {
		ws.SendError(conn, ws.MessageType_League_Portfolios, err.Error())
		return fmt.Errorf("failed to get player portfolios: %v", err)
	}

	// Marshal the portfolios directly as they should be a slice
	dataJSON, err := json.Marshal(portfolios)
	if err != nil {
		ws.SendError(conn, ws.MessageType_League_Portfolios, "Failed to serialize response")
		return fmt.Errorf("serialization error: %v", err)
	}

	response := ws.WebsocketMessage{
		Type: ws.MessageType_League_Portfolios,
		Data: json.RawMessage(dataJSON),
	}

	// Send the response to the requesting connection
	if err := conn.Ws.WriteJSON(response); err != nil {
		return fmt.Errorf("failed to send response: %v", err)
	}

	// Broadcast to all connections in the league
	responseBytes, err := json.Marshal(response)
	if err != nil {
		return fmt.Errorf("failed to marshal websocket message: %v", err)
	}
	ws.Manager.BroadcastToLeague(request.LeagueID, responseBytes)

	return nil
}

func (h *LeagueHandler) GetAllLeagues(conn *ws.Connection, rawData json.RawMessage) error {
	// Step 2: Process business logic (use the service layer)
	leagues, err := h.service.GetAllLeagues()
	if err != nil {
		ws.SendError(conn, ws.MessageType_League_GetAllLeagues, err.Error())
		return fmt.Errorf("failed to retrieve leagues: %v", err)
	}

	// Step 3: Marshal the leagues into JSON
	leaguesJSON, err := json.Marshal(leagues)
	if err != nil {
		ws.SendError(conn, ws.MessageType_League_GetAllLeagues, "Failed to serialize leagues")
		return fmt.Errorf("serialization error: %v", err)
	}

	// Step 4: Send success response back via WebSocket
	response := ws.WebsocketMessage{
		Type: ws.MessageType_League_GetAllLeagues,
		Data: json.RawMessage(leaguesJSON),
	}
	if err := conn.Ws.WriteJSON(response); err != nil {
		return fmt.Errorf("failed to send response: %v", err)
	}

	return nil
}

// QueueUp handles a player's queue-up action via WebSocket.
func (h *LeagueHandler) QueueUp(conn *ws.Connection, rawData json.RawMessage) error {
	var request struct {
		LeagueID uint `json:"league_id" binding:"required"`
		PlayerID uint `json:"player_id" binding:"required"`
	}

	if err := json.Unmarshal(rawData, &request); err != nil {
		ws.SendError(conn, ws.MessageType_League_QueueUp, "Invalid input: "+err.Error())
		return fmt.Errorf("invalid input: %v", err)
	}

	if err := h.service.QueueUpPlayer(request.LeagueID, request.PlayerID, conn); err != nil {
		ws.SendError(conn, ws.MessageType_League_QueueUp, err.Error())
		return fmt.Errorf("failed to queue up player: %v", err)
	}

	responseData := gin.H{"message": "Player queued up successfully"}
	dataJSON, err := json.Marshal(responseData)
	if err != nil {
		ws.SendError(conn, ws.MessageType_League_QueueUp, "Failed to serialize response")
		return fmt.Errorf("serialization error: %v", err)
	}

	response := ws.WebsocketMessage{
		Type: ws.MessageType_League_QueueUp,
		Data: json.RawMessage(dataJSON),
	}

	if err := conn.Ws.WriteJSON(response); err != nil {
		return fmt.Errorf("failed to send response: %v", err)
	}

	return nil
}

// SendCurrentDraftState sends the current draft state to a newly connected client
func (s *LeagueService) SendCurrentDraftState(leagueID uint, conn *ws.Connection) error {
	// Get the current player on clock
	// This would need to be tracked in your LeagueService
	// For now, we'll just broadcast the league details

	// Broadcast league details to this specific connection
	league, err := s.repo.GetLeagueDetails(leagueID)
	if err != nil {
		return fmt.Errorf("failed to get league details: %w", err)
	}

	// Prepare the data
	data := gin.H{
		"id":             league.ID,
		"league_name":    league.LeagueName,
		"start_date":     league.StartDate,
		"end_date":       league.EndDate,
		"league_state":   league.LeagueState,
		"max_players":    league.MaxPlayers,
		"league_players": league.LeaguePlayers,
	}

	// Marshal the data into JSON
	dataJSON, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("serialization error: %w", err)
	}

	// Create the WebSocket message
	response := ws.WebsocketMessage{
		Type: ws.MessageType_League_GetDetails,
		Data: json.RawMessage(dataJSON),
	}

	// Send to this specific connection
	if err := conn.Ws.WriteJSON(response); err != nil {
		return fmt.Errorf("failed to send response: %v", err)
	}

	return nil
}

func (h *LeagueHandler) SubscribeToLeague(conn *ws.Connection, rawData json.RawMessage) error {
	var request struct {
		LeagueID uint `json:"league_id" binding:"required"`
	}

	if err := json.Unmarshal(rawData, &request); err != nil {
		ws.SendError(conn, ws.MessageType_League_SubscribeToLeague, "Invalid input: "+err.Error())
		return fmt.Errorf("invalid input: %v", err)
	}

	// Add this connection to the league's subscriptions
	conn.Subscriptions[request.LeagueID] = true

	h.service.BroadcastLeagueDetails(request.LeagueID)

	// Get current league state to inform the client
	league, _, err := h.service.GetLeagueDetails(request.LeagueID)
	if err != nil {
		ws.SendError(conn, ws.MessageType_League_SubscribeToLeague, "Failed to get league details: "+err.Error())
		return fmt.Errorf("failed to get league details: %v", err)
	}

	// If the league is in draft mode, send the current draft state
	if league.LeagueState == models.InDraft {
		// Send the current player on clock information
		if err := h.service.SendCurrentDraftState(request.LeagueID, conn); err != nil {
			log.Printf("Error sending draft state: %v", err)
		}
	}

	// Send success response
	responseData := gin.H{"message": "Subscribed to league successfully"}
	dataJSON, err := json.Marshal(responseData)
	if err != nil {
		ws.SendError(conn, ws.MessageType_League_SubscribeToLeague, "Failed to serialize response")
		return fmt.Errorf("serialization error: %v", err)
	}

	response := ws.WebsocketMessage{
		Type: ws.MessageType_League_SubscribeToLeague,
		Data: json.RawMessage(dataJSON),
	}

	if err := conn.Ws.WriteJSON(response); err != nil {
		return fmt.Errorf("failed to send response: %v", err)
	}

	return nil
}

func (h *LeagueHandler) UnsubscribeToLeague(conn *ws.Connection, rawData json.RawMessage) error {
	var request struct {
		LeagueID uint `json:"league_id" binding:"required"`
	}

	if err := json.Unmarshal(rawData, &request); err != nil {
		ws.SendError(conn, ws.MessageType_League_UnsubscribeToLeague, "Invalid input: "+err.Error())
		return fmt.Errorf("invalid input: %v", err)
	}

	// Remove this connection from the league's subscriptions
	delete(conn.Subscriptions, request.LeagueID)

	// Send success response
	responseData := gin.H{"message": "Unsubscribed from league successfully"}
	dataJSON, err := json.Marshal(responseData)
	if err != nil {
		ws.SendError(conn, ws.MessageType_League_UnsubscribeToLeague, "Failed to serialize response")
		return fmt.Errorf("serialization error: %v", err)
	}

	response := ws.WebsocketMessage{
		Type: ws.MessageType_League_UnsubscribeToLeague,
		Data: json.RawMessage(dataJSON),
	}

	if err := conn.Ws.WriteJSON(response); err != nil {
		return fmt.Errorf("failed to send response: %v", err)
	}

	return nil
}

func (h *LeagueHandler) HandleDisconnect(leagueID uint, conn *ws.Connection) error {
	// Log the disconnection
	log.Printf("Player disconnected from league %d", leagueID)

	// Remove this connection's subscription
	delete(conn.Subscriptions, leagueID)

	return nil
}
