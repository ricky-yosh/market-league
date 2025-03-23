package league

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	ws "github.com/market-league/api/websocket"
	leagueportfolio "github.com/market-league/internal/league_portfolio"
	"github.com/market-league/internal/models"
	"github.com/market-league/internal/portfolio"
	"github.com/market-league/internal/user"
)

// LeagueService handles the business logic for managing leagues.
type LeagueService struct {
	repo                   *LeagueRepository
	userRepo               *user.UserRepository           // Reference to UserRepository
	portfolioRepo          *portfolio.PortfolioRepository // Reference to PortfolioRepository
	leaguePortfolioService *leagueportfolio.LeaguePortfolioService
	activeDraftChannels    map[uint]chan uint // activeDraftChannels maps leagueID to a channel that receives a drafted stockID.
	mu                     sync.Mutex         // Protect concurrent access to activeDraftChannels.
}

// NewLeagueService creates a new instance of LeagueService.
func NewLeagueService(
	repo *LeagueRepository,
	userRepo *user.UserRepository,
	portfolioRepo *portfolio.PortfolioRepository,
	leaguePortfolioService *leagueportfolio.LeaguePortfolioService,
) *LeagueService {
	return &LeagueService{
		repo:                   repo,
		userRepo:               userRepo,
		portfolioRepo:          portfolioRepo,
		leaguePortfolioService: leaguePortfolioService,
		activeDraftChannels:    make(map[uint]chan uint),
	}
}

// SetLeaguePortfolioService allows setting the dependency after initialization.
func (s *LeagueService) SetLeaguePortfolioService(lpService *leagueportfolio.LeaguePortfolioService) {
	s.leaguePortfolioService = lpService
}

// LeagueResponse represents the response with sanitized users.
type LeagueResponse struct {
	ID         uint                   `json:"id"`
	LeagueName string                 `json:"league_name"`
	StartDate  time.Time              `json:"start_date"`
	EndDate    time.Time              `json:"end_date"`
	Users      []models.SanitizedUser `json:"users"`
}

// CreateLeague creates a new league with the given details.
// Since a league starts with only one user (the owner),
// it adds the owner to the Users slice and creates a LeaguePlayer record for them.
func (s *LeagueService) CreateLeague(leagueName string, ownerUser uint, startDate, endDate string) (*LeagueResponse, error) {
	// Parse start and end dates into time.Time
	start, err := time.Parse(time.RFC3339, startDate)
	if err != nil {
		return nil, fmt.Errorf("invalid start date format: %v", err)
	}
	end, err := time.Parse(time.RFC3339, endDate)
	if err != nil {
		return nil, fmt.Errorf("invalid end date format: %v", err)
	}

	// Fetch the owner user by ID
	owner, err := s.userRepo.GetUserByID(ownerUser)
	if err != nil {
		return nil, fmt.Errorf("failed to find owner user: %v", err)
	}

	// Create a new league instance with the owner in the Users slice.
	league := &models.League{
		LeagueName: leagueName,
		StartDate:  start,
		EndDate:    end,
		Users:      []models.User{*owner},
	}

	// Save the league to the repository.
	if err := s.repo.CreateLeague(league); err != nil {
		return nil, fmt.Errorf("failed to create league: %v", err)
	}

	// Create a LeaguePlayer record for the owner.
	lp := models.LeaguePlayer{
		LeagueID:    league.ID,
		PlayerID:    owner.ID,
		DraftStatus: models.DraftNotReady, // default status
	}
	if err := s.repo.CreateLeaguePlayer(&lp); err != nil {
		return nil, fmt.Errorf("failed to create league player for owner %d: %v", owner.ID, err)
	}
	league.LeaguePlayers = []models.LeaguePlayer{lp}

	// Sanitize user data for the response.
	sanitizedUsers := SanitizeUsers(league.Users)

	// Return the league response with sanitized users.
	return &LeagueResponse{
		ID:         league.ID,
		LeagueName: league.LeagueName,
		StartDate:  league.StartDate,
		EndDate:    league.EndDate,
		Users:      sanitizedUsers,
	}, nil
}

// AddUserToLeague associates a user with a league and creates a LeaguePlayer record.
func (s *LeagueService) AddUserToLeague(userID, leagueID uint) error {
	// First, add the user to the league via the join table.
	if err := s.repo.AddUserToLeague(userID, leagueID); err != nil {
		return fmt.Errorf("failed to add user to league: %v", err)
	}

	// Next, create a LeaguePlayer record for the new user.
	lp := models.LeaguePlayer{
		LeagueID:    leagueID,
		PlayerID:    userID,
		DraftStatus: models.DraftNotReady, // default draft status
	}
	if err := s.repo.CreateLeaguePlayer(&lp); err != nil {
		return fmt.Errorf("failed to create league player record: %v", err)
	}

	return nil
}

