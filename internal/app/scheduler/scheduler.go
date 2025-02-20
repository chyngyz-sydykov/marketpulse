package scheduler

import (
	"log"
	"sync"
	"time"

	"github.com/chyngyz-sydykov/marketpulse/config"
	"github.com/chyngyz-sydykov/marketpulse/internal/app"
	"github.com/chyngyz-sydykov/marketpulse/internal/infrastructure/binance"

	"github.com/go-co-op/gocron"
)

func StartScheduler() {
	scheduler := gocron.NewScheduler(time.UTC)

	scheduler.Every(1).Hour().Do(func() {
		log.Println("hourly scheduler...")
		executeForAllCurrencies(func(currency string) {
			records, err := binance.FetchKline(currency, config.ONE_HOUR, 1000)
			if err != nil {
				log.Printf("Error fetching data for %s: %v\n", currency, err)
				return
			}
			//err = app.App.MarketDataService.StoreData(currency, records[0])
			err = app.App.MarketDataService.UpsertBatchData(currency, records)
			if err != nil {
				log.Printf("%sError: %s %s\n", config.COLOR_RED, err, config.COLOR_RED)
			}
		})
	})
	// run every 4 hours at 5 minutes past the hour (00:05, 04:05, 08:05, etc.)
	scheduler.Cron("5 */4 * * *").Do(func() {
		log.Println("4 hour scheduler...")
		executeForAllCurrencies(func(curr string) {
			err := app.App.MarketDataService.StoreGroupedRecords(curr, config.FOUR_HOUR)
			if err != nil {
				log.Printf("%sError: %s %s\n", config.COLOR_RED, err, config.COLOR_RED)
			}

			err = app.App.MarketDataService.StoreGroupedRecords(curr, config.ONE_DAY)
			if err != nil {
				log.Printf("%sError: %s %s\n", config.COLOR_RED, err, config.COLOR_RED)
			}
		})

	})
	// run every day at 00:10
	// scheduler.Cron("10 0 * * *").Do(func() {
	// 	log.Println("daily scheduler...")
	// 	executeForAllCurrencies(func(curr string) {
	// 		err := app.App.MarketDataService.StoreGroupedRecords(curr, config.FOUR_HOUR)
	// 		if err != nil {
	// 			log.Printf("%sError: %s %s\n", config.COLOR_RED, err, config.COLOR_RED)
	// 		}

	// 	})
	// })

	scheduler.StartAsync()
}
func executeForAllCurrencies(callback func(curr string)) {
	var wg sync.WaitGroup
	for _, currency := range config.DefaultCurrencies {
		wg.Add(1)
		go func(curr string) {
			defer wg.Done()
			callback(curr) // ✅ Calls the provided function
		}(currency)
	}
	wg.Wait()
}
