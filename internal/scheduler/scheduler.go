package scheduler

import (
	"log"
	"sync"
	"time"

	"github.com/chyngyz-sydykov/marketpulse/config"
	"github.com/chyngyz-sydykov/marketpulse/internal/binance"
	"github.com/chyngyz-sydykov/marketpulse/internal/marketdata"

	"github.com/go-co-op/gocron"
)

// StartScheduler initializes a gocron job to fetch market data every hour
func StartScheduler() {
	scheduler := gocron.NewScheduler(time.UTC)
	marketDataService := marketdata.NewMarketDataService()

	scheduler.Every(1).Hour().Do(func() {
		log.Println("hourly scheduler...")
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

				err = marketDataService.StoreData(curr, record)
				if err != nil {
					log.Printf("Error storing data for %s: %v\n", curr, err)
					return
				}

			}(currency)
		}
	})
	scheduler.Every(4).Hour().Do(func() {
		log.Println("4 hour scheduler...")
		var wg sync.WaitGroup
		for _, currency := range config.DefaultCurrencies {
			wg.Add(1)
			go func(curr string) {

				defer wg.Done()
				err := marketDataService.Store4HourRecords(curr)
				if err != nil {
					log.Printf("Error storing data for %s: %v\n", curr, err)
					return
				}

			}(currency)
		}
	})

	scheduler.StartAsync() // Start the scheduler in async mode
}
