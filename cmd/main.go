package main

import (
	"fmt"
	"log"
	"time"

	"github.com/chyngyz-sydykov/marketpulse/internal/database"
)

func main() {
	err := database.ConnectDB()
	if err != nil {
		log.Fatal("❌ Failed to connect to database:", err)
	}
	defer database.DB.Close()

	// Start the scheduler to fetch market data every hour
	//scheduler.StartScheduler()
	// in := indicator.NewIndicator()

	// timestamp1, _ := time.Parse(time.RFC3339, "2025-02-13T05:00:00+00:00")
	// timestamp2, _ := time.Parse(time.RFC3339, "2025-02-13T06:00:00+00:00")
	// timestamp3, _ := time.Parse(time.RFC3339, "2025-02-13T07:00:00+00:00")
	// timestamp4, _ := time.Parse(time.RFC3339, "2025-02-13T08:00:00+00:00")

	// in.ComputeAndStore("btc", []binance.RecordDto{
	// 	{Close: 100, High: 102, Low: 98, Timeframe: "1h", Timestamp: timestamp1},
	// 	{Close: 102, High: 106, Low: 100, Timeframe: "1h", Timestamp: timestamp2},
	// 	{Close: 103, High: 105, Low: 100, Timeframe: "1h", Timestamp: timestamp3},
	// 	{Close: 104, High: 108, Low: 101, Timeframe: "1h", Timestamp: timestamp4},

	//marketDataService := marketdata.NewMarketDataService()

	// Prevent main from exiting
	fmt.Println("MarketPulse is running...")
	for {
		time.Sleep(1 * time.Hour)
	}
}
