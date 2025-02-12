package utils

import (
	"fmt"

	"github.com/market-league/internal/models"
)

// * Helper Functions *
// Safely retrieve first item in stock list
func FirstStock(slice []models.Stock) (*models.Stock, error) {
	firstIndex := 0
	if len(slice) > 0 {
		return &slice[firstIndex], nil
	}
	return nil, fmt.Errorf("stocks list is empty")
}
