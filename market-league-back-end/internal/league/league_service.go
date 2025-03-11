package league

import (
	"encoding/json"
	"fmt"
	"log"
	"sync"
	"time"

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

	// Fetch the user by username
	owner, err := s.userRepo.GetUserByID(ownerUser)
	if err != nil {
		return nil, fmt.Errorf("failed to find owner user: %v", err)
	}

	// Create a new league instance and add the owner to the Users slice
	league := &models.League{
		LeagueName: leagueName,
		StartDate:  start,
		EndDate:    end,
		Users:      []models.User{*owner}, // Add the owner to the Users slice
	}

	// Save the league to the repository
	err = s.repo.CreateLeague(league)
	if err != nil {
		return nil, fmt.Errorf("failed to create league: %v", err)
	}

	// Use the existing SanitizeUsers function to sanitize the user data
	sanitizedUsers := SanitizeUsers(league.Users)

	// Return the league response with sanitized users
	return &LeagueResponse{
		ID:         league.ID,
		LeagueName: league.LeagueName,
		StartDate:  league.StartDate,
		EndDate:    league.EndDate,
		Users:      sanitizedUsers,
	}, nil
}

// AddUserToLeague associates a user with a specific league.
func (s *LeagueService) AddUserToLeague(userID, leagueID uint) error {
	// Delegate the logic to the repository
	err := s.repo.AddUserToLeague(userID, leagueID)
	if err != nil {
		return fmt.Errorf("%v", err)
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
	if err := s.repo.RemovePortfolioStocksByLeagueID(tx, leagueID); err != nil {
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

// QueueUpPlayer marks a player as queued and checks if all players are ready.
// If all are ready, it updates the league state and broadcasts the update.
func (s *LeagueService) QueueUpPlayer(leagueID uint, playerID uint, conn *ws.Connection) error {
	// 1. Update the player's queue status.
	if err := s.repo.QueueUpPlayer(leagueID, playerID); err != nil {
		return err
	}

	// 2. Subscribe this connection to the league.
	conn.Subscriptions[leagueID] = true

	// 3. Check if all players are queued.
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

		// Get all the portfolios in the league for the draft
		portfolios, err := s.portfolioRepo.GetPortfoliosForLeague(leagueID)
		if err != nil {
			return err
		}

		// Broadcast a "DraftStarted" message.
		data, err := json.Marshal(map[string]interface{}{
			"league_portfolios": portfolios,
		})
		if err != nil {
			return err
		}
		response := ws.WebsocketMessage{
			Type: ws.MessageType_DraftStarted, // Ensure this constant is defined.
			Data: json.RawMessage(data),
		}
		responseBytes, err := json.Marshal(response)
		if err != nil {
			return err
		}
		ws.Manager.BroadcastToLeague(leagueID, responseBytes)

		// Start the drafting loop in its own goroutine.
		go s.startDraftLoop(leagueID)
	}

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
		log.Printf("Draft turn for player %d in league %d\n", currentPlayer, leagueID)

		timer := time.NewTimer(draftTurnDuration)

		// Wait for player's selection on the channel or timer expiration.
		stockID, selectionReceived := s.waitForPlayerSelection(currentPlayer, leagueID, timer, selectionChannel)

		if selectionReceived {
			err := s.leaguePortfolioService.DraftStock(leagueID, currentPlayer, stockID)
			if err != nil {
				log.Printf("Error processing selection for player %d: %v", currentPlayer, err)
			} else {
				s.broadcastDraftPick(leagueID, currentPlayer, stockID, false)
			}
		} else {
			autoStockID, err := s.autoSelectStock(currentPlayer)
			if err != nil {
				log.Printf("Auto-select error for player %d: %v", currentPlayer, err)
			} else {
				err := s.leaguePortfolioService.DraftStock(leagueID, currentPlayer, autoStockID)
				if err != nil {
					log.Printf("Error processing auto-selection for player %d: %v", currentPlayer, err)
				} else {
					s.broadcastDraftPick(leagueID, currentPlayer, autoStockID, true)
				}
			}
		}
		timer.Stop()
		currentPlayerIndex = (currentPlayerIndex + 1) % len(players)
		if updatedLeague, err := s.repo.GetLeague(leagueID); err == nil {
			league = updatedLeague
		}
	}

	league.LeagueState = models.PostDraft
	if err := s.repo.UpdateLeague(league); err != nil {
		log.Println("Error updating league to PostDraft:", err)
	}
	data, err := json.Marshal(map[string]interface{}{
		"league":  league,
		"message": "Draft complete",
	})
	if err == nil {
		response := ws.WebsocketMessage{
			Type: ws.MessageType_DraftComplete,
			Data: json.RawMessage(data),
		}
		respBytes, err := json.Marshal(response)
		if err == nil {
			ws.Manager.BroadcastToLeague(leagueID, respBytes)
		}
	}
}

// waitForPlayerSelection waits for the current player's selection.
// This is a stub: in your real implementation, hook this up to your WebSocket messaging.
func (s *LeagueService) waitForPlayerSelection(playerID, leagueID uint, timer *time.Timer, selectionChan chan uint) (uint, bool) {
	select {
	case stockID := <-selectionChan:
		return stockID, true
	case <-timer.C:
		return 0, false
	}
}

// broadcastDraftPick sends a message to all connections subscribed to the league,
// informing them of the current pick (or auto-pick).
func (s *LeagueService) broadcastDraftPick(leagueID, playerID, stockID uint, auto bool) {
	payload := map[string]interface{}{
		"player_id": playerID,
		"stock_id":  stockID,
		"auto":      auto,
	}
	data, err := json.Marshal(payload)
	if err != nil {
		log.Println("broadcastDraftPick: error marshaling payload:", err)
		return
	}

	response := ws.WebsocketMessage{
		Type: ws.MessageType_DraftPick, // Use a different type if you want to differentiate auto-picks.
		Data: json.RawMessage(data),
	}

	responseBytes, err := json.Marshal(response)
	if err != nil {
		log.Println("broadcastDraftPick: error marshaling response:", err)
		return
	}
	ws.Manager.BroadcastToLeague(leagueID, responseBytes)
}

func (s *LeagueService) GetDraftSelectionChannel(leagueID uint) chan uint {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.activeDraftChannels[leagueID]
}

// isDraftComplete checks whether the draft is complete.
// You need to implement the logic (for example, after a fixed number of rounds).
func (s *LeagueService) isDraftComplete(league *models.League) bool {
	// Implement your draft completion logic.
	return false
}

// autoSelectStock returns an auto-selected stock for the given player.
// Implement your own logic here.
func (s *LeagueService) autoSelectStock(playerID uint) (uint, error) {
	// For example, simply return a fixed stock ID.
	return 1, nil
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
