package api

import (
    "log"
    "time"

    "github.com/market-league/internal/services"
)

func StartDailyTask() {
    go func() {
        for {
            now := time.Now()
            nextRun := time.Date(now.Year(), now.Month(), now.Day(), 18, 28, 0, 0, now.Location()) // Set to 6:28 PM
            if now.After(nextRun) {
                nextRun = nextRun.Add(24 * time.Hour)
            }

            // Wait until the next scheduled time
            time.Sleep(time.Until(nextRun))

            // Execute your task
            log.Println("Running daily stock data fetch...")
            quote, err := services.GetTestStock()
            if err != nil {
                log.Printf("Error fetching stock data: %v", err)
            } else {
                log.Printf("Fetched stock data: %+v", quote)
            }

            // Wait a day for the next execution
            time.Sleep(24 * time.Hour)
        }
    }()
}
