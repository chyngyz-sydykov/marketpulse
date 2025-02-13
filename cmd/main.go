package main

import (
	"fmt"
	"log"
	"time"

	"github.com/chyngyz-sydykov/marketpulse/internal/database"
	"github.com/chyngyz-sydykov/marketpulse/internal/scheduler"
)

func main() {
	err := database.ConnectDB()
	if err != nil {
		log.Fatal("❌ Failed to connect to database:", err)
	}
	defer database.DB.Close()

	// Start the scheduler to fetch market data every hour
	scheduler.StartScheduler()

	//marketDataService := marketdata.NewMarketDataService()

	// Prevent main from exiting
	fmt.Println("MarketPulse is running...")
	for {
		time.Sleep(1 * time.Hour)
	}
}