// GetLeagueDetails retrieves details for a specific league by ID.
func (s *LeagueService) GetLeagueDetails(leagueID uint) (*models.League, []models.SanitizedUser, error) {
	// Fetch the league details from the repository
	league, err := s.repo.GetLeagueDetails(leagueID)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to fetch league details: %v", err)
	}

	// Convert the league's users to sanitized DTOs
	sanitized_users := SanitizeUsers(league.Users)

	return league, sanitized_users, nil
}

// Helper Functions
func SanitizeUser(user models.User) models.SanitizedUser {
	return models.SanitizedUser{
		ID:        user.ID,
		Username:  user.Username,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
	}
}

// Sanitize User Object to make object more secure
func SanitizeUsers(users []models.User) []models.SanitizedUser {
	sanitized_users := make([]models.SanitizedUser, len(users))
	for i, user := range users {
		sanitized_users[i] = SanitizeUser(user)
	}
	return sanitized_users
}

// GetLeaderboard retrieves the leaderboard for a specific league.
func (s *LeagueService) GetLeaderboard(leagueID uint, portfolioService *portfolio.PortfolioService) ([]models.LeaderboardEntry, error) {
	// Delegate the leaderboard retrieval to the repository and pass the portfolio service for calculations
	return s.repo.GetLeaderboard(leagueID, portfolioService)
}

// RemoveLeague removes a league and all associated data in a transaction
func (s *LeagueService) RemoveLeague(leagueID uint) error {
	// Start a transaction
	tx := s.repo.db.Begin()
	if err := tx.Error; err != nil {
		return err
	}

	// Execute deletions in the correct order
	if err := s.repo.RemoveLeaguePlayerByLeagueID(tx, leagueID); err != nil {
		tx.Rollback()
		return err
	}

	if err := s.repo.RemoveOwnershipHistoriesByLeagueID(tx, leagueID); err != nil {
		tx.Rollback()
		return err
	}

	if err := s.repo.RemovePortfolioStocksByLeagueID(tx, leagueID); err != nil {
		tx.Rollback()
		return err
	}

	if err := s.repo.RemovePortfolioPointsHistoryByLeagueID(tx, leagueID); err != nil {
		tx.Rollback()
		return err
	}

	if err := s.repo.RemovePortfoliosByLeagueID(tx, leagueID); err != nil {
		tx.Rollback()
		return err
	}

	if err := s.repo.RemoveTradesByLeagueID(tx, leagueID); err != nil {
		tx.Rollback()
		return err
	}

	if err := s.repo.RemoveLeaguePortfolioStocksByLeagueID(tx, leagueID); err != nil {
		tx.Rollback()
		return err
	}

	if err := s.repo.RemoveLeaguePortfolioByLeagueID(tx, leagueID); err != nil {
		tx.Rollback()
		return err
	}

	if err := s.repo.RemoveUserLeaguesByLeagueID(tx, leagueID); err != nil {
		tx.Rollback()
		return err
	}

	if err := s.repo.RemoveLeague(tx, leagueID); err != nil {
		tx.Rollback()
		return err
	}

	// Commit the transaction
	return tx.Commit().Error
}

// Add this function to your LeagueService struct
func (s *LeagueService) GetPlayerPortfoliosInLeague(leagueID uint) ([]models.Portfolio, error) {
	// Fetch all portfolios for the league
	portfolios, err := s.portfolioRepo.GetPortfoliosForLeague(leagueID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch portfolios: %w", err)
	}
	return portfolios, nil
}

func (s *LeagueService) GetAllLeagues() ([]models.League, error) {
	// Delegate the league retrieval to the repository
	return s.repo.GetAllLeagues()
}

