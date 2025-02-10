package scheduler

import (
	"log"
	"time"

	"github.com/chyngyz-sydykov/marketpulse/internal/binance"
	"github.com/chyngyz-sydykov/marketpulse/internal/repository"

	"github.com/go-co-op/gocron"
)

// StartScheduler initializes a gocron job to fetch market data every hour
func StartScheduler() {
	scheduler := gocron.NewScheduler(time.UTC)

	// Schedule market data fetching every hour
	scheduler.Every(1).Hour().Do(func() {
		log.Println("Fetching market data...")
		data, err := binance.FetchMarketData()
		if err != nil {
			log.Println("Error fetching market data:", err)
			return
		}
		err = repository.StoreMarketData(data)
		if err != nil {
			log.Println("Error storing market data:", err)
		}
	})

	scheduler.StartAsync() // Start the scheduler in async mode
}
