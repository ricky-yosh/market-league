package services


import (
    "context"
    "log"
    "os"
    finnhub "github.com/Finnhub-Stock-API/finnhub-go/v2"
)

func GetTestStock() (finnhub.Quote, error) {
    // Set up the API client
    apiKey := os.Getenv("FINNHUB_API_KEY")  
    cfg := finnhub.NewConfiguration()
    cfg.AddDefaultHeader("X-Finnhub-Token", apiKey)
    client := finnhub.NewAPIClient(cfg).DefaultApi

    // Define the stock symbol and time range
    symbol := "AAPL"

    // Call the stock candle endpoint
    quote, _, err := client.Quote(context.Background()).Symbol(symbol).Execute()
    if err != nil {
        log.Fatalf("Error fetching stock lastBidAsk: %v", err)
    }
    
    // Output the candle data
    return quote, nil
}
