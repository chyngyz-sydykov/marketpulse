package main

import (
	"fmt"
	"log"
	"time"

	"github.com/chyngyz-sydykov/marketpulse/config"
	"github.com/chyngyz-sydykov/marketpulse/internal/binance"
	"github.com/chyngyz-sydykov/marketpulse/internal/database"
	"github.com/chyngyz-sydykov/marketpulse/internal/repository"
)

func main() {
	// Connect to the PostgreSQL database
	err := database.ConnectDB()
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Start the scheduler to fetch market data every hour
	//scheduler.StartScheduler()
	for _, currency := range config.DefaultCurrencies {
		record, err := binance.FetchKline(currency)
		if err != nil {
			log.Fatal("Failed to connect to database:", err)
		}

		err = repository.StoreMarketData(currency, record)
		if err != nil {
			log.Fatal("Failed to connect to database:", err)
		}
	}

	// Prevent main from exiting
	fmt.Println("MarketPulse is running...")
	fmt.Println("Salam...")
	for {
		time.Sleep(time.Hour)
	}
}
