package main

// import (
//     "encoding/json"
//     "fmt"
//     "io/ioutil"
//     "log"
//     "net/http"
// )

// // Structure to hold the API response data
// type StockData struct {
//     Close  []float64 `json:"c"`
//     High   []float64 `json:"h"`
//     Low    []float64 `json:"l"`
//     Open   []float64 `json:"o"`
//     Status string    `json:"s"`
//     Time   []int64   `json:"t"`
//     Volume []int64   `json:"v"`
// }

// func main() {
//     apiKey := "cs32qd1r01qk1hurusfgcs32qd1r01qk1hurusg0"
//     symbol := "AAPL"
//     url := fmt.Sprintf("https://finnhub.io/api/v1/stock/candle?symbol=%s&resolution=D&from=1609459200&to=1640995200&token=%s", symbol, apiKey)

//     // Make the HTTP request
//     resp, err := http.Get(url)
//     if err != nil {
//         log.Fatalf("Error making request to Finnhub: %v", err)
//     }
//     defer resp.Body.Close()

//     // Read the response body
//     body, err := ioutil.ReadAll(resp.Body)
//     if err != nil {
//         log.Fatalf("Error reading response body: %v", err)
//     }

//     // Check the status code of the response
//     if resp.StatusCode != http.StatusOK {
//         log.Fatalf("Error: Finnhub API returned status code %d", resp.StatusCode)
//     }

//     // Parse the JSON response
//     var stockData StockData
//     if err := json.Unmarshal(body, &stockData); err != nil {
//         log.Fatalf("Error parsing JSON: %v", err)
//     }

//     // Print the stock data
//     fmt.Printf("Stock data for %s: %+v\n", symbol, stockData)
// }


import (
    "context"
    "fmt"
    "log"
    finnhub "github.com/Finnhub-Stock-API/finnhub-go/v2"
)

func getTestStock() {
    // Set up the API client
    apiKey := "cs32qd1r01qk1hurusfgcs32qd1r01qk1hurusg0"
    cfg := finnhub.NewConfiguration()
    cfg.AddDefaultHeader("X-Finnhub-Token", apiKey)
    client := finnhub.NewAPIClient(cfg).DefaultApi

    // Define the stock symbol and time range
    symbol := "AAPL"
    from := int64(1609459200) // Example: Jan 1, 2021 (in Unix timestamp)
    to := int64(1640995200)   // Example: Jan 1, 2022 (in Unix timestamp)

    // Call the stock candle endpoint
    candles, _, err := client.StockCandles(context.Background()).Symbol(symbol).Resolution("D").From(from).To(to).Execute()
    if err != nil {
        log.Fatalf("Error fetching stock candles: %v", err)
    }

    // Output the candle data
    fmt.Printf("Stock data for %s: %+v\n", symbol, candles)
}