// QueueUpPlayer marks a player as queued and checks if all players are ready.
// If all are ready, it updates the league state and broadcasts the update.
func (s *LeagueService) QueueUpPlayer(leagueID uint, playerID uint, conn *ws.Connection) error {
	// 1. Update the player's queue status.
	if err := s.repo.QueueUpPlayer(leagueID, playerID); err != nil {
		return err
	}

	// 2. Subscribe this connection to the league.
	conn.Subscriptions[leagueID] = true

	// 3. Check if all players are queued
	allReady, err := s.repo.AllPlayersReady(leagueID)
	if err != nil {
		return err
	}

	if allReady {
		// All players are ready. Retrieve the league.
		league, err := s.repo.GetLeague(leagueID)
		if err != nil {
			return err
		}

		// Update league state to indicate the draft is live.
		league.LeagueState = models.InDraft
		if err := s.repo.UpdateLeague(league); err != nil {
			return err
		}

		// Broadcast the updated league details with its new state
		if err := s.BroadcastLeagueDetails(leagueID); err != nil {
			return err
		}

		// Start the drafting loop in its own goroutine.
		go s.startDraftLoop(leagueID)
	} else {
		// Not all players are ready, broadcast the current league state
		if err := s.BroadcastLeagueDetails(leagueID); err != nil {
			return err
		}
	}

	return nil
}

// BroadcastLeagueDetails broadcasts the league details to all subscribers
func (s *LeagueService) BroadcastLeagueDetails(leagueID uint) error {
	// Get updated league details to broadcast
	league, err := s.repo.GetLeagueDetails(leagueID)
	if err != nil {
		return fmt.Errorf("failed to get league details: %w", err)
	}

	// Prepare the data for broadcast
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

	// Broadcast to all connections subscribed to this league
	responseBytes, err := json.Marshal(response)
	if err != nil {
		return fmt.Errorf("failed to marshal websocket message: %w", err)
	}
	ws.Manager.BroadcastToLeague(leagueID, responseBytes)

	return nil
}

// Duration for each player's draft turn.
const draftTurnDuration = 30 * time.Second

// startDraftLoop handles the turn-based drafting process.
func (s *LeagueService) startDraftLoop(leagueID uint) {
	// Create a channel for receiving the draft selection.
	selectionChannel := make(chan uint)

	// Lock before updating the map.
	s.mu.Lock()
	s.activeDraftChannels[leagueID] = selectionChannel
	s.mu.Unlock()

	// Ensure the channel is removed when the draft loop completes.
	defer func() {
		s.mu.Lock()
		delete(s.activeDraftChannels, leagueID)
		s.mu.Unlock()
	}()

	league, err := s.repo.GetLeague(leagueID)
	if err != nil {
		log.Println("startDraftLoop: error getting league:", err)
		return
	}

	players := s.getOrderedDraftPlayers(league)
	if len(players) == 0 {
		log.Println("startDraftLoop: no players available for drafting")
		return
	}

	currentPlayerIndex := 0

	for !s.isDraftComplete(league) {

		currentPlayer := players[currentPlayerIndex]
		log.Printf("Draft turn for player %d in league %d", currentPlayer, leagueID)

		timer := time.NewTimer(draftTurnDuration)

		// Add logging before waiting for selection
		log.Printf("Waiting for player %d selection in league %d", currentPlayer, leagueID)

		// Wait for player's selection on the channel or timer expiration.
		stockID, selectionReceived := s.waitForPlayerSelection(currentPlayer, leagueID, timer, selectionChannel)

		// Add logging after selection is received or timer expires
		if selectionReceived {
			log.Printf("Player %d selection received: stockID=%d", currentPlayer, stockID)
		} else {
			log.Printf("Timer expired for player %d, no selection received", currentPlayer)
		}

		if selectionReceived {
			err := s.leaguePortfolioService.DraftStock(leagueID, currentPlayer, stockID)
			if err != nil {
				log.Printf("Error processing selection for player %d: %v", currentPlayer, err)
			} else {
				s.broadcastDraftPick(leagueID, currentPlayer, stockID)
			}
		} else {
			autoStockID, err := s.autoSelectStock(leagueID)
			if err != nil {
				log.Printf("Auto-select error for player %d: %v", currentPlayer, err)
			} else {
				err := s.leaguePortfolioService.DraftStock(leagueID, currentPlayer, autoStockID)
				if err != nil {
					log.Printf("Error processing auto-selection for player %d: %v", currentPlayer, err)
				} else {
					s.broadcastDraftPick(leagueID, currentPlayer, autoStockID)
				}
			}
		}
		timer.Stop()
		currentPlayerIndex = (currentPlayerIndex + 1) % len(players)
		if updatedLeague, err := s.repo.GetLeague(leagueID); err == nil {
			league = updatedLeague
		}
	}

	// Update league state to PostDraft
	league.LeagueState = models.PostDraft
	if err := s.repo.UpdateLeague(league); err != nil {
		log.Println("Error updating league to PostDraft:", err)
	}

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
		log.Println("Failed to serialize league data:", err)
	}

	// Construct the WebSocket message
	response := ws.WebsocketMessage{
		Type: ws.MessageType_League_GetDetails,
		Data: json.RawMessage(dataJSON),
	}

	// Broadcast the message to all users in the league
	respBytes, err := json.Marshal(response)
	if err != nil {
		log.Println("Failed to serialize WebSocket message:", err)
	}

	ws.Manager.BroadcastToLeague(leagueID, respBytes)
}

