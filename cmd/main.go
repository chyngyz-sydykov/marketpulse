package main

import (
	"fmt"
	"log"

	"github.com/chyngyz-sydykov/marketpulse/internal/database"
	"github.com/chyngyz-sydykov/marketpulse/internal/scheduler"
)

func main() {
	// Connect to the PostgreSQL database
	err := database.ConnectDB()
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Start the scheduler to fetch market data every hour
	scheduler.StartScheduler()

	// Prevent main from exiting
	fmt.Println("MarketPulse is running...")
	fmt.Println("Salam...")
	select {}
}
