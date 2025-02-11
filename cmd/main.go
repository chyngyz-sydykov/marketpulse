package main

import (
	"fmt"
	"log"
	"sync"

	"github.com/chyngyz-sydykov/marketpulse/config"
	"github.com/chyngyz-sydykov/marketpulse/internal/binance"
	"github.com/chyngyz-sydykov/marketpulse/internal/database"
	"github.com/chyngyz-sydykov/marketpulse/internal/marketdata"
)

func main() {
	err := database.ConnectDB()
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer database.DB.Close()

	marketDataService := marketdata.NewMarketDataService()

	// Start the scheduler to fetch market data every hour
	//scheduler.StartScheduler()
	var wg sync.WaitGroup
	for _, currency := range config.DefaultCurrencies {
		wg.Add(1)
		go func(curr string) {

			defer wg.Done()

			record, err := binance.FetchKline(curr)
			if err != nil {
				log.Printf("Error fetching data for %s: %v\n", curr, err)
				return
			}

			err = marketDataService.StoreMarketData(curr, record)
			if err != nil {
				log.Printf("Error storing data for %s: %v\n", curr, err)
				return
			}

		}(currency)
	}

	// Prevent main from exiting
	fmt.Println("MarketPulse is running...")
	fmt.Println("Salam..")
	wg.Wait()
}