func (s *LeagueService) waitForPlayerSelection(playerID, leagueID uint, timer *time.Timer, selectionChannel chan uint) (uint, bool) {
	// Notify all clients that this player is now on the clock
	s.notifyPlayerOnClock(playerID, leagueID)

	log.Printf("waitForPlayerSelection: Waiting for selection from player %d in league %d", playerID, leagueID)

	// Simply use select directly in the main function, no need for a goroutine
	select {
	case stockID := <-selectionChannel:
		// Player made a selection within the time limit
		log.Printf("waitForPlayerSelection: Player %d made selection (stock ID: %d) in league %d",
			playerID, stockID, leagueID)
		timer.Stop() // Stop the timer
		return stockID, true

	case <-timer.C:
		// Timer expired, player did not make a selection in time
		log.Printf("waitForPlayerSelection: Timer expired for player %d in league %d",
			playerID, leagueID)
		return 0, false
	}
}

func (s *LeagueService) notifyPlayerOnClock(playerID, leagueID uint) {
	// Create notification message
	data, err := json.Marshal(map[string]interface{}{
		"leagueID":      leagueID,
		"playerID":      playerID,
		"remainingTime": int(draftTurnDuration.Seconds()),
	})

	if err != nil {
		log.Printf("Error marshalling player on clock notification: %v", err)
		return
	}

	// Create websocket message
	response := ws.WebsocketMessage{
		Type: ws.MessageType_League_DraftUpdate,
		Data: json.RawMessage(data),
	}

	// Marshal the response
	respBytes, err := json.Marshal(response)
	if err != nil {
		log.Printf("Error marshalling websocket message: %v", err)
		return
	}

	// Broadcast to all members of the league
	ws.Manager.BroadcastToLeague(leagueID, respBytes)
}

// broadcastDraftPick sends a message to all connections subscribed to the league,
// informing them of the current pick (or auto-pick).
func (s *LeagueService) broadcastDraftPick(leagueID, playerID, stockID uint) error {
	// Broadcast draft pick
	if err := s.broadcastDraftPickMessage(leagueID, playerID, stockID); err != nil {
		log.Println("Error broadcasting draft pick:", err)
		return err
	}

	// Broadcast all portfolios in the league
	if err := s.broadcastLeaguePortfolios(leagueID); err != nil {
		log.Println("Error broadcasting league portfolios:", err)
		return err
	}

	// Broadcast the updated portfolio for the player who made the draft pick
	if err := s.broadcastLeaguePortfolio(leagueID); err != nil {
		log.Println("Error broadcasting league portfolio:", err)
		return err
	}

	return nil
}

// Helper function to broadcast the draft pick message
func (s *LeagueService) broadcastDraftPickMessage(leagueID, playerID, stockID uint) error {
	payload := map[string]interface{}{
		"league_id": leagueID,
		"player_id": playerID,
		"stock_id":  stockID,
	}

	data, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("error marshaling draft pick payload: %w", err)
	}

	response := ws.WebsocketMessage{
		Type: ws.MessageType_League_DraftPick,
		Data: json.RawMessage(data),
	}

	responseBytes, err := json.Marshal(response)
	if err != nil {
		return fmt.Errorf("error marshaling draft pick response: %w", err)
	}

	ws.Manager.BroadcastToLeague(leagueID, responseBytes)
	return nil
}

// Helper function to broadcast all portfolios in the league
func (s *LeagueService) broadcastLeaguePortfolios(leagueID uint) error {
	// Get all the portfolios in the league for the draft
	portfolios, err := s.portfolioRepo.GetPortfoliosForLeague(leagueID)
	if err != nil {
		return fmt.Errorf("error getting portfolios for league: %w", err)
	}

	data, err := json.Marshal(portfolios)
	if err != nil {
		return fmt.Errorf("error marshaling portfolios: %w", err)
	}

	response := ws.WebsocketMessage{
		Type: ws.MessageType_League_Portfolios,
		Data: json.RawMessage(data),
	}

	responseBytes, err := json.Marshal(response)
	if err != nil {
		return fmt.Errorf("error marshaling portfolios response: %w", err)
	}

	ws.Manager.BroadcastToLeague(leagueID, responseBytes)
	return nil
}

