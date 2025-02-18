package main

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/chyngyz-sydykov/marketpulse/config"
	"github.com/chyngyz-sydykov/marketpulse/internal/app"
	"github.com/chyngyz-sydykov/marketpulse/internal/app/scheduler"
	"github.com/chyngyz-sydykov/marketpulse/internal/infrastructure/binance"
	"github.com/chyngyz-sydykov/marketpulse/internal/infrastructure/database"
	"github.com/chyngyz-sydykov/marketpulse/internal/infrastructure/redis"
)

func main() {

	fmt.Println("MarketPulse is running...")
	err := database.ConnectDB()
	if err != nil {
		log.Fatal("❌ Failed to connect to database:", err)
	}
	defer database.DB.Close()

	err = redis.ConnectRedis()
	if err != nil {
		log.Fatal("❌ Failed to connect to redis:", err)
	}
	defer redis.Redis.Close()

	app.NewContainer()

	// Start the scheduler to fetch market data every hour
	scheduler.StartScheduler()

	startEventListeners()
	//getAlotOfData()

	// get data from binance every hour
	// save 1h data to db
	// save 4h data and indicators to db if 4h passed
	// if not enough 1h data, group by 1h data to 4h data with indicators

	// every day group 1h data to 1d data and save to redis
	// every day group 1h data to 7d data and save to redis
	for {
		time.Sleep(1 * time.Hour)
	}
}

func startEventListeners() {
	redisService := redis.NewRedisService(redis.Redis)

	// Example: Publish an event
	ctx := context.Background()

	// Example: Subscribe to events
	go redisService.SubscribeToEvent(ctx, "NewRecordAdded", func(event redis.Event) {
		log.Printf("Received event: %s from %s", event.Name, event.Source)
		//indicatorService := indicator.NewIndicator()
		var wg sync.WaitGroup
		for _, currency := range config.DefaultCurrencies {
			wg.Add(1)
			go func(curr string) {
				defer wg.Done()

				err := app.App.MarketDataService.StoreGroupedRecords(curr, config.ONE_HOUR)
				if err != nil {
					log.Printf("Error storing records for %s: %v", curr, err)
				}

				err = app.App.MarketDataService.StoreGroupedRecords(curr, config.ONE_HOUR)
				if err != nil {
					log.Printf("Error storing records for %s: %v", curr, err)
				}

				//indicatorService.ComputeAndStore(currency, config.FOUR_HOUR)

			}(currency)
		}
		wg.Wait()
	})
}

func getAlotOfData() {
	startTime := time.Now() // Start timing execution
	log.Println("get a lot of data from binance...")
	var wg sync.WaitGroup
	for _, currency := range config.DefaultCurrencies {
		wg.Add(1)
		go func(curr string) {

			defer wg.Done()

			records, err := binance.FetchKline(curr, config.ONE_HOUR, 100)
			if err != nil {
				log.Printf("Error fetching data for %s: %v\n", curr, err)
				return
			}

			err = app.App.MarketDataService.UpsertBatchData(curr, records)
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
