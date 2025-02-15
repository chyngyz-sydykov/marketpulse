package main

import (
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/chyngyz-sydykov/marketpulse/config"
	"github.com/chyngyz-sydykov/marketpulse/internal/binance"
	"github.com/chyngyz-sydykov/marketpulse/internal/database"
	"github.com/chyngyz-sydykov/marketpulse/internal/marketdata"
	"github.com/chyngyz-sydykov/marketpulse/internal/scheduler"
)

func main() {

	fmt.Println("MarketPulse is running...")
	err := database.ConnectDB()
	if err != nil {
		log.Fatal("❌ Failed to connect to database:", err)
	}
	defer database.DB.Close()

	//getAlotOfData()
	// Start the scheduler to fetch market data every hour
	scheduler.StartScheduler()
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

	for {
		time.Sleep(1 * time.Hour)
	}
}

func getAlotOfData() {
	startTime := time.Now() // Start timing execution
	log.Println("get a lot of data from binance...")
	var wg sync.WaitGroup
	marketDataService := marketdata.NewMarketDataService()
	for _, currency := range config.DefaultCurrencies {
		wg.Add(1)
		go func(curr string) {

			defer wg.Done()

			records, err := binance.FetchKline(curr, config.ONE_HOUR, 1000)
			if err != nil {
				log.Printf("Error fetching data for %s: %v\n", curr, err)
				return
			}

			err = marketDataService.StoreBatchData(curr, records)
			if err != nil {
				log.Printf("%sError: %s %s\n", config.COLOR_RED, err, config.COLOR_RED)
				return
			}

		}(currency)
	}
	wg.Wait()
	totalDuration := time.Since(startTime) // Measure total execution time
	log.Printf("🚀 time: %v\n", totalDuration)
}
