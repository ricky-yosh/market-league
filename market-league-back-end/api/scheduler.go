package api

import (
    "log"
    "time"

    "github.com/market-league/internal/services"
)

func StartDailyTask() {
    go func() {
        for {
            location, err := time.LoadLocation("America/New_York")
            if err != nil {
                log.Printf("Error loading time location: %v", err)
                return
            }

            now := time.Now().In(location)
            nextRun := time.Date(now.Year(), now.Month(), now.Day(), 19, 30, 0, 0, now.Location()) // Set to 9:31 AM in New York
            if now.After(nextRun) {
                nextRun = nextRun.Add(24 * time.Hour)
            }

            // Wait until the next scheduled time
            time.Sleep(time.Until(nextRun))
            
            // need to get from DB
            companies := []string {"AAPL",  "MSFT",  "GOOGL",  "AMZN",  "META",  "TSLA",  "NFLX",  "NVDA",  "JPM",  "BAC",  "DIS",  "V",  "MA",  "UNH",  "HD",  "PG",  "KO",  "PEP",  "CSCO",  "CMCSA",  "ORCL",  "INTC",  "IBM",  "TXN",  "UPS"}
            
            // Max 30 API/sc
            for _, company := range companies {
                quote, err := services.GetTestStock(company)
                if err != nil {
                    log.Printf("Error fetching stock data: %v", err)
                } else {
                    log.Printf("Fetched stock data: %+v", quote)
                    // NEED TO OPEN STOCK INFO
                    UpdateStockCurrentPrice(company,quote)
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
func (r *scheduler) changeCurrentPrice(symbol string, price float64) {
    company, err := FindCompany(symbol)
    return r.db.Model(&models.User{}).Where("TickerSymbol = ?", symbol).Update("CurrentPrice", price).Error
}
func (r *scheduler) FindCompany(symbol string) {
    var foundCompany models.Stock
    err := r.db.Where("TickerSymbol = ?", symbol).First(&foundCompany).Error
	if err != nil {
		return nil, err
	}
	return &foundCompany, nil
}