// Helper function to broadcast a specific player's portfolio
func (s *LeagueService) broadcastLeaguePortfolio(leagueID uint) error {
	// Get the detailed portfolio
	portfolio, err := s.leaguePortfolioService.GetLeaguePortfolioInfo(leagueID)
	if err != nil {
		return fmt.Errorf("error getting portfolio details: %w", err)
	}

	data, err := json.Marshal(portfolio)
	if err != nil {
		return fmt.Errorf("error marshaling portfolio: %w", err)
	}

	response := ws.WebsocketMessage{
		Type: ws.MessageType_LeaguePortfolio_GetLeaguePortfolioInfo,
		Data: json.RawMessage(data),
	}

	responseBytes, err := json.Marshal(response)
	if err != nil {
		return fmt.Errorf("error marshaling portfolio response: %w", err)
	}

	ws.Manager.BroadcastToLeague(leagueID, responseBytes)
	return nil
}

func (s *LeagueService) GetDraftSelectionChannel(leagueID uint) chan uint {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.activeDraftChannels[leagueID]
}

// isDraftComplete checks whether the draft is complete by verifying
// if each player has drafted 5 stocks.
func (s *LeagueService) isDraftComplete(league *models.League) bool {
	// Get all player portfolios for this league
	playerPortfolios, err := s.portfolioRepo.GetPortfoliosForLeague(league.ID)

	// Custom logging with field details
	// portfoliosJSON, _ := json.MarshalIndent(playerPortfolios, "", "  ")
	// log.Printf("Player portfolios:\n%s", string(portfoliosJSON))
	// for i, portfolio := range playerPortfolios {
	// 	log.Printf("Portfolio[%d]: ID=%d, UserID=%d, LeagueID=%d, Points=%d, StockCount=%d",
	// 		i, portfolio.ID, portfolio.UserID, portfolio.LeagueID, portfolio.Points, len(portfolio.Stocks))

	// 	// Optional: Print details about each stock
	// 	for j, stock := range portfolio.Stocks {
	// 		log.Printf("  Stock[%d]: ID=%d, Symbol=%s", j, stock.ID, stock.TickerSymbol)
	// 	}
	// }
	// End of custom logging

	if err != nil {
		log.Printf("Error checking if draft is complete: %v", err)
		return false
	}

	// Get the number of players in the league
	numPlayers := len(league.Users)
	if numPlayers == 0 {
		return false
	}

	// Check if each player has drafted 5 stocks
	for _, portfolio := range playerPortfolios {
		if len(portfolio.Stocks) < 5 {
			// At least one player has not drafted 5 stocks yet
			return false
		}
	}

	// All players have drafted 5 stocks each, so draft is complete
	return true
}

// autoSelectStock returns an auto-selected stock for the given player.
// Implement your own logic here.
func (s *LeagueService) autoSelectStock(leagueID uint) (uint, error) {
	// Get the league portfolio for the given league ID
	leaguePortfolio, err := s.leaguePortfolioService.GetLeaguePortfolioInfo(leagueID)
	if err != nil {
		return 0, fmt.Errorf("failed to get league portfolio: %w", err)
	}

	// Check if there are any stocks in the portfolio
	if len(leaguePortfolio.Stocks) == 0 {
		return 0, fmt.Errorf("no stocks available in the league portfolio")
	}

	// Use the newer random number generation approach
	randomIndex := rand.Intn(len(leaguePortfolio.Stocks))

	// Return the ID of the randomly selected stock
	return leaguePortfolio.Stocks[randomIndex].ID, nil
}

// getOrderedDraftPlayers returns a slice of player IDs for the league,
// ordered in the sequence you want for the draft.
// For now, it simply uses the order of LeaguePlayers as stored in the league.
func (s *LeagueService) getOrderedDraftPlayers(league *models.League) []uint {
	var players []uint
	for _, lp := range league.Users {
		players = append(players, lp.ID)
	}
	return players
}

func (s *LeagueService) HandlePlayerDisconnect(leagueID uint, conn *ws.Connection) {
	// Remove this connection from all league subscriptions
	for subLeagueID := range conn.Subscriptions {
		if subLeagueID == leagueID {
			delete(conn.Subscriptions, subLeagueID)
		}
	}

	// Note: We don't need to stop the draft loop here
	// The draft loop should continue and auto-select if needed
}
