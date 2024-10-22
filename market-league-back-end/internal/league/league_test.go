package league

import (
	"errors"
	"testing"
)

// Mock repository to simulate database operations
type MockLeagueRepository struct{}

func (m *MockLeagueRepository) Save(league *League) (*League, error) {
	if league.Name == "" {
		return nil, errors.New("name cannot be empty")
	}
	league.ID = 1
	return league, nil
}

func TestCreateLeague(t *testing.T) {
	mockRepo := &MockLeagueRepository{}
	service := NewLeagueService(mockRepo)

	// Test case 1: Successfully creating a league
	league, err := service.CreateLeague("Fantasy League")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if league == nil || league.ID == 0 {
		t.Errorf("Expected valid league, got nil or unassigned ID")
	}

	// Test case 2: Attempting to create a league with an empty name
	_, err = service.CreateLeague("")
	if err == nil {
		t.Errorf("Expected an error, got none")
	}
}
