package api

import (
    "log"
    "time"

	// "github.com/market-league/internal/models"
    "github.com/market-league/internal/services"
	"github.com/market-league/internal/stock"
	"gorm.io/gorm"
)

type scheduler struct {
	db *gorm.DB
	StockService *stock.StockService
    stockRepo stock.StockRepository
}

func (s *scheduler) StartDailyTask() {
    go func() {
        for {

            location, err := time.LoadLocation("America/New_York")
            if err != nil {
                log.Printf("Error loading time location: %v", err)
                return
            }

            now := time.Now().In(location)
            nextRun := time.Date(now.Year(), now.Month(), now.Day(), 14, 56, 0, 0, now.Location()) // Set to 9:31 AM in New York

            if now.After(nextRun) {
                nextRun = nextRun.Add(24 * time.Hour)
            }

            // Wait until the next scheduled time
            time.Sleep(time.Until(nextRun))
            
            // need to get from DB
            companies, err := s.stockRepo.GetAllStocks()
            // companies := []string {"AAPL",  "MSFT",  "GOOGL",  "AMZN",  "META",  "TSLA",  "NFLX",  "NVDA",  "JPM",  "BAC",  "DIS",  "V",  "MA",  "UNH",  "HD",  "PG",  "KO",  "PEP",  "CSCO",  "CMCSA",  "ORCL",  "INTC",  "IBM",  "TXN",  "UPS"}
            
            // Max 30 API/sc
            for _, company := range companies {
                quote, err := services.GetTestStock(company.TickerSymbol)
                if err != nil {
                    log.Printf("Error fetching stock data: %v", err)
                } else {
                    log.Printf("Fetched stock data: %+v", quote)
                    // NEED TO OPEN STOCK INFO
                    // Here function not made yet
                    // h.StockService.UpdateStockPrice(stockID uint, newPrice float64, timestamp *time.Time) in stock package
                    // err := s.StockService.UpdateStockPrice(company.ID, quote.c, &now)
                    // if err != nil {
                    //     log.Printf("Failed to update stock price for company %s with quote %f: %v", company.TickerSymbol, quote, err)
                    // }
                }
                time.Sleep(40 * time.Millisecond)
            }
            // Calculate time for the next execution
            now = time.Now().In(location)
            nextRun = time.Date(now.Year(), now.Month(), now.Day(), 19, 30, 0, 0, location).Add(24 * time.Hour)
            time.Sleep(time.Until(nextRun))
        }
    }()
}
// r * scheduler doesnt work like you think it does
// func (r *scheduler) UpdateStockCurrentPrice(symbol string, price float64) {
//     // company, err := FindCompany(symbol)
//     if err != nil:
//     {
//         return err
//     }else {
//         // here update history price
//         return r.db.Model(&models.User{}).Where("TickerSymbol = ?", symbol).Update("CurrentPrice", price).Error
//     }
// }
// func (r *scheduler) FindCompany(symbol string) {
//     var foundCompany models.Stock
//     err := r.db.Where("TickerSymbol = ?", symbol).First(&foundCompany).Error
// 	if err != nil {
// 		return nil, err
// 	}
// 	return &foundCompany, nil
// }

