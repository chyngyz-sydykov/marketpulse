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

				records, err := binance.FetchKline(curr, config.ONE_HOUR, 1)
				if err != nil {
					log.Printf("Error fetching data for %s: %v\n", curr, err)
					return
				}

				err = marketDataService.StoreData(curr, records[0])
				if err != nil {
					log.Printf("%sError: %s %s\n", config.COLOR_RED, err, config.COLOR_RED)
					return
				}

			}(currency)
		}
		wg.Wait()
	})
	scheduler.Every(4).Hour().Do(func() {
		log.Println("4 hour scheduler...")
		var wg sync.WaitGroup
		for _, currency := range config.DefaultCurrencies {
			wg.Add(1)
			go func(curr string) {

				defer wg.Done()
				err := marketDataService.StoreGroupedRecords(curr, config.FOUR_HOUR)
				if err != nil {
					log.Printf("%sError: %s %s\n", config.COLOR_RED, err, config.COLOR_RED)
					return
				}

			}(currency)
		}
		wg.Wait()

	})

	// scheduler.Every(1).Day().Do(func() {
	// 	log.Println("1 day scheduler...")
	// 	var wg sync.WaitGroup
	// 	for _, currency := range config.DefaultCurrencies {
	// 		wg.Add(1)
	// 		go func(curr string) {

	// 			defer wg.Done()
	// 			err := marketDataService.StoreGroupedRecords(curr, config.ONE_DAY)
	// 			if err != nil {
	// 				log.Printf("%sError: %s %s\n", config.COLOR_RED, err, config.COLOR_RED)
	// 				return
	// 			}

	// 		}(currency)
	// 	}
	// 	wg.Wait()
	// })

	// scheduler.Every(1).Week().Do(func() {
	// 	log.Println("1 week scheduler...")
	// 	var wg sync.WaitGroup
	// 	for _, currency := range config.DefaultCurrencies {
	// 		wg.Add(1)
	// 		go func(curr string) {

	// 			defer wg.Done()
	// 			err := marketDataService.StoreGroupedRecords(curr, config.ONE_WEEK)
	// 			if err != nil {
	// 				log.Printf("%sError: %s %s\n", config.COLOR_RED, err, config.COLOR_RED)
	// 				return
	// 			}

	// 		}(currency)
	// 	}
	// 	wg.Wait()
	// })

	scheduler.StartAsync() // Start the scheduler in async mode
}
