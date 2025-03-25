package api

import (
	"fmt"
	"log"
	"time"

	// "github.com/market-league/internal/models"
	"github.com/market-league/internal/models"
	ownership_history "github.com/market-league/internal/ownership_history"
	"github.com/market-league/internal/portfolio"
	"github.com/market-league/internal/services"
	"github.com/market-league/internal/stock"
	"gorm.io/gorm"
)

type Scheduler struct {
	db                      *gorm.DB
	StockService            *stock.StockService
	stockRepo               *stock.StockRepository
	ownershipHistoryService ownership_history.OwnershipHistoryServiceInterface
	portfolioService        *portfolio.PortfolioService
}

func (s *Scheduler) StartDailyTask() {
	go func() {
		defer func() {
			if r := recover(); r != nil {
				log.Printf("Scheduler recovered from panic: %v", r)
			}
		}()

		for {
			// Load the timezone
			location, err := time.LoadLocation("America/Chicago")
			if err != nil {
				log.Printf("Error loading time location: %v", err)
				return
			}

			// Get the current time
			now := time.Now().In(location)

			// Check if today is a weekday (Monday=1, Sunday=7)
			if now.Weekday() == time.Saturday || now.Weekday() == time.Sunday {
				log.Printf("Skipping task execution as today is a weekend: %s", now.Weekday())
				// Sleep until the next day
				time.Sleep(24 * time.Hour)
				continue
			}

			// Calculate the next run time (set to the next day at 9:00 AM)
			nextRun := time.Date(now.Year(), now.Month(), now.Day(), 9, 0, 0, 0, location)
			if now.After(nextRun) {
				nextRun = nextRun.Add(24 * time.Hour)
			}

			log.Printf("Current time: %s, Next run at: %s", now.Format("15:04:05"), nextRun.Format("15:04:05"))

			// Wait until the next scheduled time
			time.Sleep(time.Until(nextRun))

			// Update league statuses
			s.updateLeagueStatuses(location)

			// Fetch companies from the database
			companies, err := s.stockRepo.GetAllStocks()
			if err != nil {
				log.Printf("Error fetching stocks from database: %v", err)
				continue
			}

			log.Printf("Total companies to process: %d", len(companies))

			// Process each company
			for _, company := range companies {
				quote, err := services.GetTestStock(company.TickerSymbol)
				if err != nil {
					log.Printf("Error fetching stock data for %s: %v", company.TickerSymbol, err)
					continue
				}

				log.Printf("Fetched stock data for %s: Current Price: %.2f", company.TickerSymbol, *quote.C)

				// Update stock price in the database
				updated_time := time.Now().In(location)
				err = s.StockService.UpdateStockPrice(company.ID, float64(*quote.C), &updated_time)
				if err != nil {
					log.Printf("Failed to update stock price for %s: %v", company.TickerSymbol, err)
				}
			}
			err = s.ownershipHistoryService.UpdateActiveOwnershipHistoryCurrentPrices()
			if err != nil {
				fmt.Printf("unable to update active ownershipHistory prices! %v", err)
			}

			err = s.portfolioService.CalculateAllPortfolioTotalValues()
			if err != nil {
				fmt.Printf("unable to update total portfolio values! %v", err)
			}

			log.Printf("Task completed. Waiting for the next interval.")
		}
	}()
}

func (s *Scheduler) updateLeagueStatuses(location *time.Location) {
	// Get the current date (without time component)
	now := time.Now().In(location)
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, location)

	// Query to find leagues whose end date is today
	var leagues []struct {
		ID          int
		LeagueName  string
		EndDate     time.Time
		LeagueState models.LeagueState
	}

	result := s.db.Table("leagues").
		Select("id, league_name, end_date, league_state").
		Where("DATE(end_date) <= DATE(?)", today).
		Find(&leagues)

	if result.Error != nil {
		log.Printf("Error finding leagues to update: %v", result.Error)
		return
	}

	// Update leagues that end today
	for _, league := range leagues {
		if league.LeagueState != models.Completed {
			log.Printf("Updating league '%s' (ID: %d) status to 'completed' as its end date is today",
				league.LeagueName, league.ID)

			// Update the league status
			updateResult := s.db.Table("leagues").
				Where("id = ?", league.ID).
				Update("league_state", models.Completed)

			if updateResult.Error != nil {
				log.Printf("Error updating league status: %v", updateResult.Error)
			} else {
				log.Printf("Successfully updated league '%s' status to 'completed'", league.LeagueName)
			}
		}
	}
}
