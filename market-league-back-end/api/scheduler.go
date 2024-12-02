package api

import (
	"log"
	"time"

	// "github.com/market-league/internal/models"
	"github.com/market-league/internal/services"
	"github.com/market-league/internal/stock"
	"gorm.io/gorm"
)

type Scheduler struct {
	db           *gorm.DB
	StockService *stock.StockService
	stockRepo    *stock.StockRepository
}

func (s *Scheduler) StartDailyTask() {
	go func() {
		for {
			location, err := time.LoadLocation("America/Chicago")
			if err != nil {
				log.Printf("Error loading time location: %v", err)
				return
			}

			now := time.Now().In(location)
			nextRun := time.Date(now.Year(), now.Month(), now.Day(), 21, 10, 0, 0, now.Location()) // Set to 9:31 AM in New York

			if now.After(nextRun) {
				nextRun = nextRun.Add(5 * time.Minute)
			}

			// Wait until the next scheduled time
			time.Sleep(time.Until(nextRun))

			// need to get from DB
			companies, err := s.stockRepo.GetAllStocks()
            // log.Println(len(companies))
            firstElement := companies[0]
            log.Println("First element:", firstElement)
    
			if err != nil {
				log.Printf("Error with GetAllStocks call: %v", err)
			}
            
			// Max 30 API/sc
			for _, company := range companies {
                log.Println(company.TickerSymbol)
				quote, err := services.GetTestStock(company.TickerSymbol)
				if err != nil {
					log.Printf("Error fetching stock data: %v", err)
				} else {
					// log.Printf("Fetched stock data: %s: %f", company.TickerSymbol, *quote.C)
                    
					err := s.StockService.UpdateStockPrice(company.ID, float64(*quote.C), &now)
					if err != nil {
					    log.Printf("Failed to update stock price for company %s with quote %f: %v", company.TickerSymbol, *quote.C, err)
					}
				}
				time.Sleep(1 * time.Second)
			}
			// Calculate time for the next execution
			now = time.Now().In(location)
			nextRun = time.Date(now.Year(), now.Month(), now.Day(), 21, 10, 0, 0, location).Add(5 * time.Minute)
			time.Sleep(time.Until(nextRun))
		}
	}()
}