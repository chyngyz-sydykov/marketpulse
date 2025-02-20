package event

import (
	"context"
	"log"
	"sync"

	"github.com/chyngyz-sydykov/marketpulse/config"
	"github.com/chyngyz-sydykov/marketpulse/internal/core/indicator"
	"github.com/chyngyz-sydykov/marketpulse/internal/core/marketdata"
	"github.com/chyngyz-sydykov/marketpulse/internal/infrastructure/redis"
)

type EventListener struct {
	MarketDataService *marketdata.MarketDataService
	IndicatorService  *indicator.IndicatorService
	RedisService      redis.RedisServiceInterface
}

func NewEventListener(
	marketDataService *marketdata.MarketDataService,
	inidicatorService *indicator.IndicatorService,
	redis redis.RedisServiceInterface) *EventListener {
	return &EventListener{
		MarketDataService: marketDataService,
		IndicatorService:  inidicatorService,
		RedisService:      redis,
	}
}

func (el *EventListener) Listen() {

	ctx := context.Background()

	el.subscribeToEvent(ctx, config.EVENT_NEW_DATA_ADDED, func(currency, eventName, source string) {
		err := el.MarketDataService.StoreGroupedRecords(currency, config.FOUR_HOUR)
		if err != nil {
			log.Printf("Error storing 4H records for %s: %v", currency, err)
		}

		err = el.MarketDataService.StoreGroupedRecords(currency, config.ONE_DAY)
		if err != nil {
			log.Printf("Error storing 1D records for %s: %v", currency, err)
		}
	})
	el.subscribeToEvent(ctx, config.EVENT_NEW_GROUP_DATA_ADDED, func(currency, eventName, eventSource string) {
		err := el.IndicatorService.ComputeAndUpsertBatch(currency, config.FOUR_HOUR)
		if err != nil {
			log.Printf("Error Indicator for 4H records %s: %v", currency, err)
		}

		err = el.IndicatorService.ComputeAndUpsertBatch(currency, config.ONE_DAY)
		if err != nil {
			log.Printf("Error Indicator for 1D records %s: %v", currency, err)
		}
	})
	el.subscribeToEvent(ctx, config.EVENT_NEW_INDICATOR_ADDED, func(currency, eventName, source string) {
		log.Printf("EVENT_NEW_INDICATOR_ADDED")
	})
}
func (el *EventListener) subscribeToEvent(ctx context.Context, eventName string, callback func(curr, eventName, eventSource string)) {
	go el.RedisService.SubscribeToEvent(ctx, eventName, func(event redis.Event) {
		log.Printf("Event received: %s from %s", event.Name, event.Source)

		var wg sync.WaitGroup
		for _, currency := range config.DefaultCurrencies {
			wg.Add(1)
			go func(curr string) {
				defer wg.Done()
				callback(curr, event.Name, event.Source)
			}(currency)
		}
		wg.Wait()
	})
}
